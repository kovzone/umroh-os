# Delivery Plan — Two Developers (Sequence-First)

## Context

This document breaks down work for **two people** using a **sequence-first** principle:

- **No calendar commitment** as a delivery promise (duration is flexible).
- The **phase order** below is the engineering commitment: what comes first, what comes later, and which gates must pass before moving on.

This filename does **not** imply a fixed week count (not “8w”). Old links to `04-delivery-plan-2p-8w.md` should point to this filename instead.

## Related overview docs (04 / 05 / 06)

| File | Role |
|------|------|
| **This file (04)** | **Phase & slice order** — two-person context, MH-* / SH / CH priority, ownership, **no** calendar commitment. |
| [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md) | **Gates + task codes** — joint/L/E checklist before coding, `Sx-J-01` format, **Phase 6 — Depth backlog** (`S1-E-05`, …). |
| [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md) | **Operational backlog** — `BL-*` rows, `Exec seq`, Slice/Task Code columns pointing to **05**; product behavior remains in `docs/06-features/*`. |

In short: **04** = map of *what first*; **05** = SOP & task numbers per slice; **06** = integration ticket list.

## Priority Labels (Required)

All backlog and feature-detail items in this repo use:

- **MH-MVP**: must finish for the earliest functional phase (usually core business flow)
- **MH-V1**: must finish for v1 completeness (after the core is stable)
- **SH**: important; scheduled after all `MH-MVP` + `MH-V1` are safe
- **CH**: optional; only if capacity allows

Execution order: **MH-MVP → MH-V1 → SH → CH**.

## How “Shared” Works

“Shared” does **not** mean two people code the same task at once.

Every task must have:

- **A (Accountable):** final decider
- **R (Responsible):** primary executor
- **C (Consulted):** mandatory reviewer

Rules:

1. Only one person is **R** per task
2. Reviewers do not change scope without a 15-minute daily sync
3. If an API contract changes, update the service API doc before continuing implementation

## Primary Ownership

### Lutfi (Marketing owner)

- F10 Marketing/CRM/Agent (primary)
- F4 Booking flow side (channel, attribution, UX flow)
- B2C/B2B funnel UI
- F11 sales/marketing dashboards (phase 2)

### Elda (Finance/Inventory owner)

- F5 Payment core (invoice/VA/webhook/reconcile)
- F9 Finance core (journal, AR/AP, basic reports)
- F8 Logistics/Fulfillment core
- F11 finance/inventory dashboards (phase 2)

### Shared Foundation (Lutfi + Elda)

- F1 auth/role/audit minimum
- F2 package/departure/seat minimum
- Cross-service event contracts (booking/payment/finance/logistics/crm)

## Vertical Slices (End-to-End Order)

Prefer slices that deliver business value quickly and reduce integration risk.

### Slice 1 — Discover & Book Draft

Flow: login → browse packages → pick departure → create draft booking.

### Slice 2 — Get Paid

Flow: draft booking → issue VA → webhook received → paid/partial status.

### Slice 3 — Fulfillment + Accounting Minimum

Flow: paid → trigger fulfillment → basic journal posting → visible status.

### Slice 4 — Growth Loop

Flow: attribution → lead tracking → basic commission → initial sales dashboard.

### Slice 5 — Hardening & Go-Live Readiness

Flow: reliability, audit, permissions, UAT, priority bugfixes.

**Per-slice engineering dependency checklist + task codes (Lutfi/Elda/Joint):** [`05-slice-engineering-checklist-and-task-codes.md`](./05-slice-engineering-checklist-and-task-codes.md).  
**Feature → backlog → task code mapping:** [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md).

## Execution Order (Sequence-First)

### Phase 0 — Engineering bootstrap (S0)

**Goal:** contracts, merge conventions, DoR/DoD definitions, `docs/contracts/*` templates, **plus** a shared technical baseline so later slices do not retrofit observability, CI, health checks, or migrations ad hoc.

That baseline includes: **repo scaffold sweep** (layout matches ADR-0004), **OpenTelemetry** export wired for services (minimal trace + log correlation conventions), **CI path filters** so unrelated changes do not run the full matrix, **migration pipeline** agreed and documented (see ADR-0007), **`/livez` and `/readyz` on every Go service** (gateway-svc, iam-svc, catalog-svc, booking-svc, payment-svc, finance-svc, logistics-svc, ops-svc, visa-svc, crm-svc, jamaah-svc), optional **`/diagnostics/db-tx`** (or equivalent) pattern documented where DB readiness must be distinguished from process liveness, and an **e2e skeleton** (Playwright per ADR-0008: minimal spec + how CI will run smoke).

**Joint gate:** complete `S0-*` checklist in `05-slice-engineering-checklist-and-task-codes.md` (authoritative row list).

### Phase 1 — Discover + draft booking (S1)

**Goal:** B2C can browse catalog → detail → form → **draft booking**; internal auth minimum for routes staff need.

