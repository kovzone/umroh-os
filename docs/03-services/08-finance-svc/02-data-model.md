# finance-svc — Data Model

## Tables (planned)

### `chart_of_accounts`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| code | text unique | numeric COA code |
| name | text | |
| kind | account_kind enum | asset / liability / equity / revenue / expense |
| parent_id | uuid null fk | tree |
| is_active | boolean | |
| created_at | timestamptz | |

### `journal_entries`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| code | text unique | |
| entry_date | date | |
| description | text | |
| source_kind | source_kind enum | manual / payment / logistics / crm / closing |
| source_id | text null | reference to originating entity |
| created_by | uuid | |
| created_at | timestamptz | |

### `journal_lines`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| journal_entry_id | uuid fk | |
| account_id | uuid fk → chart_of_accounts | |
| debit | numeric(15,2) default 0 | |
| credit | numeric(15,2) default 0 | |
| currency | text | |
| fx_rate | numeric(15,6) | snapshot |
| description | text | |
| job_order_id | uuid null | for job-order costing |

### `ar_balances`, `ap_balances`
Materialized views computed from journal_lines, refreshed periodically. Or alternatively maintained tables updated on each entry. TBD in implementation session.

### `tax_records`
| col | type |
|---|---|
| id | uuid pk |
| journal_entry_id | uuid fk |
| tax_kind | tax_kind enum |
| base_amount | numeric(15,2) |
| tax_amount | numeric(15,2) |
| period | date |

### `fx_rates`
| col | type |
|---|---|
| id | uuid pk |
| from_currency | text |
| to_currency | text |
| rate | numeric(15,6) |
| effective_date | date |
| unique(from_currency, to_currency, effective_date) |  |

### `job_order_costs`
| col | type |
|---|---|
| id | uuid pk |
| package_departure_id | uuid |
| revenue | numeric(15,2) |
| cost | numeric(15,2) |
| margin | numeric(15,2) |
| computed_at | timestamptz |

## Enums

```sql
CREATE TYPE account_kind AS ENUM ('asset', 'liability', 'equity', 'revenue', 'expense');
CREATE TYPE source_kind AS ENUM ('manual', 'payment', 'logistics', 'crm', 'closing');
CREATE TYPE tax_kind AS ENUM ('pph21', 'pph23', 'ppn');
```

## Notes

- Double-entry enforcement: a check constraint or trigger on `journal_entries` ensures the sum of `journal_lines.debit` equals `journal_lines.credit` per entry.
- All amounts stored in original currency plus a snapshot `fx_rate` so reports can re-aggregate to base currency.
