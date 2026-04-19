---
id: Q029
title: Physical passport chain-of-custody process (module #96)
asked_by: session 2026-04-15 F6 draft
asked_date: 2026-04-15
blocks: F6
status: answered
---

# Q029 — Physical passport chain-of-custody

## Context

PRD module #96 Log Fisik Paspor (line 331) tracks where the physical passport is at any moment: `Agen | Pusat | Provider Visa | Kedutaan | returned`. The PRD describes this at the one-sentence level with no workflow, no SLA, and no signer/evidence requirements.

The concerns are operational:
1. **Lost passports** are disasters — they block the trip, require embassy replacement (weeks of delay), and embarrass the agency. Chain-of-custody documentation is the insurance policy.
2. **Handoff fraud** — a rare but documented risk where a passport is "misplaced" in the chain and reappears in someone else's hands (theft, identity fraud).
3. **SLA monitoring** — if a passport is sitting at the provider for 30 days without movement, that's a problem; if at the embassy for 7 days that might be normal.

Need stakeholder input to pin the states, SLAs, and evidence requirements.

## The question

1. **State enum.** PRD names 5 states; are these sufficient, or are there more (e.g., "in transit between agency locations," "with courier," "held at airport counter")?

2. **Required signers per state transition.**
   - Who signs when agent hands to pusat?
   - Who signs when pusat hands to provider?
   - Does the provider sign an intake receipt we record?
   - What happens on return — signature from jamaah, or physical presence?

3. **Evidence attachments.**
   - Photo of the physical passport at each handoff?
   - Scanned receipts?
   - Just text logs?

4. **SLA per leg.**
   - Agent → Pusat: 48 hours?
   - Pusat → Provider: same day?
   - Provider → Embassy: varies; provider's SLA?
   - Embassy → Returning: varies; visa-issuance timeline?

5. **Alerts for overdue passports.** What state + time triggers an ops alert? Dashboard surface or WhatsApp?

6. **Lost-passport workflow.** Distinct from the movement log — what happens procedurally when a passport is reported lost? Who's notified? What's the replacement process?

## Options considered

- **Option A — Full chain-of-custody with signer + photo per transition.** Every state change requires a signer name + photo upload; system writes to `passport_movements` with evidence URL.
  - Pros: forensic-quality audit trail; strongest protection against fraud and loss claims.
  - Cons: ops friction; may slow day-to-day workflow.
- **Option B — Minimal chain: states + timestamps, optional signer/photo.** System tracks transitions; signer and photo optional fields for sensitive transitions (e.g., handing to provider) but not required for routine ones (e.g., agent → pusat within the same office).
  - Pros: pragmatic balance; ops choose which transitions warrant full evidence.
  - Cons: inconsistency across ops staff.
- **Option C — Integrate with a courier service.** Use a tracked delivery service (JNE, Pos Indonesia) for inter-office transfers with tracking-number lookup; state transitions auto-populate from courier webhooks.
  - Pros: offloads tracking to professionals.
  - Cons: not all handoffs go through a courier (many are in-person between offices); courier integration is its own integration project.

## Recommendation

**Option B — minimal chain with required signer, optional photo, SLA alerts.**

Defaults:
- **States** (enum extends PRD's 5): `received_from_jamaah | at_agen_branch | at_pusat | dispatched_to_provider | at_provider | dispatched_to_embassy | at_embassy | returning | returned_to_jamaah | lost_reported`.
- **Required fields on every transition**: `from_state`, `to_state`, `signer_name` (text), `signer_role` (enum: `agen_staff | cs | ops | provider_rep | courier | jamaah`), `moved_at`, `notes` optional.
- **Photo optional**, but **required** for transitions involving the provider or embassy (highest-risk legs).
- **Courier tracking number optional** — captured as a `courier_reference` text field; no automated webhook integration in MVP.

SLA alerts per leg (defaults, agency-editable):
| From state | To state | SLA | Alert action |
|---|---|---|---|
| `received_from_jamaah` | `at_pusat` | 48h | Yellow flag |
| `at_pusat` | `dispatched_to_provider` | 72h | Yellow flag |
| `at_provider` | `dispatched_to_embassy` | 7d | Yellow; 14d → red + ops WA |
| `at_embassy` | `returning` | 14d | Yellow; 30d → red + ops WA |

Lost-passport workflow:
- Ops marks `to_state: lost_reported`; system auto-opens a high-severity incident (per Q024 matrix, `category: logistics, severity: high`).
- Escalation to agency owner within 30 minutes.
- Audit trail of every state transition for the affected passport is auto-attached to the incident.
- Jamaah is contacted within 1 hour for identity reverification and embassy replacement workflow.

Reasoning: Option A's per-transition photo upload is overkill for in-office handoffs that happen dozens of times a day; Option B concentrates evidence where risk is highest (provider, embassy) while keeping routine movements fast. Courier-integration (Option C) is a later optimization when the agency has a preferred courier partner.

Reversibility: tightening to Option A requires a UI tweak + making optional fields required; integrating a courier in Option C fashion is additive. Low commitment.

## Answer

**Decided:** **Option B** extended enum + **required signer fields** on all transitions + **photo mandatory** provider/embassy legs + **SLA alerts** per table + **lost → high incident** auto-route. **Courier:** tracking number optional text; **no webhook MVP**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
