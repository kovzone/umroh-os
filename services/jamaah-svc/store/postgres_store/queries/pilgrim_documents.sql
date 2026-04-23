-- pilgrim_documents.sql — queries for jamaah.pilgrim_documents table.
-- S3-E-02 scaffold / BL-DOC-001..003.

-- name: InsertPilgrimDocument :one
INSERT INTO jamaah.pilgrim_documents (
    jamaah_id,
    booking_id,
    doc_type,
    file_path,
    status
) VALUES (
    $1, -- jamaah_id UUID
    $2, -- booking_id UUID
    $3, -- doc_type   TEXT
    $4, -- file_path  TEXT
    'pending'
)
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by;

-- name: GetPilgrimDocumentByID :one
SELECT id, jamaah_id, booking_id, doc_type, file_path, status,
       rejection_reason, uploaded_at, reviewed_at, reviewed_by
FROM jamaah.pilgrim_documents
WHERE id = $1;

-- name: ApprovePilgrimDocument :one
UPDATE jamaah.pilgrim_documents
SET
    status      = 'approved',
    reviewed_at = now(),
    reviewed_by = $2
WHERE id = $1
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by;

-- name: RejectPilgrimDocument :one
UPDATE jamaah.pilgrim_documents
SET
    status           = 'rejected',
    rejection_reason = $3,
    reviewed_at      = now(),
    reviewed_by      = $2
WHERE id = $1
RETURNING id, jamaah_id, booking_id, doc_type, file_path, status,
          rejection_reason, uploaded_at, reviewed_at, reviewed_by;
