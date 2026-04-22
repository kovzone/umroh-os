# Observability

UmrohOS uses the full OTel stack from day one. Logs, metrics, and traces all flow into Grafana so you can pivot between them with a single trace ID.

## The pipeline

```
Go service ──zerolog (stdout JSON)── Fluent-Bit ── Loki ──┐
           │                                              │
           │──OTLP gRPC── OpenTelemetry Collector ── Tempo─┤── Grafana
           │                                              │
           │──/metrics (admin HTTP)── Prometheus scraper ──┘
```

All three converge in Grafana. Loki, Tempo, and Prometheus are configured as data sources; Grafana correlates them via the `trace_id` field.

Per ADR 0009, downstream services are gRPC-only for business calls. The `/metrics` endpoint lives on a **minimal admin HTTP handler** on each backend — single handler, no OpenAPI spec, no business routes. Liveness / readiness probes use the **standard gRPC health protocol** (`grpc.health.v1.Health/Check`); docker-compose and Kubernetes health checks use the `grpc_health_probe` binary. `gateway-svc` keeps a full REST surface (it is the only external-facing service) and serves its own `/metrics` on the same REST port.

## Tracing

- Every service uses `util/tracing.InitTracer` at startup. See `baseline/go-backend-template/demo-svc/util/tracing/`.
- Every API method, every service method, and every store method opens a span.
- Span attributes record: operation name, key input parameters (no PII), error details on failure.
- Trace context propagates across gRPC calls automatically via OTel interceptors.
- (Reserved — Temporal workflow trace propagation returns in F6 when Temporal is reintroduced per ADR 0006.)

## Logging

- zerolog with JSON output to stdout.
- **Always** use `logging.LogWithTrace(ctx, logger)` to get a logger that includes the current `trace_id`. Never use a raw `logger` inside a request handler or service method.
- Standard fields on every log line: `level`, `time`, `service`, `trace_id`, `span_id`, `caller`.
- No PII in logs (no full passport numbers, no full payment data — log only IDs).
- Log levels:
  - `debug` — verbose, off in prod
  - `info` — service-method entry/exit, key state changes
  - `warn` — recoverable problems
  - `error` — failures that returned an error
  - `fatal` — startup failures only

## Metrics

- Prometheus scrapes each service's `/metrics` endpoint — on the gateway's REST port for `gateway-svc`, on a dedicated admin HTTP port for every downstream service (per ADR 0009).
- Standard metrics from the template: gateway HTTP request count/latency by route+status, gRPC request count/latency, DB pool stats, panic counter.
- Add custom business metrics via `util/monitoring`. Examples for UmrohOS to add over time:
  - `bookings_created_total`
  - `payments_received_total{provider}`
  - `visa_pipeline_duration_seconds`
  - `commission_calculated_total{level}`

## Dashboards

- Grafana home dashboard shows: per-service request rate, error rate, p95 latency, DB pool saturation.
- Per-service dashboards: scaffolded under `grafana/dashboards/<svc>.json` when the service is built.

## Reading a trace

1. Hit Grafana → Explore → Tempo.
2. Search by trace ID (from a log line) or by service+operation+latency filter.
3. Click a span to see its attributes (params, error).
4. Use the "logs for this span" link to jump to Loki and see all log lines tagged with that trace ID.

This is the workflow the reviewer should use when reviewing a session: pick a sample request, find its trace, walk it through every service, and confirm each log line is meaningful.

## What to verify on every new service

The `scaffold-service` skill ensures these are wired. Per ADR 0009, downstream services (everything except `gateway-svc`) are gRPC-only:

For a non-gateway service:
- [ ] `google.golang.org/grpc/health` standard health service registered; `grpc_health_probe -addr=:<grpc-port>` returns SERVING
- [ ] Minimal admin HTTP handler exposing `/metrics` on its configured admin port; Prometheus target is UP
- [ ] zerolog initialized with `LogWithTrace`
- [ ] OTel tracer initialized in `cmd/start.go`; `otelgrpc.StatsHandler` on the gRPC server
- [ ] Service appears in Grafana with at least one trace after a smoke gRPC request (e.g. from another service or via `grpcurl`)
- [ ] Logs for the smoke request appear in Loki tagged with the right service name

For `gateway-svc`:
- [ ] `/metrics`, `/livez`, `/readyz` on the REST port return 200
- [ ] zerolog, OTel REST middleware, and outbound `otelgrpc` client handlers wired
- [ ] Service appears in Grafana with at least one trace after a smoke REST request
