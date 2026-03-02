# PR Merge Instructions - UPDATED

**Status**: ⚠️ OUTDATED - See CURRENT_STATUS.md for latest  
**Date**: March 1, 2026  
**Current Open PRs**: 13 (all draft status)

> **Note**: This file contains outdated merge instructions. The actual current state is:
> - All 13 open PRs are in DRAFT status
> - PRs referenced in original instructions may have been merged or closed
> - See `CURRENT_STATUS.md` for accurate, up-to-date information

---

## Step-by-Step Instructions

### 1️⃣ Merge This PR First (CRITICAL)
```bash
# On GitHub UI:
# 1. Go to PR #136
# 2. Click "Ready for review" (if still draft)
# 3. Wait for checks to pass
# 4. Click "Squash and merge"
# 5. Confirm merge
```

**Why First?** This PR fixes the empty `pr-build-validation.yml` file that's blocking all other PR merges.

---

### 2️⃣ Merge Foundation PRs (Phase 3 - Health Check)

#### PR #110: Health Checker Tests
```bash
# On GitHub UI:
# 1. Go to PR #110
# 2. It's already "Ready for review" (not draft)
# 3. Click "Squash and merge"
# 4. Close issue #31
```

#### PR #112: Health Checker Service
```bash
# On GitHub UI or command line:
# 1. Rebase PR #112 from main:
git checkout copilot/implement-health-checker-service
git pull origin copilot/implement-health-checker-service
git rebase main
git push -f origin copilot/implement-health-checker-service

# 2. On GitHub UI:
# 3. Go to PR #112
# 4. Wait for checks to pass
# 5. Click "Squash and merge"
# 6. Close issue #33
```

---

### 3️⃣ Merge Documentation PRs (Parallel)

For each of these PRs, the process is the same:

#### Process for PRs #126, #128, #130, #134:
```bash
# On GitHub UI:
# 1. Go to PR #[NUMBER]
# 2. Click "Ready for review"
# 3. Wait for checks to pass
# 4. Click "Squash and merge"
# 5. Close associated issue
```

**PR #126**: Copilot instructions (no issue)
**PR #128**: Swagger UI → Close issue #54
**PR #130**: README docs → Close issue #66
**PR #134**: Deployment docs → Close issue #63

---

### 4️⃣ Merge Infrastructure PR

#### PR #124: K8s Service Manifest
```bash
# On GitHub UI:
# 1. Review deployments/k8s/service.yaml
# 2. Verify manifest is correct
# 3. Click "Ready for review"
# 4. Click "Squash and merge"
# 5. Close issue #60
```

---

### 5️⃣ Handle Feature PRs (Need Review/Completion)

#### PR #132: Command Model
```bash
# On GitHub UI:
# 1. Review internal/homeassistant/types.go
# 2. Verify implementation and tests
# 3. If good: Click "Ready for review" → Merge
# 4. Close issue #44
```

#### PR #118: Command Tests
```bash
# On GitHub UI:
# 1. Review TDD tests in internal/homeassistant/service_test.go
# 2. If good: Click "Ready for review" → Merge
# 3. Close issue #42
```

#### PR #114: Services Discovery (Needs Work)
```bash
# Option A: Complete it yourself
git checkout copilot/implement-services-discovery-endpoint
git pull
git rebase main
# Complete the implementation
# Add tests
git push

# Option B: Assign to agent
# Close PR and create new issue for agent

# Then merge and close issue #35
```

#### PR #116: CORS Middleware (Needs Work)
```bash
# Option A: Complete it yourself
git checkout copilot/implement-cors-middleware
git pull  
git rebase main
# Complete the implementation
# Add tests
git push

# Option B: Assign to agent
# Close PR and create new issue for agent

# Then merge and close issue #38
```

---

### 6️⃣ Close Redundant PR

#### PR #135: Build Validation Workflow
```bash
# On GitHub UI:
# 1. Go to PR #135
# 2. Click "Close pull request"
# 3. Add comment: "Closed as redundant - fix applied in PR #136"
```

---

## Merge Checklist

Copy this to track progress:

```
Wave 1: Critical Blocker
[ ] PR #136 - THIS PR (monitor/manage PRs)

Wave 2: Foundation  
[ ] PR #110 - Health checker tests
[ ] PR #112 - Health checker service (rebase first)

Wave 3: Documentation
[ ] PR #126 - Copilot instructions
[ ] PR #128 - Swagger UI
[ ] PR #130 - README
[ ] PR #134 - Deployment docs

Wave 4: Infrastructure
[ ] PR #124 - K8s service manifest

Wave 5: Features (review/complete first)
[ ] PR #132 - Command model
[ ] PR #118 - Command tests
[ ] PR #114 - Services discovery (complete)
[ ] PR #116 - CORS middleware (complete)

Cleanup
[ ] PR #135 - CLOSE (redundant)
```

---

## After All PRs Merged

### Close Associated Issues
```
✅ Issue #31 (PR #110)
✅ Issue #33 (PR #112)
✅ Issue #35 (PR #114)
✅ Issue #38 (PR #116)
✅ Issue #42 (PR #118)
✅ Issue #44 (PR #132)
✅ Issue #54 (PR #128)
✅ Issue #60 (PR #124)
✅ Issue #63 (PR #134)
✅ Issue #66 (PR #130)
```

### Phase 2 Dispatch
Once T001-T019 are complete (Phase 0 & 1):
```bash
# Dispatch Phase 2 agents for:
# - Batch 2: Issues #27-34 (US1 - Device Management)
# - Use GitHub issue assignment/comments
# - Follow DISPATCH_SCHEDULE.md
```

---

## Need Help?

See detailed analysis in `PR_STATUS_REPORT.md`

---

**Quick Stats**:
- 13 open PRs
- 10 can merge now or after simple review
- 2 need completion
- 1 should be closed
- 10 issues to close

**Time Estimate**: 30-60 minutes to merge all ready PRs

**Goal**: All 97 issues from DISPATCH_SCHEDULE.md completed
