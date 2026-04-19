---
id: Q012
title: Refund penalty policy matrix — per package type, per cancellation timing
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F5, F4
status: answered
---

# Q012 — Refund penalty policy matrix

## Context

Module #100 Administrasi Refund & Pinalti (PRD line 341) confirms: "sistem paperwork otomatis jika ada jamaah yang batal berangkat... lengkap dengan kalkulasi potongan/pinalti untuk biaya tiket atau visa yang sudah hangus." The system must compute the refund net of:

1. **Timing-based penalty** — cancel at H-90 costs less than cancel at H-7.
2. **Sunk costs** — if the airline ticket has been issued and is non-refundable, that's a hard floor on the refund. Same for filed visas.

PRD gives no numbers. Every agency has their own policy matrix, and it's printed in the jamaah contract at booking time — this isn't engineering, it's legal. But we need the shape so F5 can compute previews (F5 W8).

Existing **Q004** covers seat-return OWNERSHIP (immediate vs gated) — that's about timing, not amount. This question covers the AMOUNT.

## The question

1. **Penalty bands by timing.** What's the matrix? Example shape:

   | Timing | Penalty (% of package price) |
   |---|---|
   | H-90+ | 10% |
   | H-60 to H-89 | 25% |
   | H-30 to H-59 | 50% |
   | H-14 to H-29 | 75% |
   | H-0 to H-13 | 100% (no refund) |

   Please confirm the bands and the percentages. Or is it a different shape (e.g. only flat 25% regardless of timing)?

2. **Does it vary by package kind?** Haji Furoda allotments are locked harder — often 100% non-refundable the moment a seat is assigned. Umroh Reguler may have more flexibility.

3. **Sunk-cost subtraction** — when the ticket is issued / visa is filed mid-cycle, the NET refund becomes `(package_price × (1 − penalty%)) − sunk_costs`. If that goes negative, customer owes the agency. Policy: do we pursue that, forgive it, or refuse the cancellation?

4. **Force majeure / departure cancelled by agency** — penalty waived. Confirm this is the policy.

## Options considered

- **Option A — explicit tiered matrix per package kind.** A config table with `(package_kind, h_days_before_departure_min, penalty_pct)` rows. CS can see the preview before confirming cancel. Finance audits the matrix annually.
  - Pros: transparent, defensible in the jamaah contract, easy to adjust.
  - Cons: initial definition requires legal/commercial input; mis-set values hurt customer goodwill.
- **Option B — flat penalty + sunk costs only.** Single flat 25% penalty + whatever sunk costs have accrued. Simpler.
  - Pros: easy to understand; no timing-cliff surprises for the customer.
  - Cons: same penalty for H-90 as H-7 doesn't match how costs actually behave; customer feels aggrieved cancelling early.
- **Option C — case-by-case (CS discretion with manager approval).** No rule; CS drafts a refund, manager approves.
  - Pros: flexibility for hardship cases.
  - Cons: no contract-grade predictability; susceptible to inconsistency / unfairness; doesn't scale.

## Recommendation

**Option A — tiered matrix per package kind,** driven by a config table. Ship with a conservative default matrix (like the example above, plus `haji_furoda` at 25/50/75/100/100 and `umroh_reguler` at 10/25/50/75/100), reviewed and signed off by agency legal before customers see it in the contract.

Force majeure / departure-cancelled-by-agency: penalty waived, sunk costs still netted (because airlines don't refund those either — it's an industry-level pass-through).

Negative net refund: agency absorbs the loss rather than billing the customer — chasing a refund debt from a jamaah who already cancelled is a reputation risk and rarely recoverable. Document this tolerance in the config as `min_refund_amount: 0` — refund floors at 0.

Reasoning: tiered matrix matches how costs actually behave over the lead time; per-package-kind lets the differently-priced programs have different policies without forcing the agency into a one-size rule; flooring at 0 net protects customer relationships.

Reversibility: matrix is a config table, editable by Super Admin; changes apply to new cancellations only (existing refunds honour their original matrix via snapshot at cancellation time). Moderate commitment — the jamaah contract references the matrix, so public-facing changes need PR management.

## Answer

**Decided:** **Option A** — **configurable tier matrix per `package_kind`** with **snapshot at cancellation time** baked into the booking/cancel record. **Ship defaults** as in Recommendation example (`umroh_reguler` vs `haji_furoda` bands). **Force majeure / agency-cancelled departure:** **penalty waived**; **non-refundable vendor/government sunk costs** still netted from refund (pass-through). **Negative net refund:** **floor at 0** to customer (agency absorbs residual loss; no debt collection from jamaah in MVP). **Legal:** printed contract / Syarat & Ketentuan must **mirror** the matrix — **lawyer review** before customer-facing publish.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
