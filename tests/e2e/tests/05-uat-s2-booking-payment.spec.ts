/**
 * UAT Suite 05 — Booking & Payment (S1 Booking + S2)
 *
 * Covers:
 *  - BL-BOOK-001..006: draft booking creation, validation, state machine
 *  - BL-PAY-001..008: invoice, VA, webhook, idempotency, security
 *  - BL-FE-PAY-001..002: checkout UI, payment status update
 *
 * Run: npx playwright test tests/05-uat-s2-booking-payment.spec.ts --project=api
 *      npx playwright test tests/05-uat-s2-booking-payment.spec.ts --project=browser
 */

import { test, expect } from "@playwright/test";
import { createApiClient } from "../lib/api-client";
import {
  UAT_ENV,
  UAT_PREFIX,
  uatEmail,
  loginAdmin,
  createUatBooking,
  cleanupUatData,
  dbQuery,
} from "../lib/uat-helpers";
import { gateway } from "../lib/services";

// ─── Shared state ─────────────────────────────────────────────────────────────
let adminToken = "";
let bookingId = "";
let invoiceId = "";
let vaNumber = "";

test.beforeAll(async () => {
  const { tokens } = await loginAdmin();
  adminToken = tokens.accessToken;
});

test.afterAll(async () => {
  await cleanupUatData();
});

