---
id: F4
title: Booking Creation & Room/Bus Allocation
status: draft
last_updated: 2026-04-15
moscow_profile: Must Have core (B2C self-booking, B2B agent-booking, CS closing, smart grouping)
prd_sections:
  - "A. B2C — Self-Service Booking (lines 69–99)"
  - "B. B2B — Website Replika + Dompet Komisi (lines 117–145)"
  - "C.4 Sales Closing — payment link + e-Approval (lines 217–227)"
  - "E. Smart Grouping (lines 315–325)"
  - "E.5 Cancellation (lines 339–341)"
  - "Alur Logika 5.1–5.4 (lines 1171–1177), 6.4–6.5 (lines 1191–1193)"
modules:
  - "#16 Self-Booking Engine, #17 Guest Data Form, #18 Payment Gateway B2C, #92 Smart Room Allocation, #93 Transport Allocation, #95 ID Card & Luggage Tag, #100 Refund & Pinalti"
depends_on: [F1, F2, F3]
open_questions:
  - Q004 — cancellation → seat return ownership (existing)
  - Q005 — mahram validation at booking (existing)
  - Q006 — minimum docs required to submit (existing)
  - Q010 — VA TTL and draft booking expiry
  - Q014 — partial cancellation + reopening cancelled bookings
  - Q015 — Smart Grouping trigger timing + override authority
  - Q016 — booking on behalf of a minor (no KTP)
  - Q017 — Paket Tabungan interaction with booking state
  - Q019 — abandoned checkout attribution
---

# F4 — Booking Creation & Room/Bus Allocation

## Purpose & personas

This is the moment a pilgrim commits — where an "interested person" becomes a paying customer on a specific departure. Owns the booking record, the draft→paid status machine, and the reservation of seats from catalog inventory. Smart room/bus allocation (the operationally-critical "keep families together" algorithm) runs later in ops, but **F4 owns the data shape** that feeds it.

Three closing paths — all produce the same booking record, differentiated only by `channel` and `attribution`:

Primary personas:
- **Calon jamaah** — self-books on the B2C portal for themselves or a family group.
- **Agent** — closes on behalf of their customer via the B2B replicated site (`domain.com/id/<agent-slug>`). Commission attribution is automatic via the referral link.
- **CS** — closes from a WhatsApp conversation using the internal ERP; generates payment link and sends via WA.
- **Ops manager** — reads bookings for manifest generation, smart grouping runs, verification queue triage.
- **Downstream consumers** — `payment-svc` (issues VA), `finance-svc` (journal), `logistics-svc` (kit dispatch on lunas), `crm-svc` (commission attribution), `visa-svc` (visa readiness), `ops-svc` (manifest + grouping).

## Sources

- PRD Section A self-booking (lines 69–99)
- PRD Section B B2B agent closing (lines 117–145)
- PRD Section C.4 sales closing tools (lines 217–227)
- PRD Section E smart grouping (lines 315–325) and cancellation (lines 339–341)
- PRD Alur Logika 5.1–5.4 (B2C funnel) and 6.4–6.5 (manifest + cancellation)
- Modules #16, #17, #18, #92, #93, #95, #100

## User workflows

### W1 — B2C self-booking (Alur Logika 5.2)

1. Calon jamaah browses catalog, lands on package detail page, sees Seat Tracker `Sisa X Kursi!` (real-time from F2).
2. Clicks **Book Now** → `/booking/<package_id>`. Self-Booking Engine presents: number of pax, room-type config (Double / Triple / Quad), optional add-ons.
3. Guest Data Form collects **identity only** — full name, email, WhatsApp, domicile. Deep docs (KTP, passport) come later via Customer Portal self-upload (F3).
4. Clicks **Checkout** → `/checkout/<booking_id>`. At this point the booking is persisted as `draft` with `channel = 'b2c_self'`.
5. Frontend requests VA issuance from F5 (F5 W1). Booking transitions `draft → pending_payment` once VA is issued.
6. Jamaah pays to the VA. Webhook lands at F5, triggers `MarkBookingPaid`, booking transitions `pending_payment → partially_paid | paid_in_full` depending on amount received vs total.
7. Jamaah gains Customer Portal access (billing dashboard, self-upload, logistics, boarding pass).

### W2 — B2B agent closing (Alur Logika 5.2 + Section B)

