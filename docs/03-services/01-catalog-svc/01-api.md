# catalog-svc — API

## Current status (S1-E-02 / BL-CAT-001)

The two public read endpoints below are **live** (merged 2026-04-22). All other endpoints are planned/stub and land via later cards (`BL-CAT-002` for departure detail, `BL-CAT-005..011` for admin CRUD + bulk, `S1-E-03` for the seat RPCs).

| Method | Path | Status | Purpose |
|---|---|---|---|
| `GET` | `/v1/packages` | **live** (BL-CAT-001) | List active packages. Filter by `kind`, `departure_from/to`, `airline_code`, `hotel_id`; cursor-paginated. |
| `GET` | `/v1/packages/{id}` | **live** (BL-CAT-001) | Active package detail with eager master refs + upcoming open/closed departures. |
| `GET` | `/v1/package-departures/{id}` | planned (BL-CAT-002) | Departure detail with live `remaining_seats` + `vendor_readiness`. |
| `GET` | `/v1/hotels` | planned | List hotels (admin). |
| … | | | |

Wire shapes are pinned in `docs/contracts/slice-S1.md § Catalog` — the contract is the source of truth; this file is a summary for navigation.

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/packages` | List packages (filterable) — **public** (no Bearer) |
| `POST` | `/v1/packages` | Create package — **staff only** (Bearer + `catalog.package.manage`; see `slice-S1.md` § Catalog — internal write MVP / backlog `BL-CAT-014`) |
| `GET` | `/v1/packages/{id}` | Get package detail — **public** |
| `PATCH` | `/v1/packages/{id}` | Update package — **staff only** |
| `DELETE` | `/v1/packages/{id}` | Soft-delete package — **staff only** |
| `GET` | `/v1/packages/{id}/departures` | List departures for package |
| `POST` | `/v1/packages/{id}/departures` | Create departure — **staff only** |
| `GET` | `/v1/package-departures/{id}` | Get departure detail — **public** |
| `PATCH` | `/v1/package-departures/{id}` | Update departure — **staff only** |
| `GET` | `/v1/hotels` | List hotels |
| `POST` | `/v1/hotels` | Create hotel |
| `GET` | `/v1/hotels/{id}` | Get hotel detail |
| `PATCH` | `/v1/hotels/{id}` | Update hotel |
| `GET` | `/v1/airlines` | List airlines |
| `POST` | `/v1/airlines` | Create airline |
| `GET` | `/v1/muthawwif` | List muthawwif |
| `POST` | `/v1/muthawwif` | Create muthawwif |
| `POST` | `/v1/packages/import` | Bulk import via CSV |
| `GET` | `/v1/packages/export` | Bulk export to CSV |

## gRPC methods (planned)

`CatalogService`:
- `GetPackage(GetPackageRequest) → GetPackageResponse`
- `GetPackageDeparture(...)` — used by booking-svc and ops-svc
- `ReserveSeats(ReserveSeatsRequest) → ReserveSeatsResponse` — atomic reservation, called by booking saga
- `ReleaseSeats(ReleaseSeatsRequest) → ReleaseSeatsResponse` — compensating action
- `GetHotel(...)`
- `GetAirline(...)`
- `GetMuthawwif(...)`

**Payload alignment:** departure + pricing responses expose `list_amount`, `list_currency`, and `settlement_currency` per `02-data-model.md` (F2 / **Q001**). Consumers compute **IDR** payable only at booking/invoice lock in `payment-svc`.
