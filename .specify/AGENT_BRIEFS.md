# Agent Brief Templates: Home Lab API Service Implementation

This directory contains brief templates for dispatching work to Copilot Cloud Agents. Each batch is self-contained and can be executed independently after prerequisites are met.

---

## Batch 1: Phase 0 - Research & Validation

**Issues**: #27-30 (4 tasks)
**Duration**: 1-2 days
**Priority**: P0 (Blocking for all other phases)

### Agent Brief: Batch 1

```markdown
# Agent Brief: Phase 0 - Research & Validation

## Overview
Validate technology choices for the Home Lab API Service through benchmarking and research.

## Prerequisites
✅ Go 1.24+ installed
✅ Gin, jsoniter, and sync module documentation available

## Linked Issues
- [#27] T001: Benchmark Gin framework basic routing
- [#28] T002: Benchmark jsoniter vs stdlib encoding/json
- [#29] T003: Benchmark sync.Map vs RWMutex for device storage
- [#30] T004: Complete research.md with all benchmark results

## Tasks
1. Create `research/gin_benchmark_test.go` - Compare Gin routing vs stdlib net/http
2. Create `research/json_benchmark_test.go` - Test Device struct encoding (50 devices)
3. Create `research/storage_benchmark_test.go` - Concurrent reads (100 goroutines)
4. Document all findings in `.specify/features/001-homelab-api-service/research.md`

## Acceptance Criteria
- [ ] All benchmark files created and runnable
- [ ] Results show <10ms health check target is achievable
- [ ] jsoniter performance improvement documented (target: 2-3x)
- [ ] sync.Map vs RWMutex comparison completed
- [ ] research.md complete with recommendations
- [ ] PR created and ready for review

## Testing
No tests required for research phase (pure benchmarks)

## Git Workflow
```bash
git checkout -b batch-1-phase-0-research
# Create benchmark files and documentation
git commit -m "batch-1: Phase 0 research and validation"
git push origin batch-1-phase-0-research
# Create PR with benchmark results
```

## Review Checklist
- [ ] All benchmarks completed
- [ ] Results documented clearly
- [ ] Recommendations provided
- [ ] Performance targets validated

## Notes
- Benchmarks must be portable (work on Linux, Mac, Windows)
- Include detailed output examples in research.md
- Document Go version used for benchmarks
```

---

## Batch 2: Phase 1 - Foundation Setup

**Issues**: #31-39 (10 tasks)
**Duration**: 1-2 days
**Depends On**: Batch 1 ✅
**Priority**: P1 (Blocking for all feature work)

### Agent Brief: Batch 2

```markdown
# Agent Brief: Phase 1 - Foundation Setup

## Overview
Create complete project foundation including dependencies, directory structure, models, middleware, and server setup.

## Prerequisites
✅ Research complete (Batch 1 merged)
✅ Go 1.24+ installed
✅ Project can be initialized

## Linked Issues
- [#31] T010: Initialize Go module and dependencies
- [#32] T011: Create project directory structure
- [#33] T012: Define core data models
- [#34] T013: Create service interfaces (placeholder)
- [#35] T014: Implement request ID middleware
- [#36] T015: Implement logging middleware with slog
- [#37] T016: Implement panic recovery middleware
- [#38] T017: Setup Gin server with graceful shutdown
- [#39] T018: Create response helper utilities
- [#40] T019: Create Makefile with build targets

## Tasks

### Task 1-2: Setup (Sequential)
1. **T010**: Update go.mod to Go 1.24 and add all dependencies
   - gin-gonic/gin v1.10.0
   - swaggo packages (swag, gin-swagger, files)
   - json-iterator/go v1.1.12
   - stretchr/testify v1.9.0
   - Run `go mod tidy`

2. **T011**: Create complete directory structure
   - cmd/api/
   - internal/{homeassistant,cluster,health,middleware,handlers,models,server}/
   - tests/{integration,load}/
   - api/
   - deployments/{k8s}/

### Task 3: Models (Parallelizable with middleware)
3. **T012**: Create core data models in internal/models/
   - device.go: Device struct with all required fields
   - error.go: ErrorResponse struct
   - health.go: HealthStatus struct
   - Add proper JSON tags for all exported fields

### Task 4-6: Middleware (Parallelizable)
4. **T014**: Request ID middleware
   - Generate UUID per request
   - Add to context and X-Request-ID header

5. **T015**: Logging middleware with slog
   - Log request start/end with duration
   - Extract request ID from context
   - Log method, path, status, response time

6. **T016**: Panic recovery middleware
   - Catch panics and log with stack trace
   - Return 500 error response

### Task 7-9: Server & Helpers (Sequential)
7. **T017**: Gin server setup with graceful shutdown
   - internal/server/server.go
   - internal/server/shutdown.go
   - cmd/api/main.go entry point
   - Handle SIGTERM and SIGINT

8. **T018**: Response helper utilities
   - JSONSuccess(c *gin.Context, code int, data)
   - JSONError(c *gin.Context, code int, error, details)
   - JSONCreated(c *gin.Context, data)

9. **T019**: Makefile with build targets
   - build, test, lint, run, clean, swagger targets
   - Coverage report generation

## Acceptance Criteria
- [ ] go.mod properly initialized with all dependencies
- [ ] All directories created with proper Go package structure
- [ ] Models defined with JSON tags
- [ ] All middleware implemented with unit tests
- [ ] Server initializes and handles shutdown gracefully
- [ ] Response helpers work with all common scenarios
- [ ] Makefile targets all working
- [ ] `make test` passes
- [ ] Coverage report generated

## Testing
```bash
make test          # Run all tests
make lint          # Check code quality
make build         # Verify binary builds
make run           # Test server starts
```

## Git Workflow
```bash
git checkout -b batch-2-phase-1-foundation
# Complete all tasks
# Ensure all tests pass
git commit -m "batch-2: Phase 1 foundation setup"
git push origin batch-2-phase-1-foundation
# Create PR
```

## Review Checklist
- [ ] All dependencies installed
- [ ] Project structure correct
- [ ] All middleware tested
- [ ] Server setup complete
- [ ] Makefile working
- [ ] Coverage report included
- [ ] Ready for Swagger setup and feature development

## Notes
- Foundation is CRITICAL - no feature work can proceed until merged
- All middleware must be properly ordered
- Graceful shutdown testing is important (SIGTERM/SIGINT)
- Set conservative defaults for server (3s shutdown timeout, 8080 port)
```

