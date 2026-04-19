---
id: F9
title: Finance & Accounting (PSAK-compliant)
status: written
last_updated: 2026-04-18
moscow_profile: 15 Must Have / 7 Should Have (highest Must Have count in the catalogue)
prd_sections:
  - "G. Finance & Accounting (lines 429–493)"
  - "Alur Logika 9.1–9.4 (lines 1227–1235) — one-line bullets only; F9 fills the procedural detail"
  - "Alur Logika 1.1 global finance variables (kurs / PPN / PPh 23) — lines 1273–1325"
  - "Sitemap Keuangan & Akuntansi (lines 813–829); Menu (lines 1093–1109)"
  - "Cross-refs: Auto-AP trigger line 389, Lunas-Trigger line 411, Commission dashboard lines 143–147, Refund paperwork line 341"
modules:
  - "#129 Penagihan Otomatis, #130 Integrasi Bank, #131 Buku Pembantu Piutang, #132 Kwitansi Digital, #133 Pembayaran Manual & Cek Muka"
  - "#134 Database Vendor, #135 Buku Pembantu Hutang, #136 Otorisasi Pembayaran, #137 Kas Kecil & Bon Sementara"
  - "#138 Pembiayaan Berbasis Proyek, #139 Laba Rugi Keberangkatan, #140 Anggaran vs Aktual"
  - "#141 Jurnal Otomatis (Single Input Double Entry), #142 Pengakuan Pendapatan, #143 Multi-Currency, #144 Manajemen Aset Tetap, #145 Integrasi Pajak"
  - "#146 Pencairan Komisi Agen, #147 Laporan Keuangan Real-Time, #148 Dashboard Arus Kas, #149 Peringatan Umur Hutang/Piutang, #150 Jejak Audit & Hak Akses"
depends_on: [F1, F4, F5, F6, F8, F10]
open_questions: []
---

# F9 — Finance & Accounting (PSAK-compliant)

## Purpose & personas

F9 is the **accounting truth layer** — every peso of value that moves through UmrohOS (jamaah payment, vendor invoice, commission, refund, FX gain, fixed-asset depreciation) lands here as a double-entry journal. It consumes events from F5 (payment), F8 (logistics/GRN), F10 (commission), F4 (booking lifecycle), and F6 (visa) to produce PSAK-compliant financial statements: Neraca, Laba Rugi, Arus Kas, Perubahan Ekuitas.

This is the **heaviest feature in the catalogue** — 22 modules (15 Must Have / High) covering AR, AP, journal engine, tax, FX, fixed assets, commission payout, reporting, and governance. PRD Section G is heavy on intent and light on procedural detail (Alur 9.1–9.4 are one-liners) — the spec here fills procedural detail; linked **Q001 / Q012 / Q017–Q018 / Q036–Q038 / Q042–Q053** are **answered** as of **2026-04-18** (see `docs/07-open-questions/`).

Primary personas:

