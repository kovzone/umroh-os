// proxy_logistics.go — gateway REST handlers for logistics-svc RPCs (S3 Wave 2 / ISSUE-018).
//
// Route topology (all bearer-protected):
//   GET  /v1/fulfillment-tasks           — ListFulfillmentTasks (ISSUE-018)
//   POST /v1/logistics/ship              — ShipFulfillmentTask
//   POST /v1/logistics/pickup-qr         — GeneratePickupQR
//   POST /v1/logistics/pickup-qr/redeem  — RedeemPickupQR
//
// Per ADR-0009: gateway is the single REST entry-point; logistics-svc is pure gRPC.
package rest_oapi

import (
	"errors"
	"strconv"

	"gateway-svc/adapter/logistics_grpc_adapter"
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

// ShipFulfillmentTaskBody is the JSON body for POST /v1/logistics/ship.
type ShipFulfillmentTaskBody struct {
	BookingID string `json:"booking_id"`
	Carrier   string `json:"carrier"`
	Notes     string `json:"notes"`
}

// ShipFulfillmentTaskResponseData is the response for ShipFulfillmentTask.
type ShipFulfillmentTaskResponseData struct {
	ShipmentID     string `json:"shipment_id"`
	TrackingNumber string `json:"tracking_number"`
	Status         string `json:"status"`
}

// GeneratePickupQRBody is the JSON body for POST /v1/logistics/pickup-qr.
type GeneratePickupQRBody struct {
	BookingID string `json:"booking_id"`
}

// GeneratePickupQRResponseData is the response for GeneratePickupQR.
type GeneratePickupQRResponseData struct {
	PickupTokenID string `json:"pickup_token_id"`
	Token         string `json:"token"`
	ExpiresAt     string `json:"expires_at"`
}

// RedeemPickupQRBody is the JSON body for POST /v1/logistics/pickup-qr/redeem.
type RedeemPickupQRBody struct {
	Token string `json:"token"`
}

// RedeemPickupQRResponseData is the response for RedeemPickupQR.
type RedeemPickupQRResponseData struct {
	Redeemed    bool   `json:"redeemed"`
	BookingID   string `json:"booking_id"`
	TaskID      string `json:"task_id"`
	ErrorReason string `json:"error_reason,omitempty"`
}

// ---------------------------------------------------------------------------
// ListFulfillmentTasks — GET /v1/fulfillment-tasks (bearer) — ISSUE-018
// ---------------------------------------------------------------------------

// FulfillmentTaskResponseItem is a single task in the API response.
type FulfillmentTaskResponseItem struct {
	ID             string `json:"id"`
	BookingID      string `json:"booking_id"`
	DepartureID    string `json:"departure_id"`
	Status         string `json:"status"`
	TrackingNumber string `json:"tracking_number,omitempty"`
	ShippedAt      string `json:"shipped_at,omitempty"`
	DeliveredAt    string `json:"delivered_at,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func (s *Server) ListFulfillmentTasks(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListFulfillmentTasks"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "GET /v1/fulfillment-tasks"))
	logger.Info().Str("op", op).Msg("")

	limit := int32(50)
	offset := int32(0)
	if l := c.QueryInt("limit", 50); l > 0 {
		limit = int32(l)
	}
	if o, err := strconv.Atoi(c.Query("offset", "0")); err == nil && o >= 0 {
		offset = int32(o)
	}

	result, err := s.svc.ListFulfillmentTasks(ctx, &logistics_grpc_adapter.ListFulfillmentTasksParams{
		StatusFilter:      c.Query("status", ""),
		DepartureIDFilter: c.Query("departure_id", ""),
		Limit:             limit,
		Offset:            offset,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeLogisticsError(c, span, err)
	}

	items := make([]FulfillmentTaskResponseItem, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		items = append(items, FulfillmentTaskResponseItem{
			ID:             t.ID,
			BookingID:      t.BookingID,
			DepartureID:    t.DepartureID,
			Status:         t.Status,
			TrackingNumber: t.TrackingNumber,
			ShippedAt:      t.ShippedAt,
			DeliveredAt:    t.DeliveredAt,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": items,
		"total": result.Total,
	})
}

// ---------------------------------------------------------------------------
// ShipFulfillmentTask — POST /v1/logistics/ship (bearer)
// ---------------------------------------------------------------------------

func (s *Server) ShipFulfillmentTask(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ShipFulfillmentTask"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/logistics/ship"))
	logger.Info().Str("op", op).Msg("")

	var body ShipFulfillmentTaskBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.ShipFulfillmentTask(ctx, body.BookingID, body.Carrier, body.Notes)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeLogisticsError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": ShipFulfillmentTaskResponseData{
			ShipmentID:     result.ShipmentID,
			TrackingNumber: result.TrackingNumber,
			Status:         result.Status,
		},
	})
}

// ---------------------------------------------------------------------------
// GeneratePickupQR — POST /v1/logistics/pickup-qr (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GeneratePickupQR(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GeneratePickupQR"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/logistics/pickup-qr"))
	logger.Info().Str("op", op).Msg("")

	var body GeneratePickupQRBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.GeneratePickupQR(ctx, body.BookingID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeLogisticsError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": GeneratePickupQRResponseData{
			PickupTokenID: result.PickupTokenID,
			Token:         result.Token,
			ExpiresAt:     result.ExpiresAt,
		},
	})
}

// ---------------------------------------------------------------------------
// RedeemPickupQR — POST /v1/logistics/pickup-qr/redeem (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RedeemPickupQR(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RedeemPickupQR"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/logistics/pickup-qr/redeem"))
	logger.Info().Str("op", op).Msg("")

	var body RedeemPickupQRBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RedeemPickupQR(ctx, body.Token)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeLogisticsError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": RedeemPickupQRResponseData{
			Redeemed:    result.Redeemed,
			BookingID:   result.BookingID,
			TaskID:      result.TaskID,
			ErrorReason: result.ErrorReason,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeLogisticsError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan logistik sementara tidak tersedia"
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
