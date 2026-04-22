# ADR 0009 — REST only on gateway; gRPC-only backends; single-point auth at gateway

**Status:** Accepted
**Date:** 2026-04-22
**Supersedes (in part):** The dual-transport convention implied by `docs/04-backend-conventions/05-api-design.md` § "When REST vs gRPC" (lines 101–106 in the pre-0009 revision) and the per-service REST + gRPC shape listed in `docs/01-architecture/02-service-map.md`. Those docs are rewritten to match this ADR.

## Context

The repository scaffolded every service with **both** REST (via `oapi-codegen`) and gRPC (via `protoc`) surfaces, modeled on `baseline/go-backend-template/demo-svc/`. `gateway-svc` was scaffolded as a REST-only edge proxy that calls downstream services via REST adapters (`services/gateway-svc/adapter/<svc>_rest_adapter/`). Authentication was intended as defense-in-depth: gateway would eventually run a bearer-validation middleware, and every backend service **also** ran `RequireBearerToken` middleware that called `iam-svc.ValidateToken` via gRPC on every authenticated request (`docs/04-backend-conventions/05-api-design.md:62`; `docs/03-services/00-iam-svc/00-overview.md`).

Two facts emerged from running this shape for S0–S1:

1. **The gRPC surfaces on non-gateway services were scaffolded but unused.** Only `iam-svc` exposed real RPCs (`ValidateToken`, `CheckPermission`, `RecordAudit`). Every other service carried a placeholder `Healthz` RPC. Real business RPCs (`ReserveSeats`, `ReleaseSeats`, `GetPackageDeparture`, etc.) would land per slice as consumers materialized — leaving the scaffolded gRPC listener unused for weeks to months per service. `gateway-svc` had no gRPC adapters at all.
2. **The dual-validation auth pattern duplicated logic.** Even though both layers call the same `iam-svc.ValidateToken`, every new authenticated route in a backend service needed its own wiring of `RequireBearerToken`. The "belt and braces" defense-in-depth benefit was accepted abstractly but not load-bearing in practice — the gateway edge layer was still a dependency.

The question we faced: what is the right transport split and auth model for a two-developer team shipping an MVP at 10k-jamaah/year scale?

## Decision

Two coupled decisions, accepted together.

### D1. Transport split

**REST lives only on `gateway-svc`.** All downstream services (`iam-svc`, `catalog-svc`, `booking-svc`, `payment-svc`, `jamaah-svc`, `visa-svc`, `ops-svc`, `logistics-svc`, `finance-svc`, `crm-svc`) are **gRPC-only** for business and domain calls. Their `api/rest_oapi/` package is removed.

- Client apps (browser, mobile, B2B agent portal) talk REST to `gateway-svc` only.
- `gateway-svc` proxies every client request to the appropriate backend service via a `<svc>_grpc_adapter` (one per backend, modeled on `services/booking-svc/adapter/iam_grpc_adapter/`).
- Service-to-service east-west traffic is gRPC as before.

