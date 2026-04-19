---
id: Q043
title: Revenue recognition mechanics — "terbang" event + refund reversal
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9, F4, F5
status: answered
---

# Q043 — Revenue recognition mechanics

## Context

PRD line 473 (module #142) states the policy clearly: jamaah funds are held as liability until the jamaah **"has flown"** — at that point revenue is recognized. This is a PSAK 72 point-in-time performance obligation satisfaction. But the PRD is silent on:

- **Which operational event counts as "flown"?** — wheels-up per flight tracking, boarding scan at gate, departure-date auto-fire by calendar, or something else.
- **What happens on post-departure refund, claim, or trip abortion** (jamaah dies in-trip, expelled by Saudi immigration, mass cancellation due to political situation).
- **Partial performance** — jamaah boards but doesn't complete full itinerary.
- **Multi-jamaah booking** where one jamaah flies but another doesn't.

These are not edge cases — the right answer here defines how revenue stabilizes on the P&L and whether a reversal is ever needed. Getting the event wrong creates an audit disaster (revenue recognized but service never delivered).

## The question

1. **Which event triggers revenue recognition?** Wheels-up detection, boarding scan at gate (F7 W10), departure-date auto-fire at 11:59 PM of the scheduled departure, or first Saudi arrival scan?
2. **Per-jamaah or per-departure?** — if recognition is per-jamaah, each pax on the manifest recognized independently; per-departure recognizes when the departure is marked completed.
3. **Partial trip handling** — jamaah flies out but misses flight back; jamaah quarantined mid-trip and doesn't reach Madinah. Does revenue stay recognized, partially reverse, or fully reverse?
4. **Post-recognition refund** — agency grants a courtesy refund after departure (unusual but possible for goodwill / legal reasons). Does revenue reverse? What's the offsetting entry?
5. **Failed departure** (e.g. package cancelled after jamaah already paid but before travel) — revenue stays at zero and refund hits Hutang Jamaah per W2; confirm this cleanly.
6. **Jamaah death in trip** — sensitive; accounting-wise, is revenue still recognized? (Service substantially delivered; insurance/refund separate question.)

## Options considered

- **Option A — Recognition at wheels-up (per-jamaah, from boarding scan event).** F7 W10 (Smart Bus Boarding) or ALL System check-in triggers `jamaah.flown_at`; finance-svc aggregates to per-booking revenue recognition.
  - Pros: most accurate — revenue only recognized when service actually begins.
  - Cons: depends on F7 field-app scans being reliable; weekend recognition when scans sync late is finicky.
- **Option B — Recognition at departure-date auto-fire (per-departure, 23:59 of scheduled date).** booking-svc fires `departure.completed` via cron at end of scheduled departure day; finance-svc recognizes all bookings on that departure.
  - Pros: deterministic; no field-app dependency; easy to reconcile.
  - Cons: recognizes revenue for jamaah who no-showed; requires a reversal workflow for actual non-travelers.
- **Option C — Recognition at first Saudi arrival scan (per-jamaah from F7 W7 ALL System / immigration arrival).** Waits until jamaah confirmed in Saudi.
  - Pros: maximal conservatism; highest delivery confidence.
  - Cons: multi-day lag from departure; introduces cross-border operational dependency.
- **Option D — Recognition over-time across the 9-day pilgrimage.** Per PSAK 72 over-time recognition; revenue amortized daily over trip duration.
  - Pros: smoother P&L; more PSAK 72 textbook.
  - Cons: Umrah is a discrete performance obligation, not time-based; harder to close monthly; complicates refund reversal.

## Recommendation

**Option A — per-jamaah recognition at F7 boarding scan (ALL System check-in or Smart Bus Boarding), with fallback to departure-date auto-fire (Option B) if scans don't arrive within 24h of scheduled departure.**

Per PSAK 72, Umrah is a point-in-time performance obligation and the earliest unambiguous delivery-start signal is boarding. Wheels-up (narrower) is ideal but requires flight-tracking integration that's out of MVP scope. F7's ALL System scan is already happening at the airport in the MVP scope — finance-svc subscribes to the same event. Departure-date auto-fire (Option B) as a fallback handles the edge case of scan failure or weekend-sync delay, with a reconciliation cron that reverses recognition if a specific jamaah is confirmed to have not flown (via no-show report).

Defaults to propose: recognition per-jamaah on `handling.all_system_checkin` or `handling.bus_boarding` (whichever fires first), captured via `booking-svc.RecognizeRevenue(booking_id, jamaah_id)` gRPC to finance-svc. Amount recognized = paid-to-date for that jamaah's share of the booking. Partial trip: service substantially delivered once boarding occurs; no partial reversal for mid-trip issues (handle as goodwill refund via Q053 path — separate journal, doesn't touch revenue). Post-recognition refund: reverse the revenue journal (compensating entry), then post the refund via W18. Failed departure (before any jamaah flown): revenue never recognized; refund simple cash-out per W2. Jamaah death: revenue stays recognized (service substantially delivered); insurance/refund handled outside revenue accounting. T+1 reconciliation cron compares flown jamaah to recognized revenue and flags discrepancies.

Reversibility: switching to Option B (simpler departure-date) later is easy (change the event subscription). Switching to Option D (over-time) is harder — requires per-line recognition lot tracking.

## Answer

**Decided:** **Option A** primary — **per-jamaah** recognition on first qualifying **F7 handling event** (`all_system_checkin` or `bus_boarding`); **fallback Option B** auto-fire **T+24h after scheduled departure** if no scan; **reconciliation job** reverses mistaken recognition. **Partial trip:** no automatic partial reversal post-boarding; **goodwill refunds** via **Q053**. **Pre-travel cancel:** no recognition. **Death in-trip:** recognition stands; refunds/insurance separate.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
