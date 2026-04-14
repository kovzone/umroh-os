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
