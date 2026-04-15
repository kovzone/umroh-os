---
id: Q031
title: Tasreh issuance authority — who creates the document?
asked_by: session 2026-04-15 F6 draft
asked_date: 2026-04-15
blocks: F6
status: open
---

# Q031 — Tasreh issuance authority

## Context

Tasreh = entry permit for restricted Saudi sites, most commonly for **Raudhah** (the area between the Prophet's minbar and his tomb in Masjid Nabawi) but also for other capped-entry locations. PRD line 359 describes a "digital archive of surat izin (Tasreh)" — the agency stores the tasreh and jamaah can present it at entry (PRD L551).

What's unclear from the PRD is **who creates/issues the tasreh in the first place**. Three plausible sources:

1. **Nusuk (Saudi government platform) downloads.** Since ~2023, Saudi MoHRD has moved tasreh issuance onto Nusuk; travelers book appointments and download their permits. This may require the agency or the jamaah to log into Nusuk.
2. **Agency-generated.** Some agencies generate their own tasreh-like documents internally for agency operational use (these aren't legally binding Saudi permits; they're more like agency-side manifests).
3. **Provider-supplied.** Third-party visa providers (Sajil, TafweejIT) sometimes bundle tasreh issuance into their service package.

The answer affects the data model (who owns the upload action?), the security model (Nusuk credentials = big deal), and the jamaah-facing experience (do they download from our portal or from Nusuk's app?).

## The question

1. **What's the source of tasreh in this agency's current workflow?**
   - Ops manually downloads from Nusuk and uploads to UmrohOS
   - Ops generates an agency PDF
   - Provider delivers as part of their visa package
   - Mixed (different for different package kinds)

2. **Who can upload a tasreh to the UmrohOS system?**
   - Ops admin only
   - Ops admin + muthawwif (for field corrections)
   - Jamaah themselves (if they downloaded from Nusuk app)

3. **Does Nusuk-integration (future) replace the manual upload?** If we eventually integrate with Nusuk's API, tasreh would be auto-fetched. Design the MVP to accept manual upload + be extensible to auto-fetch.

4. **Per-departure vs per-jamaah tasreh.** Tasreh is per-jamaah per-entry-appointment. A single jamaah on a 2-week trip might have multiple tasreh (e.g., 2 Raudhah visits). How does the system track multiple tasreh per jamaah?

5. **Tasreh expiry.** Tasreh is typically valid for a specific date/time slot. The system should refuse to count an expired tasreh as "used" at F7 W11 scan.

6. **Can tasreh be transferred?** No — each tasreh is tied to a specific passport. This should be enforced as a soft warning if someone attempts to scan one tasreh for a different jamaah.

## Options considered

- **Option A — Manual upload by ops, MVP; Nusuk auto-fetch deferred.** Ops downloads tasreh from Nusuk (agency-side) and uploads to UmrohOS with a `source: nusuk_download` tag. Jamaah sees the tasreh in their portal + muthawwif app.
  - Pros: doesn't require Nusuk API access; matches current Indonesian agency practice; simple data model.
  - Cons: ops-side upload friction; potential for upload delays.
- **Option B — Jamaah self-upload.** Jamaah downloads tasreh from their personal Nusuk app and uploads to UmrohOS.
  - Pros: scales without ops burden.
  - Cons: many jamaah are elderly or non-tech-savvy; would require training; risk of upload delays closer to appointment time.
- **Option C — Agency generates own tasreh-like document.** Not legally valid for Saudi entry; purely agency-side operational reference. Easy to implement but low value (won't help at the Raudhah gate).
- **Option D — Hybrid: ops uploads by default + jamaah can upload too + future Nusuk API integration auto-fetches.**
  - Pros: flexibility; scales well; future-compatible.
  - Cons: multiple source-of-truth potential if both ops and jamaah upload different versions.

## Recommendation

**Option D — ops uploads by default; jamaah can upload optionally; Nusuk auto-fetch as Phase 2 feature flag.**

Schema:
- `tasreh_records { id, jamaah_id, kind (enum: raudhah | other), issued_for (timestamp of the scheduled entry), storage_path, source (enum: nusuk_download | agency_generated | provider_supplied | nusuk_api_auto), uploaded_by (user_id), uploaded_at, valid_until (expiry timestamp), used_at nullable }`
- Multiple tasreh per jamaah is just multiple rows.
- Expiry check at F7 W11 scan: `valid_until > now()`; expired tasreh returns an explicit "expired" state, not a generic scan failure.
- Transfer prevention: at scan time, verify `tasreh.jamaah_id` matches the scanned jamaah; mismatch returns a warning and muthawwif override is required per Q022 authority model.

Upload sources:
- **Ops admin**: primary path; uploads on behalf of jamaah.
- **Jamaah self-upload**: optional, from the customer portal; useful for edge cases (jamaah already downloaded from Nusuk via their phone).
- **Conflict resolution** (both ops and jamaah uploaded different versions): most recent upload wins; older versions remain in the audit trail; ops can manually reconcile.

Phase 2 (Nusuk API integration, per Q026 access path): when API is available, add `NusukApiAdapter` that auto-fetches tasreh per jamaah at scheduled intervals. `source: nusuk_api_auto` rows created automatically.

Reasoning: manual-upload MVP is operationally feasible today; jamaah self-upload covers the long tail of "my ops admin is slow" complaints; API auto-fetch is the scaling endpoint. Agency-generated tasreh (Option C) is off the table because it has no legal value at the Saudi entry point.

Reversibility: adding auto-fetch via adapter is additive; removing jamaah self-upload later is a UI change with grandfathering.

## Answer

TBD — awaiting stakeholder input. Deciders: agency ops lead (current workflow) + muthawwif representative (field reality of tasreh presentation) + tech lead if Nusuk API access is being negotiated.
