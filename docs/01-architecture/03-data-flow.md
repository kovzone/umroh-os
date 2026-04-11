# Data Flow — Key End-to-End Scenarios

Three flows that touch most services. Use these as the canonical examples when designing new features.

## Flow 1 — Booking creation & payment

```
Calon jamaah                Frontend            gateway-svc       booking-svc       catalog-svc       payment-svc       broker-svc
     │                         │                     │                 │                 │                 │                 │
     │  browse & pick package  │                     │                 │                 │                 │                 │
     ├────────────────────────►│  GET /v1/packages   │                 │                 │                 │                 │
     │                         ├────────────────────►│ gRPC ListPkg    │                 │                 │                 │
     │                         │                     ├────────────────────────────────►│                 │                 │
     │                         │                     │◄────────────────────────────────┤                 │                 │
     │                         │◄────────────────────┤                 │                 │                 │                 │
     │  click "book"           │                     │                 │                 │                 │                 │
     ├────────────────────────►│  POST /v1/bookings  │                 │                 │                 │                 │
     │                         ├────────────────────►│ gRPC StartBookingSaga             │                 │                 │
     │                         │                     ├──────────────────────────────────────────────────────────────────────►│
     │                         │                     │                 │                 │                 │  Temporal:      │
     │                         │                     │                 │                 │                 │  reserve seat   │
     │                         │                     │                 │                 │◄────────────────┤  (catalog)      │
     │                         │                     │                 │                 │                 │  create booking │
     │                         │                     │                 │◄────────────────────────────────────┤  (booking)     │
     │                         │                     │                 │                 │                 │  issue VA       │
     │                         │                     │                 │                 │                 │  (payment)──────┤
     │                         │◄──────────────────────────────────────────────────────────────────────────────  return VA   │
     │  pay to VA              │                     │                 │                 │                 │                 │
     ├────────────────────────►│                     │                 │                 │                 │                 │
     │                         │                     │                 │                 │ webhook         │                 │
     │                         │                     │                 │                 │◄────────────────┤  (gateway)      │
     │                         │                     │                 │                 │  signal saga    │                 │
     │                         │                     │                 │                 ├────────────────►│                 │
     │                         │                     │                 │  mark paid      │                 │                 │
     │                         │                     │                 │◄────────────────────────────────────┤               │
     │                         │                     │                 │  emit fulfilment trigger            │               │
     │                         │                     │                 │────────────────────────────────────►│ → logistics   │
```

Key points:
- The booking saga is a single Temporal workflow. If any step fails, compensating activities undo the prior steps.
- The HTTP request from the frontend returns as soon as the saga is started. The frontend polls or uses a websocket for status.
- Each service owns its own tables. No service writes to another service's DB.

## Flow 2 — Visa pipeline

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
- The visa workflow may run for days. Temporal handles the long durability and retries.
- The polling is on the activity side; the workflow is otherwise idle.
- `visa-svc` is the only service that calls MOFA/Sajil — adapter pattern hides the protocol.

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
