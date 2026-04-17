---
id: Q081
title: F12 MVP scope carve — what ships vs what defers to Phase 2
asked_by: session 2026-04-17 F12 draft
asked_date: 2026-04-17
blocks: F12
status: open
---

# Q081 — F12 MVP scope carve

## Context

F12 has **zero Must Haves** across all 9 owned modules (#190–#198). Three are Should / Medium (#190, #191, #194, #199); the rest are Could / Low. The agency should consciously decide which modules ship with the MVP alongside F1–F11, vs which defer to Phase 2.

The stakes are real:
- Deferring too much loses the "Daily App" engagement narrative.
- Shipping too much balloons MVP scope for features with no Must-Have priority.
- Several modules (#196 fatwa desk, #197 forum, #195 articles) bring non-trivial content / moderation / licensing work that would be disproportionate to their Could-Have priority.

## The question

1. **Which modules ship in MVP?**
2. **Which defer to Phase 2?**
3. **Is full F12 deferral (zero modules shipped) acceptable?**
4. **If partial, what's the minimum-viable carve?**

## Options considered

- **Option A — Full MVP: ship all 9 modules (#190–#198).** Comprehensive Daily App + community + fatwa launch.
  - Pros: feature-complete; no Phase 2 gap; alumni engagement from day one.
  - Cons: significant content + moderation + licensing work for Could-Have priority modules; diverts MVP energy.
- **Option B — Lean MVP carve: ship client-side utilities + alumni shell only.**
  - **IN**: #190 prayer times, #191 qibla, #194 manasik (read-only), alumni portal shell hosting F10-owned referral entry (#199).
  - **OUT (Phase 2)**: #192 Quran (licensing), #193 dzikir, #195 articles/kajian, #196 fatwa desk (moderation + authority), #197 forum (moderation), #198 reuni board (event management).
  - Pros: smallest footprint; client-side-heavy modules ship fast; manasik is the highest-value of the Should-Have modules.
  - Cons: no community / fatwa surface on launch day.
- **Option C — Full F12 deferral: nothing in MVP; all of F12 is Phase 2.**
  - Pros: zero MVP cost.
  - Cons: zero post-trip engagement at launch; alumni immediately drift to third-party tools.
- **Option D — Minimum MVP: ship only the alumni shell + F10 referral entry; no Daily Worship utilities.** Even leaner than B.
  - Pros: absolutely minimal F12 footprint.
  - Cons: loses the Daily Worship value entirely; barely worth having F12 in MVP.

## Recommendation

**Option B — lean MVP carve: ship prayer times, qibla, manasik read-only, alumni shell. Defer Quran, dzikir, articles, fatwa, forum, reuni.**

Option A over-commits to content + moderation + licensing work for Could-Have modules during a period where every dev-hour should push Must-Have path (F1–F10 implementation). Option C is defensible but gives up a meaningful differentiator on launch day — other Umrah agencies do offer at least prayer times + manasik as a value-add. Option D is too lean; dropping Daily Worship loses the Daily App branding justification.

Option B is the carve that matches the MoSCoW priorities: the three Should-Have / Medium modules (#190, #191, #194) ship; everything Could-Have waits. This preserves the Daily App promise via the three modules that are device-sensor + static-content heavy (all client-side work, minimal backend), while deferring the modules with non-trivial stakeholder decisions (Q077 moderation, Q078 fatwa authority, Q079 content licensing at scale).

Defaults to propose:

**IN for MVP (ship alongside F1–F10 implementation):**
- #190 Prayer times + adzan push (Should / Medium)
- #191 Qibla compass (Should / Medium)
- #194 Manasik encyclopedia read-only (Should / Medium)
- Alumni portal shell (auth + landing) hosting F10-owned referral entry (#199 — F10-owned mechanics)
- Engagement metrics aggregation for F11 dashboard (minimal — just the 3 IN modules)

**OUT for Phase 2 (revisit after MVP is operating):**
- #192 Quran digital — needs Q079 licensing finalization + offline-audio strategy
- #193 Dzikir & doa — needs Q079 content licensing + religious advisor edition approval
- #195 Articles + kajian — needs content-authoring pipeline + kajian schedule management
- #196 Fatwa Desk — needs Q078 ustadz panel + authority + SLA commitments
- #197 Forum grup angkatan — needs Q077 moderation infrastructure + community manager
- #198 Reuni board — needs event management + RSVP + check-in infrastructure

**Phase 2 kick-off trigger**: revisit F12 OUT-of-MVP modules when (a) MVP is live, (b) at least 3 departure cohorts have completed (enough alumni to justify community features), (c) Q077, Q078, Q079 are answered by stakeholders.

Reversibility: Phase 2 modules are additive — shipping them later doesn't affect IN modules. Bringing an OUT module forward into MVP is a scope change; reviewer's call.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (MVP scope appetite), CTO (MVP dev-capacity), marketing lead (launch-day engagement expectation).
