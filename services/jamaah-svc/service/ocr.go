// ocr.go — TriggerOCR and GetOCRStatus service-layer implementations.
//
// BL-DOC-002 / S3-E-02 — Passport OCR stub.
//
// TriggerOCR — generates stub OCR results and confidence score:
//   - doc_type = "passport" → confidence = 0.85 (above threshold)
//   - all other types       → confidence = 0.70 (below threshold)
//
// Confidence < 0.8: status stays 'pending' (manual review queue).
// Confidence >= 0.8: status updates to 'ocr_complete'.
//
// GetOCRStatus — reads the current OCR state of a document.
//
// In a real implementation, TriggerOCR would call an OCR provider API
// (e.g. Google Vision, AWS Textract) and parse the response.

package service

import (
	"context"
	"encoding/json"
	"fmt"

	"jamaah-svc/store/postgres_store/sqlc"
	"jamaah-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

const ocrConfidenceThreshold = 0.8

// ---------------------------------------------------------------------------
// TriggerOCR
// ---------------------------------------------------------------------------

// TriggerOCRParams holds the inputs for triggering OCR on a document.
type TriggerOCRParams struct {
	DocumentID string
}

// TriggerOCRResult holds the result of a TriggerOCR call.
type TriggerOCRResult struct {
	DocumentID string
	Status     string  // "ocr_complete" or "pending"
	Confidence float64
	OcrResult  map[string]string
}

// TriggerOCR runs a stub OCR pass on the given document.
// For passport documents, confidence = 0.85; for other types, 0.70.
// If confidence >= 0.8, the document status is updated to 'ocr_complete';
// otherwise it stays 'pending' for manual review.
func (svc *Service) TriggerOCR(ctx context.Context, params *TriggerOCRParams) (*TriggerOCRResult, error) {
	const op = "service.Service.TriggerOCR"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("document_id", params.DocumentID),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.DocumentID == "" {
		return nil, fmt.Errorf("%s: document_id is required", op)
	}

	// Fetch document to determine doc_type.
	doc, err := svc.store.GetPilgrimDocumentByID(ctx, params.DocumentID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get document: %w", op, err)
	}

	// --- Stub OCR logic ---
	// Passport docs have a higher confidence score than other doc types.
	// In a real implementation, this would call an external OCR API.
	var confidence float64
	var ocrResult map[string]string

	if doc.DocType == "passport" {
		confidence = 0.85
		ocrResult = map[string]string{
			"name":        "STUB PILGRIM NAME",
			"passport_no": "A1234567",
			"dob":         "1980-01-01",
			"nationality": "INDONESIA",
			"expiry":      "2030-12-31",
			"gender":      "M",
		}
	} else {
		confidence = 0.70
		ocrResult = map[string]string{
			"raw_text": "OCR stub — low confidence, manual review required",
			"doc_type": doc.DocType,
		}
	}

	// Determine new status based on confidence threshold.
	newStatus := "pending"
	if confidence >= ocrConfidenceThreshold {
		newStatus = "ocr_complete"
	}

	// Persist OCR result.
	updated, err := svc.store.UpdateDocumentOCR(ctx, sqlc.UpdateDocumentOCRParams{
		ID:            params.DocumentID,
		OcrResult:     ocrResult,
		OcrConfidence: confidence,
		Status:        newStatus,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: update document OCR: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("document_id", updated.ID).
		Str("doc_type", updated.DocType).
		Float64("confidence", confidence).
		Str("new_status", newStatus).
		Msg("OCR stub completed")

	span.SetAttributes(
		attribute.Float64("ocr.confidence", confidence),
		attribute.String("ocr.status", newStatus),
	)
	span.SetStatus(otelCodes.Ok, newStatus)

	return &TriggerOCRResult{
		DocumentID: updated.ID,
		Status:     newStatus,
		Confidence: confidence,
		OcrResult:  ocrResult,
	}, nil
}

// ---------------------------------------------------------------------------
// GetOCRStatus
// ---------------------------------------------------------------------------

// GetOCRStatusParams holds the inputs for retrieving OCR status.
type GetOCRStatusParams struct {
	DocumentID string
}

// GetOCRStatusResult holds the OCR state of a document.
type GetOCRStatusResult struct {
	DocumentID string
	Status     string
	Confidence float64
	OcrResult  map[string]string
}

// GetOCRStatus returns the current OCR state of a document including
// the extracted fields and confidence score, if OCR has been run.
func (svc *Service) GetOCRStatus(ctx context.Context, params *GetOCRStatusParams) (*GetOCRStatusResult, error) {
	const op = "service.Service.GetOCRStatus"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("document_id", params.DocumentID),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.DocumentID == "" {
		return nil, fmt.Errorf("%s: document_id is required", op)
	}

	doc, err := svc.store.GetDocumentWithOCR(ctx, params.DocumentID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get document: %w", op, err)
	}

	// Parse OCR result JSON from the stored text, if present.
	var ocrResult map[string]string
	var confidence float64

	if doc.OcrResult.Valid && doc.OcrResult.String != "" {
		parseErr := json.Unmarshal([]byte(doc.OcrResult.String), &ocrResult)
		if parseErr != nil {
			logger.Warn().
				Str("op", op).
				Str("document_id", doc.ID).
				Err(parseErr).
				Msg("failed to parse ocr_result JSON")
		}
	}

	if doc.OcrConfidence.Valid {
		f, _ := doc.OcrConfidence.Float64Value()
		confidence = f.Float64
	}

	logger.Info().
		Str("op", op).
		Str("document_id", doc.ID).
		Str("status", doc.Status).
		Float64("confidence", confidence).
		Msg("OCR status retrieved")

	span.SetStatus(otelCodes.Ok, doc.Status)
	return &GetOCRStatusResult{
		DocumentID: doc.ID,
		Status:     doc.Status,
		Confidence: confidence,
		OcrResult:  ocrResult,
	}, nil
}
