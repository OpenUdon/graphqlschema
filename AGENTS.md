# AGENTS.md

## Purpose

`graphqlschema` is the OpenUdon-owned native metadata package for GraphQL
schema artifacts. It will parse GraphQL SDL and introspection JSON into
dependency-conscious Go metadata that downstream packages can use for
source-aware authoring, review, packaging, and validation.

Module path:

```text
github.com/OpenUdon/graphqlschema
```

This package treats GraphQL schemas and documents as untrusted metadata. It may
describe schema identity, object/input types, fields, arguments, directives,
queries, mutations, subscriptions, and operation-variable compatibility, but it
must not execute GraphQL operations, choose endpoints, resolve credentials, or
select production accounts.

## Start Here

Before substantial changes, read these in order:

1. [memory-bank/product.md](memory-bank/product.md)
2. [memory-bank/architecture.md](memory-bank/architecture.md)
3. [memory-bank/tech-stack.md](memory-bank/tech-stack.md)
4. [memory-bank/milestone.md](memory-bank/milestone.md)
5. The relevant permanent status file, such as
   [memory-bank/status-M01.md](memory-bank/status-M01.md)

Keep the memory bank short and module-specific. Do not copy catalog, CLI, UWS,
OpenUdon, or runtime-plan concerns from sibling repositories.

This project exposes [GOAL.md](GOAL.md), one optional protocol for goal requests
that span multiple status files. Follow it only when a request names it.

A `GOAL.md` run is a deliberate exception to the row-level commit rule below.
For that run, `COMMIT_POLICY: none` — the protocol default — means no commits,
while `COMMIT_POLICY: task` keeps the usual one-commit-per-row cadence.
Precedence is the request, then `GOAL.md`, then this file; only commits are
delegated, and only during the run.

## Boundary

- `../graphqlschema` owns native GraphQL schema parsing, introspection/SDL
  normalization, type and operation metadata, variable/input metadata, selector
  aliases, and selector resolution.
- `../apitools` owns catalog discovery, import, cache, refresh,
  classification, provider metadata, and source materialization that may
  reference GraphQL schema artifacts.
- `../uws` owns public workflow source-type semantics if GraphQL graduates into
  a UWS source description type.
- `../openudon` owns package artifacts, review evidence, source-bound workflow
  generation, validation wrappers, and trusted-runner handoff.
- `../udon` owns runtime lowering and any concrete GraphQL execution support.

Rule of thumb: if a change parses or normalizes GraphQL schema metadata, it
belongs here. If it executes GraphQL, chooses credentials, turns GraphQL into
UWS, or manages catalogs/packages, it belongs downstream.

## Essential Commands

```bash
GOWORK=off go test ./...
GOWORK=off go vet ./...
git diff --check
```

When changing exported APIs, run dependent checks in sibling consumers when
available:

```bash
(cd ../apitools && go test ./...)
(cd ../openudon && go test ./internal/synthesize ./internal/icot/elicitor)
```

## Hard Rules

- Treat all GraphQL schemas and operation documents as untrusted input.
- Keep the package dependency-light until a parser milestone justifies a
  dependency.
- Do not execute GraphQL operations or introspection requests.
- Do not resolve credentials, fetch tokens, sign requests, or cache secrets.
- Do not fetch remote schemas or external references unless a future milestone
  explicitly defines safe fetch policy.
- Preserve exported API compatibility unless the breaking change is intentional
  and documented.
- Run required verification before claiming a change is done.

## Work Cadence

- Update memory-bank files in the same change as implementation work they
  describe.
- Keep milestones small and parser-focused.
- Keep one permanent, zero-padded `status-<LANE><NN>.md` file for every
  indexed milestone. Never reuse an ID or create aggregate `status.md`.
- Keep later candidates unnumbered until fresh scope and dependency review
  promotes them.
- Write `Item | State | Notes` ledgers with a backticked state marker in the
  second column: `` `[ ]` ``, `` `[+]` ``, `` `[~]` ``, `` `[!]` ``, or
  `` `[X]` ``.
- Treat each row as a commit unit. Register another lane only for independent
  ownership, and document prerequisites and downstream impact before parallel
  work.
- Add evolution snapshots only when the public contract, package boundary, or
  parser strategy materially changes.
