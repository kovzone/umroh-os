// proxy_webhook.go — gateway REST handlers for payment webhook routes (ISSUE-007/008).
//
// Route topology (public — no bearer, signature-protected at payment-svc level):
//   POST /v1/webhooks/midtrans        — WebhookMidtrans
//   POST /v1/webhooks/xendit          — WebhookXendit
//   POST /v1/webhooks/mock/trigger    — WebhookMockTrigger (dev only, MOCK_GATEWAY=true)
//
// These handlers read the raw request body + gateway signature header, then
// forward to payment-svc.ProcessWebhook via gRPC. payment-svc performs all
// signature verification and business logic.
//
// Per ADR-0009: gateway is the single REST entry-point; payment-svc is pure gRPC.

package rest_oapi

import (
	"encoding/json"
	"errors"

	"gateway-svc/adapter/payment_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// WebhookMidtrans — POST /v1/webhooks/midtrans (public)
// ---------------------------------------------------------------------------

// WebhookMidtrans receives a Midtrans payment notification.
// Reads raw body + X-Callback-Token header, forwards to payment-svc via gRPC.
func (s *Server) WebhookMidtrans(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.WebhookMidtrans"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/webhooks/midtrans"))
	logger.Info().Str("op", op).Msg("")

	signature := c.Get("X-Callback-Token", "")
	rawBody := c.Body()

	result, err := s.svc.ProcessWebhook(ctx, &payment_grpc_adapter.ProcessWebhookParams{
		Gateway:   "midtrans",
		Payload:   rawBody,
		Signature: signature,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeWebhookError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"invoice_id":    result.InvoiceID,
		"invoice_status": result.NewStatus,
		"replayed":      result.Replayed,
	})
}

// ---------------------------------------------------------------------------
// WebhookXendit — POST /v1/webhooks/xendit (public)
// ---------------------------------------------------------------------------

// WebhookXendit receives a Xendit payment callback.
// Reads raw body + X-Callback-Token header, forwards to payment-svc via gRPC.
func (s *Server) WebhookXendit(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.WebhookXendit"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/webhooks/xendit"))
	logger.Info().Str("op", op).Msg("")

	signature := c.Get("X-Callback-Token", "")
	rawBody := c.Body()

	result, err := s.svc.ProcessWebhook(ctx, &payment_grpc_adapter.ProcessWebhookParams{
		Gateway:   "xendit",
		Payload:   rawBody,
		Signature: signature,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeWebhookError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"invoice_id":    result.InvoiceID,
		"invoice_status": result.NewStatus,
		"replayed":      result.Replayed,
	})
}

// ---------------------------------------------------------------------------
// WebhookMockTrigger — POST /v1/webhooks/mock/trigger (public, dev only)
// ---------------------------------------------------------------------------

// MockTriggerBody is the JSON body for POST /v1/webhooks/mock/trigger.
type MockTriggerBody struct {
	InvoiceID string  `json:"invoice_id"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
}

// WebhookMockTrigger is a dev-only endpoint that synthesises a mock gateway callback.
// Only registered when MOCK_GATEWAY=true. Constructs a JSON payload from the body
// and calls payment-svc.ProcessWebhook with gateway="mock".
func (s *Server) WebhookMockTrigger(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.WebhookMockTrigger"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/webhooks/mock/trigger"))
	logger.Info().Str("op", op).Msg("")

	var body MockTriggerBody
	if err := c.BodyParser(&body); err != nil {
		return writeWebhookError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	if body.InvoiceID == "" {
		return writeWebhookError(c, span, errors.Join(apperrors.ErrValidation, errors.New("invoice_id diperlukan")))
	}

	// Construct a synthetic mock payload that payment-svc can parse.
	payload, _ := json.Marshal(map[string]interface{}{
		"invoice_id": body.InvoiceID,
		"status":     body.Status,
		"amount":     body.Amount,
	})

	result, err := s.svc.ProcessWebhook(ctx, &payment_grpc_adapter.ProcessWebhookParams{
		Gateway:   "mock",
		Payload:   payload,
		Signature: "", // mock gateway skips signature verification
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeWebhookError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"invoice_id":    result.InvoiceID,
		"invoice_status": result.NewStatus,
		"replayed":      result.Replayed,
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeWebhookError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrUnauthorized):
		// signature mismatch → 401
		httpStatus = fiber.StatusUnauthorized
		code = "invalid_signature"
		message = "invalid or missing webhook signature"
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
		message = err.Error()
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		httpStatus = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan pembayaran sementara tidak tersedia"
	default:
		httpStatus = fiber.StatusInternalServerError
		code = "internal_error"
		message = "terjadi kesalahan tidak terduga"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
