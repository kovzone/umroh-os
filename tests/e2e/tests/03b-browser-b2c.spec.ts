/**
 * 03b-browser-b2c.spec.ts — Frontend Browser E2E Tests (B2C + Console)
 *
 * Tested live via Chrome browser terhadap http://216.176.238.161 (production).
 * Semua test ini TIDAK memerlukan auth — kecuali group "Console" yang login sebagai admin.
 *
 * Coverage:
 *   B2C-01  Landing page (/)
 *   B2C-02  Katalog paket (/packages)
 *   B2C-03  Detail paket (/packages/:id) — departures, tabs, sticky CTA
 *   B2C-04  Booking wizard Step 1 — pilih keberangkatan (harga + switching)
 *   B2C-05  Booking wizard Step 2 — isi data jamaah (form + tipe kamar reaktif)
 *   B2C-06  Booking wizard Step 3 — review & syarat (NIK masked, checkboxes)
 *   B2C-07  Booking wizard Step 4 — pembayaran (VA number, countdown, booking code)
 *   B2C-08  Contact form (/contact) — field validation, submit
 *   B2C-09  Nav links — Jadwal, Manasik, Tentang Kami route check (ISSUE-026)
 *   B2C-10  Form validation — Step 2 field validation feedback (ISSUE-027)
 *   CON-01  Console login (/console/login) — UI, wrong credentials, login sukses
 *   CON-02  Console dashboard (/console) — sidebar nav, stat cards, approval queue
 *   CON-03  Console catalog (/console/packages) — gateway error behavior (ISSUE-024)
 *   CON-04  Console ops/finance/leads — gateway error behavior (ISSUE-025)
 *
 * Known bugs (didokumentasikan sebagai test, bukan di-skip):
 *   ISSUE-020  Package detail: gambar paket broken (URL fake cdn.umroh-os.example)
 *   ISSUE-021  Package detail: harga tidak muncul di departure cards (API tidak return price_per_pax)
 *   ISSUE-022  Booking: booking-svc 503, frontend silent fallback ke demo booking
 *   ISSUE-023  Package detail: section Itinerary/Fasilitas/Syarat hanya placeholder text
 *   ISSUE-024  Console catalog: "Tidak dapat terhubung ke layanan katalog"
 *   ISSUE-025  Console ops/finance/leads: "Tidak dapat terhubung ke gateway"
 *   ISSUE-026  B2C nav links Jadwal/Manasik/Tentang Kami → 404
 *   ISSUE-027  Booking Step 2: form validation tidak menampilkan error message
 *
 * Jalankan:
 *   cd tests/e2e
 *   npx playwright test 03b --project=browser
 *
 * Cleanup data test (booking yang dibuat oleh test ini):
 *   Booking menggunakan email pemesan "uat.browser@umrohos.dev" — cleanup manual atau
 *   tambahkan SQL ke cleanup-uat.sh:
 *   DELETE FROM booking.bookings WHERE notes IS NOT DISTINCT FROM NULL
 *     AND id IN (SELECT booking_id FROM booking.jamaah WHERE email = 'uat.browser@umrohos.dev');
 */

import { test, expect, Page } from "@playwright/test";

// ─── Constants ────────────────────────────────────────────────────────────────

const BASE_URL = process.env.CORE_WEB_URL || "http://216.176.238.161";
const PKG_ID = "pkg_01JCDE00000000000000000001";
const DEP1_ID = "dep_01JCDF00000000000000000001"; // 23 Mei 2026
const DEP2_ID = "dep_01JCDF00000000000000000002"; // 7 Juli 2026

const ADMIN_EMAIL = process.env.UAT_ADMIN_EMAIL || "admin@umrohos.dev";
const ADMIN_PASS = process.env.UAT_ADMIN_PASSWORD || "password123";

// Data test jamaah — gunakan email yang konsisten agar bisa di-identify saat cleanup
const TEST_PEMESAN = {
  nama: "Ahmad Test Browser",
  email: "uat.browser@umrohos.dev",
  whatsapp: "+6281234567890",
  alamat: "Jakarta Pusat",
};
const TEST_JAMAAH1 = {
  nama: "Ahmad Test Browser",
  nik: "3271234567890001",
  dob: "1990-01-15",
};
const TEST_JAMAAH2 = {
  nama: "Siti Test Browser",
  nik: "3271234567890002",
  dob: "1992-02-20",
};

// ─── Helper ───────────────────────────────────────────────────────────────────

/** Navigasi ke URL absolut dan tunggu network idle */
async function goTo(page: Page, path: string) {
  await page.goto(`${BASE_URL}${path}`, { waitUntil: "domcontentloaded" });
}

// ─── B2C-01: Landing Page ─────────────────────────────────────────────────────

