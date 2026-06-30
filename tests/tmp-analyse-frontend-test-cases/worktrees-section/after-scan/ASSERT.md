---
label: slow, ui-automation
---

## Expected

- `worktrees-tree` is present after scan completes.
- At least one `worktree-repo-row` exists when git repos are on the machine.
- Repo row shows non-empty size text.
- `worktree-main-badge` must **not** appear.
- Main-only repos may have zero `worktree-row` children.

## Errors

- SKIP when no git repos discovered (zero repo rows after timeout).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM worktrees-tree: MISSING") {
		t.Fatal("worktrees tree missing after scan")
	}
	if strings.Contains(resp.Output, "COUNT worktree-repo-rows: 0") {
		t.Skip("no worktree repo rows found on this machine")
	}
	if strings.Contains(resp.Output, `WORKTREE_REPO_SIZE: ""`) {
		t.Fatal("expected non-empty repo row size")
	}
	if strings.Contains(resp.Output, "ELEM worktree-main-badge: present") {
		t.Fatal("worktree-main-badge must not appear after omit-main change")
	}
}
```