---
slice: S4
title: Slice S4 — Integration Contract
status: draft
last_updated: 2026-04-23
pr_owner: Lutfi
reviewer: Elda
task_codes:
  - S4-J-01
  - S4-J-02
---

# Slice S4 — Integration Contract

> Slice S4 = "Basic growth loop" — lead capture with UTM attribution → CS round-robin assignment → booking conversion → CRM events fan-out. This file is the wire-level agreement between `crm-svc`, `booking-svc`, and the gateway for the lead-to-booking user journey.
>
> **Incremental build:** only sections that have a landed `S4-J-*` card are filled. Amendments after merge follow the Bump-versi rule in `docs/contracts/README.md` (Changelog append for additive, `-v2.md` for breaking).

## Purpose

S4 is the growth-loop slice. The flow: a lead arrives (organic, WhatsApp, Instagram, Facebook, referral, agent, or direct) → UTM parameters are captured immutably → crm-svc assigns a CS via round-robin → CS nurtures and optionally converts the lead to a booking draft → booking events fan back to crm-svc for attribution and status sync. Services bound in S4: `crm-svc` (owns leads, assignment, attribution), `booking-svc` (downstream consumer for conversion, upstream event emitter for paid/created events), `gateway` (public POST /v1/leads endpoint).

S4 does **not** introduce a message queue or Temporal workflow. All service-to-service propagation is **direct gRPC calls** per ADR-0006.

## Scope

**In scope for S4 contracts:**

- `leads` table schema with UTM snapshot fields — `S4-J-01`.
- UTM capture rules and attribution snapshot — `S4-J-01` (resolves Q019 + Q057).
- REST endpoints: `POST /v1/leads`, `GET /v1/leads`, `GET /v1/leads/{id}`, `PUT /v1/leads/{id}`, `POST /v1/leads/{id}/convert` — `S4-J-01`.
- CS round-robin auto-assignment algorithm — `S4-J-01`.
- Lead status state machine — `S4-J-01`.
- `booking-svc → crm-svc.OnBookingCreated` gRPC contract — `S4-J-02`.
- `booking-svc → crm-svc.OnBookingPaidInFull` gRPC contract — `S4-J-02`.
- Attribution copy rule: lead UTM snapshot → booking on conversion — `S4-J-02`.

**Out of scope for S4 contracts (deferred to Phase 6 or later):**

- Commission calculation and override chain (`CalculateCommission` — Phase 6 `BL-CRM-012`).
- Agent onboarding, E-KYC, E-Signature, replica site (Phase 6 `S4-L-02`).
- Lead nurturing drip campaigns, broadcast hub, bot filter (Phase 6 `6.A`).
- SLA timer auto-reassignment on breach (Q066 default 10 min — out of MVP S4 scope).
- Lead transfer / ownership handoff (Q064 — Phase 6).
- Referral code attribution (Phase 6 `BL-CRM-017+`).
- Full commission wallet and `commission_ledger` (Phase 6 `BL-CRM-012`).
- Ads Manager Lite, ROAS, retargeting, landing-page builder (Phase 6 `6.A`).
- Alumni referral, ZISWAF, community threads (Phase 6 `6.H`).

---

## § S4-J-01 — Lead Schema + UTM Snapshot + Attribution Contract

*(Landed with `S4-J-01`.)*

### Background: lead capture flow

A lead enters crm-svc through one of two surfaces:

1. **Public REST endpoint** `POST /v1/leads` — called by the gateway when a calon jamaah submits a public form or clicks a tracked WhatsApp link on a landing page or replica site. No bearer token required.
2. **CS manual entry** (same endpoint, authenticated) — CS creates a lead after an inbound WhatsApp or phone call.

On creation, crm-svc:

1. Validates and stores the lead row with the UTM snapshot from the request.
2. Auto-assigns a CS via round-robin algorithm (see below).
3. Returns the created lead.

```
calon jamaah / agent site
  │
  │  POST /v1/leads  (public, no auth)
  │  { name, phone, source, utm_*, interest_package_id? }
  ▼
gateway → crm-svc
  │
  ├── INSERT leads (with utm_* snapshot, status='new')
  ├── Round-robin assign → leads.assigned_cs_id
  └── 201 Created { lead }
```

### Lead table schema

