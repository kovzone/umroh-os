# Logging & Tracing

UmrohOS uses **zerolog** for structured JSON logging and **OpenTelemetry** for distributed tracing. Every log line is correlated with the active trace via `trace_id`.

## Logging

### Setup

`util/logging.NewLogger` is initialized in `cmd/start.go` and stored on the service struct.

### The golden rule

**Every log call inside a request handler or service method must use `LogWithTrace`:**

```go
logger := logging.LogWithTrace(ctx, s.logger)
logger.Info().
    Str("email", params.Email).
    Msg("creating user")
```

This injects the current trace_id and span_id into every log line. Without it, logs cannot be correlated with traces in Grafana — and that's the whole point of having both.

### Log levels

| Level | When to use |
|---|---|
| `Debug` | Verbose flow tracing. Off in prod. Use for "I called this branch" details. |
| `Info` | Service-method entry/exit. Key state changes. Lifecycle events (server start). |
| `Warn` | Recoverable problems. Retries. Validation failures from clients. |
| `Error` | Failures returned as errors. Always log at the *catching* level, not on every wrap. |
| `Fatal` | Startup failures only. Exits the process. |

### What to log

- **Service-method entry:** `logger.Info().Interface("params", params).Msg("CreateUser called")`
- **Service-method success:** `logger.Info().Str("user_id", result.UserID).Msg("CreateUser succeeded")`
- **Service-method error:** `logger.Error().Err(err).Msg("CreateUser failed")`
- **External call:** `logger.Info().Str("upstream", "midtrans").Str("op", "create_va").Msg(...)`

### What NOT to log

- **PII.** No full passport numbers. No full payment account numbers. No raw OCR contents. Log IDs, not values.
- **Bearer tokens, API keys, passwords.** Ever.
- **Repeated heartbeats.** Use metrics, not logs.
- **Successful debug-level details in prod.** Set the level appropriately.

## Tracing

### Setup

`util/tracing.InitTracer` initializes the OTel tracer in `cmd/start.go`. Exporter target is configured via the OTel collector address in `config.json`.

### The pattern

Every service method opens a span:

```go
func (s *Service) CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error) {
    ctx, span := tracing.GetTracer().Start(ctx, "Service.CreateUser")
    defer span.End()

    logger := logging.LogWithTrace(ctx, s.logger)
    logger.Info().Interface("params", params).Msg("CreateUser called")

    // ... do work ...

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return CreateUserResult{}, err
    }

    span.SetStatus(codes.Ok, "")
    return result, nil
}
```

### Span attributes

- Name spans `<Layer>.<Method>` — e.g. `Service.CreateUser`, `Store.GetUserByID`.
- Set attributes for key inputs:
  ```go
  span.SetAttributes(attribute.String("email", params.Email))
  ```
- Don't put PII in attributes. Same rules as logs.

### Cross-service propagation

- gRPC propagates trace context automatically via OTel interceptors (configured in the template).
- HTTP requests entering via Fiber middleware extract the W3C `traceparent` header.
- (Reserved — Temporal trace propagation via the OTel interceptor returns in F6 when Temporal is reintroduced per ADR 0006.)

## Verifying the wiring

After scaffolding any service:
1. Send a smoke request.
2. Find the `trace_id` in the response or in stdout.
3. Open Grafana → Tempo → search by trace ID.
4. Confirm the span tree shows the request through all layers.
5. Click "logs for this trace" in Grafana → confirm log lines appear in Loki tagged with the same `trace_id`.

If any of those steps fail, the service is not properly instrumented — fix it before moving on.
