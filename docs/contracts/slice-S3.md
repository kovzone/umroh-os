---
slice: S3
title: Slice S3 — Integration Contract
status: draft
last_updated: 2026-04-23
pr_owner: Lutfi
reviewer: Elda
task_codes:
  - S3-J-01
  - S3-J-02
  - S3-J-03
---

# Slice S3 — Integration Contract

> Slice S3 = "Fulfillment + minimum post-pay finance" — booking transitions to `paid_in_full` → kit dispatch queue created → double-entry journal posted for deferred revenue. This file is the wire-level agreement between services for the payment-to-fulfillment-to-journal user journey.
>
> **Incremental build:** only sections that have a landed `S3-J-*` card are filled. Amendments after merge follow the Bump-versi rule in `docs/contracts/README.md` (Changelog append for additive, `-v2.md` for breaking).

## Purpose

S3 is the post-payment slice. The flow: a booking moves from `paid_in_full` (established in S2) → `booking-svc` fans out to `logistics-svc` (kit dispatch queue) and `finance-svc` (deferred-revenue journal). Services bound in S3: `payment-svc` (upstream trigger), `booking-svc` (orchestrator of fan-out), `logistics-svc` (F8 fulfillment task), `finance-svc` (F9 double-entry journal).

S3 does **not** introduce a message queue or Temporal workflow. All propagation is **direct gRPC calls** per ADR-0006. The event payload shapes contracted here define the argument structures for those calls; they also document the canonical `payment.received` and `booking.paid_in_full` event fields for future async consumer compatibility.

## Scope

**In scope for S3 contracts:**

- `payment.received` and `booking.paid_in_full` event payload shapes — `S3-J-01`.
- Direct gRPC fan-out call chain from booking-svc on `paid_in_full` — `S3-J-01`.
- `logistics-svc.OnBookingPaid` gRPC signature + fulfillment task schema — `S3-J-02`.
- Fulfillment task state machine — `S3-J-02`.
- Kit dispatch gate: `paid_in_full` only, never `partially_paid` — `S3-J-02`.
- `finance-svc.OnPaymentReceived` gRPC signature — `S3-J-03`.
- Double-entry journal for payment receipt (deferred revenue) — `S3-J-03`.
- Double-entry journal for revenue recognition on departure — `S3-J-03`.
- Minimal Chart of Accounts (hardcoded MVP) — `S3-J-03`.
- Journal idempotency contract — `S3-J-03`.

**Out of scope for S3 contracts (deferred to later slices or Phase 6):**

- Refund accounting and pinalti split (Q053 — Phase 6).
- AP disbursement approval ladder (BL-FIN-010 — Phase 6 `S3-E-07`).
- AR / AP aging alerts (BL-FIN-011 — Phase 6 `S3-E-07`).
- Revenue recognition journal (trigger from departure event) — contract shape frozen here in `§ S3-J-03`; implementation in `BL-FIN-004` (`S3-E-03`).
- GRN → Auto-AP sync (`BL-FIN-002`, `BL-LOG-011` — Phase 6).
- Courier integration, self-pickup QR, size sync (BL-LOG-002..003 — `S3-E-02`).
- Digital receipt PDF + WhatsApp delivery (deferred from S2 — still deferred S3).
- Reconciliation cron for expired VAs (deferred from S2 — still deferred S3).
- Manual / offline payment recording (deferred from S2 — still deferred S3).

---

## § S3-J-01 — Event Payload Contract

*(Landed with `S3-J-01`.)*

S3 introduces two canonical event shapes emitted during the payment-to-booking-status fan-out. In MVP these are **not published to a message broker** — they travel as gRPC arguments. The shapes are contracted here so all consumers (`logistics-svc`, `finance-svc`) share the same payload contract regardless of future transport changes.

### Background: call chain

After `payment-svc` processes a gateway webhook and updates the invoice:

1. `payment-svc` calls `booking-svc.MarkBookingPaid` (contracted in S2 `§ S2-J-03`).
2. When `booking-svc` transitions a booking to `paid_in_full`, it immediately calls (synchronously, direct gRPC):
   - `logistics-svc.OnBookingPaid(...)` — trigger kit dispatch queue.
   - `finance-svc.OnPaymentReceived(...)` — post deferred-revenue journal.
3. Both downstream calls use the payloads defined below.

