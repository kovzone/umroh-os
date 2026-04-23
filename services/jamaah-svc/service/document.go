// document.go — UploadDocument and ReviewDocument service-layer implementations.
//
// S3-E-02 scaffold / F3-W1 / BL-DOC-001..003.
//
// UploadDocument — stores document metadata in jamaah.pilgrim_documents with
// status='pending'. File storage (GCS) is handled by the REST layer or a
// separate worker; this service layer only persists the record.
//
// ReviewDocument — approve or reject a document. Every action writes a
// structured audit log entry.
//
// Audit: every state-changing operation is logged via zerolog. Full IAM
// audit-log write (iam-svc.RecordAudit) is a TODO once the iam adapter is
// wired to jamaah-svc.

package service

import (
	"context"
	"fmt"

	"jamaah-svc/store/postgres_store/sqlc"
	"jamaah-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// UploadDocument
// ---------------------------------------------------------------------------

// UploadDocumentParams holds the inputs for uploading a document.
type UploadDocumentParams struct {
	JamaahID   string
	BookingID  string
	DocType    string // ktp | passport | photo | other
	FilePath   string // GCS path
	UploadedBy string // staff user ID if uploaded on behalf; empty for self-upload
}

// UploadDocumentResult holds the result of an UploadDocument call.
type UploadDocumentResult struct {
	DocumentID string
	Status     string
}

// UploadDocument stores a document metadata record in status='pending'.
func (svc *Service) UploadDocument(ctx context.Context, params *UploadDocumentParams) (*UploadDocumentResult, error) {
	const op = "service.Service.UploadDocument"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("jamaah_id", params.JamaahID),
		attribute.String("booking_id", params.BookingID),
		attribute.String("doc_type", params.DocType),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.JamaahID == "" {
		return nil, fmt.Errorf("%s: jamaah_id is required", op)
	}
	if params.BookingID == "" {
		return nil, fmt.Errorf("%s: booking_id is required", op)
	}
	if params.FilePath == "" {
		return nil, fmt.Errorf("%s: file_path is required", op)
	}
	if !isValidDocType(params.DocType) {
		return nil, fmt.Errorf("%s: invalid doc_type %q (must be ktp|passport|photo|other)", op, params.DocType)
	}

	doc, err := svc.store.InsertPilgrimDocument(ctx, sqlc.InsertPilgrimDocumentParams{
		JamaahID:  params.JamaahID,
		BookingID: params.BookingID,
		DocType:   params.DocType,
		FilePath:  params.FilePath,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: insert document: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("jamaah_id", params.JamaahID).
		Str("booking_id", params.BookingID).
		Str("doc_type", params.DocType).
		Str("document_id", doc.ID).
		Str("uploaded_by", params.UploadedBy).
		Msg("document uploaded")

	span.SetStatus(otelCodes.Ok, "uploaded")
	return &UploadDocumentResult{
		DocumentID: doc.ID,
		Status:     doc.Status,
	}, nil
}

// ---------------------------------------------------------------------------
// ReviewDocument
// ---------------------------------------------------------------------------

// ReviewDocumentParams holds the inputs for reviewing (approving/rejecting) a document.
type ReviewDocumentParams struct {
	DocumentID      string
	Action          string // approve | reject
	RejectionReason string // required when action = reject
	ReviewedBy      string // staff user ID (required)
}

// ReviewDocumentResult holds the result of a ReviewDocument call.
type ReviewDocumentResult struct {
	DocumentID string
	Status     string
}

// ReviewDocument approves or rejects a document.
// Every action writes a structured audit log.
func (svc *Service) ReviewDocument(ctx context.Context, params *ReviewDocumentParams) (*ReviewDocumentResult, error) {
	const op = "service.Service.ReviewDocument"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("document_id", params.DocumentID),
		attribute.String("action", params.Action),
		attribute.String("reviewed_by", params.ReviewedBy),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.DocumentID == "" {
		return nil, fmt.Errorf("%s: document_id is required", op)
	}
	if params.ReviewedBy == "" {
		return nil, fmt.Errorf("%s: reviewed_by is required", op)
	}
	if params.Action != "approve" && params.Action != "reject" {
		return nil, fmt.Errorf("%s: action must be 'approve' or 'reject', got %q", op, params.Action)
	}
	if params.Action == "reject" && params.RejectionReason == "" {
		return nil, fmt.Errorf("%s: rejection_reason is required when action=reject", op)
	}

	var docID, docStatus string
	var err error

	switch params.Action {
	case "approve":
		var doc sqlc.PilgrimDocumentRow
		doc, err = svc.store.ApprovePilgrimDocument(ctx, sqlc.ApprovePilgrimDocumentParams{
			ID:         params.DocumentID,
			ReviewedBy: params.ReviewedBy,
		})
		if err == nil {
			docID, docStatus = doc.ID, doc.Status
		}

	case "reject":
		var doc sqlc.PilgrimDocumentRow
		doc, err = svc.store.RejectPilgrimDocument(ctx, sqlc.RejectPilgrimDocumentParams{
			ID:              params.DocumentID,
			ReviewedBy:      params.ReviewedBy,
			RejectionReason: params.RejectionReason,
		})
		if err == nil {
			docID, docStatus = doc.ID, doc.Status
		}
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: review document (%s): %w", op, params.Action, err)
	}

	// Audit log
	logger.Info().
		Str("op", op).
		Str("document_id", docID).
		Str("action", params.Action).
		Str("reviewed_by", params.ReviewedBy).
		Str("new_status", docStatus).
		Str("rejection_reason", params.RejectionReason).
		Msg("document reviewed")

	// TODO: call iam-svc.RecordAudit once iam adapter is wired to jamaah-svc.

	span.SetStatus(otelCodes.Ok, params.Action)
	return &ReviewDocumentResult{
		DocumentID: docID,
		Status:     docStatus,
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func isValidDocType(t string) bool {
	switch t {
	case "ktp", "passport", "photo", "other":
		return true
	}
	return false
}
