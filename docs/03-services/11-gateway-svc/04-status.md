# gateway-svc — Status

## Implementation checklist

- [~] Scaffolded as a stateless REST edge proxy (no DB, no gRPC server)
- [~] Wired into `docker-compose.dev.yml` and `monitoring/prometheus.yml`
- [ ] First REST adapter: `iam_rest_adapter` + `GET /v1/iam/system/live` proof
- [ ] REST adapters for the remaining 9 backends (catalog, booking, jamaah, payment, visa, ops, logistics, finance, crm) — added per backend as routing needs each one
- [ ] Auth middleware (validate token via iam-svc, propagate branch scope) — F1.7 dependency
- [ ] Per-route forwarding rules (`/v1/<svc>/...` → backend) — added with each adapter
- [ ] Grafana dashboard `gateway-svc.json` in the `UmrohOS Services` folder
- [ ] Verified by reviewer in `testing-guide.md`

## Current status

**Scaffolded** — service lives at `services/gateway-svc/`; REST on `4000`; **no database, no gRPC server**. The decision to use the **adapter pattern over REST** for backend calls (with a future gRPC swap path) was confirmed 2026-04-17.

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

## Adapter pattern (REST variant)

gateway-svc introduces the **REST adapter** pattern as a first-class shape alongside the baseline's `*_grpc_adapter/`. Each backend service the gateway calls gets its own package under `services/gateway-svc/adapter/<svc>_rest_adapter/`:

- `adapter.go` — `Adapter` struct (logger, tracer, `*http.Client` wrapped with `otelhttp` for trace propagation, `baseURL`); `NewAdapter(logger, tracer, baseURL) *Adapter`.
- `<topic>.go` — typed methods that call the backend's REST endpoints and decode the response envelope into typed structs.

This commit only scaffolds the gateway shell — the first adapter (`iam_rest_adapter`) and its proxy proof endpoint land in the next commit.

## Next

- Wire `iam_rest_adapter` + `GET /v1/iam/system/live` proxy proof.
- Add `grafana/dashboards/gateway-svc.json`.
- Add the remaining 9 REST adapters as routing needs each one (not at scaffold time).
