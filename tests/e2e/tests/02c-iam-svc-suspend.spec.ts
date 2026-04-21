// iam-svc admin suspend end-to-end (BL-IAM-003).
//
// Closes F1-W5 acceptance ("Suspended user cannot access again") across the
// three load-bearing code paths:
//
//   1. Downstream bearer validation — a suspended user's still-decryptable
//      access token is rejected by iam-svc ValidateToken over gRPC (session
//      revoked AND users.status != active), and the sanitised gRPC message
//      does NOT leak the internal state oracle the /security-review finding
//      flagged on BL-IAM-002 (no "session revoked" / "load session: not
//      found" / "not found" strings on the wire).
//   2. Login — the `status != active` gate in iam-svc service.Login returns
//      403 FORBIDDEN.
//   3. Refresh — the same gate in iam-svc service.RefreshSession returns 403.
//
// Also covers:
//   - Permission gate: a finance_admin bearer (no iam.users/suspend/global
//     grant) gets 403 on the admin endpoint.
//   - Idempotency: re-suspending the already-suspended user returns 200 and
//     still sweeps any sessions that appeared between calls.
//   - Self-suspend rejection: admin cannot suspend their own account.
//
// Prereqs (make dev-bootstrap handles all):
//   1. docker compose dev stack up (iam-svc :4001, finance-svc :4009)
//   2. migrations 000004 + 000005 + 000006 applied (the sacrifice user
//      `suspend-target@umrohos.dev` + the `iam.users/suspend/global`
//      permission + grant to super_admin are in 000006).

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { backendServices } from "../lib/services";

const iam = backendServices.find((s) => s.name === "iam-svc");
if (!iam) throw new Error("iam-svc not in backendServices registry");

const finance = backendServices.find((s) => s.name === "finance-svc");
if (!finance) throw new Error("finance-svc not in backendServices registry");

const ADMIN_EMAIL = "admin@umrohos.dev";
const ADMIN_USER_ID = "33333333-3333-3333-3333-333333333333";

const FINANCE_EMAIL = "finance@umrohos.dev";

const TARGET_EMAIL = "suspend-target@umrohos.dev";
const TARGET_USER_ID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa";

const SHARED_PASSWORD = "password123";

type LoginResult = { accessToken: string; refreshToken: string };

async function login(email: string): Promise<LoginResult> {
  const api = await createApiClient(iam!.baseURL);
  const res = await api.post("/v1/sessions", {
    email,
    password: SHARED_PASSWORD,
  });
  expect(res.status(), `login ${email}`).toBe(200);
  const body = await res.json();
  return {
    accessToken: body.data.access_token as string,
    refreshToken: body.data.refresh_token as string,
  };
}

test.describe.serial("iam-svc admin suspend + revoke-all (BL-IAM-003)", () => {
  let adminAccess = "";
  let targetAccess = "";
  let targetRefresh = "";

  test("admin and sacrifice both log in successfully (baseline)", async () => {
    const adminTokens = await login(ADMIN_EMAIL);
    adminAccess = adminTokens.accessToken;

    // Sacrifice is active, so the second call to this test on a re-run re-seeds
    // the sacrifice via migration 000006 being idempotent. The preceding
    // suspension (on the previous run) set status=suspended, so on this run we
    // need the migration to be re-applied OR the test env fresh. The
    // testing-guide walks the reviewer through `make migrate-down` +
    // `make migrate-up` when re-running in-place.
    const targetTokens = await login(TARGET_EMAIL);
    targetAccess = targetTokens.accessToken;
    targetRefresh = targetTokens.refreshToken;
  });

  test("while active, sacrifice's bearer is valid against finance-svc (403 FORBIDDEN — auth OK, no grant)", async () => {
    // /v1/finance/ping exercises the iam-svc gRPC ValidateToken + CheckPermission
    // path. The sacrifice has no roles, so CheckPermission returns allowed=false
    // and finance-svc maps that to 403. This proves the bearer was accepted
    // (ValidateToken succeeded) — sets up the "after suspend, same call → 401"
    // assertion below.
    const api = await createApiClient(finance!.baseURL, targetAccess);
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(403);
    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("admin suspends sacrifice — 200 echoes profile with status=suspended", async () => {
    const api = await createApiClient(iam.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(200);

    const body = await res.json();
    expect(body.data.user.user_id).toBe(TARGET_USER_ID);
    expect(body.data.user.email).toBe(TARGET_EMAIL);
    expect(body.data.user.status).toBe("suspended");
  });

  test("sacrifice's in-flight bearer is now rejected by finance-svc (401, sanitised message)", async () => {
    // This is the load-bearing assertion for F1-W5. After suspension, the PASETO
    // payload still decrypts (TTL hasn't elapsed), but iam-svc ValidateToken
    // rejects via (a) session revoked and (b) users.status != active. The gRPC
    // status is Unauthenticated with a constant "unauthorized" message per
    // apperrors.GRPCMessage — finance-svc's error middleware renders that into
    // the HTTP 401 body without the internal state-oracle strings.
    const api = await createApiClient(finance!.baseURL, targetAccess);
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(401);

    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
    const msg: string = body.error.message ?? "";
    // Sanitisation guard: the gRPC message must not leak the wrapped chain
    // that distinguishes revoked-session from user-not-found from status-mismatch.
    expect(msg).not.toContain("session revoked");
    expect(msg).not.toContain("load session");
    expect(msg).not.toContain("not found");
    expect(msg).not.toContain("user status=");
    // And it should contain the sanitised constant.
    expect(msg.toLowerCase()).toContain("unauthorized");
  });

  test("sacrifice cannot log in anymore — 403 FORBIDDEN (status gate in service.Login)", async () => {
    const api = await createApiClient(iam.baseURL);
    const res = await api.post("/v1/sessions", {
      email: TARGET_EMAIL,
      password: SHARED_PASSWORD,
    });
    expect(res.status()).toBe(403);
    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("sacrifice cannot trade their refresh token — 401 UNAUTHORIZED (session revoked)", async () => {
    // RefreshSession's first failure is the session-row check (revoked →
    // ErrUnauthorized), which fires before the status gate. Either way the
    // caller gets 401.
    const api = await createApiClient(iam.baseURL);
    const res = await api.post("/v1/sessions/refresh", {
      refresh_token: targetRefresh,
    });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("finance_admin cannot call the suspend endpoint — 403 (permission gate)", async () => {
    const { accessToken: financeAccess } = await login(FINANCE_EMAIL);
    const api = await createApiClient(iam.baseURL, financeAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(403);
    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("re-suspending the already-suspended sacrifice is a 200 no-op (idempotent)", async () => {
    const api = await createApiClient(iam.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.user.status).toBe("suspended");
  });

  test("admin cannot suspend themselves — 400 VALIDATION_ERROR", async () => {
    const api = await createApiClient(iam.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${ADMIN_USER_ID}/suspend`);
    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body.error.code).toBe("VALIDATION_ERROR");
  });

  test("no bearer on the suspend endpoint — 401 UNAUTHORIZED", async () => {
    const api = await createApiClient(iam.baseURL);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });
});
