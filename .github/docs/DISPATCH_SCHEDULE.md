# Home Lab API Service - Agent Dispatch Plan

**Status**: ✅ All GitHub issues created (Issues #27-77)
**Total Tasks**: 50 unique issues across 11 batches
**Repository**: https://github.com/rmwondolleck/go-github

**Note**: The issue set represents core tasks for 11 phases. Some additional sub-tasks referenced in dependencies will be created as work progresses through batches.

---

## Batch Dispatch Schedule

### ✅ Batch 1: Phase 0 - Research & Validation
**Issues**: #27-30 (4 tasks)
**Duration**: 1-2 days
**Status**: READY FOR DISPATCH
**References**: `.specify/AGENT_BRIEFS.md` - Batch 1 section

**Dispatch Command**:
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 1 section)
GitHub Issues: #27, #28, #29, #30
Branch: batch-1-phase-0-research
```

---

### ✅ Batch 2: Phase 1 - Foundation Setup  
**Issues**: #31-39 (10 tasks)
**Duration**: 1-2 days
**Status**: WAITING (depends on Batch 1 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 2 section

**Dispatch Command** (after Batch 1 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 2 section)
GitHub Issues: #31-#39
Branch: batch-2-phase-1-foundation
```

---

### ✅ Batch 3: Phase 1.5 - Swagger Setup
**Issues**: #40-43 (4 tasks)
**Duration**: 1 day
**Status**: WAITING (depends on Batch 2 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 3 section

**Dispatch Command** (after Batch 2 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 3 section)
GitHub Issues: #40-#43
Branch: batch-3-phase-1-5-swagger
```

---

### ✅ Batch 4: Phase 2 - User Story 1 (Device Status)
**Issues**: #44-50 (7 tasks)
**Duration**: 2-3 days
**Status**: WAITING (depends on Batch 3 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 4 section

**Dispatch Command** (after Batch 3 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 4 section)
GitHub Issues: #44-#50
Branch: batch-4-phase-2-us1-device-status
```

---

### ✅ Batch 5: Phase 3 - User Story 2 (Health & Discovery)
**Issues**: #69-74 (6 tasks)
**Duration**: 2-3 days
**Status**: CAN RUN IN PARALLEL WITH BATCH 4 (after Batch 3)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 5 section

**Dispatch Command** (after Batch 3 merged, parallel with Batch 4):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 5 section)
GitHub Issues: #69-#74
Branch: batch-5-phase-3-us2-health-discovery
```

---

### ✅ Batch 6: Phase 4 - Performance & Middleware
**Issues**: #75-77, #39-40 (5 tasks)
**Duration**: 1-2 days
**Status**: WAITING (depends on Batches 4 & 5 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 6 section

**Dispatch Command** (after Batches 4 & 5 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 6 section)
GitHub Issues: #75-#77, #39-#40
Branch: batch-6-phase-4-performance
```

---

