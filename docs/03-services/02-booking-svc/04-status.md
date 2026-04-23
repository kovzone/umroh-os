# booking-svc — Status

## Implementation checklist

- [ ] Scaffolded from baseline template
- [ ] Wired into `docker-compose.dev.yml`
- [x] **ADR 0009 refactor** — `api/rest_oapi/` removed; `util/monitoring` rewritten from Prometheus SDK to OTel SDK (OTLP metrics push to the collector's Prometheus exporter); `util/monitoring/panic.go` dropped; REST port `4003` retired from compose (`BL-REFACTOR-002` / S1-E-13, 2026-04-23)
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

**Scaffolded, gRPC-only.** Per ADR 0009 booking-svc exposes gRPC on `50053` only; the REST package (`api/rest_oapi/`), its scaffold endpoints (`/system/live`, `/system/ready`, `/system/diagnostics/db-tx`), and the `/metrics` HTTP endpoint are gone. `util/monitoring` pushes gauges + instruments to the OpenTelemetry Collector via OTLP gRPC (the collector re-exports them on its Prometheus endpoint; see catalog-svc G7 for the canonical pattern). The gateway no longer carries a `booking_rest_adapter` or `/v1/booking/system/live` route. DB pool metrics (`db_connections_{acquired,idle,total}`) keep flowing. Real business RPCs still pending.
