---
label: slow, ui-automation
---

## Expected

- `worktree-show-under-10m` checkbox exists and is **unchecked** by default.
- No visible `worktree-repo-size` or `worktree-row-size` parses strictly under 10 MiB.
- Script logs `CHECK worktree-filter-default: pass`.

## Errors

- SKIP when no worktree repo rows on machine.

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
	if strings.Contains(resp.Output, "ELEM worktree-show-under-10m: MISSING") {
		t.Fatal("worktree-show-under-10m checkbox missing from section header")
	}
	if !strings.Contains(resp.Output, "CHECKBOX worktree-show-under-10m: unchecked") {
		t.Fatal("worktree-show-under-10m should be unchecked by default")
	}
	if strings.Contains(resp.Output, "UNDER_10M_VISIBLE: yes") {
		t.Fatal("worktrees under 10 MiB must be hidden when checkbox is unchecked")
	}
	if !strings.Contains(resp.Output, "CHECK worktree-filter-default: pass") {
		t.Fatal("expected worktree filter default check to pass")
	}
}
```