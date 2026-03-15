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
cmd/api/main.go              # Single binary entrypoint — starts HTTP API + MCP concurrently
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
.github/docs/                 # Project status, ops, and history documents
.github/agents/               # Copilot agent instructions and custom agents
.github/prompts/              # speckit prompt templates
.github/workflows/            # Agentic and CI workflows
```

## Commands

```bash
make build    # Build single binary to bin/homelab-api
make test     # Run all tests with race detection and coverage
make run      # Run binary — starts HTTP API (port 8080) + MCP stdio concurrently
make lint     # Run golangci-lint
make swagger  # Generate Swagger docs
```

## Code Style

Go: Follow official Go Code Review Comments guidelines. Use gofmt, go vet, golangci-lint. Handle all errors explicitly. Use context.Context for cancellation. Structured logging via log/slog. Table-driven tests preferred.

## Key Patterns

- **Single command**: `./bin/homelab-api` starts **both** HTTP API (port 8080) and MCP stdio server concurrently — no subcommand needed
- **Shared providers**: Device and service data accessed via `internal/homeassistant/devices.go` and `internal/services/provider.go` — used by both HTTP handlers and MCP handlers
- **No data duplication**: MCP resources read from the same mock data stores as HTTP endpoints
- **Graceful shutdown**: Both HTTP and MCP modes share a context cancelled on SIGINT/SIGTERM via errgroup
- **MCP logging**: All MCP diagnostic output goes to stderr; stdout is reserved for JSON-RPC messages

## Documentation

All status, ops, and history documents live in `.github/docs/`. Do not create new `.md` files in the repository root — use `.github/docs/` for operational docs and `specs/` for feature specifications.

| Document | Path |
|---|---|
| Current project status | `.github/docs/CURRENT_STATUS.md` |
| Doc index | `.github/docs/DOC_INDEX.md` |
| PR status | `.github/docs/PR_STATUS_REPORT.md` |
| Mission / history | `.github/docs/MISSION_COMPLETE.md` |

## Recent Changes

- 002-mcp-server: MCP server integration into single binary. Both HTTP API and MCP stdio start concurrently with `./bin/homelab-api`. New dependency: mcp-go. New packages: internal/mcp/, internal/services/. Refactored: shared device provider in internal/homeassistant/devices.go.
- Housekeeping: All non-core `.md` files moved from root to `.github/docs/`.

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
