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

Authentication is validated **once, at `gateway-svc`** (per ADR 0009). Backend services do not extract bearer tokens.

- `gateway-svc`'s middleware extracts `Authorization: Bearer`, calls `iam-svc.ValidateToken` via gRPC (through gateway's own `iam_grpc_adapter`), and forwards the validated request to the downstream backend's gRPC method on success. On failure: 401. On `iam-svc` unreachable: 502 fail-closed.
- Public endpoints: explicitly tagged `security: []` in the gateway OpenAPI spec.
- Authenticated endpoints: tagged with `bearerAuth` security scheme in the gateway OpenAPI spec.
- **Backends do NOT re-validate the bearer.** Per-service `RequireBearerToken` is removed. Backends trust that a gRPC call reached them via gateway in MVP; the gateway↔backend trust contract (signed header / mTLS) is scheduled as `BL-GTW-100` for a later slice.
- `iam-svc.CheckPermission` (authorization) and `iam-svc.RecordAudit` (audit log) still live at the backend's service layer — they are not authentication and they need business context.
- Diagnostic endpoints on the gateway (`/metrics`, `/livez`, `/readyz`): protected by `RequireDiagnosticKey` middleware (shared secret in config) if exposed publicly; otherwise internal network only.

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

Per ADR 0009:

- **REST lives only on `gateway-svc`.** It is the sole external-facing surface. Client apps (browser, mobile, B2B portals) talk REST to gateway and nowhere else.
- **All downstream services are gRPC-only.** Their `api/rest_oapi/` package is removed. Business and domain calls — both north-south (gateway → backend) and east-west (backend → backend) — are gRPC.
- `gateway-svc` proxies every client request to a downstream service via a `<svc>_grpc_adapter` that wraps the gRPC client and decodes the request into typed params.

**Admin/observability endpoints** on downstream services are not "REST API" in the sense this rule targets:
- Liveness / readiness use the standard `grpc.health.v1.Health/Check` protocol. `grpc_health_probe` is the canonical probe binary for docker-compose and Kubernetes.
- `/metrics` (Prometheus scrape) lives on a minimal admin HTTP endpoint on each backend — single handler, no OpenAPI spec, no business routes. Alternatively, services may push metrics via OTLP to the OTel Collector.
- `/system/diagnostics/db-tx` (WithTx trace verification) is a gRPC method (`DiagnosticsDbTx`). Gateway proxies the public REST route `/v1/<svc>/system/diagnostics/db-tx` to it.

**Every new client-facing endpoint must include its gateway side in the same card.** A backend gRPC method reachable from a browser requires its gateway REST route + `_grpc_adapter` proxy in the same branch/PR. Half-shipped is not shipped.

**Scaffolding:** new non-gateway services scaffold from `baseline/go-backend-template/demo-grpc-svc/` (gRPC-only). Only `gateway-svc` uses the REST shape (`baseline/go-backend-template/demo-svc/`).

## Backwards compatibility

- Adding fields: safe.
- Removing fields: breaking. Bump version (REST) or rename (gRPC).
- Renaming fields: breaking. Don't.
- Changing types: breaking. Don't.
- Adding endpoints: safe.
- Adding required request fields: breaking. Make them optional with sane defaults.
