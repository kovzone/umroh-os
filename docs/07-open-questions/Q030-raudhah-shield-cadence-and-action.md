---
id: Q030
title: Raudhah Shield polling cadence and alert action
asked_by: session 2026-04-15 F6 draft
asked_date: 2026-04-15
blocks: F6
status: open
---

# Q030 — Raudhah Shield cadence + alert action

## Context

Module #105 Raudhah Shield (PRD lines 25, 359) is UmrohOS's anti-fraud feature: it polls Nusuk to detect whether a jamaah's visa is being used by someone else during the Saudi window ("pembajakan visa"). The PRD describes **what** it does but not **how often** or **what happens** on detection.

Two decisions to make:

1. **Cadence** — how often to poll Nusuk. Too often = Nusuk rate-limits or blocks us; too rarely = fraud window opens.
2. **Alert action** — on detected mismatch, what does the system actually do? Notify whom, block what, escalate where?

The answer intersects with Q024 (incident escalation) and F7 W11 (tasreh scan at Raudhah entry, which could be gated by Raudhah Shield status).

## The question

1. **Nusuk API access.** Does the agency have Nusuk API access? Auth model, rate limits, endpoint documentation? (Parallel to Q026 for MOFA/Sajil — may be same answer path.)

2. **Polling cadence during the Saudi window.**
   - Every hour? Every 6 hours? Every 12 hours?
   - Does it vary — higher frequency around Raudhah entry appointments, lower during rest periods?
   - Any per-jamaah tuning (VIP jamaah gets more frequent polls)?

3. **What constitutes a "detected mismatch"?**
   - Visa status shows as used in Nusuk when no tasreh entry scan (F7 W11) has been recorded by our system?
   - Visa linked to a different passport number than expected?
   - Unexpected travel pattern (e.g., visa activated in a city we didn't book)?

4. **On detection, who gets notified?**
   - Ops pusat (central operations)
   - Tour leader in Saudi
   - Muthawwif assigned to the jamaah
   - Agency owner (if severity is high)
   - The jamaah themselves (may panic them; depends on policy)

5. **What gets blocked?**
   - Nothing; just an alert for investigation
   - Future tasreh scans (F7 W11) for the flagged jamaah require ops approval
   - Booking + refund triggered automatically
   - The Nusuk-linked visa itself (we can't actually block it from the Saudi side)

6. **False positive handling.**
   - Clear process for dismissing false positives?
   - Audit trail of dismissals to detect patterns?

## Options considered

- **Option A — Conservative: 6h cadence, notify ops + tour leader, no auto-blocking.** System detects and alerts; humans investigate.
  - Pros: balances Nusuk rate limits with coverage; notification-only means false positives don't stall legitimate pilgrims.
  - Cons: fraud window is 6h wide; a legitimate fraud could complete a Raudhah entry before alerted.
- **Option B — Aggressive: 1h cadence, notify ops + tour leader + muthawwif + optional jamaah, gate F7 W11 scan until cleared.**
  - Pros: narrowest fraud window; active protection.
  - Cons: Nusuk rate limits may not permit 1h polling; false positive gating tasreh scan is operationally disruptive.
- **Option C — Adaptive cadence: 12h normal, 1h around known Raudhah appointment windows.** Requires knowing tasreh appointments (possible via module #167 Dompet Dokumen Digital).
  - Pros: high coverage when it matters; low load otherwise.
  - Cons: more complex scheduling logic; depends on tasreh data being captured.

## Recommendation

**Option C — adaptive cadence + notify ops and tour leader + muthawwif + audit trail for dismissals.**

Cadence defaults:
- **12 hours** during the Saudi window outside known Raudhah appointments.
- **1 hour** for the 6-hour window surrounding each scheduled Raudhah tasreh entry time (visible from tasreh records).
- **One-shot immediate poll** triggered by F7 W11 tasreh scan event (cross-verify at scan time).

Detection criteria (_(Inferred)_, to be validated against real Nusuk response shape):
- Visa status shows "used" in Nusuk with a timestamp that doesn't match a tasreh scan event in our system → **flag** (most likely real fraud signal).
- Visa linked to a different passport number than on file → **flag** (configuration or fraud).
- Status in Nusuk diverges significantly from our expected state (e.g., "cancelled" when we expect "active") → **flag**.

Alert action:
- **Immediate notification** within 5 minutes of detection to: ops pusat (WA + in-app), tour leader (WA + in-app push), muthawwif (in-app push with alert sound).
- **Open a `security` category incident** (per Q024) with high severity.
- **Do NOT notify jamaah automatically** — false positives would cause unnecessary distress; ops decides after investigation.
- **Do NOT auto-block the F7 W11 tasreh scan** — legitimate pilgrims deserve entry even if our system is wrong. But surface the Raudhah Shield alert at scan time so the muthawwif can verify identity manually.

Dismissal workflow:
- Ops reviews within 30 minutes; can mark `resolved: false_positive` with required reason.
- Dismissal audit feeds a pattern-detection daily report (repeated false positives for the same jamaah or muthawwif might indicate system config issues).

Nusuk-access dependency: like Q026 for MOFA/Sajil, Nusuk API access may not be trivially available. **Phase 1 fallback**: Raudhah Shield runs in manual mode (ops manually checks Nusuk portal for a subset of high-value bookings on a schedule); automated polling enabled once API access is secured.

Reasoning: adaptive cadence concentrates coverage where fraud is most likely (around the actual Raudhah entry); notify-don't-block preserves operational flow; ops-gated dismissal creates accountability; no-auto-jamaah-notification avoids customer trauma from false positives.

Reversibility: cadence is config; alert routing is config; adding auto-block later is a feature flag. Low commitment.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner + ops lead + legal advisor (Saudi regulatory intersect on Nusuk use + UU PDP on jamaah location data).
