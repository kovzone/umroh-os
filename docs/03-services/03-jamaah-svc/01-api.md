# jamaah-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/jamaah` | List jamaah |
| `POST` | `/v1/jamaah` | Register new jamaah |
| `GET` | `/v1/jamaah/{id}` | Get jamaah profile |
| `PATCH` | `/v1/jamaah/{id}` | Update jamaah |
| `GET` | `/v1/jamaah/{id}/documents` | List documents |
| `POST` | `/v1/jamaah/{id}/documents` | Upload document (multipart) |
| `GET` | `/v1/documents/{id}` | Get document with signed URL |
| `GET` | `/v1/documents/{id}/ocr` | Get OCR result |
| `GET` | `/v1/family-units` | List family units |
| `POST` | `/v1/family-units` | Create family unit |
| `POST` | `/v1/family-units/{id}/members` | Add jamaah with mahram relation |

## gRPC methods (planned)

`JamaahService`:
- `GetJamaah(...)`
- `BatchGetJamaah(...)` — used by ops/manifest
- `ValidateMahram(jamaah_id, group_jamaah_ids)` — returns boolean + reason
- `GetPassportData(jamaah_id)` — used by visa-svc

> All endpoints are stubs.
