# gateway-svc — Status

## Implementation checklist

- [~] Scaffolded as a stateless REST edge proxy (no DB, no gRPC server)
- [~] Wired into `docker-compose.dev.yml` and `monitoring/prometheus.yml`
- [x] First adapter: `iam_rest_adapter` + `GET /v1/iam/system/live` proof *(interim; replaced by `iam_grpc_adapter` in `BL-GTW-001` per ADR 0009)*
- [ ] **`iam_grpc_adapter`** — gateway's gRPC client to `iam-svc` (replaces the interim REST adapter per ADR 0009) — `BL-GTW-001` / S1-E-09
- [ ] **Bearer-auth middleware** — extracts `Authorization: Bearer`, calls `iam.v1.IamService/ValidateToken`, fail-closed 502 on iam-svc unreachable — `BL-GTW-001` / S1-E-09
- [x] **`catalog_grpc_adapter`** + public REST routes (`GET /v1/packages`, `GET /v1/packages/{id}`, `GET /v1/package-departures/{id}`); e2e migrated to `gateway-svc:4000` — `BL-GTW-002` / S1-E-10 (merged 2026-04-22, PR #48)
- [x] **`catalog_rest_adapter` retired** + `/v1/catalog/system/live` removed + `external.catalog_svc.address` dropped from config — `BL-REFACTOR-001` / S1-E-11 (2026-04-23). catalog-svc is now gRPC-only; operators probe via `grpc_health_probe`.
- [x] **Seven scaffold REST adapters retired** (`booking/crm/jamaah/logistics/ops/payment/visa`) + seven `/v1/<svc>/system/live` routes removed + seven `external.<svc>_svc.address` fields dropped — `BL-REFACTOR-002..008` / S1-E-13 (2026-04-23). The 7 backends now run gRPC-only; operators probe via `grpc_health_probe`.
- [x] **Iam client-facing REST** (`/v1/sessions*`, `/v1/me*`, `/v1/users/{id}/suspend`) proxied to iam gRPC via `iam_grpc_adapter` + bearer middleware mounted on every protected route — `BL-IAM-018` / S1-E-12 (2026-04-23). `iam_rest_adapter` + `/v1/iam/system/live` + `/v1/iam/system/diagnostics/db-tx` + `external.iam_svc.address` all removed.
- [x] **finance-svc ADR 0009 realignment** — `/v1/finance/ping` moved to gateway behind new `RequirePermission("journal_entry","read","global")` middleware + `finance_grpc_adapter`; `/v1/finance/system/live` deleted (no per-backend REST liveness proxy survives); `finance_rest_adapter` directory + `external.finance_svc.address` config removed — `BL-IAM-019` / S1-E-14 (2026-04-24). Last REST adapter on the gateway; ADR 0009 refactor arc is functionally complete (G10 trust contract deferred).
- [x] **Scaffold-RPC cleanup** — `rpc Healthz` + messages removed from every backend proto (10 services + gateway catalog mirror); `rpc DiagnosticsDbTx` + handler removed from catalog; `service.DbTxDiagnostic` + sqlc `InsertDbTxDiagnostic` / `GetDbTxDiagnostic` dropped across all 10 backends. Container-level liveness now uniformly served by `grpc.health.v1.Health` (BL-MON-001) — rolled into `BL-IAM-019` / S1-E-14 (2026-04-24).
- [x] **Scaffold `Ping` RPC uniform across all 10 backends** — every backend proto exposes `rpc Ping(PingRequest) returns (PingResponse)` returning `{message:"ok"}` (renamed from finance's `FinancePing` per the reviewer's scaffold-uniformity principle). Keeps the permission-gate smoke (gateway → finance gRPC via `Ping`) and gives every future slice a ready-made reachability target — `BL-IAM-019` / S1-E-14 (2026-04-24).
- [x] **`GET /v1/system/backends` aggregate health endpoint** — new `health_check_adapter` holds one `grpc.health.v1.Health` client per backend; the REST handler fans `Health.Check` out concurrently (2 s per-backend timeout) and returns `[{name, status, error?}]` sorted alphabetically. Replaces the retired per-backend `/v1/<svc>/system/live` proxies with a single uniform endpoint driven by the standard gRPC health protocol. `cmd/start.go` now dials all 10 backends at startup; the same conn map backs per-route business adapters and the health-check adapter. core-web's status page rewired to poll this endpoint (one request per tick instead of ten) — `BL-IAM-019` / S1-E-14 (2026-04-24).
- [ ] Per-backend gRPC adapters for the remaining services (booking, jamaah, payment, visa, ops, logistics, finance, crm) — opened as each consumer slice lands
- [ ] **Trust contract (gateway↔backend)** — signed header or mTLS, closes the defense-in-depth gap — `BL-GTW-100` (deferred, later slice)
- [ ] Grafana dashboard `gateway-svc.json` in the `UmrohOS Services` folder
- [ ] Verified by reviewer in `testing-guide.md`

## Current status

**Scaffolded** — service lives at `services/gateway-svc/`; REST on `4000`; **no database, no gRPC server for inbound calls** (gateway is the edge; it is itself only a REST server). Outbound to backends was originally REST-over-adapters; per ADR 0009 (2026-04-22) it transitions to **gRPC-over-adapters** — backends expose gRPC only, gateway carries a `<svc>_grpc_adapter` per backend it calls. The interim `iam_rest_adapter` is retired as `BL-GTW-001` / S1-E-09 lands.

**Scaffold deliverables (this commit only):**

- `GET /system/live` — liveness probe.
- `GET /system/ready` — readiness probe (process-local at scaffold time; later iterations may aggregate per-adapter health).
- Prometheus `/metrics` scraped at `gateway-svc:4000`.

**Pilot strip (vs. the iam-svc shape):**

Removed — gateway is stateless and REST-only:

- `store/` — no DB.
- `api/grpc_api/` — no gRPC server (gateway is the edge; everything goes out as REST for now).
- `cmd/conn.go` — no Postgres pool.
- `util/monitoring/dbpool.go` — no DB to track.
- `util/apperrors/grpc.go` — no gRPC surface.
- `/system/diagnostics/db-tx` route + handler — no DB.
- `Store` config block + `POSTGRES_*` env bindings.

Kept (and reused as-is from iam-svc):

- `util/{config,logging,tracing,monitoring,apperrors}` (minus the gRPC + DB-pool slices noted above).
- `cmd/{main.go,start.go,server.go}` — trimmed to the REST-only shape.
- `api/rest_oapi/{server.go,system.go,middleware/error.go}` — system probes only.

## Assigned ports

| Surface          | Port    |
|------------------|---------|
| REST (Fiber)     | `4000`  |

## Adapter pattern (REST in, gRPC out — per ADR 0009)

Gateway is the only REST surface in UmrohOS. For each backend it calls, it carries a **gRPC adapter** under `services/gateway-svc/adapter/<svc>_grpc_adapter/`:

- `adapter.go` — `Adapter` struct (logger, tracer, `pb.<Svc>ServiceClient` wrapped with `otelgrpc.NewClientHandler()`); `NewAdapter(logger, tracer, conn) *Adapter`.
- `<topic>.go` — typed methods that call the backend's gRPC RPCs and translate gRPC status codes to `apperrors` sentinels for consistent Fiber rendering.
- `pb/<svc>.proto` — local copy of the backend's proto (per ADR 0004).

The original REST-adapter pattern (`services/gateway-svc/adapter/<svc>_rest_adapter/`) was an interim shape used during S0 scaffolding. Each backend's REST adapter retires as its `BL-REFACTOR-*` card lands; `catalog_rest_adapter` went first (2026-04-23, `BL-REFACTOR-001`); the **seven scaffold adapters** (`booking/crm/jamaah/logistics/ops/payment/visa`) followed in a single sweep as `BL-REFACTOR-002..008` / S1-E-13 (2026-04-23). `iam_rest_adapter` still serves the interim system-probe + WithTx-diagnostic routes and retires with `BL-IAM-018` / S1-E-12. `finance_rest_adapter` still serves `/v1/finance/system/live` + the BL-IAM-002 `/v1/finance/ping` permission-gate demo and retires with `BL-IAM-019` / S1-E-14.

## Next

- Land `BL-GTW-001` / S1-E-09 — `iam_grpc_adapter` + Bearer middleware.
- Land `BL-GTW-002` / S1-E-10 — `catalog_grpc_adapter` + public GET routes; e2e migrates to gateway.
- Land `BL-IAM-018` / S1-E-12 — iam client-facing REST routes on gateway.
- Add `grafana/dashboards/gateway-svc.json`.
- Per-backend gRPC adapters for booking / jamaah / payment / visa / ops / logistics / finance / crm — opened as each consumer slice needs them.

## 2026-04-21 — S0-J-05 OpenTelemetry baseline fix

- `cmd/server.go` — `app.Use(otelfiber.Middleware(...))` wired as the first middleware after CORS; span name formatter prefixes with `<service-name-tracer>` for readable Tempo output. Gateway is the trace origin for edge requests and a trace continuer when upstream clients propagate.
- `util/tracing/tracing.go` — `otel.SetTextMapPropagator(NewCompositeTextMapPropagator(TraceContext{}, Baggage{}))` set globally. Without this, `otelhttp.NewTransport` in every `*_rest_adapter` was silently not injecting `traceparent`.
- `api/rest_oapi/proxy_iam.go` + other `proxy_*.go` — spans start from `c.UserContext()` (otelfiber's inbound-span context).
- **New traced cross-service path:** `GET /v1/iam/system/diagnostics/db-tx` — proxies iam-svc's WithTx diagnostic. `openapi.yaml` + regenerated `api.gen.go`; `service/proxy_iam.go` + `service/service.go` add `GetIamSystemDbTxDiagnostic`; `adapter/iam_rest_adapter/system.go` adds `GetSystemDbTxDiagnostic` with `DbTxDiagnosticResult` type. This is the verification endpoint S0-J-05's acceptance criterion flows through.
- `go.mod` — added `github.com/gofiber/contrib/otelfiber/v2 v2.2.3`.
