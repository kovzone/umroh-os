# Feature → Backlog Mapping (priority owned by feature spec)

This document is the operational layer between:

- `docs/06-features/*` (feature detail, acceptance, edge cases), and
- slice task codes in `[05-slice-engineering-checklist-and-task-codes.md](./05-slice-engineering-checklist-and-task-codes.md)`.

Phase order **0 → 6** follows `[04-delivery-plan-2p-sequence-first.md](./04-delivery-plan-2p-sequence-first.md)` (team & slice context); per-slice gates and task-code definitions are in **05**.

Core principles:

1. `Priority` (`MH-MVP`, `MH-V1`, `SH`, `CH`) belongs to the **feature/backlog item**.
2. `Phase` + `Exec seq` + `Slice` + `Task Code` describe **integration order** and **work packages** (not calendar deadlines).
3. One feature detail may split into several backlog rows to keep PR size small.
4. **`Blocked by gate`** — lists the **backlog ids** for the **Engineering freeze** set in **[05](./05-slice-engineering-checklist-and-task-codes.md)** (see **Freeze backlog id index** below). Spans such as `BL-JNT-001..004` mean **inclusive** ids `BL-JNT-001` through `BL-JNT-004`. The **05** tables still name the same gates with **task codes** (for example `S1-J-01..S1-J-04; S1-E-01; S1-L-01`); **06** uses the `BL-*` ids for execution tracking.
5. **Joint freeze (`Sx-J-*`):** each Joint card in **05** has a matching row here — domain **`JNT`**, ids **`BL-JNT-*`** (global sequence **001–015** across S1→S5).
6. **Engineering / Lead freeze (`Sx-E-01`, `Sx-L-01`):** where **05** lists a separate **`Sx-E-01`** or **`Sx-L-01`** row, it maps to **`BL-EGV-*`** (engineering gate review) or **`BL-LGV-*`** (lead / UX gate). **S3** has **`S3-E-01`** only; **S4** / **S5** have **Joint rows only** in **05** (no **`BL-EGV-*` / `BL-LGV-*`** for those slices in this index).

**Freeze backlog id index (S1–S5):**

| Slice | `BL-JNT-*` (Joint contracts) | `BL-EGV-*` | `BL-LGV-*` |
| ----- | ---------------------------- | ------------ | ------------ |
| **S1** | `001..004` | `BL-EGV-001` (`S1-E-01`) | `BL-LGV-001` (`S1-L-01`) |
| **S2** | `005..008` | `BL-EGV-002` (`S2-E-01`) | `BL-LGV-002` (`S2-L-01`) |
| **S3** | `009..011` | `BL-EGV-003` (`S3-E-01`) | — |
| **S4** | `012..013` | — | — |
| **S5** | `014..015` | — | — |

**Dashboard rows (`BL-DASH-*`, F11 widgets under Phase 6 slice codes such as `S3-L-05`, `S4-L-03`, `S5-L-02`):** these rows are **UI placement and product scope** (which screen owns which widget). **Data aggregation, cache TTLs, read-replica usage, and optional `dashboard-svc`** follow the engineering decision in **`docs/07-open-questions/Q066-dashboard-aggregation-architecture.md`** — not a second architecture.

Coverage notes:

- Phases **0–5** list backlog **already mapped** to primary slice/task codes.
- Phase **6** includes all **202 module numbers** from `docs/Modul UmrohOS - MosCoW.csv` as `BL-*` rows (depth expansion), plus umbrella **6.A–6.E**. Work order follows **`Exec seq`** up to **854** after **6.E** (see Phase 6 blockquote).
- **CSV `No` index → Phase 6 subsection**: `#1–24`→**6.L**; `#25–70` & `#199–202`→**6.H**; `#71–86`→**6.G**; `#87–128`→**6.K**; `#129–150`→**6.J**; `#151–163`→**6.M**; `#164–176`→**6.N**; `#177–178` & `#187–188`→**6.F**; `#179–186` & `#189`→**6.I**; `#190–198`→**6.O**.
- For module rows present in `docs/Modul UmrohOS - MosCoW.csv`, the CSV `MoSCoW` column is the default **backlog priority** (`Must Have` → `MH-V1`, `Should Have` → `SH`, `Could Have` → `CH`) unless `docs/06-features/*` decides otherwise.
- **Phase 6 `Slice` / `Task Code`** aligns with **[05](./05-slice-engineering-checklist-and-task-codes.md)** (**Phase 6 — Depth backlog**). Phase 6 table rows carry domain task codes (not one row = one PR).

---

## Backlog ID format

Pattern:

`BL-{DOMAIN}-{NNN}`

Examples:

- `BL-IAM-001`
- `BL-CAT-004`
- `BL-CAT-014`
- `BL-FE-CAT-001`
- `BL-BOOK-007`
- `BL-PAY-003`
- `BL-B2C-001`
- `BL-JMJ-001`
- `BL-JNT-001` … **`JNT`** = Joint slice-freeze contract (**05** `Sx-J-*`); ids **001–015** span S1→S5 (see freeze index above).
- `BL-EGV-001` … **`EGV`** = Engineering gate review (**05** `Sx-E-01` only where that row exists).
- `BL-LGV-001` … **`LGV`** = Lead / UX gate (**05** `Sx-L-01` only where that row exists).

Suggested workflow status:

`todo -> in_progress -> in_review -> done`

### `Exec seq` (numeric order)

For **Phases 0–5**, numbers usually follow a loose **Phase × 100 + small sequence** pattern:

- **Phase 0** → `000–099` (example: `005`)
- **Phase 1** → `100–199` (example: `110`)
- **Phase 2** → `200–299`
- and so on.

**Phase 6 (depth):** uses **`600+`** through **854** as a dedicated track after core; subsections **6.G–6.O** do not have to follow `Phase × 100` per row — per-row order follows **`Exec seq`**.

Smaller numbers run first. Multiple rows may share the same number when safe to parallelize across domains.

---

## Phase 0 — Engineering bootstrap (S0)


| Phase | Feature ref | Summary                                                                | Priority | Exec seq | Slice | Task Code | Backlog ID    | Owner | Status | Blocked by gate | Acceptance (short)                                                       |
| ----- | ----------- | ---------------------------------------------------------------------- | -------- | -------- | ----- | --------- | ------------- | ----- | ------ | --------------- | ------------------------------------------------------------------------ |
| 0     | S0-J-01     | Contract templates + README `docs/contracts/`                          | MH-MVP   | 005      | S0    | S0-J-01   | BL-ENG-001    | J     | done   | —               | `slice-Sx.md` template + README merged                                   |
| 0     | S0-J-02     | Branch strategy + merge ownership                                      | MH-MVP   | 006      | S0    | S0-J-02   | BL-ENG-002    | J     | done   | BL-ENG-001      | Written merge rules                                                      |
| 0     | S0-J-03     | DoR/DoD per PR                                                         | MH-MVP   | 007      | S0    | S0-J-03   | BL-ENG-003    | J     | done   | BL-ENG-001      | Short DoR/DoD published                                                  |
| 0     | S0-L-01     | Public vs internal UI route matrix (S1)                                | MH-MVP   | 008      | S0    | S0-L-01   | BL-FE-PLN-001 | L     | done   | BL-ENG-001      | Matrix + contract gaps in `docs/contracts/slice-S1.md` § UI route matrix |
| 0     | S0-E-01     | Service ownership S1–S2 (PR owner)                                     | MH-MVP   | 009      | S0    | S0-E-01   | BL-ENG-004    | E     | done   | BL-ENG-001      | Service vs owner table                                                   |
| 0     | S0-J-04     | Repo scaffold sweep vs ADR-0004                                        | MH-MVP   | 010      | S0    | S0-J-04   | BL-ENG-005    | J     | done   | BL-ENG-001      | Drift fixed or exceptions documented                                     |
| 0     | S0-J-05     | OpenTelemetry baseline (all Go services)                               | MH-MVP   | 011      | S0    | S0-J-05   | BL-ENG-006    | J     | done   | BL-ENG-001      | Traces/logs correlated; collector path documented                        |
| 0     | S0-J-06     | CI path filters for workflows                                          | MH-MVP   | 012      | S0    | S0-J-06   | BL-ENG-007    | J     | done   | BL-ENG-001      | Expensive jobs skip unrelated paths                                      |
| 0     | S0-J-07     | Migration pipeline + docs (ADR-0007)                                   | MH-MVP   | 013      | S0    | S0-J-07   | BL-ENG-008    | J     | done   | BL-ENG-001      | `make`/CI migrate path documented                                        |
| 0     | S0-J-08     | `/livez` + `/readyz` (+ optional db-tx diagnostic) on every Go service | MH-MVP   | 014      | S0    | S0-J-08   | BL-ENG-009    | J     | done   | BL-ENG-001      | All 11 services healthy endpoints in dev                                 |
| 0     | S0-J-09     | E2E skeleton (Playwright, ADR-0008)                                    | MH-MVP   | 015      | S0    | S0-J-09   | BL-ENG-010    | J     | done   | BL-ENG-001      | Minimal spec + run/CI notes merged                                       |


---

## Phase 1 — Discover + draft booking (S1)

Slice **S1** engineering-freeze rows from **05** are listed first (`BL-JNT-001`–`004`, then `BL-EGV-001`, `BL-LGV-001`). Dependencies chain Joint cards; **E** / **L** rows depend on the Joint ids shown.

