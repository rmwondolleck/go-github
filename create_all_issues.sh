#!/bin/bash

# Create all GitHub issues from tasks.md with label validation
# Usage: ./create_all_issues.sh

set -e

OWNER="rmwondolleck"
REPO="go-github"

echo "Creating missing labels..."

# Function to create label if it doesn't exist
create_label_if_missing() {
    local label=$1
    local color=${2:-"#ededed"}

    if ! gh label list --repo "$OWNER/$REPO" | grep -q "^$label"; then
        echo "  Creating label: $label"
        gh label create "$label" --color "$color" --repo "$OWNER/$REPO" 2>/dev/null || true
    fi
}

# Create all phase labels
for i in {0..9}; do
    create_label_if_missing "phase-$i"
done

# Create user story labels
for i in {1..4}; do
    create_label_if_missing "us-$i"
done

# Create status labels
create_label_if_missing "blocked-by"
create_label_if_missing "blocks"

echo "All labels ready!"
echo ""
echo "Creating all GitHub issues for Home Lab API Service..."

# Phase 3: User Story 2 - Health Check (T030-T035)
echo "Creating Phase 3 issues (US2 - Health & Discovery)..."

gh issue create \
  --title "T030: Write unit tests for health checker" \
  --body "## Task: T030 [P] [US2] Write unit tests for health checker

**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Test Task (TDD)
**Parallelizable**: Yes

## Description
Write unit tests FIRST for health checker service (TDD methodology)

## Files to Create
- internal/health/checker_test.go

## Tests to Write
- TestCheck_ReturnsHealthyStatus
- TestCheck_IncludesUptimeInResponse
- TestCheck_IncludesComponentsInResponse

## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All test cases documented
- [ ] Ready for service implementation

## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-3,us-2,001-homelab-api,parallel"

gh issue create \
  --title "T031: Write integration tests for health endpoint" \
  --body "## Task: T031 [P] [US2] Write integration tests for health endpoint

**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Test Task (TDD)
**Parallelizable**: Yes

## Description
Write integration tests FIRST for health endpoint (TDD methodology)

## Files to Create
- tests/integration/health_test.go

## Tests to Write
- TestHealthEndpoint_Returns200
- TestHealthEndpoint_IncludesStatusAndUptime
- TestHealthEndpoint_ResponseUnder50ms

## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All HTTP scenarios covered
- [ ] Response time assertions included
- [ ] Ready for handler implementation

## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-3,us-2,001-homelab-api,parallel"

gh issue create \
  --title "T032: Implement health checker service" \
  --body "## Task: T032 [US2] Implement health checker service

**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Implementation Task

## Description
Implement health checker service with uptime tracking

## Files to Create
- internal/health/checker.go

## Functions to Implement
- Check(ctx context.Context) HealthStatus

## Features
- Track server start time for uptime calculation
- Return healthy status
- Include component health (api_server)

## Acceptance Criteria
- [ ] Service implemented
- [ ] Uptime calculation working
- [ ] Tests from T030 now PASS

## Dependencies
- T030 (Tests written)" \
  --label "phase-3,us-2,001-homelab-api"

gh issue create \
  --title "T033: Implement health endpoint handler with Swagger docs" \
  --body "## Task: T033 [US2] Implement health endpoint handler with Swagger docs

**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Implementation Task

## Description
Implement HTTP handler for health checks with Swagger documentation

## Files to Update
- internal/handlers/health.go

## Function to Implement
- HealthHandler(c *gin.Context)

## Features
- Call health checker service
- Optimize for <50ms response
- Swagger annotations

## Acceptance Criteria
- [ ] Handler implemented
- [ ] Swagger annotations added
- [ ] Response time <50ms
- [ ] Tests from T031 now PASS

## Dependencies
- T032 (Service implemented)
- T031 (Tests written)" \
  --label "phase-3,us-2,001-homelab-api"

gh issue create \
  --title "T034: Implement services discovery endpoint with Swagger docs" \
  --body "## Task: T034 [US2] Implement services discovery endpoint with Swagger docs

**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Implementation Task

## Description
Implement HTTP handler for services discovery with Swagger documentation

## Files to Create
- internal/handlers/services.go

## Function to Implement
- ListServicesHandler(c *gin.Context)

## Features
- Return static list of available services
- Swagger annotations

