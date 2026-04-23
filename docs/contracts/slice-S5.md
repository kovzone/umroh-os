---
slice: S5
title: Slice S5 — Integration Contract
status: draft
last_updated: 2026-04-23
pr_owner: Lutfi
reviewer: Elda
task_codes:
  - S5-J-01
  - S5-J-02
---

# Slice S5 — Integration Contract

> Slice S5 = "Hardening" — UAT execution, bug triage, security & performance checklist, and platform freeze before production cut-over. This file is the wire-level agreement between the team for the quality gates, finance report endpoints, and severity classification that govern S5 completion.
>
> **Incremental build:** only sections that have a landed `S5-J-*` card are filled. Amendments after merge follow the Bump-versi rule in `docs/contracts/README.md` (Changelog append for additive, `-v2.md` for breaking).

## Purpose

S5 is the hardening slice. The flow: all S1–S4 features are code-complete → UAT scenarios are executed by Lutfi (B2C/agent journeys) and Elda (payment/finance/logistics journeys) → bugs are classified by severity and resolved within contracted SLAs → the system is declared production-ready when all MVP gates pass. S5 also introduces two read-only finance report endpoints (`GET /v1/finance/summary` and `GET /v1/finance/journals`) so the finance team can audit the double-entry journal accumulated in S3 before cut-over.

S5 does **not** introduce new booking flows, payment methods, or CRM features. It extends only `finance-svc` (two read endpoints) and adds no new gRPC services or event shapes.

## Scope

**In scope for S5 contracts:**

- `GET /v1/finance/summary` — account balance summary endpoint — `S5-J-01`.
- `GET /v1/finance/journals` — journal entry list with pagination + date filter — `S5-J-01`.
- Bug severity matrix + fix SLA table — `S5-J-02`.
- UAT mandatory scenarios reference + MVP completion gates — `S5-J-01`.

**Out of scope for S5 contracts (deferred to Phase 6 or later):**

- Finance write endpoints (AP disbursement, manual journal correction — Phase 6 `S3-E-07`).
- Revenue recognition trigger from departure event (Phase 6 `BL-FIN-004`).
- Automated performance load testing suite (Phase 6 `BL-INFRA-*`).
- Penetration testing / external security audit (post-MVP).
- Executive dashboard widgets (`S5-L-02` — Phase 6 `BL-DASH-001..005`).
- Daily prayer / qibla companion app (`S5-L-03` — Phase 6 `BL-PLG-001..009`).

---

## § S5-J-01 — Finance Report Endpoints

*(Landed with `S5-J-01`.)*

S5 adds two read-only endpoints to `finance-svc` (exposed via gateway) so that finance staff can inspect the double-entry journal built up during S3. Both endpoints are **read-only** and pull directly from the `journal_entries` and `journal_lines` tables contracted in S3 `§ S3-J-03`. No writes, no side effects.

### Background: data source

Finance data is stored in `finance-svc` using the schema contracted in S3 `§ S3-J-03`:

```
journal_entries  (id, idempotency_key, created_at)
journal_lines    (id, journal_entry_id, account_code, debit_idr, credit_idr)
```

The report endpoints aggregate these tables. `net_idr` is always computed as `total_debit_idr - total_credit_idr` — a positive net means the account has a net debit balance (asset/expense convention); a negative net means a net credit balance (liability/revenue convention).

All monetary values are **integer IDR** (no decimals, no currency conversion). This is consistent with S2 `§ S2-J-01` and S3 `§ S3-J-03`.

---

### `GET /v1/finance/summary` — Account balance summary

**Purpose:** return the aggregated debit, credit, and net balance for every account code that appears in `journal_lines`. Used by finance staff to verify the trial balance before production cut-over.

**Auth:** `Authorization: Bearer <token>` — IAM JWT with role `staff` or `finance`. Gateway enforces `RequirePermission("finance:read")` per ADR-0009 / S1-E-14.

**Request:** no query parameters. Always returns the full trial balance across all time.

**Response `200 OK`:**

```json
{
  "generated_at": "2026-04-23T15:00:00Z",
  "accounts": [
    {
      "code": "1001",
      "name": "Bank",
      "total_debit_idr": 50000000,
      "total_credit_idr": 0,
      "net_idr": 50000000
    },
    {
      "code": "2001",
      "name": "Pilgrim Liability",
      "total_debit_idr": 0,
      "total_credit_idr": 50000000,
      "net_idr": -50000000
    }
  ]
}
```