| Phase | Feature ref | Summary                                                | Priority | Exec seq | Slice | Task Code | Backlog ID     | Owner | Status | Blocked by gate                    | Acceptance (short)                                                  |
| ----- | ----------- | ------------------------------------------------------ | -------- | -------- | ----- | --------- | -------------- | ----- | ------ | ---------------------------------- | ------------------------------------------------------------------- |
| 1     | S1-J-01     | Public catalog API contract (list/detail + seats read) | MH-MVP   | 101      | S1    | S1-J-01   | BL-JNT-001     | J     | done   | —                                  | `slice-S1.md` § Catalog + **§ Changelog** (2026-04-20, `S1-J-01`) — matches **05** |
| 1     | S1-J-02     | Draft booking API contract `POST /v1/bookings`         | MH-MVP   | 102      | S1    | S1-J-02   | BL-JNT-002     | J     | done   | BL-JNT-001                         | `slice-S1.md` § Booking + **§ Changelog** (2026-04-21, `S1-J-02`) — matches **05** |
| 1     | S1-J-03     | `ReserveSeats` / `ReleaseSeats` contract               | MH-MVP   | 103      | S1    | S1-J-03   | BL-JNT-003     | J     | done   | BL-JNT-001; BL-JNT-002             | `slice-S1.md` § Inventory + **§ Changelog** (2026-04-21, `S1-J-03`) — matches **05** |
| 1     | S1-J-04     | S1 booking states + Q006 / document MVP gate           | MH-MVP   | 104      | S1    | S1-J-04   | BL-JNT-004     | J     | done   | BL-JNT-001; BL-JNT-002; BL-JNT-003 | `slice-S1.md` § Booking States + Q006 pin + **§ Changelog** (2026-04-21, `S1-J-04`) — matches **05** |
| 1     | S1-E-01     | Seat concurrency + DB transaction review (for `S1-J-03`) | MH-MVP   | 105      | S1    | S1-E-01   | BL-EGV-001     | E     | done   | BL-JNT-003                         | `slice-S1.md` § Inventory — **S1-E-01 engineering review** approved (2026-04-22) |
| 1     | S1-L-01     | S1 wireframe / screen list (URL + main components)     | MH-MVP   | 106      | S1    | S1-L-01   | BL-LGV-001     | L     | done   | BL-JNT-001                         | `slice-S1.md` § **S1-L-01 — UI screen inventory** + `apps/core-web` `(b2c)/` routes — matches **05** (no Figma) |
| 1     | F1-W1       | Internal login + refresh token flow                    | MH-MVP   | 110      | S1    | S1-E-04   | BL-IAM-001     | E     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Login succeeds, refresh works, unauthorized → 401                   |
| 1     | F1-W3       | `CheckPermission` middleware for internal routes       | MH-MVP   | 111      | S1    | S1-E-04   | BL-IAM-002     | E     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Finance routes denied for non-finance roles                         |
| 1     | F1-W5       | Basic suspend/revoke session                           | MH-MVP   | 112      | S1    | S1-E-04   | BL-IAM-003     | E     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Suspended user cannot access again                                  |
| 1     | F1-AC       | Audit log for state-changing calls                     | MH-MVP   | 113      | S1    | S1-E-04   | BL-IAM-004     | E     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Create/update booking actions recorded in audit                     |
| 1     | F2-W2       | Read model active package + departure                  | MH-MVP   | 120      | S1    | S1-E-02   | BL-CAT-001     | E     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | `list/detail` only shows valid package/departure                    |
| 1     | F2-W3       | Departure data includes seat cap + status              | MH-MVP   | 121      | S1    | S1-E-02   | BL-CAT-002     | E     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Departure detail exposes consistent remaining seats                 |
| 1     | F2-W6       | Atomic `ReserveSeats` + `ReleaseSeats`                 | MH-MVP   | 122      | S1    | S1-E-03   | BL-CAT-003     | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Last-seat race safe, no oversell                                    |
| 1     | F2-AC       | Public catalog endpoints for B2C                       | MH-MVP   | 123      | S1    | S1-E-02   | BL-CAT-004     | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | B2C can browse packages without contract errors                     |
| 1     | F2-W12      | Staff catalog write REST MVP (package + departure; Bearer + permission) | MH-MVP   | 124      | S1    | S1-E-07   | BL-CAT-014     | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-IAM-001..004 | Mutating catalog routes gated; request/response shapes align with public read models |
| 1     | F2-W13      | Internal console — catalog CRUD MVP (packages + departures)            | MH-MVP   | 125      | S1    | S1-L-06   | BL-FE-CAT-001  | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-IAM-001..004; BL-CAT-014 | Staff can maintain catalog via FE (replaces ad-hoc DB edits for MVP ops) |
| 1     | F1-W1       | Internal console login page (`/console/login`) + session entry flow     | MH-MVP   | 126      | S1    | S1-L-07   | BL-FE-CONSOLE-001 | L  | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-IAM-001..004 | Valid credentials enter console; invalid credentials stay on login with error feedback |
| 1     | F1-W3       | Internal console shell (`/console`) with full sidemenu                  | MH-MVP   | 127      | S1    | S1-L-08   | BL-FE-CONSOLE-002 | L  | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-IAM-001..004; BL-FE-CONSOLE-001 | Authenticated staff sees admin shell |
| 1     | F4-W1       | Create draft booking from B2C flow                     | MH-MVP   | 130      | S1    | S1-E-03   | BL-BOOK-001    | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Booking `draft` saved with minimum fields                           |
| 1     | F4-W2       | Stamp channel attribution (`b2c_self`/`b2b_agent`)     | MH-MVP   | 131      | S1    | S1-E-03   | BL-BOOK-002    | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Booking stores `channel` + `agent_id` when present                  |
| 1     | F4-W4       | State machine through `pending_payment`/`expired`      | MH-MVP   | 132      | S1    | S1-E-03   | BL-BOOK-003    | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Valid status transitions without skipping states                    |
| 1     | F4-W8       | Basic submit validation (active package, enough seats) | MH-MVP   | 133      | S1    | S1-E-03   | BL-BOOK-004    | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Hard fail if seats gone / departure invalid                         |
| 1     | F4-W8       | Minimum document gate (Q006) on submit                 | MH-MVP   | 134      | S1    | S1-E-03   | BL-BOOK-005    | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Missing doc → clear error per pilgrim/doc kind                      |
| 1     | F4-W6       | `ValidateMahram` integration as warning                | MH-MVP   | 135      | S1    | S1-E-03   | BL-BOOK-006    | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Mahram result stored, does not block submit                         |
| 1     | F4-UI       | Catalog → detail → booking form UI                     | MH-MVP   | 140      | S1    | S1-L-02   | BL-FE-BOOK-001 | L     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | User can reach booking form from catalog                            |
| 1     | F4-UI       | FE integration to catalog API                          | MH-MVP   | 141      | S1    | S1-L-03   | BL-FE-BOOK-002 | L     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | FE calls list/detail per contract                                   |
| 1     | F4-UI       | FE integration create draft booking                    | MH-MVP   | 142      | S1    | S1-L-04   | BL-FE-BOOK-003 | L     | done   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Form submit creates draft booking successfully                      |
| 1     | S1-E-08     | Monitoring migration — register `grpc.health.v1.Health` on all 10 backends; install `grpc_health_probe` in Dockerfiles; wire docker-compose healthchecks (gateway keeps REST healthcheck). Metrics split out — `/metrics` stays on REST port until each service's REST package is removed per `BL-REFACTOR-*` | MH-MVP | 143 | S1 | S1-E-08 | BL-MON-001 | E | done | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | All 10 backends respond SERVING to `grpc_health_probe`; `docker compose up` brings the stack healthy; gateway REST healthcheck green |
| 1     | F1-W7       | Edge auth at `gateway-svc` — `iam_grpc_adapter` + Bearer middleware + fail-closed 502 on `iam-svc` unreachable | MH-MVP | 144 | S1 | S1-E-09 | BL-GTW-001 | E | done | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Valid bearer → 200; invalid → 401; `iam-svc` down → 502; trace shows gateway→iam `ValidateToken` span |
| 1     | F2-W2       | Gateway `catalog_grpc_adapter` + REST routes for `/v1/packages*` and `/v1/package-departures/{id}`; e2e migrates to `gateway-svc:4000` | MH-MVP | 145 | S1 | S1-E-10 | BL-GTW-002 | E | done | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-MON-001; BL-GTW-001 | Public catalog reads succeed through gateway; `02e-catalog-svc-read.spec.ts` passes against `gateway-svc:4000` |
| 1     | S1-E-11     | Remove `catalog-svc` REST package (ADR 0009); migrate `/system/diagnostics/db-tx` to a gRPC method | MH-MVP | 146 | S1 | S1-E-11 | BL-REFACTOR-001 | E | done | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-GTW-002 | `services/catalog-svc/api/rest_oapi/` is gone; gateway-routed e2e still green; catalog runs gRPC-only |
| 1     | S1-E-12     | Move iam client-facing REST (`/v1/sessions*`, `/v1/me*`, `/v1/users*`) to `gateway-svc`; add matching iam gRPC RPCs | MH-MVP | 147 | S1 | S1-E-12 | BL-IAM-018 | E | done | BL-JNT-001..004; BL-EGV-001; BL-LGV-001; BL-MON-001; BL-GTW-001 | Login/logout/refresh/me flow through `gateway-svc:4000`; `services/iam-svc/api/rest_oapi/` removed |
| 1     | S1-E-13     | Remove `booking-svc` REST package (ADR 0009); retire gateway `booking_rest_adapter` + `/v1/booking/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 148 | S1 | S1-E-13 | BL-REFACTOR-002 | E | done | BL-REFACTOR-001 | `services/booking-svc/api/rest_oapi/` is gone; booking runs gRPC-only; gateway no longer carries a booking REST adapter |
| 1     | S1-E-13     | Remove `crm-svc` REST package (ADR 0009); retire gateway `crm_rest_adapter` + `/v1/crm/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 149 | S1 | S1-E-13 | BL-REFACTOR-003 | E | done | BL-REFACTOR-001 | `services/crm-svc/api/rest_oapi/` is gone; crm runs gRPC-only; gateway no longer carries a crm REST adapter |
| 1     | S1-E-13     | Remove `jamaah-svc` REST package (ADR 0009); retire gateway `jamaah_rest_adapter` + `/v1/jamaah/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 150 | S1 | S1-E-13 | BL-REFACTOR-004 | E | done | BL-REFACTOR-001 | `services/jamaah-svc/api/rest_oapi/` is gone; jamaah runs gRPC-only; gateway no longer carries a jamaah REST adapter |
| 1     | S1-E-13     | Remove `logistics-svc` REST package (ADR 0009); retire gateway `logistics_rest_adapter` + `/v1/logistics/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 151 | S1 | S1-E-13 | BL-REFACTOR-005 | E | done | BL-REFACTOR-001 | `services/logistics-svc/api/rest_oapi/` is gone; logistics runs gRPC-only; gateway no longer carries a logistics REST adapter |
| 1     | S1-E-13     | Remove `ops-svc` REST package (ADR 0009); retire gateway `ops_rest_adapter` + `/v1/ops/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 152 | S1 | S1-E-13 | BL-REFACTOR-006 | E | done | BL-REFACTOR-001 | `services/ops-svc/api/rest_oapi/` is gone; ops runs gRPC-only; gateway no longer carries an ops REST adapter |
| 1     | S1-E-13     | Remove `payment-svc` REST package (ADR 0009); retire gateway `payment_rest_adapter` + `/v1/payment/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 153 | S1 | S1-E-13 | BL-REFACTOR-007 | E | done | BL-REFACTOR-001 | `services/payment-svc/api/rest_oapi/` is gone; payment runs gRPC-only; gateway no longer carries a payment REST adapter |
| 1     | S1-E-13     | Remove `visa-svc` REST package (ADR 0009); retire gateway `visa_rest_adapter` + `/v1/visa/system/live`; migrate `util/monitoring` to OTLP metrics push | MH-MVP | 154 | S1 | S1-E-13 | BL-REFACTOR-008 | E | done | BL-REFACTOR-001 | `services/visa-svc/api/rest_oapi/` is gone; visa runs gRPC-only; gateway no longer carries a visa REST adapter |
| 1     | S1-E-14     | finance-svc ADR 0009 realignment — delete local REST + bearer-auth middleware + local `iam_grpc_adapter`; scaffold a uniform `Ping` RPC on every downstream backend; add gateway `finance_grpc_adapter`; move `/v1/finance/ping` to gateway REST behind new `RequirePermission` middleware (gateway calls `iam.CheckPermission`); delete `/v1/finance/system/live`; migrate finance metrics to OTLP push; migrate `02b`/`02c` e2e specs to `gateway-svc:4000`; drop scaffold-noise `rpc Healthz` + `rpc DiagnosticsDbTx` + `service.DbTxDiagnostic` + related sqlc queries from every backend; add `GET /v1/system/backends` on gateway that proxies `grpc.health.v1.Health.Check` to every backend + rewire core-web's status page to poll it | MH-MVP | 155 | S1 | S1-E-14 | BL-IAM-019 | E | done | BL-IAM-018; BL-REFACTOR-002..008 | finance-svc runs gRPC-only; BL-IAM-002 permission-gate coverage passes through `gateway-svc:4000`; `services/finance-svc/api/rest_oapi/` + `services/gateway-svc/adapter/finance_rest_adapter/` removed; every backend exposes a uniform `Ping` RPC; `Healthz` / `DiagnosticsDbTx` scaffold-noise RPCs no longer exist; core-web status page polls one `GET /v1/system/backends` aggregate endpoint |


