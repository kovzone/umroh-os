# ADR 0008 — E2E testing with Playwright

**Status:** Accepted
**Date:** 2026-04-17

## Context

UmrohOS has 10+ microservices behind a shared observability stack (Prometheus, Grafana, Loki, Tempo). We need a way to:

1. **Verify service health** — confirm every scaffolded service's liveness and readiness probes respond correctly.
2. **Generate real traffic** — exercise the monitoring/Grafana dashboards with actual HTTP requests, not just synthetic checks.
3. **Test cross-service workflows** — as endpoints are built, verify multi-step business scenarios (create jamaah → book package → submit payment) end-to-end.

The team has prior experience with Playwright-based API testing in the ganesa project, where a scenario-driven, serial test framework proved effective for backend-only E2E tests against REST APIs.

## Decision

Use **Playwright** (API mode, no browser) for E2E testing. Tests live in `tests/e2e/` and run via `make e2e`.

### Structure

```
tests/e2e/
├── lib/
│   ├── api-client.ts     ← thin wrapper around Playwright's APIRequestContext
│   ├── services.ts       ← service registry (name + base URL per service)
│   ├── state.ts          ← JSON file for persisting IDs/tokens across ordered tests
│   └── types.ts          ← TypeScript interfaces for API response envelopes
├── tests/
│   ├── 01-system.spec.ts ← health probes across all services
│   ├── 02-iam-auth.spec.ts
│   └── ...               ← numbered for execution order
├── .env.example
├── package.json
├── playwright.config.ts
└── tsconfig.json
```

### Key design choices

- **Serial execution** (`workers: 1`, `fullyParallel: false`). Tests run in numbered order. Later tests depend on state from earlier tests (e.g., test 02 logs in and saves a token; test 03+ uses it). This trades parallelism for simplicity and determinism.
- **Multi-service base URLs.** Unlike ganesa (single `BASE_URL`), UmrohOS has many services on different ports. The `services.ts` registry maps each service to its URL. The `01-system.spec.ts` test loops over all registered services.
- **State file** (`.state.json`). Persists tokens, IDs, and other artifacts between test specs. Cleared at the start of each full run.
- **No browser.** Playwright is used purely for its `APIRequestContext` — HTTP client with built-in retry, timeout, and assertion helpers. No DOM rendering.

### Execution

```bash
make e2e-install   # npm install (first time / CI)
make e2e           # run all e2e tests
```

Tests require the dev environment to be running (`make dev-bootstrap` or `make dev-up`).

## Rationale

1. **Proven pattern.** The ganesa project uses this exact approach (Playwright API tests, serial execution, file-based state). The team is already familiar with it.
2. **Low overhead.** No test database fixtures, no Docker-in-Docker. Tests hit live services via HTTP — same as a real client.
3. **Traffic for observability.** Running `make e2e` generates real request traffic across all services, which exercises the Prometheus metrics, Grafana dashboards, Loki log aggregation, and Tempo traces. This is valuable even before business endpoints exist.
4. **Incremental growth.** Start with health probes (01-system). Add auth flows when iam-svc endpoints land. Add booking scenarios when booking-svc is built. Each new service just adds a line to `services.ts`.

## Why not k6?

k6 tests already exist at `tests/k6/` for load/stress testing. E2E and load testing serve different purposes:

- **E2E (Playwright):** functional correctness, scenario verification, sequential multi-step flows. "Does the system do the right thing?"
- **Load (k6):** performance under concurrency. "Does the system hold up under traffic?"

Both coexist. E2E runs first (verify correctness), then k6 hammers the verified endpoints.

## Consequences

- **New dev dependency.** `npm` and `@playwright/test` in `tests/e2e/`. Installed via `make e2e-install`.
- **Requires running services.** E2E tests are not unit tests — they need `make dev-up` or `make dev-bootstrap` first.
- **Service registry maintenance.** When a new service is scaffolded, add it to `tests/e2e/lib/services.ts` and `.env.example`.
- **Test ordering matters.** Numbered prefixes (`01-`, `02-`, ...) are load-bearing. New specs must be numbered to respect dependencies.
