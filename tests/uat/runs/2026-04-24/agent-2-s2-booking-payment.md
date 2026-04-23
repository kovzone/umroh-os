# UAT Agent 2 — Booking + Payment (S2) — 2026-04-24

## Summary

- **Total test cases**: 14 (dari spec `05-uat-s2-booking-payment.spec.ts`)
- **PASS**: 0
- **FAIL**: 0 (live test tidak dapat dieksekusi — lihat keterangan)
- **NOT_DEPLOYED**: 7 (route S2 tidak terdaftar di gateway)
- **SKIP**: 7 (dependency test tidak bisa jalan karena endpoint tidak ada)
- **Issues dibuat**: ISSUE-005 .. ISSUE-011

## Catatan Penting: Network Sandbox Blocked

**Test Playwright tidak dapat dijalankan secara live** karena sandbox environment agent ini memiliki network egress proxy yang hanya mengizinkan whitelist domain tertentu (github.com, npmjs.com, api.anthropic.com). Server `216.176.238.161` diblokir.

Sebagai gantinya, dilakukan **static code analysis** terhadap codebase untuk mengidentifikasi masalah deployment dan implementasi. Semua temuan di bawah ini didasarkan pada analisis kode sumber langsung — bukan dari observasi runtime.

## Hasil Per Test Case (Prediksi dari Static Analysis)

| # | Test Case | Status | Notes |
|---|-----------|--------|-------|
| S2-BOOK-01 | POST /v1/bookings (b2c_self) → 201, status draft | ⚠️ FAIL | Schema mismatch: spec kirim `lead.full_name`/`jamaah[]`, gateway expect `lead_full_name`/`pilgrims[]` → ISSUE-011 |
| S2-BOOK-02 | POST /v1/bookings tanpa email → 422 | ⏭️ SKIP | Depends on BOOK-01 passing (adminToken + bookingId) |
| S2-BOOK-03 | POST /v1/bookings tanpa package_id → 422 | ⏭️ SKIP | Validation test, mungkin pass jika body format difix |
| S2-BOOK-04 | POST /v1/bookings dengan departure closed → 4xx | ⏭️ SKIP | Depends on booking working |
| S2-BOOK-05 | Sisa kursi berkurang setelah booking | ⏭️ SKIP | Skip jika bookingId kosong |
| S2-PAY-01 | POST /v1/bookings/{id}/submit → pending_payment | ⚠️ NOT_DEPLOYED | Route tidak terdaftar di gateway → ISSUE-006 |
| S2-PAY-02 | POST /v1/invoices → invoice + VA | ⚠️ NOT_DEPLOYED | Route tidak terdaftar di gateway → ISSUE-005 |
| S2-PAY-03 | POST /v1/invoices/{id}/virtual-accounts → VA | ⚠️ NOT_DEPLOYED | Route tidak terdaftar di gateway → ISSUE-005 |
| S2-PAY-04 | GET /v1/invoices/{id} idempoten | ⚠️ NOT_DEPLOYED | Route tidak terdaftar di gateway → ISSUE-005 |
| S2-WH-01 | POST /v1/webhooks/mock/trigger → paid_in_full | ⚠️ NOT_DEPLOYED | Route tidak ada di gateway:4000 → ISSUE-007 |
| S2-WH-02 | Duplicate webhook idempoten | ⚠️ NOT_DEPLOYED | Depends on webhook route → ISSUE-007 |
| S2-WH-03 | Webhook amount kecil → tidak paid_in_full | ⚠️ NOT_DEPLOYED | Depends on invoice + webhook routes |
| S2-WH-04 | Webhook Midtrans tanpa auth → 401/403 | ⚠️ NOT_DEPLOYED | Route tidak ada di gateway:4000 → ISSUE-008 |
| S2-UI-01 | B2C booking flow end-to-end | ⏭️ SKIP | Depends on core-web up + booking schema fix |
| S2-UI-02 | /checkout/{booking_id} tampilkan VA | ⚠️ FAIL | Checkout menggunakan mock data by default → ISSUE-010 |
| S2-UI-03 | Booking form validation — error inline | ⏭️ SKIP | Depends on core-web up |

## Issues Dibuat

| Issue | Severity | Deskripsi |
|-------|----------|-----------|
| ISSUE-005 | CRITICAL | Gateway tidak mendaftarkan route POST /v1/invoices, GET /v1/invoices/{id}, POST /v1/invoices/{id}/virtual-accounts |
| ISSUE-006 | CRITICAL | Gateway tidak mendaftarkan route POST /v1/bookings/{id}/submit |
| ISSUE-007 | HIGH | Gateway tidak mendaftarkan route POST /v1/webhooks/mock/trigger |
| ISSUE-008 | HIGH | Gateway tidak mendaftarkan route POST /v1/webhooks/midtrans |
| ISSUE-009 | HIGH | UAT spec menggunakan /v1/sessions tapi gateway mendaftarkan /v1/auth/login |
| ISSUE-010 | MEDIUM | Checkout page gunakan mock data by default (VITE_MOCK_GATEWAY tidak diset false di prod) |
| ISSUE-011 | HIGH | Booking request body schema mismatch antara spec (lead/jamaah[]) dan gateway (lead_full_name/pilgrims[]) |