```sql
CREATE TABLE leads (
  id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  source                TEXT        NOT NULL
                          CHECK (source IN (
                            'organic','whatsapp','instagram',
                            'facebook','referral','agent','direct'
                          )),

  -- UTM snapshot (immutable after INSERT)
  utm_source            TEXT,
  utm_medium            TEXT,
  utm_campaign          TEXT,
  utm_content           TEXT,
  utm_term              TEXT,

  -- Lead identity
  name                  TEXT        NOT NULL,
  phone                 TEXT        NOT NULL,     -- E.164 format, e.g. +6281234567890
  email                 TEXT,

  -- Interest
  interest_package_id   UUID        REFERENCES catalog.packages(id),
  interest_departure_id UUID        REFERENCES catalog.departures(id),

  -- Lifecycle
  status                TEXT        NOT NULL DEFAULT 'new'
                          CHECK (status IN (
                            'new','contacted','qualified','converted','lost'
                          )),

  -- Assignment
  assigned_cs_id        UUID        REFERENCES iam.users(id),

  -- Notes
  notes                 TEXT,

  -- Timestamps
  created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for round-robin assignment query
CREATE INDEX leads_assigned_cs_status_idx
  ON leads (assigned_cs_id, status)
  WHERE status NOT IN ('converted', 'lost');

-- Phone uniqueness is advisory (warn on duplicate, not rejected) — dedup via fingerprint
CREATE INDEX leads_phone_idx ON leads (phone);
```

**Field definitions:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `id` | UUID | auto | PK. UUID v4. No prefix convention (differs from S1–S3 ULID pattern — crm-svc uses UUID; cross-ref with crm-svc data-model doc if prefix desired). |
| `source` | enum text | yes | Channel that produced the lead: `organic \| whatsapp \| instagram \| facebook \| referral \| agent \| direct`. |
| `utm_source` | text | no | Raw UTM parameter from URL at lead capture time. NULL if not present. |
| `utm_medium` | text | no | Raw UTM parameter. NULL if not present. |
| `utm_campaign` | text | no | Raw UTM parameter. NULL if not present. |
| `utm_content` | text | no | Raw UTM parameter. Often carries `<agent_code>` for agent-sourced leads. NULL if not present. |
| `utm_term` | text | no | Raw UTM parameter. NULL if not present. |
| `name` | text | yes | Lead's full name. Min 2 chars. |
| `phone` | text | yes | E.164 format. Validated by crm-svc (must start with `+`, digits only after). |
| `email` | text | no | Optional. Validated format if present. |
| `interest_package_id` | UUID | no | FK to catalog package. Nullable — calon jamaah may not have selected a package yet. |
| `interest_departure_id` | UUID | no | FK to catalog departure slot. Nullable. |
| `status` | enum text | yes | Lifecycle status. Default `new`. Transitions enforced by service layer (see state machine below). |
| `assigned_cs_id` | UUID | no | FK to IAM user with role `CS`. Set by round-robin on creation. Nullable if no active CS available. |
| `notes` | text | no | CS-authored notes on the lead. Plain text. |
| `created_at` | timestamptz | auto | Set on INSERT. Never updated. |
| `updated_at` | timestamptz | auto | Updated on every mutation. |

### UTM snapshot rules (Q019 + Q057 resolution)

**Rule 1 — Immutable snapshot at lead creation**

UTM parameters are captured **once** at the moment a lead is created (`POST /v1/leads`). The `utm_*` fields in the `leads` row are set on INSERT and never updated by any subsequent event (form re-submit, WA re-click, drip campaign open, etc.).

This implements the "last-touch at lead-entry" snapshot: the UTM present when the lead first contacts the system is stored as-is. It reflects the last tracked touchpoint that actually caused the lead to convert to a form submission or first contact.

**Rule 2 — Attribution at booking conversion**

When a lead is converted to a booking (via `POST /v1/leads/{id}/convert` or via `booking-svc.OnBookingCreated` with a `lead_id`), the booking row copies the lead's UTM snapshot:

```
bookings.utm_source    ← leads.utm_source
bookings.utm_medium    ← leads.utm_medium
bookings.utm_campaign  ← leads.utm_campaign
bookings.utm_content   ← leads.utm_content
bookings.utm_term      ← leads.utm_term
```

This copy is performed by `booking-svc` when the booking is created. The `booking.utm_*` fields must be added to the `bookings` table schema (backlog `BL-CRM-001`, `S4-E-02`). The lead's own UTM snapshot remains unchanged.

**Rule 3 — Attribution model (Q057 decision)**

MVP uses **last-touch** for commission routing: the UTM snapshot on the booking at conversion time is the authoritative attribution. Both `first_touch_utm` and `last_touch_utm` are intended per Q057 Option C for ROAS analytics, but in S4 MVP only the lead-entry snapshot is stored (single touch). Full dual-capture (first + last) is a Phase 6 enhancement (`BL-CRM-002`).

**Rule 4 — Attribution window (Q019 + Q057 — 30 days)**

Attribution is valid for **30 days** from `leads.created_at`. If a lead converts more than 30 days after creation, the attribution is treated as organic (no agent commission) for ROAS purposes. The UTM snapshot is still copied to the booking for audit, but the commission routing engine (Phase 6) must enforce the 30-day cap.

**Rule 5 — Direct visit does not erase agent touch (Q057)**

