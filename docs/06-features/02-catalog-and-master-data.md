---
id: F2
title: Product Catalog & Master Data
status: written
last_updated: 2026-04-18
moscow_profile: 8 Must Have / 7 Should Have / 1 Could Have
prd_sections:
  - "D. Master Product & Inventory (lines 241–293)"
  - "A. B2C Front-End — Product Catalog + Detail (lines 53–86)"
  - "Alur Logika 1.1 — Global FX & HPP rules (line 1283)"
modules:
  - "#71–86"
depends_on: [F1]
open_questions: []
---

# F2 — Product Catalog & Master Data

## Purpose & personas

The catalog is the **single source of truth** for everything UmrohOS sells: pilgrimage packages (Umrah / Hajj / Badal), hotels, airlines, muthawwif (tour guides), itineraries, seat inventory, add-ons, plus the financial and retail product lines that orbit the core packages. It drives the B2C website, the B2B agent replicated sites, the mobile catalog, CS quoting, ops manifest generation, and the dashboards that track vendor readiness for upcoming departures.

Primary personas:
- **Product admin / CS supervisor** — creates master data, assembles packages, opens departures, runs bulk updates.
- **Agent** — reads the catalogue through their replicated landing page; shares auto-watermarked flyers.
- **Calon jamaah** — browses packages and departures on the B2C portal.
- **Ops manager** — views the vendor-readiness dashboard for upcoming departures.
- **Finance admin** — references HPP (cost of goods) pulled from catalog + FX when building invoices (F5) and journals (F9).

## Sources

- PRD Section D in full (lines 241–293)
- PRD Section A B2C catalog + product detail (lines 53–86)
- PRD Alur Logika 1.1 FX and HPP rules (around line 1283)
- Modules #71–86 from `docs/Modul UmrohOS - MosCoW.csv`

## User workflows

### W1 — Maintain hotel / airline / muthawwif master data

1. Product admin opens the relevant Master Data screen (Hotel / Airline / Muthawwif).
2. Creates or edits a record with:
   - Hotel: name, city (Mecca / Medina / transit), star rating, walking distance to Masjidil Haram or Masjid Nabawi, Google Maps location, 360° tour URL, room-type configurations, photo URLs.
   - Airline: code, name, logo, direct/transit classification, default class. _(Inferred)_ Haramain high-speed rail is also modelled here under an `operator_kind` enum (`airline` / `rail` / `bus`) to keep transport atomic.
   - Muthawwif: name, sanad (education history), lecture video URLs, transparent-background portrait URL for flyer composition.
3. Saves; any **linked packages keep a live reference** — later edits propagate on read _(Inferred)_ because the PRD's "otomatis tertarik" wording implies live-link, not snapshot. If the reviewer wants snapshot-on-link later, a `link_mode` enum can be introduced without data migration.

### W2 — Create a package variant

