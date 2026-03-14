---
# Shared guidelines for PR integration workflows
---

## PR Integration Guidelines

### Integration Order Heuristics

When determining the order to integrate PRs, follow these priority rules:

1. **Infrastructure first**: Middleware, server setup, core models
2. **Dependencies before dependents**: If PR B imports from PR A, integrate A first
3. **Smaller changes first**: PRs with fewer file changes are less likely to conflict
4. **Foundation → Features → Tests → Docs → Deployment**

### Typical Integration Order for Go Projects

```
1. go.mod / dependency changes
2. Core models and types (internal/models/)
3. Service interfaces (internal/*/service.go)
4. Middleware (internal/middleware/)
5. Response helpers (internal/handlers/response.go)
6. Service implementations (internal/*/service.go)
7. Handlers (internal/handlers/*.go)
8. Server/router setup (internal/server/)
9. Unit tests (internal/*_test.go)
10. Integration tests (tests/integration/)
11. API documentation (api/, Swagger)
12. Deployment configs (deployments/, Dockerfile, k8s/)
13. Documentation (README.md, docs/)
```

### File Conflict Resolution Matrix

| File Type | Strategy | Notes |
|-----------|----------|-------|
| `go.mod` | Merge all `require` blocks | Use highest version for duplicates |
| `*.go` (imports) | Union of all imports | Deduplicate and sort |
| `*.go` (functions) | Append new functions | Check for name collisions |
| `*_test.go` | Append test functions | Merge test tables for same function |
| `*.yaml` / `*.yml` | Deep merge | Arrays are concatenated |
| `Dockerfile` | Use most comprehensive | Verify all stages present |
| `README.md` | Section-based merge | Maintain TOC consistency |
| `Makefile` | Append new targets | Check for target name conflicts |

### Go-Specific Integration Checks

After integrating all PRs, verify:

- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] No duplicate function names across packages
- [ ] All new routes are registered in the router
- [ ] Middleware chain is correctly ordered
- [ ] All interfaces are properly implemented
- [ ] Import paths are consistent

### Commit Message Format

Each integration commit should follow this format:

```
Integrate PR #{number}: {original PR title}

Changes from: {head_branch}
Files integrated: {file_count}
Conflicts resolved: {conflict_count}
```

### Post-Integration Checklist

- [ ] All PR changes are present in the epic branch
- [ ] No functionality was lost during integration
- [ ] Tests pass for all integrated changes
- [ ] Documentation is consistent and complete
- [ ] Swagger docs are regenerated (if API changes)
- [ ] Docker build succeeds (if deployment changes)
