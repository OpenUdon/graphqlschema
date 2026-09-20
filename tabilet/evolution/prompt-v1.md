# Prompt v1

Create a standalone `github.com/OpenUdon/graphqlschema` Go module, matching the
engineering harness pattern used by `../asyncapi`, `../googlediscovery`, and
`../awssmithy`.

The package should become a metadata-only parser for GraphQL SDL and
introspection JSON. It should expose schema, operation, variable, and selector
metadata for OpenUdon/apitools without executing GraphQL operations.
