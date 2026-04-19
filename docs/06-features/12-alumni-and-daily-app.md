---
id: F12
title: Alumni Hub & Daily App
status: written
last_updated: 2026-04-18
moscow_profile: 0 Must / 3 Should / 6 Could (9 modules in F12 scope; 4 ZISWAF/referral modules owned by F10)
prd_sections:
  - "J. Daily App & Alumni Hub (lines 623–661)"
  - "Alur Logika 10.1–10.3 (lines 1237–1243) — headings only; F12 fills the narrative"
  - "Sitemap: lines 707–713 (/jamaah/certificate, /jamaah/community, /jamaah/fatwa, /jamaah/ziswaf)"
  - "Ekosistem Alumni block: lines 945–957"
  - "Category-to-service mapping: docs/00-overview/03-module-list.md line 68 (crm-svc owns alumni/community)"
modules:
  - "#190 Jadwal Shalat & Adzan, #191 Kompas Arah Kiblat, #192 Al-Quran Digital, #193 Dzikir & Kumpulan Doa (Daily Worship)"
  - "#194 Ensiklopedia Manasik, #195 Artikel & Kajian Rutin, #196 Tanya Jawab Agama / Fatwa Desk (Islamic Knowledge)"
  - "#197 Forum Grup Angkatan, #198 Papan Informasi Reuni (Alumni community)"
  - "#199 Pusat Referral Alumni — **owned by F10**; F12 hosts UI entry-point only"
  - "#200–#202 ZISWAF trio — **owned by F10** per Q059 pass-through"
depends_on: [F1, F3, F4, F10]
open_questions: []
---

# F12 — Alumni Hub & Daily App

## Purpose & personas

F12 is **post-pilgrimage engagement** — where jamaah who have completed their Umrah stay connected with the agency, use daily-worship tools, consume Islamic knowledge, participate in alumni community, and (via F10-owned mechanics) refer new jamaah + donate.

This is **the lightest feature in the catalogue** by both MoSCoW weight (zero Must Haves) and PRD narrative depth (Alur Logika 10.x is one-line bullets; no step-by-step flow). The reviewer may choose to defer parts of F12 to Phase 2 — see **Q081** for the MVP scope-carve discussion.

F12 also has a clean boundary with F10: F10 already owns modules #199 (referral mechanics) and #200–#202 (ZISWAF pass-through per Q059). F12 hosts the **UI entry-points** to those F10 surfaces inside the alumni shell but does not re-implement the mechanics. F12 owns #190–#198 — nine modules across Daily Worship, Islamic Knowledge, and Alumni community.

Primary personas:

