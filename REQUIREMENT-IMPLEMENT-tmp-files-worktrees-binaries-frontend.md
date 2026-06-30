# REQUIREMENT-IMPLEMENT — Git Worktrees & Binary Files on Tmp Files page

## Context

Tmp Files Analyse page (`/tmp-analyse`) needs a new **Repository Scans** section group after Developer Tools with:

1. **Git Worktrees** — independent SSE scan, tree UI (repo → worktrees with sizes)
2. **Binary files** — independent SSE scan, tree UI grouped by `repoPath`, select/delete with total size bar

CLI binary scanning already exists in `tmpfiles/tmpfiles.go`. Worktree metadata exists in `scan_repo` with `ListWorktrees=true`.

## Confirmed decisions

- Placement: Repository Scans after Developer Tools
- Binary tree: group by `repoPath`
- Binary ops: leaf + repo-parent checkboxes, selection bar, confirm modal, delete
- Delete scope: only paths from current scan session (server validates)
- Scan root: hidden, always `~`
- Worktree tree: includes main checkout (`isMain=true`)
- Platform: all platforms (no macOS gate for these scans)
- Worktrees: read-only (no delete)

## Tests are sealed — do not modify

Sealed trees (git added):

```
tests/tmp-worktrees-scan/          (9 leaves)
tests/tmp-binaries-scan/           (8 leaves)
tests/tmp-binaries-delete/         (7 leaves)
tests/tmp-analyse-frontend-test-cases/  (8 new leaves)
```

**Do not modify any file under these test directories.**

## Implementation tasks

### Backend

1. **`server/tmp_worktrees.go`**
   - `WorktreeHit`, `WorktreeScanSummary` types
   - `HandleTmpWorktreesScan` — SSE: `repo`, `worktree`, `summary`, `done`, `server_error`
   - `scan_repo.Scan` with `ListWorktrees=true`, default root `~`, size each worktree checkout
   - Stream each worktree immediately; flush after each event
   - Context cancellation on client disconnect
   - `parseWorktreesSSE` helpers if referenced by tests (may be in test tree only — check DOCTEST.md)

2. **`server/tmp_binaries.go`**
   - `HandleTmpBinariesScan` — SSE: `binary`, `summary`, `done`, `server_error`
   - Reuse `tmpfiles` scan logic (extract shared function if needed)
   - Default root `~`, stream each `BinaryHit` immediately

3. **`server/tmp_binaries_delete.go`** (or same file)
   - `DeleteBinariesRequest`, `DeleteBinariesResult`, `DeleteFailure`
   - `HandleTmpBinariesDelete` — POST JSON
   - Validate: path in current scan session, still a binary (`classifyFile`), regular file, not CloudStorage
   - Partial delete OK

4. **`server/server.go`** — register routes:
   - `GET /api/tmp-worktrees-scan`
   - `GET /api/tmp-binaries-scan`
   - `POST /api/tmp-binaries-delete`

5. **`tmpfiles/tmpfiles.go`** — optional refactor to export scan callback for server reuse

### Frontend (`disk-usage-analyser-react/src/TmpFilesAnalyse.tsx`)

- Add **Repository Scans** section group after Developer Tools
- **Git Worktrees** subsection: scan/stop buttons, Ant Design Tree, SSE client
- **Binary files** subsection: scan/stop buttons, Tree with checkboxes, selection bar, delete modal
- `data-testid` hooks per tests:
  - `section-repository-scans-heading`
  - `worktrees-section`, `worktrees-scan-btn`, `worktrees-stop-btn`
  - `binaries-section`, `binaries-scan-btn`, `binaries-stop-btn`
  - `binary-selected-total`, `binary-delete-btn`, `binary-delete-confirm`
  - tree row testids for worktrees and binaries
- Independent EventSource refs; scans don't block page-level cache scan

## Verify commands

```sh
doctest vet ./tests/tmp-worktrees-scan
doctest vet ./tests/tmp-binaries-scan
doctest vet ./tests/tmp-binaries-delete
doctest vet ./tests/tmp-analyse-frontend-test-cases

doctest test ./tests/tmp-worktrees-scan
doctest test ./tests/tmp-binaries-scan
doctest test ./tests/tmp-binaries-delete
doctest test ./tests/tmp-analyse-frontend-test-cases/repository-scans/renders
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/repository-scans/...
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/worktrees-section/...
doctest test --label 'slow && ui-automation' ./tests/tmp-analyse-frontend-test-cases/binaries-section/...
doctest test ./tests/tmp-analyse-frontend-test-cases/independent-scan-controls

doctest test ./...   # no regressions
```

## Reference

- Requirement: `REQUIREMENT-DESIGN-tmp-files-worktrees-binaries-frontend.md`
- Existing SSE pattern: `server/tmp_analyse.go`
- Existing CLI scan: `tmpfiles/tmpfiles.go`
- Worktree enrichment: `dot-pkgs-scan/go-pkgs/git/scan_repo/`