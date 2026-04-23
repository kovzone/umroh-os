# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: 04-uat-s1-auth-catalog.spec.ts >> S1 Auth — Login & Token (BL-IAM-001 / BL-IAM-018) >> S1-AUTH-01: login valid → 200 + PASETO token + user profile
- Location: tests/04-uat-s1-auth-catalog.spec.ts:38:7

# Error details

```
Error: Login harus 200

expect(received).toBe(expected) // Object.is equality

Expected: 200
Received: 404
```

# Test source

```ts
  1   | /**
  2   |  * UAT Suite 04 — Auth (IAM) + Catalog (S1)
  3   |  *
  4   |  * Covers:
  5   |  *  - BL-IAM-001..004: login, refresh, /v1/me, permissions, audit, suspend
  6   |  *  - BL-CAT-001..004, BL-CAT-014: public read + staff write
  7   |  *  - BL-FE-CONSOLE-001..002: console login UI + shell
  8   |  *  - BL-FE-CAT-001: console catalog CRUD UI
  9   |  *  - BL-FE-BOOK-001..003: B2C catalog → booking form UI
  10  |  *
  11  |  * Run: npx playwright test tests/04-uat-s1-auth-catalog.spec.ts --project=api
  12  |  *      npx playwright test tests/04-uat-s1-auth-catalog.spec.ts --project=browser
  13  |  */
  14  | 
  15  | import { test, expect } from "@playwright/test";
  16  | import { createApiClient } from "../lib/api-client";
  17  | import {
  18  |   UAT_ENV,
  19  |   UAT_PREFIX,
  20  |   loginAdmin,
  21  |   createUatPackage,
  22  |   createUatDeparture,
  23  |   deleteUatPackage,
  24  |   cleanupUatData,
  25  | } from "../lib/uat-helpers";
  26  | import { gateway } from "../lib/services";
  27  | 
  28  | // ─── Shared state across API tests ───────────────────────────────────────────
  29  | let adminToken = "";
  30  | let createdPackageId = "";
  31  | let createdDepartureId = "";
  32  | 
  33  | // ═══════════════════════════════════════════════════════════════════════════════
  34  | // 1. AUTH — API Tests (project: api)
  35  | // ═══════════════════════════════════════════════════════════════════════════════
  36  | 
  37  | test.describe.serial("S1 Auth — Login & Token (BL-IAM-001 / BL-IAM-018)", () => {
  38  |   test("S1-AUTH-01: login valid → 200 + PASETO token + user profile", async () => {
  39  |     const api = await createApiClient(gateway.baseURL);
  40  |     const res = await api.post("/v1/auth/login", {
  41  |       email: UAT_ENV.adminEmail,
  42  |       password: UAT_ENV.adminPassword,
  43  |     });
  44  | 
> 45  |     expect(res.status(), "Login harus 200").toBe(200);
      |                                             ^ Error: Login harus 200
  46  |     const body = await res.json();
  47  |     expect(body.data.access_token, "access_token harus PASETO v2.local.*").toMatch(/^v2\.local\./);
  48  |     expect(body.data.refresh_token).toHaveLength(64);
  49  |     expect(body.data.user.email).toBe(UAT_ENV.adminEmail);
  50  |     expect(body.data.user.user_id).toBe(UAT_ENV.adminUserId);
  51  |     expect(body.data.user.status).toBe("active");
  52  |     expect(new Date(body.data.access_expires_at).getTime()).toBeGreaterThan(Date.now());
  53  | 
  54  |     adminToken = body.data.access_token;
  55  |   });
  56  | 
  57  |   test("S1-AUTH-02: login password salah → 401, bukan 500", async () => {
  58  |     const api = await createApiClient(gateway.baseURL);
  59  |     const res = await api.post("/v1/auth/login", {
  60  |       email: UAT_ENV.adminEmail,
  61  |       password: "password-yang-salah",
  62  |     });
  63  | 
  64  |     expect(res.status(), "Password salah harus 401").toBe(401);
  65  |     const body = await res.json();
  66  |     expect(body.error.code).toBe("UNAUTHORIZED");
  67  |   });
  68  | 
  69  |   test("S1-AUTH-03: login email tidak dikenal → 401, bukan 404 (tidak leak user existence)", async () => {
  70  |     const api = await createApiClient(gateway.baseURL);
  71  |     const res = await api.post("/v1/auth/login", {
  72  |       email: "ghost@umrohos.dev",
  73  |       password: "apapun",
  74  |     });
  75  | 
  76  |     expect(res.status(), "Email tidak ada harus 401, bukan 404").toBe(401);
  77  |   });
  78  | 
  79  |   test("S1-AUTH-04: GET /v1/me dengan token valid → 200 + user data", async () => {
  80  |     const api = await createApiClient(gateway.baseURL, adminToken);
  81  |     const res = await api.get("/v1/me");
  82  | 
  83  |     expect(res.status(), "/v1/me dengan token valid harus 200").toBe(200);
  84  |     const body = await res.json();
  85  |     expect(body.data.email).toBe(UAT_ENV.adminEmail);
  86  |   });
  87  | 
  88  |   test("S1-AUTH-05: GET /v1/me tanpa token → 401", async () => {
  89  |     const api = await createApiClient(gateway.baseURL);
  90  |     const res = await api.get("/v1/me");
  91  | 
  92  |     expect(res.status(), "Request tanpa token harus 401").toBe(401);
  93  |   });
  94  | 
  95  |   test("S1-AUTH-06: GET /v1/me dengan token garbage → 401", async () => {
  96  |     const api = await createApiClient(gateway.baseURL, "Bearer ini-bukan-token");
  97  |     const res = await api.get("/v1/me");
  98  | 
  99  |     expect(res.status(), "Token garbage harus 401").toBe(401);
  100 |   });
  101 | 
  102 |   test("S1-AUTH-07: POST /v1/auth/refresh → token baru", async () => {
  103 |     // Login dulu untuk dapat refresh token
  104 |     const api = await createApiClient(gateway.baseURL);
  105 |     const loginRes = await api.post("/v1/auth/login", {
  106 |       email: UAT_ENV.adminEmail,
  107 |       password: UAT_ENV.adminPassword,
  108 |     });
  109 |     expect(loginRes.status()).toBe(200);
  110 |     const loginBody = await loginRes.json();
  111 |     const refreshToken = loginBody.data.refresh_token;
  112 | 
  113 |     // Gunakan refresh token
  114 |     const refreshRes = await api.post("/v1/auth/refresh", {
  115 |       refresh_token: refreshToken,
  116 |     });
  117 | 
  118 |     expect(refreshRes.status(), "Refresh token harus 200").toBe(200);
  119 |     const refreshBody = await refreshRes.json();
  120 |     expect(refreshBody.data.access_token).toMatch(/^v2\.local\./);
  121 |     expect(refreshBody.data.access_token).not.toBe(loginBody.data.access_token);
  122 |   });
  123 | 
  124 |   test("S1-AUTH-08: iam-svc down → gateway harus 502, bukan 200 palsu (fail-closed)", async () => {
  125 |     // Test ini hanya bisa dilakukan dengan mematikan iam-svc — skip di UAT normal
  126 |     // Tandai sebagai informational
  127 |     test.skip(true, "Membutuhkan akses ke docker compose untuk stop iam-svc");
  128 |   });
  129 | });
  130 | 
  131 | // ═══════════════════════════════════════════════════════════════════════════════
  132 | // 2. PERMISSIONS & AUDIT (BL-IAM-002..004)
  133 | // ═══════════════════════════════════════════════════════════════════════════════
  134 | 
  135 | test.describe.serial("S1 Auth — Permissions & Protected Routes (BL-IAM-002)", () => {
  136 |   test("S1-PERM-01: route staff tanpa token → 401", async () => {
  137 |     const api = await createApiClient(gateway.baseURL);
  138 |     // POST /v1/packages adalah staff-only route
  139 |     const res = await api.post("/v1/packages", { name: "test" });
  140 |     expect(res.status(), "Staff route tanpa token harus 401").toBe(401);
  141 |   });
  142 | 
  143 |   test("S1-PERM-02: DELETE /v1/auth/logout dengan token valid → 200", async () => {
  144 |     // Login sementara
  145 |     const api = await createApiClient(gateway.baseURL);
```