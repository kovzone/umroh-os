# UAT Agent 3 — CRM/Lead (S4) + Ops/Finance (S3) — 2026-04-24

## Summary

- **Total test cases**: 20 (dalam spec file `06-uat-s3s4-crm-ops.spec.ts`)
- **PASS**: 0 (network unreachable — tidak ada live test berhasil)
- **FAIL (code analysis)**: 5
- **NOT_DEPLOYED**: 3
- **SKIP**: 12 (all dependent on auth/network)
- **Issues dibuat**: ISSUE-012 .. ISSUE-019

### Metode Testing

Server `http://216.176.238.161` tidak dapat dijangkau dari session ini (ENETUNREACH pada TCP level — sama seperti ISSUE-001 yang ditemukan Agent 1). Playwright test gagal semua karena koneksi diblokir.

Sebagai alternatif, dilakukan **static code analysis** terhadap:
- Services: `crm-svc`, `finance-svc`, `logistics-svc`, `gateway-svc`
- Migration: `000013`, `000014`
- Test spec: `06-uat-s3s4-crm-ops.spec.ts`
- Test helpers: `tests/e2e/lib/uat-helpers.ts`
- Contracts: `docs/contracts/slice-S3.md`, `docs/contracts/slice-S4.md`

---

## Hasil Per Test Case

| # | Skenario | Status | Notes |
|---|----------|--------|-------|
| S4-LEAD-01 | POST /v1/leads (public) → 201, lead tersimpan | ⚠️ NOT_DEPLOYED | Route ada (`POST /v1/leads`) — kemungkinan PASS jika server online, tapi loginAdmin() (beforeAll) menggunakan /v1/sessions yang 404 → beforeAll fail |
| S4-LEAD-02 | POST /v1/leads dengan UTM params → UTM tersimpan di DB | ⚠️ SKIP | Bergantung pada S4-LEAD-01; DB unreachable |
| S4-LEAD-03 | POST /v1/leads phone invalid → 422 | ❌ FAIL | Phone validation error return HTTP 500 bukan 400/422 — ISSUE-012 |
| S4-LEAD-04 | POST /v1/leads tanpa UTM → 201, UTM fields null | ✅ LIKELY PASS | Logic OK: UTM fields nullable di schema dan `textOrNull("")` = NULL |
| S4-LEAD-05 | GET /v1/leads (admin) → status 'new' | ⚠️ SKIP | Bergantung pada createdLeadId dari S4-LEAD-01 |
| S4-LEAD-06 | CS round-robin — 3 lead baru → distribute ke CS berbeda | ⚠️ SKIP | Bergantung pada S4-LEAD-01; round-robin logic ada (GetLeastLoadedCS) |
| S4-LEAD-07 | PUT /v1/leads/{id} → update status ke 'contacted' | ⚠️ SKIP | Bergantung pada createdLeadId |
| S4-LEAD-08 | Transisi tidak valid new → converted → 400/422 | ❌ FAIL | `validTransition()` mengizinkan new→converted — ISSUE-013 |
| S4-LEAD-09 | Status 'lost' adalah terminal | ✅ LIKELY PASS | Code benar: `if from == "lost" { return false }` |
| S4-LEAD-10 | GET /v1/leads?status=new → filter leads | ✅ LIKELY PASS | ListLeads dengan StatusFilter ada di code |
| S4-LINK-01 | Booking dengan email sama → lead status update otomatis | ⚠️ SKIP | Bergantung pada S2 booking + S4 lead; network unreachable |
| S3-OPS-01 | Setelah paid, fulfillment task muncul otomatis | ❌ FAIL | `/v1/invoices` NOT_DEPLOYED + expected status 'pending' tapi actual 'queued' — ISSUE-014 |
| S3-OPS-02 | GET /console/ops → ops board tampil | ⚠️ SKIP | Core web unreachable (ISSUE-002) |
| S3-FIN-01 | Journal entry terbuat otomatis setelah payment | ❌ FAIL | `/v1/invoices` NOT_DEPLOYED + query column salah (debit_amount vs debit) — ISSUE-015 |
| S3-FIN-02 | Journal idempoten — tidak ada duplikat | 🔍 PARTIAL | DB query benar (idempotency_key check), tapi DB unreachable untuk verify |
| S4-UI-01 | /contact form tampil | ⚠️ SKIP | Core web unreachable (ISSUE-002) |
| S4-UI-02 | Isi form lead → submit → pesan sukses | ⚠️ SKIP | Core web unreachable |
| S4-UI-03 | Lead muncul di /console/leads | ⚠️ SKIP | Core web unreachable |
| S4-UI-04 | UTM params dari URL tersimpan | ⚠️ SKIP | Core web unreachable |
| S5-UI-01 | /console/finance → finance view tampil | ⚠️ SKIP | Core web unreachable |

