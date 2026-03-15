# Specification Analysis Report — Feature 002-mcp-server + README.md

**Analyst**: GitHub Copilot (speckit.analyze)  
**Date**: March 14, 2026  
**Scope**: Read-only cross-artifact consistency and quality analysis  
**Artifacts analysed**:
- `specs/002-mcp-server/spec.md`
- `specs/002-mcp-server/plan.md`
- `specs/002-mcp-server/tasks.md`
- `specs/002-mcp-server/quickstart.md`
- `specs/002-mcp-server/data-model.md`
- `specs/002-mcp-server/research.md`
- `specs/002-mcp-server/contracts/mcp-protocol.md`
- `specs/002-mcp-server/checklists/requirements.md`
- `README.md`
- `.github/agents/copilot-instructions.md`
- `cmd/api/main.go` (implementation ground truth)
- `Makefile` (implementation ground truth)

---

## Executive Summary

The `002-mcp-server` feature is **fully implemented and consistent** across `spec.md`, `tasks.md`, `quickstart.md`, `data-model.md`, `copilot-instructions.md`, and the actual source code. All 42 tasks are marked `[X]`. The architecture (single parallel-goroutine binary, `homelab://` URIs, no subcommand, no args in `.vscode/mcp.json`) is consistently expressed across all artifacts **except** for two artifacts that contain stale or contradictory content: `plan.md` + `research.md` (describe an old subcommand architecture never implemented) and `contracts/mcp-protocol.md` (states `.vscode/mcp.json` requires `args: ["mcp"]`). Additionally, `README.md` has a severely stale header and `CURRENT_STATUS.md` reflects an obsolete project state from ~March 1, 2026.

**Critical issues: 3 | High: 4 | Medium: 5 | Low: 4**

---

## Findings Table

