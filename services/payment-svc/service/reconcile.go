// reconcile.go — ReconcileInvoices service method delegates to the usecase/reconcile package.
//
// This thin wrapper satisfies the IPaymentService interface so the gRPC layer and
// the cron can both call reconciliation through a single entry point.

package service

import (
	"context"

	"payment-svc/adapter/booking_grpc_adapter"
	"payment-svc/adapter/gateway"
	"payment-svc/usecase/reconcile"
)

// ReconcileResult mirrors usecase/reconcile.ReconcileResult to avoid exposing
// the usecase package in the IPaymentService interface signature.
type ReconcileResult struct {
	InvoicesChecked   int
	BackfilledEvents  int
	ExpiredVAsMarked  int
	PartialReconciled int
	Errors            []string
}

// ReconcileInvoices runs one reconciliation pass.
// Implements IPaymentService.ReconcileInvoices.
func (s *PaymentService) ReconcileInvoices(ctx context.Context) (*ReconcileResult, error) {
	rec := reconcile.NewReconciler(reconcile.Config{
		Logger:         s.logger,
		Tracer:         s.tracer,
		Store:          s.store,
		PrimaryGateway: s.primaryGateway,
		XenditGateway:  s.xenditGateway,
		MockGateway:    s.mockGateway,
		BookingAdapter: s.bookingAdapter,
	})

	out, err := rec.Run(ctx)
	if err != nil {
		return nil, err
	}
	return &ReconcileResult{
		InvoicesChecked:   out.InvoicesChecked,
		BackfilledEvents:  out.BackfilledEvents,
		ExpiredVAsMarked:  out.ExpiredVAsMarked,
		PartialReconciled: out.PartialReconciled,
		Errors:            out.Errors,
	}, nil
}

// StartReconcileCron launches the hourly reconciliation ticker in a goroutine.
// Call from cmd/start.go after all dependencies are wired.
func (s *PaymentService) StartReconcileCron(ctx context.Context) {
	rec := reconcile.NewReconciler(reconcile.Config{
		Logger:         s.logger,
		Tracer:         s.tracer,
		Store:          s.store,
		PrimaryGateway: s.primaryGateway,
		XenditGateway:  s.xenditGateway,
		MockGateway:    s.mockGateway,
		BookingAdapter: s.bookingAdapter,
	})
	go rec.StartCron(ctx)
}

// Ensure PaymentService satisfies IPaymentService at compile time.
// (The interface check below will fail at compile time if any method is missing.)
var _ IPaymentService = (*PaymentService)(nil)

// Adapter type aliases to avoid import cycle in reconcile.go
// (these types are referenced in payment_service.go fields)
type _ = booking_grpc_adapter.Adapter
type _ = gateway.GatewayAdapter
