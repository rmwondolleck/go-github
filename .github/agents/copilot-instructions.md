# go-github Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-14

## Active Technologies

- Go 1.25.0
- Gin (HTTP framework)
- github.com/mark3labs/mcp-go (MCP protocol SDK — 002-mcp-server)
- testify (test assertions)
- log/slog (structured logging)

## Project Structure

```text
cmd/api/main.go              # Single binary entrypoint (HTTP API default, MCP via subcommand)
internal/handlers/            # HTTP API handlers (Gin)
internal/mcp/                 # MCP server package (resources, tools, prompts)
internal/homeassistant/       # Home automation types and shared device provider
internal/services/            # Shared service provider
internal/cluster/             # Cluster service logic
internal/health/              # Health checker
internal/models/              # Domain models (Device, Service, HealthStatus, ErrorResponse)
internal/middleware/           # HTTP middleware (CORS, logging, recovery, request ID)
internal/server/              # HTTP server setup
api/                          # OpenAPI/Swagger docs
tests/                        # Integration tests
specs/                        # Feature specifications
```

## Commands

```bash
make build       # Build single binary to bin/homelab-api
make test        # Run all tests with race detection and coverage
make run         # Run HTTP API server
make mcp-run     # Run MCP stdio server (bin/homelab-api mcp)
make mcp-build   # Alias for make build
make lint        # Run golangci-lint
make swagger     # Generate Swagger docs
```

## Code Style

Go: Follow official Go Code Review Comments guidelines. Use gofmt, go vet, golangci-lint. Handle all errors explicitly. Use context.Context for cancellation. Structured logging via log/slog. Table-driven tests preferred.

## Key Patterns

- **MCP mode**: `./bin/homelab-api mcp` — starts stdio MCP server, logs to stderr only
- **HTTP mode**: `./bin/homelab-api` (default) — starts HTTP API on PORT (default 8080)
- **Shared providers**: Device and service data accessed via `internal/homeassistant/devices.go` and `internal/services/provider.go` — used by both HTTP handlers and MCP handlers
- **No data duplication**: MCP resources read from the same mock data stores as HTTP endpoints

## Recent Changes

- 002-mcp-server: MCP server integration as subcommand of single binary. New dependency: mcp-go. New packages: internal/mcp/, internal/services/. Refactored: shared device provider in internal/homeassistant/devices.go.

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