| ID | Category | Severity | Location(s) | Summary | Recommendation |
|----|----------|----------|-------------|---------|----------------|
| I1 | Inconsistency | **CRITICAL** | `plan.md` Summary, Project Structure, Constitution Check VII; `research.md` R1, R6, R8 | `plan.md` Summary says *"Integrate as a subcommand: `homelab-api mcp`"*. `research.md` R1 Decision: *"Use a subcommand/flag to select the mode of operation"*, R6 shows `os.Args` `mcp` check, R8 mentions `make mcp-run` running `bin/homelab-api mcp`. The actual implementation (and `tasks.md`, `quickstart.md`, `spec.md` FR-008, `copilot-instructions.md`) runs **both servers concurrently — no subcommand**. The plan and research describe architecture that was explicitly **rejected and not implemented**. | Update `plan.md` Summary, Structure diagram (`cmd/api/main.go` comment), and Constitution VII entry to describe the parallel-goroutine architecture. Update `research.md` R1, R6, R8 to reflect the final decision (rejected subcommand; chose concurrent goroutines). |
| I2 | Inconsistency | **CRITICAL** | `contracts/mcp-protocol.md` L400–416 (IDE Configuration section) | The contract document states `.vscode/mcp.json` should have `"args": ["mcp"]` and JetBrains **Args** should be `"mcp"`. The actual `.vscode/mcp.json` (per `CURRENT_STATUS.md`, `quickstart.md`, `tasks.md` T037, `spec.md` SC-008, `copilot-instructions.md`, `README.md`) has **no `args`** — the MCP server starts automatically. This means a developer following the contract document will misconfigure their IDE, adding a `mcp` subcommand arg that is ignored (or worse, breaks the server if it were ever to be handled). | Update `contracts/mcp-protocol.md` IDE Configuration section: remove `"args": ["mcp"]` from the VS Code example and remove `"Args": "mcp"` from the JetBrains section. No args are needed. |
| I3 | Inconsistency | **CRITICAL** | `plan.md` L10–14 (Summary), `data-model.md` entity diagram | `plan.md` Summary says *"When invoked as `homelab-api mcp`, the application starts an MCP stdio server"*. The `data-model.md` entity diagram (bottom of file) also shows `if "mcp" subcmd: runMCP()` vs `else: runHTTPServer()` — a mutually-exclusive branch, not concurrent goroutines. Both contradict `spec.md` FR-008 (MUST start automatically alongside HTTP, no subcommand), `tasks.md` T034, the actual `cmd/api/main.go`, and every other artifact. | Update `data-model.md` entity diagram to show the concurrent goroutine pattern: `errgroup` launching both `srv.Run(port)` and `mcp.Run(ctx)` from a single `main()`. |
| I4 | Inconsistency | **HIGH** | `research.md` R8 (Makefile targets) | R8 says the Makefile will have `make mcp-run` running `bin/homelab-api mcp` and `make mcp-build` as an alias for `build`. The actual Makefile has **neither** `mcp-run` nor `mcp-build` targets. `make run` starts both HTTP + MCP. The R8 rationale references a quickstart that also doesn't use those targets. | Update `research.md` R8 to reflect the actual Makefile. `mcp-run` and `mcp-build` do not exist; `make run` is the single entry point. |
| I5 | Inconsistency | **HIGH** | `plan.md` Project Structure `cmd/api/main.go` note | The plan's Project Structure table says `main.go` is *"Modified: add 'mcp' subcommand dispatch"*. The actual change was: remove subcommand dispatch entirely, add concurrent goroutine launch via `errgroup`. | Update the file change description for `cmd/api/main.go` in `plan.md` to say *"Modified: launch HTTP server and MCP server as concurrent goroutines under shared errgroup context (no subcommand)"*. |
| I6 | Inconsistency | **HIGH** | `contracts/mcp-protocol.md` L4 (Transport header) | The contract header says **Binary**: `bin/homelab-api mcp` (integrated single binary, **subcommand mode**). This directly contradicts the implemented architecture (no subcommand). A developer reading this will attempt `bin/homelab-api mcp` which is not the correct invocation. | Update contract header: **Binary**: `bin/homelab-api` (integrated single binary, **concurrent HTTP + MCP mode — no subcommand needed**). |
| S1 | Staleness | **HIGH** | `README.md` L9–28 (Current Status section) | README header says *"~40% Complete"*, *"March 2, 2026"*, references *"11 Active Agents"*, *"3 PRs Ready to Merge"*, and phase completion from a March 2 snapshot. Feature 002-mcp-server is fully complete (all 42 tasks `[X]`, per `CURRENT_STATUS.md`). The status section is at least 12 days stale and describes a project state that no longer exists. | Update README Current Status section: set date to March 14, 2026; status to ~100% feature-complete for 002-mcp-server; remove agent/PR counts that are no longer relevant; reference `.github/docs/CURRENT_STATUS.md` for detail. |
| S2 | Staleness | **HIGH** | `README.md` L23–24 (Links), L27 | Links `MISSION_COMPLETE.md` and `AGENT_MONITORING_DASHBOARD.md` point to root-relative paths (`./MISSION_COMPLETE.md`, `./AGENT_MONITORING_DASHBOARD.md`). Both files have moved to `.github/docs/`. These are broken links. | Update both links to `.github/docs/MISSION_COMPLETE.md` and `.github/docs/AGENT_MONITORING_DASHBOARD.md`. |
| S3 | Staleness | **MEDIUM** | `README.md` L135–149 (Planned Features) | "Planned Features" section lists as `🔄` (in progress): *HomeAssistant device query endpoints, device control endpoints, Service discovery endpoint, MCP tool wrappers*. All four are now implemented and the `CURRENT_STATUS.md` confirms completion. Additionally, "CORS support preparation" is listed as current but CORS is fully implemented. | Move implemented items to "Current Features". Remove `🔄 MCP tool wrappers` entirely (replaced by the full MCP server section already present in README). Consider removing the Planned Features section or retitling it to reflect only truly future work (live HA integration, auth). |
| S4 | Staleness | **MEDIUM** | `README.md` L631–669 (Project Structure) | Project Structure tree is missing: `internal/mcp/` (entire package), `internal/services/` (new package), `.github/docs/` (ops/history docs), `.github/agents/` (agent instructions), `.vscode/mcp.json` (IDE config). It also still lists `service.yaml` under `deployments/k8s/` but that file is not confirmed present. | Add `internal/mcp/`, `internal/services/`, `.github/docs/`, `.github/agents/`, `.vscode/mcp.json` to the Project Structure tree. Verify `deployments/k8s/service.yaml` existence; if absent, remove it. |
| S5 | Staleness | **MEDIUM** | `README.md` L194–220 (HomeAssistant Endpoints section header) | The heading says *"HomeAssistant Endpoints (Planned)"* and describes endpoints as planned. These are now implemented. The `GET /api/v1/homeassistant/devices` endpoint is live. | Rename heading to remove "(Planned)". Update endpoint descriptions to reflect they are implemented with mock data. |
| S6 | Staleness | **MEDIUM** | `.github/docs/CURRENT_STATUS.md` | The document header says *"Last Updated: March 1, 2026"* and describes the project as `~40%` complete with agents dispatched, phases in-progress, etc. The top section correctly says `002-mcp-server COMPLETE`, but the bulk of the document (Lines 82–285) is a March 1 snapshot with stale phase/PR/agent data that is contradicted by the top section. | The "Feature 002 COMPLETE" top section is accurate. The remainder of the document (Phase 6–9 status, PR table, agent dispatch table) is from a prior run and should be archived or replaced. Recommend archiving the old content under a heading like `## Archive: Pre-002 Status (March 1, 2026)`. |
| U1 | Underspecification | **LOW** | `tasks.md` T037 | T037 says *"Create `.vscode/mcp.json`"* and is marked `[X]`. However, searching the workspace reveals **no `.vscode/mcp.json` file** exists in the repository. T042 says *"mark completed tasks"* and the `CURRENT_STATUS.md` says `.vscode/mcp.json` was created, but the file is absent. This is either a tracking error or the file was not committed. | Verify whether `.vscode/mcp.json` was created. If not, it should be created (content is fully specified in `quickstart.md`, `contracts/mcp-protocol.md`, and `README.md`). If created but not committed, commit it. This is a deliverable required by `SC-008` and `FR-009`. |
| U2 | Underspecification | **LOW** | `research.md` R1 | The *"Alternatives Considered"* section says *"Concurrent HTTP + MCP in one process: Rejected — MCP stdio requires exclusive control of stdout"*. But the **implemented solution IS concurrent HTTP + MCP** — it resolves the stdout conflict by directing HTTP logs to stderr via `slog` and using only stdout for MCP JSON-RPC (as stated in `spec.md` and `tasks.md`). The research doc records the wrong resolution. | Update R1 Alternatives Considered: the "Concurrent" option was not rejected — it is the chosen solution. The stdout conflict is resolved by routing all non-MCP output to stderr. Remove the subcommand rationale and document why concurrency was chosen over subcommand dispatch. |
| D1 | Duplication | **LOW** | `mcp/server.go` `Run()` + `cmd/api/main.go` | `internal/mcp/server.go`'s `Run()` function emits `slog.Info("mcp server started", "transport", "stdio")`. `cmd/api/main.go` also emits the same log line immediately before calling `internalmcp.Run()`. This log line will appear twice on startup. Minor duplication; no functional impact. | Remove the duplicate `slog.Info` from `internal/mcp/server.go`'s `Run()` since the caller (`main.go`) already logs the startup message. Or remove it from `main.go` and keep it in `Run()`. Prefer keeping it in `Run()` as it is the canonical location. |
| D2 | Duplication | **LOW** | `quickstart.md` + `README.md` MCP section | The MCP section in `README.md` (L382–457) and `quickstart.md` (full document) duplicate the same content: build command, run command, smoke test, VS Code config, JetBrains config, available resources table. This is intentional (one is the authoritative quickstart, one is the README summary), but the duplication may create drift. | This is acceptable duplication for discoverability. Add a note at the top of `quickstart.md` linking to `README.md` MCP section as the summary, or vice versa. No immediate action required. |

