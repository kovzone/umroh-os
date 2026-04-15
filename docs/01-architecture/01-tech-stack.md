# Tech Stack

Locked in the working environment setup session (2026-04-09). Changes to this stack require an ADR.

## The stack

| Layer | Choice | Source |
|---|---|---|
| Language | Go 1.22+ | baseline template |
| HTTP framework (REST) | Fiber v2 | baseline |
| Internal RPC | gRPC + protobuf | baseline |
| Database | PostgreSQL 15 | baseline (overrides PRD's MySQL hint) |
| DB driver | pgx/v5 | baseline |
| Data access codegen | sqlc | baseline |
| REST API codegen | OpenAPI 3 + oapi-codegen | baseline |
| Trace export | OpenTelemetry → Tempo | baseline |
| Metrics | Prometheus | baseline |
| Logs | zerolog (JSON) → Fluent-Bit → Loki | baseline |
| Dashboards | Grafana (logs/metrics/traces unified) | baseline |
| Config | Viper (JSON + env override) | baseline |
| Auth tokens | PASETO/JWT via `util/token` | baseline |
| Containerization | Docker + docker-compose.dev.yml | baseline |
| Unit testing | testify/mock + testify/require | baseline |
| Load testing | k6 | baseline |
| Frontend framework | Svelte 5 (runes mode) | ADR-0005 |
| Frontend build tool | Vite | ADR-0005 |
| Frontend best-practices guide | `.claude/skills/svelte-core-bestpractices/` | ADR-0005 |
| Frontend conventions | `docs/05-frontend-conventions/` | ADR-0005 |

## Why this stack

The baseline `go-backend-template` already implements this entire stack end-to-end with codegen and observability wired up. The team's strongest skills are on the Go / Fiber / gRPC / Postgres side. Choosing the template wholesale is the highest-leverage decision: it removes weeks of bootstrap work and gives every future session a known-good starting point.

The PRD suggests Node.js/MySQL. We override that — the Go template was supplied as the explicit coding benchmark, and the team's strongest skills align with it. The PRD is product, not tech.

See:
- `adr/0001-go-microservices.md` — why Go
- `adr/0002-postgres-over-mysql.md` — why Postgres
- `adr/0003-temporal-for-workflows.md` — **deferred for MVP; see ADR 0006**
- `adr/0004-monorepo-layout.md` — why one repo for all services
- `adr/0005-svelte-frontend.md` — why Svelte 5 (over the originally assumed React+Vite)
- `adr/0006-defer-temporal-to-f6.md` — why Temporal is deferred; MVP uses in-process saga coordination

## Frontend stack

**Svelte 5 (runes mode) + Vite.** Conventions live in `docs/05-frontend-conventions/`; the authoritative runtime skill is `.claude/skills/svelte-core-bestpractices/`. See ADR-0005 for why the earlier React+Vite assumption was replaced.

## What's deferred
- **Temporal.io workflow orchestration.** Deferred for MVP per ADR 0006. In-process saga coordination in the orchestrating service + reconciliation cron is the MVP pattern for F4/F5. Temporal is brought back when F6 visa pipeline is implemented (multi-day durable workflow — the use case that genuinely needs it).
- **CI/CD pipeline.** No GitHub Actions / CircleCI configured yet. Add when needed.
- **Production deployment.** Docker Compose is dev-only. Production target (Kubernetes? Cloud Run? Bare VMs?) is undecided and will need an ADR.
- **Kafka / event streaming.** Considered and rejected for MVP — gRPC handles synchronous calls and in-process sagas handle orchestration. Revisit if a use case appears that genuinely needs pub/sub fanout.

## What you may NOT introduce without an ADR

- A second database engine
- A different web framework
- A different ORM or query approach
- A different observability vendor
- A different config or logging library

If you think the stack is wrong for a specific need, write an ADR and discuss with the reviewer before changing anything.
