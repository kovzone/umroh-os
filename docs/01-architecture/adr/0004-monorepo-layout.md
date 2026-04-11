# ADR 0004 — Monorepo layout

**Status:** Accepted
**Date:** 2026-04-09

## Context

UmrohOS has 10+ services. They share conventions, codegen tooling, observability config, and protobuf definitions. AI agents work across the codebase in short sessions and need to find related code quickly. Repository fragmentation would force agents to clone and search across many repos per session.

## Decision

Use a **monorepo** layout, mirroring `baseline/go-backend-template`. All services live as top-level directories in the root of `umroh-os/`. Each service is its own Go module (own `go.mod`).

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
├── gateway-svc/
├── iam-svc/
├── catalog-svc/
├── booking-svc/
├── jamaah-svc/
├── payment-svc/
├── visa-svc/
├── ops-svc/
├── logistics-svc/
├── finance-svc/
├── crm-svc/
├── broker-svc/
├── proto/                    ← shared .proto files (cross-service contracts)
├── temporal/                 ← Temporal server config
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

1. **Atomic cross-service changes.** Adding a field to a shared proto and updating all consumers happens in one commit/PR.
2. **Single source of truth for tooling.** One Makefile, one docker-compose, one observability config. Each service inherits.
3. **AI agent context.** Future Claude sessions can grep the entire codebase from one root. No multi-repo context loading.
4. **Template alignment.** The baseline template is itself a monorepo of services.

## Consequences

- Each service is its own Go module (`go.mod` per service) for independent dependency control. The root is not a Go module.
- CI must understand which services changed in a given diff (path filters) to avoid rebuilding everything on every commit.
- Service ownership is enforced by code review, not by repo boundaries. Convention docs are non-negotiable.
- The `proto/` directory at the root holds shared message types and service contracts. Each service generates its own client/server code from these.

## Alternatives considered

- **Polyrepo (one repo per service)** — better isolation but worse for cross-service changes and worse for AI agent context. Rejected.
- **Single Go module at the root** — simpler dependency story but couples every service's dependencies. Rejected — independent `go.mod` per service is the template's choice.
