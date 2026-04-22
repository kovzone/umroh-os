# catalog-svc — Data Model

## ID convention

Catalog entity IDs are **ULID strings with a type prefix** per `docs/contracts/slice-S1.md § Catalog` "Response ID format". This deviates from the repo-wide default (`uuid` + `gen_random_uuid()`) set in `docs/04-backend-conventions/04-database-and-sqlc.md` — the contract-first rule wins. Implementers generate the ULID app-side and store it as `TEXT PRIMARY KEY`. Prefixes in use:

| Prefix | Table |
|---|---|
| `pkg_` | `packages` |
| `dep_` | `package_departures` |
| `pkgpr_` | `package_pricing` |
| `itn_` | `itinerary_templates` |
| `htl_` | `hotels` |
| `arl_` | `airlines` |
| `mtw_` | `muthawwif` |
| `addon_` | `addons` |

A `CHECK (id LIKE '<prefix>_%')` constraint is applied on each table; consumers must treat IDs as opaque strings (no parsing) per the contract.

## Tables

Landed by migration `000008_add_catalog_packages_and_masters` (S1-E-02 / BL-CAT-001). All objects live in the `catalog` schema per ADR 0007.

### `packages`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `pkg_`-prefixed |
| kind | catalog.package_kind | 7-value enum (see Enums) |
| name | text | |
| description | text | default `''` |
| highlights | text[] | default `{}` |
| cover_photo_url | text | default `''` |
| itinerary_id | text fk → catalog.itinerary_templates.id | nullable |
| airline_id | text fk → catalog.airlines.id | nullable |
| muthawwif_id | text fk → catalog.muthawwif.id | nullable |
| status | catalog.package_status | default `'draft'` |
| created_at, updated_at | timestamptz | |
| deleted_at | timestamptz | soft-delete nullable |

Indexes:
- `packages_status_kind_idx (status, kind) WHERE deleted_at IS NULL`
- `packages_active_id_idx (id) WHERE status='active' AND deleted_at IS NULL` — hot path for cursor pagination
- `packages_airline_id_idx (airline_id) WHERE deleted_at IS NULL`

### `package_departures`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `dep_`-prefixed |
| package_id | text fk → catalog.packages.id ON DELETE CASCADE | |
| departure_date | date | |
| return_date | date | CHECK `return_date >= departure_date` |
| total_seats | int | CHECK `> 0` |
| reserved_seats | int | CHECK `>= 0`; default 0 |
| status | catalog.departure_status | default `'open'` |
| created_at, updated_at | timestamptz | |

Constraints:
- `CHECK (reserved_seats <= total_seats)` — defence-in-depth; the atomic SQL guard on `ReserveSeats` is the primary safety (lands via S1-E-03).

### `package_pricing`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `pkgpr_`-prefixed |
| package_departure_id | text fk → catalog.package_departures.id ON DELETE CASCADE | |
| room_type | catalog.room_type | `double`/`triple`/`quad` |
| list_amount | numeric(18, 4) | display list in `list_currency`; extra precision for USD/SAR inputs. Wire layer casts to BIGINT. |
| list_currency | char(3) | ISO 4217; IDR or USD |
| settlement_currency | char(3) | CHECK `= 'IDR'` per **Q001** |

Unique: `(package_departure_id, room_type)`.

### `hotels`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `htl_`-prefixed |
| name | text | |
| city | text | `mecca`/`medina`/etc. (free-form) |
| star_rating | smallint | CHECK `BETWEEN 0 AND 5` |
| walking_distance_m | int | |

### `airlines`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `arl_`-prefixed |
| code | text unique | IATA code |
| name | text | |
| operator_kind | catalog.operator_kind | `airline`/`rail`/`bus` (Haramain HSR modelled here) |

### `muthawwif`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `mtw_`-prefixed |
| name | text | |
| portrait_url | text | |

### `itinerary_templates`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `itn_`-prefixed |
| name | text | default `''` |
| days | jsonb | array of `{day, title, description, photo_url?}` |
| public_url | text | shareable micro-web URL |

### `addons`
| col | type | notes |
|---|---|---|
| id | text pk | ULID `addon_`-prefixed |
| name | text | |
| list_amount | numeric(18, 4) | same semantics as package_pricing |
| list_currency | char(3) | |
| settlement_currency | char(3) | CHECK `= 'IDR'` |

### `package_hotels` (join)
| col | type |
|---|---|
| package_id | text fk → catalog.packages.id ON DELETE CASCADE |
| hotel_id | text fk → catalog.hotels.id ON DELETE RESTRICT |
| sort_order | smallint |
| pk(package_id, hotel_id) |  |

### `package_addons` (join)
| col | type |
|---|---|
| package_id | text fk → catalog.packages.id ON DELETE CASCADE |
| addon_id | text fk → catalog.addons.id ON DELETE RESTRICT |
| pk(package_id, addon_id) |  |

## Enums

```sql
CREATE TYPE catalog.package_kind AS ENUM (
    'umrah_reguler', 'umrah_plus',
    'hajj_furoda', 'hajj_khusus',
    'badal', 'financial', 'retail'
);
CREATE TYPE catalog.package_status AS ENUM ('draft', 'active', 'archived');
CREATE TYPE catalog.departure_status AS ENUM (
    'open', 'closed', 'departed', 'completed', 'cancelled'
);
CREATE TYPE catalog.room_type AS ENUM ('double', 'triple', 'quad');
CREATE TYPE catalog.operator_kind AS ENUM ('airline', 'rail', 'bus');
```

Public read endpoints surface only `departure_status ∈ {open, closed}` — the other three are hidden server-side and return 404 to public callers. `package_status = 'active'` is the only status visible on `/v1/packages*`.

## Notes

- **Atomic seat reservation** uses `UPDATE catalog.package_departures SET reserved_seats = reserved_seats + $1 WHERE id = $2 AND reserved_seats + $1 <= total_seats` (single statement, zero rows returned ⇒ `insufficient_capacity`). Implementation lands in S1-E-03 (booking saga) / the ReserveSeats gRPC handler.
- Photos and videos are stored in GCS; the DB only holds URLs.
- Bulk import (CSV) goes through a staging table to allow validation before commit — deferred to BL-CAT-010.
- **List vs settlement (F2, Q001):** catalog stores **list** amounts + currency for B2C/B2B display. **No** customer invoice totals here — `payment-svc` / `booking-svc` lock **IDR** payable at VA issuance using `fx_snapshot` (see payment-svc data model). Rounding to nearest **Rp 1,000** (half-up) applies once on that payable IDR total, not per catalog row.
- **Bahasa-only MVP (Q003):** `name`, `description`, itinerary day labels, highlights — all `id-ID`. No `translations` column in MVP schema.
