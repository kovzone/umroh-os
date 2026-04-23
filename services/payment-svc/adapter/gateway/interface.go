// Package gateway defines the GatewayAdapter interface and shared request/response
// types used by payment-svc to communicate with external payment gateways.
//
// Gateway selection rule (Q013 / F5 spec):
//   - Primary:  Midtrans
//   - Fallback: Xendit — only on 5xx responses or timeout > 10s
//   - No failover on 4xx (client error = don't retry with another gateway)
//
// Each concrete gateway adapter lives in its own sub-package:
//   adapter/gateway/midtrans/   — Midtrans implementation (future)
//   adapter/gateway/xendit/     — Xendit implementation (future)
//   adapter/gateway/mock/       — mock for local dev (MOCK_GATEWAY=true)
//
// All adapters implement this interface; the service layer is gateway-agnostic.
package gateway

import (
	"context"
	"time"
)

// GatewayAdapter is the contract every payment gateway adapter must satisfy.
type GatewayAdapter interface {
	// IssueVA creates a virtual account on the gateway.
	// idempotencyKey must be set to invoice.id so replays are safe.
	IssueVA(ctx context.Context, req IssueVARequest) (*IssueVAResponse, error)

	// GetPaymentStatus queries the gateway for the current payment status of a VA.
	// Used by the reconciliation cron (W5) to catch missed webhooks.
	GetPaymentStatus(ctx context.Context, gatewayVAID string) (*PaymentStatus, error)

	// Refund initiates a refund on the gateway.
	Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error)

	// VerifyWebhookSignature validates the HMAC signature of an incoming webhook
	// payload. Returns nil on success, non-nil on bad/missing signature.
	// The handler must return 401 and skip all business logic on non-nil error.
	VerifyWebhookSignature(payload []byte, signature string) error
}

// ---------------------------------------------------------------------------
// IssueVA
// ---------------------------------------------------------------------------

// IssueVARequest carries the parameters for creating a virtual account.
type IssueVARequest struct {
	// IdempotencyKey must equal invoice.id (UUID string).
	IdempotencyKey string
	// Amount in IDR (already rounded per Q001).
	AmountIDR float64
	// BookingID for gateway reference / display.
	BookingID string
	// BankCode is the preferred bank (empty = gateway default).
	BankCode string
	// ExpiryDuration: how long the VA should remain active (default 24h per Q010).
	ExpiryDuration time.Duration
}

// IssueVAResponse carries the issued VA details.
type IssueVAResponse struct {
	// GatewayVAID is the gateway's internal identifier for this VA.
	GatewayVAID string
	// AccountNumber is the bank account number the customer transfers to.
	AccountNumber string
	// BankCode identifies the issuing bank (e.g. "BCA", "BNI").
	BankCode string
	// ExpiresAt is when the VA will stop accepting payments.
	ExpiresAt time.Time
	// Gateway is the name of the adapter that issued the VA ("midtrans"|"xendit"|"mock").
	Gateway string
}

// ---------------------------------------------------------------------------
// GetPaymentStatus
// ---------------------------------------------------------------------------

// PaymentStatus is the gateway's view of a VA's settlement state.
type PaymentStatus struct {
	// GatewayVAID ties this status to the VA.
	GatewayVAID string
	// Paid is true when the gateway considers the VA settled.
	Paid bool
	// AmountPaidIDR is the total received by the gateway (may differ from invoice
	// if partial payments were made).
	AmountPaidIDR float64
	// GatewayTxnID is the gateway's transaction identifier for dedup purposes.
	GatewayTxnID string
	// Expired is true when the gateway has invalidated the VA.
	Expired bool
}

// ---------------------------------------------------------------------------
// Refund
// ---------------------------------------------------------------------------

// RefundRequest carries the parameters for initiating a refund.
type RefundRequest struct {
	// RefundID is payment-svc's refund.id; used as idempotency key.
	RefundID string
	// GatewayVAID is the VA that received the original payment.
	GatewayVAID string
	// AmountIDR is the amount to refund (may be less than total paid).
	AmountIDR float64
	// Reason is a short human-readable description (for gateway reporting).
	Reason string
}

// RefundResponse is returned by the gateway after initiating a refund.
type RefundResponse struct {
	// GatewayRefundID is the gateway's identifier for this refund.
	GatewayRefundID string
	// Status is the initial status from the gateway ("processing"|"completed").
	Status string
}
