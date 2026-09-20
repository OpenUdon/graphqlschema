# Status M01 - Standalone Harness Bootstrap

## Current State

M01 is implemented locally. The standalone module, README, license, package
docs, planning harness, and symlinks are in place. Parser implementation is
deferred to M02.

## Task State

| Item | State | Notes |
|---|---|---|
| Engineering harness | `[+]` | Added `AGENTS.md`, memory-bank files, and evolution files under `../tofu/graphqlschema`, exposed through public-repo symlinks. |
| Module bootstrap | `[+]` | Added README, license, module file, and package docs. |
| Parser implementation | `[X]` | Deferred to M02. |
| Verification | `[+]` | Standalone tests, vet, and diff checks pass. |

## Verification

- Passed: `GOWORK=off go test ./...`
- Passed: `GOWORK=off go vet ./...`
- Passed: `git diff --check`
