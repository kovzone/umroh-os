---
id: Q077
title: Community moderation policy + posting authority
asked_by: session 2026-04-17 F12 draft
asked_date: 2026-04-17
blocks: F12
status: answered
---

# Q077 — Community moderation policy

## Context

F12 W8 (forum grup angkatan) and W9 (reuni board) are the two surfaces where alumni + agency post user-generated content. PRD does not specify:

- Who can post (alumni only, alumni + guests, staff only).
- How moderation works (pre vs post).
- Escalation when content is reported.
- Retention policy for moderated-out content.
- Liability posture (UU ITE) for agency hosting UGC.

Too-loose moderation creates UU ITE liability + reputational risk. Too-strict (pre-moderation on every post) destroys forum velocity and alumni engagement.

## The question

1. **Posting authority per surface:**
   - Forum (#197): alumni of that cohort only? Open to all alumni?
   - Reuni board (#198): admin-only or alumni-initiated?
2. **Moderation model**: pre-moderation (publish after review) or post-moderation (publish immediately, review on flag)?
3. **Moderator role**: who (CS supervisor, dedicated community manager, marketing admin)?
4. **Flagging threshold**: how many user flags auto-hide?
5. **Appeal / reinstate**: can moderated-out posts be appealed?
6. **Content categories prohibited** (standard social-media TOS areas: hate, harassment, spam, commercial self-promotion, off-topic religious debate, etc.).
7. **Retention of removed content** for audit.
8. **Staff posting authority** — can ustadz / staff post in forums visible as special role?

## Options considered

- **Option A — Pre-moderation for all forum + reuni posts.** Moderator approves before public display.
  - Pros: zero-exposure for inappropriate content; tightest UU ITE defense.
  - Cons: destroys forum velocity; alumni post-then-wait experience is frustrating; moderator backlog risk.
- **Option B — Post-moderation (publish immediately; flag + review).** Community flags → auto-hide at N flags → moderator reviews.
  - Pros: real-time feel; scales with community; matches social-media norms.
  - Cons: brief exposure window to inappropriate content.
- **Option C — Hybrid: post-moderation default + pre-moderation for flagged users + reuni board admin-only.** Layered approach.
  - Pros: velocity preserved for good actors; tight control where warranted.
  - Cons: multi-state moderation system.

## Recommendation

**Option C — post-moderation default + auto-hide at 3 flags + pre-moderation gate for users with repeated violations + reuni board admin-only posting.**

Option A kills forum velocity — alumni expect Instagram / WA speed for their cohort chat, not email-moderation-queue delays. Option B is the standard pattern but offers no defense against repeat offenders (same user posts inappropriate content repeatedly). Option C starts fast for everyone, tightens on demonstrated bad actors, and keeps reuni board (which is event announcements) admin-curated so it doesn't become a general chatter surface.

Defaults to propose:

**Forum (#197):** alumni of the cohort only (gated by `package_departure_id` membership). Post-moderation default. Auto-hide at 3 distinct user flags. Moderator reviews within 24h SLA; approved → re-publish; rejected → `moderation_status=rejected` + content preserved for audit. Users with 3 rejected posts in 30 days → flagged, future posts pre-moderated for 90 days.

**Reuni board (#198):** admin-only posting (agency CS supervisor + marketing admin roles). Alumni consume + RSVP only.

**Moderator role:** `COMMUNITY_MOD` role in F1 — held by CS supervisor + marketing admin in MVP; dedicated community manager role can be added later.

**Prohibited content** (policy doc):
- Hate speech, harassment, discriminatory content
- Commercial self-promotion / competitor promotion (agency's own agent-to-jamaah sales is allowed only on agent-replica sites, not forum)
- Political campaigning
- Off-topic religious debate (directed to fatwa desk instead)
- Spam / repetitive content
- Personal information of third parties (phone, address without consent)

**Retention of removed content**: 12 months preservation; permanent deletion with audit log thereafter.

**Staff posting**: staff can post in forums with `STAFF` badge visibly distinguishing; not pre-moderated.

Reversibility: moderation mode (pre / post / hybrid) is config per surface. Tightening via emergency mode (all-pre) can be flipped fast in incident response.

## Answer

**Decided:** **Option C** — cohort-gated forum, **post-publish moderation** + **3-flag auto-hide**, **repeat offender pre-mod**, **reuni admin-only posts**, `COMMUNITY_MOD` role, **12mo removed retention**, **emergency all-pre mode** toggle.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
