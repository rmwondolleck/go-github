---
description: |
  This agent integrates multiple open pull requests into a single epic/feature branch.
  It reads all open PRs, understands their context, identifies integration points
  and conflicts, then creates a consolidated mergeable branch with all changes.
timeout-minutes: 30
strict: true
on:
  workflow_dispatch:
  issues:
    types: [labeled]
  pull_request:
    types: [labeled]
permissions:
  contents: read
  issues: read
  pull-requests: read
tools:
  github:
    lockdown: false
    toolsets: [repos, pull_requests, issues]
safe-outputs:
  add-comment: {}
  create-pull-request: {}
concurrency:
  group: pr-integration-${{ github.repository }}
  cancel-in-progress: false
engine: copilot
---

# PR Integration Agent

## Tone and Style

{% include ".github/workflows/shared/mood.md" %}

## Purpose

Integrate multiple open pull requests into a single, consolidated epic/feature branch.
This eliminates the manual rework required to rebase and merge many PRs individually by
understanding the context, dependencies, and integration points across all open PRs.

## Trigger Conditions

This workflow can be triggered in three ways:

### Manual Dispatch (workflow_dispatch)
Process all open PRs in the repository and create a consolidated integration branch.

### Issue Label Trigger (issues: labeled)
**Only process** when the label added is:
- `epic-integration`
- `ready-for-integration`

**Skip and exit with noop** if:
- The label is anything other than `epic-integration` or `ready-for-integration`
- There are fewer than 2 open PRs to integrate
- An integration PR is already open (title starts with "Epic:")
- The triggering issue is labeled `blocked`, `wip`, `on-hold`, or `needs-discussion`

### PR Label Trigger (pull_request: labeled)
**Only process** when the label added is:
- `epic-integration`
- `ready-for-integration`

**Skip and exit with noop** if:
- The label is anything other than `epic-integration` or `ready-for-integration`
- There are fewer than 2 open PRs to integrate
- An integration PR is already open (title starts with "Epic:")
- The triggering PR is labeled `blocked`, `wip`, `on-hold`, or `needs-discussion`

## Task

When triggered, follow these steps to integrate all open pull requests:

### Step 1: Discover and Catalog Open PRs

1. List all open pull requests in the repository
2. For each PR, collect:
   - PR number, title, description, and labels
   - Head branch name
   - Files changed (using the PR files API)
   - The diff content for context
3. Filter out any PRs that are:
   - Draft PRs (unless they have significant progress)
   - PRs labeled `wip` or `do-not-integrate`
   - The integration PR itself (if re-running)
4. Sort PRs by creation date (oldest first) as a baseline order

### Step 2: Cross-Reference Open Issues

Before analyzing PR dependencies, gather context from open issues to ensure the integration respects issue priorities, relationships, and statuses:

1. **List all open issues** in the repository
2. For each open issue, collect:
   - Issue number, title, body, and labels
   - Assignees and milestone (if any)
   - Linked pull requests (from issue body references like `#123` or PR cross-links)
3. **Build a PR-to-Issue map**: For each PR being integrated, identify:
   - Which issue(s) the PR addresses (from PR body `Closes #X`, `Fixes #X`, or `Resolves #X`)
   - Whether the linked issue is still open or has been closed
   - The issue's labels (e.g., `bug`, `feature`, `enhancement`, `blocked`)
4. **Check issue-based constraints**:
   - **Skip PRs** whose linked issue is labeled `blocked`, `wip`, or `on-hold`
   - **Skip PRs** whose linked issue is labeled `question` or `needs-discussion` (requires human input first)
   - **Prioritize PRs** whose linked issue is labeled `bug` (bug fixes should integrate first)
   - **Prioritize PRs** whose linked issue has a milestone (milestone work takes precedence)
   - **Flag PRs with no linked issue** — these may be orphaned and need review before integration
5. **Check for issue dependencies**:
   - If an issue body mentions "depends on #X" or "blocked by #X", ensure the dependency PR is integrated first
   - If a blocking issue is still open with no PR, flag the dependent PR for manual review
6. **Update the integration candidate list**: Remove or reorder PRs based on issue context

Include the issue context in the integration plan output (Step 4).

### Step 3: Analyze Dependencies and Conflicts

Analyze the collected PR data to build a dependency and conflict map:

