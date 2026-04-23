// pilgrim_documents.sql.go — hand-written sqlc-style query implementations for
// jamaah-svc pilgrim document queries.
//
// Run `make generate` (sqlc generate) to regenerate from pilgrim_documents.sql
// once sqlc is configured to target the jamaah schema.
//
// S3-E-02 scaffold / BL-DOC-001..003.

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row type
// ---------------------------------------------------------------------------

// PilgrimDocumentRow mirrors a row from jamaah.pilgrim_documents.
type PilgrimDocumentRow struct {
	ID              string             `json:"id"`
	JamaahID        string             `json:"jamaah_id"`
	BookingID       string             `json:"booking_id"`
	DocType         string             `json:"doc_type"`
	FilePath        string             `json:"file_path"`
	Status          string             `json:"status"`
	RejectionReason pgtype.Text        `json:"rejection_reason"`
	UploadedAt      pgtype.Timestamptz `json:"uploaded_at"`
	ReviewedAt      pgtype.Timestamptz `json:"reviewed_at"`
	ReviewedBy      pgtype.Text        `json:"reviewed_by"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// InsertPilgrimDocumentParams holds inputs for InsertPilgrimDocument.
type InsertPilgrimDocumentParams struct {
	JamaahID  string
	BookingID string
	DocType   string
	FilePath  string
}

// ApprovePilgrimDocumentParams holds inputs for ApprovePilgrimDocument.
type ApprovePilgrimDocumentParams struct {
	ID         string
	ReviewedBy string
}

// RejectPilgrimDocumentParams holds inputs for RejectPilgrimDocument.
type RejectPilgrimDocumentParams struct {
	ID              string
	ReviewedBy      string
	RejectionReason string
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const insertPilgrimDocument = `-- name: InsertPilgrimDocument :one
INSERT INTO jamaah.pilgrim_documents (
    jamaah_id, booking_id, doc_type, file_path, status
) VALUES ($1, $2, $3, $4, 'pending')
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by`

// InsertPilgrimDocument inserts a new document record in status='pending'.
func (q *Queries) InsertPilgrimDocument(ctx context.Context, arg InsertPilgrimDocumentParams) (PilgrimDocumentRow, error) {
	row := q.db.QueryRow(ctx, insertPilgrimDocument,
		arg.JamaahID, arg.BookingID, arg.DocType, arg.FilePath,
	)
	var r PilgrimDocumentRow
	err := row.Scan(
		&r.ID, &r.JamaahID, &r.BookingID, &r.DocType, &r.FilePath, &r.Status,
		&r.RejectionReason, &r.UploadedAt, &r.ReviewedAt, &r.ReviewedBy,
	)
	return r, err
}

const getPilgrimDocumentByID = `-- name: GetPilgrimDocumentByID :one
SELECT id, jamaah_id, booking_id, doc_type, file_path, status,
       rejection_reason, uploaded_at, reviewed_at, reviewed_by
FROM jamaah.pilgrim_documents WHERE id = $1`

// GetPilgrimDocumentByID fetches a document by its ID.
func (q *Queries) GetPilgrimDocumentByID(ctx context.Context, id string) (PilgrimDocumentRow, error) {
	row := q.db.QueryRow(ctx, getPilgrimDocumentByID, id)
	var r PilgrimDocumentRow
	err := row.Scan(
		&r.ID, &r.JamaahID, &r.BookingID, &r.DocType, &r.FilePath, &r.Status,
		&r.RejectionReason, &r.UploadedAt, &r.ReviewedAt, &r.ReviewedBy,
	)
	return r, err
}

const approvePilgrimDocument = `-- name: ApprovePilgrimDocument :one
UPDATE jamaah.pilgrim_documents
SET status = 'approved', reviewed_at = now(), reviewed_by = $2
WHERE id = $1
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by`

// ApprovePilgrimDocument sets a document to approved and records the reviewer.
func (q *Queries) ApprovePilgrimDocument(ctx context.Context, arg ApprovePilgrimDocumentParams) (PilgrimDocumentRow, error) {
	row := q.db.QueryRow(ctx, approvePilgrimDocument, arg.ID, arg.ReviewedBy)
	var r PilgrimDocumentRow
	err := row.Scan(
		&r.ID, &r.JamaahID, &r.BookingID, &r.DocType, &r.FilePath, &r.Status,
		&r.RejectionReason, &r.UploadedAt, &r.ReviewedAt, &r.ReviewedBy,
	)
	return r, err
}

const rejectPilgrimDocument = `-- name: RejectPilgrimDocument :one
UPDATE jamaah.pilgrim_documents
SET status = 'rejected', rejection_reason = $3, reviewed_at = now(), reviewed_by = $2
WHERE id = $1
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by`

// RejectPilgrimDocument sets a document to rejected with a reason.
func (q *Queries) RejectPilgrimDocument(ctx context.Context, arg RejectPilgrimDocumentParams) (PilgrimDocumentRow, error) {
	row := q.db.QueryRow(ctx, rejectPilgrimDocument, arg.ID, arg.ReviewedBy, arg.RejectionReason)
	var r PilgrimDocumentRow
	err := row.Scan(
		&r.ID, &r.JamaahID, &r.BookingID, &r.DocType, &r.FilePath, &r.Status,
		&r.RejectionReason, &r.UploadedAt, &r.ReviewedAt, &r.ReviewedBy,
	)
	return r, err
}
