# Cross-slice Integration Contracts

This folder is the **wire-level agreement** between services, cut per delivery slice (`S0` → `S5` and beyond). One Markdown file per slice, with JSON snippets inline. Both developers read and review these files; no slice starts coding until its contract is merged. This file is the folder's index + the shared rulebook that every `slice-Sx.md` honors.

## Why this exists

UmrohOS is built across many short AI sessions by two full-stack developers working parallel slices (see `docs/00-overview/04-delivery-plan-2p-sequence-first.md`). Without a single pinned-down contract per slice, a session rewriting `booking-svc` can silently diverge from a session rewriting `catalog-svc`. Contracts are the pin. Anything an AI session changes about cross-service shape must land here first, and must be reviewed by the non-executor owner before code follows.

## How this folder relates to the rest of the docs

```
docs/06-features/          ← user POV: what the feature does, for which role
         │
         │ shared middle layer — read by both devs
         ▼
docs/contracts/            ← THIS folder: cross-service wire-level agreement, per slice
         │                   (request/response shape, gRPC params, event payloads, error codes)
         ▼
docs/03-services/<svc>/    ← per-service tech spec: how this one service honors its side of the contract
         │
         ▼
services/<svc>/             ← the code itself (openapi.yaml, sqlc queries, handlers)
```

- Read `06-features/` to understand *why* a contract exists.
- Read `contracts/slice-Sx.md` to understand *what the wire looks like*.
- Read `03-services/<svc>/01-api.md` to understand *how one service's endpoint is built*.

## File layout

- `README.md` (this file) — folder intro, rules, forward-pointers to DoR/DoD + branch strategy + ownership appendix (those land via `S0-J-02`, `S0-J-03`, `S0-E-01`).
- `slice-Sx.md` — **template** for new slice contracts. Copy to `slice-S1.md`, `slice-S2.md`, etc. Delete sections that don't apply to the slice.
- `slice-S{n}.md` — actual per-slice contracts (`slice-S1.md`, `slice-S2.md`, ...) created as each slice's `S{n}-J-*` cards land.

## Conventions for every `slice-Sx.md`

1. **Frontmatter**: slice code, status (`draft` / `frozen`), `last_updated` date, PR-owner + reviewer.
2. **Sections**: one per integration point — Catalog, Booking, Inventory, Webhook, Events, etc. Only include the ones that apply to this slice.
3. **Every REST endpoint** documents: method + path, request body shape (JSON), response shape (JSON), error codes, idempotency key (if any), auth requirements.
4. **Every gRPC method** documents: service + method name, params, return shape, failure codes, compensation (if part of a saga).
5. **Every event** documents: event name, trigger, payload shape, producer, consumer(s).
6. **Always** include a `§ Changelog` section at the bottom of the file — one line per contract change with date + rationale.

## Bump-versi rule (contract change protocol)

When a merged contract needs to change after downstream services have already built against it, use **one** of these two mechanisms:

- **Changelog append** (default, for additive or backwards-compatible changes): add a dated entry at the bottom of the existing `slice-Sx.md` `§ Changelog` section. Example: *2026-05-03 — added optional `promo_code` field to `POST /v1/bookings` request body; existing clients unaffected.*
- **Bump to v2** (for breaking changes): copy the file to `slice-Sx-v2.md`, state the break explicitly at the top, keep `slice-Sx.md` intact so in-flight services can migrate. Deprecate `slice-Sx.md` once all consumers cut over.

Either way, the change must be reviewed by the non-executor owner (Lutfi reviews Elda's contract changes and vice versa) before any service is allowed to rely on the new shape. See `docs/00-overview/05-slice-engineering-checklist-and-task-codes.md` § "Short tips" #3 for the rationale.

## Authoring / review expectations

- Contract PRs are small by design — one slice, one PR, under an hour of review.
- Both devs read every contract PR, even when only one executes.
- Contract code follows `docs/08-commit-conventions.md` for commit format (`docs:` type).
- Branch naming: `feat/<slice>-<owner>-<seq>-<slug>` (e.g. `feat/s1-j-01-catalog-contract`) — see branch-strategy appendix below (lands via `S0-J-02`).

## Appendix slots (filled by later S0 cards)

This folder's README deliberately omits three items that other S0 cards will land here once they merge:

- **Branch strategy + merge ownership** → appended by `S0-J-02` (short paragraph covering branch naming rule confirmation, short-lived-branch policy, who-merges-whose-PR, protected-main + no-force-push).
- **DoR / DoD per PR** → appended by `S0-J-03` (short scannable table of Definition-of-Ready and Definition-of-Done columns).
- **Service ownership matrix (S1–S2)** → appended by `S0-E-01` (table mapping each backend service touched in S1–S2 to a PR-owner and a code-reviewer).

Do not add those sections here yet — they belong to those cards. This README is intentionally a scaffold, expanded as the S0 chain completes.

## Related references

- Task codes + slice definitions: `docs/00-overview/05-slice-engineering-checklist-and-task-codes.md`
- Backlog mapping (`BL-*` IDs): `docs/00-overview/06-feature-to-backlog-mapping.md`
- Commit message format: `docs/08-commit-conventions.md`
- ADR 0006 (in-process sagas, Temporal deferred to F6): `docs/01-architecture/adr/0006-defer-temporal-to-f6.md`
- ADR 0007 (migration-based schema, single-DB multi-schema): `docs/01-architecture/adr/0007-migration-based-schema.md`
