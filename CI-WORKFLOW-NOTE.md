# CI workflow note

## Status

**added** (pending push metadata below after `git push`)

## Branch / remote / push

- Branch: `master-2026-08-06-use-go-best-practice-to-review-current-project`
- Remote: `ssh://git@github.com/xhd2015/disk-usage-analyser.git`
- Commit SHA: _(filled after push)_
- Push result: _(filled after push)_

## Paths changed

- `.github/workflows/test.yml` (new)
- `script/ci/coverage-package-table.py` (new; doctest-style package table for this module)
- `CI-WORKFLOW-NOTE.md` (this file)

## How to view Actions for this push

- Repo: https://github.com/xhd2015/disk-usage-analyser
- Actions: https://github.com/xhd2015/disk-usage-analyser/actions
- Branch filter: `master-2026-08-06-use-go-best-practice-to-review-current-project`
- Workflow name: **Test** (`.github/workflows/test.yml`)

## How this differs from doctest’s workflow

| Aspect | doctest reference | this repo |
|--------|-------------------|-----------|
| Module / `COVERPKG` | `github.com/xhd2015/doctest/...` | `disk-usage-analyser/...` |
| Install doctest | `go install ./cmd/doctest` | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Labels | `!e2e` then `e2e` | `!slow && !ui-automation` then `slow && !ui-automation` (no `e2e` labels here) |
| UI leaves | n/a | `ui-automation` not run (no browser stack on ubuntu-latest) |
| Embed | n/a | Placeholder under `disk-usage-analyser-react/dist` so `//go:embed` compiles |
| Git identity | n/a | Global user for doctests that invoke git |
| Package table | `script/ci/coverage-package-table.py` | Same helper, adapted to module `disk-usage-analyser` |
| Coverage profiles | gotest + discovery + e2e | gotest + discovery + slow |
