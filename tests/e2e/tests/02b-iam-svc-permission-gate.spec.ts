// iam permission-gate end-to-end (BL-IAM-002 + BL-IAM-018 + BL-IAM-019).
//
// Proves the "finance routes denied for non-finance roles" F1 acceptance
// across the iam-svc gRPC hops. Every call targets gateway-svc:4000:
// login hits gateway REST (iam REST was retired in BL-IAM-018) and
// /v1/finance/ping now hits gateway too (moved from finance-svc:4009 in
// BL-IAM-019 / S1-E-14). Authorization is enforced by the gateway's
// RequireBearerToken + RequirePermission("journal_entry","read","global")
// middleware chain before finance-svc is called over gRPC.
//
//   1. finance@umrohos.dev logs in (gateway REST → iam gRPC) → GET
//      /v1/finance/ping on the gateway passes both middleware checks and
//      returns 200 with an envelope assembled from the bearer's identity.
//   2. cs@umrohos.dev logs in → same call → 403 FORBIDDEN from
//      RequirePermission (no grant).
//   3. No bearer → 401 UNAUTHORIZED (RequireBearerToken fails first).
//   4. Garbage bearer → 401 UNAUTHORIZED (iam-svc.ValidateToken rejects).
//
// Prereqs (make dev-bootstrap handles them):
//   1. docker compose dev stack up
//   2. migrations 000004 + 000005 applied (fixture users/roles/permission)
//   3. gateway-svc on :4000

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { gateway } from "../lib/services";

const FINANCE_EMAIL = "finance@umrohos.dev";
const FINANCE_USER_ID = "77777777-7777-7777-7777-777777777777";
const CS_EMAIL = "cs@umrohos.dev";
const SHARED_PASSWORD = "password123";

async function login(email: string): Promise<string> {
  const api = await createApiClient(gateway.baseURL);
  const res = await api.post("/v1/sessions", {
    email,
    password: SHARED_PASSWORD,
  });
  expect(res.status(), `login ${email}`).toBe(200);
  const body = await res.json();
  return body.data.access_token as string;
}

test.describe.serial("gateway /v1/finance/ping — permission gate (BL-IAM-002)", () => {
  test("finance_admin can hit /v1/finance/ping (200, identity echoed back)", async () => {
    const token = await login(FINANCE_EMAIL);
    const api = await createApiClient(gateway.baseURL, token);

    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(200);

    const body = await res.json();
    expect(body.data.message).toBe("ok");
    expect(body.data.user_id).toBe(FINANCE_USER_ID);
    expect(body.data.roles).toEqual(expect.arrayContaining(["finance_admin"]));
  });

  test("cs_agent is denied /v1/finance/ping (403 FORBIDDEN)", async () => {
    const token = await login(CS_EMAIL);
    const api = await createApiClient(gateway.baseURL, token);

    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(403);

    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("no bearer → 401 UNAUTHORIZED (middleware fails closed before iam-svc)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(401);

    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("garbage bearer → 401 UNAUTHORIZED (iam-svc rejects token)", async () => {
    const api = await createApiClient(gateway.baseURL, "not-a-real-paseto-token");
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(401);

    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });
});
