// proxy_dispatch_payment.go — gateway service dispatch for payment-svc routes.
//
// Covers:
//   - payment-svc: ReissuePaymentLink (BL-PAY-020)
//
// Each method is a thin delegation to the payment_grpc_adapter.
// No business logic lives here; all logic lives in payment-svc.

package service

import (
	"context"

	"gateway-svc/adapter/payment_grpc_adapter"
)

// ReissuePaymentLink delegates to payment-svc.ReissuePaymentLink.
// CS-facing: retrieves or re-issues the VA link for an existing booking.
func (s *Service) ReissuePaymentLink(ctx context.Context, params *payment_grpc_adapter.ReissuePaymentLinkParams) (*payment_grpc_adapter.ReissuePaymentLinkResult, error) {
	return s.adapters.paymentGrpc.ReissuePaymentLink(ctx, params)
}