If a lead was created via `utm_content=<agent_code>` (agent replica site) and the same calon jamaah later visits the central site directly, the original UTM snapshot is retained. crm-svc does not update `utm_*` fields after the initial insert; the lead deduplication fingerprint (`SHA256(lower(phone))`) identifies the returning visitor and the existing lead record is surfaced rather than creating a duplicate.

**Rule 6 — UTM validation / sanitization**

crm-svc truncates each `utm_*` field to **255 characters** on ingest. Unknown or malformed values are stored as-is (no rejection). Agent code validation (`utm_content` cross-checked against `agents.agent_code`) is a Phase 6 feature — MVP stores the raw value.

### Lead status state machine

```
             POST /v1/leads
                  │
                  ▼
               [ new ]
                  │
                  │  CS first contact (PUT status=contacted)
                  ▼
            [ contacted ]
                  │
                  │  CS qualifies interest (PUT status=qualified)
                  ▼
            [ qualified ]
                  │         │
                  │         │  unqualified / no interest
                  │         │  (PUT status=lost from any active state)
                  │         ▼
                  │      [ lost ] ──────────────────── (terminal)
                  │
                  │  POST /v1/leads/{id}/convert
                  │  or booking-svc.OnBookingCreated(lead_id)
                  ▼
           [ converted ] ─────────────────────────── (terminal)
```

**Allowed transitions (enforced by crm-svc service layer):**

| From | To | Trigger |
| --- | --- | --- |
| `new` | `contacted` | CS updates lead via `PUT /v1/leads/{id}` |
| `new` | `lost` | CS marks as lost |
| `contacted` | `qualified` | CS qualifies |
| `contacted` | `lost` | CS marks as lost |
| `qualified` | `converted` | `POST /v1/leads/{id}/convert` or `OnBookingCreated` with `lead_id` |
| `qualified` | `lost` | CS marks as lost |
| `converted` | _(none)_ | Terminal — no further transitions |
| `lost` | _(none)_ | Terminal — no further transitions |

**Constraint:** `converted` and `lost` are terminal. crm-svc returns `FAILED_PRECONDITION lead_already_terminal` for any transition attempt from these states.

> Note: F10-W3 defines a richer `cold → warm → hot → converted → lost` lifecycle (per PRD). The S4 MVP states above (`new → contacted → qualified → converted → lost`) are the **contracted engineering states** that map to the PRD lifecycle as follows: `new` ~ `cold`; `contacted` ~ `warm`; `qualified` ~ `hot`. Tags (Tanya / Janji Bayar / Closed / No Response / Bot) from PRD module #36 are stored in a `lead_tags` array column (Phase 6 `BL-CRM-001`) — not contracted in S4.

### CS round-robin auto-assignment

**Trigger:** every `POST /v1/leads` call (regardless of whether caller is authenticated).

**Algorithm:**

1. Query all IAM users with role `CS` and `status = active`.
2. For each active CS, count their currently **active leads**: `COUNT(*) FROM leads WHERE assigned_cs_id = cs.id AND status NOT IN ('converted', 'lost')`.
3. Assign the lead to the CS with the lowest active-lead count.
4. **Tie-breaking:** if two or more CS users have the same minimum active-lead count, assign to the CS user with the earliest `users.created_at` (senior-first tie-break).
5. If no active CS users exist, set `assigned_cs_id = NULL` (lead enters the queue unassigned; an admin must assign manually).

**Implementation note:** crm-svc must hold a short advisory lock or use a `SELECT FOR UPDATE SKIP LOCKED` pattern on the assignment query to prevent two concurrent lead creations both selecting the same CS. A row-level lock on a `cs_assignment_state` auxiliary table (keyed by `cs_user_id`) is the recommended approach for MVP.

**Acceptance criterion (F10-W4, BL-CRM-003):** no CS receives more than `⌈average_active_leads + 2⌉` leads in any 24-hour window (eventual consistency on the count is acceptable; the constraint is a fairness bound, not a hard cap per request).

### REST endpoints

#### `POST /v1/leads` — Create lead (public, no auth required)

**Request:**

```json
{
  "source": "instagram",
  "utm_source": "instagram",
  "utm_medium": "cpc",
  "utm_campaign": "ramadhan-2026",
  "utm_content": "agt_LUTFI01",
  "utm_term": null,
  "name": "Ahmad Fauzan",
  "phone": "+6281234567890",
  "email": "ahmad@example.com",
  "interest_package_id": "550e8400-e29b-41d4-a716-446655440000",
  "interest_departure_id": null
}
```

**Request field rules:**

| Field | Required | Validation |
| --- | --- | --- |
| `source` | yes | One of: `organic \| whatsapp \| instagram \| facebook \| referral \| agent \| direct` |
| `utm_source` | no | Max 255 chars |
| `utm_medium` | no | Max 255 chars |
| `utm_campaign` | no | Max 255 chars |
| `utm_content` | no | Max 255 chars |
| `utm_term` | no | Max 255 chars |
| `name` | yes | Min 2 chars, max 255 chars |
| `phone` | yes | E.164 format (`+` followed by 7–15 digits). crm-svc normalizes to E.164 |
| `email` | no | Valid email format if present |
| `interest_package_id` | no | Valid UUID if present; FK existence checked |
| `interest_departure_id` | no | Valid UUID if present; FK existence checked |

