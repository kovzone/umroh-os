---
id: Q033
title: Courier integration policy — single vs multi, fallback, routing
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8
status: answered
---

# Q033 — Courier integration policy

## Context

PRD module #123 *Integrasi Ekspedisi* (line 413) says fulfillment uses a courier API to print labels and push tracking numbers to jamaah WhatsApp — but doesn't name the courier(s), the selection rule, or the fallback when a courier's API is unavailable. F8 W11 (shipment) can't ship without this.

Indonesian ERP shipments typically go via one of JNE / J&T / SiCepat / AnterAja / Pos Indonesia. Each has different coverage, pricing, and reliability characteristics per region. Agencies today typically have 1–2 preferred couriers and pick based on destination + package type.

## The question

1. **Single courier or multi?** — One locked-in courier per agency, or multi-courier with routing logic?
2. **Routing rules** — if multi, what drives selection? Destination region, package weight, jamaah tier (B2C vs B2B agent shipment), cost, SLA?
3. **Fallback chain** — when primary courier API is down, what's the cascade? Queue and retry? Pick alternative automatically? Alert ops for manual pick?
4. **Courier API down for > N minutes** — do dispatch tasks queue, or do we fail fast and page ops?
5. **Courier selection per dispatch** — is this an ops decision at each dispatch, system-automatic, or jamaah-selectable at booking time?
6. **Inter-island / remote area handling** — some couriers don't cover all of Indonesia; does the system refuse to ship to non-covered addresses, or surface a warning?
7. **Agent pickup vs jamaah home delivery** — is an agent always the courier destination for B2B bookings, with the agent handling last-mile?

## Options considered

- **Option A — Single primary courier, one integration.** Pick one national courier (e.g. JNE), integrate its API, done. Manual workaround via phone for outages.
  - Pros: simplest integration; fewer adapters to maintain.
  - Cons: single point of failure; no cost optimization; coverage gaps.
- **Option B — Multi-courier with ops-selected-per-dispatch.** Integrate 2–3 couriers. Ops picks at dispatch time per their judgment (cost, coverage, recent performance).
  - Pros: flexibility; operator can route around outages.
  - Cons: manual per-dispatch decision; inconsistent outcomes.
- **Option C — Multi-courier with automatic routing + fallback chain.** System picks primary based on destination zone + weight; falls back to secondary on API error; alerts ops only when all fail.
  - Pros: fastest fulfillment; resilient; minimal manual overhead.
  - Cons: routing logic is non-trivial; adapter code per courier.

## Recommendation

**Option C — multi-courier with automatic routing (destination-zone primary) + deterministic fallback chain.**

Option A's simplicity isn't worth the outage exposure: a single-courier integration paralyzes fulfillment every time the provider has an AWS incident. Option B puts operational judgment calls on the warehouse team every day, which is wasted friction for 95% of deliveries. Option C front-loads the routing decision into config and removes it from the daily flow — the adapter layer is 3 small courier wrappers behind one `CourierAdapter` interface, the routing is a table (zone → primary / secondary / tertiary).

Defaults to propose: JNE primary for Java + major cities; J&T fallback; SiCepat for remote areas where JNE coverage is thin. Fallback cascade triggers on HTTP 5xx or > 30s timeout from primary; all-fail alerts ops with manual-resolve option. Weights > 10kg forced to JNE Cargo. B2B shipments default to the agent's registered address, not end-jamaah. No jamaah-selectable courier at booking time in MVP (adds too much complexity for marginal gain; revisit if jamaah complain).

Reversibility: courier selection is config. Adding a new courier is a new adapter + routing table row. Switching primary is a config change.

## Answer

**Decided:** **Option C** — **zone → primary/secondary/tertiary routing table**; **fallback on 5xx or >30s timeout**; **queue with admin alert if all fail**; **>10kg → cargo profile**; **B2B dispatch defaults to agent branch address**; **no jamaah courier pick** in MVP. **Defaults:** JNE primary Java/metro, J&T secondary, SiCepat tertiary remote (editable).

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
