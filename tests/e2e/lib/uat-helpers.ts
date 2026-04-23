/**
 * UAT Helpers — factory, cleanup, dan DB utils untuk UAT testing
 *
 * Semua data yang dibuat menggunakan prefix [UAT] agar cleanup SQL bisa targeted.
 */

import { APIRequestContext, request as playwrightRequest } from "@playwright/test";
import { Client } from "pg";
import { createApiClient, ApiClient } from "./api-client";

// ─── Environment ─────────────────────────────────────────────────────────────

export const UAT_ENV = {
  gatewayUrl: process.env.GATEWAY_SVC_URL || "http://216.176.238.161",
  coreWebUrl: process.env.CORE_WEB_URL || "http://216.176.238.161",
  pgUrl:
    process.env.PG_URL ||
    "postgres://postgres:IDL4Ssfdo9ettSaFfleZp4M+3vKA8wX2@216.176.238.161:5432/umrohos?sslmode=disable",
  adminEmail: process.env.UAT_ADMIN_EMAIL || "admin@umrohos.dev",
  adminPassword: process.env.UAT_ADMIN_PASSWORD || "password123",
  // Seed data IDs (dari migration 000009)
  activePkgId:
    process.env.UAT_ACTIVE_PKG_ID || "pkg_01JCDE00000000000000000001",
  draftPkgId:
    process.env.UAT_DRAFT_PKG_ID || "pkg_01JCDE00000000000000000002",
  activeDepId:
    process.env.UAT_ACTIVE_DEP_ID || "dep_01JCDF00000000000000000001",
  adminUserId:
    process.env.UAT_ADMIN_USER_ID || "33333333-3333-3333-3333-333333333333",
};

/** Tag prefix wajib untuk semua data yang dibuat selama UAT */
export const UAT_PREFIX = "[UAT]";

/** Generate email UAT dengan timestamp agar tidak konflik antar run */
export function uatEmail(role: string): string {
  return `uat.${role}.${Date.now()}@umrohos.dev`;
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}

/**
 * Login sebagai admin dan kembalikan access token.
 * Gunakan ini di beforeAll untuk mendapatkan token sekali, lalu share ke tests.
 */
export async function loginAdmin(): Promise<{
  tokens: AuthTokens;
  api: ApiClient;
}> {
  const api = await createApiClient(UAT_ENV.gatewayUrl);
  const res = await api.post("/v1/auth/login", {
    email: UAT_ENV.adminEmail,
    password: UAT_ENV.adminPassword,
  });

  if (res.status() !== 200) {
    const body = await res.text();
    throw new Error(
      `Login admin gagal: HTTP ${res.status()} — ${body}\n` +
        `Cek apakah server up di ${UAT_ENV.gatewayUrl} dan migration seed sudah jalan.`
    );
  }

  const body = await res.json();
  const tokens: AuthTokens = {
    accessToken: body.data.access_token,
    refreshToken: body.data.refresh_token,
  };

  const authedApi = await createApiClient(UAT_ENV.gatewayUrl, tokens.accessToken);
  return { tokens, api: authedApi };
}

// ─── DB Client ───────────────────────────────────────────────────────────────

/** Buat PostgreSQL client. Caller harus panggil connect() dan end(). */
export function createDbClient(): Client {
  return new Client({ connectionString: UAT_ENV.pgUrl });
}

/** Helper: jalankan query dan kembalikan rows */
export async function dbQuery<T = Record<string, unknown>>(
  sql: string,
  params: unknown[] = []
): Promise<T[]> {
  const client = createDbClient();
  await client.connect();
  try {
    const result = await client.query(sql, params);
    return result.rows as T[];
  } finally {
    await client.end();
  }
}

// ─── Test Data Factory ───────────────────────────────────────────────────────

export interface UatPackage {
  id: string;
  name: string;
}

export interface UatDeparture {
  id: string;
  packageId: string;
}

