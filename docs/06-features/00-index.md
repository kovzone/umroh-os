# Feature Specs — Index

This directory is the **middle layer** between the PRD and technical implementation. Each file describes one product feature from the user's point of view with enough detail that either developer (backend or frontend) can work from it without re-reading the 1,620-line PRD.

## Who this is for

Both developers read these specs. Both developers can propose edits. Each dev's private tracker then lists the tasks their side needs to do for each feature.

## What a feature spec contains

1. **Purpose & persona** — what this feature does, for which role (calon jamaah, agent, CS, ops admin, finance admin, etc.)
2. **Sources** — pinpointed PRD sections + module numbers from `docs/Modul UmrohOS - MosCoW.csv` (cite, don't restate)
3. **User workflow** — actor-by-actor, step-by-step happy path
4. **Acceptance criteria** — the bullet list that defines "done"
5. **Edge cases & error paths** — what could go wrong, what the system does about it
6. **Data & state implications** — which entities are touched, what statuses change
7. **API surface (high-level)** — what endpoints/RPCs need to exist; contract details live in `docs/03-services/<svc>/01-api.md`
8. **Dependencies** — which other features/services must be live first
9. **Backend notes** / **Frontend notes** — short sections calling out side-specific concerns
10. **Open questions** — links to `docs/07-open-questions/` entries for anything that needs stakeholder input

## How to handle gaps in the PRD

When the PRD doesn't answer a question during spec writing:

- **If the answer is obvious from common practice or a clear industry standard:** infer it, write it into the spec, and prefix the line with `_(Inferred)_` so the reviewer can spot it.
- **If there's a real product decision to make:** do not guess. Open a question in `docs/07-open-questions/` framed as meeting material. Reference the question ID from the spec. Leave that section marked **TBD — see Q-NN**.

## Feature list

| # | Feature | MoSCoW profile | Status |
|---|---|---|---|
| F1 | [Identity, Access, Audit](01-identity-and-access.md) | 10 Must / 3 Should — all High | **Written** |
| F2 | [Product Catalog & Master Data](02-catalog-and-master-data.md) | 8 Must / 7 Should / 1 Could | **Draft** — 4 open Qs (Q001–Q004) |
| F3 | [Pilgrim Profile & Documents](03-pilgrim-and-documents.md) | Must #87, #88 (mahram) | **Draft** — 4 open Qs (Q005–Q008) |
| F4 | [Booking Creation & Allocation](04-booking-and-allocation.md) | Must — B2C + B2B booking | **Draft** — 6 open Qs (Q004, Q005, Q006, Q010, Q014, Q015, Q016, Q017, Q019) |
| F5 | [Payment & Reconciliation](05-payment-and-reconciliation.md) | 4 Must Haves | **Draft** — 5 open Qs (Q001, Q004, Q011, Q012, Q013) |
| F6 | [Visa Pipeline & Raudhah Shield](06-visa-pipeline.md) | Must #97 | **Draft** — 9 open Qs (Q005, Q007, Q008, Q026–Q031) |
| F7 | [Operations: Verification, Grouping, Manifests](07-operations-handling.md) | 6 Must Haves | **Draft** — 8 open Qs (Q012, Q015, Q020–Q025) |
| F8 | [Warehouse, Procurement, Fulfillment](08-warehouse-and-fulfillment.md) | 6 Must / 12 Should / 2 Could | **Draft** — 10 open Qs (Q032–Q041) |
| F9 | [Finance & Accounting (PSAK)](09-finance-and-accounting.md) | **15 Must Haves — highest** | Stub |
| F10 | [Marketing, CRM, Agent Network](10-marketing-crm-agents.md) | 8 Must Haves | Stub |
| F11 | [Dashboards & Reporting](11-dashboards.md) | 4 Must Haves | Stub |
| F12 | [Alumni Hub & Daily App](12-alumni-and-daily-app.md) | 0 Must Haves | Stub |

Stubs contain: title, PRD section pointers, module numbers, and a "TBD" marker per section. A session fills them in on demand when work starts on that feature — writing all specs upfront is wasted motion; the product details evolve.

## How specs interact with other docs

```
PRD + MoSCoW CSVs
       │
       ▼
docs/06-features/ (this dir)          ← user-POV workflows, acceptance criteria
       │
       ├──────▶ docs/03-services/     ← per-service technical specs (Go side)
       │
       ├──────▶ docs/91-progress/     ← this developer's tasks (private)
       │
       └──────▶ other dev's tracker   ← the other developer's tasks (their own setup)
```

## When to write a new spec vs update an existing one

- **New feature not in the table above:** create a new numbered file, add the row to the index, add the module references.
- **Existing feature, new scope:** update the existing spec and bump its "Last updated" date. Don't shard one feature across multiple files.
- **Cross-cutting concern (e.g. observability, audit)** — these are conventions, not features. They live in `docs/04-backend-conventions/` (Go services) or `docs/05-frontend-conventions/` (Svelte 5 + Vite). Either developer may edit either dir; the split is about which *codebase* the convention applies to, not which *human* owns it.
