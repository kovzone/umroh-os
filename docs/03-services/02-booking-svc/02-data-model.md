# booking-svc — Data Model

## Tables (planned)

### `bookings`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| code | text unique | human-friendly booking number |
| package_departure_id | uuid | reference to catalog (no FK across services) |
| branch_id | uuid | scope |
| agent_id | uuid null | reference to crm |
| created_by | uuid | reference to iam.users |
| status | booking_status enum | |
| total_amount | numeric(15,2) | contractual **IDR** total after pricing lock (matches invoice `amount_total`; **Q001**) |
| currency | char(3) | MVP: always `IDR` for payable totals |
| paid_amount | numeric(15,2) default 0 | denormalized; **must equal** primary invoice `paid_amount` after each payment event (single source of truth remains `payment-svc`) |
| fx_snapshot | jsonb | copy of invoice / issuance FX lock for offline reads; set when VA issued (**F4** / **Q001**) |
| notes | text | |
| created_at, updated_at, cancelled_at | timestamptz | |

### `booking_items`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| booking_id | uuid fk | |
| jamaah_id | uuid | reference to jamaah |
| room_type | room_type enum | mirror of catalog enum |
| price | numeric(15,2) | **locked IDR** line amount at submit (catalog list × locked FX + room split); not re-priced on FX drift |
| addons | jsonb | snapshot of selected addons |
| status | booking_item_status enum | active / cancelled |
| created_at | timestamptz | |

### `room_allocations`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| booking_id | uuid fk | |
| hotel_id | uuid | reference to catalog |
| room_number | text | |
| occupants | uuid[] | array of jamaah_id |
| created_at | timestamptz | |

### `bus_allocations`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| package_departure_id | uuid | |
| bus_number | text | |
| seat_assignments | jsonb | { seat: jamaah_id } |
| created_at | timestamptz | |

### `booking_status_history`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| booking_id | uuid fk | |
| from_status | booking_status null | |
| to_status | booking_status | |
| changed_by | uuid null | |
| reason | text | |
| created_at | timestamptz | append-only |

## Enums

```sql
CREATE TYPE booking_status AS ENUM (
    'draft', 'pending_payment', 'partially_paid', 'paid_in_full', 'cancelled', 'completed'
);
CREATE TYPE booking_item_status AS ENUM ('active', 'cancelled');
```

## Notes

- Cross-service IDs (`package_departure_id`, `jamaah_id`, `hotel_id`) are stored as plain UUIDs without DB-level FKs. Referential integrity is the service layer's job, validated via gRPC reads at write time.
- `booking_status_history` is append-only — write a new row on every transition.
- **FX:** `fx_snapshot` + `total_amount`/`currency` align with `payment-svc.invoices` at saga completion; catalog list currencies are **not** duplicated here beyond line snapshots.
