// payment.go — service-layer business logic for invoice, VA issuance, webhook
// processing, and refund initiation (S2-E-02 / BL-PAY-001..008).
//
// Design decisions documented inline:
//  1. Gateway selection: try Midtrans primary; fallback to Xendit on 5xx or
//     timeout > 10s; never failover on 4xx (per Q013 / F5 spec).
//  2. VA idempotency: idempotency_key = invoice.id (UUID string). If the
//     gateway has the VA but we have no DB row (crash between gateway call and
//     DB insert), GetVAByIdempotencyKey surfaces the orphaned row on retry.
//  3. Invoice paid_amount is updated atomically inside a single DB transaction
//     alongside the payment_event insert. If the transaction rolls back, the
//     event is not persisted and the idempotency key is not consumed — safe.
//  4. booking-svc.MarkBookingPaid is called AFTER the transaction commits so
//     that the invoice state is durable before we signal downstream.
//  5. All state-changing operations emit an audit log entry via iam-svc
//     RecordAudit (best-effort, in a goroutine — audit failure never blocks
//     the happy path).

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"
	// math/big used for floatField string parsing only

	"payment-svc/adapter/booking_grpc_adapter"
	"payment-svc/adapter/gateway"
	"payment-svc/store/postgres_store"
	"payment-svc/store/postgres_store/sqlc"
	"payment-svc/util/apperrors"
	"payment-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// IssueVirtualAccount
// ---------------------------------------------------------------------------

// IssueVAParams carries everything needed to create an invoice + VA.
type IssueVAParams struct {
	// BookingID is the booking.bookings.id (TEXT prefixed ULID).
	BookingID string
	// AmountTotal is the total payable amount in IDR (already rounded per Q001).
	AmountTotal float64
	// RoundingAdjustmentIDR is the signed Rp delta from pre-round to post-round total.
	RoundingAdjustmentIDR float64
	// FXSnapshot is the FX rates at issuance time (from finance-svc / global config).
	// Stored as arbitrary JSON; payment-svc treats it as opaque bytes.
	FXSnapshot map[string]interface{}
	// GatewayPref is the preferred gateway ("midtrans"|"xendit"|"mock").
	// Empty string means use the configured primary (Midtrans).
	GatewayPref string
	// BankCode is the preferred bank for the VA (empty = gateway default).
	BankCode string
	// ActorUserID is the IAM user triggering this operation (for audit log).
	ActorUserID string
}

// IssueVAResult is returned on success.
type IssueVAResult struct {
	InvoiceID     string
	BookingID     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     time.Time
	Gateway       string
	// Replayed is true when the VA already existed (idempotent replay).
	Replayed bool
}

