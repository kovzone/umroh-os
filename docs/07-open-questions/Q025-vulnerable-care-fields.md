---
id: Q025
title: Vulnerable Care manifest fields — categories, source of flag, sensitivity handling
asked_by: session 2026-04-15 F7 draft
asked_date: 2026-04-15
blocks: F7, F3
status: open
---

# Q025 — Vulnerable Care fields

## Context

PRD module #94 Handling Manifest (line 323) describes a special manifest for airport handling + muthawwif flagging jamaah who need special care — lansia (elderly), wheelchair, etc. The benefit is real: ground staff can prioritize, muthawwif can plan special-care routing, family can be seated together on flights.

But "vulnerable" intersects with **medical privacy**, **UU PDP** (Q008), and cultural sensitivities. Operators need enough information to help; too much, and we're storing sensitive medical data with compliance implications we haven't thought through.

## The question

1. **What categories of "vulnerable" does the system track?** Candidate enum:
   - `lansia` (elderly — defined by age threshold, e.g., ≥ 70)
   - `wheelchair` / `mobility_impaired`
   - `pregnant`
   - `visually_impaired`
   - `hearing_impaired`
   - `dietary` (diabetic, celiac, etc.)
   - `medical` (ongoing condition requiring attention)
   - `mental_support_needed`
   - others?

2. **Source of the flag.**
   - Jamaah self-declares during F3 registration (check-boxes).
   - Agent or CS tags during booking (on behalf of the jamaah).
   - Medical certificate uploaded (for `medical` and `dietary`).

3. **Who can see the flag?** Ops admin, tour leader, muthawwif, airport handler — all of them? Or scoped (e.g., `dietary` visible only to ops; `medical` visible to tour leader but not all handling staff)?

4. **Sensitivity handling.** If these flags are personal health data (UU PDP intersects):
   - Are they consent-required on collection?
   - Retention policy post-departure (default is 2 years per Q008; do medical flags need faster purge)?
   - Can jamaah request erasure of these flags specifically?

5. **Operational expression.** What does each flag do concretely?
   - `wheelchair` → flight row priority, bus door-proximate seat, hotel ground-floor room.
   - `lansia` → similar treatment + medical-kit awareness.
   - `dietary` → meal-plan flag for hotel + airline.
   - `medical` → muthawwif briefed; nearby medical facility mapped.

## Options considered

- **Option A — Conservative minimal categories, self-declared + optional medical cert, ops-visible only for medical.** Start with 4 categories (lansia, wheelchair, pregnant, dietary); jamaah self-declares; medical cert uploaded for any flag with health implications; flag visible to ops, tour leader, muthawwif, airport handler (need-to-know basis for ops-visible non-medical).
  - Pros: small attack surface; manageable compliance; covers the most common cases.
  - Cons: may miss edge cases agency sees frequently.

- **Option B — Full category set, multi-party tagging, tiered visibility.** 8+ categories; jamaah + agent + ops can all tag; visibility tiered by role and category.
  - Pros: captures everything the agency might need.
  - Cons: UX complexity; tagging disputes (who's right?); compliance surface expands.

- **Option C — Just `requires_special_handling: boolean + notes: text`.** No categorical enum; free-text notes.
  - Pros: simplest.
  - Cons: unstructured data; can't aggregate for airline/hotel planning; invites free-text PHI leakage.

## Recommendation

**Option A — 4 core categories with clear operational meaning, self-declared at registration, medical cert upload where applicable.**

Schema:
- `jamaah.vulnerable_flags` — jsonb array of structured entries:
  ```
  [
    { "category": "wheelchair", "declared_by": "<user_id>", "declared_at": "...", "cert_document_id": null },
    { "category": "dietary", "notes": "diabetic", "declared_by": "...", "declared_at": "...", "cert_document_id": "<doc>" }
  ]
  ```
- Categories enum (v1): `lansia | wheelchair | pregnant | dietary`.
- Extension path: add `visually_impaired`, `hearing_impaired`, `medical`, `mental_support_needed` as agency demand emerges.

Sources:
- **Jamaah self-declares** at registration via checkbox (portal or CS-assisted).
- **CS / agent can tag on behalf** with jamaah consent noted in audit.
- **Medical certificate** required for `dietary: diabetic` and any future `medical` category; upload as a `documents` entry in F3 with kind = `medical_certificate`.

Visibility:
- **Ops admin**: sees all flags + certs.
- **Tour leader**: sees flags + category-level notes, but NOT the medical cert contents.
- **Muthawwif**: sees flags + category, no notes.
- **Airport handler**: sees flags only (for boarding + handling priority), no category detail beyond "needs assistance".

UU PDP (Q008) alignment:
- Declaring a vulnerable flag captures explicit consent at F3 registration (opt-in checkbox explaining that the data flows to handling partners for service delivery).
- Retention: same 2-year-post-departure purge as F3 documents. Alumni retention is jamaah-opt-in per Q008.
- DSR erasure requests: structured field removal is straightforward; medical certs are erased alongside.

Operational hooks:
- `wheelchair` → ops-svc auto-flags the jamaah in W4 Vulnerable Care manifest; airline manifest generator includes in priority-handling list; room allocation prefers ground-floor room where available.
- `lansia` → Vulnerable Care manifest flag; medical-kit awareness for muthawwif.
- `dietary` → meal plan flag surfaced to hotel + airline bookings (via catalog-svc package details on departure).
- `pregnant` → similar to wheelchair (priority treatment); system can optionally compute gestational age if declared-at + trimester info given.

Reasoning: 4 categories covers 90%+ of real cases in Umroh/Hajj operations; self-declaration is the cleanest UU PDP story; tiered visibility respects the compliance gradient (tour leader needs actionable info, airport handler needs minimum info); extension path keeps us flexible.

Reversibility: enum extensions and visibility rule tightening are code-level changes; the jsonb column tolerates additive categories without migration.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner + ops lead + legal (UU PDP compliance for health-adjacent data) + CS lead (real-case coverage).
