# Service Map

The full inventory of services in UmrohOS, their bounded contexts, ports, and key dependencies.

Per ADR 0009: **`gateway-svc` is the only REST surface**; all downstream services expose gRPC only. The REST port column below applies to `gateway-svc` only. Every other service has a gRPC port for business calls and a small admin HTTP port for Prometheus `/metrics` scraping (liveness/readiness go through `grpc.health.v1.Health`).

## Services

| # | Service | gRPC port | REST port | Admin `/metrics` | Bounded context | Owns (data) |
|---|---|---|---|---|---|---|
| 0 | `gateway-svc` | — | 4000 | 4000 (same as REST) | Edge | none — proxies client REST to downstream gRPC |
| 1 | `iam-svc` | 50051 | — | e.g. 9001 | Identity, RBAC, audit | users, roles, permissions, branches, sessions, audit_logs |
| 2 | `catalog-svc` | 50052 | — | e.g. 9002 | Master product & inventory | packages, hotels, airlines, muthawwif, itineraries, seat_inventory |
| 3 | `booking-svc` | 50053 | — | e.g. 9003 | Bookings & reservations | bookings, booking_items, room_allocations, bus_allocations, manifests |
| 4 | `jamaah-svc` | 50054 | — | e.g. 9004 | Pilgrim profile & documents | jamaah, family_graph, mahram_relations, documents, ocr_results |
| 5 | `payment-svc` | 50055 | — | e.g. 9005 | Payments & reconciliation | invoices, virtual_accounts, payment_events, refunds |
| 6 | `visa-svc` | 50056 | — | e.g. 9006 | Visa & Raudhah | visa_applications, visa_status_history, e_visas, tasreh |
| 7 | `ops-svc` | 50057 | — | e.g. 9007 | Operational handling | document_verification_queue, luggage_tags, airport_handling_events |
| 8 | `logistics-svc` | 50058 | — | e.g. 9008 | Warehouse & procurement | stock_items, warehouses, purchase_orders, grn, kits, shipments |
| 9 | `finance-svc` | 50059 | — | e.g. 9009 | PSAK accounting | journal_entries, chart_of_accounts, ar, ap, tax_records, fx_rates |
| 10 | `crm-svc` | 50060 | — | e.g. 9010 | Marketing, CRM, agents | leads, campaigns, agents, commission_ledger, broadcasts, alumni_threads |
| 11 | `broker-svc` | 50099 | — | — | Cross-service workflows — **DEFERRED (ADR 0006); reserved for F6 visa pipeline; not in MVP** | none |

Admin port numbers are proposed conventions and finalized in `BL-MON-001`. The existing REST ports (4001–4010) are retired as each backend's `api/rest_oapi/` package is removed per `BL-REFACTOR-001..010`.

## Dependency edges

These are the gRPC call directions allowed at the API surface. Cross-context **writes** are coordinated in-process by the orchestrating service (per ADR 0006), with compensations in code and a reconciliation cron catching mid-saga crashes.

| From | To | Why |
|---|---|---|
| `gateway-svc` | every service | Edge proxy. |
| every service | `iam-svc` | Token validation, permission checks. |
| `booking-svc` | `catalog-svc` | Resolve package/seat at booking time; reserve/release seats as saga coordinator. |
| `booking-svc` | `jamaah-svc` | Validate jamaah identity & mahram. |
| `booking-svc` | `payment-svc` | Issue VA during submit saga. |
| `payment-svc` | `booking-svc` | Mark booking paid (on webhook); trigger refund flow as saga coordinator. |
| `payment-svc` | `catalog-svc` | Release seats during refund saga. |
| `payment-svc` | `logistics-svc`, `finance-svc`, `crm-svc` | Emit event on payment state changes. |
| `visa-svc` | `jamaah-svc` | Read passport / document data. |
| `ops-svc` | `booking-svc`, `jamaah-svc` | Build manifests. |
| `logistics-svc` | `booking-svc` | Trigger kit fulfillment when paid. |
| `finance-svc` | `payment-svc`, `logistics-svc`, `crm-svc` | Pull events for journaling. |
| `crm-svc` | `iam-svc`, `booking-svc` | Lead → agent attribution, commission calc. |
| `broker-svc` | (deferred) | Temporal activities will call back into services when F6 reintroduces the service. |

## Cross-service workflows

### MVP — coordinated in-process (ADR 0006)

| Workflow | Orchestrator | Steps |
|---|---|---|
| Booking saga | `booking-svc` | reserve seat (catalog) → create booking (booking) → issue VA (payment); compensations on failure unwind prior steps |
| Payment reconciliation | `payment-svc` + reconciliation cron | webhook (payment) → mark booking paid (booking) → emit events (finance journal, logistics kit dispatch, crm commission) |
| Refund flow | `payment-svc` | refund request → cancel booking → release seat → release stock → reverse journal; compensations on each step failure |

### Reserved for F6 (Temporal, reintroduced when visa pipeline is implemented)

| Workflow | Orchestrator | Steps |
|---|---|---|
| Visa pipeline | `broker-svc` (Temporal) | docs ready (jamaah) → submit (visa) → long-poll MOFA/Sajil (days) → on issued, attach to booking |

## Dependency policy

- **Synchronous reads:** allowed via gRPC.
- **Synchronous writes across services:** disallowed as a single distributed transaction. The orchestrating service coordinates per-step calls with explicit compensation; a reconciliation cron catches mid-saga crashes. See ADR 0006.
- **No service may import another service's `store` package.** Always go through gRPC.
- **No service may write to another service's database.** Each service owns its tables exclusively.
