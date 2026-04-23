# payment-svc — Status

## Implementation checklist

- [ ] Scaffolded from baseline template
- [ ] Wired into `docker-compose.dev.yml`
- [x] **ADR 0009 refactor** — `api/rest_oapi/` removed; `util/monitoring` rewritten from Prometheus SDK to OTel SDK (OTLP metrics push to the collector's Prometheus exporter); `util/monitoring/panic.go` dropped; REST port `4005` retired from compose (`BL-REFACTOR-007` / S1-E-13, 2026-04-23)
- [ ] Initial DDL written (`_init/payment_db/`)
- [ ] sqlc queries
- [ ] OpenAPI spec
- [ ] Invoice CRUD
- [ ] Midtrans adapter (issue VA, verify webhook, refund)
- [ ] Xendit adapter (same surface)
- [ ] Webhook handlers with signature verification
- [ ] Idempotent payment_event ingestion
- [ ] Reconciliation cron job
- [ ] gRPC methods for saga
- [ ] Unit tests
- [ ] Integration tests with mocked gateway
- [ ] Verified by reviewer

## Current status

**Scaffolded, gRPC-only.** Per ADR 0009 payment-svc exposes gRPC on `50055` only; the REST package (`api/rest_oapi/`) + scaffold endpoints + `/metrics` HTTP endpoint are gone. `util/monitoring` pushes gauges to the OpenTelemetry Collector via OTLP gRPC. The gateway no longer carries a `payment_rest_adapter` or `/v1/payment/system/live` route. Real business RPCs pending.
