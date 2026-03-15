# Research: MCP Server Integration

**Feature**: 002-mcp-server  
**Date**: March 14, 2026  
**Status**: Complete

---

## R1: Integrated Architecture — Single Binary, Concurrent HTTP + MCP

**Decision**: Use a single binary (`homelab-api`) that starts **both** the HTTP API server and the MCP stdio server **concurrently** on every launch. Both modes run as goroutines under a shared `errgroup` context and shut down together on `SIGINT`/`SIGTERM`. No subcommand is required.

**Rationale**: The user constraint is explicit: "the MCP server and API server should run on the same codebase/application." A concurrent approach satisfies this fully while keeping the developer experience simple (one command, both modes active). The stdout conflict (HTTP logs vs MCP JSON-RPC output) is resolved by routing all non-MCP output — including HTTP server logs — to `stderr` via `log/slog`. Only MCP JSON-RPC messages go to `stdout`, preserving protocol correctness.

**Alternatives Considered**:
- **Separate `cmd/mcp/main.go` binary**: Rejected — violates the user constraint of a single application. Would require two build targets and two binaries.
- **Subcommand dispatch** (`homelab-api mcp`): Considered and prototyped. Rejected — unnecessarily splits a single logical server into two modes, complicates IDE configuration (args required), and makes k8s deployment reasoning harder. The stdout conflict is solvable without mode separation.
- **Environment variable mode selection** (`MODE=mcp`): Rejected — does not allow both modes to run simultaneously; same drawbacks as subcommand.

**Implementation**: Modify `cmd/api/main.go` to use `golang.org/x/sync/errgroup` to launch both `srv.Run(port)` (HTTP) and `internalmcp.Run(ctx)` (MCP) concurrently. All logs — from both modes — go to `stderr`. Only MCP JSON-RPC protocol messages go to `stdout`.

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

**Implementation**: Configure `slog` with a handler writing to `os.Stderr` at startup. Both the HTTP server and the MCP server write all diagnostic logs to stderr. The SDK handles JSON-RPC framing on stdout exclusively.

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

## R6: Concurrent Launch — errgroup Pattern

**Decision**: Use `golang.org/x/sync/errgroup` to manage both the HTTP server and MCP server goroutines under a shared cancellable context.

**Rationale**: `errgroup` propagates the first non-nil error from either goroutine to the caller and cancels the shared context, ensuring both modes shut down together. This is idiomatic Go for managing multiple concurrent long-running goroutines. No CLI framework is needed — `os.Args` is not inspected.

**Implementation**:
```go
g, gctx := errgroup.WithContext(ctx)

g.Go(func() error {
    slog.Info("http server started", "port", port)
    return srv.Run(port)
})

g.Go(func() error {
    slog.Info("mcp server started", "transport", "stdio")
    return internalmcp.Run(gctx)
})

g.Go(func() error {
    <-gctx.Done()
    return srv.GracefulShutdown(shutdownCtx)
})

return g.Wait()
```

**Alternatives Considered**:
- **cobra/urfave/cli**: Rejected — overkill. Constitution says "minimal external dependencies."
- **os.Args subcommand check**: Prototyped but rejected in favour of concurrent launch (see R1).

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

**Decision**: Since MCP and HTTP share a single binary and start together, no MCP-specific Makefile targets are needed. The existing `make build` and `make run` targets serve both modes.

**Targets**:
- `make build` — builds `bin/homelab-api` (unchanged; same binary serves both HTTP and MCP)
- `make run` — runs `./bin/homelab-api`; starts HTTP API on port 8080 **and** MCP stdio server concurrently
- `make test` — runs all tests including `internal/mcp/` package

**Rationale**: There is no separate MCP mode to invoke — both servers start automatically. Adding `mcp-run` or `mcp-build` aliases would imply a mode separation that does not exist and would confuse developers. The `.vscode/mcp.json` IDE configuration points to `bin/homelab-api` with **no args**.
