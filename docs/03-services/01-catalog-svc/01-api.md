# catalog-svc — API

Per ADR 0009 catalog-svc is **gRPC-only**. All client-facing REST paths live on `gateway-svc`; the gateway proxies to the RPCs below via `gateway-svc/adapter/catalog_grpc_adapter`. The legacy catalog REST package (`api/rest_oapi/`) was removed in `BL-REFACTOR-001 / S1-E-11` on 2026-04-23 — there is no longer a backend REST surface.

## gRPC methods (CatalogService)

Defined in `services/catalog-svc/api/grpc_api/pb/catalog.proto`.

| RPC | Status | Caller | Purpose |
|---|---|---|---|
| `Healthz` | live | reflection probes (legacy) | Pilot placeholder; real infra probes use the standard `grpc.health.v1.Health` protocol registered in `cmd/server.go`. |
| `ListPackages` | **live** (BL-GTW-002 / S1-E-10) | `gateway-svc` → public `GET /v1/packages` | Filter by `kind`, `departure_from/to`, `airline_code`, `hotel_id`; cursor-paginated. |
| `GetPackage` | **live** (BL-GTW-002 / S1-E-10) | `gateway-svc` → public `GET /v1/packages/{id}` | Active package detail with eager master refs + upcoming open/closed departures. |
| `GetPackageDeparture` | **live** (BL-GTW-002 / S1-E-10) | `gateway-svc` → public `GET /v1/package-departures/{id}` | Departure detail with live `remaining_seats`, pricing per room type, and `vendor_readiness`. |
| `DiagnosticsDbTx` | **live** (BL-REFACTOR-001 / S1-E-11) | dev-only (invoke via `grpcurl`) | State-mutating WithTx reference path. No gateway REST route — the equivalent cross-service trace demo is `/v1/iam/system/diagnostics/db-tx` on iam-svc (S0-J-05). |
| `ReserveSeats` | planned (S1-E-03) | `booking-svc` | Atomic reservation, called by the booking saga. |
| `ReleaseSeats` | planned (S1-E-03) | `booking-svc` | Compensating action. |
| `GetHotel` / `GetAirline` / `GetMuthawwif` | planned | internal | Master-data lookups. |

## Public REST paths (on gateway-svc, not here)

For completeness — these are served by `gateway-svc:4000` and proxy to the RPCs above:

| Method | Path | Backing RPC |
|---|---|---|
| `GET` | `/v1/packages` | `ListPackages` |
| `GET` | `/v1/packages/{id}` | `GetPackage` |
| `GET` | `/v1/package-departures/{id}` | `GetPackageDeparture` |

Admin / staff-write paths (create / patch / delete packages, hotels, airlines, muthawwif, bulk CSV import/export) are not live yet — they land with `BL-CAT-014` and the bulk cards. When added, they ship as paired gRPC methods on this service + gateway REST routes per ADR 0009.

## Payload alignment

Departure + pricing responses expose `list_amount`, `list_currency`, and `settlement_currency` per `02-data-model.md` (F2 / **Q001**). Consumers compute **IDR** payable only at booking/invoice lock in `payment-svc`.

## Contract

Wire shapes are pinned in `docs/contracts/slice-S1.md § Catalog` + `§ Gateway` — the contract is the source of truth; this file is a summary for navigation.
