---
id: Q020
title: Manifest format per airline / immigration regulator
asked_by: session 2026-04-15 F7 draft
asked_date: 2026-04-15
blocks: F7, F6
status: answered
---

# Q020 — Manifest format per airline / regulator

## Context

PRD module #91 Manifest Imigrasi (line 313) says manifests are "Pencetakan manifes resmi keberangkatan kolektif **sesuai format standar**" — **sesuai format standar** (per standard format) is asserted but no specific standard is named. Who is the regulator (Kemenag, Imigrasi, Saudi, the airline) and what does their form actually require?

F7 W5 (manifest generation) needs this locked before we can implement the renderer. F6 also consumes manifest state for visa filings, so it's dual-blocked.

## The question

1. **Which regulator(s) does the manifest serve?**
   - Indonesian Imigrasi (pre-departure passenger manifest) — mandatory.
   - Indonesian Kemenag / SISKOPATUH — mandatory for PPIU-licensed Umroh operations.
   - Airline operational manifest — per-airline format.
   - Saudi immigration arrival manifest — if applicable.

2. **Per regulator, what's the exact layout?**
   - Columns required (Name, Passport #, DOB, Gender, Nationality, Visa #, ...).
   - Signatures / agency stamp placement.
   - Header with agency license details (PPIU number, NIB, etc.).
   - File format — PDF only, or Excel also?
   - Filename convention (e.g., `Manifest_<DepartureDate>_<AirlineCode>_<AgencyCode>.pdf`).

3. **Do airlines require their own format?** Garuda, Saudia, Lion Air, Emirates, etc. — each may expect a different layout for ground-handling.

4. **Generation cadence.** Is this a single generate-at-H-X document, or does it need to regenerate as pax list changes (up to H-24 lock)?

## Options considered

- **Option A — One canonical UmrohOS layout, exported to PDF + Excel.** Cover the Indonesian regulator requirements (columns + header + agency details); airlines accept this generic format as operational reference. Regenerate on demand until H-24 lock.
  - Pros: one renderer to maintain, fewer per-airline edge cases.
  - Cons: if a major airline rejects the generic format, we scramble.
- **Option B — Per-regulator + per-airline profiles.** Manifest generator takes a `profile` parameter (`imigrasi` | `kemenag` | `garuda` | `saudia` | `lion`) and renders the right layout. Ops picks per departure.
  - Pros: handles reality where each regulator/airline has preferences.
  - Cons: more templates to maintain; per-template ops input needed; slower to ship.
- **Option C — Start with A, add B-profiles as airlines complain.** MVP ships one canonical layout; add airline-specific profiles reactively.
  - Pros: pragmatic.
  - Cons: risk of an early airline rejection forcing urgent profile work mid-season.

## Recommendation

**Option C — canonical layout first, per-airline profiles reactively.** Build one rock-solid Imigrasi + Kemenag-compliant layout with a complete column set (name, passport #, DOB, gender, nationality, visa #, boarding-pass number, room number, bus number, muthawwif name, emergency contact). Export PDF + Excel both. If an airline later rejects or requests modifications, add a per-airline profile at that point.

Rationale: airlines in practice accept reasonably-formatted passenger manifests for operational use; the hard requirements are the Indonesian regulator ones (PPIU licensing, immigration pre-departure filing). Ship the regulator-correct format first; harden per-airline as real operational feedback lands.

Reversibility: adding a `profile` column to the `manifests` table later is additive. Template engine choice (HTML-to-PDF, same worker pool as F2 flyer / F5 receipts) is flexible across layout variants.

## Answer

**Decided:** **Option C** — one **canonical regulator-first** PDF + Excel template (Imigrasi + Kemenag/PPIU column superset + agency header fields); **regenerate on demand until H-24 lock** then versioned freeze. **Per-airline variants:** add **reactively** when a carrier rejects the generic layout.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
