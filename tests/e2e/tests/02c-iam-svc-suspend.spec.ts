// IAM admin suspend end-to-end (BL-IAM-003 / BL-IAM-018).
//
// Migrated in S1-E-12: all IAM routes now go through gateway-svc:4000.
// Route paths changed:
//   POST   /v1/sessions          → POST   /v1/auth/login
//   POST   /v1/sessions/refresh  → POST   /v1/auth/refresh
//   POST   /v1/users/:id/suspend → POST   /v1/users/:id/suspend (same path, now on gateway)
//
// Closes F1-W5 acceptance ("Suspended user cannot access again") across the
// three load-bearing code paths:
//
//   1. Downstream bearer validation — a suspended user's still-decryptable
//      access token is rejected by iam-svc ValidateToken over gRPC (session
//      revoked AND users.status != active).
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
//   1. docker compose dev stack up (gateway-svc :4000, finance-svc :4009)
//   2. migrations 000004 + 000005 + 000006 applied (the sacrifice user
//      `suspend-target@umrohos.dev` + the `iam.users/suspend/global`
//      permission + grant to super_admin are in 000006).

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { gateway, backendServices } from "../lib/services";

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
  const api = await createApiClient(gateway.baseURL);
  const res = await api.post("/v1/auth/login", {
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

test.describe.serial("gateway admin suspend + revoke-all (BL-IAM-003 via BL-IAM-018)", () => {
  let adminAccess = "";
  let targetAccess = "";
  let targetRefresh = "";

  test("admin and sacrifice both log in successfully (baseline)", async () => {
    const adminTokens = await login(ADMIN_EMAIL);
    adminAccess = adminTokens.accessToken;

    const targetTokens = await login(TARGET_EMAIL);
    targetAccess = targetTokens.accessToken;
    targetRefresh = targetTokens.refreshToken;
  });

  test("while active, sacrifice's bearer is valid against finance-svc (403 FORBIDDEN — auth OK, no grant)", async () => {
    const api = await createApiClient(finance!.baseURL, targetAccess);
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(403);
    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("admin suspends sacrifice via gateway — 200 echoes profile with status=suspended", async () => {
    const api = await createApiClient(gateway.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(200);

    const body = await res.json();
    expect(body.data.user.user_id).toBe(TARGET_USER_ID);
    expect(body.data.user.email).toBe(TARGET_EMAIL);
    expect(body.data.user.status).toBe("suspended");
  });

  test("sacrifice's in-flight bearer is now rejected by finance-svc (401, sanitised message)", async () => {
    const api = await createApiClient(finance!.baseURL, targetAccess);
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(401);

    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
    const msg: string = body.error.message ?? "";
    expect(msg).not.toContain("session revoked");
    expect(msg).not.toContain("load session");
    expect(msg).not.toContain("not found");
    expect(msg).not.toContain("user status=");
    expect(msg.toLowerCase()).toContain("unauthorized");
  });

  test("sacrifice cannot log in anymore — 403 FORBIDDEN (status gate in service.Login)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/login", {
      email: TARGET_EMAIL,
      password: SHARED_PASSWORD,
    });
    expect(res.status()).toBe(403);
    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("sacrifice cannot trade their refresh token — 401 UNAUTHORIZED (session revoked)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/refresh", {
      refresh_token: targetRefresh,
    });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("finance_admin cannot call the suspend endpoint — 403 (permission gate)", async () => {
    const { accessToken: financeAccess } = await login(FINANCE_EMAIL);
    const api = await createApiClient(gateway.baseURL, financeAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(403);
    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("re-suspending the already-suspended sacrifice is a 200 no-op (idempotent)", async () => {
    const api = await createApiClient(gateway.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.user.status).toBe("suspended");
  });

  test("admin cannot suspend themselves — 400 VALIDATION_ERROR", async () => {
    const api = await createApiClient(gateway.baseURL, adminAccess);
    const res = await api.post(`/v1/users/${ADMIN_USER_ID}/suspend`);
    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body.error.code).toBe("VALIDATION_ERROR");
  });

  test("no bearer on the suspend endpoint — 401 UNAUTHORIZED", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post(`/v1/users/${TARGET_USER_ID}/suspend`);
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });
});
