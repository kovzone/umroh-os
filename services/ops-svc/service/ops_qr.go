// ops_qr.go — HMAC-signed QR token generation and verification (BL-OPS-003).
//
// Token format:
//   base64url(json({jamaah_id, departure_id, card_type, issued_at_unix}))
//   + "."
//   + hex(HMAC-SHA256(payload, secret))
//
// Secret is read from env QR_SECRET (default: "default-qr-secret-change-in-prod").
//
// GenerateIDCard: signs token and stores in ops.id_card_issuances (upsert).
// VerifyIDCard:   decodes token, verifies HMAC, checks DB record exists.

package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"ops-svc/store/postgres_store/sqlc"
	"ops-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// qrPayload is the JSON structure embedded in the token.
type qrPayload struct {
	JamaahID    string `json:"jamaah_id"`
	DepartureID string `json:"departure_id"`
	CardType    string `json:"card_type"`
	IssuedAtUnix int64 `json:"issued_at_unix"`
}

// qrSecret returns the HMAC secret from environment.
func qrSecret() []byte {
	s := os.Getenv("QR_SECRET")
	if s == "" {
		s = "default-qr-secret-change-in-prod"
	}
	return []byte(s)
}

// buildToken constructs token = base64url(payload_json) + "." + hex(hmac).
func buildToken(payload qrPayload) (token string, qrData string, err error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("marshal payload: %w", err)
	}
	qrData = base64.RawURLEncoding.EncodeToString(raw)

	mac := hmac.New(sha256.New, qrSecret())
	mac.Write([]byte(qrData))
	sig := hex.EncodeToString(mac.Sum(nil))

	token = qrData + "." + sig
	return token, qrData, nil
}

// verifyToken decodes and verifies token; returns payload on success.
func verifyToken(token string) (qrPayload, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return qrPayload{}, fmt.Errorf("invalid token format")
	}
	qrData, sig := parts[0], parts[1]

	// Verify HMAC.
	mac := hmac.New(sha256.New, qrSecret())
	mac.Write([]byte(qrData))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return qrPayload{}, fmt.Errorf("invalid HMAC signature")
	}

	// Decode payload.
	raw, err := base64.RawURLEncoding.DecodeString(qrData)
	if err != nil {
		return qrPayload{}, fmt.Errorf("decode payload: %w", err)
	}
	var p qrPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return qrPayload{}, fmt.Errorf("unmarshal payload: %w", err)
	}
	return p, nil
}

// ---------------------------------------------------------------------------
// GenerateIDCard
// ---------------------------------------------------------------------------

// GenerateIDCardParams holds inputs for GenerateIDCard.
type GenerateIDCardParams struct {
	JamaahID      string
	DepartureID   string
	CardType      string // "id_card" | "luggage_tag"
	JamaahName    string
	DepartureName string
}

// GenerateIDCardResult holds the result of GenerateIDCard.
type GenerateIDCardResult struct {
	Token    string
	QrData   string
	IssuedAt time.Time
}

// GenerateIDCard generates (or refreshes) an HMAC-signed ID card / luggage tag
// token and stores it in ops.id_card_issuances.
func (svc *Service) GenerateIDCard(ctx context.Context, params *GenerateIDCardParams) (*GenerateIDCardResult, error) {
	const op = "service.Service.GenerateIDCard"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("jamaah_id", params.JamaahID),
		attribute.String("departure_id", params.DepartureID),
		attribute.String("card_type", params.CardType),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.JamaahID == "" {
		return nil, fmt.Errorf("%s: jamaah_id is required", op)
	}
	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}
	if params.CardType != "id_card" && params.CardType != "luggage_tag" {
		return nil, fmt.Errorf("%s: card_type must be 'id_card' or 'luggage_tag', got %q", op, params.CardType)
	}

	issuedAtUnix := time.Now().Unix()
	payload := qrPayload{
		JamaahID:     params.JamaahID,
		DepartureID:  params.DepartureID,
		CardType:     params.CardType,
		IssuedAtUnix: issuedAtUnix,
	}

	token, qrData, err := buildToken(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: build token: %w", op, err)
	}

	issuance, err := svc.store.UpsertIDCardIssuance(ctx, sqlc.UpsertIDCardIssuanceParams{
		JamaahID:    params.JamaahID,
		DepartureID: params.DepartureID,
		CardType:    params.CardType,
		Token:       token,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: upsert id card issuance: %w", op, err)
	}

	issuedAt := issuance.IssuedAt.Time
	if !issuance.IssuedAt.Valid {
		issuedAt = time.Unix(issuedAtUnix, 0)
	}

	logger.Info().
		Str("op", op).
		Str("jamaah_id", params.JamaahID).
		Str("departure_id", params.DepartureID).
		Str("card_type", params.CardType).
		Str("issuance_id", issuance.ID).
		Msg("ID card token generated")

	span.SetStatus(otelCodes.Ok, "success")
	return &GenerateIDCardResult{
		Token:    token,
		QrData:   qrData,
		IssuedAt: issuedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// VerifyIDCard
// ---------------------------------------------------------------------------

// VerifyIDCardParams holds inputs for VerifyIDCard.
type VerifyIDCardParams struct {
	Token string
}

// VerifyIDCardResult holds the result of VerifyIDCard.
type VerifyIDCardResult struct {
	Valid       bool
	JamaahID    string
	DepartureID string
	CardType    string
	ErrorReason string
}

// VerifyIDCard verifies an HMAC-signed ID card token.
// Returns valid=false with ErrorReason on failure; does not return error for
// tamper/not-found cases (only for infrastructure failures).
func (svc *Service) VerifyIDCard(ctx context.Context, params *VerifyIDCardParams) (*VerifyIDCardResult, error) {
	const op = "service.Service.VerifyIDCard"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
	)

	if params.Token == "" {
		return &VerifyIDCardResult{Valid: false, ErrorReason: "token is required"}, nil
	}

	// Verify HMAC first.
	payload, err := verifyToken(params.Token)
	if err != nil {
		span.SetStatus(otelCodes.Ok, "invalid token")
		return &VerifyIDCardResult{
			Valid:       false,
			ErrorReason: "invalid or tampered token: " + err.Error(),
		}, nil
	}

	// Check DB record exists (audit trail).
	_, err = svc.store.GetIDCardIssuanceByToken(ctx, params.Token)
	if err != nil {
		span.SetStatus(otelCodes.Ok, "not found")
		return &VerifyIDCardResult{
			Valid:       false,
			ErrorReason: "token not found in issuance records",
		}, nil
	}

	span.SetStatus(otelCodes.Ok, "valid")
	return &VerifyIDCardResult{
		Valid:       true,
		JamaahID:    payload.JamaahID,
		DepartureID: payload.DepartureID,
		CardType:    payload.CardType,
	}, nil
}
