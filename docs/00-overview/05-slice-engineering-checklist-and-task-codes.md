# Slice Engineering — Checklist & Task Codes (Lutfi + Elda)

This document complements [`04-delivery-plan-2p-sequence-first.md`](./04-delivery-plan-2p-sequence-first.md) (phase & slice order). Backlog rows `BL-*` + `Exec seq` live in [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md).

**Goal:** finish **engineering dependencies** *before* slice coding starts, so parallel work has fewer collisions and less rework.

---

## Task code format (for AI / peer instructions)

```
S{slice}-{owner}-{seq}
```

| Segment | Meaning | Values |
|---------|---------|--------|
| `S{slice}` | Vertical slice | `S0` … `S5` (see below) |
| `{owner}` | Execution owner | `L` = Lutfi, `E` = Elda, `J` = Joint (both must align / merge together) |
| `{seq}` | Two-digit sequence | `01`, `02`, … |

**Example prompt to an AI:**

> Implement **S1-L-04** and **S1-E-03**. Follow the contract in `docs/contracts/slice-S1.md` (when it exists). Do not change the contract without explicit approval.

**Optional but useful:**

- One contract file per slice: `docs/contracts/slice-Sx.md` (Markdown + JSON examples) — *the `docs/contracts/` folder may be created when S0-J-01 completes.*
- Issue tracker (GitHub/Linear): issue title = task code, 1:1 with this checklist.

---

## Priority vs slice (correct mental model)

Use these labels on every **feature/backlog item**:

- `MH-MVP`: required for fastest MVP
- `MH-V1`: required for v1 release (after MVP)
- `SH`: important; after all `MH-MVP` + `MH-V1` are safe
- `CH`: optional; only if capacity allows

Principles:

- **Priority belongs to the feature/backlog item**, not the slice.
- **Slice is delivery order** (end-to-end container).
- **Task code** is a technical work package that may cover several backlog items.

Consequences:

- One task code may include backlog items with different priority labels; prefer splitting (`-a`, `-b`) for focus.
- `MH-MVP` / `MH-V1` / `SH` / `CH` tagging happens in the backlog mapping doc (feature level), not a default per slice.

---

## Slice definitions

| Code | Slice name | User journey (short) |
|------|------------|----------------------|
| **S0** | Engineering bootstrap | Repo conventions, CI, contract format, merge ownership |
| **S1** | Discover + draft booking | B2C: catalog → detail → booking form → **draft** (no pilgrim login); staff/internal uses auth per F1 |
| **S2** | Get paid | Draft → invoice/VA → webhook → `pending_payment` / partial / paid |
| **S3** | Fulfill + minimum finance | Paid → fulfillment tasks → basic journal / finance |
| **S4** | Growth loop | Lead + attribution + commission/CRM hook |
| **S5** | Hardening | UAT, perf, security checklist, freeze |

---

## “Slice may start coding” rules

A slice `Sx` (for `x ≥ 1`) is **READY TO BUILD** only if:

1. Every row in that slice’s **Engineering freeze** table is **checked** (no open gates).
2. Contract artifacts for the slice exist (Markdown minimum) and are **reviewed** by the non-executor owner (Lutfi reviews Elda and vice versa).
3. No open **[GATE]** items from earlier slices (unless explicitly **waived** with a dated note).

---

# S0 — Engineering bootstrap (before S1)

**Goal:** one language for contracts, merges, and quality gates — **and** a minimal **platform baseline** (observability, CI discipline, health endpoints, migrations, e2e entry point) so S1+ work does not accumulate hidden debt.

## Checklist — Joint **[GATE]**

