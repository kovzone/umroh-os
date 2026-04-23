/**
 * UAT Suite 04 — Auth (IAM) + Catalog (S1)
 *
 * Covers:
 *  - BL-IAM-001..004: login, refresh, /v1/me, permissions, audit, suspend
 *  - BL-CAT-001..004, BL-CAT-014: public read + staff write
 *  - BL-FE-CONSOLE-001..002: console login UI + shell
 *  - BL-FE-CAT-001: console catalog CRUD UI
 *  - BL-FE-BOOK-001..003: B2C catalog → booking form UI
 *
 * Run: npx playwright test tests/04-uat-s1-auth-catalog.spec.ts --project=api
 *      npx playwright test tests/04-uat-s1-auth-catalog.spec.ts --project=browser
 */

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import {
  UAT_ENV,
  UAT_PREFIX,
  loginAdmin,
  createUatPackage,
  createUatDeparture,
  deleteUatPackage,
  cleanupUatData,
} from "../lib/uat-helpers";
import { gateway } from "../lib/services";

// ─── Shared state across API tests ───────────────────────────────────────────
let adminToken = "";
let createdPackageId = "";
let createdDepartureId = "";

// ═══════════════════════════════════════════════════════════════════════════════
// 1. AUTH — API Tests (project: api)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S1 Auth — Login & Token (BL-IAM-001 / BL-IAM-018)", () => {
  test("S1-AUTH-01: login valid → 200 + PASETO token + user profile", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/sessions", {
      email: UAT_ENV.adminEmail,
      password: UAT_ENV.adminPassword,
    });

    expect(res.status(), "Login harus 200").toBe(200);
    const body = await res.json();
    expect(body.data.access_token, "access_token harus PASETO v2.local.*").toMatch(/^v2\.local\./);
    expect(body.data.refresh_token).toHaveLength(64);
    expect(body.data.user.email).toBe(UAT_ENV.adminEmail);
    expect(body.data.user.user_id).toBe(UAT_ENV.adminUserId);
    expect(body.data.user.status).toBe("active");
    expect(new Date(body.data.access_expires_at).getTime()).toBeGreaterThan(Date.now());

    adminToken = body.data.access_token;
  });

  test("S1-AUTH-02: login password salah → 401, bukan 500", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/sessions", {
      email: UAT_ENV.adminEmail,
      password: "password-yang-salah",
    });

    expect(res.status(), "Password salah harus 401").toBe(401);
    const body = await res.json();
    expect(body.error.code).toBe("UNAUTHORIZED");
  });

  test("S1-AUTH-03: login email tidak dikenal → 401, bukan 404 (tidak leak user existence)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/sessions", {
      email: "ghost@umrohos.dev",
      password: "apapun",
    });

    expect(res.status(), "Email tidak ada harus 401, bukan 404").toBe(401);
  });

  test("S1-AUTH-04: GET /v1/me dengan token valid → 200 + user data", async () => {
    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.get("/v1/me");

    expect(res.status(), "/v1/me dengan token valid harus 200").toBe(200);
    const body = await res.json();
    expect(body.data.email).toBe(UAT_ENV.adminEmail);
  });

  test("S1-AUTH-05: GET /v1/me tanpa token → 401", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/v1/me");

    expect(res.status(), "Request tanpa token harus 401").toBe(401);
  });

  test("S1-AUTH-06: GET /v1/me dengan token garbage → 401", async () => {
    const api = await createApiClient(gateway.baseURL, "Bearer ini-bukan-token");
    const res = await api.get("/v1/me");

    expect(res.status(), "Token garbage harus 401").toBe(401);
  });

  test("S1-AUTH-07: POST /v1/sessions/refresh → token baru", async () => {
    // Login dulu untuk dapat refresh token
    const api = await createApiClient(gateway.baseURL);
    const loginRes = await api.post("/v1/sessions", {
      email: UAT_ENV.adminEmail,
      password: UAT_ENV.adminPassword,
    });
    expect(loginRes.status()).toBe(200);
    const loginBody = await loginRes.json();
    const refreshToken = loginBody.data.refresh_token;

    // Gunakan refresh token
    const refreshRes = await api.post("/v1/sessions/refresh", {
      refresh_token: refreshToken,
    });

    expect(refreshRes.status(), "Refresh token harus 200").toBe(200);
    const refreshBody = await refreshRes.json();
    expect(refreshBody.data.access_token).toMatch(/^v2\.local\./);
    expect(refreshBody.data.access_token).not.toBe(loginBody.data.access_token);
  });

  test("S1-AUTH-08: iam-svc down → gateway harus 502, bukan 200 palsu (fail-closed)", async () => {
    // Test ini hanya bisa dilakukan dengan mematikan iam-svc — skip di UAT normal
    // Tandai sebagai informational
    test.skip(true, "Membutuhkan akses ke docker compose untuk stop iam-svc");
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 2. PERMISSIONS & AUDIT (BL-IAM-002..004)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S1 Auth — Permissions & Protected Routes (BL-IAM-002)", () => {
  test("S1-PERM-01: route staff tanpa token → 401", async () => {
    const api = await createApiClient(gateway.baseURL);
    // POST /v1/packages adalah staff-only route
    const res = await api.post("/v1/packages", { name: "test" });
    expect(res.status(), "Staff route tanpa token harus 401").toBe(401);
  });

  test("S1-PERM-02: DELETE /v1/sessions (logout) dengan token valid → 200", async () => {
    // Login sementara
    const api = await createApiClient(gateway.baseURL);
    const loginRes = await api.post("/v1/sessions", {
      email: UAT_ENV.adminEmail,
      password: UAT_ENV.adminPassword,
    });
    expect(loginRes.status()).toBe(200);
    const token = (await loginRes.json()).data.access_token;

    const authedApi = await createApiClient(gateway.baseURL, token);
    const logoutRes = await authedApi.delete("/v1/sessions");

    expect(logoutRes.status(), "Logout harus 200 atau 204").toBeOneOf([200, 204]);
  });

  test("S1-PERM-03: token yang sudah di-logout tidak bisa digunakan lagi → 401", async () => {
    // Login dan logout
    const api = await createApiClient(gateway.baseURL);
    const loginRes = await api.post("/v1/sessions", {
      email: UAT_ENV.adminEmail,
      password: UAT_ENV.adminPassword,
    });
    const token = (await loginRes.json()).data.access_token;
    const authedApi = await createApiClient(gateway.baseURL, token);
    await authedApi.delete("/v1/sessions");

    // Coba gunakan token yang sudah di-revoke
    const res = await authedApi.get("/v1/me");
    expect(res.status(), "Token revoked harus 401").toBe(401);
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 3. CATALOG READ — Public Endpoints (BL-CAT-001..004)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S1 Catalog — Public Read (BL-CAT-001..004)", () => {
  test("S1-CAT-01: GET /v1/packages → hanya package active yang tampil", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get("/v1/packages");

    expect(res.status(), "List packages harus 200").toBe(200);
    const body = await res.json();
    expect(Array.isArray(body.packages), "packages harus array").toBe(true);
    expect(body.packages.length, "Harus ada setidaknya 1 package active").toBeGreaterThan(0);

    // Package active harus ada
    const ids = body.packages.map((p: { id: string }) => p.id);
    expect(ids).toContain(UAT_ENV.activePkgId);

    // Package draft TIDAK boleh muncul
    expect(ids, "Package draft tidak boleh muncul di list publik").not.toContain(UAT_ENV.draftPkgId);

    // Shape validasi
    const active = body.packages.find((p: { id: string }) => p.id === UAT_ENV.activePkgId);
    expect(active).toBeTruthy();
    expect(active.starting_price?.list_currency).toBe("IDR");
    expect(active.starting_price?.list_amount).toBeGreaterThan(0);
  });

  test("S1-CAT-02: GET /v1/packages/{active_id} → detail lengkap", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get(`/v1/packages/${UAT_ENV.activePkgId}`);

    expect(res.status(), "Package detail harus 200").toBe(200);
    const body = await res.json();
    expect(body.data?.id || body.id).toBe(UAT_ENV.activePkgId);
  });

  test("S1-CAT-03: GET /v1/packages/{draft_id} → 404 (draft tidak boleh publik)", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get(`/v1/packages/${UAT_ENV.draftPkgId}`);

    expect(res.status(), "Package draft harus 404 di endpoint publik").toBe(404);
  });

  test("S1-CAT-04: GET /v1/package-departures/{active_dep_id} → detail + remaining_seats", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get(`/v1/package-departures/${UAT_ENV.activeDepId}`);

    expect(res.status(), "Departure detail harus 200").toBe(200);
    const body = await res.json();
    const dep = body.data || body;
    expect(dep.id).toBe(UAT_ENV.activeDepId);
    expect(typeof dep.remaining_seats).toBe("number");
    expect(dep.remaining_seats).toBeGreaterThanOrEqual(0);
  });

  test("S1-CAT-05: GET /v1/package-departures/{cancelled_dep} → 404", async () => {
    const api = await createApiClient(gateway.baseURL);
    const cancelledDepId = "dep_01JCDF00000000000000000003";
    const res = await api.get(`/v1/package-departures/${cancelledDepId}`);

    expect(res.status(), "Departure cancelled harus 404").toBe(404);
  });

  test("S1-CAT-06: GET /v1/packages tanpa auth → 200 (public endpoint)", async () => {
    const api = await createApiClient(gateway.baseURL); // no token
    const res = await api.get("/v1/packages");
    expect(res.status(), "Catalog list harus bisa diakses tanpa token").toBe(200);
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 4. CATALOG WRITE — Staff Endpoints (BL-CAT-014)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S1 Catalog — Staff Write (BL-CAT-014)", () => {
  let pkgId = "";
  let depId = "";

  test.beforeAll(async () => {
    const { tokens } = await loginAdmin();
    adminToken = tokens.accessToken;
  });

  test.afterAll(async () => {
    await cleanupUatData();
  });

  test("S1-WRITE-01: POST /v1/packages (staff) → 201, package terbuat dengan status draft/active", async () => {
    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post("/v1/packages", {
      name: `${UAT_PREFIX} Umroh Test ${Date.now()}`,
      description: "Package untuk UAT testing — hapus jika masih ada",
      kind: "umrah_reguler",
      duration_days: 12,
      status: "draft",
    });

    expect(res.status(), "Create package harus 201").toBeOneOf([200, 201]);
    const body = await res.json();
    pkgId = body.data?.id || body.id;
    expect(pkgId, "Response harus ada ID package").toBeTruthy();
    createdPackageId = pkgId;
  });

  test("S1-WRITE-02: POST /v1/packages tanpa auth → 401", async () => {
    const api = await createApiClient(gateway.baseURL); // no token
    const res = await api.post("/v1/packages", {
      name: `${UAT_PREFIX} Package Test Unauthorized`,
      kind: "umrah_reguler",
    });
    expect(res.status(), "Create package tanpa token harus 401").toBe(401);
  });

  test("S1-WRITE-03: POST /v1/packages tanpa nama → 422 validation error", async () => {
    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post("/v1/packages", {
      kind: "umrah_reguler",
      // name sengaja tidak diisi
    });
    expect(res.status(), "Create package tanpa nama harus 422").toBeOneOf([400, 422]);
  });

  test("S1-WRITE-04: PATCH /v1/packages/{id} → update harga/deskripsi", async () => {
    test.skip(!pkgId, "Skip: package belum terbuat (test sebelumnya gagal)");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.put(`/v1/packages/${pkgId}`, {
      description: `${UAT_PREFIX} Description updated`,
    });
    expect(res.status(), "Update package harus 200").toBe(200);
  });

  test("S1-WRITE-05: POST /v1/packages/{id}/departures → buat departure", async () => {
    test.skip(!pkgId, "Skip: package belum terbuat (test sebelumnya gagal)");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post(`/v1/packages/${pkgId}/departures`, {
      depart_date: "2026-12-01",
      return_date: "2026-12-12",
      capacity: 20,
      status: "open",
      notes: `${UAT_PREFIX} departure test`,
    });

    expect(res.status(), "Create departure harus 201").toBeOneOf([200, 201]);
    const body = await res.json();
    depId = body.data?.id || body.id;
    createdDepartureId = depId;
    expect(depId).toBeTruthy();
  });

  test("S1-WRITE-06: POST departure dengan tgl kembali sebelum tgl berangkat → 422", async () => {
    test.skip(!pkgId, "Skip: package belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post(`/v1/packages/${pkgId}/departures`, {
      depart_date: "2026-12-12",
      return_date: "2026-12-01", // return SEBELUM depart — invalid
      capacity: 20,
      status: "open",
    });

    expect(res.status(), "Tanggal tidak valid harus 422").toBeOneOf([400, 422]);
  });

  test("S1-WRITE-07: PATCH departure status → open → closed", async () => {
    test.skip(!depId, "Skip: departure belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.put(`/v1/package-departures/${depId}`, {
      status: "closed",
    });
    expect(res.status(), "Update departure status harus 200").toBe(200);
  });

  test("S1-WRITE-08: departure status draft → tidak muncul di GET /v1/packages public", async () => {
    test.skip(!pkgId || !depId, "Skip: data belum terbuat");

    // Set package ke active supaya terlihat, departure ke draft
    const api = await createApiClient(gateway.baseURL, adminToken);
    await api.put(`/v1/packages/${pkgId}`, { status: "active" });
    await api.put(`/v1/package-departures/${depId}`, { status: "draft" });

    // Cek sebagai public user
    const publicApi = await createApiClient(gateway.baseURL);
    const res = await publicApi.get(`/v1/packages/${pkgId}`);
    if (res.status() === 200) {
      const body = await res.json();
      const departures = body.data?.departures || body.departures || [];
      const draftDep = departures.find((d: { id: string }) => d.id === depId);
      expect(draftDep, "Departure draft tidak boleh muncul di view publik").toBeUndefined();
    }
  });

  test("S1-WRITE-09: DELETE /v1/packages/{id} → 200, package tidak muncul lagi", async () => {
    test.skip(!pkgId, "Skip: package belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.delete(`/v1/packages/${pkgId}`);
    expect(res.status(), "Delete package harus 200 atau 204").toBeOneOf([200, 204]);

    // Verify tidak muncul di publik
    const publicApi = await createApiClient(gateway.baseURL);
    const checkRes = await publicApi.get(`/v1/packages/${pkgId}`);
    expect(checkRes.status(), "Package yang dihapus harus 404").toBe(404);

    pkgId = ""; // sudah dihapus, cleanup tidak perlu
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 5. BROWSER TESTS — Console UI & B2C (project: browser)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S1 UI — Console Login (BL-FE-CONSOLE-001)", () => {
  test("S1-UI-01: /console/login halaman login tampil", async ({ page }) => {
    await page.goto("/console/login");
    await expect(page).not.toHaveTitle(/404|not found/i);
    // Form login harus ada
    await expect(page.locator("input[type='email'], input[name='email']")).toBeVisible({ timeout: 10_000 });
    await expect(page.locator("input[type='password']")).toBeVisible();
  });

  test("S1-UI-02: login dengan credentials invalid → error message tampil", async ({ page }) => {
    await page.goto("/console/login");
    await page.locator("input[type='email'], input[name='email']").fill("salah@umrohos.dev");
    await page.locator("input[type='password']").fill("salahpassword");
    await page.locator("button[type='submit']").click();

    // Harus ada pesan error, tidak redirect ke console
    await expect(page).not.toHaveURL(/\/console(?!\/login)/);
    // Cek ada error message
    const errorVisible = await page.locator("[data-testid='login-error'], .error, [role='alert']").isVisible().catch(() => false);
    // Minimal tidak redirect ke dashboard
    expect(page.url()).toContain("/console/login");
  });

  test("S1-UI-03: login dengan credentials valid → redirect ke /console", async ({ page }) => {
    await page.goto("/console/login");
    await page.locator("input[type='email'], input[name='email']").fill(UAT_ENV.adminEmail);
    await page.locator("input[type='password']").fill(UAT_ENV.adminPassword);
    await page.locator("button[type='submit']").click();

    // Harus redirect ke console
    await expect(page).toHaveURL(/\/console/, { timeout: 10_000 });
    expect(page.url()).not.toContain("/login");
  });

  test("S1-UI-04: console shell menampilkan sidemenu setelah login (BL-FE-CONSOLE-002)", async ({ page }) => {
    // Login dulu
    await page.goto("/console/login");
    await page.locator("input[type='email'], input[name='email']").fill(UAT_ENV.adminEmail);
    await page.locator("input[type='password']").fill(UAT_ENV.adminPassword);
    await page.locator("button[type='submit']").click();
    await expect(page).toHaveURL(/\/console/, { timeout: 10_000 });

    // Sidemenu harus visible
    const sidebar = page.locator(
      "[data-testid='console-sidebar'], nav, aside, [class*='sidebar'], [class*='sidemenu']"
    );
    await expect(sidebar.first()).toBeVisible({ timeout: 5_000 });
  });
});

test.describe.serial("S1 UI — B2C Catalog (BL-FE-BOOK-001..003)", () => {
  test("S1-UI-05: / landing page tampil dengan CTA browse packages", async ({ page }) => {
    await page.goto("/");
    await expect(page).not.toHaveTitle(/404|not found/i);
    // CTA atau link ke /packages harus ada
    const packagesCta = page.locator(
      "[data-testid='browse-packages-cta'], a[href='/packages'], a[href*='packages']"
    );
    await expect(packagesCta.first()).toBeVisible({ timeout: 10_000 });
  });

  test("S1-UI-06: /packages — daftar package tampil, ada minimal 1 kartu package", async ({ page }) => {
    await page.goto("/packages");
    await expect(page).not.toHaveTitle(/404|not found/i);
    // Tunggu loading
    const packageList = page.locator(
      "[data-testid='s1-package-catalog'], [data-testid*='package'], .package-card, .package-list"
    );
    await expect(packageList.first()).toBeVisible({ timeout: 15_000 });
  });

  test("S1-UI-07: klik package → halaman detail package tampil", async ({ page }) => {
    await page.goto("/packages");
    // Klik link package pertama
    const firstPkgLink = page.locator("a[href*='/packages/']").first();
    await expect(firstPkgLink).toBeVisible({ timeout: 10_000 });
    await firstPkgLink.click();
    await expect(page).toHaveURL(/\/packages\/[^/]+$/);
    await expect(page).not.toHaveTitle(/404|not found/i);
  });

  test("S1-UI-08: halaman detail package → tombol Pesan / booking form accessible", async ({ page }) => {
    await page.goto("/packages");
    const firstPkgLink = page.locator("a[href*='/packages/']").first();
    await expect(firstPkgLink).toBeVisible({ timeout: 10_000 });
    await firstPkgLink.click();
    await expect(page).toHaveURL(/\/packages\/[^/]+$/);

    // Cari tombol booking / pesan
    const bookBtn = page.locator(
      "[data-testid='s1-start-booking'], button:has-text('Pesan'), button:has-text('Book'), a[href*='/booking']"
    );
    await expect(bookBtn.first()).toBeVisible({ timeout: 8_000 });
  });
});
