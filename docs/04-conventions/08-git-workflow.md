# Git Workflow

UmrohOS is built across many short Claude sessions. Git is how we capture the trail. The conventions below keep history readable and reviewable.

## Branching

- `main` — protected. Always green. The reviewer merges.
- `feature/<short-name>` — one branch per task in `progress.md`. Branch from `main`, merge back via PR.
- No long-lived feature branches. A session's work either merges or is parked as a draft PR by the end of the session.

## Commit messages

Format:

```
<type>(<scope>): <subject>

<optional body>

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
```

| Type | Use |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code change without behavior change |
| `docs` | Docs only |
| `chore` | Tooling, config, dependencies |
| `test` | Tests only |

Scope is the service name (`iam-svc`, `catalog-svc`) or `docs`, `infra`, `proto`.

Examples:

```
feat(iam-svc): add CreateUser endpoint
docs: scaffold service docs for all 11 services
chore(infra): wire iam-svc into docker-compose
```

## Pull requests

Every PR includes:
- **Summary:** 1-3 bullets on what changed and why.
- **Linked task:** the matching item in `progress.md`.
- **Verification:** link to the verification block in `testing-guide.md`.
- **Test plan:** what was tested locally; what the reviewer should run.

PRs are small. One task = one PR. If a session produces more than one task's worth of work, split it.

## Reviewer expectations

- The reviewer walks the verification block in `testing-guide.md` before merging.
- The reviewer promotes `[~]` to `[x]` in `progress.md` after verification.
- The reviewer merges to `main`.

## Hard rules

- **Never commit `config.json`** (secrets). Only `config.json.sample`.
- **Never commit generated code that's reproducible from source** unless the project requires it for some downstream tool (e.g. some sqlc setups commit `sqlc.gen.go` for editor support — confirm per project).
- **Never `git push --force` to `main`.**
- **Never amend a commit on `main`** after it has been pushed.
- **Never bypass hooks** (`--no-verify`) — fix the underlying issue.
- **Never commit unless the reviewer explicitly asks.** Sessions leave the diff for review by default.
