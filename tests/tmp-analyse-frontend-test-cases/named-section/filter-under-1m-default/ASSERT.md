---
label: slow, ui-automation
---

## Expected

- `node-modules-show-under-1m` checkbox exists and is **unchecked** by default.
- After scan, no visible `node-modules-row` has parsed size strictly under 1 MiB (1048576 bytes).
- Script logs `CHECK named-filter-default: pass`.

## Errors

- SKIP when no node_modules found on machine.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
		t.Skip("no node_modules on machine")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-show-under-1m: MISSING") {
		t.Fatal("node-modules-show-under-1m checkbox missing from section header")
	}
	if !strings.Contains(resp.Output, "CHECKBOX node-modules-show-under-1m: unchecked") {
		t.Fatal("node-modules-show-under-1m should be unchecked by default")
	}
	if strings.Contains(resp.Output, "UNDER_1M_VISIBLE: yes") {
		t.Fatal("node_modules under 1 MiB must be hidden when checkbox is unchecked")
	}
	if !strings.Contains(resp.Output, "CHECK named-filter-default: pass") {
		t.Fatal("expected named filter default check to pass")
	}
}
```
