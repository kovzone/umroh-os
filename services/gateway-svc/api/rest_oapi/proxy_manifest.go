// proxy_manifest.go — gateway REST handler for jamaah-svc departure manifest
// (Wave 1A / Phase 6).
//
// Route topology (bearer-protected):
//   GET /v1/manifest/:departure_id → GetDepartureManifest
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

// ManifestJamaahData is the JSON representation of one pilgrim row.
type ManifestJamaahData struct {
	BookingID     string `json:"booking_id"`
	Name          string `json:"name"`
	NIK           string `json:"nik"`
	Phone         string `json:"phone"`
	RoomType      string `json:"room_type"`
	BookingStatus string `json:"booking_status"`
	DocStatus     string `json:"doc_status"`
}

// ManifestResponseData is the response for GET /v1/manifest/:departure_id.
type ManifestResponseData struct {
	DepartureID string               `json:"departure_id"`
	TotalJamaah int32                `json:"total_jamaah"`
	LunasPaid   int32                `json:"lunas_paid"`
	DocComplete int32                `json:"doc_complete"`
	JamaahList  []ManifestJamaahData `json:"jamaah_list"`
}

// ---------------------------------------------------------------------------
// GetDepartureManifest — GET /v1/manifest/:departure_id (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetDepartureManifest(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetDepartureManifest"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/manifest/:departure_id"),
		attribute.String("departure_id", departureID),
	)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	result, err := s.svc.GetDepartureManifest(ctx, departureID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeManifestError(c, span, err)
	}

	list := make([]ManifestJamaahData, 0, len(result.JamaahList))
	for _, j := range result.JamaahList {
		list = append(list, ManifestJamaahData{
			BookingID:     j.BookingID,
			Name:          j.Name,
			NIK:           j.NIK,
			Phone:         j.Phone,
			RoomType:      j.RoomType,
			BookingStatus: j.BookingStatus,
			DocStatus:     j.DocStatus,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": ManifestResponseData{
			DepartureID: result.DepartureID,
			TotalJamaah: result.TotalJamaah,
			LunasPaid:   result.LunasPaid,
			DocComplete: result.DocComplete,
			JamaahList:  list,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeManifestError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan jamaah sementara tidak tersedia"
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
