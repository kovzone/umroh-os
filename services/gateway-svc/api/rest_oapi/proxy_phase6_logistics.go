// proxy_phase6_logistics.go — gateway REST handlers for logistics Phase 6 routes
// (BL-LOG-010..012): purchase requests, GRN with QC, and kit assembly.
//
// Route topology (all bearer-protected):
//   POST /v1/logistics/purchase-requests              → CreatePurchaseRequest
//   PUT  /v1/logistics/purchase-requests/:id/decision → ApprovePurchaseRequest
//   POST /v1/logistics/grn-qc                         → RecordGRNWithQC
//   POST /v1/logistics/kit-assembly                   → CreateKitAssembly
//
// Per ADR-0009: gateway is the single REST entry-point; logistics-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/logistics_grpc_adapter"
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

type CreatePurchaseRequestBody struct {
	DepartureID    string `json:"departure_id"`
	RequestedBy    string `json:"requested_by,omitempty"`
	ItemName       string `json:"item_name"`
	Quantity       int32  `json:"quantity"`
	UnitPriceIdr   int64  `json:"unit_price_idr"`
	BudgetLimitIdr int64  `json:"budget_limit_idr,omitempty"`
}

type ApprovePurchaseRequestBody struct {
	ApprovedBy string `json:"approved_by"`
	Approved   bool   `json:"approved"`
	Notes      string `json:"notes,omitempty"`
}

type RecordGRNWithQCBody struct {
	GrnID       string `json:"grn_id"`
	DepartureID string `json:"departure_id"`
	AmountIdr   int64  `json:"amount_idr"`
	QcPassed    bool   `json:"qc_passed"`
	QcNotes     string `json:"qc_notes,omitempty"`
}

type KitItemBody struct {
	ItemName string `json:"item_name"`
	Quantity int32  `json:"quantity"`
}

type CreateKitAssemblyBody struct {
	DepartureID    string        `json:"departure_id"`
	AssembledBy    string        `json:"assembled_by,omitempty"`
	Items          []KitItemBody `json:"items"`
	IdempotencyKey string        `json:"idempotency_key,omitempty"`
}

// ---------------------------------------------------------------------------
// CreatePurchaseRequest — POST /v1/logistics/purchase-requests (bearer)
// ---------------------------------------------------------------------------

func (s *Server) CreatePurchaseRequest(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreatePurchaseRequest"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreatePurchaseRequestBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsPhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreatePurchaseRequest(ctx, &logistics_grpc_adapter.CreatePurchaseRequestParams{
		DepartureID:    body.DepartureID,
		RequestedBy:    body.RequestedBy,
		ItemName:       body.ItemName,
		Quantity:       body.Quantity,
		UnitPriceIdr:   body.UnitPriceIdr,
		BudgetLimitIdr: body.BudgetLimitIdr,
	})
	if err != nil {
		return writeLogisticsPhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"pr_id":           result.PrID,
			"status":          result.Status,
			"total_price_idr": result.TotalPriceIdr,
		},
	})
}

// ---------------------------------------------------------------------------
// ApprovePurchaseRequest — PUT /v1/logistics/purchase-requests/:id/decision (bearer)
// ---------------------------------------------------------------------------

func (s *Server) ApprovePurchaseRequest(c *fiber.Ctx, prID string) error {
	const op = "rest_oapi.Server.ApprovePurchaseRequest"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("pr_id", prID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("pr_id", prID).Msg("")

	var body ApprovePurchaseRequestBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsPhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.ApprovePurchaseRequest(ctx, &logistics_grpc_adapter.ApprovePurchaseRequestParams{
		PrID:       prID,
		ApprovedBy: body.ApprovedBy,
		Approved:   body.Approved,
		Notes:      body.Notes,
	})
	if err != nil {
		return writeLogisticsPhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"pr_id":      result.PrID,
			"new_status": result.NewStatus,
		},
	})
}

// ---------------------------------------------------------------------------
// RecordGRNWithQC — POST /v1/logistics/grn-qc (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RecordGRNWithQC(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordGRNWithQC"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordGRNWithQCBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsPhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordGRNWithQC(ctx, &logistics_grpc_adapter.RecordGRNWithQCParams{
		GrnID:       body.GrnID,
		DepartureID: body.DepartureID,
		AmountIdr:   body.AmountIdr,
		QcPassed:    body.QcPassed,
		QcNotes:     body.QcNotes,
	})
	if err != nil {
		return writeLogisticsPhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"grn_id":         result.GrnID,
			"qc_status":      result.QcStatus,
			"journal_posted": result.JournalPosted,
		},
	})
}

// ---------------------------------------------------------------------------
// CreateKitAssembly — POST /v1/logistics/kit-assembly (bearer)
// ---------------------------------------------------------------------------

func (s *Server) CreateKitAssembly(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateKitAssembly"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreateKitAssemblyBody
	if err := c.BodyParser(&body); err != nil {
		return writeLogisticsPhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	items := make([]*logistics_grpc_adapter.KitItem, 0, len(body.Items))
	for _, it := range body.Items {
		items = append(items, &logistics_grpc_adapter.KitItem{
			ItemName: it.ItemName,
			Quantity: it.Quantity,
		})
	}

	result, err := s.svc.CreateKitAssembly(ctx, &logistics_grpc_adapter.CreateKitAssemblyParams{
		DepartureID:    body.DepartureID,
		AssembledBy:    body.AssembledBy,
		Items:          items,
		IdempotencyKey: body.IdempotencyKey,
	})
	if err != nil {
		return writeLogisticsPhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"assembly_id": result.AssemblyID,
			"status":      result.Status,
			"idempotent":  result.Idempotent,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeLogisticsPhase6Error(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan logistics sementara tidak tersedia"
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
