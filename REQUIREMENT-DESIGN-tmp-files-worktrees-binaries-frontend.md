# Requirement — Git Worktrees & Binary Files sections on Tmp Files page

## Summary

Extend the **Tmp Files Analyse** frontend page (`/tmp-analyse`) with a new **Repository Scans** section group (after Developer Tools) containing two independent subsections:

1. **Git Worktrees** — discover git repos, list worktrees per main repo, show on-disk size per worktree checkout; tree UI (repo → worktrees, main checkout included).
2. **Binary files** — scan git repos for Go/Mach-O/ELF binaries (same logic as `disk-usage-analyser tmp-files scan`); tree UI grouped by `repoPath`; user can **select** binaries and **delete** them, with a running **total size to clear** for the selection.

Each subsection has its own **Start / Stop Scan** button beside the section title. Scans run independently of the existing page-level cache-location scan and of each other. Results stream to the frontend as soon as items are found (SSE).

---

## Confirmed Decisions

| # | Decision |
|---|----------|
| Q1 | New **"Repository Scans"** group after Developer Tools |
| Q2 | Binary tree grouped by **`repoPath`**; selectable binaries with delete + selected-total size bar |
| Q3 | Scan root hidden in v1 — always scan **`~`** |
| Q4 | Worktree tree includes **main checkout** as first child (`isMain=true`) |
| Q5 | **All platforms** — no macOS-only gate for these scans |

---

## Current State

| Layer | Exists today |
|-------|--------------|
| CLI `tmp-files scan` | ✅ Binary discovery in git repos (`tmpfiles/tmpfiles.go`) |
| `scan_repo` worktree metadata | ✅ `ListWorktrees` + `git worktree list --porcelain` |
| Frontend Tmp Files page | ✅ Cache/temp location cards + single page-level scan |
| Binary delete API | ❌ |
| API for binaries/worktrees scan | ❌ |
| Frontend worktree/binary UI | ❌ |

---

## Data Models

### Backend (Go)

```go
// Worktree scan — streamed per worktree as sizing completes
type WorktreeHit struct {
    RepoPath   string `json:"repoPath"`
    RepoName   string `json:"repoName"`
    Path       string `json:"path"`       // worktree checkout path (~/...)
    Head       string `json:"head"`       // branch or detached SHA
    IsMain     bool   `json:"isMain"`
    Size       int64  `json:"size"`
    SizeHuman  string `json:"sizeHuman"`
    FileCount  int64  `json:"fileCount"`
}

// Binary scan — reuses existing BinaryHit from tmpfiles
type BinaryHit struct {
    Path      string `json:"path"`
    Size      int64  `json:"size"`
    SizeHuman string `json:"sizeHuman"`
    Kind      string `json:"kind"`      // go | macho | elf
    TypeDesc  string `json:"typeDesc"`
    RepoPath  string `json:"repoPath"`
    RepoName  string `json:"repoName"`
}

type WorktreeScanSummary struct {
    Repos, Worktrees int
    TotalSize        int64
    TotalHuman       string
}

type BinaryScanSummary struct {
    Repos, Binaries int
    TotalSize       int64
    TotalHuman      string
}

// Binary delete
type DeleteBinariesRequest struct {
    Paths []string `json:"paths"` // absolute or ~/ paths from scan results
}

type DeleteBinariesResult struct {
    Deleted     []string `json:"deleted"`
    Failed      []DeleteFailure `json:"failed"`
    FreedSize   int64    `json:"freedSize"`
    FreedHuman  string   `json:"freedHuman"`
}

type DeleteFailure struct {
    Path  string `json:"path"`
    Error string `json:"error"`
}
```

### Frontend (TypeScript)

```ts
interface BinarySelection {
    selectedPaths: Set<string>;       // binary file paths
    selectedTotalSize: number;      // sum of selected BinaryHit.size
}
```

Tree nodes:

| Section | Tree shape |
|---------|------------|
| Git Worktrees | `repo` parent → children = worktree rows (main first) |
| Binary files | `repoPath` parent (checkbox: select all in repo) → children = binary leaves (individual checkbox) |

### Storage

- Scan results: in-memory React state; cleared on new scan start.
- Selection state: in-memory; cleared on new scan start.
- No persistence.

---

## API

### SSE scan endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /api/tmp-worktrees-scan` | Discover repos + worktrees, size each checkout |
| `GET /api/tmp-binaries-scan` | Binary discovery (wraps `tmpfiles` logic) |

Query params hidden in v1 (server defaults): `root=~`, `maxDepth=0`.

#### Worktrees SSE events

| Event | Payload |
|-------|---------|
| `repo` | `{repoPath, repoName}` |
| `worktree` | `WorktreeHit` — stream immediately |
| `summary` | `WorktreeScanSummary` |
| `done` | `{status:"complete"}` |
| `server_error` | `{error}` |

#### Binaries SSE events

| Event | Payload |
|-------|---------|
| `binary` | `BinaryHit` — stream immediately |
| `summary` | `BinaryScanSummary` |
| `done` | `{status:"complete"}` |
| `server_error` | `{error}` |

#### Cancellation

Frontend closes `EventSource` → server aborts via `r.Context().Done()`.

