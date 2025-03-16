package discord

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"twitchclipsearch/internal/metrics"
)

// WebhookError represents a Discord webhook error
type WebhookError struct {
	StatusCode int
	Message    string
	RetryAfter int
}

func (e *WebhookError) Error() string {
	return fmt.Sprintf("webhook error: %s (status: %d)", e.Message, e.StatusCode)
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	if webhookErr, ok := err.(*WebhookError); ok {
		return webhookErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// handleWebhookResponse processes the webhook response and returns appropriate error
func handleWebhookResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	webhookErr := &WebhookError{
		StatusCode: resp.StatusCode,
		Message:    fmt.Sprintf("request failed with status %d", resp.StatusCode),
	}

	// Handle rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		// Record rate limit metric
		metrics.RecordRateLimitHit("discord_webhook")
		
		// Parse retry-after header
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			webhookErr.Message = "rate limit exceeded"
			// Add retry delay logic here if needed
		}
	}

	return webhookErr
}

// retryWithBackoff implements exponential backoff for retrying failed requests
func retryWithBackoff(ctx context.Context, fn func() error, maxAttempts int) error {
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			
			// Record retry metric
			metrics.RecordRetryAttempt("discord_webhook")

			// If it's a rate limit error, use the retry-after value
			if IsRateLimitError(err) {
				time.Sleep(time.Second * 5) // Default retry delay for rate limits
				continue
			}

			// Exponential backoff
			if attempt < maxAttempts-1 {
				backoff := time.Second * time.Duration(1<<uint(attempt))
				time.Sleep(backoff)
			}
		}
	}

	return fmt.Errorf("max retry attempts reached: %w", lastErr)
}