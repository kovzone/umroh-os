---
id: Q075
title: Executive landing widget composition (top 8–12 KPIs)
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: open
---

# Q075 — Executive landing widget composition

## Context

PRD lines 25–27 and 485 describe the executive home (`/erp/dashboard`) as a *Helicopter View* where the owner sees the agency's state in "one phone screen." But the PRD **never lists** the specific widgets — no top-10 KPI list, no prioritized ordering.

Asking the owner to pick is the right move (they know what they care about); but a good default is needed to ship MVP without blocking on an opinion call.

## The question

1. **Which widgets appear on the executive home?** Top 8-12 KPIs.
2. **Ordering / priority** — most important at top (mobile-first constraint).
3. **Per-persona landing** — same home for CEO, CFO, CMO, COO, or different homes per role?
4. **Time-range default** — MTD, today, last 7 days, custom?
5. **Filter controls** — branch, date range, product — on the home screen?
6. **Widget sizes** — all equal, or hero widgets (revenue big) + ticker widgets (alerts small)?
7. **Cold start** — empty-state messaging on a brand-new agency with no data yet.

## Options considered

- **Option A — Universal home + role-aware sub-landings.** One agency-wide home (owner-curated KPIs); CFO / CMO / COO have their own per-role "my dashboard" accessible from a role-switcher.
  - Pros: unified home for agency owner; specialized access for functional leaders.
  - Cons: needs curation of the universal home; role sub-landings add surface.
- **Option B — Single universal home for all director+.** Everyone sees the same thing.
  - Pros: simplest; one curation effort.
  - Cons: CFO doesn't care about CS performance; CMO doesn't care about AR aging.
- **Option C — Personalized landing per user** (building on Q073's Phase 2 customization).
  - Pros: each user sees what they care about.
  - Cons: requires Q073 customization (deferred to Phase 2).

## Recommendation

**Option A — universal executive home (owner-curated KPI set) + role-aware sub-landings accessible from a role-switcher menu; MVP ships with a sensible default for the executive home pending stakeholder curation.**

Option B is simpler but fights how executives actually work — finance, marketing, ops leaders have different priorities and wasting their screen-real-estate on other-domain metrics reduces dashboard utility. Option C is customization territory (Q073) and deferred. Option A ships with a curated default (this question's deliverable) and lets each director-tier role go to their own sub-landing for deeper work.

Defaults to propose (this is a *provisional* ordering pending stakeholder agreement — expect curation):

| Priority | Widget | Size | Source |
|---|---|---|---|
| 1 | Revenue MTD vs target | Large (hero) | F5 + F9 |
| 2 | Cash balance (consolidated) | Large (hero) | F9 W16 |
| 3 | Bookings MTD + conversion rate | Medium | F4 + F10 |
| 4 | Departures next 30 days (count + vendor-readiness status) | Medium | F7 Eksekusi Vendor |
| 5 | Paid-unshipped queue (count + oldest age) | Medium | F8 + F10 |
| 6 | AR aging (> 30 days sum) | Medium | F9 W17 |
| 7 | Open incidents (count by severity) | Small | F7 incidents |
| 8 | Critical stock count | Small | F8 |
| 9 | CPL (today vs rolling avg) | Small | F10 ads |
| 10 | Active campaigns ROAS summary | Small | F10 |

Role-aware sub-landings: `/erp/dashboard/finance` (CFO default) emphasizes W6 (financial health); `/erp/dashboard/marketing` (CMO default) emphasizes W3 (sales board); `/erp/dashboard/operations` (COO default) emphasizes W2 + W4; `/erp/dashboard/saudi` (during-trip ops lead) emphasizes field view.

Time-range default = MTD (month-to-date); user can toggle (Today / 7d / MTD / 30d / QTD / YTD / Custom). Filter controls on home: branch (if scope allows), date range. Empty-state: on cold-start tenant, each widget shows "Data will appear once transactions begin" with a link to relevant onboarding (e.g. "Create your first package" link from revenue widget).

Reversibility: widget list + ordering + sizes are fully config; no code change to swap them.

## Answer

TBD — awaiting stakeholder input. **Stakeholder curation explicitly needed** — the default widget list above is a starting point, not a final answer. Deciders: agency owner (primary consumer of the home), CFO / CMO / COO (sub-landing composition), UX / design lead (mobile-first composition).
