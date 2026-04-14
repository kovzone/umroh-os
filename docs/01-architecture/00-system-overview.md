# System Overview

UmrohOS is a Go microservices backend behind a public REST gateway, with a Temporal-orchestrated workflow service for cross-service business processes, and a full observability stack from day one.

## High-level diagram

```
                         ┌──────────────────────┐
                         │   Frontend (React+   │   ← deferred
                         │   Vite, future)      │
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
                         │    broker-svc        │   ← Temporal workflows
                         │    (workflows)       │     (booking saga, visa pipeline)
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

1. **One service per bounded context.** No service owns data from another's context. Cross-context reads happen via gRPC; cross-context writes happen via Temporal workflows orchestrated by `broker-svc`.

2. **Three-layer architecture per service.** API → Service → Store, with interfaces between layers. See `docs/04-backend-conventions/01-three-layer-architecture.md`. This is non-negotiable and enforced by the baseline template.

3. **Code generation everywhere.** sqlc for the data layer, oapi-codegen for REST, protoc for gRPC. We do not hand-write transport boilerplate.

4. **Observability is not optional.** Every request gets a trace ID; every log line includes it via `LogWithTrace`. Metrics, traces, and logs are unified in Grafana. Wire it on day one — not as a Phase 2 retrofit.

5. **Temporal for workflows, not for everything.** Synchronous internal calls go via gRPC. Long-running, multi-step, retryable, cross-service business processes (the booking saga, the visa pipeline, refund flows) live in `broker-svc` as Temporal workflows.

6. **Single tenant, multi-branch.** No tenant isolation in the data layer. Branch scoping is a column on every table that needs it (`branch_id`), enforced by the service layer.

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
| Cross-service workflows (Temporal) | `broker-svc` |

See `02-service-map.md` for the full table with ports and dependencies.