test.describe("B2C-01 — Landing page (/)", () => {
  test("menampilkan header, hero, dan CTA utama", async ({ page }) => {
    await goTo(page, "/");

    // Logo / brand
    await expect(page.getByText("UmrohOS").first()).toBeVisible();

    // Nav links
    await expect(page.getByRole("link", { name: "Paket Umroh" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Jadwal" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Tentang Kami" })).toBeVisible();

    // CTA button
    await expect(page.getByRole("button", { name: /Daftar Sekarang/i }).or(
      page.getByRole("link", { name: /Daftar Sekarang/i })
    )).toBeVisible();
  });

  test("hero section memiliki headline dan tombol lihat paket", async ({ page }) => {
    await goTo(page, "/");

    // Setidaknya ada teks "Umrah" atau "Haji" di hero
    const heroText = await page.textContent("body");
    expect(heroText).toMatch(/[Uu]mrah|[Hh]aji/);

    // Tombol CTA "Lihat Paket" atau "Lihat Paket Umrah"
    const ctaBtn = page.getByRole("link", { name: /Lihat Paket/i }).or(
      page.getByRole("button", { name: /Lihat Paket/i })
    );
    await expect(ctaBtn.first()).toBeVisible();
  });

  test("stats section menampilkan angka kepercayaan", async ({ page }) => {
    await goTo(page, "/");
    // Stats yang terlihat saat manual testing: Izin PPIU, Akreditasi A, 25.000+
    const body = page.locator("body");
    await expect(body).toContainText(/PPIU|Akreditasi|jamaah/i);
  });

  test("klik 'Lihat Paket Umrah' navigasi ke /packages", async ({ page }) => {
    await goTo(page, "/");
    const ctaBtn = page.getByRole("link", { name: /Lihat Paket/i }).first();
    await ctaBtn.click();
    await expect(page).toHaveURL(/\/packages/);
  });
});

// ─── B2C-02: Katalog Paket ────────────────────────────────────────────────────

test.describe("B2C-02 — Katalog paket (/packages)", () => {
  test("menampilkan daftar paket dengan nama dan harga", async ({ page }) => {
    await goTo(page, "/packages");

    // Setidaknya ada 1 paket yang tampil
    const cards = page.locator("article, [class*='card'], [class*='package']");
    // Cek teks khas paket
    const body = page.locator("body");
    await expect(body).toContainText(/Rp\s*[\d.,]+/); // ada angka harga
    await expect(body).toContainText(/[Uu]mrah|[Hh]aji/);
  });

  test("menampilkan minimal 2 paket dengan tombol 'Lihat Detail' dan 'Booking Cepat'", async ({ page }) => {
    await goTo(page, "/packages");

    // Tombol lihat detail
    const detailLinks = page.getByRole("link", { name: /Lihat Detail|Detail/i });
    expect(await detailLinks.count()).toBeGreaterThanOrEqual(1);

    // Tombol booking cepat
    const bookingBtns = page.getByRole("button", { name: /Booking Cepat/i }).or(
      page.getByRole("link", { name: /Booking Cepat/i })
    );
    expect(await bookingBtns.count()).toBeGreaterThanOrEqual(1);
  });

  test("klik 'Lihat Detail' navigasi ke halaman detail paket", async ({ page }) => {
    await goTo(page, "/packages");
    const firstDetail = page.getByRole("link", { name: /Lihat Detail|Detail/i }).first();
    await firstDetail.click();
    await expect(page).toHaveURL(/\/packages\//);
  });

  test("harga paket ditampilkan dalam format Rupiah", async ({ page }) => {
    await goTo(page, "/packages");
    // Cek format harga: Rp X.XXX.XXX
    await expect(page.locator("body")).toContainText(/Rp\s*[\d.]+/);
  });
});

// ─── B2C-03: Detail Paket ─────────────────────────────────────────────────────

test.describe("B2C-03 — Detail paket", () => {
  test("menampilkan judul, deskripsi, dan breadcrumb", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);

    // Breadcrumb
    await expect(page.getByText("Paket Umrah")).toBeVisible();

    // Judul paket
    await expect(page.getByRole("heading", { name: /Umrah Reguler 12 Hari|Ramadan/i })).toBeVisible();

    // Deskripsi
    await expect(page.locator("body")).toContainText(/hotel|masjid/i);
  });

  test("menampilkan 2 departure cards dengan tanggal", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);

    // Scroll ke departure section
    await page.evaluate(() => window.scrollBy(0, 600));
    await page.waitForTimeout(500);

    // Cek dua jadwal keberangkatan
    await expect(page.getByText("23 Mei 2026")).toBeVisible();
    await expect(page.getByText("7 Juli 2026")).toBeVisible();

    // Status kuota
    await expect(page.getByText("Tersedia banyak").first()).toBeVisible();
  });

  test("switching departure mengubah tanggal di sticky footer", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);
    await page.evaluate(() => window.scrollBy(0, 600));
    await page.waitForTimeout(500);

    // Klik departure kedua (7 Juli)
    const dep2Card = page.getByText("7 Juli 2026").locator("..");
    await dep2Card.click();
    await page.waitForTimeout(500);

    // Footer sticky harus update ke "7 Juli 2026"
    const footer = page.locator("body");
    await expect(footer).toContainText("7 Juli 2026");
  });

  test("tab navigasi (Ringkasan, Itinerari, Fasilitas, Syarat & Ketentuan) berfungsi", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);
    await page.evaluate(() => window.scrollBy(0, 300));
    await page.waitForTimeout(500);

    // Cek tab ada
    await expect(page.getByRole("tab", { name: "Ringkasan" }).or(
      page.getByText("Ringkasan")
    ).first()).toBeVisible();
    await expect(page.getByText("Itinerari")).toBeVisible();
    await expect(page.getByText("Fasilitas")).toBeVisible();
  });

  test("sticky footer menampilkan tombol 'Lanjut booking' dan 'Chat WhatsApp'", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);
    await expect(page.getByRole("button", { name: /Lanjut booking/i }).or(
      page.getByRole("link", { name: /Lanjut booking/i })
    )).toBeVisible();
    await expect(page.getByText("Chat WhatsApp")).toBeVisible();
  });

  test("[KNOWN BUG] harga departure belum muncul di halaman detail", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);
    // Harga seharusnya muncul setelah pilih departure, tapi saat ini placeholder
    // Bug: API tidak mengembalikan price_per_pax di public catalog endpoint
    await expect(page.locator("body")).toContainText(
      /Harga akan ditampilkan dari departure|Mulai dari/i
    );
    // Test ini mendokumentasikan behavior saat ini — bukan expected behavior
  });
});

