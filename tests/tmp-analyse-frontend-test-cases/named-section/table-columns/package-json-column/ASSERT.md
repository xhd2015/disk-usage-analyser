---
label: slow, ui-automation
explanation: Playwright scan wait; compile and link dev server
---

## Expected

- `node-modules-tree` and `node-modules-column-header` are present after scan completes.
- Column header text includes `package.json` (`PKGJSON_HEADER: ok`).
- When node_modules rows exist: at least one `node-modules-pkgjson-*` checkbox cell is present.
- When no node_modules: `node-modules-empty-state` is shown instead.

## Errors

- `PKGJSON_HEADER: fail` or missing pkgjson cells when rows exist.

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
	if strings.Contains(resp.Output, "PKGJSON_HEADER: fail") {
		t.Fatalf("column header missing package.json label\nOutput:\n%s", resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM node-modules-pkgjson: MISSING") {
		t.Fatal("node_modules package.json checkbox column missing")
	}
	if !strings.Contains(resp.Output, "PKGJSON_HEADER: ok") {
		t.Fatalf("expected PKGJSON_HEADER: ok in output\n%s", resp.Output)
	}
}
```