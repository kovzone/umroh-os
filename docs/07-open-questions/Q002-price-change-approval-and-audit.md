---
id: Q002
title: Price-change approval thresholds and audit policy
asked_by: session 2026-04-14 F2 draft
asked_date: 2026-04-14
blocks: F2
status: open
---

# Q002 — Price-change approval thresholds and audit policy

## Context

F2 mass-update workflow (module #77) lets an admin change price or status across many packages in one action. The PRD does not specify approval gates for this — yet Section H ("Admin & Security") and module #65 (e-Approval for discounts) both imply that sensitive financial changes should not be self-approved.

The risk is a single admin inadvertently setting all active packages to 0 IDR or cancelling every departure. We need to pin down the threshold and approval pattern before shipping the mass-update endpoint.

## The question

1. Should mass price/status updates require **second-admin approval** above a configurable threshold? If yes, what are the default thresholds — e.g. >10 packages affected, or >5% price delta, or absolute-value change > Rp X?
2. Is the audit requirement **price-history row per affected package** (my default assumption), or a lighter single "mass-update batch" record that links to all affected rows?
3. Does **single-package price edit** also need approval, or only mass updates? And if single edits are unrestricted, does the threshold for "mass" start at 2 packages or higher?

## Options considered

- **Option A — Threshold-based approval (Recommended default).** Single-package edits: self-approve. Mass updates affecting >10 packages OR >5% price delta: require second-admin approval, queued as `pending_approval`. Audit: per-package price_history row in all cases.
  - Pros: matches Indonesian ERP norms, protects against accidental mass changes.
  - Cons: slightly more UI complexity (approval queue).
- **Option B — All price changes require approval.** Every edit, single or mass, needs a second admin.
  - Pros: tightest control.
  - Cons: operational drag; admins will route around with "approval parties".
- **Option C — No approval, strong audit only.** Trust admins; rely on audit trail + daily anomaly alerts (module #158).
  - Pros: fastest iteration.
  - Cons: no safety net for mass mistakes.

## Recommendation

**Option A — threshold-based approval with per-package price-history.** Single-package edits are self-approved (normal daily work shouldn't be gated); mass updates crossing the threshold require a second admin (catches the "oops I just zeroed out every package" case). Per-package price-history rows — not a batched mass-update record — because finance, agents, and customers all ask "why did *this* package's price change?" and pointing them at a batch record is a worse answer than a direct row.

Default thresholds: **10 packages** or **5% absolute price delta**, whichever triggers first. Both reviewer-configurable via the system config. These numbers are starting points, not convictions — tune after two months of real operator behaviour.

Option B (all changes gated) is tempting but Indonesian ERP users will route around it with "approval parties" that devalue the control. Option C (no gate, audit-only) is too optimistic for a system with multi-currency money in play.

Reversibility: thresholds are config values, trivially tunable. Moving from Option A to B later is a one-flag flip. Moving from A to C (removing the gate) is also trivial. Low commitment.

## Answer

TBD — awaiting stakeholder discussion.
