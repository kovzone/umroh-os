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
- [x] gRPC service for `ValidateToken` + `CheckPermission` — **BL-IAM-002** on 2026-04-21 (`feat/s1-e-04-iam-middleware`). `GetUser` stays deferred; F1-W3 acceptance did not require it.
- [x] Admin suspend + revoke-all-sessions flow — **BL-IAM-003** on 2026-04-21 (`feat/s1-e-04-iam-suspend`). Includes the ride-along `/security-review` remediation from BL-IAM-002 and a defense-in-depth `users.status` guard in `ValidateToken`.
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

**Next:** one remaining sibling card under `S1-E-04`:

1. **`BL-IAM-004`** (`feat/s1-e-04-iam-audit`) — `RecordAudit` gRPC + state-changing handlers in iam-svc + booking-svc start writing audit rows. The suspend + self-suspend events captured by `service.SuspendUser` already carry `ActorUserID` + `TargetUserID` in-context; wiring the audit write is a one-line call inside the existing `WithTx` when that card lands.

Known follow-up (not blocking any card — tracked for S1-E-06 admin CRUD): iam-svc's REST bearer middleware `RequireBearerToken` still only verifies the PASETO signature, not the session row or `users.status`. Consumer services (finance-svc et al.) fail closed because their middleware calls `iam-svc.ValidateToken` over gRPC (which BL-IAM-003 hardened with the status guard), but iam-svc's own `/v1/me*`, `DELETE /v1/sessions`, and `POST /v1/users/{id}/suspend` routes accept any unexpired-signature bearer — so a suspended super-admin's still-valid access token can in principle reach those handlers until its TTL elapses. Fidelity gap, not a F1-W5 product-security gap (the PRD's "cannot access" is about downstream business surfaces, all of which are gated). Fix path when S1-E-06 lands: reshape `RequireBearerToken` to call `service.ValidateToken` in-process, same pattern finance-svc uses, and surface the identity envelope (not the raw PASETO payload) via `c.Locals`.

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

## 2026-04-21 — BL-IAM-002 internal gRPC permission resolution

The placeholder `pb.IamService` ships its first two real RPCs so downstream services can stop faking auth in tests. Shape matches the § Booking auth row in `docs/contracts/slice-S1.md` (permission strings as `resource` / `action` / `scope` tuples, e.g. `booking.create_on_behalf`); no new § IAM section needed in the slice contract — the wire is producer-owned.

- `api/grpc_api/pb/iam.proto` — added `ValidateToken(ValidateTokenRequest)` returning `user_id` / `branch_id` / `session_id` / `roles[]` / `expires_at_unix`, and `CheckPermission(CheckPermissionRequest)` returning `allowed`. Regenerated via the root `make genpb` (`iam.pb.go` + `iam_grpc.pb.go`).
- `api/grpc_api/server.go` — two new methods on `*Server` mirroring the `Healthz` pattern (tracer span, `LogWithTrace`, delegate to service). Errors map to gRPC codes via `apperrors.GRPCCode`.
- `service/permissions.go` (new) — `ValidateToken` verifies via `tokenMaker.VerifyToken` (PASETO v2 local), reloads the session row (rejects revoked / expired), and fetches role names fresh from the DB on every call so role changes propagate without waiting for the token to roll. `CheckPermission` validates the scope against the `iam.permission_scope` enum at the service boundary before the DB hit; `allowed=false` is a valid outcome (not an error).
- `store/postgres_store/queries/permissions.sql` — added `UserHasPermission :one` (EXISTS-based join over `user_roles × role_permissions × permissions`; pgx bound params, index-backed).
- `store/postgres_store/queries/user_roles.sql` — added `ListRoleNamesForUser :many` (joins `user_roles × roles`, ORDER BY role name).
- `internal/mocks/istore.go` — three new overrides for the service-test doubles: `GetSessionByID`, `ListRoleNamesForUser`, `UserHasPermission`.
- `service/permissions_test.go` (new) — 11 cases: `ValidateToken` happy / empty-token / malformed-token / revoked-session / expired-session; `CheckPermission` allow / deny / unknown-scope / missing-field / malformed-uuid / store-error. Uses a real PASETO maker for the happy-path token so the verifier round-trip is exercised.

