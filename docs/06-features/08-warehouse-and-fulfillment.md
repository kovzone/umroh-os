---
id: F8
title: Warehouse, Procurement, Fulfillment
status: written
last_updated: 2026-04-18
moscow_profile: 6 Must Have / 12 Should Have / 2 Could Have
prd_sections:
  - "F. Inventory & Logistics (lines 367–428)"
  - "Alur Logika 7.1–7.5 (lines 1199–1209)"
  - "Sitemap Logistik & Gudang (lines 799–811); Menu (lines 1083–1091)"
  - "Cross-refs: USP line 9, B2C info line 89, B2B portal line 703, Finance vendor DB line 449, ID Card + Luggage Tag line 325, ALL System line 347"
modules:
  - "#109 Permintaan Pembelian, #110 Sinkronisasi Anggaran, #111 Persetujuan Berjenjang, #112 Otomasi Vendor"
  - "#113 Goods Receipt, #114 Quality Control, #115 Pemicu Hutang Otomatis (auto-AP)"
  - "#116 Pelabelan Barcode/SKU, #117 Multi-Warehouse, #118 Peringatan Stok Kritis, #119 Perakitan Paket (kitting), #120 Stock Opname Digital"
  - "#121 Sinkronisasi Ukuran (size sync), #122 Pemicu Pengiriman (Lunas-Trigger), #123 Integrasi Ekspedisi, #124 Pengambilan Mandiri (self-pickup), #125 Retur & Penukaran"
  - "#126 Inventory Health dashboard, #127 Fulfillment & PO Monitor dashboard, #128 Laporan Kerusakan (damage & loss)"
depends_on: [F1, F2, F4, F5, F9]
open_questions: []
---

# F8 — Warehouse, Procurement, Fulfillment

## Purpose & personas

F8 is the **physical-goods backbone**: stock in warehouses, procurement to replenish it, kit assembly against paid bookings, and shipment (or self-pickup) to jamaah. It is the ERP's supply chain — upstream of F7's field-execution scans, downstream of F2's catalog and F4/F5's paid booking.

The USP on PRD line 9 frames F8's reason for existing: **eliminate the silo** between CS (who captures uniform size at registration), Gudang (who preps the kit), and Saudi-side ops (who uses the goods). The old world is three spreadsheets; UmrohOS is one pipeline. Size flows **CS → Gudang → dispatch** without re-entry.

Primary personas:

- **Warehouse staff (Staff_Gudang)** — receives stock, performs QC, assembles kits, ships or hands over for self-pickup.
- **Procurement officer (Staff_Pengadaan)** — files Purchase Requisitions, converses with vendors, closes POs after GRN.
- **Approver (Manager / Director)** — signs off PR / PO per the approval ladder (Q032).
- **Warehouse supervisor** — runs stock opname, investigates variance, approves damage writeoffs.
- **Vendor (external)** — receives PO via Email + WhatsApp, delivers goods for receipt at warehouse (not a UmrohOS user; integrated, not onboarded-as-user).
- **Jamaah (indirect)** — receives the kit at home or picks it up at the office; sees tracking in `/jamaah/logistik`.

Consumers:

