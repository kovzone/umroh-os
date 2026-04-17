---
id: Q058
title: Alumni referral reward economics
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: open
---

# Q058 — Alumni referral reward economics

## Context

PRD module #66 + #199 mention alumni referral program — alumni issues a code, new jamaah uses it, alumni gets rewarded. PRD line 227 describes reward as *"poin diskon atau cashback"* — either/or, not specified. Economics (amount, cap, funding source, expiry) are silent.

Alumni referral is distinct from agent commission:
- Agent commission is a commercial relationship; alumni referral is gratitude + viral marketing.
- Agent commission is ongoing income; alumni referral is a one-time reward per successful referral.
- Agent commission has tax implications (PPh 21); alumni referral is typically marketing expense.

Getting the economics wrong creates either (a) low uptake (reward not compelling) or (b) abuse risk (gaming the system for large rewards).

## The question

1. **Reward format** — cashback (cash to bank), discount voucher (redeemable on next Umrah), points (accumulate, redeem later), or combination?
2. **Reward amount** — flat IDR per referral? Percentage of referred booking value? Sliding scale?
3. **Funding source** — marketing budget (expense), reduction of referred jamaah's price (less revenue), or agency pocket (net reduction of margin)?
4. **Trigger point** — reward issued on referred jamaah's booking creation, DP, full payment, or departure?
5. **Cap per alumni** — maximum rewards per year? Lifetime cap? Minimum interval between rewards?
6. **Alumni qualification** — any past jamaah, or only those who completed trip + gave testimoni?
7. **Referral expiry** — issued code valid for how long?
8. **Abuse detection** — same KTP/phone/address self-referral; friend-ring gaming; merchant fraud (alumni trading codes for other rewards).
9. **Stacking with agent commission** — if a referred jamaah books through an agent-replica site using alumni code, do both rewards apply?

## Options considered

- **Option A — Flat cashback (500K IDR) applied to alumni's next Umrah booking; 3/year cap; 2-year expiry.**
  - Pros: simple; alumni-friendly; stickiness (reward tied to booking with agency).
  - Cons: cashback-as-credit can be confusing; requires future booking to redeem (some alumni never book again).
- **Option B — Cash transfer to alumni bank account (200K IDR net after tax).**
  - Pros: universal value; no strings.
  - Cons: triggers tax implications (possibly PPh 21); lower amount vs credit form; no stickiness.
- **Option C — Points-based system: 1 point per successful referral, redeemable for merchandise / discount tiers.**
  - Pros: gamification; flexible redemption.
  - Cons: requires merchandise/reward catalog; accounting for liability of unredeemed points.
- **Option D — Flat discount on referred jamaah's booking (500K off) + 250K credit to alumni for next booking.**
  - Pros: dual incentive (alumni + referred jamaah); immediate win for both.
  - Cons: discount reduces booking revenue (affects F9 revenue recognition); coordination across two parties.

## Recommendation

**Option A — 500K IDR cashback credit to alumni's next Umrah booking + 3/year cap + 2-year expiry; funded as marketing expense; triggered at referred jamaah's paid_in_full.**

Option B's cash transfer adds tax complexity for relatively small amounts and loses stickiness (no return-booking incentive). Option C's points system requires a redemption catalog + merchandise procurement + points-liability accounting — too much overhead for MVP. Option D's dual-reward is nice but complicates the referred jamaah's booking price (revenue accounting messiness). Option A is the "Dropbox free-storage for referral" pattern adapted to Umrah — credit on next booking creates a booking-back loop.

Defaults to propose: **reward** = 500K IDR credit applied to alumni's next Umrah booking (stored in `alumni_referral_credits` ledger; consumed on next booking automatically). **Trigger**: referred jamaah reaches `paid_in_full` status (Q045's commission trigger alignment — reduces reversal risk). **Cap**: 3 successful referrals per alumni per calendar year; lifetime no cap. **Expiry**: credit valid 24 months from issue; unused credits auto-forfeit. **Qualification**: any past jamaah who completed one Umrah trip (departure-date in the past) is eligible; no testimoni requirement (avoid gating reward on favorable review). **Funding**: marketing expense — Dr Beban Promosi (6.x.x) / Cr Piutang Kredit Alumni (1.1.x) at reward issuance. **Abuse detection**: same KTP/phone/address match between alumni and referred jamaah blocks issuance (auto-detected, reviewed case-by-case if supervisor flags); friend-ring patterns (multiple alumni all referring each other) flagged via graph analysis in Phase 2. **Stacking with agent commission**: yes — if referred jamaah books through an agent-replica site using alumni code, both rewards apply independently (alumni gets credit, agent gets commission). Double-reward is a feature (viral + agent-channel both incentivized).

Reversibility: reward amount, cap, expiry, and trigger are all config in `alumni_referral_policy` table — editable by marketing admin with audit.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (marketing-budget + repeat-booking strategy), CRM lead, finance director (marketing-expense posture).
