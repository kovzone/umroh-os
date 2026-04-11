# jamaah-svc — Data Model

## Tables (planned)

### `jamaah`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| family_unit_id | uuid null fk | |
| full_name | text | |
| birth_date | date | |
| gender | gender enum | |
| nik | text unique | KTP national ID |
| phone | text | |
| email | text null | |
| address | text | |
| status | jamaah_status enum | calon / active / alumni |
| branch_id | uuid | |
| created_at, updated_at, deleted_at | timestamptz | |

### `family_units`
| col | type |
|---|---|
| id | uuid pk |
| code | text unique (K-Family Code) |
| name | text |
| created_at | timestamptz |

### `mahram_relations`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| jamaah_id | uuid fk | |
| related_jamaah_id | uuid fk | |
| relation | mahram_relation enum | husband / father / brother / son / etc |
| verified | boolean | |
| created_at | timestamptz | |

### `documents`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| jamaah_id | uuid fk | |
| kind | document_kind enum | ktp / passport / vaccine / family_book |
| storage_path | text | gs://bucket/key |
| status | document_status enum | uploaded / processing / verified / rejected |
| uploaded_by | uuid | |
| created_at, updated_at | timestamptz | |

### `ocr_results`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| document_id | uuid fk | |
| extracted | jsonb | structured fields (mrz, nik, name, dob, ...) |
| confidence | numeric(3,2) | |
| created_at | timestamptz | |

## Enums

```sql
CREATE TYPE gender AS ENUM ('male', 'female');
CREATE TYPE jamaah_status AS ENUM ('calon', 'active', 'alumni');
CREATE TYPE mahram_relation AS ENUM ('husband', 'father', 'son', 'brother', 'grandfather', 'grandson', 'uncle', 'nephew', 'father_in_law', 'son_in_law', 'mother', 'daughter', 'sister', 'wife');
CREATE TYPE document_kind AS ENUM ('ktp', 'passport', 'vaccine', 'family_book', 'photo');
CREATE TYPE document_status AS ENUM ('uploaded', 'processing', 'verified', 'rejected');
```

## Notes

- Mahram check uses a recursive CTE on `mahram_relations` to walk the graph.
- `ocr_results.extracted` is JSONB so the schema can evolve as new document types arrive.
