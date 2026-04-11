# ADR 0002 — PostgreSQL 15 (overriding the PRD's MySQL hint)

**Status:** Accepted
**Date:** 2026-04-09

## Context

The PRD references MySQL throughout (mysqldump, schema patterns). The baseline `go-backend-template` is built on PostgreSQL 15 with pgx/v5 and sqlc. UmrohOS has complex requirements: PSAK double-entry accounting (precise decimals), multi-currency journaling, mahram family-graph queries, audit trail with immutable history, and JSON document storage for OCR results.

## Decision

Use **PostgreSQL 15** for every service. Override the PRD's MySQL hint.

## Rationale

1. **Decimal precision for accounting.** Postgres `numeric` is strict. MySQL has historically had quirks with `DECIMAL` precision. Critical for PSAK compliance.
2. **JSONB for OCR / document data.** Postgres JSONB with GIN indexes is mature. MySQL JSON is functional but weaker for indexed queries.
3. **Recursive CTEs for the family graph.** Mahram validation needs to walk a family tree. Postgres `WITH RECURSIVE` is well-supported.
4. **Row-level locking maturity.** PSAK journaling requires strict transaction isolation. Postgres MVCC is better understood by the team and by the template.
5. **Template alignment.** sqlc is configured for Postgres in the template. Switching to MySQL would mean a different sqlc engine and rewriting every example query.
6. **Operator preference.** The team's primary database is Postgres; MySQL is familiar but not the default.

## Consequences

- Each service owns one logical database (`iam_db`, `catalog_db`, …) on a shared Postgres instance in dev. Production may run separate clusters.
- Migrations are managed via alphabetically-ordered DDL files in `_init/<db>/` and `<svc>/store/postgres_store/schema/`, per the template convention.
- All sqlc queries use Postgres dialect.
- The PRD's `mysqldump`-flavored examples translate to `pg_dump` workflows.

## Alternatives considered

- **MySQL 8** — matches PRD hint, but the entire template would need re-tooling, decimal handling is weaker, and the team prefers Postgres.
- **One DB per service on separate clusters from day one** — rejected as premature. Logical separation on a shared cluster is enough until production sizing is decided.
