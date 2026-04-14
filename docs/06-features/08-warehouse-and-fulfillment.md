---
id: F8
title: Warehouse, Procurement, Fulfillment
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 6 Must Haves
prd_sections:
  - "F. Inventory & Logistics"
modules:
  - "#109–128"
depends_on: [F1, F4, F5]
---

# F8 — Warehouse, Procurement, Fulfillment

## Purpose & personas

TBD — multi-warehouse stock, procurement (PR → approval → PO → vendor → GRN with QC), kit assembly, and shipment with courier integration. Triggered by paid bookings.

## Sources

- PRD Section F in full
- Modules #109–128

## User workflows

TBD:
- W1: Warehouse staff files PR; multi-level approval; PO dispatched
- W2: Goods received → GRN → QC → auto-trigger AP entry in finance
- W3: Booking paid → kit dispatch → shipment created with courier label
- W4: Low-stock alert
- W5: Stock opname (periodic inventory count)

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD.

## Data & state implications

See `docs/03-services/07-logistics-svc/02-data-model.md`.

## API surface (high-level)

See `docs/03-services/07-logistics-svc/01-api.md`.

## Dependencies

- F1 (IAM), F4 (booking), F5 (payment status)

## Backend notes

TBD.

## Frontend notes

TBD.

## Open questions

None yet.
