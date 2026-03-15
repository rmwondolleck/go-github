# Tasks: MCP Server Integration

**Feature**: 002-mcp-server  
**Branch**: `002-mcp-server`  
**Created**: March 14, 2026  
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md) | **Contract**: [contracts/mcp-protocol.md](./contracts/mcp-protocol.md)

---

## Overview

Integrate a Model Context Protocol (MCP) server into the existing `homelab-api` binary so
that **both the HTTP API server and the MCP stdio server start in parallel with a single
command**. There is no subcommand — running `./bin/homelab-api` launches both modes
concurrently via goroutines under a shared context. The MCP server listens on stdin/stdout
for local IDE clients (VS Code Copilot, JetBrains AI) while the HTTP API continues serving
Kubernetes traffic on port 8080, both shutting down cleanly on `SIGINT`/`SIGTERM`.

The implementation spans four user stories (P1 × 2, P2 × 1, P3 × 1):

| Phase | Scope | Stories |
|---|---|---|
| 1 | Setup — dependency, module, Makefile | — |
| 2 | Foundational — shared data providers, HomeAssistant devices provider, refactor HTTP handlers | — |
| 3 | US1 — Session handshake & capability discovery (MCP server core) | US1 |
| 4 | US2 — AI reads live home lab state (resource handlers) | US2 |
| 5 | US3 — AI controls smart home devices (tool handler) | US3 |
| 6 | US4 — AI uses pre-built prompt templates (prompt handlers) | US4 |
| 7 | Polish — binary dispatch, Makefile targets, IDE config, documentation | — |

**Total tasks**: 42  
**Parallelisable tasks**: 22 (marked `[P]`)  
**Test tasks**: included throughout (table-driven, `testify`, `≥ 80% coverage` target)

---

## Implementation Strategy

**MVP scope (ship first)**: Complete Phases 1 → 2 → 3 → 4 in sequence. This delivers a
fully working, read-only MCP server that Copilot can discover and query — independently
testable and immediately valuable.

**Incremental delivery**:
1. Phases 1–2 unlock all subsequent work (no other phase can start before Phase 2 is done).
2. Phase 3 (US1) and Phase 4 (US2) are tightly coupled; implement together for the first
   working end-to-end demo.
3. Phases 5 and 6 are independent of each other once Phase 3 is complete; they can be
   worked in parallel by separate contributors.
4. Phase 7 is pure integration and polish; it completes the feature but does not block
   manual testing of earlier phases.

---

## Dependency Graph

```
Phase 1 (T001–T002)
    └── Phase 2 (T003–T010)
            ├── Phase 3 (T011–T016)  [US1 — Discovery]
            │       └── Phase 4 (T017–T023)  [US2 — Read]
            │               ├── Phase 5 (T024–T028)  [US3 — Control]   ─┐
            │               └── Phase 6 (T029–T033)  [US4 — Prompts]   ─┤
            │                                                             │
            └── Phase 7 (T034–T042)  [Polish] ◄────────────────────────-┘
```

**Key dependency rules**:
- T003 (`go get mcp-go`) must complete before any `internal/mcp/` file is written.
- T005 (`internal/homeassistant/devices.go`) must complete before T008 (refactor HTTP handler).
- T006 (`internal/services/provider.go`) must complete before T009 (refactor services handler).
- T011 (`internal/mcp/server.go`) must complete before any resource/tool/prompt file.
- T034 (`cmd/api/main.go` dispatch) requires T011 to compile; should be the last integration step.

---

## Phase 1 — Setup

> Add the one new external dependency and verify the module compiles cleanly.

- [X] T001 Add `github.com/mark3labs/mcp-go` dependency: run `go get github.com/mark3labs/mcp-go@latest` and commit updated `go.mod` and `go.sum`
- [X] T002 [P] Verify module integrity: run `go mod tidy && go build ./...` and confirm zero errors before writing any new source files

---