**Response `201 Created`:**

```json
{
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "source": "instagram",
    "utm_source": "instagram",
    "utm_medium": "cpc",
    "utm_campaign": "ramadhan-2026",
    "utm_content": "agt_LUTFI01",
    "utm_term": null,
    "name": "Ahmad Fauzan",
    "phone": "+6281234567890",
    "email": "ahmad@example.com",
    "interest_package_id": "550e8400-e29b-41d4-a716-446655440000",
    "interest_departure_id": null,
    "status": "new",
    "assigned_cs_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "notes": null,
    "created_at": "2026-04-23T10:15:00Z",
    "updated_at": "2026-04-23T10:15:00Z"
  }
}
```

**Failure codes:**

| HTTP | `error.code` | When |
| --- | --- | --- |
| `400` | `invalid_request` | Missing required fields, invalid `source` enum, invalid phone format, invalid email format, malformed UUID fields |
| `404` | `package_not_found` | `interest_package_id` does not exist in catalog |
| `404` | `departure_not_found` | `interest_departure_id` does not exist in catalog |
| `500` | `internal_error` | DB failure or unexpected error |

> **No deduplication rejection:** a second lead with the same phone is accepted and stored (no `409`). crm-svc logs a warning and sets a `dedup_fingerprint` for future matching (Phase 6 `BL-CRM-001`). CS will see both leads in list view and can merge manually.

---

#### `GET /v1/leads` — List leads with filter (bearer, CS role required)

**Required header:** `Authorization: Bearer <token>` (IAM JWT, role `CS` or `admin`).

**Query parameters:**

| Parameter | Type | Description |
| --- | --- | --- |
| `status` | string | Filter by `status` enum value |
| `assigned_cs_id` | UUID | Filter by assigned CS (CS can only filter own leads unless admin) |
| `source` | string | Filter by `source` enum value |
| `from` | RFC3339 date | Filter `created_at >= from` |
| `to` | RFC3339 date | Filter `created_at <= to` |
| `page` | int | 1-indexed page number (default 1) |
| `per_page` | int | Results per page (default 20, max 100) |

**Response `200 OK`:**

```json
{
  "data": [
    {
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "source": "instagram",
      "utm_source": "instagram",
      "utm_medium": "cpc",
      "utm_campaign": "ramadhan-2026",
      "utm_content": "agt_LUTFI01",
      "utm_term": null,
      "name": "Ahmad Fauzan",
      "phone": "+6281234567890",
      "email": "ahmad@example.com",
      "interest_package_id": "550e8400-e29b-41d4-a716-446655440000",
      "interest_departure_id": null,
      "status": "new",
      "assigned_cs_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "notes": null,
      "created_at": "2026-04-23T10:15:00Z",
      "updated_at": "2026-04-23T10:15:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 142
  }
}
```

**RBAC rule:** a user with role `CS` may only retrieve leads where `assigned_cs_id = caller.user_id` OR leads with `assigned_cs_id = NULL`. Admin and supervisor roles may retrieve all leads.

---

#### `GET /v1/leads/{id}` — Lead detail (bearer required)

**Required header:** `Authorization: Bearer <token>`.

**Path param:** `id` — UUID of the lead.

**Response `200 OK`:** same shape as a single lead object in the list response above.

**Failure codes:**

| HTTP | `error.code` | When |
| --- | --- | --- |
| `401` | `unauthorized` | Missing or invalid bearer token |
| `403` | `forbidden` | CS user requesting a lead not assigned to them |
| `404` | `lead_not_found` | Lead with given UUID does not exist |

---

#### `PUT /v1/leads/{id}` — Update lead status / notes / assignment (bearer, CS role)

**Required header:** `Authorization: Bearer <token>` (role `CS` or `admin`).

**Request (all fields optional — PATCH semantics on a PUT endpoint; only provided fields are updated):**

```json
{
  "status": "contacted",
  "notes": "Interested in Ramadhan package, asked about payment installment.",
  "assigned_cs_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "interest_package_id": "550e8400-e29b-41d4-a716-446655440000",
  "interest_departure_id": "660e8400-e29b-41d4-a716-446655440111"
}
```

**Mutable fields via `PUT /v1/leads/{id}`:**

