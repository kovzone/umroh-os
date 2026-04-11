# broker-svc — Events

broker-svc does not emit business events. It orchestrates other services that do. The signals it consumes are workflow signals from Temporal, not domain events.

## Signals consumed (Temporal signals, not events)

| Signal | Sent by | Workflow | Effect |
|---|---|---|---|
| `payment_received` | payment-svc on webhook | booking saga or payment reconciliation | unblocks the saga waiting for payment |
| `documents_ready` | jamaah-svc | visa pipeline | starts visa submission |
| `cancel_request` | booking-svc | any active workflow for the booking | triggers compensation |

## Activities → other services

Activities call other services synchronously via gRPC. See `01-api.md` for the activity list.
