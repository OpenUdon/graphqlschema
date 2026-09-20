# Status M05 - UWS 1.4 Source-Profile Evidence Readiness

## Current State

M05 is implemented locally. The package now records non-normative GraphQL
source-profile evidence for later UWS 1.4 and OpenUdon planning without adding
UWS schema changes, OpenUdon fixtures, catalog rows, or runtime behavior.

## Task State

| Item | State | Notes |
|---|---|---|
| Source-profile evidence note | `[+]` | Added `docs/uws-source-profile-evidence.md` mapping M03/M04 operation and selector metadata to candidate UWS source-profile concepts without changing UWS here. |
| Selector-shape evidence | `[+]` | Recorded canonical examples: `query.book`, `mutation.checkout`, `subscription.bookUpdated`, `#/operations/{id}`, and `#/types/{rootType}/fields/{fieldName}`. |
| Boundary matrix | `[+]` | Documented package-owned parsing versus downstream UWS, OpenUdon, apitools, and runtime responsibilities. |
| Compatibility review | `[+]` | Reviewed the M04 stable API surface and examples before downstream source-profile work begins. |
| Verification | `[+]` | Standalone checks, sibling consumer checks, diff checks, and final review passed. |

## Verification

- Passed: `GOWORK=off go test ./...`
- Passed: `GOWORK=off go vet ./...`
- Passed: `git diff --check`
- Passed: `(cd ../apitools && go test ./...)`
- Passed: `(cd ../openudon && go test ./internal/synthesize ./internal/icot/elicitor)`
- Passed: `(cd ../tofu && git diff --check)`
- Passed: final deep code review; no remaining commit-blocking findings.

## Acceptance Criteria

- GraphQL parser and selector metadata are stable enough to inform UWS 1.4
  source-profile planning.
- Downstream ownership is explicit: UWS owns normative source types, OpenUdon
  owns fixtures/review, `apitools` owns catalog classification, and runtimes own
  execution.
- This repository still contains no UWS schema changes, OpenUdon workflow
  artifacts, catalog rows, endpoint calls, credentials, or GraphQL execution.