**Exceptions — admin/observability endpoints, not "REST API":**
- Each backend keeps a minimal HTTP handler exposing only `/metrics` (Prometheus scrape) on a dedicated admin port. This is a single-handler admin endpoint, not an OpenAPI-spec'd REST API. Alternatively, metrics can push via OTLP to the OTel Collector (larger migration, deferred unless Option A is insufficient).
- Liveness / readiness probes use the **standard gRPC health protocol** (`grpc.health.v1.Health/Check` via `google.golang.org/grpc/health`). `grpc_health_probe` is the canonical probe binary for docker-compose / Kubernetes health checks.
- The existing `/system/diagnostics/db-tx` REST endpoint (S0-J-05's WithTx trace verification) migrates to a `DiagnosticsDbTx` gRPC RPC; `gateway-svc` keeps the REST route `/v1/<svc>/system/diagnostics/db-tx` and proxies to the gRPC method.

### D2. Single-point authentication

**Authentication is validated once, at `gateway-svc`.** No backend service re-validates the bearer.

- `gateway-svc` runs a Fiber middleware on authenticated routes: extract `Authorization: Bearer`, call `iam-svc.ValidateToken` via gRPC (gateway grows its own `iam_grpc_adapter`), attach the identity envelope to the request context, forward via the appropriate `_grpc_adapter` to the downstream gRPC method. Fail closed: if `iam-svc` is unreachable, return 502 with apperrors envelope.
- Backend services **do not** extract bearer tokens, run `RequireBearerToken`, or call `iam-svc.ValidateToken`. That middleware is removed from every service.
- Authorization (`iam-svc.CheckPermission`) and audit (`iam-svc.RecordAudit`) remain at the backend's service layer — they are not authentication and often require business context only available on the backend.

**Accepted trade:** defense-in-depth is reduced. In MVP, backends trust that a gRPC call came from gateway because the internal docker-compose network is the only path to them. This is deliberately weaker than "gateway AND backend both validate." The compensating control is **deferred as `BL-GTW-100`** — a gateway↔backend trust contract (signed header or mTLS) for a later slice — and recorded as a known gap.

### D3. Every new client-facing endpoint must include its gateway side in the same card

A backend gRPC method that is reachable from a client browser MUST ship with its gateway REST route + `_grpc_adapter` proxy in the same branch/PR/card. Half-shipped is not shipped. E2E tests always target `gateway-svc:4000`; no Playwright spec hits a backend port directly.

## Consequences

### Services

- `baseline/go-backend-template/demo-svc/` (dual-transport) is the template **for `gateway-svc` only**. All other services scaffold from `baseline/go-backend-template/demo-grpc-svc/` (gRPC-only) going forward.
- Existing scaffolded services (all 10 non-gateway) need their `api/rest_oapi/` package removed as part of the transition. This is tracked as `BL-REFACTOR-001..010` per service; each refactor card pairs with the gateway-side routing card that replaces the public surface.
- `services/iam-svc/api/rest_oapi/` — auth endpoints (`POST /v1/sessions`, `GET /v1/me`, `/v1/users/*`, TOTP flows) move to `gateway-svc`'s REST surface, proxied to new iam-svc gRPC methods. Tracked as `BL-IAM-018`.

### Gateway

- `services/gateway-svc/` grows one `<svc>_grpc_adapter/` per backend it needs to call. The current `iam_rest_adapter/` etc. are retired (the `iam_rest_adapter` used for `/v1/iam/system/live` trace-propagation proof from S0-J-05 is replaced with a gRPC-based equivalent).
- `services/gateway-svc/` is now the sole owner of `api/rest_oapi/openapi.yaml`. All public routes, public schemas, and the bearer-auth middleware live here.

### Monitoring

- `grpc.health.v1.Health` registered in each backend's `cmd/start.go`.
- `docker-compose.dev.yml` health-checks use `grpc_health_probe -addr=localhost:<grpc-port>` instead of `curl`.
- Prometheus continues to scrape `/metrics`, but on each backend's admin port (not the removed REST port) unless the OTLP push path is adopted. See `BL-MON-001`.
- Placeholder `Healthz` RPCs in every service proto are retired in favor of the standard health service.
- OTel tracing is unaffected — `otelgrpc` handlers are already wired.

### Trust / security

- `BL-GTW-100` (deferred) designs the gateway↔backend trust contract (signed header or mTLS) to close the defense-in-depth gap introduced by D2. Until then, MVP relies on internal-network isolation as the only authentication guarantee below the gateway.

### Documentation

- `docs/04-backend-conventions/05-api-design.md` § "When REST vs gRPC" and § "Auth" are rewritten to match.
- `docs/01-architecture/00-system-overview.md` and `docs/01-architecture/02-service-map.md` drop the "REST + gRPC" shape for non-gateway services.
- `docs/01-architecture/03-data-flow.md` gateway→backend arrow relabeled from "REST" to "gRPC".
- `docs/01-architecture/04-observability.md` § "What to verify on every new service" is rewritten for the new health + metrics approach.
- `docs/06-features/01-identity-and-access.md` grows **F1-W7 "Edge auth at `gateway-svc`"**; W3 ("`CheckPermission` middleware for internal routes") loses the "each service re-validates bearer" language.
- `docs/contracts/slice-S1.md` grows `§ Gateway` with the routes table and single-point-auth rule.

## Rationale

1. **Use the tools we scaffold.** Carrying gRPC surfaces that no caller reaches for months is waste. Making gateway→backend traffic gRPC means every scaffolded gRPC surface is immediately load-bearing.
2. **Single responsibility, single place to reason about auth.** "Gateway validates; backends trust" is easier to hold in your head than "both validate; reconcile if one is ever bypassed." The explicit trust boundary at gateway is the security architecture, not a TODO.
3. **Consistent transport boundary matches operational experience.** All internal traffic speaks gRPC (with `otelgrpc` tracing); all external traffic speaks REST (at gateway). Two transports, two jobs, one place each.
4. **Monitoring concern we investigated during the decision is addressed.** Liveness / readiness use `grpc.health.v1`; `/metrics` keeps being scraped via a tiny admin HTTP endpoint; OTel traces already work for gRPC. No observability capability is lost.
5. **E2E testing through gateway forces the full path to be correct.** Tests that bypass gateway to hit a backend port would pass while gateway routing, adapter, or auth is broken — hiding real regressions.

## Alternatives considered and rejected

- **Keep dual-transport backends, build gateway REST→backend REST forwarding per plan.** Rejected. The gRPC surfaces stay unused; the duplicate auth-validation layer adds wiring cost without adding load-bearing security value at MVP scale.
- **Remove gRPC entirely; make everything REST (including east-west).** Rejected. This was briefly considered when the user asked "what's the point of gRPC." Cost is large: rewrite iam-svc's `ValidateToken`/`CheckPermission`/`RecordAudit` as REST, rewrite `booking-svc` and `finance-svc` iam adapters, supersede ADR 0004 on proto ownership, lose protobuf type safety and streaming capability (F2's planned `CatalogUpdatedSubscribe`), and the monitoring rationale it rested on did not hold up (probes are already REST regardless of whether gRPC exists). If simplicity-over-performance were compelling enough we'd take the hit, but REST-only-at-gateway accomplishes the same "one transport per direction" goal without those costs.
- **Defense-in-depth auth (gateway validates AND each backend revalidates).** Rejected in favor of single-point at gateway (D2). This is the more contested trade of the ADR. The user explicitly chose single-responsibility; the defense-in-depth layer is scheduled as a backlog item (`BL-GTW-100`) via the trust contract rather than a per-service middleware.
- **Keep REST on every backend for client compatibility during the transition.** Partially accepted as a sequencing choice, not an endpoint state: `G6` (catalog-svc gateway routing) lands before `G7` (catalog-svc REST removal) so the gateway path is proven in e2e before REST is deleted. But the end state is gRPC-only.

