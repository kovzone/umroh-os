# Service Map

The full inventory of services in UmrohOS, their bounded contexts, ports, and key dependencies.

## Services

| # | Service | Style | Port | Bounded context | Owns (data) |
|---|---|---|---|---|---|
| 0 | `gateway-svc` | REST (Fiber) | 4000 | Edge | none — proxies to internal services |
| 1 | `iam-svc` | REST + gRPC | 4001 / 50051 | Identity, RBAC, audit | users, roles, permissions, branches, sessions, audit_logs |
| 2 | `catalog-svc` | REST + gRPC | 4002 / 50052 | Master product & inventory | packages, hotels, airlines, muthawwif, itineraries, seat_inventory |
| 3 | `booking-svc` | REST + gRPC | 4003 / 50053 | Bookings & reservations | bookings, booking_items, room_allocations, bus_allocations, manifests |
| 4 | `jamaah-svc` | REST + gRPC | 4004 / 50054 | Pilgrim profile & documents | jamaah, family_graph, mahram_relations, documents, ocr_results |
| 5 | `payment-svc` | REST + gRPC | 4005 / 50055 | Payments & reconciliation | invoices, virtual_accounts, payment_events, refunds |
| 6 | `visa-svc` | REST + gRPC | 4006 / 50056 | Visa & Raudhah | visa_applications, visa_status_history, e_visas, tasreh |
| 7 | `ops-svc` | REST + gRPC | 4007 / 50057 | Operational handling | document_verification_queue, luggage_tags, airport_handling_events |
| 8 | `logistics-svc` | REST + gRPC | 4008 / 50058 | Warehouse & procurement | stock_items, warehouses, purchase_orders, grn, kits, shipments |
| 9 | `finance-svc` | REST + gRPC | 4009 / 50059 | PSAK accounting | journal_entries, chart_of_accounts, ar, ap, tax_records, fx_rates |
| 10 | `crm-svc` | REST + gRPC | 4010 / 50060 | Marketing, CRM, agents | leads, campaigns, agents, commission_ledger, broadcasts, alumni_threads |
| 11 | `broker-svc` | gRPC + Temporal | 4099 / 50099 | Cross-service workflows | none — orchestrates other services |

## Dependency edges

These are the gRPC call directions allowed at the API surface. Anything not listed should not exist; cross-context concerns should go through `broker-svc` workflows.

| From | To | Why |
|---|---|---|
| `gateway-svc` | every service | Edge proxy. |
| every service | `iam-svc` | Token validation, permission checks. |
| `booking-svc` | `catalog-svc` | Resolve package/seat at booking time. |
| `booking-svc` | `jamaah-svc` | Validate jamaah identity & mahram. |
| `payment-svc` | `booking-svc` | Mark booking paid; trigger fulfillment. |
| `visa-svc` | `jamaah-svc` | Read passport / document data. |
| `ops-svc` | `booking-svc`, `jamaah-svc` | Build manifests. |
| `logistics-svc` | `booking-svc` | Trigger kit fulfillment when paid. |
| `finance-svc` | `payment-svc`, `logistics-svc`, `crm-svc` | Pull events for journaling. |
| `crm-svc` | `iam-svc`, `booking-svc` | Lead → agent attribution, commission calc. |
| `broker-svc` | every service | Temporal activities call back into services. |

## Workflows owned by `broker-svc`

| Workflow | Steps |
|---|---|
| Booking saga | reserve seat (catalog) → create booking (booking) → issue VA (payment) → notify (crm) |
| Payment reconciliation | webhook (payment) → mark booking paid (booking) → trigger kit dispatch (logistics) → journal entry (finance) |
| Visa pipeline | docs ready (jamaah) → submit (visa) → poll status → on issued, attach to booking (booking) |
| Refund flow | refund request (payment) → cancel booking (booking) → release seat (catalog) → release stock (logistics) → reverse journal (finance) |

## Dependency policy

- **Synchronous reads:** allowed via gRPC.
- **Synchronous writes across services:** disallowed. Use a Temporal workflow.
- **No service may import another service's `store` package.** Always go through gRPC.
- **No service may write to another service's database.** Each service owns its tables exclusively.
