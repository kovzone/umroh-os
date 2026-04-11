# iam-svc — API

## REST endpoints (planned)

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/v1/sessions` | Login — exchange credentials for token |
| `POST` | `/v1/sessions/refresh` | Refresh access token |
| `DELETE` | `/v1/sessions` | Logout — revoke current session |
| `GET` | `/v1/me` | Current user profile |
| `POST` | `/v1/me/2fa/enroll` | Enroll TOTP |
| `POST` | `/v1/me/2fa/verify` | Verify TOTP code |
| `GET` | `/v1/users` | List users (admin) |
| `POST` | `/v1/users` | Create user (admin) |
| `GET` | `/v1/users/{id}` | Get user (admin) |
| `PATCH` | `/v1/users/{id}` | Update user (admin) |
| `DELETE` | `/v1/users/{id}` | Soft-delete user (admin) |
| `POST` | `/v1/users/{id}/roles` | Assign role to user |
| `DELETE` | `/v1/users/{id}/roles/{role_id}` | Revoke role |
| `GET` | `/v1/roles` | List roles |
| `POST` | `/v1/roles` | Create role |
| `GET` | `/v1/permissions` | List all permissions (catalog) |
| `GET` | `/v1/branches` | List branches |
| `POST` | `/v1/branches` | Create branch |
| `GET` | `/v1/audit-logs` | Query audit log |

## gRPC methods (planned)

`IamService`:
- `ValidateToken(ValidateTokenRequest) → ValidateTokenResponse` — every service calls this on every authenticated request
- `CheckPermission(CheckPermissionRequest) → CheckPermissionResponse` — resource + action + scope
- `GetUser(GetUserRequest) → GetUserResponse` — read-only user lookup
- `RecordAudit(RecordAuditRequest) → RecordAuditResponse` — append to audit log from any service

> All endpoints are stubs. Spec lives in `iam-svc/api/rest_oapi/openapi.yaml` once scaffolded.
