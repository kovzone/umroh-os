# Coding Style

Go code in UmrohOS follows the standard Go style with a few project-specific extras. The baseline template is the canonical reference — when in doubt, copy from `baseline/go-backend-template/demo-svc/`.

## Standard Go style

- `gofmt` (always — usually via `goimports`)
- `go vet` clean
- Pass `golangci-lint` with the template's config
- Standard import grouping: stdlib, external, internal — separated by blank lines

## Naming

- **Files:** `snake_case.go` (e.g. `auth_handler.go`, `user_store.go`)
- **Packages:** lowercase, single word, no underscores or camelCase (e.g. `service`, `apperrors`, `tracing`)
- **Types:** `PascalCase`. Avoid `Type` suffix.
- **Interfaces:** `PascalCase` with `I` prefix only when there's a struct of the same name (`IService` + `Service`). Otherwise just the noun (e.g. `Maker`).
- **Receivers:** short — `s` for service, `st` for store, `srv` for server. Consistent across the file.
- **Test files:** `*_test.go`. Test functions: `Test_<Subject>_<Scenario>`.

## Method signatures

Service-layer methods always take `(ctx context.Context, params <Method>Params)` and return `(<Method>Result, error)`. Even if both structs are empty, define them — it makes refactoring safe and signature changes don't ripple. See `01-three-layer-architecture.md`.

```go
type CreateUserParams struct {
    Email string
    Name  string
}

type CreateUserResult struct {
    UserID string
}

func (s *Service) CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error) {
    // ...
}
```

## Comments

- Public types and functions get a Go-doc comment.
- Comments explain *why*, not *what*. The code is the *what*.
- Avoid commented-out code. Delete it; git remembers.

## Errors

- Use the `apperrors` sentinel pattern. Wrap with `fmt.Errorf("%w: …", apperrors.ErrNotFound)` or `errors.Is`/`errors.As`. Details in `02-error-handling.md`.

## Dependencies

- Don't add a new module to `go.mod` without an ADR.
- Prefer the standard library. Prefer what's already in the template.

## What you do NOT do

- No global state. No package-level mutable variables.
- No init() functions for business logic. Wire dependencies explicitly in `cmd/start.go`.
- No raw `panic` outside of unrecoverable startup errors.
- No `interface{}` / `any` parameters in service-layer methods.
- No `time.Now()` directly inside service methods — inject a clock for testability if you need to assert timing. (Acceptable for non-test-critical paths.)
