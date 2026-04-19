---
id: F11
title: Dashboards & Reporting
status: written
last_updated: 2026-04-18
moscow_profile: 4 Must / 9 Should / 0 Could (13 modules in scope)
prd_sections:
  - "K (Section J in PRD numbering — Executive Dashboards, lines 575–621)"
  - "Alur Logika 11.1–11.2 (lines 1245–1249) — headings only; no expanded narrative"
  - "Sitemap /erp/dashboard/* (lines 751–759)"
  - "Menu structure (lines 1031–1039)"
  - "Data-scope rule (lines 505, 1357–1363, 1399)"
  - "Helicopter View vision (lines 9, 25, 27)"
modules:
  - "#177 Eksekusi Vendor, #178 Ketersediaan Kursi (Operational Readiness — both Must Have)"
  - "#179 Pantauan Anggaran Iklan, #180 Papan Kinerja CS (Sales & Marketing)"
  - "#181 Radar Bus, #182 Status Raudhah, #183 Pelacakan Koper, #184 Laporan Insiden (Real-Time Field View)"
  - "#185 Kesehatan Gudang, #186 Pantauan Eksekusi Logistik (Inventory Executive)"
  - "#187 Arus Kas Instan, #188 Laporan Keuangan Eksekutif (Financial Health — both Must Have)"
  - "#189 Likuiditas (Financial Health)"
depends_on: [F1, F2, F4, F5, F6, F7, F8, F9, F10]
open_questions: []
---

# F11 — Dashboards & Reporting

## Purpose & personas

F11 is the **presentation + composition layer** over every other feature's data. It does **not** own business logic; it owns the executive + operational + field viewing experience. The PRD vision (lines 9, 25, 27) is a *Helicopter View* — the owner / director sees the state of the whole agency from one phone screen without waiting for end-of-month reports.

The feature is deliberately thin in the PRD: Phase 11 Alur Logika (lines 1245–1249) has only two bullet headings and **neither is expanded** into step-by-step narrative anywhere else in the document. This signals that F11 is intentionally under-specified upfront; the real content emerges from the upstream features.

Primary personas:

- **Agency owner / Pimpinan / Direksi** — GLOBAL data scope; consumes `/erp/dashboard` home + all sub-dashboards; mobile-first viewing.
- **Direktur / CEO** — GLOBAL data scope; same as owner.
- **CFO / Finance director** — Financial Health dashboards (Arus Kas Instan, Laporan Eksekutif, Likuiditas); drives period-close and approval decisions.
- **CMO / Marketing director** — Sales & Marketing Board (CPL, ROAS, funnel, CS performance).
- **Ops lead** — Operational Readiness (vendor execution, seat inventory) + Real-Time Field View (bus radar, Raudhah, luggage, incidents).
- **Kepala Cabang (branch manager)** — BRANCH data scope; same dashboards filtered to own branch per PRD line 1357.
- **CS supervisor** — Papan Kinerja CS (team performance, conversion, response times); consumed about-the-team, not by the CS themselves.
- **Super-agent / Perwakilan** — their own Dashboard Super-View lives in F10 (line 167); F11 may cross-link but doesn't own it.

Consumer relationships (one-way; F11 consumes):

- `finance-svc` (F9 W15 + W16) → Arus Kas Instan, Laporan Eksekutif, Likuiditas.
- `logistics-svc` (F8 W15) → Kesehatan Gudang, Pantauan Eksekusi Logistik.
- `catalog-svc` (F2) → Ketersediaan Kursi.
- `booking-svc` (F4) + `payment-svc` (F5) → Sales funnel metrics.
- `visa-svc` (F6) + `ops-svc` (F7) → Eksekusi Vendor, Radar Bus, Status Raudhah, Pelacakan Koper, Laporan Insiden.
- `crm-svc` (F10) → Pantauan Anggaran Iklan, Papan Kinerja CS.

## Sources

