---
id: F1
title: Identity, Access, Audit
status: spec written, implementation not started
last_updated: 2026-04-14
moscow_profile: 10 Must Have / 3 Should Have (all High)
prd_sections:
  - "H. Admin & Security — RBAC, audit trail, system config"
modules:
  - "#151–163 (Admin & Security)"
depends_on: []
---

# F1 — Identity, Access, Audit

## Purpose & personas

UmrohOS is single-tenant per travel agency but serves many kinds of users inside one company: field staff, branch managers, marketing agents, customer service reps, warehouse staff, finance admins, jamaah themselves, and the company owner. Each role sees different data and can take different actions.

This feature establishes **who a user is**, **what they're allowed to do**, **which data they can see**, and **a tamper-resistant record of what they did**. Every other feature in the system depends on it.

Primary personas:
- **Super Admin / IT Manager** — configures roles, permissions, branches; manages API keys and backups
- **Branch Manager** — manages staff and agents within their branch only
- **Staff (ops / finance / warehouse / CS)** — day-to-day operators, permissions scoped to their job
- **Agent (reseller)** — logs in to a B2B portal to register jamaah; sees only their own leads
- **Jamaah** — logs in to a customer portal; sees only their own bookings and documents

## Sources

- PRD Section H in full
- Module #151 Pembuatan Peran Jabatan (Must Have / High)
- Module #152 Pemetaan Izin Spesifik (Must Have / High)
- Module #153 Hierarki Visibilitas Data (Must Have / High)
- Module #154 Pendaftaran Akun Staf (Must Have / High)
- Module #155 Kontrol Status Pengguna (Must Have / High)
- Module #156 Keamanan Akun & Sandi (Must Have / High)
- Module #157 Log Aktivitas Terpusat (Must Have / High)
- Module #158 Peringatan Anomali (Should Have / Medium)
- Module #159 Riwayat Sesi Pengguna (Should Have / Medium)
- Module #160 Konfigurasi Integrasi API (Must Have / High)
- Module #161 Manajemen Template Komunikasi (Should Have / Medium)
- Module #162 Konfigurasi Global Variabel (Must Have / High)
- Module #163 Pencadangan Database (Must Have / High)

## User workflows

### W1 — Staff login with 2FA

1. Staff user opens the internal console login page.
2. Submits email + password.
3. System validates credentials. On first login after enrollment, prompts for TOTP enrollment.
4. On subsequent logins, prompts for TOTP code.
5. On success: issues an access token (short-lived, 15 min) and a refresh token (7 days, rotated on use).
6. Console shows the sidebar menu filtered to the user's assigned role and branch scope.

**_(Inferred)_** Token format is PASETO v4 local symmetric. Rationale: simpler key management for single-tenant, and the baseline template already provides `util/token` for this.

### W2 — Admin grants a permission to a role

1. Super Admin opens "Roles & Permissions" screen.
2. Selects an existing role (e.g. `finance_admin`).
3. Sees the permission matrix: rows = resources (`booking`, `invoice`, `journal_entry`, ...), columns = actions (`read`, `write`, `edit`, `delete`), and a separate scope selector (`global`, `branch`, `personal`).
4. Toggles a cell; confirmation dialog asks for admin password re-entry (mini step-up auth for sensitive config).
5. Change is written; audit log records the old matrix, the new matrix, the actor, and the IP.
6. Any user with that role has their cached permissions invalidated within 60 seconds.

### W3 — Permission check on an API call

