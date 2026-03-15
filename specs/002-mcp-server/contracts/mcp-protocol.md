# MCP Server Contract: JSON-RPC 2.0 over stdio

**Feature**: 002-mcp-server  
**Protocol**: Model Context Protocol (MCP) via JSON-RPC 2.0  
**Transport**: stdin/stdout (line-delimited JSON)  
**Binary**: `bin/homelab-api` — three modes:
- `./bin/homelab-api` — HTTP API (port 8080) + MCP stdio concurrently (default / k8s)
- `./bin/homelab-api mcp` — MCP stdio only, no HTTP port bound (IDE / local dev)

---

## Transport

- **Input**: JSON-RPC 2.0 messages read from `stdin`, one per line
- **Output**: JSON-RPC 2.0 responses written to `stdout`, one per line
- **Diagnostics**: All logging goes to `stderr` (never stdout)
- **Lifecycle**: Server runs until stdin is closed or a fatal error occurs

## Server Identity

```json
{
  "name": "go-github-homelab",
  "version": "1.0.0"
}
```

## Capabilities Declared

```json
{
  "resources": { "subscribe": true, "listChanged": true },
  "tools": { "listChanged": true },
  "prompts": { "listChanged": true }
}
```

---

## Initialize Handshake

### Request

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "copilot",
      "version": "1.0.0"
    }
  }
}
```

### Response

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "resources": { "subscribe": true, "listChanged": true },
      "tools": { "listChanged": true },
      "prompts": { "listChanged": true }
    },
    "serverInfo": {
      "name": "go-github-homelab",
      "version": "1.0.0"
    }
  }
}
```

---

## Resources

### List Resources

**Method**: `resources/list`

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "resources": [
      {
        "uri": "homelab://devices",
        "name": "devices",
        "description": "Home automation device catalogue with current state",
        "mimeType": "application/json"
      },
      {
        "uri": "homelab://services",
        "name": "services",
        "description": "Available home lab services",
        "mimeType": "application/json"
      },
      {
        "uri": "homelab://cluster/services",
        "name": "cluster_services",
        "description": "Kubernetes cluster services and endpoints",
        "mimeType": "application/json"
      },
      {
        "uri": "homelab://health",
        "name": "health",
        "description": "System health status and uptime",
        "mimeType": "application/json"
      }
    ]
  }
}
```

### Read Resource — homelab://devices

**Method**: `resources/read`  
**Request params**: `{ "uri": "homelab://devices" }`

**Response**: Array of device objects with `id`, `name`, `type`, `state`, `attributes`, `controllable`.

### Read Resource — homelab://health

**Method**: `resources/read`  
**Request params**: `{ "uri": "homelab://health" }`

**Response**: Health object with `status`, `uptime`, `components`.

---

## Tools

### List Tools

**Method**: `tools/list`

**Response**: Single tool `execute_command` with input schema requiring `device_id` (string) and `action` (string), optional `parameters` (object).

### Call Tool — execute_command

**Method**: `tools/call`  
**Request params**:
```json
{
  "name": "execute_command",
  "arguments": {
    "device_id": "light.living_room",
    "action": "turn_on",
    "parameters": {}
  }
}
```

**Success response**: `content[0].text` contains JSON-marshalled `CommandResult`.

**Error responses**:
- Device not found: `isError: true`, text `"device not found: <id>"`
- Not controllable: `isError: true`, text `"device is not controllable: <id>"`

---

## Prompts

### List Prompts

**Method**: `prompts/list`

**Response**: Two prompts — `device_control` (arg: `device_name`) and `service_status` (arg: `service_name`).

### Get Prompt — device_control

**Request params**: `{ "name": "device_control", "arguments": { "device_name": "Living Room Light" } }`

**Response**: Rendered prompt message guiding AI to read `homelab://devices` then call `execute_command`.

### Get Prompt — service_status

**Request params**: `{ "name": "service_status", "arguments": { "service_name": "prometheus" } }`

**Response**: Rendered prompt message guiding AI to check `homelab://services` and `homelab://cluster/services` for the named service.

### Get Prompt — Unknown Prompt

**Error Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "error": {
    "code": -32602,
    "message": "prompt not found: unknown_prompt"
  }
}
```

---

## IDE Configuration

### VS Code (`.vscode/mcp.json`)

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

> Uses `mcp` arg — MCP-only mode. The IDE gets a clean stdio pipe with no HTTP server binding port 8080.

### JetBrains (GoLand / IntelliJ)

GitHub Copilot in JetBrains reads MCP configuration from a file — not the Settings UI.

**Config file**:
- Windows: `%LOCALAPPDATA%\github-copilot\intellij\mcp.json`
- macOS: `~/Library/Application Support/github-copilot/intellij/mcp.json`
- Linux: `~/.config/github-copilot/intellij/mcp.json`

```json
{
  "servers": {
    "go-github-homelab": {
      "type": "stdio",
      "command": "C:/Users/<you>/GolandProjects/go-github/bin/homelab-api",
      "args": ["mcp"]
    }
  }
}
```

> Use the **absolute path** — JetBrains does not resolve `${workspaceFolder}`. Restart the IDE after saving.
