---
id: F6
title: Visa Pipeline & Raudhah Shield
status: draft
last_updated: 2026-04-15
moscow_profile: 1 Must Have / 3 Should Have / 1 Could Have — plus module #105 Should Have (cross-cuts to F7)
prd_sections:
  - "E. Operational & Handling — Pelacakan Paspor & Visa (lines 327–337)"
  - "E.5 Field Execution — Raudhah Shield & Tasreh Digital (line 359)"
  - "Tech spec — Modul Document Vault & Visa Tracker (lines 1517–1619)"
  - "Alur Logika 6.1, 6.2, 6.7 (lines 1185–1197) and 8.3 (line 1217)"
modules:
  - "#96 Log Fisik Paspor — Should Have / Medium"
  - "#97 Visa Progress Tracker — Must Have / High"
  - "#98 E-Visa Repository — Should Have / Medium"
  - "#99 Integrasi API Eksternal Lanjutan (GDS + Saudi Visa Provider) — Could Have / Low"
  - "#105 Raudhah Shield & Tasreh Digital — Should Have / Medium (cross-refs F7)"
depends_on: [F1, F3, F4]
open_questions:
  - Q005 — mahram qualifying relations (existing; gates visa submission per Alur 6.7)
  - Q007 — KTP ↔ passport name mismatch (existing; name sent to MOFA/Sajil)
  - Q008 — UU PDP retention (existing; e-visa + passport scan retention)
  - Q026 — MOFA / Sajil API access — credentials, sandbox, contracts
  - Q027 — Provider selection per package kind
  - Q028 — Visa rejection handling policy (retry / escalate / refund)
  - Q029 — Physical passport chain-of-custody process (module #96)
  - Q030 — Raudhah Shield polling cadence + alert action
  - Q031 — Tasreh issuance authority (who creates the document)
---

# F6 — Visa Pipeline & Raudhah Shield

> **Temporal reintroduction:** This feature is the canonical use case that brings **Temporal + broker-svc** back into the stack per ADR 0006. The visa pipeline is a genuine multi-day durable workflow (docs verified → provider submission → polling over days → issuance or rejection). Raudhah Shield is a recurring scheduled workflow that spans the Saudi field window. Implementation of F6 will scaffold `broker-svc` and wire Temporal into `docker-compose.dev.yml` — this is part of F6's work, not a prerequisite.

## Purpose & personas

The pilgrimage cannot happen without a valid Saudi visa. F6 owns the end-to-end visa lifecycle: from "all documents verified" through provider submission (MOFA / Sajil) and multi-day polling to issuance and attachment to the booking. It also runs **Raudhah Shield** — a Nusuk-backed anti-fraud monitor that detects whether a jamaah's issued visa is being used by someone else during the Saudi window.

Primary personas:
- **Ops admin** — triggers bulk submission per departure, monitors the progress tracker, handles rejections, logs physical passport movement.
- **Muthawwif** — uses tasreh digital wallet at Raudhah entry (cross-refs F7 W11).
- **Jamaah** — receives e-visa download link and tasreh via their portal; receives WhatsApp notifications on status changes.
- **Downstream consumers** — `booking-svc` (visa attached → booking readiness), `ops-svc` (manifest generation depends on visa state), `crm-svc` (notifications to agent on status change), `finance-svc` (visa cost becomes a sunk cost for refund calc via Q012).

## Sources

- PRD Section E Pelacakan Paspor & Visa (lines 327–337).
- PRD Section E.5 Raudhah Shield (line 359).
- PRD tech spec — Modul Document Vault & Visa Tracker (lines 1517–1619) — including the authoritative state-machine and bulk-submit endpoint.
- PRD Alur Logika 6.1, 6.2, 6.7, 8.3.
- Modules #96, #97, #98, #99, #105.

## User workflows

### W1 — Document readiness gate (auto-transition)

1. Trigger: F7 W1 verification approves the last required document for a jamaah (KTP + Passport + Photo per PRD L1583).
2. `visa-svc` subscribes to `document.verified` events; when it sees the completing event for a jamaah, queries remaining required docs via `jamaah-svc.GetRequiredDocStatus(jamaah_id)`.
3. If all required docs are `verified`, `visa-svc` transitions the jamaah's `visa_applications.status` from `WAITING_DOCS` → `READY` (PRD L1605, L1583). Records the transition in `visa_status_history`.
4. No Temporal here — direct gRPC + DB write in visa-svc.
5. Emits `visa.ready(jamaah_id, booking_id)` event.

### W2 — Bulk submit to provider

1. Ops admin opens the Visa Progress dashboard, filters by a `package_departure_id`, sees jamaah in `READY` status.
2. Selects a group (typically all jamaah on a departure), picks a `provider_id` (Sajil or MOFA — per Q027), clicks **Submit**.
3. Validation before submit:
   - All selected jamaah in `READY` status
   - Passport expiry rule: `expiry_date − departure_date ≥ 180 days` for every jamaah (PRD L1611 hard rule; HTTP 400 otherwise)
   - Mahram validation (cross-refs Q005): for female jamaah < age threshold, mahram graph must validate — warning at booking is now a hard block at visa submission per Alur 6.7 (L1197)
4. Backend calls `PATCH /v1/visas/bulk-submit` (endpoint 3.3 per PRD L1587, L1591) with `{ paket_id, jamaah_ids[], new_status: 'SUBMITTED', provider_id }`.
5. **Transactional**: all-or-nothing per the PRD's bulk rollback requirement (L1597). visa-svc persists the state transition for every selected jamaah atomically.
6. _(Inferred)_ Submission to the provider happens as an **activity in a Temporal workflow** started by broker-svc: for each jamaah, call the provider API (MOFA or Sajil adapter in visa-svc), attach `provider_ref`, transition to `SUBMITTED`. On any activity failure, the workflow retries per Temporal's policy; compensating action rolls the jamaah back to `READY` if unrecoverable.
7. Emits `visa.submitted(application_id)` event per jamaah.

### W3 — Visa status poll (Temporal-scheduled activity)

1. A long-running Temporal workflow per submitted visa polls provider status until a terminal state.
2. _(Inferred — per Q030 for exact cadence)_ Default cadence: **every 4 hours during business hours** (08:00–22:00 WIB), **every 12 hours off-hours**. Exponential backoff on provider errors; alert to ops if no status change in 14 days.
3. Each poll calls the provider adapter's `GetApplicationStatus(provider_ref)`. Status response mapped to internal state (`SUBMITTED` stays `SUBMITTED`; provider says issued → fetch e-visa PDF; provider says rejected → transition to `REJECTED_BY_EMBASSY`).
4. Writes a `visa_status_history` row per poll (status snapshot + reason text).
5. On **Could Have / Module #99 automation** (direct API): the poll is automatic. Before #99 ships, ops manually marks `ISSUED` / `REJECTED` after checking the provider portal.

### W4 — E-visa received and attached to booking

1. Provider poll returns `ISSUED` status with the e-visa document URL (or fetch via secondary provider API).
2. visa-svc downloads the PDF, uploads to GCS (`e_visas/<jamaah_id>/<application_id>.pdf`), saves path to `e_visas.storage_path`.
3. Updates `visa_applications.status: SUBMITTED → ISSUED`, `issued_date`, `e_visa_url`.
4. _(Inferred)_ Calls `booking-svc.AttachVisa(booking_id, visa_application_id, e_visa_url)` directly via gRPC (per ADR 0006 — in-process coordination; the visa pipeline workflow terminates with this attach as its final activity).
5. Sends WhatsApp notification to jamaah + linked agent: "Your visa is ready." _(Inferred — PRD only explicitly mandates WA on document reject per L1585; extending the pattern to visa issuance for customer satisfaction.)_
6. Emits `visa.issued(application_id, booking_id)` event.

### W5 — Rejection handling

1. Provider poll returns `REJECTED_BY_EMBASSY` with a `reject_reason` (e.g., incomplete docs, suspicious document, mahram issue).
2. visa-svc updates `visa_applications.status: SUBMITTED → REJECTED_BY_EMBASSY`, writes reason to `visa_status_history`.
3. Temporal workflow terminates with a rejection-handling path.
4. **Policy per Q028** — default inferred:
   - **Immediate:** WhatsApp notification to jamaah + linked agent with reject_reason (extends PRD L1585 doc-reject pattern).
   - **Ops action required:** ops reviews the reason, decides on resubmission (new docs / different provider) OR refund flow (triggers F5 W8 refund via `booking-svc.CancelBooking` path with the "visa rejected" reason code; Q012 penalty matrix applies).
   - **Auto-resubmission is NOT done** — too risky; a rejected visa with no human review can result in repeated rejections + provider fee waste.
5. Emits `visa.rejected(application_id, reason)` event.

### W6 — Physical passport movement log (module #96)

1. Ops admin opens the passport's `physical_log` and records the next movement: destination, signer name, timestamp, optional photo of handover.
2. States (per PRD L331, _(Inferred)_ formalised into an enum): `received_from_jamaah | at_pusat | dispatched_to_provider | at_provider | dispatched_to_embassy | at_embassy | returning | returned_to_jamaah`.
3. Each state transition writes a `passport_movements` row.
4. _(Inferred — per Q029)_ SLA per leg: courier hand-off within 48 hours of state change; ops dashboard shows red flag on overdue.
5. This log is independent of the visa-application state — a passport may be physically at the embassy while its visa is still `SUBMITTED`.

### W7 — Raudhah Shield poll cycle (module #105)

Active only during the Saudi field window (from flight-departure day to return-flight day per the booking's `package_departure`).

1. Temporal scheduled workflow per jamaah spins up on departure day.
2. _(Inferred — per Q030)_ Default cadence: **every 6 hours** during the Saudi window. Shorter cadence (every 1 hour) around Raudhah entry appointments if ops knows them.
3. Each poll calls the Nusuk adapter's `GetVisaStatus(passport_number, visa_number)`. Response:
   - **Matches expected jamaah** (name, passport match) → write snapshot, continue.
   - **Mismatch detected** (different holder, unexpected usage pattern) → **alert** per Q030 — default inferred: notify ops pusat + tour leader + muthawwif within 5 minutes via WhatsApp + in-app push; log as an `incidents` row of category `security` (per Q024 escalation matrix).
4. Workflow terminates on return-flight day.

### W8 — Tasreh issuance and storage

1. Tasreh (entry permit, most commonly for Raudhah but also used for other restricted Saudi sites) is uploaded to `tasreh_records`.
2. **Source per Q031** — default inferred: tasreh is **downloaded from Nusuk** by ops admin (most common pattern since Saudi MoHRD issues tasreh electronically via Nusuk now); alternative sources (agency-generated, provider-supplied) supported via a `source` enum.
3. Stored with `kind` (raudhah / other), `issued_for` (scheduled entry time), `storage_path` (GCS), `used_at` nullable.
4. Surfaces in jamaah's Dompet Dokumen Digital (F3 portal) so they can present at entry (PRD L551).
5. Consumed by F7 W11 (tasreh scan at Raudhah entry) which marks `used_at`.

### W9 — Passport expiry gate

Enforced at two points:
1. **At verification** (F3 W3): if `expiry_date − max_future_departure < 180` days, surface warning but allow the document to verify (jamaah might still book a short-term departure).
2. **At visa submission** (W2 above): hard block with HTTP 400 "Paspor tidak valid. Masa berlaku kurang dari 6 bulan" per PRD L1611.
3. _(Inferred)_ No ops override. The 6-month rule is enforced by Saudi; bypassing it would result in embassy rejection or immigration denial.

### W10 — Mahram validation at visa submission

1. Before W2 submit commits, visa-svc calls `jamaah-svc.ValidateMahram(jamaah_id, departure_id)` (F3 W6) for any female jamaah < age threshold (per Q005).
2. Return `is_valid: false` → hard block with reject reason surfaced in the UI (per Alur 6.7, L1197). Unlike the booking-time mahram check (which is a warning), visa submission blocks.
3. Return `is_valid: true` → proceed with submission.

## Acceptance criteria

- `visa_applications` state machine transitions are immutable once committed (history-row per transition; no update-in-place).
- `WAITING_DOCS → READY` transition is automatic when the document-verification gate completes (no manual step).
- Passport 6-month rule is enforced at both verification warning and visa-submit hard block.
- Mahram validation at submit is a hard block; booking-time was a warning.
- Bulk submission is transactional — all selected jamaah commit together or none do (PRD L1597).
- Every status transition writes `iam.audit_logs` via F1 `RecordAudit` and `visa_status_history`.
- E-visa PDFs live in GCS; URLs are 15-min V4 signed (PRD L1613 pattern).
- WhatsApp notification fires within 60 seconds of status change for ISSUED and REJECTED.
- Raudhah Shield alerts reach ops + tour leader + muthawwif within 5 minutes of mismatch detection.
- Temporal workflows survive process restart; poll cadence resumes from last-known state.

## Edge cases & error paths

- **Provider API outage during submit.** Temporal retries with exponential backoff; after N retries, flags for ops manual intervention. Compensating action reverts affected jamaah to `READY` if the batch is abandoned.
- **Provider returns an intermediate state** (e.g., "pending review"). Maps to `SUBMITTED` internally; poll continues.
- **E-visa download fails after `ISSUED` reported.** Workflow state stays at `issued-pending-download`; retries with backoff; ops alerted if persistent.
- **Passport expired mid-pipeline.** Discovered at W2 submit attempt (hard block) or at ops periodic scan (flag to jamaah for renewal).
- **Rejection with ambiguous reason.** Ops interprets; may request free-form details from provider or resubmit after correcting the most likely issue.
- **Raudhah Shield false positive.** Mismatch is flagged but ops+tour leader verify before taking action; incident can be marked `resolved: false_positive` with audit.
- **Nusuk outage during Raudhah Shield poll.** Continue retrying; after N retries in a window, flag that Raudhah Shield coverage is degraded for the affected jamaah. Don't auto-block Raudhah entry (would break legitimate pilgrims).
- **Multi-submission (resubmission after reject).** Create a new `visa_applications` row rather than mutating the rejected one; maintains clean audit trail. Link to prior via `resubmission_of` FK.
- **Tasreh used but jamaah claims they never entered.** F7 W11 scan event is the source of truth; disputes go to ops + incident workflow.
- **6-month rule edge (179 days vs 180 days).** Strict ≥ 180. No half-day fudge.

## Data & state implications

Owned by `visa-svc`. Full schema in `docs/03-services/05-visa-svc/02-data-model.md`. Key additions from this spec:

- `visa_applications` — status enum `waiting_docs | docs_ready | submitted | issued | rejected_by_embassy` (PRD L1605). Add `provider_id`, `provider_ref`, `resubmission_of` nullable FK.
- `visa_status_history` — append-only per transition; reason text.
- `e_visas` — per issued visa; stores GCS path + visa number + validity dates.
- `tasreh_records` — `kind`, `source` enum (`nusuk_download | agency_generated | provider_supplied`), `issued_for`, `storage_path`, `used_at` nullable.
- `raudhah_monitoring` — snapshot per poll: application_id, snapshot jsonb, polled_at.
- **New:** `passport_movements` — state log per module #96: `passport_id` (from F3 documents), `from_state`, `to_state`, `signer_name`, `signer_role`, `photo_url` nullable, `moved_at`.
- **New:** `provider_submissions` — per bulk-submit batch: batch_id, paket_id, provider_id, submitted_by, committed_at, rollback_at nullable.

## API surface (high-level)

Full contracts in `docs/03-services/05-visa-svc/01-api.md`. Key surfaces:

**REST:**
- `GET /v1/visas` — list / filter (by status, departure, jamaah)
- `GET /v1/visas/{id}` — detail + status history
- `PATCH /v1/visas/bulk-submit` — ops bulk submit (PRD L1587)
- `GET /v1/visas/{id}/e-visa` — download via 15-min signed URL
- `GET|POST /v1/passport-movements` — log entries per passport
- `GET|POST /v1/tasreh` — tasreh records
- `GET /v1/raudhah-monitoring?application_id=...` — latest snapshots

**gRPC:**
- `ValidateVisaReadiness(jamaah_id)` — used by booking to gate pre-departure checks
- `GetVisaStatus(application_id)` — used by ops-svc and booking-svc
- `AttachVisaToBooking(booking_id, application_id)` — called on W4 issuance (direct gRPC per ADR 0006)
- `MarkTasrehUsed(tasreh_id)` — called by F7 W11 on entry scan

**Temporal (reintroduced with F6):**
- Workflow `VisaSubmissionWorkflow(application_id)` — long-running: submit → poll → terminal state → attach.
- Workflow `RaudhahShieldWorkflow(application_id)` — scheduled, active during Saudi window.
- Activities in `broker-svc` wrap visa-svc gRPC methods.

## Dependencies

- **F1** — audit, auth, permission checks.
- **F3** — jamaah + passport data (ValidateMahram, GetPassportData, document verification events).
- **F4** — booking records (AttachVisa, CancelBooking on rejection path).
- **F5** — refund flow (on rejection-led cancellation per Q012 matrix).
- **F7** — ops-svc receives tasreh-scan events and Raudhah Shield incident escalation.
- **External:** MOFA / Sajil (per Q026 access), Nusuk (per Q030 access), WhatsApp Business API, GCP Cloud Storage (e-visa + tasreh PDFs).

## Backend notes

- **visa-svc stays Temporal-agnostic** as the baseline pattern. It exposes pure gRPC activities that broker-svc's Temporal workflows call. This preserves the clean separation from ADR 0003 / ADR 0006.
- **Provider adapter pattern** — `visa-svc/adapter/mofa_adapter/` and `visa-svc/adapter/sajil_adapter/` with a common `VisaProviderAdapter` interface. Adding a third provider is a new adapter, not a service-layer change.
- **Bulk submission transactionality** uses `WithTx` on visa-svc's own DB for the state-flip commit; the provider API calls happen inside the Temporal activities (not inside `WithTx`) so network failures don't hold the DB transaction open. If the Temporal workflow can't complete, the application is rolled back to `READY` in a compensation activity — not inside the original `WithTx`.
- **Polling cadence** is a workflow-level concern, configurable per-environment (staging can poll every 15 minutes for testing; production uses Q030 cadence).
- **Nusuk integration** — a separate adapter (`nusuk_adapter`) owned by visa-svc. Raudhah Shield workflow calls it.
- **E-visa / tasreh PDFs** follow the 15-min V4 signed URL pattern from F3 (PRD L1613). External-party access (e.g., airline check-in desk) uses `GenerateExternalAccessURL(document_id, ttl, purpose)` per F3 W9.

## Frontend notes

- **Ops Visa Progress dashboard** — table with filter-by-departure, status columns, bulk-select + bulk-submit action, expandable row showing status history + provider reference + e-visa download.
- **Passport physical log** — per-passport timeline view, "add movement" button with signer field + photo upload.
- **Tasreh management** — upload or download-from-Nusuk flow; preview + share button.
- **Jamaah portal** — visa status badge (PROSES / ISSUED / REJECTED), e-visa download button, tasreh download button.
- **Raudhah Shield alert** surfaces as a high-severity incident on the ops dashboard; tour-leader and muthawwif apps receive push notification for their assigned pax.

## Open questions

See `docs/07-open-questions/`:

**Existing, referenced:**
- **Q005** — mahram qualifying relations (hard block at visa submit per W10)
- **Q007** — KTP ↔ passport name mismatch (name sent to provider)
- **Q008** — UU PDP retention (e-visa + passport scan storage)

**New, filed with this draft:**
- **Q026** — MOFA / Sajil API access (credentials, sandbox, contracts, rate limits)
- **Q027** — Provider selection per package kind (Umroh vs Hajj vs Badal)
- **Q028** — Visa rejection handling policy (retry rules, refund trigger, customer communication)
- **Q029** — Physical passport chain-of-custody process (module #96 states, SLAs, signer requirements)
- **Q030** — Raudhah Shield polling cadence + alert action (who's notified, what's blocked)
- **Q031** — Tasreh issuance authority (Nusuk download vs agency-generated vs provider-supplied)

**Inferred (pending reviewer confirmation):**
- E-visa → booking attach: direct gRPC call (`booking-svc.AttachVisaToBooking`) per ADR 0006
- WhatsApp notification matrix for status changes: fire on ISSUED + REJECTED (extends PRD L1585 doc-reject pattern)
- Polling cadence: 4h business hours / 12h off-hours with exponential backoff on errors (until Q030 confirms)
- 6-month rule: no ops override; absolute per PRD L1611
- Bulk-submit transaction boundary: DB transactional via `WithTx`; provider API calls happen in Temporal activities outside the DB transaction
