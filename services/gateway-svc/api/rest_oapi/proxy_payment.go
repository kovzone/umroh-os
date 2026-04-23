// proxy_payment.go — gateway REST handlers for payment-svc RPCs (BL-PAY-020).
//
// Route topology (bearer-protected):
//   POST /v1/payments/link — ReissuePaymentLink (CS closing tool)
//
// Per ADR-0009: gateway is the single REST entry-point; payment-svc is pure gRPC.

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/payment_grpc_adapter"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

// ReissuePaymentLinkBody is the JSON body for POST /v1/payments/link.
type ReissuePaymentLinkBody struct {
	BookingID   string `json:"booking_id"`
	BankCode    string `json:"bank_code"`
	GatewayPref string `json:"gateway_pref"`
}

// ReissuePaymentLinkResponseData is the JSON response for ReissuePaymentLink.
type ReissuePaymentLinkResponseData struct {
	InvoiceID     string  `json:"invoice_id"`
	BookingID     string  `json:"booking_id"`
	AccountNumber string  `json:"account_number"`
	BankCode      string  `json:"bank_code"`
	AmountTotal   float64 `json:"amount_total"`
	ExpiresAt     string  `json:"expires_at"`
	Gateway       string  `json:"gateway"`
	IsNew         bool    `json:"is_new"`
}

// ---------------------------------------------------------------------------
// ReissuePaymentLink — POST /v1/payments/link (bearer)
// ---------------------------------------------------------------------------

// ReissuePaymentLink retrieves or re-issues the VA link for an existing booking.
// CS tool: used when a pilgrim's VA has expired or was never claimed.
func (s *Server) ReissuePaymentLink(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ReissuePaymentLink"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/payments/link"))
	logger.Info().Str("op", op).Msg("")

	var body ReissuePaymentLinkBody
	if err := c.BodyParser(&body); err != nil {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	if body.BookingID == "" {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, errors.New("booking_id diperlukan")))
	}

	// Extract actor_user_id from the bearer token identity injected by middleware.
	var actorUserID string
	if id, ok := c.Locals(middleware.IdentityKey).(*middleware.Identity); ok && id != nil {
		actorUserID = id.UserID
	}

	result, err := s.svc.ReissuePaymentLink(ctx, &payment_grpc_adapter.ReissuePaymentLinkParams{
		BookingID:   body.BookingID,
		BankCode:    body.BankCode,
		GatewayPref: body.GatewayPref,
		ActorUserID: actorUserID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writePaymentError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": ReissuePaymentLinkResponseData{
			InvoiceID:     result.InvoiceID,
			BookingID:     result.BookingID,
			AccountNumber: result.AccountNumber,
			BankCode:      result.BankCode,
			AmountTotal:   result.AmountTotal,
			ExpiresAt:     result.ExpiresAt,
			Gateway:       result.Gateway,
			IsNew:         result.IsNew,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writePaymentError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
		message = "invoice tidak ditemukan untuk booking tersebut"
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		message = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		message = "akses ditolak"
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
