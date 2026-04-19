---
id: Q013
title: Dual-gateway (Midtrans / Xendit) selection and fallback rule
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F5
status: answered
---

# Q013 — Dual-gateway selection and fallback

## Context

PRD line 1263 lists both Midtrans and Xendit as supported gateways ("contoh: Midtrans/Xendit"). F5 exposes an adapter per gateway so adding / switching is routine. But the PRD doesn't say how the system chooses between them when both are configured.

Three realistic modes:

1. **Single primary** — only one gateway is active at a time; admin switches via config. The "Xendit support" is latent until someone flips it.
2. **Active/passive failover** — primary handles everything; secondary kicks in on primary failure (timeout or 5xx).
3. **Load split** — transactions split across gateways by percentage or round-robin, typically for cost optimisation or availability.

## The question

1. **Which mode** is the platform default?
2. If failover: on what failure conditions does it trigger (timeout only? 5xx? 4xx from Midtrans too)?
3. Should the gateway choice be **visible to the customer** (they pick bank from a drop-down, and the gateway underneath is invisible) or **transparent** (the customer just sees "pay via VA" with one choice)?
4. Gateway fees — does the agency prefer Midtrans or Xendit for cost reasons? This informs which is primary.

## Options considered

- **Option A — Single primary (Midtrans), Xendit disabled.** One gateway live at a time. Switch by changing config.
  - Pros: simplest. No runtime fallback logic. Easiest to debug.
  - Cons: Midtrans outage → no VA issuance → no new bookings. Hard dependency.
- **Option B — Active/passive failover (Midtrans primary, Xendit fallback).** Try Midtrans first; on transient error, retry on Xendit. Webhooks from both are processed.
  - Pros: resilience to single-gateway outage; customers never see a "payment unavailable" page.
  - Cons: operational complexity (managing two gateway accounts, two sets of keys, two webhook handlers); settlement reconciliation across two providers; slightly higher cost if Xendit fees differ.
- **Option C — Customer chooses bank, system picks gateway.** Customer sees a bank list (BCA, Mandiri, BNI, ...); the system routes to whichever gateway has the best rate / availability for that specific bank. Fully transparent to the customer.
  - Pros: optimises fees; best customer experience per-bank.
  - Cons: routing logic complexity; requires a per-bank gateway-fee table; non-trivial to test.

## Recommendation

**Option B — Midtrans primary, Xendit fallback.** Start with Midtrans as primary (it has deeper Indonesian bank coverage and is the more mature product in the market). Wire Xendit as a hot standby — keys in config, adapter ready, but only activated on Midtrans failure (5xx or timeout > 10s). Don't fall back on 4xx — those are genuine rejections (bad VA config, invalid amount) that would fail on Xendit too.

Gateway choice is **transparent** to the customer — they see one "Pay via Virtual Account" button; the system picks the gateway. Hiding it avoids confusing the customer with unnecessary choice.

Reasoning: Option A's single-gateway dependency is too brittle — a Midtrans incident (they do have them every few months) would mean zero new bookings for the duration, which on a popular departure is real revenue lost; Option C's per-bank routing is over-optimisation for a two-dev team to maintain and the bank-routing table would rot without dedicated attention. Option B hits the sweet spot: resilience without runtime complexity.

Settlement reconciliation: finance dashboard shows two tracks (Midtrans settlements + Xendit settlements) with an aggregate view. Periodic finance work to match each gateway's settlement report to the internal `payment_events` — already part of the reconciliation cron (F5 W5).

Reversibility: changing primary from Midtrans to Xendit is a config flip. Disabling fallback is a config flip. Expanding to Option C later is an adapter-layer refactor, not a data migration.

## Answer

**Decided:** **Option B** — **Midtrans primary**, **Xendit hot-standby**; failover on **5xx or timeout > 10s**; **no failover on 4xx** (treat as config/validation failure). **Customer-transparent** single pay UX. **Settlement:** finance reconciles per provider statements against `payment_events` (existing F5 pattern).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
