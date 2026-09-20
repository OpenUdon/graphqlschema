# Status M03 - Operation And Selector Metadata

## Current State

M03 is implemented locally. The package exposes selectable query, mutation, and
subscription root fields with canonical operation IDs and local selector
resolution, while preserving the no-execution boundary.

## Task State

| Item | State | Notes |
|---|---|---|
| M02 review follow-ups | `[+]` | Preserved SDL interface possible-type parity, rejected trailing introspection JSON tokens, and rejected malformed or unknown introspection type-ref kinds deterministically. |
| Operation summaries | `[+]` | Added query/mutation/subscription root-field summaries with canonical IDs shaped as `<kind>.<field>`. |
| Variable/input metadata | `[+]` | Exposed root-field argument metadata for request mapping review: type refs, defaults, descriptions, and directive uses. |
| Selector aliases | `[+]` | Supported canonical operation IDs, `#/operations/{id}`, and `#/types/{rootType}/fields/{fieldName}` fragments. |
| Selector resolution | `[+]` | Resolved selectors to operation/root-field metadata without executing GraphQL operations. |
| Fixtures and tests | `[+]` | Covered query, mutation, subscription, invalid selectors, SDL/introspection parity, and review follow-up hardening cases. |
| Documentation | `[+]` | Updated README, package docs, architecture, and product notes for selector APIs added in M03. |
| Verification | `[+]` | Standalone tests, vet, diff checks, sibling consumer checks, and final review passed. |

## Verification

- Passed: `GOWORK=off go test ./...`
- Passed: `GOWORK=off go vet ./...`
- Passed: `git diff --check`
- Passed: `(cd ../apitools && go test ./...)`
- Passed: `(cd ../openudon && go test ./internal/synthesize ./internal/icot/elicitor)`
- Passed: `(cd ../tofu && git diff --check)`
- Passed: final deep code review; no remaining commit-blocking findings.

## Acceptance Criteria

- Downstream tools can enumerate selectable GraphQL root operations.
- Canonical operation IDs and JSON Pointer aliases resolve deterministically.
- Root-field argument metadata is sufficient for request mapping review.
- No GraphQL execution, endpoint discovery, credential behavior, UWS generation,
  or OpenAPI lowering is introduced.
