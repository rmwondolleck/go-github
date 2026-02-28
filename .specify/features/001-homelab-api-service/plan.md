# Implementation Plan: Home Lab API Service

**Branch**: `001-homelab-api-service` | **Date**: 2026-02-28 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `.specify/features/001-homelab-api-service/spec.md`

## Summary

Build a high-performance REST API service for home lab monitoring with mocked HomeAssistant integration. Core focus: sub-50ms health check responses, sub-200ms endpoint latency using Gin framework with optimized in-memory data structures. Backend services designed to be MCP-tool-ready with zero HTTP coupling in business logic.

**Key Technical Decisions**:
- Gin framework for performance (vs stdlib net/http)
- In-memory sync.Map for thread-safe device storage with O(1) lookups
- Pre-allocated response structs with object pooling for zero-allocation responses
- Structured logging with slog, middleware-based request ID propagation
- JSON encoding optimized with jsoniter for 2-3x faster serialization

## Technical Context

**Language/Version**: Go 1.24+  
**Primary Dependencies**: 
- `gin-gonic/gin` v1.10+ (HTTP framework, optimized routing)
- `swaggo/swag` (Swagger/OpenAPI documentation generation)
- `swaggo/gin-swagger` (Gin middleware for Swagger UI)
- `swaggo/files` (Swagger static file handler)
- `json-iterator/go` (faster JSON encoding/decoding)
- `stretchr/testify` (testing assertions)
- Standard library: `log/slog`, `context`, `sync`

**Storage**: In-memory only (sync.Map for device storage, no persistence)  
**Testing**: Go standard `testing` package + testify assertions, table-driven tests  
**Target Platform**: Linux container (Docker) → Kubernetes cluster  
**Project Type**: Web service (REST API)  
**Performance Goals**: 
- Health check: <50ms p99
- All endpoints: <200ms p99
- Throughput: 100+ concurrent requests without degradation
- Memory: <100MB footprint

**Constraints**: 
- Zero database - all data in-memory
- No external service calls (mocked data only)
- MCP-ready architecture (services usable without HTTP layer)
- 80%+ test coverage mandatory

**Scale/Scope**: 
- POC for ~50 mocked devices
- 6 REST endpoints
- Single service integration (HomeAssistant)
- 100 requests/minute rate limit per client

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

✅ **Go 1.24+ Compliance**: Using Go 1.24, follows Code Review Comments guidelines  
✅ **Test-First Mandate**: TDD workflow required, 80% coverage enforced  
✅ **Error Handling**: All errors explicitly handled, no ignored errors  
✅ **Structured Logging**: Using log/slog with request IDs  
✅ **Project Structure**: Follows constitution conventions (/cmd, /internal, /api, /tests)  
✅ **API Design**: RESTful, JSON, versioned at /api/v1, OpenAPI documented  
✅ **Security**: No hardcoded secrets, environment variables for config  
✅ **Graceful Shutdown**: SIGTERM/SIGINT handling required  
✅ **Observability**: Health endpoint + structured logging mandatory  

**Compliance Status**: ✅ PASS - No violations, all constitution requirements met

## Project Structure

### Documentation (this feature)

```text
.specify/features/001-homelab-api-service/
├── spec.md              # Feature specification (completed)
├── plan.md              # This file - implementation plan
├── research.md          # Phase 0: Technology research (next)
├── data-model.md        # Phase 1: Data structures & interfaces
├── api-design.md        # Phase 1: OpenAPI specification
├── performance.md       # Phase 1: Performance optimization strategies
├── contracts/           # Phase 1: Service contracts & interfaces
│   ├── homeassistant.go
│   └── health.go
└── tasks.md             # Phase 2: Task breakdown (via /speckit.tasks)
```

### Source Code (repository root)