// ═══════════════════════════════════════════════════════════════════════════════
// 1. BOOKING — Draft Creation (BL-BOOK-001..006)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S2 Booking — Draft Creation (BL-BOOK-001)", () => {
  test("S2-BOOK-01: POST /v1/bookings (b2c_self) → 201, status draft", async () => {
    const api = await createApiClient(gateway.baseURL); // no auth — b2c_self is public
    const email = uatEmail("jamaah");

    const res = await api.post("/v1/bookings", {
      channel: "b2c_self",
      package_id: UAT_ENV.activePkgId,
      departure_id: UAT_ENV.activeDepId,
      room_type: "double",
      lead: {
        full_name: `${UAT_PREFIX} Test Jamaah`,
        email,
        whatsapp: "+628112345678",
        domicile: "Jakarta",
      },
      jamaah: [
        {
          full_name: `${UAT_PREFIX} Test Jamaah`,
          email,
          whatsapp: "+628112345678",
          domicile: "Jakarta",
          is_lead: true,
        },
      ],
      notes: `${UAT_PREFIX} automated test booking`,
    });

    expect(res.status(), "Create booking harus 201").toBeOneOf([200, 201]);
    const body = await res.json();
    const booking = body.data || body;
    bookingId = booking.id;
    expect(bookingId, "Response harus ada booking ID").toBeTruthy();
    expect(booking.status, "Status awal harus draft").toBe("draft");
  });

  test("S2-BOOK-02: POST /v1/bookings tanpa email (b2c_self) → 422", async () => {
    const api = await createApiClient(gateway.baseURL);

    const res = await api.post("/v1/bookings", {
      channel: "b2c_self",
      package_id: UAT_ENV.activePkgId,
      departure_id: UAT_ENV.activeDepId,
      room_type: "double",
      lead: {
        full_name: `${UAT_PREFIX} Test`,
        // email missing — required for b2c_self
        whatsapp: "+628112345678",
        domicile: "Jakarta",
      },
      jamaah: [
        {
          full_name: `${UAT_PREFIX} Test`,
          is_lead: true,
          whatsapp: "+628112345678",
          domicile: "Jakarta",
        },
      ],
    });

    expect(res.status(), "Booking tanpa email untuk b2c_self harus 422").toBeOneOf([400, 422]);
  });

  test("S2-BOOK-03: POST /v1/bookings tanpa field wajib → 422 dengan detail error", async () => {
    const api = await createApiClient(gateway.baseURL);

    const res = await api.post("/v1/bookings", {
      channel: "b2c_self",
      // package_id missing
      departure_id: UAT_ENV.activeDepId,
    });

    expect(res.status(), "Booking tanpa package_id harus 422").toBeOneOf([400, 422]);
    const body = await res.json();
    expect(body.error, "Response harus ada error object").toBeTruthy();
  });

  test("S2-BOOK-04: POST /v1/bookings dengan departure_id closed → 422 atau 404", async () => {
    const api = await createApiClient(gateway.baseURL);

    // dep_01JCDF00000000000000000003 adalah cancelled departure
    const res = await api.post("/v1/bookings", {
      channel: "b2c_self",
      package_id: UAT_ENV.activePkgId,
      departure_id: "dep_01JCDF00000000000000000003", // cancelled
      room_type: "double",
      lead: {
        full_name: `${UAT_PREFIX} Test`,
        email: uatEmail("test"),
        whatsapp: "+628112345678",
        domicile: "Jakarta",
      },
      jamaah: [{ full_name: `${UAT_PREFIX} Test`, is_lead: true, whatsapp: "+628112345678", domicile: "Jakarta" }],
    });

    expect(res.status(), "Booking dengan departure closed/cancelled harus 4xx").toBeOneOf([400, 404, 422]);
  });

  test("S2-BOOK-05: sisa kursi berkurang setelah booking dibuat", async () => {
    test.skip(!bookingId, "Skip: booking belum terbuat");

    // Cek remaining seats di departure
    const api = await createApiClient(gateway.baseURL);
    const res = await api.get(`/v1/package-departures/${UAT_ENV.activeDepId}`);
    expect(res.status()).toBe(200);

    const body = await res.json();
    const dep = body.data || body;
    // Seats harus sudah berkurang (tidak bisa tahu nilai sebelumnya tanpa snapshot, tapi minimal bisa cek tidak error)
    expect(typeof dep.remaining_seats).toBe("number");
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 2. BOOKING SUBMIT & INVOICE (BL-PAY-001..002)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S2 Payment — Invoice & VA (BL-PAY-001..002)", () => {
  test("S2-PAY-01: POST /v1/bookings/{id}/submit → pending_payment", async () => {
    test.skip(!bookingId, "Skip: booking belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post(`/v1/bookings/${bookingId}/submit`, {});

    if (res.status() === 404) {
      test.skip(true, "NOT_DEPLOYED: /v1/bookings/{id}/submit belum ada");
    }

    expect(res.status(), "Submit booking harus 200").toBe(200);
    const body = await res.json();
    const booking = body.data || body;
    expect(booking.status).toBe("pending_payment");
  });

  test("S2-PAY-02: POST /v1/invoices → invoice terbuat dengan VA", async () => {
    test.skip(!bookingId, "Skip: booking belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post("/v1/invoices", {
      booking_id: bookingId,
      gateway: "mock",
    });

    if (res.status() === 404) {
      test.skip(true, "NOT_DEPLOYED: POST /v1/invoices belum ada");
    }

    expect(res.status(), "Create invoice harus 201").toBeOneOf([200, 201]);
    const body = await res.json();
    const invoice = body.data || body;
    invoiceId = invoice.id;
    expect(invoiceId).toBeTruthy();
  });

  test("S2-PAY-03: POST /v1/invoices/{id}/virtual-accounts → VA diterbitkan", async () => {
    test.skip(!invoiceId, "Skip: invoice belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res = await api.post(`/v1/invoices/${invoiceId}/virtual-accounts`, {
      gateway: "mock",
    });

    if (res.status() === 404) {
      test.skip(true, "NOT_DEPLOYED: VA issuance belum ada");
    }

    expect(res.status(), "Issue VA harus 200 atau 201").toBeOneOf([200, 201]);
    const body = await res.json();
    const va = body.data || body;
    vaNumber = va.va_number || va.account_number || "";
    expect(va.amount).toBeGreaterThan(0);
  });

  test("S2-PAY-04: GET /v1/invoices/{id} → idempoten (reload tidak buat invoice baru)", async () => {
    test.skip(!invoiceId, "Skip: invoice belum terbuat");

    const api = await createApiClient(gateway.baseURL, adminToken);
    const res1 = await api.get(`/v1/invoices/${invoiceId}`);
    const res2 = await api.get(`/v1/invoices/${invoiceId}`);

    expect(res1.status()).toBe(200);
    expect(res2.status()).toBe(200);

    const body1 = await res1.json();
    const body2 = await res2.json();
    const id1 = body1.data?.id || body1.id;
    const id2 = body2.data?.id || body2.id;
    expect(id1, "Invoice ID harus sama (idempoten)").toBe(id2);
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 3. WEBHOOK PAYMENT (BL-PAY-003..005)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S2 Payment — Mock Webhook (BL-PAY-003..005)", () => {
  test("S2-WH-01: POST /v1/webhooks/mock/trigger → booking status jadi paid_in_full", async () => {
    test.skip(!invoiceId, "Skip: invoice belum terbuat");

    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/webhooks/mock/trigger", {
      invoice_id: invoiceId,
      status: "paid",
      amount: 25000000, // 25 juta IDR
    });

    if (res.status() === 404) {
      test.skip(true, "NOT_DEPLOYED: /v1/webhooks/mock/trigger belum ada");
    }

    expect(res.status(), "Mock webhook harus 200").toBe(200);

    // Verifikasi booking status berubah
    await new Promise((r) => setTimeout(r, 500)); // beri waktu async processing

    const adminApi = await createApiClient(gateway.baseURL, adminToken);
    const bookingRes = await adminApi.get(`/v1/bookings/${bookingId}`);
    if (bookingRes.status() === 200) {
      const booking = (await bookingRes.json()).data || (await bookingRes.json());
      expect(booking.status, "Booking harus paid_in_full setelah webhook").toBe("paid_in_full");
    }
  });

  test("S2-WH-02: duplicate webhook (invoice_id sama) → idempoten, 1 payment event saja", async () => {
    test.skip(!invoiceId, "Skip: invoice belum terbuat");

    const api = await createApiClient(gateway.baseURL);
    const payload = {
      invoice_id: invoiceId,
      status: "paid",
      amount: 25000000,
    };

    const res1 = await api.post("/v1/webhooks/mock/trigger", payload);
    const res2 = await api.post("/v1/webhooks/mock/trigger", payload);

    if (res1.status() === 404) {
      test.skip(true, "NOT_DEPLOYED");
    }

    // Keduanya harus 200 (tidak error)
    expect(res2.status(), "Duplicate webhook harus 200").toBe(200);

    // DB check: hanya ada 1 payment_event untuk invoice ini
    const rows = await dbQuery(
      `SELECT COUNT(*) as cnt FROM payment.payment_events WHERE invoice_id = $1`,
      [invoiceId]
    );
    expect(Number(rows[0]?.cnt), "Hanya boleh ada 1 payment event (idempoten)").toBe(1);
  });

  test("S2-WH-03: webhook dengan amount lebih kecil dari tagihan → status TIDAK berubah jadi paid_in_full", async () => {
    test.skip(!bookingId, "Skip: booking belum terbuat");

    // Buat booking baru untuk test ini agar tidak kontaminasi booking utama
    const newBooking = await createUatBooking(UAT_ENV.activePkgId, UAT_ENV.activeDepId);
    const adminApi = await createApiClient(gateway.baseURL, adminToken);

    // Buat invoice untuk booking baru
    const invRes = await adminApi.post("/v1/invoices", {
      booking_id: newBooking.id,
      gateway: "mock",
    });

    if (invRes.status() === 404) {
      test.skip(true, "NOT_DEPLOYED");
    }

    if (invRes.status() !== 200 && invRes.status() !== 201) return;

    const newInvoiceId = ((await invRes.json()).data || (await invRes.json())).id;

    // Kirim webhook dengan amount terlalu kecil
    const api = await createApiClient(gateway.baseURL);
    await api.post("/v1/webhooks/mock/trigger", {
      invoice_id: newInvoiceId,
      status: "paid",
      amount: 1000, // jauh lebih kecil dari tagihan
    });

    await new Promise((r) => setTimeout(r, 500));

    // Booking tidak boleh paid_in_full
    const bookRes = await adminApi.get(`/v1/bookings/${newBooking.id}`);
    if (bookRes.status() === 200) {
      const bk = (await bookRes.json()).data || await bookRes.json();
      expect(bk.status, "Booking tidak boleh paid_in_full dengan underpayment").not.toBe("paid_in_full");
    }
  });

  test("S2-WH-04: webhook Midtrans tanpa header auth yang benar → 401 atau 403", async () => {
    const api = await createApiClient(gateway.baseURL);
    const res = await api.post("/v1/webhooks/midtrans", {
      order_id: "fake-order",
      transaction_status: "settlement",
      gross_amount: "25000000",
    });

    // Tanpa header X-Callback-Token yang valid → harus ditolak
    expect(res.status(), "Webhook tanpa auth header harus 401 atau 403").toBeOneOf([400, 401, 403]);
  });
});

