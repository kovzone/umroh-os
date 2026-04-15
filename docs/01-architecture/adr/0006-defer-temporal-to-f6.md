# ADR 0006 — Defer Temporal from MVP; reintroduce for F6 visa pipeline

**Status:** Accepted
**Date:** 2026-04-15
**Supersedes (for MVP):** ADR 0003 — the rationale in 0003 remains valid for the use case that justifies Temporal, but Temporal is deferred until that use case is in scope.

## Context

ADR 0003 locked Temporal.io in as the cross-service workflow orchestrator based on the baseline template and the PRD's many long-running processes (booking saga, payment reconciliation, visa pipeline, refund flow). After drafting F4 (booking) and F5 (payment) specs, we reviewed the four saga candidates again with more concrete detail:

| Workflow | Duration | Durability need | Cross-service | Compensation needed |
|---|---|---|---|---|
| Booking saga (F4) | seconds | weak (a crash is a rare, recoverable edge) | yes (3 services) | yes, short path |
| Payment reconciliation (F5) | seconds | weak | yes (3 services) | not really — webhooks are idempotent |
| Refund flow (F5) | seconds | moderate | yes (5 services) | yes, multi-step |
| Visa pipeline (F6) | **days** | **strong** (survives process restarts) | yes (3 services) | yes, rare but real |

Only F6 genuinely benefits from Temporal's core value proposition — durable, long-running workflows. The other three are short RPC chains with explicit compensation steps; they can be coordinated in-process by the orchestrating service with a reconciliation cron catching rare mid-saga crashes.

Running Temporal in MVP for the short-lived sagas means paying the full operational cost (Temporal server + UI + its own database + SDK learning curve + search attributes + workflow-testing harness) for benefits we're barely using.

## Decision

Defer Temporal from MVP. Use **in-process saga coordination** in the orchestrating service + a **reconciliation cron** to catch mid-saga crashes.

Bring Temporal back when F6 visa pipeline is implemented — that specific feature justifies the infrastructure.

## What "in-process saga coordination" means

For F4 booking submit, the saga runs inside `booking-svc.Submit()`:

```
1. WithTx (booking-svc DB): create draft booking row
2. call catalog-svc.ReserveSeats via gRPC
     on failure: delete draft → return error
3. call payment-svc.IssueVirtualAccount via gRPC
     on failure: catalog-svc.ReleaseSeats → delete draft → return error
4. WithTx: update booking status, append status_history
5. emit booking.dp_received event (fire-and-forget to F8/F9/F10)
6. return { booking_id, va_details }
```

Each step's failure has an explicit compensation that unwinds prior steps. Compensations are code, not framework.

For F5 refund flow, same pattern — `payment-svc` orchestrates `RefundPayment → CancelBooking → ReleaseSeats → ReleaseStock → ReverseJournal` with explicit per-step compensation.

The **reconciliation cron** (already needed for missed payment webhooks per F5 W5) is extended to catch mid-saga orphans:

- Invoice exists but no booking → void invoice + release seat
- Booking draft exists but no invoice → delete after TTL (Q010)
- Booking paid but F8/F9/F10 event not delivered → re-emit from `payment_events` history

`WithTx` at the store layer (per `baseline/go-backend-template/demo-svc/store/postgres_store/tx.go`) handles single-service multi-row atomicity. It is **orthogonal** to this decision — it does not substitute for Temporal, and every service uses it regardless.

## Rationale

1. **Operational weight.** Temporal adds a server, a UI, its own database, an SDK, and a new mental model (workflows vs activities, signals, queries, search attributes). For a two-dev MVP team, that's real daily overhead.
2. **Short-lived sagas don't benefit.** The booking saga completes in seconds. The durability guarantee Temporal provides is rarely exercised; a crash during a 3-step RPC chain is a narrow failure window that a reconciliation cron catches cheaply.
3. **Temporal's sweet spot is F6.** The visa pipeline runs for days, polls external systems, must survive process restarts, and needs explicit compensation when a visa is rejected after the pipeline has already progressed halfway. That's exactly what Temporal is built for.
4. **Migration cost later is low.** F4 and F5 in-process sagas do not need to change when Temporal is reintroduced — Temporal can coexist with in-process coordination. Different workflows use different mechanisms.
5. **`WithTx` was already in every service's toolkit.** Intra-service atomicity is handled by Postgres transactions, always. Deferring Temporal does not leave us without atomicity for single-service work.

## Consequences

- `broker-svc` is **deferred** — its service docs and folder remain as planning material, but it is not scaffolded in MVP. Marked with a "deferred — reserved for F6" banner.
- Compose file for MVP does not run Temporal containers.
- `booking-svc` gains in-process saga orchestration responsibility.
- `payment-svc` gains in-process refund-flow orchestration responsibility.
- The F5 reconciliation cron scope expands to include mid-saga orphan cleanup.
- ADR 0003 is annotated with a "deferred" banner but the body stays — its rationale is valid for when Temporal is reintroduced.
- Feature specs F4 and F5 are revised to describe in-process sagas.
- F6 visa-pipeline spec flags "Temporal is the planned mechanism for this feature; introducing Temporal back into the stack is part of F6 implementation."
- Cross-cutting docs (tech stack, service map, system overview, data-flow, observability, CLAUDE.md) reflect the defer.

## When to revisit

Triggers for bringing Temporal back (any one is sufficient):

1. **F6 visa pipeline implementation starts.** Non-negotiable; do not build F6 without durable workflows.
2. **A new long-running workflow appears** (multi-day, poll-based, needs durability) — likely candidates: Hajj quota allotment negotiation, multi-currency settlement reconciliation, alumni-retention nurture sequences.
3. **The reconciliation cron becomes complex enough to justify replacement** — e.g. if we end up encoding state machines into the cron that Temporal would express more naturally.
4. **Operational pain from hand-rolled compensations** — if we catch production bugs where a compensation path was missed or wrong, Temporal's explicit-saga pattern may justify itself.

## Alternatives considered and rejected

- **Keep Temporal for MVP regardless** — rejected: operational overhead, learning curve, overkill for short sagas. The visa pipeline trigger alone doesn't justify introducing it for F4/F5 where it adds more weight than it removes.
- **Never use Temporal; build our own state machine + cron for F6** — rejected: F6 workflow state persistence, retries, resume-from-crash, visibility — all things Temporal does well. Reinventing them worse is waste.
- **Use Kafka + consumer groups for event-driven coordination** — rejected: doesn't solve multi-day polling workflows; different problem shape. Still considered separately for event fan-out if that need arises.
- **Use a lighter workflow engine (River, Asynq, Faktory)** — rejected for now: they're fine for background jobs but weaker than Temporal for multi-step sagas with compensation and long-running polling. If F6 arrives and Temporal feels too heavy, a second ADR can revisit this.

## Notes

- `WithTx` is not an alternative to Temporal — they solve different problems. This ADR is NOT a choice between them.
- The other developer's `AGENTS.md` also references Temporal in its architectural rules. That file is outside this repo's editing scope; flag to the other dev separately so they can update their own file.
