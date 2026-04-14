---
id: Q004
title: Cancellation → seat return ownership (catalog auto vs F5 refund-gated)
asked_by: session 2026-04-14 F2 draft
asked_date: 2026-04-14
blocks: F2, F4, F5
status: open
---

# Q004 — Cancellation → seat return ownership

## Context

When a booking is cancelled, the reserved seats on its `package_departure` must eventually be returned to the available pool so someone else can book them. But **when** the return happens — and which service drives it — changes the architecture meaningfully.

Two candidates:
- `catalog-svc` auto-returns seats the moment a booking hits `cancelled` status.
- `payment-svc` drives the return as part of the refund saga in `broker-svc`: seats only come back when the refund is actually processed, not when the cancellation is merely requested.

The difference matters in real life: a DP-paid customer requesting cancellation on a popular departure might lose their seat to someone else before the refund is settled, or might cause overcapacity if seats are released optimistically and the cancellation is later reversed.

## The question

1. On booking cancellation, when are seats returned to the available pool?
   - **Immediate** — at the moment the booking status becomes `cancelled`
   - **On refund settled** — only after `payment-svc` confirms the refund has been disbursed to the customer
   - **Conditional** — immediate for DP-only cancellations, on-refund for fully-paid cancellations
2. If the refund later fails (e.g. bank rejection), does the seat go back to "reserved" state or stay available?
3. Who owns the reversal path if a cancellation is reversed (customer changes mind within the grace window)?

## Options considered

- **Option A — Immediate return (Recommended default).** Seats return the moment a booking is cancelled. Refund follows its own flow asynchronously. If a cancellation is reversed, the system attempts to re-reserve; if seats are no longer available, the reversal is denied.
  - Pros: simple state machine; maximises seat utilisation on popular departures.
  - Cons: a flaky cancellation-then-reversal can cause friction; accounting has a window where seats are free but money hasn't been refunded yet.
- **Option B — Refund-gated return.** Seats stay reserved until the refund is actually disbursed. Prevents the timing mismatch.
  - Pros: seat state perfectly tracks money state.
  - Cons: popular departures can sit with ghost-reserved seats for days; complicates the saga.
- **Option C — Conditional.** Immediate for draft / DP-unpaid bookings (no money in play); refund-gated for DP-paid and fully-paid bookings.
  - Pros: best of both — no ghost seats for free bookings, tight sync for paid ones.
  - Cons: two code paths; slightly more test surface.

## Recommendation

**Option A — Immediate return, owned by `catalog-svc`.** Pilgrimage departures have fixed capacity and short windows; holding seats as "ghost-reserved" through a multi-day refund settlement is the most expensive failure mode — someone who could have filled that seat walks away. Money and seats genuinely do decouple in real operations: a cancellation decision is made quickly, while refund disbursement is paced by bank timing. Modelling them as a single coupled state overfits to an edge case at the expense of the common case.

Implementation: the saga in `broker-svc` signals `catalog.ReleaseSeats` as a compensation step the moment a booking enters `cancelled` status. Refund proceeds independently in F5 on its own timeline. Reversal within a grace window (suggest 24 hours) is handled by `catalog.ReserveSeats` attempting to re-reserve; if seats are gone, the reversal errors out and the customer is informed in the UI — this is honest behaviour and matches how airlines already operate.

Option B (refund-gated) is the "correct" answer in an accounting sense but it creates real operational harm on popular departures. Option C (conditional) is a hedge that doubles the code paths without solving the underlying tension; a cleaner conditional can be added later via policy config if real operators complain.

Reversibility: switching from Option A to Option B later means the saga waits for a `payment.refund_completed` signal before calling `ReleaseSeats` — that's a few lines in `broker-svc`, not a data migration. Low commitment; start with the operationally-friendlier default.

## Answer

TBD — awaiting stakeholder discussion.
