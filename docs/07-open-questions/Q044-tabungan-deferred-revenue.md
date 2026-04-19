---
id: Q044
title: Multi-year Tabungan deferred-revenue accounting
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9, F4, F5
status: answered
---

# Q044 — Multi-year Tabungan deferred-revenue accounting

## Context

Paket Tabungan is a savings-style product: jamaah enrolls and makes monthly installments over 12–60 months, accumulating funds toward a future Umrah booking. Q017 (existing) captures the booking-state interaction question. Q044 is the **finance-side complement**: how do Tabungan deposits sit on the balance sheet, how do they become revenue, and what happens when a Tabungan matures without a concrete departure booking?

PSAK 72 treats Tabungan as a **contract with a customer** where the performance obligation is future Umrah service. Until the obligation is satisfied (jamaah flies per Q043), deposited funds are a **contract liability**. But the nuance: a Tabungan not-yet-converted-to-a-booking is a contract with a vague performance date ("someday"); a converted Tabungan becomes a concrete booking with a scheduled departure.

## The question

1. **Tabungan deposit received** — what's the journal entry?
   - Dr Bank / Cr Hutang Tabungan Jamaah (deferred revenue / contract liability)?
   - Or Dr Bank / Cr Pendapatan Diterima Dimuka?
   - Account classification: current vs long-term liability based on expected departure date?
2. **Tabungan conversion to booking** — does the Hutang Tabungan move to Hutang Jamaah (normal booking liability), or does the Tabungan liability stay and the booking is created on top?
3. **Stale Tabungan** — jamaah stops paying after 6 months and never books. After how long does the accumulated balance become something else (forfeited revenue? refundable on demand? escheat)?
4. **Tabungan maturity without departure** — matured Tabungan balance exceeds the booked departure cost; is the excess refunded, applied to surcharges, or rolled to next season?
5. **Interest / profit-share** — does Tabungan earn anything? If yes, is it a simple growth accrual or an investment-linked return? Either has different PSAK treatment.
6. **Classification as current vs long-term liability** — if Tabungan is, say, a 5-year plan, is year-1's expected-departure balance current liability and years 2–5 long-term?

## Options considered

- **Option A — All Tabungan balances as current liability; no long-term split.** Simple, single account `Hutang Tabungan Jamaah`. Aging report shows how long each balance has sat.
  - Pros: simple; doesn't require expected-departure-date forecast.
  - Cons: doesn't reflect true balance-sheet structure; misleads liquidity analysis.
- **Option B — Split current vs long-term based on expected departure date.** Tabungan balance with expected departure ≤ 12 months = current; > 12 months = long-term. Requires tracking `expected_departure_month` on each Tabungan enrollment.
  - Pros: PSAK-1 compliant presentation; more informative balance sheet.
  - Cons: "expected departure" is a soft commitment — jamaah can push dates; reclassification cadence needed.
- **Option C — Option B + explicit stale-Tabungan escalation workflow.** Balances untouched > 24 months go to a "stale" bucket; CS contacts jamaah; if no response after 90 days, balance flagged for review (refund, rollover, forfeit per company policy).
  - Pros: handles the real-world dormancy problem; clean audit trail.
  - Cons: requires company policy (Q044 sub-question) on what happens to unclaimed balances.

## Recommendation

**Option C — split current vs long-term by expected departure + explicit stale-balance workflow.**

Option A leaves a potentially large pile of deposits undifferentiated on the balance sheet — confusing for any reviewer and non-compliant with PSAK 1 current/non-current split. Option B gets the presentation right but doesn't address the operational reality that 5-15% of Tabungan accounts go dormant (jamaah changes job, moves, passes away). Option C adds the dormancy workflow — which is a policy question we're surfacing but can default.

Defaults to propose: Tabungan deposit Dr Bank / Cr Hutang Tabungan Jamaah. Each Tabungan enrollment carries `target_departure_season` (soft commitment); balance reclassifies current/long-term monthly based on whether target_departure is ≤ 12 months. Conversion to booking: Hutang Tabungan → Hutang Jamaah via transfer journal (no revenue recognition yet — Q043 still applies). No interest / profit-share in MVP (keeps out of PSAK 71 investment-contract complexity). Stale escalation: balance untouched > 24 months enters `stale` status, CS contacts, 90-day response window, unresolved balances escalated to finance-director for case-by-case decision (refund, rollover, forfeit — company policy decision).

Reversibility: splitting the single Hutang Tabungan account into current/non-current sub-accounts later requires reclassifying existing balances — a migration, not trivial but doable. Adding interest / profit-share later is a new feature (PSAK 71 contracts) — additive.

## Answer

**Decided:** **Option C** — **Dr Bank / Cr Hutang Tabungan**; **monthly current vs non-current split** using `target_departure_season`; **conversion journal** to **Hutang Jamaah**; **no profit-share MVP**; **dormancy:** >24mo inactive → `stale` + CS workflow + **90d** then **finance director case file** (refund/rollover/forfeit **policy text still needs legal**).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
