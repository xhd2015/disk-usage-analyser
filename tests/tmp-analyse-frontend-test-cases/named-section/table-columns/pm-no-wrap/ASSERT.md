---
label: slow, ui-automation
explanation: Playwright scan wait; compile and link dev server
---

## Expected

- `node-modules-tree` and `node-modules-column-header` are present after scan completes.
- When node_modules rows exist: every PM cell has computed `white-space: nowrap` (`PM_NOWRAP: ok`).
- When no node_modules: `node-modules-empty-state` is shown instead.

## Errors

- PM cells must not be MISSING when rows exist.
- `PM_NOWRAP: fail` indicates text wrapping in the PM column.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM node-modules-tree: MISSING") {
		t.Fatal("node_modules tree missing after scan")
	}
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		if strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
			return
		}
		t.Skip("no node_modules dirs on machine and no empty state yet")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-column-header: MISSING") {
		t.Fatal("node_modules column header row missing")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-pkgmgr: MISSING") {
		t.Fatal("node_modules package manager column cell missing")
	}
	if strings.Contains(resp.Output, "PM_NOWRAP: fail") {
		t.Fatalf("PM column cells do not have white-space: nowrap\nOutput:\n%s", resp.Output)
	}
	if !strings.Contains(resp.Output, "PM_NOWRAP: ok") {
		t.Fatalf("expected PM_NOWRAP: ok in output\n%s", resp.Output)
	}
}
```