- `finance-svc` (F9) — consumes Auto-AP on GRN (#115), inventory valuation for *Neraca*.
- `ops-svc` (F7) — shipped kit includes the koper that gets a luggage-tag QR at airport handling.
- `crm-svc` (F10) — shipment dispatched / delivered events → WhatsApp notifications.
- `booking-svc` (F4) — read shipping address + size per booking.

## Sources

- PRD Section F in full (lines 367–428): Procurement & Purchasing, Inbound & QC, Warehouse Management, Fulfillment & Shipping, Executive Dashboard slice.
- PRD Alur Logika 7.1–7.5 (lines 1199–1209) — workflow phase layer.
- PRD sitemap Logistik & Gudang (lines 799–811) and menu structure (lines 1083–1091).
- Cross-refs: USP line 9, B2C koper-status line 89, B2B jamaah portal line 703, Finance Vendor DB line 449.
- Module list: #109–#128 per `docs/Modul UmrohOS - MosCoW.csv` (20 modules).

## User workflows

### W1 — Purchase Requisition with Budget Sync (modules #109, #110)

1. Procurement officer logs into `/erp/logistik/pengadaan/pr-po` and clicks **New PR**.
2. Picks line items from the SKU catalog (item, quantity, unit, target vendor if known); system pulls last-known unit price as a suggestion.
3. On submit, `logistics-svc` validates against the **production budget** per package (set at catalog-svc during package cost definition — PRD line 375). Budget source: `catalog.packages.budget_breakdown` (planned).
4. **Budget gate is hard** — over-budget PR cannot submit; requires either a budget revision (catalog-svc admin action) or line-item reduction. Error message surfaces `budget_remaining` so the procurement officer can scope down.
5. On success, PR status = `DRAFT → SUBMITTED`, enters the approval queue (W2).
6. Audit log: created_by, timestamp, line items snapshot, budget_snapshot.

### W2 — Multi-level PR / PO approval (module #111)

1. Approver (typically Manager, escalating to Director above a threshold — **Q032**) opens the approval inbox on dashboard or mobile.
2. Sees PR summary, line items, vendor, total IDR, budget remaining after approval.
3. Actions: **Approve**, **Reject** (with reason), **Request Changes** (returns to DRAFT with comments).
4. Per **Q032** (answered): **≤10M Manager**, **10–50M Director**, **>50M Director+CEO**; capex ≥10M → Director minimum; emergency Manager path with Director post-hoc **48h**; **>50M** approval requires **MFA**; multi-line PR routes by highest line threshold.
5. On approval, PR is promoted to PO: status `SUBMITTED → APPROVED`, PO code assigned (`PO-YYYYMMDD-NNNN` format), triggers W3 auto-dispatch.
6. Audit log: approver, timestamp, prior/new status, reason.

### W3 — PO dispatch to vendor (module #112)

1. On W2 approval, `logistics-svc` renders a PO PDF (agency letterhead, digital stamp per PRD line 379).
2. Vendor master lookup (**Q036** — who owns it, finance-svc or logistics-svc): resolves email + WhatsApp number.
3. System sends PO via **Email** (attachment) **and WhatsApp** (PDF attachment + short confirmation link).
4. Dispatch status written to PO: `APPROVED → DISPATCHED`, with email_message_id + whatsapp_message_id for traceability.
5. Failure modes: vendor contact missing, email bounce, WhatsApp number invalid → PO stays DISPATCHED-pending with error, procurement officer re-sends manually via console.

### W4 — Goods receipt + QC + Auto-AP (modules #113, #114, #115)

1. Vendor delivers physical goods to the warehouse receiving bay.
2. Warehouse staff opens `/erp/logistik/gudang/grn` and scans each box barcode or enters SKU + qty against the open PO.
3. System validates received qty against PO qty:
   - **Qty matches** → stock auto-increments on acceptance (PRD line 385).
   - **Qty less than PO** → partial-receipt policy (**Q040** default: accept partial; PO remains `PARTIALLY_RECEIVED`; short-shipped lines stay open for a later GRN).
   - **Qty more than PO** → reject excess; log variance.
4. QC step (module #114): inspector marks each SKU `qc_passed: true` or `qc_passed: false` with defect notes + optional photos.
5. On **QC fail** → system auto-generates a vendor return claim message (PRD line 387) via the same Email + WhatsApp channel as PO dispatch; the failed qty does NOT increment stock.
6. On **QC pass** + stock increment, **Auto-AP** fires (PRD line 389, #115) — `logistics-svc` calls `finance-svc.RecordPayable(po_id, vendor_id, amount, grn_id)` via gRPC. Posting cadence: **Q038** (default: synchronous within the GRN transaction; failure rolls back GRN so AP never drifts from stock).
7. GRN close condition: all PO lines have status `received` or `rejected`. PO status → `RECEIVED` or `CLOSED_SHORT`.

### W5 — SKU barcode labeling (module #116)

1. On first stock entry for a new SKU (or when re-labeling), system generates a unique internal barcode code (format: `SKU-<CATEGORY>-<NNNN>`, Code128 symbology).
2. Warehouse staff prints labels from the console; physically affixes to each unit or carton.
3. Scans at any downstream event (kit assembly, opname, shipment) resolve via this SKU barcode.
4. **Out of scope:** the jamaah-side luggage tag QR from F7 module #95 is a different artefact. **Q037** resolves the coexistence of SKU barcode + luggage tag QR on a shipped koper.

### W6 — Multi-warehouse stock view + low-stock alerts (modules #117, #118)

1. `logistics-svc` tracks stock per `(warehouse_id, sku)` tuple. Warehouse types: `central` (Gudang Pusat), `branch` (Gudang Cabang), `agent_consignment` (stok dititipkan, PRD line 397).
2. **Agent consignment** — stock physically at agent premises is logistics-svc's record of record. Write authority: central ops only; agents submit consumption reports that ops converts into stock-out entries. _(Inferred)_ Reconciliation cadence: monthly opname per agent (module #120).
3. Low-stock alert (#118): per-SKU `reorder_level` config. When `quantity_available ≤ reorder_level`, event `logistics.stock_low` is emitted.
4. **Reorder-point math** (**Q040**): default inferred static per-SKU `reorder_level` set at SKU creation, editable by warehouse supervisor. Demand-forecast-driven auto-reorder (reading upcoming `package_departures` pax counts) is a Phase 2 enhancement.
5. Alert consumers: procurement officer dashboard (badge + notification), daily digest email.

### W7 — Kit assembly (Perakitan Paket, module #119)

1. Kit definition — a named bundle of SKUs (e.g. "Paket Silver 2026": 1× koper, 1× ihram, 1× buku doa, 1× travel kit, 1× kain seragam-L).
2. **Kit composition ownership** (**Q034**): default inferred — kit *definitions* live in `catalog-svc` (tied to package tier); kit *instances* (actual assembled units) live in `logistics-svc`.
3. Kit assembly flow: ops or warehouse supervisor creates a batch (e.g. "50 Silver kits for Ramadhan 2026 departure batch").
4. System validates component availability across the selected warehouse; if any SKU is short, returns `insufficient_stock` with per-SKU gap.
5. On confirmation, **atomic decrement** of all component SKUs in a transaction (PRD line 401); creates `kit_instance` rows, each with a unique kit_instance_code.
6. Failure mode: transaction rollback on any partial failure; no kit with partial components exists in the system.

### W8 — Stock opname (module #120)

1. Supervisor opens an opname session: `/erp/logistik/gudang/opname`, picks warehouse + optional SKU filter.
2. System freezes the digital counts for the selected scope (snapshot).
3. Warehouse staff physically counts; enters actual counts per SKU (scan-assisted).
4. System computes variance per SKU; highlights > X% discrepancies.
5. Supervisor reviews, adds investigation notes, commits adjustments. Every adjustment writes an immutable `stock_adjustment` row (PRD line 403) with reason (opname / damage / loss / reclassification).
6. **Immutable audit trail**: adjustments cannot be deleted, only reversed by a new counter-adjustment with a reason.
7. _(Inferred)_ Opname cadence: monthly per warehouse, quarterly full-inventory. Configurable per warehouse.

### W9 — Size sync from registration (module #121)

1. CS captures jamaah uniform sizes at registration (F3 — `jamaah.uniform_size` field). Sizes per kit: koper, ihram (male), seragam-atas, seragam-bawah, sajadah.
2. On booking-paid-in-full (W10 trigger), `logistics-svc` reads `jamaah.uniform_size` from `jamaah-svc` via gRPC (or via the payload of the dispatch event).
3. Size is attached to the kit instance reserved for that jamaah: `kit_instance.reserved_for_jamaah_id = <id>` with `size_breakdown` materialized.
4. **No re-entry at warehouse** — warehouse staff sees the size breakdown pre-filled when preparing the kit (PRD line 9, USP).
5. If a jamaah has no recorded size → warning flag on the fulfillment queue; ops follows up with CS.

### W10 — Lunas-Trigger fulfillment (Pemicu Pengiriman, module #122)

1. Fulfillment queue — `/erp/logistik/gudang/pengiriman` — lists bookings in status `paid_in_full` only (PRD line 411). `partially_paid` does **not** appear. This is a security gate, not a soft filter.
2. Trigger mechanism: payment-svc emits `booking.paid_in_full` event; logistics-svc subscribes and creates a `dispatch_task { booking_id, jamaah_ids, package_kit_id, address, contact }`.
3. Alternatively (per ADR 0006), payment-svc calls `logistics-svc.DispatchKit(...)` via gRPC synchronously as part of the paid-booking fan-out. The gRPC path is preferred; the event is the async backup.
4. Compensating action: on booking cancellation before dispatch, logistics-svc releases the reserved kit instance (`kit_instance.reserved_for_jamaah_id = NULL`, stock re-available) via `ReleaseStock(booking_id)` gRPC.
5. **After dispatch**: cancellation becomes a **returns-from-jamaah** workflow (**Q035**) — not an automatic stock release.

### W11 — Courier integration / shipment (module #123)

1. Warehouse staff opens a dispatch task; confirms kit components are assembled and boxed; weighs the package.
2. Picks courier per the courier-selection policy (**Q033**) — default inferred: primary courier per region, fallback chain if primary API is down.
3. Calls courier API to generate shipping label + tracking number. Writes `shipments { booking_id, courier, tracking_number, weight, dimensions, created_at }`.
4. Prints label; affixes to package; scans SKU barcodes out of warehouse (stock decrement for non-kit-reserved items if any) and `dispatch_task.status = DISPATCHED`.
5. Tracking number pushed to jamaah via WhatsApp (module #123, PRD line 413) using crm-svc adapter.
6. Courier webhook receives delivery updates → `logistics.shipment_delivered` event → crm-svc notifies jamaah + agent.
7. **Failure modes**: courier API down → retry with backoff; if all configured couriers fail, task remains pending with operational alert.

### W12 — Self-pickup at office (Pengambilan Mandiri, module #124)

1. Alternative to courier: jamaah opts for office pickup at booking time (or ops flips the dispatch method).
2. Instead of courier label, system generates a **QR receipt** (PRD line 415) encoded with the pickup_token.
3. Jamaah receives the QR via WhatsApp; presents at office counter.
4. Office staff scans the QR via an internal console; system verifies the token + pickup_token status + identity match (jamaah KTP photo from F3 shown alongside for staff confirmation).
5. On confirm, `dispatch_task.status = PICKED_UP_SELF`; audit log records staff_id + scanned_at.
6. **QR security model** (**Q041**) — default inferred: signed HMAC token containing `{ dispatch_task_id, jamaah_id, valid_until }`, single-use, 30-day expiry from generation, revocable by ops.

### W13 — Returns & exchanges (module #125, Could Have)

1. Pre-departure window: jamaah reports wrong-size uniform or damaged koper to CS.
2. CS opens a return ticket; logistics-svc creates a `return_request { shipment_id, reason, replacement_needed }`.
3. Courier pickup of the old item OR jamaah drops off at branch.
4. On receipt at warehouse, warehouse staff scans in; QC inspects; routes to restock (if saleable) or writeoff (if damaged beyond use) via W14.
5. Replacement kit reserved + dispatched (new W11 flow).
6. **Stock-out audit**: the replacement is logged as an official stock-out against the original booking, not as a new kit instance (PRD line 417 — preserves traceability).
7. Edge: post-shipment, pre-departure swap of koper size after already printed luggage tags (F7) → old luggage tag invalidated, new one reprinted (F7 W6 flow).

### W14 — Damage & Loss Report (Laporan Kerusakan, module #128, Could Have)

1. Any field-side loss or damage (field ops reports via F7 W13 incident) that involves goods comes back as a damage/loss report to logistics.
2. Warehouse supervisor opens the report, identifies the SKU(s) involved, quantifies loss.
3. Commits as a stock_adjustment with reason = `damage` or `loss` (PRD line 427).
4. Auto-expense posting to finance (F9) — expense category: inventory writeoff.
5. **Returns-from-trip** (**Q035**): post-pilgrimage, some items return from Saudi (walkie-talkies, communication receivers — PRD line 567) — their restock-or-writeoff flow lives here.

### W15 — Executive dashboards (modules #126, #127)

1. **Inventory Health (#126)**: total asset value (IDR) across all warehouses, critical-stock visualization, aging inventory. Values computed per inventory valuation method (**Q038** — FIFO vs weighted average).
2. **Fulfillment & PO Monitor (#127)**: paid-but-unshipped queue (SLA: ship within 7 days of paid-in-full), outstanding PO count, overdue PRs, GRN backlog.
3. Served by `logistics-svc` read endpoints; consumed by F11 per **Q066**: thin **`dashboard-svc`** + service `/v1/metrics/*`; Svelte exec UI + **Grafana on read-replica** for analyst dashboards where configured.

## Acceptance criteria

- **PR budget gate** is hard — over-budget PR cannot submit (W1).
- **Approval ladder** is configurable by role + threshold (Q032); defaults enforced; any PO dispatch requires a valid approval signature chain.
- **PO → vendor dispatch** uses Email + WhatsApp both; at least one successful delivery required before `DISPATCHED` status.
- **GRN** is idempotent on `(po_id, grn_code)`; duplicate GRN submissions do not double-increment stock.
- **QC fail → no stock increment** and vendor return claim auto-generated.
- **Auto-AP** posts to finance-svc in the same transaction as the GRN (Q038 default); failure rolls back both.
- **Atomic kit assembly** — all component decrements succeed or all rollback.
- **Stock opname adjustments** are immutable (no delete; reversal-via-counter-adjustment only).
- **Multi-warehouse stock counts** are per-warehouse; cross-warehouse views are sums, never masks for missing per-warehouse entries.
- **Lunas-Trigger** — fulfillment queue is empty for bookings not in `paid_in_full`. No exception path.
- **Size sync** — no warehouse re-entry of jamaah size; flag missing sizes before dispatch.
- **Courier API failures** — retry chain per Q033; tasks never silently lost.
- **Self-pickup QR** — single-use, signed, expiring (Q041 default 30 days).
- **Damage/loss reports** post matching expense entries to finance; no orphan stock adjustments.
- **Every stock movement** has an audit row linking to origin (PO, GRN, dispatch_task, return_request, stock_adjustment).

## Edge cases & error paths

- **Partial PO receipt** — Q040 default: accept partial, keep PO `PARTIALLY_RECEIVED`; remaining lines stay open.
- **Vendor contact missing on PO dispatch** — PO stays in `APPROVED-pending-dispatch` with operational alert; procurement officer resolves.
- **QC fail after stock already incremented** (race or operator error) — compensating `stock_adjustment` with reason `qc_late_fail`, AP entry counter-journaled via `finance-svc.ReversePayable`.
- **Auto-AP synchronous failure** — GRN transaction rolls back; stock NOT incremented; operator sees error. Alternative async mode (Q038) documented in Backend notes.
- **Kit component stock insufficient at assembly time** — pre-flight check returns `insufficient_stock` with per-SKU gap; no partial assembly.
- **Duplicate dispatch for same booking** — idempotency on `(booking_id)` unique per active dispatch_task.
- **Cancellation after dispatch** — Q035 default: CS reaches out to jamaah for return courier; kit is not auto-released until physically received back.
- **Courier webhook missing delivery confirmation > 14 days after shipment** — escalation alert to ops.
- **Self-pickup QR expired** — ops regenerates with new expiry; old token marked revoked.
- **Agent consignment drift** (agent says 5 on hand, system says 7) — opname variance workflow investigates; adjustment requires supervisor approval.
- **Backorder situation** (#119 kit assembly blocked on one SKU) — partial shipment vs wait-full (Q040 default: wait-full in MVP; partial-ship policy is Phase 2).
- **Jamaah has no recorded size** — fulfillment queue flags the booking; CS re-engages jamaah; no default-size substitution (avoids wrong-size shipments).
- **Luggage tag QR printed before kit shipped** — F7 module #95 reprints after size/booking changes; F8 shipment carries kit, F7 attaches luggage tag at airport (Q037 confirms handoff point).
- **Returns received during busy pre-departure week** — queue aging visible in Fulfillment & PO Monitor dashboard; SLA configurable.

## Data & state implications

Owned by `logistics-svc` (schema defined in `docs/03-services/07-logistics-svc/02-data-model.md`). Key tables referenced / added in this spec:

- `warehouses` — existing; add `warehouse_type` enum (`central | branch | agent_consignment`).
- `stock_items` — existing; add `reorder_level`, `sku_barcode`, `valuation_method` (if per-SKU overrides global Q038 default).
- `stock_ledger` — new: immutable ledger of every stock movement. `{ id, sku, warehouse_id, delta, reason_kind enum, reason_id (FK to grn/po/dispatch_task/return_request/stock_adjustment), by_user, at }`.
- `purchase_orders` — existing; add `budget_source_id` (FK to catalog package budget), `dispatched_via_email_at`, `dispatched_via_whatsapp_at`, approval_chain jsonb.
- `po_lines` — existing.
- `goods_received_notes` — existing; add `ap_posted_at`, `ap_entry_ref`.
- `grn_lines` — new: per-SKU received qty + qc status + notes + defect photos (gcs URLs).
- `kit_definitions` — existing; add `owner_service` enum (`catalog | logistics`) per Q034 resolution.
- `kit_instances` — new: materialized kit, per-jamaah reservation. `{ id, kit_definition_id, assembled_at, reserved_for_jamaah_id nullable, reserved_for_booking_id nullable, status enum (assembled/reserved/dispatched/returned) }`.
- `dispatch_tasks` — new: `{ booking_id, jamaah_ids[], dispatch_method enum (courier/self_pickup), status enum, shipment_id nullable, self_pickup_token nullable }`.
- `shipments` — existing; add `weight_grams`, `dimensions_cm`, `label_url`, `delivered_at`, `last_webhook_at`.
- `return_requests` — new: `{ shipment_id, reason, replacement_needed, inbound_received_at, outcome enum (restocked/writeoff/replaced) }`.
- `stock_adjustments` — new: immutable; `{ sku, warehouse_id, delta, reason enum (opname/damage/loss/reclassification/qc_late_fail), notes, by_user, at }`.
- `self_pickup_tokens` — new: `{ dispatch_task_id, token_hmac, valid_until, used_at nullable, revoked_at nullable }`.

New enums:

- `po_status`: existing — extend with `partially_received`, `closed_short`, `dispatched`.
- `dispatch_status`: `pending`, `assembled`, `dispatched`, `delivered`, `picked_up_self`, `returned`, `cancelled`.
- `warehouse_type`: `central`, `branch`, `agent_consignment`.
- `stock_adjustment_reason`: `opname`, `damage`, `loss`, `reclassification`, `qc_late_fail`, `field_loss`.

## API surface (high-level)

Full contracts in `docs/03-services/07-logistics-svc/01-api.md` — spec already planned. Key surfaces this draft confirms:

**REST (ops console):**
- `POST /v1/purchase-requisitions` + `GET /v1/purchase-requisitions` (filterable by status)
- `POST /v1/purchase-requisitions/{id}/approve|reject|request-changes`
- `GET /v1/purchase-orders` + `POST /v1/purchase-orders/{id}/dispatch`
- `POST /v1/grn` (create) + `POST /v1/grn/{id}/qc` (pass/fail per line)
- `GET /v1/stock-items?warehouse_id=` + `PATCH /v1/stock-items/{id}`
- `POST /v1/kits` (kit definition) + `POST /v1/kit-instances/assemble` (batch)
- `POST /v1/opname-sessions` + `POST /v1/opname-sessions/{id}/commit`
- `GET /v1/dispatch-tasks?status=` + `POST /v1/dispatch-tasks/{id}/ship|pickup-confirm`
- `POST /v1/shipments/{id}/label` + courier webhook `POST /v1/shipments/webhook`
- `POST /v1/returns` + `POST /v1/returns/{id}/receive`
- `POST /v1/stock-adjustments` (writeoffs, damage/loss)

**gRPC (service-to-service):**
- `DispatchKit(booking_id, jamaah_ids, address)` — called by payment-svc on paid-in-full (ADR 0006 saga).
- `ReleaseStock(booking_id)` — compensating action for cancelled-before-dispatch.
- `CheckStock(sku, warehouse_id)` — pre-flight for ops actions.
- `ReservePrerequisites(booking_id, package_kit_id)` — optional pre-reservation on booking creation (not on paid) — Phase 2.
- `GetVendor(vendor_id)` — if Q036 puts vendor master in finance-svc, logistics-svc calls finance-svc here.

**Events emitted** (per `03-events.md`):
- `logistics.po_approved`, `logistics.grn_recorded`, `logistics.shipment_dispatched`, `logistics.shipment_delivered`, `logistics.stock_low`.
- Add: `logistics.kit_assembled`, `logistics.return_received`, `logistics.stock_adjusted`.

## Dependencies

- **F1** (IAM) — RBAC for procurement / warehouse / approver roles; audit logs.
- **F2** (catalog) — package definitions, kit definitions (Q034), budget source for PR validation.
- **F4** (booking) — shipping address, jamaah list per booking, booking status (paid_in_full trigger).
- **F5** (payment) — `booking.paid_in_full` event / gRPC trigger; cancellation → release stock.
- **F7** (ops) — luggage tag QR coexistence on shipped kits (Q037); damage/loss reports from field (W14).
- **F9** (finance) — auto-AP posting on GRN (#115); inventory valuation (Q038); vendor master ownership boundary (Q036); expense posting for damage/loss.
- **F10** (CRM) — WhatsApp dispatch notifications and delivery confirmations.
- **External** — courier APIs (selection per Q033), Email SMTP, WhatsApp Business API.

## Backend notes

- **Atomic stock movements**: every GRN, kit assembly, dispatch, return, adjustment is wrapped in a single DB transaction using `WithTx` per backend conventions. No partial-state mid-transaction visible to readers.
- **Ledger-style stock model**: `stock_ledger` is the source of truth; `stock_items.quantity` is a materialized cache. Reconciliation query compares sum(ledger.delta) vs cache per SKU+warehouse; drift triggers alert.
- **Auto-AP contract** (Q038 default: sync): `logistics-svc` wraps `finance-svc.RecordPayable` gRPC call inside the GRN transaction via a saga-step pattern. On remote failure, local GRN rolls back — behaves like an in-process transaction per ADR 0006. Async alternative uses an outbox pattern; cadence documented per Q038 answer.
- **Kit assembly** uses the same `WithTx` + sqlc pattern; components decrement via `UPDATE stock_items SET quantity = quantity - $n WHERE sku = $sku AND quantity >= $n RETURNING quantity` for each SKU within the transaction. Missing row → rollback.
- **Courier adapter** — abstract `CourierAdapter` interface with per-courier implementations (JNE, J&T, SiCepat, etc.). Selection logic (Q033) chooses at shipment creation time; fallback chain retried with backoff.
- **Self-pickup QR** — signed HMAC tokens with server-side secret (same pattern as F7 luggage tag signing; shared config).
- **Size sync** — read-through to `jamaah-svc` via gRPC; cached per booking for the duration of the dispatch task (size doesn't change post-registration in normal flow).
- **Budget sync** — read-through to `catalog-svc` via gRPC on PR submit; re-validated at approval time to catch budget changes between submit and approve.
- **Agent consignment write model** — only ops can write; agents submit consumption reports via a separate "agent portal" endpoint that creates pending stock-out entries requiring supervisor approval.
- **Stock opname** — atomic snapshot via `SELECT ... FOR UPDATE` on the SKU rows for the scope; prevents mid-opname stock movements.
- **PSAK valuation** (Q038) — pluggable calculator in `logistics-svc/util/valuation/` with FIFO and weighted-average implementations; default configurable.

## Frontend notes

- **Procurement console** — PR creation form with budget-remaining ticker; approval inbox with inline approve/reject; PO dispatch status surface.
- **Warehouse console** — GRN scanning UI (barcode input + qty confirmation), QC step with photo upload, kit assembly batch screen, opname session with scan-assisted counting.
- **Dispatch queue** — Lunas-Trigger gate visible (empty state explains "awaiting paid-in-full"); per-task size breakdown, courier selector, label print, tracking number display.
- **Self-pickup counter console** — QR scanner, jamaah KTP photo side-by-side, single-confirm button.
- **Returns & damage** — claim entry form with reason picker + photo upload; supervisor approval queue for writeoffs.
- **Dashboards** — Inventory Health tile (treemap by warehouse + category), Fulfillment & PO Monitor (paid-unshipped aging chart, overdue PR counter).
- **Mobile approver UI** (Svelte) — approval inbox on phone for approvers on-the-go; biometric confirm for high-value approvals (Q032).

## Open questions

None blocking — **Q032–Q041** answered **2026-04-18** (`docs/07-open-questions/`). Operational defaults in the body reflect those answers unless superseded by Super Admin config.

**Residual engineering defaults (verify during build):**

- PR approval thresholds align with **Q032** (see W3 step 4).
- Kit definitions live in `catalog-svc`, kit instances in `logistics-svc` — Q034 may collapse to single owner.
- Partial PO receipt accepted by default; remaining lines stay open — Q040.
- Reorder point is static per-SKU, editable by supervisor; demand-forecast is Phase 2 — Q040.
- Self-pickup QR is HMAC-signed, single-use, 30-day expiry, ops-revocable — Q041.
- Stock opname monthly per warehouse, quarterly full-inventory — config-driven, no hard PRD source.
- Auto-AP is synchronous within GRN transaction by default (ADR 0006 in-process saga) — Q038 may switch to async outbox.
- FIFO valuation default for PSAK compliance — Q038 may switch to weighted average.
- Returns reclassify as restock (if saleable) vs writeoff; supervisor decides — Q035 may formalize.
