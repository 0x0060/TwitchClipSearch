# Configuration Management

This directory contains configuration-related files and utilities for the TwitchClipSearch application.

## Structure

- `config.go`: Core configuration loading and validation
- `development.yaml`: Development environment configuration
- `production.yaml`: Production environment configuration
- `test.yaml`: Test environment configuration

## Usage

Configuration is loaded based on the `APP_ENV` environment variable:

```bash
# Development (default)
export APP_ENV=development

# Production
export APP_ENV=production

# Testing
export APP_ENV=test
```

## Configuration Parameters

- `database`: Database connection settings
- `twitch`: Twitch API credentials and settings
- `discord`: Discord webhook configuration
- `server`: HTTP server settings
- `metrics`: Prometheus metrics configuration