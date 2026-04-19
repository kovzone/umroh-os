---
id: Q074
title: Field radar transport + GPS source for bus tracking
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11, F7
status: answered
---

# Q074 — Field radar transport + GPS source

## Context

F11 W4 / F7 Real-Time Field View has four modules (#181–#184): Radar Bus (#181), Status Raudhah (#182), Pelacakan Koper (#183), Laporan Insiden (#184). Two of these (#181 bus, #184 incidents) are genuinely real-time — jamaah-safety and operational-awareness dependent. Two (#182 Raudhah, #183 koper) can tolerate 1–5 min freshness.

PRD line 603 mentions "instant push" for incidents (the only streaming hint in the whole doc). Everything else is aspirational "real-time" without transport specificity.

Also: **where does bus GPS come from?** Tour-leader's phone app polling location? Dedicated bus-mounted GPS tracker? Third-party GPS service? Each has different cost / reliability / battery tradeoffs.

## The question

1. **Transport for streaming dashboards** — Websocket, Server-Sent Events (SSE), long-poll, or HTTP polling with short cadence?
2. **Streaming scope** — which modules are streaming (#181 + #184), which are polling (#182 + #183)?
3. **Bus GPS source** — tour-leader app GPS (phone), dedicated bus-mounted tracker (hardware), third-party fleet-management API (Traccar, etc.), or combination?
4. **GPS update cadence** — 30s? 1min? 5min?
5. **Location retention** — how long is bus-location history kept (1 week? 1 month? trip duration)?
6. **Offline / dead-zone handling** — bus in Saudi tunnel or remote area: last-known-position visible with age badge?
7. **Privacy** — jamaah-level location tracking vs bus-level only?
8. **Incident push channel** — in-app notification only, or also WA / email to subscribed recipients?

## Options considered

- **Option A — Websocket for streaming (#181 + #184); HTTP polling (5min) for #182 + #183; bus GPS from tour-leader app.** Classic hybrid.
  - Pros: websocket handles the true-streaming needs; polling handles the slower ones; tour-leader app already has connectivity + logins (shared device).
  - Cons: tour-leader phone battery + app-uptime dependency; no redundancy if tour-leader phone off.
- **Option B — All polling (30s minimum); no websocket; tour-leader app GPS.** Simpler infra.
  - Pros: no websocket infra; simpler backend.
  - Cons: 30s latency on incidents is unacceptable (jamaah safety); polling storm on many connected dashboards.
- **Option C — Hybrid (Option A) + dedicated bus-mounted GPS tracker for redundancy.** Hardware + app both push.
  - Pros: redundant GPS; tracker works even when tour-leader phone off.
  - Cons: hardware cost (per bus), installation logistics, partner management (which GPS platform).

## Recommendation

**Option A — websocket for streaming (#181 bus-radar + #184 incidents); HTTP polling (5min) for #182 Raudhah + #183 koper; bus GPS source = tour-leader app (foreground GPS while trip active); Phase 2 revisits hardware GPS for redundancy if tour-leader-phone reliability becomes an issue.**

Option B's polling-only approach misses the latency bar for incidents — a medical emergency alert needs to reach HQ in seconds, not half a minute. Option C's hardware tracker is the right long-term answer but adds procurement + partnership + installation scope that's not MVP-justified. Option A uses the existing F7 tour-leader app (one device already in the bus) as the GPS source and accepts its single-point-of-failure for MVP; hardware redundancy is a Phase 2 option when failure modes become real.

Defaults to propose: **Transport** — websocket server with `/v1/dashboard/live` endpoint; clients subscribe to specific channels (`bus-radar/<departure_id>`, `incidents`, etc.). Fallback: if websocket drops, UI shows disconnect badge and polls at 30s until socket reconnects. **Streaming channels** — #181 (bus position every 30s push), #184 (incidents on event). **Polling channels** — #182 (Raudhah 5min), #183 (koper 5min). **Bus GPS source** — tour-leader F7 app reports location via websocket every 30s while trip active (`trip.start_date <= now() <= trip.end_date`); app keeps foreground + requests location permission on trip start. **Update cadence** — 30s while moving; 5min if stationary > 10min (battery saver). **Location retention** — full bus trajectory per trip retained 6 months; summarized to key stops beyond; detail deleted after 18 months (privacy + storage). **Offline** — last-known-position shown with "(as of HH:MM)" badge; after 10min missing, red badge "offline — last seen…". **Privacy** — bus-level tracking only (marker is "Bus 1" not individual jamaah); jamaah-individual location not tracked. **Incident push channel** — dashboard streaming + WA push to ops-lead + optionally agency-owner per subscription role; severity-filtered.

Reversibility: transport choice websocket vs SSE vs polling can be swapped at connection layer; GPS source can add hardware tracker as parallel input without breaking app-based source.

## Answer

**Decided:** **Option A** — **websocket** bus+incidents, **5m poll** Raudhah+koper, **TL phone GPS** 30s moving / 5m stationary, **6mo trail detail then aggregate**, **offline last-known badge**, **bus-level privacy only**, **incidents WA+stream**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
