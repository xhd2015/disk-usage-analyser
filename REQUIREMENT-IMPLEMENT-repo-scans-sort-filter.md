# REQUIREMENT-IMPLEMENT — Repository Scans sort + filter

## Context

Sort both sections by size DESC (dynamic resort on SSE). Frontend-only size filters:
- Binaries: `☐ <1M` default unchecked → hide <1 MiB
- Worktrees: `☐ <10M` default unchecked → hide <10 MiB
- Checked = show all sizes

## Tests sealed — do not modify

```
tests/repository-scans-layout/           (8 leaves)
tests/tmp-analyse-frontend-test-cases/binaries-section/filter-*
tests/tmp-analyse-frontend-test-cases/binaries-section/sort-by-size-desc
tests/tmp-analyse-frontend-test-cases/worktrees-section/filter-*
tests/tmp-analyse-frontend-test-cases/worktrees-section/sort-by-size-desc
tests/tmp-analyse-frontend-test-cases/DOCTEST.md (index)
```

Also do not modify `layout-harness.ts` or fixture JSON in sealed tree.

## Implementation

### 1. `disk-usage-analyser-react/src/repositoryScansLayout.ts`

Export per REQUIREMENT-DESIGN:
- `ONE_MIB = 1048576`, `TEN_MIB = 10485760`
- `sortWorktreeRepos`, `sortLinkedWorktrees`
- `filterWorktreeRepos(repos, linkedByRepo, showUnder10M)`
- `sortBinaryRepos`, `filterBinaryRepos(byRepo, showUnder1M)`
- Resort: re-apply sort on every render/state update

Filter rules:
- Binary: hide file <1M; hide repo if remaining total <1M
- Worktree: hide linked <10M; hide repo if main size <10M

### 2. `TmpFilesAnalyse.tsx`

- State: `showBinaryUnder1M` default false, `showWorktreeUnder10M` default false
- Checkboxes: `data-testid="binary-show-under-1m"` label `<1M`, `worktree-show-under-10m` label `<10M`
- Apply sort+filter before render; re-run on each SSE append
- Place checkboxes in section header near Scan/Stop

Rebuild: `cd disk-usage-analyser-react && bun run build`

## Verify

```sh
doctest test ./tests/repository-scans-layout
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/binaries-section/filter-under-1m-default
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/binaries-section/filter-show-under-1m
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/binaries-section/sort-by-size-desc
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/worktrees-section/filter-under-10m-default
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/worktrees-section/filter-show-under-10m
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/worktrees-section/sort-by-size-desc
doctest test ./tests/tmp-analyse-frontend-test-cases/binaries-section/...
doctest test ./tests/tmp-analyse-frontend-test-cases/worktrees-section/...
doctest test ./...
```