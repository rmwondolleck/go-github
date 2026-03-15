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
        "description": "Available home lab services and their status",
        "mimeType": "application/json"
      },
      {
        "uri": "homelab://cluster/services",
        "name": "cluster_services",
        "description": "Kubernetes cluster services with endpoints",
        "mimeType": "application/json"
      },
      {
        "uri": "homelab://health",
        "name": "health",
        "description": "System health status with uptime and component states",
        "mimeType": "application/json"
      }
    ]
  }
}
```

### Read Resource — Devices

**Method**: `resources/read`

**Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/read",
  "params": { "uri": "homelab://devices" }
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "contents": [
      {
        "uri": "homelab://devices",
        "mimeType": "application/json",
        "text": "[{\"id\":\"device-001\",\"name\":\"Living Room Light\",\"type\":\"light\",\"state\":\"off\",\"attributes\":{\"brightness\":0},\"controllable\":true},{\"id\":\"readonly-sensor-001\",\"name\":\"Temperature Sensor\",\"type\":\"sensor\",\"state\":\"72\",\"attributes\":{\"unit\":\"°F\"},\"controllable\":false}]"
      }
    ]
  }
}
```

### Read Resource — Services

**Request params**: `{ "uri": "homelab://services" }`

**Response `text` field** (JSON array of `models.Service`): Current service list from shared provider.

### Read Resource — Cluster Services

**Request params**: `{ "uri": "homelab://cluster/services" }`

**Response `text` field** (JSON array of `cluster.ServiceInfo`): Current cluster service list.

### Read Resource — Health

**Request params**: `{ "uri": "homelab://health" }`

**Response `text` field** (JSON object `models.HealthStatus`): Current health status including uptime and components.

### Read Resource — Unknown URI

**Request params**: `{ "uri": "homelab://unknown" }`

**Error Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "error": {
    "code": -32602,
    "message": "resource not found: homelab://unknown"
  }
}
```

---

## Tools

### List Tools

**Method**: `tools/list`

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "result": {
    "tools": [
      {
        "name": "execute_command",
        "description": "Execute a control command on a home automation device. Requires the device ID, an action name, and optional parameters.",
        "inputSchema": {
          "type": "object",
          "properties": {
            "device_id": {
              "type": "string",
              "description": "The ID of the device to control"
            },
            "action": {
              "type": "string",
              "description": "The action to perform (e.g., turn_on, turn_off, set)"
            },
            "parameters": {
              "type": "object",
              "description": "Optional parameters for the action (e.g., brightness level)"
            }
          },
          "required": ["device_id", "action"]
        }
      }
    ]
  }
}
```

### Call Tool — Success

**Method**: `tools/call`

**Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "device_id": "device-001",
      "action": "turn_on",
      "parameters": {}
    }
  }
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"status\":\"success\",\"device_id\":\"device-001\",\"action\":\"turn_on\"}"
      }
    ]
  }
}
```

### Call Tool — Device Not Found

**Response** (isError=true):
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "content": [
      { "type": "text", "text": "device not found: invalid-id" }
    ],
    "isError": true
  }
}
```

### Call Tool — Device Not Controllable

**Response** (isError=true):
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "content": [
      { "type": "text", "text": "device is not controllable: readonly-sensor-001" }
    ],
    "isError": true
  }
}
```

### Call Tool — Missing Required Arguments

**Response** (isError=true):
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "content": [
      { "type": "text", "text": "device_id is required" }
    ],
    "isError": true
  }
}
```

---

## Prompts

### List Prompts

**Method**: `prompts/list`

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "result": {
    "prompts": [
      {
        "name": "device_control",
        "description": "Guide an AI assistant to control a named home automation device",
        "arguments": [
          { "name": "device_name", "description": "Name of the device to control", "required": true }
        ]
      },
      {
        "name": "service_status",
        "description": "Guide an AI assistant to report on a named service status",
        "arguments": [
          { "name": "service_name", "description": "Name of the service to check", "required": true }
        ]
      }
    ]
  }
}
```

### Get Prompt — device_control

**Method**: `prompts/get`

**Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "prompts/get",
  "params": {
    "name": "device_control",
    "arguments": { "device_name": "Living Room Light" }
  }
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "result": {
    "description": "Guide an AI assistant to control a named home automation device",
    "messages": [
      {
        "role": "user",
        "content": {
          "type": "text",
          "text": "You are a home automation assistant. The user wants to control the device named \"Living Room Light\".\n\n1. First, read the device list from homelab://devices to find the device and its current state.\n2. Determine the device ID and whether it is controllable.\n3. If controllable, use the execute_command tool with the appropriate device_id, action, and parameters.\n4. Report the result to the user.\n\nIf the device is not found or not controllable, explain why the action cannot be performed."
        }
      }
    ]
  }
}
```

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
