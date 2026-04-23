/**
 * Payment domain types for S2 checkout flow.
 * Aligned with F5 spec: docs/06-features/05-payment-and-reconciliation.md
 */

export type Gateway = 'midtrans' | 'xendit';

export type InvoiceStatus = 'unpaid' | 'partially_paid' | 'paid' | 'void' | 'refunded';

export type VirtualAccount = {
  bank_code: string;
  account_number: string;
  amount_total: number;
  expires_at: string; // ISO 8601
};

/** Returned by POST /v1/invoices */
export type InvoiceWithVA = {
  invoice_id: string;
  va: VirtualAccount;
};

/** Returned by GET /v1/invoices/{id} */
export type Invoice = {
  id: string;
  booking_id: string;
  status: InvoiceStatus;
  amount_total: number;
  paid_amount: number;
  currency: string;
  va: VirtualAccount | null;
  created_at: string;
  updated_at: string;
};

/** UI-level checkout status (derived from Invoice.status + VA expiry) */
export type CheckoutStatus =
  | 'loading'
  | 'waiting_payment'
  | 'partially_paid'
  | 'paid'
  | 'expired'
  | 'error';