| Field | Mutable | Notes |
| --- | --- | --- |
| `status` | yes | Must follow allowed state-machine transitions; terminal states reject update |
| `notes` | yes | Full replace (not append) |
| `assigned_cs_id` | yes (admin only) | CS cannot reassign to another CS; admin can |
| `interest_package_id` | yes | Can update or set to null |
| `interest_departure_id` | yes | Can update or set to null |
| `utm_*` fields | **NO** | Immutable after INSERT — crm-svc returns `400 invalid_request utm_immutable` if caller attempts to update any `utm_*` field |
| `phone` | **NO** | Phone identity is immutable; changes require a new lead |
| `source` | **NO** | Source is immutable |

**Response `200 OK`:** updated lead object (same shape as GET response).

**Failure codes:**

| HTTP | `error.code` | When |
| --- | --- | --- |
| `400` | `invalid_transition` | Status transition not allowed by state machine |
| `400` | `lead_already_terminal` | Lead is in `converted` or `lost` — no further updates |
| `400` | `utm_immutable` | Caller attempted to modify `utm_*` fields |
| `400` | `invalid_request` | Malformed field values |
| `403` | `forbidden` | CS attempting to reassign to another CS (admin-only operation) |
| `404` | `lead_not_found` | Lead UUID does not exist |

---

#### `POST /v1/leads/{id}/convert` — Convert lead to booking draft (bearer, CS role)

**Required header:** `Authorization: Bearer <token>` (role `CS` or `admin`).

**Path param:** `id` — UUID of the lead to convert.

**Request body:**

```json
{
  "package_id": "550e8400-e29b-41d4-a716-446655440000",
  "departure_id": "660e8400-e29b-41d4-a716-446655440111",
  "jamaah_count": 2
}
```

**Request field rules:**

| Field | Required | Notes |
| --- | --- | --- |
| `package_id` | yes | UUID of the catalog package to book |
| `departure_id` | yes | UUID of the departure slot; must belong to `package_id` |
| `jamaah_count` | yes | Integer ≥ 1 |

**Behavior:**

1. crm-svc validates lead exists and is NOT in `converted` or `lost` state.
2. crm-svc calls `booking-svc.CreateBookingDraft` gRPC with `{ lead_id, package_id, departure_id, jamaah_count, utm_* }` (internal gRPC — shape contracted in booking-svc API doc, not here).
3. booking-svc creates the booking draft, copies UTM from lead, returns `booking_id`.
4. crm-svc transitions lead status → `converted`.
5. crm-svc returns the `booking_id` in response.

**Response `200 OK`:**

```json
{
  "data": {
    "lead_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "booking_id": "bkg_01JCDE...",
    "status": "converted"
  }
}
```

**Failure codes:**

| HTTP | `error.code` | When |
| --- | --- | --- |
| `400` | `lead_already_terminal` | Lead is already `converted` or `lost` |
| `400` | `invalid_request` | Missing required fields or invalid UUIDs |
| `404` | `lead_not_found` | Lead UUID does not exist |
| `404` | `package_not_found` | `package_id` does not exist |
| `404` | `departure_not_found` | `departure_id` does not exist or does not belong to package |
| `409` | `insufficient_seats` | No available seats for `departure_id` (propagated from booking-svc) |
| `500` | `internal_error` | Booking creation failed; lead status NOT changed (no partial state) |

> **Atomicity rule:** crm-svc transitions lead to `converted` ONLY after booking-svc confirms successful draft creation. If `CreateBookingDraft` returns an error, the lead status remains unchanged. crm-svc does NOT retry autonomously — the caller retries.

### `booking-svc` additions for attribution

booking-svc must add `utm_*` columns to the `bookings` table:

```sql
ALTER TABLE bookings ADD COLUMN utm_source   TEXT;
ALTER TABLE bookings ADD COLUMN utm_medium   TEXT;
ALTER TABLE bookings ADD COLUMN utm_campaign TEXT;
ALTER TABLE bookings ADD COLUMN utm_content  TEXT;
ALTER TABLE bookings ADD COLUMN utm_term     TEXT;
ALTER TABLE bookings ADD COLUMN lead_id      UUID; -- nullable FK to crm.leads
```

These columns are **write-once** on booking creation (same immutability as leads UTM snapshot). booking-svc service layer rejects updates to `utm_*` fields post-creation.

### Honored by implementation

- `S4-E-02` (`BL-CRM-001`) — crm-svc: lead table creation, `POST /v1/leads`, UTM capture, CS round-robin assignment.
- `S4-E-02` (`BL-CRM-002`) — crm-svc: UTM attribution copy to booking on convert, lead status transition `→ converted`.
- `S4-E-02` (`BL-CRM-003`) — crm-svc: CS assignment algorithm, tie-breaking, `GET /v1/leads` with RBAC filter, `PUT /v1/leads/{id}` state machine enforcement.
- `S4-E-02` — booking-svc: `utm_*` + `lead_id` columns added to `bookings` table; `CreateBookingDraft` accepts and stores these fields.

---

## § S4-J-02 — Booking → CRM Events Contract

*(Landed with `S4-J-02`.)*

