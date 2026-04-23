---
slice: S2
title: Slice S2 — Integration Contract
status: draft
last_updated: 2026-04-23
pr_owner: Elda
reviewer: Lutfi
task_codes:
  - S2-J-01
  - S2-J-02
  - S2-J-03
  - S2-J-04
---

# Slice S2 — Integration Contract

> Slice S2 = "Payment slice" — booking checkout → VA issuance → webhook ingestion → booking status transitions. This file is the wire-level agreement between services for the payment user journey.
>
> **Incremental build:** only sections that have a landed `S2-J-*` card are filled. Amendments after merge follow the Bump-versi rule in `docs/contracts/README.md` (Changelog append for additive, `-v2.md` for breaking).

## Purpose

S2 is the payment slice. The flow: a calon jamaah submits a booking (`draft → pending_payment`, run in S1's `POST /v1/bookings/{id}/submit`), payment-svc issues a Virtual Account, the jamaah transfers money, the payment gateway POSTs a webhook, payment-svc processes it and calls booking-svc to flip booking status. Services bound in S2: `gateway-svc` (REST surface), `booking-svc` (submit saga + status receiver), `payment-svc` (invoice + VA issuance, webhook ingestion, F4 callback).

## Scope

**In scope for S2 contracts:**

- Invoice creation + Virtual Account issuance REST + gRPC interfaces — `S2-J-01`.
- Invoice state machine (full, not just draft) — `S2-J-01`.
- FX snapshot rules at invoice issuance — `S2-J-01`.
- Gateway selection + failover rules — `S2-J-01`.
- Webhook endpoints for Midtrans and Xendit — `S2-J-02`.
- `payment_events` schema freeze — `S2-J-02`.
- Webhook latency budget — `S2-J-02`.
- `MarkBookingPaid` gRPC callback from payment-svc to booking-svc — `S2-J-03`.
- Booking state transitions driven by payment events — `S2-J-03`.
- `MOCK_GATEWAY` dev toggle + internal trigger endpoint — `S2-J-04`.

**Out of scope for S2 contracts (deferred to later slices):**

- Digital receipt PDF generation + WhatsApp delivery (W3 — S3).
- Reconciliation cron (W5 — S3).
- Manual / offline payment recording (W6 — S3).
- e-Approval discount workflow (W7 — S3).
- Refund flow (W8 — S3).
- Finance journal entries / AR subledger (S3 — `S3-J-*`).
- Customer Portal billing dashboard (S3 — `S3-L-*`).
- CS payment link generator UI (S3 — `S3-L-*`).

---

## § Gateway (S2 additions)

S2 adds the following routes to `gateway-svc`'s REST surface. All routes follow the ADR 0009 single-point-REST pattern: client → REST → `gateway-svc` → gRPC → `payment-svc`. Full gateway auth rules are documented in `slice-S1.md § Gateway`.

### New routes owned by `gateway-svc` (S2 surface)

| Method | Path | Auth | Proxies to | Landing card |
| --- | --- | --- | --- | --- |
| `POST` | `/v1/invoices` | staff or booking-saga (internal) | `payment.v1.PaymentService/CreateInvoice` | `S2-J-01` |
| `GET` | `/v1/invoices/{id}` | staff | `payment.v1.PaymentService/GetInvoice` | `S2-J-01` |
| `POST` | `/v1/invoices/{id}/virtual-accounts` | staff | `payment.v1.PaymentService/IssueVirtualAccount` | `S2-J-01` |
| `POST` | `/v1/webhooks/midtrans` | **public** (signature-protected, no bearer) | `payment-svc` webhook handler (direct, no gRPC relay) | `S2-J-02` |
| `POST` | `/v1/webhooks/xendit` | **public** (signature-protected, no bearer) | `payment-svc` webhook handler (direct, no gRPC relay) | `S2-J-02` |
| `POST` | `/v1/webhooks/mock/trigger` | **dev only** (blocked in prod) | `payment-svc` mock handler | `S2-J-04` |

> **Webhook routing note.** Webhook endpoints are NOT proxied via a gRPC adapter — they are routed directly to `payment-svc`'s HTTP handler on the internal network. The signature verification + dedup logic must complete within `payment-svc` itself; `gateway-svc` acts as a transparent HTTP reverse-proxy for these paths only. No `Authorization: Bearer` is checked on webhook paths — the signature IS the authentication.

---

## § S2-J-01 — Invoice + VA Issuance Contract

*(Landed with `S2-J-01`.)*

The core payment setup flow: booking-svc's submit saga creates an invoice and issues a VA via payment-svc. Covers the REST endpoint that the saga calls through gateway, the internal gRPC method on payment-svc, the invoice state machine, FX snapshot rules, idempotency, and gateway selection.

### Call flow (submit saga)

```
booking-svc (submit saga)
  │
  │  gRPC: PaymentService/IssueVirtualAccount
  ▼
payment-svc
  │
  │  1. CreateInvoice (internal, transactional)
  │  2. IssueVA → Midtrans (primary) or Xendit (fallback)
  │  3. Persist virtual_accounts row
  ▼
returns { va_details } to booking-svc saga
```

The saga does NOT call `POST /v1/invoices` (REST) — it calls `PaymentService/IssueVirtualAccount` directly over gRPC. The REST endpoint `POST /v1/invoices` is exposed for staff-initiated invoice creation (CS generating a payment link for an existing booking without running the full submit saga).

### REST — `POST /v1/invoices`

**Auth:** staff bearer (F1 session). CS-initiated invoice creation for an already-submitted booking.

**Request body:**

```json
{
  "booking_id": "bkg_01JCDE...",
  "amount_total": 77000000,
  "currency": "IDR",
  "gateway_pref": "midtrans"
}
```

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string (ULID) | yes | Must reference an existing booking in `pending_payment` status. |
| `amount_total` | integer | yes | Payable IDR amount **after** Q001 rounding (nearest Rp 1,000 half-up). Must be ≥ 1. |
| `currency` | string | yes | MVP: must be `"IDR"`. Any other value → `422`. |
| `gateway_pref` | string enum | no | `"midtrans"` (default) or `"xendit"`. Treated as a hint; failover still applies per Q013. |

**Response — `201 Created`:**

```json
{
  "invoice": {
    "id": "inv_01JCDF...",
    "code": "INV-2026-00042",
    "booking_id": "bkg_01JCDE...",
    "status": "unpaid",
    "amount_total": 77000000,
    "paid_amount": 0,
    "rounding_adjustment_idr": 0,
    "currency": "IDR",
    "fx_snapshot": {
      "mode": "lock_rate",
      "usd_to_idr": 16250.00,
      "sar_to_idr": 4333.00,
      "locked_at": "2026-04-23T08:00:00Z"
    },
    "due_date": "2026-04-24",
    "created_at": "2026-04-23T08:00:00Z",
    "virtual_account": {
      "id": "va_01JCDG...",
      "gateway": "midtrans",
      "account_number": "8831234567890",
      "bank_code": "BCA",
      "status": "active",
      "expires_at": "2026-04-24T08:00:00Z"
    }
  }
}
```

**Errors:**

| Status | Body `error.code` | When |
| --- | --- | --- |
| `400` | `invalid_json` | Request body is not valid JSON. |
| `401` | `unauthorized` | Missing or invalid staff bearer token. |
| `403` | `forbidden` | Staff token lacks `payment.invoice.create` permission. |
| `404` | `booking_not_found` | `booking_id` does not exist. |
| `409` | `invoice_already_exists` | An active (non-void) invoice already exists for this `booking_id`. Body includes `existing_invoice_id`. |
| `422` | `validation_failed` | Field validation failure. `error.details[]` carries `{field, code, message}` per violated field. |
| `422` | `invalid_booking_status` | Booking is not in `pending_payment` status. Cannot issue invoice for a `draft`, `paid_in_full`, `cancelled`, etc. booking. |
| `502` | `gateway_unavailable` | Both Midtrans and Xendit failed (5xx or timeout). Payment-svc returns a structured error; gateway-svc maps it to 502. |
| `500` | `internal_error` | Unexpected server error. |

### REST — `GET /v1/invoices/{id}`

**Auth:** staff bearer.

**Response — `200 OK`:**

```json
{
  "invoice": {
    "id": "inv_01JCDF...",
    "code": "INV-2026-00042",
    "booking_id": "bkg_01JCDE...",
    "status": "partially_paid",
    "amount_total": 77000000,
    "paid_amount": 15400000,
    "rounding_adjustment_idr": 0,
    "currency": "IDR",
    "fx_snapshot": {
      "mode": "lock_rate",
      "usd_to_idr": 16250.00,
      "sar_to_idr": 4333.00,
      "locked_at": "2026-04-23T08:00:00Z"
    },
    "due_date": "2026-04-24",
    "created_at": "2026-04-23T08:00:00Z",
    "virtual_account": {
      "id": "va_01JCDG...",
      "gateway": "midtrans",
      "account_number": "8831234567890",
      "bank_code": "BCA",
      "status": "active",
      "expires_at": "2026-04-24T08:00:00Z"
    },
    "payment_events": [
      {
        "id": "pe_01JCDH...",
        "kind": "va_created",
        "amount": 0,
        "gateway": "midtrans",
        "gateway_txn_id": "MT-TXN-001",
        "received_at": "2026-04-23T08:00:00Z"
      },
      {
        "id": "pe_01JCDI...",
        "kind": "payment_received",
        "amount": 15400000,
        "gateway": "midtrans",
        "gateway_txn_id": "MT-TXN-002",
        "received_at": "2026-04-23T10:15:00Z"
      }
    ]
  }
}
```

**Errors:**

| Status | Body `error.code` | When |
| --- | --- | --- |
| `401` | `unauthorized` | Missing or invalid staff bearer. |
| `403` | `forbidden` | Token lacks `payment.invoice.read` permission. |
| `404` | `invoice_not_found` | No invoice matches `{id}`. |
| `500` | `internal_error` | Unexpected server error. |

### gRPC — `payment.v1.PaymentService/IssueVirtualAccount`

**Called by:** `booking-svc` submit saga only. Not exposed as a REST endpoint directly.

**Proto-style signature:**

```protobuf
rpc IssueVirtualAccount(IssueVirtualAccountRequest) returns (IssueVirtualAccountResponse);

message IssueVirtualAccountRequest {
  string booking_id    = 1; // ULID — the booking being submitted
  int64  amount_total  = 2; // payable IDR amount in integer rupiah (post-rounding)
  string currency      = 3; // must be "IDR" in MVP
  string gateway_pref  = 4; // "midtrans" | "xendit" | "" (empty = use global default)
}

message IssueVirtualAccountResponse {
  string invoice_id     = 1; // ULID of the created invoice
  string va_id          = 2; // ULID of the created virtual_account row
  string account_number = 3; // e.g. "8831234567890"
  string bank_code      = 4; // e.g. "BCA"
  string gateway        = 5; // "midtrans" | "xendit" (actual gateway used)
  int64  amount_total   = 6; // echo of the locked IDR amount
  string expires_at     = 7; // RFC3339 — VA TTL (issued_at + 24h per Q010)
}
```

**Field rules:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string | yes | ULID with `bkg_` prefix. payment-svc validates it is a known booking and status is `pending_payment`. |
| `amount_total` | int64 | yes | Integer IDR (e.g. `77000000` = Rp 77.000.000). Must be ≥ 1000 (minimum VA amount on Midtrans). |
| `currency` | string | yes | `"IDR"` only in MVP. |
| `gateway_pref` | string | no | Hint. Empty string → use global config default (`midtrans`). Failover applies regardless. |

**Failure codes (gRPC `status.Code`):**

| Code | `error.code` | When |
| --- | --- | --- |
| `INVALID_ARGUMENT` | `invalid_request` | Missing required fields, `amount_total ≤ 0`, `currency ≠ "IDR"`, malformed `booking_id`. |
| `NOT_FOUND` | `booking_not_found` | `booking_id` does not exist in booking-svc (payment-svc calls `booking.v1.BookingService/GetBooking` to verify). |
| `FAILED_PRECONDITION` | `invalid_booking_status` | Booking exists but is not in `pending_payment` status. |
| `ALREADY_EXISTS` | `invoice_already_exists` | An active invoice already exists for this `booking_id`; response message contains `existing_invoice_id`. |
| `UNAVAILABLE` | `gateway_unavailable` | Both Midtrans and Xendit returned 5xx or timed out. The saga should surface this as a terminal transient failure; the cron can retry later. |
| `INTERNAL` | `internal_error` | DB-layer failure or unexpected error. |

**Idempotency:** calling `IssueVirtualAccount` twice for the same `booking_id` — whether identical or not — returns the same invoice and VA (no second VA is created). On a replay, the response carries the original VA details. Payment-svc checks for an existing active invoice on the booking before creating a new one; if found, it returns `ALREADY_EXISTS` with the existing `invoice_id`. The saga MUST treat `ALREADY_EXISTS` as a success path (idempotent replay) and use the returned `existing_invoice_id` to proceed.

### Invoice state machine

```
unpaid
  │
  │  first payment_received event with amount ≥ amount_total
  │──────────────────────────────────────────────────────────► paid
  │
  │  first payment_received event with amount < amount_total
  │──────────────────────────────────────────────────────────► partially_paid
  │                                                                │
  │                                                                │ cumulative paid_amount ≥ amount_total
  │                                                                ▼
  │                                                              paid
  │
  │  VA expires without any payment
  │──────────────────────────────────────────────────────────► void
  │
  │  admin explicitly voids (POST /v1/invoices/{id}/void)
  │──────────────────────────────────────────────────────────► void
  │
paid ──────────────────────────────────────────────────────────► refunded
  │                                                     (on refund completion)
partially_paid ────────────────────────────────────────────────► refunded
```

**Allowed transitions (enforced by payment-svc service layer):**

| From | To | Trigger |
| --- | --- | --- |
| `unpaid` | `partially_paid` | `payment_received` event where `amount < amount_total` |
| `unpaid` | `paid` | `payment_received` event where cumulative `paid_amount ≥ amount_total` |
| `unpaid` | `void` | VA TTL expired OR admin explicit void |
| `partially_paid` | `paid` | Subsequent `payment_received` brings cumulative `paid_amount ≥ amount_total` |
| `partially_paid` | `void` | Admin explicit void only (partially-paid → void requires finance approval in later slices; MVP: admin-only gate) |
| `paid` | `refunded` | Refund saga completion (S3) |
| `partially_paid` | `refunded` | Refund saga completion (S3) |

**Terminal states:** `paid`, `void`, `refunded` — no further transitions.

**`void` cannot move to any other state.** A booking that needs a fresh invoice after void must get a new invoice row (re-issue flow, future S3).

### FX snapshot rules (Q001)

1. **At issuance time:** payment-svc calls `finance-svc.GetExchangeRate` (gRPC) to fetch the current locked USD→IDR and SAR→IDR rates. These rates are persisted as `invoices.fx_snapshot` (JSONB) with the timestamp.

2. **Rounding:** the `amount_total` stored on `invoices` is the **post-rounding IDR total** — rounded once to the nearest **Rp 1,000** (half-up). The signed difference between the pre-round total and the post-round total is stored in `invoices.rounding_adjustment_idr` for finance GL reconciliation.

   Example: pre-round total = Rp 77.350.000 → rounds to Rp 77.350.000 (already at Rp 1,000 boundary). Pre-round total = Rp 77.350.450 → rounds to Rp 77.350.000; `rounding_adjustment_idr = -450`.

3. **Immutability:** once `invoices.paid_amount > 0` (first cent of any payment lands), `invoices.fx_snapshot` and `invoices.amount_total` are **immutable**. Payment-svc enforces this at the service layer. A DB-level constraint (trigger or check) SHOULD also protect against direct writes. FX rate changes in global config after this point do NOT reprice the invoice.

4. **Multi-currency display:** catalog may show `list_currency = "USD"` prices (Q001 hybrid display). The `IssueVirtualAccount` caller (booking-svc saga) is responsible for converting the displayed USD amount to IDR using the same locked rate before passing `amount_total`. Payment-svc validates that `amount_total` is in IDR and does not perform FX conversion itself.

5. **Min DP validation (Q011):** when the booking saga invokes `IssueVirtualAccount` for a **partial payment** (DP), it MUST pass `amount_total` equal to the full invoice amount. The minimum DP (20%) is enforced at the booking-svc layer when the first `payment_received` event arrives — NOT at invoice issuance. The invoice always represents the full payable amount.

### Gateway selection (Q013)

1. **Primary:** Midtrans. Payment-svc calls Midtrans's Create VA API first.
2. **Failover:** if Midtrans returns **5xx** or the call **times out (> 10s)**, payment-svc retries once on Xendit.
3. **No failover on 4xx:** a 4xx from Midtrans (invalid amount, misconfigured bank, etc.) returns `INVALID_ARGUMENT` or `FAILED_PRECONDITION` to the caller without retrying on Xendit — the same payload would fail on Xendit too.
4. **Idempotency key to gateway:** payment-svc sends `invoice.id` as the idempotency key on every gateway call, so retries (including cross-gateway retry after failover) are safe.
5. **`gateway_pref` field:** if the caller passes `gateway_pref = "xendit"`, payment-svc starts with Xendit instead and would NOT failover to Midtrans (the preferred gateway is treated as primary for that request). Primarily used in dev/test scenarios.
6. **`MOCK_GATEWAY` mode:** when `MOCK_GATEWAY=true` (env var), neither Midtrans nor Xendit is called. See `§ S2-J-04`.

### Honored by implementation

- `S2-E-01` — payment-svc: invoice creation handler, VA issuance handler, gRPC `IssueVirtualAccount` implementation, FX snapshot logic, gateway adapter (Midtrans + Xendit).
- `S2-L-02` — checkout UI (frontend): renders VA details returned by booking-svc saga (bank code, account number, amount, expiry countdown, copy-to-clipboard).

---

## § S2-J-02 — Webhook Contract

*(Landed with `S2-J-02`.)*

Payment gateways POST asynchronous payment notifications to F5. These endpoints are the only external-facing HTTP paths in S2 that do NOT require a bearer token — authentication is via gateway signature only.

### Endpoints

| Method | Path | Gateway | Auth |
| --- | --- | --- | --- |
| `POST` | `/v1/webhooks/midtrans` | Midtrans | HMAC-SHA512 signature in `X-Callback-Token` header |
| `POST` | `/v1/webhooks/xendit` | Xendit | HMAC-SHA256 signature in `X-CALLBACK-TOKEN` header |

Both endpoints are on `gateway-svc:4000` (public) and reverse-proxied to `payment-svc`'s internal webhook handler without gRPC translation. `payment-svc` owns the full handling logic for both paths.

### Signature verification

**Midtrans:**

- Header: `X-Callback-Token`
- Algorithm: SHA512 of `(order_id + status_code + gross_amount + server_key)`
- Secret source: Viper key `MIDTRANS_SERVER_KEY`
- On bad/missing signature → `401 Unauthorized`, body `{ "error": { "code": "invalid_signature" } }`. **No business logic runs.**

**Xendit:**

- Header: `X-CALLBACK-TOKEN`
- Algorithm: The token is compared directly against the configured Xendit callback token (Xendit uses a static token model, not per-request HMAC in the webhook verification step).
- Secret source: Viper key `XENDIT_CALLBACK_TOKEN`
- On bad/missing token → `401 Unauthorized`, body `{ "error": { "code": "invalid_signature" } }`. **No business logic runs.**

### Idempotency

Idempotency key: `(gateway, gateway_txn_id)` — enforced as a `UNIQUE` constraint on `payment_events(gateway, gateway_txn_id)`.

**Duplicate webhook** (same `gateway_txn_id` for the same gateway): return `200 OK` immediately, body `{ "replayed": true }`. No business logic, no booking status update, no receipt trigger.

**First receipt** of a `gateway_txn_id`: process fully (see Webhook processing pipeline below).

### Webhook processing pipeline

The entire pipeline must complete within **< 500ms p95** end-to-end (signature verify → dedupe check → DB persist → booking-svc gRPC signal → 200 response). Receipt generation and notification are async and do NOT block the 200.

```
POST /v1/webhooks/{gateway}
  │
  1. Verify signature (< 5ms)
  │   Bad signature → 401 immediately. Stop.
  │
  2. Parse gateway payload → extract:
  │     gateway_txn_id, amount, invoice_id (from order_id / reference)
  │
  3. Dedupe check:
  │     SELECT 1 FROM payment_events
  │     WHERE gateway = $gateway AND gateway_txn_id = $txn_id
  │   Found → 200 { "replayed": true }. Stop.
  │
  4. Resolve invoice:
  │     SELECT * FROM invoices WHERE id = $invoice_id FOR UPDATE
  │   Not found → 200 (log anomaly, no business logic). Stop.
  │
  5. BEGIN TRANSACTION
  │     a. INSERT INTO payment_events (...)
  │     b. UPDATE invoices SET paid_amount = paid_amount + $amount,
  │                            status = <new_status>
  │        where new_status per state machine rules above
  │   COMMIT
  │
  6. gRPC: booking.v1.BookingService/MarkBookingPaid(
  │           booking_id, amount, invoice_status)
  │   (see § S2-J-03)
  │
  7. Async (goroutine, non-blocking):
  │     Emit receipt generation task (S3)
  │
  8. Return 200 OK
```

### Midtrans webhook payload (minimal required fields)

```json
{
  "transaction_id": "MT-TXN-002",
  "order_id": "inv_01JCDF...",
  "payment_type": "bank_transfer",
  "transaction_status": "settlement",
  "gross_amount": "15400000.00",
  "status_code": "200",
  "signature_key": "<sha512>"
}
```

Fields payment-svc extracts and maps:

| Gateway field | Mapped to |
| --- | --- |
| `transaction_id` | `payment_events.gateway_txn_id` |
| `order_id` | `payment_events.invoice_id` (must match an `invoices.id`) |
| `gross_amount` | `payment_events.amount` (decimal string → `numeric`) |
| `transaction_status` | Maps to `payment_event_kind`: `settlement` → `settlement_received`; `pending` or `capture` → `payment_received` |
| full body | `payment_events.raw_payload` (JSONB) |

### Xendit webhook payload (minimal required fields)

```json
{
  "id": "XENDIT-TXN-003",
  "external_id": "inv_01JCDF...",
  "status": "PAID",
  "paid_amount": 15400000,
  "payment_channel": "BCA"
}
```

Fields payment-svc extracts and maps:

| Gateway field | Mapped to |
| --- | --- |
| `id` | `payment_events.gateway_txn_id` |
| `external_id` | `payment_events.invoice_id` |
| `paid_amount` | `payment_events.amount` (integer IDR) |
| `status` | `PAID` → `payment_received`; `SETTLED` → `settlement_received` |
| full body | `payment_events.raw_payload` (JSONB) |

### `payment_events` schema (frozen in S2)

```sql
CREATE TABLE payment_events (
  id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  invoice_id     UUID        NOT NULL REFERENCES invoices(id),
  gateway        gateway     NOT NULL,   -- enum: 'midtrans' | 'xendit' | 'mock'
  gateway_txn_id TEXT        NOT NULL,
  kind           payment_event_kind NOT NULL,
  amount         NUMERIC(15,2) NOT NULL,
  raw_payload    JSONB       NOT NULL,
  received_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT uq_payment_events_gateway_txn UNIQUE (gateway, gateway_txn_id)
);

CREATE TYPE payment_event_kind AS ENUM (
  'va_created',         -- VA row created; amount = 0
  'payment_received',   -- gateway notifies of a payment (may precede settlement)
  'settlement_received',-- final settlement confirmation from gateway
  'refund_issued',      -- refund disbursed (S3)
  'manual'              -- CS-recorded offline payment (S3)
);

CREATE TYPE gateway AS ENUM ('midtrans', 'xendit', 'mock');
```

> `'mock'` is added to the `gateway` enum to support `MOCK_GATEWAY` dev mode (see `§ S2-J-04`). It is a valid enum value in the DB; the service layer rejects `gateway = 'mock'` on any production-path code path.

### Latency budget

The webhook pipeline (steps 1–6 above, through the `MarkBookingPaid` gRPC call return) MUST complete in **< 500ms p95**.

Breakdown targets (all p95):

| Step | Budget |
| --- | --- |
| Signature verification | < 5ms |
| Dedupe SELECT | < 10ms |
| Invoice SELECT FOR UPDATE | < 15ms |
| INSERT + UPDATE transaction | < 20ms |
| gRPC `MarkBookingPaid` (round-trip) | < 50ms |
| Other overhead (parsing, serialization, network hops) | < 400ms |
| **Total** | **< 500ms** |

If `MarkBookingPaid` takes > 50ms p95 in load testing, the booking-svc team must investigate DB indexes on `bookings.id`.

### Error responses (webhook endpoints)

| Status | `error.code` | When |
| --- | --- | --- |
| `401` | `invalid_signature` | Signature missing, malformed, or does not match. |
| `200` | — (body `{ "replayed": true }`) | Duplicate `gateway_txn_id`. |
| `200` | — (body `{ "ok": true }`) | Successful first processing. |
| `500` | `internal_error` | DB failure or unexpected error during steps 4–6. Gateway will retry; must be idempotent. |

> **Do not return `4xx` for business-logic errors** (unknown invoice, wrong amount, etc.). Always return `200` after a valid signature — log the anomaly, handle it via the reconciliation cron (S3). The gateway interprets `4xx` as a permanent failure and may stop retrying.

### Honored by implementation

- `S2-E-02` — payment-svc: Midtrans webhook handler, Xendit webhook handler, signature verification, dedupe logic, `payment_events` INSERT, invoice status update transaction.
- `S2-E-03` — payment-svc: outbound `MarkBookingPaid` gRPC call to booking-svc.

---

## § S2-J-03 — Booking Callback Contract

*(Landed with `S2-J-03`.)*

After payment-svc processes a webhook and updates the invoice, it must signal booking-svc to flip the booking status. This is a synchronous, idempotent gRPC call from payment-svc to booking-svc — no Temporal, no event bus (per ADR 0006, direct gRPC for MVP).

### gRPC — `booking.v1.BookingService/MarkBookingPaid`

**Called by:** `payment-svc` (after successful `payment_events` INSERT + invoice status update).
**Served by:** `booking-svc`.

**Proto-style signature:**

```protobuf
rpc MarkBookingPaid(MarkBookingPaidRequest) returns (MarkBookingPaidResponse);

message MarkBookingPaidRequest {
  string booking_id     = 1; // ULID
  int64  amount         = 2; // integer IDR — the amount of this specific payment event
  string invoice_status = 3; // current invoice status after this payment: "partially_paid" | "paid"
}

message MarkBookingPaidResponse {
  string booking_id     = 1;
  string booking_status = 2; // new booking status after transition
  bool   replayed       = 3; // true if this call was a no-op replay
}
```

**Field rules:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string | yes | Must reference an existing booking. |
| `amount` | int64 | yes | The IDR amount of the **current** payment event (not the cumulative total). Used for audit trail only — booking-svc does not re-compute the invoice total from this field. |
| `invoice_status` | string | yes | `"partially_paid"` or `"paid"`. Booking-svc uses this to determine the target booking status. |

### Booking status transitions triggered by `MarkBookingPaid`

| `invoice_status` received | Current booking status | New booking status |
| --- | --- | --- |
| `"partially_paid"` | `pending_payment` | `partially_paid` |
| `"paid"` | `pending_payment` | `paid_in_full` |
| `"paid"` | `partially_paid` | `paid_in_full` |
| `"partially_paid"` | `partially_paid` | `partially_paid` (no-op, `replayed = true`) |
| any | `paid_in_full` | `paid_in_full` (no-op, `replayed = true`) |
| any | `cancelled` / `expired` | `FAILED_PRECONDITION invalid_booking_status` — booking is terminal; payment-svc should log the anomaly and surface to finance dashboard |

**Min DP enforcement (Q011):** booking-svc MUST verify that the first `partially_paid` transition carries at least 20% of `bookings.total_amount`. If `amount < 0.20 × total_amount`, booking-svc returns `FAILED_PRECONDITION dp_below_minimum`. Payment-svc surfaces this as an anomaly for finance review — the invoice still records the payment event (money was received regardless), but the booking does not transition.

### Idempotency

`MarkBookingPaid` is **idempotent**. Booking-svc checks whether the booking's `paid_amount` already reflects the `amount` from this call (by tracking the associated `payment_event_id` if passed as metadata, or by comparing current status vs incoming `invoice_status`):

- If the booking is already in the target status and no further transition is needed → return `replayed = true`, current `booking_status`, `OK`.
- Booking-svc MUST NOT double-count a payment amount if the same `MarkBookingPaid` is called twice.

**Implementation note:** payment-svc should pass the `payment_event_id` as gRPC metadata (key: `x-payment-event-id`) so booking-svc can use it as a dedup key if it chooses to implement a `processed_payment_events` dedup table.

### Failure codes

| Code | `error.code` | When |
| --- | --- | --- |
| `INVALID_ARGUMENT` | `invalid_request` | Missing fields, invalid `invoice_status` value, `amount ≤ 0`. |
| `NOT_FOUND` | `booking_not_found` | `booking_id` does not exist. |
| `FAILED_PRECONDITION` | `invalid_booking_status` | Booking is in a terminal status (`cancelled`, `expired`) — cannot receive payment signal. |
| `FAILED_PRECONDITION` | `dp_below_minimum` | First partial payment is < 20% of `bookings.total_amount` (Q011). |
| `INTERNAL` | `internal_error` | DB failure. |

### No Temporal (ADR 0006)

This is a **direct gRPC call** from payment-svc to booking-svc. There is no Temporal workflow, no event bus. The call is synchronous:

- payment-svc waits for booking-svc to return before it sends `200 OK` to the gateway webhook.
- If booking-svc is unavailable: payment-svc returns `500` to the gateway. The gateway retries the webhook. On the next retry, payment-svc dedupes at step 3 of the pipeline → `200 { "replayed": true }`. The reconciliation cron (S3) also catches orphaned invoice updates.
- This means booking-svc's `/system/ready` health check MUST be tight — unhealthy booking-svc = every webhook 500 = gateway retry storm.

### `booking-svc` proto stub for payment-svc

Payment-svc MUST NOT import `booking-svc`'s pb package directly. It keeps its own vendored stub at:

```
services/payment-svc/adapter/booking_grpc_adapter/pb/booking.proto
```

The stub contains ONLY the `MarkBookingPaid` RPC.

### Honored by implementation

- `S2-E-03` — payment-svc: `booking_grpc_adapter` implementation, outbound gRPC call on webhook success.
- `S2-E-01` — booking-svc: `MarkBookingPaid` handler, booking status transition logic, min DP 20% gate, idempotency.

---

## § S2-J-04 — MOCK_GATEWAY Toggle Contract

*(Landed with `S2-J-04`.)*

Local development and CI must be able to run the full payment flow without calling real Midtrans or Xendit APIs. `MOCK_GATEWAY=true` activates an in-process mock adapter.

### Env var

| Var | Values | Effect |
| --- | --- | --- |
| `MOCK_GATEWAY` | `true` / `false` (default `false`) | When `true`, payment-svc uses the mock gateway adapter instead of Midtrans/Xendit for all VA issuance calls. |

**`MOCK_GATEWAY=true` MUST NOT be set in production.** payment-svc startup SHOULD log a `WARN` level message when it detects `MOCK_GATEWAY=true`. If `ENV=production` AND `MOCK_GATEWAY=true`, payment-svc MUST refuse to start (fatal error).

### Mock VA issuance behavior

When `MOCK_GATEWAY=true` and `IssueVirtualAccount` is called:

- Does NOT call Midtrans or Xendit.
- Returns a deterministic dummy VA:

```json
{
  "gateway": "mock",
  "bank_code": "BCA",
  "account_number": "1234567890",
  "expires_at": "<now + 24h>"
}
```

- Persists a real `virtual_accounts` row with `gateway = 'mock'` and the dummy values.
- Persists a real `payment_events` row with `kind = 'va_created'`, `gateway = 'mock'`, `gateway_txn_id = "MOCK-VA-<invoice_id>"`.
- The `invoices` row is created normally (FX snapshot, rounding — all real logic runs).

### Internal trigger endpoint (dev only)

When `MOCK_GATEWAY=true`, an additional endpoint is exposed at:

```
POST /v1/webhooks/mock/trigger
```

This endpoint is **only reachable when `MOCK_GATEWAY=true`**. In production (`MOCK_GATEWAY=false`), `gateway-svc` rejects it with `404 Not Found` before it reaches payment-svc.

**Auth:** no bearer required (dev-only endpoint). IP restriction to localhost / internal network is RECOMMENDED.

**Request body:**

```json
{
  "invoice_id": "inv_01JCDF...",
  "amount": 15400000,
  "status": "payment_received"
}
```

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `invoice_id` | string | yes | Must reference an existing invoice. |
| `amount` | integer | yes | IDR amount. Must be ≥ 1. |
| `status` | string enum | yes | `"payment_received"` or `"settlement_received"`. |

**Behavior:** payment-svc fabricates a synthetic webhook payload for the mock gateway and runs it through the **exact same webhook processing pipeline** as a real webhook (steps 3–8 in `§ S2-J-02`). This means:

- A `payment_events` row is inserted with `gateway = 'mock'`, a generated `gateway_txn_id = "MOCK-PMT-<ulid>"`.
- Invoice `paid_amount` is updated.
- `MarkBookingPaid` is called on booking-svc.
- The endpoint returns the result.

**Response — `200 OK`:**

```json
{
  "payment_event_id": "pe_01JCDH...",
  "invoice_status": "partially_paid",
  "booking_status": "partially_paid"
}
```

**Errors:**

| Status | `error.code` | When |
| --- | --- | --- |
| `400` | `mock_gateway_disabled` | `MOCK_GATEWAY=false` — should not be reachable, but defense-in-depth. |
| `404` | `invoice_not_found` | `invoice_id` does not exist. |
| `422` | `validation_failed` | `amount ≤ 0`, missing fields, unknown `status`. |
| `500` | `internal_error` | Unexpected error. |

### Documentation requirements

The following files MUST be updated when `S2-J-04` lands:

1. **`env.prod.sample`** — add `MOCK_GATEWAY=false` with comment:

   ```
   # MOCK_GATEWAY: set to true in development only.
   # When true, VA issuance skips real gateway calls and returns a dummy BCA VA.
   # MUST be false in production — payment-svc refuses to start if ENV=production and MOCK_GATEWAY=true.
   MOCK_GATEWAY=false
   ```

2. **`CONTRIBUTING.md`** — add a "Local payment testing" section explaining:

   - Set `MOCK_GATEWAY=true` in `.env.local`.
   - After creating a booking and submitting, use `POST /v1/webhooks/mock/trigger` to simulate payment.
   - Provide a `curl` example or link to the Postman collection.

### Honored by implementation

- `S2-E-01` — payment-svc: mock gateway adapter under `payment-svc/adapter/mock_gateway_adapter/`, startup guard for `ENV=production && MOCK_GATEWAY=true`.
- `S2-E-02` — payment-svc: `/v1/webhooks/mock/trigger` handler.
- `S2-L-02` — frontend: no change required; UI renders the same dummy VA details.

---

## § Booking States (S2 additions)

S2 exercises the payment-driven transitions of the booking state machine. S1 contracted only `draft`; S2 adds `pending_payment → partially_paid | paid_in_full`.

### S2 state transitions (in scope)

| From | To | Trigger | Contract |
| --- | --- | --- | --- |
| `draft` | `pending_payment` | `POST /v1/bookings/{id}/submit` — seat reservation + VA issuance saga | S1 boundary (submit endpoint contracts here; VA issuance is S2-J-01) |
| `pending_payment` | `partially_paid` | `MarkBookingPaid` with `invoice_status = "partially_paid"` | `§ S2-J-03` |
| `pending_payment` | `paid_in_full` | `MarkBookingPaid` with `invoice_status = "paid"` | `§ S2-J-03` |
| `partially_paid` | `paid_in_full` | `MarkBookingPaid` with `invoice_status = "paid"` | `§ S2-J-03` |
| `pending_payment` | `expired` | VA TTL elapses (24h per Q010); reconciliation cron signals (S3) | S3 |

### Deferred (not in S2)

- `paid_in_full → departed`, `departed → completed` — fulfillment (S3).
- `any → cancelled` with refund — cancellation flow (S3).
- `pending_payment → expired` full flow — reconciliation cron (S3).

### `bookings` table additions for S2

booking-svc adds the following columns when `S2-E-01` lands:

```sql
ALTER TABLE bookings
  ADD COLUMN invoice_id UUID NULL,         -- FK to payment_events; set when invoice issued
  ADD COLUMN total_amount NUMERIC(15,2) NULL, -- locked IDR payable at submit time
  ADD COLUMN paid_amount NUMERIC(15,2) NOT NULL DEFAULT 0; -- running sum from MarkBookingPaid calls
```

`total_amount` is set by the booking-svc submit saga after `IssueVirtualAccount` returns (it copies `invoiceResponse.amount_total`). It is **immutable** after being set.

`paid_amount` is updated by the `MarkBookingPaid` handler. It is used to evaluate the min DP gate (Q011) and to determine when the booking is `paid_in_full`.

---

## § Error envelope (shared)

All S2 REST error responses use the same envelope as S1:

```json
{
  "error": {
    "code": "<snake_case>",
    "message": "<human-readable, id-ID>",
    "trace_id": "<otel_span_hex>"
  }
}
```

Validation errors carry an additional `details` array:

```json
{
  "error": {
    "code": "validation_failed",
    "message": "Validasi gagal",
    "trace_id": "...",
    "details": [
      { "field": "amount_total", "code": "must_be_positive", "message": "Jumlah harus lebih dari 0" }
    ]
  }
}
```

`trace_id` is the OTel span ID per `docs/04-backend-conventions/03-logging-and-tracing.md`.

---

## § ID format (S2 additions)

All S2 entity IDs are **ULID** strings with a type prefix. New prefixes added in S2:

| Prefix | Entity |
| --- | --- |
| `inv_` | Invoice (`invoices.id`) |
| `va_` | Virtual account (`virtual_accounts.id`) |
| `pe_` | Payment event (`payment_events.id`) |

Consumers treat all IDs as opaque strings.

---

## § Audit trail

Every state-changing call in S2 MUST emit an `iam.audit_logs` row via `iam.v1.IamService/RecordAudit` (per F1 AC: "Every state-changing call in every service produces an audit log entry"). Minimum audit events in S2:

| Trigger | `resource` | `action` | `user_id` |
| --- | --- | --- | --- |
| Invoice created | `invoice` | `create` | Staff user (if CS-initiated) or `""` (system, for saga-initiated) |
| VA issued | `virtual_account` | `create` | Same as above |
| Webhook processed (first receipt) | `payment_event` | `create` | `""` (system) |
| `MarkBookingPaid` transitions booking | `booking` | `update_status` | `""` (system) |

---

## § Changelog

- **2026-04-23** — Initial version merged via task `S2-J-01..04`. Covers: invoice + VA issuance REST + gRPC (S2-J-01); Midtrans + Xendit webhook contracts, `payment_events` schema, 500ms latency budget (S2-J-02); `MarkBookingPaid` gRPC callback + booking state machine additions (S2-J-03); `MOCK_GATEWAY` dev toggle + `/v1/webhooks/mock/trigger` (S2-J-04).
