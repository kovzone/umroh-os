---
id: Q069
title: Drill-down depth — widget to source transaction
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: answered
---

# Q069 — Drill-down depth

## Context

Dashboards need to let the executive click a KPI and see the story. "Revenue up 20% MTD" → which packages? → which bookings? → which payments? PRD says "akses cepat ke laporan" (line 619) but doesn't specify depth.

Too shallow (1 level) = dashboards as wallpaper; executives guess at causes. Too deep (unbounded) = performance + data-scope complications.

## The question

1. **Standard drill-down depth** — 2 levels? 3 levels? Unlimited?
2. **Per-widget depth** — does every widget drill to the same depth or varies?
3. **Filter preservation** — drill-through preserves upstream filters (date range, branch, product)?
4. **Cross-feature drill** — revenue KPI (F5 payments) drills to booking (F4), which drills to jamaah (F3), which drills to document (F3 again). How deep does the drill path cross features?
5. **Source-transaction view** — does drill land at the raw entity (e.g. the specific payment receipt) or a wrapped detail page?
6. **Data-scope on drill** — branch manager drilling into a regional roll-up: stays scoped to their branch only?
7. **Permissioned fields** — some fields (jamaah NIK, salary) may be restricted even after drill; per-field RBAC?

## Options considered

- **Option A — 3-level drill standard across all widgets: widget → breakdown chart → filtered transaction list → source detail modal.** Consistent; bounded.
  - Pros: predictable UX; bounded backend complexity; covers 90% of executive needs.
  - Cons: rigid; some widgets naturally drill deeper (e.g. booking → jamaah → documents).
- **Option B — Per-widget configurable depth; some widgets drill 2 levels, some 5.** Flexible.
  - Pros: tailored UX per data domain.
  - Cons: inconsistent; every widget needs its own drill definition.
- **Option C — Unlimited drill following foreign keys (like a wiki).** Every entity page has its own drills.
  - Pros: maximally exploratory.
  - Cons: executives get lost; performance unpredictable; scope-leakage risk.

## Recommendation

**Option A — 3-level standard drill; per-widget customization allowed when the data domain justifies (like booking → jamaah → documents).**

Option B's full per-widget customization costs implementation + documentation effort not justified by typical executive usage. Option C's wiki-style unlimited navigation works for analysts but overwhelms executives — the dashboard becomes a tool for exploration rather than decision-making. Option A lands on a predictable 3-level pattern: aggregate widget → breakdown → detail, with per-domain customization as an exception (not the rule).

Defaults to propose: **Standard drill** — 3 levels, every widget. Level 1: KPI value in widget. Level 2: breakdown by dimension (time / category / agent / branch / product). Level 3: filtered transaction list with click-to-detail modal for individual records. **Filter preservation** — every drill carries upstream filters via URL parameters; back button works. **Cross-feature drill** — allowed within 3 levels; e.g. revenue → per-booking list → booking detail (covers F5 → F4). Deeper inspection (booking → jamaah → documents) happens via navigating from the booking detail page into F3's own screens (separate URL; not within dashboard drill). **Source view** — level 3 lands in a detail modal summarizing the entity; "View Full Record" link takes to entity's primary view in its owning feature. **Data-scope on drill** — fully preserved; branch manager drilling into aggregate sees only their branch's contributing records. **Permissioned fields** — F1's field-level RBAC (if implemented) applies to drill views same as primary views; restricted fields show as `[permission required]`.

Reversibility: extending drill depth later is additive; tightening is a policy decision not data-schema.

## Answer

**Decided:** **Option A** — **3 dashboard drill levels** then **deep link** to owning feature screens; **filters in querystring**; **branch scope enforced**; **field RBAC** respected.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
