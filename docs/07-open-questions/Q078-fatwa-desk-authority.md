---
id: Q078
title: Fatwa Desk ustadz authority + response SLA + liability posture
asked_by: session 2026-04-17 F12 draft
asked_date: 2026-04-17
blocks: F12
status: answered
---

# Q078 — Fatwa Desk ustadz authority + SLA + liability

## Context

Fatwa Desk (module #196) lets alumni submit religious questions to the agency's panel of ustadz. Religious-opinion content is sensitive:

- **Theological authority** — who counts as qualified to answer? Not every staff member with religious knowledge is authorized to issue fatwa-like content.
- **Liability** — "for guidance only" vs "formal fatwa" disclaimer posture matters both for UU ITE and for religious integrity.
- **Response timeliness** — a 2-week response means alumni won't use the feature; a 1-hour SLA can't be met by volunteer ustadz.
- **Mis-answered content** — once published, religious guidance is "sticky"; correction protocols matter.

PRD is silent on all of this.

## The question

1. **Ustadz panel selection** — who's on the panel? Internal agency staff, contracted ustadz, volunteer pool?
2. **Authority level** — are answers "guidance" or "fatwa"? Liability disclaimers?
3. **Response SLA** — 24h? 72h? One-week?
4. **Queue distribution** — round-robin among panel, claim-based, category-based routing?
5. **Private vs public answers** — all published to archive, or some private to submitter?
6. **Correction workflow** — if an answer is later disputed or found errant?
7. **Off-scope rejection** — what happens when a question is off-topic or inappropriate?
8. **Anonymous submission** — is anonymity supported for sensitive questions?

## Options considered

- **Option A — Curated ustadz panel (3–5 contracted ustadz); 72h SLA; "guidance" framing; public archive.** Mid-weight commitment.
  - Pros: manageable; guidance framing reduces liability vs formal fatwa; public archive grows knowledge base.
  - Cons: SLA slippage if panel under-resourced.
- **Option B — Internal staff only (single religious advisor); 7-day SLA; all private (no archive).** Minimal.
  - Pros: smallest footprint; lowest liability (no public record).
  - Cons: slow; no searchable archive reduces reuse value.
- **Option C — Open volunteer panel + editorial review; 48h SLA; public archive with editorial disclaimer.** Community model.
  - Pros: scalable; resilient to individual unavailability.
  - Cons: quality-control burden; authority disputes between volunteers.
- **Option D — Defer Fatwa Desk to Phase 2 entirely.** Module is Could Have; no loss for lean MVP.
  - Pros: avoids all complexity; matches Q081 lean-MVP recommendation.
  - Cons: loses a visible feature some alumni expect.

## Recommendation

**Option A — curated ustadz panel (3–5 contracted); 72h SLA; "panduan / guidance" framing (explicitly not formal fatwa); public archive with author attribution.**

Option B's private-only model misses the reuse value (most alumni questions duplicate previous ones; a searchable archive reduces load). Option C introduces authority disputes + quality-control overhead the agency isn't positioned to manage. Option D is also viable if the agency confirms Fatwa Desk isn't MVP-critical (Q081 lean-MVP default already puts it in OUT) — but if we ship it, Option A is the right structure.

Defaults to propose: **Panel size** = 3–5 ustadz contracted as external advisors; each dedicates ~5 hours/week to queue. **Authority framing** = clear disclaimer on every answer: *"Panduan ini bersifat informatif dan bukan fatwa resmi. Konsultasikan dengan otoritas keagamaan setempat untuk kasus spesifik."* **SLA** = 72h first answer; overdue → escalation to panel lead + agency owner. **Queue distribution** = category-based routing (fiqh ibadah → ustadz specializing in ibadah; muamalah → ustadz specializing in muamalah); fallback round-robin. **Public archive** = default; author can request private-only submission. **Correction workflow** = answer can be edited within 7 days of publish with change-log badge ("Diperbarui HH:MM"); beyond 7 days, correction posts as new answer referencing the original. **Off-scope rejection** = ustadz can flag question as off-scope with reason (redirect to proper resource, e.g. medical → doctor); submitter sees rejection + explanation. **Anonymity** = supported — submitter can toggle anonymous at submit; author identity never shown publicly, only internally to ustadz + moderator. **Credentials display** = each ustadz has a profile with qualifications (pesantren background, Al-Azhar / LIPIA degree, current affiliation) visible alongside answers.

Reversibility: moving from Option A to Option B (private only) is a config change. Expanding panel is hiring. Switching to Option C volunteer is a bigger policy change.

## Answer

**Decided:** **Option A** — **3–5 contracted ustadz**, **72h SLA**, **“bukan fatwa resmi” disclaimer**, **public archive default** with edit log **7d**, **category routing**, **anonymous allowed**, **off-scope reject** with explanation. *(If Q081 defers F12 modules, ship Fatwa only when module enabled.)*

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
