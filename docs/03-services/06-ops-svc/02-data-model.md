# ops-svc — Data Model

## Tables (planned)

### `verification_tasks`
| col | type |
|---|---|
| id | uuid pk |
| document_id | uuid (ref) |
| jamaah_id | uuid |
| status | verification_status enum |
| reviewer_id | uuid null |
| reason | text |
| created_at, updated_at | timestamptz |

### `manifests`
| col | type |
|---|---|
| id | uuid pk |
| package_departure_id | uuid |
| format | manifest_format enum |
| storage_path | text |
| jamaah_count | int |
| generated_by | uuid |
| created_at | timestamptz |

### `luggage_tags`
| col | type |
|---|---|
| id | uuid pk |
| booking_id | uuid |
| jamaah_id | uuid |
| tag_code | text unique |
| qr_payload | text |
| issued_at | timestamptz |

### `handling_events`
| col | type |
|---|---|
| id | uuid pk |
| jamaah_id | uuid |
| booking_id | uuid |
| event_type | handling_event_type enum |
| location | text |
| device_id | text |
| metadata | jsonb |
| occurred_at | timestamptz |

### `grouping_runs`
| col | type |
|---|---|
| id | uuid pk |
| package_departure_id | uuid |
| input_snapshot | jsonb |
| output | jsonb |
| ran_at | timestamptz |

## Enums

```sql
CREATE TYPE verification_status AS ENUM ('pending', 'in_review', 'approved', 'rejected');
CREATE TYPE manifest_format AS ENUM ('pdf', 'xlsx', 'csv');
CREATE TYPE handling_event_type AS ENUM ('luggage_in', 'luggage_out', 'boarding', 'arrival', 'zamzam_distribution');
```
