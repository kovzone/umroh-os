---
slice: S6
title: Slice S6 — Phase 6 MH-V1 Integration Contract
status: draft
last_updated: 2026-04-23
pr_owner: Elda
reviewer: Lutfi
task_codes:
  - S6-E-01
  - S6-E-02
  - S6-E-03
  - S6-E-04
  - S6-E-05
  - S6-L-01
---

# Slice S6 — Phase 6 MH-V1 Integration Contract

> Slice S6 = "MH-V1 Depth" — visa pipeline, field operations scanning, warehouse & procurement, finance AP/AR aging, and IAM admin/security depth. This file is the wire-level agreement between services for the Phase 6 domain expansions.
>
> **Incremental build:** only sections that have a landed `S6-*` card are filled. Amendments after merge follow the Bump-versi rule in `docs/contracts/README.md` (Changelog append for additive, `-v2.md` for breaking).

## Purpose

S6 introduces five domain expansions across new and existing services, all built on the same ADR-0009 pattern: all backends are pure gRPC; REST exposed only via `gateway-svc:4000`. The five domains are:

1. **§ 6.E — Visa Pipeline** (`visa-svc`, new service) — status-machine-driven visa application lifecycle from document collection through e-visa issuance, including bulk embassy submission.
2. **§ 6.D — Field Operations Scan** (`ops-svc`, extended) — idempotent scan event recording and bus boarding roster management for field operations staff.
3. **§ 6.C — Warehouse & Procurement** (`logistics-svc`, extended) — purchase request lifecycle with budget gate, goods receipt with QC pass/fail, and idempotent kit assembly.
4. **§ 6.B — Finance AP/AR Depth** (`finance-svc`, extended) — accounts-payable disbursement batches and AR/AP aging bucket report.
5. **§ 6.M — IAM Admin/Security Depth** (`iam-svc`, extended) — data visibility scopes, API key management, global config store, and paginated activity log search.

Several Phase 6 items are explicitly deferred (see **§ Deferred**).

S6 does **not** change existing S1–S5 REST routes or gRPC signatures. All S6 additions are strictly additive.

## Scope

**In scope for S6 contracts:**

- `visa-svc` new service: `TransitionStatus`, `BulkSubmit`, `GetApplications` RPCs + gateway routes — `S6-E-01`.
- `ops-svc` extensions: `RecordScan`, `RecordBusBoarding`, `GetBoardingRoster` RPCs + gateway routes — `S6-E-02`.
- `logistics-svc` extensions: `CreatePurchaseRequest`, `ApprovePurchaseRequest`, `RecordGRNWithQC`, `CreateKitAssembly` RPCs + gateway routes — `S6-E-03`.
- `finance-svc` extensions: `CreateDisbursementBatch`, `ApproveDisbursement`, `GetARAPAging` RPCs + gateway routes — `S6-E-04`.
- `iam-svc` extensions: `SetDataScope`, `CreateAPIKey`, `RevokeAPIKey`, `GetGlobalConfig`, `SetGlobalConfig`, `SearchActivityLog` RPCs + gateway routes — `S6-E-05`.

**Out of scope for S6 contracts (deferred — see § Deferred):**

- Agent onboarding with WA Business API + KYC provider (6.A).
- Commission wallet / CRM ledger extension (6.A BL-CRM-012).
- B2C frontend / SvelteKit work (6.L) — dedicated frontend sprint.
- Dashboard/executive widgets (6.F, 6.I).
- CRM depth: WhatsApp Business API, social media integrations, LMS (6.H).
- Pilgrim daily mobile app, offline capabilities (6.N, 6.O).

---

## § 6.E — Visa Pipeline (BL-VISA-001..003)

*(Target: `S6-E-01`)*

`visa-svc` is a new backend service that manages the visa application lifecycle for each jamaah on a departure. It follows the same conventions as all other backend services: pure gRPC internally, REST exposed via `gateway-svc`. The service currently exists as a scaffold with only `Healthz`. This section contracts the three RPCs that must be implemented.

### New gateway routes (visa-svc)

| Method | Path | Auth | Proxies to | Task |
| --- | --- | --- | --- | --- |
| `PUT` | `/v1/visas/{application_id}/status` | bearer + `visa.manage` | `visa.v1.VisaService/TransitionStatus` | `S6-E-01` |
| `POST` | `/v1/visas/bulk-submit` | bearer + `visa.manage` | `visa.v1.VisaService/BulkSubmit` | `S6-E-01` |
| `GET` | `/v1/visas` | bearer | `visa.v1.VisaService/GetApplications` | `S6-E-01` |

### DB schema — `visa-svc`

```sql
CREATE SCHEMA visa;

CREATE TABLE visa.visa_applications (
  id             TEXT PRIMARY KEY,           -- ULID "vis_..."
  jamaah_id      TEXT NOT NULL,              -- ref jamaah-svc
  booking_id     TEXT NOT NULL,              -- ref booking-svc
  departure_id   TEXT NOT NULL,              -- ref catalog-svc
  status         TEXT NOT NULL DEFAULT 'WAITING_DOCS',
    -- enum: WAITING_DOCS | READY | SUBMITTED | ISSUED | REJECTED_BY_EMBASSY | CANCELLED
  provider_id    TEXT,                       -- 'sajil' | 'mofa' | null (before submission)
  provider_ref   TEXT,                       -- provider's application reference (after SUBMITTED)
  e_visa_url     TEXT,                       -- GCS path (after ISSUED)
  issued_date    DATE,
  created_at     TIMESTAMPTZ DEFAULT NOW(),
  updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE visa.status_history (
  id             BIGSERIAL PRIMARY KEY,
  application_id TEXT NOT NULL REFERENCES visa.visa_applications(id),
  from_status    TEXT,
  to_status      TEXT NOT NULL,
  reason         TEXT,
  actor_user_id  TEXT,
  created_at     TIMESTAMPTZ DEFAULT NOW()
);
```

**ULID prefix:** `vis_` for `visa_applications.id`.

### Visa application state machine

Allowed transitions enforced by `visa-svc` service layer. Any call requesting a transition not listed here returns `invalid_transition` (400).

```
WAITING_DOCS ──────────────────────────────► READY
  (via TransitionStatus — ops verifies docs)

READY ─────────────────────────────────────► SUBMITTED
  (via BulkSubmit — bulk embassy submission)

READY ─────────────────────────────────────► CANCELLED
  (via TransitionStatus — ops cancels application)

SUBMITTED ─────────────────────────────────► ISSUED
  (via TransitionStatus — e-visa received from embassy)

SUBMITTED ─────────────────────────────────► REJECTED_BY_EMBASSY
  (via TransitionStatus — embassy rejection)

REJECTED_BY_EMBASSY ────────────────────────► READY
  (via TransitionStatus — re-submission cycle)
```

**Terminal states:** `ISSUED`, `CANCELLED` — no further transitions.

**Idempotency on same status:** if `to_status` equals the current status, the call is a no-op — returns `idempotent: true`, `200 OK`, no new `status_history` row is written.

---

### RPC — `VisaService/TransitionStatus` (BL-VISA-001)

**Purpose:** manually transition a single visa application to a new status. Enforces allowed transitions per the state machine above. Idempotent when `to_status` matches current status.

**Proto-style signature:**

```protobuf
rpc TransitionStatus(TransitionStatusRequest) returns (TransitionStatusResponse);

message TransitionStatusRequest {
  string application_id  = 1; // ULID "vis_..."
  string to_status       = 2; // target status (enum values above)
  string reason          = 3; // optional — stored in status_history.reason
  string actor_user_id   = 4; // optional — stored in status_history.actor_user_id
}

message TransitionStatusResponse {
  string application_id  = 1;
  string from_status     = 2; // the status before this transition
  string to_status       = 3; // the status after this transition (echoed)
  bool   idempotent      = 4; // true if application was already in to_status
}
```

**REST mapping — `PUT /v1/visas/{application_id}/status`:**

Request body:
```json
{
  "to_status": "READY",
  "reason": "Semua dokumen sudah diverifikasi",
  "actor_user_id": "usr_01JX..."
}
```

Response `200 OK`:
```json
{
  "application_id": "vis_01JX...",
  "from_status": "WAITING_DOCS",
  "to_status": "READY",
  "idempotent": false
}
```