```
payment-svc
  │
  │  gRPC: booking.v1.BookingService/MarkBookingPaid
  ▼
booking-svc
  │   (on invoice_status = "paid" → booking transitions to paid_in_full)
  │
  ├──► gRPC: logistics.v1.LogisticsService/OnBookingPaid
  │         payload: { booking_id, departure_id, jamaah_ids[] }
  │
  └──► gRPC: finance.v1.FinanceService/OnPaymentReceived
            payload: { booking_id, invoice_id, amount, received_at }
```

> **ADR-0006 compliance.** Both fan-out calls are synchronous within the same goroutine that handles `MarkBookingPaid`. booking-svc does NOT fork goroutines for these calls — failure in either downstream must be surfaced back to the caller so the webhook pipeline can return `500` (causing gateway retry). Both calls are idempotent so retries are safe.

### Event: `payment.received`

**Emitted by:** `payment-svc` (conceptual; in MVP the payload fields map to gRPC args for `finance-svc.OnPaymentReceived`).

**Canonical payload shape:**

```json
{
  "event": "payment.received",
  "invoice_id": "inv_01JCDF...",
  "booking_id": "bkg_01JCDE...",
  "amount": 15400000,
  "invoice_status": "partially_paid",
  "received_at": "2026-04-23T10:15:00Z"
}
```

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `event` | string | yes | Literal `"payment.received"`. Used as discriminator in future async consumers. |
| `invoice_id` | string (ULID) | yes | `inv_` prefix. Reference to the invoice row. |
| `booking_id` | string (ULID) | yes | `bkg_` prefix. Reference to the booking row. |
| `amount` | int64 | yes | Integer IDR amount of **this specific payment event** — not the cumulative total. |
| `invoice_status` | string enum | yes | Current invoice status after this payment: `"partially_paid"` or `"paid"`. |
| `received_at` | string (RFC3339) | yes | Timestamp of the gateway webhook receipt (from `payment_events.received_at`). |

**Consumer rules:**

- `finance-svc` receives this payload on every `payment_received` event, regardless of `invoice_status`. Both `partially_paid` and `paid` receipts result in a Dr Bank / Cr Pilgrim Liability journal (deferred revenue is accumulated per payment, not just on full payment).
- `logistics-svc` does **not** consume `payment.received` directly. Kit dispatch is triggered only on the subsequent `booking.paid_in_full` event (see below).

### Event: `booking.paid_in_full`

**Emitted by:** `booking-svc` (on transition to `paid_in_full` status, after `MarkBookingPaid` returns with `invoice_status = "paid"`).

**Canonical payload shape:**

```json
{
  "event": "booking.paid_in_full",
  "booking_id": "bkg_01JCDE...",
  "invoice_id": "inv_01JCDF...",
  "paid_at": "2026-04-23T10:15:00Z"
}
```

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `event` | string | yes | Literal `"booking.paid_in_full"`. |
| `booking_id` | string (ULID) | yes | `bkg_` prefix. |
| `invoice_id` | string (ULID) | yes | `inv_` prefix. The invoice that reached `paid` status triggering this transition. |
| `paid_at` | string (RFC3339) | yes | Timestamp of the final payment that tipped the invoice to `paid` (equals `received_at` of the last `payment.received`). |

**Consumer rules:**

- `logistics-svc.OnBookingPaid` receives the `booking_id` from this event plus additional fields resolved by booking-svc from its own store (`departure_id`, `jamaah_ids[]`). See `§ S3-J-02`.
- `finance-svc` does **not** consume `booking.paid_in_full` directly in S3. Revenue recognition (Dr Pilgrim Liability / Cr Revenue) is triggered by a later `booking.departure_completed` event — out of scope for S3 `§ S3-J-03`.

### `MarkBookingPaid` fan-out behavior (booking-svc)

When `booking-svc.MarkBookingPaid` transitions a booking to `paid_in_full`:

1. Atomically update `bookings.status = 'paid_in_full'` in DB.
2. **Synchronously** call `logistics-svc.OnBookingPaid` (see `§ S3-J-02`).
3. **Synchronously** call `finance-svc.OnPaymentReceived` (see `§ S3-J-03`).
4. Return `MarkBookingPaidResponse` to `payment-svc`.

Both calls in steps 2–3 happen on every `MarkBookingPaid` invocation where `invoice_status = "paid"`. Both are idempotent. If either returns a non-retryable error, `MarkBookingPaid` itself returns `INTERNAL` to `payment-svc`.