- PRD Section K (numbered J in the source document — "Executive Dashboards"), lines 575–621 — the authoritative F11 block.
- Alur Logika 11.1 (line 1245) = Vendor readiness + Seat inventory + CS performance analysis.
- Alur Logika 11.2 (line 1249) = Real-time financial reports + cash flow + aging alerts.
- Sitemap `/erp/dashboard/*` routes at lines 751–759.
- Menu structure at lines 1031–1039.
- Data-scope RBAC at lines 505, 1357–1363, 1399.
- 13 modules enumerated in frontmatter (#177–#189).

## User workflows

### W1 — Executive landing (`/erp/dashboard`)

1. Owner or Direksi opens the dashboard root on phone or desktop.
2. Top-of-screen summary: today's cash balance, today's AR/AP delta, paid-but-unshipped count, verification-queue depth, open-incident count.
3. Widget grid: 8–12 KPIs per **Q075** (answered): **owner-configurable ordered list** with Recommendation-table **v1** defaults, role sub-landings `/finance /marketing /operations /saudi`, MTD range, branch+date filters — ordering in **config JSON** without redeploy. Candidate widgets: revenue MTD, bookings MTD, conversion rate, ROAS, aging AR buckets, Raudhah success rate, bus-radar status map.
4. Each widget drillable to detail (**Q069**): widget → breakdown chart → transaction list.
5. Data-scope enforced: `scope=global` shows all; `scope=branch` filters every widget to `branch_id`.
6. Mobile-responsive per PRD line 27 ("hanya melalui satu layar smartphone").

### W2 — Operational Readiness dashboard (`/erp/dashboard/operations`, modules #177, #178)

1. **Eksekusi Vendor (#177)** — per upcoming departure checklist view: ticket issued? hotel confirmed? visa filed? muthawwif assigned? manifest generated?
2. Checklist powered by events from F6 (visa.status), F8 (po.status), F7 (grouping.run / manifest.generated), catalog-svc (package.muthawwif_assigned).
3. Risk view: departures with > N days to go and checklist < 100% complete → surfaced at top; color-coded by days-to-departure.
4. **Ketersediaan Kursi (#178)** — per-departure seat inventory vs capacity; real-time via catalog-svc.
5. Drill-down: seat status view → per-jamaah roster → per-booking detail.

### W3 — Sales & Marketing Board (`/erp/dashboard/sales`, modules #179, #180)

1. **Pantauan Anggaran Iklan (#179)** — ads spend (from F10 Ads Manager Lite pull) vs closings; cost-per-lead and cost-per-acquisition by campaign. Alert if CPL > threshold and closings < target (**Q068** pins threshold).
2. **Papan Kinerja CS (#180)** — per-CS metrics: leads assigned, leads worked, avg response time (against SLA), conversion rate, revenue closed. Leaderboard.
3. Funnel visualization: lead → engaged → hot → closed → paid → departed — step-by-step conversion %.
4. UTM source breakdown: top-converting UTMs; per-agent contribution.
5. Drill-down: each metric → contributing leads / bookings / transactions.

### W4 — Real-Time Field View / Saudi Dashboard (`/erp/dashboard/saudi`, modules #181–#184)

1. **Radar Bus (#181)** — live map of buses across Saudi; powered by F7 smart bus GPS + boarding scans (per **Q074** transport).
2. **Status Raudhah (#182)** — per-departure % of jamaah who successfully entered Raudhah vs target (from F6 tasreh scans + F7 W11).
3. **Pelacakan Koper (#183)** — aggregate luggage position per departure (Jakarta / in-transit / Saudi hotel); from F7 ALL System scan events.
4. **Laporan Insiden (#184)** — live incident feed; push-notified to HQ (only module with explicit streaming requirement per PRD line 603); filter by severity.
5. Real-time transport (**Q074**): websocket streaming for bus radar + incidents; polling 1–5 min for Raudhah + luggage (per Q067 cadence tiers).

### W5 — Inventory Executive View (modules #185, #186)

1. **Kesehatan Gudang (#185)** — total warehouse asset value (IDR), per-warehouse breakdown, critical-stock chart (items below reorder). Served by F8 logistics-svc read endpoints.
2. **Pantauan Eksekusi Logistik (#186)** — paid-but-unshipped queue (aging), outstanding PO count + total amount, GRN backlog.
3. Asset value computed per F8/F9 inventory valuation method (**Q038 binding**: weighted-average default).
4. Note: these same modules are also described in F8 W15 (duplication across PRD lines 419–425 and 609–611). F11 **references** the F8 computation; does not re-specify.

### W6 — Financial Health Dashboard (modules #187, #188, #189)

1. **Arus Kas Instan (#187)** — real-time cash position across bank accounts + petty cash + brankas; from F9 W16. Trend charts; low-balance alerts.
2. **Laporan Keuangan Eksekutif (#188)** — on-demand Neraca, Laba Rugi, Perubahan Ekuitas. From F9 W15 reports.
3. **Likuiditas (#189)** — AR aging (jamaah) + AP aging (vendors), alert buckets per F9 W17.
4. F11 **presents**; F9 owns the computation. F11 composes into a CFO-focused single-screen layout.
5. Exports to PDF/Excel for board meetings (**Q072**).

### W7 — Branch-specific dashboards

1. Branch manager's login scope = `BRANCH`; the same dashboard routes (`/erp/dashboard/*`) render filtered to `branch_id`.
2. Per PRD line 1357: every metric SQL / gRPC query applies `WHERE branch_id = $1`.
3. Central user with `GLOBAL` scope per **Q071** (answered): **total + branch breakdown + drill**; **branch compare** and **view-as** are **Phase 2**; **flat per-branch** list is MVP.
4. Kasir Cabang scope sees a further-stripped dashboard per PRD line 1399 (only branch transactions + invoices).

### W8 — Alerts & notifications

1. Threshold-driven alerts per **Q068**: default inferred — fixed defaults at launch, per-role configurable later.
2. Alert types:
   - Cash balance < X (per bank account)
   - AR overdue > 60 days > Y IDR
   - CPL > Z without matching closings
   - Critical stock items > N
   - Open incidents > M per hour
   - Paid-unshipped queue aging > 7 days
3. Alert delivery: in-app dashboard badge + WA push to configured recipient (owner / CFO / ops lead per role).
4. Snooze / acknowledge per alert; audit log.

### W9 — Drill-down navigation

1. Every KPI widget drills down per **Q069** — default inferred: **widget → breakdown chart → filtered transaction list → source entity detail**. 3 levels minimum.
2. Drill-through respects data-scope (branch manager drills only into their branch's transactions).
3. URL-encoded filters preserve back-navigation and shareable links.

### W10 — Export to PDF / Excel / image

1. Per **Q072** default inferred: Financial Health (W6) exports to PDF + Excel for board meetings; Operational / Sales / Inventory dashboards export to PDF snapshot (image) only; Real-Time Field View does not export (inherent liveness).
2. Export respects data-scope; export format matches the rendering state (current filters preserved).
3. Audit log: which dashboard exported, by whom, timestamp.

## Acceptance criteria

- **Data-scope enforcement** — every metric query filters by `branch_id` for BRANCH-scoped roles; GLOBAL scopes see aggregated + breakdown per Q071; zero leakage of cross-branch data to branch users.
- **Mobile-responsive** — executive dashboards fully functional on phone (per PRD line 27).
- **Refresh cadence** — per Q067 tiers: cash flow + incidents streaming (≤ 10s latency); seat inventory real-time (≤ 30s); sales funnel ~5min; inventory value ~hourly; Neraca / Laba Rugi on-demand.
- **Aggregation performance** — dashboard landing renders < 3s p95 (with caching); drill to detail < 2s p95.
- **Alerts** — configured thresholds fire within 60s of breach; deliveries to WA successful (or logged as failed).
- **Real-time field transport** — bus radar + incidents via websocket per Q074; polling fallback when socket drops.
- **Export** — successful PDF / Excel export within 30s for financial reports; zero missing data in export vs rendered state.
- **Audit** — every export + every threshold edit audit-logged.
- **No duplication of business logic** — F11 presents, F9/F8/etc compute; any formula is in the owning feature, not F11.

## Edge cases & error paths

- **Service down** (e.g. finance-svc unavailable): affected widgets show "temporarily unavailable" with last-known-value timestamp; no crash of the dashboard page.
- **Stale cache** (underlying data moved forward, cache not yet refreshed): show data age badge ("as of 5 min ago"); refresh-now button available.
- **Cross-branch aggregation permission conflict** (user's scope doesn't actually include some branches in the aggregate): Q071 default: central scope sees all; branch scope sees their own only; no partial scopes in MVP.
- **Drill-down beyond scope** (branch manager tries to drill into central-only report): backend rejects with 403; UI hides the drill affordance where possible.
- **Alert storm** (mass threshold breach during black-swan event): throttle notifications (per-alert-type rate limit: max 1 WA push per 15 min); daily-digest email fallback.
- **Websocket disconnect** (Q074): auto-reconnect; show disconnect badge; polling fallback maintains data freshness during outage.
- **Export with real-time data**: export captures snapshot timestamp; if data changes during export, the rendered state at export-start is preserved.
- **Deleted entity in drill path** (jamaah soft-deleted after KPI aggregation): drill-down shows "[redacted]" or "[removed]" marker, not crash.
- **Multi-currency aggregation** (some bookings in IDR, some SAR): aggregated totals converted to IDR using per-transaction fx_rate (from F9 per Q048).
- **Historical data beyond retention window** (Q070): not available; cleanly hidden with message "beyond retention period."

## Data & state implications

F11 **does not own canonical data**. Every metric traces to an owning service. The service-owned aggregation endpoints are the primary API surface.

Proposed F11-owned lightweight tables (minimal, for operational state not-already-captured elsewhere):

- `dashboard_alert_thresholds` — per-alert-type threshold values + recipient_roles + snooze_rules + last_triggered. Per Q068 default: populated with fixed defaults on install; configurable by Super Admin.
- `dashboard_alert_events` — log of threshold breaches + delivery status + ack timestamp. Retention per Q070.
- `dashboard_exports` — audit log of exports (user, dashboard, format, timestamp, filter state).
- `dashboard_preferences` — per-user saved filters + widget layout (MVP: shared layout per role; user customization in Phase 2).
- **No fact tables in F11** — aggregations executed live against owning services.

Aggregation architecture per **Q066** (answered): **service `GET /v1/metrics/*` + thin `dashboard-svc` composer** + **Redis 5m default cache** + **per-service matviews when >500ms** + **Grafana on read-replica** for analyst slices alongside Svelte exec UI.

## API surface (high-level)

F11 itself has a minimal API surface — mostly a thin aggregator + presentation layer. Most work happens upstream.

**REST** (F11 thin layer — implemented in `dashboard-svc`, per **Q066**):

- `GET /v1/dashboard/home` — executive landing; returns KPI-grid structured data.
- `GET /v1/dashboard/operations` — Operational Readiness widgets.
- `GET /v1/dashboard/sales` — Sales & Marketing Board.
- `GET /v1/dashboard/saudi` — Real-Time Field View (initial payload; websocket for updates).
- `GET /v1/dashboard/inventory` — Inventory Executive.
- `GET /v1/dashboard/finance` — Financial Health composite.
- `GET /v1/dashboard/alerts` + `POST /v1/dashboard/alerts/thresholds` (admin).
- `POST /v1/dashboard/export` — PDF/Excel/image export job (async for heavy ones).
- `WebSocket /v1/dashboard/live` — subscribes to streaming data (incidents, bus positions, cash flow ticker).

**Upstream service `/metrics` endpoints** (per-service aggregation surfaces that F11 calls):

- `finance-svc`: `GET /v1/metrics/cash-flow`, `/balance-sheet-summary`, `/ar-aging`, `/ap-aging`.
- `logistics-svc`: `GET /v1/metrics/inventory-value`, `/critical-stock`, `/po-outstanding`, `/paid-unshipped`.
- `catalog-svc`: `GET /v1/metrics/seat-inventory`, `/package-sold-mtd`.
- `booking-svc`: `GET /v1/metrics/bookings-funnel`, `/conversion-rate`.
- `payment-svc`: `GET /v1/metrics/revenue-mtd`, `/payment-method-breakdown`.
- `crm-svc`: `GET /v1/metrics/cs-performance`, `/ads-spend-roas`, `/leaderboard`.
- `visa-svc` + `ops-svc`: `GET /v1/metrics/vendor-readiness`, `/verification-backlog`, `/raudhah-success-rate`.

All `/metrics` endpoints respect branch-scope via JWT claim.

## Dependencies

- **F1** (IAM) — RBAC + branch-scope enforcement on every query.
- **F2** (catalog) — seat inventory, packages master.
- **F4** (booking) — funnel metrics, booking state.
- **F5** (payment) — revenue, collections, refunds.
- **F6** (visa) — vendor readiness (visa filing state), Raudhah Shield.
- **F7** (ops) — vendor readiness (manifest, grouping), verification backlog, field radar (bus GPS, tasreh scans, incidents, luggage).
- **F8** (logistics) — inventory health, PO monitor per F8 W15; binds to Q038 valuation.
- **F9** (finance) — financial reports W15, cash flow W16, aging W17; F11 presents, F9 computes.
- **F10** (CRM) — Super-View references, CS performance, ROAS, leaderboard.
- **External** — Grafana (if Q066 Option C), PDF renderer (same worker pool as F2 flyer / F5 receipt / F7 manifest).

## Backend notes

- **Aggregation co-location** per Q066 Option A default: each service owns its `/metrics` endpoints using sqlc queries against its own Postgres schema (per ADR 0007 single-DB-multi-schema model). Materialized views for expensive rollups (e.g. daily revenue sum); refresh cadence per Q067.
- **Caching layer** — thin Redis-like cache at F11 composition layer for hot paths (home dashboard hit every session); TTL aligned with Q067 cadence.
- **Websocket surface** — ops-svc (incidents + bus tracking) publishes to a pub/sub channel; F11 maintains WS connections and fans out to dashboard clients.
- **Branch-scope enforcement** — every `/metrics` query includes branch_id param; F1 ValidateToken resolves scope per JWT. Central roles receive `branch_id=null` meaning "all branches."
- **Drill-through** — standardized filter params across all dashboards (branch, date_range, product, agent, campaign); preserved via URL; backend `/metrics` endpoints all accept the same filter shape.
- **Export worker** — async job queue: PDF / Excel generation offloaded to worker pool (shared with F2 flyer / F5 receipt / F7 manifest); returns download URL when ready.
- **Alert engine** — small cron (every 60s) sweeps configured thresholds; compares current value against threshold; fires notification via crm-svc broadcast API.
- **Observability** — F11 queries themselves traced via OTel; slow-query detection on individual `/metrics` endpoints helps target optimization.

## Frontend notes

- **Svelte dashboard shell** — single-page-app with route-based dashboard loading; skeleton loaders while data fetches.
- **Chart library** — picked at implementation time (candidates: ApexCharts, Chart.js, Recharts-equivalent for Svelte, Layercake).
- **Map component for Radar Bus (#181)** — Leaflet or Mapbox; marker clustering for scale.
- **Real-time widgets** — websocket subscription pattern; graceful reconnect + stale-indicator badge.
- **Mobile-first** — every dashboard must render useful content on 375px width; complex layouts collapse to vertical stack.
- **Dark mode** — optional but mentioned as common executive preference; respect system setting.
- **Print CSS** — for PDF export; single-page prints for board reports.
- **Localization** — dashboards bilingual Indonesian + English; numeric formatting locale-aware (IDR comma/period conventions).

## Open questions

None blocking — **Q038, Q045, Q048** (upstream) and **Q066–Q075** answered **2026-04-18** (`docs/07-open-questions/`). Cadence, drill depth, export policy, and transport defaults in the body match those answers unless config overrides.

**Implementation notes (not open decisions):**

- Keep Grafana for analyst ad-hoc; product dashboards stay in Svelte + `dashboard-svc` per **Q066**.
