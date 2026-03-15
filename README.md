# Home Lab API Service

A RESTful API service for managing and monitoring home lab infrastructure running on Kubernetes. This service provides unified access to various home automation and infrastructure services, starting with HomeAssistant integration.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Coverage](https://img.shields.io/badge/coverage-96.1%25-brightgreen.svg)](COVERAGE_REPORT.md)
[![Constitution Compliance](https://img.shields.io/badge/compliance-passing-brightgreen.svg)](COMPLIANCE_REPORT.md)

## 🚀 Current Status (March 14, 2026)

**Project Status**: ✅ Feature Complete — HTTP API + MCP Server

### Recent Updates
- ✅ **MCP Server Integration Complete**: Full Model Context Protocol server integrated into the binary — both HTTP API and MCP stdio start with a single command
- ✅ **Constitution Compliance**: Test coverage at 93.5% (target ≥80%), `gofmt` clean, `go vet` clean
- ✅ **Performance Validated**: 100-goroutine load tests pass, P99 <300µs, no memory leaks
- ✅ **Docker & Kubernetes Ready**: Distroless image at 34.1MB, K8s manifests with liveness/readiness probes

See [`.github/docs/CURRENT_STATUS.md`](.github/docs/CURRENT_STATUS.md) for full status.

## ✅ Constitution Compliance

All project requirements have been validated. See [COMPLIANCE_REPORT.md](./COMPLIANCE_REPORT.md) for full details.

| Requirement | Status |
|-------------|--------|
| Go 1.24+ standards | ✅ Go 1.25 |
| 80%+ test coverage | ✅ 93.5% |
| All errors handled | ✅ Structured error responses |
| Structured logging | ✅ `log/slog` used throughout |
| Graceful shutdown | ✅ OS signal handling with 5s timeout |
| Code formatting | ✅ `gofmt` clean |
| Lints passing | ✅ `go vet` clean |

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [API Documentation](#api-documentation)
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
│  │  • Swagger UI (/api/docs/*)                       │       │
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
- ✅ CORS middleware
- ✅ Graceful shutdown (30s timeout)
- ✅ Thread-safe server operations
- ✅ Consistent error responses
- ✅ Mocked HomeAssistant device data and control endpoints
- ✅ Service discovery endpoint
- ✅ Cluster services endpoint
- ✅ Interactive API documentation with Swagger/OpenAPI
- ✅ **MCP Server** — AI assistant integration via Model Context Protocol (resources, tools, prompts)

### Planned Features

- 🔄 Rate limiting
- 🔄 Authentication/Authorization
- 🔄 Live HomeAssistant integration (replace mock data)
- 🔄 Additional service integrations

## 📚 API Documentation

The API is documented using OpenAPI/Swagger specification. Once the server is running, you can access:

### Swagger UI

**URL**: [http://localhost:8080/api/docs/index.html](http://localhost:8080/api/docs/index.html)

The Swagger UI provides:
- 📖 Interactive API documentation
- 🧪 "Try it out" functionality to test endpoints directly in your browser
- 📋 Full endpoint specifications with request/response schemas
- 🔍 Parameter descriptions and validation rules

### OpenAPI Specification

**URL**: [http://localhost:8080/api/docs/doc.json](http://localhost:8080/api/docs/doc.json)

The OpenAPI spec can be:
- 💾 Downloaded in JSON format
- 📤 Imported into API testing tools (Postman, Insomnia, etc.)
- 🔗 Used for client SDK generation
- 📝 Integrated into CI/CD pipelines

### Generating Swagger Documentation

After modifying API endpoints or annotations, regenerate the documentation:

```bash
# Using Make
make swagger

# Or manually
swag init -g cmd/api/main.go -o api/
```

The generated files (`api/docs.go`, `api/swagger.json`, `api/swagger.yaml`) are gitignored as they are build artifacts.

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

### HomeAssistant Endpoints

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
# Check health endpoint
curl http://localhost:8080/health

# Access Swagger UI in your browser
open http://localhost:8080/api/docs/index.html
```

## 🤖 MCP Server (AI Assistant Integration)

The `homelab-api` binary runs **both** the HTTP API server and an MCP (Model Context Protocol) stdio server **concurrently** with a single command. No subcommand is needed — just run the binary and both modes start automatically.

### How It Works

The binary supports **three modes**:

| Command | HTTP API | MCP stdio | Use case |
|---|---|---|---|
| `./bin/homelab-api` | ✅ port 8080 | ✅ stdin/stdout | Default / Kubernetes |
| `./bin/homelab-api mcp` | ❌ | ✅ stdin/stdout | Local IDE / AI assistant |

Both modes handle `SIGINT`/`SIGTERM` cleanly. In Kubernetes the MCP server runs but idles harmlessly (stdin is `/dev/null`).

### Build & Run

```bash
# Build the dual-mode binary
make build
# or: go build -o bin/homelab-api ./cmd/api

# Mode 1 — default: starts BOTH HTTP API and MCP stdio concurrently
./bin/homelab-api

# Mode 2 — MCP only: no HTTP port bound (what your IDE uses)
./bin/homelab-api mcp
```

### MCP Quick Smoke Test

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1"}}}' \
  | ./bin/homelab-api mcp 2>/dev/null
```

Expected response contains `"name":"go-github-homelab"`.

### VS Code Copilot Configuration

The `.vscode/mcp.json` file is already included in the repository:

```json
{
  "servers": {
    "go-github-homelab": {
      "type": "stdio",
      "command": "${workspaceFolder}/bin/homelab-api",
      "args": ["mcp"]
    }
  }
}
```

VS Code spawns the binary in MCP-only mode (`mcp` arg) — no HTTP port is bound.

### JetBrains AI Configuration

In JetBrains IDEs (GoLand, IDEA, etc.) with the AI plugin:

1. Open **Settings → Tools → AI Assistant → MCP Servers**
2. Click **+** to add a new server
3. Set:
   - **Name**: `go-github-homelab`
   - **Command**: path to `bin/homelab-api`
   - **Args**: `mcp`
   - **Transport**: `stdio`
   - **Transport**: `stdio`

### Available Resources, Tools, and Prompts

| Type | Name | URI / Description |
|------|------|-------------------|
| Resource | Devices | `homelab://devices` — all HA smart home devices |
| Resource | Services | `homelab://services` — homelab services (prometheus, grafana, etc.) |
| Resource | Cluster Services | `homelab://cluster/services` — Kubernetes cluster services |
| Resource | Health | `homelab://health` — API health and uptime |
| Tool | execute_command | Execute a control command on a device (`device_id`, `action`) |
| Prompt | device_control | Rendered prompt for controlling a named device |
| Prompt | service_status | Rendered prompt for checking a service's status |

---

## 🧪 Testing

### Coverage Summary

| Package | Coverage |
|---------|----------|
| `internal/cluster` | 100.0% |
| `internal/handlers` | 87.8% |
| `internal/health` | 100.0% |
| `internal/homeassistant` | 100.0% |
| `internal/middleware` | 100.0% |
| `internal/server` | 100.0% |
| **Total** | **96.1%** ✅ |

> ✅ Coverage exceeds the 80% target. See [COVERAGE_REPORT.md](COVERAGE_REPORT.md) for the full report.

### Run All Tests

```bash
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

### Swagger Verification

Verify Swagger UI and API documentation:

```bash
# Run verification script
./tests/verify_swagger.sh
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

3. **Verify deployment**:
   ```bash
   curl http://localhost:8080/health
   ```

For detailed deployment instructions, see [deployments/README.md](deployments/README.md).

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

For complete Kubernetes deployment guide, see [deployments/README.md](deployments/README.md).

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
│   └── api/
│       └── main.go              # Application entry point — launches HTTP API + MCP concurrently
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   └── response.go          # Response helpers
│   ├── health/                  # Health check service
│   │   └── checker.go
│   ├── homeassistant/           # HomeAssistant integration + shared device provider
│   │   └── devices.go           # GetDevices(), GetDevice(), ExecuteCommand()
│   ├── services/                # Shared service provider
│   │   └── provider.go          # GetServices()
│   ├── mcp/                     # MCP server (resources, tools, prompts)
│   │   ├── server.go            # NewMCPServer(), Run()
│   │   ├── resources.go         # Resource handlers (devices, services, cluster, health)
│   │   ├── tools.go             # Tool handler (execute_command)
│   │   └── prompts.go           # Prompt handlers (device_control, service_status)
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
│   ├── integration/             # Integration tests
│   └── verify_swagger.sh        # Swagger UI verification
├── deployments/
│   ├── Dockerfile               # Docker image definition (Alpine)
│   ├── Dockerfile.distroless    # Docker image definition (Distroless, 34.1MB)
│   ├── README.md                # Deployment documentation
│   └── k8s/
│       ├── deployment.yaml      # Kubernetes deployment
│       └── configmap.yaml       # Kubernetes config
├── api/                         # API documentation (Swagger)
│   ├── docs.go                  # Generated Swagger docs
│   ├── swagger.json             # Generated OpenAPI spec
│   └── swagger.yaml             # Generated OpenAPI spec
├── specs/                       # Feature specifications
│   └── 002-mcp-server/          # MCP server feature spec, plan, tasks, contracts
├── research/                    # Research and benchmarks
├── .github/
│   ├── agents/                  # Copilot agent instructions
│   │   └── copilot-instructions.md
│   ├── docs/                    # Project status and history documents
│   │   └── CURRENT_STATUS.md
│   ├── prompts/                 # speckit prompt templates
│   └── workflows/               # CI/CD and agentic workflows
├── .vscode/
│   └── mcp.json                 # VS Code Copilot MCP server configuration
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

**Project Status**: ✅ Feature Complete — HTTP API + MCP Server integrated

For questions or issues, please open an issue on GitHub.
