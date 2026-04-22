# System Overview

UmrohOS is a Go microservices backend behind a public REST gateway. Per ADR 0009: **REST lives only on `gateway-svc`; every downstream service is gRPC-only.** Client apps talk REST to the gateway; the gateway proxies to backends over gRPC. Authentication is validated once at the gateway via `iam-svc.ValidateToken`; backends trust the forwarded gRPC call (the gateway↔backend trust contract is deferred as `BL-GTW-100`). Cross-service business processes are coordinated **in-process** by the orchestrating service for MVP (short-lived booking/payment/refund sagas); the F6 visa pipeline — the one multi-day durable workflow in the system — will bring Temporal back when it's implemented. See ADR 0006 for the full rationale. A full observability stack is on from day one.

## High-level diagram

```
                         ┌──────────────────────┐
                         │   Frontend            │   ← Svelte 5 (runes
                         │   (Svelte + Vite)     │     mode) + Vite
                         └──────────┬───────────┘
                                    │ HTTPS
                         ┌──────────▼───────────┐
                         │   gateway-svc        │   ← Fiber, edge auth, routing
                         │   (REST, port 4000)  │
                         └──────────┬───────────┘
                                    │ gRPC
        ┌─────────────┬─────────────┼─────────────┬─────────────┐
        ▼             ▼             ▼             ▼             ▼
   ┌────────┐    ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │ iam-svc│    │catalog- │   │booking- │   │payment- │   │ jamaah- │
   │        │    │  svc    │   │  svc    │   │  svc    │   │  svc    │
   └────────┘    └─────────┘   └─────────┘   └─────────┘   └─────────┘
        ▼             ▼             ▼             ▼             ▼
   ┌────────┐    ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │ visa-  │    │  ops-   │   │logistics│   │ finance │   │  crm-   │
   │  svc   │    │  svc    │   │  -svc   │   │  -svc   │   │  svc    │
   └────────┘    └─────────┘   └─────────┘   └─────────┘   └─────────┘

                         ┌──────────────────────┐
                         │    broker-svc        │   ← DEFERRED for MVP
                         │    (workflows)       │     (reserved for F6 visa pipeline;
                         │                      │      ADR 0006)
                         └──────────────────────┘

   ┌─────────────────────────── PostgreSQL 15 ───────────────────────────┐
   │   Each service owns its own logical DB: iam_db, catalog_db, ...    │
   └─────────────────────────────────────────────────────────────────────┘

   ┌─────────────── Observability stack (always on) ────────────────┐
   │  OTel Collector → Tempo (traces), Prometheus (metrics),         │
   │  Loki (logs via Fluent-Bit), Grafana (single pane of glass)     │
   └─────────────────────────────────────────────────────────────────┘
```

## Key principles

1. **Transport split (ADR 0009).** REST only on `gateway-svc`. All downstream services are gRPC-only — their `api/rest_oapi/` package is removed. Client apps call gateway REST; gateway proxies to backends via `<svc>_grpc_adapter` (one per backend). East-west service-to-service calls are also gRPC. Liveness/readiness use `grpc.health.v1.Health` via `grpc_health_probe`; `/metrics` stays on a minimal admin HTTP endpoint on each backend.

2. **Single-point authentication (ADR 0009).** `gateway-svc` runs the bearer-auth middleware: extract `Authorization: Bearer`, call `iam-svc.ValidateToken` via gRPC, forward on success. Backends do not re-validate. The gateway↔backend trust contract (signed header or mTLS) is deferred as `BL-GTW-100`; until then, internal-network isolation is the only authentication guarantee below the gateway.

3. **One service per bounded context.** No service owns data from another's context. Cross-context reads happen via gRPC; cross-context writes are **coordinated in-process by the orchestrating service** with explicit per-step compensation + a reconciliation cron catching mid-saga crashes (per ADR 0006). F6 visa pipeline is the one exception — it will use Temporal when implemented.

4. **Three-layer architecture per service.** API → Service → Store, with interfaces between layers. See `docs/04-backend-conventions/01-three-layer-architecture.md`. This is non-negotiable and enforced by the baseline template.

5. **Code generation everywhere.** sqlc for the data layer, oapi-codegen for REST (gateway only), protoc for gRPC. We do not hand-write transport boilerplate.

6. **Observability is not optional.** Every request gets a trace ID; every log line includes it via `LogWithTrace`. Metrics, traces, and logs are unified in Grafana. Wire it on day one — not as a Phase 2 retrofit.

7. **In-process sagas for MVP; Temporal for F6 only.** Short-lived cross-service orchestration (booking saga, refund flow) runs in-process in the orchestrating service (`booking-svc`, `payment-svc`) with explicit per-step compensation. A reconciliation cron catches mid-saga crashes. For the one long-running durable workflow — the visa pipeline (F6, runs for days, polls Saudi MOFA/Sajil, must survive restarts) — Temporal will be reintroduced when that feature is implemented. See ADR 0006.

8. **Single tenant, multi-branch.** No tenant isolation in the data layer. Branch scoping is a column on every table that needs it (`branch_id`), enforced by the service layer.

## What lives where

| Concern | Where |
|---|---|
| Edge HTTP, auth, rate limiting | `gateway-svc` |
| Identity, RBAC, audit log | `iam-svc` |
| Sellable products (packages, hotels, airlines, muthawwif) | `catalog-svc` |
| Bookings, manifests, room/bus allocation | `booking-svc` |
| Pilgrim profile, family graph, documents, OCR | `jamaah-svc` |
| Invoices, payments, gateway webhooks | `payment-svc` |
| Visa applications, e-visas, Raudhah Shield | `visa-svc` |
| Document verification queue, airport handling | `ops-svc` |
| Warehouse, procurement, kits, shipping | `logistics-svc` |
| Journals, AR/AP, tax, FX, job-order costing | `finance-svc` |
| Marketing campaigns, agents, commissions | `crm-svc` |
| Cross-service workflows (deferred; reserved for F6 visa pipeline) | `broker-svc` — not in MVP, see ADR 0006 |

See `02-service-map.md` for the full table with ports and dependencies.