**For `invoice_status = "partially_paid"` transitions**, only `finance-svc.OnPaymentReceived` is called (step 3 only). `logistics-svc.OnBookingPaid` is **NOT called** on partial payment — see `§ S3-J-02` for the gate rule.

### `payment-svc` stub for finance-svc and logistics-svc

`payment-svc` does not call `logistics-svc` or `finance-svc` directly. The fan-out is owned by `booking-svc`. `payment-svc` only calls `booking.v1.BookingService/MarkBookingPaid` (same as S2).

### Honored by implementation

- `S3-E-02` — booking-svc: `MarkBookingPaid` fan-out logic; synchronous calls to `logistics-svc.OnBookingPaid` and `finance-svc.OnPaymentReceived` on `paid_in_full` transition.
- `S3-E-02` — booking-svc: `finance-svc.OnPaymentReceived` call on every `partially_paid` or `paid` transition.

---

## § S3-J-02 — Fulfillment Task Contract

*(Landed with `S3-J-02`.)*

After `booking.paid_in_full`, `logistics-svc` creates a fulfillment task (dispatch_task) and enqueues it for warehouse processing. This section contracts the gRPC method, the task schema, and the state machine.

### Gate rule: `paid_in_full` only

Kit dispatch is **only triggered on `paid_in_full`**. Partial payment (`partially_paid`) does NOT trigger dispatch, regardless of the payment amount.

This is enforced at two layers:
1. **Caller layer:** booking-svc calls `logistics-svc.OnBookingPaid` only when transitioning to `paid_in_full` (not `partially_paid`).
2. **Callee layer:** `logistics-svc.OnBookingPaid` checks that the booking is in `paid_in_full` status via its own lookup before creating a dispatch task. If booking is not `paid_in_full`, returns `FAILED_PRECONDITION booking_not_paid_in_full`.

The fulfillment queue UI (`/erp/logistik/gudang/pengiriman`) shows only bookings with `dispatch_tasks` in `queued` or later states — bookings in `partially_paid` will never appear (F8 W10 AC: "Fulfillment queue excludes non-paid bookings").

### gRPC — `logistics.v1.LogisticsService/OnBookingPaid`

**Called by:** `booking-svc` on `paid_in_full` transition only.
**Served by:** `logistics-svc`.

**Proto-style signature:**

```protobuf
rpc OnBookingPaid(OnBookingPaidRequest) returns (OnBookingPaidResponse);

message OnBookingPaidRequest {
  string   booking_id    = 1; // ULID — bkg_ prefix
  string   departure_id  = 2; // ULID — dep_ prefix; the package departure this booking belongs to
  repeated string jamaah_ids = 3; // ULIDs of jamaah on this booking; min 1
}

message OnBookingPaidResponse {
  string dispatch_task_id = 1; // ULID of the created (or existing) dispatch_task row
  bool   replayed         = 2; // true if a dispatch_task already existed (idempotent replay)
}
```

**Field rules:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string | yes | ULID with `bkg_` prefix. logistics-svc validates it exists and is in `paid_in_full` status. |
| `departure_id` | string | yes | ULID with `dep_` prefix. logistics-svc resolves the kit definition from this departure's package. |
| `jamaah_ids` | []string | yes | At least 1 element. logistics-svc uses these to read uniform sizes (F8 W9 size sync) and reserve kit instances per jamaah. |

**Failure codes (gRPC `status.Code`):**

| Code | `error.code` | When |
| --- | --- | --- |
| `INVALID_ARGUMENT` | `invalid_request` | Missing fields, empty `jamaah_ids`, malformed ULIDs. |
| `NOT_FOUND` | `booking_not_found` | `booking_id` does not exist. |
| `NOT_FOUND` | `departure_not_found` | `departure_id` does not exist. |
| `FAILED_PRECONDITION` | `booking_not_paid_in_full` | Booking exists but is not `paid_in_full` (duplicate-call safety net). |
| `ALREADY_EXISTS` | `dispatch_task_exists` | A non-cancelled dispatch_task already exists for this `booking_id`. Response carries `existing_dispatch_task_id`. Callers MUST treat this as a success path (idempotent replay). |
| `INTERNAL` | `internal_error` | DB failure or unexpected error. |

**Idempotency:** `OnBookingPaid` is idempotent on `booking_id`. logistics-svc checks for an existing active dispatch_task before creating a new one. If found: returns `ALREADY_EXISTS` with `existing_dispatch_task_id`; booking-svc treats this as success.