---

## Coverage Summary

### Functional Requirements → Task Coverage

| Requirement Key | Has Task? | Task IDs | Notes |
|-----------------|-----------|----------|-------|
| FR-001 — session handshake | ✅ | T011, T012, T016 | Fully covered |
| FR-002 — list data sources | ✅ | T013, T016 | Covered |
| FR-003 — read live data by name | ✅ | T017–T023 | Covered |
| FR-004 — discover actions | ✅ | T014, T016 | Covered |
| FR-005 — execute actions, structured result | ✅ | T024–T028 | Covered |
| FR-006 — prompt templates with arg substitution | ✅ | T029–T033 | Covered |
| FR-007 — reuse existing data sources, no duplication | ✅ | T003–T010 | Covered by shared providers |
| FR-008 — auto-start with HTTP, no subcommand | ✅ | T034, T011 | Implemented correctly; plan.md contradicts |
| FR-009 — launchable via Makefile | ✅ | T036 | Covered |
| FR-010 — structured errors for all failures | ✅ | T025, T027, T029, T030 | Covered |
| FR-011 — concurrent requests safely handled | ✅ | T028 | Race test included |
| FR-012 — device catalogue from shared mock store | ✅ | T003, T004, T008 | Covered |

All 12 functional requirements have ≥1 associated task. **Coverage: 100%.**

