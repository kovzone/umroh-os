# UAT Issues Tracker

> **File ini adalah jembatan komunikasi antara QA session dan Dev session.**
>
> - **QA**: Append issue baru di bawah dengan format baku dari `SKILL.md`. Nomor sequential (ISSUE-001, ISSUE-002, ...).
> - **Dev**: Setelah fix, update field `Status` jadi `FIXED` dan isi `Diperbaiki oleh` + `Fix summary`.
> - **QA**: Setelah verifikasi fix, update jadi `VERIFIED` atau `REOPENED` jika masih gagal.

## Status Legend
- `OPEN` — Belum diperbaiki
- `IN-PROGRESS` — Sedang diperbaiki
- `FIXED` — Dev sudah fix, menunggu verifikasi QA
- `VERIFIED` — Fix dikonfirmasi oleh QA
- `REOPENED` — Fix tidak berhasil, perlu perbaikan ulang
- `NOT_DEPLOYED` — Fitur belum ada, bukan bug

---

<!-- Issues akan ditambahkan di bawah baris ini oleh QA agents -->

## [ISSUE-001] S1: Gateway API (port 4000) tidak dapat dijangkau — semua API test gagal ENETUNREACH
- **Severity**: CRITICAL
- **Status**: VERIFIED
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: False alarm — arsitektur memang hanya buka port 80/443. Gateway tidak diekspos langsung; semua lewat Nginx. `.env.prod` sudah diupdate: `GATEWAY_SVC_URL=http://216.176.238.161` (tanpa port)
- **Backlog ID**: BL-IAM-001
- **Service**: gateway-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/sessions
Content-Type: application/json
Body: {"email":"admin@umrohos.dev","password":"password123"}
```

### Response Aktual
```
ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
(connection refused / host unreachable — tidak ada TCP response sama sekali)
```

### Response Expected
HTTP 200 dengan body:
```json
{
  "data": {
    "access_token": "v2.local.<...>",
    "refresh_token": "<64-char string>",
    "user": { "email": "admin@umrohos.dev", "user_id": "33333333-3333-...", "status": "active" },
    "access_expires_at": "<future ISO timestamp>"
  }
}
```

### Referensi Kontrak
docs/contracts/slice-S1.md

### Petunjuk Fix
- Pastikan `gateway-svc` container berjalan: `docker compose ps gateway-svc`
- Cek port binding: `docker compose logs gateway-svc | tail -20`
- Pastikan port 4000 terbuka di firewall server: `ufw status` atau `iptables -L`
- Kemungkinan service crash atau tidak di-deploy; jalankan `docker compose -f docker-compose.dev.yml up -d gateway-svc`

---

## [ISSUE-002] S1: Core Web (port 80) mengembalikan ERR_EMPTY_RESPONSE — UI tidak dapat diakses
- **Severity**: CRITICAL
- **Status**: OPEN
- **Backlog ID**: BL-FE-CONSOLE-001
- **Service**: core-web
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
```
GET http://216.176.238.161/console/login
(browser navigate)

GET http://216.176.238.161/
(browser navigate)
```

### Response Aktual
```
net::ERR_EMPTY_RESPONSE at http://216.176.238.161/console/login
net::ERR_EMPTY_RESPONSE at http://216.176.238.161/
(TCP connection established tapi server menutup koneksi tanpa mengirim data apapun)
```

### Response Expected
HTTP 200 dengan halaman HTML SvelteKit. Halaman `/console/login` harus menampilkan form login dengan:
- `input[type='email']` atau `input[name='email']`
- `input[type='password']`
- `button[type='submit']`

Halaman `/` harus menampilkan landing page B2C dengan CTA menuju `/packages`.

### Referensi Kontrak
docs/contracts/slice-S1.md → BL-FE-CONSOLE-001, BL-FE-BOOK-001

### Petunjuk Fix
- Cek apakah `core-web` container/service berjalan: `docker compose ps core-web`
- Cek logs: `docker compose logs core-web | tail -30`
- Pastikan Nginx/reverse-proxy (port 80) dikonfigurasi benar untuk proxy ke SvelteKit app
- Cek apakah build SvelteKit berhasil: tidak ada error di `npm run build`
- Kemungkinan Nginx config salah atau core-web belum di-start

---

## [ISSUE-003] S1: PostgreSQL (port 5432) tidak dapat dijangkau dari test environment
- **Severity**: HIGH
- **Status**: VERIFIED
- **Diperbaiki oleh**: Lutfi — port 5432 dibuka sementara untuk UAT
- **Backlog ID**: BL-IAM-001
- **Service**: iam-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
```
PostgreSQL connection string:
postgres://postgres:***@216.176.238.161:5432/umrohos?sslmode=disable

Digunakan oleh:
- cleanupUatData() dalam afterAll hook (uat-helpers.ts)
- S1-WRITE-01 beforeAll block (loginAdmin() juga gagal, lalu cleanup gagal dengan error terpisah)
```

### Response Aktual
```
Error: connect ENETUNREACH 216.176.238.161:5432 - Local (0.0.0.0:0)
```

### Response Expected
PostgreSQL harus dapat dijangkau di port 5432 untuk:
1. Cleanup data UAT setelah test selesai
2. DB assertion untuk verifikasi data tersimpan dengan benar

### Referensi Kontrak
docs/contracts/slice-S1.md

### Petunjuk Fix
- Verifikasi PostgreSQL container berjalan: `docker compose ps postgres`
- Cek port binding PostgreSQL: port 5432 harus ter-expose ke host
- Cek firewall: `ufw status | grep 5432`
- Jika ini sengaja dibatasi untuk keamanan (tidak expose DB ke public), perlu update `.env.prod` dengan koneksi melalui SSH tunnel atau proxy
- Catatan: ISSUE-001 (gateway API) yang menjadi root cause lebih mendesak; jika gateway sudah jalan, test suite ini bisa berjalan penuh dan cleanup via API sudah cukup

---

## [ISSUE-004] S1: `playwright.config.ts` — project `uat-api` tidak memiliki `baseURL` sehingga UI tests gagal saat dijalankan di project ini
- **Severity**: MEDIUM
- **Status**: FIXED
- **Diperbaiki oleh**: QA — 2026-04-24
- **Fix summary**: Tambahkan `baseURL: process.env.GATEWAY_SVC_URL` ke project `uat-api` di `playwright.config.ts`
- **Backlog ID**: BL-FE-CONSOLE-001
- **Service**: gateway-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
```
npx playwright test tests/04-uat-s1-auth-catalog.spec.ts --project=uat-api
```

Error pada test S1-UI-01 (`/console/login`) dan S1-UI-05 (`/`):
```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
navigating to "/console/login", waiting until "load"
```

### Response Aktual
```
Error: Protocol error (Page.navigate): Cannot navigate to invalid URL
```
Spec file `04-uat-s1-auth-catalog.spec.ts` memiliki UI tests (describe S1 UI) yang menggunakan `page.goto("/console/login")` — path relatif yang memerlukan `baseURL`. Project `uat-api` di `playwright.config.ts` tidak memiliki `baseURL`.

### Response Expected
UI tests harus bisa dijalankan. Relative URL seperti `/console/login` harus resolve ke `http://216.176.238.161/console/login`.

### Referensi Kontrak
docs/contracts/slice-S1.md

