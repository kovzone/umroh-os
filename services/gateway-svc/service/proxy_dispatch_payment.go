// proxy_dispatch_payment.go — gateway service dispatch for payment-svc routes.
//
// Covers:
//   - ReissuePaymentLink   (BL-PAY-020)
//   - IssueVirtualAccount  (BL-PAY-001 / ISSUE-005) POST /v1/invoices
//   - GetInvoiceByID       (BL-PAY-001 / ISSUE-005) GET /v1/invoices/:id
//   - ProcessWebhook       (BL-PAY-003/004 / ISSUE-007/008) POST /v1/webhooks/*
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

// IssueVirtualAccount delegates to payment-svc.IssueVirtualAccount.
// Creates an invoice and a virtual account for a booking (POST /v1/invoices).
func (s *Service) IssueVirtualAccount(ctx context.Context, params *payment_grpc_adapter.IssueVirtualAccountParams) (*payment_grpc_adapter.IssueVirtualAccountResult, error) {
	return s.adapters.paymentGrpc.IssueVirtualAccount(ctx, params)
}

// GetInvoiceByID delegates to payment-svc.GetInvoiceByID.
// Fetches a single invoice by UUID (GET /v1/invoices/:id).
func (s *Service) GetInvoiceByID(ctx context.Context, params *payment_grpc_adapter.GetInvoiceByIDParams) (*payment_grpc_adapter.GetInvoiceByIDResult, error) {
	return s.adapters.paymentGrpc.GetInvoiceByID(ctx, params)
}

// ProcessWebhook delegates to payment-svc.ProcessWebhook.
// Forwards a raw gateway webhook payload (POST /v1/webhooks/*).
func (s *Service) ProcessWebhook(ctx context.Context, params *payment_grpc_adapter.ProcessWebhookParams) (*payment_grpc_adapter.ProcessWebhookResult, error) {
	return s.adapters.paymentGrpc.ProcessWebhook(ctx, params)
}
