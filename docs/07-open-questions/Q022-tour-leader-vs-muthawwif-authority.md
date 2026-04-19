---
id: Q022
title: Tour leader vs muthawwif override authority boundary
asked_by: session 2026-04-15 F7 draft
asked_date: 2026-04-15
blocks: F7
status: answered
---

# Q022 — Tour leader vs muthawwif override authority

## Context

Q015 already covers the high-level authority question (Smart Grouping trigger + ops/tour-leader override with muthawwif changes routed through tour leader). This finer question asks: **on the day of travel, when scanners fail or unusual situations arise, what can each role do independently?**

Real scenarios:
- Jamaah's ID card gets wet on the bus; QR won't scan. Tour leader manually marks bus boarding. Can a muthawwif do the same on a Saudi-side bus?
- Zamzam distribution scanner dies. Muthawwif manually marks jamaah X received their share. Is this logged differently from a scanned issuance?
- Room Check-In Kamar Cepat (W14): muthawwif reconciles Saudi hotel room assignments against the Smart Grouping output. Can tour leader do this too?
- Incident escalation: a muthawwif reporting a medical emergency should probably not need tour leader's permission.

The tension: tour leaders are agency-employees (Indonesia-based, company-trained, agency-loyal) while muthawwif may be Saudi-licensed contractors who work with multiple agencies. Authority boundaries are not just operational but also a trust/contract boundary.

## The question

For each of these in-field actions, which roles can perform them **independently**, which require **tour leader acknowledgement**, and which require **ops approval**?

| Action | Tour leader | Muthawwif |
|---|---|---|
| Scan jamaah ID card at bus boarding | Yes | ? |
| Manually mark boarding when scanner fails | Yes (with reason log) | ? |
| Zamzam distribution scan | ? | Yes |
| Manual zamzam issuance when scanner fails | ? | Yes (with reason log) |
| Tasreh scan at Raudhah | ? | Yes |
| Check-In Kamar Cepat (room reconciliation) | ? | Yes |
| Incident report | Yes | ? |
| Medical emergency escalation | Yes | ? |
| Bus time-stationary override (cancel alert) | Yes | ? |

Fill in the question marks.

Secondary: do manual overrides need **ops approval** post-action (async), or are they self-signed events audited without review?

## Options considered

- **Option A — Tour leader superset, muthawwif scoped to Saudi-side field actions.** Tour leader can do everything. Muthawwif can do Saudi-specific actions (zamzam, tasreh, room reconciliation, incident in Saudi). Manual overrides from either role are logged to audit; ops reviews anomalies async but doesn't gate.
  - Pros: clean authority model; matches Indonesian-agency-employee vs Saudi-contractor distinction; no real-time ops gating in-field.
  - Cons: a rogue muthawwif could theoretically mark false zamzam distributions (though audit + quota tracking mitigates).

- **Option B — Symmetric authority for field actions, with all overrides async-reviewed by ops.** Both roles can perform any field action; ops reviews overrides after the fact with a soft "acknowledge / investigate" workflow.
  - Pros: simpler code (no per-role RBAC in the field app); maximum operational resilience.
  - Cons: blurs the accountability line; harder to assign blame if something goes wrong.

- **Option C — Strict separation with real-time ops approval for overrides.** Tour leader and muthawwif have separate scoped abilities; manual overrides require ops approval before commit (push notification to ops dashboard, approve/deny in-app).
  - Pros: tightest control.
  - Cons: field apps fail in poor-connectivity moments; ops staff would be overwhelmed by real-time approval requests.

## Recommendation

**Option A — Tour leader superset, muthawwif scoped to Saudi field actions, manual overrides logged asynchronously.**

Specific abilities:

| Action | Tour leader | Muthawwif |
|---|---|---|
| Scan jamaah ID card at bus boarding | Yes | Yes (Saudi-side buses) |
| Manually mark boarding when scanner fails | Yes (with reason log + photo if available) | Yes (Saudi-side, with reason log) |
| Zamzam distribution scan | Yes (coverage when muthawwif unavailable) | Yes (primary) |
| Manual zamzam issuance | Yes (with reason log) | Yes (with reason log) |
| Tasreh scan at Raudhah | No (not typically at Raudhah) | Yes |
| Check-In Kamar Cepat | Yes (if on-site; rare) | Yes (primary) |
| Incident report | Yes | Yes |
| Medical emergency escalation | Yes | Yes |
| Bus time-stationary override | Yes | No (tour leader's bus, not muthawwif's concern) |

Manual overrides are **async-audited, not gated**:
- Every override writes a `handling_events` or `incidents` row with `override: true, override_reason: '<text>', override_evidence_url: '<optional photo>'`.
- Ops dashboard surfaces a daily digest of overrides for review — not a real-time gate.
- Patterns that look abusive (e.g., one muthawwif marking 100% manual zamzam issuances for a departure) trigger an alert per Q024 escalation matrix.

Reasoning: airport WiFi and Mecca hotel connectivity make real-time approval gates operationally dangerous (Option C); symmetric authority undermines the employee-vs-contractor accountability distinction (Option B). Option A preserves organizational roles while keeping the field app usable.

Reversibility: role-scoped action lists are config per role (in F1's permission matrix); tightening or loosening later is a permission flip. Changing from async audit to real-time approval would require UI + workflow surgery but is an additive change if ever needed.

## Answer

**Decided:** **Option A** table + **async audit** (no real-time ops gate). Matches Recommendation matrix; **manual overrides** require **reason + optional photo**; **daily ops digest** + **Q024** escalation on abuse patterns.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
