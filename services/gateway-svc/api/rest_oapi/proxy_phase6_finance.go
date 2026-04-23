// proxy_phase6_finance.go — gateway REST handlers for finance Phase 6 routes
// (BL-FIN-010/011): AP disbursement and AR/AP aging.
//
// Route topology (all bearer-protected):
//   POST /v1/finance/disbursements              → CreateDisbursementBatch
//   PUT  /v1/finance/disbursements/:id/decision → ApproveDisbursement
//   GET  /v1/finance/aging                      → GetARAPAging
//
// Per ADR-0009: gateway is the single REST entry-point; finance-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/finance_grpc_adapter"
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

type DisbursementItemBody struct {
	VendorName  string `json:"vendor_name"`
	Description string `json:"description"`
	AmountIdr   int64  `json:"amount_idr"`
	Reference   string `json:"reference,omitempty"`
}

type CreateDisbursementBatchBody struct {
	Description string                 `json:"description"`
	Items       []DisbursementItemBody `json:"items"`
	CreatedBy   string                 `json:"created_by,omitempty"`
}

type ApproveDisbursementBody struct {
	ApprovedBy string `json:"approved_by"`
	Approved   bool   `json:"approved"`
	Notes      string `json:"notes,omitempty"`
}

// ---------------------------------------------------------------------------
// CreateDisbursementBatch — POST /v1/finance/disbursements (bearer)
// ---------------------------------------------------------------------------

func (s *Server) CreateDisbursementBatch(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateDisbursementBatch"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreateDisbursementBatchBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinancePhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	items := make([]*finance_grpc_adapter.DisbursementItemInput, 0, len(body.Items))
	for _, it := range body.Items {
		items = append(items, &finance_grpc_adapter.DisbursementItemInput{
			VendorName:  it.VendorName,
			Description: it.Description,
			AmountIdr:   it.AmountIdr,
			Reference:   it.Reference,
		})
	}

	result, err := s.svc.CreateDisbursementBatch(ctx, &finance_grpc_adapter.CreateDisbursementBatchParams{
		Description: body.Description,
		Items:       items,
		CreatedBy:   body.CreatedBy,
	})
	if err != nil {
		return writeFinancePhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"batch_id":         result.BatchID,
			"total_amount_idr": result.TotalAmountIdr,
			"item_count":       result.ItemCount,
			"status":           result.Status,
		},
	})
}

// ---------------------------------------------------------------------------
// ApproveDisbursement — PUT /v1/finance/disbursements/:id/decision (bearer)
// ---------------------------------------------------------------------------

func (s *Server) ApproveDisbursement(c *fiber.Ctx, batchID string) error {
	const op = "rest_oapi.Server.ApproveDisbursement"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("batch_id", batchID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("batch_id", batchID).Msg("")

	var body ApproveDisbursementBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinancePhase6Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.ApproveDisbursement(ctx, &finance_grpc_adapter.ApproveDisbursementParams{
		BatchID:    batchID,
		ApprovedBy: body.ApprovedBy,
		Approved:   body.Approved,
		Notes:      body.Notes,
	})
	if err != nil {
		return writeFinancePhase6Error(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"batch_id":          result.BatchID,
			"status":            result.Status,
			"journal_entry_ids": result.JournalEntryIDs,
		},
	})
}

// ---------------------------------------------------------------------------
// GetARAPAging — GET /v1/finance/aging (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetARAPAging(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetARAPAging"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	agingType := c.Query("type", "both") // "AR" | "AP" | "both"
	asOfDate := c.Query("as_of_date")    // "YYYY-MM-DD"

	result, err := s.svc.GetARAPAging(ctx, &finance_grpc_adapter.GetARAPAgingParams{
		Type:     agingType,
		AsOfDate: asOfDate,
	})
	if err != nil {
		return writeFinancePhase6Error(c, span, err)
	}

	mapBuckets := func(b *finance_grpc_adapter.AgingBuckets) fiber.Map {
		if b == nil {
			return nil
		}
		return fiber.Map{
			"current": b.Current,
			"days_30": b.Days30,
			"days_60": b.Days60,
			"days_90": b.Days90,
			"over_90": b.Over90,
			"total":   b.Total,
		}
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"ar":           mapBuckets(result.AR),
			"ap":           mapBuckets(result.AP),
			"generated_at": result.GeneratedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeFinancePhase6Error(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan finance sementara tidak tersedia"
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
