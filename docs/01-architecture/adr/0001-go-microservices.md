# ADR 0001 — Go microservices (overriding the PRD's Node.js hint)

**Status:** Accepted
**Date:** 2026-04-09

## Context

The PRD repeatedly hints at a Node.js + MySQL backend (mysqldump references, JWT/middleware terminology, etc.). UmrohOS has ~11 distinct functional domains (catalog, booking, payment, visa, ops, logistics, finance, etc.), strong domain boundaries, real-time field operations, long-running workflows (visa pipeline takes days), and projected scale of 10,000+ jamaah/year with thousands of concurrent B2C users.

The reviewer supplied `baseline/go-backend-template` as the coding benchmark — a production-ready Go monorepo with three-layer clean architecture, sqlc, Fiber, gRPC, Temporal, and a full OTel observability stack.

## Decision

Build the backend in **Go**, scaffolding every service from `baseline/go-backend-template`. Override the PRD's Node.js hint.

## Rationale

1. **Operator skill match.** The team has ~8 years of Go experience (Fiber, Gin, gRPC), microservices, event-driven systems, and Postgres. The PRD's tech hint reflects general industry conventions, not the actual operator's strengths. Code is reviewed and maintained by the same engineer who will operate it — the language choice should optimize for that engineer's velocity and confidence.
2. **Template-already-exists leverage.** The baseline template is feature-complete (codegen, observability, Temporal wired, three-layer arch enforced). Choosing it saves weeks of bootstrap and gives every future AI session a known-good starting point. Choosing Node.js would require building all of that from scratch with no benchmark to copy from.
3. **Concurrency and footprint.** Go's goroutine model and small per-instance memory footprint suit a microservice mesh of 10+ services running on modest infrastructure.
4. **Type safety end-to-end.** sqlc + protobuf + OpenAPI codegen give compile-time guarantees from DB to API. Comparable Node.js stacks (Prisma + tRPC) exist but require more bespoke wiring.
5. **Workflow-first design.** The PRD has many multi-step processes (visa pipeline, booking saga, refund flow). Temporal is available in the template for the genuinely long-running cases; see ADR 0006 for the MVP decision to defer Temporal for short sagas and reintroduce it for F6 visa pipeline.

## Consequences

- Frontend (Svelte 5 + Vite — see ADR-0005) talks to the Go gateway via REST.
- The PRD's Node.js-flavored examples need translating into Go conventions. This is fine — the PRD is product, not tech.
- Any AI agent unfamiliar with Go must read `baseline/go-backend-template/README.md` and `docs/04-backend-conventions/` before writing code.
- We accept Go's slower iteration speed vs. dynamic languages in exchange for type safety, performance, and operator alignment.

## Alternatives considered

- **Node.js + NestJS + Prisma** — matches PRD hint, but doesn't match the team's strongest skills, and there's no benchmark template.
- **Java + Spring Boot** — the team has Java/Spring experience too, but the supplied template is Go and matches better with the gRPC + Temporal + Fiber stack.
- **Modular monolith in Go** — considered. Rejected because the PRD's domains have clean boundaries that benefit from independent deployability and database ownership, and because the template already does microservices well.
