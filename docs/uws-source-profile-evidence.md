# UWS Source-Profile Evidence

This note records GraphQL source-profile evidence from `graphqlschema` for
later UWS 1.4 and OpenUdon planning. It is not a UWS specification, schema
change, validator rule, OpenUdon fixture, catalog entry, or runtime contract.

## Candidate Source Shape

A future UWS GraphQL source profile can treat the GraphQL schema artifact as
the source document and UWS operations as selectors into that document.

Candidate source description evidence:

| Concept | Evidence from this package |
| --- | --- |
| Source family | GraphQL SDL or introspection JSON parsed from local bytes. |
| Source operation ID | Canonical root-field IDs such as `query.book`, `mutation.checkout`, and `subscription.bookUpdated`. |
| Source operation ref | Local JSON Pointer selectors such as `#/operations/query.book` and `#/types/Query/fields/book`. |
| Request mapping shape | Root-field arguments exposed as `Operation.Arguments` with type refs, defaults, descriptions, and directive uses. |
| Result metadata | Root-field return type exposed as `Operation.Type`. |
| Review metadata | Root-field descriptions, directive uses, root type names, type counts, and selector aliases. |

The candidate UWS type name, JSON Schema validation rules, and normative
selector requirements remain UWS-owned decisions.

## Selector Evidence

Representative selectors from the package fixtures:

| Operation | Canonical ID | JSON Pointer aliases |
| --- | --- | --- |
| Query `book` field | `query.book` | `#/operations/query.book`, `#/types/Query/fields/book` |
| Mutation `checkout` field | `mutation.checkout` | `#/operations/mutation.checkout`, `#/types/Mutation/fields/checkout` |
| Subscription `bookUpdated` field | `subscription.bookUpdated` | `#/operations/subscription.bookUpdated`, `#/types/Subscription/fields/bookUpdated` |

Selector behavior is metadata-only. Resolving a selector returns parsed schema
metadata; it does not execute a query, mutation, or subscription.

## Boundary Matrix

| Area | Owner | Notes |
| --- | --- | --- |
| GraphQL SDL parsing | `graphqlschema` | Uses `gqlparser` locally and does not contact servers. |
| GraphQL introspection JSON parsing | `graphqlschema` | Parses exported JSON only; it does not run introspection queries. |
| Operation IDs and selector aliases | `graphqlschema` | Exposes stable local metadata for root query/mutation/subscription fields. |
| Provider catalog classification | `apitools` | Should classify local GraphQL artifacts as native GraphQL metadata and avoid generic `POST /graphql` OpenAPI overlays. |
| Public source description type and validation | `../uws` | Owns any future normative `sourceDescriptions[].type` value, selector rules, and schema/model changes. |
| Workflow fixtures, review, and packaging | OpenUdon | Owns UWS examples, eval fixtures, package evidence, and validation wrappers after UWS semantics are scoped. |
| Execution and credentials | Trusted runtime packages | Own GraphQL transport, endpoint selection, authentication, retries, and side effects. |

## Compatibility Review

The M2-M4 public API surface is stable enough for downstream planning:

- Parse entrypoints: `Parse`, `ParseSDL`, `ParseIntrospection`,
  `ParseIntrospectionMap`.
- Metadata types: `Model`, `Operation`, `SelectorTarget`, `TypeDefinition`,
  `FieldDefinition`, `InputValueDefinition`, `EnumValueDefinition`,
  `DirectiveDefinition`, `DirectiveUse`, `TypeRef`.
- Lookup and selector helpers: `TypeByName`, `OperationByID`,
  `SelectorAliases`, `ResolveSelector`.

Future changes should preserve these names and meanings unless a later
milestone explicitly records a compatibility break and verifies sibling
consumers.

## Non-Goals

This evidence does not add:

- UWS schema or Go model changes;
- OpenUdon workflow or eval fixtures;
- `apitools` catalog rows or artifact materialization;
- GraphQL endpoint discovery, introspection requests, query execution,
  credentials, or runtime transport behavior;
- OpenAPI lowering for GraphQL services.