For every authenticated request (per ADR 0009 — authentication is single-point at gateway; permission checks remain at the backend's service layer):

1. Gateway extracts the bearer token and calls `iam-svc.ValidateToken` (gRPC) — see W7 for the full edge-auth flow. On success, gateway forwards to the target backend via gRPC with the validated identity envelope.
2. Service layer on the target backend calls `iam-svc.CheckPermission(user_id, resource, action, scope)` (gRPC).
3. If the permission allows the action, the request proceeds. Otherwise 403 is returned (via `apperrors.ErrForbidden`).
4. The resolution is cached in the target service for a short TTL (e.g. 60s) to avoid hammering `iam-svc`.
5. Every successful state-changing call emits an audit log entry via `iam-svc.RecordAudit`.

### W4 — Data scope filtering

A `branch_manager` at "Jakarta Pusat" querying `/v1/bookings` sees only bookings whose `branch_id` matches theirs. A `super_admin` sees all. A `cs_agent` with `personal` scope sees only bookings they personally created.

The scope comes from the user's permission grant and is attached to the authenticated context. The service layer, not the API layer, applies the scope as a `WHERE` clause predicate in sqlc queries.

### W5 — Account suspended or password reset forced

1. Super Admin opens a user's page; clicks "Suspend" or "Force password reset".
2. User's sessions are immediately revoked (all refresh tokens marked revoked).
3. Audit log records the action.
4. On next login attempt, the user either sees "Account suspended" or is required to reset the password before proceeding.

### W6 — Audit log review

1. Super Admin or auditor opens "Audit Log" screen with filters (user, date range, resource, action).
2. System returns append-only records: who, what, old value, new value, when, IP.
3. Anomaly alerts (Should Have, module #158) surface unusually bulk operations, off-hours activity, or permission escalations.
4. Session history (Should Have, module #159) shows per-user active sessions with user-agent and IP.

### W7 — Edge auth at `gateway-svc`

Per ADR 0009 (REST only at gateway; gRPC-only backends; single-point auth). This wire defines how `Authorization: Bearer` is validated exactly once, at the edge, before any backend gRPC call is made.

For every authenticated request arriving at the gateway:

1. Client sends an HTTP request to `gateway-svc` with `Authorization: Bearer <token>`.
2. Gateway's Fiber middleware extracts the bearer; on missing or malformed header → 401 immediately.
3. Middleware calls `iam-svc.ValidateToken` (gRPC, via gateway's `iam_grpc_adapter`). Typical latency: single-digit ms on the internal network.
4. On `ValidateToken` = OK: attach the returned identity envelope (`user_id`, `branch_id`, `session_id`, `roles[]`) to the Fiber context; continue to the route handler.
5. On `ValidateToken` = Unauthenticated: 401 with apperrors envelope `{error: {code: "unauthorized", message: ...}}`. Never default to allow.
6. On `iam-svc` unreachable, timeout, or any other transport-level failure: 502 with `auth_unavailable` error code — **fail closed**. The client retries; gateway does not silently trust anything.
7. Route handler forwards the request to the target backend via the appropriate `<svc>_grpc_adapter`. Backends **do not** re-validate the token (see ADR 0009 rationale and the `BL-GTW-100` trust-contract deferral).

Routes tagged `security: []` in the gateway OpenAPI spec (public GETs — catalog browsing, etc.) skip this middleware entirely. See `docs/contracts/slice-S1.md` § Gateway for the route-by-route auth matrix.

**Scope attached to context:** `branch_id` from the envelope is the default scope filter used by W4. Route handlers pass it down to backend gRPC calls via request fields or gRPC metadata (finalized in `BL-GTW-001` / S1-E-09 implementation).

**Revocation latency (W5 interaction):** `ValidateToken` MUST check the `sessions` table, not only the PASETO signature, so a suspended-or-revoked session fails within the TTL of any inline cache in `iam-svc` (≤60s — matches the W3 permission cache guidance). Gateway does not cache auth results — every request is validated fresh. This keeps W5's 60-second revocation SLA honest.

## Acceptance criteria

- A user can log in with email + password + TOTP and receive a bearer token.
- A user's role grants specific (resource, action, scope) tuples; an `rbac_test` suite proves forbidden combinations are denied.
- Every state-changing call in every service produces an audit log entry with trace_id correlation.
- Suspending a user invalidates all their active sessions within 60 seconds.
- Bulk permission changes require admin re-auth (step-up).
- Branch-scoped queries never leak data outside the user's branch.
- Audit log rows cannot be mutated (enforced at the DB role level or via trigger).
- All token secrets are loaded from config (env vars in production), never hardcoded.
- Edge-auth (W7) is the sole bearer-validation point; backends never re-validate. No protected route is reachable without first passing gateway's `iam-svc.ValidateToken` check.

## Edge cases & error paths

- **Token expired** — API returns 401; client refreshes via refresh token; if refresh also expired, redirect to login.
- **Gateway cannot reach `iam-svc`** — fail closed at the edge (502 `auth_unavailable`). Never default to allow. (Per ADR 0009; was 401 pre-2026-04-22 — kept distinct from 401 so clients can distinguish "token bad" from "auth layer down.")
- **Internal call bypasses gateway** — since MVP does not implement a gateway↔backend trust contract (deferred as `BL-GTW-100`), any caller that can reach a backend's gRPC port directly is trusted. Mitigation in MVP is internal-network isolation (docker-compose network or Kubernetes NetworkPolicy); hardening is tracked in `BL-GTW-100`.
- **User's role changes mid-session** — cached permissions auto-invalidate within 60s; very-sensitive actions (permission grants, API key rotation) force a fresh permission fetch.
- **TOTP device lost** — admin-assisted reset: Super Admin clears the `totp_secret` column on the user row; user re-enrolls on next login. Record the reset in audit.
- **Refresh token replay attempt** — rotation-on-use detects reuse; all of the user's sessions are revoked and a security alert is logged.

## Data & state implications

Owned tables in `iam-svc` (details in `docs/03-services/00-iam-svc/02-data-model.md`):

- `branches`, `users`, `roles`, `permissions`, `role_permissions`, `user_roles`
- `sessions` — active refresh tokens, rotation state
- `audit_logs` — append-only

Status transitions:
- `users.status` — `active` ↔ `suspended`, `pending` → `active` after first login+TOTP
- `sessions.revoked_at` — null → set

## API surface (high-level)

Full contracts live in `docs/03-services/00-iam-svc/01-api.md`. Per ADR 0009 all REST routes are hosted by `gateway-svc` and proxied to `iam-svc` over gRPC (delivered by `BL-IAM-018` / S1-E-12). Summary here:

REST routes on `gateway-svc` (internal console + B2C portal + B2B portal):
- `POST /v1/sessions` — login
- `POST /v1/sessions/refresh` — rotate
- `DELETE /v1/sessions` — logout
- `GET /v1/me` — current user
- `POST /v1/me/2fa/enroll`, `POST /v1/me/2fa/verify`
- `GET|POST|PATCH|DELETE /v1/users`, `/v1/roles`, `/v1/permissions`, `/v1/branches`
- `POST /v1/users/{id}/roles`, `DELETE /v1/users/{id}/roles/{role_id}`
- `GET /v1/audit-logs`

gRPC methods on `iam-svc` (service-to-service; also called by gateway to serve the REST routes above):
- `ValidateToken`
- `CheckPermission`
- `GetUser`
- `RecordAudit`
- `Login`, `RefreshSession`, `Logout`, `GetMe`, `EnrollTotp`, `VerifyTotp`, `ListUsers`, `SuspendUser`, ... — added by `BL-IAM-018` to serve the gateway REST routes.

## Dependencies

None. F1 is the root.

## Backend notes

- Per ADR 0009, `iam-svc` is gRPC-only; its scaffold baseline is `baseline/go-backend-template/demo-grpc-svc`. PASETO token handling (in `util/token`) is reused. The existing `api/rest_oapi/` package from the original `demo-svc` scaffold is removed as part of `BL-IAM-018` / S1-E-12 when client-facing auth REST routes move to `gateway-svc`.
- Permission matrix is hot-path — benchmark the `CheckPermission` gRPC at ~5ms p95 under load. Cache per-request.
- Audit log writes are **synchronous and atomic** with the business action they describe (revised 2026-04-21 with `BL-IAM-004`, superseding the earlier fire-and-forget/buffered-channel plan). Two emission paths:
  - **In-process** (a service auditing its own action): emit `q.InsertAuditLog(...)` inside the same `WithTx` closure that performs the state change. Business-success ↔ audit-success are atomic; a failed tx rolls back both the state change and its audit row. This is the reference pattern — see `iam-svc.SuspendUser` in `services/iam-svc/service/admin.go`.
  - **Cross-service**: call `iam.v1.IamService/RecordAudit` over gRPC (one RPC per row). The wire is synchronous and returns only after the row is durably inserted. Callers that need non-blocking semantics wrap the RPC in a goroutine on their side — the write itself remains atomic per call. See `docs/contracts/slice-S1.md § IAM` for the wire shape.
  - Append-only is enforced at the DB layer (`iam.audit_logs` BEFORE UPDATE / BEFORE DELETE trigger raises `insufficient_privilege`); compliance-driven user deletion cascades `user_id` / `branch_id` to `NULL` via the narrowed trigger (migration 000007), leaving resource / action / old_value / new_value / ip / created_at frozen.
- TOTP enrollment: store the secret encrypted at the application layer (AES-256, key from config) before INSERT. Never log the secret.

## Frontend notes

- Login screen needs three states: credentials, TOTP prompt, post-login menu.
- Sidebar menu is filtered client-side based on the permissions included in the `/v1/me` response.
- Role & permission management UI is the most complex admin screen — a matrix grid with toggles and a scope selector. Expect iteration.
- Token storage: access token in memory (never localStorage), refresh token in httpOnly cookie. _(Inferred — if the frontend lead disagrees, see Q-??? to be opened.)_

## Open questions

None currently. Add entries in `docs/07-open-questions/` as they arise and link here.
