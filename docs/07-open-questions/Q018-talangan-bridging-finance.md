---
id: Q018
title: Talangan (bridging finance) — process, data model, and accounting treatment
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F5, F9
status: open
---

# Q018 — Talangan (bridging finance) process

## Context

PRD line 259 lists "Talangan" under Modul Produk Finansial & Retail alongside Tabungan. Talangan is an Indonesian term for short-term bridging finance — in the pilgrimage context, it typically means: the agency (or a partner institution) fronts money so the customer can book now and pay back over time.

The PRD doesn't spell out how it works — which party extends the credit, what the interest rate / fee is, how the system models the receivable, what happens on customer default.

This is a Sharia-compliant product for Islamic finance contexts (some structures like Ijarah or Murabahah exist). The agency's specific arrangement needs stakeholder clarification.

## The question

1. **Who provides the Talangan credit?**
   a. The agency itself (agency absorbs the risk).
   b. A partner bank or cooperative (BMT, Koperasi Syariah) the agency refers the customer to.
   c. Both — agency offers in-house for small amounts; partner for larger.
2. **Does Talangan mark a booking as `paid_in_full` immediately, or as `partially_paid` until the customer finishes paying back?**
3. **Interest / fee structure.** Is there a markup (Murabahah-style) or a fixed service fee? What's the typical cost to the customer?
4. **Collections.** If the customer defaults after travelling, what's the recovery process? Does the agency block future bookings?
5. **Accounting** (F9 territory). Is the receivable from the customer on the agency's books, or off-balance-sheet (passed to the partner)?

## Options considered

- **Option A — off-balance-sheet partner referral (recommended for MVP).** Agency partners with one or more Sharia-compliant financiers. Customer applies to the partner through the B2B/B2C portal; partner approves; partner disburses to the agency; agency marks booking paid_in_full; customer repays the partner directly. The agency's books show a normal booking payment; the loan is on the partner's balance sheet.
  - Pros: zero credit risk for the agency; simpler data model (booking looks like a normal cash payment); scales without tying up agency capital.
  - Cons: dependent on partner approval timing; customer UX includes an extra application step; agency gets a referral fee instead of direct financing margin.
- **Option B — in-house financing.** Agency books the receivable itself; booking marked `partially_paid`; customer pays the agency installments post-departure.
  - Pros: more margin to the agency; simpler customer journey (one relationship).
  - Cons: credit risk on agency books; requires AR collections capability; regulatory complexity for Sharia compliance.
- **Option C — hybrid by amount.** In-house for < Rp 5M bridging; partner referral for larger amounts.
  - Pros: balances risk and margin.
  - Cons: two product flows; two accounting treatments.

## Recommendation

**Option A — off-balance-sheet partner referral.** For MVP, treat Talangan as a **customer-facing financial product surfaced via the portal** where the customer applies to an external Sharia-compliant partner, partner approves, partner disburses to the agency as a normal cash payment. The booking is marked paid like any normal payment. The agency receives a referral fee (negotiated with the partner, typically 1–3%) treated as revenue.

Data model:
- Booking has a `payment_source` enum extending `va | qris | card | bank_transfer | manual_cash | talangan_partner`.
- On `talangan_partner`, the `invoice` carries a reference to the partner's loan reference ID for audit.
- No receivable on agency books — payment_events looks identical to a cash payment from the agency's perspective.

Accounting (F9): referral fee flows as service revenue; no loan-related journal entries; no collection workflow.

Reasoning: MVP credit-risk-free; regulatory complexity (Indonesian Sharia finance rules, OJK compliance) deferred to partners who are licensed for it; if in-house financing is commercially attractive later, Option B can be added as a second payment source without disrupting the partner-based flow.

Reversibility: adding in-house financing later is additive (new payment source + new AR / collections workflow); switching from A to B wholesale would require a policy change but no data migration. Medium commitment.

## Answer

TBD — awaiting stakeholder input. Likely decider: agency owner + finance lead + Sharia compliance advisor (if in-house financing is considered).
