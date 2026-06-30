---
label: slow, ui-automation
---

## Expected

- `binary-delete-confirm-modal` appears before delete executes.
- After confirm, binary row count decreases by at least one.
- Script logs `CHECK binary-row-removed: true`.

## Errors

- SKIP when no deletable binary rows exist.

## Side Effects

- Deletes a real binary file from the developer machine (selected test binary).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP binaries-delete") {
		t.Skip("no binary rows on this machine")
	}
	if strings.Contains(resp.Output, "ELEM binary-delete-confirm-modal: MISSING") {
		t.Fatal("confirmation modal must appear before delete")
	}
	if !strings.Contains(resp.Output, "CHECK binary-row-removed: true") {
		t.Fatalf("expected deleted binary removed from tree\n%s", resp.Output)
	}
}
```