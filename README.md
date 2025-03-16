# TwitchClipSearch

[![Go Report Card](https://goreportcard.com/badge/github.com/0x0060/twitchclipsearch)](https://goreportcard.com/report/github.com/0x0060/twitchclipsearch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/0x0060/twitchclipsearch)](https://golang.org/)

A powerful service that automatically fetches and manages Twitch clips, providing Discord notifications for new content from your favorite streamers.

## Features

- Real-time Twitch clip monitoring
- Automatic Discord notifications
- RESTful API for clip management
- Configurable per-streamer webhooks
- Prometheus metrics integration
- Multi-environment configuration support

## Prerequisites

| Requirement | Version |
|------------|----------|
| Go | >= 1.16 |
| SQLite | 3.x |
| Docker (optional) | >= 20.10 |
| Kubernetes (optional) | >= 1.19 |

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/0x0060/twitchclipsearch.git
   cd twitchclipsearch
   ```

2. Set up environment variables:
   ```bash
   export TWITCH_CLIENT_ID="your_client_id"
   export TWITCH_CLIENT_SECRET="your_client_secret"
   export DISCORD_WEBHOOK_URL="your_webhook_url"
   ```

3. Run the application:
   ```bash
   make run
   ```

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|----------|
| TWITCH_CLIENT_ID | Twitch API Client ID | Yes | - |
| TWITCH_CLIENT_SECRET | Twitch API Client Secret | Yes | - |
| DISCORD_WEBHOOK_URL | Discord Webhook URL | Yes | - |
| DB_PATH | Database file path | No | clips.db |
| LOG_OUTPUT | Log output destination | No | stdout |

### Configuration Files

The application supports multiple environment-specific configuration files:

- `config/development.yaml`: Development environment settings
- `config/production.yaml`: Production environment settings
- `config/test.yaml`: Test environment settings

## 🔧 Development

### Project Structure

```
├── cmd/                 # Application entrypoints
├── config/             # Configuration files
├── deployment/         # Deployment configurations
├── docs/               # Documentation
├── internal/           # Internal packages
└── githooks/           # Git hooks
```

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make test` | Run tests |
| `make lint` | Run linters |
| `make docker` | Build Docker image |

## Deployment

### Docker

```bash
docker-compose -f deployment/docker/docker-compose.yml up
```

### Kubernetes

```bash
kubectl apply -k deployment/kubernetes/overlays/production
```

## Metrics

Prometheus metrics are available at `/metrics` endpoint with the following key metrics:

| Metric | Type | Description |
|--------|------|-------------|
| clips_processed_total | Counter | Total number of processed clips |
| webhook_requests_total | Counter | Total number of webhook requests |
| api_request_duration_seconds | Histogram | API request duration |

## Documentation

- [Architecture Overview](docs/architecture/README.md)
- [Development Guidelines](docs/development/README.md)
- [Monitoring Guide](docs/monitoring/README.md)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Twitch API](https://dev.twitch.tv/docs/api/)
- [Discord Webhooks](https://discord.com/developers/docs/resources/webhook)
- [Prometheus](https://prometheus.io/)