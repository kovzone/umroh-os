---
id: Q011
title: Minimum DP percentage, maximum installment count, installment cadence
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F5
status: open
---

# Q011 — DP and installment rules

## Context

PRD confirms the DP → cicilan → lunas pattern exists (lines 79, 85, 435) but leaves the actual numbers to commercial policy. Three separate parameters control this:

1. **Minimum DP percentage** — the first payment a jamaah must make to move a booking from `pending_payment` to `partially_paid`.
2. **Maximum installment count** — how many partial payments are allowed between DP and lunas.
3. **Installment cadence** — equal monthly / flexible at customer's pace / agency-scheduled.

These are customer-visible commercial rules; any of them wrong is a friction point on every sale.

## The question

1. **Minimum DP %.** Candidates: 10%, 20%, 30%, 50%. Does it vary by package kind? (E.g. Haji Furoda might need higher DP since allotments are locked earlier.)
2. **Maximum installments.** Candidates: unlimited (pay in any slices), 3, 6, 12. Unlimited is simplest; limits protect AR aging.
3. **Cadence.** Fixed monthly with auto-reminder? Or flexible — jamaah pays when they can, as long as lunas before the pelunasan deadline (Q010)?

## Options considered

- **Option A — flexible, simple.** Min DP 20% flat, unlimited installments, flexible cadence (jamaah pays when they can as long as lunas by pelunasan deadline). System sends WA reminders at H-60, H-30, H-14, H-7.
  - Pros: low friction; customer-friendly; no installment scheduling to build.
  - Cons: some customers under-pay for a long time then scramble; AR aging less predictable.
- **Option B — structured monthly.** Min DP 30%, maximum 6 installments, equal monthly auto-scheduled from DP to pelunasan deadline.
  - Pros: predictable AR; good for finance reports; jamaah know exactly what's owed when.
  - Cons: rigid; doesn't fit customers with irregular cashflow.
- **Option C — tiered by package kind.** Standard Umroh gets Option A (flexible 20%); Haji Furoda gets Option B (structured 30%+); Tabungan is its own long-term product.
  - Pros: business fit — premium programs justify stricter rules.
  - Cons: more rules to explain to CS, more UI conditional logic.

## Recommendation

**Option A — flexible, 20% DP, unlimited installments, reminders at 4 milestones.** Keep the rules simple for v1; the Tabungan product (Q017) is a separate pre-booking vehicle anyway, and Haji Furoda's multi-month lead time means customers who commit to it are financially serious already. Flexibility in installment count maps to real Indonesian customer cashflow (variable income from small businesses, UMKM, agricultural cycles) better than a rigid monthly.

WA reminder schedule: **H-60, H-30, H-14, H-7** before pelunasan deadline. Missed-payment escalation at H-7 is a CS task, not automated.

Reasoning: low friction = more sales; Option B's predictability benefit isn't worth losing sales for; we can always add structured plans later as a SEPARATE product offering (jamaah opts in to a structured plan for a discount) without forcing it on everyone.

Reversibility: min DP and installment count are config values; reminder schedule is a cron/template change. Moving from Option A to Option B is code work but not a data migration. Low-to-medium commitment.

## Answer

TBD — awaiting stakeholder input. Likely decider: agency owner + finance lead (AR management preferences).