---

## Output Playwright (uat-api project)

```
Running 20 tests using 1 worker

  ✘   1 [uat-api] › tests/06-uat-s3s4-crm-ops.spec.ts:47:7 › S4 CRM — Lead Capture (BL-CRM-001) › S4-LEAD-01 (0ms)
  -   2 .. 20  [uat-api] — did not run (beforeAll failed)

Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
  POST http://216.176.238.161:4000/v1/sessions (loginAdmin → beforeAll)
  
Error: connect ENETUNREACH 216.176.238.161:5432 - Local (0.0.0.0:0)
  (DB connection juga blocked)

1 failed — 19 did not run
```

Root cause: QA session ini tidak memiliki akses network ke 216.176.238.161. Ini adalah pembatasan network di environment session, bukan masalah server (lihat ISSUE-001 untuk status server).

---

## DB Integrity Findings (Static Analysis)

### Schema Verification (dari migration files)

#### `crm.leads` — migration 000014
- UTM columns: `utm_source`, `utm_medium`, `utm_campaign`, `utm_content`, `utm_term` — **TERSEDIA** ✅
- Status lifecycle: `'new','contacted','qualified','converted','lost'` — **TERSEDIA** ✅
- `assigned_cs_id` untuk round-robin: **TERSEDIA** ✅
- Index pada `status`, `assigned_cs_id`, `created_at`: **TERSEDIA** ✅

#### `logistics.fulfillment_tasks` — migration 000013
- Status default: `'queued'` — **MISMATCH dengan spec yang expect 'pending'** ❌ → ISSUE-014
- Valid statuses: `queued, processing, shipped, delivered, cancelled` (tidak ada 'pending') ❌

#### `finance.journal_entries` — migration 000013
- `idempotency_key` UNIQUE: **TERSEDIA** ✅ — mencegah duplicate journal
- `source_id UUID` (bukan `source_ref`): cleanup di uat-helpers.ts menggunakan `source_ref` → BUG ❌ → ISSUE-016

#### `finance.journal_lines` — migration 000013
- Columns: `debit NUMERIC(15,2)`, `credit NUMERIC(15,2)` — **bukan debit_amount/credit_amount** ❌ → ISSUE-015
- CHECK constraint: `(debit > 0 AND credit = 0) OR (debit = 0 AND credit > 0)` — **double-entry enforced** ✅
- Balance check juga enforced di service layer (`service/journal.go`) **sebelum INSERT** ✅

### Service Logic Verification

#### CRM Phone Validation (crm-svc/service/leads.go:169)
```go
err := fmt.Errorf("%s: phone must be at least 8 characters", op)
// MASALAH: tidak wrap dengan ErrValidation
// → gRPC returns codes.Internal → gateway returns HTTP 500
```
**BUG** ❌ → ISSUE-012

#### Status Transition Logic (crm-svc/service/leads.go:134–151)
```go
func validTransition(from, to string) bool {
    if from == to { return true }
    if from == "converted" || from == "lost" { return false }
    valid := map[string]bool{"new":true,"contacted":true,"qualified":true,"converted":true,"lost":true}
    return valid[to]  // MASALAH: allows any-to-any skip (e.g. new→converted)
}
```
**BUG** ❌ → ISSUE-013

