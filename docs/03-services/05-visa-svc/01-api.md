# visa-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/visas` | List visa applications |
| `GET` | `/v1/visas/{id}` | Get application detail |
| `POST` | `/v1/visas` | Create application (typically by saga) |
| `GET` | `/v1/visas/{id}/status-history` | Status transition log |
| `GET` | `/v1/visas/{id}/e-visa` | Download e-visa (signed URL) |
| `GET` | `/v1/tasreh` | List tasreh records |
| `POST` | `/v1/tasreh` | Register tasreh |
| `GET` | `/v1/raudhah-monitoring` | Latest Nusuk snapshot |

## gRPC methods (planned)

`VisaService`:
- `CreateApplication(...)`
- `SubmitApplication(...)` — calls MOFA/Sajil
- `GetApplicationStatus(...)`
- `AttachEVisa(...)` — uploads issued visa
- `PollNusukStatus(...)` — used by Raudhah Shield cron

> All endpoints are stubs.
