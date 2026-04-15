---
id: Q023
title: Zamzam distribution quota and policy per jamaah
asked_by: session 2026-04-15 F7 draft
asked_date: 2026-04-15
blocks: F7
status: open
---

# Q023 — Zamzam distribution quota

## Context

Module #107 Distribusi Zamzam (PRD line 363) describes QR-scan at distribution to enforce one-jamaah-one-share with no double-dipping. But the PRD doesn't specify the actual quota, the handling of edge cases (spilled container, sick jamaah's share redistributed to family), or whether muthawwif can waive.

This is a short but real stakeholder question — affects agency operational policy + muthawwif field procedure.

## The question

1. **What's the quota per jamaah per trip?** Typical Umroh packages include zamzam as a souvenir — ~5 liters in a sealed container is industry standard but the agency sets its own number.

2. **Is it one-time issuance or multi-stage?** Some packages distribute in Mecca (for immediate use) + shipped home separately. Others issue all at once.

3. **Can a jamaah's share be transferred to another jamaah?** Example: elderly jamaah cannot carry the container; spouse carries on their behalf. Real-world; needs a mechanism.

4. **What happens on scanner failure?** Muthawwif marks manual issuance — already covered by Q022. Here: is the jamaah then locked from further issuance, or does manual marking not count against the quota (since verification weakened)?

5. **Over-quota handling.** What if, for operational reasons, a jamaah needs an extra container (previous one confiscated at airport, broken, etc.)? Muthawwif override + audit, or ops approval required?

6. **Package-level zamzam cap.** Some agencies buy a bulk zamzam allocation from the agency supplier; does the system enforce a per-departure total alongside per-jamaah totals?

## Options considered

- **Option A — Simple: 5L per jamaah, one-time issuance, transfer allowed within same booking, muthawwif override with audit for edge cases.** Single numeric quota, clear default, minor override path.
  - Pros: easy to implement; matches most Indonesian Umroh conventions.
  - Cons: doesn't explicitly handle the "package bulk cap" question.

- **Option B — Configurable per package: quota + issuance stages + transfer rules.** Each package defines its zamzam policy; ops can set different rules for premium vs economy packages.
  - Pros: flexibility for differentiated packages.
  - Cons: more config surface; risk of misconfiguration.

- **Option C — Defer: manual-only issuance, no system quota enforcement.** Scanner records issuance; quota is a CS/muthawwif mental model, not enforced by the system.
  - Pros: simplest MVP; defers the policy conversation.
  - Cons: doesn't use the scanner for what it's specifically designed for; undermines module #107's value proposition.

## Recommendation

**Option A — 5 liters per jamaah, one-time issuance, in-booking transfer allowed, muthawwif override with audit.**

Specifics:
- Default quota: **5 liters per jamaah per trip**. Stored as a system config variable, editable by Super Admin.
- Distribution is **one-time**: once a jamaah's `handling_events` shows a `zamzam_distribution` with `status: completed`, subsequent scans return a soft warning and require override.
- **Transfer within same booking** is allowed: one member scans their own QR; muthawwif can tag the scan with `transferred_to: <other_jamaah_id in same booking>`. Both records written.
- **Override path** for legitimate edge cases (spilled, replaced at airport, etc.): muthawwif provides a free-text reason; logged to audit; daily digest to ops per Q024.
- **Package-level bulk cap is out of scope for MVP.** If an agency orders 500L for a 100-jamaah departure but a family of 4 transfers all of theirs to one member, we don't block that — the supplier-side quota is a procurement concern, not an in-app enforcement concern yet.

Reasoning: 5L is the Saudi-side tradition and what most Umroh packages deliver; making it configurable at the global level (not per package) keeps the policy simple to explain to CS and jamaah; transfer handling matches real family dynamics; package-level caps add complexity for an edge case.

Reversibility: quota value is a config variable; enabling per-package override is additive; adding package-level caps later is a new column + validator, no migration.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (commercial policy on zamzam allocation) + ops lead (field workflow realism).
