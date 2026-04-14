---
id: F9
title: Finance & Accounting (PSAK-compliant)
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 15 Must Have / 7 Should Have — highest Must Have count in the catalogue
prd_sections:
  - "G. Finance & Accounting"
modules:
  - "#129–150"
depends_on: [F1, F5, F8, F10]
---

# F9 — Finance & Accounting (PSAK-compliant)

## Purpose & personas

TBD — PSAK-compliant double-entry accounting. Consumes events from payment, logistics, and crm to write journal entries. Reports balance sheet, P&L, cash flow. Job-order costing per departure.

Primary personas: finance admin (journaling, approvals), CFO / owner (reports), auditor (audit trail).

## Sources

- PRD Section G in full
- Modules #129–150 (22 total, 15 of them Must Have)

## User workflows

TBD:
- W1: Event consumer records AR from payment received
- W2: Event consumer records AP from GRN
- W3: Manual journal entry (adjustments)
- W4: Monthly close with revenue recognition
- W5: Commission payout to agents
- W6: Balance sheet / P&L report generation
- W7: Tax filing support (PPh 21/23, PPN)

## Acceptance criteria

TBD. Non-negotiable: double-entry integrity (sum debits = sum credits per entry, enforced at DB).

## Edge cases & error paths

TBD. Critical: FX gain/loss on foreign-currency transactions; revenue recognition timing for unearned income; period close after correcting entries.

## Data & state implications

See `docs/03-services/08-finance-svc/02-data-model.md`.

## API surface (high-level)

See `docs/03-services/08-finance-svc/01-api.md`.

## Dependencies

- F1 (IAM), F5 (payment events), F8 (logistics events), F10 (commission events)

## Backend notes

TBD. Most complex domain in the system. PSAK compliance means specific account codes, specific categorization rules, specific reporting formats. Heavy involvement from finance stakeholders expected — many questions likely go into `07-open-questions/`.

## Frontend notes

TBD.

## Open questions

None yet. Expected: COA template selection, revenue recognition policy per package type, default tax codes, job-order close rules.
