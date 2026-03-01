.PHONY: help build test run clean lint swagger docker bench dev

# Default target: show help
.DEFAULT_GOAL := help

# Display help information
help:
	@echo "Home Lab API Service - Available Make Targets:"
	@echo ""
	@echo "  make build   - Build the application binary to bin/homelab-api"
	@echo "  make test    - Run tests with race detection and coverage report"
	@echo "  make run     - Run the application"
	@echo "  make clean   - Clean build artifacts (bin/, coverage files)"
	@echo "  make lint    - Run golangci-lint code linter"
	@echo "  make swagger - Generate Swagger API documentation"
	@echo "  make docker  - Build Docker image (homelab-api:latest)"
	@echo "  make bench   - Run benchmarks for research code"
	@echo "  make dev     - Run with hot reload using air"
	@echo ""

# Build the application
build:
	go build -o bin/homelab-api ./cmd/api

# Run tests with coverage
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run the application
run:
	go run ./cmd/api

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Run linter
lint:
	golangci-lint run

# Generate Swagger docs
swagger:
	swag init -g cmd/api/main.go -o api/

# Build Docker image
docker:
	docker build -t homelab-api:latest -f deployments/Dockerfile .

# Run benchmarks
bench:
	go test -bench=. -benchmem ./research/...

# Development: run with hot reload (requires air)
dev:
	air