1. Agent shares their replicated landing page link or an auto-watermarked flyer to WhatsApp/Instagram. Link carries a **referral token** (`ref=<agent_code>`).
2. Customer clicks through → B2C flow (W1) but with `channel = 'b2b_agent'` and `agent_id` pre-stamped on the booking.
3. On DP receipt, agent's Commission Wallet shows the entry as **Pending**; on lunas, transitions to **Confirmed** (F10 commission ledger).
4. Alternative: agent fills the form on behalf of the customer directly in the B2B portal, then pastes the VA link into their own WhatsApp for delivery.

### W3 — CS closing (Section C.4)

1. CS handles a WhatsApp inquiry using the internal ERP Leads module (F10).
2. When the jamaah is close to commit, CS opens Modul Pembuat Link Pembayaran inside the booking screen → fills pax / package / room type / optional discount via e-Approval.
3. System creates booking (`channel = 'cs'`, `created_by = <cs_user_id>`) in `draft`, immediately issues VA via F5, returns a WhatsApp-shareable link.
4. CS pastes the link into the WA conversation. Jamaah pays; remainder identical to W1.

### W4 — Status machine

```
draft
  │
  │ (submit + VA issued)
  ▼
pending_payment ──┐
  │               │ (VA expires without payment, see Q010)
  │               ▼
  │             expired ───┐
  │                        │
  │ (first payment received)
  ▼                        │
partially_paid             │
  │                        │
  │ (cumulative ≥ total)   │
  ▼                        │
paid_in_full               │
  │                        │
  │ (departure begins)     │
  ▼                        │
departed                   │
  │                        │
  │ (return + alumni move) │
  ▼                        │
completed                  │
                           │
cancelled ◄────────────────┘
  ↑
  └─────────── (from any pre-departure state, via refund flow)
```

`expired` is **a terminal state distinct from `cancelled`** — it means the VA timed out without payment; different from a customer-initiated cancel (no penalty logic runs). Q010 pins the TTL.

Transitions emit events (see §"Events").

### W5 — Smart room/bus allocation (Modules #92, #93, #95)

**Trigger timing is Q015.** Per PRD, Smart Grouping is under Section E (Operational & Handling), not Section A booking — it runs AFTER docs are verified and close to departure, not at submit time. At submit, F4 only **reserves capacity** (decrements catalog seat count, locks room-type count). Specific bed / seat / bus assignments come later via `ops-svc.RunSmartGrouping` which reads bookings + jamaah + family graph and produces:

