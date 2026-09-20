# Product

`graphqlschema` provides a small Go package for parsing GraphQL schema artifacts
into native metadata.

## Users

- `apitools`, for catalog/import/review workflows that need to classify and
  inspect GraphQL schema artifacts.
- OpenUdon, for source-aware authoring, workflow synthesis, validation, package
  review, and handoff checks if GraphQL becomes a UWS source family.
- Trusted runtime packages that need GraphQL schema metadata without pulling in
  catalog tooling or OpenUdon workflow logic.

## Capabilities

- Implemented: parse GraphQL SDL and introspection JSON from bytes or decoded
  maps.
- Implemented: preserve schema identity, root operation types,
  object/input/interface/union/enum/scalar types, fields, arguments, directives,
  descriptions, and type refs.
- Implemented: summarize query, mutation, and subscription entry points for
  source-aware operation selection.
- Implemented: expose variable/input metadata for request mapping review.
- Implemented: resolve local selectors for schema operations and fields where a
  downstream UWS source contract needs stable references.
- Implemented: document consumer-readiness boundaries for `apitools`, OpenUdon,
  and later UWS source-profile evidence work.
- Implemented: record non-normative UWS 1.4 GraphQL source-profile evidence,
  selector examples, boundary ownership, and API compatibility notes.

## Non-Goals

- GraphQL operation execution.
- Endpoint discovery, account selection, or credential resolution.
- UWS document generation.
- UWS schema/model changes.
- OpenAPI lowering as the primary contract.
- Remote introspection fetches.
- Full GraphQL validation beyond parser/index needs until a milestone scopes it.

## Reference Model

GraphQL is a source family with semantics that should not be flattened into a
REST-shaped approximation by default. The schema owns object types, input types,
field arguments, nullability, enums, directives, query/mutation/subscription
roots, and selection constraints. Downstream workflow tools own selected
operations, variable values, data flow, review evidence, and trusted execution.
