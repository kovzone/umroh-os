# PRD Pointer

The canonical Product Requirements Document for UmrohOS lives at:

```
docs/UmrohOS - Product Requirements Document.docx.md
```

It is a markdown export of the original `.docx`. Approximately 1,620 lines, mostly Indonesian. Production-grade detail (sitemap, menu structure, user flows, Figma-level UI specs for some pages).

## Section index

The PRD is structured roughly as follows:

1. **Product overview & USPs** — what UmrohOS is, why it exists, key differentiators.
2. **Module & feature platform architecture** — the 11 functional sections (A through K):
   - **A. B2C Front-End** — public catalog, self-booking
   - **B. B2B Front-End** — agent/reseller portal, replicated websites, commission wallets
   - **C. Marketing & Sales (CRM)** — campaigns, lead nurturing, ROAS
   - **D. Master Product & Inventory** — packages, hotels, flights, muthawwif
   - **E. Operational & Handling** — documents, visa, manifests, room/bus allocation
   - **F. Inventory & Logistics** — warehouse, procurement, kits, shipping
   - **G. Finance & Accounting** — PSAK-compliant journal, AR/AP, tax, multi-currency
   - **H. Admin & Security** — RBAC, audit, system config
   - **I. Jamaah Journey** — pilgrim mobile app, field operations, Raudhah Shield
   - **J. Daily App & Alumni Hub** — post-pilgrimage engagement
   - **K. Executive Dashboards** — operational/financial/inventory KPIs
3. **Sitemap & URL routing** — public and authenticated routes.
4. **Menu structure** — sidebar/menu hierarchy for each role.
5. **User workflows** — 11 phases of operation, end-to-end.
6. **Detailed workflow examples (Alur Logika)** — step-by-step actor flows.
7. **UI/UX page-level specs** — Figma-level component breakdowns for select pages.

## How to use the PRD

- **Don't try to read it cover-to-cover in one session.** It's too long. Read the section relevant to your current task.
- **Search by Indonesian term** — most module names and features are Indonesian. Use the glossary in `01-glossary.md` to find the right term.
- **Treat the PRD as the source of truth** for product behavior. If a doc in this repo disagrees with the PRD, the PRD wins — but flag the conflict in your session note.
- **The tech stack hints in the PRD (Node.js, MySQL) are overridden** by the locked stack in `docs/01-architecture/01-tech-stack.md`. The PRD is product, not tech.

## How section maps to services

| PRD section | Primary service(s) |
|---|---|
| A. B2C Front-End | Svelte web (app structure TBD — see Q009) consuming `catalog-svc`, `booking-svc`, `payment-svc` |
| B. B2B Front-End | Svelte web (app structure TBD — see Q009) consuming `crm-svc`, `catalog-svc`, `booking-svc` |
| C. Marketing & Sales | `crm-svc` |
| D. Master Product & Inventory | `catalog-svc` |
| E. Operational & Handling | `ops-svc`, `jamaah-svc`, `visa-svc` |
| F. Inventory & Logistics | `logistics-svc` |
| G. Finance & Accounting | `finance-svc` |
| H. Admin & Security | `iam-svc` |
| I. Jamaah Journey | `jamaah-svc`, `broker-svc` (workflows), Svelte field app (app structure TBD — see Q009) |
| J. Daily App & Alumni Hub | `crm-svc` (community), `jamaah-svc` (profile) |
| K. Executive Dashboards | Svelte web (app structure TBD — see Q009) consuming all services + observability stack |
