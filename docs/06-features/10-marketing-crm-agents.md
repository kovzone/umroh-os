---
id: F10
title: Marketing, CRM, Agent Network
status: draft
last_updated: 2026-04-17
moscow_profile: 11 Must / 14 Should / 13 Could (of the 38 modules in scope)
prd_sections:
  - "B. B2B Front-End (lines 101–171)"
  - "C. Marketing & Sales / CRM (lines 173–239)"
  - "J. Alumni Hub / ZISWAF (lines 647–661, partial)"
  - "Alur Logika 3.1–3.5 (lines 1149–1157), Alur 4.1–4.4 (lines 1161–1167), Alur 10.3 (line 1243)"
  - "Sitemap: lines 713 (B2C ZISWAF), 715–739 (B2B portal), 959–1025 (agent routes), 755/761–775/1035–1055 (CRM)"
  - "Cross-refs: Payroll Komisi line 481, WA template module line 533, PPh 23 line 1289, RBAC scope line 1361"
modules:
  - "#25–27 Pendaftaran Keagenan / E-KYC / E-Signature"
  - "#28–34 Website Replika + Digital Sales Tools (flyer, kartu nama, bank konten, watermark, galeri, tracking codes)"
  - "#35–36 Leads Tracker + Reminder/Tagging"
  - "#37–39 Dompet Komisi (saldo, notifikasi, pencairan)"
  - "#40–44 Academy (LMS, quiz, scripts, leaderboard, push)"
  - "#45–47 Perwakilan (Super-View, Leveling, Overriding)"
  - "#48–54 Digital Marketing (Ads Manager Lite, UTM, Landing Page + A/B, Content Planner/Publisher, Omni, Analytics)"
  - "#55–59 Lead Nurturing (Bot Filter, Pesan Berantai, Pemicu Momen, Segmentasi, Pusat Siaran)"
  - "#60–62 Lead Management (Distribusi Adil, SLA trigger, Rekam Jejak + Tagging)"
  - "#63–66 Sales Closing (Penawaran, Payment Link, e-Approval, Referral Alumni)"
  - "#67–70 Analytics (CS dashboard, ROAS, Retargeting, Radar Prospek)"
  - "#199 Pusat Referral Alumni; #200–202 ZISWAF trio"