```text
go-github/
├── cmd/
│   └── api/
│       └── main.go           # Entry point, server initialization
├── internal/
│   ├── homeassistant/
│   │   ├── service.go        # HomeAssistant business logic (MCP-ready)
│   │   ├── service_test.go   # Service unit tests
│   │   ├── mock_data.go      # Mocked device data
│   │   └── types.go          # Domain types (Device, Command)
│   ├── cluster/
│   │   ├── service.go        # Cluster info service (MCP-ready)
│   │   ├── service_test.go
│   │   └── types.go
│   ├── health/
│   │   ├── checker.go        # Health check logic
│   │   └── checker_test.go
│   ├── middleware/
│   │   ├── logging.go        # Request logging with IDs
│   │   ├── ratelimit.go      # Rate limiting middleware
│   │   ├── cors.go           # CORS headers
│   │   └── recovery.go       # Panic recovery
│   ├── handlers/
│   │   ├── health.go         # Health endpoint handler
│   │   ├── health_test.go    # Handler tests
│   │   ├── homeassistant.go  # HomeAssistant endpoint handlers
│   │   ├── homeassistant_test.go
│   │   ├── services.go       # Services list endpoint
│   │   └── response.go       # Response helper utilities
│   ├── models/
│   │   ├── device.go         # Shared Device model
│   │   ├── error.go          # Error response model
│   │   └── health.go         # Health status model
│   └── server/
│       ├── server.go         # Gin server setup & routing
│       ├── server_test.go
│       └── shutdown.go       # Graceful shutdown handler
├── api/
│   └── openapi.yaml          # OpenAPI 3.0 specification
├── tests/
│   ├── integration/
│   │   ├── health_test.go    # Integration tests for health
│   │   ├── devices_test.go   # Integration tests for devices
│   │   └── helper.go         # Test helpers & setup
│   └── load/
│       └── concurrent_test.go # Concurrent load tests
├── deployments/
│   ├── Dockerfile            # Multi-stage Docker build
│   └── k8s/
│       ├── deployment.yaml   # K8s deployment manifest
│       ├── service.yaml      # K8s service manifest
│       └── configmap.yaml    # Configuration
├── Makefile                  # Build targets
├── go.mod
├── go.sum
└── README.md
```

**Structure Decision**: Single Go project following constitution standards. Backend services in `/internal` are HTTP-agnostic and MCP-ready. Clean separation: handlers (HTTP thin wrappers) → services (pure business logic) → models (data structures).

## Phase 0: Research & Technology Validation

**Objective**: Validate technology choices and performance characteristics before design.

### Research Tasks

1. **Gin Framework Performance Baseline**
   - Benchmark basic Gin "Hello World" endpoint
   - Measure routing overhead vs stdlib
   - Validate <50ms health check achievable
   - Document memory allocation patterns

2. **JSON Serialization Optimization**
   - Benchmark stdlib encoding/json vs jsoniter
   - Test pre-allocated struct encoding
   - Measure sync.Pool effectiveness for response objects
   - Determine optimal buffer sizes

3. **In-Memory Storage Strategy**
   - Compare sync.Map vs sync.RWMutex + map for read-heavy workloads
   - Benchmark device lookup performance (target: <1ms)
   - Test concurrent read/write scenarios
   - Measure memory overhead per device

4. **Rate Limiting Options**
   - Evaluate gin middleware options (e.g., gin-contrib/rate)
   - Test token bucket vs leaky bucket algorithms
   - Measure rate limiter overhead (<5ms)
   - Validate IP-based vs token-based limiting

5. **Structured Logging Performance**
   - Benchmark log/slog with different output formats
   - Test async logging vs synchronous
   - Measure request ID propagation overhead
   - Determine optimal log level defaults

**Deliverable**: `research.md` with benchmarks, recommendations, and validated tech choices.

## Phase 1: Design & Architecture

**Objective**: Define data models, API contracts, and service interfaces before implementation.

### Design Artifacts

1. **Data Model (`data-model.md`)**
   - Device struct with optimized field layout (cache-friendly)
   - Command struct for device control
   - HealthStatus struct with component breakdown
   - Error response model with consistent fields
   - JSON tags for all exported fields

