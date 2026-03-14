# Feature Specification: MCP Server Integration

**Feature Branch**: `002-mcp-server`  
**Created**: March 14, 2026  
**Status**: Draft

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - AI Assistant Discovers Home Lab Capabilities (Priority: P1)

A developer connects GitHub Copilot (via VS Code or JetBrains) or any MCP-compatible client to the home lab service. The AI assistant automatically discovers what data and actions are available — devices, services, cluster status, health — without the developer writing any custom integration code.

**Why this priority**: Discovery is the entry point for all other functionality. Without it, the AI assistant cannot read data or take actions. It is the foundation on which all other stories depend.

**Independent Test**: A developer launches the MCP server binary, connects a standard MCP client, and verifies the client lists all available resources and tools. Delivers value immediately: a developer can see exactly what an AI assistant can do with their home lab.

**Acceptance Scenarios**:

1. **Given** the MCP server is running, **When** an AI client connects and initiates a session, **Then** the server confirms the connection and declares its capabilities
2. **Given** a connected session, **When** the AI client requests a list of available data sources, **Then** the server returns at least: home automation devices, available services, cluster services, and system health — each with a name and description
3. **Given** a connected session, **When** the AI client requests a list of available actions, **Then** the server returns at least the ability to execute device commands, with a description of required inputs
4. **Given** a connected session, **When** the AI client requests a list of available prompt templates, **Then** the server returns pre-defined templates the AI can use to guide interactions
5. **Given** a client that has not yet completed the connection handshake, **When** it attempts to request data or execute actions, **Then** the server rejects the request with a clear error indicating the session is not ready

---

### User Story 2 - AI Assistant Reads Live Home Lab State (Priority: P1)

A user asks their AI assistant a natural language question about their home lab — "what lights are on?", "is prometheus running?", "what's the health of the system?" — and the AI fetches live data directly from the home lab service to answer accurately.

**Why this priority**: Read access to live data is the core value of MCP integration. Without it, the AI assistant cannot answer real questions about the home lab. P1 alongside discovery because together they represent a fully useful read-only AI assistant.

**Independent Test**: A user asks the AI assistant to describe the current state of home automation devices and cluster services. The AI retrieves and accurately reports live data without any manual API calls. Delivers standalone value: users get AI-powered read access to their home lab.

**Acceptance Scenarios**:

1. **Given** a connected session, **When** the AI requests home automation device data, **Then** the server returns current device information sourced live from the home automation system
2. **Given** a connected session, **When** the AI requests available services, **Then** the server returns the current list of home lab services
3. **Given** a connected session, **When** the AI requests cluster service information with an optional name filter, **Then** the server returns matching services, or all services if no filter is provided
4. **Given** a connected session, **When** the AI requests system health, **Then** the server returns current health status including uptime and component states
5. **Given** a connected session, **When** the AI requests a data source that does not exist, **Then** the server returns a clear error identifying the unknown resource

---

### User Story 3 - AI Assistant Controls Smart Home Devices (Priority: P2)

A user instructs their AI assistant to take an action — "turn on the living room lights", "set the thermostat to 72", "toggle the front door lock" — and the AI executes the command against the home automation system on the user's behalf.

**Why this priority**: Action execution is what makes the integration agentic rather than purely informational. P2 because a read-only AI assistant already delivers significant value; device control adds the agentic layer on top.

**Independent Test**: A user asks the AI to turn on a controllable device. The AI executes the command and reports success or failure. The device state changes as expected.

**Acceptance Scenarios**:

1. **Given** a connected session, **When** the AI executes a device command with a valid device ID, action, and parameters, **Then** the server executes the command and returns a success result
2. **Given** a connected session, **When** the AI attempts to execute a command on a read-only device (e.g., a sensor), **Then** the server returns an error indicating the device cannot be controlled
3. **Given** a connected session, **When** the AI executes a command with missing required inputs (device ID or action), **Then** the server returns an error describing what is missing
4. **Given** a connected session, **When** the AI attempts to execute an action that does not exist, **Then** the server returns a clear error identifying the unknown action

---

### User Story 4 - AI Assistant Uses Pre-Built Prompt Templates (Priority: P3)

A developer building an AI-powered home lab dashboard uses pre-defined prompt templates from the MCP server to standardise how the AI interacts with home lab data — removing the need for each developer to craft their own prompts from scratch.

**Why this priority**: Prompt templates improve consistency and developer experience but are not required for core functionality. The system is fully valuable without them.

**Independent Test**: A developer retrieves a prompt template by name with specific arguments (e.g., a device name), and receives a fully rendered prompt ready for use with an AI model.

**Minimum required prompt templates**:
- `device_control` — guides an AI to control a named device; required argument: `device_name`
- `service_status` — guides an AI to report on a named service; required argument: `service_name`

