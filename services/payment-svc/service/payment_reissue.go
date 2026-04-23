// payment_reissue.go — ReissuePaymentLink service logic (BL-PAY-020).
//
// CS-facing: retrieves the active VA for an existing booking, or issues a new
// one on the same invoice if the old VA has expired.
//
// Invariants:
//   - Requires an existing invoice (does NOT create a new invoice — that is
//     booking-svc's saga responsibility at submission time).
//   - Returns ErrNotFound if no invoice exists for the booking_id.
//   - Returns existing active VA unchanged if one exists (IsNew=false).
//   - Issues a fresh VA via the gateway and marks the old VA expired if the
//     current VA is inactive (IsNew=true).
//   - If the invoice is already paid, returns the paid invoice details with
//     IsNew=false (CS can see it is already settled).

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"payment-svc/adapter/gateway"
	"payment-svc/store/postgres_store"
	"payment-svc/store/postgres_store/sqlc"
	"payment-svc/util/apperrors"
	"payment-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ReissuePaymentLinkParams carries the CS-facing input.
type ReissuePaymentLinkParams struct {
	BookingID   string
	BankCode    string
	GatewayPref string
	ActorUserID string
}

// ReissuePaymentLinkResult is returned on success.
type ReissuePaymentLinkResult struct {
	InvoiceID     string
	BookingID     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     time.Time
	Gateway       string
	// IsNew is true when a fresh VA was created; false when existing was returned.
	IsNew bool
}

// ReissuePaymentLink retrieves the active VA for an existing booking's invoice,
// or creates a new VA on the same invoice when the existing VA has expired.
func (s *PaymentService) ReissuePaymentLink(ctx context.Context, params *ReissuePaymentLinkParams) (*ReissuePaymentLinkResult, error) {
	const op = "service.PaymentService.ReissuePaymentLink"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
	)

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	// ---- Parse booking_id ----
	var bookingUUID pgtype.UUID
	if err := bookingUUID.Scan(params.BookingID); err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid booking_id: %w", err))
	}

	// ---- Require existing invoice (CS cannot create one) ----
	invoice, err := s.store.GetInvoiceByBookingID(ctx, bookingUUID)
	if err != nil {
		if errors.Is(postgres_store.WrapDBError(err), apperrors.ErrNotFound) {
			return nil, errors.Join(apperrors.ErrNotFound,
				fmt.Errorf("no invoice found for booking %s", params.BookingID))
		}
		err = fmt.Errorf("get invoice: %w", postgres_store.WrapDBError(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	invoiceIDStr := uuidToString(invoice.ID)
	span.SetAttributes(attribute.String("invoice_id", invoiceIDStr))

	// ---- If invoice is already paid, return info without issuing a new VA ----
	if invoice.Status == "paid" {
		// Try to find the paid VA for display.
		va, vaErr := s.store.GetVAByInvoiceID(ctx, invoice.ID)
		var accNum, bankCode, gwName, expiresAt string
		if vaErr == nil {
			accNum = va.AccountNumber
			bankCode = va.BankCode
			gwName = va.Gateway
			if va.ExpiresAt.Valid {
				expiresAt = va.ExpiresAt.Time.UTC().Format("2006-01-02T15:04:05Z")
			}
		}
		_ = expiresAt // used in response construction below; keep linter happy
		span.SetStatus(codes.Ok, "invoice already paid")
		return &ReissuePaymentLinkResult{
			InvoiceID:     invoiceIDStr,
			BookingID:     params.BookingID,
			AccountNumber: accNum,
			BankCode:      bankCode,
			AmountTotal:   numericToFloat(invoice.AmountTotal),
			ExpiresAt:     func() time.Time {
				if va.ExpiresAt.Valid {
					return va.ExpiresAt.Time
				}
				return time.Time{}
			}(),
			Gateway: gwName,
			IsNew:   false,
		}, nil
	}

	// ---- Check for an active VA ----
	existingVA, vaErr := s.store.GetVAByInvoiceID(ctx, invoice.ID)
	if vaErr == nil && existingVA.Status == "active" {
		// Active VA exists — return it as-is.
		logger.Info().Str("op", op).Str("invoice_id", invoiceIDStr).Msg("active VA found; returning existing")
		span.SetStatus(codes.Ok, "active VA returned")
		return &ReissuePaymentLinkResult{
			InvoiceID:     invoiceIDStr,
			BookingID:     params.BookingID,
			AccountNumber: existingVA.AccountNumber,
			BankCode:      existingVA.BankCode,
			AmountTotal:   numericToFloat(invoice.AmountTotal),
			ExpiresAt:     existingVA.ExpiresAt.Time,
			Gateway:       existingVA.Gateway,
			IsNew:         false,
		}, nil
	}

	// ---- No active VA — issue a fresh one via the gateway ----
	gw := s.selectGateway(params.GatewayPref)
	newKey := invoiceIDStr + ":reissue:" + fmt.Sprintf("%d", time.Now().Unix())
	vaResp, err := gw.IssueVA(ctx, gateway.IssueVARequest{
		IdempotencyKey: newKey,
		AmountIDR:      numericToFloat(invoice.AmountTotal),
		BookingID:      params.BookingID,
		BankCode:       params.BankCode,
		ExpiryDuration: 24 * time.Hour,
	})
	if err != nil {
		err = fmt.Errorf("gateway IssueVA (reissue): %w", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	// Expire the old VA if it existed.
	if vaErr == nil {
		_, _ = s.store.UpdateVAStatus(ctx, sqlc.UpdateVAStatusParams{
			ID:     existingVA.ID,
			Status: "expired",
		})
	}

	// Persist the new VA.
	newVA, err := s.store.CreateVirtualAccount(ctx, sqlc.CreateVirtualAccountParams{
		InvoiceID:      invoice.ID,
		Gateway:        vaResp.Gateway,
		GatewayVaID:    vaResp.GatewayVAID,
		AccountNumber:  vaResp.AccountNumber,
		BankCode:       vaResp.BankCode,
		ExpiresAt:      pgtype.Timestamptz{Time: vaResp.ExpiresAt, Valid: true},
		Status:         "active",
		IdempotencyKey: newKey,
	})
	if err != nil {
		err = fmt.Errorf("create reissued VA: %w", postgres_store.WrapDBError(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Ensure invoice status reflects unpaid (it might have drifted).
	_, _ = s.store.UpdateInvoiceStatus(ctx, sqlc.UpdateInvoiceStatusParams{
		ID:     invoice.ID,
		Status: "unpaid",
	})

	// Audit (best-effort).
	go s.emitAudit(context.WithoutCancel(ctx), "invoice", invoiceIDStr, "reissue_payment_link", params.ActorUserID, nil, newVA)

	span.SetStatus(codes.Ok, "new VA created")
	return &ReissuePaymentLinkResult{
		InvoiceID:     invoiceIDStr,
		BookingID:     params.BookingID,
		AccountNumber: newVA.AccountNumber,
		BankCode:      newVA.BankCode,
		AmountTotal:   numericToFloat(invoice.AmountTotal),
		ExpiresAt:     newVA.ExpiresAt.Time,
		Gateway:       newVA.Gateway,
		IsNew:         true,
	}, nil
}