S4 introduces two gRPC calls from `booking-svc` to `crm-svc` that keep the CRM synchronized with booking lifecycle events. Per ADR-0006, these are synchronous direct gRPC calls — no message broker.

### Background: call chain

After `booking-svc` processes a booking state change:

```
booking-svc
  │
  │  (on booking draft created)
  ├──► gRPC: crm.v1.CrmService/OnBookingCreated
  │         payload: { booking_id, lead_id?, package_id, departure_id, jamaah_count, created_at }
  │
  │  (on booking transitions to paid_in_full, per S3 MarkBookingPaid)
  └──► gRPC: crm.v1.CrmService/OnBookingPaidInFull
            payload: { booking_id, lead_id?, paid_at }
```

> **ADR-0006 compliance.** Both calls are synchronous within the respective handler. booking-svc does NOT fork goroutines — failure in either downstream MUST be surfaced to the caller (the `POST /v1/leads/{id}/convert` handler for `OnBookingCreated`, and the payment-svc → `MarkBookingPaid` pipeline for `OnBookingPaidInFull`). Both calls are idempotent so webhook retries are safe.

> **`OnBookingCreated` error behavior exception:** if `POST /v1/leads/{id}/convert` creates the booking successfully in booking-svc but `OnBookingCreated` to crm-svc fails, booking-svc must still return success to the REST caller (the booking is created). crm-svc will re-sync via `OnBookingPaidInFull` when payment lands. This relaxed behavior is acceptable for S4 MVP; a full saga or outbox pattern is Phase 6.

### gRPC — `crm.v1.CrmService/OnBookingCreated`

**Called by:** `booking-svc` when a booking draft is created (either via `POST /v1/leads/{id}/convert` or direct B2C booking creation without a lead).
**Served by:** `crm-svc`.

**Proto-style signature:**

```protobuf
rpc OnBookingCreated(OnBookingCreatedRequest) returns (OnBookingCreatedResponse);

message OnBookingCreatedRequest {
  string   booking_id    = 1; // ULID — bkg_ prefix
  string   lead_id       = 2; // UUID of the CRM lead; empty string if no lead (B2C direct)
  string   package_id    = 3; // UUID of the catalog package
  string   departure_id  = 4; // UUID of the departure slot
  int32    jamaah_count  = 5; // number of jamaah seats booked
  string   created_at    = 6; // RFC3339 — booking draft creation timestamp
}

message OnBookingCreatedResponse {
  bool replayed = 1; // true if booking_id was already recorded (idempotent replay)
}
```

**Field rules:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string (ULID) | yes | `bkg_` prefix. Used as idempotency key in crm-svc. |
| `lead_id` | string (UUID) | no | Empty string `""` if booking was created without a lead (B2C direct). crm-svc treats empty string as NULL. |
| `package_id` | string (UUID) | yes | crm-svc records for attribution and commission context. |
| `departure_id` | string (UUID) | yes | crm-svc records for attribution and commission context. |
| `jamaah_count` | int32 | yes | Must be ≥ 1. |
| `created_at` | string (RFC3339) | yes | Booking draft creation timestamp. |

**crm-svc behavior on `OnBookingCreated`:**

1. Check idempotency: if a `booking_id` record already exists in crm-svc's `booking_events` table, return `replayed = true` (no-op).
2. If `lead_id` is non-empty: record the booking association in `crm_booking_links { booking_id, lead_id, package_id, departure_id, jamaah_count, created_at }`.
3. If `lead_id` is empty: record the booking as unattributed (B2C direct) — still stored for reporting.
4. Do NOT transition lead status to `converted` here — lead was already transitioned by the `POST /v1/leads/{id}/convert` REST call (if lead exists). `OnBookingCreated` is informational for crm-svc bookkeeping only.

**Failure codes (gRPC `status.Code`):**

| Code | `error.code` | When |
| --- | --- | --- |
| `INVALID_ARGUMENT` | `invalid_request` | Missing `booking_id`, `package_id`, `departure_id`; `jamaah_count < 1`; malformed ULIDs/UUIDs; invalid `created_at` |
| `NOT_FOUND` | `lead_not_found` | `lead_id` is non-empty but does not exist in crm-svc leads table |
| `INTERNAL` | `internal_error` | DB failure or unexpected error |

**Idempotency:** on `booking_id` key. If the same `OnBookingCreated` is called twice, crm-svc returns `replayed = true` on the second call without creating duplicate records.

---

### gRPC — `crm.v1.CrmService/OnBookingPaidInFull`

**Called by:** `booking-svc` when a booking transitions to `paid_in_full` (i.e., after `MarkBookingPaid` with `invoice_status = "paid"`, as contracted in S3 `§ S3-J-01`).
**Served by:** `crm-svc`.

