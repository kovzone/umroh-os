# finance-svc — Status

## Implementation checklist

- [ ] Scaffolded
- [ ] Wired into compose
- [ ] **ADR 0009 realignment** — REST + local bearer-auth middleware + local `iam_grpc_adapter` removed; `/v1/finance/ping` migrated to gateway REST with `RequirePermission` middleware; `FinancePing` gRPC RPC added; `02b`/`02c` e2e specs routed through `gateway-svc:4000` (`BL-IAM-019` / S1-E-14 — fenced off G9 sweep to keep BL-IAM-002 coverage intact)
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

**Not started.** This is the most complex domain — likely the last MVP service.
