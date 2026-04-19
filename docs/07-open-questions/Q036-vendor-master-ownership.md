---
id: Q036
title: Vendor master ownership + onboarding & rating
asked_by: session 2026-04-17 F8 draft
asked_date: 2026-04-17
blocks: F8, F9
status: answered
---

# Q036 — Vendor master ownership + onboarding & rating

## Context

PRD line 449 lists *Database Vendor* under Section G (Finance). PRD module #112 *Otomasi Vendor* (line 379) in Section F says POs are dispatched to vendors via Email + WhatsApp — implying logistics-svc also needs vendor contact data. This is a classic microservices boundary question: who owns the vendor master? Both services need vendor data; one of them authoritatively owns it.

Separately, the PRD is silent on **vendor onboarding** (how a new vendor enters the system, what KYC is required) and **vendor rating** (performance tracking, blocklist, preferred-vendor tier). Both are operationally important — the agency ends up working with 20–50 vendors across categories (koper supplier, kain ihram supplier, printing, freight, hotel, airline ground handlers, etc.) and needs to track who's reliable.

## The question

1. **Which service owns the vendor master record?** `finance-svc` (per PRD line 449) or `logistics-svc`?
2. **If finance-svc owns:** how does logistics-svc access vendor contact/dispatch data? gRPC read-through, cached event replication, or shared read-only DB view?
3. **Single master or split?** — is there one `vendors` table, or separate `logistics.vendors` (suppliers of physical goods) and `finance.vendors` (including service vendors like printing, utilities)?
4. **Vendor onboarding workflow** — what fields are required (NPWP, bank account, PKP status for tax), who approves, what's the audit trail?
5. **Vendor rating** — does the system track delivery punctuality, QC pass rate, payment history? Is there a "preferred vendor" flag or tier?
6. **Vendor blocklist** — can a vendor be marked as do-not-use? What prevents a PR from selecting them?

## Options considered

- **Option A — finance-svc owns the vendor master; logistics-svc reads via gRPC.** Single source of truth matches the PRD's Database Vendor placement under Finance. Logistics calls `finance-svc.GetVendor(vendor_id)` at PO dispatch time.
  - Pros: one master; PRD-aligned; finance's accounting needs (NPWP, bank, AP aging) are central.
  - Cons: logistics adds a critical-path gRPC dependency for every dispatch; cross-service call in hot path.
- **Option B — logistics-svc owns the operational vendor record; finance-svc owns the financial vendor record; they share a `vendor_id` namespace.** Two tables, same ID. Logistics holds contact + delivery + rating. Finance holds bank + tax + AP.
  - Pros: each service has local data for its hot path; no cross-service read at PO dispatch.
  - Cons: dual schema to keep in sync on onboarding; risk of drift.
- **Option C — logistics-svc owns everything vendor-related (contact, bank, tax, rating).** Finance reads via gRPC when posting AP.
  - Pros: simpler; single owner.
  - Cons: contradicts PRD's Database Vendor placement; puts bank/tax data (sensitive) in logistics-svc's scope.

## Recommendation

**Option B — split ownership at the bounded-context seam, shared ID.**

Option A is clean on paper but puts a synchronous gRPC call from logistics-svc to finance-svc on every PO dispatch — an unnecessary hot-path dependency for data that barely changes. Option C contradicts the PRD and concentrates sensitive financial data in the wrong service. Option B respects bounded contexts: logistics needs contact + rating + delivery history (operational), finance needs bank + tax + AP aging (financial). Same `vendor_id` (UUID) lets them cross-reference without forcing one side to own both concerns.

Defaults to propose: onboarding is a workflow that writes to both services atomically — a new `vendor_id` allocated by a shared sequence or coordinated via a tiny bootstrap call. Onboarding fields: name, NPWP, PKP status, bank account (finance), vendor category, primary contact name + email + WhatsApp, delivery address (logistics), approver + approval date (audit). Rating lives in logistics-svc: delivery punctuality (GRN-date vs PO-date), QC pass rate (grn_lines.qc_passed), active-PO count. Blocklist: `status: active | blocked | deprecated` in logistics, mirrored to finance. PR's vendor picker filters out `blocked`. Vendor onboarding requires finance director approval (second tier — vendor onboarding is a procurement-side governance concern).

Reversibility: collapsing the two tables into one later is a migration; extending either is additive. Shared `vendor_id` means no data loss on restructure.

## Answer

**Decided:** **Option B split** — **shared `vendor_id`**: **finance-svc** = tax/bank/AP core; **logistics-svc** = ops contact + delivery performance + blocklist mirror; **onboarding atomic write** to both; **finance director approves** new vendor; ratings stay logistics.

**Date decided:** 2026-04-18  
**Decided by:** Documentation session 2026-04-18 (AI-assisted product defaults)
