# Error Handling

Errors in UmrohOS are **transport-agnostic**. The service layer raises domain sentinel errors; the API layer maps them to HTTP/gRPC responses via middleware. Never let HTTP status codes leak into business logic.

## The pattern

1. Define sentinel errors in `<svc>/util/apperrors/` (mirror the baseline template):
   ```go
   var (
       ErrNotFound          = errors.New("not found")
       ErrAlreadyExists     = errors.New("already exists")
       ErrInvalidInput      = errors.New("invalid input")
       ErrUnauthorized      = errors.New("unauthorized")
       ErrForbidden         = errors.New("forbidden")
       ErrConflict          = errors.New("conflict")
       ErrPreconditionFail  = errors.New("precondition failed")
       ErrInternal          = errors.New("internal error")
       ErrUpstreamUnavail   = errors.New("upstream unavailable")
   )
   ```

2. Wrap with context using `fmt.Errorf`:
   ```go
   user, err := s.store.GetUserByEmail(ctx, email)
   if err != nil {
       if errors.Is(err, sql.ErrNoRows) {
           return Result{}, fmt.Errorf("user with email %q: %w", email, apperrors.ErrNotFound)
       }
       return Result{}, fmt.Errorf("get user by email: %w", err)
   }
   ```

3. The API layer middleware maps sentinels to status codes:
   - `ErrNotFound` → 404
   - `ErrInvalidInput` → 400
   - `ErrUnauthorized` → 401
   - `ErrForbidden` → 403
   - `ErrAlreadyExists`, `ErrConflict` → 409
   - `ErrPreconditionFail` → 412
   - `ErrUpstreamUnavail` → 503
   - everything else → 500

   For gRPC: same sentinels mapped via `apperrors.ToGRPC` to canonical gRPC codes.

## Rules

- **Never hardcode HTTP status codes** in handlers. Use `apperrors`.
- **Always wrap with `%w`** so `errors.Is` works up the chain.
- **Always provide context** in the wrap message — at minimum what you were trying to do, ideally the relevant ID.
- **Never log and return.** Log at one level (usually the service method that caught it). The middleware will log the final mapping. Double-logging clutters Loki.
- **Never expose internal error messages to API clients.** The middleware emits a generic message + error code; the full error goes only to logs.
- **Database errors stay in the store.** Wrap them as `ErrNotFound` (for `sql.ErrNoRows`), `ErrAlreadyExists` (for unique violation), or `ErrInternal` (for everything else) before returning to the service layer.

## Error response shape (REST)

```json
{
  "error": {
    "code": "not_found",
    "message": "User not found"
  }
}
```

The middleware produces this from the sentinel. Never hand-build error responses in handlers.

## Panic policy

- No panics in business logic.
- Panics in goroutines must be recovered with logging and span error annotation. The template's `util/monitoring` provides `RecoverPanic` middleware.
- A panic in startup is acceptable — fail fast.
