---
slice: S1
title: Slice S1 — Integration Contract
status: draft
last_updated: 2026-04-22
pr_owner: Elda
reviewer: Elda (solo-exec S0 per § Current operating mode)
task_codes:
  - S1-J-01
  - S1-J-02
  - S1-J-03
  - S1-J-04
---

# Slice S1 — Integration Contract

> Slice S1 = "Discover + draft booking" — B2C catalog browsing → booking form → draft booking; staff routes auth'd via F1. This file is the wire-level agreement between services for that user journey.
>
> **Incremental build:** only sections that have a landed `S1-J-*` card are filled. Unfilled sections are added when their corresponding Joint card ships. Amendments after merge follow the Bump-versi rule in `docs/contracts/README.md` (Changelog append for additive, `-v2.md` for breaking).

## Purpose

S1 is the first user-facing slice. The B2C flow: a calon jamaah opens the catalog, browses packages, picks a departure, and submits a draft booking. Staff-facing flows (CS closing, admin review) are auth'd via F1. S1 services: `catalog-svc` (read), `booking-svc` (draft + saga), `iam-svc` (staff auth), `payment-svc` (scaffolded but invoice work is S2).

## Scope

**In scope for S1 contracts** (landed incrementally via `S1-J-01..04`):

- Public catalog read endpoints (list, detail, departure) — `S1-J-01`, this card.
- Draft booking REST endpoint (`POST /v1/bookings`) — `S1-J-02`.
- `ReserveSeats` / `ReleaseSeats` gRPC (internal, booking→catalog) — `S1-J-03`.
- S1 booking-state decision paragraph (only `draft` needed in S1; doc completeness per Q006) — `S1-J-04`.

**Out of scope for S1 contracts** (deferred to later slices):

- Admin write endpoints on `catalog-svc` (create / update / soft-delete packages).
- Payment/invoice/VA (S2 — `S2-J-01..04`).
- Fulfillment / finance journal (S3 — `S3-J-*`).
- Webhooks, events, CRM, dashboards (later slices).

---

## § Catalog

*(Landed via `S1-J-01`, 2026-04-20.)*

Public read-only endpoints for B2C catalog browsing. All endpoints are unauthenticated — the catalog is a public surface. Auth'd admin/staff catalog endpoints land in a later slice, not here.

### Endpoints


| Method | Path                          | Auth   | Purpose                                                                                    |
| ------ | ----------------------------- | ------ | ------------------------------------------------------------------------------------------ |
| `GET`  | `/v1/packages`                | public | List active packages, filterable.                                                          |
| `GET`  | `/v1/packages/{id}`           | public | Package detail incl. itinerary + departure summary list + master-data refs.                |
| `GET`  | `/v1/package-departures/{id}` | public | Departure detail incl. `remaining_seats`, pricing per room type, vendor-readiness summary. |


Paths are stable and match `docs/03-services/01-catalog-svc/01-api.md` row-for-row for the public-read subset. Route prefix `/v1/` is mounted under `catalog-svc`'s REST adapter (port 4002 in dev, proxied by `gateway-svc` at 4000 for public traffic).

### Conventions honored

- **Q001 (currency).** Responses carry `list_amount` (integer), `list_currency` (`IDR` or `USD`), and `settlement_currency` (always `"IDR"` in MVP). Catalog never commits to a payable amount — that lock happens in `payment-svc` at invoice creation with a snapshotted FX rate. Consumers must treat `list_amount` + `list_currency` as display-only.
- **Q003 (single-language MVP).** All human-readable strings (`name`, `description`, `itinerary` day labels) are returned in Bahasa Indonesia (`id-ID`). No `{lang}` query param, no `translations` array in the response — a future multi-language card will evolve the contract via Bump-versi.
- **No auth** on any endpoint in this section. `gateway-svc` passes these through without token validation.

### `GET /v1/packages`

**Query params:**


| Param            | Type                   | Required | Notes                                                                                                                       |
| ---------------- | ---------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------- |
| `kind`           | string enum            | no       | One of `umrah_reguler`, `umrah_plus`, `hajj_furoda`, `hajj_khusus`, `badal`, `financial`, `retail`. Omit to list all kinds. |
| `departure_from` | string (ISO-8601 date) | no       | Lower bound on departure date (earliest departure ≥ this date).                                                             |
| `departure_to`   | string (ISO-8601 date) | no       | Upper bound on departure date.                                                                                              |
| `airline_code`   | string                 | no       | IATA airline code filter.                                                                                                   |
| `hotel_id`       | string (ULID)          | no       | Filter to packages referencing this hotel.                                                                                  |
| `cursor`         | string (opaque)        | no       | Pagination cursor from previous response's `page.next_cursor`.                                                              |
| `limit`          | integer                | no       | Page size, 1–100, default 20.                                                                                               |


`status` is **not** a query param — the public list returns only `packages.status = 'active'` packages; `draft` / `archived` are filtered server-side.

**Response** (`200 OK`):

```json
{
  "packages": [
    {
      "id": "pkg_01JCDE...",
      "kind": "umrah_reguler",
      "name": "Umrah Reguler 12 Hari — Ramadan",
      "description": "Paket umrah reguler 12 hari dengan keberangkatan di pertengahan Ramadan.",
      "cover_photo_url": "https://cdn.umroh-os.example/pkg/01JCDE.../cover.jpg",
      "starting_price": {
        "list_amount": 38500000,
        "list_currency": "IDR",
        "settlement_currency": "IDR"
      },
      "next_departure": {
        "id": "dep_01JCDF...",
        "departure_date": "2026-05-12",
        "return_date": "2026-05-23",
        "remaining_seats": 12
      }
    }
  ],
  "page": {
    "next_cursor": "eyJsYXN0X2lkIjoicGtnXzAxSkNERS4uLiJ9",
    "has_more": true
  }
}
```

**Errors:**


