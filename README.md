# UmrohOS

An end-to-end ERP platform for Umrah and Hajj travel agencies.

UmrohOS brings package booking, agent networks, marketing, document and visa processing, warehouse logistics, finance and accounting, and jamaah field operations together in one place — so a pilgrimage travel agency can run its entire operation from a single system instead of juggling a dozen fragmented tools.

---

## Current state

**Phase: scaffolding complete, no business features implemented yet.** If you clone and run, you will get a full running monorepo of empty services wired together end-to-end — but nothing you can book, pay for, or file a visa with. Those land in subsequent per-feature commits.

### What's in the repository today

| Area | What's there | What's NOT there |
|---|---|---|
| **Backend services** (`services/`) | 11 Go microservices: `gateway-svc` (REST edge) + 10 domain services (`iam-svc`, `catalog-svc`, `booking-svc`, `jamaah-svc`, `payment-svc`, `visa-svc`, `ops-svc`, `logistics-svc`, `finance-svc`, `crm-svc`). Each exposes `/system/live`, `/system/ready`, `/system/diagnostics/db-tx` via REST and a placeholder `Healthz` RPC via gRPC (plus the standard gRPC health protocol). | Real business endpoints. Auth. Per-service DDL (only the shared `public.diagnostics` table exists). `broker-svc` (deferred per ADR 0006). |
| **Frontend apps** (`apps/`) | 1 app: `core-web` — SvelteKit + TypeScript + Svelte 5 (runes). Two routes: `/` (landing page) and `/system/status` (live service-health dashboard polling the 10 backends through gateway-svc). | Auth / login. Real admin pages. B2C / B2B / field-app surfaces. |
| **Database** (`migration/`) | Single shared `umrohos_dev` Postgres database; golang-migrate at repo root. Two migrations: `000001_init` (placeholder) + `000002_scaffold_services` (shared `public.diagnostics` table). | Per-service schemas and tables — land per feature task (`000003_add_iam_users_and_roles` is next, etc.). |
| **Observability** (`monitoring/`, `grafana/`) | Prometheus scrapes all 11 services; Loki aggregates logs via fluent-bit; Tempo receives OTel traces; Grafana auto-provisions **one dashboard per service** in the "UmrohOS Services" folder. | Alerting rules. SLO dashboards. Frontend RUM / Sentry. |
| **Tests** (`tests/`) | Playwright e2e at `tests/e2e/` — two projects (`api` + `browser`), **50 tests** when the full stack is up (see **`tests/e2e/README.md`**; Windows: **`scripts/e2e-local.ps1`**). k6 load stubs at `tests/k6/`. | Per-service Go unit tests (scaffolds only; real tests land with feature work). Vitest suite in `core-web` has one trivial spec as a placeholder. |
| **Docs** (`docs/`) | ADRs, architecture overview, service specs, domain glossary, backend + frontend conventions, feature specs (F1–F12), 80+ open stakeholder questions, commit + git workflow conventions. | Everything is documented; that's the point. |

### Ports you'll see running

| Service / app | REST port | gRPC port | Notes |
|---|---:|---:|---|
| `core-web`       | 3001 | — | Landing + status dashboard |
| `gateway-svc`    | 4000 | — | REST-only edge; proxies `/v1/<svc>/…` to backends |
| `iam-svc`        | 4001 | 50051 | |
| `catalog-svc`    | 4002 | 50052 | |
| `booking-svc`    | 4003 | 50053 | |
| `jamaah-svc`     | 4004 | 50054 | |
| `payment-svc`    | 4005 | 50055 | |
| `visa-svc`       | 4006 | 50056 | |
| `ops-svc`        | 4007 | 50057 | |
| `logistics-svc`  | 4008 | 50058 | |
| `finance-svc`    | 4009 | 50059 | |
| `crm-svc`        | 4010 | 50060 | |
| Postgres         | 5432 | — | Single `umrohos_dev` database |
| Grafana          | 3000 | — | admin / admin |
| Prometheus       | 9090 | — | |
| Jaeger UI        | 16686 | — | (also wired to Tempo) |
| Loki             | 3100 | — | |
| Tempo            | 3200 | — | |

---

## How to run

### Prerequisites

- **Docker** (+ Docker Compose plugin) running.
- **Node.js 20+** — required for `core-web` and the Playwright e2e suite.
- **golang-migrate** CLI — `brew install golang-migrate` on macOS. Needed for `make migrate-*` targets (ADR 0007).
- **Go 1.22+** — only needed if you'll build services outside Docker.

### Bootstrap

```bash
# 1. Seed per-service runtime configs from their .sample counterparts
#    (config.json is gitignored — each contributor creates their own)
for svc in services/*/; do
  [ -f "$svc/config.json.sample" ] && cp "$svc/config.json.sample" "$svc/config.json"
done
cp apps/core-web/.env.example apps/core-web/.env

# 2. Bring the full stack up and apply all migrations
make dev-bootstrap

# 3. Confirm 20 containers Up (8 infra + 11 backend svc + 1 frontend app)
make dev-ps
```

### Sanity checks

```bash
# Landing page (prerendered)
open http://localhost:3001/

# Service-health dashboard (polls 10 backends via gateway every 5s)
open http://localhost:3001/system/status

# Gateway itself
curl -sS http://localhost:4000/system/live        # → {"data":{"ok":true}}
curl -sS http://localhost:4000/v1/iam/system/live # → proxied to iam-svc

# Observability
open http://localhost:3000/                       # Grafana, login admin/admin
open http://localhost:9090/targets                # Prometheus — 11 jobs UP
```

### Run the e2e suite

```bash
make e2e-install   # one-time — npm install + npx playwright install chromium + per-app deps
make e2e           # expect: 47 passed (43 api + 4 browser)
```

### Teardown

```bash
make dev-down-v    # stop containers + remove volumes
make dev-rm-all    # remove built service/app images (idempotent)
```

### Other useful Make targets

```bash
make help              # list every target with its docstring
make generate          # run sqlc + oapi-codegen + protoc + openapi-typescript for all components
make migrate-create NAME=add_iam_users_and_roles
make migrate-up        # apply pending migrations against localhost:5432/umrohos_dev
make test              # unit tests for every Go service
make dev-rebuild SVC=iam-svc
make web-dev           # rebuild + up core-web
```

---

## Where things live

- **Vision, glossary, module list:** `docs/00-overview/`
- **Architecture + ADRs:** `docs/01-architecture/`
- **Domain / bounded contexts:** `docs/02-domain/`
- **Per-service specs:** `docs/03-services/<NN>-<svc>/`
- **Backend conventions:** `docs/04-backend-conventions/`
- **Frontend conventions:** `docs/05-frontend-conventions/`
- **Feature specs (the PRD → code middle layer):** `docs/06-features/`
- **Open stakeholder questions:** `docs/07-open-questions/`
- **Commit message conventions:** `docs/08-commit-conventions.md`
- **Repeatable agent skills:** `.claude/skills/`
- **Reference Go service template (read-only):** `baseline/go-backend-template/`

The PRD itself lives at `docs/UmrohOS - Product Requirements Document.docx.md` (~1,620 lines, mostly Indonesian terminology).

---

## Contributing

Two developers work on this codebase, both full-stack. Work is sliced by feature (booking, visa, finance, …), not by frontend-vs-backend. Commit messages follow [`docs/08-commit-conventions.md`](docs/08-commit-conventions.md); branching and PR rules are in [`docs/04-backend-conventions/08-git-workflow.md`](docs/04-backend-conventions/08-git-workflow.md).

---

## Author

Elda Mahaindra ([faith030@gmail.com](mailto:faith030@gmail.com))
