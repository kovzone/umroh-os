# ADR 0005 — Svelte 5 (runes mode) + Vite for the frontend

**Status:** Accepted
**Date:** 2026-04-15

## Context

Earlier docs (ADR-0001, `01-tech-stack.md`, `00-system-overview.md`) assumed the frontend would be **React + Vite**. That was a placeholder — the original PRD didn't specify a frontend framework, and "React" was the industry-default fallback. No frontend code had been written yet.

The other developer, who is picking up the frontend slice, has brought in a Svelte 5 best-practices skill (`.claude/skills/svelte-core-bestpractices/`) and confirmed the frontend will be built in **Svelte 5 with runes mode**. This ADR records the switch so the earlier assumption doesn't keep propagating.

## Decision

Build the frontend in **Svelte 5 (runes mode) + Vite**. SvelteKit vs. plain Svelte+Vite is still open (see `docs/07-open-questions/`).

Conventions live in `docs/05-frontend-conventions/`, mirroring the layout of `docs/04-backend-conventions/`. The canonical runtime guide is the `svelte-core-bestpractices` skill, which is loaded automatically in Svelte-flavored sessions.

## Rationale

1. **Developer preference.** The developer taking the frontend slice has chosen Svelte and brought a ready-to-use best-practices skill. Matching the tool to the operator's velocity outweighs any abstract framework comparison.
2. **Svelte 5 maturity.** Runes (`$state`, `$derived`, `$effect`, `$props`) give fine-grained reactivity without the `useMemo`/`useCallback` discipline React demands. `{@attach}` (5.29+) handles third-party DOM libraries cleanly. `createContext` (5.40+) gives type-safe context without string keys. The skill covers all of these and matches current svelte.dev guidance.
3. **Smaller output.** Svelte compiles to surgical DOM updates with no virtual DOM runtime — a material win for the multi-page admin/jamaah/CRM surface UmrohOS needs to render on modest hardware.
4. **No pre-existing React code.** Switching costs are zero; nothing has been written that depends on React semantics.
5. **Template-mirror pattern holds.** The backend has `baseline/go-backend-template/` + `docs/04-backend-conventions/` as the guardrail pair. The frontend now has `.claude/skills/svelte-core-bestpractices/` + `docs/05-frontend-conventions/` — same shape, same intent.

## Consequences

- All frontend code is **Svelte 5 runes mode**. Legacy Svelte 3/4 features (`$:`, `export let`, `<slot>`, stores, `use:action`, `on:click`) are banned — see `docs/05-frontend-conventions/00-coding-style.md`.
- The Go gateway's OpenAPI still defines the API contract. How the Svelte side consumes it (manual, openapi-ts, orval, hand-written client) is an open question.
- Testing framework (Vitest? Playwright?), file layout (SvelteKit vs. plain Svelte+Vite), and component/file naming are still open — tracked in `docs/07-open-questions/`.
- The `docs/05-frontend-conventions/` directory is the shared middle layer both developers read before touching any `.svelte` file.

## Alternatives considered

- **React + Vite** — originally assumed. Rejected because (a) the operator taking this codebase prefers Svelte, (b) no React code exists yet so switching cost is zero, and (c) Svelte 5's reactivity model removes an entire class of React footguns (stale closures, needless re-renders, `useEffect` dependency arrays) without adding a runtime.
- **Vue 3 + Vite** — viable but not proposed. The developer didn't bring Vue tooling.
- **SolidJS** — closer to Svelte 5 in spirit, but smaller ecosystem and no skill/template already prepared.
