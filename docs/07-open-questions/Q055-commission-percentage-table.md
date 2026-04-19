---
id: Q055
title: Commission % table (per level × per product)
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10, F9
status: answered
---

# Q055 — Commission % table (per level × per product)

## Context

Commission is the economic engine of the agent program. PRD describes accrual mechanics (Pending / Confirmed / Withdrawal) and override formula but **names no percentage values**. Every F9 journal entry (accrual, payout, tax withholding) needs a concrete number.

Commission varies by two axes:

1. **Agent tier** — Silver / Gold / Platinum / Cabang (per Q054). Higher tiers earn higher %.
2. **Product type** — Umroh Reguler, Umroh Plus, Haji Furoda, Badal, retail add-ons (insurance, merchandise), Tabungan. Different product margins support different commission rates; Haji Furoda and Badal might carry different structures entirely.

## The question

1. **Base commission % per tier × per product** — fill the matrix.
2. **Is commission % on gross package price or net (after tax)?**
3. **Promotional / override campaigns** — can tier's base % be temporarily uplifted for a campaign (e.g. double-commission month)?
4. **Non-commissionable items** — which add-ons, fees, or package components don't earn commission (tax, visa fee, deposit, specific surcharges)?
5. **B2B agent closing a B2B sub-agent's jamaah** — does the super-agent earn direct or override?
6. **Commission on Tabungan** — agent recruits a Tabungan enrollment; do they earn at enrollment time, or only when Tabungan converts to a booking (Q044 cross-ref)?
7. **Minimum commission floor** — is there a minimum IDR amount per closing (e.g. 100K IDR floor), or strict percentage?

## Options considered

- **Option A — Flat per-tier % across all products.** Silver 3%, Gold 5%, Platinum 8%. Simple, single table.
  - Pros: easy to understand; easy to implement.
  - Cons: doesn't account for product-margin variance; Haji Furoda's higher margin justifies higher %.
- **Option B — Tier × product matrix.** Full matrix (3 tiers × 4–6 products) with per-cell %.
  - Pros: reflects product-margin reality; flexible.
  - Cons: 12–18 numbers to finalize; maintenance overhead.
- **Option C — Tier × product-category (not full product list).** Categorize products into ~3 commission-categories (standard / premium / specialty) and set tier × category matrix. Individual products map to a category.
  - Pros: middle ground; manageable.
  - Cons: categorization itself is a decision; add-on products may straddle.

## Recommendation

**Option C — tier × product-category matrix (3 tiers × 3 categories = 9 cells); commission on net package price; non-commissionable items explicitly listed.**

Option A's flat-rate hides the margin reality — a Silver agent earning 3% on a high-margin Haji Furoda closes the same as a high-margin Paket Silver, which under-motivates them on the higher-lift sale. Option B's full matrix is correct but creates maintenance friction (every new product needs a commission number). Option C categorizes products once (standard = Umroh Reguler / Plus; premium = Haji Furoda / Badal; specialty = Tabungan / retail add-ons) and lets the matrix handle the rest.

Defaults to propose (values are placeholders pending stakeholder calibration):

| | Standard | Premium | Specialty |
|---|---|---|---|
| Silver | 3% | 4% | 1% |
| Gold | 5% | 7% | 2% |
| Platinum | 7% | 10% | 3% |

Cabang / Perwakilan override % (Q056) sits on top.

Commission base: **net package price** (gross − PPN − non-commissionable add-ons). Non-commissionable explicitly listed: tax amounts (PPN / PPh), government fees (visa fee, tasreh fee), insurance premium, hard costs pass-through (e.g. airline fuel surcharge). Promotional uplift: campaign can add a flat X% to all tiers for duration of campaign (e.g. +1% for Ramadhan Flash Sale); uplift tracked separately in ledger. B2B agent closing a sub-agent's jamaah: the closing agent (the one at the bottom of the hierarchy receiving the jamaah) gets direct; upline gets override per Q056. Commission on Tabungan: zero at enrollment (Specialty 1-3% above is at Tabungan → booking conversion, not enrollment — avoids paying commission on funds that may never book). Minimum commission floor: 50K IDR per closing (prevents humiliating sub-1K commissions on cheap retail items).

Reversibility: commission table is config in a single row of `commission_policy` table — editable anytime by finance-director. Historical bookings retain the rate at their `paid_in_full` timestamp (snapshot per entry).

## Answer

**Decided:** **Option C** matrix (9 cells) as **starting defaults** in Recommendation table; **commissionable base = net of taxes + listed non-commissionable pass-throughs**; **Tabungan: no commission until booking conversion**; **floor Rp50k per closing**; **Cabang override rate** separate config (not in tier table). **Snapshot rates** at accrual event (**Q045**).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
