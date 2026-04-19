---
id: Q028
title: Visa rejection handling policy (retry, escalate, refund, customer comms)
asked_by: session 2026-04-15 F6 draft
asked_date: 2026-04-15
blocks: F6, F5, F4
status: answered
---

# Q028 — Visa rejection handling policy

## Context

PRD acknowledges that visa applications can be rejected (state `REJECTED_BY_EMBASSY` per L1605) but doesn't prescribe the handling workflow. The question matters because rejections are **expensive**: a rejected visa typically means the jamaah can't travel, the provider fee is often non-refundable, and the booking enters the refund flow (which touches F4, F5, F12 penalty matrix).

There are also **policy choices** the system has to encode:
- Does the system auto-resubmit (with new docs) or always escalate to ops?
- What does the customer see? How fast? In what language?
- Does rejection auto-trigger refund or is it ops-gated?
- If the agency believes the rejection is erroneous, is there an appeal path?

## The question

1. **Auto-resubmission policy.** If the rejection reason is tractable (e.g., "Paspor buram" → request clear photo → resubmit), does the system auto-request a re-upload from the jamaah, then auto-resubmit once docs are re-verified? Or does every rejection go through ops?

2. **Customer communication.**
   - Channel: WhatsApp (consistent with doc-reject pattern per PRD L1585) + email + in-app?
   - Timing: immediate on rejection detection, or after ops reviews?
   - Language: Bahasa by default?
   - Tone/template: provider-specific rejection message forwarded as-is, or agency-crafted explanation?

3. **Refund / cancellation trigger.**
   - Auto-trigger F5 refund flow on rejection?
   - Ops reviews first, then decides (refund, resubmit, appeal)?
   - Timeline before auto-trigger if ops doesn't act?

4. **Appeal path.** If the agency believes the rejection is erroneous (common for mahram edge cases, name-mismatch disputes), is there a manual provider re-submission or appeal process? Does the system track appeals as first-class entities?

5. **Provider fee handling.** Rejected application often means provider fee is consumed. Does this flow into the refund calc per Q012 as a sunk cost, or is it agency-absorbed?

## Options considered

- **Option A — Ops-gated everywhere.** All rejections pause in an ops queue; ops reviews, decides action (resubmit / appeal / refund). Customer is notified only after ops decides. No auto-resubmit, no auto-refund.
  - Pros: tight control; no customer confusion from reactive automation; consistent messaging.
  - Cons: latency; requires ops responsiveness; bottleneck in busy season.
- **Option B — Two-tier automation + ops review.**
  - Immediate: customer gets a "we're looking into this" WA notification (no decision, just acknowledgement).
  - Ops reviews reason and decides: resubmit (if fixable), escalate to appeal, or trigger refund.
  - Final outcome communicated to customer via WA when ops commits.
  - Pros: customer feels acknowledged fast; ops still controls the substantive path.
  - Cons: two-touch customer comm is slightly more template work.
- **Option C — Category-driven auto-handling.**
  - Classify rejection reasons into buckets: `fixable_docs`, `mahram_issue`, `expired_passport`, `name_mismatch`, `insufficient_quota`, `other`.
  - `fixable_docs` → auto-request re-upload + auto-resubmit after re-verification.
  - `expired_passport` → auto-trigger refund (no recovery path).
  - `other` categories → ops review.
  - Pros: automates the common case.
  - Cons: category classification is fragile; provider error messages are in Arabic/broken English.

## Recommendation

**Option B — two-tier: fast acknowledgement WA + ops review for substantive decision.**

Workflow:

1. On `REJECTED_BY_EMBASSY` transition, immediately send WA + email to jamaah + linked agent: *"Your visa application for {package} departure {date} has been returned by the embassy. Our team is reviewing the reason and will contact you within 24 hours with next steps. Application ID: {id}."* In Bahasa; template-editable by Super Admin.

2. System opens an ops review task (new `visa_review_tasks` table) with the rejection reason, provider response blob, and a suggested-action dropdown (`resubmit` / `appeal` / `refund` / `need_more_info`).

3. Ops reviews within 24h SLA. On decision:
   - `resubmit` → creates a new `visa_applications` row linked via `resubmission_of`; if docs need update, WA to jamaah with specific instruction; tracks through W1 again.
   - `appeal` → logs an appeal record; manual provider communication by ops; outcome recorded.
   - `refund` → calls `booking-svc.CancelBooking(reason: 'visa_rejected')` which kicks off F5 refund flow with Q012 penalty matrix (visa-fee sunk cost applied).
   - `need_more_info` → free-text message to jamaah requesting clarification; task stays open.

4. Customer receives a second WA with the ops decision + next steps.

5. No auto-resubmit. The false-positive rate on rejections is high enough (provider systems are noisy) that automated loops could waste provider fees. Ops review is cheap insurance.

Provider fee handling per Q012: visa fee is a sunk cost in the refund calculation. Agency can waive case-by-case via Q012's override path.

Reasoning: customer anxiety is the operational risk; 24h silence on a rejection is much worse than 10 minutes of "we're looking." Ops review adds human judgment to edge cases without blocking the acknowledgement path. Option C's category auto-handling is tempting but provider error messages don't categorize cleanly enough for production reliability.

Reversibility: auto-acknowledgement template is config; escalating to auto-resubmit for some categories later is additive; never the other way (can't retract an auto-action).

## Answer

**Decided:** **Option B** — **immediate WA+email acknowledgement** (no legal commitment), **ops task** with suggested action, **no auto-resubmit** in MVP, **24h ops SLA** for substantive reply; **refund path** calls booking cancel + **Q012** sunk-cost rules; **provider fees** pass-through per Q012.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
