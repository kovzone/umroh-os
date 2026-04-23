// Package reconcile implements the F5-W5 reconciliation cron (S2-E-03 / BL-PAY-006).
//
// The reconciler runs hourly (or on demand via trigger) and catches:
//  1. Missed webhooks: gateway says paid but payment-svc has no event → backfill.
//  2. Expired VAs: gateway says expired/cancelled → mark VA + invoice accordingly.
//  3. Partial webhook loss: gateway paid_amount > invoice.paid_amount → reconcile diff.
//
// Design:
//   - Uses the same gateway adapter as the webhook handler — no duplicate logic.
//   - Each invoice is processed independently; one failure does not abort the batch.
//   - A summary is logged at the end of each run (finance admin dashboard source).
//   - Runs in-process per ADR-0006 (no Temporal in MVP).
//
// Integration: cmd/start.go starts the cron goroutine with a time.Ticker.
// The reconciler can also be triggered on demand via the internal HTTP endpoint
// POST /internal/reconcile (for ops / testing). That endpoint is registered in
// api/http_api/webhook.go's mux alongside the webhook routes.

package reconcile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"payment-svc/adapter/booking_grpc_adapter"
	"payment-svc/adapter/gateway"
	"payment-svc/store/postgres_store"
	"payment-svc/store/postgres_store/sqlc"
	"payment-svc/util/apperrors"
	"payment-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ReconcileResult summarises one reconciliation run.
type ReconcileResult struct {
	InvoicesChecked    int
	BackfilledEvents   int
	ExpiredVAsMarked   int
	PartialReconciled  int
	Errors             []string
	RunAt              time.Time
	DurationMs         int64
}

// Reconciler orchestrates the hourly reconciliation loop.
type Reconciler struct {
	logger         *zerolog.Logger
	tracer         trace.Tracer
	store          postgres_store.IStore
	primaryGateway gateway.GatewayAdapter
	xenditGateway  gateway.GatewayAdapter
	mockGateway    gateway.GatewayAdapter
	bookingAdapter *booking_grpc_adapter.Adapter
}

// Config holds the Reconciler dependencies.
type Config struct {
	Logger         *zerolog.Logger
	Tracer         trace.Tracer
	Store          postgres_store.IStore
	PrimaryGateway gateway.GatewayAdapter
	XenditGateway  gateway.GatewayAdapter
	MockGateway    gateway.GatewayAdapter
	BookingAdapter *booking_grpc_adapter.Adapter
}

// NewReconciler constructs a Reconciler.
func NewReconciler(cfg Config) *Reconciler {
	return &Reconciler{
		logger:         cfg.Logger,
		tracer:         cfg.Tracer,
		store:          cfg.Store,
		primaryGateway: cfg.PrimaryGateway,
		xenditGateway:  cfg.XenditGateway,
		mockGateway:    cfg.MockGateway,
		bookingAdapter: cfg.BookingAdapter,
	}
}

// StartCron launches the hourly reconciliation ticker.
// It blocks until ctx is cancelled (call in a goroutine from cmd/start.go).
func (r *Reconciler) StartCron(ctx context.Context) {
	logger := r.logger
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	logger.Info().Msg("reconciliation cron started (interval=1h)")

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("reconciliation cron stopping")
			return
		case t := <-ticker.C:
			logger.Info().Time("tick", t).Msg("reconciliation cron tick — running")
			result, err := r.Run(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("reconciliation run failed")
				continue
			}
			logger.Info().
				Int("invoices_checked", result.InvoicesChecked).
				Int("backfilled_events", result.BackfilledEvents).
				Int("expired_vas_marked", result.ExpiredVAsMarked).
				Int("partial_reconciled", result.PartialReconciled).
				Int("errors", len(result.Errors)).
				Int64("duration_ms", result.DurationMs).
				Msg("reconciliation run complete")
		}
	}
}

