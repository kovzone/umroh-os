---
id: Q048
title: FX policy — rate source, transaction-date vs settlement-date, month-end revaluation
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: answered
---

# Q048 — FX policy

## Context

PRD Alur Logika 1.1 (lines 1273–1325) establishes FX fundamentals:
- Supported currencies: IDR (base), USD, SAR.
- Rate source configurable: BI API + markup OR manual lock rate per season.
- Rate non-retroactivity: changes do not restate already-DP'd/Lunas invoices.

But silences remain: **for every single journal line in a non-IDR currency, which rate is snapshotted?** And what about the PSAK 10 requirement to revalue open foreign-currency balances monthly?

Q001 (existing, open) covers the operating-currency question at the company level. Q048 is the F9-specific mechanics layer: how do multi-currency transactions post, and how does FX gain/loss materialize.

## The question

1. **Rate source and cadence.** BI middle rate, tax rate (Kurs Menteri Keuangan), company-locked rate, or commercial bank rate? Fetched daily, weekly?
2. **Transaction-date vs settlement-date.** When vendor invoice is received in USD on date X and paid on date Y, which rate is used to snapshot the AP journal at date X vs at date Y? Difference hits FX gain/loss.
3. **Month-end revaluation (PSAK 10).** Open foreign-currency AR and AP balances revalued to closing rate monthly? If yes, the revaluation journal posts gain/loss to P&L.
4. **Locked rate per season for HPP** (PRD line 1281) — how does this interact with transaction-date rates? Is the locked rate used for HPP budgets only, or also for actual cost bookings?
5. **FX gain/loss account structure** — single `FX Gain/Loss` (4.9.x) account or separate Gain + Loss accounts?
6. **Multi-leg FX** — if a SAR vendor invoice is paid via a USD-denominated account (bank account in USD), two FX events happen: SAR→IDR at invoice and SAR→USD at settlement.

## Options considered

- **Option A — BI middle rate daily; transaction-date for AR/AP snapshot; month-end revaluation; single FX Gain/Loss account.** Standard PSAK 10 treatment.
  - Pros: textbook PSAK-compliant; clean audit trail.
  - Cons: BI API dependency; rate differences vs commercial rates actually used for remittance.
- **Option B — KMK (Kurs Menteri Keuangan) weekly; transaction-date snapshot; month-end revaluation.** Uses the tax-department official rate.
  - Pros: aligns with tax reporting; simpler cadence (weekly).
  - Cons: KMK lags commercial rates; FX gain/loss accrues vs commercial realities.
- **Option C — Company-locked rate per season for budget; BI daily rate for actual transactions; month-end revaluation.** Combines PRD's lock-rate-for-HPP with daily-rate for real flows.
  - Pros: matches PRD's season-lock intent for budgeting; real FX gain/loss against actual rates.
  - Cons: two-rate system is more complex to explain; auditor may question consistency.

## Recommendation

**Option C — company-locked rate per season for HPP/budget purposes; BI middle rate (daily API fetch) for all actual transaction snapshots; month-end revaluation per PSAK 10; single combined FX Gain/Loss account split into sub-accounts if reporting needs require.**

Option A is textbook but misses the PRD's clear intent to lock rates for HPP-and-pricing purposes (line 1281). Option B's KMK is tax-reporting-convenient but doesn't reflect commercial reality for remittances. Option C threads the needle: the locked rate is a **budget / pricing rate** (never actually posts to GL), and the BI daily rate is the **actual transaction rate** (what posts to GL). Month-end revaluation handles the gap per PSAK 10. This mirrors how Indonesian travel agencies actually operate: they lock a rate for brochure pricing so customer quotes stay stable, but actual vendor payments happen at whatever the rate is that day.

Defaults to propose: FX rate source = BI middle rate (fetched daily at 9am WIB); fallback = previous-day rate on API failure; configurable per agency. Transaction-date snapshot = every journal_line.fx_rate captures the date's rate at posting time. Month-end cron (25th of month, recalibrates): revalues open AR/AP foreign-currency balances to month-end rate; posts FX revaluation journal (Dr FX Loss / Cr AR for unfavorable moves on AR, etc.). Separate FX Gain and FX Loss accounts per PSAK 1 minimum-line-item presentation (combined or split per Q042's Option A/B decision). Multi-leg FX scenarios: two separate journal entries, one per leg. HPP lock rate = catalog-svc attribute (`season_locked_rate`); used only for package pricing and HPP budget computation; never posts to GL directly.

Reversibility: changing rate source (e.g. BI → KMK) is a config change + historical rate backfill. Switching month-end revaluation on/off is additive.

## Answer

**Decided:** **Option C** — **season lock for pricing/HPP only (non-GL)**; **BI middle daily ~09:00 WIB** for **transaction GL rates** with **previous day fallback**; **transaction-date** initial measurement; **payment-date** differences → **FX P&L**; **month-end PSAK 10 reval** on open FC monetary items; **separate `FX Gain` and `FX Loss` accounts**; **multi-leg** = separate entries per leg. **Customer invoice snapshots** follow **Q001** lock at booking.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
