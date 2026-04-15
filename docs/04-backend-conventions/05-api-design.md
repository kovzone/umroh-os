# API Design

OpenAPI 3 is the source of truth for REST. Protobuf is the source of truth for gRPC. Code generators produce the Go interfaces — handlers are written against those interfaces.

## REST (OpenAPI + oapi-codegen)

### Source of truth

`<svc>/api/rest_oapi/openapi.yaml`. Edit this file first; regenerate; then implement.

### Versioning

Every public API is versioned in the URL path: `/v1/users`, `/v2/users`. Bump the version on breaking changes only. Add new fields freely under the same version.

### URL conventions

- Plural nouns: `/v1/users`, `/v1/bookings`
- Kebab-case: `/v1/package-departures`, `/v1/virtual-accounts`
- Resource IDs in the path: `/v1/users/{id}`
- Sub-resources: `/v1/bookings/{id}/items`, `/v1/jamaah/{id}/documents`
- Filters as query params: `/v1/bookings?status=paid_in_full&branch_id=...`
- Pagination: `?limit=50&cursor=<opaque>`. Cursor-based, not offset (PSAK reports might be huge).

### Methods

| Method | Use |
|---|---|
| `GET` | Read. Idempotent. Cacheable. |
| `POST /resources` | Create. |
| `POST /resources/{id}/<action>` | Non-CRUD action (e.g. `/bookings/{id}/cancel`). Prefer this over PATCH for state transitions. |
| `PATCH /resources/{id}` | Partial update. JSON Merge Patch. |
| `PUT /resources/{id}` | Full replacement. Rare. |
| `DELETE /resources/{id}` | Soft delete (sets `deleted_at`). |

### Response shape

Success:
```json
{
  "data": { ... },
  "meta": { "next_cursor": "..." }
}
```

Error (built by middleware from `apperrors`):
```json
{
  "error": {
    "code": "not_found",
    "message": "User not found"
  }
}
```

### Status codes

The middleware maps `apperrors` sentinels — never set a status code in a handler. See `02-error-handling.md` for the mapping.

### Auth

- Public endpoints: explicitly tagged `security: []` in OpenAPI.
- Authenticated endpoints: bearer token via the `RequireBearerToken` middleware. Token is validated by calling `iam-svc.ValidateToken` via gRPC.
- Diagnostic endpoints (`/healthz`, `/metrics`): protected by `RequireDiagnosticKey` middleware (shared secret in config).

### Schema design

- Use `$ref` aggressively. Avoid duplicating message shapes.
- Required fields explicit on every schema.
- Use `format` for dates (`date`, `date-time`, `uuid`).
- Use `enum` for status fields. Match the Postgres enum values exactly.

### Codegen

```sh
cd <svc>/api/rest_oapi && oapi-codegen --config .oapi-codegen.yaml openapi.yaml
# or
make oapi-<svc>
```

Generated file: `api.gen.go`. **Never edit it.** If something is wrong, fix the spec and regenerate.

## gRPC (protobuf)

### Source of truth

Each service **owns** its own proto file at `<svc>/api/grpc_api/pb/<svc>.proto`. There is no shared root-level `proto/` directory — see ADR 0004 "Proto ownership".

When service A needs to call service B, A carries a **local copy** of B's proto under `A/adapter/B_grpc_adapter/pb/B.proto`. A generates client code from that copy and wraps the gRPC client in an adapter that hides proto types from A's service layer.

Edit the owning service's proto first; regenerate; then implement. For contract changes that affect consumers, update every consumer's adapter copy in the same commit — the "N+1 commit" is deliberate, not a workflow problem (see ADR 0004).

### Conventions

- Service name: `<Domain>Service` (e.g. `IamService`, `CatalogService`).
- RPC names: `PascalCase` verbs (`GetUser`, `CreateBooking`, `ListPackages`).
- Request/response messages: `<Method>Request`, `<Method>Response`. Even if empty, define them — don't reuse `google.protobuf.Empty`.
- Field names: `snake_case`.
- Use `field_mask` for partial updates if needed.
- **Common message types are scoped to the service that owns them.** If two services genuinely need the same message shape, each defines its own copy. Resist the urge to create a cross-service "commons" package — that's the shared-proto pattern we rejected in ADR 0004.

### When REST vs gRPC

- **REST:** any external-facing endpoint. The frontend talks REST.
- **gRPC:** any service-to-service call. Faster, type-safe, streaming-capable.

A service can expose both — the template's `iam-svc` does. The REST API serves the frontend; the gRPC API serves other backend services.

## Backwards compatibility

- Adding fields: safe.
- Removing fields: breaking. Bump version (REST) or rename (gRPC).
- Renaming fields: breaking. Don't.
- Changing types: breaking. Don't.
- Adding endpoints: safe.
- Adding required request fields: breaking. Make them optional with sane defaults.
