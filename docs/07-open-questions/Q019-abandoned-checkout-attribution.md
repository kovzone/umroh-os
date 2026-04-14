---
id: Q019
title: Abandoned checkout attribution — if CS rescues a B2C-abandoned cart, who gets commission?
asked_by: session 2026-04-15 F4/F5 draft
asked_date: 2026-04-15
blocks: F4, F10
status: open
---

# Q019 — Abandoned checkout attribution

## Context

PRD Module #57 Pemicu Momen (line 201) mentions "abandoned interest" triggers — when a customer drops out of checkout without paying, the system can re-engage them. Module #70 Radar Prospek Lama (line 239) tracks old leads who return. But neither says what happens to **commission attribution** if a salesperson rescues the sale that someone else (or the B2C self-service funnel) started.

Three realistic rescue scenarios:

1. B2C self-checkout → abandoned → CS follows up via WhatsApp → CS closes. Does commission go to CS, or to the "referring channel" (which was... nothing? direct traffic?).
2. Agent's replicated site → abandoned → CS rescues. Does the agent lose commission? Split?
3. Campaign ad → abandoned → second campaign → closed. Which campaign gets attribution for ROAS?

These decisions affect agent retention, CS motivation, and marketing analytics. Getting them wrong produces visible friction (agents complaining about lost commissions).

## The question

1. **Attribution model.** First-touch (whoever first brought the lead wins), last-touch (whoever closed wins), shared split, or time-bounded (last touch within N days wins, else first)?
2. **CS rescue of agent's lead.** Does CS get the commission? Split with the agent? Default to agent (since they referred)?
3. **Campaign → direct-close attribution.** If a jamaah saw a Meta ad, didn't buy, then later walked in off a WhatsApp broadcast — which channel gets ROAS credit?
4. **Grace window.** UTM cookies typically last 30 days. Is that our window, or different?

## Options considered

- **Option A — last-touch wins with 30-day referral window.** Whoever most recently interacted with the customer closes the attribution. If a customer's last-touch was the agent's replicated site within 30 days, agent gets commission even if CS was the one who hit Send. Attribution captures the cookie/UTM at the time of booking.
  - Pros: matches marketing industry norm; respects agent referral work (they got the customer to the door); CS is salaried, not commission-based, so their motivation isn't tied to this.
  - Cons: unfair if CS does heavy lifting on a stale referral — a bad agent who doesn't follow up still gets paid.
- **Option B — first-touch wins (within 30-day window).** The channel that first brought the customer owns the attribution for 30 days regardless of who closes. Protects agent referrals from CS poaching.
  - Pros: agent protection; clear rule.
  - Cons: a lead that went cold for 29 days and was actively rescued by a different channel still gets attributed to the original; feels unfair to the rescuer.
- **Option C — shared split when both channels touch.** If both an agent and CS touched the lead within 30 days, split commission 50/50 (or some config'd ratio).
  - Pros: operationally fair; both get something.
  - Cons: complex; commission math and reporting gets messy; split percentages invite endless negotiation.
- **Option D — no attribution; CS is salaried, agents only get commission on self-closed-via-replica deals.** Sharp boundary: agent gets commission ONLY if the booking is made through the agent's replicated site without CS involvement.
  - Pros: simplest rule.
  - Cons: discourages agents from generating leads that CS might need to close (e.g. callers with complex questions); CS and agent might even compete for the same customer.

## Recommendation

**Option A — last-touch wins with 30-day referral window.** When a booking is created, F10 computes attribution as:

1. If the booking carries a `ref=<agent_code>` in the URL (agent replicated site), and this ref was captured within the last 30 days, agent gets commission.
2. Else if there's a `utm_source=<campaign>` captured within 30 days, campaign gets ROAS attribution (marketing report), no commission (campaigns don't pay commission).
3. Else no attribution (direct booking or CS-closed without a referral trail).
4. CS is always responsible for the booking regardless — their dashboard reflects the work (Module #67 Dashboard Kinerja CS), but compensation is salaried, not commission-driven.

30-day window is the industry standard for UTM tracking and matches Meta / Google Ads defaults; implementing anything non-standard is friction for no gain.

Handling the "stale agent referral" complaint: Option A's 30-day window already caps exposure. If 29 days pass without the agent following up and CS rescues, the agent still gets paid (per rule) — but the agency can adjust the agent's LEVEL (Module #46 Leveling Otomatis) if they consistently leave referrals to rot.

Reasoning: agent commissions are a retention mechanism — protecting them (even on arguably-rescued sales) preserves the agent network's incentive to generate leads; 30 days is familiar to agents and campaigns alike; CS stays salaried so they don't need this.

Reversibility: attribution rules are a service-layer policy; changing the window from 30 days to 14 or 60 is a config value. Moving from last-touch to first-touch later is a logic change in F10's attribution service, no data migration. Low commitment.

## Answer

TBD — awaiting stakeholder input. Likely decider: agency owner + CS lead + representative agent (sample input from the agent network).
