# Open Questions — How to Use

This directory holds product questions that came up during spec-writing or implementation and couldn't be answered from the PRD or common practice alone. The markdown file is the **only place the answer lives**. Discussion with stakeholders happens offline (meetings, Notion, email, messaging — whatever works) and the reviewer writes the decision back into the Q-file's Answer section. AI sessions then pick up the answered Q-file and update any feature specs that depended on it.

## File convention

One question per file. Filename: `QNNN-short-slug.md` (zero-padded, e.g. `Q001-booking-cancel-refund-timing.md`). Numbers are sequential and never reused.

## Question template

```markdown
---
id: QNNN
title: <one-line question>
asked_by: <session name or reviewer>
asked_date: YYYY-MM-DD
blocks: <feature id(s) this blocks, e.g. F4, F5>
status: open | answered | deferred
---

# QNNN — <title>

## Context

What is being designed, and why the question arose. Reference the feature spec
and PRD section that triggered it.

## The question

A single, specific, answerable question. Not "how should we handle refunds" —
"should partial refunds be allowed when only one jamaah cancels from a
multi-pilgrim booking, or do we cancel the entire booking?"

## Options considered

- **Option A** — one-line description + pros/cons
- **Option B** — ...

## Recommendation

Which option the AI session actively recommends, and **why** — the reasoning
that makes this a genuine position, not just a fallback. Stakeholders engage
with reasoning, not defaults. This is the pick that will be inferred into the
spec (marked `_(Inferred)_`) if no answer comes back in time, but the reviewer
should read it as an opinion to agree, disagree, or override.

Format: one paragraph. Lead with the option name. Follow with the 1–3
tradeoffs that made this the pick over the alternatives. End with the
reversibility note — how cheap or expensive it would be to change later.

## Answer (filled in after the meeting)

TBD → write here once decided, with the date and who decided.
Then update the relevant feature spec and change status to `answered`.
```

## When to use this vs inferring

Use this directory when:
- The answer changes **user-visible behavior** (refund rules, commission calculation edge cases, visa rejection flow)
- The answer requires **regulatory or compliance input** (PSAK rules, Kemenag requirements)
- The answer requires **commercial judgment** (pricing rules, tier thresholds, discount approval limits)
- Two reasonable interpretations of the PRD exist

Infer (and mark `_(Inferred)_`) when:
- The answer is a well-known industry convention (HTTP status codes, JWT expiry windows, OCR field mapping)
- The answer is an implementation detail that doesn't change user behavior (which library, which port, which index)
- The answer has one obvious sensible default and reversing it later is cheap

## Recording an answer

Discussion with stakeholders happens **outside this repo** — meetings, Notion, email, a WhatsApp thread, whatever the reviewer chooses. Only the decision matters for the record; the route to it doesn't.

When an answer is in:

1. **Reviewer** opens `QNNN.md` and writes into the `## Answer` section:
   - The decided option (e.g. "Option B")
   - Date decided (`YYYY-MM-DD`)
   - Who decided (names or roles — e.g. "Religious advisor (Ustadz X) + agency owner")
   - Any amendments to the Recommendation (if stakeholders picked an option but adjusted a parameter — e.g. "Option A accepted, but threshold changed from 10 packages to 20")
2. **Reviewer** flips frontmatter `status: open → answered`.
3. **Reviewer** tells the next AI session, e.g. *"Q001 is answered."*
4. **AI session** reads the Q-file and:
   - Replaces every `TBD — see Q-NN` marker in the blocked feature specs with the concrete behavior.
   - For `_(Inferred)_` lines that the decision explicitly confirmed, removes the marker (it's now a decided fact).
   - For `_(Inferred)_` lines that the decision overrode, replaces them with the new behavior.
   - Updates the feature spec's `last_updated` date.
   - If this answer unblocks a feature spec entirely (all its questions are now answered), flips that spec's frontmatter `status: draft → written`.
5. The AI never edits the Q-file's Answer content — the reviewer owns that. AI only propagates the decision downstream into feature specs.

The question file stays in the dir as historical record — never delete answered questions. Even once every Q on a feature is answered, the Q-files remain as the audit trail of how the product decisions were made.

## Moving a question to "deferred"

Sometimes a question can't be answered yet (waiting on external input that isn't coming soon, or the feature slice has been pushed out). Set `status: deferred`, note the reason in the Answer section, and revisit when the blocking dependency resolves.

## Index