#### Finance Journal Double-Entry (finance-svc/service/journal.go)
- `OnPaymentReceived` posts Dr 1001 (Bank) / Cr 2001 (Pilgrim Liability) ✅
- Idempotency key: `"payment:" + invoiceID` ✅
- Balance check before INSERT: `if debitTotal != creditTotal { return error }` ✅
- Amounts: stored as `int64 * 100` with `Exp=-2` in NUMERIC(15,2) ✅
- Transaction: INSERT entry + 2 lines in single TX ✅

#### Finance Duplicate Check (S3-FIN-02)
The spec S3-FIN-02 query correctly checks:
```sql
SELECT COUNT(*) FROM finance.journal_entries
WHERE idempotency_key ILIKE 'payment:%'
GROUP BY idempotency_key
HAVING COUNT(*) > 1
```
This is correct and would return 0 rows if idempotency works — **logically correct** ✅

---

## Temuan Penting

### 1. CRITICAL: Network Unreachable (ISSUE-001, 002, 003) — Pre-existing
Server `216.176.238.161` tidak dapat dijangkau dari QA session. Semua live tests blocked.

### 2. HIGH: Phone Validation Returns HTTP 500 (ISSUE-012)
Service layer `crm-svc/service/leads.go:169` menggunakan bare `fmt.Errorf` untuk phone validation failure. Fix mudah: gunakan `errors.Join(apperrors.ErrValidation, ...)`.

### 3. HIGH: Lead Status Transition Bypass (ISSUE-013)
`validTransition()` di `crm-svc/service/leads.go` mengizinkan `new → converted` (seharusnya hanya `new → contacted → qualified → converted`). **Security/integrity risk**: CS bisa langsung mark lead sebagai converted tanpa proses nurturing.

### 4. HIGH: UAT Cleanup Akan Fail (ISSUE-016)
`uat-helpers.ts` cleanup SQL menggunakan `source_ref` yang tidak ada di `finance.journal_entries` (kolom sebenarnya `source_id`). Finance data UAT akan tertinggal setelah test run.

### 5. MEDIUM: Fulfillment Task Status Mismatch (ISSUE-014)
DB schema defaultnya `'queued'` tapi spec expects `'pending'`. Perlu sinkronisasi antara spec dan schema.

### 6. MEDIUM: Finance Journal Lines Query Column Salah (ISSUE-015)
Spec query `debit_amount`/`credit_amount` tidak ada di DB (`debit`/`credit`). Akan runtime error saat test dijalankan.

### 7. NOT_DEPLOYED: Endpoints Hilang
- `POST /v1/invoices` — tidak ada di gateway (ISSUE-005 dari Agent 2)
- `POST /v1/webhooks/mock/trigger` — tidak ada di gateway (ISSUE-007 dari Agent 2)
- `GET /v1/fulfillment-tasks` — tidak ada di gateway (ISSUE-018)

---

## Test Data Dibuat

Tidak ada — semua tes gagal karena network unreachable. Cleanup tidak diperlukan.

---

## Issues Baru (Agent 3)

| Issue | Severity | Status | Ringkasan |
|-------|----------|--------|-----------|
| ISSUE-012 | HIGH | OPEN | Phone validation return 500 bukan 400/422 — crm-svc service/leads.go:169 |
| ISSUE-013 | HIGH | OPEN | new→converted transition diizinkan (seharusnya blocked) — validTransition() bug |
| ISSUE-014 | MEDIUM | OPEN | Fulfillment task default status 'queued' vs spec expect 'pending' |
| ISSUE-015 | MEDIUM | OPEN | journal_lines query column salah: debit_amount → debit, credit_amount → credit |
| ISSUE-016 | HIGH | OPEN | uat-helpers.ts cleanup pakai source_ref (tidak ada), seharusnya source_id |
| ISSUE-017 | LOW | OPEN | Gateway CreateLeadBody tidak accept field 'message'/'source_note' dari spec |
| ISSUE-018 | MEDIUM | NOT_DEPLOYED | GET /v1/fulfillment-tasks belum terdaftar di gateway |
| ISSUE-019 | LOW | OPEN | Curl instruction pakai /v1/finance/journal-entries, route sebenarnya /v1/finance/journals |
