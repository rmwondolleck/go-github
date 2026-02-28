# Tasks: Home Lab API Service

**Branch**: `001-homelab-api-service` | **Date**: 2026-02-28
**Input**: Design documents from `.specify/features/001-homelab-api-service/`
**Prerequisites**: spec.md ✅, plan.md ✅

**Tests**: Tests are MANDATORY per constitution (80%+ coverage). TDD workflow required.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

Single Go project at repository root:
- `cmd/api/` - Entry points
- `internal/` - Internal packages
- `tests/integration/` - Integration tests
- `api/` - OpenAPI specs

---

## Phase 0: Research & Validation

**Purpose**: Validate technology choices and performance characteristics

- [ ] T001 [P] Benchmark Gin framework basic routing (target: <10ms overhead)
  - Create `research/gin_benchmark_test.go`
  - Compare against stdlib net/http
  - Document findings in `.specify/features/001-homelab-api-service/research.md`

- [ ] T002 [P] Benchmark jsoniter vs stdlib encoding/json
  - Create `research/json_benchmark_test.go`
  - Test Device struct encoding (50 devices)
  - Measure allocation overhead
  - Document findings in research.md

- [ ] T003 [P] Benchmark sync.Map vs RWMutex for device storage
  - Create `research/storage_benchmark_test.go`
  - Test concurrent reads (100 goroutines)
  - Test mixed read/write workload
  - Document findings in research.md

- [ ] T004 Complete research.md with all benchmark results and recommendations

**Checkpoint**: Research complete - performance targets validated, ready for design phase

---

## Phase 1: Foundation (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T010 Initialize Go module and dependencies
  - Update `go.mod` to Go 1.24
  - Add gin-gonic/gin v1.10.0
  - Add swaggo/swag v1.16.3
  - Add swaggo/gin-swagger v1.6.0
  - Add swaggo/files v1.0.1
  - Add json-iterator/go v1.1.12
  - Add stretchr/testify v1.9.0
  - Run `go mod tidy`

- [ ] T011 Create project directory structure
  - `cmd/api/`
  - `internal/{homeassistant,cluster,health,middleware,handlers,models,server}/`
  - `tests/integration/`
  - `tests/load/`
  - `api/`
  - `deployments/{k8s}/`

- [ ] T012 [P] Define core data models in `internal/models/`
  - Create `internal/models/device.go` with Device struct
  - Create `internal/models/error.go` with ErrorResponse struct
  - Create `internal/models/health.go` with HealthStatus struct
  - Add JSON tags for all exported fields
  - Write data-model.md documentation

- [ ] T013 [P] Create service interfaces in `internal/contracts/`
  - Create directory `internal/homeassistant/` (interfaces defined in service files)
  - Document contracts in `.specify/features/001-homelab-api-service/contracts/`

- [ ] T014 [P] Implement request ID middleware
  - Create `internal/middleware/requestid.go`
  - Generate UUID for each request
  - Add to context and response headers
  - Unit test: `internal/middleware/requestid_test.go`

- [ ] T015 [P] Implement logging middleware with slog
  - Create `internal/middleware/logging.go`
  - Log request start/end with duration
  - Extract request ID from context
  - Include method, path, status, duration
  - Unit test: `internal/middleware/logging_test.go`

- [ ] T016 [P] Implement panic recovery middleware
  - Create `internal/middleware/recovery.go`
  - Catch panics, log with stack trace
  - Return 500 error response
  - Unit test: `internal/middleware/recovery_test.go`

- [ ] T017 Setup Gin server with graceful shutdown
  - Create `internal/server/server.go`
  - Initialize Gin router with middleware
  - Create `internal/server/shutdown.go` for SIGTERM/SIGINT handling
  - Create `cmd/api/main.go` entry point
  - Unit test: `internal/server/server_test.go`

- [ ] T018 [P] Create response helper utilities
  - Create `internal/handlers/response.go`
  - Functions: JSONSuccess, JSONError, JSONCreated
  - Consistent error response structure
  - Unit test: `internal/handlers/response_test.go`

