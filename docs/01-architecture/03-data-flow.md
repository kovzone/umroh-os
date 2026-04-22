# Data Flow — Key End-to-End Scenarios

Three flows that touch most services. Use these as the canonical examples when designing new features.

Per ADR 0009: the app ↔ gateway hop is REST; every gateway ↔ backend and backend ↔ backend hop is gRPC. Bearer validation happens once at the gateway via `iam-svc.ValidateToken` before any forwarded gRPC call leaves gateway; backends do not re-validate.

## Flow 1 — Booking creation & payment (in-process saga — ADR 0006)

```
Calon jamaah          Frontend          gateway-svc      booking-svc     catalog-svc     payment-svc
     │                    │                  │                │               │               │
     │  browse & pick     │                  │                │               │               │
     ├───────────────────►│  GET /v1/pkg     │                │               │               │
     │                    ├─────────────────►│ gRPC ListPkg   │               │               │
     │                    │                  ├──────────────────────────────►│               │
     │                    │◄─────────────────┤                │               │               │
     │  click "book"      │                  │                │               │               │
     ├───────────────────►│  POST /v1/bookings                │               │               │
     │                    ├─────────────────►│ gRPC Submit    │               │               │
     │                    │                  ├───────────────►│ (saga starts — booking-svc coordinates)
     │                    │                  │                │ WithTx: create draft booking  │
     │                    │                  │                ├──────────────►│ ReserveSeats  │
     │                    │                  │                │◄──────────────┤               │
     │                    │                  │                ├──────────────────────────────►│ IssueVA
     │                    │                  │                │◄──────────────────────────────┤
     │                    │                  │                │ WithTx: update status, history│
     │                    │                  │                │ emit booking.dp_received      │
     │                    │◄─────────────────┤ return VA details│             │               │
     │  pay to VA         │                  │                │               │               │
     ├───────────────────►│                  │                │               │               │
     │                    │                  │                │               │ webhook       │
     │                    │                  │                │               │◄──────────────┤ (gateway)
     │                    │                  │                │ MarkPaid (gRPC)│              │
     │                    │                  │                │◄──────────────┤               │
     │                    │                  │                │ emit events → F8/F9/F10       │
```

Key points:

- **No Temporal in MVP.** `booking-svc.Submit()` is the saga coordinator — it calls `catalog.ReserveSeats` then `payment.IssueVirtualAccount` in sequence, with explicit compensations on failure: seat fails → return 409; VA fails → call `catalog.ReleaseSeats` + delete draft → return error.
- **Durability against mid-saga crash** is handled by the F5 reconciliation cron — e.g. "invoice exists but no booking" is cleaned up on the next cycle.
- `WithTx` wraps the booking-svc local DB writes so they commit atomically within that service.
- `payment-svc` webhook handler calls `booking-svc.MarkBookingPaid` synchronously via gRPC — no Temporal signal, no broker indirection.
- On payment, booking-svc emits events (`booking.dp_received`, `booking.paid_in_full`) that logistics / finance / crm subscribe to for their downstream work.
- Temporal returns when F6 visa pipeline is implemented (Flow 2 below).

## Flow 2 — Visa pipeline (Temporal — F6; brought back for this feature per ADR 0006)

```
ops-svc            jamaah-svc          visa-svc          External (MOFA/Sajil)        broker-svc
   │                  │                   │                       │                        │
   │ verify documents │                   │                       │                        │
   ├─────────────────►│ docs.ready event  │                       │                        │
   │                  ├──────────────────────────────────────────────────────────────────►│
   │                  │                   │                       │                        │ workflow start
   │                  │                   │ create application    │                        │
   │                  │                   │◄──────────────────────────────────────────────┤
   │                  │                   │ submit                │                        │
   │                  │                   ├──────────────────────►│                        │
   │                  │                   │  poll status          │                        │
   │                  │                   ├──────────────────────►│                        │
   │                  │                   │  ISSUED               │                        │
   │                  │                   │◄──────────────────────┤                        │
   │                  │                   │ attach e-visa         │                        │
   │                  │                   ├──────────────────────────────────────────────►│
   │                  │                   │                       │                        │ booking.visa_attached
```

Key points:
- The visa workflow may run for days. Temporal handles the long durability and retries — **this is the use case that justifies bringing Temporal back (F6 implementation)**, per ADR 0006.
- The polling is on the activity side; the workflow is otherwise idle.
- `visa-svc` is the only service that calls MOFA/Sajil — adapter pattern hides the protocol.
- `broker-svc` is a deferred service in MVP; introducing it (+ Temporal in the compose stack) is part of F6's implementation scope.

## Flow 3 — Pre-departure manifest generation

```
ops-svc                booking-svc            jamaah-svc           catalog-svc
   │                       │                       │                     │
   │ build manifest        │                       │                     │
   ├──────────────────────►│ ListBookings(dep_id)  │                     │
   │◄──────────────────────┤                       │                     │
   │ for each jamaah:      │                       │                     │
   ├──────────────────────────────────────────────►│ GetJamaah(id)       │
   │◄──────────────────────────────────────────────┤                     │
   │ for each room/bus:    │                       │                     │
   ├──────────────────────────────────────────────────────────────────► │ GetHotel/GetFlight
   │◄────────────────────────────────────────────────────────────────── │
   │ run smart grouping algorithm (ops-svc local)                        │
   │ persist manifest                                                    │
```

Key points:
- The smart room/bus grouping algorithm is **business logic**, not data — it lives in `ops-svc/service/`.
- All data is fetched read-only via gRPC. `ops-svc` only writes to its own tables.
- The mahram constraint is checked at this step, calling `jamaah-svc` for the family graph.
