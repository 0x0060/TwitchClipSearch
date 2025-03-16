// Package metrics provides instrumentation and monitoring capabilities
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics for the application
type Metrics struct {
	ClipsProcessed   *prometheus.CounterVec
	ClipsFailed      *prometheus.CounterVec
	ProcessingTime   *prometheus.HistogramVec
	APIRequestsTotal *prometheus.CounterVec
	APILatency       *prometheus.HistogramVec
	QueueSize        *prometheus.GaugeVec
	WorkerUtilization *prometheus.GaugeVec
}

// New creates and registers all application metrics
func New(namespace string) *Metrics {
	return &Metrics{
		ClipsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "clips_processed_total",
				Help:      "Total number of clips processed",
			},
			[]string{"streamer", "status"},
		),
		ClipsFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "clips_failed_total",
				Help:      "Total number of clip processing failures",
			},
			[]string{"streamer", "error_type"},
		),
		ProcessingTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "clip_processing_duration_seconds",
				Help:      "Time spent processing clips",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"streamer"},
		),
		APIRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "api_requests_total",
				Help:      "Total number of Twitch API requests",
			},
			[]string{"endpoint", "status"},
		),
		APILatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "api_request_duration_seconds",
				Help:      "Twitch API request latencies",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"endpoint"},
		),
		QueueSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "worker_queue_size",
				Help:      "Current size of the worker queue",
			},
			[]string{"pool"},
		),
		WorkerUtilization: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "worker_utilization",
				Help:      "Current worker pool utilization",
			},
			[]string{"pool"},
		),
	}
}

// RecordClipProcessed increments the clips processed counter
func (m *Metrics) RecordClipProcessed(streamer, status string) {
	m.ClipsProcessed.WithLabelValues(streamer, status).Inc()
}

// RecordClipFailed increments the clips failed counter
func (m *Metrics) RecordClipFailed(streamer, errorType string) {
	m.ClipsFailed.WithLabelValues(streamer, errorType).Inc()
}

// ObserveProcessingTime records clip processing duration
func (m *Metrics) ObserveProcessingTime(streamer string, duration float64) {
	m.ProcessingTime.WithLabelValues(streamer).Observe(duration)
}

// RecordAPIRequest increments the API request counter
func (m *Metrics) RecordAPIRequest(endpoint, status string) {
	m.APIRequestsTotal.WithLabelValues(endpoint, status).Inc()
}

// ObserveAPILatency records API request latency
func (m *Metrics) ObserveAPILatency(endpoint string, duration float64) {
	m.APILatency.WithLabelValues(endpoint).Observe(duration)
}

// SetQueueSize sets the current queue size metric
func (m *Metrics) SetQueueSize(pool string, size float64) {
	m.QueueSize.WithLabelValues(pool).Set(size)
}

// SetWorkerUtilization sets the worker utilization metric
func (m *Metrics) SetWorkerUtilization(pool string, utilization float64) {
	m.WorkerUtilization.WithLabelValues(pool).Set(utilization)
}