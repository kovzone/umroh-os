# iam-svc — Status

## Implementation checklist

- [ ] Scaffolded from baseline template
- [ ] Wired into `docker-compose.dev.yml`
- [ ] Initial DDL written (`_init/iam_db/`)
- [ ] sqlc queries for users, roles, permissions
- [ ] OpenAPI spec for auth endpoints
- [ ] Auth handlers (login, refresh, logout, me)
- [ ] gRPC service for `ValidateToken`, `CheckPermission`
- [ ] 2FA (TOTP) flow
- [ ] Audit log write path
- [ ] User CRUD endpoints (admin)
- [ ] Role CRUD endpoints
- [ ] Branch CRUD endpoints
- [ ] Unit tests for service layer
- [ ] Integration tests for auth flow
- [ ] Verified by reviewer in `testing-guide.md`

## Current status

**Not started.** Service does not yet exist on disk. First task: scaffold from baseline template (see Suggested Next Steps in `progress.md`).
