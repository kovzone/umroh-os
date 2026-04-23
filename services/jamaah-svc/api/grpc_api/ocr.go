// ocr.go — gRPC handlers for TriggerOCR and GetOCRStatus RPCs (BL-DOC-002).
//
// TriggerOCR — stub OCR: generates fake extracted fields + confidence score.
// In the real implementation, this would call an OCR provider API.
// Stub behavior:
//   - confidence = 0.85 for passport docs (>= 0.8 threshold → 'ocr_complete')
//   - confidence = 0.70 for other types   (< 0.8 threshold  → 'pending', manual review)
//
// GetOCRStatus — returns current OCR state including extracted fields.

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

// TriggerOCR handles the TriggerOCR RPC.
// Runs a stub OCR pass on the given document. If confidence >= 0.8, status
// becomes 'ocr_complete'; otherwise stays 'pending' for manual review.
func (s *Server) TriggerOCR(ctx context.Context, req *pb.TriggerOCRRequest) (*pb.TriggerOCRResponse, error) {
	const op = "grpc_api.Server.TriggerOCR"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("document_id", req.GetDocumentId()),
	)

	logger.Info().
		Str("op", op).
		Str("document_id", req.GetDocumentId()).
		Msg("")

	if req.GetDocumentId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "document_id is required")
	}

	result, err := s.svc.TriggerOCR(ctx, &service.TriggerOCRParams{
		DocumentID: req.GetDocumentId(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("document_id", req.GetDocumentId()).
			Err(err).
			Msg("TriggerOCR failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to trigger OCR: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("document_id", result.DocumentID).
		Str("status", result.Status).
		Float64("confidence", result.Confidence).
		Msg("TriggerOCR succeeded")

	span.SetStatus(otelCodes.Ok, result.Status)
	return &pb.TriggerOCRResponse{
		DocumentId: result.DocumentID,
		Status:     result.Status,
		Confidence: result.Confidence,
		OcrResult:  result.OcrResult,
	}, nil
}

// GetOCRStatus handles the GetOCRStatus RPC.
// Returns the current OCR state and extracted fields for a document.
func (s *Server) GetOCRStatus(ctx context.Context, req *pb.GetOCRStatusRequest) (*pb.GetOCRStatusResponse, error) {
	const op = "grpc_api.Server.GetOCRStatus"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("document_id", req.GetDocumentId()),
	)

	logger.Info().
		Str("op", op).
		Str("document_id", req.GetDocumentId()).
		Msg("")

	if req.GetDocumentId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "document_id is required")
	}

	result, err := s.svc.GetOCRStatus(ctx, &service.GetOCRStatusParams{
		DocumentID: req.GetDocumentId(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("document_id", req.GetDocumentId()).
			Err(err).
			Msg("GetOCRStatus failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to get OCR status: %v", err)
	}

	span.SetStatus(otelCodes.Ok, result.Status)
	return &pb.GetOCRStatusResponse{
		DocumentId: result.DocumentID,
		Status:     result.Status,
		Confidence: result.Confidence,
		OcrResult:  result.OcrResult,
	}, nil
}
