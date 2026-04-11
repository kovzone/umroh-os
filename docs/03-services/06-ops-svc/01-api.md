# ops-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/verification-tasks` | List pending verification |
| `POST` | `/v1/verification-tasks/{id}/approve` | Approve doc |
| `POST` | `/v1/verification-tasks/{id}/reject` | Reject with reason |
| `POST` | `/v1/manifests/generate` | Generate manifest for a departure |
| `GET` | `/v1/manifests/{id}` | Get manifest |
| `GET` | `/v1/manifests/{id}/download` | Download (PDF/Excel) |
| `POST` | `/v1/grouping/run` | Run smart room/bus allocation |
| `GET` | `/v1/luggage-tags` | List tags for a departure |
| `POST` | `/v1/luggage-tags/issue` | Issue tags for a booking |
| `POST` | `/v1/handling-events` | Record a scan event |

## gRPC methods (planned)

`OpsService`:
- `RunSmartGrouping(...)` — used by booking-svc on submit
- `BuildManifest(...)` — used by broker-svc / scheduler

> All endpoints are stubs.