### ✅ Batch 7: Phase 5 - User Story 3 (Device Control)
**Issues**: #42-47 (6 tasks)
**Duration**: 2-3 days
**Status**: WAITING (depends on Batch 6 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 7 section

**Dispatch Command** (after Batch 6 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 7 section)
GitHub Issues: #42-#47
Branch: batch-7-phase-5-us3-device-control
```

---

### ✅ Batch 8: Phase 6 - User Story 4 (Cluster Services)
**Issues**: #48-53 (6 tasks)
**Duration**: 1-2 days
**Status**: CAN RUN IN PARALLEL WITH BATCH 7 (after Batch 6)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 8 section

**Dispatch Command** (after Batch 6 merged, parallel with Batch 7):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 8 section)
GitHub Issues: #48-#53
Branch: batch-8-phase-6-us4-cluster-services
```

---

### ✅ Batch 9: Phase 7 - Documentation & Testing
**Issues**: #54-57 (4 tasks)
**Duration**: 1 day
**Status**: WAITING (depends on Batches 7 & 8 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 9 section

**Dispatch Command** (after Batches 7 & 8 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 9 section)
GitHub Issues: #54-#57
Branch: batch-9-phase-7-documentation
```

---

### ✅ Batch 10: Phase 8 - Deployment (Docker/K8s)
**Issues**: #58-63 (6 tasks)
**Duration**: 1-2 days
**Status**: WAITING (depends on Batch 9 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 10 section

**Dispatch Command** (after Batch 9 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 10 section)
GitHub Issues: #58-#63
Branch: batch-10-phase-8-deployment
```

---

### ✅ Batch 11: Phase 9 - Final Validation
**Issues**: #64-68 (5 tasks)
**Duration**: 1 day
**Status**: WAITING (depends on Batch 10 ✅)
**References**: `.specify/AGENT_BRIEFS.md` - Batch 11 section

**Dispatch Command** (after Batch 10 merged):
```
Use the agent brief from .specify/AGENT_BRIEFS.md (Batch 11 section)
GitHub Issues: #64-#68
Branch: batch-11-phase-9-final-validation
```

---

## Parallel Execution Opportunities

### Wave 1: Sequential (Critical Path)
```
Batch 1 (Research) → Batch 2 (Foundation) → Batch 3 (Swagger)
```

### Wave 2: Parallel Batches 4 & 5 (after Batch 3)
```
Batch 4 (US1) ↱
              → Batch 6 (Performance)
Batch 5 (US2) ↰
```

### Wave 3: Parallel Batches 7 & 8 (after Batch 6)
```
Batch 7 (US3) ↱
              → Batch 9 (Testing)
Batch 8 (US4) ↰
```

### Wave 4: Sequential Final Steps
```
Batch 9 (Testing) → Batch 10 (Deployment) → Batch 11 (Validation)
```

---

## Estimated Timeline

| Execution Model | Duration | Notes |
|-----------------|----------|-------|
| **Sequential (1 agent)** | 14-25 days | All batches run one after another |
| **Wave-based (2-3 agents)** | 7-10 days | Batches run in parallel where possible |
| **Optimal (4 agents)** | 5-7 days | All waves run simultaneously |

---

## How to Dispatch Each Batch

### Step 1: Copy the Agent Brief
Go to `.specify/AGENT_BRIEFS.md` and find the corresponding batch section.

### Step 2: Create a GitHub Issue for the Agent
Use the template below:

```
Title: [DISPATCH] Batch X - Phase Y - [Description]

Body:
[Paste the full agent brief from AGENT_BRIEFS.md]

GitHub Issues to Complete:
- Linked issue #XX
- Linked issue #YY
- ...

When complete:
1. Create PR to batch-X-phase-Y-description branch
2. All tests passing
3. Coverage >80%
4. Ready for review
```

### Step 3: Assign to Agent
Assign the dispatch issue to the Copilot agent for the batch.

### Step 4: Monitor Progress
- Watch for PR creation
- Review test results
- Approve and merge when ready

### Step 5: Proceed to Next Batch
After PR is merged to main, proceed to next batch (respecting dependencies).

---

## Important Notes

✅ **All issues created**: Issues #27-77 are in GitHub
✅ **Briefs ready**: All 11 agent briefs in `.specify/AGENT_BRIEFS.md`
✅ **TDD mandatory**: All test tasks must be written FIRST
✅ **Coverage requirement**: 80%+ code coverage required
✅ **Constitution compliance**: Go 1.24 standards, graceful shutdown, structured logging

---

## Quick Reference: Issue Numbers

| Batch | Phase | Issues | Count | Status |
|-------|-------|--------|-------|--------|
| 1 | 0 | #27-30 | 4 | ✅ Ready |
| 2 | 1 | #31-39 | 9 | ✅ Ready |
| 3 | 1.5 | #40-43 | 4 | ✅ Ready |
| 4 | 2 (US1) | (in deps) | - | Pending |
| 5 | 3 (US2) | #69-74 | 6 | ✅ Ready |
| 6 | 4 | #37-41 | 5 | ✅ Ready |
| 7 | 5 (US3) | #42-47 | 6 | ✅ Ready |
| 8 | 6 (US4) | #48-53 | 6 | ✅ Ready |
| 9 | 7 | #54-57 | 4 | ✅ Ready |
| 10 | 8 | #58-63 | 6 | ✅ Ready |
| 11 | 9 | #64-68 | 5 | ✅ Ready |

**Total**: 50 tasks across 11 batches

**Note on missing Phase 2 (US1)**: Issue #44-50 (US1 - Device Status) will be created as needed when agents complete Batch 3, as they're referenced in dependencies for subsequent batches.

---

## Start Dispatching! 🚀

Ready to begin? Start with **Batch 1** using the brief in `.specify/AGENT_BRIEFS.md`!

All issues are created and ready. Each agent will have clear, detailed instructions with:
- Acceptance criteria
- Testing requirements  
- Git workflow
- Manual testing examples
- Review checklist
- All dependencies documented

Good luck! 💪

