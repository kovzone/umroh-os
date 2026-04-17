# ADR 0007 — Migration-based schema management

**Status:** Accepted
**Date:** 2026-04-16

## Context

The baseline template (`go-backend-template`) manages database schema in two places:

1. `_init/<db>/` — DDL files mounted into the Postgres container at startup via `docker-entrypoint-initdb.d`.
2. `<svc>/store/postgres_store/schema/` — the same DDL, duplicated so that sqlc can read it for codegen.

This dual-location approach has problems at scale:

- **Drift risk.** Two copies of the schema must stay in sync manually. When they diverge, sqlc codegen silently disagrees with the running database.
- **No versioning.** DDL files are loaded once at container init. Adding a column means tearing down the volume and rebuilding from scratch — there is no incremental migration path.
- **Scales badly.** With 10+ services, the `_init/` directory and per-service `schema/` directories multiply the maintenance burden.

The existing `docs/04-backend-conventions/04-database-and-sqlc.md` anticipated this: *"A future ADR will introduce a real migration tool when production deployment is in scope."* This is that ADR.

## Decision

### Single database, multiple schemas

All services share one Postgres database (`umrohos_dev` in dev). Each service owns a Postgres schema namespace (`iam`, `catalog`, `booking`, etc.) within that database. This simplifies connection management (one `LOCAL_DB_URL`), dev setup (one `POSTGRES_DB` in docker-compose), and allows cross-service reads without cross-database wiring.

### Flat `migration/` directory at repo root

Use **golang-migrate** with a single `migration/` directory at the repo root containing numbered up/down SQL migration pairs:

```
migration/
├── 000001_init.up.sql
├── 000001_init.down.sql
├── 000002_scaffold_services.up.sql
├── 000002_scaffold_services.down.sql
├── 000003_add_iam_users_and_roles.up.sql
├── 000003_add_iam_users_and_roles.down.sql
└── ...
```

Migrations create schemas and tables within them (`CREATE SCHEMA iam; CREATE TABLE iam.users ...`). All services' `sqlc.yaml` point to the same `migration/` directory — sqlc reads all `.up.sql` files to understand the full database schema.

### sqlc integration

```yaml
# services/iam-svc/store/postgres_store/sqlc.yaml
schema: "../../../../migration"
```

sqlc reads all `.up.sql` files from `migration/` to understand the full schema, then validates each service's queries against it. One location, one source of truth.

### Migration execution

Migrations are applied via the `golang-migrate` CLI, invoked from Makefile targets:

- `make migrate-up` — apply all pending migrations.
- `make migrate-down STEPS=N` — roll back N migrations.
- `make migrate-create NAME=<name>` — create a new up/down migration pair.
- `make migrate-version` — show current migration version.
- `make migrate-force VERSION=N` — force-set version (recovery from dirty state).
- `make dev-bootstrap` — start the dev environment and run all migrations.

### No `_init/` directory

The `_init/` directory is removed entirely. The single database is created by Postgres via the `POSTGRES_DB` environment variable in docker-compose. Schema is applied by `make migrate-up`.

## Rationale

1. **Single source of truth.** Migration files define the schema once. sqlc reads from the same files. No duplication.
2. **Incremental changes.** Adding a column is a new migration file, not a volume teardown. Rollback is `make migrate-down`.
3. **Proven in production.** golang-migrate is mature (v4), well-maintained, and already used in team projects (m7-bms-backend, ganesa).
4. **Root-level centralization.** All schema lives in `migration/` at the repo root, outside any service directory. Easy to grep, easy for AI agents to find.
5. **Single database simplicity.** One connection string, one `make migrate-up`, one Postgres instance with schema namespaces for logical isolation.

## Why a flat directory (not nested `migration/migrations/`)?

The ganesa project uses `migration/migrations/` because the parent `migration/` package holds Go code (`embed.go`, `migrator.go`) for embedding SQL files into the binary. UmrohOS doesn't need Go migration code — we use the `golang-migrate` CLI from the Makefile. A flat `migration/` directory with only SQL files is simpler.

## Why not embed.FS?

Each service is its own Go module (`services/<svc>/go.mod` per ADR 0004). Go's `embed` directive can only embed files within or below the module directory. Migration files live at the repo root — outside any service module. Workarounds (root Go module, symlinks) are fragile and violate ADR 0004. The Makefile + CLI approach is simpler and fits the monorepo layout.

## Why not per-database subdirectories?

A single shared database with schema namespaces is simpler than per-service databases. It avoids creating multiple databases, multiple connection strings, and parameterized Makefile targets. Logical isolation via Postgres schemas (`iam.*`, `catalog.*`) is sufficient — each service queries only its own schema.

## Consequences

- **New dev dependency.** Developers must install `golang-migrate` CLI (`brew install golang-migrate` on macOS).
- **`_init/` is removed.** Database creation is handled by `POSTGRES_DB` in docker-compose. Schema is applied via `make migrate-up`.
- **Per-service `schema/` directories are removed.** sqlc reads from `migration/` at the repo root.
- **All services share schema visibility.** sqlc sees the full schema across all services. Each service's queries should only reference tables in its own Postgres schema namespace.
- **Volume teardown still works.** `make dev-down-v` + `make dev-bootstrap` gives a clean-slate dev environment. But incremental changes no longer require it.
- **Scaffolding procedure updated.** The `scaffold-service` skill no longer creates `_init/<db>/` or per-service `schema/` directories.

## Supersedes

- The dual-location schema pattern from the baseline template (`_init/<db>/` + `<svc>/store/postgres_store/schema/`).
- The "future ADR" note in `docs/04-backend-conventions/04-database-and-sqlc.md`.
