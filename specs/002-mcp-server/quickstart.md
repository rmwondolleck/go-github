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

```bash
# A single command starts BOTH the HTTP API server (port 8080) AND the MCP stdio server
./bin/homelab-api

# Or via Make
make run
```

Both the HTTP API and MCP server start concurrently under a shared context. Both shut down
cleanly on `SIGINT`/`SIGTERM`. The MCP server reads JSON-RPC messages from stdin and writes
responses to stdout. Diagnostic logs from both modes go to stderr.

### Quick Smoke Test

Send an initialize request by piping JSON to the binary (no subcommand needed):

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1.0"}}}' | ./bin/homelab-api 2>/dev/null
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
      "command": "${workspaceFolder}/bin/homelab-api"
    }
  }
}
```

> No `args` needed — the MCP server starts automatically alongside the HTTP API when the binary runs.

3. Open the project in VS Code
4. GitHub Copilot will auto-discover the MCP server configuration
5. Open Copilot Chat and ask: *"What devices are available in my home lab?"*

### JetBrains (GoLand / IntelliJ)

1. Build the binary: `make build`
2. Open Settings → Tools → AI Assistant → MCP Servers
3. Click **Add** (+) and configure:
   - **Name**: `go-github-homelab`
   - **Command**: `./bin/homelab-api` (relative to project root)
   - **Args**: *(leave empty)*
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

The MCP server is integrated into the existing single binary. Running `./bin/homelab-api`
starts **both** the HTTP API server (port 8080) and the MCP stdio server concurrently — no
subcommand or separate process required.

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
- **VS Code**: Check `.vscode/mcp.json` exists and the binary path is correct; ensure there are **no `args`** (MCP starts automatically)
- **JetBrains**: Verify the MCP server entry in Settings → Tools → AI Assistant → MCP Servers; ensure Args is empty
- Check stderr output: run `./bin/homelab-api 2>mcp-debug.log` and inspect `mcp-debug.log`

### No response from server
- The MCP server communicates via stdin/stdout — it won't produce output without input
- Use the smoke test above to verify basic functionality
- Check that no other process is consuming stdin

## Dependencies

This feature adds one new dependency:

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/mark3labs/mcp-go` | MCP protocol SDK (server, transport, types) | latest stable |

