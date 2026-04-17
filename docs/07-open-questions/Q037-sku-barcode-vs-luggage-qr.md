---
id: Q037
title: SKU barcode vs F7 luggage-tag QR coexistence on shipped kits
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8, F7
status: open
---

# Q037 — SKU barcode vs F7 luggage-tag QR coexistence

## Context

Two different barcoded artefacts land on a jamaah's koper:

- **F8 SKU barcode (module #116)** — identifies the item as a stock SKU (e.g. `SKU-KOPER-0042`). Applied at warehouse, used for stock movements, GRN scans, kit assembly, dispatch-out scans. Item-scoped.
- **F7 luggage-tag QR (module #95)** — identifies the owner jamaah and booking (e.g. signed token encoding `{jamaah_id, booking_id, tag_seq}`). Applied for airport handling, ALL System check-in, Smart Bus Boarding. Jamaah-scoped, signed, changes per person.

Timeline: F8 ships the koper to the jamaah weeks before departure (with SKU barcode attached). F7 prints luggage tags typically at the airport-handling stage, or at the jamaah's hand-over to ground crew. Both end up on the same koper — but when?

## The question

1. **Does the SKU barcode persist on the koper after shipment, or get removed before dispatch?**
2. **When is the F7 luggage tag physically attached — at warehouse pre-ship, at airport handoff, or somewhere between?**
3. **If both stay on the koper, does the F7 tag override the F8 barcode at airport scan?** Or does handling staff know to scan the F7 QR and ignore the F8 barcode?
4. **If F7 tag is generated pre-ship and included in the shipment:** does the warehouse print both? Does F8's `dispatch_task` know about the luggage-tag code?
5. **Luggage-tag reprint after allocation change** — F7 spec already notes that ID cards / luggage tags can reprint after room or bus reassignment. If the tag is already on a shipped koper, how does reprint work?
6. **What's the handoff point** in the code — does F8 emit `shipment_dispatched` and F7 listens? Is there a `logistics-svc.EmitLuggageTagsForDispatch(...)` call?

## Options considered

- **Option A — Keep SKU barcode; F7 luggage tag applied at airport only.** Koper ships with SKU barcode. Jamaah brings koper to the airport. Ground crew prints + applies luggage tag at the departure counter. Two artefacts never coexist pre-airport.
  - Pros: clean separation; no coordination needed between F8 dispatch and F7 tag generation.
  - Cons: extra step at already-busy airport; risk of last-minute tag printer failure.
- **Option B — Generate F7 luggage tag at F8 dispatch; attach both to koper pre-ship.** Warehouse prints SKU sticker for internal ops + F7 luggage tag for jamaah. Both travel home with the jamaah. At airport, ground crew scans F7 QR only.
  - Pros: tag is ready before airport crunch; no airport printer dependency.
  - Cons: two artefacts on the same koper (visual clutter, potential for ground crew to scan the wrong one); locks jamaah to allocation at ship-time; reprint-after-reshuffle requires a new sticker sent later.
- **Option C — SKU barcode replaced by luggage tag at dispatch.** Warehouse removes / covers the SKU barcode before dispatch and applies the F7 luggage tag. Only one artefact on the koper post-ship.
  - Pros: single artefact, no confusion.
  - Cons: SKU traceability lost on the shipped unit; if returned, warehouse can't scan it as the original SKU.

## Recommendation

**Option A — SKU barcode stays on the koper; F7 luggage tag is attached at airport handling.**

Option B sounds nice on paper but the real airport surface already has all the infrastructure for tag printing (portable label printers, scanners, trained ground crew). Locking the jamaah to a pre-ship luggage allocation breaks F7's reprint-after-reshuffle story (which is a real need per F7 edge case for mid-trip reallocation). Option C loses inbound traceability on returns and creates a weird "cover the old sticker" workflow.

Option A respects the bounded contexts: F8's SKU barcode is for inbound/internal ops; F7's luggage tag is for airport + field; they don't need to coexist on the same physical artefact. The airport-handling workflow is already designed around tag printing at check-in (F7 W7 ALL System).

Defaults to propose: SKU barcode stays permanently on the koper (small, corner-placed, not visually dominant). F7 luggage tag printed at airport check-in, applied over / alongside SKU barcode (doesn't matter — ground crew scan the QR, stock team scan the SKU). On return post-trip, F7 tag is removed/detached, SKU barcode is still there for warehouse restock. Reprints-after-reshuffle work because no tag exists pre-airport — reallocation before check-in just changes the data behind the tag that gets printed. No coordination needed between F8 dispatch event and F7 tag generation.

Reversibility: switching to Option B later is a process change, not a data-model change — the SKU+tag relationship would be explicit then, but nothing stored now prevents the switch.

## Answer

TBD — awaiting stakeholder input. Deciders: airport-handling team lead, warehouse supervisor, ops lead.