### Fulfillment task schema

```sql
CREATE TABLE dispatch_tasks (
  id              TEXT        PRIMARY KEY,           -- ULID, prefix 'dtask_'
  booking_id      TEXT        NOT NULL UNIQUE,       -- bkg_ ULID; one active task per booking
  departure_id    TEXT        NOT NULL,              -- dep_ ULID
  jamaah_ids      TEXT[]      NOT NULL,              -- array of jamaah ULIDs
  status          dispatch_status NOT NULL DEFAULT 'queued',
  dispatch_method dispatch_method,                  -- set when staff selects courier vs self_pickup
  tracking_number TEXT,                             -- NULL until shipped
  shipped_at      TIMESTAMPTZ,                      -- NULL until shipped
  delivered_at    TIMESTAMPTZ,                      -- NULL until confirmed delivered
  cancelled_at    TIMESTAMPTZ,                      -- NULL unless cancelled
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE dispatch_status AS ENUM (
  'queued',       -- created on paid_in_full; awaiting warehouse processing
  'processing',   -- warehouse staff has opened the task; kit assembly underway
  'shipped',      -- shipped via courier; tracking_number populated
  'delivered',    -- courier webhook confirms delivery
  'cancelled'     -- booking cancelled before dispatch (ReleaseStock triggered)
);

CREATE TYPE dispatch_method AS ENUM ('courier', 'self_pickup');
```

**Fulfillment task fields (contracted):**

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string (ULID) | `dtask_` prefix. |
| `booking_id` | string (ULID) | `bkg_` prefix. Unique constraint — one active task per booking. |
| `departure_id` | string (ULID) | `dep_` prefix. Resolved from booking. |
| `status` | enum | `queued \| processing \| shipped \| delivered \| cancelled`. |
| `tracking_number` | string? | NULL until `shipped`. Populated by warehouse staff from courier API. |
| `shipped_at` | RFC3339? | NULL until `shipped`. |

Additional fields (`dispatch_method`, `delivered_at`, `cancelled_at`, `jamaah_ids`) are implementation fields not contracted as API surface in S3 — they are schema-frozen here so downstream consumers (ops UI, `booking-svc` status reads) can rely on their presence.

### Fulfillment task state machine

```
queued
  │
  │  warehouse staff opens and starts kit assembly
  │──────────────────────────────────────────────────────────► processing
  │                                                                │
  │  booking cancelled before task opened                         │  kit assembled + courier label printed
  │──────────────────────────────────────────► cancelled          ▼
  │                                          (ReleaseStock)     shipped ──────────► delivered
  │                                                                                (courier webhook)
  │  booking cancelled while processing
  │──────────────────────────────────────────────────────────► cancelled
                                                               (ReleaseStock, manual kit unpack)
```

**Allowed transitions (enforced by logistics-svc service layer):**

| From | To | Trigger |
| --- | --- | --- |
| `queued` | `processing` | Warehouse staff opens dispatch task in console |
| `queued` | `cancelled` | `booking-svc.ReleaseStock(booking_id)` called (booking cancellation) |
| `processing` | `shipped` | Staff confirms shipment + tracking number recorded |
| `processing` | `cancelled` | `booking-svc.ReleaseStock(booking_id)` called during assembly |
| `shipped` | `delivered` | Courier webhook `logistics.shipment_delivered` |

**Terminal states:** `delivered`, `cancelled` — no further transitions.

> `shipped → cancelled` is **not** a valid transition in S3. After shipment, cancellation becomes a returns-from-jamaah workflow (F8 W13, Phase 6). `ReleaseStock` on a `shipped` task returns `FAILED_PRECONDITION already_shipped`.

### Compensating action: `ReleaseStock`

When a booking is cancelled before the kit is dispatched, `booking-svc` calls:

```protobuf
rpc ReleaseStock(ReleaseStockRequest) returns (ReleaseStockResponse);

message ReleaseStockRequest {
  string booking_id = 1;
}

message ReleaseStockResponse {
  bool released = 1; // true if kit was released; false if no task existed (no-op)
  bool replayed = 2; // true if task was already cancelled
}
```

Behavior: transitions the dispatch_task to `cancelled`, releases any reserved kit_instance rows (`kit_instance.reserved_for_booking_id = NULL`). Idempotent — calling twice returns `replayed = true`.

