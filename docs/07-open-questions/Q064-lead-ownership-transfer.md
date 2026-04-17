---
id: Q064
title: Lead ownership transfer between agents / CS
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: open
---

# Q064 — Lead ownership transfer

## Context

A lead can touch multiple owners during its lifecycle: agent A captured via replica site, passed to CS for closing, CS hands to agent B for follow-up, CS admin reassigns to another CS on SLA miss. Each transfer has commission, attribution, and accountability implications.

PRD doesn't specify transfer rules. Questions:
- Who can transfer a lead (self-owner, supervisor, admin)?
- Does a transfer reset SLA timers?
- Does a transfer change commission attribution (Q057 binding)?
- What audit is captured?
- Can the transferred-from owner object?

This is adjacent to Q057 (attribution) — but Q057 is about attribution decisions at booking; Q064 is about operational ownership during the lead lifecycle.

## The question

1. **Who initiates transfer** — lead's current owner, supervisor, admin, or all three depending on scenario?
2. **Self-transfer** — can a CS or agent pass their own lead to someone else? Is consent of receiver required?
3. **Forced transfer** — can a supervisor force a transfer against the current owner's wishes?
4. **Auto-transfer triggers** — SLA miss (W4), CS unavailable, CS off-shift, owner deactivation — do these auto-reassign?
5. **Effect on attribution** — does transfer update `last_touch` attribution? Or does attribution freeze to the first-touch source (Q057 depends)?
6. **Commission implications** — if agent A's lead is transferred to agent B who closes it, who gets commission?
7. **Audit** — what's captured (transfer reason, timestamp, before/after owner)?
8. **Transfer cool-down** — can the same lead be transferred repeatedly, or a limit?

## Options considered

- **Option A — Self-initiated transfer only (consent-based); supervisor-override for deactivation / off-shift.** Agent/CS explicitly transfers to a named receiver who accepts.
  - Pros: respects ownership; clean audit trail.
  - Cons: owner may sit on dead leads; no efficient reassignment.
- **Option B — Supervisor can force transfer anytime with audit; SLA miss auto-reassigns.** Central control.
  - Pros: operational efficiency; dead-lead cleanup.
  - Cons: erodes ownership feel; potential abuse (supervisor favors one CS).
- **Option C — Hybrid: self-transfer (consent) + supervisor force-transfer with reason + auto-transfer on SLA/off-shift.**
  - Pros: balances autonomy with operational control.
  - Cons: more states; UI surface complexity.

## Recommendation

**Option C — hybrid with four transfer types: self-transfer (consent required), supervisor force-transfer (reason captured), auto-transfer (SLA/deactivation), and consent-requested transfer (sent, receiver accepts/declines).**

Option A leaves too many stale leads. Option B erodes agent/CS morale (ownership is motivation). Option C's four-type model covers the operational scenarios without creating a free-for-all.

Defaults to propose: **Self-transfer with consent** — owner proposes transfer to named receiver (agent or CS); receiver accepts/declines in-app; accepted transfer fires. **Supervisor force-transfer** — CS supervisor or marketing admin can transfer any lead with required reason (dropdown + free-text); audit-logged. **Auto-transfer** — SLA breach (W4 per Q071), CS goes off-shift with open leads (transfer to rotation queue), agent deactivated (transfer to CS pool). **Consent-requested transfer** — supervisor initiates a request to transfer between agents; both parties notified; receiver has 24h to accept. **Attribution** — per Q057: first_touch_utm frozen, last_touch_utm updated; commission routing uses last-touch; transfer doesn't retroactively change attribution. **Commission on closed transferred lead** — current owner at `paid_in_full` time earns direct commission; override chain per Q056; original touches affect attribution report but not direct commission. **Audit** — captures: timestamp, from_user, to_user, reason enum (SLA miss / handoff / supervisor / agent deactivation / consent), free-text notes. **Cool-down** — no hard limit on transfers per lead; supervisor alerts if > 3 transfers in 7 days (signal of unclear ownership).

Reversibility: transfer rules are config. Tightening/loosening who-can-transfer is a role change.

## Answer

TBD — awaiting stakeholder input. Deciders: CRM lead, CS supervisor, agency owner (ownership-vs-efficiency balance).
