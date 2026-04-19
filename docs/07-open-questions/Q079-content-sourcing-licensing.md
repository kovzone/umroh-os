---
id: Q079
title: Manasik + Quran + daily-content sourcing & licensing
asked_by: session 2026-04-17 F12 draft
asked_date: 2026-04-17
blocks: F12
status: answered
---

# Q079 — Content sourcing & licensing

## Context

F12 serves three kinds of Islamic content with distinct sourcing + licensing concerns:

- **Quran text + translation + audio** (module #192) — Arabic text, Indonesian translation, murottal recitation. Licensing required per source.
- **Dzikir + doa** (module #193) — shorter texts; public-domain classical sources but specific editions may be copyrighted.
- **Manasik encyclopedia** (module #194) — step-by-step ritual guide; agency-authored or licensed from external (Kemenag, pesantren, specific ustadz curriculum).
- **Articles + kajian** (module #195) — agency-authored + external contributions.

Shipping unlicensed content creates copyright liability and potentially religious-authority issues (e.g. using someone's translation without permission erodes trust).

## The question

1. **Quran Arabic text source** — Tanzil.net (Creative Commons), King Fahd Complex corpus (licensed), or other?
2. **Indonesian translation** — Kemenag translation (commonly licensed for Indonesian projects), other?
3. **Audio murottal** — which reciters, which license source? (Qurancdn / Everyayah / direct licensing deals)
4. **Dzikir texts** — Hisnul Muslim, Al-Ma'thurat — licensed editions vs public-domain-derivatives?
5. **Manasik content ownership** — agency-authored (signed-off by religious advisor), licensed from Kemenag Manasik curriculum, or external religious publisher?
6. **Attribution display** — how do we credit sources visibly?
7. **Update cadence** — if Kemenag releases new Manasik guide, how do updates propagate?
8. **Translation rights** — if Indonesian translation is used, is it for display-only or also embedded in downloadable offline bundles?

## Options considered

- **Option A — All licensed content + clear attribution.** Tanzil Quran + Kemenag translation + Qurancdn murottal + licensed Hisnul Muslim + agency-authored manasik with religious advisor sign-off.
  - Pros: clean legal posture; high content quality; clear attribution chain.
  - Cons: licensing fees + maintenance; subject to license renewals.
- **Option B — Third-party API passthrough (Quran.com API etc.) instead of hosting content.** Use external APIs for Quran + dzikir; no local copy.
  - Pros: no licensing burden (API provider handles it); content always up-to-date.
  - Cons: offline capability broken (no local copy); API dependency for a core feature.
- **Option C — Hybrid: license Quran + translation locally (for offline) + third-party API for dzikir + agency-authored manasik.** Mix based on feature needs.
  - Pros: balance of offline capability + minimized licensing surface.
  - Cons: mixed model; more complex.

## Recommendation

**Option A — license everything + ship with local copies + document attribution chain + maintain license renewals in agency governance.**

Option B's API-only approach breaks offline usage (a core Daily App value proposition per Q076); a jamaah in Saudi with expensive roaming can't rely on an API fetch for Fajr prayer's Quran recitation. Option C is architecturally OK but introduces dependency on external APIs for frequently-used features (dzikir is used daily; API downtime on external provider kills feature). Option A owns the licensing cost in exchange for capability control.

Defaults to propose: **Quran Arabic text** = Tanzil.net under their Creative Commons license (attribution required). **Indonesian translation** = Kemenag translation via the Kemenag Quran API or licensed copy; direct license application from Kementerian Agama Republik Indonesia (free for non-commercial + agency-attributed use, per their standard terms). **Murottal audio** = Qurancdn.com (Creative Commons reciters) + Sheikh Mishary / Sheikh Sudais via direct licensing deals (agency cost: ~5M IDR/year estimated). **Dzikir / doa** = Hisnul Muslim public-domain Indonesian translation (careful of edition copyright — use pre-1970 translations where possible); agency's religious advisor validates authenticity. **Manasik** = agency-authored with religious advisor sign-off (Option A variant — owns the content); alternative = license from Kemenag Manasik Haji & Umrah curriculum (free for Umrah agencies with PPIU license). **Attribution** = visible on every content surface (small attribution footer); "License" page listing all sources.

Reversibility: switching Quran source later is content-replace (no schema change). Switching manasik source is mid-weight (content delete + re-import).

## Answer

**Decided:** **Option A** — **Tanzil + Kemenag translation + licensed audio** path; **agency-authored manasik** with advisor sign-off **or** licensed Kemenag pack per tenant choice; **attribution footer + `/licenses` page**; **renewal calendar owned by ops**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
