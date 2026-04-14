---
id: F3
title: Pilgrim Profile & Documents
status: draft
last_updated: 2026-04-14
moscow_profile: 2 Must Have / 2 Should Have (mahram algorithm is the domain headline)
prd_sections:
  - "E. Operational & Handling — Document Vault (L1517–1619)"
  - "A. B2C Customer Portal — Self-Upload Dokumen (L87, L701, L939–943)"
  - "Fase 6.1 & 6.7 — mahram verification for visa submission (L1185, L1197)"
modules:
  - "#20 (Self-Upload Dokumen) — Should / Medium"
  - "#87 (Penyimpanan Kolektif) — Must / High"
  - "#88 (OCR Paspor & Mahram Logic) — Must / High"
  - "#89 (Progress Tracker & Expiry Alert) — Should / Medium"
depends_on: [F1]
open_questions:
  - Q005 — mahram qualifying relation set + age threshold (religious/legal)
  - Q006 — minimum documents required to submit a booking
  - Q007 — KTP ↔ passport name mismatch handling
  - Q008 — UU PDP compliance (consent, retention, DSR)
---

# F3 — Pilgrim Profile & Documents

## Purpose & personas

The authoritative record of **who** each pilgrim is. Stores jamaah biodata, the family graph that powers mahram validation, and the scanned documents (KTP, passport, vaccine card, family book, photo) that the visa pipeline and ops manifest generation depend on. Runs the OCR pipeline that turns a passport photo into structured data. Never trusted blindly — every document lands in a staff verification queue before being accepted as truth.

Primary personas:
- **Calon jamaah** — uploads their own documents via the B2C portal (replaces WhatsApp photo submission).
- **Agent** — uploads on behalf of their registered jamaah via the B2B portal.
- **CS** — registers a jamaah during a phone booking; uploads scanned docs from the physical paperwork the caller mails or brings in.
- **Ops admin / Staff_Ops** — verifies uploaded documents, approves or rejects, triggers mahram re-check, monitors the readiness dashboard.
- **Consumers (other services)** — `visa-svc` (reads passport data for MOFA/Sajil), `booking-svc` (mahram validation at submit), `ops-svc` (manifest & verification queue).

## Sources

- PRD "Modul Document Vault & Visa Tracker" at L1517–1619 — the authoritative tech spec
- PRD L87, L301 (vault contents), L319 (K-Family Code), L701, L939–943 (self-upload flow)
- PRD L1185, L1197 (mahram references in visa submission — Fase 6.1 and 6.7)
- Modules #20, #87, #88, #89 from `docs/Modul UmrohOS - MosCoW.csv`
- Existing service data model: `docs/03-services/03-jamaah-svc/02-data-model.md`

## User workflows

### W1 — Jamaah self-upload from B2C portal

1. Calon jamaah logs in to `/jamaah/dokumen`. Authentication via F1 (jamaah role, personal scope).
2. Picks a document kind: `ktp | passport | vaksin | family_book | photo`.
3. Uploads an image (JPEG/PNG) or PDF. _(Inferred)_ MIME allowlist `image/jpeg, image/png, application/pdf`; size cap **10 MB**; reject anything else with 415 before touching storage.
4. Backend streams the file to GCS bucket `umrohos-private-docs` under path `jamaah/{jamaah_id}/{kind}/{document_id}` _(Inferred path convention)_. Inserts a `documents` row with `status: uploaded`.
5. If kind is `passport`, enqueues the OCR pipeline (W4). For other kinds, the row moves directly to the ops verification queue (W3).
6. Jamaah sees the upload listed in their portal with a real-time status badge (`uploaded → processing → needs_review → verified | rejected`).

### W2 — CS / agent upload on behalf

1. Ops admin or agent opens the jamaah's profile in the internal console.
2. Uses the same endpoint as W1 but with `jamaah_id` in the multipart body and `uploaded_by = <staff user id>` recorded.
3. Remaining flow identical to W1.

### W3 — Ops verification queue

