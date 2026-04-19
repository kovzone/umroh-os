---
id: Q042
title: PSAK scope — full PSAK vs ETAP + which standards apply
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: answered
---

# Q042 — PSAK scope

## Context

PRD Section G (line 429) and the sitemap (line 813) both cite *Keuangan & Akuntansi (Sesuai PSAK)* — but name no specific PSAK standards, and don't distinguish between **full PSAK** (for public companies and larger private entities) and **PSAK untuk Entitas Tanpa Akuntabilitas Publik (ETAP)** (simpler standard for SMEs). This is the most foundational accounting decision in F9 — it sets the grammar for every other question.

Umrah agencies vary widely in size: a boutique operator handling 500 jamaah/year can comfortably apply ETAP; an agency pushing 10,000+ jamaah/year with bank financing or any audit obligation typically applies full PSAK. The difference cascades into disclosure depth, revaluation requirements, impairment testing, deferred tax accounting (PSAK 46), and many presentation decisions.

## The question

1. **Full PSAK or ETAP?** — single most important choice.
2. **Which specific PSAK standards govern key transactions?** Candidates:
   - PSAK 1 — presentation (mandatory).
   - PSAK 2 — cash flow (direct or indirect method?).
   - PSAK 10 — FX (applicable to our USD / SAR transactions).
   - PSAK 14 — inventory (perlengkapan at F8).
   - PSAK 16 — fixed assets + depreciation.
   - PSAK 23 (older) or PSAK 72 (newer) for revenue from contracts — PSAK 72 is mandatory since 2020.
   - PSAK 46 — income tax accounting (current + deferred).
   - PSAK 57 — provisions (refund obligations, commission accruals).
   - PSAK 71 — financial instruments (agent wallet, Talangan loan receivable, ECL).
   - PSAK 73 — leases (offices, warehouses).
3. **Audit posture** — is the agency audited annually? By whom? This affects compliance rigor.
4. **Entity structure** — single entity or consolidation? (Consolidated entities = PSAK 65.)

## Options considered

- **Option A — Full PSAK across all areas.** Apply full standards including PSAK 72, 71, 46, disclosure requirements. Plan for audit-grade reporting from day one.
  - Pros: audit-ready; matches large-agency trajectory; no painful migration later if company grows.
  - Cons: significantly more implementation effort (ECL models, deferred tax, revaluation); may be over-kill for current scale.
- **Option B — ETAP with explicit upgrade path.** Apply PSAK ETAP for MVP; document what full-PSAK features are deferred; migration plan when the agency triggers scale or audit requirement.
  - Pros: faster MVP; matches current operational scale of most clients.
  - Cons: migration later is real work; some features (deferred revenue for multi-year Tabungan) harder to add retroactively.
- **Option C — Hybrid: ETAP presentation with PSAK 72 revenue + PSAK 71 financial instruments.** Take the two standards that affect correctness (not just disclosure) from full PSAK; use ETAP for the rest.
  - Pros: right-sized for scale; addresses the two areas where ETAP's simplification hurts (revenue timing + financial instruments).
  - Cons: non-standard posture; auditor may flag.

## Recommendation

**Option A — full PSAK.**

The reasoning is asymmetry: applying full PSAK on an ETAP-eligible agency produces slightly more work but passes every audit; applying ETAP on a full-PSAK-required agency fails audit and requires emergency remediation. Umrah agencies routinely face audits — bank financing review, Kemenag licensing review, occasional tax audit. The penalty for being under-compliant is severe; the penalty for being over-compliant is "extra detail in footnotes." Target full PSAK.

Specifically prioritize: PSAK 1 (presentation), PSAK 2 (cash flow — direct method preferred for jamaah-facing cash visibility), PSAK 10 (FX), PSAK 16 (fixed assets), **PSAK 72** (revenue from contracts — critical for Umrah performance obligation timing; see Q043 + Q044), PSAK 57 (refund + commission provisions), **PSAK 71** (financial instruments — critical if Talangan is a loan; see Q052), PSAK 46 (deferred tax — if full PSAK and auditable). Defer PSAK 73 leases until the agency signs a leased office/warehouse; not relevant until then.

Reversibility: starting with full PSAK and scaling down to ETAP is easy (drop disclosures, simplify ECL). Starting with ETAP and scaling up is hard (retrofit lot-level data for PSAK 71 ECL, restate prior periods for PSAK 72). Get this right up front.

## Answer

**Decided (engineering default):** Treat **`accounting_profile` = `etap`** as the **baseline implementation scope** for private SME tenants **unless** the tenant is **statutorily / bank-audited as full PSAK** — then **`full_psak` profile** flips deeper disclosures (PSAK 46 deferred until profile on). **Always implement revenue + FX + inventory + provisions to audit-grade mechanics regardless of profile label:** **PSAK 1, 2, 10, 14, 16, 57, 72** (+ **71** when Talangan in-house exists). **Rationale:** avoids over-building every PSAK disclosure for tiny agencies while not compromising revenue/FX correctness. **External auditor may require `full_psak`** — that is a **config + reporting pack** change, not a surprise rewrite if journals were honest.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
