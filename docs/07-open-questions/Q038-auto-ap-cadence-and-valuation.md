---
id: Q038
title: Auto-AP posting cadence + PSAK inventory valuation
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8, F9
status: answered
---

# Q038 — Auto-AP posting cadence + inventory valuation method

## Context

Two finance-integration concerns around F8 that the PRD leaves under-specified:

**Auto-AP on GRN** (PRD line 389, module #115, Alur 9.2 line 1231): "GRN immediately posts an AP entry (*Hutang Usaha*) to Finance without re-entry." — but synchronous within the GRN transaction, or async (outbox / event) with eventual consistency? The difference is significant for failure handling.

**Inventory valuation** (module #126 Inventory Health, Neraca Otomatis line 487): the PRD mentions total warehouse asset value (IDR) on the executive dashboard but doesn't specify the **cost method**. PSAK 14 allows FIFO or weighted average; LIFO is prohibited in Indonesia. The choice affects every GRN, every kit-assembly cost-out, every stock-adjustment writeoff, and the closing inventory balance on the Neraca.

## The question

**Part A — Auto-AP posting cadence:**
1. Is the AP entry posted **synchronously** within the GRN transaction (both succeed or both fail), or **asynchronously** via an outbox / event pattern (GRN commits, AP posts shortly after)?
2. On sync-failure: does the GRN roll back, or does it commit with a pending-ap flag?
3. On async-failure: how long can the gap persist before finance diverges from stock? What monitoring catches it?
4. Does the same pattern apply to other finance-integrations (stock-adjustment expense, returns refund offset)?

**Part B — Inventory valuation method:**
5. **FIFO or weighted average?** PSAK 14 permits either.
6. Is the choice per-SKU, per-category, or agency-wide?
7. When is the valuation computed — at write time (per GRN) or at read time (per dashboard render)?
8. How does the method interact with stock-adjustments (writeoff at current method value)?

## Options considered

**Part A — posting cadence:**

- **A1 — Synchronous, in-transaction.** Per ADR 0006 in-process saga style: logistics-svc calls finance-svc.RecordPayable inside the GRN transaction via TxWrapper; failure rolls back the GRN. Stock and AP never diverge.
  - Pros: zero drift; simplest mental model; matches ADR 0006.
  - Cons: hot-path dependency on finance-svc; a finance-svc outage blocks GRN.
- **A2 — Async via outbox pattern.** GRN commits locally; emits `logistics.grn_recorded` event; finance-svc subscribes and posts AP. Reconciliation cron checks drift.
  - Pros: GRN isn't blocked by finance-svc; classic event-sourcing pattern.
  - Cons: eventual consistency; drift window if events lost; more moving parts.

**Part B — valuation:**

- **B1 — Weighted average, agency-wide.** Single moving-average cost per SKU across all warehouses. Computed per-GRN: `new_avg = (old_qty * old_avg + new_qty * new_unit_cost) / (old_qty + new_qty)`.
  - Pros: PSAK-compliant; industry-standard for small-ish inventories; smooth cost curves.
  - Cons: loses FIFO's traceability; unit cost is a synthetic number, not a real invoice.
- **B2 — FIFO per-SKU.** Each stock lot retains its original unit cost; consumed oldest-first.
  - Pros: matches real cost history; stronger audit trail.
  - Cons: requires lot-tracking schema; more complex write path.
- **B3 — Configurable per SKU category.** Bulk fast-moving commodities (koper, ihram) use weighted average; high-value or volatile items (electronics, emas) use FIFO.
  - Pros: best of both.
  - Cons: more config; more logic; reporting complexity.

## Recommendation

**A1 (synchronous auto-AP within GRN transaction) + B1 (weighted average agency-wide).**

For Part A: ADR 0006 already commits us to in-process saga coordination for cross-service writes that must not diverge. Auto-AP-on-GRN is the textbook case for this pattern — stock incremented without the matching liability is a Neraca lie, period. The cost of a brief finance-svc outage blocking GRN (operator retries in 30 seconds) is much smaller than the cost of silent drift between two services' ledgers. If finance-svc is chronically flaky, that's a platform problem to solve directly, not a reason to tolerate drift.

For Part B: weighted average is the pragmatic PSAK-compliant default for agency inventories. FIFO is correct for regulated industries and volatile-cost inventories but introduces lot-tracking complexity that F8's SKU set doesn't justify — koper, ihram, buku doa, travel bags are commodities with stable costs. Revisit FIFO later only if there's a category (e.g. tech devices, gold souvenirs) that genuinely needs it.

Defaults to propose: AP sync via ADR 0006 saga step. Inventory valuation weighted average agency-wide, computed per GRN, cached on `stock_items.weighted_avg_cost` for fast dashboard queries. Returns and writeoffs consume at current avg cost. Reconciliation cron runs nightly comparing logistics stock-value to finance inventory-GL balance; drift > 0.5% alerts ops.

Reversibility: Part A's sync → async is additive (add outbox, keep sync as fallback). Part B's weighted-avg → FIFO requires a migration (every SKU needs historical lot data, which doesn't exist if you didn't track it) — so commit thoughtfully.

## Answer

**Decided:** **A1 synchronous AP** with finance in GRN saga **when finance healthy**; **if finance repeatedly unavailable**, **temporary ops flag** may switch service to **outbox async** (documented break-glass, reconciled nightly). **Inventory:** **weighted-average agency-wide** (FIFO deferred until a SKU class truly needs lots). **Nightly inventory vs GL drift check**.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
