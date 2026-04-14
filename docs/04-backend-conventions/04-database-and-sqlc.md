# Database & sqlc

PostgreSQL 15 is the only database. **sqlc** is the only path to it. Hand-rolled SQL outside `queries/` is forbidden.

## Schema management

Schema lives in two parallel locations (per the baseline template's convention):

1. **`_init/<svc>_db/`** — alphabetically-ordered DDL files mounted into the dev Postgres container at startup. This is the canonical source for fresh dev environments.
2. **`<svc>/store/postgres_store/schema/`** — same DDL, used by sqlc to validate query types. Must stay in sync with `_init/<svc>_db/`.

> A future ADR will introduce a real migration tool (golang-migrate or atlas) when production deployment is in scope. For now, ordered DDL is the convention.

### File ordering

DDL files are loaded in alphabetical order by Postgres. Use prefixes:

```
_init/iam_db/
├── 00_extensions.sql        ← CREATE EXTENSION ...
├── 10_types.sql             ← CREATE TYPE ... (enums, domains)
├── 20_tables.sql            ← CREATE TABLE ...
├── 30_indexes.sql           ← CREATE INDEX ...
├── 40_functions.sql         ← CREATE FUNCTION ...
└── 90_seed.sql              ← INSERT ... (dev seed data only)
```

### Table conventions

- Plural snake_case (`users`, `audit_logs`).
- Primary key: `id UUID PRIMARY KEY DEFAULT gen_random_uuid()` (use the `pgcrypto` extension).
- Timestamps: `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`, `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`. Updated via trigger or service code.
- Soft-delete: `deleted_at TIMESTAMPTZ` (nullable). Filtering happens in queries, not in views.
- Foreign keys: always defined. `ON DELETE` policy explicit (`CASCADE`, `SET NULL`, `RESTRICT`).
- Branch scope: `branch_id UUID NOT NULL REFERENCES branches(id)` on every table that's branch-scoped.
- Audit fields: `created_by UUID`, `updated_by UUID` referencing `users(id)` where applicable.

### Status enums

Use Postgres enum types for status fields:

```sql
CREATE TYPE booking_status AS ENUM (
    'draft',
    'pending_payment',
    'paid_in_full',
    'partially_paid',
    'cancelled',
    'completed'
);
```

Always lowercase snake_case. See `docs/02-domain/02-ubiquitous-language.md` for the canonical list.

## sqlc

### Query files

- Path: `<svc>/store/postgres_store/queries/<domain>.sql`
- One file per domain area (`users.sql`, `roles.sql`, `bookings.sql`, etc.).
- Each query is annotated:
  ```sql
  -- name: GetUserByEmail :one
  SELECT id, email, name, branch_id, created_at
  FROM users
  WHERE email = $1 AND deleted_at IS NULL;

  -- name: CreateUser :one
  INSERT INTO users (email, name, branch_id)
  VALUES ($1, $2, $3)
  RETURNING id, email, name, branch_id, created_at;
  ```

### Annotation types

| Annotation | Returns |
|---|---|
| `:one` | exactly one row, error if zero |
| `:many` | slice of rows |
| `:exec` | no return |
| `:execrows` | affected row count |

### Codegen

`cd <svc>/store/postgres_store && sqlc generate` (or `make sqlc-<svc>`). The generated code goes in `<svc>/store/postgres_store/sqlc/`. Do not edit it.

## Transactions

Use the template's `WithTx` helper for any write that touches multiple rows or multiple tables:

```go
err := s.store.WithTx(ctx, func(q sqlc.Querier) error {
    if _, err := q.CreateUser(ctx, ...); err != nil {
        return err
    }
    if _, err := q.AssignRole(ctx, ...); err != nil {
        return err
    }
    return nil
})
```

`WithTx` provides:
- Automatic retry on serialization failure / deadlock
- Tracing (the transaction is its own span)
- Logging

Default isolation: `READ COMMITTED`. Use `WithTxOptions` for stricter levels.

## Hard rules

- **No raw `db.Exec` / `db.Query` in service code.** Always via sqlc.
- **No SQL strings in Go files outside `queries/`.**
- **No cross-database queries.** Each service owns its database.
- **No SELECT \***, except in query annotations where sqlc needs the column list. Prefer named columns for stability.
- **Always parameterize.** Never string-concatenate values into SQL. (sqlc enforces this anyway.)
- **Schema must be in sync** between `_init/<svc>_db/` and `<svc>/store/postgres_store/schema/`. Drift causes sqlc codegen to silently disagree with the running DB.
