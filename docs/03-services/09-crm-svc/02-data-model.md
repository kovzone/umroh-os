# crm-svc — Data Model

## Tables (planned)

### `leads`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| name | text | |
| phone | text | |
| email | text null | |
| source | text | utm_source |
| campaign_id | uuid null fk | |
| agent_id | uuid null fk | |
| status | lead_status enum | cold / warm / hot / converted / lost |
| converted_booking_id | uuid null | |
| created_at, updated_at | timestamptz | |

### `campaigns`
| col | type |
|---|---|
| id | uuid pk |
| name | text |
| channel | text |
| start_date, end_date | date |
| budget | numeric(15,2) |
| utm_source | text |
| utm_medium | text |
| utm_campaign | text |
| status | campaign_status enum |

### `agents`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| user_id | uuid | reference to iam.users |
| code | text unique | |
| level | agent_level enum | silver / gold / platinum |
| parent_agent_id | uuid null fk | hierarchy |
| branch_id | uuid | |
| commission_pct | numeric(5,2) | |
| status | agent_status enum | active / suspended |

### `commission_ledger`
| col | type |
|---|---|
| id | uuid pk |
| agent_id | uuid fk |
| booking_id | uuid |
| amount | numeric(15,2) |
| kind | commission_kind enum |
| status | commission_status enum |
| paid_at | timestamptz null |
| created_at | timestamptz |

### `broadcasts`
| col | type |
|---|---|
| id | uuid pk |
| channel | broadcast_channel enum |
| segment | jsonb |
| body | text |
| sent_count | int |
| failed_count | int |
| status | broadcast_status enum |
| created_at, completed_at | timestamptz |

### `alumni_threads`, `referral_codes`, `ziswaf_transactions`

Additional tables — schemas TBD when this area is implemented.

## Enums

```sql
CREATE TYPE lead_status AS ENUM ('cold', 'warm', 'hot', 'converted', 'lost');
CREATE TYPE campaign_status AS ENUM ('draft', 'active', 'paused', 'ended');
CREATE TYPE agent_level AS ENUM ('silver', 'gold', 'platinum');
CREATE TYPE agent_status AS ENUM ('active', 'suspended');
CREATE TYPE commission_kind AS ENUM ('direct', 'override');
CREATE TYPE commission_status AS ENUM ('pending', 'approved', 'paid', 'cancelled');
CREATE TYPE broadcast_channel AS ENUM ('whatsapp', 'email', 'sms');
CREATE TYPE broadcast_status AS ENUM ('queued', 'sending', 'completed', 'failed');
```
