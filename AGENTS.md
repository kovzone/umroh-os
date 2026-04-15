# AGENTS.md — UmrohOS

Guidance for **human developers** and **coding agents** (Cursor, Claude Code, etc.). This file is intentionally short: it routes you to the right `docs/` paths instead of duplicating them. **Edit via PR** when team agreements change.

---

## What this project is

**UmrohOS** — ERP end-to-end for Umrah/Hajj travel agencies (booking, agents, documents/visa, logistics, finance, operations, CRM). See `README.md`.

---

## Working model (humans + AI, vendor-neutral)

- **Two full-stack developers**, each owning functional feature slices end-to-end across both the Go services and the Svelte 5 + Vite frontend. The split is per-workflow, not per-specialization.
- **Each developer may use a different AI coding agent** (Claude Code, Cursor, Copilot, or none) with a different working-file layout. Do not assume any specific agent's conventions. If your agent has its own instruction file (`CLAUDE.md`, `.cursor/`, etc.) it is **gitignored and private to that developer** — don't expect to find it, and don't commit your own.
- **`AGENTS.md` (this file) is the one shared, committed onboarding document** for humans and AI agents alike. Everything referenced below lives in `docs/` and is committed. If a rule in your private agent file contradicts this one, this one wins.
- **Private task trackers, session rituals, and agent-specific skills are scoped to each developer's local environment** and not part of the shared repo. The shared backlog is in `docs/06-features/` (feature specs) and `docs/07-open-questions/` (unresolved product decisions).

---

## Sources of truth (read in this order for a task)

1. **Product behavior (what to build)**  
   - Feature-level: `docs/06-features/` — start from `docs/06-features/00-index.md`, then open the feature file for your area (F1–F12).  
   - Canonical PRD (long, Indonesian): `docs/UmrohOS - Product Requirements Document.docx.md` — use section search; do not assume you have seen all of it. Pointer: `docs/00-overview/02-prd-pointer.md`.  
   - Terms: `docs/00-overview/01-glossary.md`.

2. **Bounded contexts & language**  
   - `docs/02-domain/00-bounded-contexts.md`, `docs/02-domain/02-ubiquitous-language.md`.

3. **Architecture & locked tech**  
   - Overview: `docs/01-architecture/00-system-overview.md`.  
   - Stack (authoritative over PRD tech hints): `docs/01-architecture/01-tech-stack.md`.  
   - ADRs: `docs/01-architecture/adr/` (`0001`–`0006` and future). Stack changes need a new ADR.

4. **Per-service technical specs**  
   - `docs/03-services/<service>/` — overview, API, data model, events, status as applicable.

5. **Implementation conventions**  
   - Backend: `docs/04-backend-conventions/` (especially `01-three-layer-architecture.md` — non-negotiable baseline).  
   - Frontend: `docs/05-frontend-conventions/` (Svelte 5 runes + Vite per ADR-0005).

6. **Unresolved product decisions**  
   - `docs/07-open-questions/` — see `docs/07-open-questions/00-how-to-use.md`.  
   - If a feature spec says **TBD** or points to **Qnnn**, read that question file before inventing behavior. Do not silently override an `open` question with a firm product rule.

---

## Rules agents should follow

- **Prefer the minimal doc set** for the current task: feature spec → linked service docs → linked Q files. Do not load the entire `docs/` tree into context by default.  
- **PRD vs repo docs:** If a *product* requirement in the PRD conflicts with another markdown spec, treat the PRD as authoritative for product *and* flag the inconsistency for humans to fix. **Tech choices** follow `01-tech-stack.md` and ADRs, not PRD stack hints.  
- **Open questions:** If status is `open` and no `## Answer` is filled, do not assume stakeholder sign-off; use the file’s **Recommendation** only where the template allows inference, and mark inferred behavior as the team’s convention requires (e.g. `_(Inferred)_` in specs).  
- **Microservices boundaries:** One bounded context per service. Cross-context reads go via gRPC; cross-context **writes are coordinated in-process by the orchestrating service** with explicit per-step compensations, plus a reconciliation cron catching mid-saga crashes (see ADR-0006). Temporal is deferred from MVP and reintroduced only for the F6 visa pipeline — the one multi-day durable workflow. Do not bypass this model without an ADR-level discussion.  
- **Observability:** Tracing/logging/metrics expectations are part of the baseline architecture — see architecture docs before merging “invisible” side paths.

---

## Out of scope for this file

- **Agent-specific instruction files** (`CLAUDE.md`, `.cursor/`, etc.) — gitignored, scoped to the developer who uses that agent.
- **Private per-developer task trackers, progress checklists, testing guides** — each developer keeps their own under their own conventions.
- **Credentials, production URLs, environment secrets** — never committed.
- **Duplicated content from `docs/06-features` or `docs/03-services`** — acceptance criteria, full API fields, and data models live there, not here.

---

## Maintenance

- When a new **global** rule is agreed (affects most sessions), add a **one-line pointer** here or a sentence under **Rules**, and put detail in `docs/` or a new ADR.  
- Keep this file **under ~200 lines** so it stays merge-friendly and actually read.
