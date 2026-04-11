# Module List & MoSCoW Priorities

Two companion CSV files live alongside the PRD and define the canonical product scope for UmrohOS:

| File | Purpose |
|---|---|
| `docs/Modul UmrohOS - List Modul.csv` | The full catalogue of 202 modules with Indonesian descriptions, grouped by category and sub-category. |
| `docs/Modul UmrohOS - MosCoW.csv` | The same 202 modules tagged with MoSCoW priority (Must Have / Should Have / Could Have) and a second priority axis (High / Medium / Low). |

**These files are the source of truth for what to build and in what order.** They override any priority guesses in the PRD narrative.

## How to use them

1. **When picking the next task**, start from "Must Have" modules at High priority. Work down to Should Have / Medium before touching Could Have.
2. **When proposing new `Suggested Next Steps`** in `progress.md`, name the exact module (by number or title) from the CSV so the reviewer can cross-check.
3. **When writing a verification block** in `testing-guide.md`, link the work to the module number.
4. **When a task touches multiple modules**, list all of them in the session note so future sessions can see the coverage.

The `prd-researcher` sub-agent (`.claude/agents/prd-researcher.md`) knows how to query these files alongside the PRD — delegate to it when you need to find which module a feature belongs to.

## Summary by priority

Total: **202 modules**

| MoSCoW | Count |
|---|---|
| Must Have | 77 |
| Should Have | 89 |
| Could Have | 36 |

## Summary by top-level category

| Category | Total | Must Have |
|---|---:|---:|
| B2C Front-End | 24 | 11 |
| B2B Front-End | 23 | 8 |
| Marketing & Sales | 23 | 3 |
| Master Product | 16 | 8 |
| Operational & Handling | 42 | 15 |
| Finance | 22 | 15 |
| Admin & Security | 13 | 10 |
| Jamaah Journey | 13 | 3 |
| Dashboard | 13 | 4 |
| Fitur Pelengkap & Daily App (alumni/daily) | 13 | 0 |

Notable signals:
- **Admin & Security** is almost entirely Must Have (10/13). `iam-svc` is the natural first service.
- **Finance** has the highest Must Have count (15) — the finance module is non-negotiable for PSAK compliance.
- **Operational & Handling** has the most modules overall (42) — it's the dominant surface area and maps to several services: `ops-svc`, `logistics-svc`, `jamaah-svc`, `visa-svc`.
- **Fitur Pelengkap & Daily App** has zero Must Haves — alumni hub / daily worship features can wait until the core platform is live.

## Mapping module categories to services

| CSV category | Primary service(s) |
|---|---|
| B2C Front-End | (frontend, deferred) consuming `catalog-svc`, `booking-svc`, `payment-svc` |
| B2B Front-End | (frontend, deferred) consuming `crm-svc`, `catalog-svc`, `booking-svc` |
| Marketing & Sales | `crm-svc` |
| Master Product | `catalog-svc` |
| Operational & Handling → Document Vault / Letter Engine / Smart Grouping / Visa | `jamaah-svc`, `ops-svc`, `visa-svc` |
| Operational & Handling → Terminal Hub / Field Execution | `ops-svc` + `jamaah-svc` (field app) |
| Operational & Handling → Procurement / Inbound & QC / Warehouse / Fulfillment | `logistics-svc` |
| Operational & Handling → Cancellation Management | `payment-svc` + `booking-svc` |
| Finance | `finance-svc` |
| Admin & Security | `iam-svc` |
| Jamaah Journey | `jamaah-svc` + field app frontend (deferred) |
| Dashboard | (frontend, deferred) + observability stack |
| Fitur Pelengkap & Daily App | `crm-svc` (alumni/community), `jamaah-svc` (profile) |

## Implementation order (derived from MoSCoW)

Based on Must Have counts, service dependencies, and vertical slice needs, the suggested build order is:

1. **`iam-svc`** — unblocks everything (every service calls it)
2. **`catalog-svc`** — master data for packages, hotels, airlines, muthawwif
3. **`jamaah-svc`** — pilgrim identity, documents, OCR pipeline
4. **`booking-svc`** — the first revenue-generating surface
5. **`payment-svc`** — close the booking loop
6. **`broker-svc`** — needed to orchestrate the booking saga once catalog + booking + payment exist
7. **`visa-svc`** — parallel with booking maturity
8. **`ops-svc`** — back-office workflows and manifest generation
9. **`logistics-svc`** — kit dispatch on paid bookings
10. **`finance-svc`** — PSAK accounting, consumes events from all upstream services
11. **`crm-svc`** — marketing, agent network, commissions

This ordering is not a hard rule — the reviewer may re-order based on business urgency. But the dependencies listed in `docs/01-architecture/02-service-map.md` must still be respected.

## Do not try to read the whole CSV in one session

The CSVs are large (202 rows). When you need to look up a specific module or category, either:
- Use `grep` / `awk` on the CSV directly, or
- Delegate to the `prd-researcher` sub-agent with a targeted query.

Never dump the entire CSV into a session's context — it wastes tokens on modules irrelevant to the current task.