## When to revisit

Triggers for reversing or amending this ADR:

1. **Trust contract (`BL-GTW-100`) proves operationally infeasible.** If mTLS between gateway and backends is a burden (cert rotation, debugging friction) and a signed-header approach also fails, we may want the per-service `RequireBearerToken` back as the trust mechanism — with the performance cost of double-validation on every request.
2. **gRPC monitoring gaps emerge.** If `grpc_health_probe`, OTel, or Prometheus scraping on admin ports produces gaps in Grafana that we can't close cheaply, reconsider the admin `/metrics` approach (switch to OTLP push) or, in extremis, reconsider REST on backends for observability.
3. **Client surfaces materialize that benefit from direct gRPC (e.g., gRPC-Web).** If frontend clients eventually talk gRPC-Web, "REST only at gateway" may need nuancing.
4. **Streaming needs emerge that gateway REST can't serve cleanly** (e.g., catalog updates subscription F2). At that point gateway may need a gRPC-Web or SSE bridge; ADR extension, not reversal.

## Notes

- `WithTx` at the store layer remains unchanged — this ADR is about transport and auth, not transactions.
- ADR 0004 (proto ownership) is unaffected. Proto files stay in `<svc>/api/grpc_api/pb/<svc>.proto`; consumers still carry local copies under `<consumer>/adapter/<dep>_grpc_adapter/pb/`.
- ADR 0006 (in-process sagas for MVP) is unaffected. Saga orchestration remains in-process; sagas just talk gRPC to collaborators instead of REST.
- The other developer's `AGENTS.md` may reference the old dual-transport shape; flag separately for their update.
