// xendit_adapter.go — Xendit gateway adapter stub for payment-svc.
//
// This is a STUB implementation. All methods return errors until real API keys
// are configured and the implementation is completed.
//
// To make this production-ready:
//  1. Set XENDIT_CALLBACK_TOKEN env var (or config.gateway.xendit_callback_token)
//  2. Set XENDIT_SECRET_KEY env var (for API calls, distinct from callback token)
//  3. Set XENDIT_BASE_URL (default: https://api.xendit.co)
//  4. Implement IssueVA using Xendit Fixed VA API
//     Docs: https://developers.xendit.co/api-reference/#create-fixed-virtual-account
//  5. Implement GetPaymentStatus using Xendit Get VA API
//     Docs: https://developers.xendit.co/api-reference/#get-fixed-virtual-account
//  6. Implement Refund using Xendit Refund API (if applicable)
//  7. VerifyWebhookSignature: Xendit uses a static callback token comparison
//     (not per-request HMAC). Compare X-CALLBACK-TOKEN header against configured token.
//
// TODO(REAL_API_KEY): Requires XENDIT_CALLBACK_TOKEN and XENDIT_SECRET_KEY to function.
// Used as fallback gateway when Midtrans returns 5xx or times out (Q013).

package gateway

import (
	"context"
	"errors"
	"fmt"
)

// XenditAdapter implements GatewayAdapter for Xendit.
// It is the fallback gateway (Q013: only on Midtrans 5xx or timeout).
type XenditAdapter struct {
	// callbackToken is used to verify incoming webhooks from Xendit.
	// Source: env var XENDIT_CALLBACK_TOKEN.
	// TODO(REAL_API_KEY): inject from config.
	callbackToken string

	// secretKey is used for outgoing API calls to Xendit.
	// Source: env var XENDIT_SECRET_KEY.
	// TODO(REAL_API_KEY): inject from config.
	secretKey string

	// baseURL is the Xendit API base URL.
	// Default: https://api.xendit.co
	baseURL string
}

// NewXenditAdapter constructs a XenditAdapter.
// callbackToken is required for webhook verification.
// secretKey is required for API calls (VA issuance, status, refund).
// TODO(REAL_API_KEY): call with config.Gateway.XenditCallbackToken and XenditSecretKey from cmd/start.go.
func NewXenditAdapter(callbackToken, secretKey, baseURL string) GatewayAdapter {
	if baseURL == "" {
		baseURL = "https://api.xendit.co"
	}
	return &XenditAdapter{
		callbackToken: callbackToken,
		secretKey:     secretKey,
		baseURL:       baseURL,
	}
}

// IssueVA creates a fixed virtual account on Xendit.
// TODO(REAL_API_KEY): implement using Xendit Fixed Virtual Account API.
// Endpoint: POST /callback_virtual_accounts
func (x *XenditAdapter) IssueVA(_ context.Context, req IssueVARequest) (*IssueVAResponse, error) {
	if x.secretKey == "" {
		return nil, errors.New("xendit: XENDIT_SECRET_KEY not configured — set env var and restart")
	}
	// TODO(REAL_API_KEY): implement Xendit VA creation.
	// Steps:
	//   1. Build Xendit Fixed VA request (external_id=req.IdempotencyKey, bank_code, expected_amount)
	//   2. POST to x.baseURL + "/callback_virtual_accounts" with Basic auth (secretKey:)
	//   3. Parse response: account_number, expiration_date
	//   4. Return IssueVAResponse with gateway="xendit"
	return nil, fmt.Errorf("xendit IssueVA: not yet implemented (TODO: use secretKey, idempotency_key=%q)", req.IdempotencyKey)
}

// GetPaymentStatus queries Xendit for the current payment status of a VA.
// TODO(REAL_API_KEY): implement using Xendit Get Fixed Virtual Account API.
// Endpoint: GET /callback_virtual_accounts/{id}
func (x *XenditAdapter) GetPaymentStatus(_ context.Context, gatewayVAID string) (*PaymentStatus, error) {
	if x.secretKey == "" {
		return nil, errors.New("xendit: XENDIT_SECRET_KEY not configured")
	}
	// TODO(REAL_API_KEY): implement Xendit status check.
	// Steps:
	//   1. GET x.baseURL + "/callback_virtual_accounts/" + gatewayVAID
	//   2. Check status field: "ACTIVE" → not paid, "INACTIVE" / "EXPIRED" → Expired=true
	//   3. Check paid_amount in payment history
	return nil, fmt.Errorf("xendit GetPaymentStatus: not yet implemented (gatewayVAID=%q)", gatewayVAID)
}

// Refund initiates a refund on Xendit.
// TODO(REAL_API_KEY): implement using Xendit Refund API.
// Note: Xendit refund availability depends on payment method used.
func (x *XenditAdapter) Refund(_ context.Context, req RefundRequest) (*RefundResponse, error) {
	if x.secretKey == "" {
		return nil, errors.New("xendit: XENDIT_SECRET_KEY not configured")
	}
	// TODO(REAL_API_KEY): implement Xendit refund.
	return nil, fmt.Errorf("xendit Refund: not yet implemented (refund_id=%q gateway_va_id=%q)", req.RefundID, req.GatewayVAID)
}

// VerifyWebhookSignature validates the Xendit webhook callback token.
//
// Xendit uses a static token model: the incoming X-CALLBACK-TOKEN header value
// is compared directly against the configured XENDIT_CALLBACK_TOKEN.
// Unlike Midtrans, there is no per-request HMAC computation.
//
// Per S2-J-02 contract: bad/missing token → return error → handler returns 401.
func (x *XenditAdapter) VerifyWebhookSignature(_ []byte, signature string) error {
	if x.callbackToken == "" {
		return errors.New("xendit: XENDIT_CALLBACK_TOKEN not configured — cannot verify webhook")
	}
	if signature == "" {
		return errors.New("xendit: missing X-CALLBACK-TOKEN header")
	}
	// Xendit: direct constant-time comparison of the callback token.
	// The comparison is not HMAC-based — just token equality.
	if signature != x.callbackToken {
		return errors.New("xendit: invalid callback token — webhook rejected")
	}
	return nil
}
