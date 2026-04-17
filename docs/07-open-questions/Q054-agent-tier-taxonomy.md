---
id: Q054
title: Agent tier taxonomy + qualification thresholds + demotion rules
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: open
---

# Q054 — Agent tier taxonomy + qualification thresholds + demotion rules

## Context

PRD line 169 offers *"misal: Silver → Gold"* as an **illustrative** tier progression — explicitly flagged as an example. The actual tier taxonomy, the qualification thresholds, and the demotion rules are not specified. This question pins them.

Commission math (Q055), override mechanics (Q056), and academy gating (W12) all depend on a concrete tier ladder. The existing data-model placeholder (`agent_level enum ('silver', 'gold', 'platinum')`) uses three tiers as a guess; PRD line 169 only names two (Silver, Gold).

## The question

1. **How many tiers** — 2 (Silver/Gold), 3 (Silver/Gold/Platinum), or more?
2. **Tier names** — Silver/Gold/Platinum, or Indonesian (Mitra Baru / Mitra Utama / Mitra Premier), or something else?
3. **Qualification thresholds** per tier — PRD mentions three dimensions (line 169): education completion, activity, and closings. What are the numbers?
   - Silver → Gold: e.g. 5 closings + Academy Level 1 passed + N days active?
   - Gold → Platinum: e.g. 20 closings/year + Academy Level 2 + leadership of > N downline?
4. **Demotion rules** — if a Gold agent's closings drop below threshold, does the tier demote? Grace period?
5. **Initial tier at onboarding** — all new agents start at the lowest tier, or is there an admin override?
6. **Special categories** — Cabang / Perwakilan (PRD line 171) — is this a tier or a separate role?
7. **Branch-of-tier** — different tiers behave differently in which aspects (commission %, override %, visibility, broadcast access)?

## Options considered

- **Option A — 3 tiers (Silver/Gold/Platinum) + Cabang/Perwakilan as separate role.** All agents start Silver. Promotion thresholds: Gold at 5 closings + Academy L1; Platinum at 20 closings/year + Academy L2 + maintaining activity. Cabang is admin-assigned, not promotion-earned.
  - Pros: simple ladder; clear separation between sales-earned tiers and admin-appointed representation.
  - Cons: thresholds may be off; need real-world calibration.
- **Option B — 2 tiers (Silver/Gold) + Cabang/Perwakilan role.** Matches PRD line 169 literal. Minimal taxonomy.
  - Pros: smallest surface; ships fastest.
  - Cons: no progression beyond Gold limits agent motivation; risk of flat hierarchy.
- **Option C — 4+ tiers with Performance Points accumulator.** Points earned for closings + academy + recruiting + retention; thresholds at Bronze/Silver/Gold/Platinum/Diamond.
  - Pros: gamifies agent engagement; more motivation levers.
  - Cons: complex to maintain; points accounting is its own sub-system.

## Recommendation

**Option A — 3 tiers (Silver/Gold/Platinum) + Cabang/Perwakilan as a distinct admin-appointed role.**

Option B is under-spec'd — two tiers don't provide enough progression for a network that should scale to thousands of agents. Option C is over-engineered for MVP; points-based leveling requires a scoring engine plus monthly reconciliation plus appeals — features a mature program earns into, not a starter. Option A matches the existing data-model stub (three enum values), gives two progressions for motivation, and separates the earn-based ladder from the appoint-based representation role.

Defaults to propose: **3 tiers** (Silver / Gold / Platinum). All new agents start Silver after KYC approval. **Silver → Gold**: 5 closings in any 6-month window + Academy Level 1 (fiqh + sales basics). **Gold → Platinum**: 20 closings in any 12-month window + Academy Level 2 + > 3 active sub-agents under supervision. **Demotion**: 6-month grace period — if a Gold agent has zero closings for 6 months, demote to Silver (Platinum → Gold on same rule). Demotion paused if agent's downline is active (Gold/Platinum agents recruiting sub-agents still count as "active" even without direct closings). **Cabang / Perwakilan**: admin-appointed regional representative; highest override rate; sees super-view of all downline in their region; not a promotable tier — appointed for business reasons. **Tier affects** commission % (Q055), override % (Q056), broadcast segment membership, leaderboard visibility, and academy unlock levels.

Reversibility: tier count is schema (enum extension) + config. Adding a tier is additive; renaming tiers is migration. Threshold numbers are config, tunable anytime.

## Answer

TBD — awaiting stakeholder input. Deciders: agency owner (agent-program strategy), CRM lead, existing top-agent feedback if retained.