**Response field definitions:**

| Field | Type | Notes |
| --- | --- | --- |
| `generated_at` | string (RFC3339) | Server-side timestamp at the moment the query was executed. Not cached. |
| `accounts` | array | One entry per distinct `account_code` that appears in at least one `journal_lines` row. Ordered by `code` ASC. |
| `accounts[].code` | string | Account code, e.g. `"1001"`. Matches `journal_lines.account_code`. |
| `accounts[].name` | string | Human-readable account name from the hardcoded MVP Chart of Accounts (contracted in S3 `§ S3-J-03`). |
| `accounts[].total_debit_idr` | integer | `SUM(debit_idr)` across all journal lines for this account. Always ≥ 0. |
| `accounts[].total_credit_idr` | integer | `SUM(credit_idr)` across all journal lines for this account. Always ≥ 0. |
| `accounts[].net_idr` | integer | `total_debit_idr - total_credit_idr`. May be negative. |

**Computed field rule — `net_idr`:**

```
net_idr = total_debit_idr - total_credit_idr
```

A positive `net_idr` indicates a net debit balance (normal for asset / expense accounts: `1001 Bank`). A negative `net_idr` indicates a net credit balance (normal for liability / revenue accounts: `2001 Pilgrim Liability`, `4001 Deferred Revenue`).

**Failure codes:**

| HTTP | `error.code` | When |
| --- | --- | --- |
| `401` | `unauthorized` | Missing or invalid bearer token |
| `403` | `forbidden` | Valid token but role lacks `finance:read` permission |
| `500` | `internal_error` | DB failure or unexpected error |

**Implementation note (for Elda — `S5-E-01`):** the query is a single `GROUP BY account_code` aggregation on `journal_lines`. There is no pagination — the MVP Chart of Accounts has at most ~10 codes. If the `journal_lines` table grows beyond 10 M rows before Phase 6, add a materialized view refresh (out of scope for S5).

---

### `GET /v1/finance/journals` — Journal entry list

**Purpose:** return a paginated, date-filtered list of journal entries with their lines. Used by finance staff to audit individual double-entry postings.

**Auth:** `Authorization: Bearer <token>` — IAM JWT with role `staff` or `finance`. Gateway enforces `RequirePermission("finance:read")`.

**Query parameters:**

| Parameter | Type | Required | Default | Notes |
| --- | --- | --- | --- | --- |
| `from` | string (YYYY-MM-DD) | no | (no lower bound) | Filter `journal_entries.created_at >= from` (inclusive, interpreted as `T00:00:00Z` UTC). |
| `to` | string (YYYY-MM-DD) | no | (no upper bound) | Filter `journal_entries.created_at <= to` (inclusive, interpreted as `T23:59:59Z` UTC). |
| `limit` | integer | no | `50` | Number of entries per page. Max `200`. Clamped silently to 200 if caller provides a higher value. |
| `cursor` | string | no | (first page) | Opaque pagination cursor returned in the previous response's `next_cursor` field. Must not be constructed by the client — treat as opaque. |

**Response `200 OK`:**

```json
{
  "entries": [
    {
      "id": "je_01JCDE...",
      "idempotency_key": "payment:inv_01JCDE...",
      "created_at": "2026-04-23T15:00:00Z",
      "lines": [
        { "account_code": "1001", "debit_idr": 10000000, "credit_idr": 0 },
        { "account_code": "2001", "debit_idr": 0, "credit_idr": 10000000 }
      ]
    }
  ],
  "next_cursor": null
}
```

**Response field definitions:**

| Field | Type | Notes |
| --- | --- | --- |
| `entries` | array | Journal entries matching the filter, ordered by `created_at` DESC. |
| `entries[].id` | string | Journal entry ID. ULID with `je_` prefix. Matches `journal_entries.id` from S3. |
| `entries[].idempotency_key` | string | Idempotency key used when posting the entry (e.g. `"payment:inv_01JCDE..."`). Matches `journal_entries.idempotency_key`. |
| `entries[].created_at` | string (RFC3339) | Timestamp of journal entry creation. |
| `entries[].lines` | array | All `journal_lines` rows for this entry. Minimum 2 rows (double-entry invariant). |
| `entries[].lines[].account_code` | string | Account code, e.g. `"1001"`. |
| `entries[].lines[].debit_idr` | integer | Debit amount in IDR. 0 if this line is a credit line. |
| `entries[].lines[].credit_idr` | integer | Credit amount in IDR. 0 if this line is a debit line. |
| `next_cursor` | string or null | Opaque cursor for the next page. `null` if this is the last page. |