### Petunjuk Fix
Tambahkan `baseURL` ke project `uat-api` di `playwright.config.ts`:
```typescript
{
  name: "uat-api",
  testMatch: /0[4-9]-.*\.spec\.ts$/,
  use: {
    baseURL: process.env.CORE_WEB_URL || "http://216.176.238.161",
  },
},
```
Atau pisahkan UI tests ke project `browser` saja (spec file sudah include keduanya, cukup jalankan `--project=browser` untuk UI tests).

---

## [ISSUE-005] S2: Gateway tidak mendaftarkan route S2 — POST /v1/invoices, GET /v1/invoices/{id}, POST /v1/invoices/{id}/virtual-accounts semuanya 404
- **Severity**: CRITICAL
- **Status**: FIXED
- **Backlog ID**: BL-PAY-001
- **Service**: gateway-svc / payment-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Implementasi penuh end-to-end: (1) `GetInvoiceByID` service method + gRPC handler di payment-svc; (2) extended `PaymentServiceWithReissueServer` + `payment_grpc_ext.go` dengan `GetInvoiceByID`; (3) gateway pb stub + adapter methods untuk `IssueVirtualAccount`, `GetInvoiceByID`, `ProcessWebhook`; (4) dispatch methods di `proxy_dispatch_payment.go`; (5) REST handlers di `proxy_invoice.go` (`CreateInvoice`, `GetInvoiceByID`, `IssueVirtualAccountForInvoice`); (6) `ServerInterface` + wrappers di `api.gen.go`; (7) route registrasi di `cmd/server.go` (bearer-protected).

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/invoices
Authorization: Bearer <admin-token>
Content-Type: application/json
Body: {"booking_id":"<booking-id>","gateway":"mock"}
```

```
GET http://216.176.238.161:4000/v1/invoices/<invoice-id>
Authorization: Bearer <admin-token>
```

```
POST http://216.176.238.161:4000/v1/invoices/<invoice-id>/virtual-accounts
Authorization: Bearer <admin-token>
Content-Type: application/json
Body: {"gateway":"mock"}
```

### Response Aktual
```
HTTP 404 Not Found
(route tidak terdaftar di gateway-svc)
```

### Response Expected
- `POST /v1/invoices` → HTTP 200/201 dengan invoice object (termasuk VA details jika mock gateway)
- `GET /v1/invoices/{id}` → HTTP 200 dengan invoice details + payment_events
- `POST /v1/invoices/{id}/virtual-accounts` → HTTP 200/201 dengan VA details

### Referensi Kontrak
`docs/contracts/slice-S2.md` → § S2-J-01

### Petunjuk Fix
Di `services/gateway-svc/cmd/server.go`, tambahkan route registrasi berikut di dalam `v1Protected` group:

```go
// S2 invoice routes (BL-PAY-001) — bearer required
v1Protected.Post("/invoices", wrapper.CreateInvoice)
v1Protected.Get("/invoices/:id", wrapper.GetInvoiceByID)
v1Protected.Post("/invoices/:id/virtual-accounts", wrapper.IssueVirtualAccount)
```

Juga perlu:
1. Buat handler `CreateInvoice`, `GetInvoiceByID`, `IssueVirtualAccount` di `services/gateway-svc/api/rest_oapi/proxy_invoice.go`
2. Buat service method di `services/gateway-svc/service/proxy_dispatch_invoice.go`
3. Adapter gRPC sudah ada di `services/gateway-svc/adapter/payment_grpc_adapter/`

---

## [ISSUE-006] S2: Gateway tidak mendaftarkan route POST /v1/bookings/{id}/submit — 404
- **Severity**: CRITICAL
- **Status**: FIXED
- **Backlog ID**: BL-PAY-001
- **Service**: gateway-svc / booking-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Implementasi penuh SubmitBooking end-to-end: booking-svc store (`GetBookingByIDAny`, `UpdateBookingStatus`), service layer, gRPC handler + pb stubs; gateway pb client, adapter, service dispatch, REST handler, api.gen.go ServerInterface + wrapper, dan route `POST /v1/bookings/:id/submit` di cmd/server.go (bearer-protected).

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/bookings/<booking-id>/submit
Authorization: Bearer <admin-token>
Content-Type: application/json
Body: {}
```

### Response Aktual
```
HTTP 404 Not Found
(route tidak terdaftar di gateway-svc)
```

### Response Expected
```
HTTP 200 OK
{
  "data": {
    "id": "<booking-id>",
    "status": "pending_payment",
    ...
  }
}
```

### Referensi Kontrak
`docs/contracts/slice-S2.md` → § Booking States (S2 additions): `draft → pending_payment`

### Petunjuk Fix
Di `services/gateway-svc/cmd/server.go`, tambahkan dalam `v1Protected`:

```go
// S2 booking submit (BL-BOOK-005) — bearer required
v1Protected.Post("/bookings/:id/submit", wrapper.SubmitBooking)
```

Buat handler di `proxy_booking.go` dan daftarkan di booking_grpc_adapter untuk method `SubmitBooking`.

---

