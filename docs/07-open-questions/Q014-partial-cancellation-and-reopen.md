---
id: Q014
title: Partial cancellation (one jamaah out of many) + reopening a cancelled booking
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F4, F5
status: answered
---

# Q014 — Partial cancellation and reopening

## Context

PRD describes cancellation at the booking level (Module #100, Alur Logika 6.5) but doesn't explicitly cover two common real-world cases:

1. **Partial cancellation** — one jamaah in a multi-pilgrim booking can no longer travel (illness, passport issue, family emergency). The rest want to proceed. What happens to the booking?
2. **Reopening** — a jamaah cancels, then later wants back in. Is the original booking reopened, or do they start a new booking?

Both are policy + implementation decisions that affect the status machine design and the refund calculation.

## The question

1. **Partial cancellation — the model.** Three candidates:
   a. **Shrink existing booking** — cancelled jamaah's `booking_item` becomes `cancelled`, booking pax count decreases, per-jamaah refund calculated against the cancelled jamaah's share of the booking total (penalty + sunk costs per Q012).
   b. **Split into two** — original booking keeps the continuing jamaah; a new `cancelled`-status booking is created for the leaving jamaah for record-keeping.
   c. **Forbidden** — partial cancellation not allowed; the whole booking cancels and the rest re-book fresh.
2. **Room/bus reallocation after partial cancel.** If Smart Grouping has already run (Q015), the leaving jamaah leaves a hole. Does the system auto-rerun grouping, leave the hole, or ask ops to handle manually?
3. **Reopening a cancelled booking.** Three candidates:
   a. **Allowed with original price/seat** if still available — customer literally gets their old seat back, pays back what they got refunded plus any penalty.
   b. **Allowed as a new booking only** — forbid reopening the old record; require a fresh booking at current price/availability. Old record stays cancelled for audit.
   c. **Allowed within a grace window** (e.g. 48h after cancellation, if it was a mistake) — reopen; beyond that, fresh booking.

## Options considered

Combining the two sub-questions:

- **Option A — Shrink + fresh re-book.** Partial cancel shrinks the existing booking (1a). Reopening forces a fresh booking (2b). Smart Grouping re-runs automatically on shrink if it had already been run.
  - Pros: clean data model; no reopen-with-stale-price edge cases; auto-regrouping keeps ops tidy.
  - Cons: customers who "changed their mind quickly" after cancelling face full re-booking friction even when nothing else has moved.
- **Option B — Shrink + grace-window reopen.** Partial cancel shrinks (1a). Reopening allowed within 48h as "cancel reversal" (2c); after 48h, fresh booking.
  - Pros: accommodates the "oops, meant to cancel a different trip" case; operationally simple to implement (just a status flip if seats still available).
  - Cons: complicates the state machine (a cancelled booking can become active again); commission attribution during the reversal window is ambiguous.
- **Option C — Forbid partial cancel + fresh re-book.** Full booking cancellation only (1c); reopening is a fresh booking (2b).
  - Pros: simplest model; no per-jamaah refund math.
  - Cons: operationally painful for real-life scenarios — one jamaah in a family of five can't travel, the agency doesn't want to cancel the whole trip.

## Recommendation

**Option B — Shrink existing booking on partial cancel + grace-window reopen (48 hours).** On partial cancel: the leaving jamaah's `booking_item` status becomes `cancelled`, per-jamaah refund is calculated via Q012's matrix on the jamaah's pro-rata share of the booking total, and the booking's `total_amount` + `paid_amount` reduce proportionally.

Room/bus reallocation after shrink: **leave the hole, flag to ops**. Auto-rerunning the grouping algorithm during a live departure cycle risks reshuffling other jamaah's assignments they've already been told about (room numbers, seat assignments on ID cards). Let the ops admin review + manually adjust.

Reopening within 48h: allowed if (a) seat capacity is still available, (b) price hasn't changed, (c) no refund has been paid out yet (refund still in `requested` or `processing`, not `completed`). If reopening: status returns to pre-cancel, refund request is voided, commission state reverses.

Reopening beyond 48h or if any of a/b/c fails: force fresh booking at current price.

Reasoning: shrink + grace window matches real human behaviour (people change their minds quickly, then settle); the 48h window is short enough that operational impact is bounded; forcing fresh bookings after 48h aligns customer price expectations with current market.

Reversibility: shrink policy is implementation, not easily reversible after launch (customers will have expectations set). Reopen grace window is a config value and can shift. Medium commitment on the partial-cancel model; low on the reopen window.

## Answer

**Decided:** **Option B** — **shrink booking** on partial cancel (per-jamaah `booking_item` cancelled + pro-rata refund per Q012); **reopen within 48h** only if **capacity**, **price unchanged**, and **no refund payout completed**; else **new booking**. **Smart Grouping after shrink:** **do not auto-rerun** — **flag to ops** (avoids surprise reshuffles).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
