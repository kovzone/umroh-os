# catalog-svc — Data Model

## Tables (planned)

### `packages`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| code | text unique | sellable SKU |
| name | text | |
| kind | package_kind enum | umrah / hajj / badal |
| duration_days | int | |
| description | text | |
| highlights | text[] | |
| cover_image_url | text | GCS signed URL |
| status | package_status enum | draft / active / archived |
| created_at, updated_at, deleted_at | timestamptz | |

### `package_departures`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| package_id | uuid fk | |
| departure_date | date | |
| return_date | date | |
| airline_id | uuid fk | |
| muthawwif_id | uuid fk | |
| total_seats | int | |
| reserved_seats | int | atomic counter |
| status | departure_status enum | open / closed / departed / completed / cancelled |
| created_at, updated_at | timestamptz | |

### `package_pricing`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| package_departure_id | uuid fk | |
| room_type | room_type enum | double / triple / quad |
| list_amount | numeric(18,4) | commercial list figure in `list_currency` (extra precision for USD/SAR inputs) |
| list_currency | char(3) | ISO 4217, e.g. `IDR` or `USD` — **display / quote** basis (F2 / **Q001**) |
| settlement_currency | char(3) | MVP: always `IDR` (CHECK); contractual invoice/VA amounts are always IDR at booking lock |

### `hotels`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| name | text | |
| city | text | mecca / medina / etc |
| star_rating | smallint | |
| distance_to_mosque_meters | int | |
| photos | jsonb | array of urls |
| tour_360_url | text null | |
| address | text | |
| created_at, updated_at | timestamptz | |

### `airlines`
| col | type |
|---|---|
| id | uuid pk |
| code | text unique |
| name | text |
| logo_url | text |

### `muthawwif`
| col | type |
|---|---|
| id | uuid pk |
| name | text |
| photo_url | text |
| bio | text |
| video_url | text null |
| languages | text[] |

### `itinerary_templates`
| col | type |
|---|---|
| id | uuid pk |
| name | text |
| days | jsonb |

### `addons`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| name | text | |
| list_amount | numeric(18,4) | same semantics as `package_pricing.list_amount` |
| list_currency | char(3) | ISO 4217; default `IDR` if not multi-currency |

### `package_addons`
| col | type |
|---|---|
| package_id | uuid fk |
| addon_id | uuid fk |
| pk(package_id, addon_id) |  |

## Enums

```sql
CREATE TYPE package_kind AS ENUM ('umrah', 'hajj', 'badal');
CREATE TYPE package_status AS ENUM ('draft', 'active', 'archived');
CREATE TYPE departure_status AS ENUM ('open', 'closed', 'departed', 'completed', 'cancelled');
CREATE TYPE room_type AS ENUM ('double', 'triple', 'quad');
```

## Notes

- Seat reservation is atomic via `UPDATE package_departures SET reserved_seats = reserved_seats + $1 WHERE id = $2 AND reserved_seats + $1 <= total_seats`.
- Photos and videos are stored in GCS; the DB only holds URLs.
- Bulk import goes through a staging table to allow validation before commit.
- **List vs settlement (F2, Q001):** catalog stores **list** amounts + currency for B2C/B2B display. **No** customer invoice totals here — `payment-svc` / `booking-svc` lock **IDR** payable at VA issuance using `fx_snapshot` (see payment-svc data model). Rounding to nearest **Rp 1,000** (half-up) applies once on that payable IDR total, not per catalog row.