2. **API Design (`api-design.md`)**
   - Complete OpenAPI 3.0 specification
   - Request/response schemas for all endpoints
   - Error code enumeration
   - Rate limit headers
   - Example requests/responses

3. **Performance Optimization Strategy (`performance.md`)**
   - Response object pooling with sync.Pool
   - Zero-allocation JSON encoding techniques
   - Pre-computed JSON for static responses (health check)
   - Connection pooling considerations
   - Memory profiling plan

4. **Service Contracts (`contracts/`)**
   ```go
   // contracts/homeassistant.go
   package contracts
   
   import "context"
   
   // HomeAssistantService defines MCP-ready business logic
   type HomeAssistantService interface {
       ListDevices(ctx context.Context) ([]Device, error)
       GetDevice(ctx context.Context, id string) (*Device, error)
       ExecuteCommand(ctx context.Context, id string, cmd Command) error
   }
   
   // contracts/health.go
   type HealthChecker interface {
       Check(ctx context.Context) HealthStatus
   }
   ```

5. **Quick Start Guide (`quickstart.md`)**
   - Local development setup
   - Running with `go run`
   - Building Docker image
   - Deploying to K8s
   - Testing endpoints with curl examples

**Deliverable**: Complete design documentation ready for task breakdown.

## Phase 2: Task Generation

**Objective**: Break down implementation into atomic, testable tasks.

**Process**: Run `/speckit.tasks` command to generate `tasks.md` from this plan and spec.

**Task Categories** (preview):

1. **Foundation** (P1)
   - Set up project structure
   - Configure Gin server with graceful shutdown
   - Implement request ID middleware
   - Add structured logging

2. **Core Services** (P1)
   - Implement HomeAssistant service with mock data
   - Create device storage (sync.Map)
   - Add health checker
   - Build response helpers

3. **HTTP Handlers** (P1)
   - Health endpoint handler
   - Device list handler
   - Device detail handler
   - Services discovery handler

4. **Testing** (P1)
   - Unit tests for all services (80%+ coverage)
   - Integration tests for endpoints
   - Table-driven handler tests

5. **Performance & Middleware** (P2)
   - Rate limiting middleware
   - CORS middleware
   - Response pooling optimization
   - JSON encoding optimization

6. **Device Control** (P2)
   - Command execution service
   - Device command handler
   - Command validation

7. **Deployment** (P3)
   - Dockerfile with multi-stage build
   - K8s manifests
   - Environment configuration
   - Documentation

8. **Load Testing** (P3)
   - Concurrent request tests
   - Latency benchmarks
   - Memory profiling

## Performance Optimization Strategy

### Target Metrics
- Health check: <50ms p99 (goal: <10ms p99)
- Device list: <100ms p99 (50 devices)
- Device detail: <50ms p99
- Command execution: <100ms p99
- Memory: <100MB under 100 concurrent requests

### Optimization Techniques

1. **Zero-Allocation Response Handling**
   ```go
   var deviceListPool = sync.Pool{
       New: func() interface{} {
           return &DeviceListResponse{
               Devices: make([]Device, 0, 50),
           }
       },
   }
   ```

2. **Pre-computed JSON Responses**
   - Health check response pre-marshaled
   - Services list static response
   - Update on server start only

3. **Efficient In-Memory Lookups**
   - sync.Map for O(1) device access
   - Pre-indexed by device ID
   - No locks for read operations

4. **Fast JSON Encoding**
   - Use jsoniter for 2-3x faster encoding
   - Pre-allocate buffers for known response sizes
   - Disable HTML escaping when safe

5. **Connection Optimization**
   - HTTP/2 support enabled
   - Keep-alive with reasonable timeout
   - Connection pooling for internal calls (if added)

6. **Middleware Ordering**
   - Recovery → RequestID → Logging → RateLimit → CORS → Handler
   - Minimize allocations in hot path

### Profiling Plan
1. Baseline benchmarks before optimization
2. CPU profiling with `pprof` during load tests
3. Memory profiling for allocation hotspots
4. Continuous benchmarking in CI

