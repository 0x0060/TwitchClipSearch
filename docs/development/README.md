# Development Guidelines

This directory contains development setup instructions and guidelines for the TwitchClipSearch application.

## Development Setup

1. Prerequisites
   - Go 1.x
   - Docker
   - Make
   - Git

2. Local Environment Setup
   ```bash
   # Clone the repository
   git clone https://github.com/0x0060/TwitchClipSearch.git
   
   # Install dependencies
   make deps
   
   # Set up development configuration
   cp config/development.yaml config/config.yaml
   
   # Run tests
   make test
   ```

## Development Standards

1. Code Style
   - Follow Go standard formatting (gofmt)
   - Use meaningful variable and function names
   - Write comprehensive comments
   - Include unit tests for new features

2. Git Workflow
   - Use feature branches
   - Write descriptive commit messages
   - Keep commits atomic and focused
   - Rebase before merging

3. Testing
   - Write unit tests for new code
   - Maintain test coverage above 80%
   - Include integration tests for API endpoints
   - Use table-driven tests where appropriate

4. Documentation
   - Update API documentation for endpoint changes
   - Document configuration changes
   - Keep README files current
   - Include code examples where helpful

## Development Tools

1. Recommended IDE Setup
   - VSCode with Go extension
   - Delve for debugging
   - golangci-lint for linting

2. Useful Commands
   ```bash
   # Run linter
   make lint
   
   # Run tests with coverage
   make test-coverage
   
   # Start development server
   make run-dev
   ```

## Troubleshooting

1. Common Issues
   - Database connection problems
   - API rate limiting
   - Authentication errors

2. Debugging Tips
   - Use structured logging
   - Enable debug mode in config
   - Check application metrics
   - Review error logs