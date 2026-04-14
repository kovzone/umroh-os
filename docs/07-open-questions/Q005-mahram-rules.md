---
id: Q005
title: Mahram qualifying relation set, age threshold, and same-departure rule
asked_by: session 2026-04-14 F3 draft
asked_date: 2026-04-14
blocks: F3, F4, F6
status: open
---

# Q005 — Mahram qualifying relation set, age threshold, and same-departure rule

## Context

This is the single biggest religious/legal unknown in the entire system. The mahram algorithm decides whether a female pilgrim is allowed to travel in a given booking group. Getting it wrong either (a) rejects a valid traveller and loses a sale, or (b) accepts a traveller who will be turned away at Saudi immigration — costing both the visa fee and the agency's reputation.

The PRD says exactly one thing about mahram (L1617, paraphrased): **"For a female jamaah whose age is under 45, when her documents are marked VERIFIED, the system checks that `mahram_id` is not null. If null, warn ops."**

Everything else is silent. Fase 6.1 and 6.7 (L1185, L1197) mention "mahram verification" and "Pohon Keluarga / Relasi Nasab" as part of visa submission but refer the rules to "Aturan Arab Saudi" without enumerating them.

The Saudi mahram rules themselves have shifted. Nusuk (the Saudi Ministry of Hajj platform) has, since 2024, permitted women 18+ to travel for Umrah without a mahram under certain visa classes, subject to agency sponsorship. Our rule set needs to be **policy-configurable**, not hardcoded, because it will change again.

This question blocks F3 (the algorithm itself), F4 (booking submit validation), and F6 (visa submission to MOFA/Sajil).

## The question

Four linked sub-questions:

1. **Qualifying relation set.** Which relations in our `mahram_relations` enum count as a valid mahram for a travelling woman? The standard Islamic jurisprudence set is: husband, father, son, paternal grandfather, maternal grandfather, grandson (son's son or daughter's son), brother, paternal uncle, maternal uncle, nephew (brother's or sister's son), father-in-law, son-in-law, stepfather (with consummated marriage to mother), stepson. Do we use the full set, a subset, or a more permissive set?
2. **Age threshold.** PRD says <45 triggers the mahram requirement. Is this still current policy for this agency? Should it be a config value (system-wide or per-package)?
3. **Same-departure rule.** Must the mahram be a booked jamaah on the *same departure* as the subject, or is "registered as mahram" (even if travelling separately) acceptable?
4. **Proof of relation requirements.** For each qualifying edge, which document proves it? (E.g. Buku Nikah proves husband/wife; Kartu Keluarga proves parent/child/sibling.) When is staff manual verification allowed vs. required?

## Options considered

- **Option A — Conservative classical jurisprudence.** Full traditional set, age threshold 45, must travel in same departure, strict proof-of-relation per edge.
  - Pros: matches the most conservative reading; lowest risk of visa rejection.
  - Cons: rejects bookings that current Saudi rules would admit (e.g. women 18+ without mahram on eligible visa types); may frustrate customers.
- **Option B — Current Saudi Nusuk rules (2026).** Allow women 18+ without mahram under the agency-sponsored visa type; fall back to Option A for older visa types. Same-departure required when mahram IS the qualifying mechanism.
  - Pros: matches what Saudi actually accepts today; aligns with customer expectations.
  - Cons: rules keep changing; requires ongoing policy updates; ambiguous for non-Nusuk visa paths.
- **Option C — Policy-configurable, starts conservative.** Ship with Option A as default; expose the qualifying relation set, age threshold, same-departure flag, and proof-doc-per-relation map as system config that the agency can adjust without a code deploy.
  - Pros: safest default; adapts to changing Saudi rules without engineering; matches how other travel compliance systems work.
  - Cons: more config surface to build; still needs a religious/legal decision on the initial values.

## Recommendation

**Option C — policy-configurable, ship with Option A as the conservative default.** The Saudi rules have moved twice in three years and will move again; hardcoding any snapshot is a guaranteed bug. Model the mahram policy as a single config document:

```yaml
mahram_policy:
  enabled: true
  age_threshold_years: 45
  qualifying_relations:
    - husband
    - father
    - son
    - brother
    - grandfather
    - grandson
    - uncle
    - nephew
    - father_in_law
    - son_in_law
  require_same_departure: true
  proof_docs_required:
    husband: marriage_book
    wife: marriage_book
    father: family_book
    mother: family_book
    # ... etc
  allow_staff_override: true
  override_requires_reason: true
```

The religious/legal stakeholder picks the initial values; the agency can tune via the Admin console without engineering work. The algorithm itself is policy-agnostic — it walks the graph and checks against whatever the current config says. If Saudi changes the rules in 2027, someone updates the config and nobody deploys code.

Ship the default as Option A (conservative / classical). Document the defaults in this question's Answer section so the decision is auditable.

Reversibility: config-first means policy changes are near-zero-cost. The only migration risk is if a future rule requires a relation *type* not yet in the enum (e.g. a half-sibling distinction) — extending the enum is a routine schema change.

## Answer

TBD — awaiting religious/legal stakeholder input. Needs sign-off from:
- An Islamic scholar familiar with Saudi MOFA mahram rules
- The agency owner (business policy defaults)
- Likely the legal advisor (UU PDP intersects with family-data handling, see Q008)
