.PHONY: all build test test-coverage lint clean run

all: build

build:
	go build -v ./...

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -f coverage.out coverage.html
	go clean

run:
	go run ./cmd/server/main.go

.DEFAULT_GOAL := all