// proxy_finance_correction.go — gateway REST handler for CorrectJournal RPC
// (BL-FIN-006).
//
// Route topology (bearer-protected):
//   POST /v1/finance/journals/:id/correct → CorrectJournal
//
// Per ADR-0009: gateway is the single REST entry-point; finance-svc is pure gRPC.

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/finance_grpc_adapter"
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

// CorrectJournalBody is the JSON body for POST /v1/finance/journals/:id/correct.
type CorrectJournalBody struct {
	Reason string `json:"reason"`
}

// CorrectJournalResponseData is the JSON response for CorrectJournal.
type CorrectJournalResponseData struct {
	CorrectionEntryID string `json:"correction_entry_id"`
	OriginalEntryID   string `json:"original_entry_id"`
	Idempotent        bool   `json:"idempotent"`
}

// ---------------------------------------------------------------------------
// CorrectJournal — POST /v1/finance/journals/:id/correct (bearer)
// ---------------------------------------------------------------------------

// CorrectJournal posts a reversing counter-entry for an existing journal entry.
// The original entry is preserved; only the reversal is inserted.
// Idempotent: a second call with the same entry ID returns the existing correction.
func (s *Server) CorrectJournal(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CorrectJournal"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	entryID := c.Params("id")
	span.SetAttributes(
		attribute.String("endpoint", "POST /v1/finance/journals/:id/correct"),
		attribute.String("entry_id", entryID),
	)
	logger.Info().Str("op", op).Str("entry_id", entryID).Msg("")

	if entryID == "" {
		return writeFinanceCorrectionError(c, span, errors.Join(apperrors.ErrValidation, errors.New("entry id is required")))
	}

	var body CorrectJournalBody
	if err := c.BodyParser(&body); err != nil {
		// Body parse failure is non-fatal; reason just defaults to empty.
		body.Reason = ""
	}

	// Extract actor_user_id from bearer identity.
	var actorUserID string
	if id, ok := c.Locals(middleware.IdentityKey).(*middleware.Identity); ok && id != nil {
		actorUserID = id.UserID
	}

	result, err := s.svc.CorrectJournal(ctx, &finance_grpc_adapter.CorrectJournalParams{
		EntryID:     entryID,
		Reason:      body.Reason,
		ActorUserID: actorUserID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeFinanceCorrectionError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": CorrectJournalResponseData{
			CorrectionEntryID: result.CorrectionEntryID,
			OriginalEntryID:   result.OriginalEntryID,
			Idempotent:        result.Idempotent,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeFinanceCorrectionError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "jurnal entri tidak ditemukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		message = "jurnal tidak dapat dihapus; gunakan correction untuk membuat counter-entry"
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		message = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		httpStatus = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan keuangan sementara tidak tersedia"
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
