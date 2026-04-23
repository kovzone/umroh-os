ALTER TABLE jamaah.pilgrim_documents
    ADD COLUMN IF NOT EXISTS ocr_result      JSONB,
    ADD COLUMN IF NOT EXISTS ocr_confidence  NUMERIC(4,3),
    ADD COLUMN IF NOT EXISTS ocr_requested_at TIMESTAMPTZ;

-- Update status check to include ocr_complete
ALTER TABLE jamaah.pilgrim_documents
    DROP CONSTRAINT IF EXISTS pilgrim_documents_status_check;

ALTER TABLE jamaah.pilgrim_documents
    ADD CONSTRAINT pilgrim_documents_status_check
    CHECK (status IN ('pending', 'ocr_complete', 'approved', 'rejected'));
