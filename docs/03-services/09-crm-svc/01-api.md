# crm-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/leads` | List leads |
| `POST` | `/v1/leads` | Create lead |
| `PATCH` | `/v1/leads/{id}` | Update lead status |
| `GET` | `/v1/campaigns` | List campaigns |
| `POST` | `/v1/campaigns` | Create campaign |
| `GET` | `/v1/agents` | List agents |
| `POST` | `/v1/agents` | Onboard agent |
| `GET` | `/v1/agents/{id}` | Agent profile |
| `GET` | `/v1/agents/{id}/commissions` | Commission ledger |
| `GET` | `/v1/agents/{id}/downline` | Sub-agents |
| `POST` | `/v1/broadcasts` | Send broadcast |
| `GET` | `/v1/broadcasts/{id}` | Broadcast status |
| `GET` | `/v1/alumni/threads` | Community threads |
| `POST` | `/v1/referral-codes` | Issue referral code |
| `GET` | `/v1/ziswaf` | List ZISWAF transactions |
| `POST` | `/v1/ziswaf` | Record ZISWAF |

## gRPC methods (planned)

`CrmService`:
- `GetAgentByUserId(...)`
- `CalculateCommission(booking_id)` — called by broker-svc on booking completion
- `RecordLeadConversion(lead_id, booking_id)`

> All endpoints are stubs.
