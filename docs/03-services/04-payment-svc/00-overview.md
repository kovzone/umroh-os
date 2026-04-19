# payment-svc — Overview

## Purpose

Invoices, virtual accounts, gateway webhook ingestion, settlements, and refunds. The only service that talks to Midtrans/Xendit.

## Bounded context

Payment. See `docs/02-domain/00-bounded-contexts.md` § 5.

## PRD source

PRD section G (Finance & Accounting — payment portion) and parts of A (B2C checkout).

## Owns (data)

- `invoices` — bills against bookings
- `virtual_accounts` — VA records issued by gateways
- `payment_events` — gateway webhook records
- `refunds`

## Boundaries (does NOT own)

- Bookings (`booking-svc`) — payment-svc signals booking on settlement
- Journal entries (`finance-svc`) — finance consumes payment events for journaling
- Tax calculation (`finance-svc`)

## Interactions

- **Inbound:** booking-svc requests an invoice/VA at booking submission; booking-svc calls into payment-svc synchronously as part of the in-process submit saga (per ADR 0006).
- **Outbound:** Midtrans/Xendit (issue VA, query status, refund), iam-svc (audit), booking-svc (mark paid).

## Notable behaviors

- **Idempotent webhook handling.** Gateway webhooks may retry; deduplicate via gateway transaction ID.
- **Reconciliation cron.** Periodic job reconciles VA status against the gateway in case a webhook was missed.
- **Multi-gateway support.** Adapter pattern — Midtrans and Xendit each have their own adapter under `adapter/` (**Q013** failover rules in F5).
- **Refund flow** is coordinated in-process inside payment-svc (per ADR 0006; no Temporal in MVP); payment-svc exposes the synchronous refund call and drives compensations on each downstream step.
- **FX lock + IDR settlement (Q001 / F5).** Each `invoice` stores `fx_snapshot` at issuance; `amount_total` is **IDR** after one-shot **Rp 1,000** half-up rounding; `rounding_adjustment_idr` feeds the finance rounding GL.