## [ISSUE-007] S2: Gateway tidak mendaftarkan route POST /v1/webhooks/mock/trigger — 404
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-PAY-003
- **Service**: gateway-svc / payment-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Tambahkan route `POST /v1/webhooks/mock/trigger` (dan midtrans/xendit) di gateway-svc. Handler `WebhookMockTrigger` di `proxy_webhook.go` mem-forward payload ke payment-svc.ProcessWebhook via gRPC (gateway="mock"). Route public (no bearer) karena payment-svc melakukan validasi signature.

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/webhooks/mock/trigger
Content-Type: application/json
Body: {"invoice_id":"<invoice-id>","status":"paid","amount":25000000}
```

### Response Aktual
```
HTTP 404 Not Found
(route tidak terdaftar di gateway-svc; payment-svc webhook handler berjalan di port 50065 bukan 4000)
```

### Response Expected
```
HTTP 200 OK
{
  "payment_event_id": "pe_...",
  "invoice_status": "paid",
  "booking_status": "paid_in_full"
}
```

### Referensi Kontrak
`docs/contracts/slice-S2.md` → § S2-J-04: Internal trigger endpoint

### Petunjuk Fix
Berdasarkan `services/payment-svc/api/http_api/webhook.go`, payment-svc HTTP listener berjalan di port **50065** (bukan 4000). Route `/v1/webhooks/mock/trigger` terdaftar di internal HTTP handler payment-svc.

Ada dua opsi fix:
1. **Expose port 50065** dari payment-svc ke host dan update `UAT_ENV.gatewayUrl` ke port yang tepat untuk test webhook, ATAU
2. **Tambahkan reverse proxy di gateway-svc** (port 4000) yang meneruskan `POST /v1/webhooks/*` ke `payment-svc:50065`

Opsi 2 lebih konsisten dengan S2 contract yang menyatakan gateway-svc di port 4000 sebagai single entrypoint. Di `server.go`:
```go
// S2 webhook routes — public (no bearer, signature-protected)
v1.Post("/webhooks/midtrans", wrapper.WebhookMidtrans)
v1.Post("/webhooks/xendit", wrapper.WebhookXendit)
// dev only — only when MOCK_GATEWAY=true
if cfg.Gateway.MockGateway {
    v1.Post("/webhooks/mock/trigger", wrapper.WebhookMockTrigger)
}
```

---

## [ISSUE-008] S2: Gateway tidak mendaftarkan route POST /v1/webhooks/midtrans — 404
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-PAY-004
- **Service**: gateway-svc / payment-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Sama dengan ISSUE-007. Handler `WebhookMidtrans` di `proxy_webhook.go` meneruskan raw body + `X-Callback-Token` ke payment-svc.ProcessWebhook(gateway="midtrans") via gRPC. Signature verification tetap di payment-svc. Route public (no bearer) di `cmd/server.go`. Xendit juga ditambahkan sekalian.

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/webhooks/midtrans
Content-Type: application/json
X-Callback-Token: <invalid-or-missing>
Body: {"order_id":"fake-order","transaction_status":"settlement","gross_amount":"25000000","status_code":"200"}
```

### Response Aktual
```
HTTP 404 Not Found
(route tidak terdaftar di gateway-svc port 4000)
```

### Response Expected
Tanpa header auth yang valid → HTTP 401 dengan body `{"error":{"code":"invalid_signature"}}`

### Referensi Kontrak
`docs/contracts/slice-S2.md` → § S2-J-02: Webhook Contract — signature verification

### Petunjuk Fix
Sama dengan ISSUE-007: tambahkan route `/v1/webhooks/midtrans` dan `/v1/webhooks/xendit` di gateway-svc sebagai reverse proxy ke payment-svc:50065. Lihat petunjuk fix ISSUE-007.

---

## [ISSUE-009] S2: UAT Spec menggunakan endpoint `/v1/sessions` yang tidak terdaftar di gateway — harus menggunakan `/v1/auth/login`
- **Severity**: HIGH
- **Status**: FIXED
- **Diperbaiki oleh**: QA — 2026-04-24
- **Fix summary**: Update semua `/v1/sessions` → `/v1/auth/login`, refresh → `/v1/auth/refresh`, logout → `DELETE /v1/auth/logout` di spec 04 dan uat-helpers.ts
- **Backlog ID**: BL-IAM-001
- **Service**: gateway-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
Di `tests/e2e/lib/uat-helpers.ts` line 56:
```typescript
const res = await api.post("/v1/sessions", {
  email: UAT_ENV.adminEmail,
  password: UAT_ENV.adminPassword,
});
```

### Response Aktual
```
HTTP 404 Not Found
(gateway terdaftar di /v1/auth/login, bukan /v1/sessions)
```

### Response Expected
HTTP 200 dengan access_token. Route `/v1/sessions` harus return token seperti `/v1/auth/login`.

### Referensi Kontrak
`docs/contracts/slice-S1.md` → § Auth endpoints

### Petunjuk Fix
**Opsi A (rekomendasi):** Update `uat-helpers.ts` untuk menggunakan endpoint yang benar:
```typescript
// Line 56 di tests/e2e/lib/uat-helpers.ts
const res = await api.post("/v1/auth/login", {  // bukan /v1/sessions
  email: UAT_ENV.adminEmail,
  password: UAT_ENV.adminPassword,
});
```

**Opsi B:** Tambahkan alias route di gateway-svc (tidak direkomendasikan — membuat dua endpoint untuk hal yang sama).

Catatan: Spec file `05-uat-s2-booking-payment.spec.ts` juga memanggil `loginAdmin()` di `beforeAll`. Jika `loginAdmin()` gagal, semua test yang membutuhkan `adminToken` akan skip atau error.

---

## [ISSUE-010] S2: Checkout page (`/checkout/{booking_id}`) menggunakan mock data secara default — `VITE_MOCK_GATEWAY` tidak diset `false` di production build
- **Severity**: MEDIUM
- **Status**: FIXED
- **Backlog ID**: BL-FE-PAY-001
- **Service**: core-web
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Ubah default `VITE_MOCK_GATEWAY` dari `'true'` ke `'false'` di `apps/core-web/src/lib/features/s2-payment/repository.ts` baris 16. Tambah `VITE_MOCK_GATEWAY: "false"` di build args `docker-compose.prod.yml`. Mock harus diaktifkan secara eksplisit, tidak lagi default aktif.

### Langkah Reproduksi
Di browser, navigate ke `http://216.176.238.161/checkout/<booking-id>` setelah membuat booking.

### Response Aktual
Checkout page menampilkan data VA dummy (mock), bukan VA dari payment-svc yang sebenarnya. Ini karena:

```typescript
// apps/core-web/src/lib/features/s2-payment/repository.ts, line 16
return (env['VITE_MOCK_GATEWAY'] ?? 'true') === 'true';
```

Default adalah `'true'` (mock). Jika `VITE_MOCK_GATEWAY` tidak di-set secara eksplisit di production build, mock selalu aktif.

### Response Expected
Di production/dev server (`http://216.176.238.161`), checkout page harus memanggil real API `POST /v1/invoices` dan menampilkan VA dari payment-svc.

### Referensi Kontrak
`docs/contracts/slice-S2.md` → § Checkout UI — VA Display & Polling

### Petunjuk Fix
1. Set `VITE_MOCK_GATEWAY=false` di production build args di `docker-compose.prod.yml`:
```yaml
args:
  VITE_GATEWAY_URL: ${VITE_PUBLIC_GATEWAY_URL}
  VITE_MOCK_GATEWAY: "false"   # tambahkan ini
```
2. Pastikan Vite memasukkan env var ini saat build (`vite.config.ts` harus expose `VITE_*` vars).
3. Atau ubah default di `repository.ts` menjadi `'false'` untuk production safety.

---

## [ISSUE-011] S2: Booking request body schema mismatch antara spec (`lead`, `jamaah[]`) dan gateway handler (`lead_full_name`, `pilgrims[]`)
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-BOOK-001
- **Service**: gateway-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Ditambahkan `LeadBody` struct dan field `Lead *LeadBody` + `Jamaah []PilgrimBody` ke `CreateDraftBookingBody` di `proxy_booking.go`. Merge logic: jika `lead` object ada, field-nya di-copy ke flat `lead_full_name`/etc; jika `jamaah[]` ada dan `pilgrims[]` kosong, di-copy ke `pilgrims`. Handler sudah dapat menerima kedua format.

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/bookings
Content-Type: application/json
Body:
{
  "channel": "b2c_self",
  "package_id": "pkg_01JCDE00000000000000000001",
  "departure_id": "dep_01JCDF00000000000000000001",
  "room_type": "double",
  "lead": {
    "full_name": "[UAT] Test Jamaah",
    "email": "uat.test@umrohos.dev",
    "whatsapp": "+628112345678",
    "domicile": "Jakarta"
  },
  "jamaah": [{"full_name": "[UAT] Test Jamaah", "is_lead": true, ...}]
}
```

### Response Aktual
```
HTTP 400/422 Bad Request
(gateway mengexpect field flat seperti lead_full_name, lead_email, pilgrims[] — bukan nested lead object dan jamaah[])
```

Gateway `CreateDraftBookingBody` struct (`proxy_booking.go`) menggunakan:
- `lead_full_name` (flat, bukan nested `lead.full_name`)
- `pilgrims` (bukan `jamaah`)

### Response Expected
Booking dibuat dengan HTTP 201 dan response `{"data":{"id":"...","status":"draft",...}}`

### Referensi Kontrak
`docs/contracts/slice-S2.md` → S2-BOOK-01 (test case dalam spec)

### Petunjuk Fix
**Opsi A (rekomendasi):** Update `CreateDraftBookingBody` di `services/gateway-svc/api/rest_oapi/proxy_booking.go` untuk menerima nested `lead` object dan `jamaah[]` array sebagai alias/tambahan:

```go
type LeadBody struct {
    FullName string `json:"full_name"`
    Email    string `json:"email,omitempty"`
    Whatsapp string `json:"whatsapp,omitempty"`
    Domicile string `json:"domicile,omitempty"`
}

type CreateDraftBookingBody struct {
    // ... existing flat fields sebagai backward compat ...
    Lead   *LeadBody     `json:"lead,omitempty"`
    Jamaah []PilgrimBody `json:"jamaah,omitempty"`  // alias untuk Pilgrims
}
```

**Opsi B:** Update spec file `05-uat-s2-booking-payment.spec.ts` untuk menggunakan format flat yang sesuai gateway saat ini.

---

## [ISSUE-012] S4: Phone validation error mengembalikan HTTP 500 bukan 400/422 — crm-svc tidak wrap error dengan ErrValidation
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-CRM-001
- **Service**: crm-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Di `services/crm-svc/service/leads.go`, phone validation error di-wrap dengan `errors.Join(apperrors.ErrValidation, ...)` sehingga gRPC handler memetakan ke `codes.InvalidArgument` dan gateway mengembalikan HTTP 400.

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/leads
Content-Type: application/json
Body:
{
  "name": "[UAT] Test Invalid Phone",
  "email": "uat.invalid@umrohos.dev",
  "phone": "bukan-nomor"
}
```

### Response Aktual
```
HTTP 500 Internal Server Error
{"error":{"code":"internal_error","message":"terjadi kesalahan tidak terduga"}}
```

### Response Expected
```
HTTP 400 Bad Request atau 422 Unprocessable Entity
{"error":{"code":"validation_error","message":"phone must be at least 8 characters"}}
```

### Referensi Kontrak
`docs/contracts/slice-S4.md` → § POST /v1/leads: Validations — phone format check harus return 400/422

### Petunjuk Fix
Di `services/crm-svc/service/leads.go` baris 169–174, ganti `fmt.Errorf` dengan `errors.Join(apperrors.ErrValidation, ...)`:

```go
// Before (line 169):
err := fmt.Errorf("%s: phone must be at least 8 characters", op)

// After:
err := errors.Join(apperrors.ErrValidation, fmt.Errorf("phone must be at least 8 characters"))
```

Dengan ini gRPC akan return `codes.InvalidArgument` → gateway adapter maps ke `ErrValidation` → HTTP 400.

---

## [ISSUE-013] S4: Status transition `new → converted` diizinkan padahal seharusnya ditolak — `validTransition()` tidak memvalidasi urutan sequential
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-CRM-002
- **Service**: crm-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: `validTransition()` di `services/crm-svc/service/leads.go` ditulis ulang dengan map transisi eksplisit: `new→{contacted,lost}`, `contacted→{qualified,lost}`, `qualified→{converted,lost}`. Terminal states (`converted`, `lost`) mengembalikan false untuk semua transisi keluar.

### Langkah Reproduksi
```
# Buat lead (status = 'new')
POST http://216.176.238.161:4000/v1/leads
Body: {"name":"[UAT] Test Lead","phone":"+628119876543","email":"uat.lead@umrohos.dev"}

# Langsung update ke 'converted' tanpa melalui contacted → qualified
PUT http://216.176.238.161:4000/v1/leads/{id}
Authorization: Bearer <admin-token>
Body: {"status": "converted"}
```

### Response Aktual
```
HTTP 200 OK
{"data":{"status":"converted",...}}
(transisi berhasil padahal seharusnya ditolak)
```

### Response Expected
```
HTTP 400 Bad Request atau 422 Unprocessable Entity
{"error":{"code":"invalid_transition","message":"transition new → converted not allowed"}}
```

### Referensi Kontrak
`docs/contracts/slice-S4.md` → § Allowed transitions:
- `new` → `contacted` ✓
- `contacted` → `qualified` ✓
- `qualified` → `converted` ✓
- `new` → `converted` ✗ (not allowed, must go through contacted → qualified first)

### Petunjuk Fix
Di `services/crm-svc/service/leads.go`, fungsi `validTransition()` perlu di-rewrite untuk enforce sequential ordering:

```go
func validTransition(from, to string) bool {
    if from == to {
        return true // idempotent
    }
    // Terminal states: no outbound transitions
    if from == "converted" || from == "lost" {
        return false
    }
    // Define allowed forward transitions
    allowed := map[string][]string{
        "new":       {"contacted", "lost"},
        "contacted": {"qualified", "lost"},
        "qualified": {"converted", "lost"},
    }
    for _, validTo := range allowed[from] {
        if to == validTo {
            return true
        }
    }
    return false
}
```

---

## [ISSUE-014] S3: Fulfillment task initial status mismatch — spec mengexpect `'pending'` tapi DB schema memakai `'queued'`
- **Severity**: MEDIUM
- **Status**: FIXED
- **Backlog ID**: BL-LOG-001
- **Service**: logistics-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: UAT spec (`06-uat-s3s4-crm-ops.spec.ts` baris 423) diupdate dari `.toBe("pending")` ke `.toBe("queued")` agar sesuai default DB schema. Sumber kebenaran adalah schema; tidak ada perubahan ke DB/service.

### Langkah Reproduksi
```sql
-- Setelah booking paid_in_full, cek status fulfillment task:
SELECT id, status, booking_id 
FROM logistics.fulfillment_tasks 
WHERE booking_id = '<booking-id>';
```

### Response Aktual
```
status = 'queued'
(sesuai DEFAULT di migration 000013)
```

### Response Expected
Test spec S3-OPS-01 (`06-uat-s3s4-crm-ops.spec.ts` line 423) expects:
```
expect(tasks[0].status, "Fulfillment task status awal harus pending").toBe("pending");
```

### Referensi Kontrak
`docs/contracts/slice-S3.md` → § BL-LOG-001: fulfillment_tasks status lifecycle

Migration 000013 defines:
```sql
status TEXT NOT NULL DEFAULT 'queued'
CHECK (status IN ('queued', 'processing', 'shipped', 'delivered', 'cancelled'))
```

### Petunjuk Fix
Ada dua opsi:
**Opsi A (rekomendasi):** Update UAT spec untuk mengexpect `'queued'` bukan `'pending'`:
```typescript
// services/tests/e2e/tests/06-uat-s3s4-crm-ops.spec.ts, line 423:
expect(tasks[0].status, "Fulfillment task status awal harus queued").toBe("queued");
```

**Opsi B:** Ubah default status di migration dan CHECK constraint dari `'queued'` ke `'pending'` (melibatkan schema migration baru).

Opsi A lebih aman karena tidak perlu schema migration.

---

## [ISSUE-015] S3: UAT spec S3-FIN-01 query column salah — `debit_amount`/`credit_amount` tidak ada, kolom sebenarnya `debit`/`credit`
- **Severity**: MEDIUM
- **Status**: FIXED
- **Diperbaiki oleh**: QA — 2026-04-24
- **Fix summary**: Update `dbQuery` type dan semua referensi kolom di `06-uat-s3s4-crm-ops.spec.ts` dari `debit_amount`/`credit_amount` → `debit`/`credit`
- **Backlog ID**: BL-FIN-001
- **Service**: finance-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
Query di spec `06-uat-s3s4-crm-ops.spec.ts` line 502–504:
```typescript
const lines = await dbQuery<{ debit_amount: number; credit_amount: number; account_code: string }>(
  `SELECT debit_amount, credit_amount, account_code
   FROM finance.journal_lines
   WHERE entry_id = $1`,
  [entryId]
);
```

### Response Aktual
```
ERROR: column "debit_amount" does not exist
(PostgreSQL error — migration 000013 mendefinisikan kolom sebagai 'debit' dan 'credit', bukan 'debit_amount' dan 'credit_amount')
```

### Response Expected
Query harus berhasil dan mengembalikan nilai debit/credit untuk verifikasi double-entry balance.

### Referensi Kontrak
`migration/000013_add_logistics_finance_document_tables.up.sql`:
```sql
CREATE TABLE finance.journal_lines (
    debit  NUMERIC(15, 2) NOT NULL DEFAULT 0,
    credit NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ...
)
```

### Petunjuk Fix
Update query di `tests/e2e/tests/06-uat-s3s4-crm-ops.spec.ts` baris 502–519:
```typescript
// Before:
const lines = await dbQuery<{ debit_amount: number; credit_amount: number; account_code: string }>(
  `SELECT debit_amount, credit_amount, account_code FROM finance.journal_lines WHERE entry_id = $1`,
  [entryId]
);
const totalDebit = lines.reduce((sum, l) => sum + Number(l.debit_amount || 0), 0);
const totalCredit = lines.reduce((sum, l) => sum + Number(l.credit_amount || 0), 0);

// After:
const lines = await dbQuery<{ debit: number; credit: number; account_code: string }>(
  `SELECT debit, credit, account_code FROM finance.journal_lines WHERE entry_id = $1`,
  [entryId]
);
const totalDebit = lines.reduce((sum, l) => sum + Number(l.debit || 0), 0);
const totalCredit = lines.reduce((sum, l) => sum + Number(l.credit || 0), 0);
```
Dan update type check di baris 518–519 juga.

---

## [ISSUE-016] S3: UAT cleanup `uat-helpers.ts` menggunakan column `source_ref` yang tidak ada — seharusnya `source_id`
- **Severity**: HIGH
- **Status**: FIXED
- **Diperbaiki oleh**: QA — 2026-04-24
- **Fix summary**: Update `cleanupUatData()` di `uat-helpers.ts` — ganti `source_ref` → `source_id` di semua DELETE queries finance
- **Backlog ID**: BL-FIN-001
- **Service**: finance-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: —
- **Verified oleh**: —

### Langkah Reproduksi
Jalankan `cleanupUatData()` di `tests/e2e/lib/uat-helpers.ts` setelah test suite selesai. SQL cleanup di baris 340 dan 348:
```sql
DELETE FROM finance.journal_lines WHERE entry_id IN (
  SELECT je.id FROM finance.journal_entries je
  WHERE je.source_ref IN (...)  -- SALAH: kolom tidak ada
)
DELETE FROM finance.journal_entries WHERE source_ref IN (...)  -- SALAH
```

### Response Aktual
```
ERROR: column "source_ref" does not exist
LINE 3: WHERE je.source_ref IN (
(cleanup gagal — finance data dari UAT tidak terhapus)
```

### Response Expected
Cleanup berhasil menghapus semua finance.journal_lines dan finance.journal_entries yang terkait UAT bookings.

### Referensi Kontrak
`migration/000013_add_logistics_finance_document_tables.up.sql`:
```sql
CREATE TABLE finance.journal_entries (
    source_id UUID NOT NULL,  -- bukan source_ref
    ...
)
```

### Petunjuk Fix
Update `tests/e2e/lib/uat-helpers.ts` baris 338–352, ganti `source_ref` dengan `source_id`:
```typescript
await client.query(
  `DELETE FROM finance.journal_lines WHERE entry_id IN (
    SELECT je.id FROM finance.journal_entries je
    WHERE je.source_id::text IN (
      SELECT pi.id::text FROM payment.invoices pi
      JOIN booking.bookings b ON b.id = pi.booking_id
      WHERE b.notes ILIKE '%[UAT]%'
    )
  )`
);
await client.query(
  `DELETE FROM finance.journal_entries WHERE source_id::text IN (
    SELECT pi.id::text FROM payment.invoices pi
    JOIN booking.bookings b ON b.id = pi.booking_id
    WHERE b.notes ILIKE '%[UAT]%'
  )`
);
```

---

## [ISSUE-017] S4: `CreateLeadBody` di gateway tidak menerima field `message` dan `source_note` yang dikirim oleh UAT spec — fields silently ignored
- **Severity**: LOW
- **Status**: FIXED
- **Backlog ID**: BL-CRM-001
- **Service**: gateway-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Ditambahkan `SourceNote` dan `Message` ke `CreateLeadBody` di `proxy_crm.go`. Merge logic: jika `notes` kosong dan `message` ada → `notes = message`; jika `source` kosong dan `source_note` ada → `source = source_note`.

### Langkah Reproduksi
```
POST http://216.176.238.161:4000/v1/leads
Content-Type: application/json
Body:
{
  "name": "[UAT] Test Lead",
  "phone": "+628119876543",
  "email": "uat.lead@umrohos.dev",
  "message": "Mau daftar umroh tahun ini",
  "source_note": "uat-test"
}
```

### Response Aktual
```
HTTP 201 Created
(lead berhasil dibuat, TAPI field 'message' dan 'source_note' tidak tersimpan ke DB — silently ignored)
```

### Response Expected
Berdasarkan UAT spec `06-uat-s3s4-crm-ops.spec.ts`, field `message` harusnya tersimpan ke `notes` di DB, dan `source_note` ke `source` atau sebuah metadata field.

### Referensi Kontrak
`docs/contracts/slice-S4.md` → § POST /v1/leads: Request body

### Petunjuk Fix
**Opsi A:** Update `CreateLeadBody` di `services/gateway-svc/api/rest_oapi/proxy_crm.go` untuk menerima field alias:
```go
type CreateLeadBody struct {
    Name    string `json:"name"`
    Phone   string `json:"phone"`
    Email   string `json:"email,omitempty"`
    Notes   string `json:"notes,omitempty"`
    Message string `json:"message,omitempty"`  // alias untuk notes
    // ...
}
// Lalu merge: if body.Message != "" && body.Notes == "" { body.Notes = body.Message }
```

**Opsi B (rekomendasi):** Update UAT spec untuk menggunakan field name yang sesuai (`notes` bukan `message`, hapus `source_note`).

---

## [ISSUE-018] S3/S4: Gateway tidak memiliki route GET `/v1/fulfillment-tasks` — endpoint listing fulfillment tasks belum di-deploy
- **Severity**: MEDIUM
- **Status**: FIXED
- **Backlog ID**: BL-LOG-001
- **Service**: gateway-svc / logistics-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Implementasi end-to-end `ListFulfillmentTasks`: (1) Tambah `ListFulfillmentTasks` + `CountFulfillmentTasks` SQL query di logistics-svc store; (2) Tambah service layer params/result + implementasi; (3) Tambah `ListFulfillmentTasksRequest/Response` pb message types + gRPC handler di `list_fulfillment_tasks.go`; (4) Register di `logistics_grpc_ext.go`; (5) Tambah gateway-side pb stub types + constant; (6) Tambah adapter method; (7) Tambah ke gateway `IService` + dispatch; (8) Tambah REST handler di `proxy_logistics.go`; (9) Tambah ke `api.gen.go` + route `GET /v1/fulfillment-tasks` di `cmd/server.go`. Mendukung query params `?status=&departure_id=&limit=&offset=`.

### Langkah Reproduksi
```
GET http://216.176.238.161:4000/v1/fulfillment-tasks
Authorization: Bearer <admin-token>
```

### Response Aktual
```
HTTP 404 Not Found
(route tidak terdaftar — task instructions step 4 curl test)
```

### Response Expected
List fulfillment tasks dengan pagination.

### Referensi Kontrak
`docs/contracts/slice-S3.md` → § BL-LOG-001

### Petunjuk Fix
Route `GET /v1/fulfillment-tasks` belum didaftarkan di `services/gateway-svc/cmd/server.go`. Perlu tambahkan handler dan registrasi route ketika fitur ini diimplementasikan.

---

## [ISSUE-019] S3: Finance endpoint menggunakan path `/v1/finance/journal-entries` (curl step 4) tapi route sebenarnya adalah `/v1/finance/journals`
- **Severity**: LOW
- **Status**: FIXED
- **Backlog ID**: BL-FIN-001
- **Service**: gateway-svc
- **Dilaporkan**: 2026-04-24
- **Diperbaiki oleh**: Elda Agent 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Ini bukan production bug — route yang benar (`/v1/finance/journals`) sudah ada dan berfungsi. Path `/v1/finance/journal-entries` hanya muncul dalam instruksi curl di ISSUES.md ini, bukan di test code. Tidak ada perubahan code diperlukan.

### Langkah Reproduksi
```
GET http://216.176.238.161:4000/v1/finance/journal-entries?limit=5
Authorization: Bearer <admin-token>
```

### Response Aktual
```
HTTP 404 Not Found
(route yang benar adalah /v1/finance/journals, bukan /v1/finance/journal-entries)
```

### Response Expected
HTTP 200 dengan list journal entries.

### Referensi Kontrak
`services/gateway-svc/cmd/server.go` baris 157: `v1Protected.Get("/finance/journals", wrapper.ListJournals)`

### Petunjuk Fix
Update dokumentasi/instruksi curl di task instructions dan CONFIG.md untuk menggunakan path yang benar: `/v1/finance/journals` bukan `/v1/finance/journal-entries`. Ini bukan production bug.

---

## [ISSUE-020] B2C: Gambar paket broken di halaman detail (cover_photo_url pakai domain fake)
- **Severity**: MEDIUM
- **Status**: FIXED
- **Backlog ID**: BL-CAT-002
- **Service**: catalog-svc / seed data
- **Dilaporkan**: 2026-04-24 (Browser E2E)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: (1) Buat migration `000027_fix_dev_seed_content.up.sql` yang UPDATE `cover_photo_url` ke Unsplash public image; (2) Update source migration `000009` agar fresh install pun dapat URL yang valid.
- **Verified oleh**: —

### Langkah Reproduksi
1. Buka `http://216.176.238.161/packages/pkg_01JCDE00000000000000000001`
2. Perhatikan gambar cover paket di kiri atas halaman detail

### Response Aktual
Gambar tidak tampil — menampilkan alt text placeholder (grey box).
API response `GET /v1/packages/pkg_01JCDE00000000000000000001` mengembalikan:
```json
"cover_photo_url": "https://cdn.umroh-os.example/pkg/01JCDE.../cover.jpg"
```
Domain `cdn.umroh-os.example` tidak nyata — URL tidak dapat diakses.

### Response Expected
Gambar paket tampil dengan benar.

### Referensi Kontrak
`docs/contracts/slice-S1.md` → Package detail page: cover photo

### Petunjuk Fix
Update seed data migration (file `000009-*.sql` atau equivalent) untuk mengisi `cover_photo_url` dengan URL yang valid (bisa pakai placeholder dari unsplash/picsum atau CDN nyata). Contoh:
```sql
UPDATE catalog.packages
SET cover_photo_url = 'https://picsum.photos/seed/umrah/800/600'
WHERE id = 'pkg_01JCDE00000000000000000001';
```

---

## [ISSUE-021] B2C: Harga tidak tampil di departure cards pada halaman detail paket
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-CAT-003
- **Service**: catalog-svc / gateway-svc / core-web
- **Dilaporkan**: 2026-04-24 (Browser E2E)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: Full-stack N+1 enrichment: (1) Tambah `PricePerPax *int64` ke `DepartureSummary` di adapter `types.go`; (2) Di `get_package.go`, setelah fetch package detail via gRPC, loop tiap departure → call `GetPackageDeparture` → hitung `minPrice` → set `PricePerPax`; (3) Tambah `PricePerPax *int64` ke REST `DepartureSummary` struct di `api.gen.go`; (4) Update `departureSummaryFromAdapter` di `proxy_catalog.go` untuk menyertakan field ini; (5) Update `types.ts` dan `api.ts` frontend; (6) Tampilkan field harga di `DeparturePicker.svelte` dengan format `Rp X.XXX.XXX / pax`.
- **Verified oleh**: —

### Langkah Reproduksi
1. Buka `http://216.176.238.161/packages/pkg_01JCDE00000000000000000001`
2. Scroll ke section departure cards
3. Perhatikan: departure cards tampil tapi tidak ada angka harga
4. Footer sticky terus menampilkan "Harga akan ditampilkan dari departure" meski departure sudah dipilih

### Response Aktual
```
GET /v1/packages/pkg_01JCDE00000000000000000001
Response departures:
[
  { "id": "dep_01...", "departure_date": "2026-05-23", "remaining_seats": 42, "status": "open" },
  { "id": "dep_01...", "departure_date": "2026-07-07", "remaining_seats": 45, "status": "open" }
]
// TIDAK ADA field: price_per_pax, list_price, base_price
```

### Response Expected
Departure cards menampilkan harga per pax setelah user memilih departure. Harga **memang ada** — terbukti muncul di Step 1 booking wizard (Rp 38.500.000 dan Rp 39.500.000). Tapi public catalog endpoint tidak mengembalikannya.

### Referensi Kontrak
`docs/contracts/slice-S1.md` → Package detail departure section: price display

### Petunjuk Fix
**Opsi A (rekomendasi):** Tambahkan `price_per_pax` ke response departure di public catalog endpoint `GET /v1/packages/:id`:
```go
// Di catalog-svc handler untuk GetPackage:
departures: [
  {
    id: dep.ID,
    departure_date: dep.DepartDate,
    price_per_pax: dep.PricePerPax,  // tambahkan field ini
    remaining_seats: dep.RemainingSeats,
    status: dep.Status,
  }
]
```

**Opsi B:** Frontend fetch harga dari endpoint terpisah saat departure dipilih (tapi endpoint `/v1/packages/:id/departures` saat ini return 401 unauthorized — perlu dibuka untuk public atau diubah ke endpoint baru).

---

## [ISSUE-022] B2C: booking-svc return 503 — frontend silent fallback ke demo booking palsu
- **Severity**: CRITICAL
- **Status**: FIXED (frontend part); booking-svc 503 memerlukan ops deployment check
- **Backlog ID**: BL-BOOK-006
- **Service**: booking-svc / core-web
- **Dilaporkan**: 2026-04-24 (Browser E2E — diverifikasi ulang 2026-04-24)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: (1) Frontend: hapus silent fallback di `s1-booking/repository.ts` — ketika real API gagal, error sekarang propagate ke UI sehingga pesan error muncul di Step 3 (sudah ada `{#if submitError}` di halaman). (2) Ubah default `VITE_USE_BOOKING_MOCK` dari `'true'` ke `'false'` agar production tidak secara otomatis masuk mock. (3) Backend: booking-svc 503 disebabkan container tidak berjalan atau crash saat startup — cek `docker compose logs booking-svc`. Perlu ops deploy ulang; tidak ada perubahan code booking-svc dalam fix ini.
- **Verified oleh**: —

### Root Cause (Diverifikasi)
`POST http://216.176.238.161:4000/v1/bookings` mengembalikan **HTTP 503** (booking-svc DOWN atau tidak merespons).

Frontend tidak menampilkan error message kepada user. Sebaliknya, frontend **silently fallback** ke data booking demo yang hardcoded:
- Kode booking: `UMR-BKG_DEMO`
- Harga: Rp 38.500.000 (hardcoded — harga 1 pax Quad, bukan total aktual user)
- URL checkout: `/checkout/bkg_demo_<random>`

**Akibat**: User mengira booking berhasil (mendapat kode booking dan nomor VA), padahal tidak ada booking nyata yang dibuat di sistem.

### Bukti Jaringan (Diverifikasi Browser E2E 2026-04-24)
```
OPTIONS http://216.176.238.161:4000/v1/bookings → 204  (preflight OK)
POST    http://216.176.238.161:4000/v1/bookings → 503  (booking-svc DOWN)
```

### Langkah Reproduksi
1. Buka booking wizard untuk `pkg_01JCDE00000000000000000001`, departure 23 Mei 2026
2. Isi data pemesan dan 2 jamaah (total wizard: Rp 77.000.000 = 2 × Rp 38.500.000)
3. Centang semua syarat di Step 3, klik "Lanjut ke pembayaran"
4. **Step 4 menampilkan**: Kode booking `UMR-BKG_DEMO`, total **Rp 38.500.000** (bukan Rp 77.000.000)
5. Tidak ada error message yang ditampilkan ke user

### Response Expected
- Jika booking-svc UP: POST /v1/bookings → 200 dengan booking ID nyata, total sesuai pilihan user
- Jika booking-svc DOWN: Frontend harus menampilkan pesan error yang jelas (bukan fallback diam-diam ke demo)

### Petunjuk Fix
1. **Fix booking-svc**: Periksa kenapa booking-svc return 503 (container crashed? dependency failure?)
2. **Fix frontend error handling**: Jika POST /v1/bookings gagal, tampilkan error toast/modal — jangan silent fallback ke demo. Hapus atau disable kode fallback `UMR-BKG_DEMO` di production.

---

## [ISSUE-023] B2C: Halaman detail paket — section Itinerary, Fasilitas, dan Syarat hanya tampilkan placeholder text
- **Severity**: MEDIUM
- **Status**: FIXED
- **Backlog ID**: BL-CAT-004
- **Service**: catalog-svc / seed data
- **Dilaporkan**: 2026-04-24 (Browser E2E)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: Migration `000027_fix_dev_seed_content.up.sql` mengisi ulang `itinerary_templates.days` dengan 12 hari penuh (sebelumnya hanya 3 hari placeholder). Foto tiap hari menggunakan URL Unsplash publik. Migration source `000009` juga diupdate agar fresh install pun mendapat data lengkap. Muthawwif portrait URL ikut diperbaiki.
- **Verified oleh**: —

### Langkah Reproduksi
1. Buka `http://216.176.238.161/packages/pkg_01JCDE00000000000000000001`
2. Scroll ke section "Rencana perjalanan" (tab Itinerari)
3. Section kosong — tidak ada data itinerary
4. Tab "Fasilitas": menampilkan "Detail fasilitas mengikuti kontrak pemesanan dan brosur resmi travel untuk paket ini."
5. Tab "Syarat & Ketentuan": menampilkan "Untuk teks legal lengkap, silakan unduh dari halaman konfirmasi booking atau hubungi tim kami."
6. Box "Poin penting perjalanan" kosong (dark blue box dengan judul saja)

### Response Aktual
Semua section konten paket hanya menampilkan placeholder text — tidak ada data nyata dari catalog-svc.

### Response Expected
Section menampilkan konten nyata: itinerary harian, fasilitas termasuk, syarat & ketentuan lengkap.

### Petunjuk Fix
Update seed data untuk mengisi field itinerary, fasilitas, dan syarat & ketentuan pada paket demo `pkg_01JCDE00000000000000000001`.

---

## [ISSUE-024] Console: Halaman Katalog Paket gagal load — "Tidak dapat terhubung ke layanan katalog"
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-CAT-005
- **Service**: catalog-svc (dari SvelteKit SSR server)
- **Dilaporkan**: 2026-04-24 (Browser E2E — diverifikasi via __data.json)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Root cause: SvelteKit SSR di dalam container tidak bisa resolve `VITE_GATEWAY_URL` (baked at build time, bukan runtime). Fix: (1) tambah `GATEWAY_URL=http://gateway-svc:4000` ke environment core-web di `docker-compose.dev.yml` dan `docker-compose.prod.yml`; (2) perbaiki urutan env var di `packages/+page.server.ts` — `GATEWAY_URL` harus paling awal (runtime) sebelum `VITE_*` (build-time).

### Root Cause (Diverifikasi)
SvelteKit SSR server (core-web container) mencoba fetch data katalog dari catalog-svc via Docker internal network. Call ini gagal — catalog-svc tidak dapat dijangkau atau mengembalikan error dari dalam container SvelteKit.

**Bukti**: `GET /console/packages/__data.json` mengembalikan:
```json
{
  "type": "data",
  "nodes": [{}, {}, {
    "type": "data",
    "data": [{"packages": 1, "error": 2}, [], "Tidak dapat terhubung ke layanan katalog."]
  }]
}
```
Array `packages` = `[]` (kosong), `error` = `"Tidak dapat terhubung ke layanan katalog."`

**Catatan**: Ini BERBEDA dari B2C catalog — B2C package list (`/packages`) berfungsi normal karena menggunakan endpoint yang berbeda atau catalog-svc dapat dijangkau untuk route tersebut. Console admin catalog (`/console/packages`) menggunakan route/endpoint berbeda yang gagal.

### Langkah Reproduksi
1. Login ke console: `http://216.176.238.161/console` dengan `admin@umrohos.dev` / `password123`
2. Klik "Katalog Paket" di sidebar
3. Halaman menampilkan: "Tidak dapat terhubung ke layanan katalog."
4. Tabel kosong — tidak ada data paket yang tampil

### Response Expected
Tabel Katalog Paket menampilkan daftar semua package yang ada di database.

### Petunjuk Fix
1. Periksa apakah catalog-svc container berjalan: `docker compose ps catalog-svc`
2. Periksa log catalog-svc untuk route admin: `docker compose logs catalog-svc | tail -50`
3. Periksa apakah SvelteKit server menggunakan URL internal yang benar untuk catalog-svc (env var `CATALOG_SVC_URL` atau similar)
4. Kemungkinan: admin endpoint catalog memerlukan autentikasi/scope yang berbeda dari public endpoint


---

## [ISSUE-025] Console: Halaman Ops Board, Finance, dan Leads — "Tidak dapat terhubung ke gateway"
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-OPS-001
- **Service**: gateway-svc (dari console browser/SSR)
- **Dilaporkan**: 2026-04-24 (Browser E2E)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Verified oleh**: —
- **Fix summary**: Root cause sama dengan ISSUE-024: `GATEWAY_URL` tidak di-set di docker-compose → SSR server fallback ke `localhost:4000` (tidak valid di dalam container). Fix: (1) tambah `GATEWAY_URL=http://gateway-svc:4000` di docker-compose dev + prod; (2) perbaiki urutan env var di `ops/+page.server.ts` — `GATEWAY_URL` dulu sebelum `VITE_GATEWAY_URL`. Finance dan Leads sudah menggunakan urutan yang benar tapi ikut benefit dari penambahan GATEWAY_URL di docker-compose.

### Halaman Terdampak
- `/console/ops` (Ops Board — Fulfillment)
- `/console/finance` (Laporan Keuangan)
- `/console/leads` (Leads)

Semua menampilkan banner error merah: **"Tidak dapat terhubung ke gateway."** dengan data kosong (0 booking, 0 lead, 0 journal).

### Pola Error
Sama seperti ISSUE-024 (Katalog Paket), ketiga halaman ini mengalami kegagalan koneksi ke gateway saat load data. UI sudah diimplementasikan lengkap (filter, search, status tabs) tapi tidak ada data yang berhasil dimuat.

### Response Expected
- Ops Board: menampilkan daftar booking kitting berdasarkan status fulfillment
- Finance: menampilkan ringkasan akun dan riwayat jurnal transaksi
- Leads: menampilkan daftar prospek dengan status pipeline

### Petunjuk Fix
Periksa apakah gateway endpoint untuk ops/finance/leads sudah terdaftar dan dapat diakses dari SvelteKit server atau browser:
- `GET /v1/ops/fulfillment`
- `GET /v1/finance/accounts` dan `GET /v1/finance/journals`
- `GET /v1/leads`

Kemungkinan: service logistics-svc, finance-svc, atau crm-svc belum berjalan atau belum terdaftar di gateway router.

---

## [ISSUE-026] B2C: Nav links "Jadwal", "Manasik", dan "Tentang Kami" → 404 Not Found
- **Severity**: HIGH
- **Status**: FIXED
- **Backlog ID**: BL-FE-NAV-001
- **Service**: core-web
- **Dilaporkan**: 2026-04-24 (Browser E2E)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: False alarm — `MarketingTopNav.svelte` sudah menggunakan link yang valid: "Jadwal" → `href="/packages"`, "Manasik" → `href="/#proses-booking"` (anchor di landing page), "Tentang Kami" → `href="/"`. Tidak ada route `/jadwal`, `/manasik`, atau `/about` di nav — link test mungkin mengacu versi lama kode. Tidak ada perubahan kode diperlukan.
- **Verified oleh**: —

### Langkah Reproduksi
1. Buka `http://216.176.238.161/`
2. Klik salah satu dari tiga link di navigasi header B2C:
   - "Jadwal" → `/jadwal`
   - "Manasik" → `/manasik`
   - "Tentang Kami" → `/about`
3. Ketiga URL mengembalikan halaman 404 Not Found dari SvelteKit

### Response Aktual
```
GET http://216.176.238.161/jadwal     → 404 Not Found (SvelteKit "page not found")
GET http://216.176.238.161/manasik    → 404 Not Found
GET http://216.176.238.161/about      → 404 Not Found
```
Link tersebut ada dan terlihat di navigasi header, tapi route SvelteKit-nya belum dibuat.

### Response Expected
Ketiga halaman harus ada (bisa minimal placeholder) karena:
- Link sudah terlihat di navigasi — user akan klik dan mendapat 404
- Ini merupakan navigasi utama situs B2C

### Referensi
`apps/core-web/src/routes/` — tidak ada folder `jadwal/`, `manasik/`, atau `about/`

### Petunjuk Fix
Buat SvelteKit route untuk ketiga halaman tersebut. Minimal berupa halaman "Coming Soon" / placeholder:
```
apps/core-web/src/routes/jadwal/+page.svelte
apps/core-web/src/routes/manasik/+page.svelte
apps/core-web/src/routes/about/+page.svelte
```
Atau sembunyikan/disable link tersebut di nav header sampai halaman siap.

---

## [ISSUE-027] B2C: Form booking Step 2 tidak menampilkan pesan validasi saat field wajib kosong
- **Severity**: MEDIUM
- **Status**: FIXED
- **Backlog ID**: BL-FE-BOOK-002
- **Service**: core-web
- **Dilaporkan**: 2026-04-24 (Browser E2E)
- **Diperbaiki oleh**: Lutfi — 2026-04-24
- **Fix summary**: Di `booking/[package_id]/+page.svelte`: (1) Tambah state `formTouched` — di-set `true` saat user klik "Lanjut ke Review"; (2) Per-field validation messages dengan `class:inp-error` (border merah) dan `<p class="field-err">` untuk: nama pemesan, email, WhatsApp, plus nama/NIK/tanggal-lahir tiap jamaah; (3) Tambah `<div class="validation-banner">` di atas tombol submit yang muncul saat `formTouched && !step2Valid()`; (4) Asterisk `<span class="req">*</span>` di semua label field wajib; (5) Default jumlah jamaah diubah dari 2 → 1; (6) Tombol "Kembali" dan "Ubah data" dihapus dari Step 3 (one-way flow setelah Step 2).
- **Verified oleh**: —

### Langkah Reproduksi
1. Buka booking wizard: `http://216.176.238.161/booking/pkg_01JCDE00000000000000000001?departure=dep_01JCDF00000000000000000001`
2. Di Step 1 (Pilih Paket): klik "Lanjut ke Data Jamaah" tanpa memilih room type (atau biarkan default)
3. Di Step 2 (Data Jamaah): **biarkan SEMUA field kosong**:
   - Kontak Pemesan: Nama Lengkap, Email, Nomor WhatsApp
   - Data Jamaah 1: Nama sesuai KTP, NIK, Tanggal Lahir
4. Klik tombol "Lanjut ke Review"

### Response Aktual
- Navigasi ke Step 3 **tidak terjadi** (validasi bekerja secara teknis)
- **Tidak ada pesan error yang ditampilkan kepada user**:
  - Tidak ada border merah pada field yang kosong
  - Tidak ada teks error di bawah field ("Field ini wajib diisi", dll.)
  - Tidak ada toast notification / alert
  - Form tetap terlihat sama persis seperti sebelum klik tombol

### Response Expected
Saat user klik "Lanjut ke Review" dengan field wajib kosong:
- Setiap field yang kosong harus ditandai secara visual (border merah / highlight)
- Muncul pesan error yang informatif di bawah masing-masing field
- Dan/atau muncul toast/alert yang menjelaskan apa yang perlu diisi

### Petunjuk Fix
Di komponen form Step 2 (`apps/core-web/src/lib/features/s2-booking/`), tambahkan validasi UI yang memberikan feedback visual:
1. Saat submit, set state `touched = true` untuk semua field
2. Tampilkan error message kondisional untuk setiap field yang tidak valid
3. Scroll ke field pertama yang error agar user tahu harus mengisi apa

Contoh implementasi Svelte 5 (runes):
```svelte
{#if touched && !form.leadName}
  <p class="text-red-500 text-sm mt-1">Nama lengkap wajib diisi</p>
{/if}
```
