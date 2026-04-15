# iam-svc — Events

iam-svc is mostly synchronous (gRPC reads/writes). It does not own a workflow.

## Events emitted

| Event | When | Payload | Consumed by |
|---|---|---|---|
| `iam.user.created` | Successful user creation | user_id, branch_id, role_ids | crm-svc (for agent users) |
| `iam.user.suspended` | Admin suspends a user | user_id, reason | (future) all services may subscribe to invalidate cached tokens |
| `iam.session.revoked` | Logout or admin revoke | user_id, session_id | none yet |
| `iam.role.changed` | Role assignment / revocation | user_id, role_id, action | none yet |

> Event delivery mechanism is TBD (likely a future event bus when needed; Temporal is deferred per ADR 0006). For Phase 1, iam-svc emits no events — they are notional.

## Events consumed

None. iam-svc is upstream of everything else.

## RecordAudit (gRPC)

Although not technically an event, every service writes to iam-svc's audit log via `IamService.RecordAudit`. This is the closest thing iam-svc has to inbound event consumption. The signature is logged in `01-api.md`.
