# Requirement — Repository Scans: size sort + size filters

## Summary

Enhance **Git Worktrees** and **Binary files** sections:

1. **Sort by size DESC** — repos and children reorder dynamically as SSE items arrive (largest on top).
2. **Binary filter checkbox** `☐ <1M` — frontend-only; **default unchecked** → hide items under 1 MiB.
3. **Worktree filter checkbox** `☐ <10M` — frontend-only; **default unchecked** → hide items under 10 MiB.

## Checkbox semantics

| Checkbox | Default | Unchecked (default) | Checked |
|----------|---------|---------------------|---------|
| `<1M` (binaries) | unchecked | Hide binaries/repos below 1 MiB | Show all sizes |
| `<10M` (worktrees) | unchecked | Hide worktrees/repos below 10 MiB | Show all sizes |

Label text exactly: `<1M` and `<10M`.

`data-testid`: `binary-show-under-1m`, `worktree-show-under-10m`

## Filter rules (frontend only, no backend change)

### Binary files

1. **File (leaf):** hide if `size < 1_048_576` (1 MiB) when filter active (checkbox unchecked).
2. **Repo (dir):** after leaf filter, compute repo **total size** = sum of remaining visible binaries; hide entire repo if total `< 1 MiB`.

### Git Worktrees

1. **Linked worktree (child):** hide if `size < 10_485_760` (10 MiB) when filter active.
2. **Repo row (main checkout):** hide if main `repo.size < 10 MiB`.
   - Linked children filtered independently; repo hidden based on **main checkout size** only (repo row is the main dir).

## Sort rules (frontend only)

Apply after filter, on every state update (each new SSE `repo`/`worktree`/`binary` event):

| Level | Sort key | Order |
|-------|----------|-------|
| Worktree repos | `repo.size` | DESC |
| Linked worktrees per repo | `hit.size` | DESC |
| Binary repos | sum of visible binary sizes in repo | DESC |
| Binaries per repo | `hit.size` | DESC |

Re-sort must be **stable enough** that a newly arrived larger item moves above smaller ones immediately.

## UI placement

Checkboxes in section header area (beside Scan/Stop or below title):

```
Git Worktrees          ☐ <10M   [Scan] [Stop]
Binary files           ☐ <1M    [Scan] [Stop]
```

## Implementation approach

Extract testable pure functions (recommended):

`disk-usage-analyser-react/src/repositoryScansLayout.ts`

```ts
export const ONE_MIB = 1048576;
export const TEN_MIB = 10485760;

export function sortWorktreeRepos(repos: WorktreeRepoRow[]): WorktreeRepoRow[];
export function sortLinkedWorktrees(hits: WorktreeHit[]): WorktreeHit[];
export function filterWorktreeRepos(repos, linkedByRepo, showUnder10M: boolean): { repos, linkedByRepo };
export function sortBinaryRepos(byRepo: Map<string, BinaryHit[]>): [string, BinaryHit[]][];
export function filterBinaryRepos(byRepo, showUnder1M: boolean): Map<string, BinaryHit[]>;
```

`TmpFilesAnalyse.tsx` uses these on render; checkbox state in React.

## Test plan

### 1. `tests/repository-scans-layout/` (new) — pure logic

| Leaf | Expected |
|------|----------|
| `sort/worktree-repos-desc` | Larger repo first |
| `sort/binary-repos-desc` | Larger repo total first |
| `sort/children-desc` | Children sorted DESC within repo |
| `filter/binaries-hide-under-1m` | File <1M hidden; repo hidden if total <1M |
| `filter/binaries-show-when-checked` | All sizes visible when `showUnder1M=true` |
| `filter/worktrees-hide-under-10m` | Linked <10M hidden; repo with main <10M hidden |
| `filter/worktrees-show-when-checked` | All visible when checked |
| `resort/on-new-larger-item` | Inserting larger item moves it to top |

Harness: `Run(t, req)` imports layout functions (via small Go test wrapper or duplicate constants — prefer testing TS via node script in SETUP).

Actually for doctest in this project, frontend logic tests often use Playwright. For pure functions, use Go doctest that imports... can't import TS.

**Use Go duplicate of logic in test SETUP** OR run node in SETUP. Check project patterns.

Simpler: implement logic in Go-testable form in `repositoryScansLayout.ts` and test via Playwright for integration + a new tree where ASSERT runs embedded TypeScript via node in SETUP helper.

**Designer choice:** Either `tests/repository-scans-layout/` with Node harness calling compiled helpers, or extend frontend tests only. Prefer **extracted TS module + frontend playwright** for integration; add **layout unit tests** if harness exists.

### 2. Extend `tests/tmp-analyse-frontend-test-cases/`

| Leaf | Expected |
|------|----------|
| `binaries-section/filter-under-1m-default` | Checkbox unchecked; no `binary-row` with size text showing values clearly under 1M if any exist |
| `binaries-section/filter-show-under-1m` | Checking box increases visible row count (when small binaries exist) |
| `binaries-section/sort-by-size-desc` | Repo row sizes in DOM are monotonic non-increasing |
| `worktrees-section/filter-under-10m-default` | Checkbox unchecked by default |
| `worktrees-section/filter-show-under-10m` | Toggle shows more rows when small worktrees exist |
| `worktrees-section/sort-by-size-desc` | Repo sizes monotonic non-increasing |
| `worktrees-live-stream` (update) | Order may change as rows stream (optional assert) |

Label `slow, ui-automation` where scan required.

## Out of scope

- Backend SSE order unchanged
- Persisting checkbox state across page reloads

## Verify

```sh
doctest test ./tests/repository-scans-layout   # if created
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/...
doctest test ./...
```