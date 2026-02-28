# Feature Specification: Home Lab API Service

**Feature Branch**: `001-homelab-api-service`  
**Created**: 2026-02-28  
**Status**: Draft  
**Input**: User description: "Building an API service for self-hosted K8s cluster to provide information about home-lab services. POC with mocked REST endpoints for HomeAssistant. Backend services reusable for future MCP tools. Initial implementation fully mocked."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Query HomeAssistant Device Status (Priority: P1)

As a home lab administrator, I want to query the status of HomeAssistant devices via REST API so that I can integrate home automation data into dashboards and monitoring tools.

**Why this priority**: Core MVP functionality - demonstrates API working with a concrete service and provides immediate value for monitoring.

**Independent Test**: Can be fully tested by making HTTP GET requests to `/api/v1/homeassistant/devices` and `/api/v1/homeassistant/devices/{id}` and validating JSON responses with mocked device data.

**Acceptance Scenarios**:

1. **Given** the API service is running, **When** I GET `/api/v1/homeassistant/devices`, **Then** I receive a JSON array of all HomeAssistant devices with status
2. **Given** a device exists with ID "light.living_room", **When** I GET `/api/v1/homeassistant/devices/light.living_room`, **Then** I receive JSON with device details and current state
3. **Given** a device ID that doesn't exist, **When** I GET `/api/v1/homeassistant/devices/invalid`, **Then** I receive 404 with error message

---

### User Story 2 - Health Check and Service Discovery (Priority: P1)

As a DevOps engineer, I want to check the API health status and discover available services so that I can monitor the API and understand what integrations are available.

**Why this priority**: Essential for deployment and monitoring - required for K8s liveness/readiness probes and service discovery.

**Independent Test**: Can be fully tested by making GET requests to `/health` and `/api/v1/services` endpoints without any other functionality.

**Acceptance Scenarios**:

1. **Given** the API is running, **When** I GET `/health`, **Then** I receive 200 OK with JSON containing service status and uptime
2. **Given** the API is running, **When** I GET `/api/v1/services`, **Then** I receive JSON listing all available service integrations (e.g., homeassistant)
3. **Given** the API is unhealthy, **When** I GET `/health`, **Then** I receive 503 with details about what's failing

---

### User Story 3 - Control HomeAssistant Devices (Priority: P2)

As a home lab administrator, I want to send commands to HomeAssistant devices via REST API so that I can automate control flows from external systems.

**Why this priority**: Builds on P1 read functionality, adds write capability for more advanced automation scenarios.

**Independent Test**: Can be tested by making POST requests to `/api/v1/homeassistant/devices/{id}/command` with action payloads and validating mocked command execution.

**Acceptance Scenarios**:

1. **Given** a controllable device "light.living_room", **When** I POST `{"action": "turn_on"}` to `/api/v1/homeassistant/devices/light.living_room/command`, **Then** I receive 200 with success confirmation
2. **Given** an invalid action, **When** I POST `{"action": "invalid"}`, **Then** I receive 400 with validation error
3. **Given** a read-only device "sensor.temperature", **When** I POST any command, **Then** I receive 405 Method Not Allowed

---

### User Story 4 - Query Cluster Services Info (Priority: P3)

As a home lab administrator, I want to query general information about services running in my K8s cluster so that I can get a holistic view of my infrastructure.

**Why this priority**: Nice-to-have for comprehensive monitoring, but not critical for MVP. Can be added after core HomeAssistant integration is proven.

**Independent Test**: Can be tested independently by making GET requests to `/api/v1/cluster/services` and receiving mocked cluster information.

**Acceptance Scenarios**:

1. **Given** the API is running, **When** I GET `/api/v1/cluster/services`, **Then** I receive JSON with list of K8s services and their status
2. **Given** a service name filter, **When** I GET `/api/v1/cluster/services?name=homeassistant`, **Then** I receive filtered results

---

### Edge Cases

- What happens when the API receives malformed JSON in request body?
- How does the system handle very long device IDs (>255 characters)?
- What happens if a client sends requests faster than rate limits?
- How does the API behave when request includes unsupported API version (e.g., `/api/v2/...`)?
- What happens when a device ID contains special characters or URL-unsafe characters?
- How does the system handle concurrent requests to the same device?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose RESTful HTTP endpoints at `/api/v1/*` for all service integrations
- **FR-002**: System MUST return JSON responses with consistent error structure: `{"error": "message", "code": "ERROR_CODE", "details": {...}}`
- **FR-003**: System MUST provide a `/health` endpoint returning service health status
- **FR-004**: System MUST support GET `/api/v1/homeassistant/devices` to list all devices (mocked data)
- **FR-005**: System MUST support GET `/api/v1/homeassistant/devices/{id}` to retrieve specific device details
- **FR-006**: System MUST support POST `/api/v1/homeassistant/devices/{id}/command` to send device commands (mocked execution)
- **FR-007**: System MUST validate all incoming request payloads and return 400 for invalid data
- **FR-008**: System MUST implement rate limiting on all endpoints (100 requests/minute per client IP)
- **FR-009**: System MUST use structured logging with request IDs for all operations
- **FR-010**: System MUST return appropriate HTTP status codes (200, 201, 400, 404, 405, 429, 500, 503)
- **FR-011**: System MUST implement graceful shutdown handling SIGTERM and SIGINT signals
- **FR-012**: Backend services MUST be structured as reusable packages under `/internal` that can be wrapped as MCP tools in future
- **FR-013**: System MUST use mocked data for all responses (no live integrations in POC)
- **FR-014**: System MUST include request ID in all log entries and error responses
- **FR-015**: System MUST support CORS headers for browser-based clients [NEEDS CLARIFICATION: which origins should be allowed?]

