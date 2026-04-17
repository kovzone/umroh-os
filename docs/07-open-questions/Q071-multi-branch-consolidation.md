---
id: Q071
title: Multi-branch consolidation rule for central visibility
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: open
---

# Q071 — Multi-branch consolidation rule

## Context

PRD line 1357 establishes data-scope tiers (GLOBAL / BRANCH / OWN) and confirms branch managers see only their branch. But the PRD doesn't specify what GLOBAL scope actually sees:

- **Aggregated totals only** — sum across all branches, no per-branch breakdown in UI.
- **Aggregated + per-branch breakdown with drill** — agency total + table of per-branch contributions; click to drill into branch's slice.
- **Flat union** — every transaction visible regardless of branch; filtering optional.

Each has implications for owner/director UX and performance.

## The question

1. **GLOBAL-scope default view** — aggregate only / aggregate + breakdown / flat union?
2. **Per-branch drill from GLOBAL** — can owner click into a specific branch's dashboard?
3. **Branch-comparison views** — side-by-side branch comparison (branch A vs branch B) — MVP or Phase 2?
4. **Branch hierarchy** — if branches have sub-branches (or a region > branch > sub-branch structure), how does consolidation work?
5. **Non-branch-scoped data** — some data is agency-wide (master data, package catalog, aggregate operational metrics); how does this surface?
6. **Permission for central user to "become" a branch user** (view-as) — is this allowed? For training / support scenarios?

## Options considered

- **Option A — Aggregated + per-branch breakdown with drill-to-branch.** GLOBAL sees total + per-branch table; click a branch row to see that branch's dashboard scoped.
  - Pros: clear agency-wide view + operational visibility per branch; matches executive mental model.
  - Cons: slightly more complex UI than pure aggregate.
- **Option B — Aggregated totals only; no per-branch breakdown.** Owner sees only summed values.
  - Pros: simplest.
  - Cons: owner can't see which branch is underperforming without switching scopes; poor operational value.
- **Option C — Flat union with optional filter.** Every transaction listed; owner can filter by branch.
  - Pros: maximal flexibility.
  - Cons: overwhelming for high-level view; defeats the dashboard purpose (should summarize, not list).

## Recommendation

**Option A — aggregated + per-branch breakdown with drill; branch-comparison views deferred to Phase 2; no view-as in MVP.**

Option B hides the operational insight owners actually need (which branch is bleeding money?). Option C is a raw-data surface, not a dashboard. Option A is the textbook executive-dashboard pattern: total at top, contributing segments listed below, click to drill.

Defaults to propose: **GLOBAL view** = aggregated total at widget level + per-branch breakdown table (or bar chart) showing contribution per branch. Click a branch row → navigates into the branch-scoped dashboard for that branch; scope is visually indicated ("Viewing: Cabang Jakarta"). **Branch-comparison** = Phase 2; MVP users can navigate sequentially. **Branch hierarchy** = flat for MVP (one level, no sub-branches); multi-level hierarchy Phase 2. **Non-branch-scoped data** — master data (packages, muthawwif, hotels) visible to all scopes; aggregate operational metrics (total agency revenue) shown only to GLOBAL scope, not summed at branch level to avoid double-count. **View-as** — not allowed in MVP (security posture); central users navigate with their own GLOBAL scope. Scope gating per PRD line 1357 enforced consistently.

Reversibility: consolidation rule is architectural but bounded — moving from Option A to Option C later is a UI add (new "Flat view" toggle); going backward from C to A is harder.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (visibility preference), CFO (financial-consolidation needs), ops lead (operational-branch-comparison needs).
