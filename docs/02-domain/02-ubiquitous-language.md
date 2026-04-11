# Ubiquitous Language

The mapping from human (Indonesian/English) terms to canonical code identifiers. The full Indonesian glossary is in `docs/00-overview/01-glossary.md`. This file focuses on the **code-side** vocabulary so every service uses the same names for the same things.

## Naming rules

- **Database columns:** `snake_case`. Booleans prefixed `is_` or `has_`. Timestamps suffixed `_at`. Foreign keys suffixed `_id`.
- **Go fields:** `PascalCase` for exported, `camelCase` for unexported.
- **Go types:** `PascalCase`, no `Type` suffix unless disambiguating.
- **JSON in API:** `snake_case` (matches DB columns).
- **Proto fields:** `snake_case` per protobuf style guide.
- **URL path segments:** `kebab-case`, plural nouns (`/v1/jamaah`, `/v1/package-departures`).
- **Service names:** `kebab-case` ending in `-svc`.
- **Database names:** `snake_case` ending in `_db`.

## Canonical identifiers

| Concept | Code identifier (Go/DB) | URL segment | Notes |
|---|---|---|---|
| Pilgrim | `jamaah` | `/jamaah` | Always `jamaah`, never `pilgrim`. The Indonesian term is the domain term. |
| Package | `package` | `/packages` | English. Avoid `paket` in code. |
| Package departure | `package_departure` | `/package-departures` | Specific departure date, the unit of inventory. |
| Booking | `booking` | `/bookings` | English. |
| Hotel | `hotel` | `/hotels` | |
| Airline | `airline` | `/airlines` | English, not `maskapai`. |
| Tour guide | `muthawwif` | `/muthawwif` | Indonesian (no clean English equivalent). |
| Visa | `visa` | `/visas` | |
| Tasreh | `tasreh` | `/tasreh` | Indonesian; specific to permission letters. |
| Mahram | `mahram` | n/a | Indonesian; relationship type, not an entity in URL. |
| Family unit | `family_unit` | `/family-units` | English. |
| Document | `document` | `/documents` | English. |
| OCR result | `ocr_result` | `/documents/{id}/ocr` | English. |
| Invoice | `invoice` | `/invoices` | English. |
| Virtual account | `virtual_account` | `/virtual-accounts` | English. |
| Payment event | `payment_event` | `/payment-events` | English; gateway webhook record. |
| Refund | `refund` | `/refunds` | English. |
| Branch | `branch` | `/branches` | English. |
| Agent | `agent` | `/agents` | English. Note: `agen` in Indonesian. |
| Commission | `commission` | `/commissions` | English. |
| Lead | `lead` | `/leads` | English. |
| Campaign | `campaign` | `/campaigns` | English. |
| Manifest | `manifest` | `/manifests` | English. |
| Luggage tag | `luggage_tag` | `/luggage-tags` | English. |
| Stock item | `stock_item` | `/stock-items` | English. |
| Warehouse | `warehouse` | `/warehouses` | English. |
| Purchase order | `purchase_order` | `/purchase-orders` | English. |
| Goods received note | `grn` | `/grns` | English abbreviation. |
| Kit | `kit` | `/kits` | English; pilgrim equipment kit. |
| Shipment | `shipment` | `/shipments` | English. |
| Journal entry | `journal_entry` | `/journal-entries` | English. |
| Chart of accounts | `chart_of_accounts` | `/chart-of-accounts` | English. |
| AR balance | `ar_balance` | `/ar-balances` | English. |
| AP balance | `ap_balance` | `/ap-balances` | English. |
| Tax record | `tax_record` | `/tax-records` | English. |
| FX rate | `fx_rate` | `/fx-rates` | English. |
| Job order cost | `job_order_cost` | `/job-order-costs` | English; cost center per departure. |
| User | `user` | `/users` | English. |
| Role | `role` | `/roles` | English. |
| Permission | `permission` | `/permissions` | English. |
| Audit log | `audit_log` | `/audit-logs` | English. |

## Status enums (the most common)

| Domain | Values |
|---|---|
| `booking.status` | `draft`, `pending_payment`, `paid_in_full`, `partially_paid`, `cancelled`, `completed` |
| `invoice.status` | `unpaid`, `partially_paid`, `paid`, `void`, `refunded` |
| `visa.status` | `waiting_docs`, `docs_ready`, `submitted`, `issued`, `rejected` |
| `payment_event.kind` | `va_created`, `payment_received`, `settlement_received`, `refund_issued` |
| `package_departure.status` | `open`, `closed`, `departed`, `completed`, `cancelled` |
| `verification_task.status` | `pending`, `in_review`, `approved`, `rejected` |
| `purchase_order.status` | `draft`, `submitted`, `approved`, `received`, `cancelled` |

Status values are always lowercase snake_case strings, never integers. They live as enum types in Postgres.

## When in doubt

- Prefer the English code identifier unless there's no clean translation (then keep the Indonesian).
- Pluralize URL segments. Singular Go types.
- If you introduce a new entity, add it here in the same session — don't defer.
