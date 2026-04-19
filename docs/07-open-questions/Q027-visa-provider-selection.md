---
id: Q027
title: Visa provider selection per package kind
asked_by: session 2026-04-15 F6 draft
asked_date: 2026-04-15
blocks: F6
status: answered
---

# Q027 — Provider selection per package kind

## Context

PRD line 333 names "MOFA" and "Sajil" as sequential stages of the visa pipeline, but in practice these are **distinct routes**, not stages:
- **Sajil** — newer platform for Umroh visas (most common path for Umroh Reguler / Plus).
- **MOFA / direct embassy** — used for Hajj visas, Furoda (special quota), Badal, diplomatic edge cases.
- Future: Nusuk-issued visa for self-organized pilgrims (not our primary customer base).

The system's `visa_applications.provider_id` and the `provider_visa_id` column in `Visa_Applications` (PRD L1541) let us pick per-application, but the question is the **decision rule**: who picks, and based on what?

## The question

1. **Is provider selection deterministic by package kind?**
   - `umroh_reguler`, `umroh_plus` → always Sajil?
   - `hajj_furoda`, `hajj_khusus` → always MOFA direct?
   - `badal` → which path?

2. **Or is it per-departure / per-batch**, where ops picks the provider based on current operational relationships and provider availability?

3. **Can one booking use multiple providers** (e.g., some jamaah on Sajil, some on MOFA for the same departure)? Pattern: "the family ahead got Sajil quota approvals; the latecomers go through MOFA."

4. **If Phase 1 is manual (per Q026), does provider still need to be captured in the data model**, or is it free-text notes until Phase 2?

## Options considered

- **Option A — Deterministic per package kind.** Hard-coded mapping in config. Ops doesn't pick per-batch.
  - Pros: simplest; matches typical agency workflow; no UI for selection.
  - Cons: doesn't handle provider-availability shifts or special-case bookings.
- **Option B — Ops picks per-batch.** Each bulk-submit action includes a provider selection.
  - Pros: matches operational reality where agencies juggle provider relationships.
  - Cons: a choice surface for every batch; easy to pick wrong.
- **Option C — Default by package kind + override per batch.** Config maps package kind → default provider; ops UI shows the default but allows override with reason.
  - Pros: captures both patterns; safety via default.
  - Cons: slight UI complexity.

## Recommendation

**Option C — config default per package kind + per-batch override with audit.**

Default config (agency-editable via Super Admin):
- `umroh_reguler`, `umroh_plus` → Sajil
- `hajj_furoda`, `hajj_khusus` → MOFA
- `badal` → Sajil (_(Inferred)_ since Badal is typically Umroh-framed; needs confirmation)

Per-batch override: ops picks a different provider with a required reason text. Logged to audit + `provider_submissions` batch record.

Multi-provider per booking: supported via the data model (each `visa_applications` row has its own `provider_id`), but operationally discouraged — ops console surfaces a warning if split-submission within a booking is detected.

For Phase 1 (manual, per Q026): the provider field is still captured on every submission (manual entry by ops); it populates the `provider_id` column so the audit trail is complete even when the automated API isn't live yet.

Reasoning: deterministic defaults cover the 90% case; override is surgical for the remaining 10%; multi-provider-per-booking is allowed because the real world has edge cases, but flagged because it's usually a mistake.

Reversibility: config-driven defaults are trivially editable; override path is a small UI tweak if later deprecated.

## Answer

**Decided:** **Option C** — **defaults:** `umroh_reguler`/`umroh_plus`→**Sajil**, `hajj_furoda`/`hajj_khusus`→**MOFA**, **`badal`→Sajil default** with **ops override** to MOFA when case requires; **per-batch override + reason + audit**; **split-provider within one booking** allowed in data model but **UI warns**. **Phase 1 manual:** still capture `provider_id` on each batch for audit.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