1. Admin picks a kind: `umrah_reguler`, `umrah_plus`, `hajj_furoda`, `hajj_khusus`, `badal`, `financial`, or `retail`. _(Inferred)_ Polymorphic via `package_kind` enum; all seven share a common header but only travel kinds (Umrah / Hajj / Badal) bear itinerary and departure fields.
2. Enters name, code, description, highlights, cover photo URL.
3. Links master records by reference: hotel(s), airline, muthawwif.
4. Builds itinerary by selecting an `itinerary_template` or drafting inline; renders per-day entries with photo and optional video URLs (module #80).
5. Attaches add-ons (extra-night, train upgrade, etc.) by reference.
6. Saves the package in `draft` status. Moves to `active` after internal review.

### W3 — Open a departure with a seat cap

1. Within an active travel-kind package, admin creates a `package_departure` with departure date, return date, total seats, and per-room-type **list** price (Double / Triple / Quad) in **IDR or USD** per commercial choice (**Q001**: settlement remains IDR at booking lock).
2. _(Inferred)_ Seat cap is admin-set per departure, validated against the linked transport capacity (bus seats + flight seats). A warning surfaces if the cap exceeds known capacity, but it does **not** block — vendors sometimes negotiate extra allotments mid-season.
3. Status starts `open`. Transitions: `open → closed → departed → completed` (or `cancelled` from any pre-departure state).

### W4 — Bulk import via CSV / Excel (module #76)

1. Admin uploads a spreadsheet on the Mass Input screen.
2. Smart mapper suggests column → field associations (date, price, airline code, hotel code) based on header heuristics.
3. Admin confirms or corrects the mapping.
4. System runs a **dry-run validation pass**: any row with missing references, bad enums, or duplicate codes is flagged; the admin can fix in-sheet and re-upload, or discard bad rows.
5. On confirm, all remaining rows are inserted atomically in one transaction — either the whole import commits or none does.
6. A post-import report summarises created vs skipped rows with per-row reasons.

### W5 — Mass price / status update (module #77)

1. Admin filters the package list (by kind, date range, airline, hotel, etc.).
2. Selects N packages and chooses an action: percent price delta, absolute delta, or status flip (`active ↔ archived`).
3. Per **Q002** (answered): single-package edits self-approve. **Mass update** requires second-admin approval when it affects **≥ 3 packages** **or** aggregate absolute price change **> Rp 5,000,000** **or** any line’s **relative change ≥ 5%** vs prior publish price — whichever triggers first (Super Admin–configurable defaults).
4. System writes a **price-history row** per affected package before applying the change (module #77 does not call this out but the PRD's emphasis on audit trail in Section H makes this mandatory — see Q002).
5. FX rule (Alur Logika 1.1): the new price applies only to **unbooked departures and fresh invoices** — it must not retroactively change invoices already issued, nor bookings that are already DP-paid. See "Edge cases" below.

### W6 — Real-time seat decrement on closing (module #84)

1. When a booking transitions into a paid-reservation state (F4 booking saga calls `catalog.ReserveSeats`), the catalog atomically decrements `reserved_seats` on the matching `package_departure`.
2. An event is emitted to an internal topic (see "API surface") that the B2C website, agent replicas, and mobile app subscribe to, so their "Sisa N Seat" label updates within seconds.
3. When seats hit zero, the departure auto-flips to `closed`; the label switches to "Fully Booked".

### W7 — Dynamic flyer generation (modules #78, #79, #81)

1. Sharer (admin or agent) opens the flyer generator and picks a focus theme (Price / Ustadz / Hotel).
2. Engine selects a template, composes a flyer with master data + current seat-remainder label + add-on highlights.
3. Auto-watermarks the sharer's WhatsApp number in a fixed corner per the Section B rule.
4. _(Inferred)_ Rendering is server-side via a headless-browser worker (e.g. Chromium in a render container) — guarantees identical visuals across channels and keeps the agent's watermark un-tamperable.
5. Copywriting helper (module #81, Could Have) drafts a caption from product fields; the sharer can edit before publishing.
6. Sharer clicks "Share to WhatsApp / Instagram / Facebook" which opens the native share sheet with the flyer + caption + referral link pre-filled.

### W8 — Publish interactive itinerary (module #80)

1. On package save, the system renders a public micro-web URL for the itinerary with a per-day breakdown: photos, videos, map location.
2. URL is embedded in the package detail page and the agent replicated landing page.
3. _(Inferred)_ URL is public and crawlable — SEO matters for B2C discovery. Private token-gated variants can be added later for premium-program previews if required.

### W9 — Vendor readiness checklist for a departure (module #86)

1. Ops manager opens the Vendor Readiness dashboard for an upcoming departure.
2. Sees per-departure status flags:
   - **Ticket issued** — from `airline` readiness sub-state
   - **Hotel confirmed** — per-hotel readiness sub-state
   - **Visa progress** — aggregated from F6 visa-svc
3. Risk markers surface when seats are sold-out but any sub-state is not-ready.
4. _(Inferred)_ The flag states are owned by `catalog-svc` (it holds the authoritative per-departure readiness row) but each sub-state is **updated by the owning service via gRPC** — `visa-svc.UpdateDepartureVisaStatus`, `logistics-svc` (for procurement), etc. This keeps the reader dashboard simple without centralising writes.

### W10 — Omni-channel catalog sync (modules #82, #83)

1. On any package / departure / pricing save, `catalog-svc` emits a `catalog.updated` event.
2. Consumers: the central B2C website (SSR cache invalidation), agent replicated landing pages, the mobile catalog service.
3. _(Inferred)_ For phase 1 the event is a simple HTTP fan-out webhook; adding a message broker (Kafka / Redis Streams) is a later optimisation if consumer count grows past a handful.

### W11 — Management dashboard dual view (module #85)

1. Admin opens the package page.
2. Toggle "Public" ⇄ "Management":
   - Public view shows exactly what jamaah and agents see (price, schedule, photos).
   - Management view adds the readiness checklist, HPP breakdown, price history, and current reserved-seat count.

## Acceptance criteria

- Master data (hotel / airline / muthawwif) CRUD is exposed via REST + gRPC with role-scoped permissions (F1).
- A package can be in `draft` or `active`; only `active` packages are visible on B2C / agent replicas / mobile.
- A departure has an atomic seat reservation: concurrent `ReserveSeats(n)` calls cannot oversell; verified by a k6 race test.
- The seat-remainder label on all channels reflects reality within ≤ 5 seconds of a closing event. _(Inferred target — reasonable default.)_
- Bulk CSV import is transactional: either all validated rows commit or the import rolls back; a downloadable per-row report is produced either way.
- Mass price update writes a price-history row per affected package before applying the delta.
- FX rate changes do **not** retroactively alter invoices or paid bookings; this is asserted by a regression test against the integration boundary with F5.
- Vendor readiness sub-states are queryable per departure; risk markers surface programmatically when `seats_remaining == 0` and any sub-state ≠ `ready`.
- Every state-changing call writes to `iam-svc.audit_logs` via `RecordAudit` (F1 contract).

## Edge cases & error paths

- **Concurrent seat reservation.** Under concurrent bookings, `ReserveSeats(n)` must be atomic. _(Inferred)_ Implement as a single SQL statement: `UPDATE package_departures SET reserved_seats = reserved_seats + $n WHERE id = $1 AND reserved_seats + $n <= total_seats RETURNING reserved_seats`. If zero rows returned, fail with `apperrors.ErrConflict`.
- **FX rate change mid-cycle (Alur Logika 1.1 rule).** New rate applies to: future packages, not-yet-invoiced new bookings, and unrealised P/L reports. It must **not** touch issued invoices or DP-paid bookings. Enforced by snapshotting the effective FX rate on each invoice at creation time (F5 data model).
- **Master data edited while sold packages reference it.** _(Inferred)_ Live-link means edits are visible immediately in the management dashboard, but **do not** alter any already-issued customer-facing artefact (printed ticket, PDF itinerary already delivered). The public itinerary micro-web re-renders on the next request with the new values — this is a deliberate feature, not a bug (e.g. hotel name change, room upgrade).
- **Cascading delete attempts.** A hotel or muthawwif record cannot be hard-deleted while any `active` package references it — soft-delete (`deleted_at`) blocks new selections but existing links remain intact.
- **Bulk update approval threshold breach.** When thresholds in **Q002** trip, the action is queued as `pending_approval` until a second admin approves in-console.
- **Cancellation → seat return.** Per **Q004**: **conditional** — `ReleaseSeats` runs **immediately** if the booking has **never received customer funds**; if **any** customer money was posted (DP, installment, lunas), seats stay held until the **refund saga succeeds** (or ops marks forfeiture, audited). Refund failure → seat stays **disputed / not sellable** until ops resolves. Reopen-within-grace attempts `ReserveSeats`; aligns with **Q014** (48h) when applicable.
- **Retail products in catalog vs warehouse inventory.** _(Inferred)_ Retail SKUs (#75) are a thin pointer record in catalog that links to a `logistics-svc.stock_item` by SKU code. Catalog holds the sellable description (price, photos); logistics holds the physical inventory and fulfilment. This keeps the polymorphic catalog simple and defers stock management to F8.
- **Badal product — for whom the pilgrimage is performed.** _(Inferred)_ The "beneficiary" (deceased person the Badal is performed on behalf of) is a **booking-level attribute** (F4), not a catalog field. Catalog just marks `package_kind = 'badal'`. The booking form asks for the beneficiary name + relation at booking time.

## Data & state implications

Owned by `catalog-svc`. Full schema in `docs/03-services/01-catalog-svc/02-data-model.md`. Key additions from this spec:

- `packages.kind` enum extended to `umrah_reguler | umrah_plus | hajj_furoda | hajj_khusus | badal | financial | retail`. Travel kinds (first five) require itinerary + departures; `financial` and `retail` follow a lighter shape.
- `package_departures.reserved_seats` is the hot-path counter; add a check constraint `reserved_seats <= total_seats`.
- New `package_pricing` composite: `(package_departure_id, room_type, list_amount, list_currency, settlement_currency)`. Per **Q001**: catalog may show **USD list** for clarity; **invoice / VA is always IDR** at booking lock using locked FX; payable IDR rounded **once** to nearest **Rp 1,000** (half-up) with rounding GL per F9. `settlement_currency` is always `IDR` in MVP.
- New `package_price_history` table — immutable log of price changes with actor, timestamp, old value, new value, reason.
- New `vendor_readiness` sub-records per departure for `ticket`, `hotel`, `visa` flags, updated via gRPC from owning services.
- Status enums from `02-ubiquitous-language.md`: `package_status`, `departure_status`, `room_type`. New: `vendor_readiness_state` (`not_started | in_progress | ready | blocked`).

## API surface (high-level)

Full contracts in `docs/03-services/01-catalog-svc/01-api.md`. Key surfaces:

**REST (public + internal console):**
- `GET /v1/packages` with filters (kind, branch, date range, airline, hotel, status)
- `GET /v1/packages/{id}` with eager-loaded master references
- `GET /v1/packages/{id}/departures`
- `GET /v1/package-departures/{id}` — seat counts, pricing, readiness summary
- `POST /v1/packages` / `PATCH /v1/packages/{id}` / `DELETE` (soft)
- `POST /v1/packages/import` (multipart CSV)
- `POST /v1/packages/mass-update` (filter + action body)
- `GET|POST|PATCH /v1/hotels`, `/v1/airlines`, `/v1/muthawwif`, `/v1/addons`
- `GET /v1/departures/{id}/readiness` — dashboard feed
- `GET /v1/packages/{id}/flyer?focus=price|ustadz|hotel` — rendered flyer (PNG + caption)

**gRPC (service-to-service):**
- `GetPackage`, `GetPackageDeparture` — read by booking-svc, ops-svc
- `ReserveSeats`, `ReleaseSeats` — **atomic**, called by booking-svc's in-process submit saga and payment-svc's refund saga (per ADR 0006)
- `UpdateVendorReadiness(departure_id, kind, state)` — called by visa-svc, logistics-svc
- `CatalogUpdatedSubscribe` — stream for website / agent / mobile consumers _(Inferred: stream in phase 1; move to broker later if needed)_

## Dependencies

- **F1** (Identity, Access, Audit) — required for token validation, permission checks, and audit writes.
- **Future consumers**: F4 booking (reads packages + departures), F6 visa (updates readiness), F8 logistics (retail SKU link, updates readiness), F9 finance (HPP + FX pull), F10 CRM (campaign targeting by package).

## Backend notes

- **Atomic seat reservation** is the single most performance-critical path. Prefer the single-SQL-statement approach above over pessimistic locking; verified by a k6 race-condition test.
- **Polymorphic package model.** The `kind` enum is wide; keep travel-specific columns (`itinerary_id`, `airline_id`, `muthawwif_id`) nullable and enforce "required for travel kinds" via application-level validation — not via separate tables, which would explode the join count.
- **Price history** is append-only. Add a trigger or CHECK that rejects UPDATE / DELETE.
- **FX snapshot** happens in F5 (invoice creation), not here. Catalog returns the current display price; catalog never decides invoice amount.
- **Flyer rendering** runs in a separate worker pool so a slow headless-browser render doesn't block catalog read traffic. _(Inferred.)_
- **Omni-channel event** is emitted via OTel-traced gRPC stream in phase 1. Kafka / Redis Streams are future options — ADR deferred. (Temporal is itself deferred per ADR 0006 and not a candidate for this fan-out path.)
- **Master data live-link** means reads do the join at query time. Hot paths should use Postgres `pg_hint_plan`-ready joins or materialise a denormalised view if profiling shows a problem.

## Frontend notes

- Three admin surfaces dominate: **Master Data editors** (Hotel / Airline / Muthawwif), **Package Builder** (step-by-step wizard), **Mass Update** (filter → action → preview → confirm).
- Package Builder needs a strong auto-save / versioning story — admins spend 15+ minutes per package. _(Inferred as a UX best practice; the PRD doesn't mandate.)_
- The dual-view toggle (W11) is a single component rendered across all package pages — not separate routes.
- The flyer generator is embedded in the agent B2B portal and surfaces to CS in the internal console. Pick-focus UX is a tab strip with live preview.
- Seat-remainder labels are a shared component that subscribes to the catalog stream; one implementation, reused everywhere.

## Open questions

None blocking — **Q001–Q004** answered **2026-04-18** (`docs/07-open-questions/`). Spec text above reflects those decisions; **Q003** = Bahasa-only MVP copy; **Q002** thresholds may be tuned in Super Admin config.
