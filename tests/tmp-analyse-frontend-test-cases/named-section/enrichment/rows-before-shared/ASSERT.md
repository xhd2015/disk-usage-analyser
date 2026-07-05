---
label: slow, ui-automation
---

## Expected

- At least one `node-modules-row` appears before or when `node-modules-done-badge` shows (`scan_complete`).
- Script logs `CHECK rows-at-scan-complete: true`.
- After enrichment settles, shared column has text (`0 B` or non-zero size).
- Script logs `CHECK shared-final-present: true`.

## Errors

- SKIP when no node_modules dirs on machine (row count stays 0 with no empty state).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		if !strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
			t.Skip("no node_modules dirs on machine and no empty state yet")
		}
		return
	}
	if !strings.Contains(resp.Output, "CHECK rows-at-scan-complete: true") {
		t.Fatal("expected rows visible when scan_complete (done badge) fires")
	}
	if !strings.Contains(resp.Output, "CHECK shared-final-present: true") {
		t.Fatal("expected shared column text after enrichment settles")
	}
}
```