---
id: Q046
title: Travel-agency PPN rate + PKP status + e-Faktur integration
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: open
---

# Q046 — Travel-agency PPN rate + PKP status + e-Faktur

## Context

PRD line 1285 states the default PPN rate as **11%** — but Indonesian regulation **PMK-71/PMK.03/2022** sets PPN for travel agency services (biro perjalanan) at **1.1% of gross** (DPP Nilai Lain base), not 11% of the service-value DPP. This is a special treatment for travel-related packages where the agency acts as an aggregator of third-party services (hotels, flights, visa, etc.).

If we budget / collect 11% on jamaah packages, we'll **over-collect** and owe refunds at reconciliation; if we apply 1.1% correctly, the financial math differs materially. This is one of the most high-impact ambiguities in the F9 spec.

Separately, PPN mechanics depend on whether the agency is **PKP (Pengusaha Kena Pajak, VAT-registered)**. Small agencies below the 4.8B IDR/year revenue threshold may not be PKP and issue non-PPN invoices.

## The question

1. **Confirm: is the agency PKP?** — yes/no; determines whether PPN is charged at all.
2. **If PKP: confirm PPN rate on jamaah Umrah packages.** — PMK-71 1.1% on DPP Nilai Lain (travel agency services), or 11% on a different DPP split, or 11% on full?
3. **How are bundled add-ons (insurance, optional merchandise, Raudhah Shield) taxed?** — they may be outside PMK-71's travel-agency treatment; standard 11% may apply.
4. **e-Faktur integration** — how does UmrohOS generate tax invoices that feed e-Faktur? Direct API to DJP (requires e-Faktur Desktop app installation or e-Faktur Host-to-Host), or CSV export for manual upload?
5. **Input PPN (on vendor invoices)** — hotel, airline, and other vendors issuing e-Faktur tax invoices; does the agency claim input PPN credit, or is PMK-71's special treatment exclusive (no input credit)?
6. **Output PPN on foreign-currency pricing** — when jamaah pays in IDR for a package priced from USD/SAR costs, PPN is computed on which amount?

## Options considered

- **Option A — Apply PMK-71 1.1% correctly; no input PPN claims (PMK-71 exclusive treatment).** For any PKP agency, jamaah-package PPN is 1.1% of DPP Nilai Lain (typically 10% of gross for Nilai Lain basis). Bundled non-travel add-ons may be separate at 11%.
  - Pros: regulatory-correct; no tax overpayment.
  - Cons: requires per-line classification of "is this travel-agency-PMK-71 service or not."
- **Option B — Apply 11% per PRD; treat as conservative over-collection (refund via reconciliation).** Charge 11%, file at 11%, let DJP reconcile back to 1.1% on filing.
  - Pros: matches PRD literal text.
  - Cons: over-collects from jamaah; financial over-statement; not how PMK-71 works in practice — you either apply it or don't.
- **Option C — Consult tax advisor during MVP; ship with a config flag and sensible default; finalize before production.** Config: `ppn_mode: pmk71 | standard | none`, defaults to pmk71 for PKP agencies, none for non-PKP.
  - Pros: pragmatic; allows correct configuration per client.
  - Cons: pushes decision to deploy-time instead of spec-time.

## Recommendation

**Option A with Option C's config-flag pattern as the implementation mechanism — but pinned to PMK-71 as default.**

The agency-size question (PKP yes/no) is a stakeholder input; the regulatory correct rate is not. PMK-71 is the applicable regulation for travel-agency services; applying 11% is factually wrong for a PKP travel agency. Option B's "over-collect and refund later" is not how Indonesian tax filings work — the PPN filed at 11% becomes the liability, and reconciliation to 1.1% requires an amendment process that DJP may contest. Ship Option A.

Implementation via Option C's config flag is sensible because some clients may be non-PKP (rate = 0%), some may have different service mixes. Default `ppn_mode: pmk71` for PKP travel agencies; the implementation supports `standard` (11% full) and `none` for non-PKP clients.

Defaults to propose: PKP status captured in agency-settings; default assumption is PKP (most 10K+ jamaah/year agencies are). PPN mode = `pmk71` (1.1% on DPP Nilai Lain for bundled Umrah packages, 11% on non-travel add-ons). Input PPN — **not claimable** under PMK-71 Nilai Lain treatment per current regulation interpretation (confirm with tax advisor). Output PPN on foreign-currency pricing: computed on the IDR invoice amount at the transaction-date rate (PSAK 10). e-Faktur integration: **CSV export for DJP Online upload** in MVP (fastest to ship); direct e-Faktur API integration (Host-to-Host) later if volume justifies.

Reversibility: changing PPN modes later requires re-generating historical tax records — non-trivial; get this right before going live.

## Answer

TBD — awaiting stakeholder input. **This question needs external tax advisor / accountant confirmation** — engineering's read of PMK-71 should not be the sole basis. Deciders: external tax advisor / accountant (primary), finance director, agency owner (PKP registration status).
