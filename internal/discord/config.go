package discord

import (
	"fmt"
	"time"
)

// Config represents Discord-specific configuration
type Config struct {
	WebhookURL    string        `json:"webhook_url"`
	Username      string        `json:"username"`
	RateLimit     float64       `json:"rate_limit"`
	RetryAttempts int           `json:"retry_attempts"`
	Timeout       time.Duration `json:"timeout"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.WebhookURL == "" {
		return fmt.Errorf("webhook URL is required")
	}

	if c.Username == "" {
		c.Username = "TwitchClipBot" // default username
	}

	if c.RateLimit <= 0 {
		c.RateLimit = 5 // default 5 requests per second
	}

	if c.RetryAttempts <= 0 {
		c.RetryAttempts = 3 // default 3 retry attempts
	}

	if c.Timeout <= 0 {
		c.Timeout = 10 * time.Second // default 10 second timeout
	}

	return nil
}

// NewConfig creates a new Discord configuration with default values
func NewConfig(webhookURL string) *Config {
	return &Config{
		WebhookURL:    webhookURL,
		Username:      "TwitchClipBot",
		RateLimit:     5,
		RetryAttempts: 3,
		Timeout:       10 * time.Second,
	}
}