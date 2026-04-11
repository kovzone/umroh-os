# visa-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `visa.application_created` | New application | application_id, jamaah_id, booking_id | (none yet) |
| `visa.submitted` | Sent to MOFA/Sajil | application_id, provider_ref | broker-svc |
| `visa.issued` | E-visa received | application_id, visa_number | booking-svc (attach), ops-svc |
| `visa.rejected` | Rejection received | application_id, reason | crm-svc (notify customer), booking-svc |
| `visa.raudhah_alert` | Nusuk shows misuse | application_id, snapshot | ops-svc (escalate) |

## Events consumed

| Event | From | Action |
|---|---|---|
| `jamaah.documents_ready` | jamaah-svc | trigger visa pipeline workflow |
