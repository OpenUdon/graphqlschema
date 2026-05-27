# Consumer Readiness

`graphqlschema` is ready for source-aware consumers that need local GraphQL
schema metadata without executing GraphQL operations.

For UWS 1.4 planning evidence, see
[UWS Source-Profile Evidence](uws-source-profile-evidence.md).

## Stable API Surface

The current consumer-facing surface is:

- `Parse`, `ParseSDL`, `ParseIntrospection`, and `ParseIntrospectionMap` for
  local metadata parsing.
- `Model`, `TypeDefinition`, `FieldDefinition`, `InputValueDefinition`,
  `EnumValueDefinition`, `DirectiveDefinition`, `Operation`, `SelectorTarget`,
  and `TypeRef` for native GraphQL metadata.
- `TypeByName`, `OperationByID`, `SelectorAliases`, and `ResolveSelector` for
  lookups and source-aware binding.

Selectors are intentionally local and metadata-only:

- canonical operation IDs: `query.<field>`, `mutation.<field>`, and
  `subscription.<field>`;
- operation JSON Pointers: `#/operations/{id}`;
- root-field JSON Pointers: `#/types/{rootType}/fields/{fieldName}`.

## apitools Readiness

`apitools` can classify local `.graphql`, `.graphqls`, and introspection JSON
artifacts as native GraphQL source metadata by calling this package directly.
A thin adapter should:

- parse local bytes with `Parse`;
- report `SourceKind`, root operation type names, operation IDs, type counts,
  directive counts, and selector aliases;
- treat parse errors as artifact review diagnostics;
- keep GraphQL artifacts distinct from OpenAPI, Discovery, Smithy, AsyncAPI,
  and advisory human-doc overlays.

The adapter must not generate a generic `POST /graphql` OpenAPI overlay merely
to make GraphQL fit a REST-shaped catalog entry.

## OpenUdon Readiness

OpenUdon can use this package after UWS GraphQL source-profile semantics are
scoped. The stable metadata needed for fixture planning is present:

- source operation IDs for query, mutation, and subscription root fields;
- local JSON Pointer selectors to operation/root-field metadata;
- argument metadata for request mapping review;
- return type refs and field descriptions for review evidence.

OpenUdon should continue to own workflow artifacts, package review evidence,
and UWS validation wrappers. This package only parses and indexes schema
metadata.

## Boundary

This package does not:

- execute queries, mutations, or subscriptions;
- contact GraphQL servers or perform introspection requests;
- resolve endpoints, accounts, or credentials;
- generate UWS documents, OpenAPI documents, catalog rows, or OpenUdon fixtures;
- lower GraphQL semantics into a generic HTTP runtime call.
