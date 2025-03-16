package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// RequestID adds a unique identifier to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), "request_id", requestID),
		))
	})
}

// Logger logs request details
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture status code
		ww := &responseWriter{ResponseWriter: w}

		// Process request
		next.ServeHTTP(ww, r)

		// Log request details
		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			ww.status,
			duration,
		)

		// Record metrics
		// TODO: Implement metrics recording
		// For now, we'll skip metrics recording since the metrics package is not available
	})
}

// Recover handles panics in HTTP handlers
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("panic: %v", err)

				// Return 500 Internal Server Error
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
