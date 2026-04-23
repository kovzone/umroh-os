// proxy_ops.go — gateway REST handlers for ops-svc RPCs (S3 Wave 2).
//
// Route topology (all bearer-protected):
//   POST /v1/ops/room-allocation                — RunRoomAllocation
//   GET  /v1/ops/room-allocation/:departure_id  — GetRoomAllocation
//   POST /v1/ops/id-cards                       — GenerateIDCard
//   POST /v1/ops/id-cards/verify                — VerifyIDCard
//   GET  /v1/ops/manifest/:departure_id/export  — ExportManifest
//
// Per ADR-0009: gateway is the single REST entry-point; ops-svc is pure gRPC.
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
// Request / response body types
// ---------------------------------------------------------------------------

// RunRoomAllocationBody is the JSON body for POST /v1/ops/room-allocation.
type RunRoomAllocationBody struct {
	DepartureID string   `json:"departure_id"`
	JamaahIDs   []string `json:"jamaah_ids"`
}

// RoomAssignmentData is the JSON representation of one room assignment.
type RoomAssignmentData struct {
	RoomNumber string `json:"room_number"`
	JamaahID   string `json:"jamaah_id"`
}

// RunRoomAllocationResponseData is the response for RunRoomAllocation.
type RunRoomAllocationResponseData struct {
	AllocationID string               `json:"allocation_id"`
	RoomCount    int32                `json:"room_count"`
	Assignments  []RoomAssignmentData `json:"assignments"`
}

// GetRoomAllocationResponseData is the response for GetRoomAllocation.
type GetRoomAllocationResponseData struct {
	AllocationID string               `json:"allocation_id"`
	RoomCount    int32                `json:"room_count"`
	Assignments  []RoomAssignmentData `json:"assignments"`
	Status       string               `json:"status"`
}

// GenerateIDCardBody is the JSON body for POST /v1/ops/id-cards.
type GenerateIDCardBody struct {
	JamaahID      string `json:"jamaah_id"`
	DepartureID   string `json:"departure_id"`
	CardType      string `json:"card_type"`
	JamaahName    string `json:"jamaah_name"`
	DepartureName string `json:"departure_name"`
}

// GenerateIDCardResponseData is the response for GenerateIDCard.
type GenerateIDCardResponseData struct {
	Token    string `json:"token"`
	QRData   string `json:"qr_data"`
	IssuedAt string `json:"issued_at"`
}

// VerifyIDCardBody is the JSON body for POST /v1/ops/id-cards/verify.
type VerifyIDCardBody struct {
	Token string `json:"token"`
}

// VerifyIDCardResponseData is the response for VerifyIDCard.
type VerifyIDCardResponseData struct {
	Valid        bool   `json:"valid"`
	JamaahID     string `json:"jamaah_id"`
	DepartureID  string `json:"departure_id"`
	CardType     string `json:"card_type"`
	ErrorReason  string `json:"error_reason,omitempty"`
}

// ManifestExportRowData is one row in the manifest export response.
type ManifestExportRowData struct {
	No          int32  `json:"no"`
	JamaahName  string `json:"jamaah_name"`
	PassportNo  string `json:"passport_no"`
	DocStatus   string `json:"doc_status"`
	RoomNumber  string `json:"room_number"`
}

// ExportManifestResponseData is the response for ExportManifest.
type ExportManifestResponseData struct {
	DepartureID string                  `json:"departure_id"`
	Rows        []ManifestExportRowData `json:"rows"`
}

// ---------------------------------------------------------------------------
// RunRoomAllocation — POST /v1/ops/room-allocation (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RunRoomAllocation(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RunRoomAllocation"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/ops/room-allocation"))
	logger.Info().Str("op", op).Msg("")

	var body RunRoomAllocationBody
	if err := c.BodyParser(&body); err != nil {
		return writeOpsError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RunRoomAllocation(ctx, body.DepartureID, body.JamaahIDs)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeOpsError(c, span, err)
	}

	assignments := make([]RoomAssignmentData, 0, len(result.Assignments))
	for _, a := range result.Assignments {
		assignments = append(assignments, RoomAssignmentData{
			RoomNumber: a.RoomNumber,
			JamaahID:   a.JamaahID,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": RunRoomAllocationResponseData{
			AllocationID: result.AllocationID,
			RoomCount:    result.RoomCount,
			Assignments:  assignments,
		},
	})
}

// ---------------------------------------------------------------------------
// GetRoomAllocation — GET /v1/ops/room-allocation/:departure_id (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetRoomAllocation(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetRoomAllocation"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/ops/room-allocation/:departure_id"),
		attribute.String("departure_id", departureID),
	)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	result, err := s.svc.GetRoomAllocation(ctx, departureID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeOpsError(c, span, err)
	}

	assignments := make([]RoomAssignmentData, 0, len(result.Assignments))
	for _, a := range result.Assignments {
		assignments = append(assignments, RoomAssignmentData{
			RoomNumber: a.RoomNumber,
			JamaahID:   a.JamaahID,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": GetRoomAllocationResponseData{
			AllocationID: result.AllocationID,
			RoomCount:    result.RoomCount,
			Assignments:  assignments,
			Status:       result.Status,
		},
	})
}

// ---------------------------------------------------------------------------
// GenerateIDCard — POST /v1/ops/id-cards (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GenerateIDCard(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GenerateIDCard"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/ops/id-cards"))
	logger.Info().Str("op", op).Msg("")

	var body GenerateIDCardBody
	if err := c.BodyParser(&body); err != nil {
		return writeOpsError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.GenerateIDCard(ctx, body.JamaahID, body.DepartureID, body.CardType, body.JamaahName, body.DepartureName)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeOpsError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": GenerateIDCardResponseData{
			Token:    result.Token,
			QRData:   result.QRData,
			IssuedAt: result.IssuedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// VerifyIDCard — POST /v1/ops/id-cards/verify (bearer)
// ---------------------------------------------------------------------------

func (s *Server) VerifyIDCard(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.VerifyIDCard"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/ops/id-cards/verify"))
	logger.Info().Str("op", op).Msg("")

	var body VerifyIDCardBody
	if err := c.BodyParser(&body); err != nil {
		return writeOpsError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.VerifyIDCard(ctx, body.Token)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeOpsError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": VerifyIDCardResponseData{
			Valid:       result.Valid,
			JamaahID:    result.JamaahID,
			DepartureID: result.DepartureID,
			CardType:    result.CardType,
			ErrorReason: result.ErrorReason,
		},
	})
}

// ---------------------------------------------------------------------------
// ExportManifest — GET /v1/ops/manifest/:departure_id/export (bearer)
// ---------------------------------------------------------------------------

func (s *Server) ExportManifest(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.ExportManifest"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/ops/manifest/:departure_id/export"),
		attribute.String("departure_id", departureID),
	)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	result, err := s.svc.ExportManifest(ctx, departureID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeOpsError(c, span, err)
	}

	rows := make([]ManifestExportRowData, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, ManifestExportRowData{
			No:         r.No,
			JamaahName: r.JamaahName,
			PassportNo: r.PassportNo,
			DocStatus:  r.DocStatus,
			RoomNumber: r.RoomNumber,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": ExportManifestResponseData{
			DepartureID: result.DepartureID,
			Rows:        rows,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeOpsError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan ops sementara tidak tersedia"
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
