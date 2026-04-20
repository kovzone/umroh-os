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
- Branch naming follows `§ Branch strategy + merge ownership` below.

## Branch strategy + merge ownership

Operational rules for the 2-developer repo. Contracted once here so both devs (and both agents) apply the same workflow.

### Current operating mode (S0, 2026-04-20)

For the duration of slice **S0** only, Lutfi delegated all Joint (`S0-J-*`) tasks and all S0 PR reviews to Elda. The five S0 PRs that shipped under this delegation — PR #2 (S0-J-01), #3 (S0-J-02), #4 (S0-J-03), #5 (S0-E-01), and this PR — were executed and merged by Elda.

Solo-exec workflow during S0:

- **Branches and PRs** continue as documented below — one card, one branch, one PR against `dev`.
- **Review:** Elda runs an AI-assisted code review pass (read the full diff, sanity-check against the card's DoD + the DoR/DoD table below) before clicking merge.
- **No notify-before-Joint-pickup step** — there is no other dev to notify during the delegation.
- **Scope is strictly S0.** Unless Lutfi and Elda explicitly agree to extend it, **S1+ reverts to the 2-developer rules documented in the rest of this section** (non-executor reviews, Joint pick-up-and-notify, both devs click merge on each other's PRs).

The rules that follow describe the durable 2-dev default — they are not rewritten, and they apply from S1 onward by default.

**Branch naming.** `feat/<slice>-<owner>-<seq>-<slug>` — all lowercase, hyphen-separated. `<slice>` is the slice code (`s0`, `s1`, …), `<owner>` is `j` / `e` / `l` matching the task code, `<seq>` is the two-digit sequence, `<slug>` is 2–5 words naming the card. One card = one branch = one PR.

- Examples: `feat/s0-j-02-branch-strategy`, `feat/s1-e-03-booking-draft`, `feat/s1-j-01-catalog-contract`.
- For non-slice ad-hoc work (e.g. a process-discipline card in `docs/91-progress/progress.md`), use `feat/<slug>` without the slice/owner/seq prefix.

**Short-lived branches.** A feat branch lives on the remote only as long as it takes to land its PR — target under a day for docs-only cards, under a week for code cards. Once the PR merges, delete the remote branch. The local copy can stay harmlessly until the dev prunes it with `git branch -d`.

**Merge ownership (non-executor reviews).** Elda merges Lutfi's PRs; Lutfi merges Elda's PRs. Joint (`Sx-J-*`) PRs are merged by whichever dev did not execute the card. The reviewer's merge click is the explicit sign-off — no separate approval mechanism needed.

**Protected trunks.** All PRs target `dev`, the integration branch. `main` moves only on release cuts via a dedicated `release` PR from `dev`. Neither `dev` nor `main` accepts direct pushes.

**No force-push on shared branches.** `dev`, `main`, and any `feat/*` branch that has been pushed to the remote are append-only. If history needs a fix (typo in a commit message, unwanted merge), open a follow-up commit rather than rewriting. Force-push is reserved for unpushed local branches.

## Definition of Ready / Definition of Done (per PR)

Every PR in this repo passes through these two gates. `CONTRIBUTING.md § Canonical References` cross-links here — this table is the authoritative version.

> **Note for S0 (2026-04-20):** during the solo-exec operating mode described in `§ Branch strategy + merge ownership § Current operating mode`, the DoR row "the other dev has been notified" is N/A (there is no other dev to notify), and the DoD row "reviewer signs off by merging" means Elda self-merges after an AI-assisted review pass. Both rows revert to their 2-dev wording from S1 onward by default.

| Phase | Check |
|---|---|
| **DoR** | Task card exists in a progress tracker with an explicit DoD checklist |
| **DoR** | For `Sx-E-*` / `Sx-L-*` cards: matching `Sx-J-*` contract section is merged in `slice-Sx.md` (not `TBD`) **and** any **Engineering freeze** companions listed for that slice in `docs/00-overview/05-slice-engineering-checklist-and-task-codes.md` (e.g. `S1-E-01`, `S1-L-01`) are satisfied |
| **DoR** | Branch name follows `feat/<slice>-<owner>-<seq>-<slug>` per `§ Branch strategy + merge ownership` |
| **DoR** | For Joint (`Sx-J-*`) cards: the other dev has been notified before coding starts |
| **DoD** | Code builds; tests pass; `go vet` clean (for code PRs) |
| **DoD** | `docs/92-testing/testing-guide.md` replaced with a verification block specific to this task (per each dev's private per-task convention) |
| **DoD** | `/security-review` run and findings addressed when any file outside `docs/` changes |
| **DoD** | Reviewer signs off by merging the PR (non-executor dev clicks merge per `§ Branch strategy + merge ownership`) |

DoR items are checked before coding starts; DoD items are checked before the reviewer clicks merge. A PR that fails a DoR item is sent back for prep; a PR that fails a DoD item waits for the gap to close.

## Service ownership matrix (S1–S2)

Default PR-owner and reviewer for each backend service touched in slices S1–S2. The matrix is the starting point — either dev can still pick up the other's card when needed (per `§ Branch strategy + merge ownership`, Joint cards are pick-up-and-notify). This table names the **default**, not the mandatory executor.

| Service | PR-owner | Reviewer | Active in | Notes |
|---|---|---|---|---|
| `iam-svc` | Elda | Lutfi | S1 (staff auth), F1-MIN | Per doc 04 RACI — A=Lutfi, R=Elda, C=Lutfi |
| `catalog-svc` | Elda | Lutfi | S1 (read endpoints) + S2 (seat release during refund saga) | Owns `ReserveSeats` / `ReleaseSeats` gRPC |
| `booking-svc` | Elda | Lutfi | S1 (draft + saga) + S2 (paid-state callback) | Saga orchestrator for submit/refund per ADR 0006 |
| `payment-svc` | Elda | Lutfi | S2 | Webhooks, invoice/VA issuance, reconcile cron |
| `finance-svc` | Elda | Lutfi | S3 (forward-reference) | Listed here so the ownership default extends into S3 without another card |
| `logistics-svc` | Elda | Lutfi | S3 (forward-reference) | Same reasoning as `finance-svc` |

**Cross-service gRPC contracts are Joint.** Any change to a gRPC method signature, request/response shape, or error-code convention between two services is a Joint change regardless of which service owns the implementation. Before editing:

1. Notify the other dev (WA / issue comment) that the gRPC shape will change.
2. Open a matching `Sx-J-*` contract PR against `docs/contracts/slice-Sx.md § <section>` (new or amended) **before** the implementation PR — the contract is the freeze point.
3. The implementation PR then references the merged contract PR. No silent drift.

Services not listed in the table (`gateway-svc`, `jamaah-svc`, `visa-svc`, `ops-svc`, `crm-svc`, and the deferred `broker-svc`) follow the same default (Elda as PR-owner, Lutfi as reviewer) until a future card formalises them when their slice becomes active.

## Appendix slots (filled by later S0 cards)

All three appendix slots are now filled — the S0 docs chain is complete:

- **Branch strategy + merge ownership** — ✅ landed in `§ Branch strategy + merge ownership` via `S0-J-02` (2026-04-20).
- **DoR / DoD per PR** — ✅ landed in `§ Definition of Ready / Definition of Done (per PR)` via `S0-J-03` (2026-04-20).
- **Service ownership matrix (S1–S2)** — ✅ landed in `§ Service ownership matrix (S1–S2)` via `S0-E-01` (2026-04-20).

## Related references

- Task codes + slice definitions: `docs/00-overview/05-slice-engineering-checklist-and-task-codes.md`
- Backlog mapping (`BL-*` IDs): `docs/00-overview/06-feature-to-backlog-mapping.md`
- Commit message format: `docs/08-commit-conventions.md`
- ADR 0006 (in-process sagas, Temporal deferred to F6): `docs/01-architecture/adr/0006-defer-temporal-to-f6.md`
- ADR 0007 (migration-based schema, single-DB multi-schema): `docs/01-architecture/adr/0007-migration-based-schema.md`
