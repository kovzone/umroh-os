---
slice: Sx
title: Slice Sx — Integration Contract
status: template
last_updated: 2026-04-20
pr_owner: TBD
reviewer: TBD
task_codes:
  - Sx-J-NN
---

# Slice Sx — Integration Contract

> **This is the template.** Copy this file to `slice-S1.md` / `slice-S2.md` / ... when the slice's first `Sx-J-*` card lands. Fill only the sections that apply to the slice; delete the rest. Keep the `§ Changelog` section regardless.
>
> **Who writes this:** the developer executing the slice's `Sx-J-*` contract cards (either dev, pick-up-and-notify).
> **Who reviews:** the other developer, before any `Sx-{L,E}-*` implementation work starts.

## Purpose

One-paragraph summary of what the slice does end-to-end and which services it binds together. Reference `docs/00-overview/05-slice-engineering-checklist-and-task-codes.md` for the slice's user-journey description. Reference the relevant `docs/06-features/` files for user-POV detail. Keep this section short — it's an orientation, not a spec.

## Scope

Bullet list of exactly which APIs, gRPC methods, and events this contract covers — and which it **does not** (so downstream implementers don't ship surprises). Example:

- In scope: `GET /v1/packages`, `GET /v1/packages/{id}`, `GET /v1/departures/{id}`.
- Out of scope: write endpoints on packages (deferred to admin-svc slice); agent-network pricing overrides (deferred to S4).

---

## § Catalog

_(Fill when the slice has a Catalog-facing contract. Otherwise delete this section.)_

For each REST endpoint, document:

- **Method + path**: e.g. `GET /v1/packages`
- **Purpose**: one sentence.
- **Request**: query params / path params / body shape with a JSON example.
- **Response**: body shape with a JSON example.
- **Errors**: list of error codes + when each applies.
- **Auth**: public / staff / agent / pilgrim.
- **Honors Q-NN**: link open-question IDs that govern this endpoint's shape.

```json
// example response shape — replace per slice
{
  "packages": [
    { "id": "pkg_...", "name": "...", "price_idr": 0 }
  ],
  "page": { "cursor": null }
}
```

---

## § Booking

_(Fill when the slice adds booking-write endpoints. Otherwise delete this section.)_

Same per-endpoint schema as `§ Catalog`. Flag idempotency key semantics explicitly — required for every `POST` that can retry.

---

## § Booking States

_(Fill when the slice changes the booking state machine. Otherwise delete this section.)_

List every booking status enum value valid in this slice + allowed transitions. Note which transitions require compensation events. Reference `docs/03-services/<svc>/02-data-model.md` for the persisted enum.

---

## § Inventory

_(Fill when the slice calls seat reservation or release via gRPC. Otherwise delete this section.)_

For each gRPC method, document:

- **Service + method**: e.g. `catalog.CatalogService/ReserveSeats`
- **Request message**: Protobuf-style field list with types.
- **Response message**: Protobuf-style field list with types.
- **Failure codes**: which gRPC codes mean what; how the caller should compensate.
- **Atomicity**: spell out the SQL pattern or lock semantics.
- **Compensation**: the paired `Release*` call — how and when the caller invokes it.

---

## § Webhook

_(Fill when the slice receives webhooks from an external provider — payment gateway, courier, visa provider. Otherwise delete this section.)_

Document:

- **Source**: provider name + which endpoint they POST to.
- **Signature verification**: header name + secret source (Viper key).
- **Minimal body**: fields the service needs, with a JSON example.
- **Dedupe key**: the `(provider, provider_txn_id)` or equivalent unique tuple.
- **Response codes**: 200 on first-accept, 200 on dedupe-seen, 4xx on signature fail, 5xx reserved for service bugs.

---

## § Booking integration

_(Fill when this slice pushes a state change back to `booking-svc` — e.g. S2 payment result. Otherwise delete this section.)_

Document the callback path + idempotency key + accepted state transitions. Cross-link `§ Booking States` above.

---

## § Events

_(Fill when this slice emits or consumes domain events. Otherwise delete this section.)_

For each event, document:

- **Name**: e.g. `payment.received`, `booking.paid_in_full`.
- **Trigger**: what causes the producer to emit.
- **Producer**: which service.
- **Consumer(s)**: which service(s) subscribe.
- **Payload**: JSON shape with an example.
- **Delivery guarantee**: at-least-once is the default; note exceptions.
- **Ordering**: whether consumers can assume a total order per key.

---

## § Fulfillment task

_(Fill when the slice introduces fulfillment tasks — dispatch, visa submission, etc. Otherwise delete this section.)_

Document the task schema: status enum, assignee field, due-by field, error-resolution path. Cross-link `docs/03-services/06-logistics-svc/` or the relevant service.

---

## § Journal

_(Fill when the slice posts to the finance journal. Otherwise delete this section.)_

Document the journal-entry shape the slice produces: placeholder accounts (pre-COA-seed), amount rules, debit/credit convention, reference to `docs/06-features/09-finance-and-accounting.md`.

---

## § Lead

_(Fill when the slice creates or mutates leads — S4 growth loop. Otherwise delete this section.)_

Document the lead schema, UTM-snapshot fields, attribution rules, and the booking-linkage callback.

---

## § UAT

_(Fill for S5 Hardening slice. Otherwise delete this section.)_

Enumerate mandatory UAT scenarios — one row per scenario with actor, preconditions, steps, expected outcome. Reference `docs/91-progress/progress.md` for the MVP gates list.

---

## § Changelog

_(Always keep this section.)_

One line per contract change after initial merge. Format: `YYYY-MM-DD — <what changed> — <rationale or Q/BL reference>`. Additive changes land here; breaking changes bump to `slice-Sx-v2.md` (see `README.md § Bump-versi rule`).

- 2026-MM-DD — initial version merged via task `Sx-J-01`.
