# GraphQL Schema

Dependency-light Go metadata package for GraphQL schema artifacts.

`github.com/OpenUdon/graphqlschema` will parse GraphQL SDL and introspection
JSON into native metadata for downstream authoring, packaging, validation, and
review tools. The package is intentionally schema-first: it preserves GraphQL
types, fields, arguments, directives, operation roots, and selector metadata
without executing GraphQL operations.

## Install

```bash
go get github.com/OpenUdon/graphqlschema
```

## Scope

This package is metadata-only. It does not execute queries, mutations, or
subscriptions; contact GraphQL servers; resolve credentials; or fetch remote
schema references.

The initial module contains the project harness and package boundary. Parser
APIs will be added by later milestones tracked in `memory-bank/milestone.md`.

## Verification

```bash
GOWORK=off go test ./...
GOWORK=off go vet ./...
git diff --check
```
