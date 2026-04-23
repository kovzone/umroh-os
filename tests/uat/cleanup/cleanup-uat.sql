-- UAT Cleanup Script
-- Menghapus semua data test yang dibuat dengan prefix [UAT]
-- Jalankan via: bash tests/uat/cleanup/cleanup-uat.sh
--
-- AMAN: hanya menyentuh rows yang namanya/emailnya mengandung prefix UAT.
-- Tidak akan menyentuh data real production.

BEGIN;

-- ============================================================
-- CRM: leads dari UAT
-- ============================================================
DELETE FROM crm.lead_status_history
WHERE lead_id IN (
  SELECT id FROM crm.leads
  WHERE email ILIKE 'uat.%@%'
     OR source_note ILIKE '%uat%'
     OR name ILIKE '[UAT]%'
);

DELETE FROM crm.leads
WHERE email ILIKE 'uat.%@%'
   OR source_note ILIKE '%uat%'
   OR name ILIKE '[UAT]%';

-- ============================================================
-- Payment: events, VAs, invoices terkait booking UAT
-- ============================================================
DELETE FROM payment.payment_events
WHERE invoice_id IN (
  SELECT pi.id FROM payment.invoices pi
  JOIN booking.bookings b ON b.id = pi.booking_id
  WHERE b.notes ILIKE '%[UAT]%'
     OR b.notes ILIKE '%uat%'
);

DELETE FROM payment.virtual_accounts
WHERE invoice_id IN (
  SELECT pi.id FROM payment.invoices pi
  JOIN booking.bookings b ON b.id = pi.booking_id
  WHERE b.notes ILIKE '%[UAT]%'
     OR b.notes ILIKE '%uat%'
);

DELETE FROM payment.invoices
WHERE booking_id IN (
  SELECT id FROM booking.bookings
  WHERE notes ILIKE '%[UAT]%'
     OR notes ILIKE '%uat%'
);

-- ============================================================
-- Finance: journal entries/lines terkait booking UAT
-- (idempotency key = "payment:<invoice_id>")
-- ============================================================
DELETE FROM finance.journal_lines
WHERE entry_id IN (
  SELECT je.id FROM finance.journal_entries je
  WHERE je.idempotency_key ILIKE 'payment:%'
    AND je.source_ref IN (
      SELECT pi.id::text FROM payment.invoices pi
      JOIN booking.bookings b ON b.id = pi.booking_id
      WHERE b.notes ILIKE '%[UAT]%'
    )
);

DELETE FROM finance.journal_entries
WHERE idempotency_key ILIKE 'uat.%'
   OR (source_ref IN (
      SELECT pi.id::text FROM payment.invoices pi
      JOIN booking.bookings b ON b.id = pi.booking_id
      WHERE b.notes ILIKE '%[UAT]%'
   ));

-- ============================================================
-- Logistics: fulfillment tasks terkait booking UAT
-- ============================================================
DELETE FROM logistics.fulfillment_tasks
WHERE booking_id IN (
  SELECT id FROM booking.bookings
  WHERE notes ILIKE '%[UAT]%'
     OR notes ILIKE '%uat%'
);

-- ============================================================
-- Booking: pilgrim docs, jamaah, bookings UAT
-- ============================================================
DELETE FROM booking.pilgrim_documents
WHERE jamaah_id IN (
  SELECT bj.id FROM booking.jamaah bj
  JOIN booking.bookings b ON b.id = bj.booking_id
  WHERE b.notes ILIKE '%[UAT]%'
     OR b.notes ILIKE '%uat%'
);

DELETE FROM booking.jamaah
WHERE booking_id IN (
  SELECT id FROM booking.bookings
  WHERE notes ILIKE '%[UAT]%'
     OR notes ILIKE '%uat%'
);

DELETE FROM booking.bookings
WHERE notes ILIKE '%[UAT]%'
   OR notes ILIKE '%uat%';

-- ============================================================
-- Catalog: departures dan packages UAT
-- ============================================================
DELETE FROM catalog.departures
WHERE package_id IN (
  SELECT id FROM catalog.packages
  WHERE name ILIKE '[UAT]%'
);

DELETE FROM catalog.packages
WHERE name ILIKE '[UAT]%';

-- ============================================================
-- IAM: audit logs dari UAT (opsional, bisa di-skip)
-- ============================================================
-- DELETE FROM iam.audit_logs
-- WHERE actor_id = '33333333-3333-3333-3333-333333333333'
--   AND created_at > NOW() - INTERVAL '1 day'
--   AND resource_type IN ('package', 'departure', 'booking', 'lead');

COMMIT;

-- Verifikasi cleanup
SELECT 'packages' AS tabel, COUNT(*) AS sisa_uat FROM catalog.packages WHERE name ILIKE '[UAT]%'
UNION ALL
SELECT 'bookings', COUNT(*) FROM booking.bookings WHERE notes ILIKE '%[UAT]%'
UNION ALL
SELECT 'leads', COUNT(*) FROM crm.leads WHERE email ILIKE 'uat.%@%'
UNION ALL
SELECT 'invoices', COUNT(*) FROM payment.invoices pi JOIN booking.bookings b ON b.id = pi.booking_id WHERE b.notes ILIKE '%[UAT]%';
