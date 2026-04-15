# ADR 0003 — Temporal.io for cross-service workflows

**Status:** Deferred for MVP — see ADR 0006 (2026-04-15)
**Originally Accepted:** 2026-04-09

> **Status update (2026-04-15):** Temporal is deferred from MVP per [ADR 0006](0006-defer-temporal-to-f6.md). In-process saga coordination is used for the short-lived booking / payment / refund sagas; a reconciliation cron catches mid-saga crashes. This ADR's body remains as the rationale for **eventual reintroduction** — specifically for F6 visa pipeline, which is the workflow that genuinely needs Temporal's durability and long-running semantics. Do not treat the rest of this document as current MVP direction.

## Context

UmrohOS has many long-running, multi-step business processes that span services:

- **Booking saga** — reserve seat → create booking → issue VA → notify
- **Payment reconciliation** — webhook → mark paid → fulfill kit → write journal
- **Visa pipeline** — submit → poll → on-issued → attach (runs for days)
- **Refund flow** — refund request → cancel booking → release seat → release stock → reverse journal

Each requires retries, compensating actions, durability across restarts, and visibility into state. Building this manually with cron jobs and ad-hoc state machines is error-prone.

## Decision

Use **Temporal.io** for cross-service orchestration. A dedicated `broker-svc` owns workflow definitions and activities. Other services expose plain gRPC methods that activities call into.

## Rationale

1. **Durable execution.** Workflows survive process restarts; state is persisted automatically.
2. **First-class retries and compensation.** Sagas with rollback are a built-in pattern.
3. **Long-running support.** Visa pipelines that take days don't need ad-hoc job queues.
4. **Observability.** Temporal Web UI shows every workflow's state and history. Combined with OTel, full visibility.
5. **Template alignment.** `broker-svc` already exists in `baseline/go-backend-template` as a Temporal example.
6. **Trace propagation.** OTel context propagates through Temporal activities.

## Consequences

- `broker-svc` is the only service that imports the Temporal SDK as a workflow author. Other services expose pure gRPC and remain Temporal-agnostic.
- Activities live in `broker-svc/internal/activities/`; workflows in `broker-svc/internal/workflows/`.
- Synchronous internal calls remain plain gRPC. Temporal is reserved for processes that genuinely need durability or are multi-step + cross-service.
- We accept the operational overhead of running Temporal (server + UI containers in `temporal/` already configured).

## Alternatives considered

- **Kafka + custom consumers** — more flexible for event fanout but no built-in saga support, no durable workflows, no first-class compensation. Rejected for MVP.
- **Database-backed job queue (e.g. River, Asynq)** — fine for simple background jobs but not built for multi-step sagas. Rejected.
- **No orchestration; direct gRPC chains** — fragile under partial failures, no retries, no observability. Rejected.

## When NOT to use Temporal

- Plain CRUD endpoints that touch one service. Use direct gRPC.
- Read-only cross-service joins. Use gRPC fanout.
- Notifications that don't need durability (e.g. metric increment). Use direct calls.

If a process is one service, one step, no retries needed → no Temporal.
