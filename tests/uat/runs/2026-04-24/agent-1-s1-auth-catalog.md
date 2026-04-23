# UAT Agent 1 — Auth + Catalog (S1) — 2026-04-24

## Summary
- **Total test cases**: 34 (26 API + 8 UI)
- **PASS**: 0
- **FAIL**: 6 (root test di masing-masing serial describe — cascade ke 28 skipped)
- **NOT_DEPLOYED**: 0
- **SKIP / Did Not Run**: 28 (dependency dari test pertama yang fail di setiap describe.serial)
- **Issues dibuat**: ISSUE-001, ISSUE-002, ISSUE-003, ISSUE-004

## Root Cause Utama
**Semua test gagal karena infrastruktur server tidak dapat dijangkau dari sandbox QA ini:**
1. **Gateway API port 4000** → `ENETUNREACH` — tidak dapat connect sama sekali
2. **Core Web port 80** → `ERR_EMPTY_RESPONSE` — TCP terhubung tapi server mengembalikan respons kosong
3. **PostgreSQL port 5432** → `ENETUNREACH` — tidak dapat connect

Ini adalah masalah infrastruktur/deployment, bukan bug logika bisnis. Semua test API otomatis fail karena tidak bisa connect ke `216.176.238.161:4000`.

## Hasil Per Test Case

### API Tests (project: uat-api)

| # | Test ID | Skenario | Status | Notes |
|---|---------|----------|--------|-------|
| 1 | S1-AUTH-01 | login valid → 200 + PASETO token | ❌ FAIL | ENETUNREACH :4000 → ISSUE-001 |
| 2 | S1-AUTH-02 | login password salah → 401 | ⏭️ SKIP | serial — dependen S1-AUTH-01 |
| 3 | S1-AUTH-03 | login email tidak dikenal → 401 | ⏭️ SKIP | serial |
| 4 | S1-AUTH-04 | GET /v1/me dengan token valid → 200 | ⏭️ SKIP | serial |
| 5 | S1-AUTH-05 | GET /v1/me tanpa token → 401 | ⏭️ SKIP | serial |
| 6 | S1-AUTH-06 | GET /v1/me token garbage → 401 | ⏭️ SKIP | serial |
| 7 | S1-AUTH-07 | POST /v1/sessions/refresh → token baru | ⏭️ SKIP | serial |
| 8 | S1-AUTH-08 | iam-svc down → gateway 502 | ⏭️ SKIP | test.skip (by design) |
| 9 | S1-PERM-01 | route staff tanpa token → 401 | ❌ FAIL | ENETUNREACH :4000 → ISSUE-001 |
| 10 | S1-PERM-02 | logout dengan token valid → 200 | ⏭️ SKIP | serial |
| 11 | S1-PERM-03 | token revoked tidak bisa dipakai → 401 | ⏭️ SKIP | serial |
| 12 | S1-CAT-01 | GET /v1/packages → hanya active | ❌ FAIL | ENETUNREACH :4000 → ISSUE-001 |
| 13 | S1-CAT-02 | GET /v1/packages/{active_id} → detail | ⏭️ SKIP | serial |
| 14 | S1-CAT-03 | GET /v1/packages/{draft_id} → 404 | ⏭️ SKIP | serial |
| 15 | S1-CAT-04 | GET /v1/package-departures/{id} → detail | ⏭️ SKIP | serial |
| 16 | S1-CAT-05 | GET /v1/package-departures/{cancelled} → 404 | ⏭️ SKIP | serial |
| 17 | S1-CAT-06 | GET /v1/packages tanpa auth → 200 | ⏭️ SKIP | serial |
| 18 | S1-WRITE-01 | POST /v1/packages (staff) → 201 | ❌ FAIL | ENETUNREACH :4000 (login) + :5432 (cleanup) → ISSUE-001, ISSUE-003 |
| 19 | S1-WRITE-02 | POST /v1/packages tanpa auth → 401 | ⏭️ SKIP | serial |
| 20 | S1-WRITE-03 | POST /v1/packages tanpa nama → 422 | ⏭️ SKIP | serial |
| 21 | S1-WRITE-04 | PATCH /v1/packages/{id} → update | ⏭️ SKIP | serial |
| 22 | S1-WRITE-05 | POST departures → 201 | ⏭️ SKIP | serial |
| 23 | S1-WRITE-06 | departure tgl invalid → 422 | ⏭️ SKIP | serial |
| 24 | S1-WRITE-07 | PATCH departure status → closed | ⏭️ SKIP | serial |
| 25 | S1-WRITE-08 | departure draft tidak muncul publik | ⏭️ SKIP | serial |
| 26 | S1-WRITE-09 | DELETE /v1/packages → 200 | ⏭️ SKIP | serial |