`ReleaseStock` returns `FAILED_PRECONDITION already_shipped` if the task is in `shipped` or `delivered` state. Caller (booking-svc cancellation flow) must surface this to the CS operator as a "returns required" scenario (Phase 6).

### New gateway routes (S3 additions)

No new public REST routes are introduced in S3 for the fulfillment trigger — the gRPC calls are internal service-to-service only. The following internal routes are added to `logistics-svc`:

| gRPC method | Called by | Landing card |
| --- | --- | --- |
| `logistics.v1.LogisticsService/OnBookingPaid` | booking-svc | `S3-J-02` |
| `logistics.v1.LogisticsService/ReleaseStock` | booking-svc | `S3-J-02` |

The fulfillment queue REST surface (`GET /v1/dispatch-tasks`) is out of scope for S3 contracts — it is an internal ops-console endpoint contracted in `docs/03-services/07-logistics-svc/01-api.md`.

### `booking-svc` proto stub for logistics-svc

booking-svc MUST NOT import `logistics-svc`'s pb package directly. It keeps its own vendored stub at:

```
services/booking-svc/adapter/logistics_grpc_adapter/pb/logistics.proto
```

The stub contains ONLY the `OnBookingPaid` and `ReleaseStock` RPCs.

### Honored by implementation

- `S3-E-02` (`BL-LOG-001`) — logistics-svc: `OnBookingPaid` handler, dispatch_task creation, `paid_in_full` gate check.
- `S3-E-02` — booking-svc: `logistics_grpc_adapter` implementation, fan-out call on `paid_in_full`.
- `S3-E-02` — logistics-svc: `ReleaseStock` handler, dispatch_task cancellation, kit_instance release.

---

## § S3-J-03 — Minimal Journal Contract

*(Landed with `S3-J-03`.)*

`finance-svc` receives payment events from `booking-svc` and creates double-entry journal entries. This section contracts the gRPC method, the journal entry schema, idempotency rules, and the minimal Chart of Accounts.

### Minimal Chart of Accounts (MVP — hardcoded)

The following accounts are **system-seeded at startup** (`is_system_seeded = true`). They are not user-editable in MVP. All amounts are in IDR.

| Code | Account Name | Type | Notes |
| --- | --- | --- | --- |
| `1001` | Bank / Cash | Asset | Debit-normal. Represents received IDR in bank account from payment gateway settlement. |
| `2001` | Pilgrim Liability (Deferred Revenue) | Liability | Credit-normal. Holds received-but-not-yet-recognized jamaah payments per PSAK 72. |
| `4001` | Revenue — Umroh Package | Revenue | Credit-normal. Recognized at departure per Q043. |

