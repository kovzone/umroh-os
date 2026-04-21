// iam-svc permission-gate end-to-end (BL-IAM-002).
//
// Proves the "finance routes denied for non-finance roles" F1 acceptance
// across the two new iam-svc gRPC hops and the finance-svc consumer surface:
//
//   1. finance@umrohos.dev logs in (iam-svc REST) → GET /v1/finance/ping
//      (finance-svc REST) calls iam-svc.ValidateToken + CheckPermission over
//      gRPC and returns 200.
//   2. cs@umrohos.dev logs in → same call → 403 FORBIDDEN (no grant).
//   3. No bearer → 401 UNAUTHORIZED (middleware fails closed).
//   4. Garbage bearer → 401 UNAUTHORIZED (PASETO verify fails inside iam-svc).
//
// Prereqs (make dev-bootstrap handles all three):
//   1. docker compose dev stack up
//   2. migrations 000004 + 000005 applied (fixture users/roles/permission)
//   3. iam-svc on :4001, finance-svc on :4009

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { backendServices } from "../lib/services";

const iam = backendServices.find((s) => s.name === "iam-svc");
if (!iam) throw new Error("iam-svc not in backendServices registry");

const finance = backendServices.find((s) => s.name === "finance-svc");
if (!finance) throw new Error("finance-svc not in backendServices registry");

const FINANCE_EMAIL = "finance@umrohos.dev";
const FINANCE_USER_ID = "77777777-7777-7777-7777-777777777777";
const CS_EMAIL = "cs@umrohos.dev";
const SHARED_PASSWORD = "password123";

async function login(email: string): Promise<string> {
  const api = await createApiClient(iam!.baseURL);
  const res = await api.post("/v1/sessions", {
    email,
    password: SHARED_PASSWORD,
  });
  expect(res.status(), `login ${email}`).toBe(200);
  const body = await res.json();
  return body.data.access_token as string;
}

test.describe.serial("finance-svc /v1/finance/ping — permission gate (BL-IAM-002)", () => {
  test("finance_admin can hit /v1/finance/ping (200, identity echoed back)", async () => {
    const token = await login(FINANCE_EMAIL);
    const api = await createApiClient(finance!.baseURL, token);

    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(200);

    const body = await res.json();
    expect(body.data.message).toBe("ok");
    expect(body.data.user_id).toBe(FINANCE_USER_ID);
    expect(body.data.roles).toEqual(expect.arrayContaining(["finance_admin"]));
  });

  test("cs_agent is denied /v1/finance/ping (403 FORBIDDEN)", async () => {
    const token = await login(CS_EMAIL);
    const api = await createApiClient(finance!.baseURL, token);

    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(403);

    const body = await res.json();
    expect(body.error.code).toBe("FORBIDDEN");
  });

  test("no bearer → 401 UNAUTHORIZED (middleware fails closed before iam-svc)", async () => {
    const api = await createApiClient(finance!.baseURL);
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(401);

    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("garbage bearer → 401 UNAUTHORIZED (iam-svc rejects token)", async () => {
    const api = await createApiClient(finance!.baseURL, "not-a-real-paseto-token");
    const res = await api.get("/v1/finance/ping");
    expect(res.status()).toBe(401);

    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });
});
