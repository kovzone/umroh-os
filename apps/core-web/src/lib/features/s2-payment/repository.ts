/**
 * Payment repository — adapter layer between checkout UI and payment-api.
 * Mirrors the pattern in s1-catalog/repository.ts:
 *   - VITE_MOCK_GATEWAY=true  → mock data (default during dev)
 *   - VITE_MOCK_GATEWAY=false → real gateway-svc calls with mock fallback
 */

import { issueInvoice, getInvoice } from './payment-api';
import { issueInvoiceMock, getInvoiceMock } from './mock';
import type { Invoice, InvoiceWithVA } from './types';

function useMock(): boolean {
  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const env = (import.meta as any).env ?? {};
    return (env['VITE_MOCK_GATEWAY'] ?? 'true') === 'true';
  } catch {
    return true;
  }
}

export async function issuePaymentInvoice(bookingId: string): Promise<InvoiceWithVA> {
  if (useMock()) {
    return issueInvoiceMock(bookingId);
  }

  try {
    return await issueInvoice(bookingId);
  } catch {
    return issueInvoiceMock(bookingId);
  }
}

export async function pollInvoice(invoiceId: string): Promise<Invoice> {
  if (useMock()) {
    return getInvoiceMock(invoiceId);
  }

  try {
    return await getInvoice(invoiceId);
  } catch {
    return getInvoiceMock(invoiceId);
  }
}
