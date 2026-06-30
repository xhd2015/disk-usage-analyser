# Requirement — Git Worktrees UI: omit main row + left align

## Summary

Improve the **Git Worktrees** section in Repository Scans:

1. **Omit main worktree child row** — The repo parent row (`repoPath`) represents the main checkout. Do not render or stream a separate child for `isMain=true`. Only **linked** worktrees appear as children.
2. **Left-align tree** — Worktree tree content must be flush left within the card (not centered or mis-indented).

## Confirmed decisions

| # | Decision |
|---|----------|
| Q1 | Omit main worktree row; repo dir row = main repo with its on-disk size |
| Q2 | Tree rows left-aligned |

## Data model changes

### SSE `repo` event (extend)

```json
{
  "repoPath": "~/Projects/foo",
  "repoName": "foo",
  "size": 800000000,
  "sizeHuman": "800 MB",
  "fileCount": 1234
}
```

- `size` / `fileCount` = main checkout sizing (the repo directory itself)
- Stream `repo` before linked `worktree` events for that repo

### SSE `worktree` event (linked only)

- Emit **only** when `isMain=false`
- No `worktree-main-badge` in UI
- `IsMain` field may remain on type for internal use but must not surface as a child row

### Main-only repos

- Emit `repo` event with main checkout size
- Emit **zero** `worktree` events (no linked worktrees)
- UI shows repo row only, no children

### Frontend state

```ts
interface WorktreeRepoRow {
  repoPath: string;
  repoName: string;
  size: number;
  fileCount: number;
}
// worktreeHits: linked worktrees only (grouped by repoPath)
```

Listen to SSE `repo` events → populate repo rows; `worktree` → append linked children.

Repo row size = main checkout size (from `repo` event), **not** sum of children.

## UI layout (left-aligned)

```
▼ ~/Projects/foo                    800 MB    ← main repo (repo event)
    ~/Projects/foo-wt (feature)     400 MB    ← linked only
▼ ~/Projects/solo                   50 MB     ← main-only, no children
```

- `worktrees-tree`: `width: 100%`, `textAlign: 'left'`
- Repo rows: no extra left indent
- Linked children: consistent indent (e.g. 16px)
- Remove `worktree-main-badge` entirely

`data-testid` hooks unchanged except:
- `worktree-main-badge` must **not** appear
- Add `data-testid="worktrees-tree-aligned"` or assert on `worktrees-tree` computed `text-align: left`

## Test plan

### Update `tests/tmp-worktrees-scan/`

| Leaf | Change |
|------|--------|
| `discovery/main-only` | Expect `repo` event with size; **zero** worktree events |
| `discovery/main-plus-linked` | Expect `repo` with main size; worktree events only for linked (`isMain=false`); no main worktree event |
| `streaming/repo-before-worktree` | Still valid (repo before linked worktree) |
| `sizing/non-empty-checkout` | Repo event size > 0 for main |

Add leaf if needed: `discovery/omit-main-worktree-event`

### Update `tests/tmp-analyse-frontend-test-cases/`

| Leaf | Change |
|------|--------|
| `worktrees-section/after-scan` | Remove `worktree-main-badge` expectation; allow repo-only rows (zero children) |
| `worktrees-live-stream` | Still streams repo + linked rows |
| New: `worktrees-section/left-aligned` | Playwright checks `worktrees-tree` has `text-align: left` |

## Out of scope

- Binary files section unchanged
- Delete worktrees unchanged (read-only)

## Verify

```sh
doctest test ./tests/tmp-worktrees-scan
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/worktrees-section/...
doctest test ./...
```