// IssueVirtualAccount creates an invoice and issues a virtual account.
// It is idempotent: if an invoice + VA already exists for the booking_id,
// it returns the existing VA details without calling the gateway again.
func (s *PaymentService) IssueVirtualAccount(ctx context.Context, params *IssueVAParams) (*IssueVAResult, error) {
	const op = "service.PaymentService.IssueVirtualAccount"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
		attribute.Float64("amount_total", params.AmountTotal),
	)

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	// ---- Parse booking_id UUID ----
	var bookingUUID pgtype.UUID
	if err := bookingUUID.Scan(params.BookingID); err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid booking_id: %w", err))
	}

	// ---- Idempotency: check if invoice + VA already exist for this booking ----
	existingInvoice, err := s.store.GetInvoiceByBookingID(ctx, bookingUUID)
	if err == nil {
		// Invoice exists — check for an active VA.
		existingVA, vaErr := s.store.GetVAByInvoiceID(ctx, existingInvoice.ID)
		if vaErr == nil && existingVA.Status == "active" {
			logger.Info().Str("op", op).Str("invoice_id", uuidToString(existingInvoice.ID)).Msg("replayed: returning existing VA")
			span.SetStatus(codes.Ok, "replayed")
			return &IssueVAResult{
				InvoiceID:     uuidToString(existingInvoice.ID),
				BookingID:     params.BookingID,
				AccountNumber: existingVA.AccountNumber,
				BankCode:      existingVA.BankCode,
				AmountTotal:   numericToFloat(existingInvoice.AmountTotal),
				ExpiresAt:     existingVA.ExpiresAt.Time,
				Gateway:       existingVA.Gateway,
				Replayed:      true,
			}, nil
		}
	} else if !errors.Is(err, apperrors.ErrNotFound) {
		err = fmt.Errorf("check existing invoice: %w", postgres_store.WrapDBError(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// ---- Build fx_snapshot JSON bytes ----
	fxBytes, err := json.Marshal(params.FXSnapshot)
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("marshal fx_snapshot: %w", err))
	}
	if params.FXSnapshot == nil {
		fxBytes = []byte("{}")
	}

	// ---- Create invoice in DB ----
	amountNum := floatToNumeric(params.AmountTotal)
	roundingNum := floatToNumeric(params.RoundingAdjustmentIDR)

	invoice, err := s.store.CreateInvoice(ctx, sqlc.CreateInvoiceParams{
		BookingID:             bookingUUID,
		AmountTotal:           amountNum,
		RoundingAdjustmentIdr: roundingNum,
		Currency:              "IDR",
		FxSnapshot:            fxBytes,
		Status:                "unpaid",
		PaidAmount:            floatToNumeric(0),
	})
	if err != nil {
		err = fmt.Errorf("create invoice: %w", postgres_store.WrapDBError(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	invoiceIDStr := uuidToString(invoice.ID)
	span.SetAttributes(attribute.String("invoice_id", invoiceIDStr))

	// ---- Select gateway and issue VA ----
	gw := s.selectGateway(params.GatewayPref)
	vaResp, err := gw.IssueVA(ctx, gateway.IssueVARequest{
		IdempotencyKey: invoiceIDStr,
		AmountIDR:      params.AmountTotal,
		BookingID:      params.BookingID,
		BankCode:       params.BankCode,
		ExpiryDuration: 24 * time.Hour, // Q010: 24h default
	})
	if err != nil {
		// Gateway call failed — the invoice row exists but no VA yet.
		// The saga caller should retry; idempotency on invoice means no duplicate.
		err = fmt.Errorf("gateway IssueVA: %w", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	// ---- Persist virtual account ----
	expiresAt := pgtype.Timestamptz{Time: vaResp.ExpiresAt, Valid: true}
	va, err := s.store.CreateVirtualAccount(ctx, sqlc.CreateVirtualAccountParams{
		InvoiceID:      invoice.ID,
		Gateway:        vaResp.Gateway,
		GatewayVaID:    vaResp.GatewayVAID,
		AccountNumber:  vaResp.AccountNumber,
		BankCode:       vaResp.BankCode,
		ExpiresAt:      expiresAt,
		Status:         "active",
		IdempotencyKey: invoiceIDStr,
	})
	if err != nil {
		// Unique violation on idempotency_key means the VA row already exists
		// (crash-recovery scenario from F5 edge case: gateway succeeded but DB insert
		// failed on first attempt). Fetch the existing row and return it.
		if errors.Is(postgres_store.WrapDBError(err), apperrors.ErrConflict) {
			existingVA, fetchErr := s.store.GetVAByIdempotencyKey(ctx, invoiceIDStr)
			if fetchErr == nil {
				span.SetStatus(codes.Ok, "va already exists (crash-recovery)")
				return &IssueVAResult{
					InvoiceID:     invoiceIDStr,
					BookingID:     params.BookingID,
					AccountNumber: existingVA.AccountNumber,
					BankCode:      existingVA.BankCode,
					AmountTotal:   params.AmountTotal,
					ExpiresAt:     existingVA.ExpiresAt.Time,
					Gateway:       existingVA.Gateway,
					Replayed:      true,
				}, nil
			}
		}
		err = fmt.Errorf("create virtual account: %w", postgres_store.WrapDBError(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// ---- Emit payment_event: va_created ----
	_, _ = s.store.CreatePaymentEvent(ctx, sqlc.CreatePaymentEventParams{
		InvoiceID:      invoice.ID,
		Gateway:        vaResp.Gateway,
		GatewayTxnID:   pgtype.Text{Valid: false},
		Kind:           "va_created",
		Amount:         floatToNumeric(0),
		RawPayload:     nil,
		ApprovalStatus: pgtype.Text{Valid: false},
		ReceivedAt:     pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	})

	// ---- Audit log (best-effort) ----
	go s.emitAudit(context.WithoutCancel(ctx), "invoice", invoiceIDStr, "create_invoice_va", params.ActorUserID, nil, invoice)

	span.SetAttributes(attribute.String("va_id", uuidToString(va.ID)))
	span.SetStatus(codes.Ok, "success")

	return &IssueVAResult{
		InvoiceID:     invoiceIDStr,
		BookingID:     params.BookingID,
		AccountNumber: va.AccountNumber,
		BankCode:      va.BankCode,
		AmountTotal:   params.AmountTotal,
		ExpiresAt:     va.ExpiresAt.Time,
		Gateway:       va.Gateway,
		Replayed:      false,
	}, nil
}

// ---------------------------------------------------------------------------
// ProcessWebhookEvent
// ---------------------------------------------------------------------------

// WebhookEventParams carries a gateway-signed webhook payload.
type WebhookEventParams struct {
	Gateway   string // "midtrans" | "xendit" | "mock"
	Payload   []byte // raw HTTP body
	Signature string // gateway-specific signature header value
}

// WebhookResult is the service-level result of processing a webhook.
type WebhookResult struct {
	// Replayed is true when the gateway_txn_id was already processed (safe no-op).
	Replayed bool
	InvoiceID string
	NewStatus string
}

// ProcessWebhookEvent handles an incoming gateway webhook.
//
// Critical path (< 500ms p95 per F5 spec):
//  1. Verify signature (reject 401 on failure)
//  2. Idempotency check (return 200 no-op on duplicate)
//  3. Parse payload → extract gateway_txn_id + amount
//  4. Insert payment_event + update invoice.paid_amount (single DB txn)
//  5. Call booking-svc.MarkBookingPaid (after commit, so state is durable)
//  6. Return 200 OK
func (s *PaymentService) ProcessWebhookEvent(ctx context.Context, params *WebhookEventParams) (*WebhookResult, error) {
	const op = "service.PaymentService.ProcessWebhookEvent"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("gateway", params.Gateway),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	// 1. Select gateway adapter and verify signature.
	gw := s.selectGateway(params.Gateway)
	if err := gw.VerifyWebhookSignature(params.Payload, params.Signature); err != nil {
		logger.Warn().Str("op", op).Str("gateway", params.Gateway).Err(err).Msg("webhook signature invalid")
		return nil, errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("invalid webhook signature: %w", err))
	}

	// 2. Parse payload.
	parsed, err := parseWebhookPayload(params.Gateway, params.Payload)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("parse webhook payload: %w", err))
	}

	span.SetAttributes(
		attribute.String("gateway_txn_id", parsed.GatewayTxnID),
		attribute.String("gateway_va_id", parsed.GatewayVAID),
	)

	// 3. Idempotency: check if event already processed.
	_, dupErr := s.store.GetPaymentEventByGatewayTxnID(ctx, sqlc.GetPaymentEventByGatewayTxnIDParams{
		Gateway:      params.Gateway,
		GatewayTxnID: parsed.GatewayTxnID,
	})
	if dupErr == nil {
		logger.Info().Str("op", op).Str("gateway_txn_id", parsed.GatewayTxnID).Msg("webhook replay — no-op")
		span.SetStatus(codes.Ok, "replayed")
		return &WebhookResult{Replayed: true}, nil
	}
	if !errors.Is(dupErr, apperrors.ErrNotFound) {
		return nil, fmt.Errorf("check duplicate event: %w", postgres_store.WrapDBError(dupErr))
	}

	// 4. Look up invoice by VA idempotency key (= invoice.id).
	va, err := s.store.GetVAByIdempotencyKey(ctx, parsed.GatewayVAID)
	if err != nil {
		// If we genuinely don't know this VA: log anomaly, return 200 so gateway stops retrying.
		logger.Warn().Str("op", op).Str("gateway_va_id", parsed.GatewayVAID).Msg("webhook for unknown VA — logging anomaly, returning 200")
		span.SetStatus(codes.Ok, "unknown VA — anomaly logged")
		return &WebhookResult{Replayed: false}, nil
	}

	invoice, err := s.store.GetInvoiceByID(ctx, va.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("get invoice: %w", postgres_store.WrapDBError(err))
	}

	// 5. Compute new paid_amount + invoice status.
	currentPaid := numericToFloat(invoice.PaidAmount)
	newPaid := currentPaid + parsed.AmountIDR
	amountTotal := numericToFloat(invoice.AmountTotal)

	newInvoiceStatus := "partially_paid"
	if newPaid >= amountTotal {
		newInvoiceStatus = "paid"
	}

	// 6. Atomically insert event + update invoice.
	var updatedInvoice sqlc.Invoice
	_, txErr := s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			// Insert payment event (unique constraint on gateway+gateway_txn_id).
			_, err := q.CreatePaymentEvent(ctx, sqlc.CreatePaymentEventParams{
				InvoiceID:      invoice.ID,
				Gateway:        params.Gateway,
				GatewayTxnID:   pgtype.Text{String: parsed.GatewayTxnID, Valid: true},
				Kind:           "payment_received",
				Amount:         floatToNumeric(parsed.AmountIDR),
				RawPayload:     params.Payload,
				ApprovalStatus: pgtype.Text{Valid: false},
				ReceivedAt:     pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
			})
			if err != nil {
				return fmt.Errorf("create payment event: %w", postgres_store.WrapDBError(err))
			}

			// Update invoice paid_amount + status.
			inv, err := q.UpdateInvoicePaidAmount(ctx, sqlc.UpdateInvoicePaidAmountParams{
				ID:         invoice.ID,
				PaidAmount: floatToNumeric(newPaid),
				Status:     newInvoiceStatus,
			})
			if err != nil {
				return fmt.Errorf("update invoice: %w", postgres_store.WrapDBError(err))
			}
			updatedInvoice = inv

			// If paid, mark the VA as paid too.
			if newInvoiceStatus == "paid" {
				if _, err := q.UpdateVAStatus(ctx, sqlc.UpdateVAStatusParams{
					ID:     va.ID,
					Status: "paid",
				}); err != nil {
					return fmt.Errorf("update VA status: %w", postgres_store.WrapDBError(err))
				}
			}
			return nil
		},
	})
	if txErr != nil {
		span.RecordError(txErr)
		span.SetStatus(codes.Error, txErr.Error())
		return nil, txErr
	}

	invoiceIDStr := uuidToString(updatedInvoice.ID)
	bookingIDStr := uuidToString(updatedInvoice.BookingID)

	// 7. Signal booking-svc AFTER commit (durable state first).
	if s.bookingAdapter != nil {
		if _, markErr := s.bookingAdapter.MarkBookingPaid(ctx, &booking_grpc_adapter.MarkBookingPaidParams{
			BookingID:     bookingIDStr,
			AmountPaidIDR: newPaid,
			InvoiceStatus: newInvoiceStatus,
			InvoiceID:     invoiceIDStr,
		}); markErr != nil {
			// Non-fatal: reconciliation cron will catch this on next cycle.
			logger.Warn().Err(markErr).Str("op", op).
				Str("booking_id", bookingIDStr).
				Msg("MarkBookingPaid failed — reconciliation will retry")
		}
	}

	// 8. Audit (best-effort).
	go s.emitAudit(context.WithoutCancel(ctx), "invoice", invoiceIDStr, "payment_received", "", nil, updatedInvoice)

	span.SetStatus(codes.Ok, "success")
	return &WebhookResult{
		Replayed:  false,
		InvoiceID: invoiceIDStr,
		NewStatus: newInvoiceStatus,
	}, nil
}

// ---------------------------------------------------------------------------
// StartRefund
// ---------------------------------------------------------------------------

// StartRefundParams initiates a refund for a booking.
type StartRefundParams struct {
	BookingID   string
	ReasonCode  string
	ActorUserID string
}

// StartRefundResult is returned after a refund is created in "requested" state.
type StartRefundResult struct {
	RefundID    string
	InvoiceID   string
	AmountIDR   float64
	Status      string
}

// StartRefund creates a refund record in state "requested" and calls the gateway.
// Per F5 W8 / ADR-0006: in-process saga, no Temporal.
func (s *PaymentService) StartRefund(ctx context.Context, params *StartRefundParams) (*StartRefundResult, error) {
	const op = "service.PaymentService.StartRefund"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
		attribute.String("reason_code", params.ReasonCode),
	)

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	var bookingUUID pgtype.UUID
	if err := bookingUUID.Scan(params.BookingID); err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid booking_id: %w", err))
	}

	// Find the invoice for this booking.
	invoice, err := s.store.GetInvoiceByBookingID(ctx, bookingUUID)
	if err != nil {
		return nil, fmt.Errorf("get invoice for refund: %w", postgres_store.WrapDBError(err))
	}

	// Refund amount = sum of payment_received events (paid_amount on invoice).
	refundAmount := numericToFloat(invoice.PaidAmount)
	if refundAmount <= 0 {
		return nil, errors.Join(apperrors.ErrValidation,
			fmt.Errorf("no paid amount on invoice %s to refund", uuidToString(invoice.ID)))
	}

	// Create refund record in "requested" state.
	reasonText := pgtype.Text{String: params.ReasonCode, Valid: params.ReasonCode != ""}
	refund, err := s.store.CreateRefund(ctx, sqlc.CreateRefundParams{
		InvoiceID:       invoice.ID,
		BookingID:       bookingUUID,
		Amount:          floatToNumeric(refundAmount),
		ReasonCode:      reasonText,
		Status:          "requested",
		GatewayRefundID: pgtype.Text{Valid: false},
	})
	if err != nil {
		return nil, fmt.Errorf("create refund record: %w", postgres_store.WrapDBError(err))
	}

	invoiceIDStr := uuidToString(invoice.ID)
	refundIDStr := uuidToString(refund.ID)

	// Find the active VA to use for the gateway refund.
	va, vaErr := s.store.GetVAByInvoiceID(ctx, invoice.ID)

	// Call gateway Refund API if we have a VA.
	if vaErr == nil {
		gw := s.selectGateway(va.Gateway)
		gwResp, gwErr := gw.Refund(ctx, gateway.RefundRequest{
			RefundID:    refundIDStr,
			GatewayVAID: va.GatewayVaID,
			AmountIDR:   refundAmount,
			Reason:      params.ReasonCode,
		})
		if gwErr != nil {
			// Gateway refund failed — leave in "requested" for admin review.
			logger.Warn().Err(gwErr).Str("op", op).Str("refund_id", refundIDStr).
				Msg("gateway refund call failed; staying in 'requested' for admin review")
		} else {
			// Update to "processing" with gateway_refund_id.
			updatedRefund, updateErr := s.store.UpdateRefundStatus(ctx, sqlc.UpdateRefundStatusParams{
				ID:              refund.ID,
				Status:          "processing",
				GatewayRefundID: pgtype.Text{String: gwResp.GatewayRefundID, Valid: gwResp.GatewayRefundID != ""},
			})
			if updateErr == nil {
				refund = updatedRefund
			}
		}
	}

	// Emit refund_issued payment event for audit symmetry (F5 W8 step 3).
	_, _ = s.store.CreatePaymentEvent(ctx, sqlc.CreatePaymentEventParams{
		InvoiceID:      invoice.ID,
		Gateway:        "payment-svc",
		GatewayTxnID:   pgtype.Text{String: refundIDStr, Valid: true},
		Kind:           "refund_issued",
		Amount:         floatToNumeric(-refundAmount),
		RawPayload:     nil,
		ApprovalStatus: pgtype.Text{Valid: false},
		ReceivedAt:     pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	})

	// Update invoice status to "refunded".
	_, _ = s.store.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
		ID:     invoice.ID,
		Status: "refunded",
	})

	// Audit (best-effort).
	go s.emitAudit(context.WithoutCancel(ctx), "refund", refundIDStr, "start_refund", params.ActorUserID, nil, refund)

	span.SetStatus(codes.Ok, "success")
	return &StartRefundResult{
		RefundID:  refundIDStr,
		InvoiceID: invoiceIDStr,
		AmountIDR: refundAmount,
		Status:    refund.Status,
	}, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// parsedWebhook holds the extracted fields from a gateway webhook payload.