- `room_allocations` (which jamaah share which hotel room, respecting K-Family Code from F3)
- `bus_allocations` (seat assignments per departure bus)
- `luggage_tags` with QR codes (Module #95)
- `id_cards` with QR codes

_(Inferred)_ F4 **stores** allocations once `ops-svc` computes them, but does not compute them. The allocation tables live in `ops-svc`, referenced from booking via ID only — this keeps booking lean and separates the algorithm from the data. Override authority: Q015.

### W6 — Booking-time mahram check (F3 integration)

At submit, before committing the saga, F4 calls `jamaah-svc.ValidateMahram(subject, group, departure)` for any female jamaah on the booking who falls under the mahram age threshold (Q005 pins the threshold and qualifying relations). Result is recorded on the booking as `mahram_check_result` and does NOT block submission today — it produces a **warning flag** that Ops resolves before visa submission. Blocking at booking would force chicken-and-egg flows (jamaah needs to book mahram before they can book themselves). Blocking at **visa submission** (F6) is the right gate.

### W7 — Cancellation (Alur Logika 6.5, Module #100)

1. Cancellation initiated by jamaah, CS, or agent.
2. F4 validates cancellation is allowed (booking status is pre-departure; departure hasn't happened).
3. F4 creates a `cancellation_request` record, reason code, requested-by.
4. Q004 decides when seats return — default (recommended in Q004) is **immediate** return to catalog pool upon status → `cancelled`.
5. F5 `StartRefundFlow` kicks off; refund amount calculated per Q012 (penalty matrix).
6. Booking status → `cancelled`. Downstream services receive `booking.cancelled` event and compensate (logistics stops kit prep; finance books refund journal; ops removes from manifest; crm reverses commission).
7. **Partial cancellation** (one jamaah out of many) — Q014 decides the policy. Default recommended: proportional shrink of the existing booking with per-jamaah refund calc, rather than cancel+rebook.
8. **Reopening a cancelled booking** — Q014. Default: forbidden; customer makes a fresh booking.

### W8 — Submit-time validation checklist

Before F4 creates the booking saga, it validates:

- [ ] All named jamaah exist in F3 (or are registered on-the-spot via the Guest Data Form path).
- [ ] Catalog package is `active` and departure is `open`.
- [ ] Seat availability ≥ requested pax (atomic check in the saga via F2 `ReserveSeats`).
- [ ] Minimum docs are uploaded per Q006 (default: KTP + passport per jamaah, OCR pending OK).
- [ ] If minors are on the booking, guardian linkage is present (Q016).
- [ ] Mahram check runs (W6) — warning, not blocker.
- [ ] Package pricing + add-ons total computed; total amount locked on the booking record.

If any HARD check fails (seats, active status, package validity), the submit is rejected with a specific reason code. If any SOFT check produces a warning (mahram, document completeness), the booking proceeds but with warning flags surfaced to the customer and the downstream ops queue.

## Acceptance criteria

- A booking goes through the status machine in W4 without skips or backward transitions (except the explicit `* → cancelled` arrow).
- Every booking carries `channel` (`b2c_self` | `b2b_agent` | `cs`) and, if applicable, `agent_id` and/or `cs_user_id`.
- Seat reservation is atomic under concurrent submits — verified by a k6 race test against F2 `ReserveSeats`.
- `ValidateMahram` (F3 W6) result is recorded on every booking that has female jamaah below the threshold.
- Minimum-docs validation runs at submit per Q006; missing docs produce a blocking error with the jamaah+kind that's missing.
- Every state-changing call writes to `iam.audit_logs` via F1 `RecordAudit`.
- Cancellation triggers `booking.cancelled`, which reaches catalog (seat return per Q004), F5 refund flow, F8 logistics, F9 finance, F10 crm, and F6 visa (all idempotent consumers).
- Partial cancellation (Q014) is a single transaction that reduces pax count, writes per-jamaah refund lines, and updates room/bus allocations if allocations already ran.
- Smart Grouping output is stored against the booking but the booking remains the authoritative source of WHO is on the trip; allocations are mutable, booking members are not (except via cancellation).
- Reopening a cancelled booking is rejected (Q014 default); customer must start a new booking.
- VA expiry on an unpaid booking transitions to `expired` (terminal), releases seats, and does NOT trigger the refund flow (no money was taken).

## Edge cases & error paths

- **Concurrent booking on the last seat.** Two users click Submit simultaneously when `available_seats = 1`. F2's atomic `ReserveSeats(n)` returns zero rows for the losing side; F4 fails fast with `apperrors.ErrConflict` and returns a "seat was just taken" message. Winning side proceeds normally. Verified by a k6 race scenario.
- **VA issuance fails after seat reserved.** Saga compensation: F5 returns error → F4 `ReleaseSeats(n)` to catalog → booking marked `failed` (not `cancelled`, not `expired` — this is a system fault, distinct state for ops triage). _(Inferred — `failed` is not in PRD enum; add it.)_
- **Webhook arrives for a booking already in `paid_in_full`.** Idempotent no-op — F5 deduplicates on `gateway_txn_id`. F4 simply returns current state.
- **VA expires, then customer tries to pay later.** Payment is rejected by gateway (VA expired). If we want to allow grace, commercial choice; default: customer creates a new booking.
- **Abandoned checkout re-engaged by CS** (Q019). If CS rescues a B2C-abandoned cart, attribution policy is Q019. Default: the last-touched channel wins (CS takes over).
- **Departure is cancelled mid-flight** (airline cancels, force majeure). All bookings on that departure need a bulk cancellation path. _(Inferred)_ Ops triggers a bulk operation; booking-svc cancels each and F5 refunds each with a special reason code `departure_cancelled` that bypasses the penalty matrix (Q012).
- **Mahram warning resolved mid-booking** (a male family member books into the same departure after a warning fired). F4 doesn't auto-re-check every booking when a new booking lands; instead, the mahram warning is re-evaluated at F6 visa submission (the right gate). _(Inferred.)_
- **Passport 6-month rule violation detected at booking submit** (F3 passport expiry too close to departure). HARD blocker — booking submission fails with reject reason `passport_expires_before_departure`, pointing the jamaah to renew before booking.
- **Multi-currency invoicing** (Q001) — locked at VA issuance time with FX snapshot on the booking record. Subsequent FX changes don't reprice paid bookings.
- **Paket Tabungan interaction** (Q017). Does an active savings plan hold a seat? _(Inferred default: no — savings is pre-booking, converts to a real booking when the target is hit.)_ Q017 pins this.

## Data & state implications

Owned by `booking-svc`. Full schema in `docs/03-services/02-booking-svc/02-data-model.md`. Key additions from this spec:

- `bookings.channel` enum — `b2c_self | b2b_agent | cs`.
- `bookings.status` enum — `draft | pending_payment | partially_paid | paid_in_full | departed | completed | cancelled | expired | failed`.
- `bookings.mahram_check_result` jsonb — snapshot of F3 ValidateMahram response at submit.
- `bookings.fx_snapshot` jsonb — effective rates at VA issuance time (for invoice reconciliation).
- `cancellation_requests` table — one per cancellation attempt (including partial).
- `booking_items.status` enum — `active | cancelled` (supports partial cancel Q014).
- Allocations (rooms/buses) are **not** in booking-svc — they're in ops-svc, linked by `booking_id`.

## API surface (high-level)

Full contracts in `docs/03-services/02-booking-svc/01-api.md`. Key surfaces:

**REST:**
- `GET|POST|PATCH /v1/bookings` — list, create draft, edit draft
- `POST /v1/bookings/{id}/submit` — run the saga (atomic: reserve seats + issue VA + emit events)
- `POST /v1/bookings/{id}/cancel` — initiate cancellation (amount calc via F5 Q012)
- `POST /v1/bookings/{id}/partial-cancel` — remove one or more jamaah (Q014)
- `GET /v1/bookings/{id}/status` — current state + history
- `GET /v1/departures/{id}/bookings` — manifest input (ops-svc reads this)

**gRPC (service-to-service):**
- `MarkBookingPaid` — called by F5 on webhook settlement
- `AttachVisa` — called by F6 on e-visa issued
- `CancelBooking` — called by broker-svc saga compensations
- `ListBookingsForDeparture` — used by ops-svc for manifest + grouping
- `GetBooking`, `BatchGetBooking` — read-only fan-out

## Dependencies

- **F1** — identity, permissions, audit
- **F2** — package + seat reservation (`ReserveSeats`, `ReleaseSeats`)
- **F3** — jamaah profile, mahram validation, minimum-docs gate
- **F5** — VA issuance, webhook, refund flow (interlocked saga)
- **F10** — commission attribution (via `agent_id` on booking)

Downstream consumers (events, listed below):
- F5 (webhook → MarkBookingPaid reverse direction)
- F6 (visa pipeline unlocks on paid_in_full)
- F7 (manifest + grouping)
- F8 (logistics kit dispatch on paid_in_full)
- F9 (finance — Unearned Revenue journal on DP, revenue recognition on departed)
- F10 (crm — commission state transitions)

## Backend notes

- **The saga lives in `broker-svc`**, not here. `booking-svc` exposes gRPC methods that the saga activities call; the saga controls order + compensation. See `docs/03-services/10-broker-svc/00-overview.md`.
- **Status machine** is enforced in the service layer with explicit `can_transition_to(from, to)` checks; DB does not enforce (Postgres enums can't express "this transition is allowed").
- **Fx snapshot** is set by `IssueVirtualAccount` in F5, then copied onto the booking record so the booking self-contains enough to reason about invoice amount without cross-service lookup.
- **Soft delete is not used on bookings** — cancelled is a terminal state, not a delete. Audit trail depends on booking rows being immutable beyond state transitions.
- **K-Family Code** (from F3 `jamaah.family_unit.code`) propagates into smart grouping via ops-svc read; F4 stores nothing about rooming.

## Frontend notes

- Three booking UIs, one data model: B2C checkout, B2B agent closing form, internal CS closing form. Share the same `/v1/bookings` backend endpoint with different `channel` defaults.
- Checkout page shows real-time seat tracker label from F2; if inventory drops to 0 mid-session, UI surfaces immediately and offers a date-change suggestion.
- Customer Portal billing dashboard renders from `GET /v1/bookings/{id}/status` + `GET /v1/invoices?booking_id=...` (F5).
- Cancellation UI must surface the penalty calc BEFORE commit (jamaah sees "you will receive Rp 14,200,000 back after Rp 3,800,000 penalty" and confirms) — tight coupling with F5 refund preview endpoint.

## Open questions

See `docs/07-open-questions/`:

**Existing (shared with other features):**
- **Q004** — cancellation → seat return ownership
- **Q005** — mahram qualifying relations, age threshold
- **Q006** — minimum docs required to submit a booking

**New, filed with this draft:**
- **Q010** — VA TTL and draft booking expiry (how long customer has to pay before seat returns)
- **Q014** — partial cancellation + reopening cancelled bookings
- **Q015** — Smart Grouping trigger timing + override authority
- **Q016** — booking on behalf of a minor (no KTP)
- **Q017** — Paket Tabungan interaction with booking state
- **Q019** — abandoned checkout attribution (who gets commission)
