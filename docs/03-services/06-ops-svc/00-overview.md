# ops-svc — Overview

## Purpose

Back-office operational workflows: document verification queue, smart room/bus allocation algorithm, full manifest generation, luggage tagging, airport handling events.

## Bounded context

Operations. See `docs/02-domain/00-bounded-contexts.md` § 7.

## PRD source

PRD section E — Operational & Handling.

## Owns (data)

- `verification_tasks` — documents waiting for human review
- `manifests` — generated manifests for departures
- `luggage_tags` — QR-coded tags
- `handling_events` — airport scan events
- `grouping_runs` — record of smart-allocation runs

## Boundaries (does NOT own)

- Documents themselves (`jamaah-svc`)
- Bookings (`booking-svc`) — ops reads bookings to build manifests
- Visa data (`visa-svc`)

## Interactions

- **Inbound:** UI for ops staff; broker-svc workflows.
- **Outbound:** jamaah-svc (read documents/family graph for grouping), booking-svc (read bookings), catalog-svc (read package details for manifests), visa-svc (read visa status).

## Notable behaviors

- **Smart room/bus allocation** algorithm — family-aware grouping that keeps mahram pairs together. This is the meaningful business logic that lives in ops-svc/service.
- **Verification queue** — staff approve/reject OCR results and documents.
- **Manifest generation** — produces immigration-format manifests; supports multiple output formats (PDF, Excel).
- **Airport handling** — staff scan luggage QR / boarding QR to mark events.