### Browser / UI Tests (project: browser)

| # | Test ID | Skenario | Status | Notes |
|---|---------|----------|--------|-------|
| 27 | S1-UI-01 | /console/login tampil dengan form login | ❌ FAIL | ERR_EMPTY_RESPONSE port 80 → ISSUE-002 |
| 28 | S1-UI-02 | login invalid → error message tampil | ⏭️ SKIP | serial |
| 29 | S1-UI-03 | login valid → redirect ke /console | ⏭️ SKIP | serial |
| 30 | S1-UI-04 | console shell sidemenu visible | ⏭️ SKIP | serial |
| 31 | S1-UI-05 | / landing page + CTA browse packages | ❌ FAIL | ERR_EMPTY_RESPONSE port 80 → ISSUE-002 |
| 32 | S1-UI-06 | /packages — daftar package tampil | ⏭️ SKIP | serial |
| 33 | S1-UI-07 | klik package → halaman detail | ⏭️ SKIP | serial |
| 34 | S1-UI-08 | halaman detail → tombol booking visible | ⏭️ SKIP | serial |

## Temuan Penting

### 1. Server tidak dapat dijangkau (CRITICAL)
Gateway API di `http://216.176.238.161:4000` mengembalikan `ENETUNREACH` untuk semua request. Ini bisa berarti:
- Container `gateway-svc` tidak berjalan
- Port 4000 tidak di-expose / firewall memblokir
- Server sedang down

Sampai ISSUE-001 diperbaiki, **tidak ada satupun test API yang bisa dijalankan**.

### 2. Core Web tidak merespons (CRITICAL)
`http://216.176.238.161` (port 80) mengembalikan `ERR_EMPTY_RESPONSE` — artinya TCP connection berhasil (beda dengan ENETUNREACH) tapi server menutup koneksi tanpa mengirim HTTP response. Kemungkinan:
- Nginx berjalan tapi misconfigured (tidak ada server block yang handle)
- SvelteKit app (`core-web`) belum di-build atau belum di-start
- Nginx belum di-reload setelah config berubah

### 3. Config Playwright — project `uat-api` tidak punya `baseURL` (MEDIUM)
Spec file `04-uat-s1-auth-catalog.spec.ts` berisi gabungan API tests DAN UI tests dalam satu file. Project `uat-api` di `playwright.config.ts` tidak memiliki `baseURL`, sehingga UI tests yang menggunakan relative URL (`page.goto("/console/login")`) gagal dengan "Cannot navigate to invalid URL". Detail di ISSUE-004.

## Langkah Selanjutnya untuk Dev (Elda/Lutfi)

1. **Elda**: Pastikan `gateway-svc`, `iam-svc`, `catalog-svc` container berjalan di server 216.176.238.161
   ```bash
   ssh user@216.176.238.161
   cd /path/to/umroh-os
   docker compose ps
   docker compose logs gateway-svc --tail=50
   ```

2. **Lutfi**: Pastikan `core-web` build dan Nginx berjalan
   ```bash
   docker compose logs core-web --tail=50
   nginx -t && nginx -s reload
   ```

3. **Lutfi**: Fix `playwright.config.ts` — tambahkan `baseURL` ke project `uat-api` (lihat ISSUE-004)

4. Setelah services running, minta QA Agent 1 re-run UAT ini untuk verifikasi.

## Test Data Dibuat
Tidak ada — semua tests gagal sebelum data dapat dibuat.

---

## Output Playwright Lengkap

### Run 1: project=uat-api

