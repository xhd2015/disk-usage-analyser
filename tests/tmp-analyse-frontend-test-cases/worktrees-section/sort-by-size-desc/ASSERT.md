---
label: slow, ui-automation
---

## Expected

- `worktree-repo-size` values in DOM order are monotonic non-increasing.
- Linked `worktree-row-size` values within each repo group are monotonic non-increasing.
- Script logs `CHECK worktree-sort-desc: pass`.

## Errors

- SKIP when fewer than 2 repo rows.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "COUNT worktree-repo-rows: 0") {
		t.Skip("no worktree repos on machine")
	}
	if strings.Contains(resp.Output, "SKIP insufficient-worktree-repos") {
		t.Skip("need at least 2 worktree repos to verify sort order")
	}
	if strings.Contains(resp.Output, "SORT worktree-repo-sizes: not-desc") {
		t.Fatal("worktree repo sizes must be sorted DESC")
	}
	if strings.Contains(resp.Output, "SORT worktree-child-sizes: not-desc") {
		t.Fatal("linked worktree rows must be sorted DESC within each repo")
	}
	if !strings.Contains(resp.Output, "CHECK worktree-sort-desc: pass") {
		t.Fatal("expected worktree sort DESC check to pass")
	}
}
```