---

## Phase 2 — Get paid (S2)

**S2** freeze rows (`BL-JNT-005`–`008`, `BL-EGV-002`, `BL-LGV-002`) precede implementation work.

| Phase | Feature ref | Summary                                         | Priority | Exec seq | Slice | Task Code | Backlog ID    | Owner | Status | Blocked by gate                    | Acceptance (short)                                |
| ----- | ----------- | ----------------------------------------------- | -------- | -------- | ----- | --------- | ------------- | ----- | ------ | ---------------------------------- | ------------------------------------------------- |
| 2     | S2-J-01     | Invoice + VA issue contract (`POST` invoice / VA) | MH-MVP   | 201      | S2    | S2-J-01   | BL-JNT-005    | J     | todo   | —                                  | `slice-S2.md` invoice/VA section merged + reviewed |
| 2     | S2-J-02     | Webhook contract (signature, body, dedupe, codes) | MH-MVP   | 202      | S2    | S2-J-02   | BL-JNT-006    | J     | todo   | BL-JNT-005                       | `slice-S2.md` § Webhook merged + reviewed           |
| 2     | S2-J-03     | Booking callback contract (status + idempotency) | MH-MVP   | 203      | S2    | S2-J-03   | BL-JNT-007    | J     | todo   | BL-JNT-005; BL-JNT-006           | `slice-S2.md` § Booking integration merged + reviewed |
| 2     | S2-J-04     | Stub `payment-svc` or `MOCK_GATEWAY` toggle documented | MH-MVP   | 204      | S2    | S2-J-04   | BL-JNT-008    | J     | todo   | BL-JNT-005; BL-JNT-006; BL-JNT-007 | Stub or env flag documented per **05**            |
| 2     | S2-E-01     | Payment DB tables review (`payment-svc` data model) | MH-MVP   | 205      | S2    | S2-E-01   | BL-EGV-002    | E     | todo   | BL-JNT-005; BL-JNT-006; BL-JNT-007; BL-JNT-008 | Migration reviewed per **05**                      |
| 2     | S2-L-01     | Checkout UI spec (VA/QR + polling strategy)          | MH-MVP   | 206      | S2    | S2-L-01   | BL-LGV-002    | L     | todo   | BL-JNT-005                       | Description in contract or UI note per **05**       |
| 2     | F5-W1       | Create invoice + issue VA + TTL                 | MH-MVP   | 210      | S2    | S2-E-02   | BL-PAY-001    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | VA issued, `expires_at` stored per config         |
| 2     | F5-W1       | Gateway selection + fallback rule (Q013)        | MH-MVP   | 211      | S2    | S2-E-02   | BL-PAY-002    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Failover only on timeout/5xx                      |
| 2     | F5-W2       | Webhook signature verification                  | MH-MVP   | 220      | S2    | S2-E-02   | BL-PAY-003    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Bad signature → 401, no business update           |
| 2     | F5-W2       | Webhook idempotency (`gateway_txn_id`)          | MH-MVP   | 221      | S2    | S2-E-02   | BL-PAY-004    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Webhook replay → safe no-op                       |
| 2     | F5-W2       | Update `paid_amount` + signal `MarkBookingPaid` | MH-MVP   | 222      | S2    | S2-E-02   | BL-PAY-005    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Booking status syncs partial/paid                 |
| 2     | F5-W5       | Reconciliation cron for missed webhooks         | MH-MVP   | 223      | S2    | S2-E-03   | BL-PAY-006    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Dropped webhooks recover on next cycle            |
| 2     | F5-W8       | Basic refund wired to booking cancellation      | MH-MVP   | 224      | S2    | S2-E-02   | BL-PAY-007    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Cancel → refund flow recorded and idempotent      |
| 2     | F5-W9       | FX snapshot + rounding rule Q001                | MH-MVP   | 225      | S2    | S2-E-02   | BL-PAY-008    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Snapshot immutable after first payment            |
| 2     | F5-UI       | Checkout page shows VA + countdown              | MH-MVP   | 230      | S2    | S2-L-02   | BL-FE-PAY-001 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | User sees VA account, amount, expiry in real time |
| 2     | F5-UI       | FE wiring booking → payment call                | MH-MVP   | 231      | S2    | S2-L-03   | BL-FE-PAY-002 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | From draft can proceed to issue payment           |
| 2     | F10-W14     | Payment link generator (CS closing)             | MH-V1    | 235      | S2    | S2-E-02   | BL-PAY-020    | E     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | CS can issue link/VA for existing booking         |
| 2     | F5-E2E      | End-to-end test stub → real gateway             | MH-MVP   | 240      | S2    | S2-J-05   | BL-QA-001     | J     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Draft→pay→paid scenario passes E2E                |


---

## Phase 3 — Fulfillment + minimum post-pay finance (S3)

**S3** freeze rows (`BL-JNT-009`–`011`, `BL-EGV-003`). No **`S3-L-01`** row in **05** engineering freeze.

| Phase | Feature ref | Summary                                         | Priority | Exec seq | Slice | Task Code | Backlog ID    | Owner | Status | Blocked by gate           | Acceptance (short)                                  |
| ----- | ----------- | ----------------------------------------------- | -------- | -------- | ----- | --------- | ------------- | ----- | ------ | ------------------------- | --------------------------------------------------- |
| 3     | S3-J-01     | Events `payment.received` / `booking.paid_in_full` payload contract | MH-V1    | 301      | S3    | S3-J-01   | BL-JNT-009    | J     | todo   | —                         | `slice-S3.md` event payload merged + reviewed       |
| 3     | S3-J-02     | Minimal fulfillment task contract (status, assignee) | MH-V1    | 302      | S3    | S3-J-02   | BL-JNT-010    | J     | todo   | BL-JNT-009                | `slice-S3.md` fulfillment section merged + reviewed   |
| 3     | S3-J-03     | Minimal journal contract (placeholder accounts + rules) | MH-V1    | 303      | S3    | S3-J-03   | BL-JNT-011    | J     | todo   | BL-JNT-009; BL-JNT-010    | `slice-S3.md` journal section merged + reviewed     |
| 3     | S3-E-01     | Posting load vs refund review (engineering)          | MH-V1    | 304      | S3    | S3-E-01   | BL-EGV-003    | E     | todo   | BL-JNT-009; BL-JNT-010; BL-JNT-011 | Comment on contract per **05**                      |
| 3     | F3-W1       | Upload pilgrim documents (KTP/passport/photo)   | MH-V1    | 310      | S3    | S3-E-02   | BL-DOC-001    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Documents stored + lifecycle status works           |
| 3     | F3-W4       | Passport OCR + manual review fallback           | MH-V1    | 311      | S3    | S3-E-02   | BL-DOC-002    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | OCR result stored, low confidence enters review     |
| 3     | F3-W3       | Verification queue approve/reject + audit       | MH-V1    | 312      | S3    | S3-E-02   | BL-DOC-003    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Approve/reject audited + reject notification sent   |
| 3     | F7-W5       | Manifest generator (PDF/XLSX)                   | MH-V1    | 320      | S3    | S3-E-02   | BL-OPS-001    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Manifest can be generated per departure             |
| 3     | F7-W2       | Smart grouping room allocation run + commit     | MH-V1    | 321      | S3    | S3-E-02   | BL-OPS-002    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Grouping run produces valid allocation              |
| 3     | F7-W6       | ID card + luggage tag issuance signed QR        | MH-V1    | 322      | S3    | S3-E-02   | BL-OPS-003    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | QR valid/verified, tamper rejected                  |
| 3     | F8-W10      | Trigger fulfillment only on `paid_in_full`      | MH-V1    | 330      | S3    | S3-E-02   | BL-LOG-001    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Fulfillment queue excludes non-paid bookings        |
| 3     | F8-W11      | Shipment + tracking number + WA notify          | MH-V1    | 331      | S3    | S3-E-02   | BL-LOG-002    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Shipment produces tracking + notification           |
| 3     | F8-W12      | Self-pickup QR single-use                       | MH-V1    | 332      | S3    | S3-E-02   | BL-LOG-003    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Pickup QR valid once + expiry                       |
| 3     | F9-W2       | Post payment-receipt journal (deferred revenue) | MH-V1    | 340      | S3    | S3-E-03   | BL-FIN-001    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Dr Bank / Cr pilgrim liability automatic            |
| 3     | F9-W4       | Auto-AP from GRN (synchronous)                  | MH-V1    | 341      | S3    | S3-E-03   | BL-FIN-002    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | GRN fails if AP posting fails                       |
| 3     | F9-W9       | Double-entry journal engine + idempotent source | MH-V1    | 342      | S3    | S3-E-03   | BL-FIN-003    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | No unbalanced / duplicate-source journals           |
| 3     | F9-W10      | Revenue recognition trigger on departure event  | MH-V1    | 343      | S3    | S3-E-03   | BL-FIN-004    | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Revenue recognized per Q043 trigger, not on payment |
| 3     | F7-UI       | Ops board: fulfillment + manifest summary       | MH-V1    | 350      | S3    | S3-L-02   | BL-FE-OPS-001 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | UI shows main ops status per booking                |


---

## Phase 4 — Basic growth loop (S4)

**S4** freeze rows (`BL-JNT-012`–`013`) only — no separate **`S4-E-01` / `S4-L-01`** in **05** engineering freeze.

| Phase | Feature ref | Summary                                  | Priority | Exec seq | Slice | Task Code | Backlog ID    | Owner | Status | Blocked by gate  | Acceptance (short)                      |
| ----- | ----------- | ---------------------------------------- | -------- | -------- | ----- | --------- | ------------- | ----- | ------ | ---------------- | --------------------------------------- |
| 4     | S4-J-01     | Lead schema + UTM snapshot + attribution (Q019/Q057) | MH-V1    | 401      | S4    | S4-J-01   | BL-JNT-012    | J     | todo   | —                | `slice-S4.md` lead/UTM section merged + reviewed |
| 4     | S4-J-02     | Booking → CRM events (names + payload) | MH-V1    | 402      | S4    | S4-J-02   | BL-JNT-013    | J     | todo   | BL-JNT-012       | `slice-S4.md` event section merged + reviewed |
| 4     | F10-W3      | Lead tracker + tagging (basic CRM)       | MH-V1    | 410      | S4    | S4-E-02   | BL-CRM-001    | E     | todo   | BL-JNT-012..013 | Lead stored + list/filter/tag works     |
| 4     | F10-W8      | Attribution + UTM reconciliation (basic) | MH-V1    | 411      | S4    | S4-E-02   | BL-CRM-002    | E     | todo   | BL-JNT-012..013 | UTM stored + consistent into booking    |
| 4     | F10-W4      | CS round-robin + basic SLA               | MH-V1    | 412      | S4    | S4-E-02   | BL-CRM-003    | E     | todo   | BL-JNT-012..013 | Lead distribution + minimum SLA works   |
| 4     | F10-UI      | Lead capture form (public/internal)      | MH-V1    | 420      | S4    | S4-L-01   | BL-FE-CRM-001 | L     | todo   | BL-JNT-012..013 | Form submit succeeds + basic validation |


