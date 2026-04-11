# crm-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `crm.lead_converted` | Lead → booking | lead_id, booking_id, agent_id | finance-svc |
| `crm.commission_calculated` | After booking completed | agent_id, booking_id, amount | finance-svc |
| `crm.broadcast_completed` | Broadcast finished | broadcast_id, sent, failed | none |

## Events consumed

| Event | From | Action |
|---|---|---|
| `booking.created` | booking-svc | attribute lead → booking, attach agent |
| `booking.completed` | booking-svc | calculate commission via gRPC `CalculateCommission` |
| `ops.document_rejected` | ops-svc | notify customer via WhatsApp |
| `visa.issued` | visa-svc | notify customer |
| `visa.rejected` | visa-svc | notify customer + escalate |
| `logistics.shipment_dispatched` | logistics-svc | notify customer |
