# broker-svc — Overview

> ⚠️ **DEFERRED FOR MVP — reserved for F6 visa pipeline.** Not scaffolded. The compose stack does not run Temporal in MVP. See [ADR 0006](../../01-architecture/adr/0006-defer-temporal-to-f6.md). MVP short-lived sagas (booking submit, payment refund) are coordinated in-process by the orchestrating service. The content below describes this service's **future** shape, when F6 brings Temporal back.

## Purpose

Temporal workflow orchestrator. Owns long-running, multi-step, cross-service business processes. The only service that imports the Temporal SDK as a workflow author. Other services remain Temporal-agnostic and expose plain gRPC.

## Bounded context

Workflows. See `docs/02-domain/00-bounded-contexts.md` § 11. Notably, broker-svc owns no business data — Temporal owns workflow state.

## PRD source

Cross-cutting — implements the multi-step processes described throughout the PRD (booking saga, payment reconciliation, visa pipeline, refund flow).

## Owns

- **Workflow definitions** (in `broker-svc/internal/workflows/`)
- **Activities** (in `broker-svc/internal/activities/`) — each activity wraps a gRPC call to one of the other services
- No persistent business data. Temporal handles workflow state.

## Boundaries

- Does not own any database tables (Temporal has its own database, configured separately).
- Does not contain business logic — all business decisions live in the service layer of the relevant service. broker-svc just orchestrates the order of operations and handles compensation.

## Workflows (planned)

### Booking saga
1. `ReserveSeats` (catalog-svc)
2. `CreateBookingItems` (booking-svc)
3. `IssueVirtualAccount` (payment-svc)
4. `NotifyCustomer` (crm-svc)
- Compensation on failure: release seats, cancel booking, void invoice.

### Payment reconciliation
1. `RecordPaymentEvent` (payment-svc) — triggered by webhook
2. `MarkBookingPaid` (booking-svc)
3. `DispatchKit` (logistics-svc) on full payment
4. `RecordJournalEntry` (finance-svc)

### Visa pipeline
1. `CreateApplication` (visa-svc) — triggered by `jamaah.documents_ready`
2. `SubmitApplication` (visa-svc)
3. Long poll: `GetApplicationStatus` until issued or rejected
4. `AttachVisaToBooking` (booking-svc)

### Refund flow
1. `RefundPayment` (payment-svc)
2. `CancelBooking` (booking-svc)
3. `ReleaseSeats` (catalog-svc)
4. `ReleaseStock` (logistics-svc)
5. `ReverseJournalEntry` (finance-svc)

## Notable behaviors

- **Compensation logic** — every workflow step has a corresponding undo action.
- **Trace propagation** — OTel context flows through Temporal activities.
- **Long-running** — visa pipeline may run for days; Temporal handles durability.
