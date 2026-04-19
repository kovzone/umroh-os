---
id: Q016
title: Booking on behalf of a minor (no KTP) — required docs, guardian linkage, payer distinction
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F4, F3
status: answered
---

# Q016 — Booking for minors

## Context

The PRD's booking flow assumes an adult pilgrim with a KTP (Indonesian national ID) and passport. Children under 17 don't have a KTP. But families routinely take minors on Umroh — especially babies, school-age children on holiday pilgrimages, and elderly grandparents carrying their grandchildren for a first-time pilgrimage experience.

The minimum-docs rule (Q006) says KTP + passport per jamaah before submit. For minors, KTP doesn't exist. The questions are: what substitutes for KTP? How is guardianship modelled? Who pays and who receives the receipt?

## The question

1. **Minor ID substitute.** Kartu Keluarga (KK) shows the child's name and parentage; is that the accepted substitute for KTP? Or do we need the Birth Certificate (Akta Kelahiran)?
2. **Passport requirement.** Minors need passports too (Saudi immigration requires this). Is the passport 6-month rule (Q006 / F3) applied identically to minors, or is there a different rule?
3. **Guardian linkage.** In the data model, a minor's `booking_item` must reference a guardian (`guardian_jamaah_id`). Should the guardian be another jamaah on the same booking (default), or can they be a jamaah on a different booking, or can they be a non-pilgrim (someone authorising the trip but not travelling)?
4. **Payer vs jamaah.** If the grandfather pays for the grandchild who's travelling with the mother, who receives the digital receipt? Does the booking have a `payer_id` distinct from the jamaah_ids?
5. **Age threshold.** Below what age is a jamaah treated as a minor for this purpose? Indonesian law: 18. Saudi visa rules may differ. What does the agency use as its operational threshold?
6. **Mahram implications.** Does a female minor traveller trigger the same mahram rules as a female adult (Q005), or different?

## Options considered

- **Option A — conservative: KK + Birth Certificate + same-booking guardian required.** Under-18 jamaah require BOTH KK and Akta Kelahiran as ID; guardian MUST be another jamaah on the same booking; passport-holder, 6-month rule applied identically; mahram rules for female minors differ per Q005 outcome.
  - Pros: legally defensible; Saudi-visa-friendly; matches conservative agency practice.
  - Cons: extra document upload friction; rules out scenarios where a child travels with a non-parent relative on a different booking (e.g. grandmother's booking in a different batch).
- **Option B — flexible: either KK or Akta acceptable; guardian can be on same OR a different active booking; payer can be non-traveller.** Covers more real-world scenarios.
  - Pros: operational realism; handles the "grandfather pays, mother travels with child" case explicitly.
  - Cons: rules complexity; more UX surface.
- **Option C — no minor-specific flow; treat all jamaah identically, reject bookings that fail KTP check.** Forces minors into their guardian's booking as "dependents," not as separate jamaah.
  - Pros: simplest model.
  - Cons: misrepresents the actual jamaah (the child IS a pilgrim, just a minor); complicates manifest accuracy, room assignment, and Saudi visa filings which DO recognise minor pilgrims as individuals.

## Recommendation

**Option A with one adjustment** — require KK (Kartu Keluarga); Akta Kelahiran optional / supplementary; guardian must be another jamaah on the same booking; mahram rules follow Q005's outcome (default is: same rules apply to female minors once they're old enough to be considered mature, typically 12+).

Treat minors as first-class `jamaah` records with a `date_of_birth` field and an `is_minor` computed flag. `booking_item` adds a `guardian_jamaah_id` nullable FK pointing to another `booking_item` in the same booking. Payer is distinct: a booking's `payer_user_id` is the IAM user who owns the billing (can be one of the jamaah, or a CS-proxied payer). Receipt goes to `payer_user_id`.

Age threshold: **under 18** is "minor" for guardianship purposes. Children under 2 (infants) follow additional airline lap-ticket rules captured on the booking but out of scope for this question.

Mahram: for female minors ≥ 12, apply Q005's rules identically. Below 12, no mahram requirement — minor is always under parental authority by default. _(Inferred — confirm with Q005 stakeholder input.)_

Reasoning: KK is almost universally held by Indonesian families and captures the parent-child relationship the visa desk will ultimately care about; requiring the guardian on the same booking matches actual Saudi visa mahram expectations and avoids operational complexity of cross-booking guardianship; separating `payer_user_id` from jamaah cleanly handles the "grandfather pays" case without conflating it with travel.

Reversibility: data-model additions (`guardian_jamaah_id`, `payer_user_id`, `date_of_birth`) are additive, safe to add. Loosening to cross-booking guardianship later is additive. Medium commitment on the document rule (changing to Akta-required after launch would force re-uploads).

## Answer

**Decided:** **KK required**, **Akta Kelahiran optional** (upload encouraged if name/KK edge cases); **guardian must be another `booking_item` on same booking**; **`payer_user_id` distinct** from jamaah rows (receipts to payer). **Minor = < 18** for guardianship; **passport 6-month rule same as adults**. **Mahram:** align **Q005** — females **12+** use policy engine; **<12** exempt from separate mahram requirement when travelling with declared guardian on same booking.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
