# booking-svc — Status

## Implementation checklist

- [ ] Scaffolded from baseline template
- [ ] Wired into `docker-compose.dev.yml`
- [ ] Initial DDL written (`_init/booking_db/`)
- [ ] sqlc queries for bookings, items, allocations
- [ ] OpenAPI spec
- [ ] Booking CRUD endpoints
- [ ] Submit endpoint that hands off to saga
- [ ] gRPC `MarkBookingPaid`, `AttachVisa`, `CancelBooking`, `ListBookingsForDeparture`
- [ ] Mahram validation at submit time (calls jamaah-svc)
- [ ] Status machine enforcement
- [ ] Unit tests
- [ ] Integration tests
- [ ] Verified by reviewer

## Current status

**Not started.**
