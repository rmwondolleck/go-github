# go-github API Service Constitution

## Core Principles

### I. Go Standards Compliance
- Use Go 1.24 or later
- Follow official Go Code Review Comments guidelines
- Run `gofmt`, `go vet`, and `golangci-lint` before commits
- Handle all errors explicitly - no ignored errors
- Use context.Context for cancellation and timeouts

### II. API Design
- RESTful HTTP endpoints with standard methods (GET, POST, PUT, DELETE)
- JSON request/response format
- Consistent error response structure with status codes
- API versioning via URL path (e.g., `/api/v1/resource`)
- OpenAPI/Swagger documentation required

### III. Testing (NON-NEGOTIABLE)
- Minimum 80% code coverage for new code
- Unit tests for all business logic
- Integration tests for API endpoints
- Table-driven tests preferred for multiple scenarios
- Use `testing` package, `testify` for assertions allowed
- Tests must pass before merging to main

### IV. Observability
- Structured logging using standard library `log/slog`
- Log levels: DEBUG, INFO, WARN, ERROR
- Include request ID in all logs for tracing
- Metrics via Prometheus format (if applicable)
- Health check endpoint (`/health`) required

### V. Security & Configuration
- No hardcoded secrets - use environment variables
- Input validation on all API endpoints
- Rate limiting on public endpoints
- HTTPS in production
- Graceful shutdown handling (SIGTERM, SIGINT)

## Build & Deployment Standards

### Build Requirements
- `go.mod` and `go.sum` must be committed and up-to-date
- `make` targets for common operations: `build`, `test`, `lint`, `run`
- Docker support: Dockerfile with multi-stage build
- Binary built with version information embedded

### CI/CD Requirements
- Automated tests on all pull requests
- Build validation before merge
- No merge if tests fail or linting errors exist

## Code Organization

### Project Structure
```
/cmd          - Application entry points
/internal     - Private application code
/pkg          - Public libraries (if any)
/api          - API definitions (OpenAPI specs)
/tests        - Integration tests
```

### Dependency Management
- Minimal external dependencies
- Pin dependency versions in go.mod
- Review and justify all new dependencies
- Regular dependency updates for security patches

## Governance

This constitution represents the minimum viable standards for the go-github API service. All code must comply with these principles.

**Enforcement**:
- All pull requests require passing CI checks
- Code reviews must verify compliance
- Breaking changes require discussion and approval

**Version**: 1.0.0 | **Ratified**: 2026-02-28 | **Last Amended**: 2026-02-28

