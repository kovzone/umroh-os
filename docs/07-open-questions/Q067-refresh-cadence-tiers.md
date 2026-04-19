---
id: Q067
title: Dashboard refresh cadence tiers
asked_by: session 2026-04-17 F11 draft
asked_date: 2026-04-17
blocks: F11
status: answered
---

# Q067 — Refresh cadence tiers

## Context

PRD uses "real-time" repeatedly (lines 27, 487, 489, 617) as the dashboard promise: *detik itu juga*, *real-time*, *tanpa menunggu akhir bulan*. But "real-time" is aspirational — actual freshness depends on widget type:

- Incident push (F7 W13) must be sub-second (jamaah safety).
- Bus radar (#181) can tolerate 10–30s staleness.
- Seat inventory (#178) needs near-real-time (30s) to avoid double-booking surface conflicts.
- Sales funnel (#180) — 5-min is fine.
- Inventory asset value (#185) — hourly is fine (no one watches this second-by-second).
- Neraca / Laba Rugi (#188) — on-demand with note of last-refresh time.

Picking the wrong cadence means either over-engineering (pay infra cost for freshness no one consumes) or frustrating executives (slow updates on critical data).

## The question

1. **Cadence tiers** — how many buckets (streaming / 30s / 1-min / 5-min / hourly / on-demand)?
2. **Mapping of widgets to tiers** — which widget in which bucket?
3. **Cache TTLs per tier** — how do we implement the bucket?
4. **Streaming transport** — websocket, SSE, long-poll?
5. **Last-updated indicator** — every widget shows "as of HH:MM" timestamp?
6. **Manual refresh** — every widget or dashboard exposes a refresh button?
7. **Scheduled refresh** — does the dashboard page auto-refresh cadence match the slowest widget, or the fastest? Per-widget or page-level?

## Options considered

- **Option A — 5 tiers: streaming / 30s / 5min / hourly / on-demand.** Clear tier buckets mapped to Redis TTLs + websocket for streaming.
  - Pros: clean separation; easy implementation per tier; matches most use cases.
  - Cons: requires per-widget tier classification upfront.
- **Option B — 3 tiers: streaming / 5min / on-demand.** Coarser buckets.
  - Pros: simpler to implement; fewer edge cases.
  - Cons: forces 30s-needed widgets into 5min (too slow for seat inventory).
- **Option C — Per-widget configurable TTL.** Admin sets per-widget refresh; no fixed tiers.
  - Pros: maximally flexible.
  - Cons: operational surface; admin tuning burden; inconsistent UX.

## Recommendation

**Option A — 5 tiers (streaming / 30s / 5min / hourly / on-demand) with explicit per-widget tier assignment.**

Option B's coarse bucketing is too lossy — seat inventory genuinely needs sub-minute freshness to avoid stale double-sold scenarios, while inventory value does not. Option C's per-widget flexibility is an admin burden with no obvious operator who wants that control. Option A matches observed executive-dashboard patterns across the industry and maps naturally to Redis TTL + websocket implementation.

Defaults to propose:

| Tier | TTL / latency | Transport | Default widgets |
|---|---|---|---|
| **Streaming** | ≤ 10s | Websocket | Laporan Insiden (#184), Arus Kas Instan cash events (#187), bus position updates when in-motion |
| **30 seconds** | 30s polling or cache TTL | HTTP poll | Ketersediaan Kursi (#178), Radar Bus positional summary (#181) |
| **5 minutes** | 5min cache TTL | HTTP poll | Sales & Marketing Board (#179, #180), Papan Kinerja CS, Status Raudhah (#182), Pelacakan Koper (#183), Pantauan Eksekusi Logistik (#186) |
| **Hourly** | 1h cache TTL | HTTP | Kesehatan Gudang asset value (#185), Eksekusi Vendor readiness (#177), Likuiditas aging (#189) |
| **On-demand** | no cache; computed on request | HTTP | Neraca, Laba Rugi, Perubahan Ekuitas (#188); executive-landing home refresh |

Every widget displays an "as of HH:MM" timestamp; manual refresh button on every widget. Dashboard page auto-refreshes at its fastest non-streaming cadence (30s for the saudi dashboard; 5min for sales). Manual refresh forces recompute bypassing cache.

Reversibility: cadence tiers are config; moving a widget between tiers is a config change.

## Answer

**Decided:** **Option A** five-tier table as Recommendation; **page auto-refresh = fastest non-streaming widget on that page**; **per-widget manual refresh**; **timestamps always shown**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
