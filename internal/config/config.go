package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
)

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Twitch   TwitchConfig   `yaml:"twitch"`
	Discord  DiscordConfig  `yaml:"discord"`
	Server   ServerConfig   `yaml:"server"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path           string        `yaml:"path"`
	MaxConnections int           `yaml:"max_connections"`
	Timeout        time.Duration `yaml:"timeout_seconds"`
}

// TwitchConfig holds Twitch API configuration
type TwitchConfig struct {
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	CheckIntervalSecs int    `yaml:"check_interval_secs"`
}

// DiscordConfig holds Discord webhook configuration
type DiscordConfig struct {
	Streamers map[string]string `yaml:"streamers"`
	RateLimit int               `yaml:"rate_limit"`
	Username  string            `yaml:"username"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout_seconds"`
	WriteTimeout time.Duration `yaml:"write_timeout_seconds"`
}

// MetricsConfig holds Prometheus metrics configuration
type MetricsConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Namespace string `yaml:"namespace"`
	Endpoint  string `yaml:"endpoint"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig loads the configuration based on the environment
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	configPath := fmt.Sprintf("config/%s.yaml", env)

	once.Do(func() {
		data, err := os.ReadFile(configPath)
		if err != nil {
			config = nil
			return
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			config = nil
			return
		}

		// Convert timeout seconds to duration
		cfg.Database.Timeout *= time.Second
		cfg.Server.ReadTimeout *= time.Second
		cfg.Server.WriteTimeout *= time.Second

		config = &cfg
	})

	if config == nil {
		return nil, fmt.Errorf("failed to load configuration")
	}

	return config, nil
}