1. **File Overlap Analysis**: Identify which PRs modify the same files
2. **Import/Dependency Analysis**: For Go files, check if PRs add imports or dependencies that others rely on
3. **Logical Grouping**: Group PRs by feature area:
   - Foundation/infrastructure (middleware, server setup, models)
   - API endpoints (handlers, routes)
   - Testing (unit tests, integration tests)
   - Documentation (README, Swagger, configs)
   - Deployment (Docker, Kubernetes)
4. **Conflict Prediction**: Flag file pairs where multiple PRs modify the same lines
5. **Integration Order**: Determine the optimal merge order based on:
   - Dependencies (foundation before features, features before tests)
   - File overlap (less conflicting PRs first)
   - PR maturity (more complete PRs first)

### Step 4: Generate Integration Plan

Create a detailed integration plan as a comment on the triggering issue (or as a new issue if triggered via workflow_dispatch):

```markdown
### 🔀 PR Integration Plan

**PRs to integrate**: {count} open pull requests
**Target branch**: `epic/consolidated-{date}`
**Base branch**: `main`

#### Issue Context

| PR | Linked Issue | Issue Status | Issue Labels | Priority |
|----|-------------|-------------|-------------|----------|
| #{pr_num} | #{issue_num}: {issue_title} | Open | `bug` | 🔴 High |
| #{pr_num} | #{issue_num}: {issue_title} | Open | `feature` | 🟡 Normal |
| #{pr_num} | None (orphaned) | — | — | ⚪ Review needed |

**Skipped PRs** (due to issue status):
- PR #{num}: Linked issue #{issue_num} is labeled `blocked`
- PR #{num}: Linked issue #{issue_num} is labeled `needs-discussion`

#### Integration Order

| Order | PR | Title | Files Changed | Conflicts With |
|-------|-----|-------|---------------|----------------|
| 1 | #{pr_num} | {title} | {file_count} | None |
| 2 | #{pr_num} | {title} | {file_count} | PR #{conflict} |
| ... | ... | ... | ... | ... |

#### Dependency Graph

```mermaid
graph TD
    PR_A[PR #A: Core Models] --> PR_B[PR #B: Service Layer]
    PR_A --> PR_C[PR #C: Middleware]
    PR_B --> PR_D[PR #D: Handlers]
    PR_C --> PR_D
    PR_D --> PR_E[PR #E: Unit Tests]
    PR_D --> PR_F[PR #F: Integration Tests]
    PR_D --> PR_G[PR #G: Swagger Docs]
    PR_H[PR #H: Dockerfile] -.->|no code dependency| PR_D
```

#### Predicted Conflicts

<details>
<summary><b>View Conflict Details ({count} potential conflicts)</b></summary>

**{file_path}**:
- Modified by PR #{a} and PR #{b}
- Conflict type: {overlapping lines / structural change}
- Resolution strategy: {merge both / prefer newer / manual review needed}

</details>

#### Risk Assessment
- **Low risk**: {count} PRs with no file overlaps
- **Medium risk**: {count} PRs with non-overlapping changes in same files
- **High risk**: {count} PRs with potentially conflicting changes

I'll proceed with the integration now.
```

### Step 5: Create the Epic Branch

1. Create a new branch from `main` named `epic/consolidated-{YYYY-MM-DD}`
2. For each PR in the determined integration order:
   a. Read the full diff of the PR
   b. Apply the changes to the epic branch
   c. If conflicts arise:
      - Attempt automatic resolution based on context understanding
      - For Go files: ensure imports are merged, not duplicated
      - For test files: combine test functions, merge test tables
      - For config files: merge configurations additively
   d. Verify the changes compile (for Go: check syntax validity)
   e. Create a commit with message: `Integrate PR #{number}: {title}`
3. After all PRs are integrated:
   - Run a final consistency check across all integrated files
   - Ensure no duplicate imports, conflicting routes, or broken references
   - Add any necessary glue code (e.g., registering new routes in the server)

### Step 6: Create the Consolidated PR

Create a pull request from the epic branch to `main` with:

**Title**: `Epic: Consolidate {count} PRs into integrated feature branch`

**Body**:
```markdown
## 🏗️ Epic Integration

This PR consolidates **{count}** open pull requests into a single, integrated feature branch.

### Integrated PRs

| PR | Title | Status |
|----|-------|--------|
| #{num} | {title} | ✅ Integrated |
| #{num} | {title} | ⚠️ Integrated with conflict resolution |
| #{num} | {title} | ❌ Skipped (reason) |

### Integration Summary

**Total files changed**: {count}
**New files added**: {count}
**Files modified**: {count}

### Changes by Category

<details>
<summary><b>🔧 Infrastructure & Middleware</b></summary>

- {list of infrastructure changes from integrated PRs}

</details>

<details>
<summary><b>🌐 API Endpoints</b></summary>

- {list of endpoint changes}

</details>

<details>
<summary><b>🧪 Tests</b></summary>

- {list of test additions}

</details>

<details>
<summary><b>📦 Deployment</b></summary>

- {list of deployment changes}

</details>

### Conflict Resolutions

{describe any conflicts that were resolved and how}

### Post-Merge Cleanup

After this PR is merged, the following PRs can be closed:
{list of PR numbers that are fully integrated}

### Verification Steps

1. `go build ./...` - Verify compilation
2. `go test ./...` - Run all tests
3. `make swagger` - Regenerate Swagger docs (if endpoint changes)
4. `make docker` - Verify Docker build
5. Manual review of integrated changes

---
*Generated by PR Integration Agent* | [Workflow Run](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }})
```

### Step 7: Update Status

After creating the consolidated PR:

1. Comment on the triggering issue (if applicable) with results
2. Add a comment on each integrated PR referencing the epic PR:

```markdown
### 🔀 Integrated into Epic Branch

This PR has been integrated into the consolidated epic branch:

**Epic PR**: #{epic_pr_number}

Once the epic PR is reviewed and merged, this PR can be closed.

---
*PR Integration Agent* | [Workflow Run](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }})
```

## Context-Aware Integration Strategies

### Go Source Files (.go)

When integrating Go source files from multiple PRs:

1. **Imports**: Merge import blocks, deduplicate, and organize (stdlib first, then external)
2. **Package declarations**: Must match - flag error if they don't
3. **Init functions**: Combine if multiple PRs add init() functions
4. **Route registration**: Merge route registrations in server setup
5. **Middleware chain**: Preserve the project's middleware ordering as defined in the server setup (verify by reading the existing router configuration)
6. **Handler functions**: Add all handler functions, ensure no naming conflicts
7. **Models/Types**: Merge type definitions, check for field conflicts
8. **Interfaces**: Combine interface methods from different PRs

### Test Files (_test.go)

When integrating test files:

1. **Test functions**: Combine all test functions into the appropriate test file
2. **Table-driven tests**: Merge test case tables when testing the same function
3. **Test helpers**: Deduplicate shared test helpers
4. **Mock implementations**: Combine mock structs, ensure interface compliance
5. **TestMain**: Merge TestMain functions if multiple exist

### Configuration Files

When integrating configuration files (YAML, JSON, Dockerfile):

1. **Kubernetes manifests**: Merge environment variables, volume mounts additively
2. **Docker files**: Use the most comprehensive Dockerfile as base
3. **go.mod**: Combine all dependency additions, use highest version for conflicts
4. **Swagger docs**: Regenerate after all endpoint integrations

### Documentation Files

When integrating documentation (README.md, etc.):

1. **README**: Merge sections additively, maintain table of contents
2. **API docs**: Combine endpoint documentation
3. **Comments**: Preserve all meaningful comments from all PRs

## Conflict Resolution Priorities

When conflicts cannot be automatically resolved, use these priorities:

1. **Correctness**: Choose the version that is functionally correct
2. **Completeness**: Prefer the more complete implementation
3. **Recency**: If equal, prefer the more recent PR's approach
4. **Safety**: When in doubt, flag for manual review rather than guessing

## Safety Guidelines

1. **Never** force-push to `main` or any existing branch
2. **Never** close PRs automatically - only suggest closing after epic merge
3. **Never** delete branches
4. **Always** create a new epic branch (don't reuse existing ones)
5. **Always** preserve all functionality from all integrated PRs
6. **Always** flag unresolvable conflicts for human review
7. **Preserve** existing test coverage - never remove tests
8. **Document** every conflict resolution decision

## When to Request Human Help

Comment on the issue and **do not** proceed if:

- More than 30% of PRs have high-risk conflicts
- Core architectural decisions differ between PRs
- Security-sensitive code has conflicting implementations
- Database schema changes conflict
- The integration would require significant new code beyond merging

Use this comment format:

```markdown
### 🤚 Human Assistance Needed

I've analyzed the open PRs but need human help to proceed with integration.

**Reason**: {why automatic integration is not safe}

**Conflicting PRs**:
- PR #{a} vs PR #{b}: {description of conflict}

**Suggested Resolution**:
{recommended approach for human to take}

---
*PR Integration Agent* | [Workflow Run](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }})
```