| Status | Body `error.code`     | When                                                                               |
| ------ | --------------------- | ---------------------------------------------------------------------------------- |
| `400`  | `invalid_query_param` | Any query param fails validation (bad date, unknown `kind`, `limit` out of range). |
| `400`  | `invalid_cursor`      | `cursor` is malformed or from an incompatible list.                                |
| `500`  | `internal_error`      | Unexpected server error.                                                           |


### `GET /v1/packages/{id}`

**Path params:** `id` — package ULID.

**Response** (`200 OK`):

```json
{
  "package": {
    "id": "pkg_01JCDE...",
    "kind": "umrah_reguler",
    "name": "Umrah Reguler 12 Hari — Ramadan",
    "description": "...",
    "highlights": ["Direct Jakarta-Jeddah", "Hotel 4-star Mecca 5 min walking", "Muthawwif S3 Al-Azhar"],
    "cover_photo_url": "https://cdn.umroh-os.example/pkg/01JCDE.../cover.jpg",
    "itinerary": {
      "id": "itn_01JCDG...",
      "days": [
        { "day": 1, "title": "Keberangkatan dari Jakarta", "description": "...", "photo_url": "..." }
      ],
      "public_url": "https://umroh-os.example/itinerary/itn_01JCDG..."
    },
    "hotels": [
      { "id": "htl_01JCDH...", "name": "...", "city": "mecca", "star_rating": 4, "walking_distance_m": 300 }
    ],
    "airline": {
      "id": "arl_01JCDI...",
      "code": "GA",
      "name": "Garuda Indonesia",
      "operator_kind": "airline"
    },
    "muthawwif": {
      "id": "mtw_01JCDJ...",
      "name": "Ustadz ...",
      "portrait_url": "..."
    },
    "add_ons": [
      { "id": "addon_01JCDK...", "name": "Extra night Medinah", "list_amount": 2500000, "list_currency": "IDR", "settlement_currency": "IDR" }
    ],
    "departures": [
      {
        "id": "dep_01JCDF...",
        "departure_date": "2026-05-12",
        "return_date": "2026-05-23",
        "remaining_seats": 12,
        "status": "open"
      }
    ]
  }
}
```

**Errors:**


| Status | Body `error.code`   | When                                                                                                                                           |
| ------ | ------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `404`  | `package_not_found` | No package matches `{id}`, or the package is not `status = 'active'` (the public endpoint does not leak existence of draft/archived packages). |
| `500`  | `internal_error`    | Unexpected server error.                                                                                                                       |


### `GET /v1/package-departures/{id}`

**Path params:** `id` — departure ULID.

**Response** (`200 OK`):

```json
{
  "departure": {
    "id": "dep_01JCDF...",
    "package_id": "pkg_01JCDE...",
    "departure_date": "2026-05-12",
    "return_date": "2026-05-23",
    "total_seats": 45,
    "remaining_seats": 12,
    "status": "open",
    "pricing": [
      { "room_type": "quad", "list_amount": 38500000, "list_currency": "IDR", "settlement_currency": "IDR" },
      { "room_type": "triple", "list_amount": 41500000, "list_currency": "IDR", "settlement_currency": "IDR" },
      { "room_type": "double", "list_amount": 45500000, "list_currency": "IDR", "settlement_currency": "IDR" }
    ],
    "vendor_readiness": {
      "ticket": "ready",
      "hotel": "ready",
      "visa": "in_progress"
    }
  }
}
```

- `remaining_seats = total_seats − reserved_seats`, computed at read time. The field is eventually consistent with concurrent `ReserveSeats` calls (contracted in `§ Inventory`, lands via `S1-J-03`).
- `status` values: `open` / `closed` / `departed` / `completed` / `cancelled`. Only `open` and `closed` reach a public reader for a given departure; `departed` / `completed` / `cancelled` return `404` to public callers (staff endpoints expose them via a future slice).
- `vendor_readiness` sub-states: `not_started | in_progress | ready | blocked`.

**Errors:**


| Status | Body `error.code`     | When                                                                              |
| ------ | --------------------- | --------------------------------------------------------------------------------- |
| `404`  | `departure_not_found` | No departure matches `{id}`, or status is `departed` / `completed` / `cancelled`. |
| `500`  | `internal_error`      | Unexpected server error.                                                          |


### Error envelope (shared)

All non-2xx responses use:

```json
{
  "error": {
    "code": "<snake_case>",
    "message": "<human-readable, id-ID>",
    "trace_id": "<hex>"
  }
}
```

`trace_id` is the OTel span ID surfaced per `docs/04-backend-conventions/03-logging-and-tracing.md`.

### Response ID format

All entity IDs are **ULID** strings with a type prefix: `pkg_`, `dep_`, `itn_`, `htl_`, `arl_`, `mtw_`, `addon_`. Consumers must treat IDs as opaque strings (no parsing).

### Honored by implementation

- `S1-E-02` — catalog-svc read endpoints. Must conform exactly to the shapes above.
- `S1-L-02` — catalog UI (frontend). Consumes exactly the shapes above.

---

## § Booking

*(Landed via `S1-J-02`, 2026-04-21.)*

Creates a **draft** booking via `POST /v1/bookings`. Scope for S1 is deliberately narrow: one endpoint, one transition (`∅ → draft`). Transitioning a draft onward (`draft → pending_payment`) happens via `POST /v1/bookings/{id}/submit`, which contracts in a later slice and runs the in-process booking saga per ADR 0006 + Q006's KTP+passport gate.

### Endpoint


