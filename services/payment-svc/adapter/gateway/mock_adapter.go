// mock_adapter.go — mock GatewayAdapter for local development.
//
// Activated when env var MOCK_GATEWAY=true (or config.Gateway.Mock = true).
// Safe for use in unit tests and the local docker-compose stack where real
// Midtrans/Xendit credentials are not available.
//
// Behaviour:
//   - IssueVA         → returns a deterministic dummy VA (bank=BCA, acct=1234567890, +24h TTL)
//   - GetPaymentStatus → always reports not-paid (reconcile cron does nothing in mock mode)
//   - Refund          → acknowledges immediately with a fake gateway_refund_id
//   - VerifyWebhookSignature → always returns nil (accepts any payload in dev)
//
// NEVER deploy this adapter to production.

package gateway

import (
	"context"
	"fmt"
	"time"
)

// MockAdapter implements GatewayAdapter for local development / testing.
// It is goroutine-safe (no mutable state).
type MockAdapter struct{}

// NewMockAdapter constructs the mock adapter.
func NewMockAdapter() GatewayAdapter {
	return &MockAdapter{}
}

// IssueVA returns a hard-coded dummy virtual account.
func (m *MockAdapter) IssueVA(_ context.Context, req IssueVARequest) (*IssueVAResponse, error) {
	return &IssueVAResponse{
		GatewayVAID:   fmt.Sprintf("mock-va-%s", req.IdempotencyKey),
		AccountNumber: "1234567890",
		BankCode:      "BCA",
		ExpiresAt:     time.Now().UTC().Add(24 * time.Hour),
		Gateway:       "mock",
	}, nil
}

// GetPaymentStatus always reports "not yet paid" so the reconciliation cron
// does not accidentally flip invoices in the local stack.
func (m *MockAdapter) GetPaymentStatus(_ context.Context, gatewayVAID string) (*PaymentStatus, error) {
	return &PaymentStatus{
		GatewayVAID:   gatewayVAID,
		Paid:          false,
		AmountPaidIDR: 0,
		GatewayTxnID:  "",
		Expired:       false,
	}, nil
}

// Refund acknowledges immediately with a fake refund ID.
func (m *MockAdapter) Refund(_ context.Context, req RefundRequest) (*RefundResponse, error) {
	return &RefundResponse{
		GatewayRefundID: fmt.Sprintf("mock-refund-%s", req.RefundID),
		Status:          "processing",
	}, nil
}

// VerifyWebhookSignature always passes in mock mode — dev only.
func (m *MockAdapter) VerifyWebhookSignature(_ []byte, _ string) error {
	return nil
}
