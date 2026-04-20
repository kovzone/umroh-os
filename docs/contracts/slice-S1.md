---
slice: S1
title: Slice S1 — Integration Contract
status: draft
last_updated: 2026-04-20
pr_owner: Elda
reviewer: Elda (solo-exec S0 per § Current operating mode)
task_codes:
  - S1-J-01
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

_(Landed via `S1-J-01`, 2026-04-20.)_

Public read-only endpoints for B2C catalog browsing. All endpoints are unauthenticated — the catalog is a public surface. Auth'd admin/staff catalog endpoints land in a later slice, not here.

### Endpoints

| Method | Path | Auth | Purpose |
|---|---|---|---|
| `GET` | `/v1/packages` | public | List active packages, filterable. |
| `GET` | `/v1/packages/{id}` | public | Package detail incl. itinerary + departure summary list + master-data refs. |
| `GET` | `/v1/package-departures/{id}` | public | Departure detail incl. `remaining_seats`, pricing per room type, vendor-readiness summary. |

Paths are stable and match `docs/03-services/01-catalog-svc/01-api.md` row-for-row for the public-read subset. Route prefix `/v1/` is mounted under `catalog-svc`'s REST adapter (port 4002 in dev, proxied by `gateway-svc` at 4000 for public traffic).

### Conventions honored

- **Q001 (currency).** Responses carry `list_amount` (integer), `list_currency` (`IDR` or `USD`), and `settlement_currency` (always `"IDR"` in MVP). Catalog never commits to a payable amount — that lock happens in `payment-svc` at invoice creation with a snapshotted FX rate. Consumers must treat `list_amount` + `list_currency` as display-only.
- **Q003 (single-language MVP).** All human-readable strings (`name`, `description`, `itinerary` day labels) are returned in Bahasa Indonesia (`id-ID`). No `{lang}` query param, no `translations` array in the response — a future multi-language card will evolve the contract via Bump-versi.
- **No auth** on any endpoint in this section. `gateway-svc` passes these through without token validation.

### `GET /v1/packages`

**Query params:**

| Param | Type | Required | Notes |
|---|---|---|---|
| `kind` | string enum | no | One of `umrah_reguler`, `umrah_plus`, `hajj_furoda`, `hajj_khusus`, `badal`, `financial`, `retail`. Omit to list all kinds. |
| `departure_from` | string (ISO-8601 date) | no | Lower bound on departure date (earliest departure ≥ this date). |
| `departure_to` | string (ISO-8601 date) | no | Upper bound on departure date. |
| `airline_code` | string | no | IATA airline code filter. |
| `hotel_id` | string (ULID) | no | Filter to packages referencing this hotel. |
| `cursor` | string (opaque) | no | Pagination cursor from previous response's `page.next_cursor`. |
| `limit` | integer | no | Page size, 1–100, default 20. |

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

| Status | Body `error.code` | When |
|---|---|---|
| `400` | `invalid_query_param` | Any query param fails validation (bad date, unknown `kind`, `limit` out of range). |
| `400` | `invalid_cursor` | `cursor` is malformed or from an incompatible list. |
| `500` | `internal_error` | Unexpected server error. |

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

| Status | Body `error.code` | When |
|---|---|---|
| `404` | `package_not_found` | No package matches `{id}`, or the package is not `status = 'active'` (the public endpoint does not leak existence of draft/archived packages). |
| `500` | `internal_error` | Unexpected server error. |

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

| Status | Body `error.code` | When |
|---|---|---|
| `404` | `departure_not_found` | No departure matches `{id}`, or status is `departed` / `completed` / `cancelled`. |
| `500` | `internal_error` | Unexpected server error. |

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

## § Changelog

- **2026-04-20** — Initial version merged via `S1-J-01` (PR pending). Adds `§ Catalog` with three public read endpoints. All other sections (`§ Booking`, `§ Inventory`, `§ Booking States`) remain unfilled until their respective `S1-J-*` cards land.