**F1 gate (sequencing):** B2C catalog and draft booking follow **`S1-J-*`** as soon as those contracts are frozen. Any **internal / staff / authenticated** surface, and any automated or manual test that depends on a real session, is **not** treated as done for S1 until **F1 minimum** is in place: internal login + refresh, permission middleware on protected routes, basic session suspend/revoke, and audit on state-changing calls (**`BL-IAM-001`–`BL-IAM-004`**, implemented under **`S1-E-04`**). Waivers must be explicit and dated in the slice contract or checklist.

**Joint gate:** complete engineering freeze `S1-J-*` before large implementation.

### Phase 2 — Get paid (S2)

**Goal:** draft → invoice/VA → webhook → payment status synced to booking.

**Joint gate:** complete engineering freeze `S2-J-*` + end-to-end test `S2-J-05`.

### Phase 3 — Fulfillment + minimum post-pay finance (S3)

**Goal:** `paid_in_full` triggers minimal fulfillment/logistics work + minimum auditable finance posting.

**Joint gate:** complete engineering freeze `S3-J-*`.

### Phase 4 — Basic growth loop (S4)

**Goal:** lead + basic attribution + minimal events to CRM read model.

**Joint gate:** complete engineering freeze `S4-J-*`.

### Phase 5 — Hardening & readiness (S5)

**Goal:** UAT, reliability, permission regression, minimum operational documentation, and **`BL-GTW-100` (gateway↔backend trust contract)** to close the defense-in-depth gap ADR 0009 D2 deferred.

**Joint gate:** complete engineering freeze `S5-J-*`.

### Phase 6 — Depth expansion (after core is stable)

This section is **deliberately not time-boxed**. The list below is **product domain priority** (not a copy of row order in the mapping). Concrete order between `BL-*` rows and subsections **6.G–6.O** follows the **`Exec seq`** column in [`06-feature-to-backlog-mapping.md`](./06-feature-to-backlog-mapping.md).

Suggested domain order:

1. **F10** depth (agent onboarding, commission view, basic reporting, lead SLA)
2. **F9** depth (operational AP, AR/AP aging, minimum period-close checklist)
3. **F8** depth (warehouse/QC/reorder critical flow, minimum dispatch tracking)
4. **F7** depth (manifest/grouping most critical for operations)
5. **F6** starter (visa tracker/basic read model) after documents + booking data are stable
6. **F11** dashboards as operations need them (after metric backends exist)

#### Joint final gate (“v1 operational enough”)

- Core flow Phases 1–3 stays stable after depth features land
- Minimum operational documentation exists
- Remaining bug triage is only prioritized P2/P3

## RACI — Shared Foundation (Short)

### F1 minimum (login + basic roles)

- **A:** Lutfi
- **R:** Elda
- **C:** Lutfi

### F2 minimum (package + departure + basic seats)

- **A:** Lutfi
- **R:** Lutfi
- **C:** Elda

### F4 minimum (create draft booking)

- **A:** Lutfi
- **R:** Lutfi
- **C:** Elda

## Definition of Done (Per Task)

A task is done when all are true:

1. Endpoint/UI meets acceptance
2. Events/status are consistent across services
3. Audit log records state-changing actions
4. Role permissions do not leak
5. At least one happy-path test + one edge case pass

## Sync Rituals (Recommended)

### Short daily sync

- What finished yesterday
- What is in progress today
- Blockers
- Which API/event contracts changed

### Less frequent sync

- Demo of phase/slice in progress
- Bug triage (P0/P1/P2)
- Freeze scope for the next phase (per-phase guardrails, not per calendar)

## Scope Guardrails

To protect **integration quality** without date promises:

- Do not take `CH` items before all `MH-MVP` on the core path (Phases 1–3) is stable
- Full visa pipeline (deep F6), advanced field ops, and full alumni/daily features stay **after** core + basic depth are stable (see Phase 6)
- New scope requires a trade-off: drop something from the current phase or lower priority (`SH`/`CH`)

## External stakeholder sign-offs vs engineering assumptions

Some areas depend on **people outside engineering** (examples tracked as open questions: **Q005** mahram wording cited in contracts/T&Cs, **Q008** PDP privacy policy + breach runbook, **Q012** refund matrix in printed T&Cs, **Q026** MOFA/Sajil credentials and manual Phase 1 operations, **Q030** Nusuk/Raudhah manual-window communications).

**During development and internal UAT:** the team may use **non-binding engineering assumptions** — mocked or draft configs, placeholder policy URLs, matrices and copy labelled **`DRAFT`**, and written manual runbooks. Business rules “outside the system” may be treated as **provisional** for build and test order, as long as nothing customer-facing presents them as final.

**Customer-facing go-live** (production legal/religious copy, printed S&Cs, public marketing claims): requires **explicit human approvals** (owner + date), typically tracked outside this repository. Shipping code does not replace those approvals; treat them as a **separate go-live gate** from engineering checklist completion.