- [ ] T019 Create Makefile with build targets
  - Targets: build, test, lint, run, clean, swagger
  - Add coverage report generation
  - Add Docker build target (placeholder for Phase 4)
  - Add swagger target: `swag init -g cmd/api/main.go -o api`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 1.5: API Documentation Setup (Blocking for Handlers)

**Purpose**: Setup Swagger/OpenAPI documentation generation and hosting

**⚠️ CRITICAL**: Must complete before handler implementation to enable inline documentation

- [ ] T019a [P] Install swag CLI tool
  - Run `go install github.com/swaggo/swag/cmd/swag@latest`
  - Verify installation: `swag --version`

- [ ] T019b [P] Add Swagger general API info annotations
  - Update `cmd/api/main.go` with Swagger metadata comments
  - Add: title, version, description, host, basePath
  - Example:
    ```go
    // @title Home Lab API Service
    // @version 1.0
    // @description REST API for home lab monitoring and control
    // @host localhost:8080
    // @BasePath /
    ```

- [ ] T019c Setup Swagger UI routes
  - Update `internal/server/server.go`
  - Add route: GET /api/docs/*any → ginSwagger.WrapHandler(swaggerFiles.Handler)
  - Serve OpenAPI spec at /api/docs/swagger.json
  - Serve Swagger UI at /api/docs/index.html

- [ ] T019d Generate initial Swagger docs
  - Run `make swagger` or `swag init -g cmd/api/main.go -o api`
  - Verify `api/docs.go` and `api/swagger.json` generated
  - Commit generated files to repository

---

## Phase 2: User Story 1 - Query HomeAssistant Device Status (Priority: P1) 🎯 MVP

**Goal**: Enable REST API queries for HomeAssistant device status with mocked data

**Independent Test**: `curl http://localhost:8080/api/v1/homeassistant/devices` returns JSON array of devices

### Tests for User Story 1 (TDD - Write FIRST)

- [ ] T020 [P] [US1] Write unit tests for HomeAssistant service
  - Create `internal/homeassistant/service_test.go`
  - Test: ListDevices returns all devices
  - Test: GetDevice returns specific device by ID
  - Test: GetDevice returns error for invalid ID
  - **Ensure tests FAIL before implementation**

- [ ] T021 [P] [US1] Write integration tests for device endpoints
  - Create `tests/integration/devices_test.go`
  - Test: GET /api/v1/homeassistant/devices returns 200 with device array
  - Test: GET /api/v1/homeassistant/devices/{id} returns 200 with device details
  - Test: GET /api/v1/homeassistant/devices/invalid returns 404
  - **Ensure tests FAIL before implementation**

### Implementation for User Story 1

- [ ] T022 [US1] Create mock device data
  - Create `internal/homeassistant/mock_data.go`
  - Define mockDevices slice with 50 sample devices
  - Include lights, sensors, switches, binary_sensors
  - Use realistic entity IDs (light.living_room, sensor.temperature, etc.)

- [ ] T023 [US1] Implement HomeAssistant service with in-memory storage
  - Create `internal/homeassistant/service.go`
  - Initialize sync.Map with mock devices
  - Implement ListDevices(ctx context.Context) ([]Device, error)
  - Implement GetDevice(ctx context.Context, id string) (*Device, error)
  - **Run tests from T020 - should now PASS**

- [ ] T024 [US1] Implement device list handler with Swagger docs
  - Create `internal/handlers/homeassistant.go`
  - Implement ListDevicesHandler(c *gin.Context)
  - Support optional ?type= query parameter for filtering
  - Add Swagger annotations:
    ```go
    // @Summary List all HomeAssistant devices
    // @Description Get list of all devices with optional type filtering
    // @Tags homeassistant
    // @Accept json
    // @Produce json
    // @Param type query string false "Filter by device type (light, sensor, switch, binary_sensor)"
    // @Success 200 {object} DeviceListResponse
    // @Failure 500 {object} ErrorResponse
    // @Router /api/v1/homeassistant/devices [get]
    ```
  - Use response helpers from T018
  - Run `make swagger` to regenerate docs
  - **Run integration tests from T021 - should now PASS**

- [ ] T025 [US1] Implement device detail handler with Swagger docs
  - Add GetDeviceHandler(c *gin.Context) to `internal/handlers/homeassistant.go`
  - Extract device ID from path parameter
  - Return 404 for non-existent devices
  - Add Swagger annotations:
    ```go
    // @Summary Get device by ID
    // @Description Get detailed information about a specific device
    // @Tags homeassistant
    // @Accept json
    // @Produce json
    // @Param id path string true "Device ID (e.g., light.living_room)"
    // @Success 200 {object} Device
    // @Failure 404 {object} ErrorResponse
    // @Failure 500 {object} ErrorResponse
    // @Router /api/v1/homeassistant/devices/{id} [get]
    ```
  - Run `make swagger` to regenerate docs
  - **Run integration tests from T021 - should now PASS**

- [ ] T026 [US1] Register device endpoints in server router
  - Update `internal/server/server.go`
  - Add route: GET /api/v1/homeassistant/devices
  - Add route: GET /api/v1/homeassistant/devices/:id
  - Apply all middleware

**Checkpoint**: User Story 1 complete - device query functionality working end-to-end

---

## Phase 3: User Story 2 - Health Check and Service Discovery (Priority: P1) 🎯 MVP

**Goal**: Provide K8s-ready health endpoint and service discovery API

**Independent Test**: `curl http://localhost:8080/health` returns 200 with service status

### Tests for User Story 2 (TDD - Write FIRST)

- [ ] T030 [P] [US2] Write unit tests for health checker
  - Create `internal/health/checker_test.go`
  - Test: Check returns healthy status when all components OK
  - Test: Check includes uptime in response
  - **Ensure tests FAIL before implementation**

- [ ] T031 [P] [US2] Write integration tests for health endpoint
  - Create `tests/integration/health_test.go`
  - Test: GET /health returns 200 with JSON health status
  - Test: Response includes status, components, uptime
  - Test: Response time <50ms (performance validation)
  - **Ensure tests FAIL before implementation**

### Implementation for User Story 2

- [ ] T032 [US2] Implement health checker service
  - Create `internal/health/checker.go`
  - Track server start time for uptime calculation
  - Implement Check(ctx context.Context) HealthStatus
  - Check components: api_server (always healthy for POC)
  - **Run tests from T030 - should now PASS**

- [ ] T033 [US2] Implement health endpoint handler with Swagger docs
  - Create `internal/handlers/health.go`
  - Implement HealthHandler(c *gin.Context)
  - Call health checker service
  - Optimize: consider pre-marshaling JSON for <10ms response
  - Add Swagger annotations:
    ```go
    // @Summary Health check
    // @Description Get service health status and uptime
    // @Tags health
    // @Produce json
    // @Success 200 {object} HealthStatus
    // @Failure 503 {object} ErrorResponse
    // @Router /health [get]
    ```
  - Unit test: `internal/handlers/health_test.go`
  - Run `make swagger` to regenerate docs
  - **Run integration tests from T031 - should now PASS**

- [ ] T034 [US2] Implement services discovery endpoint with Swagger docs
  - Create `internal/handlers/services.go`
  - Implement ListServicesHandler(c *gin.Context)
  - Return static list: [{name: "homeassistant", status: "available", endpoints: [...]}]
  - Add Swagger annotations:
    ```go
    // @Summary List available services
    // @Description Get list of all integrated services
    // @Tags services
    // @Produce json
    // @Success 200 {object} ServiceListResponse
    // @Router /api/v1/services [get]
    ```
  - Unit test: `internal/handlers/services_test.go`
  - Run `make swagger` to regenerate docs

- [ ] T035 [US2] Register health and services endpoints
  - Update `internal/server/server.go`
  - Add route: GET /health (no /api/v1 prefix)
  - Add route: GET /api/v1/services

**Checkpoint**: User Story 2 complete - health checks and service discovery operational

---

## Phase 4: Performance Optimization & Middleware (Priority: P2)

**Purpose**: Add rate limiting, CORS, and performance optimizations

- [ ] T040 [P] Implement rate limiting middleware
  - Create `internal/middleware/ratelimit.go`
  - Use token bucket algorithm (500 req/min per IP)
  - Return 429 Too Many Requests when exceeded
  - Add X-RateLimit-* headers
  - Unit test: `internal/middleware/ratelimit_test.go`

- [ ] T041 [P] Implement CORS middleware
  - Create `internal/middleware/cors.go`
  - Allow origin: http://localhost:3000 (configurable via env)
  - Allow methods: GET, POST, OPTIONS
  - Allow headers: Content-Type, Authorization
  - Unit test: `internal/middleware/cors_test.go`

- [ ] T042 [P] Optimize JSON encoding with jsoniter
  - Update response helpers to use jsoniter
  - Benchmark encoding performance improvement
  - Document results in performance.md

- [ ] T043 [P] Implement response pooling for device list
  - Add sync.Pool in `internal/handlers/homeassistant.go`
  - Pre-allocate DeviceListResponse structs
  - Benchmark allocation reduction
  - Document results in performance.md

- [ ] T044 Update server router with new middleware
  - Apply rate limiting to all /api/v1/* routes
  - Apply CORS to all routes
  - Update middleware order per plan.md

**Checkpoint**: Performance optimizations complete - rate limiting and CORS active

---

## Phase 5: User Story 3 - Control HomeAssistant Devices (Priority: P2)

**Goal**: Enable device control via REST API with mocked command execution

**Independent Test**: `curl -X POST http://localhost:8080/api/v1/homeassistant/devices/light.living_room/command -d '{"action":"turn_on"}'` returns 200

### Tests for User Story 3 (TDD - Write FIRST)

- [ ] T050 [P] [US3] Write unit tests for command execution
  - Update `internal/homeassistant/service_test.go`
  - Test: ExecuteCommand succeeds for valid device and action
  - Test: ExecuteCommand fails for invalid device ID
  - Test: ExecuteCommand fails for read-only device (sensor)
  - Test: ExecuteCommand fails for invalid action
  - **Ensure tests FAIL before implementation**

- [ ] T051 [P] [US3] Write integration tests for command endpoint
  - Update `tests/integration/devices_test.go`
  - Test: POST /api/v1/homeassistant/devices/{id}/command with valid command returns 200
  - Test: POST with invalid action returns 400
  - Test: POST to sensor device returns 405
  - **Ensure tests FAIL before implementation**

### Implementation for User Story 3

- [ ] T052 [US3] Define Command model
  - Create `internal/homeassistant/types.go`
  - Define Command struct with Action and Parameters fields
  - Add JSON tags and validation

- [ ] T053 [US3] Implement command execution in service
  - Update `internal/homeassistant/service.go`
  - Add ExecuteCommand(ctx context.Context, id string, cmd Command) error
  - Validate device exists and is controllable
  - Validate action is supported for device type
  - Mock command execution (just return nil for POC)
  - **Run tests from T050 - should now PASS**

- [ ] T054 [US3] Implement command endpoint handler with Swagger docs
  - Update `internal/handlers/homeassistant.go`
  - Add ExecuteCommandHandler(c *gin.Context)
  - Parse command from request body
  - Validate JSON payload
  - Return appropriate error codes (400, 404, 405)
  - Add Swagger annotations:
    ```go
    // @Summary Execute device command
    // @Description Send command to control a device
    // @Tags homeassistant
    // @Accept json
    // @Produce json
    // @Param id path string true "Device ID"
    // @Param command body Command true "Command to execute"
    // @Success 200 {object} CommandResponse
    // @Failure 400 {object} ErrorResponse
    // @Failure 404 {object} ErrorResponse
    // @Failure 405 {object} ErrorResponse
    // @Router /api/v1/homeassistant/devices/{id}/command [post]
    ```
  - Run `make swagger` to regenerate docs
  - **Run integration tests from T051 - should now PASS**

- [ ] T055 [US3] Register command endpoint
  - Update `internal/server/server.go`
  - Add route: POST /api/v1/homeassistant/devices/:id/command

**Checkpoint**: User Story 3 complete - device control functionality working

---

## Phase 6: User Story 4 - Query Cluster Services Info (Priority: P3)

**Goal**: Provide general K8s cluster service information

**Independent Test**: `curl http://localhost:8080/api/v1/cluster/services` returns JSON with service list

### Tests for User Story 4 (TDD - Write FIRST)

- [ ] T060 [P] [US4] Write unit tests for cluster service
  - Create `internal/cluster/service_test.go`
  - Test: ListServices returns mocked service list
  - Test: Optional filtering by name works
  - **Ensure tests FAIL before implementation**

- [ ] T061 [P] [US4] Write integration tests for cluster endpoint
  - Create `tests/integration/cluster_test.go`
  - Test: GET /api/v1/cluster/services returns 200 with service array
  - Test: GET with ?name=homeassistant filters results
  - **Ensure tests FAIL before implementation**

### Implementation for User Story 4

- [ ] T062 [US4] Create cluster service models
  - Create `internal/cluster/types.go`
  - Define ServiceInfo struct (Name, Namespace, Status, Endpoints)

- [ ] T063 [US4] Implement cluster service with mock data
  - Create `internal/cluster/service.go`
  - Implement ListServices(ctx context.Context, nameFilter string) ([]ServiceInfo, error)
  - Return mocked K8s services (homeassistant, prometheus, grafana, etc.)
  - **Run tests from T060 - should now PASS**

- [ ] T064 [US4] Implement cluster services handler with Swagger docs
  - Create `internal/handlers/cluster.go`
  - Implement ListClusterServicesHandler(c *gin.Context)
  - Support optional ?name= query parameter
  - Add Swagger annotations:
    ```go
    // @Summary List cluster services
    // @Description Get list of K8s cluster services with optional name filtering
    // @Tags cluster
    // @Produce json
    // @Param name query string false "Filter by service name"
    // @Success 200 {object} ServiceListResponse
    // @Router /api/v1/cluster/services [get]
    ```
  - Unit test: `internal/handlers/cluster_test.go`
  - Run `make swagger` to regenerate docs
  - **Run integration tests from T061 - should now PASS**

- [ ] T065 [US4] Register cluster endpoint
  - Update `internal/server/server.go`
  - Add route: GET /api/v1/cluster/services

**Checkpoint**: User Story 4 complete - cluster service discovery working

---

## Phase 7: API Documentation & Testing (Priority: P2)

**Purpose**: Complete API documentation and load testing

- [ ] T070 [P] Verify Swagger UI is accessible and functional
  - Start server: `make run`
  - Access Swagger UI: http://localhost:8080/api/docs/index.html
  - Verify all endpoints are documented
  - Test "Try it out" functionality for each endpoint
  - Verify OpenAPI spec downloadable at /api/docs/swagger.json
  - Document Swagger UI URL in README.md

- [ ] T071 [P] Write load tests for concurrent requests
  - Create `tests/load/concurrent_test.go`
  - Test 100 concurrent requests to health endpoint
  - Test 100 concurrent requests to device list endpoint
  - Validate response times <200ms p99
  - Detect memory leaks

- [ ] T072 Run final coverage check
  - Execute `make test` with coverage report
  - Validate >80% coverage achieved
  - Document coverage in README.md

- [ ] T073 [P] Create API usage examples in README
  - Document local development setup
  - Add curl examples for all endpoints
  - Add environment variables documentation

**Checkpoint**: API fully documented and load tested

---

## Phase 8: Deployment (Priority: P3)

**Purpose**: Containerization and K8s deployment

- [ ] T080 Create multi-stage Dockerfile
  - Create `deployments/Dockerfile`
  - Stage 1: Build binary with Go 1.24
  - Stage 2: Runtime with alpine/distroless
  - Optimize for <50MB image size

- [ ] T081 [P] Create K8s deployment manifest
  - Create `deployments/k8s/deployment.yaml`
  - Set resource limits: 100MB memory, 200m CPU
  - Configure liveness probe: /health
  - Configure readiness probe: /health
  - Set replica count: 2

- [ ] T082 [P] Create K8s service manifest
  - Create `deployments/k8s/service.yaml`
  - Type: ClusterIP
  - Port: 80 → targetPort: 8080
  - Selector: app=homelab-api

- [ ] T083 [P] Create K8s ConfigMap for configuration
  - Create `deployments/k8s/configmap.yaml`
  - Environment variables: LOG_LEVEL, RATE_LIMIT, CORS_ORIGINS
  - Mock device data configuration

- [ ] T084 Test Docker build and run locally
  - Build image: `docker build -t homelab-api:latest .`
  - Run container: `docker run -p 8080:8080 homelab-api:latest`
  - Verify all endpoints accessible

- [ ] T085 Create deployment documentation
  - Create `deployments/README.md`
  - Document Docker build process
  - Document K8s deployment steps
  - Document environment variables
  - Add troubleshooting section

**Checkpoint**: Deployment complete - service ready for K8s cluster

---

## Phase 9: Final Validation & Documentation (Priority: P3)

**Purpose**: End-to-end validation and documentation completion

- [ ] T090 Run full integration test suite
  - Execute all integration tests
  - Validate all user stories functional
  - Document any issues

- [ ] T091 Performance profiling with pprof
  - Run CPU profiling during load test
  - Run memory profiling
  - Identify and document any bottlenecks
  - Create performance.md with results

- [ ] T092 [P] Update project README.md
  - Add project overview
  - Add architecture diagram (ASCII or link to draw.io)
  - Add API endpoint documentation
  - Add development setup instructions
  - Add deployment instructions

- [ ] T093 [P] Create quickstart guide
  - Create `.specify/features/001-homelab-api-service/quickstart.md`
  - Step-by-step local development setup
  - Example API calls with expected responses
  - Docker and K8s quickstart

- [ ] T094 Final constitution compliance check
  - Verify all Go 1.24 standards followed
  - Verify 80%+ test coverage
  - Verify all errors handled
  - Verify structured logging implemented
  - Verify graceful shutdown works
  - Document compliance in README

**Checkpoint**: Project complete and ready for production use

---

## Summary

**Total Tasks**: 97 tasks organized across 9 phases (includes Phase 1.5 for Swagger setup)

**Priority Breakdown**:
- **P1 (MVP)**: Phases 1, 1.5, 2, 3 - Foundation + Swagger + US1 + US2 (41 tasks)
- **P2**: Phases 4, 5, 7 - Performance + US3 + Documentation (24 tasks)
- **P3**: Phases 6, 8, 9 - US4 + Deployment + Validation (24 tasks)
- **P0 (Research)**: Phase 0 - Technology validation (4 tasks)

**Critical Path** (blocking dependencies):
1. Phase 0: Research (validates technology choices)
2. Phase 1: Foundation (T010-T019 must complete first)
3. Phase 1.5: Swagger setup (T019a-T019d, enables inline API documentation)
4. Phases 2-3: US1 & US2 (MVP - highest priority)
5. Phases 4-9: Can be executed in parallel or sequentially

**Swagger/OpenAPI Integration**:
- Inline annotations on all handler functions
- Auto-generated OpenAPI spec via `swag` tool
- Swagger UI hosted at `/api/docs/index.html`
- OpenAPI spec available at `/api/docs/swagger.json`
- Regenerate docs with `make swagger` after handler changes

**Test-Driven Development**:
- All test tasks explicitly marked "Write FIRST"
- Tests must FAIL before implementation
- Tests must PASS after implementation
- 80%+ coverage validated in T072

**Independent User Stories**:
Each user story (US1-US4) can be developed and tested independently after Phase 1 foundation is complete.

---

## Execution Workflow

1. **Start**: Execute Phase 0 research (4 tasks, 1-2 days)
2. **Foundation**: Execute Phase 1 serially (10 tasks, 1-2 days)
3. **MVP Iteration 1**: Execute Phase 2 (US1) + Phase 3 (US2) (17 tasks, 2-3 days)
4. **MVP Iteration 2**: Execute Phase 4 (performance) + Phase 5 (US3) (10 tasks, 1-2 days)
5. **Extended Features**: Execute Phase 6 (US4) + Phase 7 (docs/testing) (11 tasks, 1-2 days)
6. **Deployment**: Execute Phase 8 + Phase 9 (18 tasks, 1-2 days)

**Estimated Total Time**: 7-13 days for complete implementation

**Minimum Viable Product** (Phase 0-3): 4-7 days for fully functional API with device queries and health checks











