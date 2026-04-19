---
id: Q057
title: UTM attribution model — window, first vs last click, tiebreakers
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10, F4, F10
status: answered
---

# Q057 — UTM attribution model

## Context

Every booking that arrives through marketing channels (ads / replica site / campaign landing page) must be attributed to the right source — drives ROAS (F11), commission routing (F10 W9), and budget reconciliation (F9 W8).

The PRD captures UTM (Q49 already notes UTM Manager module #49) but doesn't specify:
- **Attribution window length** (how long after last touch does the click count?).
- **First-click vs last-click vs linear attribution model**.
- **Tiebreaker between multiple agents** the same lead touched.

Q019 (existing) covers abandoned-checkout attribution at a higher level; Q057 is the F10-specific attribution-logic detail.

## The question

1. **Attribution window length** — 7 days? 30 days? 90 days? Variable by campaign type?
2. **Attribution model**:
   - **First-click** — first source the lead ever came from gets all attribution.
   - **Last-click** — last source before booking gets all attribution.
   - **Linear** — equal share across all touchpoints.
   - **Position-based** — e.g. 40% first + 40% last + 20% middle.
3. **Tiebreaker when two agents touched same lead** — most recent wins? Agent with WA reply wins? Round-robin of touches?
4. **Direct-visits attribution** — lead visits `/id/<agent1>` then 2 days later visits `domain.com` directly without UTM, then books — who gets credit?
5. **Organic / no-UTM** bookings — attributed to agency (no commission) or most-recent-any-agent?
6. **Abandoned checkout mid-funnel re-engagement** (Q019 cross-ref) — if agent A captured the lead but agent B closed via WA direct reach, who wins?
7. **Historical re-attribution** — if attribution model changes in the future, does that restate past commissions?

## Options considered

- **Option A — 30-day last-click with first-click fallback; agent-touch wins over organic.** Within 30-day window, last source with identifiable agent or UTM wins. Fallback to first-click if no last-click exists within window. Organic (no UTM, no agent) booking = agency, no commission.
  - Pros: matches industry norm (Meta, Google Ads default); simple.
  - Cons: under-credits early funnel touch (brand awareness) work.
- **Option B — 30-day first-click.** First source a lead ever came from wins for all bookings within the window.
  - Pros: rewards acquisition work; clean for replica-site traffic.
  - Cons: ignores nurturing agent who closed the deal; under-motivates close-side activity.
- **Option C — 30-day first-click for attribution + last-click for commission routing.** Dual attribution: first-click wins for marketing analytics / ROAS; last-click wins for agent commission routing.
  - Pros: marketing and sales teams both see what they care about.
  - Cons: complex to explain; two records per booking.
- **Option D — 30-day last-click by default; 90-day first-click for agent-program attribution.** Different windows for different purposes.
  - Pros: generous for agent nurturing (long sales cycles for Umrah); still matches ads convention.
  - Cons: mixed models cause confusion.

## Recommendation

**Option C — 30-day window; first-click for marketing analytics + ROAS; last-click for commission routing; tiebreaker = most-recent-agent-touch within the window.**

Option A is simple but gives no credit to the replica-site / top-of-funnel agent work; Umrah has long decision cycles (30–90 days typical) and agents invest significant nurturing before close. Option B flips that, but then the closing agent (often CS) gets no recognition. Option C gives both teams honest metrics: marketing sees the source that brought the lead; the agent who closed gets the commission they earned. Each booking has two attribution records; reports render the right one.

Defaults to propose: 30-day window from lead's first-touch timestamp. Both `first_touch_utm` and `last_touch_utm` captured on every lead + booking. Commission routing uses last-touch within-window agent; if no agent in last-touch but first-touch has one, falls back to first-touch (agent gets nurture credit); if no agent anywhere, no commission (organic/agency). Tiebreaker: most recent agent touch timestamp wins. Direct-visit after agent touch: if within-window, last touched agent retains (direct visit doesn't replace UTM touch); outside window, becomes organic. Abandoned-checkout re-engagement (Q019 resolution will refine): default aligns with Q019's decision — if unified by Q019, current lead keeps original attribution; if split, each arm carries its own. Historical re-attribution: attribution snapshot frozen at booking time; model changes are prospective only.

Reversibility: window length and model are config. Switching from 30 to 90 days is a one-setting change. Switching attribution model retroactively is discouraged (breaks historical reports) but doable via snapshot re-computation.

## Answer

**Decided:** **Option C** — **30-day window**; **marketing analytics = first-touch UTM**; **commission routing = last-touch agent/ref** (consistent **Q019**); **direct visit inside window does not erase agent**; **model changes prospective only** (snapshot at booking).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
