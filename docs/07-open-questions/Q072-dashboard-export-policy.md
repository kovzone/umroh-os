---
id: Q072
title: Dashboard export policy — formats and permissions
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: answered
---

# Q072 — Dashboard export policy

## Context

PRD mentions financial reports can be exported (F9 W15 describes PDF / Excel / CSV) but doesn't specify policy for F11 dashboards: which dashboards can be exported, to what format, by whom.

Concerns:
- **Regulatory need** — Neraca / Laba Rugi for audit meetings; must export.
- **Confidentiality** — executive dashboards contain sensitive data (revenue, commission cost, agent names); uncontrolled export risks leak.
- **Real-time fit** — a "snapshot of the bus radar" doesn't make sense; live data doesn't export well.
- **Operational need** — printing paid-unshipped queue for warehouse morning standup.
- **Audit trail** — who exported what, when.

## The question

1. **Which dashboards export** — all, some, regulatory-only?
2. **Which formats** — PDF, Excel, CSV, PNG snapshot?
3. **Permission by role** — can branch managers export? CS? Only directors?
4. **Audit** — every export logged?
5. **Watermark / confidentiality footer** — include user + timestamp watermark on exports?
6. **Scheduled / recurring exports** — nightly financial report to CFO email? MVP or Phase 2?
7. **Export size limits** — max rows / pages to prevent abuse?
8. **Real-time data in exports** — if bus radar is exported, does it capture the point-in-time or refuse?

## Options considered

- **Option A — Everything exportable; format depends on widget type (financial → PDF + Excel; operational → PDF image; field-live → no export).** Broad access.
  - Pros: flexibility; minimal friction.
  - Cons: confidentiality risk; operational abuse (100MB exports).
- **Option B — Export whitelist per dashboard category; role-gated.** Only financial + inventory + sales dashboards export; by director-tier role.
  - Pros: tight control; audit-friendly.
  - Cons: operational use cases blocked (branch manager can't print paid-unshipped for morning standup).
- **Option C — Hybrid: whitelist + role-gate + audit all exports + watermark.** Export allowed on most dashboards, role-scoped, always audit-logged, always watermarked.
  - Pros: operational access preserved; tight audit + confidentiality footprint.
  - Cons: more implementation work.

## Recommendation

**Option C — hybrid: export whitelist + role-gated + audit-logged + confidentiality watermark.**

Option A accepts confidentiality risk (e.g. a CS exports a commission leaderboard and shares externally). Option B is too restrictive — paid-unshipped print is a legitimate use case. Option C maintains operational access while mitigating risk: every export is traceable to a user; watermark with user+timestamp discourages casual leaking.

Defaults to propose:

| Dashboard | Export formats | Role required |
|---|---|---|
| Financial Health (W6: #187, #188, #189) | PDF + Excel | finance director, CFO, owner |
| Inventory Executive (W5: #185, #186) | PDF + Excel | warehouse supervisor, ops lead, owner |
| Operational Readiness (W2: #177, #178) | PDF | ops lead, owner |
| Sales & Marketing Board (W3: #179, #180) | PDF + Excel | marketing admin, owner |
| Real-Time Field View (W4: #181–#184) | **no export** (live data; snapshot meaningless) | n/a |
| Executive Home (W1) | PDF snapshot | directors, owner |
| Branch-scoped dashboards (W7) | PDF | branch manager (own branch only) |

All exports auto-watermarked with `Exported by: <user> <timestamp> — CONFIDENTIAL`. All exports audit-logged in `dashboard_exports` (who, what, when, filters). Recurring scheduled exports = Phase 2. Size limit: PDF max 50 pages; Excel max 50K rows; request exceeds → error with refinement hint. Real-time field dashboards simply don't offer an export button.

Reversibility: whitelist and role mapping are config; adding PDF-only to a currently-non-exportable dashboard is trivial.

## Answer

**Decided:** **Option C** whitelist matrix as Recommendation + **watermark** + **`dashboard_exports` audit** + limits **50 pages / 50k rows** + **no live field-map exports**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
