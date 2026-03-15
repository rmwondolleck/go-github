# Code Coverage Report

**Generated**: March 14, 2026  
**Tool**: `go test -coverprofile=coverage.out ./...`  
**Go Version**: 1.25

## Summary

| Metric | Value |
|--------|-------|
| **Total Coverage** | **96.1%** ✅ |
| **Target** | ≥ 80% |
| **Status** | PASS |

## Coverage by Package

| Package | Coverage | Status |
|---------|----------|--------|
| `go-github/internal/cluster` | 100.0% | ✅ |
| `go-github/internal/handlers` | 87.8% | ✅ |
| `go-github/internal/health` | 100.0% | ✅ |
| `go-github/internal/homeassistant` | 100.0% | ✅ |
| `go-github/internal/middleware` | 100.0% | ✅ |
| `go-github/internal/models` | N/A (no test files) | ℹ️ |
| `go-github/internal/server` | 100.0% | ✅ |

## Coverage by Function

| File | Function | Coverage |
|------|----------|----------|
| `internal/cluster/service.go` | `NewService` | 100.0% |
| `internal/cluster/service.go` | `ListServices` | 100.0% |
| `internal/handlers/cluster.go` | `ListClusterServicesHandler` | 71.4% |
| `internal/handlers/homeassistant.go` | `getResponseFromPool` | 100.0% |
| `internal/handlers/homeassistant.go` | `putResponseInPool` | 100.0% |
| `internal/handlers/homeassistant.go` | `DeviceListHandler` | 100.0% |
| `internal/handlers/homeassistant.go` | `ExecuteCommandHandler` | 100.0% |
| `internal/handlers/response.go` | `JSONSuccess` | 60.0% |
| `internal/handlers/response.go` | `JSONError` | 66.7% |
| `internal/handlers/response.go` | `NotFound` | 100.0% |
| `internal/handlers/response.go` | `BadRequest` | 100.0% |
| `internal/handlers/response.go` | `InternalError` | 100.0% |
| `internal/handlers/services.go` | `ListServicesHandler` | 100.0% |
| `internal/health/checker.go` | `NewChecker` | 100.0% |
| `internal/health/checker.go` | `Check` | 100.0% |
| `internal/health/checker.go` | `formatUptime` | 100.0% |
| `internal/homeassistant/types.go` | `Validate` | 100.0% |
| `internal/middleware/cors.go` | `CORS` | 100.0% |
| `internal/middleware/cors.go` | `parseOrigins` | 100.0% |
| `internal/middleware/cors.go` | `isOriginAllowed` | 100.0% |
| `internal/middleware/logging.go` | `Logger` | 100.0% |
| `internal/middleware/recovery.go` | `Recovery` | 100.0% |
| `internal/middleware/request_id.go` | `RequestID` | 100.0% |
| `internal/server/server.go` | `New` | 100.0% |
| `internal/server/server.go` | `Run` | 100.0% |
| `internal/server/server.go` | `Router` | 100.0% |
| `internal/server/server.go` | `healthHandler` | 100.0% |
| `internal/server/server.go` | `apiRootHandler` | 100.0% |
| `internal/server/shutdown.go` | `GracefulShutdown` | 100.0% |

## Files with Partial Coverage

### `internal/handlers/cluster.go` — 71.4%

The error path in `ListClusterServicesHandler` (when `svc.ListServices` returns an error) is not reached in unit tests because the mock cluster service always returns successfully. The error branch is covered at the integration level.

### `internal/handlers/response.go` — `JSONSuccess` 60%, `JSONError` 66.7%

The fallback branches triggered when `jsoniter.Marshal` fails are not exercised in unit tests. These branches guard against an internal encoding failure that would only occur with non-serializable types — a defensive coding pattern that is correct but difficult to trigger with standard types.

## Comparison Against 80% Target

```
Total Coverage : 96.1%
Target         : 80.0%
Margin         : +16.1%
Status         : ✅ PASS
```

## How to Reproduce

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View total coverage
go tool cover -func=coverage.out | grep total

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

## Recommendations

The remaining uncovered lines are intentional defensive fallbacks (JSON encoding errors, service error paths) that are only reachable under abnormal runtime conditions. No further action is required to meet the 80% requirement. If 100% coverage is desired for these functions, consider:

1. **`ListClusterServicesHandler` error branch**: Introduce a dependency injection point or interface mock to simulate service failures.
2. **`JSONSuccess`/`JSONError` marshal error branches**: Use a custom type that fails JSON marshaling to trigger these paths.
