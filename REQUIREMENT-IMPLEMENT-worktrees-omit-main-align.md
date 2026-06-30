# REQUIREMENT-IMPLEMENT — Git Worktrees: omit main row + left align

## Context

Followup to Repository Scans Git Worktrees section. User feedback:
- Q1: Omit main worktree child row — repo path row IS the main checkout
- Q2: Tree not left-aligned — fix layout

## Tests are sealed — do not modify

```
tests/tmp-worktrees-scan/          (10 leaves, updated)
tests/tmp-analyse-frontend-test-cases/worktrees-section/  (updated + left-aligned leaf)
```

Do not edit ASSERT.md, SETUP.md, playwright scripts, or DOCTEST.md in sealed trees.

## Implementation tasks

### Backend `server/tmp_worktrees.go`

1. Extend `repo` SSE event with `size`, `sizeHuman`, `fileCount` from main checkout sizing
2. Stream `repo` event before linked worktree events
3. Emit `worktree` events **only** for linked worktrees (`isMain=false`)
4. Do NOT emit `worktree` for main checkout
5. Summary `worktrees` count = linked only; include repos with main-only

### Frontend `disk-usage-analyser-react/src/TmpFilesAnalyse.tsx`

1. Add state for repo rows from `repo` SSE events (`WorktreeRepoRow`)
2. `worktree` SSE → linked children only, grouped by `repoPath`
3. Repo row displays `repoPath` + main size (from repo event), not sum of children
4. Remove main worktree child row and `worktree-main-badge`
5. Left-align: `worktrees-tree` with `width: 100%`, `textAlign: 'left'`; repo rows flush left; linked children indented consistently
6. Main-only repos: show repo row only, no children block

### Verify

```sh
doctest test ./tests/tmp-worktrees-scan
doctest test ./tests/tmp-analyse-frontend-test-cases/worktrees-section/left-aligned
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/worktrees-section/...
doctest test ./...
```

Rebuild frontend before UI tests: `cd disk-usage-analyser-react && bun run build`