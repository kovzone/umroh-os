// ocr.sql.go — hand-written sqlc-style query implementations for OCR columns
// on jamaah.pilgrim_documents.
//
// BL-DOC-002 / S3-E-02 — Passport OCR stub.
//
// These queries write/read the ocr_result, ocr_confidence, and ocr_requested_at
// columns added by migration 000020_add_pilgrim_document_ocr.

package sqlc

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Extended row type
// ---------------------------------------------------------------------------

// PilgrimDocumentOCRRow extends PilgrimDocumentRow with OCR columns.
type PilgrimDocumentOCRRow struct {
	PilgrimDocumentRow
	OcrResult       pgtype.Text        `json:"ocr_result"`        // JSONB stored as text
	OcrConfidence   pgtype.Numeric     `json:"ocr_confidence"`
	OcrRequestedAt  pgtype.Timestamptz `json:"ocr_requested_at"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// UpdateDocumentOCRParams holds inputs for UpdateDocumentOCR.
type UpdateDocumentOCRParams struct {
	ID            string
	OcrResult     map[string]string // serialised to JSONB
	OcrConfidence float64
	Status        string // "ocr_complete" or "pending"
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const updateDocumentOCR = `-- name: UpdateDocumentOCR :one
UPDATE jamaah.pilgrim_documents
SET ocr_result       = $1::jsonb,
    ocr_confidence   = $2,
    ocr_requested_at = now(),
    status           = $3
WHERE id = $4
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by,
          ocr_result, ocr_confidence, ocr_requested_at`

// UpdateDocumentOCR writes OCR results to a pilgrim document and updates its status.
func (q *Queries) UpdateDocumentOCR(ctx context.Context, arg UpdateDocumentOCRParams) (PilgrimDocumentOCRRow, error) {
	ocrJSON, err := json.Marshal(arg.OcrResult)
	if err != nil {
		return PilgrimDocumentOCRRow{}, err
	}

	row := q.db.QueryRow(ctx, updateDocumentOCR,
		string(ocrJSON),
		arg.OcrConfidence,
		arg.Status,
		arg.ID,
	)

	var r PilgrimDocumentOCRRow
	err = row.Scan(
		&r.ID, &r.JamaahID, &r.BookingID, &r.DocType, &r.FilePath, &r.Status,
		&r.RejectionReason, &r.UploadedAt, &r.ReviewedAt, &r.ReviewedBy,
		&r.OcrResult, &r.OcrConfidence, &r.OcrRequestedAt,
	)
	return r, err
}

const getDocumentWithOCR = `-- name: GetDocumentWithOCR :one
SELECT id, jamaah_id, booking_id, doc_type, file_path, status,
       rejection_reason, uploaded_at, reviewed_at, reviewed_by,
       ocr_result, ocr_confidence, ocr_requested_at
FROM jamaah.pilgrim_documents
WHERE id = $1`

// GetDocumentWithOCR fetches a document by ID including OCR columns.
func (q *Queries) GetDocumentWithOCR(ctx context.Context, id string) (PilgrimDocumentOCRRow, error) {
	row := q.db.QueryRow(ctx, getDocumentWithOCR, id)
	var r PilgrimDocumentOCRRow
	err := row.Scan(
		&r.ID, &r.JamaahID, &r.BookingID, &r.DocType, &r.FilePath, &r.Status,
		&r.RejectionReason, &r.UploadedAt, &r.ReviewedAt, &r.ReviewedBy,
		&r.OcrResult, &r.OcrConfidence, &r.OcrRequestedAt,
	)
	return r, err
}
