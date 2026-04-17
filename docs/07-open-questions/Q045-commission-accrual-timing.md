---
id: Q045
title: Commission accrual timing — booking / paid / departure / payout
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9, F10
status: open
---

# Q045 — Commission accrual timing

## Context

PRD Section B (B2B agent network) describes a Dompet Komisi Agen with three states: **Pending** (jamaah DP), **Confirmed** (jamaah lunas), **Withdrawal** (paid out to agent bank). Module #146 (F9 line 481) covers the payout mechanics. But the PRD does not specify **when** the agent commission becomes a recognized expense / liability on the GL — a question that's foundational to F9 because it determines when Beban Komisi Agen hits the P&L vs just sits as contingent.

PSAK 57 (provisions) is the relevant standard: a liability accrues when it's probable an outflow will be required and the amount is reliably estimable. Commission on a DP'd booking is "probable" but not certain (jamaah might cancel); commission on a fully-paid booking is near-certain; commission post-departure is irreversible.

## The question

1. **At which event does Beban Komisi Agen accrue to P&L?**
   - On booking created (earliest — highest cancellation risk).
   - On jamaah DP (booking confirmed, ~50% paid).
   - On jamaah Lunas (booking paid-in-full) — matches PRD "Confirmed" state.
   - On departure / revenue recognition — matches matching principle with revenue.
   - On payout (cash basis) — simplest, non-GAAP.
2. **Clawback on cancellation / refund** — if commission already accrued and jamaah refunds, does the accrual reverse? At what stage is clawback blocked?
3. **Clawback on departure-date no-show** — edge case; policy decision.
4. **Override commission** (super-agent override on sub-agent sales, PRD B2B module) — same timing as primary commission, or different?
5. **Partial commission** — is it possible for commission to accrue on partial jamaah payment (e.g. 30% of agent commission recognized when jamaah is 30% paid)?

## Options considered

- **Option A — Accrue on booking created (earliest).** Dr Beban Komisi Agen / Cr Hutang Komisi Agen at booking_created event.
  - Pros: matches the commission economics from day one; no surprises for the agent.
  - Cons: high reversal rate as some bookings cancel pre-DP; noisy P&L; inflates Hutang Komisi balance misleadingly.
- **Option B — Accrue on jamaah Lunas (fully paid).** Matches PRD's "Confirmed" wallet state.
  - Pros: sharply reduces cancellation-reversal rate; aligns with "revenue is more probable" moment.
  - Cons: doesn't match expense recognition to revenue recognition (matching principle would prefer departure).
- **Option C — Accrue on departure (matched to revenue recognition per Q043).**
  - Pros: matching principle compliance; clean P&L-per-departure (job-order costing includes commission in the right period).
  - Cons: agent payouts typically expected faster than revenue recognition — creates UX mismatch for agents who expect commission shortly after sale.
- **Option D — Accrue at Lunas, pay out at Lunas + N days, reverse on pre-departure refund.** Hybrid practical model.
  - Pros: matches industry practice; agent satisfaction; limited reversal window.
  - Cons: still has a timing gap vs departure-date revenue recognition.

## Recommendation

**Option B — accrue at jamaah Lunas (fully paid) with a 30-day post-departure clawback window.**

Option A creates unnecessary P&L noise — roughly 15-20% of initial bookings don't reach DP in typical Umrah agency data, and reversing those commissions monthly is busy-work for finance. Option C is the textbook PSAK 72 matching-principle answer but violates the UX expectation that agents get paid shortly after their jamaah pays — delaying payout to post-departure (typically 2-3 months after sale) creates agent churn. Option B lands the accrual at the "near-certain" signal (Lunas) while keeping a narrow clawback window for genuine post-departure refunds (30 days — covers most legitimate refund scenarios).

Defaults to propose: accrue Beban Komisi Agen / Hutang Komisi Agen on `booking.paid_in_full` event. Payout runs monthly (W14) — agent receives payment at month-end of the confirmation month (pays agents roughly 15-30 days after Lunas). Clawback applies pre-payout (common) — reverse accrual, no cash movement. Post-payout refund (rare) — within 30 days of original accrual, reverse and demand-back (agent wallet negative balance netted against next payout); after 30 days, expense stays — classify as `Beban Komisi Forfeit` (distinct from regular commission expense for reporting). Override commissions follow the same timing as the underlying primary commission. No partial commissions; commission is all-or-nothing at Lunas.

Reversibility: moving to Option C (departure-based) later is a non-trivial restatement — past accrual periods would need re-classification. Moving to Option A (booking-based) is easier (just more accruals plus reversals). Q045 should be pinned early.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (agent-satisfaction vs accounting-rigor balance), CRM lead, finance director, external auditor.
