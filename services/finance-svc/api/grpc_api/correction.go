// correction.go — gRPC handlers for CorrectJournal + DeleteJournalEntry
// (BL-FIN-006).
//
// CorrectJournal:    posts reversing counter-entry, preserving audit trail.
// DeleteJournalEntry: anti-delete guard — always returns PermissionDenied.

package grpc_api

import (
	"context"
	"errors"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"
	"finance-svc/util/apperrors"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
	grpcCodes "google.golang.org/grpc/codes"
)

// CorrectJournal posts a reversing counter-entry for an existing journal entry.
// Implements pb.CorrectionHandler.
func (s *Server) CorrectJournal(ctx context.Context, req *pb.CorrectJournalRequest) (*pb.CorrectJournalResponse, error) {
	const op = "grpc_api.Server.CorrectJournal"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("entry_id", req.GetEntryId()),
	)
	logger.Info().Str("op", op).Str("entry_id", req.GetEntryId()).Msg("")

	if req.GetEntryId() == "" {
		return nil, status.Error(grpcCodes.InvalidArgument, "entry_id is required")
	}

	result, err := s.svc.CorrectJournal(ctx, &service.CorrectJournalParams{
		EntryID:     req.GetEntryId(),
		Reason:      req.GetReason(),
		ActorUserID: req.GetActorUserId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			return nil, status.Error(grpcCodes.NotFound, err.Error())
		case errors.Is(err, apperrors.ErrValidation):
			return nil, status.Error(grpcCodes.InvalidArgument, err.Error())
		case errors.Is(err, apperrors.ErrForbidden):
			return nil, status.Error(grpcCodes.PermissionDenied, err.Error())
		default:
			return nil, status.Error(grpcCodes.Internal, err.Error())
		}
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CorrectJournalResponse{
		CorrectionEntryId: result.CorrectionEntryID,
		OriginalEntryId:   result.OriginalEntryID,
		Idempotent:        result.Idempotent,
	}, nil
}

// DeleteJournalEntry is the anti-delete guard. It always returns PermissionDenied.
// Journal entries are immutable; corrections must use CorrectJournal.
// Implements pb.CorrectionHandler.
func (s *Server) DeleteJournalEntry(_ context.Context, req *pb.DeleteJournalEntryRequest) (*pb.DeleteJournalEntryResponse, error) {
	return nil, status.Errorf(grpcCodes.PermissionDenied,
		"journal entries cannot be deleted (entry_id=%s); use CorrectJournal to post a reversing counter-entry",
		req.GetEntryId())
}
