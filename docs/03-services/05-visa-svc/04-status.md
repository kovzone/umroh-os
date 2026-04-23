# visa-svc — Status

## Implementation checklist

- [ ] Scaffolded
- [ ] Wired into compose
- [x] **ADR 0009 refactor** — `api/rest_oapi/` removed; `util/monitoring` rewritten from Prometheus SDK to OTel SDK (OTLP metrics push to the collector's Prometheus exporter); `util/monitoring/panic.go` dropped; REST port `4006` retired from compose (`BL-REFACTOR-008` / S1-E-13, 2026-04-23)
- [ ] DDL
- [ ] sqlc queries
- [ ] OpenAPI spec
- [ ] MOFA/Sajil adapter
- [ ] Nusuk adapter (Raudhah Shield)
- [ ] Visa pipeline activities (in broker-svc)
- [ ] gRPC methods
- [ ] Unit tests
- [ ] Integration tests with mocked provider
- [ ] Verified by reviewer

## Current status

**Scaffolded, gRPC-only.** Per ADR 0009 visa-svc exposes gRPC on `50056` only; the REST package (`api/rest_oapi/`) + scaffold endpoints + `/metrics` HTTP endpoint are gone. `util/monitoring` pushes gauges to the OpenTelemetry Collector via OTLP gRPC. The gateway no longer carries a `visa_rest_adapter` or `/v1/visa/system/live` route. Real business RPCs pending.