| Code | Owner | Work | Output / proof of done |
|------|-------|------|-------------------------|
| S0-J-01 | J | Create folder & contract templates `docs/contracts/README.md` + `slice-Sx.md` template | Folder + template merged |
| S0-J-02 | J | Agree **branch strategy** (e.g. `main` + short-lived `feat/*`) + “who merges” rules | One paragraph in contract README or internal wiki |
| S0-J-03 | J | Define **Definition of Ready (DoR)** & **Definition of Done (DoD)** per PR | Short table in contract README |
| S0-J-04 | J | **Scaffold sweep** — repo layout, service modules, and `apps/core-web` match [`docs/01-architecture/adr/0004-monorepo-layout.md`](../../01-architecture/adr/0004-monorepo-layout.md); drift fixed or documented | Short note in contract README *or* ADR-0004 addendum listing exceptions |
| S0-J-05 | J | **OpenTelemetry** — tracing (or trace propagation) + log field conventions documented; each Go service exports OTel to the shared collector path used in dev (`docker-compose` / env doc updated) | Doc snippet + one cross-service request shows correlated trace IDs in logs |
| S0-J-06 | J | **CI path filters** — workflow(s) skip expensive jobs when only `docs/` or unrelated paths change; full matrix still runs for service and `apps/core-web` changes | PR or `CONTRIBUTING.md` pointer to workflow behavior |
| S0-J-07 | J | **Migration pipeline** — `make` targets + CI step for `migrate-up` / version check aligned with [`docs/01-architecture/adr/0007-migration-based-schema.md`](../../01-architecture/adr/0007-migration-based-schema.md) | New contributor can apply migrations from README path |
| S0-J-08 | J | **Standard health** — every Go service exposes **`GET /livez`** (process up) and **`GET /readyz`** (critical dependencies ready); document optional **`GET /diagnostics/db-tx`** (or equivalent) for DB transaction probe where needed. Services: `gateway-svc`, `iam-svc`, `catalog-svc`, `booking-svc`, `payment-svc`, `finance-svc`, `logistics-svc`, `ops-svc`, `visa-svc`, `crm-svc`, `jamaah-svc` | All listed binaries respond 200 on `/livez` and `/readyz` in dev compose; contract README or `docs/03-services/*` links pattern |
| S0-J-09 | J | **E2E skeleton** — minimal Playwright flow per [`docs/01-architecture/adr/0008-e2e-testing-with-playwright.md`](../../01-architecture/adr/0008-e2e-testing-with-playwright.md); document prerequisite (`make dev-bootstrap` or `make dev-up`) and intended CI smoke job | At least one spec merged + doc on how to run locally and in CI |
| S0-L-01 | L | List **roles + UI routes** touched in S1 (public vs internal) | Role vs URL table |
| S0-E-01 | E | List **services** touched S1–S2 + file ownership (PR owner) | Service vs owner table |

---

# S1 — Discover + draft booking

## Engineering freeze (REQUIRED before S1 feature coding) **[GATE]**

| Code | Owner | Work | Output / proof of done |
|------|-------|------|-------------------------|
| S1-J-01 | J | **Public catalog API** contract: list packages, package detail, departure detail + remaining seats (read) | `docs/contracts/slice-S1.md` § Catalog |
| S1-J-02 | J | **API** `POST /v1/bookings` (draft): required fields, error shape, idempotency key (if any) | `docs/contracts/slice-S1.md` § Booking |
| S1-J-03 | J | **Contract** `ReserveSeats` / `ReleaseSeats` (gRPC or internal REST): parameters, failure codes, compensation | `docs/contracts/slice-S1.md` § Inventory |
| S1-J-04 | J | **Booking states** used in S1: at least `draft` (full documents **not** required in S1 if MVP only needs KTP+passport at a later gate — write explicitly) | One decision paragraph + Q006 reference |
| S1-E-01 | E | Review **seat concurrency** + DB transactions (atomic statements) for S1-J-03 | “Approved” comment on contract doc or PR review |
| S1-L-01 | L | Wireframe / screen list for S1 (URL + main components) | Figma link *or* bullets in contract |

_Backlog alignment:_ rows in `docs/00-overview/06-feature-to-backlog-mapping.md` that target Slice **S1** use the gate token `S1-J-01..S1-J-04; S1-E-01; S1-L-01` — the **full** freeze set in the table above (Joint contracts + `S1-E-01` review + `S1-L-01` wireframe), not `S1-J-*` alone. The same gate is tracked there as backlog ids `BL-JNT-001..004; BL-EGV-001; BL-LGV-001` (see **06** § *Freeze backlog id index*).

## Implementation checklist (after freeze)

