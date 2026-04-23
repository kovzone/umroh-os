# UAT Config & Reference

## Environment

| Variable | Value |
|----------|-------|
| Gateway API | `http://216.176.238.161:4000` |
| Core Web (B2C) | `http://216.176.238.161` |
| PostgreSQL | `postgres://postgres:IDL4Ssfdo9ettSaFfleZp4M+3vKA8wX2@216.176.238.161:5432/umrohos?sslmode=disable` |
| Admin Email | `admin@umrohos.dev` |
| Admin Password | `password123` |

## Seed Data (dari migrations)

| ID | Keterangan |
|----|------------|
| `pkg_01JCDE00000000000000000001` | Package aktif (Umroh Reguler 12 hari) |
| `pkg_01JCDE00000000000000000002` | Package draft |
| `pkg_01JCDE00000000000000000003` | Package archived |
| `dep_01JCDF00000000000000000001` | Departure aktif (open, 40 seats) |
| `dep_01JCDF00000000000000000003` | Departure cancelled |
| `33333333-3333-3333-3333-333333333333` | Admin user ID |
| `11111111-1111-1111-1111-111111111111` | HQ branch ID |

## Scope UAT Per Slice

### ✅ S1 — Sudah Deployed (semua done)
- Auth: login, refresh, logout, /v1/me, suspend, permissions
- Catalog read: GET /v1/packages, /v1/packages/{id}, /v1/package-departures/{id}
- Catalog write: POST/PATCH/DELETE /v1/packages, POST/PATCH departures
- Booking: POST /v1/bookings (draft)
- Console UI: login, shell, catalog CRUD

### ✅ S2 — Sudah Deployed (BL-PAY-001..008 done, FE done)
- Invoice: POST /v1/invoices, GET /v1/invoices/{id}
- VA: POST /v1/invoices/{id}/virtual-accounts
- Webhook mock: POST /v1/webhooks/mock/trigger
- Checkout UI: /checkout/{booking_id}

### ⚠️ S3 — Partial (sebagian done)
- Done: Fulfillment tasks (BL-LOG-001), Finance journal (BL-FIN-001, 003, 004), Doc upload (BL-DOC-001)
- Done: Ops board UI (BL-FE-OPS-001)
- TODO: OCR, Manifest, Smart grouping, Shipment, Self-pickup → tandai NOT_DEPLOYED

### ✅ S4 — Sudah Deployed (BL-CRM-001..003, BL-FE-CRM-001 done)
- Lead: POST /v1/leads, GET /v1/leads, PUT /v1/leads/{id}, POST /v1/leads/{id}/convert
- CRM console: /console/leads

## Route Map

### Public (no auth)
```
GET  /v1/packages
GET  /v1/packages/{id}
GET  /v1/package-departures/{id}
POST /v1/bookings          (channel: b2c_self, b2b_agent)
POST /v1/sessions          (login)
POST /v1/sessions/refresh
POST /v1/leads             (public lead capture)
POST /v1/webhooks/midtrans
POST /v1/webhooks/xendit
POST /v1/webhooks/mock/trigger   (dev only)
```

### Staff (requires Bearer)
```
DELETE /v1/sessions        (logout)
GET    /v1/me
POST   /v1/packages
PATCH  /v1/packages/{id}
DELETE /v1/packages/{id}
POST   /v1/packages/{id}/departures
PATCH  /v1/package-departures/{id}
POST   /v1/bookings        (channel: cs)
POST   /v1/invoices
GET    /v1/invoices/{id}
POST   /v1/invoices/{id}/virtual-accounts
GET    /v1/leads
PUT    /v1/leads/{id}
POST   /v1/leads/{id}/convert
```

## Database Tables Relevant untuk UAT

```sql
-- Auth
iam.users, iam.sessions, iam.audit_logs

-- Catalog
catalog.packages, catalog.departures, catalog.pricing_rows

-- Booking
booking.bookings, booking.jamaah, booking.pilgrim_documents

-- Payment
payment.invoices, payment.virtual_accounts, payment.payment_events, payment.refunds

-- Finance
finance.journal_entries, finance.journal_lines

-- Logistics
logistics.fulfillment_tasks

-- CRM
crm.leads, crm.lead_status_history, crm_booking_links
```

## Test Data Tagging

Prefix wajib untuk semua UAT test data:
- Nama: `[UAT] ...`
- Email: `uat.{role}.{timestamp}@umrohos.dev`
- Notes/source: mengandung kata `uat`