export interface UatBooking {
  id: string;
  packageId: string;
  departureId: string;
}

export interface UatInvoice {
  id: string;
  bookingId: string;
  vaNumber?: string;
}

export interface UatLead {
  id: string;
  email: string;
}

/**
 * Buat package UAT baru via API.
 * Nama otomatis pakai prefix [UAT].
 */
export async function createUatPackage(
  api: ApiClient,
  overrides: Record<string, unknown> = {}
): Promise<UatPackage> {
  const name = `${UAT_PREFIX} Umroh Reguler Test ${Date.now()}`;
  const res = await api.post("/v1/packages", {
    name,
    description: "Package untuk UAT automated testing. Hapus jika masih ada.",
    kind: "umrah_reguler",
    duration_days: 12,
    status: "active",
    ...overrides,
  });

  if (res.status() !== 201 && res.status() !== 200) {
    const body = await res.text();
    throw new Error(`createUatPackage gagal: HTTP ${res.status()} — ${body}`);
  }

  const body = await res.json();
  return { id: body.data?.id || body.id, name };
}

/**
 * Buat departure UAT baru di bawah packageId.
 */
export async function createUatDeparture(
  api: ApiClient,
  packageId: string,
  overrides: Record<string, unknown> = {}
): Promise<UatDeparture> {
  const res = await api.post(`/v1/packages/${packageId}/departures`, {
    depart_date: "2026-12-01",
    return_date: "2026-12-12",
    capacity: 20,
    status: "open",
    notes: `${UAT_PREFIX} departure test`,
    ...overrides,
  });

  if (res.status() !== 201 && res.status() !== 200) {
    const body = await res.text();
    throw new Error(`createUatDeparture gagal: HTTP ${res.status()} — ${body}`);
  }

  const body = await res.json();
  return { id: body.data?.id || body.id, packageId };
}

/**
 * Buat draft booking UAT via B2C channel.
 */
export async function createUatBooking(
  packageId: string,
  departureId: string,
  overrides: Record<string, unknown> = {}
): Promise<UatBooking> {
  const api = await createApiClient(UAT_ENV.gatewayUrl);
  const email = uatEmail("jamaah");
  const res = await api.post("/v1/bookings", {
    channel: "b2c_self",
    package_id: packageId,
    departure_id: departureId,
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
    ...overrides,
  });

  if (res.status() !== 201 && res.status() !== 200) {
    const body = await res.text();
    throw new Error(`createUatBooking gagal: HTTP ${res.status()} — ${body}`);
  }

  const body = await res.json();
  return {
    id: body.data?.id || body.id,
    packageId,
    departureId,
  };
}

/**
 * Create invoice untuk booking (staff call).
 */
export async function createUatInvoice(
  api: ApiClient,
  bookingId: string
): Promise<UatInvoice> {
  const res = await api.post("/v1/invoices", {
    booking_id: bookingId,
    gateway: "mock",
  });

  if (res.status() !== 201 && res.status() !== 200) {
    const body = await res.text();
    throw new Error(`createUatInvoice gagal: HTTP ${res.status()} — ${body}`);
  }

  const body = await res.json();
  return { id: body.data?.id || body.id, bookingId };
}

/**
 * Buat lead UAT via public endpoint.
 */
export async function createUatLead(
  overrides: Record<string, unknown> = {}
): Promise<UatLead> {
  const api = await createApiClient(UAT_ENV.gatewayUrl);
  const email = uatEmail("lead");
  const res = await api.post("/v1/leads", {
    name: `${UAT_PREFIX} Test Lead`,
    email,
    phone: "+628119876543",
    message: "Test lead dari UAT automated testing",
    source_note: "uat-test",
    ...overrides,
  });

  if (res.status() !== 201 && res.status() !== 200) {
    const body = await res.text();
    throw new Error(`createUatLead gagal: HTTP ${res.status()} — ${body}`);
  }

  const body = await res.json();
  return { id: body.data?.id || body.id, email };
}

