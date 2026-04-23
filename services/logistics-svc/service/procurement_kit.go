// procurement_kit.go — logistics-svc service logic for purchase requests,
// GRN with QC, and kit assembly (BL-LOG-010..012).

package service

import (
	"context"
	"errors"
	"fmt"

	"logistics-svc/store/postgres_store/sqlc"
	"logistics-svc/util/logging"
	"logistics-svc/util/ulid"

	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Params / Results
// ---------------------------------------------------------------------------

type CreatePurchaseRequestParams struct {
	DepartureID   string
	RequestedBy   string
	ItemName      string
	Quantity      int32
	UnitPriceIdr  int64
	BudgetLimit   int64 // 0 = no limit
}

type CreatePurchaseRequestResult struct {
	PRID          string
	Status        string
	TotalPriceIdr int64
}

type ApprovePurchaseRequestParams struct {
	PRID       string
	ApprovedBy string
	Approved   bool
	Notes      string
}

type ApprovePurchaseRequestResult struct {
	PRID      string
	NewStatus string
}

type RecordGRNWithQCParams struct {
	GRNID       string
	DepartureID string
	AmountIdr   int64
	QCPassed    bool
	QCNotes     string
}

type RecordGRNWithQCResult struct {
	GRNID             string
	QCStatus          string
	JournalPosted     bool
	RequiresManualCorrection bool
}

type KitItem struct {
	ItemName string
	Quantity int32
}

type CreateKitAssemblyParams struct {
	DepartureID    string
	AssembledBy    string
	Items          []KitItem
	IdempotencyKey string
}

type CreateKitAssemblyResult struct {
	AssemblyID string
	Status     string
	Idempotent bool
}

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

var (
	ErrBudgetExceeded         = errors.New("budget_exceeded")
	ErrInvalidQuantity        = errors.New("invalid_quantity")
	ErrInvalidUnitPrice       = errors.New("invalid_unit_price")
	ErrLogisticsMissingField  = errors.New("missing_required_field")
	ErrPRNotPending           = errors.New("invalid_status_for_decision")
	ErrPRNotFound             = errors.New("not_found")
	ErrEmptyItemsList         = errors.New("empty_items_list")
	ErrInvalidItemQty         = errors.New("invalid_item_quantity")
	ErrAssemblyFailed         = errors.New("assembly_failed")
	ErrGRNNotFound            = errors.New("grn_not_found")
	ErrInvalidAmount          = errors.New("invalid_amount")
)

// ---------------------------------------------------------------------------
// CreatePurchaseRequest (BL-LOG-010)
// ---------------------------------------------------------------------------

func (s *Service) CreatePurchaseRequest(ctx context.Context, params *CreatePurchaseRequestParams) (*CreatePurchaseRequestResult, error) {
	const op = "service.CreatePurchaseRequest"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.DepartureID == "" || params.RequestedBy == "" || params.ItemName == "" {
		return nil, ErrLogisticsMissingField
	}
	if params.Quantity < 1 {
		return nil, ErrInvalidQuantity
	}
	if params.UnitPriceIdr < 1 {
		return nil, ErrInvalidUnitPrice
	}

	total := int64(params.Quantity) * params.UnitPriceIdr
	if params.BudgetLimit > 0 && total > params.BudgetLimit {
		return nil, ErrBudgetExceeded
	}

	id, err := ulid.New("pr_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertPurchaseRequest(ctx, sqlc.InsertPurchaseRequestParams{
		ID:            id,
		DepartureID:   params.DepartureID,
		RequestedBy:   params.RequestedBy,
		ItemName:      params.ItemName,
		Quantity:      params.Quantity,
		UnitPriceIdr:  params.UnitPriceIdr,
		TotalPriceIdr: total,
		BudgetLimit:   params.BudgetLimit,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert PR")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	return &CreatePurchaseRequestResult{
		PRID:          row.ID,
		Status:        row.Status,
		TotalPriceIdr: row.TotalPriceIdr,
	}, nil
}

// ---------------------------------------------------------------------------
// ApprovePurchaseRequest (BL-LOG-010)
// ---------------------------------------------------------------------------

func (s *Service) ApprovePurchaseRequest(ctx context.Context, params *ApprovePurchaseRequestParams) (*ApprovePurchaseRequestResult, error) {
	const op = "service.ApprovePurchaseRequest"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.PRID == "" {
		return nil, ErrPRNotFound
	}

	pr, err := s.store.GetPurchaseRequestByID(ctx, params.PRID)
	if err != nil {
		return nil, ErrPRNotFound
	}
	if pr.Status != "pending" {
		return nil, ErrPRNotPending
	}

	newStatus := "approved"
	if !params.Approved {
		newStatus = "rejected"
	}

	if err := s.store.UpdatePurchaseRequestDecision(ctx, sqlc.UpdatePurchaseRequestDecisionParams{
		ID:         params.PRID,
		NewStatus:  newStatus,
		ApprovedBy: params.ApprovedBy,
		Notes:      params.Notes,
	}); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("update PR decision")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	return &ApprovePurchaseRequestResult{PRID: params.PRID, NewStatus: newStatus}, nil
}

// ---------------------------------------------------------------------------
// RecordGRNWithQC (BL-LOG-011)
// ---------------------------------------------------------------------------

// RecordGRNWithQC validates QC and returns whether AP journal should be posted.
// The gateway is responsible for calling finance-svc.OnGRNReceived when
// QCPassed=true (same pattern as existing S3 /v1/finance/grn route).
func (s *Service) RecordGRNWithQC(ctx context.Context, params *RecordGRNWithQCParams) (*RecordGRNWithQCResult, error) {
	const op = "service.RecordGRNWithQC"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.GRNID == "" || params.DepartureID == "" {
		return nil, ErrLogisticsMissingField
	}
	if params.AmountIdr < 1 {
		return nil, ErrInvalidAmount
	}

	if !params.QCPassed {
		logger.Info().Str("op", op).Str("grn_id", params.GRNID).Msg("QC failed — no AP journal posted")
		return &RecordGRNWithQCResult{
			GRNID:         params.GRNID,
			QCStatus:      "rejected",
			JournalPosted: false,
		}, nil
	}

	// QC passed — gateway will post AP journal via finance-svc
	return &RecordGRNWithQCResult{
		GRNID:         params.GRNID,
		QCStatus:      "passed",
		JournalPosted: true,
	}, nil
}

// ---------------------------------------------------------------------------
// CreateKitAssembly (BL-LOG-012)
// ---------------------------------------------------------------------------

func (s *Service) CreateKitAssembly(ctx context.Context, params *CreateKitAssemblyParams) (*CreateKitAssemblyResult, error) {
	const op = "service.CreateKitAssembly"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.DepartureID == "" || params.AssembledBy == "" || params.IdempotencyKey == "" {
		return nil, ErrLogisticsMissingField
	}
	if len(params.Items) == 0 {
		return nil, ErrEmptyItemsList
	}
	for _, item := range params.Items {
		if item.Quantity < 1 {
			return nil, ErrInvalidItemQty
		}
	}

	id, err := ulid.New("kit_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	assembly, wasInserted, err := s.store.InsertKitAssemblyIdempotent(ctx, sqlc.InsertKitAssemblyParams{
		ID:             id,
		DepartureID:    params.DepartureID,
		AssembledBy:    params.AssembledBy,
		Status:         "pending",
		IdempotencyKey: params.IdempotencyKey,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert kit assembly")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	if !wasInserted {
		return &CreateKitAssemblyResult{
			AssemblyID: assembly.ID,
			Status:     assembly.Status,
			Idempotent: true,
		}, nil
	}

	// Insert items and mark fulfilled — in a real system would check stock
	for _, item := range params.Items {
		if err := s.store.InsertKitAssemblyItem(ctx, sqlc.InsertKitAssemblyItemParams{
			AssemblyID: assembly.ID,
			ItemName:   item.ItemName,
			Quantity:   item.Quantity,
			Fulfilled:  true,
		}); err != nil {
			logger.Error().Err(err).Str("op", op).Msg("insert kit assembly item")
			// Mark assembly failed
			_ = s.store.UpdateKitAssemblyStatus(ctx, "failed", assembly.ID)
			return nil, ErrAssemblyFailed
		}
	}

	// All items fulfilled → complete
	if err := s.store.UpdateKitAssemblyStatus(ctx, "completed", assembly.ID); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("mark kit assembly completed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	return &CreateKitAssemblyResult{
		AssemblyID: assembly.ID,
		Status:     "completed",
		Idempotent: false,
	}, nil
}
