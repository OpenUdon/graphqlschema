# Milestones

This file owns the roadmap, active milestone, permanent status-file index,
dependencies, candidates, scope, and acceptance criteria. Per-task state lives
in one `status-<LANE><NN>.md` file per milestone.

## Status ID Pattern

Use one uppercase domain letter and a zero-padded number from `01` through
`99`. `M` is the default lane for this focused parser. M06 explicitly
normalizes legacy M1-M5 filenames to M01-M05 without changing their identity.
Never reuse an ID, remove cancelled history, or create aggregate `status.md`.

Register another lane only for a genuinely independent domain. Lane letters
classify ownership rather than order; parallel work requires non-overlapping
files, resolved prerequisites, and named downstream impacts.

## Current Dashboard

Active milestones: none. Latest completed milestone: M06 - Parallel-Lane
Harness Migration.

## Status Files

| Milestone | Status File | State |
|---|---|---|
| M01 - Standalone Harness Bootstrap | [status-M01.md](status-M01.md) | Complete |
| M02 - SDL And Introspection Parser Bootstrap | [status-M02.md](status-M02.md) | Complete |
| M03 - Operation And Selector Metadata | [status-M03.md](status-M03.md) | Complete |
| M04 - Consumer Integration Readiness | [status-M04.md](status-M04.md) | Complete |
| M05 - UWS 1.4 Source-Profile Evidence Readiness | [status-M05.md](status-M05.md) | Complete |
| M06 - Parallel-Lane Harness Migration | [status-M06.md](status-M06.md) | Complete |

## Candidate Directions

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| Additional GraphQL schema fidelity | Broader directives and dialect-specific metadata have no current consumer requirement. | A named consumer supplies representative SDL/introspection fixtures and an acceptance gap. |
| Remote schema acquisition | Network fetching belongs in apitools and needs an explicit safety policy. | Apitools approves a bounded source-discovery contract that consumes this parser without moving network policy here. |

## M01 - Standalone Harness Bootstrap

**Goal.** Create a standalone GraphQL schema metadata module with a clear
engineering harness.

**Scope.**

- Create module `github.com/OpenUdon/graphqlschema`.
- Add README, LICENSE, package docs, AGENTS, memory-bank, and evolution files.
- Document the metadata-only boundary and initial milestone sequence.
- Verify standalone tests, vet, and diff checks.

**Acceptance.** The standalone module builds independently, documents its
metadata-only boundary, and is ready for parser implementation.

## M02 - SDL And Introspection Parser Bootstrap

**Goal.** Parse GraphQL SDL and introspection JSON into native metadata.

**Scope.**

- Add parse entrypoints for SDL and introspection JSON.
- Preserve schema metadata, root operation types, object/input/enum/directive
  summaries, fields, arguments, and descriptions.
- Add representative fixtures and malformed-input tests.

**Acceptance.** Common SDL and introspection fixtures parse deterministically
without executing queries or fetching remote schemas.

## M03 - Operation And Selector Metadata

**Goal.** Expose query, mutation, and subscription root fields as selectable
GraphQL operations for source-aware workflow tooling.

**Scope.**

- Add operation summaries for root fields on the parsed query, mutation, and
  subscription types.
- Use canonical operation IDs shaped as `<kind>.<field>`, such as
  `query.book`, `mutation.checkout`, and `subscription.bookUpdated`.
- Add selector aliases and local selector resolution for canonical operation
  IDs plus JSON Pointer fragments to root fields.
- Expose variable/input metadata from root-field arguments, including type
  refs, defaults, descriptions, and directive uses.
- Preserve SDL/introspection parity for root-field metadata and interface
  possible-type summaries before downstream consumers rely on selectors.
- Add parser-hardening follow-ups found in M02 review: reject trailing JSON
  tokens in introspection input and reject malformed or unknown type-ref kinds
  deterministically.

**Acceptance.** Downstream tools can select GraphQL root operations, resolve
their local selectors, and review variable mappings without OpenAPI lowering or
GraphQL execution.

## M04 - Consumer Integration Readiness

**Goal.** Prepare `apitools` and OpenUdon consumer integration after parser and
selector APIs stabilize.

**Scope.**

- Freeze exported metadata and selector APIs for M02/M03 behavior.
- Add package examples that show SDL and introspection parsing, operation
  summary lookup, and selector resolution.
- Add `apitools` integration planning notes or adapters for classifying local
  GraphQL schema artifacts as native GraphQL source metadata.
- Keep GraphQL-first providers out of generic `POST /graphql` OpenAPI overlays;
  consumers should use native GraphQL metadata once available.
- Verify sibling checks in `../apitools` and `../openudon`.

**Acceptance.** `apitools` can classify and summarize package-local GraphQL
schema artifacts without duplicating parser logic, and OpenUdon has enough
stable metadata to plan source-aware fixture work.

## M05 - UWS 1.4 Source-Profile Evidence Readiness

**Goal.** Produce stable GraphQL source-profile evidence for later UWS 1.4 and
OpenUdon planning without making normative UWS changes in this package.

**Scope.**

- Add release/readiness documentation that maps package metadata to candidate
  UWS GraphQL source-profile concepts.
- Record selector-shape evidence from representative query, mutation, and
  subscription schemas.
- Document what remains owned by downstream consumers: source description type
  naming, UWS validation rules, workflow request mapping, package fixtures, and
  runtime execution.
- Add compatibility checks for public API stability before downstream source
  profile work starts.
- Do not add UWS schema changes, OpenUdon workflow fixtures, catalog rows, or
  GraphQL execution behavior in this repository.

**Acceptance.** The package has stable parser and selector evidence for a UWS
1.4 GraphQL source-profile proposal, with clear downstream ownership and all
standalone and sibling verification passing. No UWS schema changes, OpenUdon
fixtures, catalog rows, or runtime execution behavior are added in this
repository.

## M06 - Parallel-Lane Harness Migration

**Goal.** Adopt the current lane-aware harness without changing GraphQL parser
behavior.

**Acceptance.** Historical IDs and task states are canonical and parseable,
candidates remain unnumbered, verification passes, and the unattended runner
finds no actionable row.
