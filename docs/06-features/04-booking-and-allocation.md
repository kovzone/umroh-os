---
id: F4
title: Booking Creation & Room/Bus Allocation
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: Must Have core (B2C self-booking + smart grouping)
prd_sections:
  - "A. B2C Front-End — Self-Service Booking"
  - "E. Smart Grouping"
modules:
  - "#16, #17, #18, #92, #93, #95"
depends_on: [F1, F2, F3]
---

# F4 — Booking Creation & Room/Bus Allocation

## Purpose & personas

TBD — the link between a jamaah and a package. Owns bookings, booking items, room allocations, bus allocations, and the saga that orchestrates creation (reserve seats → create booking → issue VA → notify).

Primary personas: calon jamaah (self-book), agent (book for client), CS (book from call), ops admin (smart grouping + manifest prep).

## Sources

- PRD Sections A, E
- Modules #16, #17, #18, #92, #93, #95

## User workflows

TBD. Main flows:
- W1: Calon jamaah self-books via B2C portal
- W2: Agent books for multiple jamaah via B2B portal
- W3: CS books from WhatsApp conversation
- W4: Smart grouping runs at submit time, keeping families together in rooms and buses
- W5: Booking cancellation

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD. Critical: concurrent seat reservation race; mahram failure at submit; partial payment then abandonment.

## Data & state implications

See `docs/03-services/02-booking-svc/02-data-model.md`. Status machine: draft → pending_payment → partially_paid → paid_in_full → completed | cancelled.

## API surface (high-level)

See `docs/03-services/02-booking-svc/01-api.md`.

## Dependencies

- F1 (IAM)
- F2 (catalog — seat inventory)
- F3 (jamaah — mahram validation)

## Backend notes

TBD. Saga orchestrated by `broker-svc`. Compensations required for every step.

## Frontend notes

TBD.

## Open questions

None yet. Candidate: what happens when only one jamaah out of a 5-person booking cancels?