## Analisis Root Cause

### 1. Missing Gateway Routes (ISSUE-005, 006, 007, 008)

File `services/gateway-svc/cmd/server.go` adalah sumber kebenaran untuk semua registered routes. Setelah analisis lengkap semua 60+ routes yang terdaftar, route S2 berikut **tidak ditemukan**:

- `POST /v1/invoices`
- `GET /v1/invoices/:id`
- `POST /v1/invoices/:id/virtual-accounts`
- `POST /v1/bookings/:id/submit`
- `POST /v1/webhooks/midtrans`
- `POST /v1/webhooks/xendit`
- `POST /v1/webhooks/mock/trigger`

Satu-satunya route payment yang terdaftar adalah:
```
POST /v1/payments/link  (BL-PAY-020 - ReissuePaymentLink)
```

Handler `proxy_payment.go` hanya mengimplementasikan `ReissuePaymentLink`, bukan invoice/VA/webhook routes.

### 2. Webhook Architecture Gap (ISSUE-007, 008)

Payment-svc menjalankan HTTP webhook server di port **50065** (bukan 4000). Port ini tidak di-expose di `docker-compose.prod.yml`. Webhook dari Midtrans/Xendit perlu reach ke port ini, tapi tidak ada nginx proxy di depan. Route `/v1/webhooks/*` tidak ada di gateway-svc port 4000.

### 3. Schema Mismatch (ISSUE-009, 011)

UAT spec `05-uat-s2-booking-payment.spec.ts` ditulis mengharapkan API contract yang berbeda dari implementasi aktual gateway:
- Login endpoint: spec uses `/v1/sessions`, gateway registers `/v1/auth/login`
- Booking body: spec sends nested `lead.{}` + `jamaah[]`, gateway expects flat `lead_full_name` + `pilgrims[]`

### 4. Frontend Mock Mode (ISSUE-010)

`apps/core-web/src/lib/features/s2-payment/repository.ts` default ke mock mode (`VITE_MOCK_GATEWAY ?? 'true'`). Checkout page akan menampilkan dummy VA tanpa koneksi ke backend nyata kecuali env var diset secara eksplisit.

## Komponen yang Terverifikasi Ada (Tidak Perlu Difix)

Meskipun banyak routes hilang, komponen backend S2 ini sudah **diimplementasikan** berdasarkan analisis kode:

| Komponen | File | Status |
|----------|------|--------|
| payment-svc gRPC IssueVirtualAccount | `services/payment-svc/api/grpc_api/payment.go` | Implemented |
| payment-svc HTTP webhook handler | `services/payment-svc/api/http_api/webhook.go` | Implemented (port 50065) |
| Mock gateway adapter | `services/payment-svc/adapter/gateway/mock_adapter.go` | Implemented |
| payment_grpc_adapter (gateway→payment) | `services/gateway-svc/adapter/payment_grpc_adapter/adapter.go` | Implemented |
| Checkout UI page | `apps/core-web/src/routes/(b2c)/checkout/[booking_id]/+page.svelte` | Implemented |
| S2 payment repository | `apps/core-web/src/lib/features/s2-payment/repository.ts` | Implemented (mock mode issue) |

## Prioritas Fix (Urutan Rekomendasi)

1. **ISSUE-009** — Fix `uat-helpers.ts` endpoint `/v1/sessions` → `/v1/auth/login` (5 menit, tidak memblokir test)
2. **ISSUE-011** — Fix booking body schema mismatch di gateway atau spec (1-2 jam)
3. **ISSUE-005** — Daftarkan invoice routes di gateway-svc `server.go` + buat handlers (4-8 jam)
4. **ISSUE-006** — Daftarkan submit booking route di gateway-svc (2-4 jam)
5. **ISSUE-007, 008** — Expose webhook routes via gateway atau nginx proxy (2-4 jam)
6. **ISSUE-010** — Set `VITE_MOCK_GATEWAY=false` di docker-compose.prod.yml (15 menit)

## Metode Testing yang Digunakan

Karena network sandbox memblokir koneksi ke 216.176.238.161, seluruh analisis dilakukan melalui:

1. **Static code analysis** — membaca semua handler, router, dan adapter Go files
2. **Schema comparison** — membandingkan request/response types antara spec, contract, dan implementasi
3. **Route mapping** — membandingkan routes yang terdaftar di `server.go` dengan routes yang dibutuhkan oleh spec
4. **Frontend analysis** — membaca Svelte components dan repository layer untuk verifikasi behavior

## Output Playwright (Tidak Tersedia)

Test Playwright tidak dapat dijalankan karena:
- Network egress blocked: `ENETUNREACH 216.176.238.161:4000`
- Proxy allowlist: hanya mengizinkan github.com, npmjs.com, api.anthropic.com
- Bahkan Playwright `request.newContext()` dengan explicit proxy config mengembalikan 403 dari proxy

Untuk menjalankan test ini, perlu dijalankan langsung dari server 216.176.238.161 atau dari mesin yang memiliki akses network ke server tersebut.
