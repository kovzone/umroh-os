# booking-svc — Overview

## Purpose

The link between jamaah and packages. Owns reservations, room/bus allocations, and lightweight manifest references.

## Bounded context

Booking. See `docs/02-domain/00-bounded-contexts.md` § 4.

## PRD source

PRD sections A (B2C self-booking), B (B2B agent booking), and parts of E (room/bus allocation).

## Owns (data)

- `bookings` — reservation header
- `booking_items` — one row per jamaah on a booking
- `room_allocations` — which jamaah share which hotel room
- `bus_allocations` — bus seating
- `booking_status_history` — audit of state transitions

## Boundaries (does NOT own)

- Package or seat inventory (`catalog-svc`)
- Jamaah biodata (`jamaah-svc`)
- Payments (`payment-svc`)
- Visa state (`visa-svc`)
- Full manifests (`ops-svc`) — booking only stores draft allocations; ops-svc generates the final manifest

## Interactions

- **Inbound:** payment-svc signals payment status changes; visa-svc signals visa attached; ops-svc reads bookings for manifest generation.
- **Outbound:** catalog-svc to read package and reserve seats; jamaah-svc to read jamaah & validate mahram; iam-svc to check permissions.

## Notable behaviors

- **Booking saga** (Temporal, owned by broker-svc) orchestrates: catalog reserve → booking create → payment VA issue → notify.
- **Mahram validation** at booking time — for women under 45, verify a mahram is on the same booking by calling jamaah-svc.
- **Status machine:** draft → pending_payment → partially_paid → paid_in_full → completed (or → cancelled at any point).
- **Smart room allocation** algorithm runs at booking creation if family grouping is requested. Lives in ops-svc; booking-svc just stores the result.
