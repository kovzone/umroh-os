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
