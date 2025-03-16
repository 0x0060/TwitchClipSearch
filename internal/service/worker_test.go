package service

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	t.Run("basic functionality", func(t *testing.T) {
		pool := NewWorkerPool(2)
		pool.Start()
		defer pool.Stop()

		var counter atomic.Int32
		for i := 0; i < 5; i++ {
			if !pool.Submit(func() {
				counter.Add(1)
				time.Sleep(10 * time.Millisecond)
			}) {
				t.Error("Failed to submit task")
			}
		}

		time.Sleep(100 * time.Millisecond)
		if counter.Load() != 5 {
			t.Errorf("Expected 5 tasks to complete, got %d", counter.Load())
		}
	})

	t.Run("error handling", func(t *testing.T) {
		var errCount atomic.Int32
		pool := NewWorkerPool(1, WithErrorHandler(func(err error) {
			errCount.Add(1)
		}))
		pool.Start()
		defer pool.Stop()

		pool.Submit(func() {
			panic("test error")
		})

		time.Sleep(50 * time.Millisecond)
		if errCount.Load() != 1 {
			t.Errorf("Expected 1 error, got %d", errCount.Load())
		}
	})

	t.Run("queue full behavior", func(t *testing.T) {
		pool := NewWorkerPool(1, WithMaxQueued(1))
		pool.Start()
		defer pool.Stop()

		pool.Submit(func() {
			time.Sleep(50 * time.Millisecond)
		})

		pool.Submit(func() {})

		if pool.Submit(func() {}) {
			t.Error("Expected task submission to fail when queue is full")
		}
	})

	t.Run("stop behavior", func(t *testing.T) {
		pool := NewWorkerPool(1)
		pool.Start()

		var completed atomic.Bool
		pool.Submit(func() {
			time.Sleep(50 * time.Millisecond)
			completed.Store(true)
		})

		pool.Stop()
		if !completed.Load() {
			t.Error("Task should complete before stop returns")
		}

		if pool.Submit(func() {}) {
			t.Error("Should not accept tasks after stop")
		}
	})
}

func BenchmarkWorkerPool(b *testing.B) {
	pool := NewWorkerPool(4)
	pool.Start()
	defer pool.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Submit(func() {
				time.Sleep(time.Microsecond)
			})
		}
	})
}