# Contributing to UmrohOS

This document is the shared contribution workflow for all developers and AI tools used in this repository.

## Goal

Keep changes merge-ready, traceable, and safe across a mixed-tool workflow (Cursor, Claude Code, or manual coding).

## Branch and PR Workflow

- Start from latest `dev` (the integration branch; `main` moves only on release cuts), then create a short-lived task branch before making any commit.
- Keep PR scope focused on one problem area.
- Rebase or merge from `dev` when your branch is stale before review.
- Do not force-push shared branches unless explicitly coordinated.

## Non-negotiable Git Rules

- MUST create a task branch from `dev` before any commit.
- MUST NOT commit directly on `dev` or `main`.
- MUST open a PR targeting `dev` for all shared changes; no direct push to protected trunks (`dev`, `main`).
- When sharing a GitHub link so the PR form opens **with title and description already filled**, use a **compare URL** with `quick_pull=1` and URL-encoded `title` / `body` (not `/pull/new/<branch>` alone). See **Pre-filled PR open links** below.
- MUST keep one PR focused on one concern/task.
- MUST include verification evidence in the PR before requesting review.
- MUST NOT force-push to `dev` or `main`, or to any `feat/*` branch once pushed.
- MUST NOT bypass required checks or repository protections.

## Quality Gate (Minimum Before Merge)

Every PR must include:

1. Clear scope:
  - What changed.
  - What intentionally did not change.
2. Risk notes:
  - Possible regression areas.
  - Rollback strategy if needed.
3. Verification evidence:
  - Relevant local commands for changed area (for example: lint, test, build).
  - Result summary for each command.
4. Docs/spec impact:
  - If behavior changes, update relevant docs/specs in the same PR.
  - If no behavior change, explicitly state `No doc impact`.
  - New or updated committed prose under `docs/` MUST be **English** — see `docs/00-overview/00-documentation-language.md`.
  - Do not commit a doc-only batch until **all** `docs/` files in that commit are fully English for the scope you are migrating (avoid mixed-language commits).

## Shared vs Local Configuration

Shared (commit to repo):

- Product/domain/architecture conventions in `docs/`.
- Team-level workflow rules in this file.
- AI onboarding and routing guidance in `AGENTS.md`.
- PR template and other repo-level collaboration files.
- Shared Cursor agent rules under `.cursor/rules/` (see `.gitignore` — other `.cursor/*` paths stay local).

Local only (do not commit):

- Agent-specific private instructions (`CLAUDE.md`, local `.cursor/` files outside `rules/`, and similar local setup).
- Personal editor preferences (theme, keymap, UI preferences).
- Local helper scripts or notes that are not part of team workflow.

## AI-Assisted Work Expectations

- AI output is a draft; developers remain responsible for correctness and review.
- Follow the authority hierarchy and domain boundaries documented in `AGENTS.md` and linked `docs/`.
- If product behavior is ambiguous and open questions exist, do not silently invent final behavior.

## Pre-filled PR open links

GitHub can open the **Create pull request** screen with **title** and **description** already populated when you use a **compare** URL with `quick_pull=1` and encoded `title` / `body` parameters. Official reference: [Using query parameters to create a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/using-query-parameters-to-create-a-pull-request).

Pattern (integration PRs target `dev`):

```text
https://github.com/<owner>/<repo>/compare/dev...<your-branch>?quick_pull=1&title=<url-encoded>&body=<url-encoded>
```

Important:

- The `body` parameter **replaces** the web UI default from `.github/pull_request_template.md`; it does not append to it. Put the **full** template sections (same headings as the file) into `body`, filled out for this PR, then URL-encode.
- Very long URLs can hit server limits (**414**). If so, open the compare page with `quick_pull=1` and title only, then paste a pre-written body, or use `gh pr create` with `--body-file`.

Cursor agents follow the same convention in `.cursor/rules/pr-prefilled-open-link.mdc`.

## Pull Request Checklist

Use the PR template in `.github/pull_request_template.md` and ensure all required sections are filled before requesting review. For compare-link workflows, embed that same filled content in the `body` query parameter (see above).

## CI and path filters

GitHub Actions workflow **`.github/workflows/ci.yml`** (card **S0-J-06**) runs on every `pull_request` and on `push` to `dev` / `main`. It uses path filters so **expensive jobs are skipped** when the diff touches neither Go services nor `apps/core-web`:

| Job | Runs when the PR / push changes any of |
| --- | --- |
| **Go unit tests** (`make test`) | `services/**`, `Makefile`, `migration/**`, `docker-compose.dev.yml`, or `.github/workflows/ci.yml` |
| **core-web** (`npm run check` + `npm test` in `apps/core-web`) | `apps/core-web/**` or `.github/workflows/ci.yml` |
| **Skip code checks** (fast no-op) | Neither of the above (for example **docs-only** or other non-code paths) |

Mixed PRs (for example `services/` + `docs/`) still match the backend and/or web filters, so the full relevant matrix runs. Editing **only** this workflow file matches both filters on purpose so CI self-validates.

## Canonical References

- Detailed Git workflow: `docs/04-backend-conventions/08-git-workflow.md`
- Detailed commit message conventions: `docs/08-commit-conventions.md`
- Per-PR Definition of Ready / Definition of Done: `docs/contracts/README.md § Definition of Ready / Definition of Done (per PR)`
- Branch strategy + merge ownership: `docs/contracts/README.md § Branch strategy + merge ownership`
- Team AI onboarding and doc authority: `AGENTS.md`