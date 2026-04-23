-- Revert status constraint to original
ALTER TABLE jamaah.pilgrim_documents
    DROP CONSTRAINT IF EXISTS pilgrim_documents_status_check;

ALTER TABLE jamaah.pilgrim_documents
    ADD CONSTRAINT pilgrim_documents_status_check
    CHECK (status IN ('pending', 'approved', 'rejected'));

-- Remove OCR columns
ALTER TABLE jamaah.pilgrim_documents
    DROP COLUMN IF EXISTS ocr_result,
    DROP COLUMN IF EXISTS ocr_confidence,
    DROP COLUMN IF EXISTS ocr_requested_at;
