# jamaah-svc — Status

## Implementation checklist

- [ ] Scaffolded from baseline template
- [ ] Wired into `docker-compose.dev.yml`
- [x] **ADR 0009 refactor** — `api/rest_oapi/` removed; `util/monitoring` rewritten from Prometheus SDK to OTel SDK (OTLP metrics push to the collector's Prometheus exporter); `util/monitoring/panic.go` dropped; REST port `4004` retired from compose (`BL-REFACTOR-004` / S1-E-13, 2026-04-23)
- [ ] Initial DDL written (`_init/jamaah_db/`)
- [ ] sqlc queries for jamaah, family_units, documents
- [ ] OpenAPI spec
- [ ] Jamaah CRUD endpoints
- [ ] Document upload + GCS integration
- [ ] OCR pipeline (GCP Vision adapter)
- [ ] Mahram validation gRPC method (recursive CTE)
- [ ] gRPC `GetJamaah`, `BatchGetJamaah`, `GetPassportData`
- [ ] Unit tests
- [ ] Integration tests
- [ ] Verified by reviewer

## Current status

**Scaffolded, gRPC-only.** Per ADR 0009 jamaah-svc exposes gRPC on `50054` only; the REST package (`api/rest_oapi/`) + scaffold endpoints + `/metrics` HTTP endpoint are gone. `util/monitoring` pushes gauges to the OpenTelemetry Collector via OTLP gRPC. The gateway no longer carries a `jamaah_rest_adapter` or `/v1/jamaah/system/live` route. Real business RPCs pending.