---

## Phase 5 — Hardening & readiness (S5)

**S5** freeze rows (`BL-JNT-014`–`015`) only — no separate **`S5-E-01` / `S5-L-01`** in **05** engineering freeze (implementation checklist still uses **`S5-E-01` / `S5-L-01`** task codes).

| Phase | Feature ref | Summary                                          | Priority | Exec seq | Slice | Task Code | Backlog ID    | Owner | Status | Blocked by gate  | Acceptance (short)                                     |
| ----- | ----------- | ------------------------------------------------ | -------- | -------- | ----- | --------- | ------------- | ----- | ------ | ---------------- | ------------------------------------------------------ |
| 5     | S5-J-01     | Mandatory UAT scenarios (MVP gates) checklist     | MH-V1    | 501      | S5    | S5-J-01   | BL-JNT-014    | J     | todo   | —                | Checklist in `slice-S5.md` merged + reviewed           |
| 5     | S5-J-02     | Bug severity matrix + fix SLA                     | MH-V1    | 502      | S5    | S5-J-02   | BL-JNT-015    | J     | todo   | BL-JNT-014       | Table in `slice-S5.md` merged + reviewed             |
| 5     | F9-W15      | Basic reports: balance sheet + P&L + cash flow   | MH-V1    | 510      | S5    | S5-E-01   | BL-FIN-005    | E     | todo   | BL-JNT-014..015 | Three main reports runnable for current period         |
| 5     | F9-W18      | Finance audit trail + anti-delete journals       | MH-V1    | 511      | S5    | S5-E-01   | BL-FIN-006    | E     | todo   | BL-JNT-014..015 | Journal delete rejected; corrections via counter-entry |
| 5     | F9-UI       | Finance view: journal/payment status summary     | MH-V1    | 520      | S5    | S5-L-01   | BL-FE-FIN-001 | L     | todo   | BL-JNT-014..015 | Finance can trace booking → invoice → journal          |
| 5     | QA          | UAT core checklist + permission/audit regression | MH-V1    | 530      | S5    | S5-L-01   | BL-QA-002     | L     | todo   | BL-JNT-014..015 | Scenario list passes + evidence                        |
| 5     | QA          | UAT payment/finance/logistics checklist          | MH-V1    | 531      | S5    | S5-E-01   | BL-QA-003     | E     | todo   | BL-JNT-014..015 | Scenario list passes + evidence                        |
| 5     | F1-W7 (defense-in-depth) | Gateway↔backend trust contract — signed per-hop header or mTLS between gateway and each downstream backend. Closes the defense-in-depth gap ADR 0009 D2 deferred; replaces `insecure.NewCredentials()` on the in-cluster gRPC dials. | SH | 540 | S5 | S5-E-02 | BL-GTW-100 | E | deferred | S1-E-14 | Gateway↔backend hops authenticate end-to-end (not just at the edge); `insecure.NewCredentials()` removed from every backend dial; compose stack still boots healthy; e2e suite green |


---

## Phase 6 — Depth expansion (after core is stable)

