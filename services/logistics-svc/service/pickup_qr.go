// pickup_qr.go — GeneratePickupQR and RedeemPickupQR service-layer implementations.
//
// BL-LOG-003 / S3-E-02.
//
// GeneratePickupQR:
//   1. Looks up fulfillment task by booking_id.
//   2. Idempotent: if a non-expired, non-used token already exists, return it.
//   3. Generates a UUID v4 token with 7-day TTL.
//   4. Inserts pickup_token row.
//
// RedeemPickupQR:
//   1. Finds token by token string.
//   2. Rejects if expired.
//   3. Rejects if already used.
//   4. Marks token as used + used_at = now().
//   5. Returns booking_id via the task.

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"logistics-svc/store/postgres_store/sqlc"
	"logistics-svc/util/logging"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

const pickupTokenTTL = 7 * 24 * time.Hour

// ---------------------------------------------------------------------------
// GeneratePickupQR
// ---------------------------------------------------------------------------

// GeneratePickupQRParams holds inputs for generating a pickup QR token.
type GeneratePickupQRParams struct {
	BookingID string
}

// GeneratePickupQRResult holds the result of a GeneratePickupQR call.
type GeneratePickupQRResult struct {
	PickupTokenID string
	Token         string
	ExpiresAt     time.Time
	Replayed      bool // true if an existing active token was returned
}

// GeneratePickupQR creates a single-use pickup token with a 7-day TTL.
// Idempotent: if a valid (non-expired, non-used) token already exists for
// this booking's task, it is returned without creating a duplicate.
func (svc *Service) GeneratePickupQR(ctx context.Context, params *GeneratePickupQRParams) (*GeneratePickupQRResult, error) {
	const op = "service.Service.GeneratePickupQR"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.BookingID == "" {
		return nil, fmt.Errorf("%s: booking_id is required", op)
	}

	// --- Lookup task ---
	task, err := svc.store.GetFulfillmentTaskByBookingID(ctx, params.BookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: no fulfillment task found for booking_id %s", op, params.BookingID)
		}
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get task: %w", op, err)
	}

	// --- Idempotency: check for existing active token ---
	existing, err := svc.store.GetActivePickupTokenByTaskID(ctx, task.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get existing token: %w", op, err)
	}
	if err == nil {
		logger.Info().
			Str("op", op).
			Str("booking_id", params.BookingID).
			Str("token_id", existing.ID).
			Bool("replayed", true).
			Msg("active pickup token already exists, returning existing")

		span.SetStatus(otelCodes.Ok, "replayed")
		return &GeneratePickupQRResult{
			PickupTokenID: existing.ID,
			Token:         existing.Token,
			ExpiresAt:     existing.ExpiresAt,
			Replayed:      true,
		}, nil
	}

	// --- Generate new token ---
	token := uuid.New().String()
	expiresAt := time.Now().UTC().Add(pickupTokenTTL)

	pt, err := svc.store.InsertPickupToken(ctx, sqlc.InsertPickupTokenParams{
		TaskID:    task.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: insert pickup token: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("task_id", task.ID).
		Str("token_id", pt.ID).
		Time("expires_at", expiresAt).
		Msg("pickup QR token generated")

	span.SetStatus(otelCodes.Ok, "created")
	return &GeneratePickupQRResult{
		PickupTokenID: pt.ID,
		Token:         pt.Token,
		ExpiresAt:     pt.ExpiresAt,
		Replayed:      false,
	}, nil
}

// ---------------------------------------------------------------------------
// RedeemPickupQR
// ---------------------------------------------------------------------------

// RedeemPickupQRParams holds inputs for redeeming a pickup QR token.
type RedeemPickupQRParams struct {
	Token string
}

// RedeemPickupQRResult holds the result of a RedeemPickupQR call.
type RedeemPickupQRResult struct {
	Redeemed    bool
	BookingID   string
	TaskID      string
	ErrorReason string
}

// RedeemPickupQR validates and marks a pickup token as used.
// Returns Redeemed=false with an ErrorReason if the token is expired or already used.
func (svc *Service) RedeemPickupQR(ctx context.Context, params *RedeemPickupQRParams) (*RedeemPickupQRResult, error) {
	const op = "service.Service.RedeemPickupQR"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("token", params.Token),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.Token == "" {
		return nil, fmt.Errorf("%s: token is required", op)
	}

	// --- Find token ---
	pt, err := svc.store.GetPickupTokenByToken(ctx, params.Token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &RedeemPickupQRResult{
				Redeemed:    false,
				ErrorReason: "token not found",
			}, nil
		}
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get token: %w", op, err)
	}

	// --- Check: already used ---
	if pt.Used {
		logger.Warn().
			Str("op", op).
			Str("token_id", pt.ID).
			Msg("pickup token already used")
		span.SetStatus(otelCodes.Ok, "already_used")
		return &RedeemPickupQRResult{
			Redeemed:    false,
			TaskID:      pt.TaskID,
			ErrorReason: "token already used",
		}, nil
	}

	// --- Check: expired ---
	if time.Now().UTC().After(pt.ExpiresAt) {
		logger.Warn().
			Str("op", op).
			Str("token_id", pt.ID).
			Time("expires_at", pt.ExpiresAt).
			Msg("pickup token expired")
		span.SetStatus(otelCodes.Ok, "expired")
		return &RedeemPickupQRResult{
			Redeemed:    false,
			TaskID:      pt.TaskID,
			ErrorReason: "token expired",
		}, nil
	}

	// --- Mark as used ---
	_, err = svc.store.MarkPickupTokenUsed(ctx, pt.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: mark token used: %w", op, err)
	}

	// --- Fetch task to return booking_id ---
	task, err := svc.store.GetFulfillmentTaskByID(ctx, pt.TaskID)
	if err != nil {
		// Non-fatal: token was marked used; just return without booking_id.
		logger.Warn().
			Str("op", op).
			Str("task_id", pt.TaskID).
			Err(err).
			Msg("could not fetch task after redeeming token")
	}

	logger.Info().
		Str("op", op).
		Str("token_id", pt.ID).
		Str("task_id", pt.TaskID).
		Str("booking_id", task.BookingID).
		Msg("pickup QR redeemed")

	span.SetStatus(otelCodes.Ok, "redeemed")
	return &RedeemPickupQRResult{
		Redeemed:  true,
		BookingID: task.BookingID,
		TaskID:    pt.TaskID,
	}, nil
}
