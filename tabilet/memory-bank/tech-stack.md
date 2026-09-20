# Tech Stack

## Runtime

- Go 1.25.x.
- `github.com/vektah/gqlparser/v2` for GraphQL SDL parsing and type-system
  validation.
- Standard-library JSON decoding for GraphQL introspection responses.

## Verification

| Command | Purpose |
|---|---|
| `GOWORK=off go test ./...` | Run tests independent of any parent workspace. |
| `GOWORK=off go vet ./...` | Run static checks independent of any parent workspace. |
| `git diff --check` | Catch whitespace and patch formatting issues. |
| `(cd ../apitools && go test ./...)` | Check catalog/tooling consumers after exported API changes. |
| `(cd ../openudon && go test ./internal/synthesize ./internal/icot/elicitor)` | Check OpenUdon source-aware consumers after exported API changes. |

## External References

- GraphQL SDL and introspection JSON are the behavior fixtures.
- The package must not execute introspection queries or fetch remote schemas.
- Callers that accept remote or user-uploaded documents should bound input
  bytes before invoking parse APIs.
## Harness Runner

```bash
../skills/harness/tackle-memory-bank-api-loop --model lane-audit .
```