Seed fixtures landed in `migration/000005_seed_iam_test_roles_and_permissions.{up,down}.sql`: one `journal_entry/read/global` permission, two roles (`finance_admin` + `cs_agent`), two dev-only users (`finance@umrohos.dev` with `finance_admin`, `cs@umrohos.dev` with `cs_agent`). `finance_admin` and `super_admin` both hold the permission; `cs_agent` holds nothing — drives the deny path. Passwords reuse the dev-only bcrypt hash from `000004_seed_initial_admin`.

Consumer-side wire-up landed alongside in `services/finance-svc/`: new `adapter/iam_grpc_adapter/` + `api/rest_oapi/middleware/bearer_auth.go` + `GET /v1/finance/ping` handler + `cmd/start.go` dial to `iam-svc:50051`. The bearer middleware fails closed on any iam error (`mapIamError` default branch → `ErrUnauthorized`) per F1 "never default to allow". e2e coverage: `tests/e2e/tests/02b-iam-svc-permission-gate.spec.ts` (4 cases — 200 / 403 / 401 / 401). (finance-svc's `04-status.md` stays pristine — its own domain checklist lands with S3-E-03 + S3-E-07; nothing in this card's scope fills any of those boxes.)

## 2026-04-21 — BL-IAM-003 admin suspend + revoke-all-sessions

Third code-adding sibling of `S1-E-04`. Closes F1-W5 acceptance ("Suspended user cannot access again") via a single admin-facing REST action plus a defense-in-depth guard on the hot-path gRPC validator. Producer-internal only — no new `§ IAM` section needed in `docs/contracts/slice-S1.md` and no Changelog entry required (iam-svc is the sole consumer of the suspend surface; downstream services still reach it indirectly via their existing `ValidateToken` gate, which already fails closed on revoked sessions).

- `api/rest_oapi/openapi.yaml` — new `POST /v1/users/{id}/suspend` (bearer + `iam.users/suspend/global` required). Response: `{ data: { user: UserProfile } }` with `user.status="suspended"`. 200 / 400 / 401 / 403 / 404 responses documented. Regenerated `api.gen.go` via `make oapi`.
- `api/rest_oapi/admin.go` (new) — handler extracts the actor id from the bearer payload, calls `service.CheckPermission(actor, "iam.users", "suspend", "global")` in-process, then `service.SuspendUser`. The permission gate runs *before* the target lookup so an un-granted bearer can't probe for the existence of arbitrary user ids.
- `cmd/server.go` — new `/v1/users` fiber group under the same `RequireBearerToken` middleware used by `/v1/me*`.
- `service/admin.go` (new) — `SuspendUser(ActorUserID, TargetUserID)` flips `iam.users.status` to `suspended` and revokes every non-revoked row in `iam.sessions` for the target inside one `WithTx`. Rejects self-suspend (`ActorUserID == TargetUserID` → `ErrValidation`) before any UUID parse or DB hit so the admin can't accidentally lock themselves out of the one seat that holds the suspend grant. Idempotent: re-suspending a `suspended` user is a no-op on status and still sweeps any sessions that raced in between. Actor id is kept on the span + log line — BL-IAM-004 will add the audit write inside the same tx.
- `service/permissions.go` — `ValidateToken` now reloads `users.status` after the session checks pass and rejects `!= active` with `ErrUnauthorized`. Adds one DB hit per authenticated call (via `GetUserByID` — index-backed on the primary key) on top of the existing `GetSessionByID` + `ListRoleNamesForUser`. Defense-in-depth: the `SuspendUser` tx already revokes every session so the session-row check alone is load-bearing in practice, but any future admin path that mutates status without also revoking sessions would otherwise let an in-flight access token keep working until its TTL elapses.
- `util/apperrors/grpc.go` — new `GRPCMessage(err) string` returning constants only (`"not found"` / `"conflict"` / `"validation error"` / `"unauthorized"` / `"forbidden"` / `"internal error"`, default `"internal error"`). `api/grpc_api/server.go:70,111` swap `err.Error()` → `apperrors.GRPCMessage(err)` on the gRPC wire; the full wrapped chain still goes to zerolog `.Err(err)` and `span.RecordError` unchanged. This is the ride-along remediation of the session-state oracle surfaced by BL-IAM-002's `/security-review` (confidence 6/10 at the time — below the fix-before-merge threshold, tracked here).
- `service/permissions_test.go` — extended the `ValidateToken_happyPath` to also mock `GetUserByID` returning `IamUserStatusActive`; new `rejectsSuspendedUserWithLiveSession` case (suspended user, live non-revoked session → `ErrUnauthorized`, `ListRoleNamesForUser` never called).
- `service/admin_test.go` (new) — six pre-tx guard cases: missing fields, self-suspend, malformed actor uuid, malformed target uuid, target not found (`pgx.ErrNoRows` → `ErrNotFound`), lookup store error (→ `ErrInternal` via `WrapDBError`). The `WithTx` body itself is exercised end-to-end in the 02c e2e spec rather than with a hand-rolled mock (same pattern as `Test_Logout_*` in `auth_test.go`).
- `util/apperrors/grpc_test.go` (new) — table-driven coverage for both `GRPCCode` and `GRPCMessage`; the message table includes two cases that explicitly verify the wrapped inner chain (`"load session: not found"`, `"session revoked"`) does NOT leak through.

