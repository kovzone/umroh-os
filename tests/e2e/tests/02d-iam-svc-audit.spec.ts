// iam-svc audit-emit end-to-end (BL-IAM-004).
//
// Closes F1 § AC acceptance ("Every state-changing call in every service
// produces an audit log entry") for the first in-process emitter: the admin
// suspend action. BL-IAM-004 adds the iam.v1.IamService.RecordAudit gRPC
// surface for cross-service callers AND wires iam-svc.service.SuspendUser to
// call q.InsertAuditLog inside its existing WithTx so business-success ↔
// audit-success are atomic.
//
// This spec covers the in-process emit end-to-end:
//
//   1. Admin suspends sacrifice via POST /v1/users/{id}/suspend (same surface
//      exercised by 02c).
//   2. One new row lands in iam.audit_logs with the shape the spec pins:
//      resource="user", action="suspend", actor=admin, resource_id=target,
//      old_value.status ∈ {"active","suspended"}, new_value.status="suspended".
//   3. Re-suspending (idempotent status-wise) still writes a second audit
//      row — every admin action is auditable, including no-op ones.
//   4. The append-only trigger on iam.audit_logs rejects UPDATE and DELETE
//      against a just-written row (F1: "audit log rows cannot be mutated").
//
// The direct-RPC wire (RecordAudit over gRPC from a Go client) is covered
// by service-layer unit tests in services/iam-svc/service/audit_test.go; the
// first cross-service consumer lands with S1-E-03 (booking-svc draft create),
// which will exercise the RPC end-to-end via its new iam_grpc_adapter.
//
// Prereqs (make dev-bootstrap handles all):
//   1. docker compose dev stack up (iam-svc :4001, postgres :5432)
//   2. migrations 000004..000006 applied (fixture admin + sacrifice target)
//   3. pg host access on localhost:5432 with postgres/changeme

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { gateway } from "../lib/services";
import { withPg } from "../lib/pg-client";

const ADMIN_EMAIL = "admin@umrohos.dev";
const ADMIN_USER_ID = "33333333-3333-3333-3333-333333333333";

const TARGET_USER_ID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa";

const SHARED_PASSWORD = "password123";

async function loginAdmin(): Promise<string> {
  // Migrated in S1-E-12: POST /v1/auth/login on gateway-svc:4000
  const api = await createApiClient(gateway.baseURL);
  const res = await api.post("/v1/auth/login", {
    email: ADMIN_EMAIL,
    password: SHARED_PASSWORD,
  });
  expect(res.status(), "admin login").toBe(200);
  const body = await res.json();
  return body.data.access_token as string;
}

// latestSuspendAudit returns the most recent iam.audit_logs row for
// (resource="user", resource_id=target, action="suspend"). Returns null when
// no row exists (test harness should not crash on empty table).
async function latestSuspendAudit(targetUserID: string) {
  return withPg(async (pg) => {
    const rows = await pg.query(
      `SELECT id, user_id, branch_id, resource, resource_id, action,
              old_value, new_value, created_at
         FROM iam.audit_logs
        WHERE resource = 'user'
          AND resource_id = $1
          AND action = 'suspend'
        ORDER BY created_at DESC
        LIMIT 1`,
      [targetUserID],
    );
    return rows.rows[0] ?? null;
  });
}

async function countSuspendAudits(targetUserID: string): Promise<number> {
  return withPg(async (pg) => {
    const r = await pg.query(
      `SELECT COUNT(*)::int AS n
         FROM iam.audit_logs
        WHERE resource = 'user'
          AND resource_id = $1
          AND action = 'suspend'`,
      [targetUserID],
    );
    return r.rows[0].n as number;
  });
}