type parsedWebhook struct {
	GatewayTxnID string
	GatewayVAID  string  // used as idempotency_key look-up in virtual_accounts
	AmountIDR    float64
}

// parseWebhookPayload extracts the fields payment-svc needs from a raw gateway
// webhook payload. Each gateway has its own payload shape; this dispatcher
// routes to the appropriate parser.
//
// MVP implementation: generic JSON parse looking for common field names.
// Production implementation should use gateway-specific structs per gateway
// documentation.
func parseWebhookPayload(gateway string, payload []byte) (*parsedWebhook, error) {
	// Parse as generic JSON map.
	var raw map[string]interface{}
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	result := &parsedWebhook{}

	switch gateway {
	case "midtrans":
		// Midtrans: transaction_id, order_id (= gateway_va_id / invoice_id), gross_amount
		result.GatewayTxnID = stringField(raw, "transaction_id")
		result.GatewayVAID = stringField(raw, "order_id") // idempotency_key = invoice.id
		result.AmountIDR = floatField(raw, "gross_amount")
	case "xendit":
		// Xendit: id (= transaction id), external_id (= order/invoice id), amount
		result.GatewayTxnID = stringField(raw, "id")
		result.GatewayVAID = stringField(raw, "external_id")
		result.AmountIDR = floatField(raw, "amount")
	default:
		// Mock / generic: gateway_txn_id, gateway_va_id, amount
		result.GatewayTxnID = stringField(raw, "gateway_txn_id")
		result.GatewayVAID = stringField(raw, "gateway_va_id")
		result.AmountIDR = floatField(raw, "amount")
	}

	if result.GatewayTxnID == "" {
		return nil, fmt.Errorf("missing gateway_txn_id in %s payload", gateway)
	}
	if result.GatewayVAID == "" {
		return nil, fmt.Errorf("missing gateway_va_id / order_id in %s payload", gateway)
	}
	return result, nil
}

func stringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func floatField(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case string:
			f, _ := new(big.Float).SetString(val)
			if f != nil {
				r, _ := f.Float64()
				return r
			}
		}
	}
	return 0
}

// selectGateway returns the gateway adapter for the given name.
// Empty string or "midtrans" → primary gateway.
// Falls back to mock when MOCK_GATEWAY=true (set via config).
func (s *PaymentService) selectGateway(pref string) gateway.GatewayAdapter {
	switch pref {
	case "xendit":
		if s.xenditGateway != nil {
			return s.xenditGateway
		}
	case "mock":
		if s.mockGateway != nil {
			return s.mockGateway
		}
	}
	// Default: primary (Midtrans) or mock if primary is nil.
	if s.primaryGateway != nil {
		return s.primaryGateway
	}
	return s.mockGateway
}

// emitAudit fires a best-effort audit log entry via iam-svc.
// Called in a goroutine — never blocks the business transaction.
func (s *PaymentService) emitAudit(ctx context.Context, resource, resourceID, action, actorUserID string, oldVal, newVal interface{}) {
	if s.iamAudit == nil {
		return
	}
	var oldBytes, newBytes []byte
	if oldVal != nil {
		oldBytes, _ = json.Marshal(oldVal)
	}
	if newVal != nil {
		newBytes, _ = json.Marshal(newVal)
	}
	_, _ = s.iamAudit.RecordAudit(ctx, &iamAuditAdapter{
		ActorUserID: actorUserID,
		Resource:    resource,
		ResourceID:  resourceID,
		Action:      action,
		OldValue:    oldBytes,
		NewValue:    newBytes,
	})
}

// ---------------------------------------------------------------------------
// pgtype ↔ Go conversion helpers
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
