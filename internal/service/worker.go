package service

import (
	"sync"
	"sync/atomic"

	"twitchclipsearch/internal/logger"
	"twitchclipsearch/internal/metrics"
)

type Task func()

type WorkerPool struct {
	workers      int
	tasks        chan Task
	shutdown     chan struct{}
	wg           sync.WaitGroup
	isRunning    atomic.Bool
	maxQueued    int
	errorHandler func(error)
}

func NewWorkerPool(workers int, opts ...WorkerPoolOption) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}

	pool := &WorkerPool{
		workers:      workers,
		tasks:        make(chan Task),
		shutdown:     make(chan struct{}),
		maxQueued:    workers * 100,
		errorHandler: func(err error) {},
	}

	for _, opt := range opts {
		opt(pool)
	}

	return pool
}

type WorkerPoolOption func(*WorkerPool)

func WithMaxQueued(max int) WorkerPoolOption {
	return func(p *WorkerPool) {
		if max > 0 {
			p.maxQueued = max
		}
	}
}

func WithErrorHandler(handler func(error)) WorkerPoolOption {
	return func(p *WorkerPool) {
		if handler != nil {
			p.errorHandler = handler
		}
	}
}

func (p *WorkerPool) Start() {
	if !p.isRunning.CompareAndSwap(false, true) {
		return
	}

	p.tasks = make(chan Task, p.maxQueued)
	p.shutdown = make(chan struct{})

	// Update worker pool metrics
	metrics.RecordGauge("worker_pool_size", float64(p.workers), "pool", "default")

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) Stop() {
	if !p.isRunning.CompareAndSwap(true, false) {
		return
	}

	close(p.shutdown)
	p.wg.Wait()
	close(p.tasks)
}

func (p *WorkerPool) Submit(task Task) bool {
	if !p.isRunning.Load() {
		logger.Error("Cannot submit task: worker pool is not running")
		return false
	}

	select {
	case p.tasks <- task:
		// Update queue size metric
		metrics.RecordQueueSize("default", float64(len(p.tasks)))
		return true
	default:
		logger.Error("Task queue is full")
		metrics.RecordError("worker_pool_queue_full")
		return false
	}
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.shutdown:
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok {
							logger.Error("Worker panic recovered", "error", err)
							metrics.RecordError("worker_panic")
							p.errorHandler(err)
						}
					}
				}()
				// Update worker utilization metric
				metrics.RecordWorkerUtilization("default", 1.0)
				task()
				metrics.RecordWorkerUtilization("default", 0.0)
			}()
			// Update queue size metric after task completion
			metrics.RecordQueueSize("default", float64(len(p.tasks)))
		}
	}
}

func (p *WorkerPool) Size() int {
	return p.workers
}

func (p *WorkerPool) IsRunning() bool {
	return p.isRunning.Load()
}
