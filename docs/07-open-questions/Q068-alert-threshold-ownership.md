---
id: Q068
title: Alert threshold ownership + default values
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: open
---

# Q068 — Alert threshold ownership

## Context

F11 W8 lists alert types (cash balance low, AR overdue, CPL high, critical stock, open incidents, paid-unshipped aging). PRD mentions alerts in two modules (#179 ad-spend alert, #189 aging) but gives no threshold numbers.

Threshold ownership matters: hard-coded defaults require code deploys to change; fully configurable per-tenant adds operational surface; per-role-configurable is flexible but complex. Without thresholds, alerts don't fire at all.

## The question

1. **Ownership model** — fixed defaults, per-tenant config, per-role config, or individual-user config?
2. **Who configures** — Super Admin only, finance director for finance alerts, ops lead for ops alerts?
3. **Default values** at install — what should they be out-of-box?
4. **Alert delivery target** — role-based (all directors), role+individual (specific person), or per-alert configurable?
5. **Snooze / mute** — can a recipient mute an alert type temporarily?
6. **Alert categories** — financial, operational, marketing, field-ops — owned by different roles?
7. **Alert audit** — every threshold change + every fire logged?

## Options considered

- **Option A — Fixed defaults; per-tenant configurable by Super Admin only.** One global set per agency, tunable via admin UI.
  - Pros: simple; one authority; audit trail clean.
  - Cons: role-specific tuning not possible (ops lead can't raise paid-unshipped threshold without Super Admin).
- **Option B — Per-role configurable by category owner.** Finance-category alerts tunable by finance director; ops-category by ops lead; etc.
  - Pros: owners can tune their own alerts without bottlenecking on Super Admin.
  - Cons: more config UI; need to map categories to owners.
- **Option C — Individual-user customizable on top of tenant defaults.** Override per-user (e.g. "I want paid-unshipped alert at 5 days, not 7").
  - Pros: maximum flexibility.
  - Cons: alert storm when everyone tunes differently; operational chaos; rare use case in practice.

## Recommendation

**Option B — per-role configurable by category owner + fixed sensible defaults at install + Super Admin override.**

Option A's one-authority bottleneck is an operational anti-pattern — finance director should tune AR alert threshold without pinging Super Admin. Option C's individual-user overrides are rarely asked-for and create chaos. Option B puts tuning authority where operational knowledge lives: finance owns financial thresholds, ops owns operational thresholds, marketing owns marketing thresholds; Super Admin can override anyone for governance.

Defaults to propose:

| Alert type | Default threshold | Owner role |
|---|---|---|
| Cash balance low | per bank account: < 50M IDR | finance director |
| AR overdue | > 60 days, balance > 10M IDR | finance director |
| AP aging overdue | > 30 days past due-date | finance director |
| CPL high | CPL > 2× rolling-7-day-avg with < 3 closings | marketing admin |
| Critical stock | SKU quantity ≤ reorder_level | warehouse supervisor |
| Open incidents | > 3 incidents in 60min same departure | ops lead |
| Paid-unshipped aging | > 7 days since paid | warehouse supervisor |
| ROAS low | ROAS < 1.5× on active campaign > 7 days running | marketing admin |

Alert delivery: category owner role receives (all members of role); one individual from role can additionally be designated primary recipient (e.g. specific finance director user). Snooze: 24h temporary mute per user per alert-type; audit logged. Thresholds editable by category owner + Super Admin; change logged in `dashboard_alert_thresholds` audit. Alert fire events logged in `dashboard_alert_events`.

Reversibility: threshold values, ownership mapping, and delivery targets are all config.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (governance model), finance director / ops lead / marketing admin (category-owner role).