---

## Batch 3: Phase 1.5 - Swagger Setup

**Issues**: #40-43 (4 tasks)
**Duration**: 1 day
**Depends On**: Batch 2 ✅
**Priority**: P1 (Blocking for handler implementation)

### Agent Brief: Batch 3

```markdown
# Agent Brief: Phase 1.5 - Swagger/OpenAPI Setup

## Overview
Setup Swagger/OpenAPI documentation infrastructure with code generation and UI hosting.

## Prerequisites
✅ Batch 2 (Foundation) merged to main
✅ Go 1.24+ installed
✅ Makefile ready

## Linked Issues
- [#41] T019a: Install swag CLI tool
- [#42] T019b: Add Swagger general API info annotations
- [#43] T019c: Setup Swagger UI routes
- [#44] T019d: Generate initial Swagger docs

## Tasks

1. **T019a**: Install swag CLI
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   swag --version  # Verify installation
   ```

2. **T019b**: Add Swagger metadata to main.go
   ```go
   // @title Home Lab API Service
   // @version 1.0
   // @description REST API for home lab monitoring and control
   // @host localhost:8080
   // @BasePath /
   ```

3. **T019c**: Setup Swagger UI routes in internal/server/server.go
   - Add imports for ginSwagger and swaggerFiles
   - Route: GET /api/docs/* → Swagger UI
   - Auto-serve swagger.json at /api/docs/swagger.json

4. **T019d**: Generate initial documentation
   ```bash
   make swagger
   # Verify api/docs.go and api/swagger.json created
   git add api/docs.go api/swagger.json
   ```

## Acceptance Criteria
- [ ] swag CLI installed and verified
- [ ] Swagger metadata in main.go
- [ ] Swagger UI routes configured
- [ ] Initial docs generated successfully
- [ ] Swagger UI accessible at http://localhost:8080/api/docs/index.html
- [ ] OpenAPI spec downloadable at /api/docs/swagger.json
- [ ] `make swagger` target working
- [ ] Generated files committed to repository

## Testing
```bash
make run                              # Start server
curl http://localhost:8080/api/docs/swagger.json  # Test spec access
# Browser: http://localhost:8080/api/docs/index.html
```

## Git Workflow
```bash
git checkout -b batch-3-phase-1-5-swagger
# Complete Swagger setup
git commit -m "batch-3: Phase 1.5 Swagger/OpenAPI setup"
git push origin batch-3-phase-1-5-swagger
# Create PR
```

## Review Checklist
- [ ] swag CLI installed
- [ ] Swagger annotations added
- [ ] Routes configured
- [ ] Docs generated
- [ ] UI accessible in browser
- [ ] OpenAPI spec valid
- [ ] Make target working

## Notes
- This is CRITICAL for handler work - handlers will include inline Swagger annotations
- Generated files (api/docs.go, api/swagger.json) should be committed
- After this batch, handlers will include Swagger comments
- Each handler task will regenerate docs with `make swagger`
```

---

## Batch 4: Phase 2 - User Story 1 (Device Status)

**Issues**: #44-50 (7 tasks)
**Duration**: 2-3 days
**Depends On**: Batch 3 ✅
**Priority**: P1 (MVP Core)

### Agent Brief: Batch 4

```markdown
# Agent Brief: Phase 2 - User Story 1: Query HomeAssistant Device Status

## Overview
Implement REST API endpoints for querying HomeAssistant device status with mocked data.

## Prerequisites
✅ Batch 3 (Swagger) merged to main
✅ Swagger UI running locally

## Linked Issues
- [#51] T020: Write unit tests for HomeAssistant service
- [#52] T021: Write integration tests for device endpoints
- [#53] T022: Create mock device data
- [#54] T023: Implement HomeAssistant service
- [#55] T024: Implement device list handler with Swagger
- [#56] T025: Implement device detail handler with Swagger
- [#57] T026: Register device endpoints

## User Story
**As a** home lab administrator
**I want to** query the status of HomeAssistant devices via REST API
**So that** I can integrate home automation data into dashboards

## Tasks (TDD: Write tests FIRST)

### Test-Driven Development Phase
1. **T020**: Write unit tests for service (FAIL intentionally)
   - Create internal/homeassistant/service_test.go
   - Tests: ListDevices, GetDevice, error cases

2. **T021**: Write integration tests for endpoints (FAIL intentionally)
   - Create tests/integration/devices_test.go
   - Tests: GET /api/v1/homeassistant/devices, GET /{id}, 404 handling

### Implementation Phase
3. **T022**: Create mock device data
   - internal/homeassistant/mock_data.go
   - 50 realistic devices (lights, sensors, switches, binary_sensors)

4. **T023**: Implement service with sync.Map storage
   - internal/homeassistant/service.go
   - ListDevices() and GetDevice() methods

5. **T024**: Device list handler with Swagger
   - internal/handlers/homeassistant.go
   - Support ?type= query parameter for filtering
   - Add Swagger annotations
   - Run `make swagger`

6. **T025**: Device detail handler with Swagger
   - Add GetDeviceHandler() to handlers
   - 404 handling for non-existent devices
   - Add Swagger annotations
   - Run `make swagger`

7. **T026**: Register endpoints in router
   - Update internal/server/server.go
   - Route: GET /api/v1/homeassistant/devices
   - Route: GET /api/v1/homeassistant/devices/:id

## Acceptance Criteria
- [ ] Unit tests written and FAIL before implementation
- [ ] Integration tests written and FAIL before implementation
- [ ] T020 tests now PASS with service implementation
- [ ] T021 tests now PASS with handler implementation
- [ ] Mock devices realistic with all types represented
- [ ] Type filtering working correctly
- [ ] 404 errors returned for non-existent devices
- [ ] Swagger annotations on both handlers
- [ ] `make swagger` regenerates docs cleanly
- [ ] All endpoints tested manually with curl
- [ ] Coverage >80% for new code

## Manual Testing
```bash
# List all devices
curl http://localhost:8080/api/v1/homeassistant/devices

# Filter by type
curl "http://localhost:8080/api/v1/homeassistant/devices?type=light"

# Get specific device
curl http://localhost:8080/api/v1/homeassistant/devices/light.living_room

# Non-existent device (should be 404)
curl http://localhost:8080/api/v1/homeassistant/devices/invalid

# View Swagger UI
# Browser: http://localhost:8080/api/docs/index.html
```

## Git Workflow
```bash
git checkout -b batch-4-phase-2-us1-device-status
# Write tests first (should FAIL)
# Implement service
# Tests now PASS
# Implement handlers with Swagger
# Register routes
git commit -m "batch-4: Phase 2 US1 - Device status queries"
git push origin batch-4-phase-2-us1-device-status
# Create PR
```

## Review Checklist
- [ ] TDD workflow followed (tests written first)
- [ ] All tests passing
- [ ] Coverage >80%
- [ ] Swagger docs generated
- [ ] Manual curl tests successful
- [ ] Type filtering implemented
- [ ] Error handling correct
- [ ] Ready for US2

## Notes
- TDD is MANDATORY - tests must be written before implementation
- Type filtering is required (light, sensor, switch, binary_sensor)
- Mock data should be realistic and representative
- Swagger annotations must be present on both handlers
- This is first MVP feature - should be production-ready
```

---

## Batch 5: Phase 3 - User Story 2 (Health & Discovery)

**Issues**: #58-64 (7 tasks)
**Duration**: 2-3 days
**Depends On**: Batch 4 ✅
**Priority**: P1 (MVP Core)

### Agent Brief: Batch 5

```markdown
# Agent Brief: Phase 3 - User Story 2: Health Check and Service Discovery

## Overview
Implement K8s-ready health endpoint and service discovery API.

## Prerequisites
✅ Batch 4 (US1) merged to main
✅ Device endpoints working

## Linked Issues
- [#65] T030: Write unit tests for health checker
- [#66] T031: Write integration tests for health endpoint
- [#67] T032: Implement health checker service
- [#68] T033: Implement health endpoint handler
- [#69] T034: Implement services discovery endpoint
- [#70] T035: Register health and services endpoints

## User Story
**As a** DevOps engineer
**I want to** check API health status and discover available services
**So that** I can monitor the API and understand available integrations

## Tasks (TDD: Write tests FIRST)

### Test-Driven Development Phase
1. **T030**: Write unit tests for health checker (FAIL intentionally)
   - Create internal/health/checker_test.go
   - Tests: Check status, uptime calculation

2. **T031**: Write integration tests for health endpoint (FAIL intentionally)
   - Create tests/integration/health_test.go
   - Tests: /health returns 200, response format, <50ms response time

### Implementation Phase
3. **T032**: Implement health checker service
   - internal/health/checker.go
   - Track uptime from server start time
   - Return healthy status with component info

4. **T033**: Health endpoint handler with Swagger
   - internal/handlers/health.go
   - Optimize for <50ms response time
   - Consider pre-marshaled JSON for <10ms
   - Add Swagger annotations
   - Run `make swagger`

5. **T034**: Services discovery endpoint with Swagger
   - internal/handlers/services.go
   - Return static list of available services
   - Add Swagger annotations
   - Run `make swagger`

6. **T035**: Register health and services endpoints
   - Update internal/server/server.go
   - Route: GET /health (no /api/v1 prefix)
   - Route: GET /api/v1/services

## Acceptance Criteria
- [ ] Unit tests written and FAIL before implementation
- [ ] Integration tests written and FAIL before implementation
- [ ] T030 tests now PASS with health checker implementation
- [ ] T031 tests now PASS with handler implementation
- [ ] Health check response time <50ms (target: <10ms)
- [ ] Uptime calculation correct
- [ ] Services list returned correctly
- [ ] Swagger annotations on both handlers
- [ ] `make swagger` regenerates docs cleanly
- [ ] K8s can use /health for liveness/readiness probes
- [ ] Coverage >80% for new code

## Manual Testing
```bash
# Health check
curl http://localhost:8080/health
# Expected response includes: status, components, uptime_seconds

# Services discovery
curl http://localhost:8080/api/v1/services
# Expected response includes: homeassistant service info

# Performance test
time curl http://localhost:8080/health  # Should be <50ms
```

## Git Workflow
```bash
git checkout -b batch-5-phase-3-us2-health-discovery
# Write tests first (should FAIL)
# Implement health checker
# Implement handlers
# Register routes
# Run `make swagger`
git commit -m "batch-5: Phase 3 US2 - Health checks and discovery"
git push origin batch-5-phase-3-us2-health-discovery
# Create PR
```

## Review Checklist
- [ ] TDD workflow followed
- [ ] All tests passing
- [ ] Health check <50ms
- [ ] Uptime calculation verified
- [ ] Services list correct
- [ ] Swagger docs generated
- [ ] K8s probe compatible
- [ ] Coverage >80%

## Notes
- Health check performance is critical - measure carefully
- Consider pre-marshaling response JSON for <10ms target
- This endpoint will be hit frequently by K8s, so optimization matters
- Services list is static for now - easy to update with actual K8s query later
- /health endpoint has NO /api/v1 prefix (K8s standard)
- /api/v1/services IS under API versioning
```

---

## Batch 6: Phase 4 - Performance Optimization & Middleware

**Issues**: #71-75 (5 tasks)
**Duration**: 1-2 days
**Depends On**: Batch 5 ✅
**Priority**: P2

### Agent Brief: Batch 6

```markdown
# Agent Brief: Phase 4 - Performance Optimization & Middleware

## Overview
Add rate limiting, CORS middleware, and performance optimizations for production readiness.

## Prerequisites
✅ Batch 5 (US2) merged to main
✅ All endpoints working

## Linked Issues
- [#76] T040: Implement rate limiting middleware
- [#77] T041: Implement CORS middleware
- [#78] T042: Optimize JSON encoding with jsoniter
- [#79] T043: Implement response pooling
- [#80] T044: Update server router with middleware

## Tasks

1. **T040**: Rate limiting middleware
   - internal/middleware/ratelimit.go
   - Token bucket: 500 req/min per IP
   - Return 429 Too Many Requests
   - Add X-RateLimit-Limit and X-RateLimit-Remaining headers
   - Unit test included

2. **T041**: CORS middleware
   - internal/middleware/cors.go
   - Allow origin: http://localhost:3000 (configurable via env)
   - Allow methods: GET, POST, OPTIONS
   - Allow headers: Content-Type, Authorization
   - Unit test included

3. **T042**: Optimize JSON encoding with jsoniter
   - Update internal/handlers/response.go
   - Replace stdlib json with jsoniter
   - Benchmark comparison (target: 2-3x improvement)
   - Document results in performance notes

4. **T043**: Response object pooling
   - Update internal/handlers/homeassistant.go
   - Add sync.Pool for DeviceListResponse
   - Pre-allocate response structs
   - Benchmark allocation reduction
   - Document results

5. **T044**: Register middleware in server router
   - Update internal/server/server.go
   - Apply rate limiting to /api/v1/* routes
   - Apply CORS to all routes
   - Correct ordering: Recovery → RequestID → Logging → RateLimit → CORS → Handler

## Acceptance Criteria
- [ ] Rate limiting working (500 req/min per IP)
- [ ] CORS headers present on responses
- [ ] jsoniter integrated with 2-3x improvement documented
- [ ] Response pooling reducing allocations
- [ ] All middleware unit tests passing
- [ ] Middleware ordering correct
- [ ] No performance regression in existing endpoints
- [ ] All endpoints still responding correctly

## Performance Benchmarking
```bash
# Before optimization
make test  # Benchmark results
# After optimization
make test  # Compare results
# Expected: 2-3x faster JSON encoding, reduced allocations
```

## Testing
```bash
make test  # All tests passing
make run   # Server starts, endpoints accessible

# Manual CORS test
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS http://localhost:8080/api/v1/homeassistant/devices

# Manual rate limit test
for i in {1..100}; do curl http://localhost:8080/health; done
# Check for 429 responses after 500 hits
```

## Git Workflow
```bash
git checkout -b batch-6-phase-4-performance
# Implement rate limiting
# Implement CORS
# Optimize JSON and response pooling
# Update server middleware ordering
git commit -m "batch-6: Phase 4 performance optimization"
git push origin batch-6-phase-4-performance
# Create PR with benchmarks
```

## Review Checklist
- [ ] Rate limiting working correctly
- [ ] CORS headers present
- [ ] jsoniter integrated
- [ ] Response pooling implemented
- [ ] Benchmarks documented
- [ ] Middleware ordering correct
- [ ] All tests passing
- [ ] No regressions

## Notes
- Rate limiting per IP address is important
- CORS origin is configurable via environment variable for flexibility
- jsoniter should be measurably faster than stdlib
- Response pooling reduces GC pressure
- Middleware ordering matters for correct header propagation
```

---

## Batch 7: Phase 5 - User Story 3 (Device Control)

**Issues**: #81-86 (6 tasks)
**Duration**: 2-3 days
**Depends On**: Batch 6 ✅
**Priority**: P2

### Agent Brief: Batch 7

```markdown
# Agent Brief: Phase 5 - User Story 3: Control HomeAssistant Devices

## Overview
Implement device control endpoints with command execution and validation.

## Prerequisites
✅ Batch 6 (Performance) merged to main
✅ Device status queries working

## Linked Issues
- [#87] T050: Write unit tests for command execution
- [#88] T051: Write integration tests for command endpoint
- [#89] T052: Define Command model
- [#90] T053: Implement command execution service
- [#91] T054: Implement command endpoint handler
- [#92] T055: Register command endpoint

## User Story
**As a** home lab administrator
**I want to** send commands to HomeAssistant devices via REST API
**So that** I can automate control flows from external systems

## Tasks (TDD: Write tests FIRST)

### Test-Driven Development Phase
1. **T050**: Write unit tests for command execution (FAIL intentionally)
   - Update internal/homeassistant/service_test.go
   - Tests: valid command, invalid device, read-only device, invalid action

2. **T051**: Write integration tests for command endpoint (FAIL intentionally)
   - Update tests/integration/devices_test.go
   - Tests: POST success, invalid action (400), read-only device (405)

### Implementation Phase
3. **T052**: Define Command model
   - Create internal/homeassistant/types.go
   - Command struct: Action, Parameters
   - JSON tags and validation

4. **T053**: Implement command execution in service
   - Update internal/homeassistant/service.go
   - ExecuteCommand() method
   - Validation: device exists, device is controllable, action valid
   - Mock execution (just return nil for POC)

5. **T054**: Command endpoint handler with Swagger
   - Update internal/handlers/homeassistant.go
   - ExecuteCommandHandler() method
   - Parse JSON body
   - Error handling: 400, 404, 405 status codes
   - Add Swagger annotations
   - Run `make swagger`

6. **T055**: Register command endpoint
   - Update internal/server/server.go
   - Route: POST /api/v1/homeassistant/devices/:id/command

## Acceptance Criteria
- [ ] Unit tests written and FAIL before implementation
- [ ] Integration tests written and FAIL before implementation
- [ ] T050 tests now PASS with service implementation
- [ ] T051 tests now PASS with handler implementation
- [ ] Command model properly defined
- [ ] Validation: device exists check
- [ ] Validation: device type controllable check (not sensor)
- [ ] Validation: action valid for device type
- [ ] Correct HTTP status codes (200, 400, 404, 405)
- [ ] Swagger annotations on handler
- [ ] `make swagger` regenerates docs cleanly
- [ ] Coverage >80% for new code

## Manual Testing
```bash
# Valid command
curl -X POST http://localhost:8080/api/v1/homeassistant/devices/light.living_room/command \
  -H "Content-Type: application/json" \
  -d '{"action":"turn_on"}'
# Expected: 200

# Invalid action
curl -X POST http://localhost:8080/api/v1/homeassistant/devices/light.living_room/command \
  -H "Content-Type: application/json" \
  -d '{"action":"invalid"}'
# Expected: 400

# Read-only device (sensor)
curl -X POST http://localhost:8080/api/v1/homeassistant/devices/sensor.temperature/command \
  -H "Content-Type: application/json" \
  -d '{"action":"turn_on"}'
# Expected: 405

# Non-existent device
curl -X POST http://localhost:8080/api/v1/homeassistant/devices/invalid/command \
  -H "Content-Type: application/json" \
  -d '{"action":"turn_on"}'
# Expected: 404
```

## Git Workflow
```bash
git checkout -b batch-7-phase-5-us3-device-control
# Write tests first (should FAIL)
# Define Command model
# Implement service logic
# Tests now PASS
# Implement handler with Swagger
# Register endpoint
git commit -m "batch-7: Phase 5 US3 - Device control"
git push origin batch-7-phase-5-us3-device-control
# Create PR
```

## Review Checklist
- [ ] TDD workflow followed
- [ ] All tests passing
- [ ] Command validation correct
- [ ] HTTP status codes correct
- [ ] Swagger docs generated
- [ ] Manual curl tests successful
- [ ] Coverage >80%
- [ ] Ready for testing phase

## Notes
- Sensor devices should return 405 (Method Not Allowed)
- Command execution is mocked (just validation for POC)
- Validate action is supported before execution
- Consider which device types support which actions
```

---

## Batch 8: Phase 6 - User Story 4 (Cluster Services)

**Issues**: #93-98 (6 tasks)
**Duration**: 1-2 days
**Depends On**: Batch 7 ✅
**Priority**: P3

### Agent Brief: Batch 8

```markdown
# Agent Brief: Phase 6 - User Story 4: Query Cluster Services Info

## Overview
Implement cluster service discovery endpoint for monitoring K8s services.

## Prerequisites
✅ Batch 7 (US3) merged to main
✅ All device endpoints working

## Linked Issues
- [#99] T060: Write unit tests for cluster service
- [#100] T061: Write integration tests for cluster endpoint
- [#101] T062: Create cluster service models
- [#102] T063: Implement cluster service
- [#103] T064: Implement cluster services handler
- [#104] T065: Register cluster endpoint

## User Story
**As a** home lab administrator
**I want to** query general information about services running in my K8s cluster
**So that** I can get a holistic view of my infrastructure

## Tasks (TDD: Write tests FIRST)

### Test-Driven Development Phase
1. **T060**: Write unit tests for cluster service (FAIL intentionally)
   - Create internal/cluster/service_test.go
   - Tests: list services, filter by name

2. **T061**: Write integration tests for cluster endpoint (FAIL intentionally)
   - Create tests/integration/cluster_test.go
   - Tests: GET /api/v1/cluster/services, filtering

### Implementation Phase
3. **T062**: Create cluster service models
   - Create internal/cluster/types.go
   - ServiceInfo struct: Name, Namespace, Status, Endpoints

4. **T063**: Implement cluster service
   - Create internal/cluster/service.go
   - ListServices() method
   - Mock K8s services (homeassistant, prometheus, grafana, etc.)

5. **T064**: Cluster services handler with Swagger
   - Create internal/handlers/cluster.go
   - ListClusterServicesHandler() method
   - Support ?name= query parameter for filtering
   - Add Swagger annotations
   - Run `make swagger`

6. **T065**: Register cluster endpoint
   - Update internal/server/server.go
   - Route: GET /api/v1/cluster/services

## Acceptance Criteria
- [ ] Unit tests written and FAIL before implementation
- [ ] Integration tests written and FAIL before implementation
- [ ] T060 tests now PASS with service implementation
- [ ] T061 tests now PASS with handler implementation
- [ ] Mock services realistic (homeassistant, prometheus, etc.)
- [ ] Service filtering working
- [ ] Swagger annotations on handler
- [ ] `make swagger` regenerates docs cleanly
- [ ] Coverage >80% for new code

## Manual Testing
```bash
# List all cluster services
curl http://localhost:8080/api/v1/cluster/services

# Filter by name
curl "http://localhost:8080/api/v1/cluster/services?name=homeassistant"

# View in Swagger UI
# Browser: http://localhost:8080/api/docs/index.html
```

## Git Workflow
```bash
git checkout -b batch-8-phase-6-us4-cluster-services
# Write tests first
# Implement service
# Tests now PASS
# Implement handler
# Register endpoint
git commit -m "batch-8: Phase 6 US4 - Cluster services discovery"
git push origin batch-8-phase-6-us4-cluster-services
# Create PR
```

## Review Checklist
- [ ] TDD workflow followed
- [ ] All tests passing
- [ ] Mock services realistic
- [ ] Filtering working
- [ ] Swagger docs generated
- [ ] Coverage >80%
- [ ] Ready for testing

## Notes
- This is P3 priority but should still follow TDD
- Mock services should be representative of typical K8s cluster
- Filtering is simple (optional - could be extended later)
```

---

## Batch 9: Phase 7 - Documentation & Testing

**Issues**: #105-108 (4 tasks)
**Duration**: 1 day
**Depends On**: Batch 8 ✅
**Priority**: P2

### Agent Brief: Batch 9

```markdown
# Agent Brief: Phase 7 - Documentation & Testing

## Overview
Complete Swagger UI verification, load testing, coverage validation, and API documentation.

## Prerequisites
✅ Batch 8 (US4) merged to main
✅ All 4 user stories implemented

## Linked Issues
- [#109] T070: Verify Swagger UI functionality
- [#110] T071: Write load tests for concurrent requests
- [#111] T072: Run final coverage check
- [#112] T073: Create API usage examples

## Tasks

1. **T070**: Verify Swagger UI Accessibility
   - Start server: `make run`
   - Access: http://localhost:8080/api/docs/index.html
   - Verify all endpoints documented
   - Test "Try it out" for each endpoint
   - Verify OpenAPI spec at /api/docs/swagger.json
   - Document URL in README

2. **T071**: Write load tests
   - Create tests/load/concurrent_test.go
   - 100 concurrent requests to /health
   - 100 concurrent requests to /api/v1/homeassistant/devices
   - Validate <200ms p99 response time
   - Run and document results

3. **T072**: Code coverage validation
   - Execute: `make test`
   - Generate coverage report
   - Validate ≥80% coverage
   - Document coverage percentage

4. **T073**: Create API usage documentation
   - Update README.md
   - Add curl examples for all endpoints
   - Document environment variables
   - Local development setup instructions

## Acceptance Criteria
- [ ] Swagger UI accessible and all endpoints visible
- [ ] "Try it out" works for sample endpoints
- [ ] Load tests created and passing
- [ ] Response times documented
- [ ] Coverage ≥80% validated
- [ ] README updated with examples
- [ ] All curl examples working
- [ ] Documentation complete and clear

## Testing
```bash
make test          # Generate coverage
make run           # Start server for manual Swagger testing
# Browser: http://localhost:8080/api/docs/index.html
go test ./tests/load -v  # Run load tests
```

## Git Workflow
```bash
git checkout -b batch-9-phase-7-documentation
# Verify Swagger UI
# Write load tests
# Validate coverage
# Update README
git commit -m "batch-9: Phase 7 documentation and testing"
git push origin batch-9-phase-7-documentation
# Create PR
```

## Review Checklist
- [ ] Swagger UI fully functional
- [ ] All endpoints documented
- [ ] Load tests complete
- [ ] Coverage ≥80%
- [ ] README updated
- [ ] Examples tested
- [ ] Documentation clear

## Notes
- Swagger UI verification is manual and important
- Load tests validate performance targets
- Coverage must be ≥80% per constitution
- README should be comprehensive for new developers
```

---

## Batch 10: Phase 8 - Deployment

**Issues**: #113-118 (6 tasks)
**Duration**: 1-2 days
**Depends On**: Batch 9 ✅
**Priority**: P3

### Agent Brief: Batch 10

```markdown
# Agent Brief: Phase 8 - Deployment

## Overview
Create Docker image, Kubernetes manifests, and deployment documentation.

## Prerequisites
✅ Batch 9 (Documentation) merged to main
✅ All code complete and tested

## Linked Issues
- [#119] T080: Create multi-stage Dockerfile
- [#120] T081: Create K8s deployment manifest
- [#121] T082: Create K8s service manifest
- [#122] T083: Create K8s ConfigMap
- [#123] T084: Test Docker build and run
- [#124] T085: Create deployment documentation

## Tasks

1. **T080**: Multi-stage Dockerfile
   - Stage 1: Go 1.24 build image
   - Build binary
   - Stage 2: Alpine or distroless runtime
   - COPY binary
   - Expose port 8080
   - Target size: <50MB

2. **T081**: K8s Deployment
   - deployments/k8s/deployment.yaml
   - Image: homelab-api:latest
   - Replicas: 2
   - Resources: 100MB memory, 200m CPU
   - Liveness probe: GET /health (10s)
   - Readiness probe: GET /health (5s)

3. **T082**: K8s Service
   - deployments/k8s/service.yaml
   - Type: ClusterIP
   - Port 80 → 8080
   - Selector: app=homelab-api

4. **T083**: K8s ConfigMap
   - deployments/k8s/configmap.yaml
   - LOG_LEVEL: info
   - RATE_LIMIT: 500
   - CORS_ORIGINS: http://localhost:3000

5. **T084**: Test Docker locally
   - Build: `docker build -t homelab-api:latest .`
   - Run: `docker run -p 8080:8080 homelab-api:latest`
   - Verify endpoints accessible

6. **T085**: Deployment documentation
   - Create deployments/README.md
   - Docker build instructions
   - K8s deployment steps
   - Environment variables
   - Troubleshooting

## Acceptance Criteria
- [ ] Dockerfile builds successfully
- [ ] Image size <50MB
- [ ] Local Docker test successful
- [ ] K8s manifests valid (kubectl apply --dry-run)
- [ ] ConfigMap includes required variables
- [ ] Probes configured for K8s
- [ ] Documentation complete
- [ ] Deployment tested locally

## Testing
```bash
# Docker
docker build -t homelab-api:latest .
docker run -p 8080:8080 homelab-api:latest
curl http://localhost:8080/health

# K8s validation
kubectl apply -f deployments/k8s/ --dry-run=client
```

## Git Workflow
```bash
git checkout -b batch-10-phase-8-deployment
# Create Dockerfile
# Create K8s manifests
# Create ConfigMap
# Test locally
git commit -m "batch-10: Phase 8 deployment"
git push origin batch-10-phase-8-deployment
# Create PR
```

## Review Checklist
- [ ] Dockerfile optimized
- [ ] Image builds and runs
- [ ] K8s manifests valid
- [ ] Probes configured
- [ ] ConfigMap complete
- [ ] Documentation thorough
- [ ] All tested locally

## Notes
- Multi-stage build keeps image size small
- Distroless runtime preferred for security
- K8s probes are important for cluster health
- ConfigMap makes configuration flexible
```

---

## Batch 11: Phase 9 - Final Validation

**Issues**: #125-129 (5 tasks)
**Duration**: 1 day
**Depends On**: Batch 10 ✅
**Priority**: P3

### Agent Brief: Batch 11

```markdown
# Agent Brief: Phase 9 - Final Validation & Compliance

## Overview
End-to-end validation, performance profiling, and constitution compliance verification.

## Prerequisites
✅ Batch 10 (Deployment) merged to main
✅ All code complete, tested, and deployed

## Linked Issues
- [#130] T090: Run full integration test suite
- [#131] T091: Performance profiling with pprof
- [#132] T092: Update README.md
- [#133] T093: Create quickstart guide
- [#134] T094: Final constitution compliance check

## Tasks

1. **T090**: Integration Test Suite
   - Execute all integration tests
   - Verify all user stories functional
   - Document any issues
   - Generate test report

2. **T091**: Performance Profiling
   - Run load test with CPU profiling
   - Run memory profiling
   - Analyze with pprof
   - Document bottlenecks
   - Create performance.md

3. **T092**: Update README.md
   - Project overview
   - Architecture diagram
   - API endpoint documentation
   - Development setup
   - Deployment instructions
   - Performance notes

4. **T093**: Create quickstart guide
   - .specify/features/001-homelab-api-service/quickstart.md
   - Step-by-step local setup
   - Example API calls
   - Docker quickstart
   - K8s quickstart

5. **T094**: Constitution Compliance
   - Verify Go 1.24+ standards
   - Verify 80%+ test coverage
   - Verify error handling
   - Verify structured logging
   - Verify graceful shutdown
   - Verify code formatting
   - Document compliance

## Acceptance Criteria
- [ ] All integration tests passing
- [ ] Performance profiling complete
- [ ] Bottlenecks identified and documented
- [ ] README comprehensive
- [ ] Quickstart guide clear
- [ ] Constitution requirements met
- [ ] Go 1.24 standards followed
- [ ] Coverage ≥80%
- [ ] All errors handled
- [ ] Structured logging implemented
- [ ] Graceful shutdown working

## Testing
```bash
make test              # Integration tests
make test-coverage    # Coverage report
go test ./tests/load -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof
```

## Git Workflow
```bash
git checkout -b batch-11-phase-9-final-validation
# Run all tests
# Profile performance
# Update documentation
# Verify compliance
git commit -m "batch-11: Phase 9 final validation and compliance"
git push origin batch-11-phase-9-final-validation
# Create PR for final review
```

## Review Checklist
- [ ] All tests passing
- [ ] Performance metrics documented
- [ ] README comprehensive
- [ ] Quickstart complete
- [ ] Constitution requirements met
- [ ] Go standards verified
- [ ] Coverage ≥80%
- [ ] Ready for production

## Notes
- This is the final phase - project completion validation
- Performance profiling will show optimization effectiveness
- Constitution compliance is non-negotiable
- Documentation should be clear for future contributors
```

---

# Dispatch Summary

## Total: 11 Batches, 97 Tasks

| Batch | Phase | Tasks | Duration | Status |
|-------|-------|-------|----------|--------|
| 1 | 0 | T001-T004 (4) | 1-2d | Research |
| 2 | 1 | T010-T019 (10) | 1-2d | Foundation |
| 3 | 1.5 | T019a-T019d (4) | 1d | Swagger |
| 4 | 2 | T020-T026 (7) | 2-3d | US1: Devices |
| 5 | 3 | T030-T035 (7) | 2-3d | US2: Health |
| 6 | 4 | T040-T044 (5) | 1-2d | Performance |
| 7 | 5 | T050-T055 (6) | 2-3d | US3: Control |
| 8 | 6 | T060-T065 (6) | 1-2d | US4: Cluster |
| 9 | 7 | T070-T073 (4) | 1d | Testing |
| 10 | 8 | T080-T085 (6) | 1-2d | Deployment |
| 11 | 9 | T090-T094 (5) | 1d | Validation |

**Total Duration**: 14-25 days (with parallel execution: 7-13 days)

## Key Dependencies

```
Batch 1 (Research) →
  Batch 2 (Foundation) →
    Batch 3 (Swagger) →
      Batches 4,5 (US1,US2 parallel) →
        Batch 6 (Performance) →
          Batches 7,8 (US3,US4 parallel) →
            Batch 9 (Testing) →
              Batch 10 (Deployment) →
                Batch 11 (Validation)
```

## Parallel Execution Opportunities

- **After Batch 3**: Dispatch Batches 4 and 5 simultaneously
- **After Batch 6**: Dispatch Batches 7 and 8 simultaneously
- **Can start after Batch 4**: Batch 9 (testing) could start early

With 3-4 agents working in parallel: **7-10 days total**

