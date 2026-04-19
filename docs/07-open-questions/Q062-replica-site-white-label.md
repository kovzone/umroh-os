---
id: Q062
title: Replica-site white-label scope
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: answered
---

# Q062 — Replica-site white-label scope

## Context

Module #28 (Website Replika Agen) auto-provisions each agent a personal landing page at `domain.com/id/<agent-code>`. The PRD's description (lines 117, 723, 283) focuses on:
- Full catalog from central.
- WA button piped to agent's number.
- Auto-watermarked flyer (#32).
- Auto-update sync (Modul #283).

**What customization the agent has control over is not specified.** The range is wide:

- Zero customization — every replica site looks identical, only WA number changes.
- Light customization — agent name, photo, tagline.
- Medium customization — agent logo, brand color, banner image.
- Heavy customization — custom subdomain, custom theme, full white-label.

## The question

1. **Customization scope** — zero / light / medium / heavy?
2. **Custom domain support** — can agent use their own domain (e.g. `umroh.agenku.com` pointing to their replica)?
3. **Logo + branding** — can agent upload their own logo and replace UmrohOS branding?
4. **Color scheme** — can agent pick colors (from palette, or free)?
5. **Testimoni filtering** — can agent choose which testimonies to feature on their site?
6. **Package filtering** — can agent hide certain packages (e.g. don't show Paket Platinum, or only show Haji Furoda)?
7. **Content ownership** — if agent adds custom content (blog post, about-me), who owns it + moderates it?
8. **Breaking isolation** — if agent customization breaks site functionality, who supports them?

## Options considered

- **Option A — Zero customization.** Replica is identical template; only agent name + photo + WA number change.
  - Pros: simplest; most control over brand consistency; easiest to support.
  - Cons: agents want identity; "just another replica" limits differentiation.
- **Option B — Light customization: agent name + photo + tagline + featured testimoni selection.** Agent can curate, not design.
  - Pros: moderate identity; bounded engineering; moderation surface small.
  - Cons: still feels templated.
- **Option C — Medium customization: color scheme from palette + banner image upload + package hide/show.** Deeper personalization, constrained to safe dimensions.
  - Pros: meaningful personalization; still bounded; brand integrity preserved via palette.
  - Cons: moderation of banner uploads; palette maintenance.
- **Option D — Heavy customization / white-label: custom domain + logo + full theme.** Agent feels site is theirs.
  - Pros: strongest agent-identity feeling; supports super-agents running as mini-agencies.
  - Cons: substantial engineering (DNS, SSL, theme system); moderation costs; brand dilution risk.

## Recommendation

**Option B — light customization (agent name, photo, tagline, featured testimoni selection); MVP only; Option C as Phase 2; Option D (full white-label) deferred.**

Option A's zero-customization kills agent engagement — motivated Super-Agens want something they can call theirs. Option D is a product of its own (reseller platform); distracts from MVP scope and brings DNS + SSL + theme infrastructure costs that don't match the Could Have priority. Option B gives agents a sense of ownership (their face, their words) without adding theme-system complexity; Option C's package-hide is a reasonable Phase 2 add. White-label path (Option D) stays available for super-agents with bigger business, treated as a separate upsell.

Defaults to propose: **MVP customization** — agent name (from KTP), agent photo (uploaded at KYC, replaceable via portal), tagline (free-text 140 char, moderated pre-display on first save — Q063 reuse), featured testimoni selection (up to 3 approved testimonies to highlight, from central-moderated pool). **Custom domain** — not supported in MVP; all replica sites under `domain.com/id/<agent-code>`. **Logo / color scheme** — no change in MVP; central branding only. **Package filtering** — not in MVP; full catalog always shown. **Content ownership** — agent owns their tagline + testimoni selection; moderation by central staff on first save. **Breaking isolation** — moot in Option B; no agent-deployable code. **Phase 2 path** — Option C (palette-constrained color + banner + package hide). **Phase 3 path** — Option D (custom domain, full theme) as super-agent tier entitlement.

Reversibility: expanding customization later is purely additive — new config columns on `agents` table, new UI pages. Contracting later is the hard direction.

## Answer

**Decided:** **Option B MVP** — name/photo/tagline + **≤3 featured approved testimoni**; **no custom domain / palette / logo** until Phase 2 (**Option C**) and **Option D** later; **first-save moderation** on tagline (**Q063**).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
