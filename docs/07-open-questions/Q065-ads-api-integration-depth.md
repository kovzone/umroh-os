---
id: Q065
title: Ads API integration depth (Meta / Google)
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: open
---

# Q065 — Ads API integration depth

## Context

Module #48 *Ads Manager Lite* says "pull cost / impression / click from Meta/Google via API" — but "via API" can mean several levels of integration:

1. **Pull-only** — daily ingest of spend/clicks/impressions per campaign for ROAS + budget-reconciliation.
2. **Pull + account linking** — agency connects Meta/Google Ads account via OAuth; UmrohOS reads granular data.
3. **Push creative** — UmrohOS uploads ad creative / copy to Meta/Google.
4. **Push campaign management** — UmrohOS creates/edits/pauses campaigns directly.
5. **Full-stack ads platform** — UmrohOS is the ads manager; native Meta/Google dashboards not needed.

Each level is a progressively bigger engineering commitment. Module #48 is tagged **Should Have / Medium**, so the right depth is not "the maximum."

## The question

1. **Integration depth** — which of the 5 levels above is the right MVP scope?
2. **OAuth flow** — who connects the Meta/Google account (agency owner, marketing admin)? Refresh token management?
3. **Pull cadence** — hourly? Daily? Real-time?
4. **Data granularity** — per-campaign? Per-ad-set? Per-ad?
5. **Reconciliation with finance** — pulled spend joins to F9's marketing expense journal; at what posting cadence?
6. **Multiple ads accounts** — if agency runs separate Meta accounts per branch, does UmrohOS aggregate?
7. **Retargeting (#69)** — requires push (custom audience upload); upgrade path from pull-only?
8. **Fallback** — if API fails, manual CSV upload from Meta/Google Ads dashboard?

## Options considered

- **Option A — Pull-only, daily cron, per-campaign granularity, single account.** Lightweight.
  - Pros: smallest engineering; covers ROAS + budget use cases; low maintenance.
  - Cons: no push; retargeting requires manual upload; granular optimization requires native dashboards.
- **Option B — Pull-only + push custom audience (retargeting upload).** Lightweight plus retargeting.
  - Pros: covers #48 + #69 core use cases; still lightweight.
  - Cons: custom audience API is a separate OAuth scope; incremental complexity.
- **Option C — Full-stack integration: pull + push creative + campaign management.** Heavy.
  - Pros: agency works inside UmrohOS only; one interface.
  - Cons: substantial engineering; maintaining parity with Meta/Google native features is a losing game.
- **Option D — Pull-only MVP; CSV upload for retargeting; upgrade to push in Phase 2.** Defer the harder part.
  - Pros: fastest to ship; retargeting works via manual process.
  - Cons: manual uploads are friction for marketing team.

## Recommendation

**Option D — pull-only integration (daily cron, per-campaign granularity, OAuth via marketing admin) for MVP; retargeting via CSV export + manual upload; upgrade to Option B (custom audience API push) in Phase 2.**

Option A is right for MVP scope but doesn't address retargeting at all. Option B solves retargeting elegantly but adds OAuth scope complexity + moderation of what gets pushed. Option C is a "replace Meta/Google Ads" scope creep. Option D ships pull-only fast and defers the push-side to Phase 2 with a clear manual-upload workaround in the interim — preserves the Should Have priority intent of #48 without over-committing.

Defaults to propose: **MVP integration** — pull-only. OAuth flow: marketing admin connects Meta + Google Ads via Ads API OAuth; refresh tokens managed by crm-svc adapter. **Pull cadence** — daily at 6am WIB (covers previous-day close); admin-triggered manual refresh available. **Granularity** — per-campaign in MVP; per-ad-set in Phase 2 (more granular ROAS). **Reconciliation** — daily ads costs join to F9 via campaign's UTM key → marketing expense journal (Dr Beban Promosi / Cr Hutang Meta Ads), posted T+1 batch. **Multiple ads accounts** — MVP single-account per platform; multi-account in Phase 2. **Retargeting (#69)** — MVP exports warm-lead list to CSV for marketing team to manually upload to Meta/Google; Phase 2 pushes via Custom Audiences API. **Fallback** — manual CSV upload of ads spend data accepted via admin UI; same journal path as auto-pull.

Reversibility: integration depth expansion is strictly additive (pull → pull+push → more push); nothing in Option D forecloses upgrades. OAuth scopes expand per feature, not retroactively.

## Answer

TBD — awaiting stakeholder input. Deciders: marketing admin (ad-ops workflow), agency owner (integration investment), engineering (scope constraint).
