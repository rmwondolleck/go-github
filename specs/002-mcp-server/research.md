# Research: MCP Server Integration

**Feature**: 002-mcp-server  
**Date**: March 14, 2026  
**Status**: Complete

---

## R1: Integrated Architecture — Single Binary with MCP Subcommand

**Decision**: Use a single binary (`homelab-api`) with a subcommand/flag to select the mode of operation: `homelab-api` (default: HTTP API), `homelab-api mcp` (MCP stdio server). Both modes share the same codebase and compiled binary.

**Rationale**: The user constraint is explicit: "the MCP server and API server should run on the same codebase/application — they should be integrated into a single application, not separate binaries." A subcommand approach satisfies this while keeping MCP's stdio requirement (stdin/stdout communication) separate from the HTTP server's port-based communication. Both modes are incompatible at runtime (HTTP writes to stdout for logs; MCP uses stdout for JSON-RPC), so they run as separate modes of the same binary.

**Alternatives Considered**:
- **Separate `cmd/mcp/main.go` binary**: Rejected — violates the user constraint of a single application. Would require two build targets and two binaries.
- **Concurrent HTTP + MCP in one process**: Rejected — MCP stdio requires exclusive control of stdout. Running an HTTP server simultaneously would corrupt MCP's JSON-RPC output stream.
- **Environment variable mode selection** (`MODE=mcp`): Considered viable but less idiomatic than a subcommand. Subcommands are the standard Go pattern (e.g., `go build`, `go test`).

**Implementation**: Modify `cmd/api/main.go` to check `os.Args` for a `mcp` subcommand. If present, start the MCP stdio server instead of the HTTP server. The `internal/mcp/` package provides the server construction. All logs in MCP mode go to stderr.

---

## R2: mcp-go SDK API Patterns

**Decision**: Use `github.com/mark3labs/mcp-go` SDK with the `server.NewMCPServer()`, `server.NewStdioServer()`, and the `mcp` types package for resources, tools, and prompts.

**Rationale**: The mcp-go SDK is the most mature Go implementation of MCP. It handles JSON-RPC framing, session lifecycle, capability negotiation, and stdio transport — all concerns that would be error-prone to implement from scratch. The SDK is actively maintained by mark3labs.

**Key API Surface** (from SDK research):

```go
// Create MCP server with capabilities
s := server.NewMCPServer(
    "server-name", "1.0.0",
    server.WithResourceCapabilities(true, true),
    server.WithToolCapabilities(true),
    server.WithPromptCapabilities(true),
)

// Register resource
s.AddResource(mcp.Resource{
    URI:         "homelab://devices",
    Name:        "devices",
    Description: "Home automation device catalogue",
    MIMEType:    "application/json",
}, handlerFunc)
// Handler signature: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)

// Register tool
s.AddTool(mcp.Tool{
    Name:        "execute_command",
    Description: "Execute a device command",
    InputSchema: mcp.ToolInputSchema{...},
}, toolHandlerFunc)
// Handler signature: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

// Register prompt
s.AddPrompt(mcp.Prompt{
    Name:        "device_control",
    Description: "Guide AI to control a device",
    Arguments:   []mcp.PromptArgument{{Name: "device_name", Required: true}},
}, promptHandlerFunc)
// Handler signature: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)

// Start stdio transport
stdio := server.NewStdioServer(s)
stdio.Listen(ctx, os.Stdin, os.Stdout)
```

**Alternatives Considered**:
- **Raw JSON-RPC implementation**: Rejected — substantial boilerplate for protocol framing, session management, and capability negotiation.
- **Other MCP Go libraries**: None as mature as mcp-go as of this date.

---

## R3: Shared Data Access — Extracting Mock Data into Reusable Providers

**Decision**: Extract mock data (devices, services) from HTTP handlers into shared provider packages under `internal/`. The MCP handlers and HTTP handlers both consume these providers.

**Rationale**: The spec requires FR-007 ("serve data from existing home lab data sources without duplicating that data") and FR-012 ("device list sourced from the same mock device store used by the device command action"). Currently, mock data lives inside `internal/handlers/homeassistant.go` (`mockDevices` map) and `internal/handlers/services.go` (inline services slice). These need to be extracted to shared packages so MCP handlers can access them without importing the HTTP handler package.

