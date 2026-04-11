# visa-svc — Data Model

## Tables (planned)

### `visa_applications`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| jamaah_id | uuid | reference |
| booking_id | uuid | reference |
| status | visa_status enum | |
| provider | visa_provider enum | mofa / sajil |
| provider_ref | text null | external reference |
| submitted_at | timestamptz null | |
| issued_at | timestamptz null | |
| created_at, updated_at | timestamptz | |

### `visa_status_history`
| col | type |
|---|---|
| id | uuid pk |
| visa_application_id | uuid fk |
| from_status | visa_status null |
| to_status | visa_status |
| reason | text |
| created_at | timestamptz |

### `e_visas`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| visa_application_id | uuid fk | |
| storage_path | text | gs:// |
| visa_number | text | |
| valid_from | date | |
| valid_to | date | |

### `tasreh_records`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| visa_application_id | uuid fk | |
| kind | tasreh_kind enum | raudhah / other |
| issued_for | timestamptz | scheduled entry time |
| storage_path | text | |
| used_at | timestamptz null | scan time at entry |

### `raudhah_monitoring`
| col | type |
|---|---|
| id | uuid pk |
| visa_application_id | uuid fk |
| snapshot | jsonb |
| polled_at | timestamptz |

## Enums

```sql
CREATE TYPE visa_status AS ENUM ('waiting_docs', 'docs_ready', 'submitted', 'issued', 'rejected');
CREATE TYPE visa_provider AS ENUM ('mofa', 'sajil');
CREATE TYPE tasreh_kind AS ENUM ('raudhah', 'other');
```
