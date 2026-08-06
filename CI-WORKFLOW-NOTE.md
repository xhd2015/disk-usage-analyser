# CI workflow note

## Status

**added**

Branch was missing `.github/workflows`; full doctest-style Test workflow was added and pushed.

## Branch / remote / push

| Field | Value |
|-------|--------|
| Branch | `master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Remote URL | `ssh://git@github.com/xhd2015/disk-usage-analyser.git` |
| Commit SHA | `e889ff89669ffeec6e7297c3f1880f2d5a6d127a` (short: `e889ff8`) |
| Push result | **success** — `455ee3e..e889ff8` → `origin/master-2026-08-06-use-go-best-practice-to-review-current-project` (upstream set) |

## Paths changed

- `.github/workflows/test.yml` (new)
- `script/ci/coverage-package-table.py` (new; module-aware package coverage table)
- `CI-WORKFLOW-NOTE.md` (this file)

## How to view Actions for this push

- Repo: https://github.com/xhd2015/disk-usage-analyser
- Actions (all): https://github.com/xhd2015/disk-usage-analyser/actions
- Branch runs: https://github.com/xhd2015/disk-usage-analyser/actions?query=branch%3Amaster-2026-08-06-use-go-best-practice-to-review-current-project
- Commit: https://github.com/xhd2015/disk-usage-analyser/commit/e889ff89669ffeec6e7297c3f1880f2d5a6d127a
- Workflow file: **Test** → `.github/workflows/test.yml`

## Workflow contents (aligned with doctest pattern)

- `on: push` + `pull_request`
- `actions/setup-go@v5` from `go.mod`
- Embed placeholder under `disk-usage-analyser-react/dist` (gitignored empty tree otherwise breaks `//go:embed`)
- Git identity for doctests that use git
- `go test ./...` with `-coverpkg=disk-usage-analyser/...` → `coverage-gotest.out`
- `go install github.com/xhd2015/doctest/cmd/doctest@latest`
- Doctest discovery: `--label '!slow && !ui-automation'` → `coverage-doctest-discovery.out`
- Doctest slow (non-UI): `--label 'slow && !ui-automation'` → `coverage-doctest-slow.out`
- xgo install + `xgo tool coverage merge`
- `GITHUB_STEP_SUMMARY` total + `script/ci/coverage-package-table.py`
- `actions/upload-artifact@v4` for coverage profiles

## How this differs from doctest’s workflow

| Aspect | doctest reference | this repo |
|--------|-------------------|-----------|
| Module / `COVERPKG` | `github.com/xhd2015/doctest/...` | `disk-usage-analyser/...` |
| Install doctest | `go install ./cmd/doctest` | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Labels | `!e2e` then `e2e` | `!slow && !ui-automation` then `slow && !ui-automation` (no `e2e` labels in this tree) |
| UI leaves | n/a | `ui-automation` **not** run (no browser stack on ubuntu-latest) |
| Embed | n/a | CI creates placeholder `dist` files |
| Git identity | n/a | Configured for git-using doctests |
| Package table | same idea | Adapted to module `disk-usage-analyser`; skips `script/` and `cmd/` |
| Coverage profiles | gotest + discovery + e2e | gotest + discovery + slow |

## Note on pre-commit hook

`git-hook-github-workflow-test` warns that the file differs from its **minimal** template (golang container + `doctest --label-all`). That template is intentionally **not** used here: this branch follows the fuller **doctest coverage** pattern requested by the user.
