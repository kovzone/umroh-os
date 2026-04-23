// disbursement.go — gRPC handlers for AP disbursement and AR/AP aging RPCs
// (BL-FIN-010/011).
//
// CreateDisbursementBatch: server-side total + bulk item insert.
// ApproveDisbursement:     transactional journal posting per item on approval.
// GetARAPAging:            aging bucket report parameterized by date + type.

package grpc_api

import (
	"context"
	"errors"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// CreateDisbursementBatch (BL-FIN-010)
// ---------------------------------------------------------------------------

func (s *Server) CreateDisbursementBatch(ctx context.Context, req *pb.CreateDisbursementBatchRequest) (*pb.CreateDisbursementBatchResponse, error) {
	const op = "grpc_api.Server.CreateDisbursementBatch"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("created_by", req.GetCreatedBy()))
	logger.Info().Str("op", op).Str("created_by", req.GetCreatedBy()).Msg("")

	var items []service.DisbursementItemInput
	for _, it := range req.GetItems() {
		items = append(items, service.DisbursementItemInput{
			VendorName:  it.VendorName,
			Description: it.Description,
			AmountIdr:   it.AmountIdr,
			Reference:   it.Reference,
		})
	}

	result, err := s.svc.CreateDisbursementBatch(ctx, &service.CreateDisbursementBatchParams{
		Description: req.GetDescription(),
		Items:       items,
		CreatedBy:   req.GetCreatedBy(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrEmptyDisbursementItems):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "empty_items_list")
		case errors.Is(err, service.ErrInvalidItemAmount):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_item_amount")
		case errors.Is(err, service.ErrDisbursementMissingField):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "missing_required_field")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.CreateDisbursementBatchResponse{
		BatchId:        result.BatchID,
		TotalAmountIdr: result.TotalAmountIdr,
		ItemCount:      result.ItemCount,
		Status:         result.Status,
	}, nil
}

// ---------------------------------------------------------------------------
// ApproveDisbursement (BL-FIN-010)
// ---------------------------------------------------------------------------

func (s *Server) ApproveDisbursement(ctx context.Context, req *pb.ApproveDisbursementRequest) (*pb.ApproveDisbursementResponse, error) {
	const op = "grpc_api.Server.ApproveDisbursement"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("batch_id", req.GetBatchId()))
	logger.Info().Str("op", op).Str("batch_id", req.GetBatchId()).Msg("")

	if req.GetBatchId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "batch_id is required")
	}

	result, err := s.svc.ApproveDisbursement(ctx, &service.ApproveDisbursementParams{
		BatchID:    req.GetBatchId(),
		ApprovedBy: req.GetApprovedBy(),
		Approved:   req.GetApproved(),
		Notes:      req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrDisbursementNotFound):
			return nil, grpcStatus.Error(grpcCodes.NotFound, "not_found")
		case errors.Is(err, service.ErrDisbursementNotPending):
			return nil, grpcStatus.Error(grpcCodes.FailedPrecondition, "invalid_status_for_decision")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.ApproveDisbursementResponse{
		BatchId:         result.BatchID,
		Status:          result.Status,
		JournalEntryIds: result.JournalEntryIDs,
	}, nil
}

// ---------------------------------------------------------------------------
// GetARAPAging (BL-FIN-011)
// ---------------------------------------------------------------------------

func (s *Server) GetARAPAging(ctx context.Context, req *pb.GetARAPAgingRequest) (*pb.GetARAPAgingResponse, error) {
	const op = "grpc_api.Server.GetARAPAging"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("type", req.GetType()))
	logger.Info().Str("op", op).Str("type", req.GetType()).Msg("")

	result, err := s.svc.GetARAPAging(ctx, &service.GetARAPAgingParams{
		Type:     req.GetType(),
		AsOfDate: req.GetAsOfDate(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrInvalidAgingType):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_type")
		case errors.Is(err, service.ErrInvalidDateFormat):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_date_format")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	resp := &pb.GetARAPAgingResponse{GeneratedAt: result.GeneratedAt}
	if result.AR != nil {
		resp.Ar = &pb.AgingBucketsPb{
			Current: result.AR.Current,
			Days30:  result.AR.Days30,
			Days60:  result.AR.Days60,
			Days90:  result.AR.Days90,
			Over90:  result.AR.Over90,
			Total:   result.AR.Total,
		}
	}
	if result.AP != nil {
		resp.Ap = &pb.AgingBucketsPb{
			Current: result.AP.Current,
			Days30:  result.AP.Days30,
			Days60:  result.AP.Days60,
			Days90:  result.AP.Days90,
			Over90:  result.AP.Over90,
			Total:   result.AP.Total,
		}
	}

	return resp, nil
}
