# logistics-svc — Overview

## Purpose

Warehouse, procurement, kit assembly, and shipping. Triggered by booking payment status to fulfill jamaah kits.

## Bounded context

Logistics. See `docs/02-domain/00-bounded-contexts.md` § 8.

## PRD source

PRD section F — Inventory & Logistics.

## Owns (data)

- `warehouses`
- `stock_items` (per warehouse)
- `purchase_orders` and `po_lines`
- `goods_received_notes`
- `kit_definitions` and `kit_components`
- `shipments` and `shipment_items`
- `vendors` (logistics vendors only — finance owns financial vendor records)

## Boundaries (does NOT own)

- Bookings (`booking-svc`)
- Vendor financial records (`finance-svc` — same vendor IDs may exist in both)
- Visa providers / payment gateways

## Interactions

- **Inbound:** payment-svc calls DispatchKit via gRPC as part of the paid-booking event fan-out (per ADR 0006).
- **Outbound:** courier APIs (label printing, tracking), booking-svc (read shipping address), finance-svc (PO journal entries).

## Notable behaviors

- **Procurement workflow** — digital PR → multi-level approval → auto-dispatch to vendors.
- **GRN with QC** — goods received against PO, with quality check.
- **Multi-warehouse** — stock tracked separately per warehouse.
- **Kit assembly** — bundle of stock items per kit definition; pulled from warehouse on fulfillment.
- **Shipping integration** — label printing, real-time tracking via courier APIs.