**Side effects:** on every non-idempotent call, `visa-svc` inserts a row into `visa.status_history` recording `from_status`, `to_status`, `reason`, `actor_user_id`, and `created_at`.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_transition` | `FAILED_PRECONDITION` | Transition not allowed by state machine |
| `400` | `invalid_status_value` | `INVALID_ARGUMENT` | `to_status` is not a known enum value |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `visa.manage` permission |
| `404` | `not_found` | `NOT_FOUND` | `application_id` does not exist |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `VisaService/BulkSubmit` (BL-VISA-002)

**Purpose:** atomically transition all specified jamaah's visa applications for a departure from `READY` to `SUBMITTED`. All-or-nothing — any validation failure rolls back the entire batch.

**Proto-style signature:**

```protobuf
rpc BulkSubmit(BulkSubmitRequest) returns (BulkSubmitResponse);

message BulkSubmitRequest {
  string          departure_id = 1; // ULID — ref catalog-svc
  repeated string jamaah_ids   = 2; // list of jamaah_id to submit; must be non-empty
  string          provider_id  = 3; // "sajil" | "mofa"
}

message BulkSubmitResponse {
  int32           submitted_count  = 1;
  repeated string application_ids  = 2; // vis_... IDs that were transitioned
}
```

**REST mapping — `POST /v1/visas/bulk-submit`:**

Request body:
```json
{
  "departure_id": "dep_01JX...",
  "jamaah_ids": ["jmh_01JX...", "jmh_01JY..."],
  "provider_id": "sajil"
}
```

Response `200 OK`:
```json
{
  "submitted_count": 2,
  "application_ids": ["vis_01JX...", "vis_01JY..."]
}
```

**Validation rules (all checked before any DB write; any failure aborts entire batch):**

1. `jamaah_ids` must be non-empty — at least 1 entry.
2. `departure_id` must reference an existing departure in `catalog-svc` (via gRPC `CatalogService/GetDeparture`).
3. For each `jamaah_id`: a `visa_applications` row must exist for `(jamaah_id, departure_id)` with `status = 'READY'`. If any jamaah is not in `READY` status → `not_all_ready`.
4. For each jamaah: passport expiry date (fetched from jamaah-svc) must be ≥ 180 days from the departure date. If any jamaah fails → `passport_expiry_violation`.

**Transactional guarantee:** all `status` updates and `status_history` inserts are executed within a single DB transaction. If the transaction fails, no applications are transitioned.

**Side effects:** each transitioned application receives:
- `status` updated to `SUBMITTED`
- `provider_id` set to the submitted `provider_id`
- A `status_history` row with `from_status = 'READY'`, `to_status = 'SUBMITTED'`, `actor_user_id = ''` (system)

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `not_all_ready` | `FAILED_PRECONDITION` | One or more jamaah not in `READY` status; `error.details` lists which `jamaah_id` values are non-ready |
| `400` | `passport_expiry_violation` | `FAILED_PRECONDITION` | One or more passports expire within 180 days of departure; `error.details` lists `jamaah_id` + `expiry_date` |
| `400` | `empty_jamaah_list` | `INVALID_ARGUMENT` | `jamaah_ids` is empty |
| `400` | `invalid_provider` | `INVALID_ARGUMENT` | `provider_id` is not `"sajil"` or `"mofa"` |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `visa.manage` permission |
| `404` | `departure_not_found` | `NOT_FOUND` | `departure_id` does not exist in catalog-svc |
| `500` | `internal_error` | `INTERNAL` | DB transaction failure or unexpected error |

---

### RPC — `VisaService/GetApplications` (BL-VISA-003)

**Purpose:** list visa applications for a departure, with full status history embedded. Supports optional status filter.

**Proto-style signature:**

```protobuf
rpc GetApplications(GetApplicationsRequest) returns (GetApplicationsResponse);

message GetApplicationsRequest {
  string departure_id   = 1; // required
  string status_filter  = 2; // optional — filter by status enum value; empty = all statuses
}

message GetApplicationsResponse {
  repeated ApplicationRecord applications = 1;
}

message ApplicationRecord {
  string                      id           = 1; // vis_...
  string                      jamaah_id    = 2;
  string                      status       = 3;
  string                      provider_ref = 4; // empty before SUBMITTED
  string                      issued_date  = 5; // "YYYY-MM-DD" or empty
  repeated StatusHistoryEntry history      = 6;
}