**Ordering:** always `created_at` DESC (newest entries first). Within the same `created_at` timestamp, ordering is by `id` DESC (ULID sort is stable).

**Pagination:** cursor-based. The cursor encodes the `(created_at, id)` of the last entry on the current page. `finance-svc` uses a keyset pagination query:

```sql
WHERE (created_at, id) < (:last_created_at, :last_id)
ORDER BY created_at DESC, id DESC
LIMIT :limit
```

The cursor value is an opaque base64-encoded string — clients MUST NOT parse or construct it.

**Date filter behavior:**

- `from` only: return all entries from that date to the most recent.
- `to` only: return all entries up to and including that date.
- Both: return entries in the inclusive date range `[from, to]`.
- Neither: return all entries (paginated).

**Failure codes:**

| HTTP | `error.code` | When |
| --- | --- | --- |
| `400` | `invalid_date_format` | `from` or `to` is not a valid `YYYY-MM-DD` string |
| `400` | `invalid_date_range` | `from` is after `to` |
| `400` | `invalid_cursor` | `cursor` value is malformed or tampered |
| `400` | `invalid_limit` | `limit` is not a positive integer |
| `401` | `unauthorized` | Missing or invalid bearer token |
| `403` | `forbidden` | Valid token but role lacks `finance:read` permission |
| `500` | `internal_error` | DB failure or unexpected error |

### Honored by implementation

- `S5-E-01` (`BL-FIN-005`) — `finance-svc`: `GET /v1/finance/summary` aggregation query + gateway route + `RequirePermission("finance:read")`.
- `S5-E-01` (`BL-FIN-006`) — `finance-svc`: `GET /v1/finance/journals` keyset pagination + date filter + gateway route.

---

## § S5-J-02 — Bug Severity Matrix

*(Landed with `S5-J-02`.)*

All bugs discovered during UAT are classified using the four-level severity matrix below. The label (`P0`–`P3`) is applied at triage and governs the fix SLA and merge path.

| Severity | Label | Definisi | Fix SLA |
| --- | --- | --- | --- |
| P0 | Critical | Data loss, security breach, payment error, system down (service unreachable), double-entry imbalance, incorrect VA amount, seat count corruption | Fix dalam **4 jam** jam kerja; hotfix branch langsung ke `main`; post-mortem wajib |
| P1 | High | Fitur utama broken tapi ada workaround; data tidak akurat tapi tidak hilang; booking stuck di status salah; assignment CS salah tapi bisa diperbaiki manual | Fix dalam **1 hari kerja**; PR ke `main` via review normal |
| P2 | Medium | Fitur sekunder error; UI mismatch vs desain; UX buruk tapi tidak menghalangi alur utama; pesan error tidak informatif; pagination off-by-one | Fix dalam **sprint berikutnya**; PR normal |
| P3 | Low | Typo; kosmetik; label/terjemahan kurang tepat; improvement minor yang tidak memengaruhi fungsi | Masuk **backlog**; prioritas rendah; boleh di-batch |

**Triage rules:**

1. Any bug that triggers a data inconsistency in `journal_entries` / `journal_lines` (debit ≠ credit on any entry) is automatically **P0**, regardless of how it was discovered.
2. Any bug that allows a payment to be accepted by the gateway but not reflected in `invoices` or `bookings` is automatically **P0**.
3. Security findings (unauthenticated access to bearer-protected endpoints, IDOR, exposed secrets) are automatically **P0**.
4. P0 fixes bypass the normal sprint cycle — the hotfix is branched from `main`, reviewed by both Lutfi and Elda, and merged directly to `main` with a post-merge smoke run.
5. P1 and above bugs block the S5 "done" gate — see `§ UAT` below.

---

## § UAT

*(Landed with `S5-J-01`.)*

### UAT checklist reference

The full UAT scenario checklist is maintained in:

```
docs/UAT-Checklist-S1-S4.docx
```

This document contains one scenario row per critical user journey, covering:

- **S1** — catalog browse, departure detail, booking draft creation (public B2C + staff console).
- **S2** — invoice issuance, VA display, mock Midtrans/Xendit webhook, status transitions to `paid`.
- **S3** — fulfillment task creation on `paid_in_full`, double-entry journal posting (Dr 1001 Bank / Cr 2001 Pilgrim Liability), ops board display.
- **S4** — lead capture with UTM params, CS round-robin assignment, lead status transitions, lead→booking conversion, `OnBookingCreated` and `OnBookingPaidInFull` CRM events.

