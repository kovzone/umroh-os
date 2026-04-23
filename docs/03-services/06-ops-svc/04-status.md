# ops-svc — Status

## Implementation checklist

- [ ] Scaffolded
- [ ] Wired into compose
- [x] **ADR 0009 refactor** — `api/rest_oapi/` removed; `util/monitoring` rewritten from Prometheus SDK to OTel SDK (OTLP metrics push to the collector's Prometheus exporter); `util/monitoring/panic.go` dropped; REST port `4007` retired from compose (`BL-REFACTOR-006` / S1-E-13, 2026-04-23)
- [ ] DDL
- [ ] sqlc queries
- [ ] OpenAPI spec
- [ ] Verification queue endpoints
- [ ] Manifest generator (PDF/Excel)
- [ ] Smart room/bus allocation algorithm
- [ ] Luggage tag QR generation
- [ ] Handling event capture
- [ ] gRPC methods
- [ ] Unit tests (especially for grouping algorithm)
- [ ] Integration tests
- [ ] Verified by reviewer

## Current status

**Scaffolded, gRPC-only.** Per ADR 0009 ops-svc exposes gRPC on `50057` only; the REST package (`api/rest_oapi/`) + scaffold endpoints + `/metrics` HTTP endpoint are gone. `util/monitoring` pushes gauges to the OpenTelemetry Collector via OTLP gRPC. The gateway no longer carries an `ops_rest_adapter` or `/v1/ops/system/live` route. Real business RPCs pending.
