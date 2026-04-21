# iam-svc — Status

## Implementation checklist

- [x] Scaffolded from baseline template (pilot — hybrid REST+gRPC binary)
- [x] Wired into `docker-compose.dev.yml`
- [x] Migration `000003_add_iam_users_and_roles` — `iam` schema + 8 tables + append-only audit trigger (landed via BL-IAM-001 scrape from abandoned `feat/f1-min`)
- [x] sqlc queries for users / roles / permissions / branches / sessions / audit_logs / join tables (8 query files, 48 queries)
- [x] OpenAPI spec for auth endpoints (6 paths — landed BL-IAM-001 2026-04-21)
- [x] Auth handlers (login, refresh, logout, me) — BL-IAM-001
- [x] Bearer-token middleware (PASETO/JWT via `util/token`) — BL-IAM-001
- [x] 2FA (TOTP) enrollment + verify endpoints — BL-IAM-001. Login-time enforcement deferred to `S1-E-06`.
- [x] Seed migration `000004_seed_initial_admin` — HQ branch + super_admin role + admin user (dev password)
- [x] Unit tests for service-layer helpers + Logout / GetMe / VerifyTOTP — BL-IAM-001
- [x] e2e spec `tests/e2e/tests/02a-iam-svc-sessions.spec.ts` — BL-IAM-001
- [ ] gRPC service for `ValidateToken`, `CheckPermission`, `GetUser` — **BL-IAM-002** (`feat/s1-e-04-iam-middleware`)
- [ ] Admin suspend + revoke-all-sessions flow — **BL-IAM-003** (`feat/s1-e-04-iam-suspend`)
- [ ] Audit log write path + `RecordAudit` gRPC — **BL-IAM-004** (`feat/s1-e-04-iam-audit`)
- [ ] User CRUD endpoints (admin) — `S1-E-06` depth card
- [ ] Role CRUD endpoints — `S1-E-06`
- [ ] Branch CRUD endpoints — `S1-E-06`
- [ ] Login-time TOTP enforcement — `S1-E-06`
- [ ] Verified by reviewer in `testing-guide.md`

## Current status

**BL-IAM-001 landed (2026-04-21)** — iam-svc now ships the first real user-facing surface: internal login + refresh + logout + `/v1/me` + self-serve TOTP enrollment. REST on `4001`, gRPC on `50051`; shares the `umrohos_dev` database with every other service (ADR 0007). The `iam` schema + 8 tables + append-only audit trigger are created by migration `000003_add_iam_users_and_roles`; `000004_seed_initial_admin` seeds a dev-only `admin@umrohos.dev / password123` user in the `HQ` branch with the `super_admin` role.

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
- ~~`util/token/`, `service/auth.go`, `api/rest_oapi/auth.go`, `api/rest_oapi/middleware/bearer_auth.go`, `/auth/token` + `/auth/me` routes, token config fields~~ — **re-added in BL-IAM-001 (2026-04-21)** under the UmrohOS-specific paths `/v1/sessions*` + `/v1/me*`; the demo-svc `/auth/token` + `/auth/me` shapes were never a product requirement.
- ~~Service-layer unit-test scaffolding (`*_test.go` under `service/`) — stale w.r.t. the trimmed interface; real tests come with F1.14.~~ **Re-added in BL-IAM-001**: `service/auth_helpers_test.go` / `crypto_test.go` / `auth_test.go` / `me_test.go` + `internal/mocks/istore.go` + `api/rest_oapi/middleware/bearer_auth_test.go`.

**Path convention:** Per ADR 0004 all Go microservices live under `services/` at the repo root. `baseline/go-backend-template/` stays untouched as the reference.

**Next:** sibling cards under `S1-E-04`:

1. **`BL-IAM-002`** (`feat/s1-e-04-iam-middleware`) — `iam.v1.IamService/ValidateToken` + `CheckPermission` gRPC + consumer middleware in other services.
2. **`BL-IAM-003`** (`feat/s1-e-04-iam-suspend`) — admin suspend + revoke-all-sessions flow.
3. **`BL-IAM-004`** (`feat/s1-e-04-iam-audit`) — `RecordAudit` gRPC + state-changing handlers in iam-svc + booking-svc start writing audit rows.

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