> This call is added to the **existing S3 fan-out chain** in `booking-svc.MarkBookingPaid`. After S4 lands, the fan-out becomes:
> ```
> booking-svc (on paid_in_full)
>   ├──► logistics-svc.OnBookingPaid    (S3-J-02)
>   ├──► finance-svc.OnPaymentReceived  (S3-J-03)
>   └──► crm-svc.OnBookingPaidInFull    (S4-J-02)  ← NEW
> ```
> All three calls remain synchronous and idempotent per ADR-0006.

**Proto-style signature:**

```protobuf
rpc OnBookingPaidInFull(OnBookingPaidInFullRequest) returns (OnBookingPaidInFullResponse);

message OnBookingPaidInFullRequest {
  string booking_id = 1; // ULID — bkg_ prefix
  string lead_id    = 2; // UUID of the CRM lead; empty string if no associated lead
  string paid_at    = 3; // RFC3339 — timestamp of full payment (= received_at of final payment event)
}

message OnBookingPaidInFullResponse {
  bool replayed = 1; // true if booking_id was already processed (idempotent replay)
}
```

**Field rules:**

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `booking_id` | string (ULID) | yes | `bkg_` prefix. Used as idempotency key. |
| `lead_id` | string (UUID) | no | Empty string `""` if booking has no lead association. crm-svc resolves from `crm_booking_links` if empty and `OnBookingCreated` was previously received. |
| `paid_at` | string (RFC3339) | yes | Timestamp of the final payment event that brought the invoice to `paid` status. |

**crm-svc behavior on `OnBookingPaidInFull`:**

1. Check idempotency: if `booking_id` already recorded as `paid_in_full` in `crm_booking_events`, return `replayed = true` (no-op).
2. Record a `crm_booking_events` row: `{ booking_id, event_type='paid_in_full', paid_at }`.
3. If `lead_id` is non-empty OR can be resolved from `crm_booking_links.booking_id`: transition associated lead status → `converted` (if not already `converted`). If lead is already `converted`, this is a no-op for the lead.
4. If `lead_id` is empty AND cannot be resolved from `crm_booking_links`: no lead update — booking is unattributed (B2C direct without lead origin).
5. Emit `crm.commission_confirmed` event context (Phase 6 `BL-CRM-012` will wire the actual commission calculation; S4 MVP just records the paid_at timestamp for future commission engine pickup).

**Failure codes (gRPC `status.Code`):**

| Code | `error.code` | When |
| --- | --- | --- |
| `INVALID_ARGUMENT` | `invalid_request` | Missing `booking_id`; malformed ULID/UUID; invalid `paid_at` format |
| `NOT_FOUND` | `lead_not_found` | `lead_id` provided but does not exist in crm-svc |
| `INTERNAL` | `internal_error` | DB failure or unexpected error |

**Idempotency:** on `booking_id` key. Duplicate calls (gateway retries) return `replayed = true` without side effects.

**Lead status transition safety:**

crm-svc applies the `converted` transition only if the lead's current status allows it (not already in `converted` or `lost`). If the lead is already `converted` (from the `POST /v1/leads/{id}/convert` call that happened earlier), `OnBookingPaidInFull` is effectively a no-op for lead status — it still records the `crm_booking_events` row and returns success.

If the lead is in `lost` status when `OnBookingPaidInFull` arrives (edge case: CS marked lead lost but the booking still completed via a different channel), crm-svc logs a warning and does NOT override the `lost` status. The booking event is recorded normally.

### `crm_booking_links` schema

```sql
CREATE TABLE crm_booking_links (
  id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id    TEXT        NOT NULL UNIQUE,   -- ULID, bkg_ prefix
  lead_id       UUID        REFERENCES leads(id),  -- nullable (B2C direct)
  package_id    UUID        NOT NULL,
  departure_id  UUID        NOT NULL,
  jamaah_count  INT         NOT NULL CHECK (jamaah_count >= 1),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE crm_booking_events (
  id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  booking_id    TEXT        NOT NULL,             -- ULID, bkg_ prefix
  event_type    TEXT        NOT NULL              -- 'booking_created' | 'paid_in_full'
                  CHECK (event_type IN ('booking_created', 'paid_in_full')),
  occurred_at   TIMESTAMPTZ NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (booking_id, event_type)                -- idempotency constraint
);
```

> The `UNIQUE (booking_id, event_type)` constraint on `crm_booking_events` is the DB-level idempotency gate. crm-svc uses `INSERT ... ON CONFLICT (booking_id, event_type) DO NOTHING` and checks whether a row was actually inserted to determine `replayed`.

### Updated fan-out in `booking-svc.MarkBookingPaid`

S4 adds one synchronous step to the fan-out contracted in S3 `§ S3-J-01`. The full sequence after S4 lands:

1. Atomically update `bookings.status = 'paid_in_full'` in DB.
2. **Synchronously** call `logistics-svc.OnBookingPaid` (S3-J-02).
3. **Synchronously** call `finance-svc.OnPaymentReceived` (S3-J-03).
4. **Synchronously** call `crm-svc.OnBookingPaidInFull` (S4-J-02). ← NEW
5. Return `MarkBookingPaidResponse` to `payment-svc`.

Steps 2–4 are all idempotent. If step 4 (`crm-svc`) fails with a non-retryable error, `MarkBookingPaid` returns `INTERNAL` to `payment-svc` (same behavior as a step 2 or 3 failure). The gateway will retry the webhook, and all three downstream calls will be replayed safely via their idempotency gates.

### `booking-svc` proto stub for crm-svc

booking-svc MUST NOT import `crm-svc`'s pb package directly. It keeps its own vendored stub at:

```
services/booking-svc/adapter/crm_grpc_adapter/pb/crm.proto
```

The stub contains ONLY the `OnBookingCreated` and `OnBookingPaidInFull` RPCs.

### Honored by implementation

- `S4-E-02` (`BL-CRM-001`) — crm-svc: `OnBookingCreated` handler, `crm_booking_links` table, idempotency gate.
- `S4-E-02` (`BL-CRM-002`) — crm-svc: `OnBookingPaidInFull` handler, `crm_booking_events` table, lead `→ converted` transition on payment, idempotency gate.
- `S4-E-02` — booking-svc: `crm_grpc_adapter` implementation; fan-out call `OnBookingCreated` on booking draft creation; fan-out call `OnBookingPaidInFull` on `paid_in_full` (added to `MarkBookingPaid` handler).

---

## § Lead Status (S4 additions)

S4 introduces the `leads` entity and its status lifecycle. For completeness, the booking status columns that S4 adds to `bookings` are noted here.

**Lead status values (contracted):**

| Status | Description | Terminal |
| --- | --- | --- |
| `new` | Lead created, not yet contacted | No |
| `contacted` | CS has made first contact | No |
| `qualified` | CS has confirmed genuine interest | No |
| `converted` | Lead converted to a booking | Yes |
| `lost` | Lead is no longer an active prospect | Yes |

**Booking table additions (S4 schema changes):**

| Column added | Type | Notes |
| --- | --- | --- |
| `lead_id` | UUID nullable | FK to crm.leads; populated if booking originated from a lead |
| `utm_source` | text nullable | Copied from lead.utm_source at booking creation |
| `utm_medium` | text nullable | Copied from lead.utm_medium at booking creation |
| `utm_campaign` | text nullable | Copied from lead.utm_campaign at booking creation |
| `utm_content` | text nullable | Copied from lead.utm_content at booking creation |
| `utm_term` | text nullable | Copied from lead.utm_term at booking creation |

---

## § ID format (S4 additions)

S4 crm-svc entities use UUID v4 (no prefix). The `booking_id` and `invoice_id` values received from booking-svc / payment-svc retain their ULID-with-prefix format from S1–S3.

| Entity | ID type | Notes |
| --- | --- | --- |
| `leads.id` | UUID v4 | crm-svc owns; no prefix convention in S4 MVP |
| `crm_booking_links.id` | UUID v4 | crm-svc internal |
| `crm_booking_events.id` | UUID v4 | crm-svc internal |

---

## § Audit trail (S4 additions)

Every state-changing call in S4 MUST emit an `iam.audit_logs` row via `iam.v1.IamService/RecordAudit` (per F1 AC). Minimum audit events in S4:

| Trigger | `resource` | `action` | `user_id` |
| --- | --- | --- | --- |
| `POST /v1/leads` creates lead | `lead` | `create` | `""` (anonymous public call) or caller's `user_id` if authenticated |
| `PUT /v1/leads/{id}` updates status | `lead` | `update_status` | Caller's `user_id` (CS or admin) |
| `PUT /v1/leads/{id}` updates assignment | `lead` | `update_assignment` | Caller's `user_id` (admin) |
| `POST /v1/leads/{id}/convert` converts lead | `lead` | `convert` | Caller's `user_id` (CS or admin) |
| `OnBookingCreated` records link | `crm_booking_link` | `create` | `""` (system) |
| `OnBookingPaidInFull` transitions lead to converted | `lead` | `update_status` | `""` (system) |

---

## § Error envelope (S4 additions)

All S4 REST error responses use the shared envelope from S1–S3:

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

- **2026-04-23** — Initial version drafted via tasks `S4-J-01` and `S4-J-02`. Covers: `leads` table schema + UTM snapshot rules (Q019 + Q057 resolution) + CS round-robin algorithm + REST endpoints `POST/GET/PUT /v1/leads` + `POST /v1/leads/{id}/convert` (S4-J-01); `OnBookingCreated` gRPC + `OnBookingPaidInFull` gRPC + `crm_booking_links` + `crm_booking_events` tables + updated `MarkBookingPaid` fan-out chain (S4-J-02).
