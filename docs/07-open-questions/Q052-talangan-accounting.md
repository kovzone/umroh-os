---
id: Q052
title: Talangan accounting — loan receivable (PSAK 71) vs booking receivable
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: answered
---

# Q052 — Talangan accounting

## Context

Q018 (existing, open) captures Talangan (bridging finance) from the product-shape perspective: the agency lends jamaah a bridging amount (e.g. covers part of the package cost) that jamaah repays in installments. Q018 focuses on the booking-state and refund-interaction aspects. Q052 is F9's accounting-side complement: **how does Talangan sit on the balance sheet, and what accounting standard governs it?**

The options are materially different in PSAK terms:
- **Booking receivable** (Piutang Jamaah) — treat Talangan-arising piutang like any other jamaah AR; governed by general revenue / liability rules.
- **Loan receivable** (Piutang Pembiayaan) — treat as a separate loan asset; governed by **PSAK 71 (Financial Instruments)**, which requires Expected Credit Loss (ECL) modeling, effective-interest amortization, and different impairment treatment.

The classification matters because PSAK 71 compliance is substantially more work than ordinary AR.

## The question

1. **Classification**: Is Talangan a loan (PSAK 71 financial asset) or a bundled booking receivable?
2. **Interest / margin**: Does the agency charge interest or margin on Talangan? Flat fee, percentage, or zero?
3. **Default / impairment**: When jamaah fails to repay Talangan, what's the accounting? ECL model (PSAK 71) or simple writeoff?
4. **Collateral**: Is Talangan secured by anything (the Umrah booking itself as collateral)?
5. **Syariah compatibility**: If Talangan must be syariah-compliant, the structure is murabaha / akad qardh / ijarah — different accounting per PSAK Syariah standards (PSAK 102, 103, etc.).
6. **Recognition of Talangan income**: If interest/margin is charged, when is it recognized (effective-interest over life, or upfront)?
7. **On jamaah refund**: if jamaah cancels before full Talangan repayment, how is the outstanding loan balance handled?

## Options considered

- **Option A — Classify Talangan as loan receivable under PSAK 71, with ECL model + effective-interest amortization.** Full financial-instrument treatment.
  - Pros: regulation-correct if Talangan is genuinely a loan; appropriate audit posture.
  - Cons: substantial implementation (ECL model, effective-interest calc, impairment testing); probably overkill for small Talangan balances.
- **Option B — Classify Talangan as enhanced booking receivable.** Treat as Piutang Jamaah with an extended repayment schedule; no interest / margin recognition separately; simple writeoff on default.
  - Pros: simple; matches typical Indonesian agency practice where Talangan is an operational accommodation, not a financial product.
  - Cons: if Talangan carries interest/margin, this misrepresents the transaction's economic substance.
- **Option C — Configurable per Talangan structure — if interest-bearing, PSAK 71; if zero-interest goodwill accommodation, booking receivable.** Classification driven by presence/absence of margin.
  - Pros: handles both patterns; honest accounting.
  - Cons: two code paths; classification change mid-life is complex.

## Recommendation

**Option B — classify Talangan as enhanced booking receivable (no separate loan asset) with zero-margin default; revisit to Option A if the agency starts charging interest/margin at a non-trivial rate.**

Option A is the regulation-correct answer if the agency operates Talangan as a for-profit lending activity (syariah or conventional). But most Indonesian Umrah agencies use Talangan as a **soft operational accommodation** — a way to bridge a jamaah's short-term cash shortfall to close a booking — not as a profit center. The margin, if any, is typically a flat service fee (< 2% of principal) covering risk, not an interest rate. Treating this as PSAK 71 lending triggers ECL modeling and effective-interest accounting that dwarfs the materiality. Option B is the pragmatic default.

Defaults to propose: Talangan-arising piutang lives in `Piutang Jamaah (1.1.2.01)` like normal jamaah AR, with a `talangan_schedule` attribute capturing the repayment schedule. Margin (if any): flat fee recognized upfront as `Pendapatan Lain-Lain (4.9.x)` when Talangan is disbursed. No interest / effective-interest amortization. Default on repayment: after 90 days overdue, escalate to finance director; after 180 days, writeoff as `Beban Kerugian Piutang (6.x.x)`. Jamaah refund before full Talangan repayment: refund reduces outstanding Talangan balance first; any net-positive refund goes back to jamaah; any net-negative (Talangan > refund entitlement) becomes a debt-collection matter. Syariah compliance: MVP is conventional; syariah version (murabaha / qardh structure) is Phase 2 if agency requires. Upgrade path to PSAK 71 full treatment if interest/margin rises above materiality threshold (e.g. > 5% of principal or > 10M IDR per Talangan): at that point, reclassify Talangan to `Piutang Pembiayaan (1.1.X)` and apply PSAK 71.

Reversibility: upgrading to Option A later requires re-classifying existing Talangan balances — non-trivial but bounded in scope (small volume initially). Downgrading (Option A → Option B) is even rarer; not a common direction.

## Answer

**Decided:** Align with **Q018 Option A (MVP partner Talangan)** — **no agency loan receivable**; **referral fee income only**. If **in-house Talangan** is later enabled: **classify as PSAK 71 financial asset** with **ECL + effective interest** (Option A path). Until then, **do not** book **Piutang Pembiayaan** for Talangan.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
