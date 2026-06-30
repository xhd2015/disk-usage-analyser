---
label: slow, ui-automation
---

## Expected

- Toggling a leaf checkbox changes `binary-selected-total` text.
- Selected total includes "to clear" with a human size.
- `binary-delete-btn` becomes visible when selection is non-empty.

## Errors

- SKIP when no binary rows after scan.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP binaries-select") {
		t.Skip("no binary rows on this machine")
	}
	if !strings.Contains(resp.Output, "CHECK selected-total-updated: true") {
		t.Fatalf("selected total did not update\n%s", resp.Output)
	}
	if !strings.Contains(resp.Output, "ELEM binary-delete-btn: visible") {
		t.Fatal("delete button should be visible when items selected")
	}
}
```