# Result v1

The module direction is a metadata-only GraphQL schema package:

- public module path `github.com/OpenUdon/graphqlschema`;
- planned SDL and introspection JSON parse entrypoints;
- native model for schema, root operation types, object/input/enum/directive
  summaries, fields, arguments, and descriptions;
- planned selector aliases and operation lookup for query, mutation, and
  subscription entry points;
- no GraphQL execution, endpoint discovery, credential behavior, or OpenAPI
  lowering as the primary contract.
