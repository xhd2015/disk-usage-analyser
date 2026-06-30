# Scenario

**Feature**: tmp-worktrees-scan shared fixture harness

```
GET /api/tmp-worktrees-scan -> scan_repo discovery -> git worktree list -> size checkout -> SSE repo/worktree/summary/done
```

## Preconditions

- Tests must never scan the real user home directory.
- Each leaf creates a temporary fake home under `t.TempDir()` and sets `HOME` before httptest.
- Real `git` is required for worktree fixtures; leaves skip when `git` is not on PATH.
- Git repositories use `git init`, an initial commit, and `git worktree add` for linked checkouts.

## Steps

1. Create a fresh fixture home for each leaf.
2. Set `req.HomeDir` to the fake home path.
3. Set `req.Op` to the handler operation for the leaf scenario.
4. Run the httptest SSE harness.

## Context

- Scan root is always `~` (the injected `HOME`).
- Main checkout size is carried on the `repo` SSE event (`size`, `fileCount`).
- Only linked worktrees appear as `worktree` SSE events (`isMain=false`).
- Main-only repos emit a `repo` event and zero `worktree` events.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.HomeDir = filepath.Join(t.TempDir(), "home")
	if err := os.MkdirAll(req.HomeDir, 0755); err != nil {
		return err
	}
	if req.Op == "" {
		req.Op = "worktrees-scan"
	}
	return nil
}

func gitAvailable(t *testing.T) bool {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available on PATH")
		return false
	}
	return true
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}

func gitInitRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func gitInitialCommit(t *testing.T, dir string) {
	t.Helper()
	readme := filepath.Join(dir, "README")
	if err := os.WriteFile(readme, []byte("init\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README")
	runGit(t, dir, "commit", "-m", "init")
}

func gitWorktreeAdd(t *testing.T, mainDir, wtDir, branch string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(wtDir), 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, mainDir, "worktree", "add", "-b", branch, wtDir)
}

func writeFile(t *testing.T, base string, rel string, data []byte) {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func repoUnderHome(t *testing.T, home string, rel string) string {
	t.Helper()
	dir := filepath.Join(home, filepath.FromSlash(rel))
	gitInitRepo(t, dir)
	return dir
}
```