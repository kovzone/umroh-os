# payment-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/invoices` | List invoices (filterable) |
| `GET` | `/v1/invoices/{id}` | Get invoice detail |
| `POST` | `/v1/invoices` | Create invoice (typically called by saga, but exposed for manual ops) |
| `POST` | `/v1/invoices/{id}/virtual-accounts` | Issue VA against invoice |
| `GET` | `/v1/virtual-accounts/{id}` | Get VA detail |
| `POST` | `/v1/webhooks/midtrans` | Midtrans webhook (no auth — verified by signature) |
| `POST` | `/v1/webhooks/xendit` | Xendit webhook |
| `POST` | `/v1/refunds` | Initiate refund |
| `GET` | `/v1/refunds/{id}` | Refund detail |

## gRPC methods (planned)

`PaymentService`:
- `CreateInvoice(...)`
- `IssueVirtualAccount(...)` — saga calls this
- `GetInvoiceStatus(...)`
- `Refund(...)`
- `ReconcileInvoice(...)` — used by reconciliation cron

> Webhook endpoints validate gateway signatures via the appropriate adapter.
