---
id: Q063
title: Testimoni moderation policy
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10, F7
status: answered
---

# Q063 — Testimoni moderation policy

## Context

Testimoni flow (rating + review) is captured in F7 Alur 8.7 (post-trip review). Testimoni content feeds:
- B2C gallery (Module #41, PRD line 1225)
- Agent replica-site galleries (#33)
- Central content library (Bank Konten #31)
- Marketing campaigns (case studies, social proof)

Unmoderated testimonies risk:
- Offensive / inappropriate content on public-facing surfaces.
- Complaints escalated as public reviews (reputational damage).
- Competitor references, spam, off-topic rants.
- UU ITE liability for hosted defamatory content.

But over-moderation risks:
- Cherry-picking only positive reviews (authenticity loss, reputational risk when discovered).
- Review-publishing delay (loses timeliness for campaigns).
- Review fatigue from jamaah ("why submit if edited heavily").

## The question

1. **Moderation model** — pre-approval required before public display? Post-approval (publish-then-review)? Hybrid?
2. **Edit authority** — can moderators edit testimoni text (typo fix), or only approve/reject as-is?
3. **Rejection reasons** — does jamaah get feedback on rejection? Can they resubmit?
4. **Negative review handling** — agency gets a 2-star review; does it get published? Hidden? Triggered for follow-up?
5. **Moderator role** — who moderates (marketing admin, ops lead, CS supervisor, customer experience lead)?
6. **SLA for moderation** — how quickly does a submitted testimoni need to be reviewed?
7. **Featured selection** — per Q062, agents pick which approved testimonies feature on their replica sites; do ops restrict choice?
8. **Media (photo / video) moderation** — different rules than text (recognizability, consent)?

## Options considered

- **Option A — Pre-approval required; moderator reviews all; rejected with reason; resubmission allowed.** Strict gate.
  - Pros: cleanest public-facing surface; liability bounded.
  - Cons: delay costs marketing timeliness; rejected reviews may strain relationship.
- **Option B — Post-approval: published immediately; moderator sweeps daily; inappropriate removed.** Speed over safety.
  - Pros: fast; real-time freshness.
  - Cons: window of exposed inappropriate content; liability.
- **Option C — Pre-approval for public-facing surfaces, direct-publish for agency-internal view.** Testimoni submitted by jamaah goes to internal dashboard immediately (ops sees), pre-approval needed before going public.
  - Pros: ops gets signal fast; public surface protected.
  - Cons: dual-state; confusing UX.
- **Option D — Score-based: 4+ stars = auto-approve + publish; 1-3 stars = moderator review + customer-service response workflow.** Use the rating as a trust signal.
  - Pros: fast for positive reviews; careful for negative ones.
  - Cons: biases public view toward positive; authenticity concerns.

## Recommendation

**Option C — pre-approval required for public-facing surfaces; direct-publish to agency-internal dashboards; negative reviews trigger CS follow-up workflow before public display decision.**

Option A's universal pre-approval is the safe default but creates a bottleneck — 24-hour moderation SLA on a flood of testimonies after a big departure delays content pipeline by days. Option B's post-approval accepts brief exposure risk and relies on a sweep cadence that's hard to maintain. Option D's score-gating is tempting but creates an "authenticity" problem — if only 4+ star reviews show publicly, savvy readers notice. Option C's split gets the internal signal fast (ops sees low ratings immediately to intervene) while keeping public display moderated.

Defaults to propose: **Submission** — jamaah submits testimoni via Daily App post-trip (F12 surface) or portal; lands in `testimonies` with status `pending`. **Internal dashboard** — all new testimonies visible to ops + marketing immediately (ops can triage negative reviews fast). **Pre-approval for public** — marketing admin reviews; 4-state: `pending`, `approved`, `rejected`, `edit_requested`. Moderator can approve as-is, request edit (sends to jamaah with reason), or reject (with reason). **Negative review (≤ 3 stars)** — not outright blocked; CS assigned a follow-up ticket (offer resolution); after CS response, marketing decides whether to publish (with context) or keep private. **SLA** — 24h SLA for initial decision; 72h total including re-submission cycle. **Moderator role** — marketing admin default; CS lead and ops lead can also moderate. **Featured selection (Q062)** — agents pick from approved pool; ops can mark certain testimonies "do not feature" (rare; e.g. privacy concerns). **Media moderation** — photos/videos pass same approval; extra check for recognizable third parties (consent gated).

Reversibility: moderation model and SLA are config. Switching between pre/post approval later is a workflow change — simpler if data model already tracks all approval states.

## Answer

**Decided:** **Option C** — **internal visibility immediate**, **public requires approval**; **≤3★ triggers CS ticket** before publish decision; **24h mod SLA**; **media same rules**; **removed content retained 12mo**; staff posts with `STAFF` badge.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
