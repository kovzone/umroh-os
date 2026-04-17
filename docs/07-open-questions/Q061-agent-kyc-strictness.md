---
id: Q061
title: Agent KYC strictness + activation thresholds
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: open
---

# Q061 — Agent KYC strictness + activation thresholds

## Context

Agent onboarding (modules #25, #26) requires KTP, NPWP, photo upload and admin verification. PRD doesn't specify how strict the verification is, what's required vs optional, and whether there are activation gates beyond admin approval (e.g. minimum training completion, first-booking quota).

KYC strictness matters for:
- **Legal posture** — agents earning commission are effectively in a business relationship; stricter KYC protects the agency (tax, compliance).
- **Tax reporting** — NPWP required for PPh 21/23 withholding at normal rates; no NPWP means 2× rate per Q047.
- **Fraud prevention** — stricter KYC reduces sock-puppet / friend-ring abuse.
- **Onboarding friction** — too strict and agents drop off; too loose and bad actors slip through.

## The question

1. **Required documents** — KTP always required. Is NPWP always required, or only above commission threshold?
2. **Bank account verification** — is name-match between KTP and bank account required? How is it verified (manual admin review, third-party e-KYC service)?
3. **Face-match / selfie liveness** — is KTP photo compared against selfie for identity match? Automated (e-KYC provider) or manual?
4. **Activation gates beyond admin approval** — must the agent complete Academy Level 1 before earning commission? First-training-session before seeing downline view?
5. **Initial tier at activation** — Silver default (per Q054), or provisional "Pending" status until first closing?
6. **Age / legal capacity** — minimum age for agent (18? 17 with guardian consent)?
7. **Re-verification cadence** — does KYC re-verify periodically (annual) or once-and-done?
8. **Suspicious pattern rejection** — automatic rejection criteria (known-bad phone, duplicated KTP, mismatched name on bank)?

## Options considered

- **Option A — Minimal KYC: KTP + photo + phone; NPWP only when commission accrued > 4.8M IDR/year (PTKP threshold).** Lowest friction.
  - Pros: fast onboarding; more agents activated.
  - Cons: weak anti-fraud; mixed NPWP state across agents complicates tax; manual tracking of the PTKP threshold.
- **Option B — Strict KYC: KTP + NPWP + bank + selfie match + Academy L1 before activation.** Heavy upfront.
  - Pros: better quality agents; cleaner tax handling; reduced fraud.
  - Cons: high drop-off at onboarding; slower network growth.
- **Option C — Tiered KYC: minimal to onboard + incremental unlocks.** Start with KTP + photo → can browse + share links. Add NPWP + bank to earn commission. Add Academy L1 to promote to Gold.
  - Pros: reduces friction gate; progressive unlock aligns requirements with value.
  - Cons: two states (onboarded-but-can't-earn vs fully-active); UI complexity.

## Recommendation

**Option C — tiered KYC with progressive unlocks; full-PPh-compliant activation requires KTP + NPWP + bank verified.**

Option A's loose KYC creates tax headaches (some agents with NPWP, some without, PPh withholding math inconsistent) and accepts anti-fraud risk in exchange for agent count. Option B's strict upfront requirements are correct but create onboarding drop-off that kills the viral part of the agent program. Option C threads the needle: low-friction signup to browse and share, but earning commission requires completing the full KYC pack — motivates agents to finish KYC quickly once they see the value.

Defaults to propose: **Signup** — KTP + selfie + phone + email; agent can browse catalog, share replica site, accumulate leads. Status = `onboarded_partial`. **Commission earning** — requires NPWP (or declaration of no-NPWP with understanding of 2× PPh withholding per Q047) + bank account + signed E-MoU. Admin verifies bank-account name-match manually against KTP (Phase 2: e-KYC provider automation). Status = `active`. **Tier promotion to Gold** — requires Academy Level 1 completion (per Q054). **Face-match / liveness** — MVP skips automated liveness check; manual admin comparison KTP-photo vs selfie. Phase 2: integrate e-KYC provider (Privy, Asli RI). **Age** — 17+ with KTP (normal Indonesian ID age); under-17 not eligible. **Re-verification** — once-and-done for MVP; annual re-verify triggered only if Bukti Potong generation fails (stale NPWP data). **Suspicious patterns** — duplicate KTP number auto-rejects; duplicate phone flags for admin review; duplicate bank account flags for admin review.

Reversibility: raising/lowering KYC strictness is config per agent-tier. Adding liveness check is additive (Phase 2 integration). Retroactively tightening KYC on existing agents is policy decision (grandfather existing, enforce new on new signups).

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (growth-vs-quality balance), CRM lead, finance director (tax compliance posture), legal (UU PDP considerations per Q008).
