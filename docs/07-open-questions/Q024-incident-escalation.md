---
id: Q024
title: Incident / issue report escalation matrix
asked_by: session 2026-04-15 F7 draft
asked_date: 2026-04-15
blocks: F7
status: answered
---

# Q024 — Incident escalation matrix

## Context

Field execution (F7 W13) produces incidents across categories: medical emergency, lost jamaah, vendor problem, logistics, security, other. PRD mentions the capability (module #571 referenced in PRD, "Pelaporan Isu Harian") but doesn't define the **routing and response** — who gets alerted per category, SLA for ops response, how incidents close.

This is operational org-chart information — the agency knows who plays each role. The system needs this encoded to route alerts.

## The question

For each incident category + severity, define:

1. **First responder** (who gets the alert within seconds).
2. **Escalation tier 2** (who gets paged if first responder doesn't ack within SLA).
3. **SLA for ack** (how long until tier-2 escalation).
4. **SLA for resolution** (how long from creation to closure).
5. **Closure authority** (who can mark resolved; audit trail requirements).

Categories × severity matrix:

| Category | Severity | First responder | Tier 2 | Ack SLA | Resolve SLA |
|---|---|---|---|---|---|
| Medical | emergency | ? | ? | ? | ? |
| Medical | non-urgent | ? | ? | ? | ? |
| Lost jamaah | (always urgent) | ? | ? | ? | ? |
| Vendor problem | high | ? | ? | ? | ? |
| Vendor problem | low | ? | ? | ? | ? |
| Logistics | high | ? | ? | ? | ? |
| Logistics | low | ? | ? | ? | ? |
| Security | any | ? | ? | ? | ? |
| Other | any | ? | ? | ? | ? |

Secondary: notification channels — WhatsApp (primary), in-app push, email, SMS? Which category/severity triggers which?

## Options considered

- **Option A — Operational-safety-first defaults, configurable by agency.** Pre-populate with reasonable defaults that the agency can edit via the Admin console. Escalation configuration lives in `incident_escalation_config` table editable at runtime.
  - Pros: reasonable start; agency customizes over time; audit of who changed rules.
  - Cons: requires the agency to think through the matrix on initial setup; misconfiguration possible.

- **Option B — Simple flat rule: all incidents go to ops + tour leader immediately, ops escalates manually.** No tiered routing in the system.
  - Pros: simplest; trusts the ops team to manage.
  - Cons: under-utilises the scanner/alert capability; risk of missed urgent incidents when ops is overwhelmed.

- **Option C — Fixed rules, not configurable (hard-code agency's policy).** Ship with one escalation policy baked into the code.
  - Pros: no config surface.
  - Cons: changes require deploy; no flexibility per season or per package type.

## Recommendation

**Option A with sensible defaults** shipped out of the box. The agency confirms or edits during onboarding.

Suggested default matrix:

| Category | Severity | First responder | Tier 2 (if no ack) | Ack SLA | Resolve SLA |
|---|---|---|---|---|---|
| Medical | emergency | Tour leader + Ops on-call + Emergency contact | Agency owner | 2 min | 1 hour (reach medical facility) |
| Medical | non-urgent | Tour leader | Ops lead | 30 min | 24 hours |
| Lost jamaah | (always urgent) | Tour leader + Ops on-call | Agency owner | 5 min | 4 hours (locate or escalate to Saudi authorities) |
| Vendor problem | high | Ops lead | Agency owner | 30 min | 24 hours |
| Vendor problem | low | Ops lead | — | 4 hours | 48 hours |
| Logistics | high | Ops lead | Agency owner | 30 min | 12 hours |
| Logistics | low | Ops lead | — | 4 hours | 48 hours |
| Security | any | Tour leader + Agency owner | Saudi authorities (if applicable) | 5 min | Per incident |
| Other | any | Ops lead | — | 4 hours | 72 hours |

Channels:
- **Emergency tier:** WhatsApp + in-app push + SMS (redundant — at least one must reach).
- **High tier:** WhatsApp + in-app push.
- **Standard tier:** In-app push + daily digest.

Closure authority:
- **Tier 1 and 2 responders** can mark resolved with required resolution note.
- **Emergency incidents** require tour leader + ops lead both acknowledgement before closure; creates a paired audit trail.

Ops dashboard shows live incident feed with status, time-since-created, tier-2-escalation countdown.

Reasoning: medical emergencies need multiple redundant channels (airport WiFi unreliable; SMS is a fallback). Agency-owner involvement on emergencies matches Indonesian agency practice where the owner is personally accountable. Security incidents loop in Saudi authorities because we're not equipped to handle ourselves.

Reversibility: all escalation routing lives in a config table; changes apply to new incidents, existing incidents continue their original routing for closure. Low commitment.

## Answer

**Decided:** **Option A** — `incident_escalation_config` with **Recommendation default matrix** + **channel policy** (emergency = WA+push+SMS; high = WA+push; standard = push+digest). **Emergency closure:** **tour leader + ops lead** dual-ack. **All rows editable** by Super Admin at onboarding.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
