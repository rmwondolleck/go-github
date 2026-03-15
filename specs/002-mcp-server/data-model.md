# Data Model: MCP Server Integration

**Feature**: 002-mcp-server  
**Date**: March 14, 2026  
**Status**: Complete

---

## Overview

This document defines the data entities for the MCP Server integration. The MCP layer does **not** introduce new persistent entities — it wraps existing domain models (`models.Device`, `models.Service`, `cluster.ServiceInfo`, `models.HealthStatus`) and exposes them via MCP protocol types (`mcp.Resource`, `mcp.Tool`, `mcp.Prompt`).

The key architectural change is **extracting shared data providers** from HTTP handlers so both the HTTP and MCP layers consume the same data sources.

---

## Existing Entities (Unchanged)

### Device (`internal/models/device.go`)

```go
type Device struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    Type         string                 `json:"type"`
    State        string                 `json:"state"`
    Attributes   map[string]interface{} `json:"attributes"`
    LastUpdated  time.Time              `json:"last_updated"`
    Controllable bool                   `json:"controllable"`
}
```

**Validation Rules**:
- `ID` must be non-empty
- `Name` must be non-empty
- `Type` must be one of: `light`, `sensor`, `switch`, `thermostat` (extensible)
- `Controllable` determines whether `execute_command` tool accepts commands for this device

### Command (`internal/homeassistant/types.go`)

```go
type Command struct {
    Action     string                 `json:"action"`
    Parameters map[string]interface{} `json:"parameters"`
}
```

**Validation Rules**:
- `Action` must be non-empty (trimmed)
- `Parameters` must not be nil

### Service (`internal/models/service.go`)

```go
type Service struct {
    Name     string `json:"name"`
    Type     string `json:"type"`
    Status   string `json:"status"`
    Endpoint string `json:"endpoint"`
}
```

### ClusterServiceInfo (`internal/cluster/types.go`)

```go
type ServiceInfo struct {
    Name      string   `json:"name"`
    Namespace string   `json:"namespace"`
    Status    string   `json:"status"`
    Endpoints []string `json:"endpoints"`
}
```

### HealthStatus (`internal/models/health.go`)

```go
type HealthStatus struct {
    Status     string            `json:"status"`
    Uptime     string            `json:"uptime"`
    Components map[string]string `json:"components"`
}
```

---

## New Shared Data Providers (Extracted from Handlers)

### DeviceProvider (`internal/homeassistant/devices.go`) — NEW FILE

Extracts `mockDevices` map from `internal/handlers/homeassistant.go`.

```go
// Functions exported for shared access by HTTP handlers and MCP handlers
func GetDevices() map[string]*models.Device    // Returns full device catalogue
func GetDevice(id string) (*models.Device, bool) // Returns single device by ID
func ExecuteCommand(deviceID string, cmd Command) (CommandResult, error) // Execute and return result
```

```go
// CommandResult represents the outcome of a device command execution
type CommandResult struct {
    Status   string `json:"status"`   // "success" or "error"
    DeviceID string `json:"device_id"`
    Action   string `json:"action"`
}
```

**State Transitions** (Device.State):
- `off` → `on` (via `turn_on` action)
- `on` → `off` (via `turn_off` action)
- `{value}` → `{new_value}` (via `set` action on controllable devices)
- Note: Current mock implementation does not persist state changes; returns success without mutation

**Error Conditions**:
- Device not found → `ErrDeviceNotFound`
- Device not controllable → `ErrDeviceNotControllable`
- Invalid command → validation error from `Command.Validate()`

### ServiceProvider (`internal/services/provider.go`) — NEW FILE

Extracts inline service list from `internal/handlers/services.go`.

```go
func GetServices() []models.Service  // Returns full service catalogue
```

---

## MCP Protocol Entities (SDK Types — Not Custom)

These are **SDK-provided types** from `github.com/mark3labs/mcp-go/mcp`. They are not custom entities but are documented here to show how domain models map to the MCP protocol.

### MCP Resource Registration

| Resource URI | Name | Description | MIME Type | Data Source |
|---|---|---|---|---|
| `homelab://devices` | `devices` | Home automation device catalogue | `application/json` | `homeassistant.GetDevices()` |
| `homelab://services` | `services` | Available home lab services | `application/json` | `services.GetServices()` |
| `homelab://cluster/services` | `cluster_services` | Kubernetes cluster services | `application/json` | `cluster.Service.ListServices()` |
| `homelab://health` | `health` | System health status | `application/json` | `health.Checker.Check()` |

### MCP Tool Registration

| Tool Name | Description | Input Schema |
|---|---|---|
| `execute_command` | Execute a command on a home automation device | `device_id` (string, required), `action` (string, required), `parameters` (object, optional) |

**Input Schema (JSON Schema)**:
```json
{
  "type": "object",
  "properties": {
    "device_id": { "type": "string", "description": "The ID of the device to control" },
    "action": { "type": "string", "description": "The action to perform (e.g., turn_on, turn_off, set)" },
    "parameters": { "type": "object", "description": "Optional parameters for the action" }
  },
  "required": ["device_id", "action"]
}
```

### MCP Prompt Registration

| Prompt Name | Description | Arguments |
|---|---|---|
| `device_control` | Guides an AI to control a named device | `device_name` (string, required) |
| `service_status` | Guides an AI to report on a named service | `service_name` (string, required) |

**Prompt Template Rendering**:

`device_control` with `device_name="Living Room Light"`:
```
You are a home automation assistant. The user wants to control the device named "Living Room Light".

1. First, read the device list from homelab://devices to find the device and its current state.
2. Determine the device ID and whether it is controllable.
3. If controllable, use the execute_command tool with the appropriate device_id, action, and parameters.
4. Report the result to the user.

If the device is not found or not controllable, explain why the action cannot be performed.
```

`service_status` with `service_name="prometheus"`:
```
You are a home lab monitoring assistant. The user wants to know the status of the service named "prometheus".

1. Read the service list from homelab://services to find the service.
2. If not found in services, check homelab://cluster/services.
3. Report the service name, type, status, and endpoint.
4. If the service is not found in either source, inform the user.
```

---

## Entity Relationship Diagram

```
┌──────────────────────────────────────────────────────────┐
│   cmd/api/main.go  (single binary, single command)       │
│                                                           │
│   errgroup.Go ──────────────────────────────────────┐    │
│       │                                             │    │
│       ▼ goroutine 1                  goroutine 2 ◀──┘    │
│  srv.Run(port)               internalmcp.Run(ctx)        │
│  (HTTP API :8080)            (MCP stdio)                 │
│       │                             │                    │
│       │ shared context (errgroup)   │                    │
│       └──────────────┬──────────────┘                    │
└──────────────────────┼───────────────────────────────────┘
                       │ both use
                       ▼
         ┌─────────────────────────┐
         │  Shared Data Providers  │
         │                         │
         │ homeassistant/devices.go │
         │ services/provider.go     │
         │ cluster/service.go       │
         │ health/checker.go        │
         └──────────┬──────────────┘
                    │ consumed by
          ┌─────────┴──────────┐
          ▼                    ▼
┌─────────────────┐  ┌──────────────────────┐
│ internal/       │  │   internal/mcp/      │
│  handlers/      │  │   server.go          │
│  (HTTP layer)   │  │   resources.go       │
│                 │  │   tools.go           │
│                 │  │   prompts.go         │
└─────────────────┘  └──────────────────────┘
```

Both the HTTP handlers and MCP handlers consume the **same shared data providers**. No data duplication.