1. Staff_Ops opens the verification queue. The queue lists all `documents` rows where `status in (needs_review, processing)`, ordered by `created_at`.
2. Previews the file via a **15-minute V4 signed URL** (module #87 requirement; strict per PRD L1613). The URL is generated per-request and not cached.
3. Reviews OCR output (if passport) side-by-side with the image. Can edit fields before approving.
4. Clicks **Approve** → `status: verified`, final OCR fields sync into the `jamaah` master biodata. **Reject** → `status: rejected`, staff writes a reject reason. Every action writes to `iam.audit_logs` via F1's `RecordAudit`.
5. On reject, a WhatsApp notification fires to the jamaah (and CC to their linked agent if any) containing the reject reason and a link back to the upload screen (PRD L1585).
6. When the full required-doc set for a jamaah is `verified`, the associated `visa_applications` row flips `WAITING_DOCS → READY` via gRPC into F6.

### W4 — Passport OCR pipeline (module #88)

1. Trigger: a `documents` row with `kind = passport` and `status = uploaded` is enqueued.
2. Worker downloads the image from GCS, sends to **GCP Cloud Vision — Text Detection**.
3. Regex parses the MRZ to extract six fields: `full_name`, `passport_number`, `date_of_birth`, `nationality`, `gender`, `expiry_date` (PRD L1563).
4. Writes `ocr_results { document_id, extracted, confidence }`. Document status becomes `needs_review` (PRD L1603).
5. **Low-confidence or MRZ regex failure** → worker returns 206 "Silakan input manual" to the trigger; document still moves to `needs_review` so staff can fill fields manually (PRD L1615).
6. **Passport 6-month rule** applies when this passport is later linked to a booking departure: if `expiry_date − departure_date < 180 days`, reject with 400 "Paspor tidak valid. Masa berlaku kurang dari 6 bulan" (PRD L1611). _(Inferred — this check lives on the booking-submit path, not in the OCR worker, because OCR doesn't yet know the departure.)_
7. No auto-approval, ever. Human-in-the-loop is mandatory before `verified`.

### W5 — Family unit formation & mahram edges

1. CS or ops admin creates a `family_unit` with a human-friendly K-Family Code (e.g. `KF-2026-0342`).
2. Attaches one or more jamaah records to the unit via `family_unit_id`.
3. For each kinship pair, inserts a row in `mahram_relations` with `(subject_jamaah_id, related_jamaah_id, relation, verified: bool)`. The `relation` enum covers: husband, wife, father, mother, son, daughter, brother, sister, grandfather, grandson, uncle, nephew, father_in_law, son_in_law.
4. _(Inferred)_ Relations are stored as **directed edges**. Inverses are auto-inserted by a DB trigger (father ↔ son/daughter, husband ↔ wife, etc.) so a recursive walk can traverse in either direction.
5. Proof-of-relation document requirements: **Buku Nikah** proves husband/wife edges; **Kartu Keluarga** proves parent/child/sibling edges; other relations verified manually by staff with `verified = false` until overridden. _(Inferred — PRD is silent on proof docs per relation; see Q005.)_

### W6 — Mahram validation at booking / visa submit

1. **Caller:** `booking-svc` at submission time (F4), and `visa-svc` at visa submission (F6), via gRPC into `jamaah-svc.ValidateMahram`.
2. **Input:** `{ subject_jamaah_id, group_jamaah_ids[], departure_id }`.
3. **Output:**
   ```
   {
     is_valid: bool,
     needs_mahram: bool,
     found_mahram: { jamaah_id, relation } | null,
     reason_code: OK | NOT_REQUIRED | NO_RELATION_IN_GROUP
                | RELATION_NOT_QUALIFYING | GENDER_MISMATCH
                | NOT_IN_SAME_DEPARTURE | AGE_UNKNOWN,
     reason_message: string
   }
   ```
4. **Logic — TBD on the details, see Q005.** Sketch of the default behaviour we'd implement pending the religious/legal answer:
   - If subject is male OR age ≥ threshold → `is_valid: true, needs_mahram: false, reason_code: NOT_REQUIRED`.
   - Otherwise walk `mahram_relations` to find a path to any male member of `group_jamaah_ids` whose `relation` is in the **qualifying set**. If found and the candidate is in the same `departure_id`, return `is_valid: true, found_mahram: {...}`.
   - Else return `is_valid: false` with the specific reason.
5. The current PRD rule (L1617) is only the binary check: "female, age < 45, `mahram_id` must not be null". Our implementation models the richer graph intentionally; see Q005 for the qualifying relation set + age threshold + same-departure rule that we need stakeholder sign-off on.

### W7 — Passport expiry alerting (module #89)

1. Daily scheduled job scans `documents` where `kind = passport`, `status = verified`, and the owning jamaah has a `booking` on a future departure.
2. For each, compute `passport.expiry_date − departure_date`.
3. If < **180 days**: flag on the jamaah's dashboard card (yellow / red progress badge per module #89), send WhatsApp + email alert to jamaah and linked agent. _(Inferred — PRD implies WA-driven; email as fallback.)_
4. If the passport expires **before** departure: alert escalated; booking auto-flagged for ops review.

### W8 — Document re-upload after rejection

1. WhatsApp reject notification links jamaah to the upload screen (W1).
2. New upload creates a **new `documents` row** for the same kind. The old rejected row stays in place for audit (PRD-aligned: never delete). _(Inferred — soft-supersede.)_
3. The new row runs through OCR (if passport) and verification (W3–W4) normally.
4. The jamaah's document completeness is computed as "at least one `verified` row per required kind" — old rejected rows don't block completeness.

### W9 — External party access (visa provider, embassy)

1. When `visa-svc` submits to MOFA/Sajil, it needs to hand the provider the passport scan.
2. _(Inferred)_ `jamaah-svc` exposes `GenerateExternalAccessURL(document_id, ttl, purpose)` — returns a signed URL with extended TTL (e.g. 1 hour) and logs the access intent with `purpose` (`visa_submission`, `embassy_review`, etc.) to audit.
3. No direct GCS access is ever handed to external parties.

## Acceptance criteria

- Documents flow: `uploaded → processing (passport only) → needs_review → verified | rejected`, with staff action required for every `verified` transition.
- GCS signed URLs are V4 with **≤ 15 minute TTL** by default; external-party URLs go through `GenerateExternalAccessURL` with explicit purpose logging.
- Passport OCR extracts all six MRZ fields (name, passport_number, DOB, nationality, gender, expiry) for at least 95% of clean input images. _(Inferred benchmark; PRD silent.)_
- MRZ regex failure returns HTTP 206 with `NEEDS_REVIEW` status — never auto-reject.
- Passport 6-month rule blocks booking submission with HTTP 400 at the F4 submit path.
- `ValidateMahram` gRPC returns in < 50 ms p95 for groups up to 20 jamaah. Heavy unit-test coverage on the graph walk.
- Document rejection triggers a WhatsApp notification within 30 s of the staff action.
- Passport expiry alerts run on a daily cron; the first alert fires ≥ 180 days before departure.
- Every state-changing action on `documents`, `jamaah`, `family_units`, `mahram_relations` writes to `iam.audit_logs` via F1's `RecordAudit`.
- Upload rejects MIME types outside the allowlist with 415 before a byte of file data is stored.
- Soft-supersede: new uploads for the same kind do not delete prior rejected rows; the rejected row remains in place for audit.

## Edge cases & error paths

- **MRZ regex fails entirely.** Status moves to `needs_review` with empty `ocr_results.extracted`; staff enters fields manually in W3.
- **MRZ extracts but passport_number collides with an existing verified passport on a different jamaah.** Flag for ops review; do NOT auto-merge. Indonesian name-reuse across family members makes this a real risk.
- **KTP name does NOT match passport name.** Very common in Indonesia (shortened names, Arabic transliteration differences). Current default: store both, flag the mismatch on the jamaah profile, require staff acknowledgement before `verified`. Final policy depends on **Q007**.
- **Jamaah uploads passport belonging to someone else.** Detected either by the collision rule above or by a mismatch between OCR name and the jamaah's registered name. Staff rejects with reason `wrong_passport`.
- **Mahram verification succeeds at booking but fails later** (e.g. the mahram cancels). The booking-svc status-change path (F4 cancellation) signals `jamaah-svc` to re-check; if now invalid, escalate to ops — do NOT auto-cancel the female's booking. _(Inferred — PRD silent on this reversal path.)_
- **Passport expires between verification and departure.** The daily cron (W7) catches it; booking is flagged red on the dashboard.
- **Two jamaah share the same KK but are flagged as husband-wife plus parent-child cycle** (reconstituted families exist). The graph permits this; mahram validation uses the first qualifying path found. No cycle detection beyond standard CTE depth limit.
- **File upload fails mid-stream.** Partial file in GCS is cleaned up by the handler on error. Document row is not inserted unless the multipart completes successfully.
- **OCR provider outage.** Document sits in `processing` with a retry counter; after 3 retries, moves to `needs_review` with a flag that OCR was skipped, so staff can enter fields manually.
- **Jamaah attempts to book before any document is uploaded.** Policy depends on **Q006** (minimum docs required). Default (pending Q006): allow draft booking with KTP only; block submit-to-saga without passport OCR verified.

## Data & state implications

Owned tables in `jamaah-svc`. Full schema in `docs/03-services/03-jamaah-svc/02-data-model.md`. Key points aligned with this spec:

- `jamaah`: biodata, `family_unit_id` nullable, `status` enum (`calon | active | alumni`). Branch-scoped.
- `family_units`: `code` (K-Family Code), human name.
- `mahram_relations`: directed edges with `relation` enum and `verified` bool. Recursive CTE supports graph traversal. _(Implementation detail, not in PRD.)_
- `documents`: `kind`, `storage_path`, `status`, `uploaded_by`. Old rejected rows persist (soft-supersede).
- `ocr_results`: JSONB `extracted`, numeric `confidence`. One row per document.
- `document_kind` enum extended from the baseline to: `ktp | passport | vaccine | family_book | photo | marriage_book` (Buku Nikah added for mahram proof).
- `document_status` enum: `uploaded | processing | needs_review | verified | rejected`. (PRD L1603 uses `NEEDS_REVIEW`; we mirror.)

## API surface (high-level)

Full contracts land in `docs/03-services/03-jamaah-svc/01-api.md`. Key surfaces:

**REST:**
- `POST /v1/jamaah` — register new jamaah (CS or self-register via portal)
- `GET|PATCH /v1/jamaah/{id}` — profile read/update
- `GET /v1/jamaah/{id}/documents` — list with current status per kind
- `POST /v1/jamaah/{id}/documents` — upload (multipart; worker triggers OCR if passport)
- `GET /v1/documents/{id}` — metadata + short-lived signed URL
- `POST /v1/documents/{id}/approve` — staff action
- `POST /v1/documents/{id}/reject` — staff action with reason
- `GET|POST /v1/family-units`, `POST /v1/family-units/{id}/members` — graph maintenance
- `GET /v1/jamaah/{id}/completeness` — doc readiness badge (green/yellow/red per module #89)

**gRPC (service-to-service):**
- `GetJamaah`, `BatchGetJamaah` — master lookup
- `GetPassportData(jamaah_id)` — for `visa-svc` consumption
- `ValidateMahram(subject, group, departure)` — **the headline method; contract detailed in W6**
- `GenerateExternalAccessURL(document_id, ttl, purpose)` — for visa-svc to hand URLs to MOFA/Sajil
- `NotifyDocumentVerified(document_id)` — emits on verify; `visa-svc` listens to flip WAITING_DOCS → READY

## Dependencies

- **F1** (Identity, Access, Audit) — authentication, RBAC, audit writes.
- **External adapters:** GCP Cloud Storage (document vault), GCP Cloud Vision (passport OCR), WhatsApp Business API (reject + expiry alerts), email (fallback alerts).
- **Consumers:** F4 booking (mahram validation at submit, passport-expiry gate), F6 visa (passport data, external URL handoff, `NotifyDocumentVerified`), F7 ops (verification queue, manifest build), F5 payment (none directly), F9 finance (none).

## Backend notes

- **Mahram graph walk** uses Postgres recursive CTE on `mahram_relations`. Cap recursion depth to a small number (e.g. 5 hops) so malformed data can't cause runaway queries. Benchmark the typical 5–15-person family group.
- **OCR worker** runs in a separate goroutine pool with its own rate limiter against GCP Vision. Avoid synchronous OCR on the upload request path — the PRD's flow allows immediate return once the file is in GCS; OCR is async.
- **Signed URL generation** always uses a fresh V4 URL per read request; never cache. External-party URLs go through a dedicated method that logs purpose + TTL to audit.
- **Document-state machine** is enforced at the service layer, not the DB — Postgres can't express "rejected rows are immutable but still exist" without triggers; the service layer is simpler and faster.
- **Name-mismatch detection** between KTP and passport is a simple Levenshtein-distance check with a configurable threshold; policy on what to do with mismatches is Q007.
- **No GCP Vision for KTP in MVP** — manual field entry by staff. KTP OCR is a later iteration once we see how often it would save time; PRD does not mandate.
- **Audit writes are async** from the user-facing request path. The rejection or verification response returns before the audit write completes; use a buffered channel with a shutdown flush.

## Frontend notes

- Three primary surfaces: **jamaah self-upload portal** (mobile-friendly, single-field-at-a-time), **CS/agent upload console** (bulk-upload-capable, shows OCR preview), **ops verification queue** (side-by-side image + editable OCR fields).
- Verification queue is the single highest-touch ops screen — quality here determines how fast a season's jamaah get visa-ready. Worth investing in keyboard shortcuts and pre-loaded previews.
- Progress tracker (module #89) is a shared component — green/yellow/red badge rendered on every jamaah list view.
- Signed-URL lifecycle is backend-managed; frontend just fetches a fresh URL on demand and handles 15-minute expiry gracefully with a re-fetch.

## Open questions

See `docs/07-open-questions/`:
- **Q005** — mahram qualifying relation set, age threshold, same-departure rule
- **Q006** — minimum documents required to submit a booking
- **Q007** — KTP ↔ passport name mismatch handling policy
- **Q008** — UU PDP compliance (consent capture, retention, DSR, DPO)
