# Status M04 - Consumer Integration Readiness

## Current State

M04 is implemented locally. Consumer-facing examples and a readiness note now
document the stable parser/selector surface for `apitools`, OpenUdon, and later
UWS planning without adding downstream integration code in this repository.

## Task State

| Item | State | Notes |
|---|---|---|
| API stability review | `[+]` | Reviewed `Operation`, `SelectorTarget`, `OperationByID`, `SelectorAliases`, and `ResolveSelector` as the M04 consumer surface. |
| Package examples | `[+]` | Added compiled examples for SDL parsing, introspection parsing, operation lookup, and selector resolution. |
| apitools readiness | `[+]` | Added readiness guidance for a thin `apitools` adapter that classifies and summarizes local GraphQL schema artifacts as native GraphQL metadata. |
| OpenUdon readiness | `[+]` | Documented stable metadata for future OpenUdon fixture planning after UWS source-profile semantics are scoped. |
| Generic overlay boundary | `[+]` | Documented that consumers should not create generic `POST /graphql` OpenAPI overlays when native GraphQL metadata is available. |
| Verification | `[+]` | Standalone checks, sibling checks, diff checks, and final review passed. |

## Verification

- Passed: `GOWORK=off go test ./...`
- Passed: `GOWORK=off go vet ./...`
- Passed: `git diff --check`
- Passed: `(cd ../apitools && go test ./...)`
- Passed: `(cd ../openudon && go test ./internal/synthesize ./internal/icot/elicitor)`
- Passed: `(cd ../tofu && git diff --check)`
- Passed: final deep code review; no remaining commit-blocking findings.

## Acceptance Criteria

- `apitools` can consume this package without duplicating GraphQL parsing.
- OpenUdon has stable metadata contracts for future fixture planning.
- Consumer-facing examples document the intended no-execution usage.
- No catalog rows, UWS schema changes, OpenUdon fixtures, or runtime execution
  behavior are implemented in this repository.
