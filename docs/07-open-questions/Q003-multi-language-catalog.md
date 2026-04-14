---
id: Q003
title: Multi-language catalog content scope
asked_by: session 2026-04-14 F2 draft
asked_date: 2026-04-14
blocks: F2
status: open
---

# Q003 — Multi-language catalog content scope

## Context

The PRD is in Bahasa Indonesia and doesn't explicitly call out multi-language product content. But pilgrimage packages often target audiences who speak Arabic (returning diaspora) or English (high-end international customers). This decision affects the catalog schema (per-field translations vs single text column) and the B2C rendering layer.

Locking this early is cheap; retrofitting after live packages exist is painful.

## The question

1. Does the catalog need to store product content (name, description, highlights, itinerary day titles) in **multiple languages** from day one?
2. If yes, which languages — Bahasa Indonesia only, Bahasa + English, or Bahasa + English + Arabic?
3. If multi-language, is it **per-package opt-in** (some packages are Bahasa-only) or **all packages must have all languages**?

## Options considered

- **Option A — Bahasa only (Recommended default for MVP).** Single `text` columns. Simplest. Matches PRD.
  - Pros: no schema complexity, fastest to ship.
  - Cons: retrofitting later means a migration + editor UI rework.
- **Option B — Bahasa + English (schema-ready from day one).** Text columns become `jsonb` with `{"id": "...", "en": "..."}` shape, or dedicated translation tables. Only `id` populated at first; English fields added gradually.
  - Pros: future-proof; marginal upfront cost.
  - Cons: editor UI complexity even if English is empty; translators needed.
- **Option C — Bahasa + English + Arabic.** Full trilingual from day one.
  - Pros: addresses all three user segments.
  - Cons: triples translation load; Arabic RTL considerations; questionable demand vs cost.

## Recommendation

**Option A — Bahasa only for MVP.** The PRD is in Bahasa, the target market is Indonesian pilgrimage agencies, and the CS/agent workflows assume Bahasa end-to-end. Shipping a multilingual editor that 100% of users ignore is waste. The right time to add English or Arabic is after UmrohOS has Indonesian market fit and a specific customer or business reason asks for it.

Hedge: store text columns as plain `text` today, but **reserve the field names** (`name`, `description`, `highlights`) so a later migration to `jsonb` with `{"id": "...", "en": "..."}` is an in-place column-type change, not a schema rework. Add a comment in the data model doc flagging this.

Option B (schema-ready day one) sounds responsible but the editor UX cost is real — every field doubles in complexity and admins have to tab through empty English fields forever. Option C (trilingual) is a different product; not now.

Reversibility: plain `text` → `jsonb` migration is routine; the cost is one migration and an editor UI iteration. Low to medium commitment — defer, but be ready.

## Answer

TBD — awaiting stakeholder discussion.
