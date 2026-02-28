#!/bin/bash
# Create remaining GitHub issues (Phases 1-9)
set -e
OWNER="rmwondolleck"
REPO="go-github"
echo "Creating remaining phases issues..."
# Phase 3: User Story 2 - Health & Discovery (T030-T035)
echo "Creating Phase 3 issues (US2 - Health & Discovery)..."
gh issue create \
  --title "T030: Write unit tests for health checker" \
  --body "## Task: T030 [P] [US2] Write unit tests for health checker
**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Test Task (TDD)
**Parallelizable**: Yes
## Description
Write unit tests FIRST for health checker service
## Files to Create
- internal/health/checker_test.go
## Tests to Write
- TestCheck_ReturnsHealthyStatus
- TestCheck_IncludesUptime
## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All test cases documented
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
Write integration tests FIRST for health endpoint
## Files to Create
- tests/integration/health_test.go
## Tests to Write
- TestHealth_Returns200WithStatus
- TestHealth_IncludesComponentsAndUptime
- TestHealth_ResponseTimeLessThan50ms
## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] Performance validation included
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
## Acceptance Criteria
- [ ] Uptime calculated correctly
- [ ] Component status returned
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
Implement health check HTTP handler with Swagger documentation
## Files to Create
- internal/handlers/health.go
- internal/handlers/health_test.go
## Function to Implement
- HealthHandler(c *gin.Context)
## Acceptance Criteria
- [ ] Handler returns <10ms response time
- [ ] Swagger annotations added
- [ ] Tests from T031 now PASS
## Dependencies
- T032 (Service implemented)" \
  --label "phase-3,us-2,001-homelab-api"
gh issue create \
  --title "T034: Implement services discovery endpoint with Swagger docs" \
  --body "## Task: T034 [US2] Implement services discovery endpoint with Swagger docs
**Phase**: Phase 3 - User Story 2
**User Story**: US2 - Health Check and Service Discovery (P1 MVP)
**Type**: Implementation Task
## Description
Implement services discovery endpoint with static service list
## Files to Create
- internal/handlers/services.go
- internal/handlers/services_test.go
## Function to Implement
- ListServicesHandler(c *gin.Context)
## Acceptance Criteria
- [ ] Handler returns service list
- [ ] Swagger annotations added
- [ ] Unit tests included
## Dependencies
- T019d (Swagger initialized)" \
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
- [ ] Routes registered
- [ ] Endpoints accessible
- [ ] All tests passing
## Dependencies
- T033, T034 (Handlers implemented)" \
  --label "phase-3,us-2,001-homelab-api"
echo "Phase 3 issues created"
# Phase 4: Performance & Middleware (T040-T044)
echo "Creating Phase 4 issues (Performance & Middleware)..."
gh issue create \
  --title "T040: Implement rate limiting middleware" \
  --body "## Task: T040 [P] Implement rate limiting middleware
**Phase**: Phase 4 - Performance Optimization & Middleware (P2)
**Type**: Implementation Task
**Parallelizable**: Yes
## Description
Implement rate limiting middleware with token bucket algorithm
## Files to Create
- internal/middleware/ratelimit.go
- internal/middleware/ratelimit_test.go
## Features
- 500 req/min per IP
- Token bucket algorithm
- X-RateLimit-* headers
- 429 Too Many Requests response
## Acceptance Criteria
- [ ] Rate limiting working
- [ ] Headers correct
- [ ] Tests passing
## Notes
Performance optimization task" \
  --label "phase-4,001-homelab-api,parallel"
gh issue create \
  --title "T041: Implement CORS middleware" \
  --body "## Task: T041 [P] Implement CORS middleware
**Phase**: Phase 4 - Performance Optimization & Middleware (P2)
**Type**: Implementation Task
**Parallelizable**: Yes
## Description
Implement CORS middleware for browser-based clients
## Files to Create
- internal/middleware/cors.go
- internal/middleware/cors_test.go
## Features
- Allow origin: http://localhost:3000 (configurable)
- Allow methods: GET, POST, OPTIONS
- Allow headers: Content-Type, Authorization
## Acceptance Criteria
- [ ] CORS headers correct
- [ ] Configuration working
- [ ] Tests passing
## Notes
Cross-origin support for web clients" \
  --label "phase-4,001-homelab-api,parallel"
gh issue create \
  --title "T042: Optimize JSON encoding with jsoniter" \
  --body "## Task: T042 [P] Optimize JSON encoding with jsoniter
**Phase**: Phase 4 - Performance Optimization & Middleware (P2)
**Type**: Implementation Task
**Parallelizable**: Yes
## Description
Optimize JSON encoding for improved performance
## Files to Update
- internal/handlers/response.go
## Features
- Replace stdlib encoding/json with jsoniter
- Benchmark performance improvement
- Document results
## Acceptance Criteria
- [ ] jsoniter integrated
- [ ] Benchmarks show 2-3x improvement
- [ ] Results documented in performance.md
## Notes
Performance optimization - 2-3x faster JSON" \
  --label "phase-4,001-homelab-api,parallel"
gh issue create \
  --title "T043: Implement response pooling for device list" \
  --body "## Task: T043 [P] Implement response pooling for device list
**Phase**: Phase 4 - Performance Optimization & Middleware (P2)
**Type**: Implementation Task
**Parallelizable**: Yes
## Description
Implement sync.Pool for zero-allocation device list responses
## Files to Update
- internal/handlers/homeassistant.go
## Features
- sync.Pool for DeviceListResponse
- Pre-allocated structs
- Benchmark allocation reduction
## Acceptance Criteria
- [ ] Pool implemented
- [ ] Benchmarks show allocation reduction
- [ ] No memory leaks
## Notes
Performance optimization - zero-allocation responses" \
  --label "phase-4,001-homelab-api,parallel"