depends_on: [F1, F2, F3, F4, F5, F7, F9]
open_questions:
  - Q019 — Abandoned checkout attribution (existing — directly touches UTM)
  - Q045 — Commission accrual timing (existing — binds #37, #47, #481)
  - Q054 — Agent tier taxonomy & qualification thresholds
  - Q055 — Commission % table (per level × per product)
  - Q056 — Overriding formula + hierarchy depth
  - Q057 — UTM attribution model (window, first vs last click)
  - Q058 — Alumni referral reward economics (points / cashback / discount)
  - Q059 — ZISWAF scope (donation platform vs pass-through)
  - Q060 — WhatsApp broadcast rate limits & quality-score handling
  - Q061 — Agent KYC strictness + activation thresholds
  - Q062 — Replica-site white-label scope
  - Q063 — Testimoni moderation policy
  - Q064 — Lead ownership transfer between agents / CS
  - Q065 — Ads API integration depth (Meta/Google)
---

# F10 — Marketing, CRM, Agent Network

## Purpose & personas

F10 owns the **top of the funnel** (marketing, landing pages, ads, lead capture) + the **B2B growth channel** (agent network with hierarchy, replica sites, commissions) + the **relationship layer** (CS workflow, nurturing, broadcasts, alumni community, ZISWAF). It's the loudest service for a growing agency — more jamaah arrive through the agent network and paid ads than through direct brand traffic.

Primary personas:

- **Agent (Agen)** — individual mitra who sells Umrah packages via replica site + WhatsApp close. Earns commission on paid jamaah.
- **Super-agent (Super-Agen)** — upline of multiple agents; earns **override commission** on downline sales (the "selisih jaringan" — PRD line 171).
- **Cabang / Perwakilan** — branch or regional representative; super-view over all downline agents; highest override.
- **CS (Customer Service)** — internal staff handling WhatsApp leads, closing via Payment Link, pushing reminders.
- **Marketing admin** — runs campaigns, configures UTM, builds landing pages, publishes content, tracks ROAS.
- **Content editor** — manages Bank Konten (shared asset library), testimoni gallery, content calendar.
- **Alumni** — past jamaah; issues referral codes, participates in community threads, contributes testimoni.
- **Calon jamaah / B2C lead** — self-service browse → WA inquire → booking conversion path.

Consumers:

- **F4 (booking)** — consumes lead → booking attribution; B2B-closed bookings inherit agent attribution.
- **F9 (finance)** — consumes commission events (`crm.commission_confirmed`) to accrue + payout (per Q045 timing).
- **F11 (dashboards)** — consumes Super-View, CS performance, ROAS, leaderboard as upstream sources.
- **F7 (ops)** — consumes testimoni capture pattern post-trip (Alur 8.7).
- **F12 (alumni daily app)** — shares alumni community surface with F10; F10 owns referral mechanics, F12 owns community UX (per Q059 cross-binding).

## Sources

- PRD Section B in full (lines 101–171) — B2B Front-End: Onboarding, Digital Sales Tools, Leads, Dompet Komisi, Academy, Perwakilan.
- PRD Section C in full (lines 173–239) — Marketing & Sales: Digital Marketing, Perawatan Prospek (Lead Nurturing), Manajemen Prospek (Lead Mgmt), Sales Closing, Analytics.
- PRD Section J (partial, lines 647–661) — Alumni Hub & ZISWAF; only the referral + ZISWAF slice is F10. Full alumni community is F12.
- PRD Alur Logika 3.1–3.5 (lines 1149–1157) + 4.1–4.4 (lines 1161–1167) + 10.3 (line 1243).
- 38 modules enumerated in frontmatter (#25–#70 + #199–#202).

## User workflows

### W1 — Agent self-onboarding + E-KYC + E-Signature (modules #25, #26, #27)

1. Calon agent visits `/mitra/register` (PRD line 971). Form: name, phone, email, NPWP (optional?, see Q061), selects partnership package (Silver / Gold / Platinum tiers per Q054).
2. Uploads KTP, NPWP (if applicable), self-photo for E-KYC (#26).
3. Admin validates uploads (identity-match workflow; may involve CS manual review).
4. On admin approval, system generates an E-Signature MoU (PRD line 111) — agent signs digitally; signed PDF archived per UU ITE requirements (cross-ref Q008).
5. Activation: agent user account created in IAM (F1) with role `AGEN`; default tier assigned per Q054; replica site auto-provisioned at `domain.com/id/<agent-code>` (#28).
6. Welcome WA broadcast: agent receives portal login, replica-site URL, academy link.

### W2 — Replica site + one-click sharing (modules #28, #29, #32)

1. Each agent has a personal page `domain.com/id/<agent-code>` (PRD line 723, 117) with full catalog + agent's WA button + agent's name/photo.
2. Catalog auto-syncs from catalog-svc (Modul #283 Auto-Update Agen): any central product change propagates immediately.
3. Agent shares a product: one-click generates social-ready asset (watermarked flyer with agent's WA number — #32) + share link with embedded UTM (`?utm_source=agen&utm_medium=wa&utm_campaign=<campaign>&utm_content=<agent_code>`).
4. Click tracking: landing on replica site logs lead (visitor → `leads` table) with source attribution.
5. **White-label scope**: **Q062** — MVP default inferred is *no custom branding* (only agent name + photo + WA number; central layout); custom logos/colors deferred to Phase 2.

### W3 — Lead tracker + tagging (modules #35, #36, #62)

1. Every lead recorded with source (organic / replica-agent / ads / direct-WA), UTM, agent_id if applicable, first-touch timestamp.
2. Lead lifecycle: `cold → warm → hot → converted → lost`. Transitions by CS / system.
3. Tags: `Tanya / Janji Bayar / Closed / No Response / Bot` (PRD line 215). Tags composable.
4. Lead detail shows: history (ad clicks, page views on replica site, prior messages), agent/CS owner, follow-up reminders.
5. Reminder + follow-up (#36): CS schedules reminders; push notification fires; task appears in CS dashboard.

### W4 — CS round-robin + SLA (modules #60, #61)

1. Incoming leads (WA inbound, landing page submit, replica-site form) enter a queue keyed by lead_source + channel.
2. Distribusi Adil (#60): round-robin auto-assignment to active CS agents respecting CS availability (online/offline flag).
3. SLA timer (#61): CS has N minutes to respond (**Q066** per Q071 default inferred: 10 minutes globally, configurable per channel).
4. Miss SLA → lead auto-reassigned to next CS in rotation; previous CS flagged (repeated misses affect performance dashboard #67).
5. Manual transfer: CS can hand off a lead to another CS or an agent (**Q064** default: requires supervisor approval + audit log).

### W5 — Lead nurturing drip campaign (modules #55, #56, #58, #57)

1. Bot filter + classifier (#55): inbound WA messages screened; classified cold / warm / hot based on keyword heuristics + prior engagement.
2. Segmentation (#58): jamaah lists built on interest signals (Umroh / Haji / Alumni / Flash-sale) for targeted drip.
3. Drip campaign (#56): D+1 / D+3 / D+7 WA sequence via crm-svc WA adapter.
4. Moment-based triggers (#57): Islamic calendar dates (bulan Rajab / awal Ramadhan), abandoned interest (Q019 cross-ref), re-engagement of dormant leads.
5. Opt-out respected: jamaah who reply STOP exits all drip sequences; marker on lead record.

### W6 — Landing pages + A/B testing (modules #50, #49)

1. Marketing admin opens Landing Page Builder; picks template; edits content blocks (hero, offer, CTA, testimoni, FAQ).
2. Page published at `domain.com/lp/<slug>`. Two variants can be published as A/B test.
3. Traffic split (**Q070** default inferred: 50/50 until statistical significance or until 500 leads accumulated; winner becomes canonical).
4. Analytics: lead conversion per variant, cost-per-lead per variant (joined with ads spend).
5. UTM parameters auto-generated per variant (#49).

### W7 — Ads Manager Lite + ROAS (modules #48, #68)

1. Integration depth — **Q065**: MVP default inferred is *pull-only API integration* with Meta Ads + Google Ads (cost, impressions, clicks per campaign/ad-set/ad). Pushing campaign creation is out of MVP (managed in native ad platforms).
2. Cost data reconciled with finance-svc for marketing-budget actuals (F9 W8 budget-vs-actual); join key = campaign UTM identifier.
3. ROAS (#68): revenue generated by UTM (joined via attribution window per Q057) ÷ ad spend.
4. Retargeting (#69): non-converted warm leads pushed back to Meta as custom audience — **Q069 variant**: Phase 2; MVP is manual export to CSV for upload.
5. Radar Prospek (#70): hot-alert when a dormant lead re-engages (visits replica site / opens broadcast message); pushed to CS dashboard.

### W8 — Attribution + UTM reconciliation (cross-refs Q019, Q057)

1. Every lead + every booking carries a `first_touch_utm` and `last_touch_utm` snapshot (per **Q057** — both models captured; attribution picks one for commission per policy).
2. Attribution window (**Q057** default inferred 30 days from first touch); outside the window, booking is "organic" or attributed to latest direct agent contact.
3. **Abandoned checkout attribution (Q019)**: if jamaah drops off mid-checkout then re-engages via a different agent's link within the window, the original agent retains attribution — resolution per Q019 to finalize.
4. Audit: every attribution decision is logged with source UTMs + decision timestamp + commission impact.

### W9 — Commission calculation + override (modules #37, #47)

1. Trigger: `booking.paid_in_full` event from payment-svc → crm-svc computes commission via `CalculateCommission(booking_id)` gRPC per ADR 0006.
2. Direct commission: closing agent's tier % × booking revenue (minus tax, etc. per Q055).
3. **Override commission (#47)** per **Q056**: upline's tier % − closing agent's tier % × booking revenue. Example: if closing Silver agent has 5% and upline Gold Super-Agen has 8%, override to Super-Agen is (8% − 5%) = 3% of booking revenue. Cabang / Perwakilan (10%) gets (10% − 8%) = 2% above Super-Agen.
4. Writes `commission_ledger` rows: one `direct` for closing agent, one `override` each for upline levels.
5. Status per Q045 default: accrue at `paid_in_full` event, commit on monthly payout run. Pre-payout clawback on refund: reverse entries; post-payout in 30-day window: reverse + mark agent debt.
6. Emits `crm.commission_confirmed` → F9 records accrual journal (W14 in F9).

### W10 — Commission payout (Dompet Komisi — modules #37, #38, #39)

1. Agent sees Saldo Pending / Confirmed / Withdrawal on portal (PRD line 143).
2. Real-time notification (#38): every VA payment from jamaah in the agent's downline triggers a portal + WA push (Pending +X).
3. Payout request (#39): agent submits `PayoutRequest { amount }`; goes to Finance payroll batch (#481).
4. F9 W14 handles the actual payout; F10 reflects status `Withdrawn`.
5. Transaction history per jamaah drill-down: which booking, which commission, when paid.

### W11 — Super-View + Leveling (modules #45, #46)

1. Super-agent / Cabang sees downline performance (dashboard at `/mitra/downline`) — list of sub-agents, their closing counts, earnings, activity (RBAC-gated per PRD line 1019).
2. Leveling otomatis (#46): per Q054 criteria (e.g. ≥ 10 closings in a season + academy quiz passed → promote Silver to Gold); system auto-flips tier with a promotion notification.
3. Demotion: reverse direction if activity drops below threshold for N months (Q054 default: 6-month grace period before demotion).
4. Leaderboard (#43): top-10 agents by closing this month (RBAC-scoped per branch).

### W12 — Academy / LMS (modules #40, #41, #42, #44)

1. Video-based courses: fiqh, sales scripts, tutorials. Agent must complete course N before advancing to tier Gold (per Q054).
2. Quiz + badges (#41): post-course quiz; badges visible on agent profile.
3. Script Jualan (#42): searchable FAQ / sales reply bank — copy to WA with one tap.
4. Push notifications (#44): flash-sale announcements to all active agents.
5. _(Inferred)_ LMS is a lightweight content surface in MVP (static video + quiz); full interactive LMS deferred.

### W13 — Broadcast hub (module #59)

1. Marketing admin composes broadcast: segment (all agents / Silver only / Hot leads / Alumni / custom segment), channel (WA / email / SMS), template from Communication Template module (PRD line 533, F7 infra).
2. Preview + send confirmation; broadcast queued.
3. Rate limiting (**Q060**): default inferred 80 messages/second (Meta's quality-tier-dependent limit) per WA business account; queue drains respecting throughput; per-recipient dedup.
4. Dead-letter / bounce handling: failed messages logged; recipient flagged; quality-score impact monitored (Meta).
5. Opt-out enforcement: recipients on opt-out list filtered out; Sustained-failure recipients auto-suppressed.

### W14 — Sales closing tools (modules #63, #64, #65)

1. Price quote generator (#63): CS or agent picks product + add-ons; system renders PDF quote + WA-friendly text; sends to lead.
2. Payment Link generator (#64): one-click VA or payment link from payment-svc; pushed to lead WA — primary conversion tool.
3. e-Approval discount (#65): if discount > threshold (role-dependent per IAM), approval request routed to manager; approved discount encoded in payment link; audit logged.
4. On successful payment → booking creation path proceeds in F4; commission attribution fires per W9.

### W15 — Alumni referral program (modules #66, #199)

1. Each alumni gets a unique referral code on trip completion.
2. Referral use: new jamaah enters code at booking (F4); attribution system records referral → alumni.
3. **Reward economics (Q058)** default inferred: successful referral (referred jamaah reaches `paid_in_full`) credits alumni with 500K IDR cashback applied to their next Umrah booking. Cap: 3 referrals per year; reward expires after 2 years of non-use.
4. Referral ledger: alumni sees referrals tracked (Pending / Confirmed / Used).
5. Referral ≠ agent commission — alumni reward is a marketing expense, not commission (no PPh 21 treatment).

### W16 — Alumni community + testimoni (modules #33, #66, #199)

1. Post-trip rating & review (Alur 8.7, PRD 1489): jamaah rates muthawwif / hotel / tour leader, writes testimoni, optional photo upload.
2. **Moderation (Q063)** default inferred: testimoni pending-review queue; ops or marketing admin approves / rejects / asks for edit. Approved testimoni surfaces to B2C gallery (Module #41 in F10) + agent replica site (#33 Galeri Dokumentasi) + central content.
3. Community threads (`alumni_threads` table): read-only forum in MVP (F12 owns full community UX); F10 surfaces referral mechanics and testimoni capture.
4. Alumni engagement: re-engagement broadcasts (W13) timed to next Umrah season.

### W17 — ZISWAF (modules #200–#202)

1. **Scope (Q059)** — pivotal product decision: MVP default inferred is *pass-through* — UmrohOS records ZISWAF intent and routes to a LAZ partner for actual donation handling. Does NOT hold donation funds (licensing complexity, legal obligations avoided).
2. Tabungan Niat Kembali (#200): auto-debit savings tied to a future Umrah booking; implementation = integration with existing Tabungan product (F5 Q017) + visual tracker.
3. Zakat Maal calculator (#201): calculator tool with current gold-nisab reference; generates "calculated obligation" number.
4. Sedekah & Infaq pagi (#202): one-tap link to LAZ partner's payment gateway; UmrohOS tracks the click + optional donation-completion receipt.
5. _(Inferred)_ No custody of donation funds in MVP. Pass-through only.

### W18 — Content planner + omnichannel (modules #51, #52, #53, #54)

1. Content calendar (#51): collaborative posting schedule across marketing team.
2. Publisher + scheduler (#52): push post to FB / IG / TikTok / WA status on schedule.
3. Omni-channel distribution (#53): one asset → website + agent replica sites + mobile catalog simultaneously.
4. Analytics (#54): per-post engagement, click-through to booking; feeds F11 marketing dashboard.
5. _(Inferred)_ MVP ships #51 + #52 as minimum; #53 / #54 deferred (Could Have priorities).

### W19 — Retargeting + radar (modules #69, #70)

1. Sinkronisasi Retargeting (#69): non-closed warm leads exported to Facebook / Meta Ads custom audience. MVP = CSV export manual; API-push deferred.
2. Radar Prospek Lama (#70): ML-lite signal — dormant lead returns to site / opens broadcast / visits replica → CS dashboard hot-alert.
3. Scoring: rule-based (recent activity + past engagement + segment) in MVP; no full ML model.

## Acceptance criteria

- **Replica site auto-provisioned** within 30 seconds of agent activation; catalog reflects central within 60 seconds of any change.
- **Commission calculation** fires on `paid_in_full` event; override chain walked per Q056; writes to `commission_ledger`; emits `crm.commission_confirmed` to F9.
- **Commission clawback** on pre-payout refund reverses ledger entries; post-payout within 30 days creates agent debt (Q045 default).
- **UTM attribution** captures first_touch + last_touch; attribution decision audit-logged per booking.
- **CS SLA timer** fires on lead arrival; breach reassigns to next CS in rotation; performance metrics updated.
- **Round-robin fairness** — no CS receives > 2× the average lead count in a 24h window.
- **Broadcast rate limit** — no WhatsApp tier's throughput cap exceeded; opt-out list honored; bounce tracking on.
- **A/B test winner selection** per Q070 criteria; loser variant deactivated with one-click; historical analytics retained.
- **Replica-site white-label scope** per Q062; agent customization boundaries enforced.
- **Agent tier promotion/demotion** per Q054; grace period honored; audit trail on every tier transition.
- **Testimoni moderation** — no unmoderated content appears on public surfaces; rejected testimoni preserved for audit.
- **Referral reward** issued only on `paid_in_full` of referred jamaah; alumni cap enforced per Q058.
- **ZISWAF pass-through** — no UmrohOS custody of donation funds per Q059 default.
- **Lead ownership transfer** — per Q064; audit log captures source → destination + reason.

## Edge cases & error paths

- **Agent disabled / suspended mid-season**: ongoing replica-site visitors see a "agent unavailable" page; existing downline commissions unaffected; pending commissions continue to settle.
- **Overriding chain orphaned** (super-agent deactivated while sub-agent active): override commission skips the deactivated level and flows to the next active upline (or disappears if no upline; confirm Q056).
- **Lead bombardment** (spam-flood): bot filter (#55) rate-limits per source IP / phone pattern; suspected spam auto-filed without CS assignment.
- **WA opt-out mid-drip**: drip campaign immediately halts for that lead; `drip_status = opted_out`.
- **Broadcast rate-limit breach**: queue drains slower; delivery ETA reported; admin warned if delivery time exceeds 60 min.
- **Meta quality-score drop**: automated alert; broadcast throughput throttled automatically.
- **Ads budget overrun**: budget-vs-actual alert (F9 W8) triggers; admin pauses campaign; audit log.
- **UTM tampering** (attacker crafts fraudulent utm_content): replica-site code validated against active agent_code; unknown codes attributed as organic + logged.
- **Attribution window edge** (jamaah books 31 days after first click): organic attribution (default 30-day window Q057).
- **Referral code abuse** (alumni refers family member counted as self-referral): detection rule — same KTP / phone / address match blocks reward; flagged for review.
- **Testimoni with offensive content**: moderation rejects; audit trail preserved per Q063.
- **A/B test with tiny sample**: winner-selection threshold (Q070 default 500 leads OR 2 weeks) prevents premature declaration.
- **Dead-letter broadcast messages**: retry N times with exponential backoff; final failure logs + quality-score impact noted.
- **Agent tier promotion during active downline season**: takes effect from next booking; historical commission calculations unchanged.
- **Lead claimed by two agents** (jamaah clicked two different replica sites): per Q057 policy (first-click vs last-click) resolves; audit logged.
- **ZISWAF click without donation completion**: tracked as intent, not actual donation; no finance journal (per Q059 pass-through default).

## Data & state implications

Owned by `crm-svc` (schema planned at `docs/03-services/09-crm-svc/02-data-model.md`). Key tables referenced / extended:

- `leads` — existing; add `first_touch_utm` + `last_touch_utm` + `attribution_decision` snapshot; `dedup_fingerprint` (phone-hash + email-hash) for unification.
- `campaigns` — existing; add `ads_platform_id_meta`, `ads_platform_id_google` for joining cost data; `ab_parent_id` for A/B variants.
- `agents` — existing; add `tier_history jsonb` (past promotions/demotions), `academy_completions jsonb`, `kyc_status` enum, `replica_site_url`, `replica_visible` bool.
- `commission_ledger` — existing; add `source_type` (`direct | override | referral`), `clawback_of` (nullable FK for reversals), `payout_batch_id` (FK to finance payout run).
- `overrides` — new (or materialized from commission_ledger): per-booking override chain snapshot `{ booking_id, ladder: [{agent_id, tier_pct, override_pct, override_amount}] }`.
- `broadcasts` — existing; add `template_id` (FK to comm template module #533), `target_segment`, `throughput_cap`, `opt_out_honored_count`.
- `broadcast_recipients` — new: per-recipient delivery status `{ broadcast_id, recipient_phone, status (queued/sent/delivered/failed/opted_out), retry_count, last_attempt_at }`.
- `drip_campaigns` + `drip_enrollments` — new: sequence definition + per-lead enrollment state.
- `referral_codes` — existing; add `alumni_id`, `reward_policy` (see Q058), `used_count`, `cap_used`.
- `alumni_threads` — existing; full community schema owned by F12.
- `testimonies` — new: `{ jamaah_id, booking_id, target_kind (muthawwif/hotel/tour_leader/general), rating, text, media_urls jsonb, moderation_status enum, moderated_by, moderated_at, published_at }`.
- `ziswaf_intents` — new: `{ alumni_id (nullable — non-alumni can donate too), kind (zakat/sedekah/tabungan_niat), amount_intent, laz_partner, click_at, completion_url }` — no custody per Q059.
- `ads_cost_pulls` — new: per-campaign daily cost ingest from Meta/Google; joined to campaigns + ROAS calc.

New enums:

- `agent_tier` — see Q054 (current stub lists silver/gold/platinum as example).
- `commission_kind` — extend with `referral` per W15.
- `attribution_model` — `first_click | last_click | linear | position_based` — Q057 picks one.
- `testimoni_moderation_status` — `pending | approved | rejected | edit_requested`.
- `lead_transfer_reason` — `cs_handoff | agent_handoff | escalation | auto_reassign_sla_breach`.

## API surface (high-level)

Full contracts in `docs/03-services/09-crm-svc/01-api.md`. Key surfaces confirmed here:

**REST (agent portal + CRM console):**

- Agent onboarding: `POST /v1/agents/register` + `POST /v1/agents/{id}/kyc` (upload docs) + `POST /v1/agents/{id}/esign-mou` + `POST /v1/agents/{id}/activate` (admin).
- Agent self-serve: `GET /v1/me/dashboard` (Dompet Komisi) + `GET /v1/me/downline` + `POST /v1/me/payout-request` + `GET /v1/me/replica-site-stats`.
- Leads: `GET /v1/leads` + `POST /v1/leads` + `PATCH /v1/leads/{id}` (tags, status) + `POST /v1/leads/{id}/transfer` (Q064).
- Campaigns: `GET /v1/campaigns` + `POST /v1/campaigns` + `POST /v1/campaigns/{id}/ab-variant`.
- Landing pages: `GET /v1/landing-pages` + `POST /v1/landing-pages` + `POST /v1/landing-pages/{id}/publish` + `POST /v1/landing-pages/{id}/ab-winner`.
- Broadcasts: `POST /v1/broadcasts` + `GET /v1/broadcasts/{id}` + `GET /v1/broadcasts/{id}/recipients`.
- Drip: `POST /v1/drip-campaigns` + `POST /v1/drip-enrollments/stop` (opt-out).
- Testimoni: `POST /v1/testimonies` (jamaah) + `POST /v1/testimonies/{id}/moderate` (admin).
- Referral: `POST /v1/referral-codes` + `GET /v1/referral-codes/{code}/redeem`.
- Ads: `GET /v1/ads-costs?campaign_id=&from=&to=` (pulled data; Q065 depth).
- Super-view: `GET /v1/super-view?agent_id=` (RBAC Perwakilan / Super-Agen).
- Leaderboard: `GET /v1/leaderboard?period=`.
- ZISWAF: `GET /v1/ziswaf-intents` + `POST /v1/ziswaf-intents` (intent-only per Q059).

**gRPC (service-to-service):**

- `GetAgentByUserId(user_id)` — F1 / F4 use this.
- `CalculateCommission(booking_id)` — F5 calls on `paid_in_full` (ADR 0006 saga).
- `RecordLeadConversion(lead_id, booking_id)` — F4 calls on booking creation.
- `ClawbackCommission(booking_id, reason)` — F5 calls on refund.
- `AttributeBooking(booking_id)` — F4 calls to pin attribution decision.
- `SendBroadcast(segment, template, vars)` — other services can trigger WA notifications via this (gateway into the broadcast hub).

**Events emitted** (`03-events.md`):

- `crm.lead_converted` — existing.
- `crm.commission_confirmed` — new (accrual event to F9).
- `crm.commission_paid_out` — new (payout completed event; triggers WA receipt).
- `crm.commission_clawed_back` — new.
- `crm.broadcast_completed` — existing.
- `crm.referral_rewarded` — new.
- `crm.testimoni_approved` — new (for marketing asset pipeline).
- `crm.ziswaf_intent_recorded` — new.

## Dependencies

- **F1 (IAM)** — agent is user with `AGEN` role + tier + data-scope (`PERSONAL/SELF` per PRD line 1361); agent hierarchy = IAM-tracked relationship.
- **F2 (catalog)** — replica sites + mobile catalog consume catalog-svc; `catalog.product_updated` events trigger replica-site cache invalidation (Modul #283 Auto-Update Agen).
- **F3 (jamaah)** — lead → jamaah conversion path; duplicated-contact fingerprinting joins leads to jamaah.
- **F4 (booking)** — booking attribution snapshot; B2B-closed bookings attach agent_id; agent-as-initiator flow for manual booking creation.
- **F5 (payment)** — `paid_in_full` event triggers commission calculation; Payment Link generator (#64) calls payment-svc; refund event triggers clawback.
- **F6 (visa)** — visa.issued / visa.rejected events → customer WA via broadcast adapter.
- **F7 (ops)** — testimoni capture pattern (Alur 8.7); document-rejection WA notifications via broadcast; broadcast template management shared.
- **F8 (logistics)** — shipment.dispatched / delivered → customer WA via broadcast.
- **F9 (finance)** — commission accrual + payout per Q045; referral rewards as marketing expense; ads budget reconciliation.
- **F11 (dashboards)** — Super-View, CS performance, ROAS, leaderboard all flow upward to F11.
- **F12 (alumni daily app)** — shares community surface; F10 owns referral mechanics, F12 owns daily-app UX.
- **External** — Meta Graph API (Ads + WA Business), Google Ads API, Meta WhatsApp Business API (primary broadcast surface), LAZ partner API (ZISWAF pass-through).

## Backend notes

- **Commission calculator** is a pure function: input (booking_revenue, tier_ladder, clawback_window, vat/pph policy) → output (per-agent amount breakdown). Unit-tested exhaustively with edge cases: partial payments, refunds mid-payout, orphaned upline (inactive super-agent), tier-promotion mid-booking.
- **Override chain walk** (`overrides` table): recursive CTE or iterative loop from closing agent up through parent_agent_id chain; stops at null parent OR inactive agent (Q056 pins behavior).
- **Attribution decider** (W8): for each booking, compute both first-touch and last-touch; Q057 decides which wins; both recorded for audit even when only one drives commission.
- **Broadcast throughput** — Meta's WA Business API has quality-tier throughput caps (1K → 10K → 100K per 24h). Track quality tier per agency; queue respects current tier; alerts on tier downgrade.
- **Dead-letter queue** on broadcast failures: retry 3× with exponential backoff (5s / 30s / 5min); after 3rd fail, recipient suppressed; alerts admin.
- **Lead deduplication** via fingerprint `SHA256(lower(phone) || lower(email))`; joins agent-lead-jamaah as stable identity.
- **A/B test winner selection** — simple rule-based: minimum 500 leads OR 14 days elapsed; winner = variant with higher conversion (lead → booking); if tie, earlier-created variant.
- **UTM injection on replica share** — server-side render ensures every outbound link from replica site carries agent_code; URL shortener (optional) preserves UTM via redirect.
- **KYC verification** — MVP: admin manual review. Automated KYC (liveness, document OCR) is Phase 2.
- **E-Signature** — MVP: in-app click-to-sign generating signed PDF + timestamp + IP; full QSCD-level E-Signature (Tanda Tangan Elektronik Tersertifikasi per UU ITE) is Phase 2 — cross-ref Q008 UU PDP.
- **Ads API integration** — per Q065 default pull-only: daily cron reads spend/clicks/impressions per campaign via Meta Ads API + Google Ads API; joins to `campaigns.ads_platform_id_*`. No push of campaign creation.
- **ZISWAF pass-through** — no custody; integration is an outbound link with tracking pixel for click-attribution; completion confirmed via LAZ partner webhook (if integration available) or self-reported.

## Frontend notes

- **Agent portal** (Svelte) — dashboard (Dompet Komisi real-time), lead tracker, replica-site preview + share tools, downline view (if Super-Agen+), academy, payout request.
- **Replica site** — public-facing; catalog browse; WA button prominent; flyer share. Minimal customization (Q062 default).
- **CRM console** (admin web app) — lead queue + assignment, CS performance dashboard, campaign builder, landing-page builder, broadcast composer, testimoni moderation queue, agent onboarding verification, referral management, ZISWAF intent log.
- **CS workbench** — inbox-style lead view with WA thread, quick-reply scripts (#42), tag controls, payment-link generator inline.
- **Marketing admin workbench** — UTM link generator, A/B test comparison view, ROAS dashboard, ads-cost pull status, content calendar.
- **Leaderboard** — gamified tile on agent portal + admin view; respects RBAC branch scope.
- **Notifications** — real-time WebSocket push for CS lead assignment, agent commission accrual, broadcast completion; WA push for agent Dompet Komisi updates (#38).
- **Mobile-first** — agent portal must work well on phone (most agents are mobile-only); CS workbench is desktop-primary but usable on tablet.

## Open questions

See `docs/07-open-questions/`.

**Existing, binding:**

- **Q019** — Abandoned checkout attribution (directly drives W8 attribution logic)
- **Q045** — Commission accrual timing (W9 + W10 use this)

**New, filed with this draft (Q054–Q065):**

- **Q054** — Agent tier taxonomy + qualification thresholds + demotion rules
- **Q055** — Commission % table (per level × per product)
- **Q056** — Overriding formula + hierarchy depth + orphaned-upline handling
- **Q057** — UTM attribution model (window, first vs last click)
- **Q058** — Alumni referral reward economics
- **Q059** — ZISWAF scope (donation platform vs pass-through)
- **Q060** — WhatsApp broadcast rate limits + quality-score handling
- **Q061** — Agent KYC strictness + activation thresholds
- **Q062** — Replica-site white-label scope
- **Q063** — Testimoni moderation policy
- **Q064** — Lead ownership transfer between agents / CS
- **Q065** — Ads API integration depth (Meta/Google)

**Inferred (pending reviewer confirmation):**

- CS SLA 10 minutes globally (per-channel configurable later) — could be filed as a thin Q but inferred for MVP.
- A/B winner criteria: 500 leads OR 14 days, higher conversion rate wins (tie → earlier created) — Q070 candidate, inferred.
- MVP KYC = admin manual review; automated KYC Phase 2.
- E-Signature = in-app click-sign with signed PDF + timestamp; QSCD-level deferred.
- Retargeting = CSV manual export for MVP; API push Phase 2.
- Content planner modules #51–#54: MVP ships #51 + #52; #53 + #54 deferred.
- Referral code abuse detection: same KTP/phone/address blocks reward.
- Agent tier example values (Silver/Gold/Platinum): placeholder only — Q054 finalizes.
- Alumni community surface: F10 owns referral mechanics; F12 owns thread/fatwa UX (split confirmed with Q059 scope).
