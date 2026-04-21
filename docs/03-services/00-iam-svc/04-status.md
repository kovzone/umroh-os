# iam-svc — Status

## Implementation checklist

- [~] Scaffolded from baseline template (pilot — hybrid REST+gRPC binary)
- [~] Wired into `docker-compose.dev.yml`
- [ ] F1.2 DDL for `iam.users`, `iam.roles`, `iam.permissions`, `iam.role_permissions`, `iam.user_roles`, `iam.branches`, `iam.sessions`, `iam.audit_logs`
- [ ] sqlc queries for users, roles, permissions
- [ ] OpenAPI spec for auth endpoints
- [ ] Auth handlers (login, refresh, logout, me)
- [ ] gRPC service for `ValidateToken`, `CheckPermission`
- [ ] 2FA (TOTP) flow
- [ ] Audit log write path
- [ ] User CRUD endpoints (admin)
- [ ] Role CRUD endpoints
- [ ] Branch CRUD endpoints
- [ ] Unit tests for service layer
- [ ] Integration tests for auth flow
- [ ] Verified by reviewer in `testing-guide.md`

## Current status

**Scaffolded** — service lives at `services/iam-svc/`; REST on `4001`, gRPC on `50051`; shares the `umrohos_dev` database with every other service (ADR 0007). The `iam` Postgres schema namespace is created by the first F1.2 migration; until then iam-svc only touches the shared `public.diagnostics` table via `/system/diagnostics/db-tx`.

Scope-per-ADR-0006: `broker-svc` + Temporal containers are commented-out in `docker-compose.dev.yml` behind an `# F6 — see ADR 0006` marker.

**Pilot deliverables (this scaffold only):**

- `GET /system/live` — liveness probe.
- `GET /system/ready` — readiness probe (`SELECT 1` against `umrohos_dev`).
- `GET /system/diagnostics/db-tx` — WithTx reference: inserts a row into `public.diagnostics` (stamped `service='iam-svc'`), reads it back inside the same transaction, returns the ID.
- `pb.IamService/Healthz` gRPC placeholder + standard gRPC health protocol (`grpc.health.v1.Health`) with `iam.v1.IamService` registered as SERVING.

**Pilot strip (vs. the baseline `demo-svc`):**

Removed — baseline reference code not applicable to a minimal scaffold:

- `adapter/broker_grpc_adapter/` — depends on `broker-svc`, deferred by ADR 0006.
- `adapter/demo_grpc_adapter/` — iam-svc doesn't call another service; it's the one everyone calls.
- `api/rest_oapi/scenarios.go` + `/scenarios/*` routes + scenarios schemas in `openapi.yaml` — Temporal demo scenarios, deferred with broker-svc.
- `util/token/`, `service/auth.go`, `api/rest_oapi/auth.go`, `api/rest_oapi/middleware/bearer_auth.go`, `/auth/token` + `/auth/me` routes, token config fields — iam-specific auth machinery that belongs to F1.5, not scaffold scope.
- Service-layer unit-test scaffolding (`*_test.go` under `service/`) — stale w.r.t. the trimmed interface; real tests come with F1.14.

**Path convention:** Per ADR 0004 all Go microservices live under `services/` at the repo root. `baseline/go-backend-template/` stays untouched as the reference.

**Next:** F1.2 — the first real iam migration (`add_iam_users_and_roles`) creating `CREATE SCHEMA iam` and the tables listed above.

## Assigned ports

| Surface            | Port    |
|--------------------|---------|
| REST (Fiber)       | `4001`  |
| gRPC (IamService)  | `50051` |

## Pilot reference for the sweep session

- **Hybrid binary** — one process exposes Fiber REST and a gRPC server side-by-side.
- **ADR 0006 enforced in compose** — Temporal + broker-svc blocks commented out with reactivation marker.
- **Single shared database (ADR 0007)** — `umrohos_dev`, created by `POSTGRES_DB` in compose, schema applied by `make migrate-up`. No `_init/` directory.
- **Shared `public.diagnostics`** — every service's `/system/diagnostics/db-tx` writes into the same table, stamped with its app name. Per-service schemas (`iam.*`, `catalog.*`, ...) land in feature-slice migrations, not in the scaffold.
- **Task-named migrations** — `000002_scaffold_services` captures the scaffolding commit; later migrations follow the same task-oriented naming (`000003_add_iam_users_and_roles`, not `000003_create_users_table`).

## 2026-04-21 — S0-J-05 OpenTelemetry baseline fix

- `cmd/server.go` — `app.Use(otelfiber.Middleware(...))` wired as the first middleware after CORS so inbound `traceparent` is extracted and handler spans continue an upstream trace.
- `util/tracing/tracing.go` — `otel.SetTextMapPropagator(NewCompositeTextMapPropagator(TraceContext{}, Baggage{}))` set globally so otelhttp outbound / otelfiber inbound share the W3C propagator.
- `api/rest_oapi/system.go` — `DbTxDiagnostic` handler now starts its span from `c.UserContext()` (otelfiber's inbound-span context) instead of `c.Context()`, so handler spans correctly inherit the trace.
- `go.mod` — added `github.com/gofiber/contrib/otelfiber/v2 v2.2.3`.
