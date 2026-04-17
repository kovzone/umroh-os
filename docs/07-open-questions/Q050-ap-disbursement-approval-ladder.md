---
id: Q050
title: AP disbursement approval threshold ladder
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: open
---

# Q050 — AP disbursement approval threshold ladder

## Context

PRD module #136 *Otorisasi Pembayaran* (line 453) says vendor payments follow tiered e-approval "dari Manajer hingga ke Direktur" — same underspecified ladder as Q032 for procurement (PR/PO approval). But **PR approval** and **AP disbursement approval** are different governance moments: PR is commit-to-spend, AP disbursement is cash-out-now. The approval ladders may share shape but tend to have different thresholds (cash-out scrutiny is typically stricter).

## The question

1. **Threshold IDR amounts per approver level for AP disbursement** — ≤X: Manager only; X–Y: Director; > Y: Director + CEO.
2. **Is the AP ladder independent from the PR ladder (Q032), or linked?** — an already-approved PR's subsequent disbursement may inherit approval, OR require a fresh sign-off at payment time.
3. **Cumulative batch approval** — when the AP officer batches 20 vendor payments totaling 500M IDR, does the approval fire against the batch total or against each payment individually?
4. **Recurring payment approval** — monthly hotel retainers, quarterly insurance premiums; do these require per-payment approval or a one-time annual approval?
5. **Urgent / emergency disbursement bypass** — same-day vendor demand (e.g. visa emergency fee); is there an expedited approval path?
6. **Dual-control requirement** — do high-value payments require two approvers (not just one director) or a single-director approval suffices?
7. **Biometric / second-factor** — like Q032's high-value 2FA, does > 50M IDR require TOTP / biometric beyond click-approve?

## Options considered

- **Option A — Same thresholds as Q032 PR ladder; PR approval implies AP approval.** Once PR is approved, the subsequent AP disbursement on that PO doesn't need fresh sign-off.
  - Pros: simple; avoids approval fatigue.
  - Cons: weak cash-out control; an approved PR's final payable may differ from approved amount (price changes, GRN variances); director loses cash-out visibility.
- **Option B — Independent AP ladder, stricter than PR.** AP-disbursement thresholds lower than PR thresholds (e.g. ≤5M Manager; 5–30M Director; >30M Director + CEO). Each payment needs fresh sign-off.
  - Pros: strong cash-out control; catches variance at payment time.
  - Cons: approval fatigue; slower payments; agent / vendor relationship risk.
- **Option C — Same thresholds as PR, but fresh AP sign-off required for variances > N%.** PR approval implies AP approval UNLESS GRN amount varies > 10% from PR approved amount; then fresh sign-off.
  - Pros: balances speed and control; targets approval friction where actually needed.
  - Cons: more complex logic; variance threshold is arbitrary.

## Recommendation

**Option C — PR approval implies AP approval by default; fresh sign-off only for variances > 10% or for AP entries not linked to an approved PR (e.g. recurring payments, one-off services without PR).**

Option A is too loose — a 100M IDR PR approved at 8M GRN should not auto-disburse at 100M; variance is exactly where governance should engage. Option B's fresh-sign-off-everywhere creates approval fatigue and slows legitimate payments; the agency owner has better things to do than click-approve every 2M catering payment when the original 100M catering PR was already approved. Option C delegates the majority of payments (matching the PR) while escalating variances and PR-less payments.

Defaults to propose: AP disbursement uses same thresholds as Q032's PR ladder (≤10M Manager; 10–50M Director; >50M Director + CEO). Variance guard: AP disbursement amount > 110% of PR approved amount → fresh Director sign-off regardless of disbursement amount. AP disbursement without a linked PR (recurring payments, emergencies, one-off services < invoice threshold) → fresh approval per the AP ladder. Batch approval: applies against batch total (so a 20-invoice 500M IDR batch needs Director + CEO); but the batch surface lists per-line amounts for scrutiny. Recurring payments (identified by `is_recurring` flag on AP entry): one-time annual approval by Director with monthly auto-disbursement up to an annual cap; cap breach = fresh approval. Emergency: 4h expedited path (Manager can initiate; Director must confirm within 24h asynchronously for audit trail). Dual-control: not required for single-director approvals below CEO tier; > 50M CEO approval IS dual-control (Director + CEO both). Biometric / TOTP for > 50M approvals.

Reversibility: approval thresholds are config. Adjusting variance % or recurring-cap is an admin action.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (approval authority), finance director, external auditor (governance adequacy).
