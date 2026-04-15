---
id: F7
title: Operations — Verification, Smart Grouping, Manifests, Airport Handling, Field Execution
status: draft
last_updated: 2026-04-15
moscow_profile: 6 Must Have / 5 Should Have / 2 Could Have (+ module #88 Must Have, cross-cutting from F3)
prd_sections:
  - "E. Operational & Handling (lines 295–365)"
  - "E.1 Cancellation — Refund & Pinalti (line 341)"
  - "E.5 Terminal Hub & Field Execution (lines 343–365)"
  - "Alur Logika 6.1–6.5 (lines 1185–1193), 8.1–8.5 (lines 1213–1221)"
  - "Tech spec for Document Vault & Verification (lines 1525–1619)"
modules:
  - "#88 OCR Paspor & Mahram Logic (cross-cutting with F3)"
  - "#90 Generator Surat Resmi, #91 Manifest Imigrasi"
  - "#92 Algoritma Penempatan Kamar, #93 Pengatur Transportasi, #94 Handling Manifest, #95 ID Card & Staff Assignment"
  - "#100 Administrasi Refund & Pinalti"
  - "#101 ALL System, #102 Penghitung Koper, #103 Broadcast Keberangkatan & Kedatangan"
  - "#104 Smart Bus Boarding, #105 Raudhah Shield & Tasreh Digital, #106 Manajemen Perangkat Audio, #107 Distribusi Zamzam, #108 Check-In Kamar Cepat"
depends_on: [F1, F3, F4]
open_questions:
  - Q015 — Smart Grouping trigger timing + override authority (existing)
  - Q012 — Refund penalty policy matrix (existing, referenced by module #100)
  - Q020 — Manifest format per airline / regulator
  - Q021 — Luggage Tag QR payload scheme + ALL system protocol
  - Q022 — Tour leader vs muthawwif override authority boundary (refines Q015)
  - Q023 — Zamzam distribution quota policy
  - Q024 — Incident / issue report escalation matrix
  - Q025 — Vulnerable Care manifest fields (module #94)
---

# F7 — Operations: Verification, Smart Grouping, Manifests, Airport Handling, Field Execution

## Purpose & personas

F7 is the operational backbone — where product becomes pilgrimage. Once a booking is paid (F4/F5) and documents uploaded (F3), F7 is what turns that data into a physically executable trip: verified documents, family-aware room and bus assignments, immigration-ready manifests, QR-coded ID cards and luggage tags, airport handling workflows, and on-the-ground tooling for tour leaders and muthawwif in Saudi Arabia.

This is also the **largest** feature in the catalogue — 15 modules spanning 5 distinct sub-areas. Expect multiple implementation sessions per sub-area.

Primary personas:

- **Ops reviewer (Staff_Ops)** — works the verification queue; approves/rejects OCR output; triggers manifest generation.
- **Ops admin** — runs Smart Grouping; edits allocations pre-departure; generates letters.
- **Tour leader** — in-field orchestrator of a departure (bus boarding, incident reports, hotel check-in reconciliation).
- **Muthawwif** — Saudi-side religious guide; scans tasreh, distributes zamzam, reports Raudhah issues.
- **Airport handler / porter** — scans ID cards and luggage at departure and arrival.
- **Tour leader and muthawwif use separate mobile UIs** but their data lives in the same backend.

## Sources

- PRD Section E in full (lines 295–365) — Document Vault (overlap with F3), Automated Letter Engine, Smart Grouping, Cancellation, Terminal Hub, Field Execution.
- PRD Alur Logika 6.1–6.5 (verification + letter + grouping + refund flows) and 8.1–8.5 (airport + field flows).
- PRD tech spec for Jamaah_Documents verification endpoints (lines 1525–1619).
- 15 modules enumerated in frontmatter.

## User workflows

### W1 — Document verification queue (module #88)

1. `jamaah-svc` OCR pipeline (F3 W4) produces a document in `needs_review` with extracted MRZ fields.
2. Ops reviewer opens the verification queue (filter by status = `needs_review`, order by `created_at`).
3. Selects a doc; previews via 15-min V4 signed URL (PRD L1613); sees OCR extraction alongside the image.
4. Edits `final_extracted_data` if the OCR needs correction.
5. Clicks **Approve** → `PATCH /v1/documents/{id}/verify` with `{ status: 'VERIFIED', final_extracted_data }`. Backend syncs corrected fields into master `jamaah` record (PRD L1569). If this completes the jamaah's required-doc set (KTP + Passport + Foto for visa purposes), `visa_applications.status` auto-flips `WAITING_DOCS → READY` (PRD L1577).
6. Or clicks **Reject** → `PATCH .../verify` with `{ status: 'REJECTED', reject_reason: 'Foto buram' }`. WhatsApp notification fires to jamaah + linked agent requesting re-upload (PRD L1585).
7. Audit log written async (old_value, new_value, actor, IP, timestamp) per PRD L1619.

Cross-cutting edge rules from F3:
- Passport expiry < 6 months before departure → HTTP 400 (PRD L1611) — hard block, not verification-level.
- Woman < 45 (per Q005 threshold) with `mahram_id` null → Warning on verify (PRD L1617), not block.

### W2 — Smart room allocation (module #92)

Trigger timing: **Q015** (open). Default inference per Q015's recommendation is "on-demand by ops admin when pax list is stable, typically ≥ H-14" — not at booking submit.

1. Ops admin opens a `package_departure` and clicks **Run Grouping**.
2. Algorithm inputs:
   - All `bookings` in `paid_in_full` or `partially_paid` status on this departure
   - `jamaah.family_unit_id` (K-Family Code, PRD L319) and `mahram_relations` (from F3) for family clustering
   - `jamaah.domicile` (for city-of-origin grouping)
   - `preferences.ac_tolerance` and similar preference flags (PRD L319)
   - Linked `hotel` room inventory from catalog (room types + counts)
3. _(Inferred — see Backend notes for algorithm)_ Produces room assignments respecting: (a) mahram constraint (male mahram shares with minor female or spouse; not with adult non-mahram female), (b) family unit co-location, (c) domicile proximity as a tie-breaker, (d) room-type capacity limits.
4. Writes `room_allocations { booking_id, hotel_id, room_number, occupants[] }` to ops-svc.
5. Returns a diff preview before commit so ops can review before publishing.

**Re-running after partial cancellation:** leaves a hole rather than auto-regrouping (Q014 default). Ops decides whether to fill manually or re-run full grouping.

### W3 — Smart bus / transport allocation (module #93)

Parallel to W2, triggered in the same ops action:

1. Algorithm reads linked `airline` + `bus` transport inventory, pax count per departure.
2. Assigns:
   - **Flight seats** — family clusters seated adjacent where the airline class allows; wheelchair and elderly flagged for special assistance rows.
   - **Bus seat manifests** — "Bus 1", "Bus 2", ... lists per leg of the journey (airport → hotel, inter-city, hotel → airport).
3. Writes `bus_allocations { package_departure_id, bus_number, seat_assignments }`.

### W4 — Vulnerable Care manifest (module #94)

Special manifest for airport handling + muthawwif flagging lansia (elderly), wheelchair-dependent, medical-needs jamaah.

1. During W2/W3 grouping, jamaah flagged with `vulnerable_category` (field set per **Q025**) are extracted into a separate manifest.
2. Handling team sees this manifest at airport for boarding priority; muthawwif sees it in-field for special-care routing.

### W5 — Manifest & letter generation (modules #90, #91)

1. Ops admin clicks **Generate Letter** or **Generate Manifest** for a package_departure.
2. For letters (module #90): picks type (Rekomendasi Paspor / Surat Izin Cuti / Surat Keterangan Berangkat), selects recipient jamaah(s), renders PDF on agency letterhead.
3. For manifests (module #91): system renders **immigration-format** manifest — **Q020** pins the exact format required. Default inferred: PDF with jamaah fields (name, passport number, DOB, gender, nationality) in a tabular layout with agency header; Excel export supported in parallel. Regenerated as pax list changes until H-24 lock.
4. Signed-URL download; audit-logged.

### W6 — ID Card + Luggage Tag issuance with QR (module #95)

1. Output of W2/W3 grouping. Each jamaah gets:
   - **ID Card** — printable card with photo, name, booking code, package, departure date, assigned muthawwif, QR payload.
   - **Luggage Tag** — printable tag, booking code, luggage sequence number (1/3, 2/3, 3/3), QR payload.
2. **QR payload** — per **Q021**. Default inferred: a signed short token (HMAC) encoding `{ jamaah_id, booking_id, tag_seq }` with a server-side secret, verifiable but not reversible without the secret. Prevents spoofing.
3. Regeneration: when a room or bus assignment changes post-issue, ID Card is reprinted with the new data; old QR token is marked invalid in audit. Luggage tags survive allocation changes (tag_seq is stable).

### W7 — ALL System scan (module #101, PRD L347)

ALL = Airport Logistics List. Agency-internal, not a third-party vendor.

1. Airport handling team uses a tablet/phone app running the ALL System UI.
2. At check-in desk, scans jamaah ID card QR. System resolves the jamaah and displays a **four-checkbox card**: presensi, paspor handover, visa handover, boarding pass handover.
3. Handler ticks each as completed; backend writes `handling_events` rows (one per event type).
4. When all four checks complete for all pax on a departure, the departure `ground_ready` flag flips true — visible to ops dashboard.
5. _(Inferred)_ **Offline fallback:** when WiFi drops, the tablet queues scans locally in IndexedDB / SQLite; syncs on reconnect with idempotency on `(device_id, scanned_at, event_kind, jamaah_id)`. Conflict resolution if two tablets scan the same jamaah: first-write wins; second returns 409 with the winning event's timestamp.

### W8 — Penghitung Koper / Luggage counter (module #102)

1. At departure airport, porter scans each luggage tag QR before loading onto the aircraft. System writes `handling_events { kind: 'luggage_departure_scan', ... }`.
2. At Saudi arrival (Jeddah / Medina), Saudi porter team scans each tag on pickup. `handling_events { kind: 'luggage_arrival_scan', ... }`.
3. Reconciliation view: for a departure, show per-jamaah luggage_departure_count vs luggage_arrival_count.
4. _(Inferred)_ **Mismatch policy:** soft warning surfaced to ops + muthawwif. Not a hard block — airlines routinely route luggage late. Ops follows up with the airline; if not resolved within 48h of arrival, incident workflow (W11) triggers.

### W9 — Broadcast Keberangkatan & Kedatangan (module #103)

1. System-triggered WhatsApp broadcasts via the same adapter as F10:
   - **H-3 jam (3 hours before departure)** — airport meeting-point location, gate info, tour-leader WA contact. Sent to jamaah + linked emergency contact.
   - **On arrival** — "welcome" + pickup gate number + hotel transfer bus number.
2. _(Inferred)_ Templates are company-wide Bahasa-Indonesia defaults editable by Super Admin via Module #161 Manajemen Template Komunikasi. No per-package override in MVP.

### W10 — Smart Bus Boarding (module #104, PRD L357)

1. Tour leader opens mobile app at bus door.
2. Scans each jamaah's ID card QR as they board. System writes `handling_events { kind: 'bus_boarding', bus_number, ... }`.
3. Roster view shows scanned / not-yet-scanned; identifies missing jamaah before the bus departs.
4. When all pax on the bus's allocation are scanned, tour leader taps **Start Trip**.
5. Bus Time Manager monitors: if the bus remains stationary > 20 minutes after Start Trip, alerts pusat ops (PRD L563).
6. _(Inferred)_ Manual override for failed scans: tour leader can mark a jamaah as boarded with a required reason ("scanner failed"; "jamaah missing ID card"); logged to audit. **Q022** pins whether muthawwif has the same override authority.

### W11 — Tasreh scan / Raudhah Shield (module #105, cross-refs F6)

1. Jamaah presents digital Tasreh (from their F3 document wallet) at Raudhah entry.
2. Muthawwif or Askar (Saudi enforcement) scans the Tasreh QR via the muthawwif app.
3. Backend verifies against F6 Raudhah Shield (which polls Nusuk for the visa holder's real-time status).
4. If Tasreh is valid and visa is active → proceed. If visa status anomalous → **Q? (cross-ref F6)** — spec in F6 covers the action (flag for review; notify legal; etc. — not owned by F7).
5. F7's role: scan event ingestion, log to `handling_events`, surface on muthawwif dashboard.

### W12 — Zamzam distribution (module #107)

1. At distribution point (Saudi hotel or shipping depot), muthawwif scans jamaah ID card.
2. System writes `handling_events { kind: 'zamzam_distribution', jamaah_id, quantity_liters, scanned_at }`.
3. Per-jamaah quota check (**Q023**): default inferred 5 liters per jamaah per trip, one-time issuance. Attempt to scan the same jamaah again returns a soft warning with a "manager override" path for legitimate cases (e.g., spilled container replacement).

### W13 — Incident / issue reporting (field execution)

1. Tour leader or muthawwif taps **Issue Report** button in their app.
2. Picks category: medical / lost jamaah / vendor problem / logistics / security / other.
3. Adds short description, photo (optional), location pin.
4. System writes `incidents { severity, category, reporter_id, location, description, created_at }`.
5. **Escalation (Q024):** who gets notified, SLA, resolution workflow — stakeholder-bound.

### W14 — Check-In Kamar Cepat (module #108, Could Have)

_(Inferred — this is Could Have / Low priority per module list.)_

1. Muthawwif opens app in hotel lobby on arrival.
2. Views side-by-side: system Room List (from W2) vs actual hotel room-number assignment (manually entered by muthawwif from the hotel's rooming list; photo upload of the printed sheet is supported for audit).
3. Reconciles mismatches: e.g., "Jamaah X was planned for Room 302 but hotel assigned Room 305" — muthawwif accepts or reassigns within the room-type constraint.
4. Writes updates to `room_allocations` with `reconciled_at` flag.
5. _(Inferred)_ A CSV-import or API integration from Saudi hotels is out of scope for MVP — manual entry suffices.

### W15 — Refund / Pinalti administration (module #100, cross-refs F4/F5)

1. On F4 cancellation (W7 in F4), F5 refund flow (W8 in F5) pulls the penalty matrix from **Q012** (pending answer).
2. F7's contribution: **ops-side paperwork**. When a refund is initiated, ops-svc generates a refund agreement PDF listing the components (ticket cost burned, visa cost burned, hotel cost burned, pinalti) and emits it to jamaah via WhatsApp + email for acknowledgement.
3. Finance (F9) handles the money; ops handles the document trail.

## Acceptance criteria

- Every document approval / rejection writes to `iam.audit_logs` via F1 `RecordAudit`.
- Verification queue latency: ops reviewer sees a new `needs_review` doc within 10 seconds of F3 OCR completion.
- Smart Grouping (W2 + W3) completes in < 30 seconds for a 100-jamaah departure.
- Mahram constraint is provably respected: no room assignment places an adult female with a non-mahram adult male. Unit-tested exhaustively.
- Room allocation's vulnerable-care manifest (W4) extracts every jamaah with a `vulnerable_category` flag set; no false negatives.
- Manifest generation produces PDF + Excel (per Q020 default) for any `package_departure` in `open`/`closed` status.
- QR-encoded ID cards and luggage tags use signed tokens (Q021 default); unsigned or tampered QR returns 401.
- ALL System scan events are idempotent on `(device_id, scanned_at, event_kind, jamaah_id)`.
- Scan events queue locally during offline (W7), survive app crashes, and sync on reconnect.
- Bus boarding (W10) alerts pusat on stationary-bus-after-start condition within 20 min.
- Zamzam per-jamaah quota (Q023 default 5L) enforced; override requires reason + audit.
- Incident report reaches pusat within 30 seconds of submission (Q024 default routing).

## Edge cases & error paths

- **OCR low confidence** → document stays in `needs_review`; staff manually fills fields (PRD L1615, HTTP 206 fallback).
- **Passport expiry violation** at verification → HTTP 400; block `VERIFIED` status (PRD L1611).
- **Mahram-warning on verification** → soft warning only (PRD L1617). Block is at booking or visa stage, not here.
- **Mid-grouping reshuffle after a partial cancel** (Q014) → leave the hole; ops decides. Don't auto-regroup a live departure.
- **Luggage tag reprint after room change** → Q021 scheme uses stable `tag_seq` so the physical tag can persist; old QR token revoked; new ID card regenerated.
- **ALL system offline during check-in** → local IndexedDB queue; sync on reconnect with idempotency keys (W7 inferred).
- **Two tablets scan same jamaah concurrently** → first-write wins; second returns 409.
- **Luggage count mismatch** → soft warning, not block; incident workflow after 48h (W8 inferred).
- **Bus stationary > 20 min** → pusat alert (module #104).
- **Tasreh visa flagged anomalous in Nusuk** → handled in F6; F7 just logs the scan event.
- **Zamzam over-quota attempt** → soft warning + override path (Q023 default).
- **Incident with no connectivity** → queue locally in muthawwif app; sync on reconnect; escalation SLA starts from sync-time, not report-time, with that fact surfaced in the audit (inferred).
- **Hotel room-number mismatch at Check-In Kamar Cepat** → muthawwif adjusts; audit trail preserves the original grouping-algorithm output + the reconciled assignment (W14 inferred).

## Data & state implications

Owned by `ops-svc`. Full schema in `docs/03-services/06-ops-svc/02-data-model.md`. Key tables referenced / extended by this spec:

- `verification_tasks` — queue entries referencing F3 `documents`.
- `room_allocations { booking_id, hotel_id, room_number, occupants[], reconciled_at }`.
- `bus_allocations { package_departure_id, bus_number, seat_assignments (jsonb) }`.
- `manifests { package_departure_id, format ('pdf' | 'xlsx' | 'vulnerable_care'), storage_path, jamaah_count, generated_at }`.
- `luggage_tags { id, booking_id, jamaah_id, tag_code, qr_payload_signed, tag_seq_total, tag_seq }`.
- `handling_events { device_id, event_kind enum, jamaah_id, booking_id, metadata jsonb, scanned_at, synced_at nullable }` — idempotent on `(device_id, scanned_at, event_kind, jamaah_id)`.
- `incidents { severity, category, reporter_id, location, description, status, escalated_to[], created_at, resolved_at }`.
- `grouping_runs { package_departure_id, input_snapshot jsonb, output jsonb, ran_at, committed_at, reconciled_at }`.

New enum additions:
- `vulnerable_category`: `lansia | wheelchair | medical | pregnant | dietary` (Q025 may refine).
- `handling_event_kind`: `luggage_departure_scan | luggage_arrival_scan | bus_boarding | zamzam_distribution | tasreh_scan | all_system_checkin`.
- `incident_category`: `medical | lost_jamaah | vendor | logistics | security | other`.

## API surface (high-level)

Full contracts live in `docs/03-services/06-ops-svc/01-api.md` — spec already planned. Key surfaces this draft confirms:

**REST (ops console + field apps):**
- `GET /v1/verification-tasks` + `POST /v1/verification-tasks/{id}/approve|reject`
- `POST /v1/grouping/run` (W2 + W3 combined)
- `POST /v1/manifests/generate` + `GET /v1/manifests/{id}/download`
- `POST /v1/luggage-tags/issue` + `GET /v1/luggage-tags?departure_id=`
- `POST /v1/handling-events` (batch endpoint for offline-sync flush from field apps)
- `POST /v1/incidents` (field app report submission)
- `POST /v1/room-allocations/{id}/reconcile` (Check-In Kamar Cepat)

**gRPC (service-to-service):**
- `RunSmartGrouping(...)` — callable by ops-svc's own scheduler or by an admin via REST
- `BuildManifest(...)` — same
- `GetVerificationStatus(jamaah_id)` — read by F6 visa to gate `WAITING_DOCS → READY`
- `NotifyDepartureGroundReady(departure_id)` — emitted when W7 all-four-checks complete per jamaah for an entire departure

## Dependencies

- **F1** — auth, RBAC, audit
- **F3** — jamaah + documents (OCR output feeds W1; mahram graph feeds W2)
- **F4** — bookings + paid status (input to grouping; refund trigger)
- **F5** — refund flow (cross-refs W15); WhatsApp adapter (W9, W1 reject notifications)
- **F6 (future)** — visa readiness gate after W1; Raudhah Shield consumes W11 scans
- **External** — GCP Vision (via F3), GCP Cloud Storage (signed URLs), WhatsApp Business API (broadcasts, reject notifications), Nusuk app (via F6)

## Backend notes

- **Smart Grouping algorithm (W2) — sketch:**

  ```
  Input: jamaah[] (with family_unit, mahram_relations, domicile, gender, age, preferences, vulnerable flags),
         hotel_rooms[] (per type: Double/Triple/Quad counts)

  Step 1: Sort jamaah by (family_unit_id, domicile) — families cluster.
  Step 2: For each family cluster:
    a. Group adults by mahram compatibility (a room may contain mahram-related individuals of both genders;
       otherwise single-gender only).
    b. Fit into available room types preferring smaller rooms for smaller clusters (Double for 2, Triple for 3, Quad for 4).
    c. If family exceeds largest room (>4), split across adjacent rooms respecting mahram rules.
  Step 3: Fill remaining rooms from non-family singletons:
    a. Group by gender + domicile + preference compatibility.
    b. Single-gender rooms only (no cross-gender without mahram).
  Step 4: Unassigned jamaah → flag for ops manual placement.
  Output: room_allocations[]
  ```

  _(Inferred algorithm. Reviewer and religious advisor may need to validate mahram compatibility rules; overlaps with Q005.)_

- **Manifest generator** runs in a separate worker pool (same as F2 flyer render, same as F5 receipt PDF). Heavy render should not block verification-queue REST traffic.
- **Scan-event idempotency** via DB unique constraint on `(device_id, scanned_at, event_kind, jamaah_id)`. Conflict → 409 with pointer to the winning event.
- **Offline-sync batch endpoint** accepts an array of scan events; per-event idempotency handled; response indicates per-event success/409 status.
- **QR signing key** lives in config (same secret-management pattern as auth tokens per F1). Rotation strategy: overlap window where old + new keys both verify, controlled by config + a force-rotation admin action.
- **Vulnerable Care flag** set at F3 jamaah registration (per Q025 fields) and carried through; ops can override per-booking.

## Frontend notes

- **Ops verification console** — verification queue is the single highest-touch ops UI. Keyboard shortcuts for approve/reject, image-zoom, OCR-field inline edit.
- **Ops admin console** — Smart Grouping runner with diff preview, manifest generator, letter templates, incident feed.
- **Tour leader mobile app (Svelte)** — scanner surface (bus boarding, ID scan), incident report form, group roster, check-in kamar cepat screen.
- **Muthawwif mobile app (Svelte)** — tasreh scanner, zamzam distribution scanner, room reconciliation, incident escalation.
- **Airport handling tablet UI (Svelte)** — ALL System four-checkbox scanner, luggage counter, manifest view.
- **Field apps need offline-first** — IndexedDB / SQLite queue, sync service, conflict resolution UI for rare 409 cases.

## Open questions

See `docs/07-open-questions/`:

**Existing, referenced:**
- **Q012** — refund penalty policy matrix (module #100 uses this)
- **Q015** — Smart Grouping trigger timing + ops/tour-leader override

**New, filed with this draft:**
- **Q020** — Manifest format per airline / immigration regulator
- **Q021** — Luggage Tag QR payload scheme + ALL System scan protocol
- **Q022** — Tour leader vs muthawwif override authority boundary (refines Q015 for in-field actions)
- **Q023** — Zamzam distribution quota policy
- **Q024** — Incident / issue report escalation matrix
- **Q025** — Vulnerable Care manifest fields (module #94)

**Inferred (pending reviewer confirmation):**
- ALL System offline fallback: local queue + sync on reconnect, idempotent
- Luggage count mismatch: soft warning, not block; incident workflow after 48h
- Check-In Kamar Cepat (module #108): manual muthawwif entry for MVP (module is Could Have)
- Manajemen Perangkat Audio (module #106, Could Have): **out of MVP scope**; track audio receivers in an external Excel workaround until a concrete demand surfaces
- WhatsApp broadcast templates (module #103): company-wide Bahasa defaults via Module #161, editable by Super Admin, no per-package override in MVP