### Key Entities *(include if feature involves data)*

- **Device**: Represents a HomeAssistant device (light, sensor, switch, etc.). Key attributes: ID (string), Name (string), Type (enum: light/sensor/switch/binary_sensor), State (string: on/off/value), Attributes (map of additional properties like brightness, temperature, etc.)

- **ServiceInfo**: Represents a K8s cluster service. Key attributes: Name (string), Namespace (string), Status (enum: running/stopped/error), Endpoints (array of URLs)

- **Command**: Represents an action to execute on a device. Key attributes: DeviceID (string), Action (string: turn_on/turn_off/set_brightness), Parameters (map of action-specific params)

- **HealthStatus**: Represents overall service health. Key attributes: Status (enum: healthy/degraded/unhealthy), Components (map of component name to status), Uptime (duration)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: API responds to health check requests in under 50ms
- **SC-002**: All REST endpoints return responses within 200ms for mocked data
- **SC-003**: System successfully handles 100 concurrent requests without errors
- **SC-004**: API provides OpenAPI/Swagger documentation accessible at `/api/docs` with interactive UI
- **SC-005**: OpenAPI spec available at `/api/docs/swagger.json` for programmatic access
- **SC-005**: All endpoints return correctly formatted JSON matching documented schema
- **SC-006**: Backend service packages can be imported and used without HTTP layer (demonstrating MCP-readiness)
- **SC-007**: API passes all automated tests with >80% code coverage
- **SC-008**: Docker image builds successfully and runs in K8s cluster with <100MB memory footprint
- **SC-009**: API gracefully shuts down within 5 seconds of receiving SIGTERM
- **SC-010**: All API operations are logged with structured logging including request IDs

## Technical Constraints *(from Constitution)*

### Go Standards
- Go 1.24 minimum version
- Follow official Go Code Review Comments guidelines
- All code must pass `gofmt`, `go vet`, and `golangci-lint`
- Explicit error handling (no ignored errors)
- Use `context.Context` for cancellation and timeouts

### Testing Requirements
- Minimum 80% code coverage
- Table-driven tests for handlers
- Integration tests for all endpoints
- Mock all external dependencies

### Architecture Principles
- Backend services in `/internal` must be independent of HTTP layer
- Services should accept interfaces, return concrete types
- Design services to be wrappable as MCP tools:
  - Pure functions where possible
  - Clear input/output contracts
  - No HTTP coupling in business logic
  - Context-aware for cancellation

### Project Structure
```
/cmd/api              - HTTP server entry point
/internal/homeassistant - HomeAssistant service logic
/internal/cluster     - Cluster service logic  
/internal/middleware  - HTTP middleware (logging, rate limiting, etc.)
/internal/handlers    - HTTP handlers (thin wrappers around services)
/internal/models      - Shared data models
/api                  - OpenAPI specifications
/tests                - Integration tests
```

## API Specification Preview

### Endpoints

```
GET  /health
GET  /api/v1/services
GET  /api/v1/homeassistant/devices
GET  /api/v1/homeassistant/devices/{id}
POST /api/v1/homeassistant/devices/{id}/command
GET  /api/v1/cluster/services
```

### Example Response Formats

**GET /api/v1/homeassistant/devices**
```json
{
  "devices": [
    {
      "id": "light.living_room",
      "name": "Living Room Light",
      "type": "light",
      "state": "on",
      "attributes": {
        "brightness": 200,
        "color_temp": 370
      },
      "last_updated": "2026-02-28T10:30:00Z"
    }
  ],
  "total": 1,
  "request_id": "req_abc123"
}
```

**Error Response**
```json
{
  "error": "Device not found",
  "code": "DEVICE_NOT_FOUND",
  "details": {
    "device_id": "light.invalid"
  },
  "request_id": "req_xyz789"
}
```

## Out of Scope

- Live HomeAssistant integration (mocked only in POC)
- Authentication/authorization (will add in future iteration)
- Persistence layer (all data in-memory mocks)
- WebSocket support for real-time updates
- Custom device types beyond light/sensor/switch
- Historical data or time-series metrics
- Device grouping or scenes
- Automation rules or workflows

## Open Questions

1. Should the API support filtering devices by type or room in the list endpoint?
2. What rate limits are appropriate for a home lab environment?
3. Should we include metrics endpoint (Prometheus format) in the POC or defer to later?
4. Do we need request/response size limits?
5. Should the mocked data be configurable via environment variables or config file?

## Dependencies

- Go 1.24+
- Standard library (net/http, encoding/json, log/slog, context)
- External: `chi` router (or similar minimal router), `testify` for test assertions
- Docker for containerization
- K8s cluster for deployment (assumes existing setup)

## Next Steps

1. Create detailed API design (OpenAPI spec)
2. Generate implementation plan with task breakdown
3. Set up project structure and boilerplate
4. Implement P1 user stories (device queries + health check)
5. Add integration tests
6. Create Dockerfile and K8s manifests
7. Implement P2 functionality (device commands)
8. Add rate limiting and monitoring
9. Documentation and deployment guide