- **Finance admin (Staff_Keuangan)** — posts manual entries, verifies manual payments, reconciles bank statements, runs monthly close.
- **Kasir / cashier** — records cash receipts + petty cash disbursements; per-branch data scope.
- **AP officer** — processes vendor payments after GRN; routes to approver per threshold ladder (Q050).
- **AP approver (Manager → Director)** — signs off vendor disbursements per Otorisasi Pembayaran (#136).
- **Finance director / CFO** — closes periods, signs off manual adjustments over threshold, reviews reports before distribution.
- **Agency owner / CEO** — consumes Dashboard Arus Kas (#148) + Neraca + Laba Rugi in real time.
- **External auditor (annual)** — audits the GL via audit trail (#150); never writes.
- **Tax officer / external accountant** — generates tax records for monthly PPN / PPh 21 / PPh 23 filings.

Consumer relationships:

- **F11 Dashboards** — executive KPIs read from F9's aggregation endpoints.
- **F5 (payment)** — emits `payment.received` → F9 records AR collection.
- **F8 (logistics)** — emits `logistics.grn_recorded` → F9 records AP liability (synchronous per Q038 default).
- **F10 (CRM/commission)** — emits `crm.commission_confirmed` → F9 accrues commission per Q045 timing.
- **F4 (booking)** — emits `booking.departure_completed` → F9 recognizes revenue per Q043 trigger.

## Sources

- PRD Section G in full (lines 429–493): Revenue & Receivables, Disbursement & Payables, Project-Based Accounting, Core Accounting & Tax, Reporting & Governance.
- PRD Alur Logika 9.1–9.4 (lines 1227–1235) — bullet-level; this spec fills the narrative.
- PRD Alur Logika 1.1 (lines 1273–1325) — global finance variables: kurs, PPN, PPh 23, FX non-retroactivity rule.
- PRD sitemap Keuangan & Akuntansi (lines 813–829) and menu structure (lines 1093–1109).
- PRD line 1417 — journal hard-delete prohibition (audit standard).
- PRD line 1343 + 1395–1409 — branch data-scope applied to finance queries.
- Module list: #129–#150 per `docs/Modul UmrohOS - MosCoW.csv` (22 modules).

## User workflows

### W1 — Auto-invoice issuance (module #129)

1. Trigger events: booking created (DP invoice), Tabungan enrollment (monthly-installment schedule invoice), booking-status-transition (final-settlement invoice).
2. `finance-svc` subscribes to F4/F5 events or pulls via scheduled cron (Phase 1 per 03-events.md note).
3. Invoice generation: picks the correct invoice template per event type, computes amount, attaches jamaah + booking + package context, assigns invoice number (`INV-YYYYMMDD-NNNN`).
4. Emits `finance.invoice_issued` → F5 (for VA issuance) + F10 (for WA notification).
5. AR sub-ledger (Buku Pembantu Piutang, #131) updates: credit jamaah_ar balance.
6. Dual-granularity grouping (PRD line 439): per-jamaah AND per-family (via group_code); reports roll up both ways.

### W2 — Payment receipt posting + digital receipt (modules #130, #131, #132)

1. `payment-svc` confirms VA payment via webhook; emits `payment.received` with `(payment_id, booking_id, jamaah_id, amount, paid_at, gateway, gateway_txn_id)`.
2. `finance-svc` creates a journal entry in the same transaction as AR update (per ADR 0006 in-process saga):
   - Dr **Bank (1.1.2.x)** — amount
   - Cr **Hutang Jamaah (2.1.x)** — amount
   - Revenue is **NOT recognized here** — per Q043 default, recognition waits for departure (PRD line 473).
3. AR sub-ledger (#131): decrements jamaah_ar balance; updates invoice status to `paid`.
4. Digital receipt (#132): `finance-svc` calls crm-svc to send WhatsApp receipt with PDF attachment (PRD line 441: "sedetik setelah" confirmation).
5. Audit log captures receipt issuance (who / when / amount / destination).

### W3 — Manual payment + admin verification (module #133)

1. Cashier at branch receives cash or manual bank transfer; opens `/erp/finance/penerimaan/manual`.
2. Records amount, jamaah + booking, payment method (cash / bank transfer / check), reference / slip number, optional receipt photo.
3. Submits for admin verification: status `pending_verification`. Cashier cannot self-verify (audit rule per #150).
4. Finance admin reviews + approves OR rejects (with reason). On approval, same journal entry as W2 fires (`Dr Bank/Kas, Cr Hutang Jamaah`).
5. Branch data-scope rule (PRD line 1343): cashier sees only own-branch; admin may be branch-scoped or central-scoped per role assignment.
6. Failed-verification: entry does not hit GL; jamaah AR unchanged; audit log shows attempt + rejection reason.

### W4 — Vendor master + AP creation (modules #134, #135)

1. Vendor onboarding — see **Q036** for the ownership boundary between finance-svc and logistics-svc. Per Q036 recommended split: logistics owns operational fields (contact, delivery, rating); finance owns financial fields (NPWP, PKP status, bank account, AP aging).
2. AP trigger: `logistics-svc` fires `RecordPayable(po_id, vendor_id, amount, grn_id)` via gRPC synchronously within the GRN transaction (per Q038 default).
3. `finance-svc` creates journal entry:
   - Dr **Persediaan (1.1.3.x)** — amount (or Beban if non-inventory service)
   - Cr **Hutang Usaha (2.1.x)** — amount
4. AP sub-ledger (#135, Buku Pembantu Hutang): credits vendor_ap balance; tracks payment history.
5. Tax withholding check — if vendor is service-type and has NPWP, PPh 23 accrued (W13); if no NPWP, 2× rate per Q047.

### W5 — AP disbursement with tiered e-approval (module #136)

1. AP officer opens `/erp/finance/pembayaran/utang` — sees open AP aged by due date.
2. Selects vendor(s) for payment; reviews AP summary (amount + PPh 23 witholding + net payable).
3. Creates payment batch; submits for approval per the Otorisasi Pembayaran ladder (**Q050** — default inferred: ≤10M Manager; 10–50M Director; >50M Director + CEO).
4. Approver reviews + approves OR rejects. Biometric / TOTP confirm for high-value tiers (Q050 default).
5. On approval: finance-svc executes bank transfer (via bank API if integrated, or manually with staff entry for MVP); journal entry:
   - Dr **Hutang Usaha** — gross amount
   - Cr **Bank** — net amount
   - Cr **PPh 23 Dipotong (2.1.y)** — witholding amount (if applicable, per Q047)
6. AP sub-ledger: clears the selected AP; vendor payment history logged.
7. Audit trail: approver_chain, approval_timestamps, two-factor evidence, disbursement_ref.

### W6 — Petty cash + field officer advance (module #137)

1. Petty cash fund seeded at each branch with monthly float (e.g. 5M IDR per branch). Managed by cashier.
2. Expense entry: cashier records disbursement (amount, category, recipient, photo of receipt). Journal:
   - Dr **Beban (6.x.x)** — amount (category-mapped account)
   - Cr **Kas Kecil (1.1.1.x)** — amount
3. Replenishment: when petty cash drops below threshold (e.g. 20%), cashier requests top-up; approved per Q050; top-up journal:
   - Dr **Kas Kecil** — top-up amount
   - Cr **Bank** — top-up amount
4. Field officer advance (bon sementara): similar flow, but disbursement lands in `Piutang Karyawan (1.1.7.x)` until the officer submits receipts; receipts clear the piutang to Beban.
5. Unsettled advances aging report via #149.

### W7 — Job-order costing per departure (modules #138, #139)

1. Each `package_departure_id` is an independent project / cost center (PRD lines 461–463). Projects allocate a unique code (`PROJ-<departure_code>`).
2. `journal_lines.job_order_id` is populated on every entry that touches a specific departure: jamaah revenue (at W10), direct costs (tickets, visa, LA, catering, bus, hotel).
3. Per-departure P&L (Laba Rugi Keberangkatan, #139):
   - Revenue: sum of recognized revenue per departure (W10)
   - Direct costs: tickets + visa + land-arrangement + bus + hotel + muthawwif honor
   - Indirect allocation (optional): marketing + ops overhead, pro-rated per pax
   - Gross margin: revenue − direct costs
4. Real-time view at `/erp/finance/project-based/laba-rugi/{departure_id}`.
5. Cost-center binding: F8 (logistics) PO lines, F6 visa fees, F10 commission all tag `job_order_id` via the departure they relate to.

### W8 — Budget vs Actual (module #140)

1. Budget source: each package's HPP breakdown at catalog-svc (established in F2) seeds per-departure budget by multiplying unit costs × pax forecast.
2. Actuals: real journal_lines with `job_order_id = <departure>`.
3. Variance view: budget − actual per category (tickets, hotel, visa, etc.); flags > X% variance.
4. Dashboard at `/erp/finance/project-based/bva` with drill-down from category → individual journal entries.

### W9 — Auto-journaling engine (module #141 — "Single Input, Double Entry")

1. Core principle: **operational transactions** (payments, GRNs, commissions, refunds, manifests) enter once in their owning service; finance-svc transforms them into PSAK-compliant journal entries without re-entry.
2. Entry map (illustrative, to be finalized with Q049 COA):
   | Source event | Debit | Credit |
   |---|---|---|
   | `payment.received` | Bank | Hutang Jamaah |
   | `booking.departure_completed` | Hutang Jamaah | Pendapatan Umroh |
   | `logistics.grn_recorded` | Persediaan / Beban | Hutang Usaha |
   | `payment.disbursed` (AP paid) | Hutang Usaha | Bank (net) + PPh 23 |
   | `crm.commission_confirmed` | Beban Komisi Agen | Hutang Komisi Agen |
   | `crm.commission_paid_out` | Hutang Komisi Agen | Bank + PPh 21 |
   | `payment.refund_completed` | Hutang Jamaah | Bank (see Q053 for pinalti split) |
   | `asset.depreciation_run` | Beban Depresiasi | Akumulasi Depresiasi |
3. Double-entry enforcement: DB constraint or trigger ensures `sum(debit) = sum(credit)` per journal_entry. Violation → transaction rollback.
4. Idempotency: `(source_kind, source_id)` unique; replay of an event does not double-post.
5. Authority: event-sourced entries are system-owned (source_kind != 'manual'). Manual entries (W10b below) follow a separate approval path.

### W10 — Revenue recognition at departure (module #142)

1. Trigger: `booking.departure_completed` event from booking-svc. **Q043** pins the exact event (wheels-up vs boarding scan vs departure-date auto-fire).
2. Per-booking revenue recognition — for each paid jamaah on the departure:
   - Dr **Hutang Jamaah** — amount paid to date
   - Cr **Pendapatan Umroh** — amount, tagged with `job_order_id = <departure>`
3. Deferred revenue contrast: Tabungan enrollments (see Q044) post to `Hutang Tabungan Jamaah` until they convert to a booking, then follow the same trigger.
4. Multi-installment bookings: total of all receipts recognized at departure (not per-installment); any arrears at departure remain as `Hutang Jamaah` outstanding.
5. Reversal on post-departure refund / claim: compensating journal (see Q053).

### W10b — Manual journal entries (adjustments)

1. Finance admin opens `/erp/finance/jurnal/manual`. Purpose: period-end accruals, error corrections, reclassifications, depreciation catch-up, FX revaluation.
2. Entry form: description + date + lines (account, debit, credit, job_order optional). Real-time validation: sum(debit) = sum(credit).
3. Approval: threshold-based (**Q051** default: any manual entry requires finance-director approval; > 50M IDR requires CFO); small adjustment entries by finance admin alone within limit.
4. Hard-delete prohibited (PRD line 1417). Correction is always via counter-entry (reversal), never by delete. Role system warns before granting Delete on Jurnal Akuntansi.
5. Audit log captures: creator, approver, timestamps, original-entry-reference (if a reversal).

### W11 — Multi-currency + FX (module #143)

1. Supported currencies: IDR (base), USD, SAR (per PRD line 1277).
2. FX rate source: **Q048** — default inferred BI middle rate + configurable markup; locked rate per season for HPP purposes (PRD line 1281).
3. Every journal_line stores `currency` + `fx_rate` snapshot. Reports aggregate in IDR using the line's snapshot.
4. FX gain/loss: on settlement of a foreign-currency AR/AP, difference between snapshot rate and settlement rate posts to `FX Gain (4.x)` or `FX Loss (6.x)`.
5. Non-retroactivity (PRD line 1313): rate changes do **not** restate invoices already DP'd / Lunas. Only new transactions use the new rate.
6. Month-end revaluation per PSAK 10: open foreign-currency AR / AP revalued to closing rate; FX gain/loss to P&L. Configurable (Q048 may defer or confirm this).

### W12 — Fixed asset register + depreciation (module #144)

1. Asset register: `/erp/finance/aset/daftar` — add asset, category (land / building / vehicle / equipment / IT), purchase cost, useful-life (months), depreciation method.
2. _(Inferred)_ Default method: straight-line per PSAK 16 — monthly depreciation = cost / useful_life_months. Other methods (declining balance, units-of-production) configurable per asset.
3. Monthly depreciation run (cron on 1st of each month, adjustable): for each active asset, post:
   - Dr **Beban Depresiasi (6.x.x)**
   - Cr **Akumulasi Depresiasi (1.2.x)**
4. Disposal / retirement: separate workflow; writes off remaining book value; captures disposal gain/loss.
5. Asset opname (annual): physical count vs register.

### W13 — Tax calculation (module #145)

1. Three Indonesian tax regimes handled:
   - **PPN** (VAT) — output on jamaah invoices (if PKP — **Q046**); input on vendor invoices with e-Faktur.
   - **PPh 21** (employee income tax) — witheld on salaries + agent commissions paid to individuals (**Q047**).
   - **PPh 23** (vendor services witholding) — 2% for services, 15% for rent (PRD line 1287–1289).
2. **Critical Q046**: PRD states PPN 11% on jamaah packages, but PMK-71/2022 sets travel-agency PPN at 1.1% of gross (DPP Nilai Lain). Default inferred: 1.1% on travel-agency bundled packages, 11% on non-travel add-ons (insurance, standalone merchandise). Confirmed with accountant via Q046.
3. Witholding on AP (Q047): vendor with NPWP → normal rate; without NPWP → 2× rate; bundled journal entry per W5.
4. Monthly tax records aggregate: `tax_records` rolls up by period + tax_kind. Export to e-Faktur format (CSV) for DJP upload.
5. Annual tax fileing support: SPT data exports for SPT PPh Badan, SPT PPh 21 Tahunan, etc. Read-only generators; the agency's external accountant files via DJP Online.

### W14 — Commission payout to agents (module #146)

1. Commission accrual per **Q045** (answered): accrue **Beban Komisi / Hutang Komisi on `paid_in_full`** (per-jamaah proration); **no partial-% accruals**; **clawback** until payout batch posted; **post-payout refund window 30d** nets wallet negative else **`Beban Komisi Forfeit`** bucket; override commissions same event.
2. Accrual journal:
   - Dr **Beban Komisi Agen (6.x.x)**, tagged with job_order
   - Cr **Hutang Komisi Agen (2.1.x)** — per-agent sub-ledger
3. Payout batch (#146): finance admin runs monthly; generates payout batch from confirmed commissions; picks bank-transfer file (CSV for upload to corporate banking).
4. Payout journal:
   - Dr **Hutang Komisi Agen** — gross
   - Cr **Bank** — net
   - Cr **PPh 21 Dipotong (2.1.y)** — witholding (per Q047 rate)
5. Payout proof (receipts, tax slip) pushed to agent via WA per crm-svc.
6. Clawback on jamaah refund: **Q045** — reverse accrual if refund before payout batch; after payout, apply **30d** post-payout window per Q045 (wallet negative or forfeit bucket).

### W15 — Financial statement reports (module #147)

1. On-demand generation at `/erp/finance/reports/*`:
   - **Neraca (Balance Sheet)** — as-of-date snapshot; assets = liabilities + equity; drill-down per COA node.
   - **Laba Rugi (P&L)** — period range; grouped per PSAK 1 (by function or nature — **Q042** to confirm).
   - **Arus Kas (Cash Flow)** — period range; direct or indirect method per Q042.
   - **Perubahan Ekuitas (Statement of Changes in Equity)** — period range.
2. Aggregation from `journal_lines` × COA tree; all amounts in IDR base (using per-line fx_rate snapshot).
3. Comparative periods (YoY, QoQ) — toggle.
4. Export formats: PDF (formatted), Excel (analyst drilldown), CSV.
5. Real-time principle (module #147, PRD line 487): no batch-reconciliation delay — report generated at click time reflects all committed journals.

### W16 — Cash flow dashboard (module #148)

1. Single-screen view: consolidated bank balances across all accounts, petty cash (per-branch + central), AR aging, AP aging, Tabungan float (jamaah liability), agent wallet balances, brankas fisik (physical cash vault).
2. Trend charts: 30-day rolling cash movement, weekly AR collection, weekly AP outflow.
3. Alerts: low-balance threshold per bank account, overdue AR, overdue AP.
4. Audience: agency owner / CFO daily check.

### W17 — AR / AP aging alerts (module #149)

1. Daily cron: for each open AR (jamaah) and AP (vendor), compute age = today − due_date.
2. Buckets: current / 1–30 / 31–60 / 61–90 / 90+ days.
3. AR alert: overdue jamaah payments trigger Collection workflow (cross-ref F5 / F10 WA reminders).
4. AP alert: approaching-due vendor payments surface in AP officer dashboard; overdue → elevated alert to finance-director.
5. Period-end aging snapshot captured as reportable artefact.

### W18 — Audit trail + anti-fraud access control (module #150)

1. Every finance action writes to `iam.audit_logs` (F1 `RecordAudit` gRPC): actor, action, entity_kind, entity_id, before/after snapshot, timestamp, IP, trace_id.
2. Hard rule: Delete on Jurnal Akuntansi is discouraged by role system — granting it triggers a warning (PRD line 1417). Corrections via counter-entry only.
3. Fraud alert triggers (PRD line 523): bulk manual-journal activity > threshold, repeated out-of-hours entries, large petty-cash withdrawals, journal edits just before period close.
4. Retention: **Q055 candidate** — default inferred 10 years per Indonesian tax/accounting law (UU KUP Pasal 28).
5. External auditor role (read-only): `/erp/finance/audit/*` endpoints with query-only access to GL and source-documents.

## Acceptance criteria

- **Double-entry integrity** — every journal entry balances (Σ debit = Σ credit); enforced at DB via trigger/constraint; never violated.
- **Idempotency on event-sourced entries** — `(source_kind, source_id)` uniquely identifies an entry; replay does not double-post.
- **Revenue recognition** — jamaah payments stay in Hutang Jamaah until the Q043 trigger event; no pre-recognition.
- **Auto-AP synchronous with GRN** — logistics GRN transaction rolls back on finance-svc failure (Q038 default).
- **Multi-currency snapshot** — every journal_line stores currency + fx_rate; non-retroactivity enforced.
- **FX gain/loss** — settlement-rate vs snapshot-rate difference posts automatically.
- **AP disbursement approval** — no disbursement without matching approval chain per Q050 ladder.
- **Manual journal approval** — above Q051 threshold requires finance-director / CFO co-sign.
- **Hard-delete prohibition** — DB role + UI both refuse deletion of journal entries.
- **Branch data-scope** — queries filter by branch_id for non-central roles (PRD line 1343).
- **Tax aggregation** — monthly tax_records rolls up PPN + PPh 21 + PPh 23 correctly; e-Faktur export valid.
- **Financial statements** — Neraca balances; Laba Rugi per period sums to retained-earnings change; Arus Kas reconciles to bank balance change.
- **Audit trail completeness** — every write has a matching audit row; no "ghost" entries.
- **Period close** — closed periods cannot be written to without re-open approval (Q051).

## Edge cases & error paths

- **Webhook late-arrival / replay** — payment.received idempotent on (gateway, gateway_txn_id); journal entry idempotent on source_id.
- **Auto-AP failure** (finance-svc down during GRN) — GRN rolls back per Q038; stock not incremented; operator retries.
- **Commission accrual on cancelled booking** — reverse-journal clears the accrual; Q045 refines timing rules.
- **Jamaah refund post-revenue-recognition** — compensating revenue reversal + refund journal; Q053 refines pinalti split.
- **FX rate source outage** — fall back to previous-day rate with flag; alert finance admin; Q048 refines.
- **Tabungan matured but no booking yet** — stays in Hutang Tabungan; periodic aging reconciliation; Q044 refines stale-balance handling.
- **Negative petty cash** — blocked by UI; cashier must request replenishment first.
- **Manual journal attempted against closed period** — blocked; Q051 refines re-open approval flow.
- **Vendor with no NPWP requiring PPh 23** — 2× rate per Q047 default; bundled in AP journal.
- **Multi-entity consolidation** (if agency has subsidiaries) — out of MVP scope; flagged for Phase 2.
- **FX on partial payment** — each receipt uses its settlement-date rate; the AR balance aggregates across multiple snapshot rates (unrealized revaluation monthly per PSAK 10).
- **Audit-trail gap** (app crashed mid-write) — reconciliation cron (same pattern as F4/F5 per ADR 0006) detects uncommitted journal ids; flags for manual resolution.
- **Tax-rate change mid-year** — new rate applies prospectively; historical entries unchanged; period-close captures the pre-change totals.
- **Pinalti accounting** — Q053 default inferred: pinalti retained from jamaah refund posts to `Pendapatan Lain-Lain (4.9.x)` (non-operating revenue), not reversed from main revenue stream.
- **Year-end closing entries** — revenue/expense accounts zero out to retained earnings on annual hard-close.

## Data & state implications

Owned by `finance-svc` (schema planned at `docs/03-services/08-finance-svc/02-data-model.md`). Key tables referenced / extended by this spec:

- `chart_of_accounts` — existing; add `is_system_seeded` flag + `coa_template_version` metadata per Q049.
- `journal_entries` — existing; source_kind extended: `manual | payment | logistics | crm | closing | depreciation | fx_revaluation | commission | tax | refund`.
- `journal_lines` — existing; add `job_order_id` (already planned), `tax_kind nullable`, `is_reversal boolean`, `reverses_line_id nullable`.
- `ar_balances`, `ap_balances` — materialized from journal_lines; per Q042/Q049 may become maintained tables for performance.
- `tax_records` — existing; add `efaktur_nomor nullable`, `vendor_npwp nullable`.
- `fx_rates` — existing; add `source` enum (`bi_api | manual | locked_hpp`).
- `job_order_costs` — existing; revisit to include budget + actual + variance columns.
- `fixed_assets` — new: `{ id, code, name, category, purchase_cost, purchase_date, useful_life_months, depreciation_method, accumulated_depreciation, book_value, status, disposed_at nullable }`.
- `depreciation_runs` — new: monthly run record `{ run_period, journal_entry_id, asset_count, total_depreciation }`.
- `periods` — new: `{ period (YYYY-MM), status enum (open/soft_closed/hard_closed), closed_at, closed_by }`; gates write attempts.
- `agent_commission_ledger` — new (or owned by F10 and mirrored): per-agent running balance for commissions.
- `tabungan_liabilities` — new: per-Tabungan-account deferred-revenue tracking (Q044).
- `manual_journal_approvals` — new: audit trail of who approved what threshold.
- `bank_accounts` — new: master table of bank accounts with `coa_account_id` mapping for posting.
- `petty_cash_accounts` — new: per-branch petty cash fund state.

New enums:

- `account_kind` — existing: `asset | liability | equity | revenue | expense`. Extend: `contra_asset` (Akumulasi Depresiasi) if needed.
- `source_kind` — see journal_entries above.
- `tax_kind` — existing: `pph21 | pph23 | ppn`. Extend: `pph4_2` (final) if Q046 introduces Final PPN treatment.
- `fx_source` — `bi_api | manual | locked_hpp`.
- `period_status` — `open | soft_closed | hard_closed`.
- `depreciation_method` — `straight_line | declining_balance | units_of_production`.
- `commission_accrual_stage` — per Q045: `pending | confirmed | paid | reversed`.

## API surface (high-level)

Full contracts in `docs/03-services/08-finance-svc/01-api.md`. Key surfaces confirmed here:

**REST (finance console):**
- `GET /v1/chart-of-accounts` (tree view) + `POST /v1/chart-of-accounts` (admin add)
- `GET /v1/journal-entries` (filterable: period, source_kind, job_order) + `POST /v1/journal-entries` (manual entry) + `GET /v1/journal-entries/{id}`
- `POST /v1/journal-entries/{id}/approve` (finance director action)
- `GET /v1/ar-balances` (aging) + `GET /v1/ap-balances` (aging)
- `POST /v1/payments/manual` (W3)
- `POST /v1/ap-disbursements` (W5: batch) + `POST /v1/ap-disbursements/{id}/approve`
- `GET /v1/fx-rates` + `POST /v1/fx-rates` (daily snapshot; source = bi_api / manual / locked_hpp)
- `GET /v1/fixed-assets` + `POST /v1/fixed-assets` + `POST /v1/fixed-assets/depreciation-run`
- `GET /v1/job-order-costs/{departure_id}` (W7 per-departure P&L + BvA)
- `GET /v1/reports/balance-sheet` + `/profit-loss` + `/cash-flow` + `/equity-changes`
- `GET /v1/tax-records?period=YYYY-MM&kind=ppn|pph21|pph23` + `POST /v1/tax-records/efaktur-export`
- `POST /v1/periods/{period}/close` + `POST /v1/periods/{period}/reopen` (CFO)
- `GET /v1/dashboards/cash-flow` (W16 real-time)

**gRPC (service-to-service):**
- `RecordJournalEntry(source_kind, source_id, lines[])` — called by payment-svc, logistics-svc, crm-svc, booking-svc to push event-based entries (per ADR 0006 in-process saga).
- `RecordPayable(po_id, vendor_id, amount, grn_id)` — called by logistics-svc W4 (Q038 synchronous).
- `GetExchangeRate(from, to, date)` — queried by any service needing FX conversion.
- `RecognizeRevenue(booking_id, departure_id)` — called by booking-svc when Q043 event fires.
- `GetARBalance(jamaah_id)` — callable from F5 / F10 for jamaah-facing views.

**Events emitted** (`03-events.md`):
- `finance.entry_recorded` — all committed journals; consumers: dashboards.
- `finance.period_closed` — monthly / annual close; consumers: reporting caches invalidation.
- Add: `finance.commission_paid_out` (trigger WA payout receipt), `finance.ap_disbursed` (trigger vendor remittance advice).

## Dependencies

- **F1 (IAM)** — RBAC for finance / kasir / AP-officer / director roles; branch-scope enforcement; audit logs.
- **F2 (catalog)** — HPP per package seeds BvA budgets; package + departure codes for job-order tagging.
- **F4 (booking)** — `departure_completed` event triggers revenue recognition (Q043); booking lifecycle informs Tabungan conversion (Q044).
- **F5 (payment)** — `payment.received`, `payment.refund_completed` events drive W2 + W18 refund accounting (Q053).
- **F6 (visa)** — visa fee as sunk cost on rejection (Q028 cross-ref); vendor invoices from visa providers feed W4.
- **F7 (operations)** — ops-side refund paperwork (W15 in F7) triggers refund journal (Q053).
- **F8 (logistics)** — Auto-AP from GRN (Q038); vendor master dual-binding (Q036); kit-assembly cost tags job_order.
- **F10 (CRM)** — commission events drive W14; agent wallet mirrors W14 accrual.
- **F11 (Dashboards)** — executive KPIs read via W15 + W16 endpoints.
- **External** — BI FX API (if used per Q048), Indonesian DJP e-Faktur (export-only for MVP), corporate banking (manual CSV upload for MVP, API integration later).

## Backend notes

- **Double-entry enforcement** is enforced at DB level via trigger on `journal_entries` checking `SUM(debit) = SUM(credit)` across all lines at commit time. sqlc queries insert lines inside a `WithTx` transaction; the trigger fires on COMMIT, rolling back the transaction if unbalanced.
- **Idempotency** on event-sourced entries via unique constraint `(source_kind, source_id)` in `journal_entries`. Replay safe.
- **Period close** implemented as a state flag on `periods` + a trigger that rejects inserts/updates to `journal_entries` where `entry_date ∈ closed_period`. Soft-close still allows adjustment with director approval; hard-close is final.
- **Depreciation run** is idempotent per `(asset_id, period)` — rerunning does not double-post.
- **FX revaluation** runs as a month-end job (cron or admin-triggered): revalues open foreign-currency AR/AP balances at closing rate; posts FX gain/loss journal; creates `journal_entries.source_kind = 'fx_revaluation'`.
- **Auto-journaling contracts** — source services push via `finance-svc.RecordJournalEntry(source_kind, source_id, lines[])` gRPC; lines must balance (validated client-side via shared lib + re-validated server-side). Failure = stored unposted with error; retry via reconciliation cron.
- **Job-order tagging** — every auto-journal rule knows how to resolve `job_order_id` from its source: payment → booking → departure; GRN → PO → (optional) linked departure; commission → booking → departure.
- **Reports** generated by aggregating `journal_lines` via sqlc queries. For large datasets, consider materialized views for Neraca + Laba Rugi with nightly refresh + on-demand refresh trigger.
- **Branch data-scope** applied via a filter clause in every query: `WHERE branch_id = $1 OR $branch_scope = 'all'`. Enforced in the service layer, not relied upon in the UI.
- **PSAK 72 revenue recognition** for Umrah contracts: point-in-time satisfaction at departure (not over-time), per Q043 default. PSAK 71 / ECL for Talangan: see Q052.
- **Tax calculation** — per-transaction helper that computes PPN + PPh based on transaction kind, vendor NPWP, customer PKP status. Output lines bundled into the same journal entry.
- **Immutable journal ledger** — no deletes ever. Role system denies Delete on journal_entries; counter-entry is the only correction path.

## Frontend notes

- **Finance console layout** — top-level modules: Penerimaan, Pembayaran, Jurnal, Laporan, Pengaturan.
- **Manual journal entry form** — live sum(debit) vs sum(credit) ticker; save disabled if unbalanced. Keyboard shortcuts for fast entry. Account picker typeahead + favorites.
- **AP disbursement batch screen** — multi-select open AP, approval routing visible, net-payable auto-computed post-witholding.
- **Report screens** — Neraca / Laba Rugi / Arus Kas / Equity. Comparative periods. Drill-down from category to journal entries to source transaction.
- **Cash flow dashboard (W16)** — single-screen consolidated view with trend charts; designed for executive glance.
- **Tax export** — one-click monthly PPN + PPh exports to e-Faktur-compatible CSV.
- **Audit viewer** — journal-entry page shows before/after diff on any edit, the approver chain, and the source event. Read-only auditor role removes edit affordances.
- **Period close wizard** — step-by-step: run FX revaluation → run depreciation → check unposted entries → confirm close. Soft-close + hard-close as distinct actions.
- **Mobile-friendly approval surface** — AP disbursement approval + manual-journal approval reachable on phone for directors on-the-go; biometric / TOTP second factor for > Q050 threshold.

## Open questions

None blocking — **Q001, Q012, Q017–Q018, Q036–Q038, Q042–Q053** answered **2026-04-18** (`docs/07-open-questions/`). Tax rates, COA seed, and PSAK framing in workflows above follow those answers; **PKP / e-Faktur** still needs **tax advisor** sign-off before production filing.

**Engineering defaults (cross-check during implementation):**

- Revenue recognition trigger = wheels-up / departure-date event, per PRD line 473 — Q043 formalizes.
- Straight-line depreciation default per PSAK 16 — deeper configurability per asset is additive.
- Audit-trail retention = 10 years per UU KUP Pasal 28 — tax/accounting law default.
- e-Faktur integration = export-to-CSV for DJP Online upload (MVP); direct DJP API later.
- Bank reconciliation = auto-match from VA statements + manual tick-off for non-VA transactions (MVP).
- Pinalti on refund posts to `Pendapatan Lain-Lain` (non-operating revenue) — Q053 refines.
- Fraud alert threshold = > 5 journal edits in one hour by same user, > 50M IDR petty cash, journal edits in the last 24h before period close — configurable.
- Multi-entity consolidation = **out of MVP scope**; single-entity only. Revisit if agency grows into holding structure.