UAT **wajib dijalankan sebelum S5 dinyatakan selesai**. Neither Lutfi nor Elda may mark S5 complete until all MVP gates below are satisfied.

### MVP completion gates

S5 is declared **DONE** only when **all** of the following are true:

1. **Semua P0 test case LULUS** — every scenario in `docs/UAT-Checklist-S1-S4.docx` that is marked as P0 criticality has a ✅ result recorded by the executing tester.
2. **Tidak ada open P0 bug** — zero bugs in the P0 severity bucket remain open in the issue tracker.
3. **Tidak ada open P1 bug** — zero bugs in the P1 severity bucket remain open (P1 bugs must be resolved, not deferred).
4. **Finance trial balance verified** — `GET /v1/finance/summary` returns a balanced trial balance (sum of all `net_idr` across all accounts = 0, within the test seed data).
5. **Journal integrity verified** — `GET /v1/finance/journals` returns entries where each entry's lines satisfy `SUM(debit_idr) = SUM(credit_idr)` for at least the mock payment scenarios executed in UAT.

P2 and P3 bugs discovered during UAT do **not** block S5 completion — they are triaged into the Phase 6 backlog.

### UAT owners

| Journey area | UAT executor | Task code |
| --- | --- | --- |
| B2C booking + agent lead capture | Lutfi | `S5-L-01` |
| Payment webhook + finance journal + logistics tasks | Elda | `S5-E-01` |

---

## § UAT Scenario Checklists (BL-QA-002 / BL-QA-003)

*(Landed with `S2-J-05` / `S5-J-02`. Evidence rows correspond to `docs/UAT-Checklist-S1-S4.docx`.)*

These checklists define the **minimum passing criteria** for S5 to be declared done. Each row is a scenario the designated executor runs manually (or via the Playwright E2E suite in `tests/e2e/tests/`). Status is recorded in `docs/UAT-Checklist-S1-S4.docx` — this document is the *definition*; the `.docx` is the *evidence*.

**Legend:** ✅ Pass | ❌ Fail | ⏭ Skip (reason required) | 🔁 Re-test after fix

---

### BL-QA-002 — Core checklist: Auth, Catalog, Booking, Permission & Audit

*Executor: Lutfi | Task code: `S5-L-01`*

| # | Scenario | Endpoint / Action | Expected result | Priority |
|---|---|---|---|---|
| C-01 | Login valid (admin) | `POST /v1/auth/login` | 200, PASETO `v2.local.*` access token + 64-char refresh token | P0 |
| C-02 | Login invalid password | `POST /v1/auth/login` | 401 `invalid_credentials` | P0 |
| C-03 | Refresh session | `POST /v1/auth/refresh` | 200, new access + refresh token pair | P0 |
| C-04 | Logout (invalidate token) | `DELETE /v1/auth/logout` | 204; subsequent `/v1/me` with old token → 401 | P0 |
| C-05 | `GET /v1/me` with valid token | Bearer auth | 200, correct `user_id`, `email`, `roles` | P0 |
| C-06 | `GET /v1/me` with no token | No Authorization header | 401 `unauthorized` | P0 |
| C-07 | Browse catalog public (no auth) | `GET /v1/packages` | 200, at least 1 active package; draft/archived excluded | P0 |
| C-08 | Package detail — active | `GET /v1/packages/{id}` | 200, pricing array non-empty, departures array present | P0 |
| C-09 | Package detail — draft → 404 | `GET /v1/packages/{draft_id}` | 404 `package_not_found` | P0 |
| C-10 | Draft booking (B2C, no auth) | `POST /v1/bookings` | 201/200, `status: "draft"`, booking_id returned | P0 |
| C-11 | Booking without required fields | `POST /v1/bookings` (missing email) | 400/422, error object present | P1 |
| C-12 | Booking with cancelled departure | `POST /v1/bookings` | 400/404, booking rejected | P1 |
| C-13 | Create package (staff) | `POST /v1/packages` with valid token | 201, package returned | P0 |
| C-14 | Create package (no permission) | `POST /v1/packages` with non-staff token | 403 `forbidden` | P0 |
| C-15 | Suspend user (super_admin) | `POST /v1/users/{id}/suspend` | 200, user status becomes `suspended` | P1 |
| C-16 | Suspended user cannot login | `POST /v1/auth/login` | 401 or 403 `account_suspended` | P0 |
| C-17 | Audit log created on state change | DB query `iam.audit_logs` after C-13 | At least 1 audit row for `resource=package action=create` | P1 |
| C-18 | Permission gate: `catalog.package.manage` | Staff with correct permission creates; staff without → 403 | Permission enforced correctly | P0 |
| C-19 | 2FA TOTP enroll | `POST /v1/me/2fa/enroll` | 200, TOTP URI returned | P2 |
| C-20 | TOTP verify | `POST /v1/me/2fa/verify` with valid OTP | 200 `verified: true` | P2 |

