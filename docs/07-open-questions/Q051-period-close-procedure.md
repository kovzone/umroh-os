---
id: Q051
title: Period close procedure + re-open authority
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: answered
---

# Q051 — Period close procedure + re-open authority

## Context

PRD doesn't describe period close. Module #147 (real-time financial reports) implies reports can be generated any time, which is correct — but doesn't describe the **periodic lock** that PSAK-compliant accounting requires. Without period close:
- Late entries can retroactively restate prior periods (audit risk).
- Month-over-month comparisons are unstable.
- Period-based reports (SPT PPN, monthly financial statements) can't be treated as final.

PSAK 1 requires comparative periods. Period close is how you make each period comparable. But close-too-hard blocks legitimate corrections; close-too-soft defeats the purpose.

## The question

1. **Close cadence** — monthly (soft close), annual (hard close)? Quarterly?
2. **Soft close vs hard close** — soft = warn but allow adjustment with approval; hard = immutable, only restatement via new period.
3. **Pre-close checklist** — what must run before close (FX revaluation, depreciation, accruals, reconciliation)?
4. **Re-open authority** — if a soft-closed period needs a correction, who can re-open? Finance director? CFO? CEO?
5. **Audit trail on re-opens** — what's captured about why the period was re-opened?
6. **Year-end hard close + retained earnings roll** — what happens on annual close (revenue/expense accounts zero out to retained earnings)?
7. **Branch-level close vs entity-wide close** — if branches report separately, close per branch or only at entity level?

## Options considered

- **Option A — Monthly soft close + annual hard close; finance-director re-open with CFO/CEO approval; strict audit trail.** Standard for Indonesian SMEs + mid-size enterprises.
  - Pros: balances speed and rigor; matches most agency practice.
  - Cons: soft close can be abused if approval is rubber-stamped.
- **Option B — Monthly hard close; correcting entries only via reversal in subsequent period.** No re-open path; errors corrected forward.
  - Pros: rigorous; no retroactive restatement.
  - Cons: when you close month 3 on April 5th and then learn a March invoice was mis-categorized, the correction lives in April — auditors don't love this.
- **Option C — Monthly soft close, mid-year (Q2, Q4) progressive locks, annual hard close.** Graduated locking.
  - Pros: middle ground; allows corrections early but tightens as year progresses.
  - Cons: more state to track; confusing for users.

## Recommendation

**Option A — monthly soft close + annual hard close; soft-close-re-open requires finance-director + CFO co-sign; hard-close re-open is effectively impossible (requires audit committee + restatement).**

Option B is too rigid for a real agency — legitimate corrections to prior periods happen (late vendor invoices, bank reconciliation adjustments, booking-status amendments). Option C adds complexity without clear return. Option A is the industry-standard pattern that balances pragmatism and rigor: monthly soft close = each month is "normally final" but can be re-opened with proper approval; year-end hard close = audit-ready lock that shouldn't move.

Defaults to propose: monthly soft close runs on 10th of following month (covers late vendor invoices). Pre-close checklist: (1) all unposted events reconciled, (2) FX revaluation run, (3) monthly depreciation run, (4) bank reconciliation complete, (5) AR/AP aging reviewed, (6) all manual journals approved. Soft-close status: `periods.status = 'soft_closed'`. Re-open requires finance director + CFO co-sign; reason required in audit log; entry of new journal posts allowed in re-opened period until re-closed. Annual hard close runs 60 days after year-end (covers annual audit). Hard-close re-open: effectively requires restatement process (not a single-button operation). Closing journals on annual: revenue / expense accounts zero to retained earnings. Branch-level close: not in MVP; entity-wide close only (branches see their slice of the entity-level close). Audit log on close: period, closer, pre-close checklist status per item, unposted-events count, total journal count.

Reversibility: close cadence (monthly vs quarterly) is a config. Soft vs hard rules are config. Making hard-close re-openable is additive (but discouraged by audit standards).

## Answer

**Decided:** **Option A** — **monthly soft close T+10**, checklist (FX, depreciation, bank recon, AR/AP review); **reopen requires finance director + CFO** reason; **annual hard close T+60** with **closing entries** to retained earnings; **entity-wide only** MVP.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
