# iam-svc — Overview

## Purpose

Identity, access control, and audit. Every other service calls iam-svc to validate tokens and check permissions.

## Bounded context

Identity & Access. See `docs/02-domain/00-bounded-contexts.md` § 1.

## PRD source

PRD section H — Admin & Security (RBAC, audit trail, system config).

## Owns (data)

- `users` — login accounts (staff, agents, jamaah)
- `roles` — named permission bundles
- `permissions` — capabilities (resource + action)
- `role_permissions` — many-to-many
- `user_roles` — many-to-many
- `branches` — company offices, scope boundary
- `sessions` — active token records
- `audit_logs` — immutable CRUD record

## Boundaries (does NOT own)

- Jamaah biodata (`jamaah-svc`)
- Agent network hierarchy (`crm-svc`) — even though agents have `users`, the hierarchy lives in crm
- Documents and OCR (`jamaah-svc`)
- Per-resource business rules (each owning service enforces its own)

## Interactions

- **Inbound:** every service calls `IamService.ValidateToken` and `IamService.CheckPermission` via gRPC.
- **Outbound:** none in the synchronous path. Audit log writes happen via local store.

## Notable behaviors

- **RBAC** with granular permission matrix (Read/Write/Edit/Delete per module).
- **Data scope** hierarchy (Global / Branch / Personal). Permission checks must take scope into account.
- **2FA (OTP)** mandatory for admin-level operations. Implement via TOTP.
- **Audit trail** is immutable — no UPDATE on `audit_logs`.
- **Token format:** PASETO via the template's `util/token`.