> These three accounts are the **minimum** needed for S3 payment receipt and revenue recognition entries. Full COA (F9 W9 / module #141, per Q049) expands these in Phase 6.

### Double-entry journal rules

#### Rule 1 — Payment receipt (on every `payment.received`)

On each payment gateway confirmation (both `partially_paid` and `paid`):

```
Dr  1001  Bank / Cash             +amount
Cr  2001  Pilgrim Liability        +amount
```

> Revenue is **not recognized here**. The debit to Bank represents cash received; the credit to Pilgrim Liability (Deferred Revenue) reflects the obligation to deliver the pilgrimage service per PSAK 72 point-in-time recognition. Revenue recognition occurs at departure (Rule 2).

#### Rule 2 — Revenue recognition (on `booking.departure_completed`)

When the booking departs (triggered by `booking.departure_completed` event, Q043 trigger):

```
Dr  2001  Pilgrim Liability        +total_paid_amount
Cr  4001  Revenue — Umroh Package  +total_paid_amount
```

> Rule 2 implementation is in `BL-FIN-004` (`S3-E-03`). The gRPC method `finance-svc.RecognizeRevenue` (contracted below) is the call surface — the shape is frozen here so `booking-svc` can implement the call even before `finance-svc` fully implements the handler.

### gRPC — `finance.v1.FinanceService/OnPaymentReceived`

**Called by:** `booking-svc` on every `MarkBookingPaid` call (both `partially_paid` and `paid_in_full` transitions).
**Served by:** `finance-svc`.

**Proto-style signature:**

```protobuf
rpc OnPaymentReceived(OnPaymentReceivedRequest) returns (OnPaymentReceivedResponse);

message OnPaymentReceivedRequest {
  string booking_id   = 1; // ULID — bkg_ prefix
  string invoice_id   = 2; // ULID — inv_ prefix
  int64  amount       = 3; // integer IDR — the amount of THIS payment event (not cumulative)
  string received_at  = 4; // RFC3339 — timestamp of the gateway payment event
}

message OnPaymentReceivedResponse {
  string journal_entry_id = 1; // ULID of the created journal_entry row
  bool   replayed         = 2; // true if idempotency key matched an existing entry (no-op)
}
```

**Field rules:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string | yes | ULID with `bkg_` prefix. Used to resolve `job_order_id` (departure) for GL tagging. |
| `invoice_id` | string | yes | ULID with `inv_` prefix. Used as part of the idempotency key (see below). |
| `amount` | int64 | yes | Integer IDR. Must be > 0. |
| `received_at` | string | yes | RFC3339. Recorded as `journal_entries.entry_date`. |

**Failure codes (gRPC `status.Code`):**

| Code | `error.code` | When |
| --- | --- | --- |
| `INVALID_ARGUMENT` | `invalid_request` | Missing fields, `amount ≤ 0`, malformed ULIDs, invalid `received_at` format. |
| `NOT_FOUND` | `booking_not_found` | `booking_id` does not exist in finance-svc's store (or booking-svc lookup). |
| `INTERNAL` | `internal_error` | DB failure, unbalanced journal entry (double-entry constraint violation), unexpected error. |

### gRPC — `finance.v1.FinanceService/RecognizeRevenue`

**Called by:** `booking-svc` when departure is completed (Q043 trigger — `booking.departure_completed`).
**Served by:** `finance-svc`.

> Shape frozen in S3 for consumer compatibility. Full implementation lands in `BL-FIN-004` (`S3-E-03`).

**Proto-style signature:**

```protobuf
rpc RecognizeRevenue(RecognizeRevenueRequest) returns (RecognizeRevenueResponse);

message RecognizeRevenueRequest {
  string booking_id    = 1; // ULID — bkg_ prefix
  string departure_id  = 2; // ULID — dep_ prefix; used for job_order tagging
  string recognized_at = 3; // RFC3339 — timestamp of departure completion event
}

message RecognizeRevenueResponse {
  string journal_entry_id = 1; // ULID of the created journal_entry row
  bool   replayed         = 2; // true if idempotency key matched (no-op)
}
```

**Behavior:** finance-svc sums all committed `journal_lines` with `account_code = '2001'` and `source_id` matching `payment` entries for this `booking_id`, then posts:
```
Dr  2001  Pilgrim Liability   +total_pilgrim_liability_for_booking
Cr  4001  Revenue             +total_pilgrim_liability_for_booking
```
tagged with `job_order_id = departure_id`.

### Journal entry schema (contracted fields)

```sql
CREATE TABLE journal_entries (
  id               TEXT        PRIMARY KEY,    -- ULID, prefix 'je_'
  source_type      source_kind NOT NULL,       -- 'payment' | 'departure' | 'refund'
  source_id        TEXT        NOT NULL,       -- invoice_id (payment) | departure_id (departure) | refund_id (refund)
  idempotency_key  TEXT        NOT NULL UNIQUE, -- source_type + ':' + source_id
  entry_date       TIMESTAMPTZ NOT NULL,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE journal_lines (
  id               TEXT        PRIMARY KEY,    -- ULID, prefix 'jl_'
  journal_entry_id TEXT        NOT NULL REFERENCES journal_entries(id),
  account_code     TEXT        NOT NULL,       -- references chart_of_accounts.code
  debit            NUMERIC(15,2) NOT NULL DEFAULT 0,
  credit           NUMERIC(15,2) NOT NULL DEFAULT 0,
  currency         TEXT        NOT NULL DEFAULT 'IDR',
  fx_rate          NUMERIC(10,4) NOT NULL DEFAULT 1.0000,
  job_order_id     TEXT        -- nullable; departure_id for revenue-tagged entries
);

-- Double-entry balance constraint (enforced at application layer + optional DB trigger)
-- SUM(debit) = SUM(credit) across all lines of a journal_entry

CREATE TYPE source_kind AS ENUM (
  'payment',    -- from OnPaymentReceived call; source_id = invoice_id
  'departure',  -- from RecognizeRevenue call; source_id = departure_id
  'refund'      -- future; source_id = refund_id (Phase 6)
);
```

**Contracted field definitions:**

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string (ULID) | `je_` prefix. |
| `source_type` | enum | `payment \| departure \| refund`. Determines the accounting rule applied. |
| `source_id` | string | `invoice_id` for `payment` entries; `departure_id` for `departure` entries; `refund_id` (future) for `refund` entries. |
| `entries[]` | array of journal_lines | Each line has `account_code`, `debit`, `credit`. `debit + credit` per line: exactly one must be non-zero (a line cannot be both debit and credit). |
| `idempotency_key` | string | Computed as `source_type + ":" + source_id`. Unique constraint on `journal_entries`. |

### Idempotency rules

**Rule:** `idempotency_key = source_type + ":" + source_id`

Examples:
- Payment receipt for invoice `inv_01JCDF…`: key = `"payment:inv_01JCDF…"`
- Revenue recognition for departure `dep_01ABCD…`: key = `"departure:dep_01ABCD…"`

**Duplicate insert behavior:** if `OnPaymentReceived` or `RecognizeRevenue` is called with the same `source_id` a second time, finance-svc detects the existing `idempotency_key` in `journal_entries` and returns:
- `journal_entry_id` = the existing entry's ID.
- `replayed = true`.
- No new journal_entry or journal_lines rows are inserted.
- HTTP equivalent: `200 OK` (not `409 Conflict`) — the operation is idempotent.

**Implementation note:** finance-svc inserts with `ON CONFLICT (idempotency_key) DO NOTHING` returning the existing row, or equivalent at service layer.

### Double-entry balance enforcement

finance-svc validates `SUM(debit) = SUM(credit)` across all lines of a journal_entry before committing. Violation → transaction rollback; `OnPaymentReceived` returns `INTERNAL internal_error`. This is enforced at the service layer; a DB-level trigger SHOULD also protect against raw writes.

**No journal entry may be deleted.** finance-svc service layer rejects DELETE on `journal_entries` and `journal_lines`. Corrections via counter-entry (reversal) only — per F9 PRD line 1417.

### Journal entry examples

#### Example A — Partial payment received (Rp 15.400.000)

```json
{
  "id": "je_01JCDH...",
  "source_type": "payment",
  "source_id": "inv_01JCDF...",
  "idempotency_key": "payment:inv_01JCDF...",
  "entry_date": "2026-04-23T10:15:00Z",
  "entries": [
    { "account_code": "1001", "debit": 15400000, "credit": 0 },
    { "account_code": "2001", "debit": 0, "credit": 15400000 }
  ]
}
```

#### Example B — Final payment received (remaining Rp 61.600.000, same invoice)

> booking-svc calls `OnPaymentReceived` for each payment event — each with a unique `invoice_id`... wait: in this system, a booking may have one invoice; multiple payments reduce `paid_amount` toward `amount_total` but all reference the same `invoice_id`.

> **Clarification:** because `idempotency_key = "payment:" + invoice_id`, and there is one invoice per booking in MVP, the idempotency key would collide on the second payment. To handle multiple payment events against the same invoice, `source_id` for `payment` entries uses `payment_event_id` (from `payment_events.id`, not `invoice_id`).

**Revised rule:**

- `source_type = "payment"` → `source_id = payment_event_id` (the `pe_` ULID from `payment_events`).
- `idempotency_key = "payment:" + payment_event_id`.

This means `OnPaymentReceivedRequest` carries `payment_event_id` instead of (or in addition to) `invoice_id`.

**Updated proto:**

```protobuf
message OnPaymentReceivedRequest {
  string booking_id        = 1; // ULID — bkg_ prefix
  string invoice_id        = 2; // ULID — inv_ prefix (for audit trail)
  string payment_event_id  = 3; // ULID — pe_ prefix; used as idempotency source_id
  int64  amount            = 4; // integer IDR
  string received_at       = 5; // RFC3339
}
```

`idempotency_key = "payment:" + payment_event_id`

booking-svc passes the `payment_event_id` from the `x-payment-event-id` gRPC metadata header (established in S2 `§ S2-J-03`) when calling `finance-svc.OnPaymentReceived`.

#### Example B (corrected) — Final payment for same invoice

```json
{
  "id": "je_01JCDI...",
  "source_type": "payment",
  "source_id": "pe_01JCDI...",
  "idempotency_key": "payment:pe_01JCDI...",
  "entry_date": "2026-04-23T12:00:00Z",
  "entries": [
    { "account_code": "1001", "debit": 61600000, "credit": 0 },
    { "account_code": "2001", "debit": 0, "credit": 61600000 }
  ]
}
```

#### Example C — Revenue recognition on departure

```json
{
  "id": "je_01JCDJ...",
  "source_type": "departure",
  "source_id": "dep_01ABCD...",
  "idempotency_key": "departure:dep_01ABCD...",
  "entry_date": "2026-05-15T04:30:00Z",
  "entries": [
    { "account_code": "2001", "debit": 77000000, "credit": 0 },
    { "account_code": "4001", "debit": 0, "credit": 77000000 }
  ]
}
```

> Note: `departure_id` is used as `source_id` for revenue recognition — one entry per departure per booking. `RecognizeRevenue` is called once per booking when the departure completes; the amount is the sum of all `2001` credits accumulated for that booking.

### `booking-svc` proto stub for finance-svc

booking-svc MUST NOT import `finance-svc`'s pb package directly. It keeps its own vendored stub at:

```
services/booking-svc/adapter/finance_grpc_adapter/pb/finance.proto
```

The stub contains ONLY the `OnPaymentReceived` and `RecognizeRevenue` RPCs.

### Honored by implementation

- `S3-E-03` (`BL-FIN-001`) — finance-svc: `OnPaymentReceived` handler, Dr Bank / Cr Pilgrim Liability journal, idempotency gate, COA seed data.
- `S3-E-03` (`BL-FIN-003`) — finance-svc: double-entry journal engine, balance enforcement, idempotent source constraint.
- `S3-E-03` (`BL-FIN-004`) — finance-svc: `RecognizeRevenue` handler, Dr Pilgrim Liability / Cr Revenue journal (at departure).
- `S3-E-02` — booking-svc: `finance_grpc_adapter` implementation; fan-out call on every `MarkBookingPaid`; passing `payment_event_id` from `x-payment-event-id` metadata.

---

## § Booking States (S3 additions)

S3 does not introduce new booking status values. It exercises the existing `paid_in_full` state (established in S2) as the trigger for downstream fan-out. Post-S3 transitions (`departed`, `completed`, `cancelled`) are deferred.

| From | To | Trigger | Contract |
| --- | --- | --- | --- |
| `paid_in_full` | _(no S3 transition)_ | — | S3 uses `paid_in_full` as input only; no booking status changes in S3 scope |

---

## § ID format (S3 additions)

All S3 entity IDs follow the ULID-with-prefix convention from S1/S2. New prefixes added in S3:

| Prefix | Entity |
| --- | --- |
| `dtask_` | Dispatch task (`dispatch_tasks.id`) |
| `je_` | Journal entry (`journal_entries.id`) |
| `jl_` | Journal line (`journal_lines.id`) |

---

## § Audit trail (S3 additions)

Every state-changing call in S3 MUST emit an `iam.audit_logs` row via `iam.v1.IamService/RecordAudit` (per F1 AC). Minimum audit events in S3:

| Trigger | `resource` | `action` | `user_id` |
| --- | --- | --- | --- |
| `OnBookingPaid` creates dispatch_task | `dispatch_task` | `create` | `""` (system) |
| `ReleaseStock` cancels dispatch_task | `dispatch_task` | `update_status` | Staff user (if CS-initiated cancel) or `""` (system) |
| `OnPaymentReceived` creates journal_entry | `journal_entry` | `create` | `""` (system) |
| `RecognizeRevenue` creates journal_entry | `journal_entry` | `create` | `""` (system) |

---

## § Error envelope (shared)

All S3 REST error responses (if any are added) use the same envelope as S1/S2:

```json
{
  "error": {
    "code": "<snake_case>",
    "message": "<human-readable, id-ID>",
    "trace_id": "<otel_span_hex>"
  }
}
```

`trace_id` is the OTel span ID per `docs/04-backend-conventions/03-logging-and-tracing.md`.

---

## § Changelog

- **2026-04-23** — Initial version drafted via task `S3-J-01..03`. Covers: `payment.received` + `booking.paid_in_full` event payload shapes + gRPC fan-out call chain (S3-J-01); `OnBookingPaid` gRPC + fulfillment task schema + dispatch state machine + `ReleaseStock` compensating action (S3-J-02); `OnPaymentReceived` + `RecognizeRevenue` gRPC + journal entry schema + minimal COA + idempotency contract (S3-J-03).
