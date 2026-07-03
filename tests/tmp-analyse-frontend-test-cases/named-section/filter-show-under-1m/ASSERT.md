---
label: slow, ui-automation
---

## Expected

- Checking `node-modules-show-under-1m` increases visible `node-modules-row` count when small hits exist.
- Script logs `CHECK named-filter-toggle: pass` when count grows, or skips when no sub-1M hits on machine.

## Errors

- Checkbox must exist.

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
		t.Fatal("node-modules-show-under-1m checkbox missing")
	}
	if strings.Contains(resp.Output, "SKIP no-sub-1m-named") {
		t.Skip("no sub-1M node_modules on machine to toggle")
	}
	if !strings.Contains(resp.Output, "CHECK named-filter-toggle: pass") {
		t.Fatal("checking node-modules-show-under-1m should show more rows when small hits exist")
	}
}
```
