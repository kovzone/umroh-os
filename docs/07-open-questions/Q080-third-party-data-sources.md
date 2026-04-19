---
id: Q080
title: Third-party data source selection (prayer times API, qibla, Quran corpus)
asked_by: session 2026-04-17 F12 draft
asked_date: 2026-04-17
blocks: F12
status: answered
---

# Q080 — Third-party data source selection

## Context

F12 Daily Worship modules rely on computed or external-source data:

- **Prayer times** (module #190) — depends on chosen calculation method (MWL, ISNA, Egyptian, Umm Al-Qura, Kemenag, Institute of Geophysics Tehran, etc.); each produces different values. Many APIs exist (Aladhan, IslamicFinder, Kemenag, etc.) but the choice of method + source matters.
- **Qibla** (module #191) — coordinates of Kaaba are fixed (21.4225°N, 39.8262°E) but the bearing computation depends on whether we use great-circle vs rhumb line; cartographic convention matters.
- **Quran corpus** (module #192) — already addressed in Q079 from licensing angle; Q080 is the technical API / data-source angle.

Cross-cuts Q079 (licensing) but is distinct: Q079 is legal/content, Q080 is technical/API.

## The question

1. **Prayer times calculation method** — which school of calculation (MWL, Kemenag Indonesia, Umm Al-Qura, Diyanet, others)?
2. **Computation locus** — client-side (JS library like adhan-js) or server-side (UmrohOS computes + caches) or third-party API?
3. **Location-based method switching** — when jamaah is in Saudi, auto-switch to Umm Al-Qura?
4. **Qibla algorithm** — great-circle bearing (standard) or rhumb line?
5. **Quran API** — Quran.com API, alquran.cloud, self-hosted from Tanzil corpus?
6. **Failover** — if chosen API is down, fallback strategy?
7. **Caching** — local computation + caching or per-request API calls?
8. **Accuracy disclaimers** — prayer times have regional + seasonal nuances that calculation methods can't fully cover; how does UI communicate?

## Options considered

- **Option A — Client-side computation via adhan-js (npm library); method = Kemenag Indonesia (default) with auto-switch to Umm Al-Qura in Saudi; qibla = great-circle bearing computed client-side; Quran from self-hosted Tanzil corpus.** Offline-first; no API dependency.
  - Pros: works offline; no rate-limit concerns; fast; consistent.
  - Cons: libraries and corpora need maintenance; bugs in offline computation are embedded until deployed fix.
- **Option B — Third-party APIs for everything (Aladhan for prayer times, Quran.com API for Quran).** Always up-to-date from source.
  - Pros: no local computation bugs; external provider handles edge cases.
  - Cons: API dependency; offline broken; rate limits; latency.
- **Option C — Client-side computation with server-side validation fallback.** Compute locally; occasionally verify against server recomputation.
  - Pros: offline + correctness cross-check.
  - Cons: complexity.

## Recommendation

**Option A — client-side computation with Kemenag-Indonesia default method (auto-switch to Umm Al-Qura in Saudi-geo-fence); qibla via great-circle client-side; Quran self-hosted from Tanzil corpus per Q079 licensing.**

Option B's API-only approach breaks the offline-first promise of Q076 PWA design. Option C's server-side-validation-fallback is over-engineered for a Daily Worship utility where stakes are low (a 1-minute prayer-time difference between methods isn't catastrophic). Option A ships the standard pattern: adhan-js (well-maintained, widely-used) for prayer times + simple trig for qibla + local Quran corpus. Aladhan API as a reconciliation tool for debugging, not primary.

Defaults to propose: **Prayer times library** = `adhan-js` (MIT-licensed, used by major Islamic apps). **Default method** = Kemenag Indonesia (Indonesian jamaah at home). **Auto-switch**: when user geolocation within 500km of Mecca, method flips to Umm Al-Qura (Saudi standard); transition transparent to user with small "Metode: Umm Al-Qura" badge. **Qibla algorithm** = great-circle bearing (spherical-earth approximation — accurate to < 1° for all reasonable locations; sufficient for hand-held compass display). **Quran content** = self-hosted Tanzil corpus (text) + licensed Kemenag translation + Qurancdn murottal per Q079. **Caching**: prayer times pre-computed for 7 days ahead on first app open; qibla computed on-demand. **Failover**: Aladhan API (free, rate-limited) as fallback if local computation library load fails. **UI accuracy disclaimer**: small info-icon on prayer times screen with "Metode perhitungan: [Method]; toleransi ±3 menit untuk lokasi pegunungan / gedung tinggi."

Reversibility: switching calculation method is config; switching to API-primary is a code path swap. Adding new methods as user-selectable is UX additive.

## Answer

**Decided:** **Option A** — **`adhan-js` client-side**, **Kemenag method domestic**, **auto Umm al-Qura within 500km of Makkah**, **great-circle qibla**, **7d prayer cache**, **Aladhan fallback**, **disclaimer tooltip**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
