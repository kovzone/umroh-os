# core-web — Status

## Implementation checklist

- [~] Scaffolded as SvelteKit + TypeScript + Svelte 5 (runes mode) app at `apps/core-web/`
- [~] Wired into `docker-compose.dev.yml` with dev target + Vite hot reload on port `3001:3000`
- [~] Typed API client generated from gateway-svc openapi.yaml via `openapi-typescript` (run `npm run gen:api`)
- [~] Landing page (`/`) with hero + capability grid + footer; prerendered at build time
- [~] Status page (`/system/status`) polls all 10 gateway proxy endpoints every 5s via `createSubscriber`
- [~] Shared `Header` + `Footer` components inside `+layout.svelte`; `Sign in` placeholder disabled until F1.5
- [~] Vitest unit tests for `ServiceStatus` state class (4 tests, green)
- [~] Playwright browser e2e spec at `tests/e2e/tests/03-core-web-status.spec.ts` (4 tests covering landing + status + footer navigation)
- [~] Production Dockerfile `prod` stage with `@sveltejs/adapter-node` (path exists; no compose service uses it yet)
- [ ] Auth wiring (enable the Sign-in button, login form, token storage, protected routes) — F1.5 dependency
- [ ] Real admin routes (users/roles/branches/audit) — land with their feature slices
- [ ] Component library / design-system extraction — lands when second app forces shared UI
- [ ] Sentry / RUM / web-vitals instrumentation — separate observability concern
- [ ] Verified by reviewer in `testing-guide.md`

## Current status

**Scaffolded** — app lives at `apps/core-web/`. REST on port `3001` (host) / `3000` (container). Talks to `gateway-svc` via the typed client generated from its OpenAPI spec. No DB, no gRPC, no `/metrics` endpoint — frontend doesn't expose Prometheus.

**Decision recap (ADR 0005 + Q009 recommendations applied where open):**

- **SvelteKit + TypeScript + Svelte 5 runes mode** — routing, layouts, and future SSR/CSR toggle baked in from day one (not "add routing later").
- **`@sveltejs/adapter-node`** — production build path exists in the multi-stage Dockerfile (`prod` target); dev compose uses the `dev` target with Vite's hot-reload server.
- **Typed API client via `openapi-typescript` + `openapi-fetch`** — gateway's `openapi.yaml` is the single source of truth; `npm run gen:api` regenerates `src/lib/api/gateway/schema.d.ts`.
- **Class-with-`$state` for shared state** — `ServiceStatus` in `src/lib/state/service-status.svelte.ts` is the reference shape. No writable stores; no module-level mutable state.
- **`createSubscriber` for external event sources** — the 5s poll loop is owned by the class's subscriber, so the interval starts on first reactive read and stops when the page unmounts. No manual `$effect` + `setInterval` lifecycle.
- **Vitest (unit/component) + Playwright (browser e2e)** — Vitest has its own `vitest.config.ts` (avoids Vite type-version collision); Playwright is extended with a `browser` project that targets `http://localhost:3001` via `devices["Desktop Chrome"]`.
- **Naming:** `PascalCase.svelte` for components, `kebab-case.ts` for modules, `camelCase` exports — per Q009-rec #5.
- **Routing split — landing vs status.** `/` is a public-facing landing page (hero + capability grid + footer); the service-health dashboard lives at `/system/status` and is linked from the footer. Decided 2026-04-18 per reviewer's request: root path should represent UmrohOS itself, not operational tooling.

**Scaffold deliverables (this commit):**

- **`GET /` landing page** — product framing: hero (`UmrohOS` name + tagline + short description) + "What's inside" capability grid (4 cards: Booking & packages, Jamaah/documents/visa, Finance PSAK, Field operations) + footer with link to the status page. Prerendered at build time (`export const prerender = true`).
- **`GET /system/status` status page** — 10 cards, one per backend the gateway fronts (`iam, catalog, booking, jamaah, payment, visa, ops, logistics, finance, crm`). Each card polls `GET <VITE_GATEWAY_URL>/v1/<shortName>/system/live` every 5s via `openapi-fetch` (typed path argument). Card states: `pending` (grey, initial), `ok` (green), `fail` (red, with error message). CSR-only (`export const ssr = false`).
- **Shared shell** in `+layout.svelte` — `Header` (brand + disabled Sign-in button titled "Sign-in lands with F1.5 — iam-svc auth") + slot + `Footer` (status-page link + version line).
- Proves the end-to-end chain: browser → host `:3001` → SvelteKit → `openapi-fetch` bundle → gateway `:4000` → `*_rest_adapter` → backend `/system/live`.

**Layers laid out (future pages plug in):**

| Layer | Location | Today |
|---|---|---|
| Routes | `src/routes/+layout.svelte` + `src/routes/+page.svelte` + `src/routes/system/status/+page.svelte` | Landing (prerendered) + status grid (CSR). Nested sub-route structure established. |
| Components | `src/lib/components/*.svelte` (PascalCase) | `Header`, `Footer`, `CapabilityCard`, `ServiceStatusCard` |
| State | `src/lib/state/*.svelte.ts` (class with `$state`) | `service-status.svelte.ts` + `service-status.test.ts` |
| API client | `src/lib/api/gateway/{client.ts, schema.d.ts}` | `gateway` client (typed); regenerate via `npm run gen:api` |
| Design tokens | `src/app.css` (CSS custom properties) | Minimal dark theme; promote to shared when second app lands |

## Assigned ports

| Surface | Port (host → container) |
|---|---|
| REST (SvelteKit / Vite dev server) | `3001` → `3000` |

## Out of scope (explicit)

- **Auth.** No login form, no token handling, no protected routes. Arrives with F1.5 real iam auth.
- **Production compose service.** `Dockerfile prod` target exists so CI can build it; no compose service runs it yet.
- **Frontend observability.** No Sentry / RUM / web-vitals / frontend Grafana dashboard in this scaffold.
- **Multi-app layout.** Only `apps/core-web/`. Siblings (`apps/storefront-web/`, `apps/agent-web/`, `apps/field-web/`) are Q009-rec #1 future work.
- **Docs rename (Q009-rec #3).** `docs/03-services/12-core-web/` sits alongside the backend service folders for now; the rename to `08-frontend-apps/` is its own doc-only commit.
- **`.claude/skills/scaffold-app/`** — a future companion skill to `scaffold-service`, not scoped here.

## Verification

Run the manual walk in `docs/92-testing/testing-guide.md` Section 9. Automated:

```bash
make dev-down-v && make dev-rm-all && make dev-bootstrap   # 20 containers Up
make e2e-install                                            # idempotent
make e2e                                                    # 47 passed (43 api + 4 browser)
cd apps/core-web && npm run check                           # 0 errors
cd apps/core-web && npm run test                            # 4 passed
cd apps/core-web && npm run build                           # adapter-node output in build/ + prerendered landing
```