> This section follows **after** Phases 1–5 pass integration gates.  
> `Exec seq` is intentionally **600+** so depth work does not “skip ahead” of foundations.  
> After **6.E** (through `Exec seq` **652**), the next ranges run in order: **6.G** **653–668** → **6.F** **669–672** → **6.H** **673–722** → **6.I** **723–731** → **6.J** **732–753** → **6.K** **754–795** → **6.L** **796–819** → **6.M** **820–832** → **6.N** **833–845** → **6.O** **846–854** (all CSV rows **#1–#202** covered).

### 6.A — Marketing/CRM depth (F10)


| Phase | Feature ref | Summary                                    | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate                    | Acceptance (short)                                |
| ----- | ----------- | ------------------------------------------ | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ---------------------------------- | ------------------------------------------------- |
| 6     | F10-W1      | Agent onboarding + E-KYC + e-MoU           | MH-V1    | 610      | S4    | S4-L-02   | BL-CRM-010 | L     | todo   | BL-JNT-012..013                   | Register → approve → active agent role flow       |
| 6     | F10-W2      | Replica site + share UTM + tracking        | MH-V1    | 611      | S4    | S4-L-02   | BL-CRM-011 | L     | todo   | BL-JNT-012..013                   | Replica renders catalog + lead tracking           |
| 6     | F10-W9      | Commission wallet (balance + basic status) | MH-V1    | 612      | S4    | S4-E-03   | BL-CRM-012 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Commission balance consistent with payment events |


Official F10 CSV module splits live in **6.H** (`BL-CRM-017`–`BL-CRM-066`). Rows **6.A** (`BL-CRM-010`–`012`) stay a short **integration package** bridging several modules at once.

### 6.B — Finance depth (F9)


| Phase | Feature ref | Summary                        | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate           | Acceptance (short)                        |
| ----- | ----------- | ------------------------------ | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ------------------------- | ----------------------------------------- |
| 6     | F9-W5       | AP disbursement ladder minimum | MH-V1    | 620      | S3    | S3-E-07   | BL-FIN-010 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Batch AP + approval + audit               |
| 6     | F9-W17      | Basic AR/AP aging alerts       | MH-V1    | 621      | S3    | S3-E-07   | BL-FIN-011 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Aging buckets visible + basic alert rules |


### 6.C — Warehouse / procurement depth (F8)


| Phase | Feature ref | Summary                 | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate           | Acceptance (short)                        |
| ----- | ----------- | ----------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ------------------------- | ----------------------------------------- |
| 6     | F8-W1       | PR + budget gate        | MH-V1    | 630      | S3    | S3-E-05   | BL-LOG-010 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Over-budget PR rejected                   |
| 6     | F8-W4       | GRN + QC + auto-AP sync | MH-V1    | 631      | S3    | S3-E-05   | BL-LOG-011 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | GRN rolls back when finance posting fails |
| 6     | F8-W7       | Kit assembly atomic     | MH-V1    | 632      | S3    | S3-E-05   | BL-LOG-012 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Assembly all-or-nothing                   |


### 6.D — Field operations depth (F7)


| Phase | Feature ref | Summary                       | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate           | Acceptance (short)         |
| ----- | ----------- | ----------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ------------------------- | -------------------------- |
| 6     | F7-W7       | ALL system scan + idempotency | MH-V1    | 640      | S3    | S3-E-04   | BL-OPS-010 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Event scan idempotent      |
| 6     | F7-W10      | Bus boarding scan + roster    | MH-V1    | 641      | S3    | S3-E-04   | BL-OPS-011 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Boarding roster consistent |


### 6.E — Visa pipeline Must (#97) (F6)


| Phase | Feature ref | Summary                                          | Priority | Exec seq | Slice | Task Code | Backlog ID  | Owner | Status | Blocked by gate           | Acceptance (short)                      |
| ----- | ----------- | ------------------------------------------------ | -------- | -------- | ----- | --------- | ----------- | ----- | ------ | ------------------------- | --------------------------------------- |
| 6     | F6-W1       | Readiness auto-transition `WAITING_DOCS → READY` | MH-V1    | 650      | S3    | S3-E-06   | BL-VISA-001 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Status transition recorded + idempotent |
| 6     | F6-W2       | Bulk submit visa all-or-nothing                  | MH-V1    | 651      | S3    | S3-E-06   | BL-VISA-002 | E     | todo   | BL-VISA-001               | Bulk submit atomic per spec             |
| 6     | F6-W3       | Poll status provider + history                   | MH-V1    | 652      | S3    | S3-E-06   | BL-VISA-003 | E     | todo   | BL-VISA-002               | Poll history persisted                  |


### 6.G — Master product CSV modules (#71–#86) (F2)

Per-row priority defaults from `docs/Modul UmrohOS - MosCoW.csv` (`No` 71–86).


| Phase | Feature ref | Summary                                        | Priority | Exec seq | Slice | Task Code | Backlog ID  | Owner | Status | Blocked by gate                    | Acceptance (short)                                              |
| ----- | ----------- | ---------------------------------------------- | -------- | -------- | ----- | --------- | ----------- | ----- | ------ | ---------------------------------- | --------------------------------------------------------------- |
| 6     | F2-CSV-71   | Hotel master data (hotel database)             | MH-V1    | 653      | S1    | S1-E-05   | BL-CAT-005  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | CRUD + hotel reference usable in product composition            |
| 6     | F2-CSV-72   | Guide / muthawwif master data (guide database) | MH-V1    | 654      | S1    | S1-E-05   | BL-CAT-006  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | CRUD + consistent assignment to product/departure               |
| 6     | F2-CSV-73   | Transport & airline masters                    | MH-V1    | 655      | S1    | S1-E-05   | BL-CAT-007  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Carrier/route/mode references available for product composition |
| 6     | F2-CSV-74   | Trip product variants (template + constraints) | MH-V1    | 656      | S1    | S1-E-05   | BL-CAT-008  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Variants do not break public read model + publish validation    |
| 6     | F2-CSV-75   | Financial & retail products (non-core addons)  | SH       | 657      | S1    | S1-E-05   | BL-CAT-009  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Addons separate from core package + clear pricing rules         |
| 6     | F2-CSV-76   | Smart bulk import (spreadsheet + validation)   | SH       | 658      | S1    | S1-E-05   | BL-CAT-010  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Partial import failure safe + per-row error report              |
| 6     | F2-CSV-77   | Bulk update (guarded bulk edit)                | SH       | 659      | S1    | S1-E-05   | BL-CAT-011  | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Bulk update preview + audit + rollback policy                   |
| 6     | F2-CSV-78   | Dynamic flyer generator                        | SH       | 660      | S4    | S4-L-02   | BL-CRM-013  | L     | todo   | BL-JNT-012..013                   | Flyer renders from template + live package data                 |
| 6     | F2-CSV-79   | Omni-flyer (multi-format/channel)              | CH       | 661      | S4    | S4-L-02   | BL-CRM-014  | L     | todo   | BL-JNT-012..013                   | Single content source → multiple output variants                |
| 6     | F2-CSV-80   | Interactive itinerary (shareable)              | SH       | 662      | S4    | S4-L-02   | BL-CRM-015  | L     | todo   | BL-JNT-012..013                   | Itinerary consistent with master itinerary + deep-link          |
| 6     | F2-CSV-81   | Copywriting automation                         | CH       | 663      | S4    | S4-L-02   | BL-CRM-016  | L     | todo   | BL-JNT-012..013                   | Output reviewable + no auto-publish without gate                |
| 6     | F2-CSV-82   | Single-door sync (catalog → agent channels)    | MH-V1    | 664      | S4    | S4-E-04   | BL-CAT-012  | E     | todo   | BL-JNT-012..013                   | Master changes propagate idempotently per agent                 |
| 6     | F2-CSV-83   | Agent auto-update (catalog version + diff)     | MH-V1    | 665      | S4    | S4-E-04   | BL-CAT-013  | E     | todo   | BL-JNT-012..013                   | Agent has version snapshot + safe upgrade path                  |
| 6     | F2-CSV-84   | Cross-channel seat tracking (agent/B2C)        | MH-V1    | 666      | S1    | S1-E-05   | BL-BOOK-007 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Seat state prevents double-sell across channels                 |
| 6     | F2-CSV-85   | Dual dashboard view (role/context switch)      | SH       | 667      | S5    | S5-L-02   | BL-DASH-005 | L     | todo   | BL-JNT-014..015                   | Two display modes consistent with permissions + filters         |
| 6     | F2-CSV-86   | Vendor readiness checklist per departure       | MH-V1    | 668      | S3    | S3-E-04   | BL-OPS-020  | E     | todo   | BL-JNT-009..011; BL-EGV-003          | Checklist item + status + minimum attachment evidence           |


### 6.F — Dashboard Must modules (#177–#178, #187–#188) (F11)

Remaining Dashboard **Should** modules (**#179–#186**, **#189**) are in **6.I** (`BL-DASH-006`–`014`). Module **#85** (*dual dashboard view*) maps to **6.G** (`BL-DASH-005`).


| Phase | Feature ref | Summary                             | Priority | Exec seq | Slice | Task Code | Backlog ID  | Owner | Status | Blocked by gate  | Acceptance (short)                                       |
| ----- | ----------- | ----------------------------------- | -------- | -------- | ----- | --------- | ----------- | ----- | ------ | ---------------- | -------------------------------------------------------- |
| 6     | F11-CSV-177 | Vendor execution (readiness widget) | MH-V1    | 669      | S5    | S5-L-02   | BL-DASH-001 | L     | todo   | BL-OPS-020       | Widget checklist consistent with vendor checklist events |
| 6     | F11-CSV-178 | Seat availability (widget)          | MH-V1    | 670      | S5    | S5-L-02   | BL-DASH-002 | L     | todo   | BL-JNT-014..015 | Real-time seat inventory widget                          |
| 6     | F11-CSV-187 | Instant cash flow (widget)          | MH-V1    | 671      | S5    | S5-L-02   | BL-DASH-003 | L     | todo   | BL-JNT-014..015 | Cash summary consistent with F9                          |
| 6     | F11-CSV-188 | Executive financial report (widget) | MH-V1    | 672      | S5    | S5-L-02   | BL-DASH-004 | L     | todo   | BL-JNT-014..015 | Executive P&L / balance sheet summary                    |


### 6.H — Marketing/CRM & alumni/ZISWAF CSV modules (#25–#70, #199–#202) (F10)

Per-row priority defaults from `docs/Modul UmrohOS - MosCoW.csv` (`MoSCoW` column). Modules **#199–#202** appear in the CSV under *Complementary features & daily app*; domain mapping follows F10 in `docs/06-features/10-marketing-crm-agents.md` (alumni referral + ZISWAF slice).


| Phase | Feature ref | Summary                               | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate                    | Acceptance (short)                                |
| ----- | ----------- | ------------------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ---------------------------------- | ------------------------------------------------- |
| 6     | F10-CSV-25  | Self-service agency registration      | MH-V1    | 673      | S4    | S4-L-02   | BL-CRM-017 | L     | todo   | BL-JNT-012..013                   | Partner register form + status pipeline           |
| 6     | F10-CSV-26  | E-KYC & verification                  | MH-V1    | 674      | S4    | S4-L-02   | BL-CRM-018 | L     | todo   | BL-JNT-012..013                   | KYC upload + verification outcome recorded        |
| 6     | F10-CSV-27  | E-signature                           | SH       | 675      | S4    | S4-L-02   | BL-CRM-019 | L     | todo   | BL-JNT-012..013                   | Digital MoU stored + audit trail                  |
| 6     | F10-CSV-28  | Replica website                       | MH-V1    | 676      | S4    | S4-L-02   | BL-CRM-020 | L     | todo   | BL-JNT-012..013                   | Replica renders catalog + agent branding          |
| 6     | F10-CSV-29  | One-click social sharing              | SH       | 677      | S4    | S4-L-02   | BL-CRM-021 | L     | todo   | BL-JNT-012..013                   | One-click share + UTM preserved                   |
| 6     | F10-CSV-30  | Digital business card                 | CH       | 678      | S4    | S4-L-02   | BL-CRM-022 | L     | todo   | BL-JNT-012..013                   | Digital card generate + share                     |
| 6     | F10-CSV-31  | Content bank                          | SH       | 679      | S4    | S4-L-02   | BL-CRM-023 | L     | todo   | BL-JNT-012..013                   | Content bank browse/search + access control       |
| 6     | F10-CSV-32  | Auto-watermark flyer & share          | SH       | 680      | S4    | S4-L-02   | BL-CRM-024 | L     | todo   | BL-JNT-012..013                   | Agent flyer watermark for WA + export             |
| 6     | F10-CSV-33  | Program documentation gallery         | CH       | 681      | S4    | S4-L-02   | BL-CRM-025 | L     | todo   | BL-JNT-012..013                   | Per-program gallery + publish controls            |
| 6     | F10-CSV-34  | Custom tracking code integration      | CH       | 682      | S4    | S4-L-02   | BL-CRM-026 | L     | todo   | BL-JNT-012..013                   | Agent tracking code installed + validation        |
| 6     | F10-CSV-35  | Leads tracker                         | MH-V1    | 683      | S4    | S4-L-02   | BL-CRM-027 | L     | todo   | BL-JNT-012..013                   | Lead intake + basic source/attribution            |
| 6     | F10-CSV-36  | Reminder & follow-up tagging          | SH       | 684      | S4    | S4-L-02   | BL-CRM-028 | L     | todo   | BL-JNT-012..013                   | Reminders + follow-up tags stored                 |
| 6     | F10-CSV-37  | Commission balance & status           | MH-V1    | 685      | S4    | S4-L-02   | BL-CRM-029 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Commission balance consistent per agent           |
| 6     | F10-CSV-38  | Real-time notifications               | SH       | 686      | S4    | S4-L-02   | BL-CRM-030 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | Real-time commission event notifications          |
| 6     | F10-CSV-39  | Transaction & payout history          | MH-V1    | 687      | S4    | S4-L-02   | BL-CRM-031 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | History + payout request + status                 |
| 6     | F10-CSV-40  | LMS                                   | SH       | 688      | S4    | S4-L-02   | BL-CRM-032 | L     | todo   | BL-JNT-012..013                   | Basic academy courses/modules                     |
| 6     | F10-CSV-41  | Quiz & badges                         | CH       | 689      | S4    | S4-L-02   | BL-CRM-033 | L     | todo   | BL-JNT-012..013                   | Quizzes + badges recorded                         |
| 6     | F10-CSV-42  | Sales scripts                         | SH       | 690      | S4    | S4-L-02   | BL-CRM-034 | L     | todo   | BL-JNT-012..013                   | Tier-searchable sales scripts                     |
| 6     | F10-CSV-43  | Leaderboard                           | CH       | 691      | S4    | S4-L-02   | BL-CRM-035 | L     | todo   | BL-JNT-012..013                   | Leaderboard per configured rules                  |
| 6     | F10-CSV-44  | Push notifications                    | SH       | 692      | S4    | S4-L-02   | BL-CRM-036 | L     | todo   | BL-JNT-012..013                   | Scheduled academy push                            |
| 6     | F10-CSV-45  | Super-view dashboard                  | MH-V1    | 693      | S4    | S4-L-02   | BL-CRM-037 | L     | todo   | BL-JNT-012..013                   | Downline aggregate super-view                     |
| 6     | F10-CSV-46  | Auto tier leveling                    | SH       | 694      | S4    | S4-L-02   | BL-CRM-038 | L     | todo   | BL-JNT-012..013                   | Automatic tier level per rules                    |
| 6     | F10-CSV-47  | Overriding commission                 | MH-V1    | 695      | S4    | S4-L-02   | BL-CRM-039 | L     | todo   | BL-JNT-012..013                   | Commission overrides computed deterministically   |
| 6     | F10-CSV-48  | Ads Manager Lite                      | SH       | 696      | S4    | S4-E-03   | BL-CRM-040 | E     | todo   | BL-JNT-012..013                   | Ad spend + campaign status sync                   |
| 6     | F10-CSV-49  | UTM & campaign attribution management | SH       | 697      | S4    | S4-E-03   | BL-CRM-041 | E     | todo   | BL-JNT-012..013                   | UTM builder + consistency to lead/booking         |
| 6     | F10-CSV-50  | Landing page builder & A/B testing    | SH       | 698      | S4    | S4-E-03   | BL-CRM-042 | E     | todo   | BL-JNT-012..013                   | LP publish + A/B variants + metrics               |
| 6     | F10-CSV-51  | Content planner & calendar            | CH       | 699      | S4    | S4-E-03   | BL-CRM-043 | E     | todo   | BL-JNT-012..013                   | Content calendar + assignment                     |
| 6     | F10-CSV-52  | Content publisher & scheduler         | CH       | 700      | S4    | S4-E-03   | BL-CRM-044 | E     | todo   | BL-JNT-012..013                   | Multi-channel publish schedule                    |
| 6     | F10-CSV-53  | Omni-channel distribution             | CH       | 701      | S4    | S4-E-03   | BL-CRM-045 | E     | todo   | BL-JNT-012..013                   | Omni-channel content distribution                 |
| 6     | F10-CSV-54  | Social media & content analytics      | CH       | 702      | S4    | S4-E-03   | BL-CRM-046 | E     | todo   | BL-JNT-012..013                   | Content analytics + basic export                  |
| 6     | F10-CSV-55  | Bot filter & auto-classification      | SH       | 703      | S4    | S4-E-03   | BL-CRM-047 | E     | todo   | BL-JNT-012..013                   | Bot filter + lead classification                  |
| 6     | F10-CSV-56  | Drip messages                         | SH       | 704      | S4    | S4-E-03   | BL-CRM-048 | E     | todo   | BL-JNT-012..013                   | Drip templates + limits                           |
| 6     | F10-CSV-57  | Moment triggers                       | CH       | 705      | S4    | S4-E-03   | BL-CRM-049 | E     | todo   | BL-JNT-012..013                   | Automated moment triggers                         |
| 6     | F10-CSV-58  | Smart database segmentation           | SH       | 706      | S4    | S4-E-03   | BL-CRM-050 | E     | todo   | BL-JNT-012..013                   | Saved segment queries + preview counts            |
| 6     | F10-CSV-59  | Mass broadcast center                 | SH       | 707      | S4    | S4-E-03   | BL-CRM-051 | E     | todo   | BL-JNT-012..013                   | Mass broadcast consent + rate limits              |
| 6     | F10-CSV-60  | Fair distribution                     | MH-V1    | 708      | S4    | S4-E-03   | BL-CRM-052 | E     | todo   | BL-JNT-012..013                   | Fair lead distribution + audit                    |
| 6     | F10-CSV-61  | Response-speed triggers               | SH       | 709      | S4    | S4-E-03   | BL-CRM-053 | E     | todo   | BL-JNT-012..013                   | SLA triggers + escalation                         |
| 6     | F10-CSV-62  | Lead trail & tagging                  | MH-V1    | 710      | S4    | S4-E-03   | BL-CRM-054 | E     | todo   | BL-JNT-012..013                   | Lead timeline + multi-tagging                     |
| 6     | F10-CSV-63  | Price quote generator                 | SH       | 711      | S4    | S4-E-03   | BL-CRM-055 | E     | todo   | BL-JNT-012..013                   | Quote PDF/link + reference number                 |
| 6     | F10-CSV-64  | Payment link builder                  | MH-V1    | 712      | S4    | S4-E-03   | BL-CRM-056 | E     | todo   | BL-PAY-020                         | Pay link issued for existing booking (CS closing) |
| 6     | F10-CSV-65  | E-approval discounts                  | SH       | 713      | S4    | S4-E-03   | BL-CRM-057 | E     | todo   | BL-JNT-012..013                   | Multi-level discount approval flow                |
| 6     | F10-CSV-66  | Alumni loyalty & referral             | CH       | 714      | S4    | S4-E-03   | BL-CRM-058 | E     | todo   | BL-JNT-012..013                   | Referral code + reward tracking                   |
| 6     | F10-CSV-67  | CS performance dashboard              | SH       | 715      | S4    | S4-E-03   | BL-CRM-059 | E     | todo   | BL-JNT-012..013                   | CS metrics dashboard consistent with SLA          |
| 6     | F10-CSV-68  | ROAS calculator                       | SH       | 716      | S4    | S4-E-03   | BL-CRM-060 | E     | todo   | BL-JNT-012..013                   | ROAS from spend + attributed revenue              |
| 6     | F10-CSV-69  | Retargeting sync                      | CH       | 717      | S4    | S4-E-03   | BL-CRM-061 | E     | todo   | BL-JNT-012..013                   | Retargeting audience sync                         |
| 6     | F10-CSV-70  | Stale prospect radar                  | CH       | 718      | S4    | S4-E-03   | BL-CRM-062 | E     | todo   | BL-JNT-012..013                   | Dormant prospect radar + tasks                    |
| 6     | F10-CSV-199 | Alumni referral hub                   | SH       | 719      | S4    | S4-E-03   | BL-CRM-063 | E     | todo   | BL-JNT-012..013                   | Alumni referral + booking attribution             |
| 6     | F10-CSV-200 | Return-intent savings                 | CH       | 720      | S4    | S4-E-03   | BL-CRM-064 | E     | todo   | BL-JNT-012..013                   | Basic return-intent savings flow                  |
| 6     | F10-CSV-201 | Zakat calculator                      | CH       | 721      | S4    | S4-E-03   | BL-CRM-065 | E     | todo   | BL-JNT-012..013                   | Zakat calculator input + result                   |
| 6     | F10-CSV-202 | Morning charity & infaq               | CH       | 722      | S4    | S4-E-03   | BL-CRM-066 | E     | todo   | BL-JNT-012..013                   | Morning charity/infaq flow + receipt              |


### 6.I — Dashboard CSV modules — remaining Should (#179–#186, #189) (F11)

Per-row priority defaults from `docs/Modul UmrohOS - MosCoW.csv` (`No` **179–186**, **189**). Must modules **#177, #178, #187, #188** are covered in **6.F** (`BL-DASH-001`–`004`); full F11 scope: `docs/06-features/11-dashboards.md`.


| Phase | Feature ref | Summary                              | Priority | Exec seq | Slice | Task Code | Backlog ID  | Owner | Status | Blocked by gate           | Acceptance (short)                                     |
| ----- | ----------- | ------------------------------------ | -------- | -------- | ----- | --------- | ----------- | ----- | ------ | ------------------------- | ------------------------------------------------------ |
| 6     | F11-CSV-179 | Ad budget monitor (widget)           | SH       | 723      | S4    | S4-L-03   | BL-DASH-006 | L     | todo   | BL-JNT-012..013          | Spend vs closings + CPL/CPA consistent with F10        |
| 6     | F11-CSV-180 | CS performance board (widget)        | SH       | 724      | S4    | S4-L-03   | BL-DASH-007 | L     | todo   | BL-JNT-012..013          | CS team metrics + SLA + leaderboard from CRM           |
| 6     | F11-CSV-181 | Live bus radar                       | SH       | 725      | S3    | S3-L-05   | BL-DASH-008 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Fleet map/status from GPS/boarding feed F7             |
| 6     | F11-CSV-182 | Raudhah status (aggregate)           | SH       | 726      | S3    | S3-L-05   | BL-DASH-009 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | % entered Raudhah per departure + drill-down           |
| 6     | F11-CSV-183 | Luggage tracking (aggregate)         | SH       | 727      | S3    | S3-L-05   | BL-DASH-010 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Aggregate luggage position per departure from F7 scans |
| 6     | F11-CSV-184 | Incident report (feed)               | SH       | 728      | S3    | S3-L-05   | BL-DASH-011 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Incident feed + severity filter + HQ notification      |
| 6     | F11-CSV-185 | Warehouse health (widget)            | SH       | 729      | S3    | S3-L-05   | BL-DASH-012 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Stock value + critical vs reorder chart (read F8)      |
| 6     | F11-CSV-186 | Logistics execution monitor (widget) | SH       | 730      | S3    | S3-L-05   | BL-DASH-013 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Paid-unshipped aging + GRN/PO backlog summary          |
| 6     | F11-CSV-189 | Liquidity — AR/AP aging snapshot     | SH       | 731      | S5    | S5-L-02   | BL-DASH-014 | L     | todo   | BL-FIN-011                | Aging buckets + alerts consistent with F9              |


CSV modules **#126–#128** (*Executive dashboard* operational under *Operational & handling*) are split again in **6.K** as `F11-CSV-126`–`128` with `BL-DASH-015`–`017` — aligned with **BL-DASH-012** / **#185** (*warehouse health*) and **BL-DASH-013** / **#186** (*logistics execution monitor*) in the table above.

### 6.J — Finance CSV modules (#129–#150) (F9)

Priority source: `docs/Modul UmrohOS - MosCoW.csv`. Domain: `docs/06-features/09-finance-and-accounting.md`. Rows **6.B** (`BL-FIN-010`, `BL-FIN-011`) remain umbrella rows; this table details **Finance** CSV modules.


| Phase | Feature ref | Summary                                | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate           | Acceptance (short)                                  |
| ----- | ----------- | -------------------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ------------------------- | --------------------------------------------------- |
| 6     | F9-CSV-129  | Automated billing                      | MH-V1    | 732      | S3    | S3-E-07   | BL-FIN-020 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Scheduled billing + consistent receivable status    |
| 6     | F9-CSV-130  | Bank integration                       | MH-V1    | 733      | S3    | S3-E-07   | BL-FIN-021 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Bank feed connected + basic reconciliation          |
| 6     | F9-CSV-131  | Accounts receivable subledger          | MH-V1    | 734      | S3    | S3-E-07   | BL-FIN-022 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | AR subledger per pilgrim/booking                    |
| 6     | F9-CSV-132  | Digital receipts                       | MH-V1    | 735      | S3    | S3-E-07   | BL-FIN-023 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Digital receipt + sequence + audit                  |
| 6     | F9-CSV-133  | Manual payment & down payment [add-on] | SH       | 736      | S3    | S3-E-07   | BL-FIN-024 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Manual/down-payment recording + evidence            |
| 6     | F9-CSV-134  | Vendor master                          | MH-V1    | 737      | S3    | S3-E-07   | BL-FIN-025 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Vendor master + AP linkage                          |
| 6     | F9-CSV-135  | Accounts payable subledger             | MH-V1    | 738      | S3    | S3-E-07   | BL-FIN-026 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | AP subledger per vendor                             |
| 6     | F9-CSV-136  | Payment authorization                  | MH-V1    | 739      | S3    | S3-E-07   | BL-FIN-027 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Multi-level authorization before disbursement       |
| 6     | F9-CSV-137  | Petty cash & temporary vouchers        | SH       | 740      | S3    | S3-E-07   | BL-FIN-028 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Petty cash + temp vouchers + period close           |
| 6     | F9-CSV-138  | Project-based costing                  | MH-V1    | 741      | S3    | S3-E-07   | BL-FIN-029 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Cost object per project/departure + allocation      |
| 6     | F9-CSV-139  | Departure P&L                          | MH-V1    | 742      | S3    | S3-E-07   | BL-FIN-030 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | P&L per departure vs initial budget                 |
| 6     | F9-CSV-140  | Budget vs actual analysis              | SH       | 743      | S3    | S3-E-07   | BL-FIN-031 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Budget vs actual variance + drill-down              |
| 6     | F9-CSV-141  | Automated journals                     | MH-V1    | 744      | S3    | S3-E-07   | BL-FIN-032 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Auto journals idempotent per source event           |
| 6     | F9-CSV-142  | Revenue recognition                    | MH-V1    | 745      | S3    | S3-E-07   | BL-FIN-033 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Revenue recognition per Q043 policy                 |
| 6     | F9-CSV-143  | Multi-currency                         | SH       | 746      | S3    | S3-E-07   | BL-FIN-034 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Multi-currency + rate snapshot consistent with Q001 |
| 6     | F9-CSV-144  | Fixed asset management                 | SH       | 747      | S3    | S3-E-07   | BL-FIN-035 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Asset cards + basic depreciation                    |
| 6     | F9-CSV-145  | Tax integration                        | SH       | 748      | S3    | S3-E-07   | BL-FIN-036 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Integrated tax (VAT/WHT) per Q046/Q047              |
| 6     | F9-CSV-146  | Agent commission payout                | MH-V1    | 749      | S3    | S3-E-07   | BL-FIN-037 | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Agent commission payout flow + audit                |
| 6     | F9-CSV-147  | Real-time financial reports            | MH-V1    | 750      | S3    | S3-E-07   | BL-FIN-038 | E     | todo   | BL-JNT-014..015          | Real-time reports (balance sheet / P&L summary)     |
| 6     | F9-CSV-148  | Cash flow dashboard                    | MH-V1    | 751      | S3    | S3-E-07   | BL-FIN-039 | E     | todo   | BL-JNT-014..015          | Cash flow dashboard + bank/petty balances           |
| 6     | F9-CSV-149  | AR/AP aging alerts                     | SH       | 752      | S3    | S3-E-07   | BL-FIN-040 | E     | todo   | BL-JNT-014..015          | Receivable/payable aging alerts + buckets           |
| 6     | F9-CSV-150  | Audit trail & access control           | MH-V1    | 753      | S3    | S3-E-07   | BL-FIN-041 | E     | todo   | BL-JNT-014..015          | Finance audit trail + RBAC anti-delete              |


### 6.K — Operational & handling CSV modules (#87–#128) (F7/F8)

Priority source: `docs/Modul UmrohOS - MosCoW.csv`. **#87–#108** map to **F7** (`docs/06-features/07-operations-handling.md` + related pilgrim docs); **#109–#125** to **F8** (`docs/06-features/08-warehouse-and-fulfillment.md`); **#126–#128** are operational *executive dashboard* widgets (`BL-DASH-015`–`017`; align naming with similar widgets in **6.I** when overlapping).

Rows **6.C** / **6.D** / **6.E** remain umbrella rows; this table details **Operational & handling** CSV modules.


| Phase | Feature ref | Summary                           | Priority | Exec seq | Slice | Task Code | Backlog ID  | Owner | Status | Blocked by gate           | Acceptance (short)                           |
| ----- | ----------- | --------------------------------- | -------- | -------- | ----- | --------- | ----------- | ----- | ------ | ------------------------- | -------------------------------------------- |
| 6     | F7-CSV-87   | Collective storage                | MH-V1    | 754      | S3    | S3-E-04   | BL-OPS-021  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Per-departure collective vault + ACL         |
| 6     | F7-CSV-88   | Passport OCR & mahram logic       | MH-V1    | 755      | S3    | S3-E-04   | BL-OPS-022  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Passport OCR + mahram flags per rules        |
| 6     | F7-CSV-89   | Progress tracker & expiry alert   | SH       | 756      | S3    | S3-E-04   | BL-OPS-023  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Document progress + expiry alerts            |
| 6     | F7-CSV-90   | Official letter generator         | MH-V1    | 757      | S3    | S3-E-04   | BL-OPS-024  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Official letters from templates              |
| 6     | F7-CSV-91   | Immigration manifest              | MH-V1    | 758      | S3    | S3-E-04   | BL-OPS-025  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Immigration manifest format + versioning     |
| 6     | F7-CSV-92   | Rooming algorithm                 | MH-V1    | 759      | S3    | S3-E-04   | BL-OPS-026  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Rooming algorithm + valid constraints        |
| 6     | F7-CSV-93   | Transport arrangement             | SH       | 760      | S3    | S3-E-04   | BL-OPS-027  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Group seat/transport assignment              |
| 6     | F7-CSV-94   | Manifest handling                 | SH       | 761      | S3    | S3-E-04   | BL-OPS-028  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Manifest handling + delta publish            |
| 6     | F7-CSV-95   | ID card & staff assignment        | MH-V1    | 762      | S3    | S3-E-04   | BL-OPS-029  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Printable ID card + staff assignment         |
| 6     | F7-CSV-96   | Physical passport log             | SH       | 763      | S3    | S3-E-04   | BL-OPS-030  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Physical passport log + handover audit       |
| 6     | F7-CSV-97   | Visa progress tracker             | MH-V1    | 764      | S3    | S3-E-04   | BL-OPS-031  | E     | todo   | BL-VISA-003               | Visa status per pilgrim + SLA                |
| 6     | F7-CSV-98   | E-visa repository                 | SH       | 765      | S3    | S3-E-04   | BL-OPS-032  | E     | todo   | BL-VISA-003               | E-visa repository + metadata                 |
| 6     | F7-CSV-99   | Advanced external API integration | CH       | 766      | S3    | S3-E-04   | BL-OPS-033  | E     | todo   | BL-VISA-003               | Advanced provider connector (optional)       |
| 6     | F7-CSV-100  | Refund & penalty administration   | SH       | 767      | S3    | S3-E-04   | BL-OPS-034  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Admin refund + penalties + approval          |
| 6     | F7-CSV-101  | ALL system                        | MH-V1    | 768      | S3    | S3-E-04   | BL-OPS-035  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | ALL system scan + idempotency                |
| 6     | F7-CSV-102  | Luggage counter                   | SH       | 769      | S3    | S3-E-04   | BL-OPS-036  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Luggage count + scan events                  |
| 6     | F7-CSV-103  | Departure & arrival broadcast     | SH       | 770      | S3    | S3-E-04   | BL-OPS-037  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Departure/arrival schedule broadcast         |
| 6     | F7-CSV-104  | Smart bus boarding                | MH-V1    | 771      | S3    | S3-E-04   | BL-OPS-038  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Bus boarding + roster                        |
| 6     | F7-CSV-105  | Raudhah shield & digital tasreh   | SH       | 772      | S3    | S3-E-04   | BL-OPS-039  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Raudhah shield + digital tasreh              |
| 6     | F7-CSV-106  | Audio device management           | CH       | 773      | S3    | S3-E-04   | BL-OPS-040  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Field audio equipment inventory              |
| 6     | F7-CSV-107  | Zamzam distribution               | SH       | 774      | S3    | S3-E-04   | BL-OPS-041  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Zamzam distribution + handover proof         |
| 6     | F7-CSV-108  | Express room check-in [add-on]    | CH       | 775      | S3    | S3-E-04   | BL-OPS-042  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Express room check-in (add-on)               |
| 6     | F8-CSV-109  | Purchase request                  | MH-V1    | 776      | S3    | S3-E-05   | BL-LOG-013  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | PR + approval + budget link                  |
| 6     | F8-CSV-110  | Budget synchronization            | MH-V1    | 777      | S3    | S3-E-05   | BL-LOG-014  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | PR budget sync vs actual                     |
| 6     | F8-CSV-111  | Tiered approvals                  | SH       | 778      | S3    | S3-E-05   | BL-LOG-015  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Tiered PR/PO approval                        |
| 6     | F8-CSV-112  | Vendor automation                 | SH       | 779      | S3    | S3-E-05   | BL-LOG-016  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Rule-based vendor selection automation       |
| 6     | F8-CSV-113  | Goods Receipt                     | MH-V1    | 780      | S3    | S3-E-05   | BL-LOG-017  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | GRN + partial + reversal policy              |
| 6     | F8-CSV-114  | Quality control                   | SH       | 781      | S3    | S3-E-05   | BL-LOG-018  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Inbound QC + reject status                   |
| 6     | F8-CSV-115  | Auto AP trigger                   | MH-V1    | 782      | S3    | S3-E-05   | BL-LOG-019  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Auto AP posting from GRN                     |
| 6     | F8-CSV-116  | Barcode/SKU labeling              | SH       | 783      | S3    | S3-E-05   | BL-LOG-020  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Barcode/SKU labels consistent with master    |
| 6     | F8-CSV-117  | Multi-warehouse                   | SH       | 784      | S3    | S3-E-05   | BL-LOG-021  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Multi-warehouse + transfers                  |
| 6     | F8-CSV-118  | Critical stock alerts             | SH       | 785      | S3    | S3-E-05   | BL-LOG-022  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Stock below reorder alerts                   |
| 6     | F8-CSV-119  | Kit assembly                      | SH       | 786      | S3    | S3-E-05   | BL-LOG-023  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Kit assembly all-or-nothing                  |
| 6     | F8-CSV-120  | Digital stocktake                 | SH       | 787      | S3    | S3-E-05   | BL-LOG-024  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Digital stocktake + variance                 |
| 6     | F8-CSV-121  | Size synchronization              | MH-V1    | 788      | S3    | S3-E-05   | BL-LOG-025  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Fulfillment size presets synced with catalog |
| 6     | F8-CSV-122  | Shipment trigger                  | MH-V1    | 789      | S3    | S3-E-05   | BL-LOG-026  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Ship trigger after paid-in-full              |
| 6     | F8-CSV-123  | Courier integration               | SH       | 790      | S3    | S3-E-05   | BL-LOG-027  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Courier integration + tracking numbers       |
| 6     | F8-CSV-124  | Self pickup                       | SH       | 791      | S3    | S3-E-05   | BL-LOG-028  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Self pickup + single-use QR                  |
| 6     | F8-CSV-125  | Returns & exchanges [add-on]      | CH       | 792      | S3    | S3-E-05   | BL-LOG-029  | E     | todo   | BL-JNT-009..011; BL-EGV-003 | Returns/exchanges (add-on)                   |
| 6     | F11-CSV-126 | Inventory health                  | SH       | 793      | S5    | S5-L-02   | BL-DASH-015 | L     | todo   | BL-JNT-014..015          | Inventory health widget (read F8)            |
| 6     | F11-CSV-127 | Fulfillment & PO monitor          | SH       | 794      | S5    | S5-L-02   | BL-DASH-016 | L     | todo   | BL-JNT-014..015          | PO + fulfillment backlog monitor             |
| 6     | F11-CSV-128 | Damage report                     | CH       | 795      | S5    | S5-L-02   | BL-DASH-017 | L     | todo   | BL-JNT-014..015          | Goods damage report (add-on)                 |


### 6.L — B2C front-end CSV modules (#1–#24)

Priority source: `docs/Modul UmrohOS - MosCoW.csv`. Scope: public B2C site — align with `docs/06-features/02-catalog-and-master-data.md`, `docs/06-features/04-booking-and-allocation.md`, `docs/06-features/05-payment-and-reconciliation.md`.


| Phase | Feature ref | Summary                                       | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate                    | Acceptance (short)                            |
| ----- | ----------- | --------------------------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ---------------------------------- | --------------------------------------------- |
| 6     | B2C-CSV-01  | Dynamic homepage module                       | MH-V1    | 796      | S1    | S1-L-05   | BL-B2C-001 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Dynamic homepage + mobile performance         |
| 6     | B2C-CSV-02  | Legitimacy validation & About Us module       | MH-V1    | 797      | S1    | S1-L-05   | BL-B2C-002 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Verified legitimacy & About Us visible        |
| 6     | B2C-CSV-03  | Gallery & social proof module                 | SH       | 798      | S1    | S1-L-05   | BL-B2C-003 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Gallery + social proof (testimonials)         |
| 6     | B2C-CSV-04  | Blog & articles module                        | CH       | 799      | S1    | S1-L-05   | BL-B2C-004 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Blog/article list + basic SEO                 |
| 6     | B2C-CSV-05  | Brand identity & white-label module           | SH       | 800      | S1    | S1-L-05   | BL-B2C-005 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | White-label theme + brand assets              |
| 6     | B2C-CSV-06  | Menu builder & navigation module              | MH-V1    | 801      | S1    | S1-L-05   | BL-B2C-006 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Configurable nav menu + ACL                   |
| 6     | B2C-CSV-07  | Smart search & filter module                  | MH-V1    | 802      | S1    | S1-L-05   | BL-B2C-007 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Package search/filter consistent with catalog |
| 6     | B2C-CSV-08  | Product detail & interactive itinerary module | MH-V1    | 803      | S1    | S1-L-05   | BL-B2C-008 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Product detail + interactive itinerary        |
| 6     | B2C-CSV-09  | Real-time availability module                 | MH-V1    | 804      | S1    | S1-L-05   | BL-B2C-009 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Real-time seats/availability from catalog     |
| 6     | B2C-CSV-10  | Savings simulation calculator module          | CH       | 805      | S1    | S1-L-05   | BL-B2C-010 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Savings simulation (add-on) + disclaimer      |
| 6     | B2C-CSV-11  | Essential info & seat tracker                 | MH-V1    | 806      | S1    | S1-L-05   | BL-B2C-011 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Essential info + public seat tracker          |
| 6     | B2C-CSV-12  | Smart accommodation specs                     | SH       | 807      | S1    | S1-L-05   | BL-B2C-012 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Short accommodation specs per package         |
| 6     | B2C-CSV-13  | Guide profile                                 | SH       | 808      | S1    | S1-L-05   | BL-B2C-013 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Guide profile shown per package               |
| 6     | B2C-CSV-14  | Micro-web itinerary                           | SH       | 809      | S1    | S1-L-05   | BL-B2C-014 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Shareable micro-web itinerary                 |
| 6     | B2C-CSV-15  | Call-to-action                                | MH-V1    | 810      | S1    | S1-L-05   | BL-B2C-015 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | CTAs to WA/booking with consistent tracking   |
| 6     | B2C-CSV-16  | Self-booking engine                           | MH-V1    | 811      | S1    | S1-L-05   | BL-B2C-016 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | End-to-end B2C self-booking flow              |
| 6     | B2C-CSV-17  | Guest data form                               | MH-V1    | 812      | S1    | S1-L-05   | BL-B2C-017 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Guest data form + field validation            |
| 6     | B2C-CSV-18  | B2C payment gateway                           | MH-V1    | 813      | S2    | S2-L-04   | BL-B2C-018 | L     | todo   | BL-JNT-005..008; BL-EGV-002; BL-LGV-002 | B2C checkout wired to payment slice           |
| 6     | B2C-CSV-19  | Payment history                               | SH       | 814      | S1    | S1-L-05   | BL-B2C-019 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Pilgrim payment transaction history           |
| 6     | B2C-CSV-20  | Self-upload documents                         | SH       | 815      | S3    | S3-L-03   | BL-B2C-020 | L     | todo   | BL-DOC-001                         | Self-service document upload + status         |
| 6     | B2C-CSV-21  | Logistics & kitting info                      | CH       | 816      | S1    | S1-L-05   | BL-B2C-021 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Logistics/kitting info (add-on)               |
| 6     | B2C-CSV-22  | Departure information board                   | SH       | 817      | S1    | S1-L-05   | BL-B2C-022 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Departure info board (read)                   |
| 6     | B2C-CSV-23  | Knowledge base                                | CH       | 818      | S1    | S1-L-05   | BL-B2C-023 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Knowledge base browse/search                  |
| 6     | B2C-CSV-24  | Floating chat                                 | SH       | 819      | S1    | S1-L-05   | BL-B2C-024 | L     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Floating chat + channel routing               |


### 6.M — Admin & security CSV modules (#151–#163) (F1)

Priority source: `docs/Modul UmrohOS - MosCoW.csv`. Domain: `docs/06-features/01-identity-and-access.md` (RBAC, audit, configuration).


| Phase | Feature ref | Summary                           | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate                    | Acceptance (short)                     |
| ----- | ----------- | --------------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ---------------------------------- | -------------------------------------- |
| 6     | F1-CSV-151  | Job role creation                 | MH-V1    | 820      | S1    | S1-E-06   | BL-IAM-005 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Job role CRUD + IAM mapping            |
| 6     | F1-CSV-152  | Specific permission mapping       | MH-V1    | 821      | S1    | S1-E-06   | BL-IAM-006 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Granular permissions per route/action  |
| 6     | F1-CSV-153  | Data visibility hierarchy         | MH-V1    | 822      | S1    | S1-E-06   | BL-IAM-007 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Data scope hierarchy (global/branch)   |
| 6     | F1-CSV-154  | Staff account registration        | MH-V1    | 823      | S1    | S1-E-06   | BL-IAM-008 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Staff account onboarding + invite      |
| 6     | F1-CSV-155  | User status control               | MH-V1    | 824      | S1    | S1-E-06   | BL-IAM-009 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Suspend/active user + audited reason   |
| 6     | F1-CSV-156  | Account & password security       | MH-V1    | 825      | S1    | S1-E-06   | BL-IAM-010 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Password policy + optional MFA         |
| 6     | F1-CSV-157  | Centralized activity log          | MH-V1    | 826      | S1    | S1-E-06   | BL-IAM-011 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Centralized searchable activity log    |
| 6     | F1-CSV-158  | Anomaly alerts                    | SH       | 827      | S1    | S1-E-06   | BL-IAM-012 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Login/action anomaly alerts (SH)       |
| 6     | F1-CSV-159  | User session history              | SH       | 828      | S1    | S1-E-06   | BL-IAM-013 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Session history + revoke               |
| 6     | F1-CSV-160  | API integration configuration     | MH-V1    | 829      | S1    | S1-E-06   | BL-IAM-014 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | API key configuration + rotation       |
| 6     | F1-CSV-161  | Communication template management | SH       | 830      | S1    | S1-E-06   | BL-IAM-015 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Centralized WA/email templates         |
| 6     | F1-CSV-162  | Global variable configuration     | MH-V1    | 831      | S1    | S1-E-06   | BL-IAM-016 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | Global config key-value + audit        |
| 6     | F1-CSV-163  | Database backup                   | MH-V1    | 832      | S1    | S1-E-06   | BL-IAM-017 | E     | todo   | BL-JNT-001..004; BL-EGV-001; BL-LGV-001 | DB backup schedule/restore (procedure) |


### 6.N — Pilgrim journey CSV modules (#164–#176) (F12)

Priority source: `docs/Modul UmrohOS - MosCoW.csv`. Domain: `docs/06-features/12-alumni-and-daily-app.md` (pilgrim experience in-trip / post-booking).


| Phase | Feature ref | Summary                         | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate           | Acceptance (short)                       |
| ----- | ----------- | ------------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ------------------------- | ---------------------------------------- |
| 6     | F12-CSV-164 | Schedule & live information     | MH-V1    | 833      | S3    | S3-L-04   | BL-JMJ-001 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Live schedule + per-trip push info       |
| 6     | F12-CSV-165 | Smart reminders                 | SH       | 834      | S3    | S3-L-04   | BL-JMJ-002 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Smart itinerary reminders                |
| 6     | F12-CSV-166 | Digital worship guide           | SH       | 835      | S3    | S3-L-04   | BL-JMJ-003 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Offline-friendly worship guide           |
| 6     | F12-CSV-167 | Digital document wallet         | MH-V1    | 836      | S3    | S3-L-04   | BL-JMJ-004 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Per-pilgrim digital document wallet      |
| 6     | F12-CSV-168 | Emergency button                | SH       | 837      | S3    | S3-L-04   | BL-JMJ-005 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Emergency button + escalation            |
| 6     | F12-CSV-169 | E-certificate generator         | CH       | 838      | S3    | S3-L-04   | BL-JMJ-006 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | E-certificate generation (add-on)        |
| 6     | F12-CSV-170 | Bus boarding attendance         | MH-V1    | 839      | S3    | S3-L-04   | BL-JMJ-007 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Bus boarding attendance + scan           |
| 6     | F12-CSV-171 | Bus time control                | SH       | 840      | S3    | S3-L-04   | BL-JMJ-008 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Bus timing control + SLA                 |
| 6     | F12-CSV-172 | Airport luggage tracking        | SH       | 841      | S3    | S3-L-04   | BL-JMJ-009 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Airport luggage status tracking          |
| 6     | F12-CSV-173 | Communication device management | CH       | 842      | S3    | S3-L-04   | BL-JMJ-010 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Communication device management (add-on) |
| 6     | F12-CSV-174 | Zamzam distribution             | SH       | 843      | S3    | S3-L-04   | BL-JMJ-011 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Zamzam distribution + proof              |
| 6     | F12-CSV-175 | Daily issue reporting           | SH       | 844      | S3    | S3-L-04   | BL-JMJ-012 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Daily issue reporting                    |
| 6     | F12-CSV-176 | Mitigation education            | CH       | 845      | S3    | S3-L-04   | BL-JMJ-013 | L     | todo   | BL-JNT-009..011; BL-EGV-003 | Mitigation education (add-on)            |


### 6.O — Complementary features & daily app CSV modules (#190–#198) (F12)

Priority source: `docs/Modul UmrohOS - MosCoW.csv`. Modules **#199–#202** (referral/ZISWAF) are in **6.H**. Rows below: daily content & community — `docs/06-features/12-alumni-and-daily-app.md`.


| Phase | Feature ref | Summary                    | Priority | Exec seq | Slice | Task Code | Backlog ID | Owner | Status | Blocked by gate  | Acceptance (short)                    |
| ----- | ----------- | -------------------------- | -------- | -------- | ----- | --------- | ---------- | ----- | ------ | ---------------- | ------------------------------------- |
| 6     | F12-CSV-190 | Prayer times & adhan       | SH       | 846      | S5    | S5-L-03   | BL-PLG-001 | L     | todo   | BL-JNT-014..015 | Prayer schedule + adhan notifications |
| 6     | F12-CSV-191 | Qibla compass              | SH       | 847      | S5    | S5-L-03   | BL-PLG-002 | L     | todo   | BL-JNT-014..015 | Accurate qibla compass per location   |
| 6     | F12-CSV-192 | Digital Quran              | CH       | 848      | S5    | S5-L-03   | BL-PLG-003 | L     | todo   | BL-JNT-014..015 | Digital Quran (add-on)                |
| 6     | F12-CSV-193 | Dhikr & prayer collection  | CH       | 849      | S5    | S5-L-03   | BL-PLG-004 | L     | todo   | BL-JNT-014..015 | Daily dhikr & supplications           |
| 6     | F12-CSV-194 | Manasik encyclopedia       | SH       | 850      | S5    | S5-L-03   | BL-PLG-005 | L     | todo   | BL-JNT-014..015 | Searchable manasik encyclopedia       |
| 6     | F12-CSV-195 | Articles & regular studies | CH       | 851      | S5    | S5-L-03   | BL-PLG-006 | L     | todo   | BL-JNT-014..015 | Regular articles/studies              |
| 6     | F12-CSV-196 | Religious Q&A              | CH       | 852      | S5    | S5-L-03   | BL-PLG-007 | L     | todo   | BL-JNT-014..015 | Religious Q&A (moderated)             |
| 6     | F12-CSV-197 | Cohort group forum         | CH       | 853      | S5    | S5-L-03   | BL-PLG-008 | L     | todo   | BL-JNT-014..015 | Cohort group forum                    |
| 6     | F12-CSV-198 | Reunion notice board       | CH       | 854      | S5    | S5-L-03   | BL-PLG-009 | L     | todo   | BL-JNT-014..015 | Reunion & announcement board          |


---

## Daily usage rules

1. When starting work, name at least: `Backlog ID` + `Task Code` (when mapped to a slice).
2. If one backlog item is too large for a single PR, split into a new item (`BL-...-next`) before continuing coding.
3. Change priority in this table first, then mirror to the board ticket.
4. `Blocked by gate` lists freeze **backlog ids** (Core principles 4–6 and the **Freeze backlog id index**) that must finish first; until they do, keep the row in `todo`.
5. Phase 6 rows use **`Slice` / `Task Code`** from the depth packages in **[05](./05-slice-engineering-checklist-and-task-codes.md)** (**Phase 6 — Depth backlog**). One task code may cover many `BL-*` rows; split sub-tasks if the PR grows too large.

Example prompt:

> Implement `BL-PAY-004` on `S2-E-02`. Do not change contracts outside `slice-S2.md`.

---

## References

- `[04-delivery-plan-2p-sequence-first.md](./04-delivery-plan-2p-sequence-first.md)` — phase & slice order
- `[05-slice-engineering-checklist-and-task-codes.md](./05-slice-engineering-checklist-and-task-codes.md)` — gates & task codes
- `docs/06-features/01-identity-and-access.md`
- `docs/06-features/02-catalog-and-master-data.md`
- `docs/06-features/03-pilgrim-and-documents.md`
- `docs/06-features/04-booking-and-allocation.md`
- `docs/06-features/05-payment-and-reconciliation.md`
- `docs/06-features/06-visa-pipeline.md`
- `docs/06-features/07-operations-handling.md`
- `docs/06-features/08-warehouse-and-fulfillment.md`
- `docs/06-features/09-finance-and-accounting.md`
- `docs/06-features/10-marketing-crm-agents.md`
- `docs/06-features/11-dashboards.md`
- `docs/06-features/12-alumni-and-daily-app.md`
- `docs/Modul UmrohOS - MosCoW.csv`

