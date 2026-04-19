---
id: Q017
title: Paket Tabungan (savings-toward-pilgrimage) — interaction with booking state
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F4, F5, F9
status: answered
---

# Q017 — Paket Tabungan interaction with booking state

## Context

PRD mentions "Paket Tabungan" / savings-toward-pilgrimage in several places:

- Line 59 B2C — Modul Kalkulator Simulasi Tabungan (calculator for savings goal)
- Line 259 Master Product — Modul Produk Finansial & Retail lists Tabungan as a product kind alongside Umroh Reguler, Haji Furoda, etc.
- Line 435 Finance — "cicilan Tabungan" payments as a recurring receipt

But the PRD doesn't say how Tabungan interacts with the booking state machine. Is an active Tabungan plan a `booking` in `partially_paid` state? Is it a separate entity that doesn't hold a seat? Does it convert to a booking when the target is hit? These are different data models.

## The question

1. **Does an active Tabungan plan hold a seat?**
   a. Yes — creating a Tabungan plan is effectively a booking in `partially_paid` status with a long pelunasan deadline (until the saver's target departure).
   b. No — Tabungan is a pre-booking savings vehicle; funds accumulate, and the customer converts to a real booking only when they choose a departure and the balance ≥ minimum DP.
2. **What happens when the savings target is reached?**
   a. System auto-creates a booking on the pre-designated departure.
   b. System notifies the customer; customer manually picks a departure and initiates a booking; accumulated balance applies as the DP/lunas amount.
3. **What if the customer drops out?**
   a. Tabungan can be withdrawn at any time (like a bank savings), minus a small admin fee.
   b. Tabungan is locked until used — withdrawal requires a refund flow similar to a booking cancellation.
4. **PSAK accounting treatment** — is Tabungan a liability (unearned income) from day one, or only recognised on booking conversion?

## Options considered

- **Option A — pre-booking savings vehicle (recommended).** Tabungan is a separate entity (`savings_plans`) that does NOT hold a seat. Funds accumulate via recurring deposits. When target reached, customer chooses a departure and creates a normal booking; the accumulated balance applies as payment. Withdrawal at any time minus fee.
  - Pros: clean separation from booking state machine; no long-running partially_paid bookings cluttering reports; accommodates customers who want to save for "Umroh some day" without committing to a specific date.
  - Cons: savers don't get the "your seat is reserved" feeling — they're saving for a future option.
- **Option B — long-running booking.** Tabungan creates a real booking in `partially_paid` against a target departure; recurring payments accumulate until lunas.
  - Pros: customer has a concrete departure committed from day one.
  - Cons: tons of `partially_paid` bookings with 2-year pelunasan deadlines; seat inventory locked without real commitment; harder to surface a customer who's actually in trouble making their payments.
- **Option C — hybrid.** Tabungan is a savings plan until the customer picks a target departure; then it converts to a real booking (Option A's mechanism) but with preferential DP terms (e.g. 10% DP instead of 20% from Q011).
  - Pros: accumulates savings flexibly; rewards planners with easier booking terms when they commit.
  - Cons: two product flows to explain.

## Recommendation

**Option A — pre-booking savings vehicle.** Tabungan plans are a separate entity from bookings, no seat held. When funds reach the customer's target, system sends a WA notification with available departures; customer picks one; accumulated savings apply as the initial payment.

Withdrawal: allowed at any time, fee = 5% of accumulated balance (configurable). Funds return to the customer within 14 business days.

Accounting (Q018 territory — for F9 finance): Tabungan funds sit as **customer-credit liability** in finance (not unearned revenue yet, because there's no booking against which to recognise revenue). On conversion to a booking, the balance moves from customer-credit to unearned revenue (standard booking treatment). This is cleaner than Option B's long-running liability-against-a-specific-departure.

Reasoning: many Indonesian pilgrims save for Umroh/Haji over years without knowing exactly when they'll go — locking them into a specific departure from day one misrepresents the commitment; seat-inventory is scarcer than money and shouldn't be held by non-committed saving; Option B creates reporting headaches (every finance report has to distinguish Tabungan-bookings from real bookings).

Reversibility: data-model choice is entrenching. Moving from Option A to Option B later would require migrating savings plans to bookings — possible but annoying. This is the most committed of the four recommendations in this batch; worth confirming stakeholder buy-in before implementation.

## Answer

**Decided:** **Option A** — Tabungan = **`savings_plan` pre-booking** (no seat); converts to normal booking with balance applied at customer-picked departure. **Withdrawal:** allowed with **5% fee (configurable 0–10%)**, **T+14 business days** payout SLA. **Accounting:** customer liability until conversion (see **Q044**). **Fee tweak** remains Super Admin config.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