| ID | Title | Blocks | Status |
|---|---|---|---|
| [Q001](Q001-operating-currency-and-fx.md) | Operating currency, FX handling modes, HPP formula | F2, F5, F9 | open |
| [Q002](Q002-price-change-approval-and-audit.md) | Price-change approval thresholds and audit policy | F2 | open |
| [Q003](Q003-multi-language-catalog.md) | Multi-language catalog content scope | F2 | open |
| [Q004](Q004-cancellation-seat-return.md) | Cancellation → seat return ownership | F2, F4, F5 | open |
| [Q005](Q005-mahram-rules.md) | Mahram qualifying relations, age threshold, same-departure rule | F3, F4, F6 | open |
| [Q006](Q006-minimum-docs-for-booking.md) | Minimum documents required to submit a booking | F3, F4 | open |
| [Q007](Q007-name-mismatch.md) | KTP ↔ passport name mismatch handling | F3, F6 | open |
| [Q008](Q008-uu-pdp-compliance.md) | UU PDP compliance — consent, retention, DSR, DPO | F3, F1 | open |
| [Q010](Q010-va-ttl-and-draft-expiry.md) | VA TTL and draft booking expiry | F4, F5 | open |
| [Q011](Q011-dp-and-installment-rules.md) | Minimum DP %, max installments, cadence | F5 | open |
| [Q012](Q012-refund-penalty-matrix.md) | Refund penalty policy matrix (per package, per timing) | F5, F4 | open |
| [Q013](Q013-gateway-selection-and-fallback.md) | Dual-gateway selection and fallback rule | F5 | open |
| [Q014](Q014-partial-cancellation-and-reopen.md) | Partial cancellation + reopening cancelled bookings | F4, F5 | open |
| [Q015](Q015-smart-grouping-trigger-and-override.md) | Smart Grouping trigger timing + override authority | F4, F7 | open |
| [Q016](Q016-booking-for-minors.md) | Booking on behalf of a minor (no KTP) | F4, F3 | open |
| [Q017](Q017-paket-tabungan-interaction.md) | Paket Tabungan interaction with booking state | F4, F5, F9 | open |
| [Q018](Q018-talangan-bridging-finance.md) | Talangan (bridging finance) process and accounting | F5, F9 | open |
| [Q019](Q019-abandoned-checkout-attribution.md) | Abandoned checkout attribution (commission routing) | F4, F10 | open |
| [Q020](Q020-manifest-format.md) | Manifest format per airline / regulator | F7, F6 | open |
| [Q021](Q021-luggage-qr-and-all-system.md) | Luggage Tag QR scheme + ALL System protocol | F7 | open |
| [Q022](Q022-tour-leader-vs-muthawwif-authority.md) | Tour leader vs muthawwif override authority boundary | F7 | open |
| [Q023](Q023-zamzam-distribution-policy.md) | Zamzam distribution quota and policy | F7 | open |
| [Q024](Q024-incident-escalation.md) | Incident / issue report escalation matrix | F7 | open |
| [Q025](Q025-vulnerable-care-fields.md) | Vulnerable Care manifest fields and sensitivity | F7, F3 | open |
| [Q026](Q026-mofa-sajil-api-access.md) | MOFA / Sajil API access — credentials, sandbox, contracts | F6 | open |
| [Q027](Q027-visa-provider-selection.md) | Visa provider selection per package kind | F6 | open |
| [Q028](Q028-visa-rejection-handling.md) | Visa rejection handling policy | F6, F5, F4 | open |
| [Q029](Q029-passport-custody.md) | Physical passport chain-of-custody process | F6 | open |
| [Q030](Q030-raudhah-shield-cadence-and-action.md) | Raudhah Shield polling cadence + alert action | F6 | open |
| [Q031](Q031-tasreh-issuance-authority.md) | Tasreh issuance authority | F6 | open |
| [Q032](Q032-pr-po-approval-ladder.md) | PR / PO approval threshold ladder | F8 | open |
| [Q033](Q033-courier-integration-policy.md) | Courier integration policy — single vs multi, fallback, routing | F8 | open |
| [Q034](Q034-kit-composition-ownership.md) | Kit composition ownership — catalog-svc vs logistics-svc | F8, F2 | open |
| [Q035](Q035-post-ship-loss-returns-protocol.md) | Post-ship loss / damage / returns-from-trip protocol | F8 | open |
| [Q036](Q036-vendor-master-ownership.md) | Vendor master ownership + onboarding & rating | F8, F9 | open |
| [Q037](Q037-sku-barcode-vs-luggage-qr.md) | SKU barcode vs F7 luggage-tag QR coexistence on shipped kits | F8, F7 | open |
| [Q038](Q038-auto-ap-cadence-and-valuation.md) | Auto-AP posting cadence + PSAK inventory valuation | F8, F9 | open |
| [Q039](Q039-saudi-side-warehouse-scope.md) | Saudi-side warehouse scope | F8 | open |
| [Q040](Q040-stock-availability-policy.md) | Stock availability policy — partial shipments + reorder-point math | F8 | open |
| [Q041](Q041-self-pickup-qr-security.md) | Self-pickup QR security model | F8 | open |
| [Q042](Q042-psak-scope.md) | PSAK scope — full PSAK vs ETAP + which standards apply | F9 | open |
| [Q043](Q043-revenue-recognition-mechanics.md) | Revenue recognition mechanics — "terbang" event + refund reversal | F9, F4, F5 | open |
| [Q044](Q044-tabungan-deferred-revenue.md) | Multi-year Tabungan deferred-revenue accounting | F9, F4, F5 | open |
| [Q045](Q045-commission-accrual-timing.md) | Commission accrual timing — booking / paid / departure / payout | F9, F10 | open |
| [Q046](Q046-travel-agency-ppn-rate.md) | Travel-agency PPN rate (PMK-71) + PKP status + e-Faktur | F9 | open |
| [Q047](Q047-pph-scheme-and-witholding.md) | PPh 21 scheme + PPh 23 witholding for no-NPWP vendors | F9 | open |
| [Q048](Q048-fx-policy.md) | FX policy — rate source, transaction vs settlement date, revaluation | F9 | open |
| [Q049](Q049-coa-template-seed.md) | Chart of Accounts template seed | F9 | open |
| [Q050](Q050-ap-disbursement-approval-ladder.md) | AP disbursement approval threshold ladder | F9 | open |
| [Q051](Q051-period-close-procedure.md) | Period close procedure + re-open authority | F9 | open |
| [Q052](Q052-talangan-accounting.md) | Talangan accounting — loan receivable (PSAK 71) vs booking receivable | F9 | open |
| [Q053](Q053-refund-pinalti-accounting.md) | Refund & pinalti accounting entries | F9, F5, F4 | open |
| [Q054](Q054-agent-tier-taxonomy.md) | Agent tier taxonomy + qualification thresholds + demotion rules | F10 | open |
| [Q055](Q055-commission-percentage-table.md) | Commission % table (per level × per product) | F10, F9 | open |
| [Q056](Q056-overriding-formula.md) | Overriding commission formula + hierarchy depth | F10 | open |
| [Q057](Q057-utm-attribution-model.md) | UTM attribution model — window, first vs last click | F10, F4 | open |
| [Q058](Q058-alumni-referral-reward.md) | Alumni referral reward economics | F10 | open |
| [Q059](Q059-ziswaf-scope.md) | ZISWAF scope — donation platform vs pass-through | F10 | open |
| [Q060](Q060-wa-broadcast-limits.md) | WhatsApp broadcast rate limits + quality-score handling | F10, F7 | open |
| [Q061](Q061-agent-kyc-strictness.md) | Agent KYC strictness + activation thresholds | F10 | open |
| [Q062](Q062-replica-site-white-label.md) | Replica-site white-label scope | F10 | open |
| [Q063](Q063-testimoni-moderation.md) | Testimoni moderation policy | F10, F7 | open |
| [Q064](Q064-lead-ownership-transfer.md) | Lead ownership transfer between agents / CS | F10 | open |
| [Q065](Q065-ads-api-integration-depth.md) | Ads API integration depth (Meta / Google) | F10 | open |
| [Q066](Q066-dashboard-aggregation-architecture.md) | Dashboard aggregation architecture — service `/metrics` vs CQRS vs OLAP | F11 | open |
| [Q067](Q067-refresh-cadence-tiers.md) | Dashboard refresh cadence tiers (streaming / polling / on-demand) | F11 | open |
| [Q068](Q068-alert-threshold-ownership.md) | Alert threshold ownership + default values | F11 | open |
| [Q069](Q069-drill-down-depth.md) | Drill-down depth — widget to source transaction | F11 | open |
| [Q070](Q070-historical-retention-window.md) | Historical data retention window for dashboards | F11 | open |
| [Q071](Q071-multi-branch-consolidation.md) | Multi-branch consolidation rule for central visibility | F11 | open |
| [Q072](Q072-dashboard-export-policy.md) | Dashboard export policy — formats and permissions | F11 | open |
| [Q073](Q073-custom-dashboard-building.md) | Custom dashboard building vs fixed catalog | F11 | open |
| [Q074](Q074-field-radar-transport.md) | Field radar transport + GPS source for bus tracking | F11, F7 | open |
| [Q075](Q075-executive-landing-widget-composition.md) | Executive landing widget composition (top 8–12 KPIs) | F11 | open |
| [Q076](Q076-daily-app-form-factor.md) | Daily App form factor (native vs PWA vs responsive web) | F12 | open |
| [Q077](Q077-community-moderation-policy.md) | Community moderation policy + posting authority | F12 | open |
| [Q078](Q078-fatwa-desk-authority.md) | Fatwa Desk ustadz authority + response SLA + liability | F12 | open |
| [Q079](Q079-content-sourcing-licensing.md) | Manasik + Quran + daily-content sourcing & licensing | F12 | open |
| [Q080](Q080-third-party-data-sources.md) | Third-party data source selection (prayer times API, qibla, Quran) | F12 | open |
| [Q081](Q081-mvp-scope-carve.md) | F12 MVP scope carve — what ships vs what defers to Phase 2 | F12 | open |
