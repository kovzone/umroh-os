---
id: Q060
title: WhatsApp broadcast rate limits + quality-score handling
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10, F7
status: open
---

# Q060 — WhatsApp broadcast rate limits + quality-score handling

## Context

F10's broadcast hub (module #59) and drip campaigns (#56) rely entirely on WhatsApp Business API (Meta). Meta imposes:

- **Messaging tier throughput limits** (1K / 10K / 100K unique recipients per 24h, per tier).
- **Template message approval** requirement (templates must be pre-approved by Meta).
- **Quality rating** that throttles or blocks senders with high report-rate / low read-rate.
- **24-hour messaging window** rule (only templated messages outside a customer-initiated 24h window).

F10 also has internal concerns: opt-out enforcement, dead-letter queue, duplicate suppression, tier-downgrade handling.

PRD doesn't specify any of this — leaves the ops layer to us.

## The question

1. **Starting messaging tier** — which Meta tier does the agency sign up at? Most start at 1K (Tier 1).
2. **Tier upgrade cadence** — quality-tier upgrades happen at Meta's discretion; do we track quality internally to anticipate them?
3. **Throughput cap strategy** — if a broadcast exceeds tier cap, queue to next day or reject?
4. **Quality-score threshold for auto-throttle** — below what quality score does UmrohOS reduce broadcast volume?
5. **Opt-out enforcement** — centralized opt-out list? One-list-per-agent? How do users opt out (reply STOP, portal toggle)?
6. **Template approval workflow** — who approves templates for Meta submission? How are rejections surfaced?
7. **Dead-letter** — failed messages after N retries — where do they go? Alert admin? Suppress recipient?
8. **Message-cost attribution** — Meta charges per conversation; does F9 journal these costs to marketing?
9. **Agent-initiated vs broadcast vs transactional** — three Meta message categories; which templates fall where?

## Options considered

- **Option A — Strict rate-limit adherence with queue + automatic throttle on quality-score drop.** Enforce tier cap per 24h; queue overflow to next day; monitor quality score daily; auto-throttle broadcasts when score drops below good rating.
  - Pros: avoids Meta penalties; automatic recovery.
  - Cons: queue can grow large during campaign launches; tier-cap feels arbitrary to admin.
- **Option B — Manual admin controls; admin handles tier + quality.** No auto-throttle; admin reads Meta dashboard; manually pauses campaigns on quality issues.
  - Pros: simple; no internal quality-score model.
  - Cons: relies on admin discipline; late discovery of quality drops.
- **Option C — Hybrid: auto-enforce rate cap + opt-out; admin alerts on quality; tier upgrade via admin request to Meta.** Middle ground.
  - Pros: automates the basics; keeps strategic choices with admin.
  - Cons: still some admin burden.

## Recommendation

**Option C — hybrid with auto rate-cap enforcement + centralized opt-out + admin alerts on quality drops; template approval via marketing-admin role.**

Option A's full-auto is good defensively but over-engineers MVP; quality-score tracking is a whole sub-system best deferred. Option B's full-manual is brittle — a single campaign blast can tank quality and block broadcasts for a week. Option C enforces the mechanical constraints (rate cap, opt-out, template approval) while leaving the judgment calls (tier upgrade timing, paused campaign resumption) to admin with alerts.

Defaults to propose: **starting tier** = Tier 1 (1K unique recipients / 24h) per Meta default for new business accounts. **Tier upgrade** = admin triggers by Meta-dashboard request; UmrohOS internal tracks `current_tier_cap` as config; operator updates when Meta upgrades. **Throughput cap** = enforce 80% of cap as internal soft limit to leave headroom; queue above soft limit to next day with ETA surfaced to admin. **Quality-score** = daily pull from Meta if API supports (Meta's WhatsApp Business Management API exposes quality ratings); alert admin on "Medium" or "Low" ratings; recommended (non-enforced) pause. **Opt-out** = centralized per-tenant list; users opt out via reply "STOP" / "BERHENTI" / "UNSUB" (keywords configurable); portal toggle for WA preferences; opt-out list filters **every** broadcast and drip. **Template approval** = marketing admin submits new template via Meta Business Manager (manual step for MVP); template registered in UmrohOS with `meta_template_id`; status checked via API pull. **Dead-letter** = 3 retries with exponential backoff; after 3rd fail, recipient phone flagged `wa_unreachable`; suppressed from future broadcasts until manually cleared; alert admin on mass failures. **Cost attribution** = daily Meta billing pull joined to broadcast campaigns; F9 W8 reconciles as marketing expense under `Beban Komunikasi WA`. **Message categories** = MVP assumes all broadcast templates are "Marketing" category (highest-priced but most flexible); "Utility" templates (transactional — booking confirmations, VA issued) separately tagged for lower-cost handling.

Reversibility: throughput cap % (80%), retry count (3), opt-out keywords — all config.

## Answer

TBD — awaiting stakeholder input. Deciders: marketing admin, CRM lead, Meta Business Manager admin (if separate), agency owner (cost posture).
