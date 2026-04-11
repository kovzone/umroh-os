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

- **Inbound:** booking-svc requests an invoice/VA at booking submission; broker-svc activities call into payment-svc.
- **Outbound:** Midtrans/Xendit (issue VA, query status, refund), iam-svc (audit), booking-svc (mark paid).

## Notable behaviors

- **Idempotent webhook handling.** Gateway webhooks may retry; deduplicate via gateway transaction ID.
- **Reconciliation cron.** Periodic job reconciles VA status against the gateway in case a webhook was missed.
- **Multi-gateway support.** Adapter pattern — Midtrans and Xendit each have their own adapter under `adapter/`.
- **Refund flow** is a Temporal saga in broker-svc; payment-svc exposes the synchronous refund call.
