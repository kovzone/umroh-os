# payment-svc — Data Model

## Tables (planned)

### `invoices`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| code | text unique | human-friendly |
| booking_id | uuid | reference |
| amount_total | numeric(15,2) | contractual **IDR** payable (post–**Q001** rounding to nearest Rp 1,000 half-up) |
| paid_amount | numeric(15,2) | running sum of `payment_received` amounts; default 0 |
| rounding_adjustment_idr | numeric(15,2) | signed delta: `amount_total` − pre-round sum; posts to finance **sales rounding** GL (F9); default 0 |
| currency | char(3) | MVP: always `IDR` for VA/QRIS/card settlement (**Q001**) |
| fx_snapshot | jsonb | locked USD/SAR→IDR (and mode metadata) at **first issuance**; **immutable** after first cent of `paid_amount` lands |
| status | invoice_status enum | |
| due_date | date null | |
| created_at, updated_at | timestamptz | |

### `virtual_accounts`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| invoice_id | uuid fk | |
| gateway | gateway enum | midtrans / xendit |
| gateway_va_id | text | gateway's VA reference |
| account_number | text | |
| bank_code | text | |
| status | va_status enum | |
| expires_at | timestamptz | |
| created_at | timestamptz | |

### `payment_events`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| invoice_id | uuid fk | |
| gateway | gateway enum | |
| gateway_txn_id | text unique | dedupe key |
| kind | payment_event_kind enum | |
| amount | numeric(15,2) | |
| raw_payload | jsonb | full webhook body |
| received_at | timestamptz | |

### `refunds`
| col | type |
|---|---|
| id | uuid pk |
| invoice_id | uuid fk |
| amount | numeric(15,2) |
| status | refund_status enum |
| reason | text |
| created_at, updated_at | timestamptz |

## Enums

```sql
CREATE TYPE invoice_status AS ENUM ('unpaid', 'partially_paid', 'paid', 'void', 'refunded');
CREATE TYPE va_status AS ENUM ('active', 'paid', 'expired', 'cancelled');
CREATE TYPE gateway AS ENUM ('midtrans', 'xendit');
CREATE TYPE payment_event_kind AS ENUM ('va_created', 'payment_received', 'settlement_received', 'refund_issued');
CREATE TYPE refund_status AS ENUM ('requested', 'approved', 'processing', 'completed', 'failed');
```

## Notes

- `payment_events` is the source of truth for what the gateway told us. Replay by re-processing rows.
- `gateway_txn_id` is unique to dedupe webhook retries.
- **Naming:** use `amount_total` / `paid_amount` consistently with F5 feature spec and booking-svc copies.
- **Money precision:** IDR VA amounts match `amount_total`; internal margin/HPP lines may use higher precision in `finance-svc` — do not round per line on the invoice boundary.
