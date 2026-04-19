---
id: Q006
title: Minimum documents required to submit a booking
asked_by: session 2026-04-14 F3 draft
asked_date: 2026-04-14
blocks: F3, F4
status: answered
---

# Q006 — Minimum documents required to submit a booking

## Context

The PRD gates document completeness at the **visa submission** stage (WAITING_DOCS → READY), not at booking submit. That means a jamaah can theoretically book a package without having uploaded anything. But in practice, agencies want *some* baseline — at minimum a KTP — before accepting money, because a jamaah with no uploaded docs is a weak commitment signal.

We need a policy before the F4 booking-submit saga can enforce it. The decision affects (a) the B2C self-booking flow (what gets blocked in the UI), (b) the booking-svc validation on submit, and (c) the commercial policy of "when can we collect a DP".

## The question

Which document(s) must be uploaded and verified before a booking can transition from `draft` to `pending_payment` (the state where a VA is issued)?

1. **Nothing** — DP collection is allowed with no documents uploaded. Verification happens later on the visa pipeline timeline.
2. **KTP only** — at minimum the jamaah has to prove they're a real Indonesian citizen before committing money.
3. **KTP + passport (scan uploaded, OCR not necessarily verified yet)** — prove identity and that a valid travel document exists. Verification of the passport MRZ can follow asynchronously.
4. **KTP + passport + the relevant mahram proof** (Kartu Keluarga for parent/sibling or Buku Nikah for spouse) when applicable — prove the whole booking group can travel together.

Secondary questions:
- If the minimum is not met, does the B2C portal **block** the submit button, or allow submit with a warning?
- If the booking is a multi-jamaah booking, is the minimum enforced **per jamaah** or **per booking** (i.e. one jamaah's docs satisfy the whole booking)?

## Options considered

Mapped to the levels above. The real tradeoff is commitment signal vs friction:

- **Level 1 (nothing)** — maximum conversion; highest risk of ghost bookings that never complete.
- **Level 2 (KTP)** — small friction, strong signal. Industry norm for Indonesian travel.
- **Level 3 (KTP + passport)** — more friction, stronger signal. Avoids the "booked without a valid passport" class of problem early.
- **Level 4 (full mahram-proof set)** — high friction; may lose customers who'd otherwise commit. Risk: the mahram proof doc requires a spouse/parent/sibling who isn't yet registered; chicken-and-egg.

## Recommendation

**Level 3 — KTP + passport scan uploaded (OCR can be pending).** Enforced per-jamaah for multi-jamaah bookings — each jamaah on the booking must have at least those two uploads. The B2C portal **blocks** the submit button with a helpful message; does not silently allow submit.

Rationale:
- KTP alone is too weak — agencies routinely report "we took a DP then discovered the customer has an expired passport and can't travel this season." One-time cost to fix it at booking is cheaper than a refund saga later (F5).
- Uploading is cheap friction; jamaah photos of their passport are sitting on their phone already. We don't require the OCR to be *verified* — just that a file exists in GCS with `kind = passport, status in (uploaded, processing, needs_review, verified)`. If it later rejects, the booking moves to a "documents need re-upload" state rather than cancelling.
- Mahram proof is *not* required at booking submit (Level 4) because it forces an order-of-operations cliff — jamaah often need to ask a relative to register before they can prove the relationship. Mahram validation happens at submit (F4) based on whatever relations are on file, and at visa submission (F6) based on fully verified relations. Booking is allowed to proceed with `needs_mahram: true, found_mahram: null` as a *flag*, not a block, as long as KTP + passport are up.

Reversibility: the minimum-doc rule is enforced by booking-svc at submit time. Tightening to Level 4 later is a one-line policy change; relaxing to Level 2 is the same. The B2C UI will need a minor copy update per change. Low commitment.

## Answer

**Decided:** **Level 3** — **KTP + passport scan per jamaah** before `draft` → `pending_payment`; OCR may still be `processing`. **B2C:** hard **block** on submit. **Multi-jamaah:** rule enforced **per jamaah**. **Mahram proof:** **not** required at booking submit (follows Q005 timeline / visa stage). **Minor without KTP:** follow **Q016** (KK substitute).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
