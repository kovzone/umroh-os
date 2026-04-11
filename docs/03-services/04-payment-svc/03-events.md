# payment-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `payment.invoice_created` | Invoice created | invoice_id, booking_id, amount | finance-svc |
| `payment.va_issued` | VA created at gateway | va_id, invoice_id, account_number | none directly (returned via saga to caller) |
| `payment.received` | Webhook reports payment | invoice_id, amount, gateway_txn_id | booking-svc (mark paid), finance-svc (journal) |
| `payment.settled` | Settlement reported | invoice_id, settlement_id | finance-svc |
| `payment.refund_completed` | Refund processed | refund_id, invoice_id, amount | booking-svc, finance-svc |

## Events consumed

| Event | From | Action |
|---|---|---|
| `booking.cancelled` | booking-svc | trigger refund flow if applicable |

> Phase 1: payment-svc emits via Temporal saga signals coordinated by broker-svc.
