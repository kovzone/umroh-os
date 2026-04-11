# finance-svc — Overview

## Purpose

PSAK-compliant double-entry accounting. Consumes events from payment, logistics, and crm to write journal entries. Owns the chart of accounts, AR/AP, tax records, FX, and job-order costing.

## Bounded context

Finance & Accounting. See `docs/02-domain/00-bounded-contexts.md` § 9.

## PRD source

PRD section G — Finance & Accounting.

## Owns (data)

- `chart_of_accounts` (COA tree)
- `journal_entries` and `journal_lines`
- `ar_balances`
- `ap_balances`
- `tax_records`
- `fx_rates`
- `job_order_costs` (cost center per departure)
- `vendors` (financial side)

## Boundaries (does NOT own)

- Payment events (`payment-svc`) — finance consumes them
- Logistics POs (`logistics-svc`) — finance consumes GRN events
- Commission calculations (`crm-svc`) — finance journals the result

## Interactions

- **Inbound:** event subscriptions from payment, logistics, crm (initially: gRPC pulls + scheduled jobs).
- **Outbound:** none in synchronous path. Reports are generated on-demand.

## Notable behaviors

- **Double-entry enforcement** — every journal entry must balance (sum of debits = sum of credits).
- **Multi-currency** with FX gain/loss tracking — `fx_rates` snapshotted daily.
- **PSAK-compliant** chart of accounts and reporting (balance sheet, P&L, cash flow).
- **Job-order costing** — each `package_departure` is its own cost center for profitability.
- **Tax calculation** — PPh 21/23 and PPN per Indonesian rules.
