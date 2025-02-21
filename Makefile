.PHONY: dev build clean test lint

# Development
dev:
	air

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Clean build files
clean:
	rm -rf bin tmp

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Setup development environment
setup: install-tools
	cp .env.example .env