---
id: F6
title: Visa Pipeline & Raudhah Shield
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 1 Must Have (Visa Progress Tracker) + Should Haves
prd_sections:
  - "E. Operational & Handling — Visa & Passport Tracking"
  - "I. Jamaah Journey — Raudhah Shield"
modules:
  - "#96, #97, #98, #99, #105"
depends_on: [F1, F3]
---

# F6 — Visa Pipeline & Raudhah Shield

## Purpose & personas

TBD — long-running Temporal workflow that submits Saudi visa applications, polls status, and attaches issued e-visas to bookings. Raudhah Shield polls the Nusuk app to detect visa misuse.

## Sources

- PRD Section E — Visa portion; Section I — Raudhah Shield
- Modules #96, #97, #98, #99, #105

## User workflows

TBD:
- W1: Ops admin submits a visa application after documents ready
- W2: System polls MOFA/Sajil for status
- W3: Visa issued → attached to booking, downloadable from portal
- W4: Raudhah Shield polls Nusuk during pilgrimage window
- W5: Tasreh (Raudhah permit) registration and scan-at-entry

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD. Critical: provider API downtime, rejection handling, tasreh double-use detection.

## Data & state implications

See `docs/03-services/05-visa-svc/02-data-model.md`.

## API surface (high-level)

See `docs/03-services/05-visa-svc/01-api.md`.

## Dependencies

- F1 (IAM), F3 (jamaah — passport data)

## Backend notes

TBD. Visa pipeline is a multi-day Temporal workflow. Long polls on activity side.

## Frontend notes

TBD.

## Open questions

None yet. Candidate: MOFA/Sajil API credentials and test environment.
