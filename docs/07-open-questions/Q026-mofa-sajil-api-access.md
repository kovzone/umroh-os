---
id: Q026
title: MOFA / Sajil API access — credentials, sandbox, contracts, rate limits
asked_by: session 2026-04-15 F6 draft
asked_date: 2026-04-15
blocks: F6
status: open
---

# Q026 — MOFA / Sajil API access

## Context

F6 W2 (bulk submit) and W3 (status poll) depend on programmatic access to Saudi's MOFA and/or Sajil visa systems. PRD line 333 names them as pipeline stages (`Input Data > MOFA > Sajil > Issued`) but gives no technical detail — no endpoint, no auth model, no rate limits, no sandbox.

This is the hardest integration in the system from an access-logistics standpoint. Saudi government systems are not public APIs with signup flows. Access typically requires:
- Agency's PPIU license (Indonesian Umroh operator authorization)
- Registration as a "Tafweej" or equivalent intermediary with Saudi MoHRD
- A designated technical contact at the provider side
- A sandbox environment that may or may not be offered

Without this pinned down, F6 can't leave the spec-drafting phase.

## The question

1. **Which system does the agency have access to today?**
   - Direct MOFA portal (manual web entry by ops staff)
   - Sajil (typically via a licensed intermediary)
   - A third-party wrapper service (TafweejIT, Maqsad, Wathiqa, etc.) that exposes MOFA/Sajil behind a unified API
   - Something else

2. **Are there API credentials available?** If yes, what auth model (API key, OAuth, mTLS, signed requests)? If no, is there a path to obtain them, and what's the timeline?

3. **Is there a sandbox / test environment?** Without a sandbox, every dev test hits production — unacceptable for iteration.

4. **Rate limits and request quotas.** Per-minute, per-day, per-batch?

5. **Schema / contract documentation.** OpenAPI spec? PDF? Swagger? Email-a-sample-request?

6. **What currently happens manually that we'd automate?** Understanding the manual baseline helps scope what the automation has to replicate.

## Options considered

- **Option A — Direct MOFA/Sajil API integration.** Best long-term if access is available; most automatable.
  - Pros: no middleman; agency controls the integration.
  - Cons: access may not be available; Saudi systems are known to change schemas without notice.
- **Option B — Third-party wrapper (TafweejIT, Maqsad, etc.).** Buy access through a specialized vendor that handles the Saudi-side quirks.
  - Pros: faster path to working integration; vendor handles breaking changes.
  - Cons: ongoing fees; another dependency; vendor's own reliability.
- **Option C — Hybrid: manual submission via ops portal, automated status poll via screen-scrape or provider API.** Ops continues to submit via the MOFA/Sajil web UI; our automation picks up status by polling.
  - Pros: immediate; no access negotiations needed.
  - Cons: screen scraping is fragile; doesn't automate the submit side which is the bigger operational burden.
- **Option D — MVP operates fully manual; automation is Module #99 (Could Have) for later.** Ops staff handles submission + status entry in the UmrohOS UI as today; the system provides UI affordance + tracking but no Saudi integration.
  - Pros: ships without blocking on access negotiations.
  - Cons: doesn't realize F6's automation value.

## Recommendation

**Two-phase approach anchored by actual access**:

**Phase 1 (MVP launch):** **Option D — fully manual workflow**. Visa Progress Tracker UI, state machine, bulk ops actions, WhatsApp notifications on status change — all of this works without provider integration. Ops staff continues to enter status into the UmrohOS UI after checking the MOFA/Sajil portal manually. The **automated poll and submit are behind a feature flag**, off by default.

**Phase 2 (when access is confirmed):** enable the feature flag, wire in the appropriate adapter:
- If direct MOFA/Sajil API access is granted → native adapter.
- If only a third-party wrapper (TafweejIT etc.) is available → adapter against that wrapper.
- Either way, visa-svc's provider-adapter pattern absorbs the choice — service layer is unchanged.

Rationale: access to Saudi government systems is political/operational, not technical. Gating MVP launch on securing it would stall the whole product. The manual workflow is already what agencies do today; the UmrohOS value in Phase 1 is the progress tracking, document organization, state machine, and audit — not the auto-submit.

For the spec: document the adapter interface (`VisaProviderAdapter` with `Submit`, `GetStatus`, `DownloadEVisa`) and ship a `ManualAdapter` as the default implementation. When a real provider is wired, it's a new adapter plus config.

Reversibility: the adapter-interface design means switching providers is additive; Phase 2 doesn't disturb Phase 1 data.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (MoU / MoA with providers) + operational head (current manual workflow details) + potentially legal (intermediary agreements).
