# iam-svc — Data Model

## Tables (planned)

### `branches`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| name | text | |
| code | text unique | short code |
| parent_id | uuid null fk | for hierarchical branches |
| created_at, updated_at | timestamptz | |

### `users`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| email | text unique | |
| password_hash | text | bcrypt |
| name | text | |
| branch_id | uuid fk → branches | |
| status | user_status enum | active / suspended / pending |
| totp_secret | text null | encrypted at rest |
| totp_verified_at | timestamptz null | |
| last_login_at | timestamptz null | |
| created_at, updated_at, deleted_at | timestamptz | |

### `roles`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| name | text unique | e.g. `ops_admin`, `branch_manager` |
| description | text | |
| created_at, updated_at | timestamptz | |

### `permissions`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| resource | text | e.g. `booking`, `package` |
| action | text | `read` / `write` / `edit` / `delete` |
| scope | permission_scope enum | `global` / `branch` / `personal` |
| unique(resource, action, scope) |  |  |

### `role_permissions`
| col | type |
|---|---|
| role_id | uuid fk |
| permission_id | uuid fk |
| pk(role_id, permission_id) |  |

### `user_roles`
| col | type |
|---|---|
| user_id | uuid fk |
| role_id | uuid fk |
| pk(user_id, role_id) |  |

### `sessions`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| user_id | uuid fk | |
| refresh_token_hash | text | |
| user_agent | text | |
| ip | inet | |
| issued_at, expires_at | timestamptz | |
| revoked_at | timestamptz null | |

### `audit_logs`
| col | type | notes |
|---|---|---|
| id | uuid pk | |
| user_id | uuid fk null | nullable for system actions |
| branch_id | uuid fk null | |
| resource | text | what was acted on |
| resource_id | text | |
| action | text | created / updated / deleted / viewed |
| old_value | jsonb null | |
| new_value | jsonb null | |
| ip | inet null | |
| created_at | timestamptz | append-only; no UPDATE |

## Enums

```sql
CREATE TYPE user_status AS ENUM ('active', 'suspended', 'pending');
CREATE TYPE permission_scope AS ENUM ('global', 'branch', 'personal');
```

## Relationships

```
branches ──< users ──< user_roles >── roles ──< role_permissions >── permissions
                  └─< sessions
                  └─< audit_logs
```

## Notes

- `audit_logs` is append-only. Enforce via DB role permissions or trigger.
- `password_hash` uses bcrypt cost 12.
- `totp_secret` is encrypted at the application layer before insert (key from config).
