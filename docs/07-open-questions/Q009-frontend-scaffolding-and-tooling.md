---
id: Q009
title: Frontend scaffolding, testing framework, file naming, and API client strategy
asked_by: session 2026-04-15 Svelte conventions extract
asked_date: 2026-04-15
blocks: all frontend feature work
status: open
---

# Q009 — Frontend scaffolding, testing, naming, and API client

## Context

ADR-0005 locks the frontend stack as **Svelte 5 (runes mode) + Vite**. Conventions for _how to write Svelte code_ are captured in `docs/05-frontend-conventions/` and the `svelte-core-bestpractices` skill. But several scaffolding-level decisions are not covered by that skill and need stakeholder input before the first frontend feature can land.

## The question

Six linked sub-questions:

1. **App architecture.** The module list identifies four frontend surfaces: B2C storefront, B2B agent portal, Executive Dashboard, and Jamaah field app. Are these **one SvelteKit app** with route-level separation, **one monorepo with multiple apps** (e.g. `apps/b2c`, `apps/b2b`, `apps/admin` + shared `packages/ui`), or **separate repos** (micro-frontend)? This is the highest-impact decision and gates the next five.
2. **SvelteKit or plain Svelte + Vite?** SvelteKit gives file-based routing, server endpoints, form actions, and SSR out of the box. Plain Svelte+Vite is lighter but we'd build routing/SSR ourselves (or skip SSR entirely). The Go gateway already handles API surface, so this is really a question about _how the frontend is structured and served_.
3. **Docs layout.** Backend services each get a folder under `docs/03-services/` (e.g. `00-iam-svc/`). Where do frontend app docs live? Options: (a) rename `docs/03-services/` → `docs/03-backend-services/` and create a sibling `docs/08-frontend-apps/`; (b) add frontend apps _inside_ `docs/03-services/` alongside the `-svc` folders; (c) keep docs lean — cover structure in `docs/05-frontend-conventions/` only, defer per-app docs until we scaffold. This partially depends on the answer to sub-question 1.
4. **Testing framework.** The skill is silent on testing. Industry defaults for Svelte 5 are **Vitest** (unit + component) + **Playwright** (E2E). Are we adopting those, something else, or deferring until we hit a bug that demands tests?
5. **Component and file naming.** The skill uses `PascalCase.svelte` in all examples (e.g. `Button.svelte`, `Table.svelte`) — that's idiomatic, but not stated as a rule. Are we locking `PascalCase.svelte` for components, `kebab-case.ts` for modules? Any exceptions?
6. **OpenAPI → Svelte client.** The backend defines every REST endpoint in `openapi.yaml` and generates Go server stubs via `oapi-codegen`. How does the Svelte side consume it? Options: hand-written `fetch` wrappers, `openapi-typescript` (types only), `openapi-fetch` (typed runtime), `orval`, or SvelteKit remote functions if we pick SvelteKit.

## Options considered

### Sub-question 1 — App architecture

- **Option A — Single SvelteKit app, route-level separation.** One codebase, one build, one deploy. Routes like `/`, `/agent`, `/admin`, `/dashboard` each serve their audience; auth middleware gates them. Pros: simplest to ship, maximum component reuse, one set of dependencies to upgrade. Cons: blast radius of a bad deploy is 100%; harder to tune bundle size per audience.
- **Option B — Monorepo with multiple apps.** `frontend/` contains `apps/b2c`, `apps/b2b`, `apps/admin`, `apps/field`, plus `packages/ui` and `packages/api-client`. Each app deploys independently. Pros: independent release cadence, per-audience bundle tuning, cleaner auth boundaries. Cons: more tooling (pnpm workspaces, turbo / nx), more CI complexity.
- **Option C — Separate repos (micro-frontend).** Each surface is its own repository. Pros: strongest isolation. Cons: hardest to share UI and API-client code, slowest to iterate across surfaces, hardest to keep design-system consistency.

### Sub-question 2 — SvelteKit vs. plain Svelte+Vite

- **Option A — SvelteKit.** File-based routing, form actions, SSR/CSR toggle per route, remote functions for typed client-server calls, first-class integration with `hydratable`. Pros: batteries-included, matches current Svelte best-practices momentum. Cons: opinionated, adds a deploy target (Node adapter, static adapter, etc.).
- **Option B — Plain Svelte + Vite.** Maximum flexibility, simplest build, deploy as static assets behind the Go gateway. Cons: we build routing (svelte-routing / tinro) and state-hydration ourselves; we lose remote functions.