| Method | Path           | Auth                                       | Purpose                                                                                                                                                                  |
| ------ | -------------- | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `POST` | `/v1/bookings` | **public OR staff** — see Auth rules below | Create a `draft` booking for B2C self / B2B agent / CS-closing flows. Persists only; does not issue VA, does not reserve seats (that's `POST /v1/bookings/{id}/submit`). |


### Auth rules

The endpoint accepts three closing channels per F4 W1–W3. The `channel` field in the request body is the source of truth for which authentication mode applies:


| `channel`   | Required auth                                                                             | How `agent_id` / `staff_user_id` are populated                                                                                                                                                                            |
| ----------- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `b2c_self`  | **public** — no token required                                                            | Neither field set.                                                                                                                                                                                                        |
| `b2b_agent` | **public** — no token required                                                            | `agent_id` comes from the request body (stamped by the agent's replicated landing page, sourced from the `ref=<agent_code>` referral link). `gateway-svc` may additionally validate the referral token before forwarding. |
| `cs`        | **staff** — F1 PASETO / JWT required; `iam-svc.CheckPermission(booking.create_on_behalf)` | `staff_user_id` is filled from the token's claims server-side, not from the request body.                                                                                                                                 |


Requests with `channel = "cs"` but no valid F1 token return `401 unauthorized`. Requests with `channel = "b2c_self"` and a token are still accepted (token ignored). Requests with `channel = "b2b_agent"` must include `agent_id`; absence returns `422 validation_failed`.

### Idempotency

Callers SHOULD send an `Idempotency-Key` HTTP header (string, ≤ 128 chars). Scope is `(channel, key)` + the fingerprint of the request body. Replay behavior:

- Same key + **identical body** within **24 h** → returns the original booking (same `id`, same `created_at`, `Idempotency-Replayed: true` header).
- Same key + **different body** within 24 h → returns `409 idempotency_conflict` with the original `booking.id` so the caller can diagnose.
- Same key after 24 h → treated as a new request (the prior record has been GC'd).
- Missing `Idempotency-Key` is allowed (server mints a row per request). B2C and B2B self-serve SHOULD include one; CS closing MUST include one.

### Request body

```json
{
  "channel": "b2c_self",
  "package_id": "pkg_01JCDE...",
  "departure_id": "dep_01JCDF...",
  "room_type": "quad",
  "lead": {
    "full_name": "Budi Santoso",
    "email": "budi@example.com",
    "whatsapp": "+62811234567",
    "domicile": "Jakarta"
  },
  "jamaah": [
    {
      "full_name": "Budi Santoso",
      "email": "budi@example.com",
      "whatsapp": "+62811234567",
      "domicile": "Jakarta",
      "is_lead": true
    },
    {
      "full_name": "Siti Aminah",
      "whatsapp": "+62811234568",
      "domicile": "Jakarta",
      "is_lead": false
    }
  ],
  "add_on_ids": ["addon_01JCDK..."],
  "agent_id": null,
  "notes": null
}
```

**Required fields** (422 if missing):


| Field                | Type          | Notes                                                                                   |
| -------------------- | ------------- | --------------------------------------------------------------------------------------- |
| `channel`            | enum          | `b2c_self`                                                                              |
| `package_id`         | string (ULID) | Must reference an `active` package; 404 if not found or `draft`/`archived`              |
| `departure_id`       | string (ULID) | Must belong to the named package and have `status = open`; 404 if not found or not open |
| `room_type`          | enum          | `double`                                                                                |
| `lead.full_name`     | string        | Non-empty, ≤ 120 chars                                                                  |
| `lead.whatsapp`      | string        | E.164 format                                                                            |
| `lead.domicile`      | string        | Non-empty                                                                               |
| `jamaah[]`           | array         | At least one entry; exactly one must have `is_lead = true` and match `lead` fields      |
| `jamaah[].full_name` | string        | Non-empty, ≤ 120 chars per entry                                                        |
| `jamaah[].domicile`  | string        | Non-empty per entry                                                                     |


**Conditionally required:**

- `lead.email` — required when `channel = b2c_self`; optional for `b2b_agent` / `cs`.
- `agent_id` (ULID) — required when `channel = b2b_agent`; must be `null` otherwise.

**Optional:**

- `add_on_ids` (array of ULIDs) — must reference add-ons listed on the package; unknown IDs return `422 validation_failed`.
- `notes` (string, ≤ 2000 chars) — free-form internal note (visible to staff only).
- `jamaah[].email` and `jamaah[].whatsapp` for non-lead jamaah — encouraged but not required at draft.

### Documents (Q006)

**No documents are required at this endpoint.** Q006 (`docs/07-open-questions/Q006-minimum-docs-for-booking.md`) mandates KTP + passport scan per jamaah **before `draft → pending_payment`** — i.e. at the submit endpoint (`POST /v1/bookings/{id}/submit`), not here. Draft creation remains lightweight so a customer can fill the form without their docs immediately at hand; docs get uploaded via the Customer Portal (F3) before they click "Submit & Pay".

Consumers of this contract must **not** require document uploads before calling `POST /v1/bookings`. The submit-side contract (future `S1-J-` or follow-up) will spell out the blocking check.

### Response — `201 Created`

```json
{
  "booking": {
    "id": "bkg_01JCDL...",
    "status": "draft",
    "channel": "b2c_self",
    "package_id": "pkg_01JCDE...",
    "departure_id": "dep_01JCDF...",
    "room_type": "quad",
    "agent_id": null,
    "staff_user_id": null,
    "lead": {
      "full_name": "Budi Santoso",
      "email": "budi@example.com",
      "whatsapp": "+62811234567",
      "domicile": "Jakarta"
    },
    "jamaah": [
      {
        "id": "bkgitem_01JCDM...",
        "full_name": "Budi Santoso",
        "email": "budi@example.com",
        "whatsapp": "+62811234567",
        "domicile": "Jakarta",
        "is_lead": true,
        "status": "active"
      },
      {
        "id": "bkgitem_01JCDN...",
        "full_name": "Siti Aminah",
        "whatsapp": "+62811234568",
        "domicile": "Jakarta",
        "is_lead": false,
        "status": "active"
      }
    ],
    "add_ons": [
      { "id": "addon_01JCDK...", "name": "Extra night Medinah", "list_amount": 2500000, "list_currency": "IDR", "settlement_currency": "IDR" }
    ],
    "pricing": {
      "list_amount": 79500000,
      "list_currency": "IDR",
      "settlement_currency": "IDR",
      "breakdown": [
        { "line": "2 × quad @ Rp 38,500,000", "list_amount": 77000000 },
        { "line": "1 × Extra night Medinah", "list_amount": 2500000 }
      ]
    },
    "notes": null,
    "created_at": "2026-04-21T02:14:32.105Z",
    "expires_at": "2026-04-21T02:44:32.105Z"
  }
}
```

- `booking.status` is always `"draft"` on this endpoint — no other state is reachable here.
- `booking.pricing` shows the **list amount** at draft time per Q001; a payable IDR amount only gets locked at `POST /v1/bookings/{id}/submit` (via `payment-svc` VA issuance with FX snapshot).
- `booking.expires_at` is a **non-binding** hint that a draft held for ≥ 30 min without being submitted may be GC'd. The hard expiry (VA TTL) only starts after submit. *(Inferred — keeps the draft table from growing unbounded.)*
- `booking.add_ons` carries the full add-on record (not just IDs) for client convenience, matching the catalog shape in `§ Catalog`.

Response header `Idempotency-Replayed: true` is set only when the request matched a prior (key, body) pair.

### Errors

All errors use the shared error envelope defined in `§ Catalog`.


| Status | Body `error.code`      | When                                                                                                                                                                                                                                          |
| ------ | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `400`  | `invalid_json`         | Request body is not valid JSON.                                                                                                                                                                                                               |
| `401`  | `unauthorized`         | `channel = cs` and no valid F1 token, or token lacks `booking.create_on_behalf` permission.                                                                                                                                                   |
| `404`  | `package_not_found`    | `package_id` not found, or package `status ≠ active`.                                                                                                                                                                                         |
| `404`  | `departure_not_found`  | `departure_id` not found, does not belong to the package, or `status ≠ open`.                                                                                                                                                                 |
| `404`  | `add_on_not_found`     | Any `add_on_ids` entry not found or not linked to the package.                                                                                                                                                                                |
| `409`  | `idempotency_conflict` | Same `Idempotency-Key` with a different request body within the 24 h window. Body includes `original_booking_id`.                                                                                                                             |
| `422`  | `validation_failed`    | Field-level validation errors. Body's `error.details` carries an array of `{field, code, message}` per failed field (e.g. `{"field": "lead.whatsapp", "code": "not_e164", "message": "WhatsApp harus dalam format internasional (+62...)"}`). |
| `500`  | `internal_error`       | Unexpected server error.                                                                                                                                                                                                                      |


Notably absent: **no `409 seats_unavailable`** at this endpoint. Seat reservation is `POST /v1/bookings/{id}/submit`'s job, per F4 W8. Creating a draft does not check seats; the submit saga may still fail with `seats_unavailable` if capacity ran out between draft and submit.

### Conventions honored

- **Q006 (min docs at submit).** Draft creation requires no documents; the requirement applies at submit. See the "Documents (Q006)" sub-section above.
- **Q001 (currency).** `pricing.list_amount` + `list_currency` + `settlement_currency = "IDR"` mirror the catalog shape. No payable commitment here — that happens at VA issuance in `payment-svc`.
- **Q003 (single-language MVP).** Validation error messages (`error.details[].message`) are in `id-ID`. `error.code` remains `snake_case` English — machine-readable, language-invariant.
- **ULID IDs** with type prefixes: `bkg_` for bookings, `bkgitem_` for booking items. Reuse `pkg_`, `dep_`, `addon_` from `§ Catalog`.

### Honored by implementation

- `S1-E-03` — `booking-svc` draft creation handler. Must conform exactly to the shape above.
- `S1-L-04` — booking UI (frontend) for all three channels. Consumes the shape above; wires the channel-specific auth behavior client-side.

---

## § Inventory

*(Landed via `S1-J-03`, 2026-04-21.)*

Internal gRPC contract for atomic seat reservation on a `package_departure`. These are the **first internal-service methods** in this slice — all REST endpoints above are public; these two are the concurrency-critical inter-service calls that keep capacity honest under parallel submits.

### Service + methods


| Service                     | Method         | Purpose                                                                              |
| --------------------------- | -------------- | ------------------------------------------------------------------------------------ |
| `catalog.v1.CatalogService` | `ReserveSeats` | Atomic decrement of available seats on a departure; idempotent per `reservation_id`. |
| `catalog.v1.CatalogService` | `ReleaseSeats` | Atomic increment, reversing a prior `ReserveSeats`; idempotent per `reservation_id`. |


Both methods target `catalog-svc`'s gRPC port (50052 in dev). Callers today: `booking-svc` (submit saga) and, after refund-flow lands, `payment-svc`. Callers MUST NOT attempt to call these over REST — there is no REST equivalent for the inventory path, on purpose: the atomic SQL guard + dedup lookup live in a single gRPC handler transaction.

### Idempotency (per reviewer option a)

Both methods take a caller-supplied `reservation_id` (ULID) that catalog-svc stores in a dedup table alongside the `departure_id` + `seats` it applied. The dedup row has a TTL (default 24 h, bounded to [1, 168]). Semantics:

- Same `reservation_id` + same `(departure_id, seats)` within TTL → the call is a **replay**: no decrement/increment happens; response carries `replayed: true` and the original `reserved_at` / `remaining_seats`.
- Same `reservation_id` + different `departure_id` or `seats` within TTL → `ALREADY_EXISTS` (`reservation_id_conflict`) — signals a programmer bug in the caller.
- Same `reservation_id` after TTL expiry → treated as a fresh request (TTL row has been GC'd).

The booking saga uses the booking's own ULID (`bkg_...`) as the `reservation_id`. That gives "one booking, one reservation, one ID" by construction, and the TTL always exceeds the saga's wall-clock budget by an order of magnitude.

### `catalog.v1.CatalogService/ReserveSeats`

**Request** (`ReserveSeatsRequest`):


| Field                   | Type     | Required | Notes                                                                                                                                                           |
| ----------------------- | -------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `reservation_id`        | `string` | yes      | ULID. Caller-supplied. Per `§ Booking`'s convention, booking-svc passes the booking's own ULID (`bkg_...`). Payment-svc later will pass a refund-specific ULID. |
| `departure_id`          | `string` | yes      | ULID of the `package_departures` row (matches `GET /v1/package-departures/{id}` in `§ Catalog`).                                                                |
| `seats`                 | `int32`  | yes      | Number of seats to reserve. Must be ≥ 1. For the booking saga, equals `len(booking.jamaah[])`.                                                                  |
| `idempotency_ttl_hours` | `int32`  | no       | Defaults to **24**. Clamps to `[1, 168]`. Optional override for long-running flows.                                                                             |


**Response** (`ReserveSeatsResponse`):


| Field             | Type          | Notes                                                                                                                                                               |
| ----------------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `reservation`     | `Reservation` | `{ reservation_id, departure_id, seats, reserved_at, expires_at }`. `expires_at` = `reserved_at + idempotency_ttl_hours`; it's the dedup-row TTL, NOT a VA timeout. |
| `remaining_seats` | `int32`       | Post-decrement count. Matches what `GET /v1/package-departures/{id}` would return next, modulo parallel reservers.                                                  |
| `replayed`        | `bool`        | `true` if this call matched a prior `(reservation_id, departure_id, seats)`; no new decrement occurred.                                                             |


**Failure codes (gRPC `status.Code`):**


| Code                  | `error.code`              | When                                                                                                                                                                         |
| --------------------- | ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `FAILED_PRECONDITION` | `insufficient_capacity`   | Atomic SQL returned zero rows — `reserved_seats + seats > total_seats` at commit time. Expected outcome; caller should surface to user as "seat just taken".                 |
| `NOT_FOUND`           | `departure_not_found`     | Unknown `departure_id`, or departure status is `departed` / `completed` / `cancelled` (inventory is frozen).                                                                 |
| `INVALID_ARGUMENT`    | `invalid_request`         | Missing required fields, `seats ≤ 0`, malformed ULID, or `idempotency_ttl_hours` out of range.                                                                               |
| `ALREADY_EXISTS`      | `reservation_id_conflict` | Same `reservation_id` previously seen with a **different** `departure_id` or **different** `seats`. Response message carries the original values for the caller to diagnose. |
| `INTERNAL`            | `internal_error`          | Catch-all; includes database transaction failures.                                                                                                                           |


### `catalog.v1.CatalogService/ReleaseSeats`

**Request** (`ReleaseSeatsRequest`):


| Field            | Type     | Required | Notes                                                                                                                                                                                                                                                         |
| ---------------- | -------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `reservation_id` | `string` | yes      | Must reference a prior `ReserveSeats` (live or already released).                                                                                                                                                                                             |
| `seats`          | `int32`  | no       | Optional partial-release override. If omitted, releases the full `seats` the reservation originally held. If specified and less than the original, a partial release is recorded (remainder stays reserved). Partial release > original → `INVALID_ARGUMENT`. |
| `reason`         | `string` | no       | Free-form audit note (≤ 256 chars): `"saga_failure"`, `"refund_settled"`, `"departure_cancelled"`, etc. Written to `iam.audit_logs` via F1 `RecordAudit`.                                                                                                     |


**Response** (`ReleaseSeatsResponse`):


| Field             | Type       | Notes                                                                            |
| ----------------- | ---------- | -------------------------------------------------------------------------------- |
| `released`        | `Released` | `{ reservation_id, departure_id, seats_released, released_at }`.                 |
| `remaining_seats` | `int32`    | Post-increment count.                                                            |
| `replayed`        | `bool`     | `true` if the reservation was already fully released; no new increment occurred. |


**Failure codes:**


| Code                  | `error.code`             | When                                                                                                                      |
| --------------------- | ------------------------ | ------------------------------------------------------------------------------------------------------------------------- |
| `NOT_FOUND`           | `reservation_not_found`  | `reservation_id` never existed or TTL-expired from the dedup table.                                                       |
| `INVALID_ARGUMENT`    | `invalid_request`        | `seats` ≤ 0, exceeds the reservation's original `seats`, or malformed ULID.                                               |
| `FAILED_PRECONDITION` | `reservation_not_active` | Reservation exists in dedup but in a terminal released state where the specific `seats` partial override would overshoot. |
| `INTERNAL`            | `internal_error`         | Catch-all.                                                                                                                |


### Atomic pattern (reference, not implementation)

`ReserveSeats` executes a single SQL statement inside the handler transaction, followed by the dedup-row write:

```sql
UPDATE package_departures
SET reserved_seats = reserved_seats + $seats
WHERE id = $departure_id
  AND reserved_seats + $seats <= total_seats
RETURNING reserved_seats, total_seats;
```

- Zero rows returned → `FAILED_PRECONDITION insufficient_capacity`.
- Rows returned → commit the dedup row (`reservation_id`, `departure_id`, `seats`, `expires_at`) in the same transaction.

This is the **repo-wide pattern for atomic capacity decrements** and is referenced by F2 `Acceptance criteria` ("concurrent `ReserveSeats(n)` cannot oversell"). The actual DDL + sqlc query for the `package_departures` row + dedup table lives in `S1-E-02` (catalog-svc scaffolding); this contract deliberately does **not** prescribe table names or index choices beyond the capacity guard.

### S1-E-01 engineering review (BL-EGV-001)

**Date:** 2026-04-22  
**Reviewer:** Engineering

**Approved** for `S1-E-01` (seat concurrency + DB transaction gate) with the following findings:

1. The contract requires a single atomic guard statement (`UPDATE ... WHERE reserved_seats + $seats <= total_seats RETURNING ...`) as the seat decrement primitive; this is sufficient to prevent oversell under concurrent writers when executed inside one DB transaction.
2. The dedup write (`reservation_id`, `departure_id`, `seats`, `expires_at`) is explicitly required to commit in the same transaction as the capacity decrement, which preserves exactly-once semantics for retries/replays.
3. Failure mapping is explicit and safe for concurrent races (`FAILED_PRECONDITION insufficient_capacity` for zero-row atomic update, `ALREADY_EXISTS reservation_id_conflict` for idempotency misuse).
4. Compensation is contractually pinned to idempotent `ReleaseSeats` and disallows policy drift (Q004 timing lives in caller saga policy, mechanism remains transactional in catalog).

**Gate note:** This review approves the `S1-J-03` contract and atomic pattern. Runtime enforcement lands in `S1-E-02` / `S1-E-03` implementation PRs.

### Compensation story

Two callers invoke this pair today; their compensation patterns are documented here so implementers don't hand-roll inconsistent retry logic.

**booking-svc submit saga** (`S1-E-03`, per ADR 0006 in-process saga model):

1. `ReserveSeats(reservation_id=bkg_..., departure_id, seats=len(jamaah))` — must succeed before any downstream step.
2. Downstream steps (VA issuance via `payment-svc`, etc.).
3. On any downstream step failure: `ReleaseSeats(reservation_id=bkg_..., reason="saga_failure")`. Idempotent — if the saga crashes mid-compensation and retries, the second call is a `replayed: true` no-op.
4. If the saga itself retries from scratch (e.g. transient network error), step 1 is a `replayed: true` no-op because the booking ULID is stable.

**payment-svc refund flow** (contracted in a later slice; summarized here for Q004 alignment): calls `ReleaseSeats` per **Q004** timing rules — which are **conditional**, not always-immediate:

- **Never-funded cancellation** (booking had no customer money posted — e.g. draft abandoned, VA expired, customer cancelled before paying): `ReleaseSeats` fires **immediately** when the booking transitions to `cancelled`. This is the saga's job, not payment's.
- **Funded cancellation** (any DP, installment, or lunas on the booking): `ReleaseSeats` fires **only after the refund saga reaches a terminal-success state**. If the refund fails (bank rejection), the seats stay in the dedup table's reserved state — a "held / disputed" slot that is NOT sellable until ops resolves manually. Re-attempting `ReleaseSeats` with the same `reservation_id` after ops clears the dispute is safe (replay or fresh release, depending on TTL).
- **Reversal within grace window** (customer changes mind; per Q014, 48 h): caller attempts `ReserveSeats` with a **new** `reservation_id` against the same departure. If capacity is gone, the reversal honestly fails and the UI surfaces "seat no longer available".

Implementers **must not** release seats purely on the booking status reaching `cancelled` — the caller decides timing based on Q004's conditional rule. This contract enforces the mechanism (gRPC methods + idempotency); the policy lives in the saga code.

### Conventions honored

- **Q004 (seat return timing).** Cited above; compensation prose is the authoritative read of Q004's conditional rule for S1 saga code.
- **ADR 0006 (in-process saga).** Contract assumes callers implement in-process compensation, not Temporal. Temporal returns for F6; no changes to this contract are needed when that lands.
- `**reservation_id` is caller-supplied.** booking-svc passes `bkg_...` ULIDs; payment-svc will pass its own ULIDs later. No server-generated reservation IDs — the caller owns the key so retries don't need round-trips.

### Honored by implementation

- `S1-E-02` — catalog-svc `ReserveSeats` / `ReleaseSeats` handlers. Must implement the atomic SQL + dedup pattern exactly as described and return the failure codes listed.
- `S1-E-03` — booking-svc submit saga. Must call `ReserveSeats` first + compensate via `ReleaseSeats` on any downstream failure.
- Later: payment-svc refund flow, once its card lands.

---

## § Booking States

*(Landed via `S1-J-04`, 2026-04-21.)*

S1 exercises **only the `draft` state** of the full booking lifecycle. Everything past `draft` (payment, fulfillment, cancellation, completion) lives in **S2 and later slices**. This section documents the decision so implementers know what to reject + what to defer, and consumers know how far the current contract reaches.

### States in scope for S1


| Status  | Entered by                                      | Exited by                                                                                                                                                     | Notes                                                                     |
| ------- | ----------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| `draft` | `POST /v1/bookings` (contracted in `§ Booking`) | No transitions are contracted in S1. `POST /v1/bookings/{id}/submit` — which moves `draft → pending_payment` — is contracted in a **later slice** (not here). | Only persistent status a booking can have at the end of any S1 code path. |


Concretely for implementers:

- `S1-E-03` (booking-svc draft handler) **MUST** persist every new booking with `status = "draft"` and no other value.
- `S1-E-03` **MUST NOT** expose any endpoint that attempts to transition a booking out of `draft` — retry semantics for aborted submits are a later-slice concern.
- Any read endpoint added in S1 that observes `bookings.status` **MUST** short-circuit on non-`draft` values and return `404` (the record is not part of S1's observable state). This is a defence-in-depth check against future schema drift where rows from S2+ leak into S1 reads.

### Forward-looking full state machine (not contracted in S1)

Reproduced from `docs/06-features/04-booking-and-allocation.md § W4` for readers' reference. **Do not implement any transition below in S1 code.** The states are named here only so callers of `POST /v1/bookings` understand where their draft will eventually flow; each transition's trigger, compensation, and error handling belong to its own `Sx-J-*` contract.

```
draft
  │
  │ (submit + VA issued)
  ▼
pending_payment ──┐
  │               │ (VA expires without payment — Q010: VA TTL 24h, global config)
  │               ▼
  │             expired ───┐
  │                        │
  │ (first payment received)
  ▼                        │
partially_paid             │
  │                        │
  │ (cumulative ≥ total)   │
  ▼                        │
paid_in_full               │
  │                        │
  │ (departure begins)     │
  ▼                        │
departed                   │
  │                        │
  │ (return + alumni move) │
  ▼                        │
completed                  │
                           │
cancelled ◄────────────────┘
  ↑
  └─────────── (from any pre-departure state, via refund flow)
```

Plus two system-fault states documented in F4's edge-cases (not reached in S1): `failed` (saga mid-flight fault; operator triage) and `expired` (VA timed out before any payment; terminal, distinct from customer-cancelled).

### Q006 gate on the `draft → pending_payment` transition (future contract)

When a future slice contracts `POST /v1/bookings/{id}/submit`, **Q006** mandates that the handler block the transition unless every jamaah on the booking has both a KTP scan and a passport scan present in the documents store (OCR may be pending/queued; presence is sufficient). See `docs/07-open-questions/Q006-minimum-docs-for-booking.md`:

- Enforced **per-jamaah** — one jamaah's documents do NOT satisfy another's.
- Enforced as a **hard block** in the B2C self-serve UI and as a **422 response** from the submit endpoint.
- `draft` creation itself (contracted in `§ Booking` above) is **not** subject to this gate — customers can park a draft while they gather documents.

This sub-section exists in S1 to freeze the intent so S2+ authors cannot quietly soften the gate or move it to a different transition. The rule is: documents gate `draft → pending_payment`, and only that transition.

### Honored by implementation

- `S1-E-03` — booking-svc draft handler. Persists only `draft`; refuses any status transition request; short-circuits reads of non-`draft` rows.
- Future slice — submit handler (`POST /v1/bookings/{id}/submit`) will enforce the Q006 gate when contracted.

---

## § S1 UI placement (core-web vs multi-app)

**Authoritative product default:** Q009 (`docs/07-open-questions/Q009-frontend-scaffolding-and-tooling.md`) targets a **monorepo of multiple SvelteKit apps** (`apps/b2c`, `apps/b2b`, `apps/admin`, `apps/field`) plus shared `packages/ui` and `packages/api-client`.

**Engineering consensus for S1 (this slice):**

1. **Implement S1 in `apps/core-web` first.** The repo currently has one scaffolded frontend (`apps/core-web`). Public catalog + draft-booking journeys ship here until a split is justified (bundle size, independent deploy, or auth boundary pain).
2. **Organize routes for a clean future split.** Use SvelteKit layout groups and/or clear URL segments so **public B2C/B2B-attributed** surfaces stay separable from **internal console** routes (see **§ UI route matrix** below). Prefer new S1 pages under a predictable subtree (for example `(public)/...` or `/packages/...` + `/booking/...`) rather than overloading `/` beyond the existing landing.
3. **When to extract `apps/b2c` (and shared packages).** Follow Q009 timing: introduce `packages/api-client` + additional apps when the second major audience or deploy cadence forces it — not as a blocking prerequisite for `S1-L-02..04`.

This section exists so `S1-L-01` / backlog gate **BL-LGV-001** can cite **code + contract** (screen inventory + file paths under `apps/core-web`) without implying the multi-app scaffold already exists.

---

## § UI route matrix (S0-L-01 / BL-FE-PLN-001)

Planning map for **`apps/core-web`** (and future staff console routes) so S1 UI, gateway auth, and `S1-E-04` middleware share one list of **public vs internal** surfaces. Path segments in `{braces}` are dynamic.

### Conventions

- **Surface** — `public`: B2C or B2B-attributed journeys without a staff F1 session. `internal`: staff console behind F1 (`docs/06-features/01-identity-and-access.md`).
- **S1 status** — `active-now`: required for slice S1 (“discover + **draft** booking”). `coming-next`: same journey but depends on a **future** `Sx-J-*` contract row or another slice (noted per row).
- **Role anchor** — Rows use F1 **personas** (Super Admin, Branch Manager, staff subtypes, Agent, Jamaah). Extra granularity is only in **Permission / scope** (for example `booking.create_on_behalf` for CS), not new role IDs.

### Matrix

| UI route pattern | Surface | S1 status | Role / persona (F1 anchor) | Permission / scope | Auth (session) | Wire anchor (`slice-S1`) | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `/` _(Inferred — may ship as `/catalog`.)_ | public | active-now | Jamaah (calon jamaah); Agent (same catalog read) | n/a | none | § Catalog — `GET /v1/packages` | Product landing stays **`/`**; package **list** ships at **`/packages`** (see **§ S1-L-01**). B2B replica may restyle only. |
| `/packages/{package_id}` | public | active-now | Jamaah; Agent | n/a | none | § Catalog — `GET /v1/packages/{id}` | “Package detail” from F2/F4 browse path. |
| `/packages/{package_id}/departures/{departure_id}` _(Inferred — may be merged into package detail.)_ | public | active-now | Jamaah; Agent | n/a | none | § Catalog — `GET /v1/package-departures/{id}` | Departure detail + seat label; optional standalone page. |
| `/booking/{package_id}` | public | active-now | Jamaah; Agent | `channel`: `b2c_self` or `b2b_agent`; B2B requires `agent_id` / referral context per F4 W2 | none (B2C/B2B) | § Booking — `POST /v1/bookings` | Self-Booking Engine / agent-stamped flow (F4 W1–W2). |
| `/id/{agent_slug}` _(Inferred host/path — PRD “replicated site”; F4 W2.)_ | public | coming-next | Agent; Jamaah | Referral / `ref` handling on gateway _(partially in § Booking auth table)_ | none | § Catalog + § Booking | **Gap:** canonical URL + `ref` validation contract for `gateway-svc` not fully specified here. |
| `/checkout/{booking_id}` | public | coming-next | Jamaah; Agent | n/a until submit is contracted | none until staff CS variant | § Booking (draft only today); submit + VA **not** in S1 | F4 W1 step 4–5: combines **submit saga** + VA (F5). **Split:** submit (`POST /v1/bookings/{id}/submit`) is a future `Sx-J-*` row in this file; VA UI is **S2** (`slice-S2.md`). |
| `/portal/login` _(Inferred.)_ | public | coming-next | Agent | B2B portal session _(product TBD)_ | Agent portal session | F1 REST (`docs/03-services/00-iam-svc/01-api.md`) — not duplicated in `slice-S1` | Alternative F4 W2 path (“agent fills on behalf”). **Gap:** portal vs replica-only scope. |
| `/console/login` _(Inferred.)_ | internal | coming-next | Staff (CS, ops, finance, …) | n/a | F1 access + refresh | F1 — `POST /v1/sessions` (service doc) | Console entry; **not** in `slice-S1` REST. |
| `/console/bookings/new` _(Inferred.)_ | internal | coming-next | Staff with CS closing duty | `booking.create_on_behalf` (named in § Booking for `channel = cs`) | F1 bearer | § Booking — `POST /v1/bookings` with `channel = cs` | F4 W3 CS closing. **Gap:** internal navigation + which booking-svc read endpoints power the form (`GET/PATCH` draft — see below). |
| `/console/bookings/{booking_id}` _(Inferred.)_ | internal | coming-next | Same as above | same | F1 bearer | § Booking — draft create only today | Read/edit draft in UI needs **GET/PATCH** contract (gap). |

### Contract gaps to close (next Joint / slice rows)

1. **`GET /v1/bookings/{id}` + `PATCH /v1/bookings/{id}` (draft read/update)** — Listed as planned in `docs/03-services/02-booking-svc/01-api.md` but **not** yet in `slice-S1.md § Booking`. Add when a `S1-J-*` card (or follow-up) lands, so public/internal booking forms share one wire shape.
2. **`POST /v1/bookings/{id}/submit`** — F4 + Q006 gate; belongs in a **future** `slice-S1` § (or `slice-S2` boundary) once the Joint card is picked up. Drives `/checkout/{booking_id}` “pay” CTA behavior.
3. **B2B replica routing** — Freeze host + path + `ref` / `agent_id` propagation in gateway + `core-web` (table above cites partial § Booking rules only).
4. **Internal console shell** — Base path (`/console/...`), SSR vs SPA, and which routes require `CheckPermission` first (`S1-E-04`) — document when internal UX slice starts (may remain outside `slice-S1` if console is a separate app; if so, duplicate a slim matrix there).

---

## § S1-L-01 — UI screen inventory (BL-LGV-001)

**Proof of done (replaces wireframe):** canonical **URL → SvelteKit route files → primary components → UI states** for every **`active-now`** public row in **§ UI route matrix**, implemented under **`apps/core-web`**. **Figma:** not used; this table + shipped shells satisfy **05** `S1-L-01` (“Figma link *or* bullets in contract”).

**Catalog list URL:** **`/packages`** — product marketing landing stays **`/`**; the matrix row that inferred package list on `/` is implemented as **`/packages`** (equivalent to “may ship as `/catalog`”).

| Canonical URL | `apps/core-web` route files | Main components | Wire anchor | States in this shell |
| --- | --- | --- | --- | --- |
| `/packages` | `(b2c)/packages/+page.ts`, `(b2c)/packages/+page.svelte` | `PageShell`, stub list cards | § Catalog — `GET /v1/packages` | **Static stub** (loading/empty/error reserved for **S1-L-03**) |
| `/packages/{package_id}` | `(b2c)/packages/[package_id]/+page.ts`, `+page.svelte` | `PageShell`, CTA row | § Catalog — `GET /v1/packages/{id}` | Stub copy; **404** from gateway TBD **S1-L-03** |
| `/packages/{package_id}/departures/{departure_id}` | `(b2c)/packages/[package_id]/departures/[departure_id]/+page.ts`, `+page.svelte` | `PageShell`, link to booking | § Catalog — `GET /v1/package-departures/{id}` | Stub copy; seat labels in **S1-L-03** |
| `/booking/{package_id}` | `(b2c)/booking/[package_id]/+page.ts`, `+page.svelte` | `PageShell`, disabled form fields | § Booking — `POST /v1/bookings` | **Disabled** draft form; **submit** enabled in **S1-L-04** |

**Global chrome:** `Header.svelte` adds a **Packages** nav link to `/packages`. Landing **`/`** adds **Browse packages** CTA (`+page.svelte`).

**Route group `(b2c)`:** Does not appear in the URL; groups public catalog/booking trees for a future `apps/b2c` extract per **§ S1 UI placement**.

---

## § Changelog

- **2026-04-22** — Added `S1-E-01` engineering approval note under `§ Inventory` (`BL-EGV-001`): reviewed seat concurrency + transaction atomicity requirements for `ReserveSeats` / `ReleaseSeats`; gate approved with implementation deferred to `S1-E-02` / `S1-E-03`.
- **2026-04-22** — Added `§ S1-L-01 — UI screen inventory` + shipped route shells in `apps/core-web` — closes backlog **BL-LGV-001** / task **S1-L-01** per **05** (contract bullets + components; no Figma).
- **2026-04-21** — Added `§ S1 UI placement (core-web vs multi-app)` — records engineering consensus: S1 ships in `apps/core-web` with route-level separation until Q009 multi-app / `packages/api-client` split is justified; aligns `S1-L-01` proof with code paths without requiring `apps/b2c` to exist yet.
- **2026-04-21** — Added `§ UI route matrix` via `S0-L-01` / `BL-FE-PLN-001` — public vs internal URL table for S1, `active-now` vs `coming-next`, F1-aligned roles + permission column; **Contract gaps** list for draft GET/PATCH, submit, B2B routing, console shell.
- **2026-04-21** — Added `§ Booking States` via `S1-J-04` — documents that S1 exercises only the `draft` state; includes the forward-looking full state machine from F4 W4 (unchanged, reproduced for readers); pins the Q006 KTP+passport gate on the future `draft → pending_payment` transition so S2+ authors cannot soften it. **All four S1 contracts are now in; the contract-first gate for S1 code (`S1-E-02`, `S1-E-03`, `S1-E-04`) is satisfied.**
- **2026-04-21** — Added `§ Inventory` via `S1-J-03` — contracts `catalog.v1.CatalogService/ReserveSeats` + `ReleaseSeats` (gRPC): caller-supplied `reservation_id` for idempotency (option a), atomic-SQL pattern reference, five failure codes (`FAILED_PRECONDITION` / `NOT_FOUND` / `INVALID_ARGUMENT` / `ALREADY_EXISTS` / `INTERNAL`), and compensation prose covering the booking saga (ADR 0006) + the payment refund flow's Q004 conditional timing.
- **2026-04-21** — Added `§ Booking` via `S1-J-02` — contracts the `POST /v1/bookings` draft endpoint (three channels, idempotency, auth table, error codes, JSON examples, Q006 documented as submit-time not draft-time).
- **2026-04-20** — Initial version merged via `S1-J-01` (PR #7, commit `6c3fda8`). Adds `§ Catalog` with three public read endpoints. All other sections remain unfilled until their respective `S1-J-*` cards land.

