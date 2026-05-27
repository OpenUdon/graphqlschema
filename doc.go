// Package graphqlschema provides metadata-only parsing support for GraphQL
// schema artifacts.
//
// Parse, ParseSDL, ParseIntrospection, and ParseIntrospectionMap convert
// GraphQL SDL or introspection JSON into native metadata for source-aware
// authoring and review tools. The package must not execute GraphQL operations,
// contact GraphQL servers, fetch remote schemas, or resolve credentials.
package graphqlschema
