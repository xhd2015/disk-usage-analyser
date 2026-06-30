---
label: slow, ui-automation
---

## Expected

- Scanning badge appears during worktrees scan.
- Repo row count reaches at least 1 before done badge (repo events stream first).
- Script logs `CHECK worktrees-scanning-seen: true` and `CHECK worktree-repo-growth: true`.
- When ≥2 repos stream, final `worktree-repo-size` order is monotonic non-increasing (size DESC).

## Errors

- SKIP when no git repos on machine (maxRepoRows stays 0).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "MAX worktree-repo-rows: 0") {
		t.Skip("no worktree repo rows streamed on this machine")
	}
	if !strings.Contains(resp.Output, "CHECK worktrees-scanning-seen: true") {
		t.Fatal("expected scanning badge during worktrees scan")
	}
	if !strings.Contains(resp.Output, "CHECK worktree-repo-growth: true") {
		t.Fatal("expected repo rows to stream during worktrees scan")
	}
	if strings.Contains(resp.Output, "FINAL_REPO_SIZES_BYTES: [") &&
		!strings.Contains(resp.Output, "FINAL_REPO_SIZES_BYTES: []") &&
		strings.Count(resp.Output, "FINAL_REPO_SIZES_BYTES:") > 0 {
		// When multiple repos present, final DOM order must be size DESC.
		if strings.Contains(resp.Output, "SORT worktree-repo-sizes-final: not-desc") {
			t.Fatal("streaming worktree repos should end sorted by size DESC")
		}
	}
}
```