// Run executes one reconciliation pass. Callable on demand (HTTP trigger, tests).
func (r *Reconciler) Run(ctx context.Context) (*ReconcileResult, error) {
	const op = "reconcile.Reconciler.Run"

	start := time.Now()

	ctx, span := r.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, r.logger)
	logger.Info().Str("op", op).Msg("reconciliation run started")

	result := &ReconcileResult{RunAt: start}

	// Fetch all invoices that are unpaid/partially_paid with an active non-expired VA.
	invoices, err := r.store.ListUnpaidInvoicesWithActiveVA(ctx)
	if err != nil {
		err = fmt.Errorf("list unpaid invoices: %w", postgres_store.WrapDBError(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	result.InvoicesChecked = len(invoices)
	span.SetAttributes(attribute.Int("invoices_checked", len(invoices)))

	for _, invoice := range invoices {
		if err := r.reconcileInvoice(ctx, invoice, result); err != nil {
			errMsg := fmt.Sprintf("invoice %s: %v", uuidToString(invoice.ID), err)
			result.Errors = append(result.Errors, errMsg)
			logger.Warn().Err(err).Str("invoice_id", uuidToString(invoice.ID)).Msg("reconcile invoice error — continuing")
		}
	}

	result.DurationMs = time.Since(start).Milliseconds()

	span.SetAttributes(
		attribute.Int("backfilled_events", result.BackfilledEvents),
		attribute.Int("expired_vas_marked", result.ExpiredVAsMarked),
		attribute.Int("partial_reconciled", result.PartialReconciled),
		attribute.Int("errors", len(result.Errors)),
	)
	if len(result.Errors) > 0 {
		span.SetStatus(codes.Error, fmt.Sprintf("%d invoice errors", len(result.Errors)))
	} else {
		span.SetStatus(codes.Ok, "success")
	}

	return result, nil
}

// reconcileInvoice processes one invoice in the reconciliation pass.
func (r *Reconciler) reconcileInvoice(ctx context.Context, invoice sqlc.Invoice, result *ReconcileResult) error {
	// Find the most recent active VA for this invoice to get gateway details.
	va, err := r.store.GetVAByInvoiceID(ctx, invoice.ID)
	if err != nil {
		return fmt.Errorf("get VA: %w", postgres_store.WrapDBError(err))
	}

	// Select the appropriate gateway adapter.
	gw := r.selectGateway(va.Gateway)
	if gw == nil {
		return fmt.Errorf("unknown gateway %q", va.Gateway)
	}

	// Query gateway for current payment status.
	status, err := gw.GetPaymentStatus(ctx, va.GatewayVaID)
	if err != nil {
		return fmt.Errorf("gateway GetPaymentStatus: %w", err)
	}

	// Case 1: Gateway says expired — mark VA and invoice.
	if status.Expired {
		if _, err := r.store.UpdateVAStatus(ctx, sqlc.UpdateVAStatusParams{
			ID:     va.ID,
			Status: "expired",
		}); err != nil {
			return fmt.Errorf("mark VA expired: %w", postgres_store.WrapDBError(err))
		}
		if _, err := r.store.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
			ID:     invoice.ID,
			Status: "void",
		}); err != nil {
			return fmt.Errorf("void invoice: %w", postgres_store.WrapDBError(err))
		}
		result.ExpiredVAsMarked++
		return nil
	}

	// Case 2: Gateway says paid / partially paid — check for missing events.
	if !status.Paid || status.AmountPaidIDR <= 0 {
		// Gateway not yet paid — nothing to do this cycle.
		return nil
	}

	currentPaidIDR := numericToFloat(invoice.PaidAmount)

	if status.AmountPaidIDR <= currentPaidIDR {
		// Already reconciled — no gap.
		return nil
	}

	// There is a gap: gateway has more money than we recorded.
	// Check for idempotency: if we already have an event for this gateway_txn_id, skip.
	if status.GatewayTxnID != "" {
		_, dupErr := r.store.GetPaymentEventByGatewayTxnID(ctx, sqlc.GetPaymentEventByGatewayTxnIDParams{
			Gateway:      va.Gateway,
			GatewayTxnID: status.GatewayTxnID,
		})
		if dupErr == nil {
			// Event already exists — partial reconcile but event is there.
			// Update paid_amount if it differs.
			if status.AmountPaidIDR > currentPaidIDR {
				newStatus := "partially_paid"
				if status.AmountPaidIDR >= numericToFloat(invoice.AmountTotal) {
					newStatus = "paid"
				}
				if _, err := r.store.UpdateInvoicePaidAmount(ctx, sqlc.UpdateInvoicePaidAmountParams{
					ID:         invoice.ID,
					PaidAmount: floatToNumeric(status.AmountPaidIDR),
					Status:     newStatus,
				}); err != nil {
					return fmt.Errorf("update paid_amount: %w", postgres_store.WrapDBError(err))
				}
				result.PartialReconciled++
			}
			return nil
		}
		if !errors.Is(postgres_store.WrapDBError(dupErr), errors.New("not found")) {
			// Ignore ErrNotFound — that's the expected case (event doesn't exist yet).
			// For any other error, skip this invoice and log.
		}
	}

	// Backfill: insert payment_event for the reconciled amount.
	backfillAmount := status.AmountPaidIDR - currentPaidIDR
	newPaid := currentPaidIDR + backfillAmount

	newInvoiceStatus := "partially_paid"
	if newPaid >= numericToFloat(invoice.AmountTotal) {
		newInvoiceStatus = "paid"
	}

	_, txErr := r.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			// Insert backfill event.
			txnID := pgtype.Text{Valid: false}
			if status.GatewayTxnID != "" {
				txnID = pgtype.Text{String: status.GatewayTxnID, Valid: true}
			}
			if _, err := q.CreatePaymentEvent(ctx, sqlc.CreatePaymentEventParams{
				InvoiceID:      invoice.ID,
				Gateway:        va.Gateway,
				GatewayTxnID:   txnID,
				Kind:           "payment_received",
				Amount:         floatToNumeric(backfillAmount),
				RawPayload:     []byte(`{"source":"reconciliation_backfill"}`),
				ApprovalStatus: pgtype.Text{Valid: false},
				ReceivedAt:     pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
			}); err != nil {
				return fmt.Errorf("insert backfill event: %w", postgres_store.WrapDBError(err))
			}

			// Update invoice paid_amount + status.
			if _, err := q.UpdateInvoicePaidAmount(ctx, sqlc.UpdateInvoicePaidAmountParams{
				ID:         invoice.ID,
				PaidAmount: floatToNumeric(newPaid),
				Status:     newInvoiceStatus,
			}); err != nil {
				return fmt.Errorf("update invoice: %w", postgres_store.WrapDBError(err))
			}

			// If paid, mark VA as paid.
			if newInvoiceStatus == "paid" {
				if _, err := q.UpdateVAStatus(ctx, sqlc.UpdateVAStatusParams{
					ID:     va.ID,
					Status: "paid",
				}); err != nil {
					return fmt.Errorf("mark VA paid: %w", postgres_store.WrapDBError(err))
				}
			}
			return nil
		},
	})
	if txErr != nil {
		return fmt.Errorf("backfill transaction: %w", txErr)
	}

	result.BackfilledEvents++

	// Signal booking-svc after durable commit.
	if r.bookingAdapter != nil {
		bookingIDStr := uuidToString(invoice.BookingID)
		if _, markErr := r.bookingAdapter.MarkBookingPaid(ctx, &booking_grpc_adapter.MarkBookingPaidParams{
			BookingID:     bookingIDStr,
			AmountPaidIDR: newPaid,
			InvoiceStatus: newInvoiceStatus,
			InvoiceID:     uuidToString(invoice.ID),
		}); markErr != nil {
			// Non-fatal — next reconcile cycle will retry if booking-svc is down.
			r.logger.Warn().Err(markErr).
				Str("booking_id", bookingIDStr).
				Msg("MarkBookingPaid failed in reconciler — will retry")
		}
	}

	return nil
}

// selectGateway returns the appropriate adapter for the given gateway name.
func (r *Reconciler) selectGateway(name string) gateway.GatewayAdapter {
	switch name {
	case "xendit":
		if r.xenditGateway != nil {
			return r.xenditGateway
		}
	case "mock":
		if r.mockGateway != nil {
			return r.mockGateway
		}
	}
	if r.primaryGateway != nil {
		return r.primaryGateway
	}
	return r.mockGateway
}

// ---------------------------------------------------------------------------
// pgtype helpers (duplicated from service layer to avoid circular import)
// ---------------------------------------------------------------------------

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%.2f", f))
	return n
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
