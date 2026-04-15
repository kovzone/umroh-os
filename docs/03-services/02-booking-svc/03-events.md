# booking-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `booking.created` | Draft booking exists | booking_id, package_departure_id, agent_id | crm-svc (lead → booking attribution) |
| `booking.submitted` | Draft → pending_payment | booking_id | payment-svc (in-process saga step) |
| `booking.paid_in_full` | Final payment received | booking_id | logistics-svc (kit dispatch), finance-svc (journal) |
| `booking.cancelled` | Cancellation | booking_id, reason | payment-svc (refund flow) |
| `booking.visa_attached` | Visa issued and attached | booking_id, visa_id | (none yet) |
| `booking.completed` | Departure completed | booking_id | crm-svc (alumni promotion), finance-svc (recognize revenue) |

## Events consumed

| Event | From | Action |
|---|---|---|
| `payment.received` | payment-svc | mark booking paid (via gRPC `MarkBookingPaid`) |
| `visa.issued` | visa-svc | attach visa (via gRPC `AttachVisa`) |
| `catalog.seats.released` | catalog-svc | (compensating) — saga handles |

> Phase 1: events are notional. Cross-service signaling happens via direct gRPC calls coordinated by the orchestrating service in-process (per ADR 0006).
