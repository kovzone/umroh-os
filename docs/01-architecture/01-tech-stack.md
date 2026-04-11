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
| Workflow orchestration | Temporal.io | baseline |
| Trace export | OpenTelemetry → Tempo | baseline |
| Metrics | Prometheus | baseline |
| Logs | zerolog (JSON) → Fluent-Bit → Loki | baseline |
| Dashboards | Grafana (logs/metrics/traces unified) | baseline |
| Config | Viper (JSON + env override) | baseline |
| Auth tokens | PASETO/JWT via `util/token` | baseline |
| Containerization | Docker + docker-compose.dev.yml | baseline |
| Unit testing | testify/mock + testify/require | baseline |
| Load testing | k6 | baseline |

## Why this stack

The baseline `go-backend-template` already implements this entire stack end-to-end with codegen, observability, and Temporal wired up. The team's strongest skills are on the Go / Fiber / gRPC / Postgres side. Choosing the template wholesale is the highest-leverage decision: it removes weeks of bootstrap work and gives every future session a known-good starting point.

The PRD suggests Node.js/MySQL. We override that — the Go template was supplied as the explicit coding benchmark, and the team's strongest skills align with it. The PRD is product, not tech.

See:
- `adr/0001-go-microservices.md` — why Go
- `adr/0002-postgres-over-mysql.md` — why Postgres
- `adr/0003-temporal-for-workflows.md` — why Temporal
- `adr/0004-monorepo-layout.md` — why one repo for all services

## What's deferred

- **Frontend stack.** React+Vite is the planned direction but the frontend standard will be provided separately. No frontend work in this repo until then.
- **CI/CD pipeline.** No GitHub Actions / CircleCI configured yet. Add when needed.
- **Production deployment.** Docker Compose is dev-only. Production target (Kubernetes? Cloud Run? Bare VMs?) is undecided and will need an ADR.
- **Kafka / event streaming.** Considered and rejected for MVP — Temporal handles workflow orchestration and gRPC handles synchronous calls. Revisit if a use case appears that genuinely needs pub/sub fanout.

## What you may NOT introduce without an ADR

- A second database engine
- A different web framework
- A different ORM or query approach
- A different observability vendor
- A different config or logging library

If you think the stack is wrong for a specific need, write an ADR and discuss with the reviewer before changing anything.
