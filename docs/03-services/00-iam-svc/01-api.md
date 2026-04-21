# iam-svc — API

## REST endpoints

### Landed (BL-IAM-001, 2026-04-21)

| Method | Path | Auth | Purpose |
|---|---|---|---|
| `POST` | `/v1/sessions` | public | Login — exchange email + password for access + refresh token pair |
| `POST` | `/v1/sessions/refresh` | public | Rotate refresh — revoke the old row, issue a new pair |
| `DELETE` | `/v1/sessions` | bearer | Logout — revoke the current session row (by session-id claim) |
| `GET` | `/v1/me` | bearer | Current user profile + TOTP enrollment state |
| `POST` | `/v1/me/2fa/enroll` | bearer | Start TOTP enrollment (generate + persist encrypted secret) |
| `POST` | `/v1/me/2fa/verify` | bearer | Verify TOTP code and stamp `users.totp_verified_at` |

OpenAPI source of truth: `services/iam-svc/api/rest_oapi/openapi.yaml`.
Bearer enforcement: middleware, not spec — see `api/rest_oapi/middleware/bearer_auth.go`.

### Planned (sibling + depth cards)

| Method | Path | Ships with |
|---|---|---|
| `GET` | `/v1/users` | `S1-E-06` (depth) |
| `POST` | `/v1/users` | `S1-E-06` |
| `GET` | `/v1/users/{id}` | `S1-E-06` |
| `PATCH` | `/v1/users/{id}` | `S1-E-06` |
| `DELETE` | `/v1/users/{id}` | `S1-E-06` |
| `POST` | `/v1/users/{id}/roles` | `S1-E-06` |
| `DELETE` | `/v1/users/{id}/roles/{role_id}` | `S1-E-06` |
| `GET` | `/v1/roles` / `POST /v1/roles` | `S1-E-06` |
| `GET` | `/v1/permissions` | `S1-E-06` |
| `GET` / `POST` | `/v1/branches` | `S1-E-06` |
| `GET` | `/v1/audit-logs` | `BL-IAM-004` |

## gRPC methods

### Planned — not in BL-IAM-001

All internal-service RPCs land in **BL-IAM-002** (`feat/s1-e-04-iam-middleware`):

- `iam.v1.IamService/ValidateToken` — every service calls this on every authenticated request.
- `iam.v1.IamService/CheckPermission` — resource + action + scope resolution.
- `iam.v1.IamService/GetUser` — read-only user lookup for cross-service user-id → profile.
- `iam.v1.IamService/RecordAudit` — append to the audit log from any service (BL-IAM-004).

The current scaffold gRPC surface is the pilot `iam.v1.IamService/Healthz` placeholder only.

---

## BL-IAM-001 wire shapes

Shapes below are the authoritative `{data:…}` / `{error:…}` envelopes returned by `services/iam-svc/api/rest_oapi/`. Error envelope is documented once at the bottom.

### `POST /v1/sessions` — login

**Request:**

```json
{
  "email": "admin@umrohos.dev",
  "password": "password123",
  "totp_code": "123456"
}
```

`totp_code` is optional. It is accepted but not validated in BL-IAM-001 — login-time TOTP enforcement lands in `S1-E-06`. Clients should start sending the code as soon as a user verifies their TOTP secret, so the switch to enforcement is not a breaking change.

**Success — `200 OK`:**

```json
{
  "data": {
    "access_token": "v2.local.xxxx...",
    "refresh_token": "<64-hex-chars>",
    "access_expires_at": "2026-04-21T14:09:51.878123+07:00",
    "refresh_expires_at": "2026-04-28T13:54:51.576722+07:00",
    "user": {
      "user_id": "33333333-3333-3333-3333-333333333333",
      "email": "admin@umrohos.dev",
      "name": "Dev Admin",
      "branch_id": "11111111-1111-1111-1111-111111111111",
      "status": "active"
    }
  }
}
```

**Errors:** `400 VALIDATION_ERROR` (malformed body), `401 UNAUTHORIZED` (wrong email or password — same code, no existence leak), `403 FORBIDDEN` (user `suspended` or `pending`).

### `POST /v1/sessions/refresh` — rotate

**Request:** `{"refresh_token": "<64-hex-chars>"}`.

**Success — `200 OK`:** same shape as login response, minus the `user` block (caller already has it from login):

```json
{
  "data": {
    "access_token": "v2.local.xxxx...",
    "refresh_token": "<64-new-hex-chars>",
    "access_expires_at": "2026-04-21T14:10:05+07:00",
    "refresh_expires_at": "2026-04-28T13:55:05+07:00"
  }
}
```

The old refresh row is revoked **inside the same transaction** as the new row's insert. Replaying the old refresh token triggers revoke-all-sessions for that user (defensive).

**Errors:** `400`, `401 UNAUTHORIZED` (unknown / expired / already-revoked refresh), `403 FORBIDDEN` (user no longer active).

### `DELETE /v1/sessions` — logout

**Request:** empty body; `Authorization: Bearer <access_token>` required.

**Success — `204 No Content`.** Idempotent — a second call on an already-revoked session row still returns 204.

**Errors:** `401 UNAUTHORIZED` (missing / malformed / invalid / expired bearer).

### `GET /v1/me` — current user

**Request:** `Authorization: Bearer <access_token>` required.

**Success — `200 OK`:**

```json
{
  "data": {
    "user": {
      "user_id": "33333333-3333-3333-3333-333333333333",
      "email": "admin@umrohos.dev",
      "name": "Dev Admin",
      "branch_id": "11111111-1111-1111-1111-111111111111",
      "status": "active"
    },
    "totp_enrolled": false,
    "totp_verified": false
  }
}
```

**Errors:** `401 UNAUTHORIZED`, `404 NOT_FOUND` (user row soft-deleted after token issuance — rare).

### `POST /v1/me/2fa/enroll` — begin TOTP enrollment

**Request:** empty body; bearer required.

**Success — `200 OK`:**

```json
{
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "otpauth_url": "otpauth://totp/UmrohOS:admin@umrohos.dev?secret=JBSWY3DPEHPK3PXP&issuer=UmrohOS"
  }
}
```

The `secret` (plaintext base32) is returned exactly once; clients must show it and then discard from memory. The server-side row stores an AES-256-GCM ciphertext of the same secret; the key comes from `config.totp.encryption_key`. Until the secret is **verified** via `/v1/me/2fa/verify`, `users.totp_verified_at` stays NULL.

**Errors:** `401 UNAUTHORIZED`, `409 CONFLICT` (user already has a verified TOTP secret — admin reset lands in `S1-E-06`).

### `POST /v1/me/2fa/verify` — verify + stamp verified_at

**Request:** `{"code": "123456"}`; bearer required.

**Success — `200 OK`:**

```json
{ "data": { "verified_at": "2026-04-21T07:10:22.145123Z" } }
```

**Errors:** `400 VALIDATION_ERROR` (malformed body), `401 UNAUTHORIZED` (invalid code OR missing/invalid bearer), `422 VALIDATION_ERROR` (TOTP not yet enrolled).

### Shared error envelope

Every non-2xx response uses:

```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "unauthorized\n<contextual detail from the server log>"
  }
}
```

Codes produced by iam-svc today: `UNAUTHORIZED`, `FORBIDDEN`, `VALIDATION_ERROR`, `NOT_FOUND`, `CONFLICT`, `INTERNAL_ERROR`.
