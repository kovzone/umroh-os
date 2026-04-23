// ocr.go — gateway adapter methods for jamaah-svc OCR RPCs (S3 Wave 2).
//
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
package jamaah_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/jamaah_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// TriggerOCRResult is the gateway-local result for TriggerOCR.
type TriggerOCRResult struct {
	DocumentID string
	Status     string
	Confidence float64
	OcrResult  map[string]string
}

// GetOCRStatusResult is the gateway-local result for GetOCRStatus.
type GetOCRStatusResult struct {
	DocumentID string
	Status     string
	Confidence float64
	OcrResult  map[string]string
}

// ---------------------------------------------------------------------------
// TriggerOCR
// ---------------------------------------------------------------------------

// TriggerOCR initiates OCR processing on a passport document.
func (a *Adapter) TriggerOCR(ctx context.Context, documentID string) (*TriggerOCRResult, error) {
	const op = "jamaah_grpc_adapter.Adapter.TriggerOCR"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "TriggerOCR"),
		attribute.String("document_id", documentID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.ocrClient.TriggerOCR(ctx, &pb.TriggerOCRRequest{
		DocumentId: documentID,
	})
	if err != nil {
		wrapped := mapJamaahError(err)
		logger.Warn().Err(wrapped).Str("document_id", documentID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &TriggerOCRResult{
		DocumentID: resp.GetDocumentId(),
		Status:     resp.GetStatus(),
		Confidence: resp.GetConfidence(),
		OcrResult:  resp.GetOcrResult(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetOCRStatus
// ---------------------------------------------------------------------------

// GetOCRStatus retrieves OCR results and confidence score for a document.
func (a *Adapter) GetOCRStatus(ctx context.Context, documentID string) (*GetOCRStatusResult, error) {
	const op = "jamaah_grpc_adapter.Adapter.GetOCRStatus"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetOCRStatus"),
		attribute.String("document_id", documentID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.ocrClient.GetOCRStatus(ctx, &pb.GetOCRStatusRequest{
		DocumentId: documentID,
	})
	if err != nil {
		wrapped := mapJamaahError(err)
		logger.Warn().Err(wrapped).Str("document_id", documentID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetOCRStatusResult{
		DocumentID: resp.GetDocumentId(),
		Status:     resp.GetStatus(),
		Confidence: resp.GetConfidence(),
		OcrResult:  resp.GetOcrResult(),
	}, nil
}