## Phase 2 — Foundational: Shared Data Providers

> Extract mock data from HTTP handlers into standalone provider packages so both the HTTP
> and MCP layers can consume the same data without coupling or duplication (FR-007, FR-012).
>
> **Independent test**: `go test ./internal/homeassistant/... ./internal/services/...` passes
> with all exported functions returning expected data.

- [X] T003 Create `internal/homeassistant/devices.go`: define `CommandResult` struct (`Status`, `DeviceID`, `Action` string fields, JSON tags); define package-level sentinel errors `ErrDeviceNotFound` and `ErrDeviceNotControllable`
- [X] T004 [P] Implement `GetDevices()`, `GetDevice(id string)`, and `ExecuteCommand(deviceID string, cmd Command)` in `internal/homeassistant/devices.go`; move the `mockDevices` map here from `internal/handlers/homeassistant.go`; import `go-github/internal/models` and `time`
- [X] T005 [P] Write table-driven unit tests in `internal/homeassistant/devices_test.go` covering: `GetDevices` returns all mock devices, `GetDevice` returns correct device or `false`, `ExecuteCommand` returns success for controllable device, `ExecuteCommand` returns `ErrDeviceNotFound` for unknown ID, `ExecuteCommand` returns `ErrDeviceNotControllable` for non-controllable device
- [X] T006 Create `internal/services/provider.go` (new package `services`): implement `GetServices() []models.Service` returning the five mock services currently inlined in `internal/handlers/services.go` (homeassistant, prometheus, grafana, node-exporter, alertmanager)
- [X] T007 [P] Write unit tests in `internal/services/provider_test.go`: `GetServices` returns a non-empty slice, contains an entry with `Name == "prometheus"`, all entries have non-empty `Status`
- [X] T008 Refactor `internal/handlers/homeassistant.go`: replace the local `mockDevices` map and inline device-lookup logic in `ExecuteCommandHandler` with calls to `homeassistant.GetDevices()`, `homeassistant.GetDevice()`, and `homeassistant.ExecuteCommand()`; map `ErrDeviceNotFound` → HTTP 404, `ErrDeviceNotControllable` → HTTP 405; keep existing handler signatures and Swagger annotations unchanged
- [X] T009 Refactor `internal/handlers/services.go`: replace the inline service slice in `ListServicesHandler` with a call to `services.GetServices()`; keep existing handler signature and Swagger annotation unchanged
- [X] T010 [P] Run `go test -race ./internal/handlers/...` and confirm all existing handler tests pass after refactoring (regression gate before writing MCP code)

---

## Phase 3 — US1: Session Handshake & Capability Discovery

> **Story goal**: An MCP client connects, completes the initialize handshake, and enumerates
> all available resources, tools, and prompts in a single session.
>
> **Independent test**: Run `echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.1"}}}' | ./bin/homelab-api 2>/dev/null` and verify the JSON response on stdout contains `serverInfo.name == "go-github-homelab"` and declared `capabilities` for resources, tools, and prompts.