**Extraction Plan**:
- `internal/homeassistant/devices.go` — export `GetDevices()`, `GetDevice(id)`, and `ExecuteCommand(deviceID, cmd)` functions; move `mockDevices` here
- `internal/services/provider.go` — extract service list with `GetServices()` function
- HTTP handlers become thin wrappers calling these providers
- MCP resource/tool handlers call the same providers

**Alternatives Considered**:
- **MCP handlers importing `internal/handlers`**: Rejected — creates coupling between MCP and HTTP layers. The handlers package uses gin-specific types.
- **Duplicating mock data**: Rejected — violates FR-007 and FR-012.

---

## R4: MCP Logging Strategy

**Decision**: All MCP server logging goes to stderr via `log/slog`. No output goes to stdout except JSON-RPC messages.

**Rationale**: MCP stdio transport requires stdout to be exclusively used for JSON-RPC communication. Any non-JSON-RPC output on stdout corrupts the protocol stream. The mcp-go SDK's `StdioServer` writes to the provided stdout writer and reads from stdin. Diagnostic logging must go to stderr.

**Implementation**: Configure `slog` with a handler writing to `os.Stderr` when running in MCP mode. The SDK already handles JSON-RPC framing on stdout.

---

## R5: Testing Strategy for MCP Components

**Decision**: Use the mcp-go SDK's in-process testing capabilities combined with standard Go testing patterns.

**Rationale**: The mcp-go SDK provides in-process server/client for testing MCP servers without real stdio pipes. This allows unit tests to send MCP requests and validate responses programmatically.

**Testing Approach**:
- **Unit tests** for each handler function (resources, tools, prompts) — test with mock data, validate response structure
- **Server integration tests** — create an `MCPServer`, register handlers, send in-process requests, validate end-to-end flow
- **Table-driven tests** — per constitution requirement, use table-driven patterns for parameterised test cases
- **Race detection** — all tests run with `-race` flag per existing `make test` configuration

**Target Coverage**: ≥80% for `internal/mcp/` package (per constitution and SC-005)

---

## R6: Subcommand Implementation — Minimal Approach

**Decision**: Use simple `os.Args` inspection for subcommand detection rather than a CLI framework.

**Rationale**: The application has exactly two modes: default (HTTP API) and `mcp`. A full CLI framework (cobra, urfave/cli) adds unnecessary dependency weight for a binary choice. The constitution emphasises "minimal external dependencies."

**Implementation**:
```go
func main() {
    if len(os.Args) > 1 && os.Args[1] == "mcp" {
        runMCP()
        return
    }
    runHTTPServer()
}
```

**Alternatives Considered**:
- **cobra/urfave/cli**: Rejected — overkill for a single subcommand. Constitution says "minimal external dependencies."
- **Flag-based** (`--mode=mcp`): Viable but less readable. Subcommands are more discoverable.

---

## R7: Resource URI Scheme

**Decision**: Use `homelab://` as the URI scheme for all MCP resources, matching the spec.

**Resources**:
| URI | Description | Data Source |
|-----|-------------|-------------|
| `homelab://devices` | Home automation devices | `internal/homeassistant/devices.go` |
| `homelab://services` | Home lab services | `internal/services/provider.go` |
| `homelab://cluster/services` | Kubernetes cluster services | `internal/cluster/service.go` |
| `homelab://health` | System health status | `internal/health/checker.go` |

---

## R8: Build Tooling — Makefile Targets

**Decision**: Since MCP and HTTP share a single binary, adjust existing Makefile targets and add MCP-specific convenience targets.

**New/Modified Targets**:
- `make build` — builds `bin/homelab-api` (unchanged, same binary serves both modes)
- `make mcp-run` — runs `bin/homelab-api mcp` (convenience target for MCP mode)
- `make mcp-build` — alias for `make build` (for quickstart compatibility; outputs same binary)

**Rationale**: The quickstart doc references `make mcp-build` and `make mcp-run`. Since it's a single binary, `mcp-build` is just an alias for `build`. The MCP configuration in `.vscode/mcp.json` will point to `bin/homelab-api` with args `["mcp"]`.
