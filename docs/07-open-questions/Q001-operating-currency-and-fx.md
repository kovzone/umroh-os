---
id: Q001
title: Operating currency, FX handling modes, and HPP formula
asked_by: session 2026-04-14 F2 draft
asked_date: 2026-04-14
blocks: F2, F5, F9
status: open
---

# Q001 — Operating currency, FX handling modes, and HPP formula

## Context

F2 (catalog) needs to decide how package prices are entered, stored, and displayed. PRD Section D treats IDR as the default sell currency but Alur Logika 1.1 introduces global FX variables (IDR ↔ USD ↔ SAR) and two handling modes:

- **Lock Rate** — admin sets a seasonal rate manually (e.g. 1 SAR = Rp 4,100) and all HPP calculations use that locked rate until rotated.
- **API + Markup** — system pulls daily Bank Indonesia kurs and adds a configurable markup (e.g. BI + Rp 150).

F5 (payment) needs to know which currency to bill in and how FX snapshots attach to invoices. F9 (finance) needs to know how to journal multi-currency transactions and compute unrealised P/L.

## The question

Three linked sub-questions. Please answer each:

1. **Sell currency.** Is IDR the **only** customer-facing currency, with all cost components (SAR hotel rates, USD airfare) converted into IDR via the operating FX rate at quote time? Or should the system support multi-currency sell prices (e.g. a Haji Furoda package priced in USD)?
2. **Default FX mode.** Which mode is the platform default — Lock Rate or API + Markup? Can individual packages override the default?
3. **HPP (cost of goods sold) formula.** Is HPP computed as the sum of foreign-currency cost components multiplied by the effective rate at HPP-lock time, plus IDR cost components? Or is there a different formula in use today?

## Options considered

- **Option A — IDR-only sell, Lock Rate default.** Simplest. All prices displayed in IDR. Admin rotates the lock rate per season. Cost components (SAR/USD) stored with their native amounts plus the locked FX rate. Matches what most Indonesian Umrah agencies do.
  - Pros: simplest admin mental model, no FX surprises for customers, clean invoice accounting.
  - Cons: requires admin discipline to rotate the lock rate; big market moves can make already-signed packages unprofitable.
- **Option B — IDR-only sell, API + Markup default.** Daily kurs auto-updates. Protects against market shifts but new bookings get slightly different prices day-to-day.
  - Pros: always current, no manual rotation.
  - Cons: customers/agents see daily price jitter; more complex to explain to CS.
- **Option C — Multi-currency sell.** Some packages (esp. Haji Furoda) priced directly in USD. Customer pays in IDR at point-of-sale using current rate.
  - Pros: matches how some premium programs are actually sold.
  - Cons: doubles the complexity in catalog + invoice + journal; PSAK journaling around multi-currency sales is intricate.

## Recommendation

**Option A — IDR-only sell with Lock Rate default.** The Indonesian Umrah market expects IDR prices, and Lock Rate matches how agencies already think about seasonal pricing (they negotiate hotel/flight allotments in SAR/USD at the start of a season, then fix a customer-facing IDR rate). API+Markup sounds tidy but creates daily price jitter that CS will have to explain away on every call — the admin discipline cost of rotating a lock rate is cheaper than the support cost of jittery prices. Option C (multi-currency sell) is the right answer for Haji Furoda eventually, but paying its complexity upfront before we know demand is premature.

Keep `currency` as a per-package column from day one so we can add Option C later without migration. Expose API+Markup as an alternative mode behind a global config flag — it's a one-line change to implement, and having it unused is cheaper than retrofitting.

HPP formula (also recommended): `HPP = Σ(cost_component_idr) + Σ(cost_component_sar × locked_sar_rate) + Σ(cost_component_usd × locked_usd_rate)`. Sell price stays an admin-set field, **not** auto-computed from HPP + markup, because margin targets vary per package kind and commercial negotiations frequently override the formula.

Reversibility: adding Option C later is schema-compatible (the `currency` column already exists). Switching from Lock Rate to API+Markup is a config flip, no migration. This recommendation is deliberately the most reversible of the three.

## Answer

TBD — awaiting stakeholder discussion.