### Non-Functional Requirements → Task Coverage

| NFR Key | Has Task? | Task IDs | Notes |
|---------|-----------|----------|-------|
| SC-001 — handshake < 2s | ✅ | T035 (smoke test) | Integration test validates |
| SC-002 — resource reads < 500ms | ⚠️ Partial | T039 | No explicit performance test; covered by general test suite |
| SC-003 — tool execution < 1s | ⚠️ Partial | T039 | No explicit benchmark; validated by smoke test only |
| SC-004 — 10 concurrent requests | ✅ | T028 | Race test with 10 goroutines |
| SC-005 — ≥80% coverage | ✅ | T039 | 83.8% achieved |
| SC-006 — clean build from checkout | ✅ | T041 | Docker build verified |
| SC-007 — existing tests pass | ✅ | T039 | All existing tests pass |
| SC-008 — IDE config in docs | ⚠️ | T037 | `.vscode/mcp.json` marked done but file not found in workspace |

### Unmapped Tasks

All 42 tasks map to at least one requirement, user story, or success criterion. No unmapped tasks found.

---

## Constitution Alignment

No constitution violations found in the **implemented** code. The `plan.md` Constitution Check table itself is accurate (all rules pass) **for the as-built implementation**. However:

- **Rule VII** (`plan.md` L35): The note says *"Entry point in `cmd/api/main.go` (subcommand)"*. The actual entry point uses parallel goroutines, not subcommand dispatch. The rule passes for the implementation but the note is inaccurate.
- All other constitution checks (Go standards, testing ≥80%, observability via slog/stderr, graceful shutdown, `go.mod`, `internal/` structure, minimal deps) are verified correct by the implementation.

---

## Metrics

| Metric | Value |
|--------|-------|
| Total Functional Requirements | 12 |
| Total Non-Functional Requirements (Success Criteria) | 8 |
| Total Tasks | 42 |
| Tasks marked [X] (complete) | 42 |
| FR Coverage (≥1 task) | 12/12 = **100%** |
| NFR Coverage (≥1 task) | 6/8 = **75%** (SC-002, SC-003 lack dedicated perf tests) |
| Ambiguity Count | 0 (no vague/unmeasurable criteria found) |
| Duplication Count | 2 (minor, D1 and D2) |
| Inconsistency Count | 6 (I1–I6) |
| Staleness Issues Count | 6 (S1–S6) |
| Underspecification Count | 2 (U1, U2) |
| **Critical Issues** | **3** (I1, I2, I3) |
| High Issues | 4 (I4, I5, I6, S1, S2) |
| Medium Issues | 4 (S3, S4, S5, S6) |
| Low Issues | 4 (U1, U2, D1, D2) |
| Total Findings | 16 |

---

## Detailed Analysis Notes

### Architecture Divergence: plan.md / research.md vs. Implementation

This is the most significant finding. The `plan.md` was written **before** the implementation decision was finalised. The final architecture (concurrent goroutines) differs from what `plan.md` specifies (subcommand dispatch). This is **expected** in iterative development — the spec evolved. However, the plan was never updated to match.

The net result is:
1. `plan.md` is an inaccurate historical record that will mislead future developers or agents reading it to understand the architecture.
2. `research.md` R1 records "concurrent HTTP + MCP: Rejected" when it is actually the chosen solution.
3. `data-model.md` entity diagram shows a branching `if mcp subcmd` structure.
4. `contracts/mcp-protocol.md` IDE Configuration section instructs users to pass `args: ["mcp"]` which is **incorrect** — the binary requires no args.

The `contracts/mcp-protocol.md` issue (I2) is the most user-facing: a developer following that document to configure VS Code or JetBrains will add an `args: ["mcp"]` that is not honoured by the binary (it will be ignored since `main.go` no longer checks `os.Args`).

### tasks.md — Architecture Alignment

`tasks.md` is the **best-aligned** document. Its Overview section explicitly says:
> *"There is no subcommand — running `./bin/homelab-api` launches both modes concurrently via goroutines"*

And T034 says:
> *"remove any `os.Args` subcommand dispatch; instead... launch both the HTTP server and the MCP server as concurrent goroutines"*

This confirms `tasks.md` was written with or after the architecture correction and matches the implementation.

### .vscode/mcp.json — Missing File

T037 is marked `[X]` and `CURRENT_STATUS.md` lists `.vscode/mcp.json` under "What Was Built". However, no `.vscode/mcp.json` file exists in the workspace. The content is fully specified in three places (quickstart.md, contracts/mcp-protocol.md, README.md). This is either:
- The file was not committed (most likely — `.vscode/` is often gitignored), or
- The task was prematurely marked complete.

