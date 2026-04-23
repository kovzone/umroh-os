/**
 * Typed API wrapper for payment endpoints.
 * All requests route through gateway-svc:4000 — no direct calls to payment-svc.
 *
 * MOCK_GATEWAY=true (server-side env) or VITE_MOCK_GATEWAY=true (client-side)
 * makes the backend return dummy VA data without hitting a real payment gateway.
 *
 * NOTE: openapi-fetch schema for payment endpoints pending S2-E-01 schema
 * regeneration. This manual wrapper bridges the gap (same pattern as
 * src/lib/features/s1-catalog/write-api.ts).
 *
 * All public functions accept an optional `baseUrl` so that server-side
 * callers (+page.server.ts) can pass process.env values instead of relying
 * on import.meta.env (Vite/browser-only).
 */

import type { Gateway, Invoice, InvoiceWithVA } from './types';

// ---- Error handling ----

export class PaymentApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string
  ) {
    super(message);
    this.name = 'PaymentApiError';
  }
}

// ---- Internal helpers ----

function resolveBase(override?: string): string {
  if (override) return override;
  // Server-side (Node): try process.env
  if (typeof process !== 'undefined' && process.env) {
    return (
      process.env['VITE_PAYMENT_API_BASE_URL'] ??
      process.env['VITE_GATEWAY_URL'] ??
      process.env['GATEWAY_URL'] ??
      'http://localhost:4000'
    );
  }
  // Browser: Vite env
  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const env = (import.meta as any).env ?? {};
    return env['VITE_PAYMENT_API_BASE_URL'] ?? env['VITE_GATEWAY_URL'] ?? 'http://localhost:4000';
  } catch {
    return 'http://localhost:4000';
  }
}

async function apiFetch<T>(
  path: string,
  method: 'GET' | 'POST',
  body?: unknown,
  baseUrl?: string
): Promise<T> {
  const url = `${resolveBase(baseUrl)}${path}`;
  const response = await fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json'
    },
    ...(body !== undefined ? { body: JSON.stringify(body) } : {})
  });

  const json = (await response.json().catch(() => ({}))) as {
    error?: { code?: string; message?: string };
  };

  if (!response.ok) {
    const code = json.error?.code ?? 'request_failed';
    const message = json.error?.message ?? `HTTP ${response.status}`;
    throw new PaymentApiError(response.status, code, message);
  }

  return json as T;
}

// ---- Public API ----

/**
 * POST /v1/invoices — issue invoice + VA for a booking.
 *
 * Idempotent: calling twice for the same booking_id returns the existing
 * invoice+VA (backend enforces this via the idempotency constraint).
 *
 * @param bookingId   The draft booking to issue payment for.
 * @param gatewayPref Optional gateway preference ('midtrans' | 'xendit').
 *                    Backend defaults to Midtrans → Xendit failover per Q013.
 * @param baseUrl     Override gateway base URL (server-side callers).
 */
export async function issueInvoice(
  bookingId: string,
  gatewayPref?: Gateway,
  baseUrl?: string
): Promise<InvoiceWithVA> {
  const body: { booking_id: string; gateway_pref?: Gateway } = { booking_id: bookingId };
  if (gatewayPref) body.gateway_pref = gatewayPref;

  return apiFetch<InvoiceWithVA>('/v1/invoices', 'POST', body, baseUrl);
}

/**
 * GET /v1/invoices/{id} — poll invoice + payment status.
 *
 * Used by the checkout page to detect status transitions:
 *   unpaid → partially_paid → paid
 * Also returns the VA so we can detect expiry from va.expires_at.
 *
 * @param invoiceId The invoice to poll.
 * @param baseUrl   Override gateway base URL (server-side callers).
 */
export async function getInvoice(invoiceId: string, baseUrl?: string): Promise<Invoice> {
  return apiFetch<Invoice>(`/v1/invoices/${invoiceId}`, 'GET', undefined, baseUrl);
}