| Code | Owner | Work | Depends on |
|------|-------|------|------------|
| S1-L-02 | L | Catalog UI + detail + path to booking form | S1-J-01 |
| S1-L-03 | L | Client integration → catalog API | S1-J-01 |
| S1-L-04 | L | Client integration → create draft booking | S1-J-02 |
| S1-E-02 | E | `catalog-svc` read endpoints per contract | S1-J-01 |
| S1-E-03 | E | `booking-svc` draft + reserve-seat orchestration per contract | S1-J-02, S1-J-03 |
| S1-E-04 | E | **Internal** auth middleware for admin/CS routes used in tests (if S1 needs it) | S0-L-01 |
| S1-E-07 | E | Staff-authenticated **catalog write** REST (package + departure MVP; mutating methods on the same path family as public `GET`, Bearer + permission per `slice-S1.md` § Catalog — internal write) | S1-J-01, S1-E-04 |
| S1-L-06 | L | Internal console: package + departure CRUD screens (`/console/packages/...`) wired to **S1-E-07** | S1-L-01, S1-E-04, S1-E-07 |
| S1-L-07 | L | Internal console login page (`/console/login`) integrated to F1 session endpoint | S1-L-01, S1-E-04 |
| S1-L-08 | L | Internal console shell (`/console`) with full sidemenu + route-level 403 state | S1-L-01, S1-E-04, S1-L-07 |

---

# S2 — Get paid (VA + webhook)

## Engineering freeze **[GATE]**

| Code | Owner | Work | Output / proof of done |
|------|-------|------|-------------------------|
| S2-J-01 | J | Contract `POST` invoice + `POST` VA issue: `amount_total`, `currency`, `fx_snapshot`, TTL | `docs/contracts/slice-S2.md` |
| S2-J-02 | J | Webhook contract: signature header, minimal body, dedupe key, response codes | `slice-S2.md` § Webhook |
| S2-J-03 | J | Callback to booking: status transition + idempotency | `slice-S2.md` § Booking integration |
| S2-J-04 | J | **Stub** `payment-svc` (responses still match contract) or `MOCK_GATEWAY` toggle | Stub merged / env flag documented |
| S2-E-01 | E | DB tables for invoice/events per `docs/03-services/04-payment-svc/02-data-model.md` | Migration reviewed |
| S2-L-01 | L | Checkout UI: show VA/QR + polling strategy | Description in contract or UI comment |

_Backlog alignment:_ Slice **S2** rows in `docs/00-overview/06-feature-to-backlog-mapping.md` use `S2-J-01..S2-J-04; S2-E-01; S2-L-01` as the gate token — the full freeze set above. Backlog ids: `BL-JNT-005..008; BL-EGV-002; BL-LGV-002` (**06** freeze index).

## Implementation checklist

| Code | Owner | Work | Depends on |
|------|-------|------|------------|
| S2-E-02 | E | Implement `payment-svc` invoice + VA + webhook | S2-J-01–J-04 |
| S2-E-03 | E | Minimal reconcile cron | S2-J-02 |
| S2-L-02 | L | Checkout page + error UX | S2-J-01 |
| S2-L-03 | L | Wire booking flow → payment calls | S2-J-03 |
| S2-L-04 | L | Deep B2C checkout (VA/QR + advanced error UX; `BL-B2C-018`) | S2-L-02 |
| S2-J-05 | J | End-to-end test: stub then real gateway | S2-E-02, S2-L-03 |

---

# S3 — Fulfillment + minimum finance

## Engineering freeze **[GATE]**

| Code | Owner | Work | Output |
|------|-------|------|--------|
| S3-J-01 | J | Event `payment.received` / `booking.paid_in_full` → payload for logistics + finance | `slice-S3.md` |
| S3-J-02 | J | Minimal fulfillment task contract (status, assignee) | `slice-S3.md` |
| S3-J-03 | J | Minimal journal contract (placeholder accounts + amount rules) | `slice-S3.md` |
| S3-E-01 | E | Review posting load vs refund | Comment on contract |

_Backlog alignment:_ Slice **S3** rows in `docs/00-overview/06-feature-to-backlog-mapping.md` use `S3-J-01..S3-J-03; S3-E-01` as the gate token — the full freeze set above. Backlog ids: `BL-JNT-009..011; BL-EGV-003` (**06** freeze index).

