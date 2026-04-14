---
id: F5
title: Payment & Reconciliation
status: draft
last_updated: 2026-04-15
moscow_profile: 4+ Must Haves (VA issuance, bank integration, digital receipts, AR subledger)
prd_sections:
  - "A. Self-Service Booking — Payment Gateway B2C (line 79)"
  - "C.4 Sales Closing — Pembuat Link Pembayaran (line 223)"
  - "G.1 Revenue & Receivables (lines 431–443)"
  - "E.5 Cancellation — Refund & Pinalti (line 341)"
  - "Alur Logika 5.3 (line 1175), 6.5 (line 1193)"
  - "Integrasi API config (lines 1255–1325)"
modules:
  - "#18 Payment Gateway B2C, #64 Pembuat Link Pembayaran, #129 Penagihan Otomatis, #130 Integrasi Bank, #131 Buku Pembantu Piutang, #132 Kwitansi Digital, #133 Pembayaran Manual, #100 Refund & Pinalti"
depends_on: [F1, F4]
open_questions:
  - Q001 — operating currency, FX handling (existing — FX-lock at invoice issuance)
  - Q004 — cancellation → seat return timing (existing — interlocks with refund flow)
  - Q011 — minimum DP %, max installment count, cadence
  - Q012 — refund penalty policy matrix (per package type, per timing bucket)
  - Q013 — dual-gateway (Midtrans / Xendit) selection and fallback rule
---

# F5 — Payment & Reconciliation

## Purpose & personas

Turns booking commitments into money in the bank. Owns invoices, virtual accounts, gateway webhook ingestion, digital receipts, and refunds. The **only** service that talks to payment gateways (Midtrans, Xendit) and direct bank APIs. Emits events that F4 (booking status), F8 (kit dispatch), F9 (accounting), and F10 (commission confirmation) depend on.

The interlock with F4 is tight: F4 submit triggers F5 VA issuance; F5 webhook triggers F4 status update. The saga that coordinates this lives in `broker-svc`.

