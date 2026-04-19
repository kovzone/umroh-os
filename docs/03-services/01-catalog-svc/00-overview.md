# catalog-svc — Overview

## Purpose

Master record for all sellable products: packages, hotels, airlines, muthawwif, itineraries, and seat inventory.

## Bounded context

Product Catalog. See `docs/02-domain/00-bounded-contexts.md` § 2.

## PRD source

PRD section D — Master Product & Inventory.

## Owns (data)

- `packages` — sellable Umrah/Hajj/Badal packages
- `package_departures` — specific departure dates with seat inventory
- `package_pricing` — per-room-type **list** price + `list_currency` + `settlement_currency` (MVP IDR settlement; **Q001**)
- `hotels` — master hotel records (with photos, 360 tours, distance to mosque)
- `airlines` — master airline records
- `muthawwif` — master tour-leader records
- `itinerary_templates` — reusable itinerary structures
- `addons` — optional extras (e.g. extra night in Medina)

## Boundaries (does NOT own)

- Bookings (`booking-svc`)
- Pricing rules per agent / per campaign (`crm-svc`)
- Stock items / kits (`logistics-svc`)
- Vendor payment accounts (`finance-svc`)

## Interactions

- **Inbound:** booking-svc reads packages and reserves seats; ops-svc reads packages for manifest generation; crm-svc reads packages for marketing materials.
- **Outbound:** none in the synchronous path.

## Notable behaviors

- **Real-time seat inventory** — when a booking is created, the saga calls `ReserveSeat`; on cancel, `ReleaseSeat`.
- **360° photos and video tours** for hotels — stored in GCS with signed URLs.
- **Bundling** — packages can include addons.
- **Bulk import/export** via Excel/CSV for catalog updates.
- **Auto-generated flyers** — uses package data to render watermarked promotional materials. May call out to a separate render worker.
