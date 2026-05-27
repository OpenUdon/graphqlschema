# GraphQL Schema

Dependency-light Go metadata package for GraphQL schema artifacts.

`github.com/OpenUdon/graphqlschema` parses GraphQL SDL and introspection JSON
into native metadata for downstream authoring, packaging, validation, and
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

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/OpenUdon/graphqlschema"
)

func main() {
	data, err := os.ReadFile("schema.graphqls")
	if err != nil {
		panic(err)
	}

	model, err := graphqlschema.Parse(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("query root %q with %d types\n", model.QueryType, len(model.Types))
	if op, ok := model.OperationByID("query.book"); ok {
		fmt.Printf("%s returns %s\n", op.ID, op.Type)
	}
}
```

## API

- `Parse` auto-detects SDL versus introspection JSON.
- `ParseSDL` parses GraphQL schema definition language.
- `ParseIntrospection` parses JSON introspection responses.
- `ParseIntrospectionMap` parses an already-decoded introspection response.
- `TypeByName` looks up parsed type metadata by GraphQL type name.
- `OperationByID` looks up selectable root fields such as `query.book`.
- `SelectorAliases` lists canonical operation IDs and local JSON Pointer
  selectors.
- `ResolveSelector` resolves operation IDs, `#/operations/{id}`, and
  `#/types/{rootType}/fields/{fieldName}` pointers.

## Consumers

See [Consumer Readiness](docs/consumer-readiness.md) for `apitools`, OpenUdon,
and UWS source-profile planning boundaries.

## Verification

```bash
GOWORK=off go test ./...
GOWORK=off go vet ./...
git diff --check
```
