---
id: Q039
title: Saudi-side warehouse scope
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8
status: open
---

# Q039 — Saudi-side warehouse scope

## Context

PRD Section F treats warehouses as Indonesia-side: *Gudang Pusat, Gudang Cabang, stok dititipkan di Agen* (line 397). But operational reality has **Saudi-side consumables** that need management:

- Zamzam (bulk water for distribution — F7 module #107 handles the per-jamaah scan, but where does the water come from? Bought in Saudi, staged at the hotel?).
- Local SIM cards and phone kits distributed on arrival.
- Receivers / audio devices used during city tours (F7 module #106, F7 references "Receiver Asset Recovery" at PRD line 567 — implies they return to Indonesia).
- Banners / agency signage positioned at hotels for jamaah way-finding.
- Food vouchers, Saudi-bought pharmaceutical items for medical kits.
- Emergency replacement inventory (spare ihram, spare koper when one is damaged mid-trip).

The PRD doesn't say whether these sit in an F8 warehouse (Saudi-type `warehouse_type`), get treated as in-field consumables outside stock management, or belong somewhere else entirely.

## The question

1. **Should Saudi-side consumables be tracked as F8 stock?** If yes, as a new `warehouse_type: saudi_hotel` or `saudi_staging`?
2. **Which items specifically?** — zamzam, receivers, banners, SIM cards, spare kit, medical supplies.
3. **If tracked as stock:** who does the GRN when supplies arrive in Saudi (muthawwif? local vendor? tour leader)? Offline support?
4. **If NOT tracked as stock:** where does the cost recognition happen — direct expense at purchase?
5. **Inventory flow across borders** — do some items ship from Indonesia to Saudi (F8 kit items the jamaah brings) vs buy-in-Saudi? Both?
6. **Consumption recording** — zamzam is per-jamaah-scanned at distribution; is that decrementing a stock quantity somewhere?
7. **Returns to Indonesia** — audio receivers specifically are mentioned as returning (PRD line 567); do they come back to a specific warehouse?

## Options considered

- **Option A — Saudi-side is out of F8 scope entirely.** Saudi consumables are treated as direct expenses at purchase (finance-svc books them). No stock tracking, no GRN, no opname. F7 field apps scan/distribute without a stock backend.
  - Pros: no schema expansion; recognizes operational reality (Saudi procurement often cash-based, no digital receipts).
  - Cons: zero stock visibility; no reorder alerts; no audit trail for asset recovery (receivers).
- **Option B — Minimal Saudi warehouse: only durable-return assets tracked (receivers, banners).** Consumables (zamzam, SIM cards) remain expenses. Assets that return to Indonesia get stock records with `warehouse_type: saudi_active`.
  - Pros: solves the receiver-recovery problem; minimal scope.
  - Cons: split logic (some Saudi items tracked, some not); expense-vs-asset boundary can be ambiguous.
- **Option C — Full Saudi warehouse support: all Saudi-side inventory tracked.** New `warehouse_type: saudi_staging` (Mecca/Medina hotel consignment); muthawwif or tour leader acts as GRN receiver via field app; opname on departure back to Indonesia.
  - Pros: complete inventory picture; supports Phase 2 analytics.
  - Cons: substantial scope; field-app complexity; offline support required; custom workflows for each consumable type.

## Recommendation

**Option B — track only durable-return assets (receivers, banners, spare koper stock) as Saudi-side warehouse; consumables are direct expense.**

Option A loses the audit trail on returning assets — the PRD explicitly calls out Receiver Asset Recovery as a workflow, so we need some Saudi-side stock notion. Option C is the right long-term answer but lands in Phase 2; MVP doesn't need full Saudi inventory. Option B threads the needle: track the 10–20% of items that are durable + returning + expensive (receivers, spare inventory); leave the 80% of consumables (zamzam, SIMs, vouchers, single-use medical) as direct-expense purchases that finance books.

Defaults to propose: `warehouse_type: saudi_active` for Saudi-hotel-staged durables. At departure from Indonesia, receivers + spare items are "dispatched" to the Saudi warehouse (stock movement between warehouses, not a jamaah dispatch). At return, a bulk-GRN reclassifies them back to Indonesia central. Damage/loss during trip: F7 incident → F8 damage report (W14). Consumables: finance-svc books expenses at purchase; F7 field apps log consumption events (zamzam distribution scans) without a decrement target — they're metrics, not stock. Receiver returns: tour leader submits return manifest from F7 app (photos of sealed box) on arrival back; warehouse GRN's against that manifest.

Reversibility: expanding to Option C later is additive — new warehouse_type values, new field flows. Nothing in B forecloses the larger version.

## Answer

TBD — awaiting stakeholder input. Deciders: ops lead (Saudi operations experience), warehouse supervisor, finance director (direct-expense vs asset treatment).
