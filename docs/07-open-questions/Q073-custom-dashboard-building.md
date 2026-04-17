---
id: Q073
title: Custom dashboard building vs fixed catalog
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: open
---

# Q073 — Custom dashboard building vs fixed catalog

## Context

PRD describes a fixed catalog of dashboards (modules #177–#189) with specific names and contents. It does **not** mention user-built dashboards (analysts constructing their own KPI mash-ups, executives rearranging tiles).

Custom dashboard building is a significant feature:
- **Drag-drop UI** for widget placement.
- **Widget catalog** users pick from.
- **Filter persistence** per user/role.
- **Sharing** (save a custom dashboard to team).
- **Template / blueprint** governance.

Teams with Looker / Metabase / Power BI expect this. PRD's silence on it is interesting — typical Indonesian SME ERPs don't offer it.

## The question

1. **Build-your-own dashboards** — MVP feature, Phase 2, or never?
2. **Layout personalization** — even without full build-your-own, can users hide / rearrange tiles on fixed dashboards?
3. **User-saved filters** — can users save a filter preset and revisit?
4. **Sharing custom dashboards** — user-shared, team-shared, public?
5. **Template governance** — if custom dashboards exist, who moderates (ensure scope gates, prevent data-leak risk)?
6. **Grafana fallback** — since Grafana is already in the stack per Q066 recommendation, does "custom dashboard" just = "use Grafana"?

## Options considered

- **Option A — Fixed catalog only for MVP; no customization; Grafana available for ad-hoc.** Strict MVP scope.
  - Pros: fastest to ship; matches typical Indonesian SME ERP norm; no moderation burden.
  - Cons: power users stuck with fixed layout.
- **Option B — Fixed catalog + user-level tile customization (hide, rearrange).** No new widgets, but users personalize their fixed dashboards.
  - Pros: moderate flexibility; low implementation cost.
  - Cons: still no new metrics.
- **Option C — Full custom dashboard builder — user picks widgets, arranges, saves.** Full BI-tool feature set.
  - Pros: power-user experience; one-stop product.
  - Cons: major implementation; moderation burden; duplicates Grafana for marginal gain.

## Recommendation

**Option A — fixed catalog only for MVP; Grafana on read-replica for ad-hoc (already recommended in Q066); per-user saved filters as a lightweight personalization in Phase 2.**

Option C's full builder is a major product feature that would dominate the F11 implementation effort — hard to justify when the modules are mostly Should Have priority. Option B's tile customization is tempting but still requires UX, layout-persistence schema, and data migration if widgets change — for marginal value. Option A accepts the fixed catalog constraint, which is what the PRD describes, and relies on Grafana for ad-hoc needs (analysts + power users get their own tool, executives get the curated product dashboards).

Defaults to propose: **MVP** = fixed dashboard catalog; no drag-drop; no widget-selection UI. **Phase 2** = per-user saved filters (date range, branch scope where applicable); user-level tile hide/show (but no rearrange). **Phase 3** or never: full custom dashboard builder — revisit if executives/analysts consistently ask (based on support tickets). **Grafana posture** (Q066) handles ad-hoc: analysts + CTO + operations can build custom Grafana dashboards against the read replica. Product dashboards (for executives) remain fixed and curated.

Reversibility: adding tile customization later is additive; adding full builder later is significant work (new feature).

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (dashboard UX maturity), CTO (implementation scope), power users if retained (analyst feedback).