- [X] T011 [US1] Create `internal/mcp/server.go`: implement `NewMCPServer()` that constructs a `server.MCPServer` (from `github.com/mark3labs/mcp-go/server`) named `"go-github-homelab"` version `"1.0.0"` with `server.WithResourceCapabilities(true, true)`, `server.WithToolCapabilities(true)`, and `server.WithPromptCapabilities(true)`; export a `Run(ctx context.Context)` function that creates a `server.NewStdioServer` and calls `Listen(ctx, os.Stdin, os.Stdout)` with slog writing to `os.Stderr`; this function is intended to be launched in a goroutine alongside the HTTP server
- [X] T012 [P] [US1] Write `internal/mcp/server_test.go`: test that `NewMCPServer()` returns a non-nil server; test that calling `NewMCPServer()` twice returns independent instances; use the mcp-go in-process client to send an `initialize` request and assert the response `serverInfo.name` equals `"go-github-homelab"` and all three capability objects are present in the result
- [X] T013 [US1] Register resource stubs in `internal/mcp/server.go` (or via `RegisterAll(s *server.MCPServer)` called from `NewMCPServer`): add four `mcp.Resource` entries with URIs `homelab://devices`, `homelab://services`, `homelab://cluster/services`, `homelab://health` and placeholder handler functions that return an empty `[]mcp.ResourceContents` — these will be replaced in Phase 4
- [X] T014 [US1] Register tool stub in `internal/mcp/server.go`: add `execute_command` tool with the JSON schema from the contract (`device_id` required string, `action` required string, `parameters` optional object); placeholder handler returns an empty `*mcp.CallToolResult` — will be replaced in Phase 5
- [X] T015 [US1] Register prompt stubs in `internal/mcp/server.go`: add `device_control` (arg: `device_name`, required) and `service_status` (arg: `service_name`, required) prompts; placeholder handlers return empty `*mcp.GetPromptResult` — will be replaced in Phase 6
- [X] T016 [P] [US1] Extend `internal/mcp/server_test.go`: use in-process client to call `resources/list` and assert exactly 4 resources are returned with URIs matching the contract; call `tools/list` and assert `execute_command` is present with `inputSchema.required` containing `device_id` and `action`; call `prompts/list` and assert both `device_control` and `service_status` are present with their required arguments

---

## Phase 4 — US2: AI Reads Live Home Lab State

> **Story goal**: An AI assistant asks for current device, service, cluster, and health data;
> the MCP server returns live data from the shared providers.
>
> **Independent test**: Using an in-process MCP client (or piped JSON), call `resources/read`
> for each of the four resource URIs and verify each returns a non-empty JSON payload;
> call `resources/read` with `homelab://unknown` and verify the error response.

- [X] T017 [US2] Create `internal/mcp/resources.go`: implement `DevicesResourceHandler(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)` — call `homeassistant.GetDevices()`, marshal the map values to JSON, return a single `mcp.TextResourceContents` with URI `homelab://devices` and `MIMEType "application/json"`
- [X] T018 [P] [US2] Implement `ServicesResourceHandler` in `internal/mcp/resources.go`: call `services.GetServices()`, marshal to JSON, return `mcp.TextResourceContents` for URI `homelab://services`
- [X] T019 [P] [US2] Implement `ClusterServicesResourceHandler` in `internal/mcp/resources.go`: instantiate `cluster.NewService()`, call `ListServices(ctx)`, marshal result to JSON, return `mcp.TextResourceContents` for URI `homelab://cluster/services`; wrap any error from `ListServices` as a resource-read error
- [X] T020 [P] [US2] Implement `HealthResourceHandler` in `internal/mcp/resources.go`: instantiate `health.NewChecker()`, call `Check(ctx)`, marshal `models.HealthStatus` to JSON, return `mcp.TextResourceContents` for URI `homelab://health`
- [X] T021 [US2] Wire real resource handlers into `internal/mcp/server.go`: replace the four placeholder handlers registered in T013 with the real functions from `resources.go`
- [X] T022 [P] [US2] Write `internal/mcp/resources_test.go` with table-driven tests covering all four handlers: each test calls the handler directly with a valid `mcp.ReadResourceRequest` and asserts (a) no error, (b) exactly one `ResourceContents` item returned, (c) the `MIMEType` is `"application/json"`, (d) the `text` field is valid JSON and non-empty; add a test that confirms the health handler returns a payload containing `"status"` key
- [X] T023 [P] [US2] Add error-path tests in `internal/mcp/resources_test.go`: verify that a handler receiving a context cancelled mid-call returns an error (not a panic); verify that marshalling a nil map produces valid JSON `"null"` or `"[]"` and does not error

---

## Phase 5 — US3: AI Controls Smart Home Devices