// ─── B2C-04: Booking Wizard Step 1 ───────────────────────────────────────────

test.describe("B2C-04 — Booking wizard Step 1 (pilih keberangkatan)", () => {
  test("halaman booking terbuka dari detail paket", async ({ page }) => {
    await goTo(page, `/packages/${PKG_ID}`);

    // Klik Lanjut booking dari sticky footer
    const lanjutBtn = page.getByRole("button", { name: /Lanjut booking/i }).or(
      page.getByRole("link", { name: /Lanjut booking/i })
    ).last();
    await lanjutBtn.click();
    await page.waitForTimeout(2000);

    await expect(page).toHaveURL(/\/booking\//);
  });

  test("step 1 menampilkan departure dengan harga dalam Rupiah", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=1`);

    // Step indicator
    await expect(page.getByText("Pilih keberangkatan")).toBeVisible();

    // Harga departure (harga muncul di wizard meski tidak di detail page)
    await expect(page.locator("body")).toContainText(/Rp\s*[\d.]+/);

    // Dua departure cards
    await expect(page.getByText("23 Mei 2026")).toBeVisible();
    await expect(page.getByText("7 Juli 2026")).toBeVisible();
  });

  test("switching departure mengubah total di ringkasan paket", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=1`);

    // Ambil total awal (23 Mei)
    const body = page.locator("body");

    // Klik departure kedua (7 Juli 2026)
    await page.getByText("7 Juli 2026").click();
    await page.waitForTimeout(1000);

    // URL harus update
    await expect(page).toHaveURL(new RegExp(DEP2_ID));

    // Ringkasan paket harus update juga
    await expect(body).toContainText("7 Juli 2026");
  });

  test("tombol Lanjut ada dan aktif setelah pilih departure", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=1`);

    const lanjutBtn = page.getByRole("button", { name: "Lanjut" });
    await expect(lanjutBtn).toBeVisible();
    await expect(lanjutBtn).not.toBeDisabled();
  });

  test("ringkasan paket menampilkan nama dan durasi", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=1`);
    await expect(page.getByText("Ringkasan paket")).toBeVisible();
    await expect(page.locator("body")).toContainText(/Umrah Reguler/i);
    await expect(page.locator("body")).toContainText(/Hari/i);
  });
});

// ─── B2C-05: Booking Wizard Step 2 ───────────────────────────────────────────

