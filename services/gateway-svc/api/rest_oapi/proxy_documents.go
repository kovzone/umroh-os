// proxy_documents.go — gateway REST handlers for jamaah-svc OCR RPCs (S3 Wave 2).
//
// Route topology (all bearer-protected):
//   POST /v1/documents/:id/ocr — TriggerOCR
//   GET  /v1/documents/:id/ocr — GetOCRStatus
//
// Per ADR-0009: gateway is the single REST entry-point; jamaah-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Response body types
// ---------------------------------------------------------------------------

// OCRResponseData is the JSON representation of an OCR result.
type OCRResponseData struct {
	DocumentID string            `json:"document_id"`
	Status     string            `json:"status"`
	Confidence float64           `json:"confidence"`
	OcrResult  map[string]string `json:"ocr_result"`
}

// ---------------------------------------------------------------------------
// TriggerDocumentOCR — POST /v1/documents/:id/ocr (bearer)
// ---------------------------------------------------------------------------

func (s *Server) TriggerDocumentOCR(c *fiber.Ctx, documentID string) error {
	const op = "rest_oapi.Server.TriggerDocumentOCR"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "POST /v1/documents/:id/ocr"),
		attribute.String("document_id", documentID),
	)
	logger.Info().Str("op", op).Str("document_id", documentID).Msg("")

	result, err := s.svc.TriggerOCR(ctx, documentID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeDocumentError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": OCRResponseData{
			DocumentID: result.DocumentID,
			Status:     result.Status,
			Confidence: result.Confidence,
			OcrResult:  result.OcrResult,
		},
	})
}

// ---------------------------------------------------------------------------
// GetDocumentOCRStatus — GET /v1/documents/:id/ocr (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetDocumentOCRStatus(c *fiber.Ctx, documentID string) error {
	const op = "rest_oapi.Server.GetDocumentOCRStatus"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/documents/:id/ocr"),
		attribute.String("document_id", documentID),
	)
	logger.Info().Str("op", op).Str("document_id", documentID).Msg("")

	result, err := s.svc.GetOCRStatus(ctx, documentID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeDocumentError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": OCRResponseData{
			DocumentID: result.DocumentID,
			Status:     result.Status,
			Confidence: result.Confidence,
			OcrResult:  result.OcrResult,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeDocumentError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = err.Error()
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
		message = "layanan dokumen sementara tidak tersedia"
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
