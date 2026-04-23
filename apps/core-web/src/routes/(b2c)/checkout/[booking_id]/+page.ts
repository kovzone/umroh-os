/**
 * S2-L-03: Wire booking → payment call.
 *
 * Server (or client) load function:
 * - Accepts booking_id from route params
 * - Issues VA via POST /v1/invoices (idempotent — safe to call on each page load)
 * - Returns invoice data for the page component
 *
 * Using +page.ts (client-side load) so the countdown timer and polling
 * work correctly without SSR hydration mismatches on time-sensitive state.
 * The initial issuePaymentInvoice call is safe to run in the browser.
 */

import type { PageLoad } from './$types';
import { issuePaymentInvoice } from '$lib/features/s2-payment/repository';
import type { InvoiceWithVA } from '$lib/features/s2-payment/types';

export const ssr = false;

export const load: PageLoad = async ({ params }) => {
  const bookingId = params.booking_id;

  let initialInvoice: InvoiceWithVA | null = null;
  let issueError: string | null = null;

  try {
    initialInvoice = await issuePaymentInvoice(bookingId);
  } catch (err) {
    const e = err as Error;
    issueError = e.message ?? 'Gagal menerbitkan Virtual Account.';
  }

  return {
    bookingId,
    initialInvoice,
    issueError
  };
};
