# ops-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `ops.document_verified` | Reviewer approves | document_id, jamaah_id | jamaah-svc (mark verified) |
| `ops.document_rejected` | Reviewer rejects | document_id, reason | crm-svc (notify customer) |
| `ops.manifest_generated` | Manifest produced | manifest_id, departure_id | (none yet) |
| `ops.luggage_scanned` | Airport scan | tag_code, event_type | (none yet — could feed dashboards) |

## Events consumed

| Event | From | Action |
|---|---|---|
| `jamaah.document_uploaded` | jamaah-svc | enqueue verification task |
| `jamaah.ocr_completed` | jamaah-svc | attach OCR to verification task |
