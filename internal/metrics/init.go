package metrics

import (
	"sync"
)

var (
	metrics *Metrics
	once    sync.Once
)

// InitMetrics initializes the metrics system
func InitMetrics() {
	once.Do(func() {
		metrics = New("twitchclipsearch")
	})
}

// RecordRateLimitHit records a rate limit hit for the specified service
func RecordRateLimitHit(service string) {
	if metrics != nil {
		metrics.APIRequestsTotal.WithLabelValues(service, "rate_limited").Inc()
	}
}

// RecordRetryAttempt records a retry attempt for the specified service
func RecordRetryAttempt(service string) {
	if metrics != nil {
		metrics.APIRequestsTotal.WithLabelValues(service, "retry").Inc()
	}
}

// RecordError records an error for the specified error type
func RecordError(errorType string) {
	if metrics != nil {
		metrics.ClipsFailed.WithLabelValues("system", errorType).Inc()
	}
}