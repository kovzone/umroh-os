---
id: Q010
title: VA TTL and draft booking expiry — how long does a customer have to pay before the seat returns?
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F4, F5
status: answered
---

# Q010 — VA TTL and draft booking expiry

## Context

When a booking is submitted and a virtual account is issued (F5 W1), the customer has some window to transfer the money. After that window, the VA expires and the reserved seat needs to return to the available pool. PRD line 535 confirms `batas waktu tunggu pelunasan invoice` is a **global config variable** — i.e. the platform expects a configurable TTL — but doesn't specify the default value.

Two related TTLs to decide:

1. **VA TTL** — how long a customer has from VA issuance to make the first payment (DP or full).
2. **Draft booking TTL** — how long a booking can sit in `draft` (form half-filled, not yet submitted) before it auto-expires and the reservation attempt is treated as abandoned.

## The question

1. **VA TTL for DP** — industry-standard bank VA windows are 24h, 48h, or 72h. Which default? Should it vary by package kind (Haji Furoda being higher-value, maybe longer grace)?
2. **Draft TTL** — how long between "user opens checkout" and "user abandons" before the tentative seat hold is released? Candidates: 15 min, 30 min, 1 hour.
3. **Should both be admin-configurable per package** or fixed platform-wide defaults with global override only?
4. **Pelunasan deadline** — separate from VA TTL (DP), how long does a `partially_paid` booking have before the remainder is due? PRD line 85 mentions "tenggat waktu pelunasan" but gives no default. Default candidates: H-30, H-45, H-60 before departure.

## Options considered

- **Option A — conservative, industry-standard defaults.** VA TTL 24h for DP; draft TTL 30 min; pelunasan deadline H-30 before departure. All admin-configurable globally but **not** per-package (simplicity).
  - Pros: matches what most Indonesian travel sites do; jamaah already expect the pattern.
  - Cons: some premium programs (Haji Furoda) might want longer; fixed defaults force a workaround (extend manually per case).

- **Option B — generous, matches premium pilgrimage behavior.** VA TTL 48h; draft TTL 1 hour; pelunasan deadline H-45. All admin-configurable globally AND per-package kind.
  - Pros: friendlier to high-value customers; premium packages get natural breathing room.
  - Cons: slower inventory turnover on popular departures; per-package config adds UI complexity.

- **Option C — tiered by package kind.** Standard Umroh gets Option A; Haji Furoda / Umroh Plus gets Option B. Draft TTL always 30 min regardless of kind.
  - Pros: business-logic fit (longer window where stakes are higher); draft TTL stays simple.
  - Cons: premium-tier logic leaks into the payment layer.

## Recommendation

**Option A** with all three TTLs exposed as global config variables editable by Super Admin via the Pengaturan panel (PRD line 535 already treats this as a global var).

- VA TTL default: **24 hours**
- Draft TTL default: **30 minutes**
- Pelunasan deadline default: **H-30** before departure

Reasoning: these are sensible industry defaults that match what jamaah already expect from other Indonesian booking platforms; starting simple (platform-wide, not per-package) cuts UI and logic surface area now; if Haji Furoda customers push back later, per-package override can be added without data migration (just a new `package_departure.va_ttl_override_hours` column). Starting lenient (Option B) risks slow inventory turnover on popular departures without concrete customer pressure saying the 24h window is too tight.

Reversibility: changing any of the three values is a config flip, no schema change. Adding per-package override is an additive column, cheap. Low commitment.

## Answer

**Decided:** **Option A** — **VA TTL 24h**, **draft TTL 30m**, **pelunasan H-30**; all **global Super Admin config** for MVP (no per-package override until needed). **Rationale:** matches common Indonesian e-commerce expectations; per-package override remains an additive column later.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
