# core-web

UmrohOS front-gate web app. Operator console — today shows live status of every backend; future home for admin (users / roles / branches / audit) and executive dashboard surfaces.

Built on **SvelteKit + TypeScript + Svelte 5 (runes mode)** per ADR 0005.

## Layout

```
src/
├── routes/            # SvelteKit file-based routing
│   ├── +layout.svelte # nav shell
│   ├── +page.svelte   # status page (CSR — polls every 5s)
│   └── +page.ts       # ssr=false, prerender=false
└── lib/
    ├── components/    # reusable PascalCase.svelte
    ├── state/         # *.svelte.ts classes with $state (no stores)
    └── api/gateway/   # typed client; schema.d.ts is generated
```

## Local quickstart

Prereqs: Node 20+ and the dev stack running (`make dev-bootstrap` from repo root).

```bash
cp .env.example .env       # idempotent; uses VITE_GATEWAY_URL=http://localhost:4000
npm install
npm run gen:api            # regenerate src/lib/api/gateway/schema.d.ts
npm run dev                # http://localhost:3000 (or 3001 via docker-compose)
```

## Scripts

- `npm run dev`     — Vite dev server, hot reload, host 0.0.0.0 port 3000.
- `npm run build`   — SvelteKit `adapter-node` production build into `build/`.
- `npm run preview` — serve the production build locally.
- `npm run check`   — `svelte-check` over the project.
- `npm run test`    — Vitest unit + component tests (jsdom).
- `npm run gen:api` — regenerate gateway types from `services/gateway-svc/api/rest_oapi/openapi.yaml`.

## Tests

- Unit / component: `npm run test` (Vitest).
- E2E (browser, against the running stack): from repo root `make e2e` — Playwright project `browser` covers `core-web`.

## Conventions

See `docs/05-frontend-conventions/` and the `svelte-core-bestpractices` skill. Highlights enforced here:

- Runes mode only — no `$:`, no `export let`, no stores.
- Class-with-`$state` for shared state (see `src/lib/state/service-status.svelte.ts`).
- `createSubscriber` for external event sources (the 5s poll loop).
- `clsx`-array `class` attribute, `onclick={...}` event handlers, keyed `{#each ... (item.id)}`.
- Components: `PascalCase.svelte`. Modules: `kebab-case.ts`.

## Out of scope today

Auth, multi-app split, design system, Sentry/RUM, production deploy. See `docs/03-services/12-core-web/04-status.md` for the full deferred list.