## Testing Strategy

### Test Pyramid

**Unit Tests (70%)**
- All services tested independently
- Mock-free pure functions preferred
- Table-driven tests for handlers
- Edge cases and error paths covered

**Integration Tests (25%)**
- Full HTTP request/response cycle
- All endpoints tested end-to-end
- Error scenarios (404, 400, 429, 500)
- Concurrent request handling

**Load Tests (5%)**
- 100 concurrent requests
- Sustained load for 1 minute
- Latency percentile validation
- Memory leak detection

### Test Infrastructure

```go
// tests/integration/helper.go
func SetupTestServer(t *testing.T) *httptest.Server {
    // Create server with all middleware
    // Pre-load mock data
    // Return test server
}

func AssertResponseTime(t *testing.T, duration time.Duration, threshold time.Duration) {
    // Assert response time under threshold
}
```

### CI/CD Integration
- All tests run on PR
- Benchmarks compared to baseline
- Coverage report generated
- Lint checks (gofmt, go vet, golangci-lint)

## Development Workflow

### Phase 0: Research (1-2 days)
1. Run benchmarks for Gin, jsoniter, sync.Map
2. Document findings in `research.md`
3. Validate performance targets achievable

### Phase 1: Design (1 day)
1. Write data models with struct field tags
2. Create OpenAPI specification
3. Define service interfaces
4. Document optimization strategies

### Phase 2: Implementation (3-5 days)
1. Foundation setup (project structure, server, middleware)
2. P1 services and handlers (HomeAssistant, health)
3. P1 testing (unit + integration)
4. P2 features (rate limiting, commands, optimization)
5. P2 testing
6. P3 deployment (Docker, K8s)

### Phase 3: Validation (1 day)
1. Load testing and profiling
2. Performance metric validation
3. Documentation review
4. Deployment to test cluster

## Dependencies & External Tools

### Build Dependencies
- Go 1.24+ installed
- Docker for containerization
- Make for build automation
- golangci-lint for code quality

### Runtime Dependencies (Go modules)
```go
// go.mod
module go-github

go 1.24

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/json-iterator/go v1.1.12
    github.com/stretchr/testify v1.9.0
    github.com/swaggo/swag v1.16.3
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/files v1.0.1
)
```

### Development Tools
- `pprof` for profiling
- `vegeta` or `hey` for load testing
- `curl` / `httpie` for manual testing
- `swagger-ui` for API docs viewing

## Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Gin performance insufficient for <50ms target | High | Phase 0 benchmarking validates; fallback to stdlib with custom optimizations |
| sync.Map doesn't scale to 100 concurrent | Medium | Phase 0 load testing; alternative: partitioned RWMutex maps |
| JSON encoding bottleneck | Medium | jsoniter + response pooling tested in Phase 0 |
| Rate limiting adds excessive latency | Low | Benchmark middleware overhead; use efficient token bucket |
| Test coverage below 80% | High | TDD mandatory; coverage checked in CI |

## Next Steps

1. ✅ **Completed**: Feature specification written
2. ✅ **Completed**: Implementation plan documented (this file)
3. **Next**: Execute Phase 0 research (`research.md`)
4. **Then**: Execute Phase 1 design (data models, API spec, contracts)
5. **Then**: Generate tasks with `/speckit.tasks` command
6. **Then**: Begin implementation following task breakdown

## Open Questions for Clarification

1. ✅ **FR-015 CORS Origins**: `["http://localhost:3000"]` for dev, configurable via env var - **APPROVED**
2. ✅ **Rate Limiting**: **500 req/min** per IP (updated from 100) - **APPROVED**
3. ✅ **Metrics Endpoint**: Defer to future iteration - **APPROVED**
4. ✅ **Device Filtering**: Support `?type=light` filtering - **APPROVED**
5. ✅ **Mock Data Configuration**: JSON file for mock devices - **APPROVED**

**Status**: ✅ All questions answered. Ready for task generation.




