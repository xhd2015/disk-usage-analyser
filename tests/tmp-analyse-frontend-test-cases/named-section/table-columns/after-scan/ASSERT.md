---
label: slow, ui-automation
explanation: Playwright scan wait; compile and link dev server
---

## Expected

- `node-modules-tree` is present after scan completes.
- When node_modules dirs exist: column header row is present.
- At least one `node-modules-pkgmgr-*` and one `node-modules-shared-*` cell with non-empty text.
- When no node_modules: `node-modules-empty-state` is shown instead.

## Errors

- Tree element must not be MISSING.
- Column cells must not be MISSING when rows exist.

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
	if strings.Contains(resp.Output, "ELEM node-modules-shared: MISSING") {
		t.Fatal("node_modules shared size column cell missing")
	}
	for _, token := range []string{"PKG_MGR_VALUE:", "SHARED_VALUE:"} {
		if !strings.Contains(resp.Output, token) {
			t.Fatalf("expected %s in output", token)
		}
	}
}
```