---
id: Q015
title: Smart Grouping trigger timing and override authority
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F4, F7
status: answered
---

# Q015 — Smart Grouping trigger timing + override authority

## Context

Modules #92 (Algoritma Penempatan Kamar) and #93 (Pengatur Transportasi) produce room / bus / plane-seat assignments. PRD Section E locates them under Operational & Handling (not Booking) and Alur Logika 6.4 treats them as pre-departure ops activity. The algorithm uses K-Family Code, domicile, and preferences (PRD line 319).

Two decisions the PRD doesn't explicitly pin:

1. **When does the algorithm run?** — at every booking submit? At lunas? On a nightly batch? On-demand by ops close to departure?
2. **Who can override an assignment?** — ops admin only? Tour leader / Muthawwif in Saudi? The jamaah themselves?

Both affect operational UX and data model (re-running the algorithm on a live departure disturbs already-known assignments).

## The question

1. **Trigger timing.** Four candidates:
   a. At booking submit — every new booking triggers a re-run for the whole departure.
   b. At lunas (paid_in_full) — only fully-paid jamaah participate in grouping.
   c. Nightly batch — the algorithm runs once per day for upcoming departures.
   d. On-demand — ops admin clicks "Run Grouping" when ready; algorithm runs then.
2. **Override authority.** Who can manually edit an auto-generated assignment?
   a. Ops admin only
   b. Ops admin + Tour Leader (in-field changes mid-trip)
   c. Ops admin + Tour Leader + Muthawwif (including Saudi-side room assignments per PRD line 365)
   d. Any authorised user (including jamaah self-service "request my preferred roommate")

## Options considered

- **Option A — On-demand + ops/tour leader override.** Ops triggers the algorithm on-demand (1d) when the pax list is stable; ops admin and tour leader can override (2b). Muthawwif edits go through the tour leader for centralised logging.
  - Pros: deterministic (no surprise reshuffles); matches real workflow (ops does a review pass when docs are ready).
  - Cons: requires ops discipline — forgetting to run before departure is a real risk.
- **Option B — Nightly batch with override.** Nightly job runs for departures in the next 60 days; ops and tour leader can override (2b).
  - Pros: auto-rerun catches late-added bookings; no ops workflow dependency.
  - Cons: assignments may shift nightly as new bookings come in, confusing jamaah who were told their room number 3 days ago.
- **Option C — At-lunas + ops-only override.** Algorithm runs when a booking hits paid_in_full (1b); ops is the only override authority (2a); tour leader request changes via ops.
  - Pros: cleaner authorization model; assignment stability increases as more jamaah reach lunas.
  - Cons: draft/pending bookings ignored in grouping until payment — potentially inefficient if most bookings reach lunas only close to departure.

## Recommendation

**Option A — On-demand trigger + ops/tour-leader override.**

Ops triggers the algorithm manually when:
- All (or most) jamaah are paid_in_full
- Docs are verified (passport OCR done, mahram resolved)
- Departure is close enough that the pax list is stable (≥ H-14 is a reasonable heuristic)

Override authority:
- **Ops admin** can edit any assignment before departure and logs the reason
- **Tour leader** can edit assignments on the day of travel or during the trip (e.g. family bed arrangement issue in Mecca) — also logs reason
- **Muthawwif** in Saudi reports room-swap needs to the tour leader who commits the change on the system (single change log via tour leader)
- **Jamaah self-service** is NOT allowed at launch; add later if real demand surfaces

Re-running the algorithm after initial run: **off by default**, but ops can explicitly click "Re-run" and review the diff before committing (surface changes to existing assignments so ops sees "Jamaah X moves from Room 302 to Room 305" and decides whether to accept).

Partial-cancel handling (see Q014): leaves a hole rather than auto-regrouping — ops decides whether to fill the hole manually or rerun.

Reasoning: on-demand matches how ops actually works (reviewing a departure when it's ready, not running algorithms nightly on half-baked data); ops-or-tour-leader-only keeps the authorization model simple while allowing the in-field reality (hotels do swap rooms).

Reversibility: trigger timing is a process decision, not data-shape — can shift from on-demand to nightly later without migration. Override authority is a permission matrix in F1 — adjustable. Low-to-medium commitment.

## Answer

**Decided:** **Option A** — **on-demand** grouping run by ops when pax list stable (heuristic **≥ H-14** + most jamaah **paid_in_full** + docs verified); **no nightly auto-rerun** in MVP. **Overrides:** **ops admin** anytime; **tour leader** in-field; **muthawwif** routes room-critical edits **through tour leader account** for audit (same net effect as Recommendation table). **Re-run:** allowed only via explicit **“Re-run with diff review”** commit. **Jamaah self-service roommate requests:** **deferred**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