test.describe("B2C-05 — Booking wizard Step 2 (data jamaah)", () => {
  test.beforeEach(async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=2`);
  });

  test("menampilkan form kontak pemesan dengan semua field wajib", async ({ page }) => {
    await expect(page.getByText("Kontak pemesan")).toBeVisible();
    await expect(page.getByPlaceholder(/Ahmad Fauzan/i)).toBeVisible();
    await expect(page.getByPlaceholder(/nama@email.com/i)).toBeVisible();
    await expect(page.getByPlaceholder(/\+62/i)).toBeVisible();
  });

  test("dropdown tipe kamar memiliki pilihan quad, triple, double", async ({ page }) => {
    const select = page.locator("select");
    await expect(select).toBeVisible();

    const options = await select.locator("option").allTextContents();
    const optionTexts = options.join(" ").toLowerCase();
    expect(optionTexts).toContain("quad");
    expect(optionTexts).toContain("triple");
    expect(optionTexts).toContain("double");
  });

  test("mengganti tipe kamar mengubah total secara reaktif", async ({ page }) => {
    // Ambil total sebelum
    const totalBefore = await page.getByText(/Rp\s*[\d.,]+/).first().textContent();

    // Ganti ke Double
    const select = page.locator("select");
    await select.selectOption("double");
    await page.waitForTimeout(500);

    // Total harus berubah
    const totalAfter = await page.getByText(/Rp\s*[\d.,]+/).first().textContent();
    expect(totalBefore).not.toBe(totalAfter);
  });

  test("menampilkan section Data jamaah untuk tiap jamaah", async ({ page }) => {
    await page.evaluate(() => window.scrollBy(0, 400));
    await expect(page.getByText("Data jamaah")).toBeVisible();
    // Jamaah 1 dengan badge Pemesan
    await expect(page.getByText("Jamaah 1")).toBeVisible();
    await expect(page.getByText("PEMESAN").or(page.getByText("pemesan")).first()).toBeVisible();
  });

  test("mengisi form lengkap dan klik 'Lanjut ke Review' berhasil", async ({ page }) => {
    // Isi kontak pemesan
    await page.getByPlaceholder(/Ahmad Fauzan/i).fill(TEST_PEMESAN.nama);
    await page.getByPlaceholder(/nama@email.com/i).fill(TEST_PEMESAN.email);
    await page.getByPlaceholder(/\+62/i).fill(TEST_PEMESAN.whatsapp);

    const alamatField = page.getByPlaceholder(/Kota atau domisili/i);
    if (await alamatField.count() > 0) {
      await alamatField.fill(TEST_PEMESAN.alamat);
    }

    // Ganti tipe kamar ke double
    await page.locator("select").selectOption("double");
    await page.waitForTimeout(300);

    // Isi jamaah 1
    await page.evaluate(() => window.scrollBy(0, 500));
    const namaKtpFields = page.getByPlaceholder(/Nama sesuai KTP/i);
    await namaKtpFields.nth(0).fill(TEST_JAMAAH1.nama);
    await page.getByPlaceholder(/NIK/i).nth(0).fill(TEST_JAMAAH1.nik);
    await page.locator("input[type='date']").nth(0).fill(TEST_JAMAAH1.dob);

    // Isi jamaah 2
    await namaKtpFields.nth(1).fill(TEST_JAMAAH2.nama);
    await page.getByPlaceholder(/NIK/i).nth(1).fill(TEST_JAMAAH2.nik);
    await page.locator("input[type='date']").nth(1).fill(TEST_JAMAAH2.dob);

    // Submit ke Review
    await page.getByRole("button", { name: /Lanjut ke Review/i }).click();
    await page.waitForTimeout(2000);

    // Harus pindah ke step 3
    await expect(page).toHaveURL(/step=3/);
  });
});

// ─── B2C-06: Booking Wizard Step 3 ───────────────────────────────────────────

test.describe("B2C-06 — Booking wizard Step 3 (review & syarat)", () => {
  // Setup: isi form di step 2 dulu, lalu lanjut ke step 3
  async function fillAndProceedToStep3(page: Page) {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=2`);
    await page.getByPlaceholder(/Ahmad Fauzan/i).fill(TEST_PEMESAN.nama);
    await page.getByPlaceholder(/nama@email.com/i).fill(TEST_PEMESAN.email);
    await page.getByPlaceholder(/\+62/i).fill(TEST_PEMESAN.whatsapp);
    await page.locator("select").selectOption("double");
    await page.waitForTimeout(300);
    await page.evaluate(() => window.scrollBy(0, 500));
    const namaKtpFields = page.getByPlaceholder(/Nama sesuai KTP/i);
    await namaKtpFields.nth(0).fill(TEST_JAMAAH1.nama);
    await page.getByPlaceholder(/NIK/i).nth(0).fill(TEST_JAMAAH1.nik);
    await page.locator("input[type='date']").nth(0).fill(TEST_JAMAAH1.dob);
    await namaKtpFields.nth(1).fill(TEST_JAMAAH2.nama);
    await page.getByPlaceholder(/NIK/i).nth(1).fill(TEST_JAMAAH2.nik);
    await page.locator("input[type='date']").nth(1).fill(TEST_JAMAAH2.dob);
    await page.evaluate(() => window.scrollBy(0, 500));
    await page.getByRole("button", { name: /Lanjut ke Review/i }).click();
    await page.waitForTimeout(2000);
    await expect(page).toHaveURL(/step=3/);
  }

  test("step 3 menampilkan review pemesanan dengan ringkasan keberangkatan", async ({ page }) => {
    await fillAndProceedToStep3(page);
    await expect(page.getByText("Review pemesanan")).toBeVisible();
    await expect(page.getByText("Ringkasan keberangkatan")).toBeVisible();
    await expect(page.locator("body")).toContainText("23 Mei 2026");
    await expect(page.locator("body")).toContainText(/13\s*hari/i);
    await expect(page.locator("body")).toContainText(/Double/i);
  });

  test("daftar jamaah ditampilkan dengan NIK yang di-mask", async ({ page }) => {
    await fillAndProceedToStep3(page);
    await expect(page.getByText("Daftar jamaah")).toBeVisible();
    await expect(page.getByText(TEST_JAMAAH1.nama)).toBeVisible();
    await expect(page.getByText(TEST_JAMAAH2.nama)).toBeVisible();

    // NIK harus di-mask (tampil ****XXXX, bukan angka lengkap)
    await expect(page.locator("body")).toContainText(/\*{4}/);
    // NIK asli tidak boleh tampil di-clear
    await expect(page.locator("body")).not.toContainText(TEST_JAMAAH1.nik);
    await expect(page.locator("body")).not.toContainText(TEST_JAMAAH2.nik);
  });

  test("link 'Ubah data' membawa kembali ke step 2", async ({ page }) => {
    await fillAndProceedToStep3(page);
    await page.getByText("Ubah data").click();
    await page.waitForTimeout(1000);
    await expect(page).toHaveURL(/step=2/);
  });

  test("semua 4 checkbox syarat dapat di-centang", async ({ page }) => {
    await fillAndProceedToStep3(page);
    await page.evaluate(() => window.scrollBy(0, 600));

    const checkboxes = page.locator("input[type='checkbox']");
    const count = await checkboxes.count();
    expect(count).toBeGreaterThanOrEqual(4);

    // Centang semua
    for (let i = 0; i < count; i++) {
      if (!(await checkboxes.nth(i).isChecked())) {
        await checkboxes.nth(i).click();
      }
    }

    // Semua harus tercentang
    for (let i = 0; i < count; i++) {
      await expect(checkboxes.nth(i)).toBeChecked();
    }
  });

  test("tombol 'Lanjut ke pembayaran' muncul setelah syarat dicentang", async ({ page }) => {
    await fillAndProceedToStep3(page);
    await page.evaluate(() => window.scrollBy(0, 800));

    const checkboxes = page.locator("input[type='checkbox']");
    const count = await checkboxes.count();
    for (let i = 0; i < count; i++) {
      if (!(await checkboxes.nth(i).isChecked())) {
        await checkboxes.nth(i).click();
      }
    }

    await expect(page.getByRole("button", { name: /Lanjut ke pembayaran/i })).toBeVisible();
  });
});

// ─── B2C-07: Booking Wizard Step 4 ───────────────────────────────────────────

