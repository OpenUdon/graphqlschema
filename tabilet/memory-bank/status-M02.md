# Status M02 - SDL And Introspection Parser Bootstrap

## Current State

M02 is implemented locally. The package parses GraphQL SDL and introspection JSON
into native metadata while preserving the no-execution boundary.

## Task State

| Item | State | Notes |
|---|---|---|
| SDL parser | `[+]` | Added `Parse` and `ParseSDL` using `github.com/vektah/gqlparser/v2`. |
| Introspection parser | `[+]` | Added `ParseIntrospection` and `ParseIntrospectionMap` for `data.__schema` and root `__schema` forms. |
| Public model | `[+]` | Added metadata types for roots, types, fields, input values, enum values, directives, directive uses, and type refs. |
| Fixtures and malformed-input tests | `[+]` | Added representative SDL/introspection fixtures and structural error tests. |
| Documentation | `[+]` | Updated README, package docs, product, architecture, and tech-stack notes. |
| Verification | `[+]` | Standalone tests, vet, and diff checks pass. |

## Verification

- Passed: `GOWORK=off go test ./...`
- Passed: `GOWORK=off go vet ./...`
- Passed: `git diff --check`
- Passed: `(cd ../apitools && go test ./...)`
- Passed: `(cd ../openudon && go test ./internal/synthesize ./internal/icot/elicitor)`
- Passed: `(cd ../tofu && git diff --check)`
