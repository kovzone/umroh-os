# finance-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/chart-of-accounts` | List COA tree |
| `POST` | `/v1/chart-of-accounts` | Add account |
| `GET` | `/v1/journal-entries` | List entries |
| `POST` | `/v1/journal-entries` | Manual entry |
| `GET` | `/v1/journal-entries/{id}` | Detail |
| `GET` | `/v1/ar-balances` | AR aging |
| `GET` | `/v1/ap-balances` | AP aging |
| `GET` | `/v1/tax-records` | Tax records |
| `GET` | `/v1/fx-rates` | FX rates |
| `POST` | `/v1/fx-rates` | Add daily snapshot |
| `GET` | `/v1/job-order-costs/{departure_id}` | Per-departure P&L |
| `GET` | `/v1/reports/balance-sheet` | Generate balance sheet |
| `GET` | `/v1/reports/profit-loss` | Generate P&L |
| `GET` | `/v1/reports/cash-flow` | Cash flow statement |

## gRPC methods (planned)

`FinanceService`:
- `RecordJournalEntry(...)` — used by other services to push events for journaling
- `GetExchangeRate(...)`

> All endpoints are stubs.
