// midtrans_adapter.go — Midtrans gateway adapter stub for payment-svc.
//
// This is a STUB implementation. All methods are skeletons that return errors
// indicating they need real API keys to function.
//
// To make this production-ready:
//  1. Set MIDTRANS_SERVER_KEY env var (or config.gateway.midtrans_server_key)
//  2. Set MIDTRANS_BASE_URL (default: https://api.midtrans.com)
//  3. Implement IssueVA using Midtrans Bank Transfer / BCA VA API
//     Docs: https://docs.midtrans.com/reference/bank-transfer-object
//  4. Implement GetPaymentStatus using Midtrans Transaction Status API
//     Docs: https://docs.midtrans.com/reference/transaction-status-object
//  5. Implement Refund using Midtrans Refund API
//     Docs: https://docs.midtrans.com/reference/refund-object
//  6. Implement VerifyWebhookSignature using SHA512 HMAC
//     Algorithm: SHA512(order_id + status_code + gross_amount + server_key)
//
// TODO(REAL_API_KEY): Requires MIDTRANS_SERVER_KEY to function.
// All non-mock production flows go through this adapter.

package gateway

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
)

// MidtransAdapter implements GatewayAdapter for Midtrans.
// It is the primary gateway (Q013 gateway selection rule).
type MidtransAdapter struct {
	// serverKey is the Midtrans server key used for signature verification and API calls.
	// Source: env var MIDTRANS_SERVER_KEY.
	// TODO(REAL_API_KEY): inject from config.
	serverKey string

	// baseURL is the Midtrans API base URL.
	// Production: https://api.midtrans.com
	// Sandbox:    https://api.sandbox.midtrans.com
	baseURL string
}

// NewMidtransAdapter constructs a MidtransAdapter.
// serverKey must be non-empty for production use.
// TODO(REAL_API_KEY): call with config.Gateway.MidtransServerKey from cmd/start.go.
func NewMidtransAdapter(serverKey, baseURL string) GatewayAdapter {
	if baseURL == "" {
		baseURL = "https://api.midtrans.com"
	}
	return &MidtransAdapter{
		serverKey: serverKey,
		baseURL:   baseURL,
	}
}

// IssueVA creates a virtual account on Midtrans (BCA or configurable bank).
// TODO(REAL_API_KEY): implement using Midtrans Bank Transfer Create API.
// Endpoint: POST /v2/charge (payment_type: bank_transfer)
func (m *MidtransAdapter) IssueVA(_ context.Context, req IssueVARequest) (*IssueVAResponse, error) {
	if m.serverKey == "" {
		return nil, errors.New("midtrans: MIDTRANS_SERVER_KEY not configured — set env var and restart")
	}
	// TODO(REAL_API_KEY): implement Midtrans VA creation.
	// Steps:
	//   1. Build Midtrans charge request (transaction_details + bank_transfer)
	//   2. POST to m.baseURL + "/v2/charge" with Basic auth (serverKey:)
	//   3. Parse response: va_numbers[0].va_number = account_number
	//   4. Return IssueVAResponse with gateway="midtrans"
	return nil, fmt.Errorf("midtrans IssueVA: not yet implemented (TODO: use serverKey=%q baseURL=%q idempotency_key=%q)",
		m.serverKey[:minInt(4, len(m.serverKey))]+"...", m.baseURL, req.IdempotencyKey)
}

// GetPaymentStatus queries Midtrans for the current payment status of a VA.
// TODO(REAL_API_KEY): implement using Midtrans Get Status API.
// Endpoint: GET /v2/{order_id}/status
func (m *MidtransAdapter) GetPaymentStatus(_ context.Context, gatewayVAID string) (*PaymentStatus, error) {
	if m.serverKey == "" {
		return nil, errors.New("midtrans: MIDTRANS_SERVER_KEY not configured")
	}
	// TODO(REAL_API_KEY): implement Midtrans status check.
	// Steps:
	//   1. GET m.baseURL + "/v2/" + gatewayVAID + "/status"
	//   2. Parse transaction_status: "settlement" → Paid=true, "expire" → Expired=true
	//   3. Return PaymentStatus
	return nil, fmt.Errorf("midtrans GetPaymentStatus: not yet implemented (gatewayVAID=%q)", gatewayVAID)
}

// Refund initiates a refund on Midtrans.
// TODO(REAL_API_KEY): implement using Midtrans Refund API.
// Endpoint: POST /v2/{order_id}/refund
func (m *MidtransAdapter) Refund(_ context.Context, req RefundRequest) (*RefundResponse, error) {
	if m.serverKey == "" {
		return nil, errors.New("midtrans: MIDTRANS_SERVER_KEY not configured")
	}
	// TODO(REAL_API_KEY): implement Midtrans refund.
	return nil, fmt.Errorf("midtrans Refund: not yet implemented (refund_id=%q gateway_va_id=%q)", req.RefundID, req.GatewayVAID)
}

// VerifyWebhookSignature validates the Midtrans webhook signature.
//
// Algorithm (per Midtrans docs):
//   signature = SHA512(order_id + status_code + gross_amount + server_key)
//
// The signature is sent in X-Callback-Token header (as hex string).
// This method validates the signature from the raw payload.
func (m *MidtransAdapter) VerifyWebhookSignature(payload []byte, signature string) error {
	if m.serverKey == "" {
		return errors.New("midtrans: MIDTRANS_SERVER_KEY not configured — cannot verify signature")
	}
	if signature == "" {
		return errors.New("midtrans: missing signature header (X-Callback-Token)")
	}

	// TODO(REAL_API_KEY): extract order_id, status_code, gross_amount from payload
	// and compute SHA512(order_id + status_code + gross_amount + serverKey).
	// For now, provide the computation skeleton.
	_ = sha512.New()
	// Real implementation:
	//   var raw map[string]interface{}
	//   json.Unmarshal(payload, &raw)
	//   orderID     := raw["order_id"].(string)
	//   statusCode  := raw["status_code"].(string)
	//   grossAmount := raw["gross_amount"].(string)
	//   h := sha512.New()
	//   h.Write([]byte(orderID + statusCode + grossAmount + m.serverKey))
	//   computed := hex.EncodeToString(h.Sum(nil))
	//   if !hmac.Equal([]byte(computed), []byte(signature)) { return error }
	_ = payload

	return fmt.Errorf("midtrans VerifyWebhookSignature: not yet implemented — set MIDTRANS_SERVER_KEY and implement signature check")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