If `.vscode/` is in `.gitignore`, the file should still exist locally and the README note about it being "already included in the repository" would be incorrect. **Verify `.gitignore`.**

### CURRENT_STATUS.md — Structural Inconsistency

`CURRENT_STATUS.md` has a **split personality**: the top ~60 lines accurately reflect the completed 002-mcp-server feature (March 14, 2026), but Lines 70–285 are verbatim content from the March 1, 2026 status document describing an entirely different project state (agents dispatched, phases in-progress, PRs needing review). The two halves contradict each other. This creates significant confusion for any agent or developer reading the file linearly.

---

## Next Actions

### Before `/speckit.implement` on any follow-on feature

3 **CRITICAL** issues must be resolved:

1. **I1/I3** — `plan.md` and `data-model.md` describe a subcommand architecture that was not implemented. These are the specification-of-record for the feature; leaving them inaccurate poisons any future agent reading them for context.
   - **Action**: Update `plan.md` Summary, Structure section, Constitution Check VII note, and `data-model.md` entity diagram to reflect the parallel goroutine architecture.

2. **I2** — `contracts/mcp-protocol.md` IDE Configuration section instructs developers to add `args: ["mcp"]` which is incorrect and will cause misconfiguration.
   - **Action**: Remove `args: ["mcp"]` from both VS Code and JetBrains examples in `contracts/mcp-protocol.md`.

### High-priority follow-up (before PR review)

3. **S1/S2** — README.md has a stale header, stale completion percentage, and broken links to `.github/docs/` files.
   - **Action**: Update README Current Status section; fix both broken doc links.

4. **U1** — `.vscode/mcp.json` may not be committed despite T037 being marked `[X]`.
   - **Action**: Run `ls -la .vscode/` to confirm. If missing, create and commit the file. If `.vscode/` is gitignored, the README claim that it is "already included in the repository" is incorrect and must be updated.

### Low-priority / can proceed without

5. **S3–S6** (Medium): README Planned Features, Project Structure, and CURRENT_STATUS.md all have stale content that should be updated for accuracy but don't block implementation or review.

6. **I4/I5/U2** (Medium/Low): `research.md` R1, R6, R8 and `plan.md` file change descriptions reflect the old architecture. Worth correcting for historical accuracy but non-blocking.

7. **D1** (Low): Duplicate startup log line. Remove one instance.

### Suggested commands

```bash
# Check whether .vscode/mcp.json exists or is gitignored
ls -la .vscode/ 2>/dev/null || echo "no .vscode directory"
grep -i vscode .gitignore 2>/dev/null

# Verify the binary truly starts without args (no subcommand handled)
grep -n "os.Args" cmd/api/main.go
# Expected: no results (confirms subcommand dispatch was removed)

# Run the full test suite to confirm all 42 tasks' outputs are green
make test
```

---

## Remediation Scope (Summary)

| Priority | Files to Edit | Nature of Change |
|----------|---------------|-----------------|
| CRITICAL | `specs/002-mcp-server/plan.md` | Rewrite Summary, file change entry for `main.go`, Constitution VII note, Project Structure diagram to describe concurrent goroutine architecture |
| CRITICAL | `specs/002-mcp-server/contracts/mcp-protocol.md` | Remove `args: ["mcp"]` from IDE Configuration section (VS Code + JetBrains) |
| CRITICAL | `specs/002-mcp-server/data-model.md` | Update entity diagram — replace `if "mcp" subcmd` branch with concurrent goroutine diagram |
| HIGH | `README.md` | Update Current Status date/percentage; fix two broken `.github/docs/` links |
| HIGH | `specs/002-mcp-server/research.md` | Correct R1 (concurrent = chosen, not rejected), R6 (no `os.Args` check), R8 (no `mcp-run`/`mcp-build` targets) |
| MEDIUM | `README.md` | Update Planned Features → Current Features; update Project Structure tree |
| MEDIUM | `.github/docs/CURRENT_STATUS.md` | Archive stale March 1 content; retain only the accurate 002-complete section |
| LOW | `.vscode/mcp.json` | Create if missing; verify commit status |
| LOW | `internal/mcp/server.go` | Remove duplicate startup `slog.Info` |

---

*This report is read-only. No files were modified during analysis.*  
*Generated by speckit.analyze — March 14, 2026*

