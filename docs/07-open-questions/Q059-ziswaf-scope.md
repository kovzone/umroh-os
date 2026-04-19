---
id: Q059
title: ZISWAF scope — donation platform vs pass-through
asked_by: session 2026-04-17 F10 draft
asked_date: 2026-04-17
blocks: F10
status: answered
---

# Q059 — ZISWAF scope

## Context

PRD modules #200, #201, #202 cover ZISWAF (Zakat, Infaq, Shadaqoh, Waqaf) — three Islamic donation / savings types:

- **Tabungan Niat Kembali (#200)** — visual auto-debit savings targeted at future Umrah.
- **Kalkulator Zakat (#201)** — tool to compute Zakat Maal obligation (per gold-nisab).
- **Sedekah & Infaq Pagi (#202)** — morning donation via LAZ partner.

The difference between "donation platform" and "pass-through" is enormous in practice:

- **Donation platform**: UmrohOS holds donor funds, issues donor receipts, reconciles with LAZ beneficiaries, files Zakat-related tax-deduction certifications. Requires LAZ license or partnership-with-custody arrangement. Regulatory-heavy.
- **Pass-through**: UmrohOS tracks donor intent and links to LAZ partner's own payment page. Donor pays LAZ directly; UmrohOS just shows the click / completion confirmation. Regulatory-light.

PRD doesn't specify which model.

## The question

1. **Custody model** — does UmrohOS hold donation funds, or pass through to a partner LAZ (Lembaga Amil Zakat)?
2. **LAZ partnership** — if pass-through, which LAZ partner(s)? BAZNAS (government), Dompet Dhuafa (nonprofit), Rumah Zakat, LAZISMU? Multiple partners or single-integration?
3. **Receipt issuance** — who issues the donor tax-deductible receipt (Bukti Pembayaran Zakat that can reduce PPh obligation per tax rules)?
4. **Custody-holding implications** — if UmrohOS holds funds, licensing (must be a LAZ or work under LAZ umbrella), reporting to government, BAZNAS registration, etc.
5. **Tabungan Niat Kembali specifics** — is this an actual third-party savings product (partner bank), or just an auto-debit reminder that routes to the agency's main booking-DP account?
6. **Sedekah & Infaq (#202) — morning donation integration** — is this a daily push notification, a curated widget, or something deeper?
7. **Religious endorsement** — given the religious nature of ZISWAF, is an ustadz review required before feature launches?

## Options considered

- **Option A — Full donation platform with LAZ partner under custody.** UmrohOS holds funds, issues receipts, reconciles. Requires formal partnership with one LAZ (e.g. Dompet Dhuafa).
  - Pros: self-contained UX; donor receipts native; potential cross-sell with Umrah packages.
  - Cons: heavy compliance; LAZ licensing; regulatory scrutiny; non-trivial implementation.
- **Option B — Pass-through with click + completion tracking.** UmrohOS shows LAZ partner page in iframe / webview / deep-link; donor pays LAZ directly; UmrohOS gets callback/webhook for confirmation.
  - Pros: regulatory-light; partner handles compliance; faster to ship.
  - Cons: weaker brand / flow; dependent on partner's UX quality.
- **Option C — Defer ZISWAF from MVP; revisit in Phase 2.** The three modules are all Could Have priority; agency owner may agree to defer.
  - Pros: removes all complexity; focuses MVP energy elsewhere.
  - Cons: PRD says it's a product; deferring risks missing a differentiator.

## Recommendation

**Option B — pass-through with click + completion tracking; single LAZ partner (TBD by agency); deferrable per Option C if agency confirms it's not MVP-critical.**

Option A's self-custody carries LAZ licensing complexity that's disproportionate to the PRD priority (all three modules are Could Have / Low). Option B matches the Could Have / Low priority — small implementation, regulatory risk managed by partner LAZ. Option C is the honest deferral if the agency confirms ZISWAF isn't actually a launch-priority feature; many agencies treat it as a Phase 2 "nice-to-have."

Defaults to propose: **custody model** — pass-through only. No UmrohOS custody of donation funds. **LAZ partner** — single-integration initially (agency picks partner — e.g. Dompet Dhuafa, LAZISMU, or the agency's existing ZIS partner). **Receipt** — issued by LAZ partner; UmrohOS relays the receipt URL to donor. **Tabungan Niat Kembali** — reminder + auto-debit + visual tracker; funds route to the agency's own Tabungan booking product (cross-ref F4 Q017 for Tabungan accounting); no external savings product. **Sedekah & Infaq (#202)** — daily-subuh push notification with one-tap link to LAZ partner's subuh-sedekah page. **Religious endorsement** — recommended: one-time ustadz review of the workflows + copy before launch; captured as a required sign-off step. **Implementation footprint**: small — three UI surfaces (calculator, donation link button, auto-debit reminder) + intent tracking + webhook-completion record. No finance journal entries (intent-only, no custody).

Reversibility: upgrading from pass-through to custody later requires LAZ licensing + funds-handling infra; non-trivial but feasible. Defer-vs-ship toggle is a feature-flag (Option C path remains open).

## Answer

**Decided:** **Option B pass-through** — **no in-app custody**; **single LAZ partner** configurable; **receipts = partner**; **Tabungan Niat** = agency Tabungan rails (**Q017**); **daily subuh deeplink**; **ustadz one-time copy review** before enable. **Feature flag** can turn module off (**Option C**) per tenant.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
