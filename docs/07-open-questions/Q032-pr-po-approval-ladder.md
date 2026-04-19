---
id: Q032
title: PR / PO approval threshold ladder
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8
status: answered
---

# Q032 — PR / PO approval threshold ladder

## Context

PRD module #111 *Persetujuan Berjenjang* (line 377) says Purchase Requisition approval is "Manager → Director, dashboard / mobile" — but gives no IDR thresholds, no category rules, no fallback for absent approver, and no SLA. F8 W2 (multi-level PR/PO approval) cannot ship without a concrete ladder.

Approval ladders are a governance concern. Under-specifying risks shadow approvals (the procurement officer implicitly self-approves because no threshold triggers). Over-specifying blocks legitimate spend (every SKU replenishment goes to the CEO). The right ladder reflects how the business actually signs off today, plus guardrails for audit.

## The question

1. **Threshold IDR amounts per approver level** — what range triggers which approver? E.g. ≤X: Manager only; X–Y: Director; > Y: Director + CEO.
2. **Category-based rules** — does approval differ for capex (new warehouse shelving) vs opex (koper restock) vs emergency (last-minute replacement)?
3. **Fallback when approver is absent** — is there a named deputy? A timeout that auto-escalates?
4. **SLA for approval** — how long can a PR sit in the queue before overdue? Auto-escalate or auto-reject?
5. **Multi-line PR splits** — if a PR has lines that straddle thresholds, does the whole PR go to the higher approver, or do lines split?
6. **High-value approval extra requirements** — does > 50M IDR require digital signature / biometric confirm / dual approval beyond a single click?

## Options considered

- **Option A — Flat two-tier: Manager → Director above a single threshold.** One IDR threshold (e.g. 10M). Below: Manager. Above: Director. Matches PRD text literally.
  - Pros: easy to implement, easy to understand.
  - Cons: CEO never involved on large capex; no category handling; fragile for high-value procurement.
- **Option B — Three-tier with category-aware routing.** ≤10M Manager; 10–50M Director; >50M Director + CEO. Capex always Director minimum regardless of amount. Emergency category bypasses Director if Manager is reachable (post-hoc Director approval within 48h).
  - Pros: reflects real governance; handles edge cases.
  - Cons: more config surface; category classification required on every PR.
- **Option C — Fully configurable admin-managed approval matrix.** Admin defines rules per (amount range, category, warehouse) → approver chain. Read at submit time.
  - Pros: maximum flexibility.
  - Cons: config complexity; misconfiguration risk; PR matching no rule is undefined.

## Recommendation

**Option B — three-tier with category routing + 48-hour SLA + delegation-to-deputy.**

Real procurement has three-tier patterns for good reason: Manager signs day-to-day replenishment, Director signs anything material, CEO sees large capex. Option A quietly pushes Director into rubber-stamp mode for 5M IDR koper restocks — or worse, keeps CEO uninvolved on 100M IDR warehouse fit-outs. Option C (fully configurable) is tempting but tends to rot: admins tune it once then never revisit, leaving gaps. Three tiers plus a category flag hit the 80/20.

Defaults to propose: ≤10M IDR Manager; 10–50M Director; >50M Director + CEO. Capex category always minimum Director. Emergency allows Manager-only initial with Director post-hoc within 48h. Approver-absent: auto-delegation to named deputy (config per role). SLA 48h then auto-escalation email to next tier. Multi-line PRs route by the highest-value line's threshold. >50M approvals require second factor (biometric or TOTP).

Reversibility: the ladder is config, not code — adjusting thresholds later is an admin action, not a deploy. Category classification must be baked into the PR schema now; adding it later is a migration.

## Answer

**Decided:** **Option B** ladder — **≤10M Manager**, **10–50M Director**, **>50M Director+CEO**; **capex ≥10M → Director minimum**; **emergency Manager path** with **Director post-hoc within 48h**; **48h approver SLA** then auto-escalate; **multi-line PR routes by highest line threshold**; **>50M requires MFA** on approval action; **named deputy** per role in config.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
