---
id: F7
title: Operations — Verification, Smart Grouping, Manifests, Airport Handling
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 6 Must Haves
prd_sections:
  - "E. Operational & Handling"
modules:
  - "#90, #91, #92, #93, #94, #95, #100–108"
depends_on: [F1, F3, F4]
---

# F7 — Operations: Verification, Grouping, Manifests, Airport Handling

## Purpose & personas

TBD — back-office workflows: document verification queue (reviews OCR output), smart room/bus allocation, manifest generation for immigration, luggage tag QR issuance, airport scan ingestion, and field-execution apps.

## Sources

- PRD Section E in full
- Modules #90, #91, #92, #93, #94, #95, #100, #101, #102, #103, #104, #105, #106, #107, #108

## User workflows

TBD:
- W1: Ops reviewer approves/rejects a scanned document
- W2: Smart grouping runs at booking submit
- W3: Manifest generated for a departure
- W4: Airport handler scans luggage QR tags
- W5: Bus tour leader scans boarding QR
- W6: Zamzam distribution recording
- W7: Daily incident reporting from the field

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD.

## Data & state implications

See `docs/03-services/06-ops-svc/02-data-model.md`.

## API surface (high-level)

See `docs/03-services/06-ops-svc/01-api.md`.

## Dependencies

- F1 (IAM), F3 (jamaah), F4 (booking)

## Backend notes

TBD. Smart grouping algorithm is a non-trivial piece of logic — heavy unit tests.

## Frontend notes

TBD. Verification queue UI and manifest preview are ops admin screens. Field apps are separate mobile surfaces.

## Open questions

None yet.
