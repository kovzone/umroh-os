---
id: F11
title: Dashboards & Reporting
status: stub — spec to be written before implementation starts
last_updated: 2026-04-14
moscow_profile: 4 Must Haves (Operational Readiness + Financial Health) + several Should Haves
prd_sections:
  - "K. Executive Dashboards"
modules:
  - "#177–189"
depends_on: [F1, F2, F4, F5, F7, F8, F9]
---

# F11 — Dashboards & Reporting

## Purpose & personas

TBD — executive dashboards across operational, sales, field, inventory, and financial domains. Backend exposes aggregation endpoints; frontend renders.

## Sources

- PRD Section K
- Modules #177–189

## User workflows

TBD:
- W1: Director opens home dashboard, sees KPIs across all domains
- W2: Ops manager checks field radar / Raudhah status / luggage location in real time
- W3: CFO reviews cash flow and aging reports
- W4: CMO reviews ROAS and CS performance

## Acceptance criteria

TBD.

## Edge cases & error paths

TBD. Critical: aggregation performance, stale caches, permission scoping (branch managers see only their branch).

## Data & state implications

Read-only aggregations across multiple services. Could use materialized views in each service's DB, or a dedicated analytics store later.

## API surface (high-level)

Each owning service exposes a `GET /v1/metrics/...` aggregation endpoint. Contracts TBD.

## Dependencies

Most other features must be live first.

## Backend notes

TBD. Default to service-owned aggregations; escalate to a dedicated analytics store if performance becomes a problem.

## Frontend notes

TBD. Heaviest UI surface in the product — charts, maps, real-time indicators.

## Open questions

None yet.
