---
id: Q070
title: Historical data retention window for dashboards
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: open
---

# Q070 — Historical retention window

## Context

Dashboards expose historical data for trend analysis, YoY comparisons, quarterly reviews. But unlimited history carries costs: Postgres table size, slow aggregation queries, expensive materialized view refreshes, long index builds. Plus Indonesian tax law already requires 10-year accounting record retention (UU KUP Pasal 28) for the underlying journal data — dashboard access is a separate question.

PRD doesn't specify retention. Candidate policies range from "1 year of dashboard data" (lean) to "all-time from installation" (complete but costly).

## The question

1. **Retention window length** — 1 year, 3 years, 7 years, all-time?
2. **Per-dashboard variation** — financial dashboards retain longer (regulatory)? Operational retain shorter (not needed)?
3. **Archive vs delete** — beyond the window, is data archived (cold store) or summarized (aggregates only) or deleted (lost)?
4. **Drill-down depth vs retention** — can executive drill to a 2-year-old transaction from the aggregate, or only see the number?
5. **Report-generator access** — can reports generate against 10-year data even if dashboards don't?
6. **YoY / trend-chart windows** — default comparison period (last year, last quarter)?

## Options considered

- **Option A — 3-year rolling window on dashboards; 10-year on reports (regulatory minimum).** Dashboard queries bounded to 3 years for performance; full reports (Neraca, Laba Rugi export) hit the full retained data.
  - Pros: bounds dashboard performance; respects regulatory retention on the report side.
  - Cons: asymmetry between dashboard and report might confuse some users.
- **Option B — All-time on everything.** No window; dashboards and reports see all data.
  - Pros: maximum visibility; no "why can't I see 4 years ago" confusion.
  - Cons: performance risk as data grows; UI choices dominated by long-tail.
- **Option C — 1 year on dashboards + aggregates beyond; reports see all.** Dashboards show last year's raw data + pre-computed aggregates for older periods.
  - Pros: bounded data volume.
  - Cons: aggregate-only is lossy for ad-hoc drill.

## Recommendation

**Option A — 3-year rolling window on dashboards (raw + drill); reports and export access full 10-year retained data.**

Option B risks performance as volume grows (imagine 30K jamaah across 5 years). Option C's aggregate-only for >1y breaks drill (users hit a wall when they click into an old period). Option A is the pragmatic default: 3 years is enough for multi-year trend analysis on dashboards (season-over-season, Ramadhan vs non-Ramadhan) while bounding query size; reports remain unlimited for regulatory needs.

Defaults to propose: **Dashboard data window** = 3 years from current date; queries `WHERE entry_date >= now() - '3 years'::interval`. **Report window** = no bound — Neraca / Laba Rugi / ad-hoc reports access all retained data. **Archive** = nothing deleted in MVP; data stays in primary tables per 10-year regulatory requirement (UU KUP Pasal 28). Phase 2: cold-storage archival for > 5-year data to optimize storage costs. **Drill access** = within 3 years, drill fully functional; beyond 3 years, drill shows message "Beyond dashboard retention — use reports for detailed view" and links to the reports section. **YoY / trend defaults** = last 12 months vs prior 12 months default comparison; user can select custom range up to 3 years back.

Reversibility: retention window is config (`dashboard.history_window = 3y`); can be extended later without data migration (data is already there).

## Answer

TBD — awaiting stakeholder input. Deciders: finance director (regulatory requirement), CTO (performance posture), agency owner (trend analysis needs).
