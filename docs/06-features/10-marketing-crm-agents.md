---
id: F10
title: Marketing, CRM, Agent Network
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 8 Must Haves
prd_sections:
  - "B. B2B Front-End"
  - "C. Marketing & Sales"
  - "J. Alumni Hub (partial)"
modules:
  - "#25–70 + #199–202"
depends_on: [F1, F4]
---

# F10 — Marketing, CRM, Agent Network

## Purpose & personas

TBD — B2B agent ecosystem (onboarding, replicated websites, commissions with overriding logic), marketing campaigns, lead nurturing, CS performance tracking, alumni referrals.

## Sources

- PRD Sections B, C, J
- Modules #25–70 and #199–202

## User workflows

TBD:
- W1: Agent self-onboarding + e-KYC + e-signature
- W2: Agent shares product via replicated site; click tracked via UTM
- W3: Lead → booking attribution; commission calculated
- W4: Super-agent's overriding commission calculated from downline closings
- W5: CS handles WhatsApp leads with round-robin assignment + response-time tracking
- W6: Campaign launch with A/B landing pages; ROAS reported
- W7: Broadcast to 5,000+ agents with rate limiting
- W8: Alumni referral code usage

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD. Critical: commission math with multi-level overrides; broadcast rate limiting; attribution windows.

## Data & state implications

See `docs/03-services/09-crm-svc/02-data-model.md`.

## API surface (high-level)

See `docs/03-services/09-crm-svc/01-api.md`.

## Dependencies

- F1 (IAM), F4 (booking — attribution source)

## Backend notes

TBD. Commission math has multiple edge cases that deserve unit-test coverage. WhatsApp adapter needs retry + dead-letter.

## Frontend notes

TBD. B2B portal is a substantial frontend surface. Replicated sites have per-agent white-labeling.

## Open questions

None yet. Expected: agent level thresholds, commission % per level, referral reward economics.