> **Story goal**: An AI assistant executes a device command via the `execute_command` tool;
> the server validates inputs, delegates to `homeassistant.ExecuteCommand`, and returns a
> structured success or error result.
>
> **Independent test**: Using an in-process client, call `tools/call` with `execute_command`,
> `device_id: "device-001"`, `action: "turn_on"`, and verify `content[0].text` parses to
> `{status:"success",...}`; call with `device_id: "readonly-sensor-001"` and verify `isError: true`.

- [X] T024 [US3] Create `internal/mcp/tools.go`: implement `ExecuteCommandHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)` — extract `device_id` and `action` string args from `req.Params.Arguments`; return `isError: true` content if either is missing or empty; call `homeassistant.ExecuteCommand(deviceID, homeassistant.Command{Action: action, Parameters: params})`
- [X] T025 [US3] Map `ExecuteCommand` errors to MCP tool results in `internal/mcp/tools.go`: `ErrDeviceNotFound` → `isError: true`, text `"device not found: <id>"`; `ErrDeviceNotControllable` → `isError: true`, text `"device is not controllable: <id>"`; success → JSON-marshal `CommandResult` into `content[0].text`
- [X] T026 [US3] Wire real tool handler into `internal/mcp/server.go`: replace the placeholder `execute_command` handler registered in T014 with `tools.ExecuteCommandHandler`
- [X] T027 [P] [US3] Write `internal/mcp/tools_test.go` with table-driven tests: success path (`device-001`, `turn_on`) → `isError` is false, text contains `"success"`; not-found path (`bad-id`, `turn_on`) → `isError` true, text contains `"device not found"`; non-controllable path (`readonly-sensor-001`, `turn_on`) → `isError` true, text contains `"not controllable"`; missing `device_id` → `isError` true, text contains `"device_id is required"`; missing `action` → `isError` true, text contains `"action is required"`
- [X] T028 [P] [US3] Add race-condition test in `internal/mcp/tools_test.go`: launch 10 goroutines concurrently each calling `ExecuteCommandHandler` with valid args; assert all return without error and with non-nil results (validates FR-011 — concurrent requests handled safely)

---

## Phase 6 — US4: AI Uses Pre-Built Prompt Templates

> **Story goal**: An AI assistant requests `device_control` or `service_status` prompt
> templates with arguments; the server returns fully rendered, ready-to-use prompt messages.
>
> **Independent test**: Using an in-process client, call `prompts/get` with
> `name: "device_control"`, `arguments: {device_name: "Living Room Light"}` and verify the
> returned message text contains `"Living Room Light"` and `"homelab://devices"`.

- [X] T029 [US4] Create `internal/mcp/prompts.go`: implement `DeviceControlPromptHandler(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error)` — extract `device_name` from `req.Params.Arguments`; return `isError`-style error if `device_name` is empty; return a `mcp.GetPromptResult` with one `mcp.PromptMessage` (role `user`) whose text content matches the template from the data model (references `homelab://devices` and `execute_command`)
- [X] T030 [US4] Implement `ServiceStatusPromptHandler` in `internal/mcp/prompts.go`: extract `service_name`; return error if empty; return rendered prompt message referencing `homelab://services` and `homelab://cluster/services`
- [X] T031 [US4] Wire real prompt handlers into `internal/mcp/server.go`: replace the placeholder `device_control` and `service_status` handlers registered in T015 with the real functions from `prompts.go`
- [X] T032 [P] [US4] Write `internal/mcp/prompts_test.go` with table-driven tests: `device_control` with `device_name="Living Room Light"` → result non-nil, message text contains `"Living Room Light"` and `"homelab://devices"`; `service_status` with `service_name="prometheus"` → message text contains `"prometheus"` and `"homelab://services"`; `device_control` with empty `device_name` → returns non-nil error; `service_status` with empty `service_name` → returns non-nil error
- [X] T033 [P] [US4] Add unknown-prompt test via in-process client in `internal/mcp/server_test.go`: call `prompts/get` with `name: "nonexistent_prompt"` and assert the response contains a JSON-RPC error (code `-32602` or equivalent SDK error), not a panic

