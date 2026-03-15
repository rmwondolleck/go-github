# Implementation Plan: MCP Server Integration

**Branch**: `002-mcp-server` | **Date**: March 14, 2026 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-mcp-server/spec.md`

## Summary

Integrate a Model Context Protocol (MCP) server into the existing go-github home lab API application as a subcommand of the single binary. When invoked as `homelab-api mcp`, the application starts an MCP stdio server (using the `github.com/mark3labs/mcp-go` SDK) that exposes home lab data (devices, services, cluster services, health) as MCP resources, device control as an MCP tool, and AI interaction templates as MCP prompts. The MCP server reuses all existing domain logic via shared data providers extracted from HTTP handlers. No separate binary is produced.

## Technical Context

**Language/Version**: Go 1.25.0  
**Primary Dependencies**: gin (HTTP, existing), github.com/mark3labs/mcp-go (MCP SDK, new)  
**Storage**: N/A (in-memory mock data, shared providers)  
**Testing**: Go `testing` package, `testify` for assertions, mcp-go in-process test utilities  
**Target Platform**: Linux server, macOS, Windows (local development)  
**Project Type**: Web service + MCP server (single binary, dual mode)  
**Performance Goals**: Handshake + enumerate < 2s, resource reads < 500ms, tool execution < 1s  
**Constraints**: MCP stdout exclusive (no HTTP output in MCP mode), stderr for all logging  
**Scale/Scope**: Single-user per MCP process (stdio), mock data

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Constitution Rule | Status | Notes |
|---|---|---|---|
| I | **Go Standards Compliance** — Go 1.24+, gofmt, go vet, golangci-lint, explicit error handling, context.Context | ✅ PASS | Go 1.25, all errors handled, context passed through SDK handlers |
| II | **API Design** — RESTful, JSON, consistent errors, versioning, OpenAPI | ✅ PASS (N/A for MCP) | MCP uses JSON-RPC, not REST. Existing HTTP API unchanged. MCP contract documented in contracts/ |
| III | **Testing (NON-NEGOTIABLE)** — 80% coverage, unit + integration, table-driven, testify allowed | ✅ PASS | Target ≥80% for internal/mcp/. Table-driven tests. testify used. |
| IV | **Observability** — log/slog, log levels, request ID, health endpoint | ✅ PASS | slog to stderr in MCP mode. Existing /health endpoint preserved. MCP health resource added. |
| V | **Security & Configuration** — no hardcoded secrets, env vars, graceful shutdown | ✅ PASS | No secrets. Graceful shutdown on stdin close. SIGTERM/SIGINT handled. |
| VI | **Build Requirements** — go.mod committed, Make targets, Docker, version info | ✅ PASS | go.mod updated with mcp-go. New Make targets added. Same binary. |
| VII | **Code Organization** — /cmd, /internal, /pkg, /api, /tests | ✅ PASS | New code in internal/mcp/. Entry point in cmd/api/main.go (subcommand). |
| VIII | **Dependency Management** — minimal deps, pinned versions, justified | ✅ PASS | One new dep: mcp-go. Justified: SDK handles MCP protocol complexity. |

**Post-Phase-1 Re-check**: All gates still pass. No new violations introduced by the design.

## Project Structure

### Documentation (this feature)

```text
specs/002-mcp-server/
├── plan.md              # This file
├── research.md          # Phase 0 output — technology decisions
├── data-model.md        # Phase 1 output — entity definitions
├── quickstart.md        # Phase 1 output — build/run/connect guide
├── contracts/
│   └── mcp-protocol.md  # Phase 1 output — JSON-RPC contract
├── checklists/
│   └── requirements.md  # Requirements tracking
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/api/
└── main.go              # Modified: add "mcp" subcommand dispatch

internal/mcp/            # NEW: MCP server package
├── server.go            # NewMCPServer(), registration, Run()
├── server_test.go       # Server creation and registration tests
├── resources.go         # Resource handlers (devices, services, cluster, health)
├── resources_test.go    # Resource handler tests
├── tools.go             # Tool handler (execute_command)
├── tools_test.go        # Tool handler tests
├── prompts.go           # Prompt handlers (device_control, service_status)
└── prompts_test.go      # Prompt handler tests

internal/homeassistant/
├── types.go             # EXISTING: Command type
├── types_test.go        # EXISTING
└── devices.go           # NEW: exported GetDevices(), GetDevice(), ExecuteCommand()

internal/services/       # NEW: shared service provider
└── provider.go          # GetServices() — extracted from handlers/services.go

internal/handlers/       # MODIFIED: refactored to use shared providers
├── homeassistant.go     # Refactored: delegates to homeassistant.GetDevice/ExecuteCommand
└── services.go          # Refactored: delegates to services.GetServices()

.vscode/
└── mcp.json             # NEW: Copilot MCP server configuration

Makefile                 # MODIFIED: add mcp-build, mcp-run targets
```

**Structure Decision**: Single project structure. MCP code lives in `internal/mcp/` alongside existing internal packages. Shared data providers are extracted from `internal/handlers/` into domain-specific packages (`internal/homeassistant/`, `internal/services/`). Entry point remains `cmd/api/main.go` with subcommand dispatch. No new top-level directories.

## Complexity Tracking

No constitution violations to justify. The design adds one new external dependency (mcp-go) and one new internal package (internal/mcp/), both well within the project's existing patterns.
