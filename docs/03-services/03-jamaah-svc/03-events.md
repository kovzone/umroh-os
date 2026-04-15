# jamaah-svc — Events

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `jamaah.registered` | New jamaah | jamaah_id, branch_id | crm-svc (alumni / lead linkage) |
| `jamaah.document_uploaded` | New document | document_id, jamaah_id, kind | ops-svc (verification queue) |
| `jamaah.ocr_completed` | OCR job done | document_id, extracted | ops-svc |
| `jamaah.documents_ready` | All required docs verified | jamaah_id, booking_id | (F6 — triggers the visa pipeline when Temporal is reintroduced per ADR 0006) |

## Events consumed

| Event | From | Action |
|---|---|---|
| (none yet) | | |

> Phase 1: notional. The OCR pipeline can start as a synchronous gRPC call to GCP Vision; the worker pattern can come later.