// ─── Cleanup ─────────────────────────────────────────────────────────────────

/**
 * Hapus package UAT via API (soft-delete / archive).
 * Gagal silent — sudah terhapus atau tidak ada = OK.
 */
export async function deleteUatPackage(
  api: ApiClient,
  packageId: string
): Promise<void> {
  try {
    await api.delete(`/v1/packages/${packageId}`);
  } catch {
    // Silent — already deleted or not found
  }
}

/**
 * Cleanup via DB langsung: hapus semua data [UAT] sekaligus.
 * Lebih efisien daripada call API satu-satu.
 * Panggil ini di afterAll test suite.
 */
export async function cleanupUatData(): Promise<void> {
  const client = createDbClient();
  await client.connect();
  try {
    await client.query("BEGIN");

    // CRM
    await client.query(
      `DELETE FROM crm.lead_status_history WHERE lead_id IN (
        SELECT id FROM crm.leads WHERE email ILIKE 'uat.%@%' OR name ILIKE '[UAT]%'
      )`
    );
    await client.query(
      `DELETE FROM crm.leads WHERE email ILIKE 'uat.%@%' OR name ILIKE '[UAT]%'`
    );

    // Payment
    await client.query(
      `DELETE FROM payment.payment_events WHERE invoice_id IN (
        SELECT pi.id FROM payment.invoices pi
        JOIN booking.bookings b ON b.id = pi.booking_id
        WHERE b.notes ILIKE '%[UAT]%'
      )`
    );
    await client.query(
      `DELETE FROM payment.virtual_accounts WHERE invoice_id IN (
        SELECT pi.id FROM payment.invoices pi
        JOIN booking.bookings b ON b.id = pi.booking_id
        WHERE b.notes ILIKE '%[UAT]%'
      )`
    );
    await client.query(
      `DELETE FROM payment.invoices WHERE booking_id IN (
        SELECT id FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
      )`
    );

    // Finance
    await client.query(
      `DELETE FROM finance.journal_lines WHERE entry_id IN (
        SELECT je.id FROM finance.journal_entries je
        WHERE je.source_id::text IN (
          SELECT pi.id::text FROM payment.invoices pi
          JOIN booking.bookings b ON b.id = pi.booking_id
          WHERE b.notes ILIKE '%[UAT]%'
        )
      )`
    );
    await client.query(
      `DELETE FROM finance.journal_entries WHERE source_id::text IN (
        SELECT pi.id::text FROM payment.invoices pi
        JOIN booking.bookings b ON b.id = pi.booking_id
        WHERE b.notes ILIKE '%[UAT]%'
      )`
    );

    // Logistics
    await client.query(
      `DELETE FROM logistics.fulfillment_tasks WHERE booking_id IN (
        SELECT id FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
      )`
    );

    // Booking
    await client.query(
      `DELETE FROM booking.pilgrim_documents WHERE jamaah_id IN (
        SELECT bj.id FROM booking.jamaah bj
        JOIN booking.bookings b ON b.id = bj.booking_id
        WHERE b.notes ILIKE '%[UAT]%'
      )`
    );
    await client.query(
      `DELETE FROM booking.jamaah WHERE booking_id IN (
        SELECT id FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
      )`
    );
    await client.query(
      `DELETE FROM booking.bookings WHERE notes ILIKE '%[UAT]%'`
    );

    // Catalog
    await client.query(
      `DELETE FROM catalog.departures WHERE package_id IN (
        SELECT id FROM catalog.packages WHERE name ILIKE '[UAT]%'
      )`
    );
    await client.query(
      `DELETE FROM catalog.packages WHERE name ILIKE '[UAT]%'`
    );

    await client.query("COMMIT");
  } catch (err) {
    await client.query("ROLLBACK");
    console.error("UAT cleanup gagal, manual cleanup diperlukan:", err);
    throw err;
  } finally {
    await client.end();
  }
}
