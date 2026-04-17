---
id: Q053
title: Refund & pinalti accounting entries
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9, F5, F4
status: open
---

# Q053 — Refund & pinalti accounting entries

## Context

Q012 (existing, open) covers the refund penalty policy matrix — how much penalty is charged per cancellation window. Q053 is the F9-specific downstream: **how do refund + pinalti journal entries look?** Multiple sub-questions arise from the cancellation mechanics:

- Pre-departure cancellation — no revenue recognized yet; refund is a simple Hutang Jamaah reduction + cash-out, but **pinalti retained** — where does it post?
- Post-departure refund (goodwill / legal) — revenue already recognized; reversal journal?
- Sunk costs already paid to vendors (visa fee, ticket deposit) — these don't come back on refund; accounting treatment?
- Pinalti retention as revenue — is it `Pendapatan Umroh` (main revenue), `Pendapatan Lain-Lain` (non-operating other income), or `Pendapatan Pinalti` (distinct account)?

These mechanics are invisible from PRD Section G; F9 must define them.

## The question

1. **Pre-departure cancellation + pinalti:** the refund journal structure?
   - Dr Hutang Jamaah — full amount paid
   - Cr Bank — refund amount (paid − pinalti)
   - Cr [???] — pinalti amount
2. **Where does pinalti post** — `Pendapatan Lain-Lain`, `Pendapatan Umroh`, `Pendapatan Pinalti`?
3. **Post-departure refund (goodwill):** is the revenue reversed?
   - Dr Pendapatan Umroh — refund amount
   - Cr Hutang Jamaah — refund amount
4. **Sunk-cost handling on pre-departure cancellation:** visa fees already paid to government (non-refundable). Agency ate the cost, but jamaah's pinalti may or may not cover it. Does the sunk cost stay as expense?
5. **Partial refund** (Q014 — one jamaah cancels from a multi-jamaah booking): per-jamaah journal entries.
6. **Refund timing:** does the refund journal fire when refund is committed (Hutang Jamaah flipped to Refund Liability), or when cash actually leaves the bank (paid)?

## Options considered

- **Option A — Pinalti to `Pendapatan Lain-Lain (4.9.x)` non-operating revenue; sunk costs stay as expense; post-departure refund reverses revenue.** Standard treatment for cancellation penalties in service industries.
  - Pros: correct matching; pinalti not muddling Umrah revenue metrics.
  - Cons: agencies may want to see pinalti alongside main revenue for P&L visibility.
- **Option B — Pinalti to main `Pendapatan Umroh`; sunk costs stay as expense; post-departure refund reverses revenue.**
  - Pros: single revenue line, simple presentation.
  - Cons: muddles actual Umrah-delivery revenue with cancellation income.
- **Option C — Distinct `Pendapatan Pinalti (4.8.x)` account; sunk costs stay as expense; post-departure refund reverses revenue and transfers recovered amount to Pendapatan Pinalti.**
  - Pros: finest granularity; reports can show pinalti impact explicitly.
  - Cons: one more account; extra complexity in the journal engine.

## Recommendation

**Option A — pinalti to `Pendapatan Lain-Lain (4.9.x)` non-operating revenue; sunk costs stay as expense (already posted at payment time); post-departure refund reverses revenue.**

Option B muddles Umrah delivery revenue with cancellation fees — a 200M IDR pinalti month looks like a 200M IDR main-revenue month, which misleads any KPI built on Pendapatan Umroh. Option C's granularity is nice but extra complexity for marginal presentation benefit; Pendapatan Lain-Lain is already the natural home for this kind of non-recurring income per PSAK 1. Option A is the cleanest.

Defaults to propose:

**Pre-departure cancellation + refund:**
```
Dr Hutang Jamaah                 — amount_paid_to_date
Cr Bank                          — refund_amount (= paid − pinalti)
Cr Pendapatan Lain-Lain (4.9.x)  — pinalti_amount
```

Sunk costs: already posted at payment-to-vendor time (e.g. `Dr Beban Visa / Cr Bank` when visa fee paid; Dr Hutang Usaha / Cr Bank when AP settled). They stay as expense — not related to the refund journal.

**Post-departure refund (rare, goodwill/legal):**
```
Dr Pendapatan Umroh              — refund_amount (revenue reversal, tagged with job_order)
Cr Hutang Jamaah                 — refund_amount
```

Then the refund-out:
```
Dr Hutang Jamaah                 — refund_amount
Cr Bank                          — refund_amount
```

**Partial refund (one jamaah of many):** per-jamaah share of paid amount; same pattern applied to the cancelling jamaah's share only; other jamaah's Hutang Jamaah balances unchanged.

**Refund timing:** two-step. When refund is **committed** (approved per ops + finance approval), flip `Hutang Jamaah` to a distinct `Refund Liability (2.1.x)` account (optional — MVP may skip this intermediate step and post directly on cash-out). When refund cash actually leaves, `Dr Hutang Jamaah / Cr Bank` (or `Dr Refund Liability / Cr Bank` if two-step).

Reversibility: reclassifying pinalti from `Pendapatan Lain-Lain` to a distinct account later is a COA-restatement — doable but non-trivial.

## Answer

TBD — awaiting stakeholder input. Deciders: finance director, external auditor, agency owner (KPI preference on pinalti visibility).