```
Running 34 tests using 1 worker

  ✘   1 [uat-api] S1 Auth — Login & Token (BL-IAM-001 / BL-IAM-018) › S1-AUTH-01 (22ms)
  -   2 [uat-api] S1-AUTH-02
  -   3 [uat-api] S1-AUTH-03
  -   4 [uat-api] S1-AUTH-04
  -   5 [uat-api] S1-AUTH-05
  -   6 [uat-api] S1-AUTH-06
  -   7 [uat-api] S1-AUTH-07
  -   8 [uat-api] S1-AUTH-08
  ✘   9 [uat-api] S1 Auth — Permissions & Protected Routes (BL-IAM-002) › S1-PERM-01 (20ms)
  -  10 [uat-api] S1-PERM-02
  -  11 [uat-api] S1-PERM-03
  ✘  12 [uat-api] S1 Catalog — Public Read (BL-CAT-001..004) › S1-CAT-01 (24ms)
  -  13..17 [uat-api] S1-CAT-02..06
  ✘  18 [uat-api] S1 Catalog — Staff Write (BL-CAT-014) › S1-WRITE-01 (0ms)
  -  19..26 [uat-api] S1-WRITE-02..09
  ✘  27 [uat-api] S1 UI — Console Login (BL-FE-CONSOLE-001) › S1-UI-01 (287ms)
  -  28..30 [uat-api] S1-UI-02..04
  ✘  31 [uat-api] S1 UI — B2C Catalog (BL-FE-BOOK-001..003) › S1-UI-05 (285ms)
  -  32..34 [uat-api] S1-UI-06..08

ERRORS:
1) S1-AUTH-01: Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
   → POST http://216.176.238.161:4000/v1/sessions

2) S1-PERM-01: Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
   → POST http://216.176.238.161:4000/v1/packages

3) S1-CAT-01: Error: apiRequestContext.get: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
   → GET http://216.176.238.161:4000/v1/packages

4) S1-WRITE-01:
   Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
   → POST http://216.176.238.161:4000/v1/sessions
   ALSO: Error: connect ENETUNREACH 216.176.238.161:5432 - Local (0.0.0.0:0)
   (cleanup via pg also failed)

5) S1-UI-01 (under uat-api project):
   Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
   navigating to "/console/login" — no baseURL configured in uat-api project

6) S1-UI-05 (under uat-api project):
   Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
   navigating to "/" — no baseURL configured in uat-api project

6 failed, 28 did not run
```

### Run 2: project=browser

```
Running 34 tests using 1 worker

  ✘   1 [browser] S1 Auth — Login & Token › S1-AUTH-01 (23ms)
  -   2..8 [browser] S1-AUTH-02..08
  ✘   9 [browser] S1 Auth — Permissions › S1-PERM-01 (19ms)
  -  10..11 [browser] S1-PERM-02..03
  ✘  12 [browser] S1 Catalog — Public Read › S1-CAT-01 (102ms)
  -  13..17 [browser] S1-CAT-02..06
  ✘  18 [browser] S1 Catalog — Staff Write › S1-WRITE-01 (0ms)
  -  19..26 [browser] S1-WRITE-02..09
  ✘  27 [browser] S1 UI — Console Login › S1-UI-01 (1.2s)
  -  28..30 [browser] S1-UI-02..04
  ✘  31 [browser] S1 UI — B2C Catalog › S1-UI-05 (1.1s)
  -  32..34 [browser] S1-UI-06..08

ERRORS:
1) S1-AUTH-01: Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)

2) S1-PERM-01: Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)

3) S1-CAT-01: Error: apiRequestContext.get: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)

4) S1-WRITE-01:
   Error: apiRequestContext.post: connect ENETUNREACH 216.176.238.161:4000 - Local (0.0.0.0:0)
   ALSO: Error: connect ENETUNREACH 216.176.238.161:5432 - Local (0.0.0.0:0)

5) S1-UI-01:
   Error: page.goto: net::ERR_EMPTY_RESPONSE at http://216.176.238.161/console/login
   (TCP connected, server returned empty response — core-web tidak merespons)

6) S1-UI-05:
   Error: page.goto: net::ERR_EMPTY_RESPONSE at http://216.176.238.161/
   (TCP connected, server returned empty response — core-web tidak merespons)

6 failed, 28 did not run
```

---
*Ditulis oleh: QA Agent 1 — 2026-04-24*
