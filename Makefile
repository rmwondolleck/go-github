# Default target: show help
.DEFAULT_GOAL := help

# Phony targets (targets that don't create files)
.PHONY: help build test run clean lint swagger docker bench dev

# Display help information
help:
	@echo "Home Lab API Service — Dual-Mode Binary (HTTP + MCP)"
	@echo ""
	@echo "  A single ./bin/homelab-api binary starts BOTH the HTTP API server"
	@echo "  (port 8080, serves Kubernetes traffic) AND the MCP stdio server"
	@echo "  (stdin/stdout, serves VS Code Copilot and JetBrains AI) concurrently."
	@echo "  Both modes shut down on SIGINT/SIGTERM."
	@echo ""
	@echo "Available Make Targets:"
	@echo ""
	@echo "  make build   - Build the dual-mode binary to bin/homelab-api"
	@echo "  make test    - Run tests with race detection and coverage report"
	@echo "  make run     - Run the binary (starts BOTH HTTP API and MCP stdio)"
	@echo "  make clean   - Clean build artifacts (bin/, coverage, Swagger docs)"
	@echo "  make lint    - Run golangci-lint code linter"
	@echo "  make swagger - Generate Swagger API documentation"
	@echo "  make docker  - Build Docker image (homelab-api:latest)"
	@echo "  make bench   - Run benchmarks for research code"
	@echo "  make dev     - Run with hot reload using air"
	@echo ""
	@echo "MCP Quick Smoke Test:"
	@echo "  echo '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"initialize\",\"params\":{\"protocolVersion\":\"2024-11-05\",\"capabilities\":{},\"clientInfo\":{\"name\":\"test\",\"version\":\"0.1\"}}}' | ./bin/homelab-api 2>/dev/null"
	@echo ""

# Build the dual-mode application (HTTP API + MCP stdio in one binary)
build:
	mkdir -p bin/
	go build -o bin/homelab-api ./cmd/api

# Run tests with coverage
test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run the application (starts BOTH HTTP API on :8080 AND MCP stdio server)
run: build
	./bin/homelab-api

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html api/docs.go api/swagger.json api/swagger.yaml

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
