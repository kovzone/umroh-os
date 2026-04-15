# ADR 0004 — Monorepo layout

**Status:** Accepted
**Date:** 2026-04-09

## Context

UmrohOS has 10+ services. They share conventions, codegen tooling, and observability config. AI agents work across the codebase in short sessions and need to find related code quickly. Repository fragmentation would force agents to clone and search across many repos per session.

**Note:** Earlier drafts of this ADR proposed a shared root-level `proto/` directory for cross-service contracts. That was reversed after verifying the baseline template — each service owns its own proto under `<svc>/api/grpc_api/pb/<svc>.proto`, and consumers carry a **local copy** of a dependency's proto inside their adapter (`<consumer>/adapter/<dep>_grpc_adapter/pb/<dep>.proto`). See the "Proto ownership" section below.

## Decision

Use a **monorepo** layout. All Go microservices live under a single `services/` directory at the repo root, each as its own Go module (own `go.mod`). Shared infrastructure (docker-compose, Makefile, `_init/`, `monitoring/`, `grafana/`, `temporal/`, `tests/`) sits at the repo root. The `baseline/go-backend-template/` directory is kept as a read-only reference template; `services/` mirrors its service-internal layout (`cmd/`, `api/`, `service/`, `store/`, `util/`) per service.

## Layout

```
umroh-os/
├── CLAUDE.md
├── docker-compose.dev.yml
├── Makefile
├── _init/
│   ├── iam_db/
│   ├── catalog_db/
│   └── ... (one per service that owns data)
├── services/                 ← all Go microservices live here
│   ├── gateway-svc/
│   ├── iam-svc/
│   ├── catalog-svc/
│   ├── booking-svc/
│   ├── jamaah-svc/
│   ├── payment-svc/
│   ├── visa-svc/
│   ├── ops-svc/
│   ├── logistics-svc/
│   ├── finance-svc/
│   ├── crm-svc/
│   └── broker-svc/           ← deferred in MVP (ADR 0006); reserved for F6
├── temporal/                 ← Temporal server config (from baseline; unused in MVP per ADR 0006)
├── monitoring/               ← OTel collector, Prometheus, Loki, Fluent-Bit configs
├── grafana/                  ← Grafana datasources + dashboards
├── tests/                    ← integration & k6 load tests
├── docs/
├── .claude/
│   └── skills/
└── baseline/                 ← reference template, do not edit
    └── go-backend-template/
```

## Rationale

1. **Atomic cross-service changes.** When a contract genuinely needs to change across services, one commit/PR can update the service's proto + every consumer's adapter copy in lockstep.
2. **Single source of truth for tooling.** One Makefile, one docker-compose, one observability config. Each service inherits.
3. **AI agent context.** Future AI sessions can grep the entire codebase from one root. No multi-repo context loading.
4. **Template alignment.** The baseline template is itself a monorepo of services, and we mirror its layout including the proto-ownership model below.

## Proto ownership

Each service **owns its own proto** at `<svc>/api/grpc_api/pb/<svc>.proto`. There is no shared root-level `proto/` directory. When service A needs to call service B, A carries a **local copy** of B's proto at `A/adapter/B_grpc_adapter/pb/B.proto`, alongside the adapter code that wraps the gRPC client.

Why this pattern:
- **Contract ownership is clear.** The service that serves the proto owns it. Consumers take a local copy they control.
- **Versioning is per-consumer.** If B's proto changes, A's adapter copy stays on the old version until A explicitly updates — no silent cross-service breakage from a commit to a shared file.
- **No distributed-monolith smell.** A shared `proto/` directory invites treating every message as implicitly common, which quietly erodes service boundaries.
- **Matches the baseline template exactly.** The baseline's `demo-svc`, `demo-grpc-svc`, and `broker-svc` all follow this pattern — see `baseline/go-backend-template/*/adapter/*/pb/*.proto`.

The same reasoning applies to **utility packages**: each service has its own `<svc>/util/` (config, logging, tracing, apperrors, token, monitoring). There is no shared root-level `util/`. Service-local utils are scaffolded with the service per the `scaffold-service` procedure.

## Consequences

- Each service is its own Go module (`go.mod` per service) for independent dependency control. The root is not a Go module.
- CI must understand which services changed in a given diff (path filters) to avoid rebuilding everything on every commit.
- Service ownership is enforced by code review, not by repo boundaries. Convention docs are non-negotiable.
- Protobuf and util packages are per-service (see "Proto ownership" above). A contract change that affects N consumers is an N+1 commit (the service + each consumer's adapter copy) — deliberate, not a workflow inefficiency.

## Alternatives considered

- **Polyrepo (one repo per service)** — better isolation but worse for cross-service changes and worse for AI agent context. Rejected.
- **Single Go module at the root** — simpler dependency story but couples every service's dependencies. Rejected — independent `go.mod` per service is the template's choice.