### Binary delete endpoint (new)

| Endpoint | Purpose |
|----------|---------|
| `POST /api/tmp-binaries-delete` | Delete selected binary files |

**Request:** `DeleteBinariesRequest` JSON body.

**Response:** `DeleteBinariesResult` JSON.

**Server-side safety rules (proposed):**

1. Each path must resolve to a regular file (not directory, not symlink to outside).
2. Each path must lie under a discovered git repo root from a prior scan session **OR** simply under user home — TBD (see open questions).
3. Re-validate file is still a binary (same `classifyFile` logic) before delete — prevents arbitrary file deletion.
4. Return per-path success/failure; partial delete is OK.

**Platform:** all platforms (no macOS gate).

**CloudStorage:** reject or skip paths on remote-backed filesystems (same `remotefs` guard).

---

## Frontend UI

### Section layout (after Developer Tools)

```
── Repository Scans ──────────────────────────────────────────

┌─ Git Worktrees ──────────────────────── [▶ Scan] [■ Stop] ─┐
│  ▼ ~/Projects/foo          1.2 GB                          │
│    ├─ main (master)        800 MB                            │
│    └─ ../foo-wt (feature)  400 MB                            │
└──────────────────────────────────────────────────────────────┘

┌─ Binary files ───────────────────────── [▶ Scan] [■ Stop] ─┐
│  Selected: 45 MB to clear          [Delete Selected]       │
│  ▼ ☐ ~/Projects/foo/          45 MB                        │
│    ├─ ☑ bin/app    go   12 MB                               │
│    └─ ☑ cmd/x      macho 33 MB                              │
└──────────────────────────────────────────────────────────────┘
```

### Binary selection & delete UX

- **Leaf checkbox** on each binary row.
- **Parent checkbox** on each `repoPath` node — toggles all binaries in that repo (indeterminate when partially selected).
- **Selection bar** (visible when `selectedPaths.size > 0`):
  - Text: `X selected · Y MB to clear` (`data-testid="binary-selected-total"`)
  - **Delete Selected** button (`data-testid="binary-delete-btn"`)
- **Delete flow:**
  1. Click Delete Selected → confirmation `Modal` listing count + total size
  2. On confirm → `POST /api/tmp-binaries-delete`
  3. On success → remove deleted paths from tree; update repo aggregates; clear selection for deleted items; show brief success message
  4. On partial failure → show which paths failed; keep failed items selected

### General UX rules

- Section scan buttons independent of page-level Start/Stop Scan.
- Incremental tree updates during scan.
- Stop leaves partial results.
- Empty state after scan with zero hits.
- `data-testid` hooks for Playwright doctests.

---

## Test Plan

### 1. `tests/tmp-worktrees-scan/` — backend

| Scenario | Expected |
|----------|----------|
| SSE ordering | `worktree` before `done` |
| Main + linked worktrees | Correct paths, heads, `isMain` |
| Worktree sizing | Non-empty checkout `size > 0` |
| Client disconnect | Scan aborts |
| No repos | Zero worktrees |

### 2. `tests/tmp-binaries-scan/` — backend

| Scenario | Expected |
|----------|----------|
| SSE streaming | `binary` before `done` |
| Classification | go / macho / elf |
| Multi-repo | Distinct `repoPath` per hit |
| Ignored dirs | `vendor/` skipped |

### 3. `tests/tmp-binaries-delete/` — backend (new)

| Scenario | Expected |
|----------|----------|
| Delete single binary | File removed; `freedSize` matches |
| Delete multiple | All removed; correct total |
| Non-binary path rejected | `failed` entry, file untouched |
| Directory path rejected | `failed` entry |
| Already deleted path | `failed` with not-found |
| Partial batch | Some `deleted`, some `failed` |

### 4. `tests/tmp-analyse-frontend-test-cases/` — UI

| Leaf | Verifies |
|------|----------|
| `repository-scans/renders` | Section group + both subsections + scan buttons |
| `worktrees-section/after-scan` | Repo → worktree tree with sizes |
| `binaries-section/after-scan` | Repo → binary tree with kind badges |
| `binaries-section/select-and-total` | Checkbox toggles update selected total |
| `binaries-section/delete-selected` | Confirm modal → files removed from tree |
| `binaries-section/repo-select-all` | Parent checkbox selects all repo binaries |
| `independent-scan-controls` | Repo scans don't block cache scan |
| `worktrees-live-stream` | Row count grows during scan |

Labels: `ui-automation`, `slow` where needed.

---

## Open Questions (delete UX — need confirmation)

1. **Confirmation** — Always show a confirmation modal before delete? *(proposed: yes)*
2. **Repo-level select** — Parent checkbox selects all binaries under that repo? *(proposed: yes)*
3. **Delete scope validation** — Only allow deleting paths that appeared in the **current scan results** (safest), or any binary under `~`? *(proposed: current scan results only)*
4. **Worktree delete** — Out of scope for v1 (read-only worktree section)? *(proposed: yes, binaries only)*

---

## Approval Gate

Reply **go ahead** (with any amendments) to proceed to test design (Phase 2).

Pending your answers on the four delete UX questions above (or accept proposed defaults).