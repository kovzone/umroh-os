---
id: Q008
title: UU PDP compliance — consent, retention, DSR, DPO
asked_by: session 2026-04-14 F3 draft
asked_date: 2026-04-14
blocks: F3, F1
status: answered
---

# Q008 — UU PDP compliance (consent, retention, data-subject rights, DPO notification)

## Context

UU No. 27/2022 (Undang-Undang Perlindungan Data Pribadi, "UU PDP") is Indonesia's personal data protection law, in full force since October 2024. It applies to UmrohOS because we process:
- **Sensitive personal data** (biometrics via OCR, health data via vaccination records)
- **General personal data** (KTP, passport, contact, family relations)

The law requires, among other things:
- **Explicit consent** at collection, with purpose disclosed
- **Legitimate basis** for processing (consent is one of six; others include contract performance, legal obligation, vital interest)
- **Retention policy** — data must be deleted or anonymised when no longer needed for the declared purpose
- **Data subject rights (DSR)** — access, rectification, erasure, portability, objection
- **Data Protection Officer (DPO)** notification within 3×24 hours of a breach affecting data subjects
- **Cross-border transfer rules** — relevant because MOFA/Sajil in Saudi Arabia receives passport data

The PRD does not address UU PDP anywhere. This is a gap that will eventually draw a regulator letter if unaddressed.

This question blocks F3 (document storage, retention, consent capture) and touches F1 (audit, access control, DSR fulfilment).

## The question

Five linked sub-questions:

1. **Legitimate basis.** Is the processing basis **explicit consent** (captured at registration) or **contract performance** (necessary for providing the pilgrimage service)? Both are defensible; the choice affects consent-flow UX.
2. **Consent capture UX.** Where and how is consent collected — at jamaah registration, at first upload, both? Is it a single bundled consent or granular (separate opt-ins for OCR, for sharing with Saudi MOFA, for WhatsApp contact, etc.)?
3. **Retention policy.** How long do we keep documents and OCR results after pilgrimage completion? Options: minimum legal retention (7 years for financial records per PSAK, overrides shorter policies), alumni-opt-in for longer retention, purge after 1 year by default.
4. **DSR fulfilment SLA.** UU PDP requires "reasonable time" for DSR requests. What's our target — 7 days, 14 days, 30 days? What's the ops workflow for fulfilling an erasure request given financial records must be kept?
5. **DPO appointment.** Is there a Data Protection Officer named for the agency? If not, who fills the role for breach-notification purposes?

## Options considered

This is really five decisions, each with multiple options. Summarising the dominant tradeoffs:

- **Basis = consent vs contract.** Consent requires renewal and can be withdrawn, which conflicts with visa / immigration record-keeping obligations. Contract performance is legally cleaner for the core processing; consent is still needed for anything beyond the contract (marketing, future broadcasts).
- **Retention = strict vs generous.** Strict (purge after 1 year) respects the principle of data minimisation. Generous (7 years to match PSAK) keeps the audit story simple. Hybrid (purge documents after 1 year but keep financial records) is the compromise.
- **DSR SLA.** The only meaningful threshold is: can we fulfil within 30 days without heroics? If yes, 30 days. If we need automation, invest in the tooling now.
- **DPO.** For a single-tenant agency ERP, the practical answer is usually "the Super Admin role doubles as DPO" — not perfect but common.

## Recommendation

**Phase 1 compliance stance (ship this with F3):**

1. **Legitimate basis: contract performance** for core data (biodata, documents, mahram relations, payment, visa) — this is data the agency literally cannot provide the service without. **Consent** for optional data uses: marketing broadcasts, retargeting pixels, extended retention beyond the legal minimum. Two separate consent toggles at registration.

2. **Consent capture:** a mandatory consent modal at first registration listing the core processing purposes + cross-border transfer to Saudi authorities, with a link to a privacy policy doc (to be drafted outside engineering). Additional opt-ins for marketing are off-by-default.

3. **Retention:** documents (`documents` table + GCS files) kept for **2 years post-departure**, then purged. Financial records (`journal_entries`, `invoices`) kept for **7 years per PSAK**. Alumni can opt in to extended document retention for return-pilgrimage workflows. The purge runs as a scheduled job with a dry-run report reviewed by ops before execution.

4. **DSR SLA: 14 days** for access/rectification, **30 days** for erasure (requires review because of the 7-year financial retention override). Build a DSR admin console in F1 that handles the common cases; complex erasure requests escalate to the DPO.

5. **DPO:** Super Admin role doubles as DPO. Add a `dpo_contact_email` field to system config for breach notifications. Document the breach-notification procedure in `docs/04-backend-conventions/` as a runbook.

Rationale: this is the minimum defensible stance. Going stricter (full consent-based, 1-year retention) would add friction and conflict with PSAK obligations. Going looser would invite regulatory risk. The split-retention model (documents short, financials long) matches how other Indonesian fintechs and travel platforms operate post-UU-PDP.

Note that this recommendation is engineering's best read on what's defensible; **it is not a substitute for legal review**. The privacy policy document itself must be drafted by a lawyer or privacy advisor, not by us.

Reversibility: all five decisions are code/config changes, not data migrations. Retention rules change the purge job's window; DSR SLA is a policy statement; DPO is a config field; legitimate basis only affects the consent modal copy.

## Answer

**Decided:** Adopt **Recommendation defaults** (contract performance for core processing; separate **marketing consent**; granular toggles off-by-default for marketing; **2-year** document retention post-departure with **7-year** financials; **14-day** DSR access/rectify / **30-day** erasure subject to legal retention; Super Admin `dpo_contact_email` as DPO contact for breach routing). **Privacy policy text** and **breach runbook** must still be **drafted/reviewed by qualified legal counsel** — engineering stance is not a substitute.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
