package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"twitchclipsearch/internal/database"

	"golang.org/x/time/rate"
)

// Client represents a Discord webhook client
type Client struct {
	webhookURL    string
	username      string
	rateLimiter   *rate.Limiter
	httpClient    *http.Client
	retryAttempts int
}

// ClientConfig holds configuration for the Discord client
type ClientConfig struct {
	WebhookURL    string
	Username      string
	RateLimit     float64
	RetryAttempts int
}

// NewClient creates a new Discord webhook client
func NewClient(config *ClientConfig) *Client {
	if config.RateLimit == 0 {
		config.RateLimit = 5 // default 5 requests per second
	}
	if config.RetryAttempts == 0 {
		config.RetryAttempts = 3 // default 3 retry attempts
	}

	return &Client{
		webhookURL:    config.WebhookURL,
		username:      config.Username,
		rateLimiter:   rate.NewLimiter(rate.Limit(config.RateLimit), 1),
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		retryAttempts: config.RetryAttempts,
	}
}

// SendClipNotification sends a clip notification to Discord
func (c *Client) SendClipNotification(clip *database.Clip) error {
	// Wait for rate limit
	// Add context import at the top of the file
	if err := c.rateLimiter.Wait(context.Background()); err != nil {
		return fmt.Errorf("rate limit wait error: %w", err)
	}

	// Create message payload
	msg := NewMessage(clip)
	msg.Username = c.username

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send with retries
	for attempt := 0; attempt < c.retryAttempts; attempt++ {
		err = c.sendWebhook(payload)
		if err == nil {
			return nil
		}

		// Wait before retry
		if attempt < c.retryAttempts-1 {
			time.Sleep(time.Second * time.Duration(attempt+1))
		}
	}

	return fmt.Errorf("failed to send webhook after %d attempts: %w", c.retryAttempts, err)
}

// sendWebhook sends the actual HTTP request to Discord
func (c *Client) sendWebhook(payload []byte) error {
	req, err := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("webhook request failed with status %d", resp.StatusCode)
	}

	return nil
}
