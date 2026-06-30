---
label: slow, ui-automation
---

## Expected

- Checking `worktree-show-under-10m` increases visible repo and/or worktree row count when small worktrees exist.
- Script logs `CHECK worktree-filter-toggle: pass`, or skips when no sub-10M items on machine.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "COUNT worktree-repo-rows-before: 0") {
		t.Skip("no worktree repos on machine")
	}
	if strings.Contains(resp.Output, "ELEM worktree-show-under-10m: MISSING") {
		t.Fatal("worktree-show-under-10m checkbox missing")
	}
	if strings.Contains(resp.Output, "SKIP no-sub-10m-worktrees") {
		t.Skip("no sub-10M worktrees on machine to toggle")
	}
	if !strings.Contains(resp.Output, "CHECK worktree-filter-toggle: pass") {
		t.Fatal("checking worktree-show-under-10m should show more rows when small worktrees exist")
	}
}
```