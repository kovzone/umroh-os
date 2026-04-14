---
id: F5
title: Payment & Reconciliation
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 4+ Must Haves (VA issuance, bank integration, digital receipts, AR subledger)
prd_sections:
  - "A. Self-Service Booking — Payment Gateway B2C"
  - "G. Finance — Revenue & Receivables (payment slice)"
modules:
  - "#18, #129, #130, #131, #132, #133"
depends_on: [F1, F4]
---

# F5 — Payment & Reconciliation

## Purpose & personas

TBD — virtual accounts, gateway webhooks, settlement reconciliation, refunds. Only service that talks to Midtrans/Xendit. Emits events that finance-svc and booking-svc subscribe to.

## Sources

- PRD Section G — payment portion; Section A — checkout
- Modules #18, #129–133

## User workflows

TBD:
- W1: Booking submit triggers VA issuance
- W2: Jamaah pays to VA; webhook fires; booking marked paid
- W3: Partial payment (DP) flow
- W4: Refund initiation
- W5: Reconciliation cron recovers missed webhooks

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD. Critical: webhook replay, signature verification failure, settlement lag, partial refund, FX on a foreign-currency invoice.

## Data & state implications

See `docs/03-services/04-payment-svc/02-data-model.md`.

## API surface (high-level)

See `docs/03-services/04-payment-svc/01-api.md`.

## Dependencies

- F1 (IAM)
- F4 (booking — invoice belongs to a booking)

## Backend notes

TBD. Idempotency on `gateway_txn_id`. Webhook signature verification mandatory. Refund flow is a saga in `broker-svc`.

## Frontend notes

TBD.

## Open questions

None yet.
