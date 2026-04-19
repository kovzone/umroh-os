---
id: Q047
title: PPh 21 scheme + PPh 23 witholding for no-NPWP vendors
asked_by: session 2026-04-17 F9 draft
asked_date: 2026-04-17
blocks: F9
status: answered
---

# Q047 — PPh 21 scheme + PPh 23 witholding for no-NPWP vendors

## Context

PRD line 1287–1289 states **PPh 23 default rate**: 2% (services) / 15% (rent). The PRD is silent on:

- **PPh 21 scheme** — Indonesian employee income tax has multiple computation models: monthly gross up, net calculation, **TER (Tarif Efektif Rata-rata)** introduced 2024 replacing older monthly schemes, per-article payments (Pasal 21 ayat 1b for pensiunan, etc.). Commission payments to agents-as-individuals are PPh 21 territory if the agent is treated as a temporary service provider.
- **PPh 23 witholding rate when vendor has no NPWP** — Indonesian regulation says vendors without NPWP face **2× the rate** (so 4% instead of 2%, or 30% instead of 15%). The PRD doesn't mention this doubling.
- **PPh 4(2) final** — some transactions (construction, rent, certain services) fall under PPh 4(2) final at different rates, outside the PPh 23 regime.

## The question

1. **PPh 21 computation scheme for employees** — TER monthly rates, or an older per-Pasal computation? Which employee types are in scope (permanent / contract / temporary)?
2. **PPh 21 on agent commission** — is an agent treated as a permanent employee (monthly PPh 21 on commission), a temporary service provider (PPh 21 single-payment Pasal 21 ayat 1), or a PPh 23 recipient (if agent is a corporate entity)?
3. **PPh 23 with no-NPWP vendor** — confirm 2× rate (4% / 30%)?
4. **PPh 4(2) final scope** — which vendor categories fall under this (rent is already cited; any others in Umrah operational contracts)?
5. **Witholding slip generation** — does UmrohOS generate the Bukti Potong PPh for each witholding (required under e-Bukti Potong regulation), or does the external accountant handle?
6. **Monthly reporting** — F9 exports PPh 21 and PPh 23 data for the agency's monthly SPT Masa filings; format?

## Options considered

- **Option A — PPh 21 TER scheme (2024+); PPh 23 with automatic 2× for no-NPWP vendors; PPh 4(2) for rent; Bukti Potong generation in-app.** Full compliance with current regulations.
  - Pros: regulatory-correct; reduces external-accountant dependency.
  - Cons: substantial implementation (TER rate tables, Bukti Potong PDF generation, e-Bukti Potong upload format).
- **Option B — Minimal F9 witholding engine; external accountant handles PPh calculations.** F9 tags transactions with "subject to PPh N" flag; external accountant re-computes and files. F9 only posts the net payment entry.
  - Pros: faster MVP; defers tax complexity.
  - Cons: harder to reconcile; accountant has to pull data out to re-compute; risk of under-witholding from the vendor/employee perspective.
- **Option C — PPh 23 + 2× for no-NPWP handled in-app; PPh 21 deferred to external payroll system.** Commissions to corporate agents withhold PPh 23; individual-agent commissions exported to the payroll system for PPh 21 handling.
  - Pros: splits complexity reasonably; matches how most agencies handle payroll (separate system).
  - Cons: dual-system setup; commission reconciliation across UmrohOS and payroll.

## Recommendation

**Option C — PPh 23 + 2× no-NPWP in-app; PPh 21 for employees exported to payroll; PPh 21 for individual-agent commission handled in-app at Pasal 21 ayat 1 single-payment rate.**

Option A is the most complete answer but blows up the MVP scope with TER rate tables + e-Bukti Potong integration. Option B pushes too much back to the accountant — the vendor AP flow needs to know witholding to compute net payable correctly, and agents receiving commission need their Bukti Potong. Option C draws a clean line: UmrohOS handles vendor-side (PPh 23) and agent-commission-side (PPh 21 Pasal 21) — the two witholdings that sit inside AP + commission payout flows natively. Employee payroll PPh 21 is a separate system (most agencies use Gadjian, Talenta, or similar); UmrohOS doesn't duplicate payroll.

Defaults to propose: PPh 23 on vendor services — 2% with NPWP, 4% without (auto-doubled based on vendor master's npwp field). PPh 4(2) at 10% for rent (if the agency has rent-type AP entries). PPh 21 on individual-agent commission — Pasal 21 ayat 1 at 5% (TER for occasional PPh 21 if under PTKP; otherwise brackets). Corporate-agent commission: PPh 23 at 2% (same as other services). Bukti Potong generation in-app as PDF (basic — not e-Bukti Potong XML) for MVP; external accountant uploads to DJP Online. Monthly PPh 21 + PPh 23 export: CSV per tax_kind per period for SPT Masa filing.

Reversibility: upgrading to full TER monthly + e-Bukti Potong later is additive. Switching PPh 21 from in-app to external payroll is a configuration change.

## Answer

**Decided:** **Option C** — **PPh 23** in-app (**2% / 4% no-NPWP**) + **PPh 4(2) rent 10%** where applicable; **employee payroll PPh 21** stays **external payroll** export; **individual agent commission PPh 21** as **single-payment** path in-app **until payroll integration exists**; **Bukti Potong PDF MVP** (not full e-Bukti Potong XML). **TER tables:** add when payroll module ships. **Rates:** require **tax advisor confirmation** annually.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
