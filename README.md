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
A RESTful API service for managing and monitoring home lab infrastructure running on Kubernetes. This service provides unified access to various home automation and infrastructure services, starting with HomeAssistant integration.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Development Setup](#development-setup)
- [Building and Running](#building-and-running)
- [Testing](#testing)
- [Deployment](#deployment)
- [Configuration](#configuration)
- [Project Structure](#project-structure)
- [Contributing](#contributing)

## 🎯 Overview

The Home Lab API Service is designed to provide a unified REST API for interacting with various services running in a self-hosted Kubernetes cluster. The initial implementation focuses on HomeAssistant integration with mocked data, providing a proof-of-concept that can be extended to other home lab services.

### Key Goals

- **Unified API**: Single REST API for multiple home lab services
- **Kubernetes Native**: Designed for deployment on K8s with health checks and graceful shutdown
- **Modular Architecture**: Reusable internal packages that can be wrapped as MCP tools
- **Production Ready**: Structured logging, request tracing, error handling, and recovery middleware
- **Mock-First**: POC uses mocked data to validate API design before live integrations

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Client Applications                      │
│            (Dashboards, Scripts, Monitoring Tools)           │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Home Lab API Service                       │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              Gin HTTP Server (Port 8080)                │ │
│  └────────────────────────────────────────────────────────┘ │
│                         │                                    │
│  ┌─────────────────────┴────────────────────────────┐       │
│  │           Middleware Chain                        │       │
│  │  • Request ID Generation                          │       │
│  │  • Structured Logging (slog)                      │       │
│  │  • Panic Recovery                                 │       │
│  └───────────────────────────────────────────────────┘       │
│                         │                                    │
│  ┌─────────────────────┴────────────────────────────┐       │
│  │              Route Handlers                       │       │
│  │  • Health Check (/health)                         │       │
│  │  • API v1 Routes (/api/v1/*)                      │       │
│  │    - HomeAssistant endpoints                      │       │
│  │    - Service discovery                            │       │
│  │    - Cluster info                                 │       │
│  └───────────────────────────────────────────────────┘       │
│                         │                                    │
│  ┌─────────────────────┴────────────────────────────┐       │
│  │         Internal Business Logic                   │       │
│  │  • Models (Device, Health, Error)                 │       │
│  │  • Health Checker                                 │       │
│  │  • Response Handlers                              │       │
│  └───────────────────────────────────────────────────┘       │
│                         │                                    │
│  ┌─────────────────────┴────────────────────────────┐       │
│  │      Service Integrations (Mocked)                │       │
│  │  • HomeAssistant Service                          │       │
│  │  • Cluster Service                                │       │
│  │  (Future: Real integrations)                      │       │
│  └───────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

- **Language**: Go 1.25
- **Web Framework**: Gin (v1.12.0)
- **Logging**: Standard library `log/slog`
- **API Documentation**: Swagger/OpenAPI (swaggo)
- **Testing**: Go testing + testify
- **Container**: Docker
- **Orchestration**: Kubernetes

## ✨ Features

### Current Features (v1)

- ✅ RESTful API with versioned endpoints (`/api/v1`)
- ✅ Health check endpoint for K8s probes
- ✅ Structured logging with request tracing
- ✅ Request ID generation and propagation
- ✅ Panic recovery middleware
- ✅ Graceful shutdown (30s timeout)
- ✅ Thread-safe server operations
- ✅ Consistent error responses
- ✅ CORS support preparation
- ✅ Mocked HomeAssistant device data

### Planned Features

- 🔄 HomeAssistant device query endpoints
- 🔄 HomeAssistant device control endpoints
- 🔄 Service discovery endpoint
- 🔄 Rate limiting
- 🔄 Authentication/Authorization
- 🔄 Live HomeAssistant integration
- 🔄 Additional service integrations
- 🔄 MCP tool wrappers

## 🔌 API Endpoints

### Health Check

**GET /health**

Returns the health status of the service.

**Response**: 200 OK
```json
{
  "status": "ok"
}
```

**Use Case**: Kubernetes liveness and readiness probes

---

### API Version 1

**Base Path**: `/api/v1`

**GET /api/v1**

Returns API version information.

**Response**: 200 OK
```json
{
  "message": "API v1"
}
```

---

### HomeAssistant Endpoints (Planned)

**GET /api/v1/homeassistant/devices**

List all HomeAssistant devices (mocked data).

**Response**: 200 OK
```json
[
  {
    "id": "light.living_room",
    "name": "Living Room Light",
    "type": "light",
    "state": "on",
    "attributes": {
      "brightness": 255,
      "color_temp": 370
    },
    "last_updated": "2026-03-01T15:00:00Z",
    "controllable": true
  }
]
```

**GET /api/v1/homeassistant/devices/{id}**

Get a specific device by ID.

**Response**: 200 OK (device found) or 404 Not Found

**POST /api/v1/homeassistant/devices/{id}/command**

Send a command to a device.

**Request Body**:
```json
{
  "action": "turn_on",
  "parameters": {}
}
```

**Response**: 200 OK (success) or 400 Bad Request (invalid command)

---

### Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

**HTTP Status Codes**:
- `200` - Success
- `201` - Created
- `400` - Bad Request (invalid input)
- `404` - Not Found
- `405` - Method Not Allowed
- `429` - Too Many Requests (rate limit)
- `500` - Internal Server Error
- `503` - Service Unavailable

## 🚀 Development Setup

### Prerequisites

- **Go 1.25+**: [Download](https://golang.org/dl/)
- **Make**: For build automation
- **Docker**: For containerization (optional)
- **golangci-lint**: For code linting (optional)
- **air**: For hot reload during development (optional)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/rmwondolleck/go-github.git
   cd go-github
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Verify installation**:
   ```bash
   go version
   go mod verify
   ```

### Environment Variables

The service supports the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |

Set environment variables:
```bash
export PORT=8080
```

## 🏃 Building and Running

### Using Make (Recommended)

The project includes a Makefile with common tasks:

```bash
# Show available commands
make help

# Build the application
make build

# Run tests
make test

# Run the application
make run

# Run with hot reload (requires air)
make dev

# Run linter
make lint

# Generate Swagger documentation
make swagger

# Build Docker image
make docker

# Clean build artifacts
make clean
```

### Manual Build

```bash
# Build binary
go build -o bin/homelab-api ./cmd/api

# Run the binary
./bin/homelab-api
```

### Running Directly

```bash
# Run without building
go run ./cmd/api
```

The service will start on port 8080 (or the port specified by `PORT` environment variable).

**Verify the service**:
```bash
curl http://localhost:8080/health
```

## 🧪 Testing

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
# Run tests with coverage
make test

# Or manually
go test -v -race -coverprofile=coverage.out ./...
```

### Run Specific Tests

```bash
# Test a specific package
go test -v ./internal/middleware

# Run a specific test
go test -v -run TestRequestID ./internal/middleware
```

### View Coverage Report

```bash
# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### Integration Tests

Integration tests are located in `tests/integration/`:

```bash
go test -v ./tests/integration/...
```

## 📦 Deployment

### Docker Deployment

1. **Build Docker image**:
   ```bash
   make docker
   # or
   docker build -t homelab-api:latest -f deployments/Dockerfile .
   ```

2. **Run container**:
   ```bash
   docker run -p 8080:8080 homelab-api:latest
   ```

### Kubernetes Deployment

The project includes Kubernetes manifests in `deployments/k8s/`.

1. **Apply deployment**:
   ```bash
   kubectl apply -f deployments/k8s/deployment.yaml
   ```

2. **Verify deployment**:
   ```bash
   kubectl get pods -l app=homelab-api
   kubectl logs -l app=homelab-api
   ```

3. **Check health**:
   ```bash
   kubectl port-forward deployment/homelab-api 8080:8080
   curl http://localhost:8080/health
   ```

### Kubernetes Configuration

The deployment includes:

- **Replicas**: 2 pods for high availability
- **Resource Limits**: 
  - Memory: 100Mi limit, 50Mi request
  - CPU: 200m limit, 100m request
- **Liveness Probe**: `/health` endpoint checked every 10s
- **Readiness Probe**: `/health` endpoint checked every 5s
- **Port**: 8080 (HTTP)

### Production Considerations

- [ ] Configure resource limits based on actual usage
- [ ] Set up horizontal pod autoscaling (HPA)
- [ ] Configure ingress for external access
- [ ] Add TLS/SSL certificates
- [ ] Set up monitoring and alerting
- [ ] Configure log aggregation
- [ ] Implement authentication/authorization
- [ ] Enable rate limiting
- [ ] Configure CORS for allowed origins

## ⚙️ Configuration

### Build Configuration

The binary is built with:
```bash
go build -o bin/homelab-api ./cmd/api
```

### Server Configuration

The server uses:
- **Framework**: Gin in release mode (production)
- **Middleware Chain**:
  1. Request ID generation
  2. Structured logging
  3. Panic recovery
- **Graceful Shutdown**: 30-second timeout for SIGTERM/SIGINT

### Logging

The service uses Go's standard `log/slog` for structured logging:

```go
slog.Info("message", "key", "value")
slog.Error("error", "error", err)
```

All logs include:
- Request ID (for tracing)
- Timestamp
- Log level
- Structured fields

## 📁 Project Structure

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
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   └── response.go          # Response helpers
│   ├── health/                  # Health check service
│   │   └── checker.go
│   ├── homeassistant/           # HomeAssistant integration
│   ├── cluster/                 # Cluster service integration
│   ├── middleware/              # HTTP middleware
│   │   ├── logging.go           # Request logging
│   │   ├── recovery.go          # Panic recovery
│   │   └── request_id.go        # Request ID generation
│   ├── models/                  # Data models
│   │   ├── device.go            # Device model
│   │   ├── error.go             # Error response model
│   │   └── health.go            # Health status model
│   └── server/                  # HTTP server
│       ├── server.go            # Server setup and routing
│       └── shutdown.go          # Graceful shutdown
├── tests/
│   └── integration/             # Integration tests
├── deployments/
│   └── k8s/
│       └── deployment.yaml      # Kubernetes deployment
├── api/                         # API documentation (Swagger)
├── research/                    # Research and benchmarks
├── .github/
│   └── workflows/               # CI/CD workflows
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
└── README.md                    # This file
```

### Internal Package Guidelines

- **`/internal`**: Private application code (not importable by external projects)
- **`/cmd`**: Application entry points
- **`/deployments`**: Deployment configurations
- **`/tests`**: Integration and end-to-end tests
- **`/api`**: API documentation and schemas

## 🤝 Contributing

### Development Workflow

1. Create a feature branch
2. Make changes with tests
3. Run tests and linter
4. Submit pull request

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new features
- Update documentation

### Testing Requirements

- Unit tests for new functions
- Integration tests for new endpoints
- Maintain >80% code coverage

### Commit Messages

Follow conventional commit format:
```
type(scope): subject

body

footer
```

Example:
```
feat(api): add HomeAssistant device listing endpoint

Implements GET /api/v1/homeassistant/devices with mocked data.
Returns array of device objects with status and attributes.

Closes #123
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Gin Web Framework](https://github.com/gin-gonic/gin)
- API documentation with [Swaggo](https://github.com/swaggo/swag)
- Designed for Kubernetes deployment

---

**Project Status**: 🚧 In Development (POC Phase)

For questions or issues, please open an issue on GitHub.
