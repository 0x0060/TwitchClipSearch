# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache make git

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o twitchclipsearch ./cmd/twitchclipsearch

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/twitchclipsearch .

# Copy configuration
COPY config/production.yaml /app/config/config.yaml

# Set environment variables
ENV GO_ENV=production

# Expose application port
EXPOSE 8080

# Run the application
CMD ["/app/twitchclipsearch"]