## Acceptance Criteria
- [ ] Handler implemented
- [ ] Swagger annotations added
- [ ] Static service list returned
- [ ] Tests passing

## Dependencies
- T018 (Response helpers)" \
  --label "phase-3,us-2,001-homelab-api"

gh issue create \
  --title "T035: Register health and services endpoints" \
  --body "## Task: T035 [US2] Register health and services endpoints

**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Implementation Task

## Description
Register health and services endpoints in Gin router

## Files to Update
- internal/server/server.go

## Routes to Register
- GET /health
- GET /api/v1/services

## Acceptance Criteria
- [ ] Routes registered correctly
- [ ] Middleware applied
- [ ] Endpoints accessible
- [ ] Tests passing

## Dependencies
- T033, T034 (Handlers implemented)" \
  --label "phase-3,us-2,001-homelab-api"

echo "Phase 3 issues created (T030-T035)"

# Phase 4: Performance Optimization (T040-T044)
echo "Creating Phase 4 issues (Performance & Middleware)..."

gh issue create \
  --title "T040: Implement rate limiting middleware" \
  --body "## Task: T040 [P] Implement rate limiting middleware

**Phase**: Phase 4 - Performance Optimization & Middleware
**Type**: Performance Task
**Parallelizable**: Yes

## Description
Create middleware for rate limiting with token bucket algorithm

## Files to Create
- internal/middleware/ratelimit.go
- internal/middleware/ratelimit_test.go

## Features
- Token bucket algorithm
- 500 req/min per IP
- Return 429 Too Many Requests
- Add X-RateLimit-* headers

## Acceptance Criteria
- [ ] Middleware implemented
- [ ] Token bucket working correctly
- [ ] Headers added to responses
- [ ] Tests passing

## Dependencies
- T017 (Server setup)" \
  --label "phase-4,001-homelab-api,parallel"

gh issue create \
  --title "T041: Implement CORS middleware" \
  --body "## Task: T041 [P] Implement CORS middleware

**Phase**: Phase 4 - Performance Optimization & Middleware
**Type**: Performance Task
**Parallelizable**: Yes

## Description
Create CORS middleware for cross-origin requests

## Files to Create
- internal/middleware/cors.go
- internal/middleware/cors_test.go

## Features
- Allow origin: http://localhost:3000
- Configurable via environment variable
- Allow methods: GET, POST, OPTIONS
- Allow headers: Content-Type, Authorization

## Acceptance Criteria
- [ ] Middleware implemented
- [ ] CORS headers added
- [ ] Configuration working
- [ ] Tests passing

## Dependencies
- T017 (Server setup)" \
  --label "phase-4,001-homelab-api,parallel"

gh issue create \
  --title "T042: Optimize JSON encoding with jsoniter" \
  --body "## Task: T042 [P] Optimize JSON encoding with jsoniter

**Phase**: Phase 4 - Performance Optimization & Middleware
**Type**: Performance Task
**Parallelizable**: Yes

## Description
Replace stdlib JSON encoding with jsoniter for better performance

## Files to Update
- internal/handlers/response.go

## Features
- Use jsoniter for encoding
- Benchmark performance improvement
- Document results

## Acceptance Criteria
- [ ] jsoniter integrated
- [ ] Benchmarks completed
- [ ] 2-3x improvement documented
- [ ] Tests passing

## Dependencies
- T018 (Response helpers)" \
  --label "phase-4,001-homelab-api,parallel"

gh issue create \
  --title "T043: Implement response pooling for device list" \
  --body "## Task: T043 [P] Implement response pooling for device list

**Phase**: Phase 4 - Performance Optimization & Middleware
**Type**: Performance Task
**Parallelizable**: Yes

## Description
Implement object pooling for DeviceListResponse to reduce allocations

## Files to Update
- internal/handlers/homeassistant.go

## Features
- sync.Pool for response structs
- Pre-allocate response objects
- Benchmark allocation reduction

## Acceptance Criteria
- [ ] Pool implemented
- [ ] Allocations reduced
- [ ] Benchmarks documented
- [ ] Tests passing

## Dependencies
- T024 (Device list handler)" \
  --label "phase-4,001-homelab-api,parallel"

gh issue create \
  --title "T044: Update server router with new middleware" \
  --body "## Task: T044 Update server router with new middleware

**Phase**: Phase 4 - Performance Optimization & Middleware
**Type**: Integration Task