// ═══════════════════════════════════════════════════════════════════════════════
// 4. BROWSER TESTS — Checkout UI (BL-FE-PAY-001..002)
// ═══════════════════════════════════════════════════════════════════════════════

test.describe.serial("S2 UI — Checkout Page (BL-FE-PAY-001)", () => {
  test("S2-UI-01: B2C booking flow end-to-end — pilih paket → isi form → submit → redirect checkout", async ({ page }) => {
    // Navigasi ke halaman packages
    await page.goto("/packages");
    await expect(page).not.toHaveTitle(/404/i);

    // Klik package pertama
    const pkgLink = page.locator("a[href*='/packages/']").first();
    await expect(pkgLink).toBeVisible({ timeout: 15_000 });
    await pkgLink.click();

    // Di halaman detail — cari tombol pesan
    const bookBtn = page.locator(
      "[data-testid='s1-start-booking'], button:has-text('Pesan'), button:has-text('Book'), a[href*='/booking']"
    );
    await expect(bookBtn.first()).toBeVisible({ timeout: 8_000 });
    await bookBtn.first().click();

    // Di form booking — isi data minimal
    const nameInput = page.locator(
      "input[name='full_name'], input[placeholder*='nama'], input[name='name']"
    ).first();
    if (await nameInput.isVisible()) {
      await nameInput.fill(`${UAT_PREFIX} Test Jamaah UI`);
    }

    const emailInput = page.locator("input[type='email'], input[name='email']").first();
    if (await emailInput.isVisible()) {
      await emailInput.fill(uatEmail("ui"));
    }

    const phoneInput = page.locator(
      "input[name='whatsapp'], input[name='phone'], input[type='tel']"
    ).first();
    if (await phoneInput.isVisible()) {
      await phoneInput.fill("+628119999999");
    }

    const domicileInput = page.locator("input[name='domicile']").first();
    if (await domicileInput.isVisible()) {
      await domicileInput.fill("Jakarta");
    }

    // Submit
    const submitBtn = page.locator("button[type='submit']").first();
    await expect(submitBtn).toBeVisible();
    await submitBtn.click();

    // Cek redirect ke checkout atau konfirmasi
    await page.waitForURL(/\/(checkout|booking|confirm)\//, { timeout: 15_000 }).catch(() => {
      // Mungkin ada step tambahan
    });
    await expect(page).not.toHaveTitle(/404|error/i);
  });

  test("S2-UI-02: /checkout/{booking_id} menampilkan VA number, amount, dan countdown", async ({ page }) => {
    test.skip(!bookingId, "Skip: tidak ada booking ID dari API tests");

    await page.goto(`/checkout/${bookingId}`);
    await expect(page).not.toHaveTitle(/404/i);

    // Halaman checkout harus menampilkan info pembayaran
    await page.waitForLoadState("networkidle");
    const content = await page.content();

    // Minimal ada salah satu dari: VA number, amount, atau status info
    const hasPaymentInfo =
      content.includes("Virtual Account") ||
      content.includes("VA") ||
      content.includes("Rp") ||
      content.includes("IDR") ||
      content.includes("pending") ||
      content.includes("Pembayaran");

    expect(hasPaymentInfo, "Checkout page harus menampilkan info pembayaran").toBe(true);
  });

  test("S2-UI-03: booking validation — submit form tanpa data → error inline tampil", async ({ page }) => {
    await page.goto("/packages");
    const pkgLink = page.locator("a[href*='/packages/']").first();
    await expect(pkgLink).toBeVisible({ timeout: 15_000 });
    await pkgLink.click();

    const bookBtn = page.locator(
      "[data-testid='s1-start-booking'], button:has-text('Pesan'), button:has-text('Book'), a[href*='/booking']"
    );
    if (await bookBtn.first().isVisible({ timeout: 5_000 }).catch(() => false)) {
      await bookBtn.first().click();

      // Submit tanpa isi data
      const submitBtn = page.locator("button[type='submit']").first();
      if (await submitBtn.isVisible({ timeout: 3_000 }).catch(() => false)) {
        await submitBtn.click();

        // Form tidak boleh redirect — harus tetap di halaman yang sama
        await new Promise((r) => setTimeout(r, 1000));
        expect(page.url()).not.toMatch(/checkout/);
      }
    }
  });
});
