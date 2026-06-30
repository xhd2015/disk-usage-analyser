# Requirement — Binary files section left align

## Summary

Mirror the Git Worktrees left-align fix for the **Binary files** section tree.

## Problem

`binaries-tree` lacks `width: 100%` and `textAlign: 'left'` (worktrees-tree has both). Tree content appears not left-aligned within the card.

## Expected UI

- `binaries-tree`: `width: 100%`, `textAlign: 'left'`
- Repo rows and binary child rows flush left; children indented consistently (e.g. 24px, matching current indent)
- Selection bar (`binary-selected-total`, `binary-delete-btn`) unchanged

## Test plan

Add leaf under `tests/tmp-analyse-frontend-test-cases/binaries-section/`:

| Leaf | Verifies |
|------|----------|
| `left-aligned` | After binaries scan, `binaries-tree` computed `text-align` is `left` |

Playwright script mirrors `worktrees-section/left-aligned/worktrees-left-aligned.js`:
- Scan binaries, wait for done badge
- Assert `getComputedStyle(binaries-tree).textAlign === 'left'`
- Label: `slow, ui-automation`

Update `tests/tmp-analyse-frontend-test-cases/DOCTEST.md` index.

## Implementation (for implementer)

`TmpFilesAnalyse.tsx`: apply same styles as worktrees-tree to `binaries-tree` and child row containers.

## Verify

```sh
doctest test ./tests/tmp-analyse-frontend-test-cases/binaries-section/left-aligned
doctest test ./...
```