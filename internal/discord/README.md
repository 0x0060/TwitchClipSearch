# Discord Integration

This package handles Discord webhook integration for the TwitchClipSearch application.

## Features

- Webhook message formatting
- Clip notification delivery
- Rate limit handling
- Error recovery

## Components

- `webhook.go`: Core webhook sending functionality
- `message.go`: Message formatting and templating
- `client.go`: Discord API client implementation
- `config.go`: Discord-specific configuration

## Usage

```go
// Initialize Discord client
discord := discord.NewClient(webhookURL)

// Send clip notification
err := discord.SendClipNotification(clip)
```

## Configuration

```yaml
discord:
  webhook_url: "https://discord.com/api/webhooks/..."
  username: "TwitchClipBot"
  rate_limit: 5 # requests per second
  retry_attempts: 3
```

## Error Handling

- Rate limit exceeded
- Network failures
- Invalid webhook URL
- Message formatting errors