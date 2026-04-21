# Playwright e2e (`tests/e2e`)

Two projects (see `playwright.config.ts`):

- **`api`** — hits backend REST ports directly (`01-*`, `02-*`).
- **`browser`** — Chromium against **`core-web`** at `CORE_WEB_URL` (default `http://localhost:3001`) for `03-*` specs.

## Prerequisites

1. **Docker Desktop** running (Linux engine).
2. **Per-service `config.json`** — `docker-compose.dev.yml` mounts `./services/<svc>/config.json` as a **file**. If that path is a **folder** (empty `config.json` directories are a common mistake on Windows) or the file is missing, Go services **restart forever** (`Config File "config" Not Found in [/app]`). From repo root, run the same seed as README bootstrap:
   - **PowerShell:** `.\scripts\ensure-service-configs.ps1`
   - **Bash:** the `for svc in services/*/` loop in root `README.md` § Bootstrap.
3. **Dev stack up** — gateway `:4000`, services `:4001+`, **`core-web` `:3001`** (see repo root `README.md` / `docker-compose.dev.yml`).
4. **Migrations applied** — `umrohos_dev` schema; same as `make migrate-up` from the Makefile (`golang-migrate` on PATH).
5. **Node.js 20+**.

### Windows / Playwright networking

- **API** specs use Node’s HTTP client. On Windows, `localhost` often resolves to **`::1` first**, while Docker Desktop publishes backend ports on IPv4, which yields `ECONNREFUSED`. Defaults in `tests/e2e/lib/services.ts` use **`127.0.0.1`** for gateway and each service.
- **Browser** specs load `core-web` via `baseURL` in `playwright.config.ts`, default **`http://localhost:3001`**. On some Docker Desktop + Windows setups, **`127.0.0.1:3001` hits the wrong process** (e.g. Grafana) while `localhost:3001` reaches core-web — set `CORE_WEB_URL` in `.env` if your machine differs.

## Linux / macOS (Makefile)

From the repository root:

```bash
make dev-bootstrap   # compose up + wait for postgres + migrate up
make e2e-install     # npm install + playwright chromium + apps/core-web npm install
make e2e             # cd tests/e2e && npm test
```

## Windows (no `make`)

From the repository root in **PowerShell**:

```powershell
.\scripts\e2e-local.ps1
```

Optional flags if the stack and DB are already prepared:

```powershell
.\scripts\e2e-local.ps1 -SkipUp       # do not run docker compose up
.\scripts\e2e-local.ps1 -SkipMigrate  # skip migrate up
```

If `migrate` is not installed, install [golang-migrate](https://github.com/golang-migrate/migrate) and ensure it is on `PATH`, or apply migrations once from an environment where `make migrate-up` works (for example Git Bash with `make`).

## Manual steps (any OS)

```bash
docker compose -f docker-compose.dev.yml up -d
# wait until postgres accepts connections, then:
migrate -source file://migration -database "postgres://postgres:changeme@localhost:5432/umrohos_dev?sslmode=disable" up
cd tests/e2e && npm install && npx playwright install chromium && npm test
```

## Environment

Copy `tests/e2e/.env.example` to `tests/e2e/.env` if you need to override URLs (defaults match `docker-compose.dev.yml` host ports).
