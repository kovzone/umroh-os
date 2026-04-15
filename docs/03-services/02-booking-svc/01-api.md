# booking-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/bookings` | List bookings (filterable by branch, status, agent) |
| `POST` | `/v1/bookings` | Create draft booking |
| `GET` | `/v1/bookings/{id}` | Get booking detail |
| `PATCH` | `/v1/bookings/{id}` | Update draft booking |
| `POST` | `/v1/bookings/{id}/submit` | Submit draft → start saga |
| `POST` | `/v1/bookings/{id}/cancel` | Cancel booking (triggers refund flow) |
| `GET` | `/v1/bookings/{id}/items` | List booking items (jamaah on booking) |
| `POST` | `/v1/bookings/{id}/items` | Add jamaah to booking |
| `DELETE` | `/v1/bookings/{id}/items/{item_id}` | Remove jamaah |
| `GET` | `/v1/bookings/{id}/room-allocations` | Get room layout |
| `PATCH` | `/v1/bookings/{id}/room-allocations` | Update room layout |
| `GET` | `/v1/bookings/{id}/bus-allocations` | Get bus seating |

## gRPC methods (planned)

`BookingService`:
- `GetBooking(...)`
- `ListBookingsForDeparture(...)` — used by ops-svc for manifest gen
- `MarkBookingPaid(...)` — called by payment-svc on settlement
- `AttachVisa(...)` — called by visa-svc when e-visa issued
- `CancelBooking(...)` — called by payment-svc refund saga compensation (in-process; ADR 0006)
- `CreateBookingItem(...)` — used by saga during atomic creation

> All endpoints are stubs.