---

## Phase 7 — Polish & Cross-Cutting Concerns

> Wire the binary dispatch, add Makefile targets, create IDE config, and ensure the full
> feature is observable, documented, and regression-safe.

- [X] T034 Modify `cmd/api/main.go`: remove any `os.Args` subcommand dispatch; instead, under the existing shared `ctx` (cancelled on `SIGINT`/`SIGTERM`), launch **both** the HTTP server and the MCP server as concurrent goroutines — `go srv.Run(port)` for HTTP (existing) and `go mcp.NewMCPServer().Run(ctx)` for MCP (new); collect errors from both via an `errgroup` or `select` on an error channel; import `go-github/internal/mcp`; add a `slog.Info("mcp server started", "transport", "stdio")` log line alongside the existing HTTP startup log so operators can confirm both modes are active
- [X] T035 [P] Write `internal/mcp/server_test.go` integration test (or new `cmd/api/main_test.go`): run the compiled binary (no subcommand arg), pipe an `initialize` JSON-RPC message to its stdin, read stdout, and assert the response is valid JSON containing `serverInfo.name == "go-github-homelab"` (end-to-end smoke test per SC-001 / SC-006); verify the HTTP server is also reachable on its port within the same process
- [X] T036 [P] Update `Makefile`: rename or alias `mcp-build` to the standard `build` target (same binary serves both modes); ensure `make run` starts `./bin/homelab-api` with both HTTP and MCP active simultaneously; update the `help` target to document that a single binary serves both HTTP and MCP; update `.PHONY` list; confirm `make build && make run` starts both modes with log lines visible for both
- [X] T037 [P] Create `.vscode/mcp.json` with the VS Code Copilot MCP server configuration: `servers["go-github-homelab"]` with `type: "stdio"`, `command: "${workspaceFolder}/bin/homelab-api"`, and **no `args`** (MCP starts automatically alongside HTTP when the binary runs)
- [X] T038 [P] Update `README.md`: add an "MCP Server" section after the existing API section documenting (a) how to build (`make build`), (b) how to run — a single `make run` or `./bin/homelab-api` starts **both** HTTP API and MCP simultaneously, (c) the quick smoke-test pipe command (no subcommand needed), (d) VS Code connection steps referencing `.vscode/mcp.json`, (e) JetBrains manual configuration steps (Name, Command, no Args, Transport: stdio), (f) available resources and tools summary table
- [X] T039 [P] Run the full test suite: `go test -v -race -coverprofile=coverage.out ./...` and confirm (a) all existing tests pass (SC-007), (b) `internal/mcp/` package coverage is ≥ 80% (SC-005), (c) zero race conditions detected
- [X] T040 [P] Run `go vet ./...` and `golangci-lint run` (if available) against all modified and new files; resolve any reported issues in `internal/homeassistant/`, `internal/services/`, `internal/mcp/`, and `cmd/api/`
- [X] T041 [P] Verify Docker build still works: `make docker` completes without error; run `docker run --rm homelab-api:latest` and confirm the binary prints a valid JSON-RPC initialize response when piped a test message on stdin — both HTTP and MCP start from the same container entrypoint (SC-006 extended to container build)
- [X] T042 [P] Create `specs/002-mcp-server/tasks.md` final status pass: mark completed tasks, update the `CURRENT_STATUS.md` in the repo root to reflect feature 002 is complete, and confirm the branch is ready for review

---

## Parallel Execution Examples

### Sprint 1 — Phases 1 & 2 (sequential foundation)

```
[single track]  T001 → T002 → T003 → T004 → T006
                                ↓         ↓
                               T005      T007
                                ↓         ↓
                T008 (needs T004+T005) → T009 (needs T006) → T010
```

### Sprint 2 — Phase 3 (US1 core)