## Implementation checklist

| Code | Owner | Work | Depends on |
|------|-------|------|------------|
| S3-E-02 | E | `logistics-svc` trigger + status | S3-J-02 |
| S3-E-03 | E | `finance-svc` basic posting | S3-J-03 |
| S3-L-02 | L | Portal “kitting” status UI (read-only OK) | S3-J-02 |

---

# S4 — Growth loop (CRM)

## Engineering freeze **[GATE]**

| Code | Owner | Work | Output |
|------|-------|------|--------|
| S4-J-01 | J | Lead schema + UTM snapshot + attribution (Q019/Q057) | `slice-S4.md` |
| S4-J-02 | J | Events from booking → CRM (event names + payload) | `slice-S4.md` |

_Backlog alignment:_ Slice **S4** rows use `S4-J-01..S4-J-02` — the **Engineering freeze** table above lists Joint cards only (no separate `S4-L-01` wireframe or `S4-E-01` review row here). Backlog ids: `BL-JNT-012..013` (**06** freeze index).

## Implementation checklist

| Code | Owner | Work | Depends on |
|------|-------|------|------------|
| S4-L-01 | L | Lead list + capture form | S4-J-01 |
| S4-E-02 | E | Lead storage endpoints + event consumption | S4-J-02 |

---

# S5 — Hardening

## Engineering freeze **[GATE]**

| Code | Owner | Work | Output |
|------|-------|------|--------|
| S5-J-01 | J | Mandatory UAT scenarios (from MVP gates) | Checklist in `slice-S5.md` |
| S5-J-02 | J | Bug severity matrix + fix SLA | Table |

_Backlog alignment:_ Slice **S5** rows use `S5-J-01..S5-J-02` — the **Engineering freeze** table above lists Joint cards only (no separate `S5-L-01` wireframe or `S5-E-01` review row here). Backlog ids: `BL-JNT-014..015` (**06** freeze index).

## Implementation checklist

| Code | Owner | Work | Depends on |
|------|-------|------|------------|
| S5-L-01 | L | UAT journeys B2C/agent | S5-J-01 |
| S5-E-01 | E | UAT payment/finance/logistics | S5-J-01 |

---

# Phase 6 — Depth backlog (after core S1–S5 is stable)

The following packages map **`BL-*` Phase 6** rows in `docs/00-overview/06-feature-to-backlog-mapping.md` to **Slice + Task Code**. One code may cover many backlog rows; split sub-tasks (`S4-E-03a`, etc.) if PR size grows.

## Implementation checklist — depth by domain

