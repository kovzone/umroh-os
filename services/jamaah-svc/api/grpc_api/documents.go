// documents.go — gRPC handlers for UploadDocument and ReviewDocument RPCs (S3-E-02).
//
// F3-W1 / F3-W3 / BL-DOC-001..003.
//
// UploadDocument — stores document metadata in jamaah.pilgrim_documents.
// ReviewDocument — approve or reject a pending document.

package grpc_api

import (
	"context"

	"jamaah-svc/api/grpc_api/pb"
	"jamaah-svc/service"
	"jamaah-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// UploadDocument handles the UploadDocument RPC.
// Stores document metadata (file_path from GCS) in jamaah.pilgrim_documents.
func (s *Server) UploadDocument(ctx context.Context, req *pb.UploadDocumentRequest) (*pb.UploadDocumentResponse, error) {
	const op = "grpc_api.Server.UploadDocument"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("jamaah_id", req.GetJamaahId()),
		attribute.String("booking_id", req.GetBookingId()),
		attribute.String("doc_type", req.GetDocType()),
	)

	logger.Info().
		Str("op", op).
		Str("jamaah_id", req.GetJamaahId()).
		Str("doc_type", req.GetDocType()).
		Msg("")

	if req.GetJamaahId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "jamaah_id is required")
	}
	if req.GetBookingId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "booking_id is required")
	}
	if req.GetFilePath() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "file_path is required")
	}
	if req.GetDocType() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "doc_type is required")
	}

	result, err := s.svc.UploadDocument(ctx, &service.UploadDocumentParams{
		JamaahID:   req.GetJamaahId(),
		BookingID:  req.GetBookingId(),
		DocType:    req.GetDocType(),
		FilePath:   req.GetFilePath(),
		UploadedBy: req.GetUploadedBy(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("jamaah_id", req.GetJamaahId()).
			Err(err).
			Msg("UploadDocument failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to upload document: %v", err)
	}

	span.SetStatus(otelCodes.Ok, "uploaded")
	return &pb.UploadDocumentResponse{
		DocumentId: result.DocumentID,
		Status:     result.Status,
	}, nil
}

// ReviewDocument handles the ReviewDocument RPC.
// Approves or rejects a pending document. Audit log is written for every action.
func (s *Server) ReviewDocument(ctx context.Context, req *pb.ReviewDocumentRequest) (*pb.ReviewDocumentResponse, error) {
	const op = "grpc_api.Server.ReviewDocument"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("document_id", req.GetDocumentId()),
		attribute.String("action", req.GetAction()),
		attribute.String("reviewed_by", req.GetReviewedBy()),
	)

	logger.Info().
		Str("op", op).
		Str("document_id", req.GetDocumentId()).
		Str("action", req.GetAction()).
		Str("reviewed_by", req.GetReviewedBy()).
		Msg("")

	if req.GetDocumentId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "document_id is required")
	}
	if req.GetAction() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "action is required (approve|reject)")
	}
	if req.GetReviewedBy() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "reviewed_by is required")
	}

	result, err := s.svc.ReviewDocument(ctx, &service.ReviewDocumentParams{
		DocumentID:      req.GetDocumentId(),
		Action:          req.GetAction(),
		RejectionReason: req.GetRejectionReason(),
		ReviewedBy:      req.GetReviewedBy(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("document_id", req.GetDocumentId()).
			Err(err).
			Msg("ReviewDocument failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to review document: %v", err)
	}

	span.SetStatus(otelCodes.Ok, result.Status)
	return &pb.ReviewDocumentResponse{
		DocumentId: result.DocumentID,
		Status:     result.Status,
	}, nil
}
