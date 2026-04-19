# Commit Message Conventions

Repo-wide — applies to both the Go services codebase and the Svelte 5 + Vite codebase, and to any AI coding agent or human committing on this repo.

## Format

```
<type>: <short message in lower case>

<optional body: longer description in normal case, paragraphs or bullets>
```

## Types


| Type       | When to use                                             |
| ---------- | ------------------------------------------------------- |
| `feat`     | New feature or capability                               |
| `fix`      | Bug fix                                                 |
| `docs`     | Documentation only                                      |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test`     | Adding or updating tests                                |
| `chore`    | Maintenance tasks (deps, configs, CI)                   |
| `build`    | Build system, Docker, Makefile changes                  |
| `perf`     | Performance improvement                                 |
| `style`    | Formatting, whitespace, no code change                  |


## Rules

1. Subject line is **all lower case** (type + message).
2. No period at the end of the subject line.
3. Keep subject line under **~72 characters**.
4. Separate subject from body with a blank line.
5. Body uses normal letter case.
6. Body is optional — use it when the *what* isn't obvious from the subject, or when the *why* matters for future readers.
7. No mandatory scope syntax (`feat(foo): ...`) — include the affected service or area in the subject prose when useful. Consistency of style matters more than parenthetical metadata.

## Examples

Simple:

```
feat: add bulk-submit endpoint to visa-svc
```

```
fix: reject booking submit when passport expiry < 180d before departure
```

```
docs: draft f6 visa pipeline spec + file q026–q031
```

With body:

```
feat: wire apperrors sentinels to http status codes in iam-svc

Maps the domain sentinels defined in docs/04-backend-conventions/02-error-handling.md
to http responses via the shared ErrorHandler middleware:

- ErrNotFound           → 404
- ErrInvalidInput       → 400
- ErrUnauthorized       → 401
- ErrForbidden          → 403
- ErrAlreadyExists      → 409
- ErrConflict           → 409
- ErrPreconditionFail   → 412
- ErrUpstreamUnavail    → 503
- everything else       → 500

The middleware also strips internal error details from the client-facing
response and logs the full error chain with trace_id.
```

```
refactor: collapse catalog-svc seat-reservation into single atomic sql

Replaces the two-step read-then-update pattern with a single
UPDATE ... WHERE reserved_seats + $n <= total_seats RETURNING ...
statement. Removes the race window under concurrent booking submits;
verified by the new k6 race-condition test at tests/load/catalog_race.js.
```

## Notes on AI-assisted commits

Each developer's AI coding agent may have its own **attribution footer** convention (e.g. `Co-Authored-By: <agent name> <email>`). Attribution footers are allowed and encouraged — they're a useful audit trail for human reviewers — but they are **not part of the repo-wide convention**. Do not fail a PR because a commit lacks (or carries) an AI attribution footer.

## Rationale

- **Lowercase subjects** are easier to scan in long logs and play nicely with GitHub's display and search.
- **No scope parens** — scopes are useful but often wrong or inconsistent (a single commit commonly touches multiple services). Prose is clearer.
- **Type prefix** keeps `git log --oneline` groupable and keeps release-note tooling (if we ever add one) workable.
- **Short-subject + optional body** is the Conventional Commits idea without the full ceremony.

## Related

- `docs/04-backend-conventions/08-git-workflow.md` — backend git workflow (branching, PRs, reviewer expectations). Defers to this file for commit message format.
- Frontend git workflow lives in `docs/05-frontend-conventions/` and follows the same commit format.