test.describe.serial("iam-svc audit emit on suspend (BL-IAM-004)", () => {
  let adminAccess = "";
  let baselineCount = 0;

  test("baseline — capture existing suspend-audit count for the target", async () => {
    adminAccess = await loginAdmin();
    baselineCount = await countSuspendAudits(TARGET_USER_ID);
  });

  test("admin suspends target — audit row lands inside the same tx", async () => {
    const api = await createApiClient(gateway.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(200);

    // Exactly one new row since baseline — the WithTx emitted one InsertAuditLog
    // after the status flip and session revoke.
    const after = await countSuspendAudits(TARGET_USER_ID);
    expect(after - baselineCount).toBe(1);

    const row = await latestSuspendAudit(TARGET_USER_ID);
    expect(row, "latest suspend audit row").not.toBeNull();

    // Actor is the admin; resource tuple is (user, target-uuid, suspend).
    expect(row.user_id).toBe(ADMIN_USER_ID);
    expect(row.resource).toBe("user");
    expect(row.resource_id).toBe(TARGET_USER_ID);
    expect(row.action).toBe("suspend");

    // JSONB payloads reflect the status transition. old_value depends on the
    // target's status at call time (active on first suspend, suspended on
    // re-suspend) — pin new_value exactly, accept either value on old_value.
    expect(row.new_value).toEqual({ status: "suspended" });
    expect(["active", "suspended"]).toContain(row.old_value.status);

    // Branch is inherited from the target user (fixture user lives in the
    // seeded branch from migration 000004). Non-null, UUID-shaped.
    expect(row.branch_id).toMatch(
      /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/,
    );
    expect(new Date(row.created_at).getTime()).toBeGreaterThan(
      Date.now() - 60_000,
    );
  });

  test("re-suspending the already-suspended target still writes a second audit row (every admin action is auditable)", async () => {
    const beforeCount = await countSuspendAudits(TARGET_USER_ID);

    const api = await createApiClient(gateway.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(200);

    const afterCount = await countSuspendAudits(TARGET_USER_ID);
    expect(afterCount - beforeCount).toBe(1);

    const row = await latestSuspendAudit(TARGET_USER_ID);
    // On the re-suspend, old_value captures the already-suspended state.
    expect(row.old_value).toEqual({ status: "suspended" });
    expect(row.new_value).toEqual({ status: "suspended" });
  });

  test("append-only trigger rejects UPDATE against a written audit row (F1 AC)", async () => {
    const row = await latestSuspendAudit(TARGET_USER_ID);
    expect(row, "audit row must exist from the previous test").not.toBeNull();

    const err = await withPg(async (pg) => {
      try {
        await pg.query(
          `UPDATE iam.audit_logs SET action = 'tampered' WHERE id = $1`,
          [row.id],
        );
        return null;
      } catch (e: unknown) {
        return e as Error & { code?: string };
      }
    });
    expect(err, "UPDATE must be rejected by the append-only trigger").not.toBeNull();
    // pgx / pg error code for insufficient_privilege is "42501".
    expect(err!.code).toBe("42501");
  });

  test("append-only trigger rejects DELETE against a written audit row (F1 AC)", async () => {
    const row = await latestSuspendAudit(TARGET_USER_ID);
    expect(row, "audit row must exist from the previous test").not.toBeNull();

    const err = await withPg(async (pg) => {
      try {
        await pg.query(`DELETE FROM iam.audit_logs WHERE id = $1`, [row.id]);
        return null;
      } catch (e: unknown) {
        return e as Error & { code?: string };
      }
    });
    expect(err, "DELETE must be rejected by the append-only trigger").not.toBeNull();
    expect(err!.code).toBe("42501");
  });

  // Migration 000007 narrows the append-only trigger so the FK cascade path
  // (iam.users deletion → ON DELETE SET NULL → UPDATE audit_logs.user_id = NULL)
  // is permitted. Exercising the cascade end-to-end would require deleting the
  // admin fixture mid-test — destructive. Instead, reproduce the cascade's
  // exact UPDATE shape directly against the row and assert it succeeds, then
  // verify a non-FK-column UPDATE on the same row is still rejected.
  test("FK cascade path (user_id → NULL only, all other columns frozen) is allowed by the narrowed trigger", async () => {
    const row = await latestSuspendAudit(TARGET_USER_ID);
    expect(row, "audit row must exist from the previous test").not.toBeNull();
    expect(row.user_id, "row must have a non-null actor to exercise the null cascade").not.toBeNull();

    const result = await withPg(async (pg) => {
      try {
        const r = await pg.query(
          `UPDATE iam.audit_logs SET user_id = NULL WHERE id = $1 RETURNING user_id`,
          [row.id],
        );
        return { ok: true, userId: r.rows[0].user_id };
      } catch (e: unknown) {
        return { ok: false, err: e as Error & { code?: string } };
      }
    });
    expect(result.ok, "user_id → NULL must be allowed by the narrowed trigger").toBe(true);
    expect(result.userId).toBeNull();
  });

  test("non-FK-column UPDATE is still rejected even with the narrowed trigger", async () => {
    const row = await latestSuspendAudit(TARGET_USER_ID);
    expect(row, "audit row must exist from the previous test").not.toBeNull();

    const err = await withPg(async (pg) => {
      try {
        await pg.query(
          `UPDATE iam.audit_logs SET action = 'tampered', user_id = NULL WHERE id = $1`,
          [row.id],
        );
        return null;
      } catch (e: unknown) {
        return e as Error & { code?: string };
      }
    });
    expect(err, "UPDATE touching action must still be rejected").not.toBeNull();
    expect(err!.code).toBe("42501");
  });
});