Seed fixtures landed in `migration/000006_seed_iam_user_suspend_permission.{up,down}.sql`:

- Permission tuple `iam.users / suspend / global` (UUID `9999...9999`) granted to `super_admin` (UUID `2222...2222`). `finance_admin` and `cs_agent` deliberately do not receive the grant — they drive the permission-deny path in the e2e.
- Dedicated fixture user `suspend-target@umrohos.dev` (UUID `aaaa...aaaa`, HQ branch, status `active`, no role) — the e2e's sacrifice user. Using a separate fixture keeps `cs@umrohos.dev` active so BL-IAM-002's `02b-iam-svc-permission-gate.spec.ts` stays green across suite re-runs. Re-runs of `02c` against a stack where the sacrifice is already suspended need a `make migrate-down STEPS=1 && make migrate-up` to re-seed — the testing-guide walks through this.

e2e coverage: `tests/e2e/tests/02c-iam-svc-suspend.spec.ts`, ten cases run serially: (1) baseline login for admin + sacrifice, (2) sacrifice's active bearer against finance-svc ping returns 403 FORBIDDEN (auth OK, no grant), (3) admin suspends sacrifice → 200 with `status=suspended`, (4) sacrifice's *in-flight* bearer against finance-svc ping → 401 with a sanitised `message` (assert the body does NOT contain `"session revoked"` / `"load session"` / `"not found"` / `"user status="` — proves the BL-IAM-002 leak remediation), (5) sacrifice login → 403 (status gate in `service.Login`), (6) sacrifice refresh → 401 (session revoked in `service.RefreshSession`), (7) finance_admin bearer cannot call the suspend endpoint → 403 (permission gate), (8) re-suspend → 200 idempotent, (9) admin self-suspend → 400 `VALIDATION_ERROR`, (10) no bearer → 401 `UNAUTHORIZED`. All 10/10 green against the live stack; `02a` + `02b` stay 14/14 green.

Contract gate: unchanged. iam-svc is the sole producer-consumer of the new permission-tuple shape; no other service pulls a proto or schema from this card. Slice-S1 contract (`docs/contracts/slice-S1.md`) was inspected — no IAM section exists, no IAM section was added, no Changelog append needed.

Security review: one pass, zero findings at confidence ≥ 8 — the BL-IAM-002 state-oracle remediation lands clean, the permission gate runs before the target lookup (no existence-oracle leak on 404), the suspend tx keeps status-flip + session-revoke atomic, and the permission seed grants the new action exclusively to `super_admin`.
