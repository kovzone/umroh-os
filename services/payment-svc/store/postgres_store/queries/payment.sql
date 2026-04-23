-- payment.sql — sqlc queries for payment-svc core tables (S2-E-02).
--
-- Tables live in the `payment` schema (000012 migration).
-- All queries use named parameters ($1, $2, ...) per pgx/v5 convention.

-- ===========================================================================
-- invoices
-- ===========================================================================

-- name: CreateInvoice :one
INSERT INTO payment.invoices (
    booking_id,
    amount_total,
    rounding_adjustment_idr,
    currency,
    fx_snapshot,
    status,
    paid_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, booking_id, amount_total, rounding_adjustment_idr, currency, fx_snapshot, status, paid_amount, created_at, updated_at;

-- name: GetInvoiceByID :one
SELECT id, booking_id, amount_total, rounding_adjustment_idr, currency, fx_snapshot, status, paid_amount, created_at, updated_at
FROM payment.invoices
WHERE id = $1;

-- name: GetInvoiceByBookingID :one
SELECT id, booking_id, amount_total, rounding_adjustment_idr, currency, fx_snapshot, status, paid_amount, created_at, updated_at
FROM payment.invoices
WHERE booking_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateInvoicePaidAmount :one
UPDATE payment.invoices
SET
    paid_amount = $2,
    status      = $3,
    updated_at  = NOW()
WHERE id = $1
RETURNING id, booking_id, amount_total, rounding_adjustment_idr, currency, fx_snapshot, status, paid_amount, created_at, updated_at;

-- name: UpdateInvoiceStatus :one
UPDATE payment.invoices
SET
    status     = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, booking_id, amount_total, rounding_adjustment_idr, currency, fx_snapshot, status, paid_amount, created_at, updated_at;

-- name: ListUnpaidInvoicesWithActiveVA :many
-- Used by reconciliation cron: invoices still open that have at least one
-- non-expired VA.
SELECT DISTINCT ON (i.id)
    i.id, i.booking_id, i.amount_total, i.rounding_adjustment_idr, i.currency,
    i.fx_snapshot, i.status, i.paid_amount, i.created_at, i.updated_at
FROM payment.invoices i
JOIN payment.virtual_accounts va ON va.invoice_id = i.id
WHERE i.status IN ('unpaid', 'partially_paid')
  AND va.status = 'active'
  AND va.expires_at > NOW()
ORDER BY i.id, i.created_at ASC;

-- ===========================================================================
-- virtual_accounts
-- ===========================================================================

-- name: CreateVirtualAccount :one
INSERT INTO payment.virtual_accounts (
    invoice_id,
    gateway,
    gateway_va_id,
    account_number,
    bank_code,
    expires_at,
    status,
    idempotency_key
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, invoice_id, gateway, gateway_va_id, account_number, bank_code, expires_at, status, idempotency_key, created_at;

-- name: GetVAByInvoiceID :one
SELECT id, invoice_id, gateway, gateway_va_id, account_number, bank_code, expires_at, status, idempotency_key, created_at
FROM payment.virtual_accounts
WHERE invoice_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetVAByIdempotencyKey :one
SELECT id, invoice_id, gateway, gateway_va_id, account_number, bank_code, expires_at, status, idempotency_key, created_at
FROM payment.virtual_accounts
WHERE idempotency_key = $1;

-- name: UpdateVAStatus :one
UPDATE payment.virtual_accounts
SET status = $2
WHERE id = $1
RETURNING id, invoice_id, gateway, gateway_va_id, account_number, bank_code, expires_at, status, idempotency_key, created_at;

-- name: ListActiveVAsByInvoiceID :many
SELECT id, invoice_id, gateway, gateway_va_id, account_number, bank_code, expires_at, status, idempotency_key, created_at
FROM payment.virtual_accounts
WHERE invoice_id = $1
  AND status = 'active';

-- ===========================================================================
-- payment_events
-- ===========================================================================

-- name: CreatePaymentEvent :one
INSERT INTO payment.payment_events (
    invoice_id,
    gateway,
    gateway_txn_id,
    kind,
    amount,
    raw_payload,
    approval_status,
    received_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, invoice_id, gateway, gateway_txn_id, kind, amount, raw_payload, approval_status, received_at, created_at;

-- name: GetPaymentEventByGatewayTxnID :one
SELECT id, invoice_id, gateway, gateway_txn_id, kind, amount, raw_payload, approval_status, received_at, created_at
FROM payment.payment_events
WHERE gateway = $1
  AND gateway_txn_id = $2;

-- name: ListPaymentEventsByInvoiceID :many
SELECT id, invoice_id, gateway, gateway_txn_id, kind, amount, raw_payload, approval_status, received_at, created_at
FROM payment.payment_events
WHERE invoice_id = $1
ORDER BY received_at ASC;

-- ===========================================================================
-- refunds
-- ===========================================================================

-- name: CreateRefund :one
INSERT INTO payment.refunds (
    invoice_id,
    booking_id,
    amount,
    reason_code,
    status,
    gateway_refund_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, invoice_id, booking_id, amount, reason_code, status, gateway_refund_id, created_at, updated_at;

-- name: UpdateRefundStatus :one
UPDATE payment.refunds
SET
    status            = $2,
    gateway_refund_id = COALESCE($3, gateway_refund_id),
    updated_at        = NOW()
WHERE id = $1
RETURNING id, invoice_id, booking_id, amount, reason_code, status, gateway_refund_id, created_at, updated_at;

-- name: GetRefundByID :one
SELECT id, invoice_id, booking_id, amount, reason_code, status, gateway_refund_id, created_at, updated_at
FROM payment.refunds
WHERE id = $1;