**Regression scenarios (permission + audit trail):**

| # | Regression | What to verify |
|---|---|---|
| R-01 | All protected routes require bearer | Hit each `v1Protected` route without token → every one returns 401 |
| R-02 | Audit log not writable by client | `POST /iam/audit_logs` (if such route existed) should be 404 or 403 |
| R-03 | `RequireBearerToken` middleware blocks all protected groups | Manually hit `/v1/me`, `/v1/leads`, `/v1/finance/summary` without auth → all 401 |
| R-04 | Token from service A not accepted by service B | Token issued by iam-svc validated by gateway middleware correctly |

---

### BL-QA-003 — Payment, Finance & Logistics checklist

*Executor: Elda | Task code: `S5-E-01`*

#### Payment (S2)

| # | Scenario | Endpoint / Action | Expected result | Priority |
|---|---|---|---|---|
| P-01 | Issue VA (mock gateway) | `POST /v1/invoices` + `POST /v1/invoices/{id}/virtual-accounts` | 200/201; `account_number`, `bank_code`, `amount`, `expires_at` present | P0 |
| P-02 | VA idempotent on same booking | Repeat `POST /v1/invoices` same `booking_id` | Same `invoice_id` returned, no duplicate invoice | P0 |
| P-03 | Mock webhook: paid full amount | `POST /v1/webhooks/mock/trigger` `status=paid amount=full` | 200; booking transitions to `paid_in_full` within 1 s | P0 |
| P-04 | Mock webhook: underpayment | `POST /v1/webhooks/mock/trigger` `amount < invoice.amount` | Booking stays `pending_payment` | P0 |
| P-05 | Duplicate webhook idempotent | Same `invoice_id` trigger twice | 200 both times; only 1 row in `payment.payment_events` | P0 |
| P-06 | Midtrans webhook without auth header | `POST /v1/webhooks/midtrans` no `X-Callback-Token` | 401/403 | P0 |
| P-07 | ReissuePaymentLink (CS closing) | `POST /v1/payments/link` with `booking_id` + bearer | 200; new VA issued if old expired (`is_new: true`); idempotent if still active | P1 |
| P-08 | ReissuePaymentLink — booking not found | `POST /v1/payments/link` bad `booking_id` | 404 `booking_not_found` | P1 |
| P-09 | VA already paid → reissue returns existing | `POST /v1/payments/link` on `paid_in_full` booking | 200, `is_new: false`, same VA details | P1 |

#### Finance (S3 + S5)

| # | Scenario | Endpoint / Action | Expected result | Priority |
|---|---|---|---|---|
| F-01 | Journal posted after payment | DB query `finance.journal_entries` after P-03 | At least 1 entry with `idempotency_key LIKE 'payment:%'` | P0 |
| F-02 | Double-entry balance invariant | `SELECT SUM(debit_idr) - SUM(credit_idr) FROM finance.journal_lines WHERE journal_entry_id = {id}` | Result = 0 for every entry | P0 |
| F-03 | Finance summary endpoint | `GET /v1/finance/summary` (with finance-role bearer) | 200; `accounts` array non-empty; every `net_idr = total_debit_idr - total_credit_idr` | P0 |
| F-04 | Finance summary: trial balance | Sum all `net_idr` across all accounts after seeded test runs | Total ≈ 0 (balanced ledger) | P0 |
| F-05 | Journal list — first page | `GET /v1/finance/journals` | 200; `entries` non-empty; each entry has ≥2 lines | P0 |
| F-06 | Journal list — date filter | `GET /v1/finance/journals?from=2026-01-01&to=2026-12-31` | Only entries within range returned | P1 |
| F-07 | Journal list — invalid date format | `GET /v1/finance/journals?from=bad-date` | 400 `invalid_date_format` | P1 |
| F-08 | Journal delete is forbidden | `DELETE /v1/finance/journals/{id}` (any bearer) | 403 `forbidden` — use CorrectJournal instead | P0 |
| F-09 | CorrectJournal: reversal posted | `POST /v1/finance/journals/{id}/correct` | 200; `correction_entry_id` returned; new journal lines reverse Dr/Cr of original | P0 |
| F-10 | CorrectJournal: idempotent second call | Same `{id}` corrected twice | 200 both times; `idempotent: true` on second call; only 1 correction entry in DB | P0 |
| F-11 | CorrectJournal: not found | `POST /v1/finance/journals/je_nonexistent/correct` | 404 `not_found` | P1 |
| F-12 | Finance summary: unauthorized | `GET /v1/finance/summary` no token | 401 | P0 |
| F-13 | P&L report | `GET /v1/finance/pl` with bearer | 200; `revenues`, `expenses`, `net_profit` fields present | P1 |
| F-14 | Balance sheet | `GET /v1/finance/balance-sheet` with bearer | 200; `assets`, `liabilities`, `equity` fields present | P1 |
| F-15 | GRN auto-journal | `POST /v1/finance/grn` → DB `finance.journal_entries` | Entry posted with `idempotency_key LIKE 'grn:%'` | P1 |

