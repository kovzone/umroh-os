// proxy_phase6_ops.go — gateway REST handlers for ops Phase 6 routes
// (BL-OPS-010/011): scan events and bus boarding.
//
// Route topology (all bearer-protected):
//   POST /v1/ops/scans                          → RecordScan
//   POST /v1/ops/bus-boarding                   → RecordBusBoarding
//   GET  /v1/ops/bus-boarding/:departure_id     → GetBoardingRoster
//
// Per ADR-0009: gateway is the single REST entry-point; ops-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/ops_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request body types
// ---------------------------------------------------------------------------

type RecordScanBody struct {
	ScanType       string `json:"scan_type"`
	DepartureID    string `json:"departure_id"`
	JamaahID       string `json:"jamaah_id"`
	ScannedBy      string `json:"scanned_by,omitempty"`
	DeviceID       string `json:"device_id,omitempty"`
	Location       string `json:"location,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

type RecordBusBoardingBody struct {
	DepartureID string `json:"departure_id"`
	BusNumber   string `json:"bus_number"`
	JamaahID    string `json:"jamaah_id"`
	ScannedBy   string `json:"scanned_by,omitempty"`
	Status      string `json:"status,omitempty"` // "boarded" | "absent"
}

// ---------------------------------------------------------------------------
// RecordScan — POST /v1/ops/scans (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RecordScan(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordScan"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordScanBody
	if err := c.BodyParser(&body); err != nil {
		return writeOpsPhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordScan(ctx, &ops_grpc_adapter.RecordScanParams{
		ScanType:       body.ScanType,
		DepartureID:    body.DepartureID,
		JamaahID:       body.JamaahID,
		ScannedBy:      body.ScannedBy,
		DeviceID:       body.DeviceID,
		Location:       body.Location,
		IdempotencyKey: body.IdempotencyKey,
	})
	if err != nil {
		return writeOpsPhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"scan_id":    result.ScanID,
			"idempotent": result.Idempotent,
		},
	})
}

// ---------------------------------------------------------------------------
// RecordBusBoarding — POST /v1/ops/bus-boarding (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RecordBusBoarding(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordBusBoarding"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordBusBoardingBody
	if err := c.BodyParser(&body); err != nil {
		return writeOpsPhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordBusBoarding(ctx, &ops_grpc_adapter.RecordBusBoardingParams{
		DepartureID: body.DepartureID,
		BusNumber:   body.BusNumber,
		JamaahID:    body.JamaahID,
		ScannedBy:   body.ScannedBy,
		Status:      body.Status,
	})
	if err != nil {
		return writeOpsPhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"boarding_id": result.BoardingID,
			"status":      result.Status,
			"idempotent":  result.Idempotent,
		},
	})
}

// ---------------------------------------------------------------------------
// GetBoardingRoster — GET /v1/ops/bus-boarding/:departure_id (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetBoardingRoster(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetBoardingRoster"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	busNumber := c.Query("bus_number") // optional filter

	result, err := s.svc.GetBoardingRoster(ctx, departureID, busNumber)
	if err != nil {
		return writeOpsPhase6Error(c, span, err)
	}

	type boardingJSON struct {
		JamaahID  string `json:"jamaah_id"`
		BusNumber string `json:"bus_number"`
		Status    string `json:"status"`
		BoardedAt string `json:"boarded_at,omitempty"`
	}
	data := make([]boardingJSON, 0, len(result.Boardings))
	for _, b := range result.Boardings {
		data = append(data, boardingJSON{
			JamaahID:  b.JamaahID,
			BusNumber: b.BusNumber,
			Status:    b.Status,
			BoardedAt: b.BoardedAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"boardings":     data,
			"total_boarded": result.TotalBoarded,
			"total_absent":  result.TotalAbsent,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeOpsPhase6Error(c *fiber.Ctx, span trace.Span, err error) error {
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
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
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
