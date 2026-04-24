# finance-svc — Status

## Implementation checklist

- [x] Scaffolded
- [x] Wired into compose
- [x] **ADR 0009 realignment** — REST + local bearer-auth middleware + local `iam_grpc_adapter` removed; `/v1/finance/ping` migrated to gateway REST with `RequirePermission` middleware; `FinancePing` gRPC RPC added; monitoring migrated from Prometheus-SDK to OTLP push (catalog G7 clone); `02b`/`02c` e2e specs routed through `gateway-svc:4000` (`BL-IAM-019` / S1-E-14)
- [ ] DDL with double-entry check
- [ ] Initial chart of accounts seed (PSAK-aligned)
- [ ] sqlc queries
- [ ] OpenAPI spec
- [ ] Manual journal entry endpoint
- [ ] Event consumers from payment / logistics / crm
- [ ] AR/AP aging reports
- [ ] Tax calculation (PPh / PPN)
- [ ] FX handling with daily rate snapshot
- [ ] Job-order cost rollup per departure
- [ ] Balance sheet / P&L / cash flow reports
- [ ] Unit tests (especially for double-entry invariant)
- [ ] Integration tests
- [ ] Verified by reviewer (PSAK compliance check)

## Current status

**Scaffold only, gRPC-only per ADR 0009** as of BL-IAM-019 / S1-E-14.
finance-svc exposes `FinancePing` (the permission-gate smoke surface
called by gateway's `/v1/finance/ping`) over gRPC and the standard
`grpc.health.v1.Health` container liveness. No REST, no direct
Prometheus scrape. Real finance work (double-entry schema, chart of
accounts, journal entries, AR/AP, tax, FX, reports) is the most complex
domain in the platform and is likely the last MVP service to start.