Primary personas:
- **Calon jamaah** — pays via VA / QRIS / card; receives digital receipt; sees billing dashboard on Customer Portal.
- **Agent** — receives commission payout (F10 handles calculation; F5 handles the disbursement transaction).
- **CS** — generates payment links via the internal console; handles manual/cash payments (module #133 Should Have).
- **Finance admin** — reconciles gateway settlement reports against booked receipts, handles refund approvals, monitors AR aging.
- **Downstream consumers** — F4 (booking status), F9 (journal entries), F8 (kit dispatch trigger on lunas), F10 (commission Pending→Confirmed).

## Sources

- PRD Section A payment gateway (line 79)
- PRD Section C.4 payment link generator (line 223), e-Approval discount (line 225)
- PRD Section G.1 Revenue & Receivables (lines 431–443) — **the authoritative section for this feature**
- PRD Section E.5 cancellation + refund (line 341)
- PRD Alur Logika 5.3 (rekonsiliasi + kwitansi) and 6.5 (refund flow)
- PRD Section Integrasi API (lines 1255–1325) — gateway config, webhook URL, 2FA rules
- Modules #18, #64, #100, #129, #130, #131, #132, #133

## User workflows

### W1 — Issue a virtual account at booking checkout

Called from the F4 booking submit saga.

1. Saga activity `IssueVirtualAccount` is invoked with `{ booking_id, amount, currency, gateway_pref }`.
2. F5 creates an `invoice` in state `unpaid` with `booking_id`, `amount_total`, `currency`, `fx_snapshot` (Q001 FX lock).
3. F5 selects gateway per Q013 (default: Midtrans primary, Xendit fallback on Midtrans failure).
4. Calls the gateway's Create VA API with idempotency key = `invoice.id`, receives `gateway_va_id`, `account_number`, `bank_code`, `expires_at` (from Q010 default TTL).
5. Persists `virtual_accounts` row; returns `{ account_number, bank_code, amount, expires_at }` to the saga.
6. Saga passes the VA details back to the caller (B2C checkout page / CS WhatsApp message / B2B agent portal).

### W2 — Pay via virtual account (Alur Logika 5.3)

1. Jamaah transfers money to the issued VA via their bank app / internet banking.
2. Gateway receives the payment, POSTs a webhook to `api.<domain>/v1/webhooks/<gateway>` with a signed payload.
3. F5 webhook handler **verifies the signature** (HMAC per gateway spec) — rejects unsigned or badly-signed webhooks with 401 (no retry).
4. F5 checks idempotency — if `gateway_txn_id` already exists in `payment_events`, no-op 200 so the gateway stops retrying.
5. Inserts `payment_events` row with `kind = 'payment_received'`, `amount`, `raw_payload`, `received_at`.
6. Updates the `invoice.paid_amount += amount`; if now ≥ `amount_total`, invoice status → `paid`; else `partially_paid`.
7. Signals the F4 saga: `broker.SignalPaymentReceived(booking_id, amount, invoice_status)`.
8. F4 side transitions booking status accordingly (`pending_payment → partially_paid | paid_in_full`).
9. F5 triggers digital receipt (W3) asynchronously.
10. Returns 200 OK to the gateway.

**Total latency budget**: signature verify + idempotency check + event insert + F4 signal should complete in < 500ms p95 so the gateway doesn't time out and retry.

### W3 — Digital receipt delivery (Module #132, line 441)

1. On every `payment_event` of kind `payment_received`, F5 async-triggers receipt generation.
2. Fetches invoice + booking + jamaah context (via F4 `GetBooking` and F3 `GetJamaah`).
3. Renders a PDF receipt (server-side via headless browser worker, same pool as the flyer generator from F2 W7).
4. Uploads PDF to GCS under `receipts/<invoice_id>/<payment_event_id>.pdf`.
5. Delivers via **WhatsApp** (primary; per PRD line 441: "sedetik setelah dana terdeteksi masuk") with a short text summary + PDF attachment.
6. Falls back to email if the WhatsApp send fails 3x.

### W4 — CS-generated payment link (Module #64, line 223)

1. CS opens a booking in the internal console, clicks **Generate Payment Link**.
2. Selects gateway, amount (full / DP per Q011 / custom with e-Approval per W7), and expiry.
3. F5 creates an invoice + VA same as W1 but without the saga detour — direct REST call since the booking already exists.
4. Returns a WhatsApp-shareable URL that opens the checkout page.
5. CS pastes the link into the WA conversation.

### W5 — Reconciliation cron (Module #131, lines 431–443)

Runs hourly. Catches missed webhooks and gateway settlement lag.

1. For every `invoice` with status `unpaid | partially_paid` and a non-expired VA, F5 queries the gateway's GetPaymentStatus API.
2. If gateway reports paid but F5 has no `payment_events` row: inserts the event (backfill). Signals F4 as in W2.
3. If gateway reports paid with amount higher than F5's `paid_amount` (partial webhook lost): reconciles.
4. If gateway reports expired/cancelled: F5 marks the invoice `void` and the VA `expired`; signals F4 which transitions the booking to `expired`.
5. Every reconciliation run writes a summary to the finance admin's dashboard (module #131 AR sub-ledger view).

### W6 — Manual / offline payment (Module #133, Should Have)

1. CS confirms a jamaah paid in cash or via direct transfer (without VA).
2. CS records the payment in the Manual Payment screen: amount, method (`cash | direct_transfer | check`), reference number, attached receipt photo.
3. F5 creates a `payment_events` row with `kind = 'manual'` and `approval_status = 'pending'`.
4. Finance admin reviews and approves in the approval queue.
5. On approval: invoice updates, F4 signalled, booking status transitions as in W2. On reject: event stays for audit, booking unaffected.

### W7 — e-Approval discount (Module #65, line 225)

1. CS or agent requests a discount on a booking at quote time.
2. Request record created with requested amount + reason + approver role.
3. Manager with `discount_approve` permission reviews; approves / rejects.
4. On approve: discount is applied to the booking's total; if invoice already issued, a new `invoice_adjustment` row nets against the outstanding amount.
5. Audit trail captures requester, approver, amount, timestamp, reason.

### W8 — Refund (Module #100, line 341; Alur Logika 6.5)

Triggered by F4 cancellation (W7 in F4) or by a departure-cancelled bulk operation.

1. F4 calls `broker.StartRefundFlow(booking_id, reason_code)`.
2. Refund flow workflow in broker-svc:
   - Reads the booking, its invoice(s), and all `payment_events` of kind `payment_received`.
   - Calculates net refund: `sum(received) - penalty(per Q012) - sunk_costs(issued_tickets + filed_visas)`.
   - Creates a `refund` record with status `requested`, shows a preview to the initiator for confirmation (F4 UI shows this — jamaah sees the net before confirming cancel).
   - On confirm: calls gateway's Refund API (for VA/QRIS/card) OR creates a manual-refund task (for cash receipts that need bank transfer out).
   - On gateway ack: status → `processing`; on final settlement: `completed`.
   - On failure: `failed` with retry after admin review.
3. Refund issues a negative `payment_events` row (`kind = 'refund_issued'`) for audit symmetry.
4. Signals F4 for booking state, F9 for reverse journal, F8 for kit-release-if-applicable.

### W9 — FX handling (per Q001)

1. At VA issuance time (W1), F5 snapshots the current FX rates for every foreign-currency cost component and stores them on the invoice.
2. Once the invoice is in state `partially_paid` or `paid`, the snapshot is **immutable** — subsequent FX rate changes in global config do not reprice this invoice (PRD line 1313 hard rule).
3. For a booking that's still `draft` or `pending_payment` with an unpaid VA, the admin can choose to void the old invoice and issue a new one at fresh rates if a customer complaint warrants it. _(Inferred — voluntary; not PRD-mandated.)_

## Acceptance criteria

- VA issuance is idempotent — calling `IssueVirtualAccount` twice for the same invoice returns the same VA, not a second one.
- Webhook handlers verify signatures; unsigned / bad-signature payloads are rejected with 401 and never touched business logic.
- Webhook idempotency is enforced via unique constraint on `(gateway, gateway_txn_id)`. Retries are no-ops returning 200.
- Webhook end-to-end latency (signature → dedupe → persist → F4 signal → 200) is < 500ms p95.
- Digital receipt (PDF + WA send) ships within 60s of `payment_received` event 99% of the time.
- Reconciliation cron runs hourly and catches missed webhooks; a simulated webhook drop in the test suite must be recovered within the next reconciliation cycle.
- Refund flow preview shows the exact net-amount calculation (received − penalty − sunk costs) to the initiator BEFORE confirmation.
- FX snapshot on an invoice is never mutated after the first payment lands (DB check constraint or trigger).
- Manual payment requires finance-admin approval before booking status transitions.
- Every state-changing call writes to `iam.audit_logs` via F1 `RecordAudit`.

## Edge cases & error paths

- **Webhook replay from gateway retry.** First receipt processes fully; subsequent receipts hit the unique constraint on `gateway_txn_id` and return 200 immediately. No booking state change on retries.
- **Webhook arrives for an unknown invoice.** Possible if the invoice was created, VA issued, but the invoice record was deleted manually (shouldn't happen but defense in depth). Return 200 to the gateway, log the anomaly, finance dashboard surfaces as an orphan payment.
- **Gateway returns success on Create VA but F5 DB insert fails** (network blip). F5 has the gateway's `gateway_va_id` but no DB row — orphan. Saga compensation: on next invocation, F5's idempotency key matches, gateway returns the existing VA, F5 retries the DB insert. Net: eventually consistent, no double-VA.
- **Partial webhook (first tranche of a multi-installment lands, second is lost).** Reconciliation cron (W5) catches it, backfills the missing event, updates invoice `paid_amount`.
- **Payment lands on an already-cancelled booking** (race between cancellation and payment settlement). Invoice is `void` but payment still came in. Log anomaly, auto-trigger refund flow for the received amount (negative penalty — full return because cancel preceded payment). Finance reviews the exceptional-case dashboard.
- **Multi-currency invoice + FX rate change mid-flight** (Q001). FX snapshot on invoice locks the rate; invoice paid_amount is always tracked in the invoice's original currency. Finance conversion to base currency for reports uses the snapshot rate, not current.
- **Refund fails on gateway side** (e.g. VA expired, card deactivated). Falls back to manual refund task routed to finance admin. Booking stays in `cancelled` state; refund record status → `failed` → admin manually transfers via bank + marks completed.
- **Partial cancellation refund (Q014)** — individual jamaah cancels from a multi-pilgrim booking. Refund = (per-jamaah share of received − per-jamaah penalty − per-jamaah sunk). Booking `total_amount` and `paid_amount` reduce accordingly. F4 handles the booking-shape change; F5 handles the money.
- **Manual-payment approval rejected after booking was already flipped** (shouldn't happen — booking only flips AFTER approval — but defense in depth). Approval is strictly gating; no booking transition without it.
- **Settlement takes days to arrive** (card payments with delayed settlement). Invoice can show `paid` from the webhook while gateway-side settlement is pending. Finance reports distinguish `received` (webhook) from `settled` (gateway confirmation).

## Data & state implications

Owned by `payment-svc`. Full schema in `docs/03-services/04-payment-svc/02-data-model.md`. Key additions from this spec:

- `invoices.fx_snapshot` jsonb — FX rates at issuance; immutable once any payment lands.
- `invoices.status` enum — `unpaid | partially_paid | paid | void | refunded`.
- `invoices.paid_amount` numeric — running tally, updated atomically on each event.
- `virtual_accounts.gateway` enum — `midtrans | xendit` (extendable).
- `payment_events.gateway_txn_id` with `UNIQUE(gateway, gateway_txn_id)` — webhook idempotency key.
- `payment_events.kind` enum — `va_created | payment_received | settlement_received | refund_issued | manual`.
- `refunds.status` enum — `requested | processing | completed | failed`.
- New `invoice_adjustments` — discount approvals, manual corrections; separate from `payment_events` to keep that table immutable and pure.

## API surface (high-level)

Full contracts in `docs/03-services/04-payment-svc/01-api.md`.

**REST:**
- `GET /v1/invoices` — filter by booking, status, date range
- `GET /v1/invoices/{id}` — detail + payment event history
- `POST /v1/invoices` — create (usually called by saga)
- `POST /v1/invoices/{id}/virtual-accounts` — issue a new VA (rotation / re-issue after expiry)
- `POST /v1/invoices/{id}/manual-payment` — record manual receipt (CS-facing)
- `POST /v1/invoices/{id}/void` — cancel an unpaid invoice (admin)
- `POST /v1/webhooks/midtrans` — public, signature-protected
- `POST /v1/webhooks/xendit` — public, signature-protected
- `POST /v1/refunds` — initiate refund (usually called by saga)
- `GET /v1/refunds/{id}` — status + settlement details
- `POST /v1/discount-approvals` — e-Approval workflow
- `GET /v1/reconciliation/reports` — finance admin view

**gRPC (service-to-service):**
- `IssueVirtualAccount`, `Refund`, `VoidInvoice`, `ReconcileInvoice` — called by broker-svc saga activities
- `GetInvoice`, `ListInvoicesByBooking` — read-only fan-out
- No gRPC from F5 pushing to F4 — all F5→F4 communication goes via **broker signals** (`SignalPaymentReceived`, `SignalRefundCompleted`, etc.) so the saga owns the order of operations

## Dependencies

- **F1** — identity, permissions, audit
- **F4** — booking records (invoice belongs to a booking)
- **External** — Midtrans and Xendit gateways, WhatsApp for receipt delivery, email as fallback, GCS for receipt PDFs, bank APIs for direct VA / settlement verification

Downstream consumers:
- F4 (status transitions on DP / lunas / expired / cancelled via broker signals)
- F9 (journal entries — AR, Unearned Revenue, tax)
- F8 (kit dispatch on lunas event)
- F10 (commission state transitions Pending → Confirmed on lunas)
- F6 (visa pipeline unlocks on lunas — passport-readiness gate)

## Backend notes

- **Adapter pattern.** Each gateway (Midtrans, Xendit, manual) is an adapter under `payment-svc/adapter/` with a common `GatewayAdapter` interface — `IssueVA`, `GetStatus`, `Refund`, `VerifyWebhookSignature`. Adding a new gateway is a new adapter; the service layer is gateway-agnostic.
- **Idempotency** at every call boundary: VA issuance uses `invoice.id` as idempotency key; webhook dedup uses `(gateway, gateway_txn_id)`; refund uses `refund.id`. This survives retries, gateway flakes, and duplicate saga runs.
- **Webhook handler is the critical path** — keep it tight: verify signature, dedupe, insert, signal, return 200. Anything heavier (receipt generation, external notifications) goes to a worker queue, not the webhook path itself.
- **Reconciliation cron** uses the same gateway adapter as the webhook — no duplicate logic. It just pulls status, feeds the same event-processing pipeline that the webhook uses.
- **FX snapshot** is computed from F9 `GetExchangeRate` at issuance time; finance owns the FX source of truth, F5 just snapshots the answer.
- **Dual-gateway selection** (Q013). Default implementation: primary/secondary config; try primary, fallback to secondary on transient error (5xx, timeout). On persistent error (4xx), return the error to the caller — don't mask a real failure.
- **PDF receipt rendering** runs in a separate goroutine pool with its own rate limiter, shared with F2 flyer rendering. Slow render doesn't block webhooks.

## Frontend notes

- Checkout page embeds the VA info (bank code, account number, amount, expiry countdown) with a copy-to-clipboard affordance.
- Customer Portal billing dashboard shows all invoices for the current jamaah, their payment history, and a "Pay remainder" button that routes back to issue a fresh VA for the outstanding balance.
- CS internal console has a dedicated "Generate Payment Link" quick action on every booking detail page.
- Finance admin reconciliation dashboard shows discrepancies (webhook says paid but F5 doesn't have it, or vice versa) as a queue to action.
- Refund preview UI shows the exact breakdown: "Diterima Rp X − Pinalti Rp Y − Biaya hangus Rp Z = Refund Rp W" before the customer (or CS on their behalf) confirms cancellation.

## Open questions

See `docs/07-open-questions/`:

**Existing (shared with other features):**
- **Q001** — operating currency, FX lock at invoice issuance
- **Q004** — cancellation → seat return timing (interlocks with refund flow)

**New, filed with this draft:**
- **Q011** — minimum DP %, max installment count, cadence
- **Q012** — refund penalty policy matrix (per package type, per timing bucket)
- **Q013** — dual-gateway (Midtrans / Xendit) selection and fallback rule
