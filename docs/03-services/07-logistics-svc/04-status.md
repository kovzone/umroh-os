# logistics-svc — Status

## Implementation checklist

- [ ] Scaffolded
- [ ] Wired into compose
- [x] **ADR 0009 refactor** — `api/rest_oapi/` removed; `util/monitoring` rewritten from Prometheus SDK to OTel SDK (OTLP metrics push to the collector's Prometheus exporter); `util/monitoring/panic.go` dropped; REST port `4008` retired from compose (`BL-REFACTOR-005` / S1-E-13, 2026-04-23)
- [ ] DDL
- [ ] sqlc queries
- [ ] OpenAPI spec
- [ ] Stock CRUD
- [ ] PR/PO workflow with multi-level approval
- [ ] GRN with QC
- [ ] Kit definition CRUD
- [ ] Kit dispatch on booking paid (gRPC)
- [ ] Courier adapter (label + tracking)
- [ ] Unit tests
- [ ] Integration tests
- [ ] Verified by reviewer

## Current status

**Scaffolded, gRPC-only.** Per ADR 0009 logistics-svc exposes gRPC on `50058` only; the REST package (`api/rest_oapi/`) + scaffold endpoints + `/metrics` HTTP endpoint are gone. `util/monitoring` pushes gauges to the OpenTelemetry Collector via OTLP gRPC. The gateway no longer carries a `logistics_rest_adapter` or `/v1/logistics/system/live` route. Real business RPCs pending.
