# Home Lab API Service

A RESTful API service built with Go and Gin framework for managing smart home devices and providing health monitoring capabilities.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Local Development Setup](#local-development-setup)
- [Environment Variables](#environment-variables)
- [Building and Running](#building-and-running)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Development Tools](#development-tools)

## Features

- RESTful API with Gin framework
- Health monitoring endpoint
- Request ID tracking for distributed tracing
- Structured logging with Go's slog
- Graceful shutdown support
- Middleware chain (Request ID, Logging, Recovery)
- Swagger API documentation support

## Prerequisites

- Go 1.25.0 or higher
- Docker (optional, for containerized deployment)
- Make (optional, for using Makefile commands)

## Local Development Setup

1. **Clone the repository**

```bash
git clone https://github.com/rmwondolleck/go-github.git
cd go-github
```

2. **Install dependencies**

```bash
go mod download
```

3. **Build the application**

```bash
make build
# Or without make:
go build -o bin/homelab-api ./cmd/api
```

4. **Run the application**

```bash
make run
# Or without make:
go run ./cmd/api
```

The server will start on port 8080 by default (or the port specified in the `PORT` environment variable).

## Environment Variables

The following environment variables can be configured:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Port number for the HTTP server | `8080` | No |

### Example Configuration

Create a `.env` file in the project root (optional):

```bash
PORT=8080
```

Or export environment variables directly:

```bash
export PORT=8080
```

## Building and Running

### Using Make

The project includes a Makefile with convenient targets:

```bash
# Show all available targets
make help

# Build the application
make build

# Run tests with coverage
make test

# Run the application
make run

# Clean build artifacts
make clean

# Run linter
make lint

# Generate Swagger documentation
make swagger

# Build Docker image
make docker

# Run with hot reload (requires air)
make dev
```

### Without Make

```bash
# Build
mkdir -p bin/
go build -o bin/homelab-api ./cmd/api

# Run
go run ./cmd/api

# Test
go test -v -race ./...

# Clean
rm -rf bin/ coverage.out coverage.html
```

### Using Docker

```bash
# Build Docker image
docker build -t homelab-api:latest -f deployments/Dockerfile .

# Run container
docker run -p 8080:8080 homelab-api:latest

# Run with custom port
docker run -p 3000:3000 -e PORT=3000 homelab-api:latest
```

## API Endpoints

### Health Check

Check the health status of the API service.

**Endpoint:** `GET /health`

**Response:**

```json
{
  "status": "ok"
}
```

**Example using curl:**

```bash
curl http://localhost:8080/health
```

**Example Response:**

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "ok"
}
```

---

### API Version Info

Get information about the API version.

**Endpoint:** `GET /api/v1`

**Response:**

```json
{
  "message": "API v1"
}
```

**Example using curl:**

```bash
curl http://localhost:8080/api/v1
```

**Example Response:**

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{
  "message": "API v1"
}
```

---

### Complete API Examples

#### 1. Health Check with Verbose Output

```bash
# Basic health check
curl -v http://localhost:8080/health

# Health check with request ID in header
curl -H "X-Request-ID: my-custom-id" http://localhost:8080/health

# Health check with formatted JSON output
curl -s http://localhost:8080/health | jq
```

#### 2. API V1 Endpoint

```bash
# Get API version info
curl http://localhost:8080/api/v1

# With verbose output
curl -v http://localhost:8080/api/v1

# With custom headers
curl -H "Content-Type: application/json" \
     -H "X-Request-ID: test-123" \
     http://localhost:8080/api/v1
```

#### 3. Testing Different HTTP Methods

```bash
# GET request (supported)
curl -X GET http://localhost:8080/health

# POST request (should return 404 for /health)
curl -X POST http://localhost:8080/health

# OPTIONS request
curl -X OPTIONS http://localhost:8080/health -v
```

#### 4. Error Handling Examples

```bash
# Non-existent endpoint (404)
curl -v http://localhost:8080/nonexistent

# Invalid path
curl -v http://localhost:8080/api/invalid
```

## Testing

### Run All Tests

```bash
make test
# Or:
go test -v -race ./...
```

### Run Specific Package Tests

```bash
# Test server package
go test -v ./internal/server/...

# Test middleware
go test -v ./internal/middleware/...

# Test handlers
go test -v ./internal/handlers/...
```

### Run Integration Tests

```bash
go test -v ./tests/integration/...
```

### Generate Coverage Report

```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html in your browser
```

### Run Benchmarks

```bash
make bench
# Or:
go test -bench=. -benchmem ./research/...
```

## Development Tools

### Hot Reload with Air

For development with automatic reloading:

```bash
# Install air
go install github.com/air-verse/air@latest

# Run with hot reload
make dev
# Or:
air
```

### Linting

```bash
# Install golangci-lint
# See: https://golangci-lint.run/usage/install/

# Run linter
make lint
# Or:
golangci-lint run
```

### Swagger Documentation

Generate and view API documentation:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
make swagger
# Or:
swag init -g cmd/api/main.go -o api/

# Docs will be available at:
# http://localhost:8080/swagger/index.html (when implemented)
```

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entry point
│       └── main.go
├── internal/
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # Gin middleware (logging, recovery, request ID)
│   ├── models/          # Data models (Device, Error, Health)
│   └── server/          # HTTP server setup and configuration
├── tests/
│   └── integration/     # Integration tests
├── api/                 # Swagger API documentation
├── deployments/         # Dockerfile and deployment configs
├── Makefile            # Build and development commands
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## Middleware

The API includes the following middleware chain:

1. **Request ID**: Generates unique IDs for each request for tracing
2. **Logger**: Structured logging of all HTTP requests with request ID, method, path, status, and duration
3. **Recovery**: Panic recovery with stack trace logging

All requests automatically include these middleware features.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues, questions, or contributions, please open an issue on GitHub.