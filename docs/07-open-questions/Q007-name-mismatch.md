---
id: Q007
title: KTP ↔ passport name mismatch handling policy
asked_by: session 2026-04-14 F3 draft
asked_date: 2026-04-14
blocks: F3, F6
status: answered
---

# Q007 — KTP ↔ passport name mismatch handling policy

## Context

Name mismatches between Indonesian KTP and passport are **endemic**. Common causes:
- Shortened name on KTP (`Muhammad Ali`) vs full name on passport (`Muhammad Ali bin Abdullah`)
- Arabic transliteration variants (`Ahmad` / `Achmad` / `Ahmed`)
- Diacritic drops in KTP fields that exist in the passport
- Name changes (marriage) reflected in one doc but not the other
- Typos in the original issuing record

Saudi MOFA and Indonesian immigration increasingly cross-check these fields. A mismatch caught at the visa stage (F6) means a rejected visa and a cascade of refund + rescheduling work. A mismatch caught at booking (F3 W3 verification) is cheap to resolve.

The policy question is: what do we do when the system detects a mismatch?

## The question

1. **Detection threshold.** At what Levenshtein distance (or other similarity metric) do we flag "mismatch"? Exact match only, or allow N character differences?
2. **Policy on detected mismatch:**
   - **Block verification** — staff cannot approve the passport until jamaah provides a supporting document (e.g. surat keterangan from Dukcapil or a legalised affidavit) explaining the difference.
   - **Allow verification with staff acknowledgement** — staff can approve by ticking an "I confirm the mismatch is explained" checkbox with a free-text reason, logged to audit.
   - **Auto-allow, log only** — mismatch detected and logged but no UX friction; visa desk handles it later.
3. **Which name flows to MOFA/Sajil?** When both are on file and differ, does F6 send the KTP name, the passport name, or both?

## Options considered

- **Option A — Strict block.** Mismatch halts verification; require supporting doc before proceeding.
  - Pros: catches everything at booking; no visa rejections for this cause.
  - Cons: rejects ~10–20% of Indonesian pilgrims at the gate for issues that are in practice accepted by Saudi authorities; operational nightmare.
- **Option B — Staff-acknowledged with reason.** Staff can approve with a mandatory reason code (`shortened_name`, `transliteration_variant`, `recent_name_change`, `other`); logged to audit. Downstream visa handling uses the passport name by default.
  - Pros: pragmatic; keeps the control without over-blocking; audit trail for later pattern analysis.
  - Cons: relies on staff judgment; inconsistent across staff members possible.
- **Option C — Auto-allow, log only.** System detects, logs, but does not gate. Visa desk handles escalations.
  - Pros: zero friction.
  - Cons: no control surface; mismatches reach visa stage where the cost is highest.

## Recommendation

**Option B — staff-acknowledged with structured reason.** The passport name is the name that Saudi MOFA sees; we send passport name always (overrides KTP at the visa submission payload level). Staff can approve a verified-with-mismatch passport by picking a reason code from a fixed enum:

- `shortened_name` — KTP has abbreviated form; passport is the full legal name
- `transliteration_variant` — Arabic-origin name spelled differently
- `recent_name_change` — marriage or legal name change
- `typo_in_ktp` — KTP issuing error (rare but real)
- `typo_in_passport` — passport issuing error (very rare)
- `other` — with mandatory free-text explanation

Detection threshold: Levenshtein distance ≤ 2 on any name token → "close match, not flagged". ≥ 3 OR different token count → "flagged, requires staff acknowledgement". _(Inferred threshold; tune after first month of real data.)_

The reason-code bucket is important — after three months we'll have enough data to know which mismatches almost always approve cleanly and can be relaxed further, which almost always fail at visa and should escalate earlier.

F6 sends the passport name to MOFA/Sajil; KTP name is stored but not transmitted unless the Saudi API has a specific "also-known-as" field (it doesn't today).

Reversibility: policy can move toward Option A (stricter) or Option C (looser) with a config flip on the booking/visa validation layer. The reason-code data is cumulative and makes future decisions evidence-driven.

## Answer

**Decided:** **Option B** — staff may approve mismatch with **mandatory structured reason codes** + audit; **passport name** is the canonical name sent to MOFA/Sajil payloads. **Detection:** tokenised normalisation + Levenshtein: **≤2** edits on longest token → *soft match* (suggest review); **≥3** or token-count mismatch → *hard flag* requiring acknowledgement. Thresholds remain **tunable from first month of production data**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
