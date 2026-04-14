---
id: F12
title: Alumni Hub & Daily App
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 0 Must Haves (all Should or Could Have) — defer until core platform ships
prd_sections:
  - "J. Daily App & Alumni Hub"
modules:
  - "#190–202"
depends_on: [F1, F3, F4]
---

# F12 — Alumni Hub & Daily App

## Purpose & personas

TBD — post-pilgrimage engagement: alumni community, daily worship tools (prayer times, qibla, Quran), referral hub, ZISWAF (charitable giving). None are Must Have; defer until the revenue-generating path (F4, F5) is live.

## Sources

- PRD Section J
- Modules #190–202

## User workflows

TBD. Most flows are light CRUD + feed surfaces.

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD.

## Data & state implications

Community forums, referral codes, ZISWAF ledger entries (routes to finance for ZISWAF).

## API surface (high-level)

Split across `crm-svc` (community, referrals) and `jamaah-svc` (profile-linked data).

## Dependencies

- F1, F3, F4

## Backend notes

TBD. Explicitly parked until MVP features are operating. No backend work in early phases.

## Frontend notes

TBD.

## Open questions

None yet.
