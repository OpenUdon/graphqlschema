# Architecture

## Package Shape

| Area | Responsibility |
|---|---|
| `doc.go` | Package documentation and boundary statement. |
| `model.go` | Public metadata model types. |
| `parse.go` | SDL and introspection parse entrypoints and SDL normalization. |
| `introspection.go` | Introspection JSON map parsing and structural validation. |
| `selectors.go` | Operation summaries, selector aliases, and local selector resolution. |
| `testdata/` | Representative SDL and introspection fixtures. |
| `docs/` | Consumer-readiness, UWS evidence, and downstream-boundary notes. |

## Planned Parse Flow

```text
GraphQL SDL or introspection JSON
  -> structural root validation
  -> schema metadata
  -> root query/mutation/subscription types
  -> object/input/enum/interface/union/scalar/directive summaries
  -> field and argument metadata
  -> query/mutation/subscription operation summaries
  -> selector aliases
  -> Model
```

The package should emit native GraphQL metadata rather than OpenAPI, UWS, or
runtime-shaped data. Downstream packages that need those views must implement
explicit adapters outside this module.

## Validation Strategy

Validation should be structural and parser-oriented first:

- malformed SDL or introspection JSON is rejected;
- unsupported root shapes produce deterministic errors;
- malformed or unknown introspection type-ref kinds are rejected;
- introspection JSON with trailing data is rejected;
- unknown GraphQL extensions are preserved where practical;
- operation execution, endpoint calls, and credential checks are excluded;
- remote introspection fetches are excluded unless a future milestone defines
  safe fetch policy.

## Dependency Strategy

SDL parsing uses `github.com/vektah/gqlparser/v2` to avoid hand-rolling
GraphQL syntax and type-system validation. Introspection JSON parsing is
implemented locally with the standard library so introspection responses remain
plain decoded metadata.

## Boundaries

- `apitools` may consume this package for catalog review, metadata summaries,
  and artifact classification using a thin local-artifact adapter.
- OpenUdon may consume this package after a UWS GraphQL source contract exists.
- Trusted runtimes may consume this package for metadata, but execution belongs
  outside this module.
- No package here may execute provider operations, resolve credentials, fetch
  remote schemas, or choose endpoints/accounts/runtimes.
## Harness Layout

Private planning uses permanent `status-<LANE><NN>.md` ledgers. `M` is the
default parser lane; another letter requires an explicit ownership and
dependency review. The unattended runner reads state from the second column of
`Item | State | Notes` tables. Candidate directions remain unnumbered.