```
T011 (server.go skeleton)
    ↓
T013 → T014 → T015  (stubs in order, same file)
    ↓
T012 + T016  [P, parallel after T011/T013/T014/T015]
```

### Sprint 3 — Phases 4 & 5 & 6 (parallel tracks)

```
Track A (US2 reads):   T017 → T018[P] → T019[P] → T020[P] → T021 → T022[P] → T023[P]
Track B (US3 tools):   T024 → T025 → T026 → T027[P] → T028[P]
Track C (US4 prompts): T029 → T030 → T031 → T032[P] → T033[P]
```
*Tracks B and C can start after T011 (server skeleton) is complete.*

### Sprint 4 — Phase 7 (polish, mostly parallel)

```
T034 (main.go dispatch — blocks T035)
T035[P] + T036[P] + T037[P] + T038[P]  (parallel once T034 compiles)
T039 → T040 → T041 → T042  (final validation sequence)
```

---

## Acceptance Criteria Reference

| Task Range | User Story | Acceptance Scenario Covered |
|---|---|---|
| T011–T016 | US1 | AS-1.1 (handshake), AS-1.2 (list resources), AS-1.3 (list tools), AS-1.4 (list prompts), AS-1.5 (pre-handshake rejection handled by SDK) |
| T017–T023 | US2 | AS-2.1 (devices), AS-2.2 (services), AS-2.3 (cluster), AS-2.4 (health), AS-2.5 (unknown resource → error) |
| T024–T028 | US3 | AS-3.1 (valid command), AS-3.2 (non-controllable), AS-3.3 (missing inputs), AS-3.4 (unknown action → ErrDeviceNotFound) |
| T029–T033 | US4 | AS-4.1 (device_control), AS-4.2 (service_status), AS-4.3 (unknown prompt → SDK error) |
| T034–T042 | All | FR-008 (standalone process), FR-009 (Makefile), SC-001 (< 2s), SC-005 (≥ 80% coverage), SC-006 (clean build), SC-007 (existing tests pass), SC-008 (IDE config) |

---

## File Change Summary

| File | Action | Phase |
|---|---|---|
| `go.mod`, `go.sum` | Modified — add `github.com/mark3labs/mcp-go` | 1 |
| `internal/homeassistant/devices.go` | **New** — `GetDevices`, `GetDevice`, `ExecuteCommand`, `CommandResult`, sentinel errors | 2 |
| `internal/homeassistant/devices_test.go` | **New** — unit tests for provider | 2 |
| `internal/services/provider.go` | **New** — `GetServices()` | 2 |
| `internal/services/provider_test.go` | **New** — unit tests for provider | 2 |
| `internal/handlers/homeassistant.go` | Modified — delegate to `homeassistant` package | 2 |
| `internal/handlers/services.go` | Modified — delegate to `services` package | 2 |
| `internal/mcp/server.go` | **New** — `NewMCPServer()`, `Run()`, registration | 3 |
| `internal/mcp/server_test.go` | **New** — server creation, capability, list tests | 3 + 6 |
| `internal/mcp/resources.go` | **New** — four resource handlers | 4 |
| `internal/mcp/resources_test.go` | **New** — resource handler tests | 4 |
| `internal/mcp/tools.go` | **New** — `execute_command` handler | 5 |
| `internal/mcp/tools_test.go` | **New** — tool handler tests incl. race test | 5 |
| `internal/mcp/prompts.go` | **New** — `device_control`, `service_status` handlers | 6 |
| `internal/mcp/prompts_test.go` | **New** — prompt handler tests | 6 |
| `cmd/api/main.go` | Modified — launch HTTP + MCP as parallel goroutines under shared context | 7 |
| `Makefile` | Modified — update `build`/`run` to document dual-mode binary; add `help` entries | 7 |
| `.vscode/mcp.json` | **New** — VS Code Copilot MCP configuration | 7 |
| `README.md` | Modified — add MCP Server section | 7 |