- **Alumni** — past jamaah who completed a trip (`booking.departure_completed` past date). Primary consumer of every F12 surface.
- **Calon jamaah / B2C guest** — non-alumni can access some Daily App utilities (prayer times, qibla, Quran) as a loss-leader engagement; community + fatwa gated to alumni.
- **Ustadz / religious advisor** — responds to Fatwa Desk Q&A (#196); authority and selection per Q078.
- **Content editor** — maintains manasik encyclopedia (#194) + article feed (#195) + reuni announcements (#198).
- **Community moderator** — moderates forum threads (#197); role per Q077.
- **Agency marketing admin** — uses F12 surfaces as brand/retention touch-point; measures DAU/MAU.

F12 is a **jamaah-facing read-mostly** product — unlike F1–F11 which are heavy transactional / ops. Engagement metrics (daily active users, retention) matter more than throughput.

## Sources

- PRD Section J in full (lines 623–661).
- Alur Logika 10.1–10.3 (lines 1237–1243) — bullet headings only.
- Sitemap jamaah routes (lines 707–713): `/jamaah/certificate`, `/jamaah/community`, `/jamaah/fatwa`, `/jamaah/ziswaf` (ziswaf is F10-owned; F12 hosts entry).
- Ekosistem Alumni block (lines 945–957).
- Category-to-service mapping: F12 co-tenants inside `crm-svc` (no dedicated alumni-svc exists in service map per `docs/01-architecture/02-service-map.md`).
- 9 modules (F12-owned) enumerated in frontmatter.

## User workflows

### W1 — Prayer times + adzan push (module #190)

1. Alumni / guest opens Daily App; app requests geolocation permission.
2. App computes today's 5 prayer times via chosen calculation method (**Q080** — e.g. Kemenag Indonesia, MWL, Umm Al-Qura for Saudi usage).
3. Next-prayer countdown visible on home screen; swipe for week ahead.
4. Push notification fires at each adhan time (opt-outable); adzan audio plays if permissions allowed.
5. Location-based — traveling to Saudi shifts prayer times automatically.
6. Offline: cached prayer times for current week; recompute if location changes.

### W2 — Qibla compass (module #191)

1. Alumni opens `/jamaah/qibla` or Daily App qibla screen.
2. App uses device gyroscope + magnetometer + location to compute bearing to Kaaba (21.4225°N, 39.8262°E).
3. Rotating compass needle shows direction; haptic tap on alignment.
4. Device-sensor heavy — requires mobile form factor (per **Q076**); fallback web view shows static Google Maps heading arrow.
5. Accuracy caveat / calibration instructions shown on first use.

### W3 — Al-Quran digital (module #192, Could)

1. Alumni browses mushaf by surah / juz / page; text with tajwid coloring.
2. Audio recitation per verse (murottal) — multiple reciter selection.
3. Translation overlay (Bahasa Indonesia default; English optional).
4. Bookmarks + last-read memory per user.
5. Content licensing per **Q079** — Tanzil / King Fahd Quran Complex source; licensing verified.
6. Offline: full Quran text cached after first load; audio streams on-demand.

### W4 — Dzikir + kumpulan doa (module #193, Could)

1. Categorized dzikir: pagi (morning), petang (evening), after-prayer, traveler's, daily.
2. Tasbih counter on-screen — tap to increment; long-tap to reset.
3. Audio-play option for each dzikir (male / female reciter).
4. Static content; no user-contributed content.

### W5 — Manasik encyclopedia (module #194, Should)

1. Step-by-step Umrah + Hajj ritual guide: ihram → talbiyah → tawaf → sa'i → tahallul → ziyarah.
2. Text + short-form video + audio per step.
3. Pre-trip jamaah reads as preparation; during-trip as reference; post-trip for understanding.
4. Content ownership per **Q079** — agency-authored (signed-off by religious advisor) or licensed external (Kemenag Manasik Umrah material).
5. Offline-first — full content downloadable for Saudi-side usage (data roaming expensive).

### W6 — Artikel & kajian rutin (module #195, Could)

1. Daily article feed — agency-curated + ustadz-authored; daily subuh email / push.
2. Kajian schedule — upcoming live-streams (YouTube / Zoom link); calendar view.
3. Past kajian archive — video-on-demand.
4. Comment section disabled in MVP (reduces moderation burden); feedback via WA only.

### W7 — Fatwa Desk (module #196, Could)

1. Alumni submits question to `/jamaah/fatwa`; categorized (fiqih ibadah, muamalah, umroh-specific, etc.); optional anonymity.
2. Ustadz reviews queue; drafts answer; publishes.
3. **Q078 defaults** — authority and liability: agency's panel of dewan asatidz (curated ustadz list) answers; disclaimers clear ("for guidance, not formal fatwa"); SLA 72 hours for answer.
4. Published fatwa visible in public archive (searchable); question author identity shown only with consent.
5. Related-fatwa suggestions shown before user submits (reduces duplication).

### W8 — Forum grup angkatan (module #197, Could)

1. Each jamaah cohort (per `package_departure_id`) gets a private forum thread.
2. Alumni post text + photo in cohort thread; tag fellow alumni.
3. **Q077 defaults** — post-approval moderation (publish fast, flag-and-review); moderator role = agency CS supervisor or dedicated community manager.
4. Reactions (emoji) + reply threading; mention notifications.
5. Reporting: flag inappropriate posts; auto-hide on 3 flags pending moderator review.
6. Retention: forums active for 2 years post-departure; archived thereafter (read-only).

### W9 — Papan informasi reuni (module #198, Could)

1. Agency posts upcoming reunion / kopdar / tabligh akbar events to a public board.
2. Alumni RSVP + invite; calendar integration.
3. Attendance tracking via QR check-in at event.
4. Past events archived with photo gallery.
5. Admin-only posting; alumni consume.

### W10 — Referral + ZISWAF entry-points (cross-ref F10)

1. Alumni shell surfaces entry-points to F10-owned pages:
   - `/jamaah/referral` → displays alumni's referral code + reward status; full mechanics in F10 W15.
   - `/jamaah/ziswaf` → entry to ZISWAF tools; pass-through to LAZ partner per Q059.
2. F12 **hosts** the entry; F10 **owns** the data + behavior.

### W11 — Engagement metrics for F11

1. F12 tracks DAU, MAU, per-module usage, retention (30d / 60d / 90d).
2. Surfaced to F11 dashboards — marketing admin sees engagement as alumni-retention KPI.
3. Metrics aggregated, privacy-respecting — no per-user tracking exposed outside F1 RBAC.

## Acceptance criteria

- **Prayer times** correct for location within 30s of GPS fix; adzan push fires within 10s of prayer time.
- **Qibla** bearing accurate to ±5° with proper device calibration.
- **Quran** text renders correctly with tajwid coloring; audio streams without stuttering on 4G.
- **Manasik** full content downloadable for offline Saudi-side use (< 500MB bundle).
- **Fatwa Desk** 72h SLA met on > 90% of submissions (Q078 default).
- **Forum** moderation flags trigger review within 24h; no inappropriate content persists > 48h.
- **Engagement metrics** surfaced to F11 with proper RBAC (aggregate only; no per-user identity).
- **F10 entry-points** correctly route to F10 screens; scope/auth preserved.
- **Offline-first** for Daily Worship + manasik — core surfaces functional without network.
- **Data-scope** per F1 — alumni sees own agency's content; no cross-agency leakage.
- **Content licensing** documented (Q079) — no unlicensed Quran / dzikir / manasik content shipped.

## Edge cases & error paths

- **GPS permission denied** — fall back to user-entered city for prayer times; show instructional tooltip.
- **Magnetometer unavailable / uncalibrated** — qibla screen shows calibration instruction ("move phone in figure-8"); fallback to static-map heading.
- **Quran audio CDN failure** — skip audio, text-only.
- **Offline bundle outdated** (manasik content updated) — sync on next network; show version badge.
- **Fatwa Desk overdue** — beyond 72h, push priority alert to ustadz panel + CS supervisor; response escalation.
- **Forum spam** — rate-limit per user (10 posts/hour); auto-flag duplicate content.
- **Forum cross-tagging** (non-cohort member tagged): tag works but cross-cohort-thread visibility not granted; no data leak.
- **Cohort merge / split** (rare — two departures merged operationally): forum threads stay separate; no automatic merge.
- **Alumni leaving agency** (switches to competitor) — forum access revoked; historical content retained per data-retention policy.
- **Jamaah deceased** — forum posts preserved; memorial-badge option for family to request.
- **Pre-trip jamaah using Daily App** — allowed for prayer times + qibla + manasik preparation; full alumni features (community, fatwa, referral) gated until first completed trip.

## Data & state implications

F12 co-tenants inside `crm-svc`. No new service; new tables:

- `forum_threads` — per-cohort `{ id, package_departure_id, created_at, last_post_at, post_count, status }`.
- `forum_posts` — `{ id, thread_id, author_jamaah_id, body, media_urls jsonb, reactions jsonb, flagged_count, moderation_status, created_at }`.
- `fatwa_questions` — `{ id, author_jamaah_id nullable (anonymous), category, body, status enum (submitted/answering/answered/rejected), submitted_at, answered_at }`.
- `fatwa_answers` — `{ id, question_id, ustadz_id, body, published_at, citations jsonb }`.
- `fatwa_categories` — seeded list: fiqih ibadah, muamalah, umroh-specific, keluarga, zakat, lainnya.
- `reuni_events` — `{ id, title, description, date, location, rsvp_count, status }`.
- `reuni_rsvps` — `{ event_id, jamaah_id, attending_count, checked_in_at nullable }`.
- `content_articles` — `{ id, title, slug, body, author, published_at, category }`.
- `content_kajian_schedule` — `{ id, title, stream_url, scheduled_at, duration, status enum }`.
- `daily_app_usage` — `{ jamaah_id, module_name, usage_date, usage_count }` aggregated daily (privacy-respecting).

No changes to existing F10 tables (referral / ZISWAF) — F12 reads them via F10's gRPC.

**Data ownership split with F10**: F12 owns community + fatwa + content tables; F10 owns `referral_codes`, `alumni_referral_credits`, `ziswaf_intents`. Clean bounded context.

## API surface (high-level)

Surfaces live inside `crm-svc/api/rest_oapi/` (co-tenant). Routes under `/v1/community/`, `/v1/fatwa/`, `/v1/daily-app/`.

**REST (jamaah portal / Daily App):**

- Daily App utilities (mostly client-side logic; backend only for content):
  - `GET /v1/daily-app/prayer-times?lat=&lng=&date=` — computed server-side or client-side with API fallback.
  - `GET /v1/daily-app/quran?surah=&ayah=` — text + audio URLs (proxying licensed source per Q079).
  - `GET /v1/daily-app/dzikir?category=` — static content.
  - `GET /v1/daily-app/manasik?step=` — content delivery with offline-bundle manifest.
- Content feeds:
  - `GET /v1/content/articles?category=&page=` — paginated article list.
  - `GET /v1/content/articles/{slug}` — detail.
  - `GET /v1/content/kajian/schedule` + `GET /v1/content/kajian/archive`.
- Fatwa Desk:
  - `POST /v1/fatwa/questions` (alumni submit).
  - `GET /v1/fatwa/questions?category=&my=` — public archive + my-submissions view.
  - `POST /v1/fatwa/questions/{id}/answer` (ustadz role) + `POST /v1/fatwa/questions/{id}/publish`.
- Forum:
  - `GET /v1/community/threads?departure_id=` — my cohort's thread.
  - `POST /v1/community/posts` (create) + `POST /v1/community/posts/{id}/flag`.
  - `POST /v1/community/posts/{id}/moderate` (moderator role).
- Reuni:
  - `GET /v1/reuni/events` + `POST /v1/reuni/events/{id}/rsvp`.
  - `POST /v1/reuni/events/{id}/checkin` (QR at event).
- Usage metrics (internal, for F11):
  - `GET /v1/daily-app/metrics?range=` — aggregated DAU/MAU/retention.

**gRPC (service-to-service):**

- `GetAlumniStatus(jamaah_id)` — F10 calls to validate alumni eligibility for referral.
- `RecordUsageEvent(jamaah_id, module, count)` — Daily App client pushes usage events for metrics.

No events emitted specifically for F12 in MVP; F11 reads metrics directly via gRPC.

## Dependencies

- **F1** — alumni auth via jamaah portal; RBAC roles `ALUMNI`, `USTADZ`, `COMMUNITY_MOD`.
- **F3** — jamaah identity; alumni status derived from completed bookings.
- **F4** — `booking.departure_completed` event flips jamaah to alumni status; cohort assignment from `package_departure_id`.
- **F10** — referral mechanics (#199) and ZISWAF pass-through (#200–#202); F12 hosts entry-points.
- **F11** — consumes F12 engagement metrics as an alumni-retention KPI.
- **External** — prayer times API (Q080; candidates Aladhan.com, IslamicFinder, self-compute from Institute of Geophysics Tehran algorithm), Quran text (Tanzil or King Fahd Complex corpus per Q079), murottal audio (Quran.com API or licensed CDN), manasik video CDN.

## Backend notes

- **Co-tenancy in crm-svc** — F12's tables live in `crm` Postgres schema (ADR 0007 single-DB-multi-schema). No new Go module. F12-specific handlers in `crm-svc/api/rest_oapi/handlers/` alongside CRM / marketing handlers.
- **Offline-first strategy** — Daily App manifest endpoint lists all cacheable resources + version; service-worker (if PWA per Q076) caches on first load; updates via ETag/version check.
- **Content delivery** — static Quran / dzikir / manasik content served from CDN (CloudFlare or similar); dynamic content (articles, fatwa, forum) from crm-svc direct.
- **Moderation** — flagged-posts queue surfaces in CS supervisor console; auto-hide at 3 flags; re-hide after moderator action.
- **Fatwa workflow** — queued questions distributed round-robin among active ustadz panel (or claim-based); SLA monitoring via cron; overdue escalations to agency owner.
- **Metrics aggregation** — daily cron rolls up `daily_app_usage` into weekly / monthly summaries for F11 consumption; raw data retained 90 days then summarized.
- **Privacy** — per PDP (Q008), per-user usage not exposed outside F1 RBAC; aggregate-only metrics to F11.
- **Internationalization** — content Bahasa-Indonesia default; Arabic for Quran; English secondary for translations. Locale negotiation per jamaah preference.
- **Daily App as jamaah-portal-internal vs separate mobile app** (Q076) — architecture impact:
  - Responsive-web-inside-portal (default Q076 recommendation): no separate build; reuses jamaah-portal auth; simpler.
  - Separate native app: requires two builds + app-store approval + separate auth flow; higher ongoing cost.
  - Svelte PWA: installable; offline-first; single codebase; middle ground.

## Frontend notes

- **Surface**: per Q076 default — responsive web within jamaah portal + Svelte PWA upgrade path; installable on phone via browser prompt; mobile-first layout.
- **Daily Worship screens** — widget-style home with prayer countdown hero; qibla launch tile; Quran / dzikir / manasik quick-access.
- **Community forum** — Twitter-style timeline per cohort; media upload inline; reactions + replies.
- **Fatwa Desk** — submission form + category picker + browseable public archive.
- **Manasik** — accordion-navigated step-by-step guide; video inline; photo gallery; download-offline toggle.
- **Reuni board** — card layout of events; RSVP buttons; event detail page with photo gallery.
- **Sensor access** — qibla requires gyroscope permission; gracefully degrades if denied.
- **Offline indicator** — clear visual when cached content is shown.
- **Push notifications** — browser-native Web Push for subuh reminders + kajian start; opt-in at first app open.

## Open questions

None blocking — **Q058, Q059, Q063, Q076–Q081** answered **2026-04-18** (`docs/07-open-questions/`). Defaults below match those answers; **Q079** licensing still needs **legal** confirmation before shipping Quran/murottal.

**Residual defaults (engineering / content ops):**

- Default prayer-times calculation method: Kemenag Indonesia (domestic); Umm Al-Qura (Saudi-side auto-switch based on location).
- Default Quran corpus: Tanzil + Indonesian Kemenag translation (requires license confirmation per Q079).
- Default qibla coordinates: Kaaba 21.4225°N, 39.8262°E.
- Forum post-moderation cadence: within 24h; auto-hide at 3 flags.
- Fatwa SLA: 72h for first answer.
- Manasik content: agency-authored with religious advisor sign-off (Phase 1); external-source licensing (Phase 2 option).
- Alumni retention for forum: 2 years post-departure active; archived thereafter.
- Engagement metrics raw retention: 90 days; aggregated retention: per Q070 (3 years on dashboards).
- Daily App form factor: responsive web + Svelte PWA upgrade; no separate native app for MVP.
- **MVP scope carve (Q081)**: IN — #190 prayer times, #191 qibla, #194 manasik read-only, alumni shell hosting F10 referral entry. OUT (Phase 2) — #192 Quran, #193 dzikir, #195 articles/kajian, #196 fatwa desk, #197 forum, #198 reuni board. Rationale: zero Must Haves + content-licensing + moderation complexity in the OUT set; IN set is client-side-heavy with minimal backend.
