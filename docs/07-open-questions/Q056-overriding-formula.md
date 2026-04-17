---
id: Q056
title: Overriding commission formula + hierarchy depth + orphaned-upline handling
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: open
---

# Q056 — Overriding commission formula

## Context

PRD line 171 describes overriding as *"kalkulator pembagian komisi otomatis antara agen yang closing dengan perwakilan di atasnya (selisih jaringan)"* — the upline earns the **difference** (selisih) between their tier % and the closing agent's tier %. This is the mechanism that makes the network model work — agents recruit and mentor because they earn off downline sales.

But the formula is under-specified:

- Does it compound across multiple levels (Silver → Super-Agen → Cabang → beyond?)
- What happens if an intermediate upline is inactive?
- Is hierarchy depth bounded?
- Does override apply to promotional campaign uplifts (Q055)?

## The question

1. **Formula** — override % = upline's tier % − downline's tier %? Or something else (e.g. fixed uplift per level)?
2. **Compounding across levels** — if a closing Silver's immediate upline is Gold and their upline is Platinum, does Platinum earn (Platinum − Gold) or (Platinum − Silver)?
3. **Hierarchy depth** — max number of levels that can earn override from a single closing. 2 levels? Unlimited?
4. **Orphaned-upline handling** — if an intermediate upline (say Super-Agen) is suspended or deactivated between closing agent and Cabang, does override:
   - Skip the deactivated level (Cabang earns (Cabang − Silver))?
   - Stop at deactivated level (Cabang earns nothing)?
   - Hold in escrow until deactivated agent reactivates?
5. **Override on promotional uplift** — does the uplifted base % apply to override calculation?
6. **Cabang / Perwakilan role** — is it the top of the hierarchy (can't be recruited by another agent) or can it be nested?
7. **Cross-tree closings** — if an agent switches upline mid-season, do past-closing overrides stay with old upline or transfer to new?

## Options considered

- **Option A — Standard selisih jaringan, compound up, skip inactive.** Override per level = (this level % − immediately downstream level %). Compounds up through all active levels. Inactive levels skipped; next active level earns (this level % − last-active-downstream level %).
  - Pros: matches literal "selisih" reading; compensates actives for managing inactive uplines; simplest math per step.
  - Cons: deactivated agents lose all earnings (including pre-deactivation rightful ones, if handled per-event).
- **Option B — Fixed uplift per level.** Each upline earns a fixed X% (say 1%) regardless of tier math. Simpler.
  - Pros: trivial math; easy to explain.
  - Cons: doesn't reward tier advancement; under-motivates upper tiers.
- **Option C — Only immediate upline earns override.** Single-level only. No compounding up.
  - Pros: simplest; limits override cost.
  - Cons: disincentivizes Super-Agen or Cabang from investing in building a deep network.

## Recommendation

**Option A — standard selisih jaringan, compounds up through all active levels (max 3 levels for hierarchy depth cap), skip inactive uplines.**

Option B flattens the motivation structure — a Platinum agent deep in a recruiter chain should earn more than the 1% for every downline closing because they invested in recruiting + mentoring. Option C limits the network depth's value, capping the payoff of investing in recruiting — Super-Agens will stop recruiting above their immediate reports. Option A gets the economics right at the cost of some complexity.

Defaults to propose: override % per level = (this level's base tier %) − (immediately downstream level's base tier %). Hierarchy cap = **3 levels** (closing agent + 2 uplines max earn — e.g. Silver → Gold Super-Agen → Cabang). Inactive handling: skip; next active upline earns (this level % − last-active-downstream level %). Promotional uplift applies to both closing and override calculations (uplifted base propagates up the chain). Cabang is the top of the explicit hierarchy; no Cabang-of-Cabang. Cross-tree transfer: agent upline change is rare and triggers an audit log; past closings stay with the historical upline at time of closing (snapshot); future closings go to new upline. Example calculation for a 10M IDR Umroh closing with Standard category (Q055 defaults):

- Closing Silver agent: 3% direct = 300K
- Gold Super-Agen upline: (5% − 3%) = 2% override = 200K
- Cabang/Perwakilan: (10% presumed Cabang rate per Q055 − 5%) = 5% override = 500K
- Total commission paid out: 1M IDR (10% of 10M revenue)

Reversibility: override rules are config. Changing hierarchy depth cap (3 → unlimited, or 3 → 2) is config, not schema. Historical commissions frozen at their calculation time.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (commission-budget implication), CRM lead (agent-program experience), finance director (payout cost impact).
