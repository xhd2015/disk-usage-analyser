---
label: slow, ui-automation
---

## Expected

- `node-modules-delete-confirm-modal` appears before delete executes.
- After confirm, node_modules row count decreases by at least one.
- Script logs `CHECK named-row-removed: true`.

## Errors

- SKIP when no deletable node_modules rows exist.

## Side Effects

- Deletes a real node_modules directory from the developer machine (selected test directory).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP named-delete") {
		t.Skip("no node_modules rows on this machine")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-delete-confirm-modal: MISSING") {
		t.Fatal("confirmation modal must appear before delete")
	}
	if !strings.Contains(resp.Output, "CHECK named-row-removed: true") {
		t.Fatalf("expected deleted node_modules removed from tree\n%s", resp.Output)
	}
}
```