### Sub-question 3 — Docs layout for frontend apps

- **Option A — Rename backend dir, add sibling frontend dir.** `docs/03-services/` → `docs/03-backend-services/`; new `docs/08-frontend-apps/` holds one folder per frontend surface (matching whatever architecture sub-question 1 lands on). Honest naming; no mixing of concerns. Cost: rename touches every cross-link.
- **Option B — Keep `03-services`, put frontend apps inside it.** E.g. add `docs/03-services/11-b2c-web/` alongside `00-iam-svc/`. Cheap today; `03-services` becomes heterogeneous (microservices + frontend apps).
- **Option C — No per-app docs yet.** Keep `docs/05-frontend-conventions/` as the only frontend-doc surface. Add per-app pages only when a specific app has enough surface area to document. Cheapest; honest that we don't know yet.

### Sub-question 4 — Testing

- **Option A — Vitest + Playwright.** Vitest is the Vite-native test runner; Playwright is the reference E2E tool. This is the path of least resistance.
- **Option B — Vitest only, E2E deferred.** Start with unit tests, add Playwright when the app has enough user flows to justify it.
- **Option C — No tests until a bug.** Matches how the backend has started (no test discipline baked into the template beyond compile-time checks). High risk as the frontend grows.

### Sub-question 5 — Naming

- **Option A — `PascalCase.svelte` for components, `kebab-case.ts` for modules, `camelCase` for exports.** Idiomatic Svelte community default; matches the skill's examples.
- **Option B — `kebab-case.svelte` everywhere.** Unusual in Svelte land; mainly seen in web-component projects.

### Sub-question 6 — OpenAPI consumption

- **Option A — `openapi-fetch` + `openapi-typescript`.** Generates types from `openapi.yaml` and a tiny typed `fetch` wrapper. Zero runtime overhead beyond `fetch`. Fits both SvelteKit and plain Svelte.
- **Option B — `orval`.** Generates typed client functions + TanStack Query hooks. Heavier runtime; the query-hook shape is React-flavored.
- **Option C — Hand-written client.** Zero dependencies; duplicates type information we already generate on the Go side.
- **Option D — SvelteKit remote functions.** Only viable if Option A under sub-question 1 is chosen. Bypasses OpenAPI entirely for internal calls — the server-side code lives alongside the component and types flow automatically.

## Recommendation

Sub-questions 1 and 3 are intertwined and should be decided together — the other four can fall out of that decision.

1. **App architecture: Option B — monorepo with multiple apps** (`apps/b2c`, `apps/b2b`, `apps/admin`, `apps/field` + `packages/ui`, `packages/api-client`). Gives independent deploy and blast-radius isolation without the overhead of separate repos. Worth the tooling cost once there's more than one audience, and UmrohOS has four.
2. **SvelteKit.** Each app in the monorepo is a SvelteKit app. File-based routing + form actions + remote functions is materially cheaper than reinventing routing, and the `hydratable` + async-Svelte story works best inside SvelteKit.
3. **Docs layout: Option A — rename `03-services` → `03-backend-services`, add sibling `08-frontend-apps/`** (one folder per app in `apps/`). Honest naming; keeps `03-backend-services` as a clean microservice surface; frontend apps get first-class docs on the same tier as backend services. Defer creating the dirs until we scaffold each app.
4. **Vitest for unit + component, Playwright for critical E2E flows (booking checkout, visa upload, payment confirm).** Don't chase coverage — test the handful of flows that would cause real money or compliance damage if they silently broke.
5. **`PascalCase.svelte` / `kebab-case.ts` / `camelCase` exports.** Codify as a lint rule once we pick a linter.
6. **`openapi-fetch` + `openapi-typescript`**, generated from the Go gateway's `openapi.yaml` via a Vite plugin or prebuild script, living in `packages/api-client`. Gives us end-to-end type safety with near-zero runtime cost. Use SvelteKit remote functions only for SvelteKit-internal endpoints (session, form actions) — not for cross-service calls.

Reversibility: the recommendations are mostly changeable later. The biggest commitments are (1) app architecture and (2) SvelteKit — route structure and form-action code would need rewriting if we later ripped SvelteKit out, and migrating from a monorepo to separate repos (or vice versa) is a one-time churn but a real one. The Vitest/Playwright/openapi-fetch/naming choices can be swapped independently. The docs layout can be revised without code changes.

## Answer

TBD — awaiting stakeholder discussion.
