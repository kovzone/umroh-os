# Three-Layer Architecture

Every UmrohOS service follows three-layer clean architecture: **API → Service → Store**. The layers communicate via Go interfaces. This is non-negotiable.

## The layers

```
┌─────────────────────────────────────────────┐
│  API layer  (api/rest_oapi/, api/grpc_api/) │  ← Thin. Parses, calls, returns.
└────────────────────┬────────────────────────┘
                     │ calls IService
┌────────────────────▼────────────────────────┐
│  Service layer  (service/)                  │  ← Business logic. Orchestrates store + adapters.
└────────────────────┬────────────────────────┘
                     │ calls IStore, IAdapter
┌────────────────────▼────────────────────────┐
│  Store layer  (store/postgres_store/)       │  ← Type-safe SQL via sqlc.
└─────────────────────────────────────────────┘
```

Dependencies flow **downward only**. The store does not know the service exists. The service does not know the API exists.

## API layer

- **Path:** `<svc>/api/rest_oapi/` (REST) or `<svc>/api/grpc_api/` (gRPC)
- **Generated from:** `openapi.yaml` (REST) or `.proto` files (gRPC). Never hand-write the interface.
- **Responsibility:**
  - Parse the request (path params, query, body)
  - Validate at the transport level (required fields, basic types)
  - Call the service layer with a `Params` struct
  - Return the service `Result` as the response
  - Map errors via the `ErrorHandler` middleware (REST) or `apperrors.ToGRPC` (gRPC)
- **Forbidden in this layer:** business logic, direct DB calls, calling adapters, building Params from strings without going through the service layer.

## Service layer

- **Path:** `<svc>/service/`
- **Hand-written.** This is where the business logic lives.
- **Structure:**
  - `service.go` — defines `IService` interface, `Service` struct, `NewService` constructor that wires dependencies (store, adapters, config)
  - One file per domain area (e.g. `auth.go`, `users.go`, `roles.go`)
- **Each method:**
  1. Open a tracing span (`util/tracing.GetTracer(...).Start(ctx, "Service.MethodName")`)
  2. Get a trace-correlated logger (`logging.LogWithTrace(ctx, s.logger)`)
  3. Log entry with input params
  4. Validate business rules
  5. Call store / adapters (typically inside a `WithTx` if multiple writes)
  6. Log result or error
  7. Set span status, return
- **Pattern:** every method takes `(ctx, <Method>Params)` and returns `(<Method>Result, error)`. Even if Params/Result are empty.

## Store layer

- **Path:** `<svc>/store/postgres_store/`
- **Generated from sqlc** based on annotated SQL queries in `queries/`.
- **Structure:**
  - `schema/` — DDL (tables, types, functions). Mirrors `_init/<svc>_db/`.
  - `queries/<domain>.sql` — annotated queries grouped by domain.
  - `sqlc/` — generated Go code (do not edit).
  - `store.go` — defines `IStore` interface (wraps `sqlc.Querier` + `WithTx`), `Store` struct, `NewStore` constructor.
  - `tx.go` — `WithTx` helper with auto-retry on serialization/deadlock, tracing, logging.
- **Forbidden:** business logic, calling other services, anything not directly mapped to SQL.

## Adapters

- **Path:** `<svc>/adapter/<other_svc>_grpc_adapter/`
- **Purpose:** wrap calls to *other* services so the service layer never sees a proto type or a transport error.
- **Pattern:** define `IAdapter` interface in the adapter package; `Adapter` struct holds the gRPC client; methods take/return plain Go types.

## Wiring

- **`cmd/start.go`** is the only place that constructs concrete types. It builds: config → logger → tracer → DB pool → store → adapters → service → server.
- All dependencies are injected via constructors. Nothing reads from package globals.
- Dependency direction: API depends on `IService`, Service depends on `IStore` and `IAdapter`. Concrete types only meet at the wiring point.

## Why interfaces between every layer

1. **Testability.** Service tests use a `mocks.IStore` from `testify/mock`, never a real DB.
2. **Swappability.** If you ever replace Postgres with something else, only the store implementation changes.
3. **Clarity of contract.** The interface tells you exactly what one layer expects from another.

## Common mistakes

- ❌ Calling sqlc-generated methods from a handler. → Always go through service.
- ❌ Importing a proto package from the service layer. → Use an adapter.
- ❌ Returning HTTP status codes from a service method. → Return `apperrors` sentinels.
- ❌ Reading config inside a service method. → Inject what you need at construction time.
- ❌ Logging without trace correlation. → Always `logging.LogWithTrace(ctx, logger)`.
- ❌ Skipping `Params`/`Result` structs because the method "only takes a string". → Define them anyway.