gh issue create \
  --title "T044: Update server router with new middleware" \
  --body "## Task: T044 Update server router with new middleware
**Phase**: Phase 4 - Performance Optimization & Middleware (P2)
**Type**: Implementation Task
## Description
Apply all middleware to Gin router in correct order
## Files to Update
- internal/server/server.go
## Middleware Order
1. Recovery
2. RequestID
3. Logging
4. RateLimit (on /api/v1/*)
5. CORS
6. Handlers
## Acceptance Criteria
- [ ] All middleware applied
- [ ] Order correct
- [ ] Tests passing
## Dependencies
- T040, T041, T042, T043 (Middleware implemented)" \
  --label "phase-4,001-homelab-api"
echo "Phase 4 issues created"
# Phase 5: User Story 3 - Device Control (T050-T055)
echo "Creating Phase 5 issues (US3 - Device Control)..."
gh issue create \
  --title "T050: Write unit tests for command execution" \
  --body "## Task: T050 [P] [US3] Write unit tests for command execution
**Phase**: Phase 5 - User Story 3 (P2)
**User Story**: US3 - Control HomeAssistant Devices
**Type**: Test Task (TDD)
**Parallelizable**: Yes
## Description
Write unit tests FIRST for command execution service
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
## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-5,us-3,001-homelab-api,parallel"
gh issue create \
  --title "T051: Write integration tests for command endpoint" \
  --body "## Task: T051 [P] [US3] Write integration tests for command endpoint
**Phase**: Phase 5 - User Story 3 (P2)
**User Story**: US3 - Control HomeAssistant Devices
**Type**: Test Task (TDD)
**Parallelizable**: Yes
## Description
Write integration tests FIRST for command endpoint
## Files to Update
- tests/integration/devices_test.go
## Tests to Write
- TestExecuteCommand_Returns200ForValidCommand
- TestExecuteCommand_Returns400ForInvalidAction
- TestExecuteCommand_Returns405ForReadOnlyDevice
## Acceptance Criteria
- [ ] Tests created and FAIL before implementation
- [ ] All HTTP scenarios covered
## Notes
TDD: Tests MUST be written before implementation" \
  --label "phase-5,us-3,001-homelab-api,parallel"
gh issue create \
  --title "T052: Define Command model" \
  --body "## Task: T052 [US3] Define Command model
**Phase**: Phase 5 - User Story 3 (P2)
**User Story**: US3 - Control HomeAssistant Devices
**Type**: Implementation Task
## Description
Create Command model with action and parameters
## Files to Create
- internal/homeassistant/types.go
## Struct to Define
- Command with Action and Parameters fields
## Acceptance Criteria
- [ ] Command struct defined
- [ ] JSON tags added
- [ ] Validation ready
## Dependencies
- T050, T051 (Tests written)" \
  --label "phase-5,us-3,001-homelab-api"
gh issue create \
  --title "T053: Implement command execution in service" \
  --body "## Task: T053 [US3] Implement command execution in service
**Phase**: Phase 5 - User Story 3 (P2)
**User Story**: US3 - Control HomeAssistant Devices
**Type**: Implementation Task
## Description
Implement command execution in HomeAssistant service
## Files to Update
- internal/homeassistant/service.go
## Function to Implement
- ExecuteCommand(ctx context.Context, id string, cmd Command) error
## Features
- Validate device exists and controllable
- Validate action supported for device type
- Mock command execution
## Acceptance Criteria
- [ ] Service method implemented
- [ ] Tests from T050 now PASS
## Dependencies
- T052 (Command model defined)" \
  --label "phase-5,us-3,001-homelab-api"
gh issue create \
  --title "T054: Implement command endpoint handler with Swagger docs" \
  --body "## Task: T054 [US3] Implement command endpoint handler with Swagger docs
**Phase**: Phase 5 - User Story 3 (P2)
**User Story**: US3 - Control HomeAssistant Devices
**Type**: Implementation Task
## Description
Implement command execution HTTP handler with Swagger documentation
## Files to Update
- internal/handlers/homeassistant.go
## Function to Implement
- ExecuteCommandHandler(c *gin.Context)
## Features
- Parse command from request body
- Validate JSON payload
- Return appropriate error codes (400, 404, 405)
- Swagger annotations
## Acceptance Criteria
- [ ] Handler implemented
- [ ] Error handling correct
- [ ] Tests from T051 now PASS
## Dependencies
- T053 (Service implemented)" \
  --label "phase-5,us-3,001-homelab-api"
gh issue create \
  --title "T055: Register command endpoint" \
  --body "## Task: T055 [US3] Register command endpoint
**Phase**: Phase 5 - User Story 3 (P2)
**User Story**: US3 - Control HomeAssistant Devices
**Type**: Implementation Task
## Description
Register device command endpoint in Gin router
## Files to Update
- internal/server/server.go
## Route to Register
- POST /api/v1/homeassistant/devices/:id/command
## Acceptance Criteria
- [ ] Route registered
- [ ] Endpoint accessible
- [ ] Tests passing
## Dependencies
- T054 (Handler implemented)" \
  --label "phase-5,us-3,001-homelab-api"
echo "Phase 5 issues created"
echo ""
echo "✅ Issues for Phases 3-5 created successfully!"
