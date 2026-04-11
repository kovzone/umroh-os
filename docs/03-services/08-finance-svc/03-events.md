# finance-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `finance.entry_recorded` | Journal entry created | entry_id, source_kind | (none yet — could feed dashboards) |
| `finance.period_closed` | Month/year close | period | none |

## Events consumed

| Event | From | Action |
|---|---|---|
| `payment.received` | payment-svc | record AR collection journal |
| `payment.refund_completed` | payment-svc | reverse journal |
| `logistics.po_approved` | logistics-svc | record AP commitment |
| `logistics.grn_recorded` | logistics-svc | record AP liability + asset |
| `crm.commission_calculated` | crm-svc | record commission expense |
| `booking.completed` | booking-svc | recognize unearned revenue |

> Phase 1: finance pulls events via gRPC from the source services on a schedule, until a real event bus exists.