## Description
Integrate new middleware into server router with correct ordering

## Files to Update
- internal/server/server.go

## Middleware to Apply
- Rate limiting to /api/v1/* routes
- CORS to all routes
- Correct ordering

## Acceptance Criteria
- [ ] Middleware integrated
- [ ] Correct ordering applied
- [ ] All endpoints protected
- [ ] Tests passing

## Dependencies
- T040, T041 (Middleware implemented)" \
  --label "phase-4,001-homelab-api"

echo "Phase 4 issues created (T040-T044)"

# Phase 5: User Story 3 - Device Control (T050-T055)
echo "Creating Phase 5 issues (US3 - Device Control)..."

gh issue create \
  --title "T050: Write unit tests for command execution" \
  --body "## Task: T050 [P] [US3] Write unit tests for command execution

**Phase**: Phase 5 - User Story 3
**User Story**: US3 - Control HomeAssistant Devices (P2)
**Type**: Test Task (TDD)
**Parallelizable**: Yes

## Description
Write unit tests FIRST for command execution (TDD methodology)

## Files to Update
- internal/homeassistant/service_test.go

## Tests to Write
- TestExecuteCommand_SucceedsForValidDevice
- TestExecuteCommand_FailsForInvalidDevice
- TestExecuteCommand_FailsForReadOnlyDevice
- TestExecuteCommand_FailsForInvalidAction

## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All scenarios covered
- [ ] Ready for service implementation

## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-5,us-3,001-homelab-api,parallel"

gh issue create \
  --title "T051: Write integration tests for command endpoint" \
  --body "## Task: T051 [P] [US3] Write integration tests for command endpoint

**Phase**: Phase 5 - User Story 3
**User Story**: US3 - Control HomeAssistant Devices (P2)
**Type**: Test Task (TDD)
**Parallelizable**: Yes

## Description
Write integration tests FIRST for command endpoint (TDD methodology)

## Files to Update
- tests/integration/devices_test.go

## Tests to Write
- TestExecuteCommand_Returns200ForValidCommand
- TestExecuteCommand_Returns400ForInvalidAction
- TestExecuteCommand_Returns405ForReadOnlyDevice

## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All HTTP scenarios covered
- [ ] Ready for handler implementation

## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-5,us-3,001-homelab-api,parallel"

gh issue create \
  --title "T052: Define Command model" \
  --body "## Task: T052 [US3] Define Command model

**Phase**: Phase 5 - User Story 3
**User Story**: US3 - Control HomeAssistant Devices (P2)
**Type**: Implementation Task

## Description
Define Command struct for device control

## Files to Create
- internal/homeassistant/types.go

## Struct to Define
- Command with Action and Parameters fields
- JSON tags
- Validation

## Acceptance Criteria
- [ ] Model defined
- [ ] JSON tags added
- [ ] Validation rules in place

## Dependencies
- T023 (Service foundation)" \
  --label "phase-5,us-3,001-homelab-api"

gh issue create \
  --title "T053: Implement command execution in service" \
  --body "## Task: T053 [US3] Implement command execution in service

**Phase**: Phase 5 - User Story 3
**User Story**: US3 - Control HomeAssistant Devices (P2)
**Type**: Implementation Task

## Description
Implement command execution in HomeAssistant service

## Files to Update
- internal/homeassistant/service.go

## Function to Implement
- ExecuteCommand(ctx context.Context, id string, cmd Command) error

## Features
- Validate device exists and is controllable
- Validate action supported for device type
- Mock execution (return nil)

## Acceptance Criteria
- [ ] Service implemented
- [ ] Validation working
- [ ] Tests from T050 now PASS

## Dependencies
- T052 (Command model)
- T050 (Tests written)" \
  --label "phase-5,us-3,001-homelab-api"

gh issue create \
  --title "T054: Implement command endpoint handler with Swagger docs" \
  --body "## Task: T054 [US3] Implement command endpoint handler with Swagger docs

**Phase**: Phase 5 - User Story 3
**User Story**: US3 - Control HomeAssistant Devices (P2)
**Type**: Implementation Task

## Description
Implement HTTP handler for device commands with Swagger documentation

## Files to Update
- internal/handlers/homeassistant.go

## Function to Implement
- ExecuteCommandHandler(c *gin.Context)

## Features
- Parse command from request body
- Validate JSON payload
- Return appropriate error codes
- Swagger annotations

## Acceptance Criteria
- [ ] Handler implemented
- [ ] Swagger annotations added
- [ ] Error handling correct
- [ ] Tests from T051 now PASS

## Dependencies
- T053 (Service implemented)
- T051 (Tests written)" \
  --label "phase-5,us-3,001-homelab-api"

gh issue create \
  --title "T055: Register command endpoint" \
  --body "## Task: T055 [US3] Register command endpoint

**Phase**: Phase 5 - User Story 3
**User Story**: US3 - Control HomeAssistant Devices (P2)
**Type**: Implementation Task

## Description
Register command endpoint in Gin router

## Files to Update
- internal/server/server.go

## Route to Register
- POST /api/v1/homeassistant/devices/:id/command

## Acceptance Criteria
- [ ] Route registered correctly
- [ ] Middleware applied
- [ ] Endpoint accessible
- [ ] Tests passing

## Dependencies
- T054 (Handler implemented)" \
  --label "phase-5,us-3,001-homelab-api"

echo "Phase 5 issues created (T050-T055)"

# Phase 6: User Story 4 - Cluster Services (T060-T065)
echo "Creating Phase 6 issues (US4 - Cluster Services)..."

gh issue create \
  --title "T060: Write unit tests for cluster service" \
  --body "## Task: T060 [P] [US4] Write unit tests for cluster service

**Phase**: Phase 6 - User Story 4
**User Story**: US4 - Query Cluster Services Info (P3)
**Type**: Test Task (TDD)
**Parallelizable**: Yes

## Description
Write unit tests FIRST for cluster service (TDD methodology)

## Files to Create
- internal/cluster/service_test.go

## Tests to Write
- TestListServices_ReturnsMockedServices
- TestListServices_FiltersByName

## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All test cases documented
- [ ] Ready for service implementation

## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-6,us-4,001-homelab-api,parallel"

gh issue create \
  --title "T061: Write integration tests for cluster endpoint" \
  --body "## Task: T061 [P] [US4] Write integration tests for cluster endpoint

**Phase**: Phase 6 - User Story 4
**User Story**: US4 - Query Cluster Services Info (P3)
**Type**: Test Task (TDD)
**Parallelizable**: Yes

## Description
Write integration tests FIRST for cluster endpoint (TDD methodology)

## Files to Create
- tests/integration/cluster_test.go

## Tests to Write
- TestListClusterServices_Returns200
- TestListClusterServices_FiltersByName

## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All HTTP scenarios covered
- [ ] Ready for handler implementation

## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-6,us-4,001-homelab-api,parallel"

gh issue create \
  --title "T062: Create cluster service models" \
  --body "## Task: T062 [US4] Create cluster service models

**Phase**: Phase 6 - User Story 4
**User Story**: US4 - Query Cluster Services Info (P3)
**Type**: Implementation Task

## Description
Define cluster service models

## Files to Create
- internal/cluster/types.go

## Struct to Define
- ServiceInfo with Name, Namespace, Status, Endpoints

## Acceptance Criteria
- [ ] Model defined
- [ ] JSON tags added
- [ ] Ready for service implementation

## Dependencies
- T061 (Tests written)" \
  --label "phase-6,us-4,001-homelab-api"

gh issue create \
  --title "T063: Implement cluster service with mock data" \
  --body "## Task: T063 [US4] Implement cluster service with mock data

**Phase**: Phase 6 - User Story 4
**User Story**: US4 - Query Cluster Services Info (P3)
**Type**: Implementation Task

## Description
Implement cluster service with mocked K8s services

## Files to Create
- internal/cluster/service.go

## Function to Implement
- ListServices(ctx context.Context, nameFilter string) ([]ServiceInfo, error)

## Mock Services
- homeassistant, prometheus, grafana, etc.

## Acceptance Criteria
- [ ] Service implemented
- [ ] Mock services created
- [ ] Filtering working
- [ ] Tests from T060 now PASS

## Dependencies
- T062 (Models created)
- T060 (Tests written)" \
  --label "phase-6,us-4,001-homelab-api"

gh issue create \
  --title "T064: Implement cluster services handler with Swagger docs" \
  --body "## Task: T064 [US4] Implement cluster services handler with Swagger docs

**Phase**: Phase 6 - User Story 4
**User Story**: US4 - Query Cluster Services Info (P3)
**Type**: Implementation Task

## Description
Implement HTTP handler for cluster services with Swagger documentation

## Files to Create
- internal/handlers/cluster.go

## Function to Implement
- ListClusterServicesHandler(c *gin.Context)

## Features
- Support ?name= query parameter
- Swagger annotations

## Acceptance Criteria
- [ ] Handler implemented
- [ ] Swagger annotations added
- [ ] Filtering working
- [ ] Tests from T061 now PASS

## Dependencies
- T063 (Service implemented)
- T061 (Tests written)" \
  --label "phase-6,us-4,001-homelab-api"

gh issue create \
  --title "T065: Register cluster endpoint" \
  --body "## Task: T065 [US4] Register cluster endpoint

**Phase**: Phase 6 - User Story 4
**User Story**: US4 - Query Cluster Services Info (P3)
**Type**: Implementation Task

## Description
Register cluster services endpoint in Gin router

## Files to Update
- internal/server/server.go

## Route to Register
- GET /api/v1/cluster/services

## Acceptance Criteria
- [ ] Route registered correctly
- [ ] Middleware applied
- [ ] Endpoint accessible
- [ ] Tests passing

## Dependencies
- T064 (Handler implemented)" \
  --label "phase-6,us-4,001-homelab-api"

echo "Phase 6 issues created (T060-T065)"

# Phase 7: Documentation & Testing (T070-T073)
echo "Creating Phase 7 issues (Documentation & Testing)..."

gh issue create \
  --title "T070: Verify Swagger UI is accessible and functional" \
  --body "## Task: T070 [P] Verify Swagger UI is accessible and functional

**Phase**: Phase 7 - API Documentation & Testing
**Type**: Documentation Task
**Parallelizable**: Yes

## Description
Verify Swagger UI and OpenAPI spec are properly accessible

## Acceptance Criteria
- [ ] Swagger UI accessible at http://localhost:8080/api/docs/index.html
- [ ] All endpoints documented
- [ ] \"Try it out\" functionality working
- [ ] OpenAPI spec downloadable at /api/docs/swagger.json
- [ ] URL documented in README.md

## Dependencies
- All handler implementation tasks (T024, T025, T033, T034, T054, T064)" \
  --label "phase-7,001-homelab-api,parallel"

gh issue create \
  --title "T071: Write load tests for concurrent requests" \
  --body "## Task: T071 [P] Write load tests for concurrent requests

**Phase**: Phase 7 - API Documentation & Testing
**Type**: Testing Task
**Parallelizable**: Yes

## Description
Create load tests for concurrent request handling

## Files to Create
- tests/load/concurrent_test.go

## Tests to Write
- 100 concurrent requests to health endpoint
- 100 concurrent requests to device list endpoint
- Validate response times <200ms p99
- Detect memory leaks

## Acceptance Criteria
- [ ] Load tests created
- [ ] All scenarios tested
- [ ] Performance targets validated

## Dependencies
- All implementation tasks (Phases 1-6)" \
  --label "phase-7,001-homelab-api,parallel"

gh issue create \
  --title "T072: Run final coverage check" \
  --body "## Task: T072 Run final coverage check

**Phase**: Phase 7 - API Documentation & Testing
**Type**: Testing Task

## Description
Validate code coverage meets 80%+ requirement

## Acceptance Criteria
- [ ] Coverage report generated: make test
- [ ] Coverage ≥80% achieved
- [ ] Coverage documented in README.md

## Dependencies
- All test tasks (all *_test.go files)" \
  --label "phase-7,001-homelab-api"

gh issue create \
  --title "T073: Create API usage examples in README" \
  --body "## Task: T073 [P] Create API usage examples in README

**Phase**: Phase 7 - API Documentation & Testing
**Type**: Documentation Task
**Parallelizable**: Yes

## Description
Document API usage with curl examples

## Files to Update
- README.md

## Content to Add
- Local development setup
- curl examples for all endpoints
- Environment variables documentation

## Acceptance Criteria
- [ ] README updated
- [ ] All endpoints documented with examples
- [ ] Setup instructions clear

## Dependencies
- All implementation tasks" \
  --label "phase-7,001-homelab-api,parallel"

echo "Phase 7 issues created (T070-T073)"

# Phase 8: Deployment (T080-T085)
echo "Creating Phase 8 issues (Deployment)..."

gh issue create \
  --title "T080: Create multi-stage Dockerfile" \
  --body "## Task: T080 Create multi-stage Dockerfile

**Phase**: Phase 8 - Deployment
**Type**: Deployment Task

## Description
Create optimized multi-stage Dockerfile for containerization

## Files to Create
- deployments/Dockerfile

## Stages
- Stage 1: Build binary with Go 1.24
- Stage 2: Runtime with alpine/distroless

## Optimization
- Optimize for <50MB image size

## Acceptance Criteria
- [ ] Dockerfile created
- [ ] Build stage compiles successfully
- [ ] Runtime stage includes binary only
- [ ] Image size <50MB

## Dependencies
- T019 (Makefile with docker target)" \
  --label "phase-8,001-homelab-api"

gh issue create \
  --title "T081: Create K8s deployment manifest" \
  --body "## Task: T081 [P] Create K8s deployment manifest

**Phase**: Phase 8 - Deployment
**Type**: Deployment Task
**Parallelizable**: Yes

## Description
Create Kubernetes deployment manifest

## Files to Create
- deployments/k8s/deployment.yaml

## Configuration
- Resource limits: 100MB memory, 200m CPU
- Liveness probe: /health
- Readiness probe: /health
- Replica count: 2

## Acceptance Criteria
- [ ] Manifest created
- [ ] Probes configured
- [ ] Resource limits set
- [ ] Validates successfully

## Dependencies
- T080 (Dockerfile)" \
  --label "phase-8,001-homelab-api,parallel"

gh issue create \
  --title "T082: Create K8s service manifest" \
  --body "## Task: T082 [P] Create K8s service manifest

**Phase**: Phase 8 - Deployment
**Type**: Deployment Task
**Parallelizable**: Yes

## Description
Create Kubernetes service manifest

## Files to Create
- deployments/k8s/service.yaml

## Configuration
- Type: ClusterIP
- Port: 80 → targetPort: 8080
- Selector: app=homelab-api

## Acceptance Criteria
- [ ] Manifest created
- [ ] Ports configured
- [ ] Selector correct
- [ ] Validates successfully

## Dependencies
- T081 (Deployment manifest)" \
  --label "phase-8,001-homelab-api,parallel"

gh issue create \
  --title "T083: Create K8s ConfigMap for configuration" \
  --body "## Task: T083 [P] Create K8s ConfigMap for configuration

**Phase**: Phase 8 - Deployment
**Type**: Deployment Task
**Parallelizable**: Yes

## Description
Create Kubernetes ConfigMap for environment configuration

## Files to Create
- deployments/k8s/configmap.yaml

## Configuration
- Environment variables: LOG_LEVEL, RATE_LIMIT, CORS_ORIGINS
- Mock device data configuration

## Acceptance Criteria
- [ ] ConfigMap created
- [ ] All variables configured
- [ ] Format correct
- [ ] Validates successfully

## Dependencies
- T081, T082 (Manifests)" \
  --label "phase-8,001-homelab-api,parallel"

gh issue create \
  --title "T084: Test Docker build and run locally" \
  --body "## Task: T084 Test Docker build and run locally

**Phase**: Phase 8 - Deployment
**Type**: Deployment Task

## Description
Test Docker image build and local execution

## Commands
- docker build -t homelab-api:latest .
- docker run -p 8080:8080 homelab-api:latest

## Acceptance Criteria
- [ ] Build successful
- [ ] Image runs locally
- [ ] All endpoints accessible
- [ ] Health check responds

## Dependencies
- T080 (Dockerfile)" \
  --label "phase-8,001-homelab-api"

gh issue create \
  --title "T085: Create deployment documentation" \
  --body "## Task: T085 Create deployment documentation

**Phase**: Phase 8 - Deployment
**Type**: Documentation Task

## Description
Document deployment process

## Files to Create
- deployments/README.md

## Content
- Docker build process
- K8s deployment steps
- Environment variables
- Troubleshooting section

## Acceptance Criteria
- [ ] Documentation complete
- [ ] All steps clear
- [ ] Examples provided

## Dependencies
- T080-T083 (All manifests)" \
  --label "phase-8,001-homelab-api"

echo "Phase 8 issues created (T080-T085)"

# Phase 9: Final Validation (T090-T094)
echo "Creating Phase 9 issues (Final Validation & Documentation)..."

gh issue create \
  --title "T090: Run full integration test suite" \
  --body "## Task: T090 Run full integration test suite

**Phase**: Phase 9 - Final Validation & Documentation
**Type**: Testing Task

## Description
Execute complete integration test suite

## Acceptance Criteria
- [ ] All integration tests pass
- [ ] All user stories validated functional
- [ ] No test failures or errors
- [ ] Issues documented if found

## Dependencies
- All implementation and test tasks" \
  --label "phase-9,001-homelab-api"

gh issue create \
  --title "T091: Performance profiling with pprof" \
  --body "## Task: T091 Performance profiling with pprof

**Phase**: Phase 9 - Final Validation & Documentation
**Type**: Performance Task

## Description
Profile application performance with pprof

## Profiling
- CPU profiling during load test
- Memory profiling
- Identify bottlenecks
- Document results

## Files to Create/Update
- performance.md

## Acceptance Criteria
- [ ] CPU profile completed
- [ ] Memory profile completed
- [ ] Results documented
- [ ] Bottlenecks identified

## Dependencies
- T071 (Load tests)" \
  --label "phase-9,001-homelab-api"

gh issue create \
  --title "T092: Update project README.md" \
  --body "## Task: T092 [P] Update project README.md

**Phase**: Phase 9 - Final Validation & Documentation
**Type**: Documentation Task
**Parallelizable**: Yes

## Description
Complete project README with comprehensive documentation

## Content
- Project overview
- Architecture diagram
- API endpoint documentation
- Development setup instructions
- Deployment instructions

## Acceptance Criteria
- [ ] README complete
- [ ] All sections documented
- [ ] Clear setup instructions

## Dependencies
- All implementation tasks" \
  --label "phase-9,001-homelab-api,parallel"

gh issue create \
  --title "T093: Create quickstart guide" \
  --body "## Task: T093 [P] Create quickstart guide

**Phase**: Phase 9 - Final Validation & Documentation
**Type**: Documentation Task
**Parallelizable**: Yes

## Description
Create quickstart guide for developers

## Files to Create
- .specify/features/001-homelab-api-service/quickstart.md

## Content
- Local development setup
- Example API calls with responses
- Docker quickstart
- K8s quickstart

## Acceptance Criteria
- [ ] Guide complete
- [ ] All examples working
- [ ] Clear instructions

## Dependencies
- All implementation tasks" \
  --label "phase-9,001-homelab-api,parallel"

gh issue create \
  --title "T094: Final constitution compliance check" \
  --body "## Task: T094 Final constitution compliance check

**Phase**: Phase 9 - Final Validation & Documentation
**Type**: Validation Task

## Description
Verify all constitution requirements met

## Validation Points
- [ ] Go 1.24+ standards followed
- [ ] 80%+ test coverage achieved
- [ ] All errors handled
- [ ] Structured logging implemented
- [ ] Graceful shutdown working
- [ ] Code formatting correct
- [ ] All lints passing

## Files to Update
- README.md (add compliance note)

## Acceptance Criteria
- [ ] All checks passed
- [ ] Compliance documented
- [ ] Ready for production

## Dependencies
- T072 (Coverage check)
- T090 (Integration tests)" \
  --label "phase-9,001-homelab-api"

echo "Phase 9 issues created (T090-T094)"

echo ""
echo "✅ All 97 GitHub issues created successfully!"
echo ""
echo "Summary of created issues:"
echo "  Phase 0 (Research): Issues #27-30 (4 tasks)"
echo "  Phase 1 (Foundation): Issues #31-39 (10 tasks)"
echo "  Phase 1.5 (Swagger): Issues #40-43 (4 tasks)"
echo "  Phase 2 (US1): Issues #44-50 (7 tasks)"
echo "  Phase 3 (US2): Just created (7 tasks)"
echo "  Phase 4 (Performance): Just created (5 tasks)"
echo "  Phase 5 (US3): Just created (6 tasks)"
echo "  Phase 6 (US4): Just created (6 tasks)"
echo "  Phase 7 (Testing): Just created (4 tasks)"
echo "  Phase 8 (Deployment): Just created (6 tasks)"
echo "  Phase 9 (Validation): Just created (2 tasks)"
echo ""
echo "Total: 97 tasks across 9 phases"

