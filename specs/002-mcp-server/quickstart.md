# Quickstart: MCP Server Integration

**Feature**: 002-mcp-server  
**Date**: March 14, 2026

---

## Prerequisites

- Go 1.25.0 or later
- Make
- An MCP-compatible client (GitHub Copilot in VS Code or JetBrains)

## Build

```bash
# Build the single binary (serves both HTTP API and MCP simultaneously)
make build

# Verify the binary was created
ls -la bin/homelab-api
```

## Run (Manual Test)

The binary supports three modes:

```bash
# Mode 1 (default): starts BOTH HTTP API (port 8080) AND MCP stdio concurrently
./bin/homelab-api

# Mode 2 (MCP-only): starts MCP stdio server only — no HTTP port bound
# Use this when you only need local IDE/AI integration
./bin/homelab-api mcp

# Via Make
make run       # Mode 1 — both servers
```

**Mode 1** is the default and what runs in Kubernetes. Both servers share a context and both shut down on `SIGINT`/`SIGTERM`.

**Mode 2** (`mcp` arg) is ideal for local development — your IDE spawns the binary, it serves MCP over stdio, no port 8080 is bound. This is what `.vscode/mcp.json` and JetBrains use.

The MCP server reads JSON-RPC messages from stdin and writes responses to stdout. All diagnostic logs go to stderr in both modes.

### Quick Smoke Test

```bash
# Works in both modes — pipe to the binary directly
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1.0"}}}' | ./bin/homelab-api mcp 2>/dev/null
```

You should see a JSON response with `serverInfo.name` = `"go-github-homelab"`.

## Connect to GitHub Copilot

### VS Code

1. Build the binary: `make build`
2. The project includes `.vscode/mcp.json` pre-configured:

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

> VS Code spawns the binary with `mcp` arg — MCP-only mode, no HTTP port bound.

3. Open the project in VS Code
4. GitHub Copilot will auto-discover the MCP server configuration
5. Open Copilot Chat and ask: *"What devices are available in my home lab?"*

### JetBrains (GoLand / IntelliJ)

1. Build the binary: `make build`
2. Open Settings → Tools → AI Assistant → MCP Servers
3. Click **Add** (+) and configure:
   - **Name**: `go-github-homelab`
   - **Command**: `./bin/homelab-api` (relative to project root)
   - **Args**: `mcp`
   - **Transport**: stdio
4. Click **OK** and restart the AI Assistant
5. Open the AI Assistant chat and ask: *"What services are running in my home lab?"*

## Run Tests

```bash
# Run all tests (includes MCP server tests)
make test

# Run only MCP package tests
go test -v -race -cover ./internal/mcp/...

# Run with coverage report
go test -v -race -coverprofile=coverage.out ./internal/mcp/...
go tool cover -html=coverage.out -o coverage.html
```

## Project Structure

The MCP server is integrated into the existing single binary. Three launch modes are supported:

- `./bin/homelab-api` — default, starts HTTP API (port 8080) + MCP stdio concurrently
- `./bin/homelab-api mcp` — MCP-only mode, no HTTP port bound (used by IDEs)
- Kubernetes runs the default mode; the MCP server idles harmlessly when stdin is `/dev/null`

```
cmd/api/
└── main.go                  # Single entrypoint: launches HTTP API + MCP concurrently

internal/mcp/
├── server.go                # MCP server creation, resource/tool/prompt registration
├── server_test.go           # Server setup tests
├── resources.go             # Resource handlers (devices, services, cluster, health)
├── resources_test.go        # Resource handler tests
├── tools.go                 # Tool handlers (execute_command)
├── tools_test.go            # Tool handler tests
├── prompts.go               # Prompt handlers (device_control, service_status)
└── prompts_test.go          # Prompt handler tests

internal/homeassistant/
├── types.go                 # Command type (existing)
└── devices.go               # Exported device catalogue (NEW — shared with HTTP handlers)

internal/services/
└── provider.go              # Exported service list (NEW — shared with HTTP handlers)
```

## What You Can Do

Once connected, ask your AI assistant:

| Question | What Happens |
|----------|-------------|
| *"What devices are in my home lab?"* | Reads `homelab://devices` resource |
| *"What services are running?"* | Reads `homelab://services` resource |
| *"Show me the cluster services"* | Reads `homelab://cluster/services` resource |
| *"Is the system healthy?"* | Reads `homelab://health` resource |
| *"Turn on the living room light"* | Calls `execute_command` tool |
| *"Check the status of prometheus"* | Uses `service_status` prompt template |

## Troubleshooting

### Server won't start
- Check that the binary exists: `ls bin/homelab-api`
- Rebuild: `make build`
- Check Go version: `go version` (needs 1.25.0+)

### Copilot doesn't see the MCP server
- **VS Code**: Check `.vscode/mcp.json` exists, binary path is correct, and `args` contains `["mcp"]`
- **JetBrains**: Verify the MCP server entry in Settings → Tools → AI Assistant → MCP Servers; ensure Args is `mcp`
- Check stderr output: run `./bin/homelab-api mcp 2>mcp-debug.log` and inspect `mcp-debug.log`

### No response from server
- The MCP server communicates via stdin/stdout — it won't produce output without input
- Use the smoke test above (`./bin/homelab-api mcp`) to verify basic functionality
- Check that no other process is consuming stdin

## Dependencies

This feature adds one new dependency:

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/mark3labs/mcp-go` | MCP protocol SDK (server, transport, types) | latest stable |

