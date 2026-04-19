# catalog-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/v1/packages` | List packages (filterable) |
| `POST` | `/v1/packages` | Create package |
| `GET` | `/v1/packages/{id}` | Get package detail |
| `PATCH` | `/v1/packages/{id}` | Update package |
| `DELETE` | `/v1/packages/{id}` | Soft-delete package |
| `GET` | `/v1/packages/{id}/departures` | List departures for package |
| `POST` | `/v1/packages/{id}/departures` | Create departure |
| `GET` | `/v1/package-departures/{id}` | Get departure detail |
| `PATCH` | `/v1/package-departures/{id}` | Update departure |
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

> All endpoints are stubs. Will be defined when catalog-svc is scaffolded.

**Payload alignment:** departure + pricing responses expose `list_amount`, `list_currency`, and `settlement_currency` per `02-data-model.md` (F2 / **Q001**). Consumers compute **IDR** payable only at booking/invoice lock in `payment-svc`.
