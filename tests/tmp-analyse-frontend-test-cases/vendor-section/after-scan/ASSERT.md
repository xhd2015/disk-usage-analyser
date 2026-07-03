---
label: slow, ui-automation
---

## Expected

- `vendor-tree` is present after scan completes.
- When vendor dirs exist: at least one `vendor-row` with size, name, path, repo columns.
- When no vendor dirs: `vendor-empty-state` is shown instead of rows.

## Errors

- Tree element must not be MISSING.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM vendor-tree: MISSING") {
		t.Fatal("vendor tree missing after scan")
	}
	if strings.Contains(resp.Output, "COUNT vendor-rows: 0") {
		if !strings.Contains(resp.Output, "VENDOR_EMPTY_STATE: present") {
			t.Skip("no vendor dirs on machine and no empty state yet")
		}
		return
	}
	for _, token := range []string{"VENDOR_ENTITY_SIZE:", "VENDOR_ENTITY_NAME:", "VENDOR_ENTITY_PATH:", "VENDOR_ENTITY_REPO:"} {
		if !strings.Contains(resp.Output, token) {
			t.Fatalf("expected %s in output", token)
		}
	}
}
```