message StatusHistoryEntry {
  string from_status = 1; // empty string for initial creation
  string to_status   = 2;
  string reason      = 3;
  string created_at  = 4; // RFC3339
}
```

**REST mapping — `GET /v1/visas?departure_id=...&status=...`:**

Query parameters:

| Parameter | Type | Required | Notes |
| --- | --- | --- | --- |
| `departure_id` | string | yes | Filter by departure |
| `status` | string | no | Filter by `status` enum value; omit to return all statuses |

Response `200 OK`:
```json
{
  "applications": [
    {
      "id": "vis_01JX...",
      "jamaah_id": "jmh_01JX...",
      "status": "SUBMITTED",
      "provider_ref": "SAJIL-REF-12345",
      "issued_date": "",
      "history": [
        {
          "from_status": "",
          "to_status": "WAITING_DOCS",
          "reason": "Application created",
          "created_at": "2026-04-01T08:00:00Z"
        },
        {
          "from_status": "WAITING_DOCS",
          "to_status": "READY",
          "reason": "Semua dokumen sudah diverifikasi",
          "created_at": "2026-04-10T10:00:00Z"
        },
        {
          "from_status": "READY",
          "to_status": "SUBMITTED",
          "reason": "",
          "created_at": "2026-04-15T09:00:00Z"
        }
      ]
    }
  ]
}
```

**Ordering:** applications ordered by `created_at` ASC; `history` entries ordered by `status_history.id` ASC.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `missing_departure_id` | `INVALID_ARGUMENT` | `departure_id` query param is absent |
| `400` | `invalid_status_filter` | `INVALID_ARGUMENT` | `status` param is not a known enum value |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

### Honored by implementation

- `S6-E-01` — `visa-svc`: `TransitionStatus` handler, state machine enforcement, `status_history` insert on non-idempotent calls.
- `S6-E-01` — `visa-svc`: `BulkSubmit` handler, transactional all-or-nothing batch transition, passport-expiry validation via jamaah-svc gRPC adapter, departure lookup via catalog-svc gRPC adapter.
- `S6-E-01` — `visa-svc`: `GetApplications` handler, left-join with `status_history`.
- `S6-E-01` — `gateway-svc`: three new route registrations for `visa-svc`, `RequirePermission("visa.manage")` on `PUT` and `POST` paths.

---

## § 6.D — Field Operations Scan (BL-OPS-010/011)

*(Target: `S6-E-02`)*

`ops-svc` gains three RPCs for scan event recording and bus boarding management. All scan operations are idempotent via a caller-supplied `idempotency_key`.

### New gateway routes (ops-svc additions)

| Method | Path | Auth | Proxies to | Task |
| --- | --- | --- | --- | --- |
| `POST` | `/v1/ops/scans` | bearer | `ops.v1.OpsService/RecordScan` | `S6-E-02` |
| `POST` | `/v1/ops/bus-boarding` | bearer | `ops.v1.OpsService/RecordBusBoarding` | `S6-E-02` |
| `GET` | `/v1/ops/bus-boarding/{departure_id}` | bearer | `ops.v1.OpsService/GetBoardingRoster` | `S6-E-02` |

### DB schema additions — `ops-svc`

```sql
CREATE TABLE ops.scan_events (
  id               TEXT PRIMARY KEY,          -- ULID "scan_..."
  scan_type        TEXT NOT NULL,             -- 'ALL' | 'bus_boarding' | 'luggage' | 'raudhah'
  departure_id     TEXT NOT NULL,
  jamaah_id        TEXT NOT NULL,
  scanned_by       TEXT NOT NULL,             -- user_id of scanner
  device_id        TEXT,
  location         TEXT,
  idempotency_key  TEXT NOT NULL UNIQUE,      -- prevents duplicate scans
  metadata         JSONB DEFAULT '{}',
  scanned_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE ops.bus_boardings (
  id              TEXT PRIMARY KEY,           -- ULID "bbd_..."
  departure_id    TEXT NOT NULL,
  bus_number      TEXT NOT NULL,
  jamaah_id       TEXT NOT NULL,
  status          TEXT NOT NULL DEFAULT 'boarded', -- 'boarded' | 'absent' | 'late'
  scan_event_id   TEXT REFERENCES ops.scan_events(id),
  boarded_at      TIMESTAMPTZ
);
```

**ULID prefixes:** `scan_` for `scan_events.id`, `bbd_` for `bus_boardings.id`.

**Idempotency convention for scans:** the caller constructs the idempotency key as `{scan_type}:{departure_id}:{jamaah_id}:{YYYY-MM-DD}`. This allows a jamaah to be re-scanned on a different day without collision, but deduplicates multiple scans of the same type for the same jamaah on the same day.

---

### RPC — `OpsService/RecordScan` (BL-OPS-010)

**Purpose:** record any scan event idempotently. If the `idempotency_key` already exists, the call is a no-op.

**Proto-style signature:**

```protobuf
rpc RecordScan(RecordScanRequest) returns (RecordScanResponse);

message RecordScanRequest {
  string scan_type       = 1; // 'ALL' | 'bus_boarding' | 'luggage' | 'raudhah'
  string departure_id    = 2;
  string jamaah_id       = 3;
  string scanned_by      = 4; // user_id of the scanner
  string device_id       = 5; // optional
  string location        = 6; // optional — freetext location label
  string idempotency_key = 7; // caller-generated; format: {scan_type}:{departure_id}:{jamaah_id}:{YYYY-MM-DD}
  bytes  metadata        = 8; // optional JSONB payload
}

message RecordScanResponse {
  string scan_id    = 1; // scan_... ULID of the inserted (or existing) row
  bool   idempotent = 2; // true if key already existed; no new row was written
}
```

**REST mapping — `POST /v1/ops/scans`:**

Request body:
```json
{
  "scan_type": "luggage",
  "departure_id": "dep_01JX...",
  "jamaah_id": "jmh_01JX...",
  "scanned_by": "usr_01JX...",
  "device_id": "scanner-07",
  "location": "Gate B — Soetta",
  "idempotency_key": "luggage:dep_01JX...:jmh_01JX...:2026-06-15",
  "metadata": {}
}
```

Response `200 OK`:
```json
{
  "scan_id": "scan_01JX...",
  "idempotent": false
}
```

**Idempotency behavior:** `ops-svc` attempts an `INSERT ... ON CONFLICT (idempotency_key) DO NOTHING`. If a row already existed, `scan_id` returns the existing row's `id` and `idempotent = true`.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_scan_type` | `INVALID_ARGUMENT` | `scan_type` is not one of the allowed enum values |
| `400` | `missing_required_field` | `INVALID_ARGUMENT` | `departure_id`, `jamaah_id`, `scanned_by`, or `idempotency_key` is empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `OpsService/RecordBusBoarding` (BL-OPS-011)

**Purpose:** record a bus boarding event for a jamaah. Creates a `scan_events` row (with `scan_type = 'bus_boarding'`) and a `bus_boardings` row atomically. Idempotent: if a boarding row already exists for `(departure_id, jamaah_id)`, the call returns the existing row with `idempotent: true`.

**Proto-style signature:**

```protobuf
rpc RecordBusBoarding(RecordBusBoardingRequest) returns (RecordBusBoardingResponse);

message RecordBusBoardingRequest {
  string departure_id = 1;
  string bus_number   = 2;
  string jamaah_id    = 3;
  string scanned_by   = 4; // user_id
  string status       = 5; // 'boarded' | 'absent' | 'late'; default 'boarded' if empty
}

message RecordBusBoardingResponse {
  string boarding_id = 1; // bbd_... ULID
  string status      = 2; // echoed status
  bool   idempotent  = 3;
}
```

**REST mapping — `POST /v1/ops/bus-boarding`:**

Request body:
```json
{
  "departure_id": "dep_01JX...",
  "bus_number": "Bus-A",
  "jamaah_id": "jmh_01JX...",
  "scanned_by": "usr_01JX...",
  "status": "boarded"
}
```

Response `200 OK`:
```json
{
  "boarding_id": "bbd_01JX...",
  "status": "boarded",
  "idempotent": false
}
```

**Atomic write:** `ops-svc` executes within a single transaction:
1. `INSERT INTO ops.scan_events` (using idempotency_key = `bus_boarding:{departure_id}:{jamaah_id}:{YYYY-MM-DD}`)
2. `INSERT INTO ops.bus_boardings` linking to the new `scan_event_id`

If a boarding row already exists for `(departure_id, jamaah_id)`, the service skips both inserts and returns the existing `boarding_id` with `idempotent: true`.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_status` | `INVALID_ARGUMENT` | `status` not one of `boarded`, `absent`, `late` |
| `400` | `missing_required_field` | `INVALID_ARGUMENT` | Any of `departure_id`, `bus_number`, `jamaah_id`, `scanned_by` is empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `OpsService/GetBoardingRoster` (BL-OPS-011)

**Purpose:** retrieve the boarding status roster for a departure, with optional bus filter.

**Proto-style signature:**

```protobuf
rpc GetBoardingRoster(GetBoardingRosterRequest) returns (GetBoardingRosterResponse);

message GetBoardingRosterRequest {
  string departure_id = 1; // required
  string bus_number   = 2; // optional — filter to a single bus
}

message GetBoardingRosterResponse {
  repeated BoardingEntry boardings     = 1;
  int32                  total_boarded = 2;
  int32                  total_absent  = 3;
}

message BoardingEntry {
  string jamaah_id  = 1;
  string bus_number = 2;
  string status     = 3; // 'boarded' | 'absent' | 'late'
  string boarded_at = 4; // RFC3339 or empty string if absent
}
```

**REST mapping — `GET /v1/ops/bus-boarding/{departure_id}?bus_number=...`:**

Response `200 OK`:
```json
{
  "boardings": [
    {
      "jamaah_id": "jmh_01JX...",
      "bus_number": "Bus-A",
      "status": "boarded",
      "boarded_at": "2026-06-15T05:30:00Z"
    },
    {
      "jamaah_id": "jmh_01JY...",
      "bus_number": "Bus-A",
      "status": "absent",
      "boarded_at": ""
    }
  ],
  "total_boarded": 1,
  "total_absent": 1
}
```

**Aggregation:** `total_boarded` = count of rows where `status IN ('boarded', 'late')`; `total_absent` = count where `status = 'absent'`.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

### Honored by implementation

- `S6-E-02` — `ops-svc`: `RecordScan` handler, `INSERT ON CONFLICT DO NOTHING`, idempotency return logic.
- `S6-E-02` — `ops-svc`: `RecordBusBoarding` handler, transactional scan + boarding insert, dedup on `(departure_id, jamaah_id)`.
- `S6-E-02` — `ops-svc`: `GetBoardingRoster` handler, aggregation query with optional bus_number filter.
- `S6-E-02` — `gateway-svc`: three new route registrations for `ops-svc` scan/boarding paths.

---

## § 6.C — Warehouse & Procurement (BL-LOG-010..012)

*(Target: `S6-E-03`)*

`logistics-svc` gains four RPCs covering purchase requests, goods-receipt-with-QC, and kit assembly.

### New gateway routes (logistics-svc additions)

| Method | Path | Auth | Proxies to | Task |
| --- | --- | --- | --- | --- |
| `POST` | `/v1/logistics/purchase-requests` | bearer + `logistics.manage` | `logistics.v1.LogisticsService/CreatePurchaseRequest` | `S6-E-03` |
| `PUT` | `/v1/logistics/purchase-requests/{id}/decision` | bearer + `logistics.approve` | `logistics.v1.LogisticsService/ApprovePurchaseRequest` | `S6-E-03` |
| `POST` | `/v1/finance/grn` | bearer + `logistics.manage` | `logistics.v1.LogisticsService/RecordGRNWithQC` | `S6-E-03` |
| `POST` | `/v1/logistics/kit-assembly` | bearer + `logistics.manage` | `logistics.v1.LogisticsService/CreateKitAssembly` | `S6-E-03` |

> **Note on `POST /v1/finance/grn`:** this route already exists in `S3` for GRN without QC. S6 extends the request body by adding the `qc_passed` and `qc_notes` fields. The route is not versioned — it is a backward-compatible additive extension. Callers that do not send `qc_passed` receive `qc_passed = true` by default (preserving S3 behavior).

### DB schema additions — `logistics-svc`

```sql
CREATE TABLE logistics.purchase_requests (
  id               TEXT PRIMARY KEY,          -- ULID "pr_..."
  departure_id     TEXT NOT NULL,
  requested_by     TEXT NOT NULL,
  item_name        TEXT NOT NULL,
  quantity         INT NOT NULL,
  unit_price_idr   BIGINT NOT NULL,
  total_price_idr  BIGINT NOT NULL,            -- quantity * unit_price_idr
  budget_limit_idr BIGINT,                     -- if set, request rejected if total > limit
  status           TEXT NOT NULL DEFAULT 'pending',
    -- pending | approved | rejected | fulfilled
  approved_by      TEXT,
  notes            TEXT,
  created_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE logistics.kit_assemblies (
  id               TEXT PRIMARY KEY,           -- ULID "kit_..."
  departure_id     TEXT NOT NULL,
  assembled_by     TEXT NOT NULL,
  status           TEXT NOT NULL DEFAULT 'pending',  -- pending | completed | failed
  idempotency_key  TEXT NOT NULL UNIQUE,
  assembled_at     TIMESTAMPTZ,
  created_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE logistics.kit_assembly_items (
  id           BIGSERIAL PRIMARY KEY,
  assembly_id  TEXT NOT NULL REFERENCES logistics.kit_assemblies(id),
  item_name    TEXT NOT NULL,
  quantity     INT NOT NULL,
  fulfilled    BOOLEAN NOT NULL DEFAULT FALSE
);
```

**ULID prefixes:** `pr_` for `purchase_requests.id`, `kit_` for `kit_assemblies.id`.

---

### RPC — `LogisticsService/CreatePurchaseRequest` (BL-LOG-010)

**Purpose:** create a purchase request with optional budget gate. Validates that `total_price_idr = quantity × unit_price_idr` and rejects if `budget_limit_idr` is set and `total_price_idr > budget_limit_idr`.

**Proto-style signature:**

```protobuf
rpc CreatePurchaseRequest(CreatePurchaseRequestRequest)
    returns (CreatePurchaseRequestResponse);

message CreatePurchaseRequestRequest {
  string departure_id      = 1;
  string requested_by      = 2; // user_id
  string item_name         = 3;
  int32  quantity          = 4; // must be >= 1
  int64  unit_price_idr    = 5; // integer IDR; must be >= 1
  int64  budget_limit_idr  = 6; // optional; 0 = no limit
}

message CreatePurchaseRequestResponse {
  string pr_id            = 1; // pr_... ULID
  string status           = 2; // always "pending" on creation
  int64  total_price_idr  = 3; // quantity * unit_price_idr
}
```

**REST mapping — `POST /v1/logistics/purchase-requests`:**

Request body:
```json
{
  "departure_id": "dep_01JX...",
  "requested_by": "usr_01JX...",
  "item_name": "Kain Ihram",
  "quantity": 100,
  "unit_price_idr": 75000,
  "budget_limit_idr": 8000000
}
```

Response `201 Created`:
```json
{
  "pr_id": "pr_01JX...",
  "status": "pending",
  "total_price_idr": 7500000
}
```

**Validation rules:**

1. `total_price_idr` is computed server-side as `quantity × unit_price_idr`. The computed value is stored; callers do not send `total_price_idr`.
2. If `budget_limit_idr > 0` and computed `total_price_idr > budget_limit_idr` → return `budget_exceeded` (400). No row is inserted.
3. `quantity` must be ≥ 1; `unit_price_idr` must be ≥ 1.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `budget_exceeded` | `FAILED_PRECONDITION` | `total_price_idr > budget_limit_idr` |
| `400` | `invalid_quantity` | `INVALID_ARGUMENT` | `quantity < 1` |
| `400` | `invalid_unit_price` | `INVALID_ARGUMENT` | `unit_price_idr < 1` |
| `400` | `missing_required_field` | `INVALID_ARGUMENT` | `departure_id`, `requested_by`, or `item_name` is empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `logistics.manage` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `LogisticsService/ApprovePurchaseRequest` (BL-LOG-010)

**Purpose:** approve or reject a purchase request. Only `pending` PRs may be approved or rejected.

**Proto-style signature:**

```protobuf
rpc ApprovePurchaseRequest(ApprovePurchaseRequestRequest)
    returns (ApprovePurchaseRequestResponse);

message ApprovePurchaseRequestRequest {
  string pr_id       = 1; // pr_... ULID
  string approved_by = 2; // user_id of approver
  bool   approved    = 3; // true = approve; false = reject
  string notes       = 4; // optional — stored in purchase_requests.notes
}

message ApprovePurchaseRequestResponse {
  string pr_id      = 1;
  string new_status = 2; // "approved" or "rejected"
}
```

**REST mapping — `PUT /v1/logistics/purchase-requests/{id}/decision`:**

Request body:
```json
{
  "approved_by": "usr_01JX...",
  "approved": true,
  "notes": "Budget dalam batas wajar"
}
```

Response `200 OK`:
```json
{
  "pr_id": "pr_01JX...",
  "new_status": "approved"
}
```

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_status_for_decision` | `FAILED_PRECONDITION` | PR is not in `pending` status |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `logistics.approve` permission |
| `404` | `not_found` | `NOT_FOUND` | `pr_id` does not exist |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `LogisticsService/RecordGRNWithQC` (BL-LOG-011)

**Purpose:** record a goods receipt note (GRN) with QC pass/fail. If `qc_passed = true`, behavior is identical to the S3 GRN flow (AP journal is posted via `finance-svc`). If `qc_passed = false`, `grn_status` is set to `rejected` and **no AP journal is posted**; if a journal was already posted (race condition), `finance-svc.CorrectJournal` must be invoked manually.

**Proto-style signature:**

```protobuf
rpc RecordGRNWithQC(RecordGRNWithQCRequest) returns (RecordGRNWithQCResponse);

message RecordGRNWithQCRequest {
  string grn_id       = 1; // logistics-svc GRN ULID
  string departure_id = 2;
  int64  amount_idr   = 3; // integer IDR
  bool   qc_passed    = 4; // true = pass, false = fail/reject
  string qc_notes     = 5; // optional — reason for failure
}

message RecordGRNWithQCResponse {
  string grn_id        = 1;
  string qc_status     = 2; // "passed" | "rejected"
  bool   journal_posted = 3; // true only when qc_passed=true and AP journal was posted
}
```

**REST mapping — `POST /v1/finance/grn` (field additions only):**

Request body additions vs. S3 contract:
```json
{
  "grn_id": "grn_01JX...",
  "departure_id": "dep_01JX...",
  "amount_idr": 7500000,
  "qc_passed": true,
  "qc_notes": ""
}
```

Response `200 OK`:
```json
{
  "grn_id": "grn_01JX...",
  "qc_status": "passed",
  "journal_posted": true
}
```

**QC logic:**

| `qc_passed` | `grn_status` set | AP journal | `journal_posted` |
| --- | --- | --- | --- |
| `true` | `received` | Posted (Dr Inventory / Cr AP — same as S3) | `true` |
| `false` | `rejected` | NOT posted | `false` |

If `qc_passed = false` and a journal was already posted (e.g., due to a prior call or race): `logistics-svc` sets `grn_status = 'rejected'` and returns `journal_posted = false` with an additional field `requires_manual_correction: true`. Finance ops must then call `POST /v1/finance/journals/{id}/correct` (S3 `CorrectJournal`) to reverse the erroneous posting.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `missing_required_field` | `INVALID_ARGUMENT` | `grn_id`, `departure_id`, or `amount_idr` missing |
| `400` | `invalid_amount` | `INVALID_ARGUMENT` | `amount_idr < 1` |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `logistics.manage` permission |
| `404` | `grn_not_found` | `NOT_FOUND` | `grn_id` does not exist in logistics-svc |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `LogisticsService/CreateKitAssembly` (BL-LOG-012)

**Purpose:** create an idempotent kit assembly record. All items must be marked `fulfilled = true` within the same transaction before the assembly status is set to `completed`. If any item cannot be fulfilled (e.g., stock shortage), the transaction rolls back and status is set to `failed`.

**Proto-style signature:**

```protobuf
rpc CreateKitAssembly(CreateKitAssemblyRequest) returns (CreateKitAssemblyResponse);

message CreateKitAssemblyRequest {
  string             departure_id    = 1;
  string             assembled_by    = 2; // user_id
  repeated KitItem   items           = 3; // must be non-empty
  string             idempotency_key = 4;
}

message KitItem {
  string item_name = 1;
  int32  quantity  = 2; // must be >= 1
}

message CreateKitAssemblyResponse {
  string assembly_id = 1; // kit_... ULID
  string status      = 2; // "completed" | "failed"
  bool   idempotent  = 3;
}
```

**REST mapping — `POST /v1/logistics/kit-assembly`:**

Request body:
```json
{
  "departure_id": "dep_01JX...",
  "assembled_by": "usr_01JX...",
  "items": [
    { "item_name": "Kain Ihram", "quantity": 50 },
    { "item_name": "Tas Tenteng", "quantity": 50 }
  ],
  "idempotency_key": "kit:dep_01JX...:2026-06-01"
}
```

Response `200 OK`:
```json
{
  "assembly_id": "kit_01JX...",
  "status": "completed",
  "idempotent": false
}
```

**Transactional behavior:**
1. On `idempotency_key` collision → return existing `assembly_id`, `status`, and `idempotent: true`. No new rows.
2. Otherwise: insert `kit_assemblies` row and all `kit_assembly_items` rows within a single transaction.
3. Mark every item `fulfilled = true`. If all succeed → set `status = 'completed'`, set `assembled_at = NOW()`.
4. If any item fails fulfillment → rollback entire transaction, set `status = 'failed'` (via a separate error-path INSERT that persists the failure record without items).

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `empty_items_list` | `INVALID_ARGUMENT` | `items` array is empty |
| `400` | `invalid_item_quantity` | `INVALID_ARGUMENT` | Any item has `quantity < 1` |
| `400` | `assembly_failed` | `FAILED_PRECONDITION` | One or more items could not be fulfilled; `error.details` lists which items; assembly persisted with `status = 'failed'` |
| `400` | `missing_required_field` | `INVALID_ARGUMENT` | `departure_id`, `assembled_by`, or `idempotency_key` is empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `logistics.manage` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

### Honored by implementation

- `S6-E-03` — `logistics-svc`: `CreatePurchaseRequest` handler, budget gate validation, total price computation.
- `S6-E-03` — `logistics-svc`: `ApprovePurchaseRequest` handler, status gate (pending only), `approved_by` + `notes` persistence.
- `S6-E-03` — `logistics-svc`: `RecordGRNWithQC` handler, QC flag routing, conditional AP journal call to `finance-svc`.
- `S6-E-03` — `logistics-svc`: `CreateKitAssembly` handler, transactional all-or-nothing fulfill loop, idempotency key dedup.
- `S6-E-03` — `gateway-svc`: four new route registrations, `RequirePermission` on `logistics.manage` and `logistics.approve` paths.

---

## § 6.B — Finance AP/AR Depth (BL-FIN-010/011)

*(Target: `S6-E-04`)*

`finance-svc` gains three RPCs: AP disbursement batch creation, approval (with automatic journal posting on approval), and an AR/AP aging bucket report.

### New gateway routes (finance-svc additions)

| Method | Path | Auth | Proxies to | Task |
| --- | --- | --- | --- | --- |
| `POST` | `/v1/finance/disbursements` | bearer + `finance.manage` | `finance.v1.FinanceService/CreateDisbursementBatch` | `S6-E-04` |
| `PUT` | `/v1/finance/disbursements/{id}/decision` | bearer + `finance.approve` | `finance.v1.FinanceService/ApproveDisbursement` | `S6-E-04` |
| `GET` | `/v1/finance/aging` | bearer + `finance.read` | `finance.v1.FinanceService/GetARAPAging` | `S6-E-04` |

### DB schema additions — `finance-svc`

```sql
CREATE TABLE finance.disbursement_batches (
  id                TEXT PRIMARY KEY,           -- ULID "dis_..."
  description       TEXT NOT NULL,
  total_amount_idr  BIGINT NOT NULL,
  status            TEXT NOT NULL DEFAULT 'pending_approval',
    -- pending_approval | approved | rejected | disbursed
  approved_by       TEXT,
  approved_at       TIMESTAMPTZ,
  created_by        TEXT NOT NULL,
  created_at        TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE finance.disbursement_items (
  id                BIGSERIAL PRIMARY KEY,
  batch_id          TEXT NOT NULL REFERENCES finance.disbursement_batches(id),
  vendor_name       TEXT NOT NULL,
  description       TEXT NOT NULL,
  amount_idr        BIGINT NOT NULL,
  reference         TEXT,
  journal_entry_id  TEXT   -- filled when batch is disbursed / approved
);
```

**ULID prefix:** `dis_` for `disbursement_batches.id`.

**`total_amount_idr` invariant:** always equals `SUM(items[].amount_idr)`. Computed server-side at creation; not provided by caller.

---

### RPC — `FinanceService/CreateDisbursementBatch` (BL-FIN-010)

**Purpose:** create an AP disbursement batch in `pending_approval` status. `total_amount_idr` is computed server-side from the sum of item amounts.

**Proto-style signature:**

```protobuf
rpc CreateDisbursementBatch(CreateDisbursementBatchRequest)
    returns (CreateDisbursementBatchResponse);

message CreateDisbursementBatchRequest {
  string                 description = 1;
  repeated DisbursementItem items    = 2; // must be non-empty
  string                 created_by  = 3;
}

message DisbursementItem {
  string vendor_name  = 1;
  string description  = 2;
  int64  amount_idr   = 3; // must be >= 1
  string reference    = 4; // optional — e.g. PO number, invoice number
}

message CreateDisbursementBatchResponse {
  string batch_id          = 1; // dis_... ULID
  int64  total_amount_idr  = 2; // SUM of items
  int32  item_count        = 3;
  string status            = 4; // always "pending_approval"
}
```

**REST mapping — `POST /v1/finance/disbursements`:**

Request body:
```json
{
  "description": "Pembayaran hotel Makkah batch Juni 2026",
  "items": [
    {
      "vendor_name": "Hotel Al-Safwa",
      "description": "Akomodasi Makkah 10 malam",
      "amount_idr": 150000000,
      "reference": "PO-2026-045"
    },
    {
      "vendor_name": "Hotel Pullman Zam Zam",
      "description": "Akomodasi Madinah 5 malam",
      "amount_idr": 80000000,
      "reference": "PO-2026-046"
    }
  ],
  "created_by": "usr_01JX..."
}
```

Response `201 Created`:
```json
{
  "batch_id": "dis_01JX...",
  "total_amount_idr": 230000000,
  "item_count": 2,
  "status": "pending_approval"
}
```

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `empty_items_list` | `INVALID_ARGUMENT` | `items` array is empty |
| `400` | `invalid_item_amount` | `INVALID_ARGUMENT` | Any item has `amount_idr < 1` |
| `400` | `missing_required_field` | `INVALID_ARGUMENT` | `description`, `created_by`, or any item's required fields are empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `finance.manage` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `FinanceService/ApproveDisbursement` (BL-FIN-010)

**Purpose:** approve or reject a disbursement batch. On approval, `finance-svc` atomically posts one double-entry journal entry per item (Dr AP / Cr Cash) and fills `disbursement_items.journal_entry_id`. On rejection, status transitions to `rejected` with no journal posting.

**Proto-style signature:**

```protobuf
rpc ApproveDisbursement(ApproveDisbursementRequest)
    returns (ApproveDisbursementResponse);

message ApproveDisbursementRequest {
  string batch_id     = 1;
  string approved_by  = 2; // user_id
  bool   approved     = 3; // true = approve; false = reject
  string notes        = 4; // optional
}

message ApproveDisbursementResponse {
  string          batch_id          = 1;
  string          status            = 2; // "approved" | "rejected"
  repeated string journal_entry_ids = 3; // one per item; empty if rejected
}
```

**REST mapping — `PUT /v1/finance/disbursements/{id}/decision`:**

Request body:
```json
{
  "approved_by": "usr_01JX...",
  "approved": true,
  "notes": ""
}
```

Response `200 OK` (approved):
```json
{
  "batch_id": "dis_01JX...",
  "status": "approved",
  "journal_entry_ids": ["je_01JX...", "je_01JY..."]
}
```

Response `200 OK` (rejected):
```json
{
  "batch_id": "dis_01JX...",
  "status": "rejected",
  "journal_entry_ids": []
}
```

**Journal posting (on approval):** for each `disbursement_items` row, `finance-svc` calls its own internal `PostJournal` logic (per S3 contract) with:

```
idempotency_key = "disbursement:{batch_id}:{item.id}"
lines:
  - account_code: "2101"  (AP — Accounts Payable)   debit_idr: item.amount_idr
  - account_code: "1001"  (Cash / Bank)              credit_idr: item.amount_idr
```

All journal inserts and the `disbursement_batches` status update are performed within a single DB transaction. If any journal insert fails, the entire transaction rolls back and the batch remains `pending_approval`.

**Idempotency on re-call:** if the batch is already in `approved` or `rejected` status, the call returns the current status with `idempotent: true` behavior (no second journal posting).

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_status_for_decision` | `FAILED_PRECONDITION` | Batch is not in `pending_approval` status |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `finance.approve` permission |
| `404` | `not_found` | `NOT_FOUND` | `batch_id` does not exist |
| `500` | `internal_error` | `INTERNAL` | DB transaction failure or unexpected error |

---

### RPC — `FinanceService/GetARAPAging` (BL-FIN-011)

**Purpose:** return AR/AP aging buckets — the outstanding balance grouped by days overdue as of a given date. Supports AR-only, AP-only, or both.

**Aging buckets:**

| Bucket | Days overdue (from due_date to `as_of_date`) |
| --- | --- |
| `current` | 0 days (not yet due) |
| `days_30` | 1–30 days |
| `days_60` | 31–60 days |
| `days_90` | 61–90 days |
| `over_90` | > 90 days |

**Proto-style signature:**

```protobuf
rpc GetARAPAging(GetARAPAgingRequest) returns (GetARAPAgingResponse);

message GetARAPAgingRequest {
  string type        = 1; // "AR" | "AP" | "both"
  string as_of_date  = 2; // "YYYY-MM-DD"; defaults to today if empty
}

message GetARAPAgingResponse {
  AgingBuckets ar           = 1; // null/empty if type = "AP"
  AgingBuckets ap           = 2; // null/empty if type = "AR"
  string       generated_at = 3; // RFC3339
}

message AgingBuckets {
  int64 current  = 1; // not yet due
  int64 days_30  = 2; // 1–30 days overdue
  int64 days_60  = 3; // 31–60 days overdue
  int64 days_90  = 4; // 61–90 days overdue
  int64 over_90  = 5; // > 90 days overdue
  int64 total    = 6; // sum of all buckets
}
```

**REST mapping — `GET /v1/finance/aging?type=AR&as_of_date=2026-06-30`:**

Query parameters:

| Parameter | Type | Required | Default | Notes |
| --- | --- | --- | --- | --- |
| `type` | string | yes | — | `AR`, `AP`, or `both` |
| `as_of_date` | string (YYYY-MM-DD) | no | today | Aging computed relative to this date |

Response `200 OK`:
```json
{
  "ar": {
    "current": 150000000,
    "days_30": 25000000,
    "days_60": 10000000,
    "days_90": 0,
    "over_90": 5000000,
    "total": 190000000
  },
  "ap": null,
  "generated_at": "2026-04-23T15:00:00Z"
}
```

**Data source:** `finance-svc` computes aging from the `invoices` table (AR) and `disbursement_batches` table (AP) joined with `due_date` or equivalent aging anchor. Only outstanding (unpaid / unapproved) amounts are bucketed.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_type` | `INVALID_ARGUMENT` | `type` not one of `AR`, `AP`, `both` |
| `400` | `invalid_date_format` | `INVALID_ARGUMENT` | `as_of_date` not a valid `YYYY-MM-DD` string |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `finance.read` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

### Honored by implementation

- `S6-E-04` — `finance-svc`: `CreateDisbursementBatch` handler, server-side total computation, items bulk insert.
- `S6-E-04` — `finance-svc`: `ApproveDisbursement` handler, transactional journal posting per item (Dr AP / Cr Cash), idempotency on already-decided batches.
- `S6-E-04` — `finance-svc`: `GetARAPAging` handler, bucket aggregation query parameterized by `as_of_date`.
- `S6-E-04` — `gateway-svc`: three new route registrations with `RequirePermission("finance.manage")`, `RequirePermission("finance.approve")`, `RequirePermission("finance.read")`.

---

## § 6.M — IAM Admin/Security Depth (BL-IAM-007/010/011/014/016)

*(Target: `S6-E-05`)*

`iam-svc` gains six RPCs covering data visibility scopes, API key management, global config, and paginated activity log search.

### New gateway routes (iam-svc additions)

| Method | Path | Auth | Proxies to | Task |
| --- | --- | --- | --- | --- |
| `PUT` | `/v1/admin/users/{id}/data-scope` | bearer + `iam.admin` | `iam.v1.IamService/SetDataScope` | `S6-E-05` |
| `POST` | `/v1/admin/api-keys` | bearer + `iam.admin` | `iam.v1.IamService/CreateAPIKey` | `S6-E-05` |
| `DELETE` | `/v1/admin/api-keys/{id}` | bearer + `iam.admin` | `iam.v1.IamService/RevokeAPIKey` | `S6-E-05` |
| `GET` | `/v1/admin/config` | bearer + `iam.admin` | `iam.v1.IamService/GetGlobalConfig` | `S6-E-05` |
| `PUT` | `/v1/admin/config/{key}` | bearer + `iam.admin` | `iam.v1.IamService/SetGlobalConfig` | `S6-E-05` |
| `GET` | `/v1/admin/activity-log` | bearer + `iam.admin` | `iam.v1.IamService/SearchActivityLog` | `S6-E-05` |

### DB schema additions — `iam-svc`

```sql
-- Data visibility scopes (BL-IAM-007)
CREATE TABLE iam.data_scopes (
  id          TEXT PRIMARY KEY,           -- ULID "dsc_..."
  user_id     TEXT NOT NULL UNIQUE,       -- one scope record per user
  scope_type  TEXT NOT NULL,              -- 'global' | 'branch' | 'own_only'
  branch_id   TEXT,                       -- only for scope_type='branch'
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- API keys (BL-IAM-014)
CREATE TABLE iam.api_keys (
  id           TEXT PRIMARY KEY,           -- ULID "akey_..."
  name         TEXT NOT NULL,
  key_hash     TEXT NOT NULL UNIQUE,       -- argon2id hash of the full key
  key_prefix   TEXT NOT NULL,              -- first 8 chars of plaintext key (shown in UI)
  created_by   TEXT NOT NULL,
  expires_at   TIMESTAMPTZ,
  last_used_at TIMESTAMPTZ,
  revoked_at   TIMESTAMPTZ,
  scopes       TEXT[] NOT NULL DEFAULT '{}',
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

-- Global config (BL-IAM-016)
CREATE TABLE iam.global_config (
  key         TEXT PRIMARY KEY,
  value       TEXT NOT NULL,
  description TEXT,
  updated_by  TEXT,
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);
```

**ULID prefixes:** `dsc_` for `data_scopes.id`, `akey_` for `api_keys.id`.

**API key security note:** the plaintext key is returned **only once** at creation (via `CreateAPIKey`). `iam-svc` stores only the argon2id hash. If the caller loses the plaintext key, it must be revoked and a new one created.

---

### RPC — `IamService/SetDataScope` (BL-IAM-007)

**Purpose:** set or replace the data visibility scope for a user. Upsert semantics — creates the record if absent, replaces if present.

**Proto-style signature:**

```protobuf
rpc SetDataScope(SetDataScopeRequest) returns (SetDataScopeResponse);

message SetDataScopeRequest {
  string user_id    = 1;
  string scope_type = 2; // 'global' | 'branch' | 'own_only'
  string branch_id  = 3; // required when scope_type = 'branch'; ignored otherwise
}

message SetDataScopeResponse {
  string user_id    = 1;
  string scope_type = 2;
}
```

**REST mapping — `PUT /v1/admin/users/{id}/data-scope`:**

Request body:
```json
{
  "scope_type": "branch",
  "branch_id": "brn_01JX..."
}
```

Response `200 OK`:
```json
{
  "user_id": "usr_01JX...",
  "scope_type": "branch"
}
```

**Validation:** if `scope_type = 'branch'` and `branch_id` is empty → `missing_branch_id` (400).

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_scope_type` | `INVALID_ARGUMENT` | `scope_type` not one of `global`, `branch`, `own_only` |
| `400` | `missing_branch_id` | `INVALID_ARGUMENT` | `scope_type = 'branch'` but `branch_id` is empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `iam.admin` permission |
| `404` | `user_not_found` | `NOT_FOUND` | `user_id` does not exist in `iam.users` |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `IamService/CreateAPIKey` (BL-IAM-014)

**Purpose:** create an API key. Returns the plaintext key exactly once. The plaintext key is never stored — only its argon2id hash is persisted.

**Proto-style signature:**

```protobuf
rpc CreateAPIKey(CreateAPIKeyRequest) returns (CreateAPIKeyResponse);

message CreateAPIKeyRequest {
  string          name       = 1;
  repeated string scopes     = 2; // list of permission strings the key carries
  string          expires_at = 3; // RFC3339; empty = no expiry
  string          created_by = 4; // user_id
}

message CreateAPIKeyResponse {
  string key_id        = 1; // akey_... ULID
  string plaintext_key = 2; // ONLY returned here — never stored; never returned again
  string key_prefix    = 3; // first 8 chars of plaintext_key (for UI display)
  string expires_at    = 4; // RFC3339 or empty
}
```

**REST mapping — `POST /v1/admin/api-keys`:**

Request body:
```json
{
  "name": "CI/CD Deploy Key",
  "scopes": ["catalog.read", "logistics.manage"],
  "expires_at": "2027-04-23T00:00:00Z",
  "created_by": "usr_01JX..."
}
```

Response `201 Created`:
```json
{
  "key_id": "akey_01JX...",
  "plaintext_key": "umroh_k1_aXq9Rm3Z...",
  "key_prefix": "umroh_k1",
  "expires_at": "2027-04-23T00:00:00Z"
}
```

**Key generation:** `iam-svc` generates a cryptographically random 32-byte key, base58-encodes it, prefixes with `umroh_k1_`. The first 8 characters (prefix segment) are stored in `key_prefix` for display. The full plaintext key is returned in `plaintext_key` and then discarded server-side.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `missing_name` | `INVALID_ARGUMENT` | `name` is empty |
| `400` | `invalid_expires_at` | `INVALID_ARGUMENT` | `expires_at` is not a valid RFC3339 string |
| `400` | `expires_in_past` | `INVALID_ARGUMENT` | `expires_at` is before `NOW()` |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `iam.admin` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `IamService/RevokeAPIKey` (BL-IAM-014)

**Purpose:** revoke an API key by setting `revoked_at`. Revoked keys are rejected by all services. Idempotent — revoking an already-revoked key returns the existing `revoked_at` timestamp.

**Proto-style signature:**

```protobuf
rpc RevokeAPIKey(RevokeAPIKeyRequest) returns (RevokeAPIKeyResponse);

message RevokeAPIKeyRequest {
  string key_id = 1; // akey_... ULID
}

message RevokeAPIKeyResponse {
  string key_id     = 1;
  string revoked_at = 2; // RFC3339
}
```

**REST mapping — `DELETE /v1/admin/api-keys/{id}`:**

Response `200 OK`:
```json
{
  "key_id": "akey_01JX...",
  "revoked_at": "2026-04-23T15:00:00Z"
}
```

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `iam.admin` permission |
| `404` | `not_found` | `NOT_FOUND` | `key_id` does not exist |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `IamService/GetGlobalConfig` (BL-IAM-016)

**Purpose:** retrieve one, several, or all global configuration values from `iam.global_config`.

**Proto-style signature:**

```protobuf
rpc GetGlobalConfig(GetGlobalConfigRequest) returns (GetGlobalConfigResponse);

message GetGlobalConfigRequest {
  repeated string keys = 1; // empty = return all keys
}

message GetGlobalConfigResponse {
  repeated ConfigEntry configs = 1;
}

message ConfigEntry {
  string key         = 1;
  string value       = 2;
  string description = 3;
  string updated_at  = 4; // RFC3339
}
```

**REST mapping — `GET /v1/admin/config?keys=key1,key2`:**

Query parameters:

| Parameter | Type | Required | Notes |
| --- | --- | --- | --- |
| `keys` | comma-separated string | no | Omit to return all keys |

Response `200 OK`:
```json
{
  "configs": [
    {
      "key": "midtrans.merchant_id",
      "value": "G123456",
      "description": "Midtrans merchant ID for VA issuance",
      "updated_at": "2026-04-01T00:00:00Z"
    }
  ]
}
```

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `iam.admin` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `IamService/SetGlobalConfig` (BL-IAM-016)

**Purpose:** create or update a global config key (upsert). Stores `updated_by` and timestamps the change.

**Proto-style signature:**

```protobuf
rpc SetGlobalConfig(SetGlobalConfigRequest) returns (SetGlobalConfigResponse);

message SetGlobalConfigRequest {
  string key         = 1;
  string value       = 2;
  string description = 3; // optional — ignored if empty on update
  string updated_by  = 4; // user_id
}

message SetGlobalConfigResponse {
  string key        = 1;
  string value      = 2;
  string updated_at = 3; // RFC3339
}
```

**REST mapping — `PUT /v1/admin/config/{key}`:**

Request body:
```json
{
  "value": "G999999",
  "description": "Updated Midtrans merchant ID",
  "updated_by": "usr_01JX..."
}
```

Response `200 OK`:
```json
{
  "key": "midtrans.merchant_id",
  "value": "G999999",
  "updated_at": "2026-04-23T15:00:00Z"
}
```

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `missing_key` | `INVALID_ARGUMENT` | `key` is empty |
| `400` | `missing_value` | `INVALID_ARGUMENT` | `value` is empty |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `iam.admin` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

---

### RPC — `IamService/SearchActivityLog` (BL-IAM-011)

**Purpose:** paginated, multi-filter search over `iam.audit_logs` (the table contracted in S1/S2 `§ Audit trail`). Supports cursor-based pagination.

**Proto-style signature:**

```protobuf
rpc SearchActivityLog(SearchActivityLogRequest)
    returns (SearchActivityLogResponse);

message SearchActivityLogRequest {
  string user_id     = 1; // optional filter
  string resource    = 2; // optional filter — e.g. "booking", "invoice"
  string action      = 3; // optional filter — e.g. "create", "update_status"
  string from        = 4; // optional — RFC3339 lower bound on created_at
  string to          = 5; // optional — RFC3339 upper bound on created_at
  int32  limit       = 6; // default 50; max 200
  string cursor      = 7; // opaque pagination cursor; empty = first page
}

message SearchActivityLogResponse {
  repeated ActivityLogEntry logs        = 1;
  string                    next_cursor = 2; // empty if last page
}

message ActivityLogEntry {
  string id          = 1;
  string user_id     = 2;
  string resource    = 3;
  string action      = 4;
  string resource_id = 5;
  string created_at  = 6; // RFC3339
}
```

**REST mapping — `GET /v1/admin/activity-log`:**

Query parameters:

| Parameter | Type | Required | Default | Notes |
| --- | --- | --- | --- | --- |
| `user_id` | string | no | — | Filter by user |
| `resource` | string | no | — | Filter by resource type |
| `action` | string | no | — | Filter by action |
| `from` | string (RFC3339) | no | — | Lower bound on `created_at` (inclusive) |
| `to` | string (RFC3339) | no | — | Upper bound on `created_at` (inclusive) |
| `limit` | integer | no | `50` | Max `200`; silently clamped |
| `cursor` | string | no | — | Opaque cursor from previous response |

Response `200 OK`:
```json
{
  "logs": [
    {
      "id": "aud_01JX...",
      "user_id": "usr_01JX...",
      "resource": "booking",
      "action": "update_status",
      "resource_id": "bkg_01JX...",
      "created_at": "2026-04-23T10:00:00Z"
    }
  ],
  "next_cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNi0wNC0yM..."
}
```

**Ordering:** `created_at` DESC, `id` DESC. Cursor encodes `(created_at, id)` of the last returned entry. Opaque base64 — clients MUST NOT parse or construct.

**Failure codes:**

| HTTP | `error.code` | gRPC code | When |
| --- | --- | --- | --- |
| `400` | `invalid_date_format` | `INVALID_ARGUMENT` | `from` or `to` is not valid RFC3339 |
| `400` | `invalid_date_range` | `INVALID_ARGUMENT` | `from` is after `to` |
| `400` | `invalid_cursor` | `INVALID_ARGUMENT` | `cursor` is malformed or tampered |
| `400` | `invalid_limit` | `INVALID_ARGUMENT` | `limit` is not a positive integer |
| `401` | `unauthorized` | `UNAUTHENTICATED` | Missing or invalid bearer token |
| `403` | `forbidden` | `PERMISSION_DENIED` | Valid token but lacks `iam.admin` permission |
| `500` | `internal_error` | `INTERNAL` | DB failure or unexpected error |

### Honored by implementation

- `S6-E-05` — `iam-svc`: `SetDataScope` handler, upsert on `iam.data_scopes`, validation for `branch` scope.
- `S6-E-05` — `iam-svc`: `CreateAPIKey` handler, crypto-random key generation, argon2id hash storage, plaintext returned once.
- `S6-E-05` — `iam-svc`: `RevokeAPIKey` handler, idempotent `revoked_at` set.
- `S6-E-05` — `iam-svc`: `GetGlobalConfig` / `SetGlobalConfig` handlers, upsert on `iam.global_config`.
- `S6-E-05` — `iam-svc`: `SearchActivityLog` handler, keyset pagination over `iam.audit_logs`, multi-filter WHERE composition.
- `S6-E-05` — `gateway-svc`: six new route registrations, all gated with `RequirePermission("iam.admin")`.

---

## § ID prefixes (S6 additions)

All S6 entity IDs are **ULID** strings with a type prefix. New prefixes introduced in S6:

| Prefix | Entity | Service |
| --- | --- | --- |
| `vis_` | Visa application (`visa.visa_applications.id`) | `visa-svc` |
| `scan_` | Scan event (`ops.scan_events.id`) | `ops-svc` |
| `bbd_` | Bus boarding (`ops.bus_boardings.id`) | `ops-svc` |
| `pr_` | Purchase request (`logistics.purchase_requests.id`) | `logistics-svc` |
| `kit_` | Kit assembly (`logistics.kit_assemblies.id`) | `logistics-svc` |
| `dis_` | Disbursement batch (`finance.disbursement_batches.id`) | `finance-svc` |
| `dsc_` | Data scope (`iam.data_scopes.id`) | `iam-svc` |
| `akey_` | API key (`iam.api_keys.id`) | `iam-svc` |

Consumers treat all IDs as opaque strings.

---

## § Error envelope (S6 additions)

All S6 REST error responses use the shared envelope established in S1–S5:

```json
{
  "error": {
    "code": "<snake_case>",
    "message": "<human-readable, id-ID>",
    "trace_id": "<otel_span_hex>"
  }
}
```

Where a failure involves multiple entities (e.g., `not_all_ready` in `BulkSubmit`, `passport_expiry_violation`), `error` carries an additional `details` array:

```json
{
  "error": {
    "code": "passport_expiry_violation",
    "message": "Beberapa paspor jamaah akan kadaluarsa dalam 180 hari sebelum keberangkatan",
    "trace_id": "...",
    "details": [
      {
        "jamaah_id": "jmh_01JX...",
        "expiry_date": "2026-08-01",
        "days_before_departure": 46
      }
    ]
  }
}
```

`trace_id` is the OTel span ID per `docs/04-backend-conventions/03-logging-and-tracing.md`.

---

## § Audit trail (S6 additions)

Every state-changing call in S6 MUST emit an `iam.audit_logs` row via `iam.v1.IamService/RecordAudit` (consistent with S2 audit trail convention). Minimum audit events:

| Service | Trigger | `resource` | `action` | `user_id` |
| --- | --- | --- | --- | --- |
| `visa-svc` | `TransitionStatus` (non-idempotent) | `visa_application` | `update_status` | `actor_user_id` from request |
| `visa-svc` | `BulkSubmit` | `visa_application` | `bulk_submit` | `""` (system) |
| `ops-svc` | `RecordScan` (non-idempotent) | `scan_event` | `create` | `scanned_by` |
| `ops-svc` | `RecordBusBoarding` (non-idempotent) | `bus_boarding` | `create` | `scanned_by` |
| `logistics-svc` | `CreatePurchaseRequest` | `purchase_request` | `create` | `requested_by` |
| `logistics-svc` | `ApprovePurchaseRequest` | `purchase_request` | `decision` | `approved_by` |
| `logistics-svc` | `RecordGRNWithQC` | `grn` | `qc_record` | from bearer token |
| `logistics-svc` | `CreateKitAssembly` (non-idempotent) | `kit_assembly` | `create` | `assembled_by` |
| `finance-svc` | `CreateDisbursementBatch` | `disbursement_batch` | `create` | `created_by` |
| `finance-svc` | `ApproveDisbursement` | `disbursement_batch` | `decision` | `approved_by` |
| `iam-svc` | `SetDataScope` | `data_scope` | `upsert` | from bearer token |
| `iam-svc` | `CreateAPIKey` | `api_key` | `create` | `created_by` |
| `iam-svc` | `RevokeAPIKey` | `api_key` | `revoke` | from bearer token |
| `iam-svc` | `SetGlobalConfig` | `global_config` | `upsert` | `updated_by` |

---

## § Deferred items

The following Phase 6 items are **out of scope for this contract** — they require external integrations or dedicated delivery sprints. Placeholders are recorded here for roadmap tracking.

### 6.A — Agent Onboarding & Commission Wallet

**Status:** Deferred to Phase 6.H / 6.O sprint.

**Reason:** requires WA Business API integration for onboarding flow and a KYC provider contract (not yet selected). Commission wallet (BL-CRM-012) requires a `payment-svc` commission ledger extension that is out of scope for MH-V1.

**Future scope:** agent registration form → KYC submission → WA Business API approval notification → commission wallet activation.

### 6.L — B2C Frontend

**Status:** Deferred — dedicated frontend sprint required.

**Reason:** all Svelte 5 / SvelteKit frontend work for B2C jamaah-facing flows (package browsing, booking, payment checkout enhancements, visa status tracker) requires a dedicated Lutfi frontend sprint. Backend contracts in §§ 6.E / 6.D / 6.C / 6.B / 6.M above provide the REST surface that the frontend will consume.

### 6.F / 6.I — Dashboard Widgets

**Status:** Deferred — requires backend aggregation services + frontend widgets.

**Reason:** executive dashboard (departure occupancy, revenue trend, visa pipeline gauge) depends on a `reporting-svc` aggregation layer not yet designed. Frontend widget components also require the 6.L frontend sprint.

### 6.H — CRM Depth

**Status:** Deferred to Phase 6.H sprint.

**Reason:** WhatsApp Business API (meta BSP), social media integrations, and LMS platform require separate vendor contracts and API onboarding.

### 6.N / 6.O — Pilgrim Daily App

**Status:** Deferred — separate mobile delivery.

**Reason:** native mobile app (pilgrim prayer times, qibla, daily guidance), offline capabilities, and push notification service are a separate delivery from the ERP web platform. Requires dedicated mobile sprint and cross-platform testing.

---

## § Changelog

- **2026-04-23** — Initial draft. Covers: § 6.E Visa Pipeline — `VisaService/TransitionStatus`, `BulkSubmit`, `GetApplications` with full state machine, DB schema (`visa.visa_applications`, `visa.status_history`), gateway routes (`S6-E-01`); § 6.D Field Operations Scan — `OpsService/RecordScan`, `RecordBusBoarding`, `GetBoardingRoster`, DB schema (`ops.scan_events`, `ops.bus_boardings`), gateway routes (`S6-E-02`); § 6.C Warehouse & Procurement — `LogisticsService/CreatePurchaseRequest`, `ApprovePurchaseRequest`, `RecordGRNWithQC`, `CreateKitAssembly`, DB schema (`logistics.purchase_requests`, `logistics.kit_assemblies`, `logistics.kit_assembly_items`), gateway routes (`S6-E-03`); § 6.B Finance AP/AR Depth — `FinanceService/CreateDisbursementBatch`, `ApproveDisbursement`, `GetARAPAging`, DB schema (`finance.disbursement_batches`, `finance.disbursement_items`), gateway routes (`S6-E-04`); § 6.M IAM Admin/Security Depth — `IamService/SetDataScope`, `CreateAPIKey`, `RevokeAPIKey`, `GetGlobalConfig`, `SetGlobalConfig`, `SearchActivityLog`, DB schema (`iam.data_scopes`, `iam.api_keys`, `iam.global_config`), gateway routes (`S6-E-05`); § Deferred placeholders for 6.A, 6.L, 6.F/6.I, 6.H, 6.N/6.O.
