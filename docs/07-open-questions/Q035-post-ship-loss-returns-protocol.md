---
id: Q035
title: Post-ship loss / damage / returns-from-trip protocol
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8
status: open
---

# Q035 — Post-ship loss, damage, and returns-from-trip protocol

## Context

PRD module #125 *Retur & Penukaran* (line 417) covers wrong-size uniform and damaged koper pre-departure. PRD module #128 *Laporan Kerusakan* (line 427) handles field damage/loss adjustments. But the PRD is silent on:

- **Lost-in-transit** — courier loses the package between warehouse and jamaah address.
- **Damaged-on-arrival** — jamaah receives a broken koper or defective ihram.
- **Returns-from-trip** — agency-owned gear (walkie-talkies mentioned at PRD line 567, "Receiver Asset Recovery" in F7 field execution; possibly thermoses, signage, banners) that come back from Saudi after the pilgrimage and need to be restocked or written off.

These are three distinct flows with different responsible parties (courier vs agency vs jamaah) and different accounting treatments (insurance claim vs warranty replacement vs routine return). F8 W14 references Q035 for the umbrella answer.

## The question

1. **Lost-in-transit** — who absorbs the cost (courier liability via insurance claim, or agency writeoff)? What's the SLA to declare "lost"? What replacement workflow fires?
2. **Damaged-on-arrival** — jamaah reports damage within N days. Does a field rep inspect? Is replacement automatic? Who pays (courier insurance, vendor warranty, agency)?
3. **Returns-from-trip** — which items are expected to return? (Agency-owned: receivers, banners, signage. Jamaah-owned: koper, ihram, kit items.) Is there a per-item manifest at return? Who receives at the warehouse? How are they sorted between restock / refurbish / writeoff?
4. **Courier insurance claim process** — does `logistics-svc` automate the claim filing, or does ops submit manually and we just record the outcome?
5. **Replacement kit triggering** — does damaged-on-arrival trigger a full re-dispatch (W11 again), or partial component replacement only?
6. **Accounting** — how does the cost flow to finance? Writeoff line? Courier refund offset? Vendor warranty credit?

## Options considered

- **Option A — Minimal: treat everything as manual stock_adjustment.** No automated claim workflows. Ops investigates each case and commits a stock_adjustment with reason + notes. Finance sees the expense.
  - Pros: minimal code; flexible for edge cases.
  - Cons: no audit trail for courier vs vendor responsibility; manual effort per case; risk of uncollected insurance.
- **Option B — Structured claim workflows per category, but manual filing.** Separate `loss_claims` table with columns for category (lost/damaged-arrival/returns-from-trip), responsible party, claim status, claim amount, resolution. Workflow forms per category. Ops files claim with courier/vendor manually; records outcome in UmrohOS.
  - Pros: audit trail per claim type; visibility into unresolved claims; separable accounting.
  - Cons: more schema and UI; still manual to external parties.
- **Option C — Full automation including courier claim API integration.** Integrate courier claim APIs where available. UmrohOS files claim on ops approval; tracks status; reconciles refund against original payment.
  - Pros: least ops friction; fastest recovery.
  - Cons: every courier's claim API is different (or non-existent); significant scope creep for MVP.

## Recommendation

**Option B — structured claim workflows per category, manual filing with external parties.**

Option A leaves money on the table — if the courier loses an 8M IDR koper shipment, the claim workflow is what prevents it quietly becoming a writeoff. Option C is over-engineered for MVP: agencies already have relationships and claim processes with their courier and vendors; automating those is scope creep without clear return. Option B gives the audit trail and visibility without trying to replace human-to-human insurance workflows.

Defaults to propose: three categories `lost_in_transit`, `damaged_on_arrival`, `returns_from_trip`. Each has a responsible-party field (courier / vendor / agency / jamaah). SLA: lost-in-transit declared at 14 days past expected delivery (per courier guideline); damaged-on-arrival accepted up to 7 days post-delivery. Replacement kit dispatched on supervisor approval. Returns-from-trip: per-departure return manifest compiled by tour leader in F7 app → handed to warehouse supervisor on return; per-item sorting (restock / refurbish / writeoff) with photos for damaged items. Accounting: stock_adjustments post with reason, claim_outcome field populated when resolved (insurance pay / warranty / writeoff). Receiver Asset Recovery (F7 referenced walkie-talkies) treated as a returns-from-trip case.

Reversibility: the claim schema is additive. Adding courier API integration later is additive (status field flips from manual to automated).

## Answer

TBD — awaiting stakeholder input. Deciders: warehouse supervisor, logistics manager, finance (insurance claim history), ops lead (field return experience).
