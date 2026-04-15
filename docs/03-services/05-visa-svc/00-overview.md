# visa-svc — Overview

> **Note on Temporal / broker-svc references below.** Temporal is currently deferred from MVP (see [ADR 0006](../../01-architecture/adr/0006-defer-temporal-to-f6.md)). This service's implementation (F6 visa pipeline) is what **brings Temporal and broker-svc back into the stack** — the long-running multi-day workflow is the genuine Temporal use case. References below to "Temporal workflow" and "broker-svc activities" describe the F6 target state; they are planned, not current MVP direction.

## Purpose

Visa application lifecycle, e-visa storage, tasreh records, and Raudhah Shield monitoring. The only service that talks to Saudi visa systems (MOFA / Sajil / Nusuk).

## Bounded context

Visa. See `docs/02-domain/00-bounded-contexts.md` § 6.

## PRD source

PRD section E (operational handling — visa portion) and the Raudhah Shield innovation in section I.

## Owns (data)

- `visa_applications` — one per jamaah per booking
- `visa_status_history` — status transitions
- `e_visas` — issued visa documents
- `tasreh_records` — Raudhah and other permission letters
- `raudhah_monitoring` — Nusuk-checked status snapshots

## Boundaries (does NOT own)

- Jamaah biodata (`jamaah-svc`) — read via gRPC
- Bookings (`booking-svc`) — signal on visa attached

## Interactions

- **Inbound:** broker-svc visa pipeline workflow drives this; ops-svc reads visa status.
- **Outbound:** jamaah-svc (passport data), MOFA/Sajil API, Nusuk API, booking-svc (attach visa).

## Notable behaviors

- **Visa pipeline** is a Temporal workflow that may run for days. Activities live in broker-svc; visa-svc exposes pure gRPC for the activities to call.
- **Raudhah Shield.** Periodically polls Nusuk for the visa holder's status to detect misuse. Anti-fraud feature.
- **Status machine:** waiting_docs → docs_ready → submitted → issued | rejected.