**Acceptance Scenarios**:

1. **Given** a connected session, **When** the AI requests the device control prompt template with a device name argument, **Then** the server returns a rendered, ready-to-use prompt
2. **Given** a connected session, **When** the AI requests the service status prompt template with a service name argument, **Then** the server returns a rendered, ready-to-use prompt
3. **Given** a connected session, **When** the AI requests a prompt template that does not exist, **Then** the server returns a clear error identifying the unknown template

---

### Edge Cases

- What happens when the MCP server receives a malformed or unparseable message? → Returns a standard error response; does not crash
- What happens when the underlying home lab data source is unavailable during a read? → Returns a clear error to the client; the server remains running
- What happens when two requests arrive from the same client in rapid succession? → Both are handled correctly and in order without corruption. Note: each server process serves exactly one client via stdin/stdout; running multiple clients requires launching multiple server processes
- What happens when the client disconnects unexpectedly mid-session? → The server cleans up the session and exits cleanly; the server process does not auto-restart — the IDE (e.g., VS Code or JetBrains with Copilot) is responsible for relaunching the process
- What happens if the server encounters an unrecoverable internal error? → The server logs the error to stderr and exits with a non-zero status code; it does not silently hang
- What happens when a device command is attempted with a device that does not exist? → Returns a clear "device not found" error

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support a standard session handshake that establishes a connection before any data or actions are accessible
- **FR-002**: The system MUST allow clients to discover all available data sources, their names, and descriptions in a single request
- **FR-003**: The system MUST allow clients to read live data from any registered data source by name
- **FR-004**: The system MUST allow clients to discover all available actions, their names, descriptions, and required inputs
- **FR-005**: The system MUST allow clients to execute any registered action and receive a structured result or error
- **FR-006**: The system MUST allow clients to discover and retrieve pre-defined prompt templates with argument substitution
- **FR-007**: The system MUST serve data from the existing home lab data sources — devices, services, cluster services, health — without duplicating that data
- **FR-008**: The system MUST be runnable as a standalone process independently of the existing HTTP API server
- **FR-009**: The system MUST be launchable via the existing project build tooling (Makefile)
- **FR-010**: The system MUST return structured, descriptive errors for all failure cases
- **FR-011**: The system MUST handle concurrent requests on the same session safely without data corruption or crashes
- **FR-012**: The system MUST expose the home automation device catalogue as a readable data source; the device list is sourced from the same mock device store used by the device command action — no separate data store is introduced

### Key Entities

- **Session**: Represents an active connection between an AI client and the MCP server; tracks whether the handshake has been completed and what the client is capable of
- **Resource**: A named, read-only data source (e.g., "home automation devices", "system health"); has a name, description, and content type
- **Tool**: A named, executable action (e.g., "device command"); has a name, description, and a definition of required and optional inputs
- **Prompt Template**: A named, parameterised text template (e.g., "device control"); has a name, description, defined arguments, and renders to a ready-to-use message

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An AI client can connect, complete the handshake, and enumerate all available resources and tools within 2 seconds of the server starting
- **SC-002**: All data sources return accurate, live results in under 500 milliseconds under normal operating conditions
- **SC-003**: Device command execution completes end-to-end in under 1 second
- **SC-004**: The server correctly handles 10 simultaneous requests without errors or data inconsistencies
- **SC-005**: Test coverage for the new MCP server component is at or above 80%
- **SC-006**: The MCP server binary builds and starts successfully from a clean checkout
- **SC-007**: All existing tests in the project continue to pass after this feature is added
- **SC-008**: A Copilot MCP server configuration snippet (`.vscode/mcp.json` for VS Code, equivalent for JetBrains) is provided in the project documentation so a developer can wire the MCP server into GitHub Copilot without additional research

---

## Assumptions

- The MCP server communicates over **standard input/output** (stdin/stdout), which is the standard connection method used by GitHub Copilot in VS Code and JetBrains, and all MCP-compatible clients for local server processes
- The MCP server runs as a **separate process** from the existing HTTP API server but shares the same codebase
- **No authentication or authorisation** is required for the initial version — the server is assumed to run locally and be trusted by the process that launches it
- The existing mock data used by the home automation and cluster services is acceptable for this version; connecting to real external systems is out of scope
- HTTP-based transport (e.g., Server-Sent Events) is **out of scope** for this feature; stdin/stdout only
- The targeted MCP protocol version is the **current stable specification** as of the feature creation date
- Each running MCP server process serves **exactly one client** via stdin/stdout; running multiple simultaneous Copilot sessions requires launching multiple server processes — this is standard MCP stdio behaviour and is managed by the IDE
- A **Copilot MCP configuration snippet** (`.vscode/mcp.json` for VS Code, equivalent for JetBrains) must be included in project documentation so developers can connect without additional research
