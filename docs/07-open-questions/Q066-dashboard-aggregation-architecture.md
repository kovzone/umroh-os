---
id: Q066
title: Dashboard aggregation architecture — service `/metrics` vs CQRS vs OLAP hybrid
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: open
---

# Q066 — Dashboard aggregation architecture

## Context

F11 dashboards compose data from 9 upstream services (F2, F4, F5, F6, F7, F8, F9, F10). How we architect the aggregation affects: latency (dashboard render time), coupling (cross-service dependencies on hot paths), complexity (infra + schema), and operational readiness of observability stacks (Grafana/Loki/Prometheus already live per ADR 0001).

PRD is silent. Three standard patterns compete:

1. **Service-owned `/metrics` endpoints** — each service owns its aggregation; F11 composes via gRPC/REST.
2. **CQRS projection** — a dedicated dashboard read-store consumes events from all services; F11 reads from there.
3. **OLAP / analytics store** — Postgres read-replica + Metabase or Grafana Postgres queries for ad-hoc; product dashboards still in Svelte.

## The question

1. **Pick a pattern** (A / B / C or hybrid) for MVP F11.
2. **`dashboard-svc` service vs gateway-svc aggregation** — should F11 have its own service, or does the gateway compose upstream calls?
3. **Caching layer** — Redis in front of `/metrics` endpoints? TTL strategy?
4. **Materialized views vs ad-hoc SQL** per-service — when do we materialize?
5. **Grafana usage** — already in the stack; do we also surface product dashboards via Grafana, or keep Grafana for observability only and build product UI in Svelte?
6. **Read-replica strategy** — does the dashboard pattern require a Postgres read-replica for heavy aggregations, or do we run against primary with caching?

## Options considered

- **Option A — Service-owned `/metrics` endpoints; F11 gateway composes; Redis cache; ad-hoc via Grafana on read-replica.** Each service exposes its own aggregation (e.g. `finance-svc GET /v1/metrics/cash-flow`); a thin composition layer in gateway-svc (or a new `dashboard-svc`) calls them. Grafana on Postgres read-replica for analyst / executive ad-hoc queries.
  - Pros: respects bounded contexts; simplest infra; no new read model; schema stays canonical; reuses existing stack.
  - Cons: N+1 query risk for multi-widget dashboards; service outages cascade to dashboards.
- **Option B — CQRS with dedicated `dashboard-svc` read-store consuming upstream events.** Dashboard-svc subscribes to all service events; builds denormalized read model optimized for dashboard queries.
  - Pros: fast reads; isolates dashboard load from hot OLTP paths; scales horizontally.
  - Cons: eventual consistency confusion; event-schema versioning risk; another Go module + another Postgres database; mismatch with ADR 0006 (which moved away from broker-heavy architectures for MVP).
- **Option C — Postgres read-replica + Grafana-only dashboards.** No custom Svelte dashboard UI; Grafana handles all.
  - Pros: zero custom code; fastest time-to-dashboards.
  - Cons: Grafana is a developer/analyst tool; not a product experience for executives; limited interactivity; mobile story weak.

## Recommendation

**Option A — service-owned `/metrics` endpoints + thin composition in a new `dashboard-svc` (or gateway-svc extension); Redis cache (5-min default TTL, per Q067 cadence tiers); materialized views per-service for expensive rollups; Grafana on read-replica for analyst + executive ad-hoc in parallel.**

Option B's CQRS projection is architecturally cleaner at scale but adds a full service + infrastructure layer for MVP that isn't justified by current volume — UmrohOS runs ~10K jamaah/year, not a billion events/day. The bounded-context respect of Option A outweighs the N+1 risk (mitigated by caching + materialized views). Option C trades product polish for speed; executives typically want a branded single-screen experience, which Grafana doesn't provide naturally. Running Grafana **in parallel** with service-owned `/metrics` (what Option A includes) is the right compromise: analysts get Grafana; executives get the product dashboard.

Defaults to propose: Each service adds a `/v1/metrics/*` REST group that returns pre-aggregated data shapes matching F11's widget needs. A new `dashboard-svc` (preferred) or gateway-svc composition layer calls these with per-request fan-out. Redis cache keyed on `(endpoint, filter_params, branch_id)`; 5-min default TTL; invalidation on event publish where feasible. Materialized views used when a single `/metrics` endpoint's raw SQL exceeds 500ms (e.g. financial reports, inventory valuation, revenue-by-campaign); nightly refresh + on-demand trigger. Grafana hits Postgres read-replica for ad-hoc; dashboards in product UI hit the `/metrics` path. If volume grows past 100K events/day, revisit CQRS (Option B).

Reversibility: adding a CQRS projection later is additive — service `/metrics` endpoints keep working while a parallel read-store stands up. Moving away from Grafana parallel path is a one-setting change.

## Answer

TBD — awaiting stakeholder input. Deciders: CTO / backend lead (architecture posture), Ops / observability lead (Grafana usage model), agency owner (dashboard UX expectation).