#### Logistics (S3)

| # | Scenario | Endpoint / Action | Expected result | Priority |
|---|---|---|---|---|
| L-01 | Ship fulfillment task | `POST /v1/logistics/ship` (bearer, paid booking) | 200; `shipment_id` returned; carrier + tracking_number present | P1 |
| L-02 | Generate pickup QR | `POST /v1/logistics/pickup-qr` (bearer) | 200; `qr_url` or `qr_token` returned | P1 |
| L-03 | Redeem pickup QR | `POST /v1/logistics/pickup-qr/redeem` (valid token) | 200; task marked as collected | P1 |
| L-04 | Redeem QR: invalid token | `POST /v1/logistics/pickup-qr/redeem` (bad token) | 400/401 `invalid_token` | P1 |
| L-05 | Room allocation | `POST /v1/ops/room-allocation` (bearer) | 200; rooms assigned; no conflict | P1 |
| L-06 | Export manifest | `GET /v1/ops/manifest/{departure_id}/export` | 200; downloadable; contains jamaah rows | P1 |

#### Finance / Payment end-to-end trace

This composite scenario verifies the full chain: booking → invoice → VA → webhook → journal → summary.

| Step | Action | Verify |
|---|---|---|
| E2E-01 | Create draft booking | `status: "draft"` |
| E2E-02 | Issue VA (mock) | `account_number` non-empty, `expires_at` in future |
| E2E-03 | Trigger mock webhook `paid` | Booking `status: "paid_in_full"` |
| E2E-04 | Confirm journal entry exists | `finance.journal_entries` has `payment:{invoice_id}` |
| E2E-05 | Confirm double-entry balance | `SUM(debit) = SUM(credit)` for entry from E2E-04 |
| E2E-06 | Pull finance summary | `1001 Bank net_idr > 0`; `2001 Pilgrim Liability net_idr < 0` |

This trace is **automated** in `tests/e2e/tests/05-uat-s2-booking-payment.spec.ts` and must pass before S5 is declared done.

---

## § Error envelope (S5 additions)

All S5 REST error responses use the shared envelope established in S1–S4:

```json
{
  "error": {
    "code": "<snake_case>",
    "message": "<human-readable, id-ID>",
    "trace_id": "<otel_span_hex>"
  }
}
```

`trace_id` is the OTel span ID per `docs/04-backend-conventions/03-logging-and-tracing.md`.

---

## § Changelog

- **2026-04-23** — Initial version drafted via tasks `S5-J-01` and `S5-J-02`. Covers: `GET /v1/finance/summary` account balance aggregation endpoint (S5-J-01); `GET /v1/finance/journals` paginated journal entry list with date filter and keyset cursor pagination (S5-J-01); bug severity matrix P0–P3 with fix SLA and triage rules (S5-J-02); UAT completion gates referencing `docs/UAT-Checklist-S1-S4.docx` (S5-J-01).
- **2026-04-23** — UAT scenario checklists added (BL-QA-002 / BL-QA-003): core auth/catalog/booking/permission/audit regression matrix (C-01..C-20, R-01..R-04); payment/finance/logistics pass criteria (P-01..P-09, F-01..F-15, L-01..L-06); end-to-end payment-to-journal trace (E2E-01..E2E-06).
