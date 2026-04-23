// IAM sessions end-to-end (BL-IAM-001 / BL-IAM-018).
//
// Migrated in S1-E-12 from iam-svc:4001 → gateway-svc:4000.
// Per ADR 0009, all client-facing IAM routes are now proxied through
// gateway-svc REST → iam-svc gRPC. Route paths changed:
//   POST   /v1/sessions        → POST   /v1/auth/login
//   POST   /v1/sessions/refresh → POST  /v1/auth/refresh
//   DELETE /v1/sessions        → DELETE /v1/auth/logout
//   GET    /v1/me              → GET    /v1/me  (unchanged)
//   POST   /v1/me/2fa/enroll   → POST   /v1/me/2fa/enroll  (unchanged)
//
// Exercises the full login → me → refresh → logout happy path plus the 401
// sad paths (no bearer, garbage bearer, expired refresh, refresh replay)
// against the dev compose stack.
//
// Prereqs (make dev-bootstrap handles all three):
//   1. docker compose dev stack up
//   2. migration 000004_seed_initial_admin applied (HQ branch + super_admin
//      role + admin@umrohos.dev/password123 user)
//   3. gateway-svc listening on :4000, iam-svc listening on :50051 (gRPC only)

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import { gateway } from "../lib/services";

const ADMIN_EMAIL = "admin@umrohos.dev";
const ADMIN_PASSWORD = "password123";
const ADMIN_USER_ID = "33333333-3333-3333-3333-333333333333";
const HQ_BRANCH_ID = "11111111-1111-1111-1111-111111111111";

test.describe.serial("gateway /v1/auth + /v1/me — IAM sessions (BL-IAM-001 via BL-IAM-018)", () => {
  // Shared state across the ordered flow; each step expects the previous one's output.
  let accessToken = "";
  let refreshToken = "";
  let priorRefreshToken = "";

  test("login returns access + refresh tokens and user profile", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/login", {
      email: ADMIN_EMAIL,
      password: ADMIN_PASSWORD,
    });
    expect(res.status()).toBe(200);

    const body = await res.json();
    expect(body.data.access_token).toMatch(/^v2\.local\./); // PASETO v2 local
    expect(typeof body.data.refresh_token).toBe("string");
    expect(body.data.refresh_token).toHaveLength(64); // 32-byte hex
    expect(body.data.user.user_id).toBe(ADMIN_USER_ID);
    expect(body.data.user.email).toBe(ADMIN_EMAIL);
    expect(body.data.user.branch_id).toBe(HQ_BRANCH_ID);
    expect(body.data.user.status).toBe("active");
    expect(new Date(body.data.access_expires_at).getTime()).toBeGreaterThan(Date.now());
    expect(new Date(body.data.refresh_expires_at).getTime()).toBeGreaterThan(Date.now());

    accessToken = body.data.access_token;
    refreshToken = body.data.refresh_token;
  });

  test("login with wrong password returns 401 UNAUTHORIZED (no user leak)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/login", {
      email: ADMIN_EMAIL,
      password: "definitely-not-the-password",
    });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("login with unknown email returns 401 UNAUTHORIZED (not 404)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/login", {
      email: "ghost@umrohos.dev",
      password: "irrelevant",
    });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("GET /v1/me with no bearer returns 401", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/v1/me");
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("GET /v1/me with garbage bearer returns 401", async () => {
    const api = await createApiClient(gateway.baseURL, "not-a-real-token");
    const res = await api.get("/v1/me");
    expect(res.status()).toBe(401);
  });

  test("GET /v1/me with the issued access token returns the admin profile", async () => {
    const api = await createApiClient(gateway.baseURL, accessToken);
    const res = await api.get("/v1/me");
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.user.user_id).toBe(ADMIN_USER_ID);
    expect(body.data.user.email).toBe(ADMIN_EMAIL);
    expect(body.data.user.branch_id).toBe(HQ_BRANCH_ID);
    // totp_enrolled / totp_verified are not asserted strictly: once an earlier
    // run of this spec's enroll step writes a secret, re-runs will observe
    // enrolled=true. Only their presence + type matters at this point.
    expect(typeof body.data.totp_enrolled).toBe("boolean");
    expect(typeof body.data.totp_verified).toBe("boolean");
  });

  test("refresh rotates tokens and invalidates the prior refresh", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/refresh", {
      refresh_token: refreshToken,
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body.data.access_token).not.toBe(accessToken);
    expect(body.data.refresh_token).not.toBe(refreshToken);

    priorRefreshToken = refreshToken;
    accessToken = body.data.access_token;
    refreshToken = body.data.refresh_token;
  });

  test("replaying the rotated-out refresh token returns 401", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/auth/refresh", {
      refresh_token: priorRefreshToken,
    });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("logout with the new access token returns 204 + subsequent refresh returns 401", async () => {
    // The active refresh after the first rotation was flagged revoked by the
    // replay step above (defensive revoke-all-for-user). Log in again to get
    // a clean session pair, then test the logout flow.
    const login = await createApiClient(gateway.baseURL);
    const loginRes = await login.post("/v1/auth/login", {
      email: ADMIN_EMAIL,
      password: ADMIN_PASSWORD,
    });
    expect(loginRes.status()).toBe(200);
    const loginBody = await loginRes.json();
    const freshAccess = loginBody.data.access_token;
    const freshRefresh = loginBody.data.refresh_token;

    const auth = await createApiClient(gateway.baseURL, freshAccess);
    const logoutRes = await auth.delete("/v1/auth/logout");
    expect(logoutRes.status()).toBe(204);

    // Refresh with the now-revoked token must 401 (replay guard also triggers
    // revoke-all-for-user, so refreshing again would also fail even if we
    // tried the prior hash).
    const anon = await createApiClient(gateway.baseURL);
    const refRes = await anon.post("/v1/auth/refresh", {
      refresh_token: freshRefresh,
    });
    expect(refRes.status()).toBe(401);
  });

  test("POST /v1/me/2fa/enroll returns a base32 secret and otpauth URL", async () => {
    // Clean login for a fresh access token (we just logged out in the previous test).
    const api = await createApiClient(gateway.baseURL);
    const loginRes = await api.post("/v1/auth/login", {
      email: ADMIN_EMAIL,
      password: ADMIN_PASSWORD,
    });
    const loginBody = await loginRes.json();
    const freshAccess = loginBody.data.access_token;

    const auth = await createApiClient(gateway.baseURL, freshAccess);
    const res = await auth.post("/v1/me/2fa/enroll");

    // The seeded admin may or may not already have an enrolled TOTP secret
    // depending on whether an earlier run of this spec persisted one.
    // Accept either 200 (first run, fresh enroll) or 409 (re-run against a
    // verified user — we don't /verify in this spec, so this 409 only
    // occurs if something else verified the secret first).
    expect([200, 409]).toContain(res.status());
    if (res.status() === 200) {
      const body = await res.json();
      expect(body.data.secret).toMatch(/^[A-Z2-7]+=*$/); // base32
      expect(body.data.otpauth_url).toMatch(/^otpauth:\/\/totp\//);
    }
  });
});