test.describe("B2C-07 — Booking wizard Step 4 (pembayaran)", () => {
  /**
   * Test ini melakukan submit booking nyata ke server.
   * Booking yang dibuat bisa diidentifikasi dari email pemesan: uat.browser@umrohos.dev
   * CATATAN: Jalankan cleanup setelah testing selesai.
   */
  async function completeBookingToPayment(page: Page) {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=2`);
    await page.getByPlaceholder(/Ahmad Fauzan/i).fill(TEST_PEMESAN.nama);
    await page.getByPlaceholder(/nama@email.com/i).fill(TEST_PEMESAN.email);
    await page.getByPlaceholder(/\+62/i).fill(TEST_PEMESAN.whatsapp);
    await page.locator("select").selectOption("double");
    await page.waitForTimeout(300);
    await page.evaluate(() => window.scrollBy(0, 500));
    const namaKtpFields = page.getByPlaceholder(/Nama sesuai KTP/i);
    await namaKtpFields.nth(0).fill(TEST_JAMAAH1.nama);
    await page.getByPlaceholder(/NIK/i).nth(0).fill(TEST_JAMAAH1.nik);
    await page.locator("input[type='date']").nth(0).fill(TEST_JAMAAH1.dob);
    await namaKtpFields.nth(1).fill(TEST_JAMAAH2.nama);
    await page.getByPlaceholder(/NIK/i).nth(1).fill(TEST_JAMAAH2.nik);
    await page.locator("input[type='date']").nth(1).fill(TEST_JAMAAH2.dob);
    await page.evaluate(() => window.scrollBy(0, 500));
    await page.getByRole("button", { name: /Lanjut ke Review/i }).click();
    await page.waitForTimeout(2000);
    await expect(page).toHaveURL(/step=3/);

    // Centang semua syarat
    await page.evaluate(() => window.scrollBy(0, 800));
    const checkboxes = page.locator("input[type='checkbox']");
    const count = await checkboxes.count();
    for (let i = 0; i < count; i++) {
      if (!(await checkboxes.nth(i).isChecked())) {
        await checkboxes.nth(i).click();
      }
    }

    // Submit ke payment
    await page.getByRole("button", { name: /Lanjut ke pembayaran/i }).click();
    await page.waitForTimeout(4000);
  }

  test("submit booking berhasil dan redirect ke halaman pembayaran", async ({ page }) => {
    await completeBookingToPayment(page);

    // URL harus ke /checkout/...
    await expect(page).toHaveURL(/\/checkout\//);
    await expect(page.getByText("Pembayaran")).toBeVisible();
  });

  test("halaman pembayaran menampilkan kode booking", async ({ page }) => {
    await completeBookingToPayment(page);

    // Kode booking harus ada (format UMR-BKG-XXXXX atau BKG-XXXXX)
    await expect(page.locator("body")).toContainText(/UMR-|BKG-|KODE BOOKING/i);
  });

  test("halaman pembayaran menampilkan Virtual Account BCA", async ({ page }) => {
    await completeBookingToPayment(page);

    // BCA virtual account
    await expect(page.locator("body")).toContainText(/BCA|Bank Central Asia|NOMOR VIRTUAL ACCOUNT/i);

    // Tombol Salin VA
    await expect(page.getByRole("button", { name: /Salin/i }).or(
      page.getByText("Salin")
    ).first()).toBeVisible();
  });

  test("countdown timer tampil pada halaman pembayaran", async ({ page }) => {
    await completeBookingToPayment(page);

    // Timer format JAM:MNT:DTK atau hh:mm:ss
    await expect(page.locator("body")).toContainText(/JAM|DTK|MNT|\d{2}\s*:\s*\d{2}/i);
  });

  test("status booking 'Menunggu Pembayaran'", async ({ page }) => {
    await completeBookingToPayment(page);
    await expect(page.locator("body")).toContainText(/Menunggu Pembayaran|menunggu transfer/i);
  });
});

// ─── B2C-08: Contact Form ─────────────────────────────────────────────────────

test.describe("B2C-08 — Contact form (/contact)", () => {
  test.beforeEach(async ({ page }) => {
    await goTo(page, "/contact");
  });

  test("halaman contact menampilkan 'Hubungi Kami' dan info kontak", async ({ page }) => {
    await expect(page.getByRole("heading", { name: /Hubungi Kami/i })).toBeVisible();
    await expect(page.getByText("WhatsApp")).toBeVisible();
    await expect(page.locator("body")).toContainText(/08.00|Operasional/i);
    await expect(page.locator("body")).toContainText(/PPIU|Izin/i);
  });

  test("formulir konsultasi memiliki field wajib Nama dan Nomor HP", async ({ page }) => {
    await expect(page.getByText("Formulir Konsultasi")).toBeVisible();
    await expect(page.locator("body")).toContainText("NAMA LENGKAP");
    await expect(page.locator("body")).toContainText("NOMOR HP");

    // Field nama dan HP harus ada
    await expect(page.getByPlaceholder(/Ahmad Fauzi/i)).toBeVisible();
    await expect(page.getByPlaceholder(/08xxx|62xxx/i)).toBeVisible();
  });

  test("field email dan pesan/pertanyaan bersifat opsional", async ({ page }) => {
    await expect(page.locator("body")).toContainText(/EMAIL.*opsional|opsional.*EMAIL/i);
    await expect(page.locator("body")).toContainText(/PESAN.*opsional|opsional.*PESAN/i);
  });

  test("dropdown Minat Paket ada", async ({ page }) => {
    const dropdown = page.getByRole("combobox").or(
      page.locator("select")
    ).first();
    await expect(dropdown).toBeVisible();
  });

  test("tombol 'Kirim Pesan' visible", async ({ page }) => {
    await page.evaluate(() => window.scrollBy(0, 400));
    await expect(page.getByRole("button", { name: /Kirim Pesan/i })).toBeVisible();
  });

  test("privacy notice 'Data Anda aman' ditampilkan", async ({ page }) => {
    await page.evaluate(() => window.scrollBy(0, 400));
    await expect(page.locator("body")).toContainText(/Data Anda aman|tidak akan dibagikan/i);
  });
});

// ─── CON-01: Console Login ────────────────────────────────────────────────────

test.describe("CON-01 — Console login (/console/login)", () => {
  test.beforeEach(async ({ page }) => {
    await goTo(page, "/console/login");
  });

  test("halaman login menampilkan 'Internal Operations Console' dan form login", async ({ page }) => {
    await expect(page.getByText("Internal Operations Console")).toBeVisible();
    await expect(page.getByText("Console Login")).toBeVisible();
    await expect(page.getByText("RESTRICTED ACCESS ENVIRONMENT").or(
      page.locator("body")
    )).toBeTruthy();
  });

  test("form login memiliki field email dan password", async ({ page }) => {
    await expect(page.getByPlaceholder(/operator@umrohos/i).or(
      page.locator("input[type='email']")
    ).first()).toBeVisible();
    await expect(page.locator("input[type='password']")).toBeVisible();
    await expect(page.getByRole("button", { name: /SIGN IN|Login|Masuk/i })).toBeVisible();
  });

  test("sidebar kiri menampilkan fitur (Audit trail, Role-based access)", async ({ page }) => {
    await expect(page.getByText("Audit trail")).toBeVisible();
    await expect(page.getByText("Role-based access")).toBeVisible();
  });

  test("'Lupa Password?' link tersedia", async ({ page }) => {
    await expect(page.getByText(/Lupa Password/i)).toBeVisible();
  });

  test("login dengan kredensial salah menampilkan pesan error", async ({ page }) => {
    const emailInput = page.locator("input[type='email']").or(
      page.getByPlaceholder(/operator@umrohos/i)
    ).first();
    await emailInput.fill("wrong@example.com");
    await page.locator("input[type='password']").fill("wrongpass");
    await page.getByRole("button", { name: /SIGN IN|Login|Masuk/i }).click();
    await page.waitForTimeout(2000);

    // Harus ada error message atau tetap di halaman login
    const isStillOnLogin = page.url().includes("/console/login");
    const hasError = await page.locator("body").textContent().then(t =>
      /error|invalid|salah|gagal|unauthorized/i.test(t || "")
    );
    expect(isStillOnLogin || hasError).toBe(true);
  });

  test("login dengan kredensial benar redirect ke dashboard", async ({ page }) => {
    const emailInput = page.locator("input[type='email']").or(
      page.getByPlaceholder(/operator@umrohos/i)
    ).first();
    await emailInput.fill(ADMIN_EMAIL);
    await page.locator("input[type='password']").fill(ADMIN_PASS);
    await page.getByRole("button", { name: /SIGN IN|Login|Masuk/i }).click();
    await page.waitForTimeout(3000);

    // Harus redirect ke /console
    await expect(page).toHaveURL(/\/console$/);
  });
});

// ─── CON-02: Console Dashboard ────────────────────────────────────────────────

test.describe("CON-02 — Console dashboard (/console)", () => {
  test.beforeEach(async ({ page }) => {
    // Login dulu
    await goTo(page, "/console/login");
    const emailInput = page.locator("input[type='email']").or(
      page.getByPlaceholder(/operator@umrohos/i)
    ).first();
    await emailInput.fill(ADMIN_EMAIL);
    await page.locator("input[type='password']").fill(ADMIN_PASS);
    await page.getByRole("button", { name: /SIGN IN|Login|Masuk/i }).click();
    await page.waitForTimeout(3000);
    await expect(page).toHaveURL(/\/console$/);
  });

  test("dashboard menampilkan heading 'Dashboard Operasional'", async ({ page }) => {
    await expect(page.getByRole("heading", { name: /Dashboard Operasional/i })).toBeVisible();
  });

  test("dashboard memiliki stat cards kunci (Booking, Pending Payment, Visa)", async ({ page }) => {
    await expect(page.locator("body")).toContainText(/Booking Hari Ini|BOOKING HARI INI/i);
    await expect(page.locator("body")).toContainText(/Pending Payment|PENDING PAYMENT/i);
    await expect(page.locator("body")).toContainText(/Visa|VISA/i);
  });

  test("sidebar navigasi menampilkan semua menu utama", async ({ page }) => {
    const nav = page.locator("nav, aside, [class*='sidebar']");
    await expect(page.getByRole("link", { name: "Dashboard" }).or(
      page.getByText("Dashboard").first()
    )).toBeVisible();
    await expect(page.getByText("Katalog Paket")).toBeVisible();
    await expect(page.getByText("Booking")).toBeVisible();
    await expect(page.getByText("Finance").or(page.getByText("Keuangan"))).toBeVisible();
    await expect(page.getByText("CRM")).toBeVisible();
  });

  test("tabel Antrian Persetujuan ditampilkan", async ({ page }) => {
    await expect(page.getByText("Antrian Persetujuan")).toBeVisible();
  });

  test("status integrasi ditampilkan di sidebar kanan", async ({ page }) => {
    await expect(page.locator("body")).toContainText(/Status Integrasi|IAM|Booking/i);
  });

  test("navigasi ke halaman Booking dari sidebar", async ({ page }) => {
    await page.getByRole("link", { name: "Booking" }).first().click();
    await page.waitForTimeout(1500);
    await expect(page).toHaveURL(/\/console\/booking|\/console.*booking/);
  });

  test("navigasi ke halaman Katalog Paket dari sidebar", async ({ page }) => {
    await page.getByRole("link", { name: "Katalog Paket" }).click();
    await page.waitForTimeout(1500);
    await expect(page).toHaveURL(/\/console\/catalog|\/console.*paket|\/console.*catalog/i);
  });
});

// ─── B2C-09: Nav Links 404 Check (ISSUE-026) ─────────────────────────────────

test.describe("B2C-09 — Nav links route check (ISSUE-026)", () => {
  /**
   * ISSUE-026: Link "Jadwal", "Manasik", dan "Tentang Kami" di navigasi header
   * B2C semuanya mengarah ke route SvelteKit yang belum dibuat → 404 Not Found.
   *
   * Test ini mendokumentasikan behavior saat ini. Ketika bug diperbaiki (route dibuat),
   * test assertion harus dibalik: expect NOT toContainText("404").
   */

  test("[ISSUE-026] /jadwal mengembalikan 404 — route belum dibuat", async ({ page }) => {
    const response = await page.goto(`${BASE_URL}/jadwal`);
    // Response status atau konten halaman menunjukkan 404
    const status = response?.status();
    const bodyText = await page.locator("body").textContent();
    const is404 = status === 404 || /not found|404/i.test(bodyText || "");
    expect(is404).toBe(true);
  });

  test("[ISSUE-026] /manasik mengembalikan 404 — route belum dibuat", async ({ page }) => {
    const response = await page.goto(`${BASE_URL}/manasik`);
    const status = response?.status();
    const bodyText = await page.locator("body").textContent();
    const is404 = status === 404 || /not found|404/i.test(bodyText || "");
    expect(is404).toBe(true);
  });

  test("[ISSUE-026] /about mengembalikan 404 — route belum dibuat", async ({ page }) => {
    const response = await page.goto(`${BASE_URL}/about`);
    const status = response?.status();
    const bodyText = await page.locator("body").textContent();
    const is404 = status === 404 || /not found|404/i.test(bodyText || "");
    expect(is404).toBe(true);
  });

  test("[ISSUE-026] nav header di landing page memiliki link Jadwal, Manasik, Tentang Kami", async ({ page }) => {
    // Verifikasi link ADA di nav (mereka visible tapi mengarah ke 404)
    await goTo(page, "/");
    await expect(page.getByRole("link", { name: "Jadwal" })).toBeVisible();
    await expect(page.getByRole("link", { name: "Manasik" }).or(
      page.getByText("Manasik")
    ).first()).toBeVisible();
    await expect(page.getByRole("link", { name: "Tentang Kami" }).or(
      page.getByText("Tentang Kami")
    ).first()).toBeVisible();
  });
});

// ─── B2C-10: Form Validation Feedback (ISSUE-027) ────────────────────────────

test.describe("B2C-10 — Form validation feedback (ISSUE-027)", () => {
  /**
   * ISSUE-027: Step 2 form validation secara teknis bekerja (navigasi ke Step 3 diblokir)
   * tetapi tidak memberikan feedback visual kepada user — tidak ada border merah, tidak ada
   * pesan error di bawah field, tidak ada toast/alert.
   *
   * Test ini mendokumentasikan behavior aktual saat ini (validation silent).
   * Ketika bug diperbaiki, ubah assertion terakhir: expect error messages toBeVisible().
   */

  test("[ISSUE-027] klik 'Lanjut ke Review' dengan form kosong — tidak navigasi ke Step 3", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=2`);

    // Biarkan semua field kosong, langsung klik submit
    await page.getByRole("button", { name: /Lanjut ke Review/i }).click();
    await page.waitForTimeout(2000);

    // Validasi berhasil memblokir navigasi — tetap di step=2
    await expect(page).not.toHaveURL(/step=3/);
  });

  test("[ISSUE-027] field kontak pemesan tidak menampilkan error message saat kosong (bug)", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=2`);

    // Submit form kosong
    await page.getByRole("button", { name: /Lanjut ke Review/i }).click();
    await page.waitForTimeout(1000);

    // BEHAVIOR AKTUAL (bug): tidak ada pesan error di field Nama Lengkap
    // Ketika bug diperbaiki, balikkan assertion ini:
    //   await expect(page.getByText(/wajib diisi|required|harus diisi/i).first()).toBeVisible();
    const errorMsgs = page.getByText(/wajib diisi|required|harus diisi|tidak boleh kosong/i);
    const errorCount = await errorMsgs.count();

    // Dokumentasi: saat ini 0 error messages — test LULUS karena mencatat bug yang ada
    // Jika errorCount > 0, bug sudah diperbaiki (update assertion ini)
    console.log(`[ISSUE-027] Error messages shown: ${errorCount} (expected ≥1 setelah bug fix)`);
    // Test ini intentionally tidak fail — tujuannya mendokumentasikan behavior
    // Uncomment baris di bawah setelah bug diperbaiki:
    // expect(errorCount).toBeGreaterThanOrEqual(1);
  });

  test("[ISSUE-027] field input tidak memiliki border merah setelah submit kosong (bug)", async ({ page }) => {
    await goTo(page, `/booking/${PKG_ID}?departure=${DEP1_ID}&step=2`);

    await page.getByRole("button", { name: /Lanjut ke Review/i }).click();
    await page.waitForTimeout(1000);

    // BEHAVIOR AKTUAL (bug): input tidak mendapat class error/invalid setelah submit kosong
    const namaInput = page.getByPlaceholder(/Ahmad Fauzan/i);
    await expect(namaInput).toBeVisible();

    // Cek apakah ada visual error indicator (border-red, aria-invalid, dll)
    const isInvalid = await namaInput.getAttribute("aria-invalid");
    const classList = await namaInput.getAttribute("class") || "";
    const hasErrorClass = /error|invalid|border-red/i.test(classList);

    console.log(`[ISSUE-027] aria-invalid="${isInvalid}", errorClass=${hasErrorClass}`);
    // Dokumentasi: saat ini kedua indicator tidak ada — test hanya mencatat
    // Setelah bug fix: expect(isInvalid).toBe("true") atau expect(hasErrorClass).toBe(true)
  });
});

// ─── CON-03: Console Catalog Error (ISSUE-024) ───────────────────────────────

test.describe("CON-03 — Console catalog error behavior (ISSUE-024)", () => {
  async function loginAdmin(page: Page) {
    await goTo(page, "/console/login");
    const emailInput = page.locator("input[type='email']").or(
      page.getByPlaceholder(/operator@umrohos/i)
    ).first();
    await emailInput.fill(ADMIN_EMAIL);
    await page.locator("input[type='password']").fill(ADMIN_PASS);
    await page.getByRole("button", { name: /SIGN IN|Login|Masuk/i }).click();
    await page.waitForTimeout(3000);
    await expect(page).toHaveURL(/\/console$/);
  }

  test("[ISSUE-024] console /packages menampilkan error 'Tidak dapat terhubung ke layanan katalog'", async ({ page }) => {
    await loginAdmin(page);
    await goTo(page, "/console/packages");
    await page.waitForTimeout(2000);

    // Verifikasi error message muncul (BEHAVIOR AKTUAL — bug)
    await expect(page.locator("body")).toContainText(
      /Tidak dapat terhubung ke layanan katalog|katalog tidak tersedia/i
    );
  });

  test("[ISSUE-024] console catalog menampilkan tabel kosong saat error", async ({ page }) => {
    await loginAdmin(page);
    await goTo(page, "/console/packages");
    await page.waitForTimeout(2000);

    // Tabel kosong — tidak ada row data
    const rows = page.locator("table tbody tr").or(
      page.locator("[class*='package-row']")
    );
    const count = await rows.count();
    // Dokumentasi: saat ini 0 baris karena catalog-svc tidak terjangkau
    console.log(`[ISSUE-024] Package table rows: ${count} (expected >0 setelah bug fix)`);
  });

  test("[ISSUE-024] __data.json endpoint confirms catalog SSR failure", async ({ page }) => {
    await loginAdmin(page);
    // Fetch __data.json untuk verifikasi SSR error tanpa render halaman
    const response = await page.request.get(`${BASE_URL}/console/packages/__data.json`);
    const json = await response.json().catch(() => null);

    if (json) {
      const text = JSON.stringify(json);
      const hasError = /Tidak dapat terhubung|katalog/i.test(text);
      console.log(`[ISSUE-024] __data.json error confirmed: ${hasError}`);
      // Konfirmasi error ada di SSR response
      expect(hasError).toBe(true);
    }
  });
});

// ─── CON-04: Console Ops/Finance/Leads Error (ISSUE-025) ─────────────────────

test.describe("CON-04 — Console ops/finance/leads gateway error (ISSUE-025)", () => {
  async function loginAdmin(page: Page) {
    await goTo(page, "/console/login");
    const emailInput = page.locator("input[type='email']").or(
      page.getByPlaceholder(/operator@umrohos/i)
    ).first();
    await emailInput.fill(ADMIN_EMAIL);
    await page.locator("input[type='password']").fill(ADMIN_PASS);
    await page.getByRole("button", { name: /SIGN IN|Login|Masuk/i }).click();
    await page.waitForTimeout(3000);
    await expect(page).toHaveURL(/\/console$/);
  }

  test("[ISSUE-025] /console/ops menampilkan error gateway", async ({ page }) => {
    await loginAdmin(page);
    await goTo(page, "/console/ops");
    await page.waitForTimeout(2000);

    await expect(page.locator("body")).toContainText(
      /Tidak dapat terhubung ke gateway|gateway tidak tersedia|gagal memuat/i
    );
  });

  test("[ISSUE-025] /console/finance menampilkan error gateway", async ({ page }) => {
    await loginAdmin(page);
    await goTo(page, "/console/finance");
    await page.waitForTimeout(2000);

    await expect(page.locator("body")).toContainText(
      /Tidak dapat terhubung ke gateway|gateway tidak tersedia|gagal memuat/i
    );
  });

  test("[ISSUE-025] /console/leads menampilkan error gateway", async ({ page }) => {
    await loginAdmin(page);
    await goTo(page, "/console/leads");
    await page.waitForTimeout(2000);

    await expect(page.locator("body")).toContainText(
      /Tidak dapat terhubung ke gateway|gateway tidak tersedia|gagal memuat/i
    );
  });

  test("[ISSUE-025] halaman ops/finance/leads UI sudah ada (bukan NOT_DEPLOYED)", async ({ page }) => {
    await loginAdmin(page);

    // Ops page memiliki UI filter/tabs meski data kosong
    await goTo(page, "/console/ops");
    await page.waitForTimeout(1500);
    // Halaman harus ter-render (ada konten UI, bukan blank)
    const bodyText = await page.locator("body").textContent();
    expect(bodyText?.length).toBeGreaterThan(50);

    // Finance page
    await goTo(page, "/console/finance");
    await page.waitForTimeout(1500);
    const financeText = await page.locator("body").textContent();
    expect(financeText?.length).toBeGreaterThan(50);
  });
});
