// procurement_kit.go — gRPC handlers for procurement and kit assembly RPCs
// (BL-LOG-010..012).
//
// CreatePurchaseRequest: budget gate + total computation.
// ApprovePurchaseRequest: status gate (pending only), approve/reject.
// RecordGRNWithQC: QC flag routing, conditional AP journal signal.
// CreateKitAssembly: idempotent kit assembly with all-or-nothing items.

package grpc_api

import (
	"context"
	"errors"

	"logistics-svc/api/grpc_api/pb"
	"logistics-svc/service"
	"logistics-svc/util/logging"

	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/attribute"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// CreatePurchaseRequest (BL-LOG-010)
// ---------------------------------------------------------------------------

func (s *Server) CreatePurchaseRequest(ctx context.Context, req *pb.CreatePurchaseRequestRequest) (*pb.CreatePurchaseRequestResponse, error) {
	const op = "grpc_api.Server.CreatePurchaseRequest"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureId()))
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureId()).Msg("")

	result, err := s.svc.CreatePurchaseRequest(ctx, &service.CreatePurchaseRequestParams{
		DepartureID:  req.GetDepartureId(),
		RequestedBy:  req.GetRequestedBy(),
		ItemName:     req.GetItemName(),
		Quantity:     req.GetQuantity(),
		UnitPriceIdr: req.GetUnitPriceIdr(),
		BudgetLimit:  req.GetBudgetLimitIdr(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrBudgetExceeded):
			return nil, grpcStatus.Error(grpcCodes.FailedPrecondition, "budget_exceeded")
		case errors.Is(err, service.ErrInvalidQuantity):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_quantity")
		case errors.Is(err, service.ErrInvalidUnitPrice):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_unit_price")
		case errors.Is(err, service.ErrLogisticsMissingField):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "missing_required_field")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.CreatePurchaseRequestResponse{
		PrId:          result.PRID,
		Status:        result.Status,
		TotalPriceIdr: result.TotalPriceIdr,
	}, nil
}

// ---------------------------------------------------------------------------
// ApprovePurchaseRequest (BL-LOG-010)
// ---------------------------------------------------------------------------

func (s *Server) ApprovePurchaseRequest(ctx context.Context, req *pb.ApprovePurchaseRequestRequest) (*pb.ApprovePurchaseRequestResponse, error) {
	const op = "grpc_api.Server.ApprovePurchaseRequest"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("pr_id", req.GetPrId()))
	logger.Info().Str("op", op).Str("pr_id", req.GetPrId()).Msg("")

	if req.GetPrId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "pr_id is required")
	}

	result, err := s.svc.ApprovePurchaseRequest(ctx, &service.ApprovePurchaseRequestParams{
		PRID:       req.GetPrId(),
		ApprovedBy: req.GetApprovedBy(),
		Approved:   req.GetApproved(),
		Notes:      req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrPRNotFound):
			return nil, grpcStatus.Error(grpcCodes.NotFound, "not_found")
		case errors.Is(err, service.ErrPRNotPending):
			return nil, grpcStatus.Error(grpcCodes.FailedPrecondition, "invalid_status_for_decision")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.ApprovePurchaseRequestResponse{
		PrId:      result.PRID,
		NewStatus: result.NewStatus,
	}, nil
}

// ---------------------------------------------------------------------------
// RecordGRNWithQC (BL-LOG-011)
// ---------------------------------------------------------------------------

func (s *Server) RecordGRNWithQC(ctx context.Context, req *pb.RecordGRNWithQCRequest) (*pb.RecordGRNWithQCResponse, error) {
	const op = "grpc_api.Server.RecordGRNWithQC"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("grn_id", req.GetGrnId()),
		attribute.Bool("qc_passed", req.GetQcPassed()),
	)
	logger.Info().Str("op", op).Str("grn_id", req.GetGrnId()).Msg("")

	result, err := s.svc.RecordGRNWithQC(ctx, &service.RecordGRNWithQCParams{
		GRNID:       req.GetGrnId(),
		DepartureID: req.GetDepartureId(),
		AmountIdr:   req.GetAmountIdr(),
		QCPassed:    req.GetQcPassed(),
		QCNotes:     req.GetQcNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrGRNNotFound):
			return nil, grpcStatus.Error(grpcCodes.NotFound, "grn_not_found")
		case errors.Is(err, service.ErrInvalidAmount):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_amount")
		case errors.Is(err, service.ErrLogisticsMissingField):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "missing_required_field")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.RecordGRNWithQCResponse{
		GrnId:         result.GRNID,
		QcStatus:      result.QCStatus,
		JournalPosted: result.JournalPosted,
	}, nil
}

// ---------------------------------------------------------------------------
// CreateKitAssembly (BL-LOG-012)
// ---------------------------------------------------------------------------

func (s *Server) CreateKitAssembly(ctx context.Context, req *pb.CreateKitAssemblyRequest) (*pb.CreateKitAssemblyResponse, error) {
	const op = "grpc_api.Server.CreateKitAssembly"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureId()))
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureId()).Msg("")

	var items []service.KitItem
	for _, it := range req.GetItems() {
		items = append(items, service.KitItem{
			ItemName: it.ItemName,
			Quantity: it.Quantity,
		})
	}

	result, err := s.svc.CreateKitAssembly(ctx, &service.CreateKitAssemblyParams{
		DepartureID:    req.GetDepartureId(),
		AssembledBy:    req.GetAssembledBy(),
		Items:          items,
		IdempotencyKey: req.GetIdempotencyKey(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrEmptyItemsList):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "empty_items_list")
		case errors.Is(err, service.ErrInvalidItemQty):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_item_quantity")
		case errors.Is(err, service.ErrAssemblyFailed):
			return nil, grpcStatus.Error(grpcCodes.FailedPrecondition, "assembly_failed")
		case errors.Is(err, service.ErrLogisticsMissingField):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "missing_required_field")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.CreateKitAssemblyResponse{
		AssemblyId: result.AssemblyID,
		Status:     result.Status,
		Idempotent: result.Idempotent,
	}, nil
}
