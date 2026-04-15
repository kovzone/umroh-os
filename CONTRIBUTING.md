# Contributing to UmrohOS

This document is the shared contribution workflow for all developers and AI tools used in this repository.

## Goal

Keep changes merge-ready, traceable, and safe across a mixed-tool workflow (Cursor, Claude Code, or manual coding).

## Branch and PR Workflow

- Work from `main` by creating a short-lived branch per task.
- Keep PR scope focused on one problem area.
- Rebase or merge from `main` when your branch is stale before review.
- Do not force-push shared branches unless explicitly coordinated.

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

## Shared vs Local Configuration

Shared (commit to repo):

- Product/domain/architecture conventions in `docs/`.
- Team-level workflow rules in this file.
- AI onboarding and routing guidance in `AGENTS.md`.
- PR template and other repo-level collaboration files.

Local only (do not commit):

- Agent-specific private instructions (`CLAUDE.md`, `.cursor/`, and similar local setup).
- Personal editor preferences (theme, keymap, UI preferences).
- Local helper scripts or notes that are not part of team workflow.

## AI-Assisted Work Expectations

- AI output is a draft; developers remain responsible for correctness and review.
- Follow the authority hierarchy and domain boundaries documented in `AGENTS.md` and linked `docs/`.
- If product behavior is ambiguous and open questions exist, do not silently invent final behavior.

## Pull Request Checklist

Use the PR template in `.github/pull_request_template.md` and ensure all required sections are filled before requesting review.