| Code | Owner | Work | Depends on | Backlog domain (short; see Phase 6 in mapping) |
|------|-------|------|------------|------------------------------------------------|
| **S1-E-05** | E | **Catalog & master data depth** — hotel/guide/transport masters, variants/addons, import & bulk edit, **cross-channel seats** | S1-E-02, S1-J-01 | `BL-CAT-005`–`011`, `BL-BOOK-007` (`BL-CAT-012`–`013` → **S4-E-04**) |
| **S1-E-06** | E | **IAM & admin platform** — granular RBAC, staff, session/MFA security, centralized logs, API keys, comm templates, global config, backup procedure | S1-E-04, S1-J-01 | `BL-IAM-005`–`017` |
| **S1-L-05** | L | **B2C site depth** — homepage through self-booking, guest form, history, logistics info, KB, chat, etc. (**not** VA checkout & document upload) | S1-L-02, S1-J-01 | `BL-B2C-001`–`017`, `019`, `021`–`024` |
| **S2-L-04** | L | **Deep B2C checkout** — gateway UX + wiring to invoice/VA | S2-L-02, S2-J-01 | `BL-B2C-018` |
| **S3-E-04** | E | **Field ops depth** — collective docs, manifest, rooming/transport, visa UI data, ALL, bus, Raudhah, Zamzam, admin refund, vendor checklist, … | S3-E-02, S3-J-02 | `BL-OPS-010`, `011`, `020`, `BL-OPS-021`–`042` |
| **S3-E-05** | E | **Warehouse & procurement depth** — PR/PO/GRN/QC, multi-warehouse stock, assembly, ship/return | S3-E-02, S3-J-02 | `BL-LOG-010`–`029` |
| **S3-E-06** | E | **Visa pipeline depth** — readiness, bulk submit, poll provider + history | S3-J-02, `BL-VISA-001` | `BL-VISA-001`–`003` |
| **S3-E-07** | E | **Finance module depth** — billing, bank, subledger, AP ladder, tax, advanced rev-rec, commission, reports & audit | S3-E-03, S3-J-03 | `BL-FIN-010`, `011`, `BL-FIN-020`–`041` |
| **S3-L-03** | L | **B2C self-upload documents** + read status | S3-J-02, `BL-DOC-001` | `BL-B2C-020` |
| **S3-L-04** | L | **Pilgrim journey** — in-trip / post-booking (schedule, document wallet, bus, Zamzam, …) | S3-L-02, S3-J-02 | `BL-JMJ-001`–`013` |
| **S3-L-05** | L | **Operational dashboard widgets** — live bus, Raudhah, luggage, incidents, warehouse health, logistics monitoring | S3-L-02, S3-J-02 | `BL-DASH-008`–`013` |
| **S4-E-03** | E | **CRM & growth API/back office depth** — campaigns, automation, lead routing, wallet/payout, ads/UTM backend, discount approval, alumni/ZISWAF data | S4-E-02, S4-J-01 | `BL-CRM-012`, `BL-CRM-040`–`066` (+ rows Owner **E** in **6.H**) |
| **S4-E-04** | E | **Catalog sync → agent channels** — version snapshot, diff, idempotent push | S4-J-02, S1-E-05 | `BL-CAT-012`, `013` |
| **S4-L-02** | L | **Agent portal & marketing assets** — onboarding, replica, content, flyer/itinerary, academy UI, super-view, … | S4-L-01, S4-J-01 | `BL-CRM-010`, `011`, `BL-CRM-013`–`039` (+ rows Owner **L** in **6.H**) |
| **S4-L-03** | L | **Ads & CS dashboard widgets** — spend vs closings, CS performance | S4-L-01, S4-J-01 | `BL-DASH-006`, `007` |
| **S5-L-02** | L | **Executive dashboard widgets** — vendor readiness, seats, cash/P&L snapshot, liquidity, inventory/PO exec, damage | S5-L-01, S5-J-01 | `BL-DASH-001`–`005`, `014`–`017` |
| **S5-L-03** | L | **Daily app & companions** — prayer, qibla, manasik/community content | S5-J-01 | `BL-PLG-001`–`009` |

---

## Backlog mapping template (feature level)

Use this template to map feature detail to task codes:

| Feature ID | Feature summary | Priority | Slice | Task Code | Backlog ID | Owner | Acceptance (short) |
|------------|-----------------|----------|-------|-----------|------------|-------|--------------------|
| F4-BOOK-001 | Create draft booking | MH-MVP | S1 | S1-E-03 | BL-BOOK-001 | E | Draft saved + atomic seat reserve |
| F5-PAY-001 | Issue VA | MH-MVP | S2 | S2-E-02 | BL-PAY-001 | E | VA issued + TTL stored |
| F10-CRM-001 | Lead capture basic | SH | S4 | S4-L-01 | BL-CRM-001 | L | Lead stored with UTM |

Fill rules:

1. `Priority` must come from feature/backlog level.
2. `Slice` and `Task Code` show when and by whom work runs.
3. One `Feature ID` may split into several `Backlog ID` rows if too large for one PR.

---

## Short tips (keep the task-code system alive)

1. **Do not start a slice without merged contracts** — that is one small PR reviewable in under an hour.
2. When asking AI to code, attach: **task code + slice contract contents + service owner** — so boundaries are not changed casually.
3. When a contract changes: **bump version** (`slice-S2-v2.md` or a *Changelog* section in the same file) so history stays clear.

---

## References

- Phase & slice order (sequence-first): [`04-delivery-plan-2p-sequence-first.md`](./04-delivery-plan-2p-sequence-first.md)
- `BL-*` mapping + `Exec seq` (including Phase 6 → depth task codes): [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md)
- Payment data model: `docs/03-services/04-payment-svc/02-data-model.md`
- Booking flow (product): `docs/06-features/04-booking-and-allocation.md`
