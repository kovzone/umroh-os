---
id: Q040
title: Stock availability policy — partial shipments + reorder-point math
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8
status: open
---

# Q040 — Stock availability policy

## Context

Two related product decisions that the PRD under-specifies:

**Partial receipts / partial shipments.** PRD module #113 (line 385) says stock increments when received qty matches PO qty — silent on short-shipped receipts. PRD module #119 (line 401) says kit assembly is atomic (all components or none) — silent on what happens when one SKU is out and the rest are in stock.

**Reorder-point math.** PRD module #118 (line 399) says alert fires when an item hits "batas minimum" — silent on how that threshold is set. Static-per-SKU is the naive answer; demand-forecasting (e.g. upcoming 30-day departures × kit components) is smarter but requires a feed from booking-svc.

Together these shape whether F8 is a "just-enough" stock system or a planning-aware one.

## The question

**Partial operations:**
1. Does GRN accept partial receipts (some line qty received, some not)? If yes, what's the PO's status — `partially_received` with open lines, or closed-short forcing a new PO?
2. Does kit assembly allow partial-shipment (kit dispatched without one component, backorder later)?
3. Does dispatch allow partial per booking (some jamaah in a multi-jamaah booking get their kit, others wait)?

**Reorder-point math:**
4. Is the reorder threshold static per-SKU (set once, edited by supervisor), or computed from demand?
5. If computed: what window (30 days? 90 days? to next peak season?), what data (upcoming departures × kit components × lead time to vendor), and where does it live (logistics-svc compute, or an analytics service)?
6. What happens when auto-reorder triggers — creates a draft PR for review, or auto-submits?

## Options considered

**Partial operations:**

- **Option A (partial) — Accept partial receipts; block partial kit assembly and partial dispatch.** GRN can close partial, leaving open lines for later receipt; kits require all-or-none at assembly time; dispatch is per-booking (all jamaah in the booking or none).
  - Pros: preserves kit integrity (jamaah don't receive incomplete kits); matches PRD line 401.
  - Cons: one missing SKU can stall many dispatches.
- **Option B (partial) — Accept partial everywhere with explicit tracking.** Partial GRN (open remaining lines), partial kit assembly (backorder the missing component as a separate shipment), partial booking dispatch (each jamaah independent).
  - Pros: throughput; nothing blocks on one SKU.
  - Cons: jamaah receives partial kit (bad UX); ops workload tracking orphan backorders.
- **Option C (partial) — Partial GRN only; kit and dispatch stay all-or-none.** Matches PRD literal reading.
  - Pros: compromise; preserves customer-facing completeness.
  - Cons: same throughput risk as A for kit assembly.

**Reorder-point math:**

- **Option X (reorder) — Static per-SKU threshold, editable by supervisor.** No demand forecast.
  - Pros: simplest; supervisor controls; fast to ship.
  - Cons: threshold drifts from reality; manual tuning required.
- **Option Y (reorder) — Demand-forecast-driven, auto-computed.** logistics-svc queries booking-svc for upcoming 30-day / 90-day pax counts, multiplies by kit components + safety stock + vendor lead time.
  - Pros: adaptive; matches real demand.
  - Cons: more complex; depends on timely booking data; harder to debug threshold decisions.
- **Option Z (reorder) — Static default with optional auto-compute override per SKU.** MVP uses static; items with high variance can flip to demand-forecast.
  - Pros: pragmatic; supervisor can choose per SKU.
  - Cons: two code paths.

## Recommendation

**Option C (partial GRN only) + Option X (static reorder threshold) for MVP.**

For partial ops: the jamaah-receiving-incomplete-kit scenario is a customer-facing disaster (jamaah arrives at airport without ihram because it backordered). Option A's all-or-none for kit assembly is the right constraint for MVP. Partial GRN is operationally common (vendor delivers 48 of 50 koper; accept what's there, keep line open) and doesn't touch the customer. Option C gives us partial-GRN without partial-dispatch — the right compromise.

For reorder: Option Y is the right long-term answer but depends on booking-svc having clean upcoming-departure pax counts, which it won't have on day one. Option X gets us live with a supervisor-tunable threshold in one table column; supervisor refines it over the first few departures based on actual burn rate. Option Y upgrades later by reading booking-svc and overriding `reorder_level` computationally (additive). Don't prematurely tangle logistics with booking-demand-forecasting before the basic stock loop works.

Defaults to propose: partial GRN yes; partial kit assembly no; partial dispatch no. Static reorder thresholds per-SKU, initial value set at SKU creation (supervisor input), editable anytime. Low-stock alert fires at threshold; creates a draft PR (not auto-submit) assigned to procurement officer for review. Phase 2: demand-forecast override for any SKU with `reorder_mode: auto`.

Reversibility: flipping reorder from static to auto is additive (new column on stock_items). Allowing partial kit assembly later is a policy change + UX (backorder tracking). Both are safe bets.

## Answer

TBD — awaiting stakeholder input. Deciders: warehouse supervisor (operational realism), ops lead (jamaah-experience threshold), procurement officer.
