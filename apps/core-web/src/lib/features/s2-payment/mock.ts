/**
 * Mock implementations for payment API.
 * Used when VITE_MOCK_GATEWAY=true (client) or MOCK_GATEWAY=true (server).
 *
 * Returns realistic dummy data so the checkout UI can be developed and tested
 * independently of payment-svc (which Elda Agent is building in parallel).
 */

import type { Invoice, InvoiceWithVA } from './types';

function randomId(prefix: string): string {
  return `${prefix}_${Math.random().toString(36).slice(2, 10)}`;
}

/** Simulate network delay */
function delay(ms = 400): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// In-memory store so repeated calls for the same bookingId are idempotent
const invoiceStore = new Map<
  string,
  { invoiceId: string; accountNumber: string; expiresAt: string }
>();

export async function issueInvoiceMock(bookingId: string): Promise<InvoiceWithVA> {
  await delay(600);

  // Idempotency: return the same VA for the same bookingId
  let stored = invoiceStore.get(bookingId);
  if (!stored) {
    const invoiceId = randomId('inv');
    const accountNumber = `8000${Math.floor(Math.random() * 1e10)
      .toString()
      .padStart(10, '0')}`;
    // VA expires in 24 hours per Q010
    const expiresAt = new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString();
    stored = { invoiceId, accountNumber, expiresAt };
    invoiceStore.set(bookingId, stored);
  }

  return {
    invoice_id: stored.invoiceId,
    va: {
      bank_code: 'BCA',
      account_number: stored.accountNumber,
      amount_total: 38_500_000,
      expires_at: stored.expiresAt
    }
  };
}

export async function getInvoiceMock(invoiceId: string): Promise<Invoice> {
  await delay(200);

  // Find stored data by invoiceId
  let foundEntry: { invoiceId: string; accountNumber: string; expiresAt: string } | undefined;
  let foundBookingId = 'bkg_demo';

  for (const [bookingId, entry] of invoiceStore) {
    if (entry.invoiceId === invoiceId) {
      foundEntry = entry;
      foundBookingId = bookingId;
      break;
    }
  }

  const accountNumber = foundEntry?.accountNumber ?? '800012345678901';
  const expiresAt = foundEntry?.expiresAt ?? new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString();

  return {
    id: invoiceId,
    booking_id: foundBookingId,
    status: 'unpaid',
    amount_total: 38_500_000,
    paid_amount: 0,
    currency: 'IDR',
    va: {
      bank_code: 'BCA',
      account_number: accountNumber,
      amount_total: 38_500_000,
      expires_at: expiresAt
    },
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  };
}
