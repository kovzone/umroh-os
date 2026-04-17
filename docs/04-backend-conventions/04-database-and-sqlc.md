# Database & sqlc

PostgreSQL 15 is the only database. **sqlc** is the only path to it. Hand-rolled SQL outside `queries/` is forbidden.

## Schema management

Schema is managed via **golang-migrate** with a flat `migration/` directory at the repo root (see ADR 0007). All services share a single database (`umrohos_dev`) with per-service Postgres schema namespaces (`iam`, `catalog`, `booking`, etc.):

```
migration/
├── 000001_init.up.sql
├── 000001_init.down.sql
├── 000002_scaffold_services.up.sql         # shared public.diagnostics table
├── 000002_scaffold_services.down.sql
├── 000003_add_iam_users_and_roles.up.sql   # F1.2 — per-service tables in iam schema
├── 000003_add_iam_users_and_roles.down.sql
└── ...
```

Filenames describe **the task the developer is performing in that commit**, not the DDL diff. `000002_scaffold_services` captures "the scaffolding step" as a whole; `000003_add_iam_users_and_roles` captures F1.2. Avoid mechanical names like `000002_create_users_table` — they hide the commit's intent.

Every service's `sqlc.yaml` points to the same migration directory:

```yaml
schema: "../../../../migration"
```

sqlc reads all `.up.sql` files to understand the full database schema. **One location, one source of truth** — no per-service `schema/` directories, no duplication.

### Creating a new migration

```bash
make migrate-create NAME=add_iam_users_and_roles
# Creates: migration/000003_add_iam_users_and_roles.up.sql
#          migration/000003_add_iam_users_and_roles.down.sql
```

The name (`add_iam_users_and_roles`) describes the task. Write DDL in the `.up.sql` file using schema-qualified names (`CREATE SCHEMA IF NOT EXISTS iam; CREATE TABLE iam.users ...`); write the matching teardown in `.down.sql`. Then run `make migrate-up` to apply and `make sqlc` to regenerate query code.

### Migration workflow

```bash
make dev-bootstrap        # Start postgres + run all migrations
make migrate-up           # Apply all pending migrations
make migrate-down STEPS=1 # Roll back one migration
make migrate-version      # Show current version
make migrate-force VERSION=1  # Fix dirty state
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
- **No cross-schema queries.** Each service owns its Postgres schema namespace. Do not join across schema boundaries.
- **No SELECT \***, except in query annotations where sqlc needs the column list. Prefer named columns for stability.
- **Always parameterize.** Never string-concatenate values into SQL. (sqlc enforces this anyway.)
- **Migrations are the schema source of truth.** sqlc reads from `migration/`. Do not create per-service `schema/` directories.
