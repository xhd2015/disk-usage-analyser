---
label: slow, ui-automation
---

## Expected

- Checking `binary-show-under-1m` increases visible `binary-row` count when small binaries exist.
- Script logs `CHECK binary-filter-toggle: pass` when count grows, or skips when no sub-1M binaries on machine.

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
	if strings.Contains(resp.Output, "BINARIES_EMPTY_STATE: present") {
		t.Skip("no binaries on machine")
	}
	if strings.Contains(resp.Output, "ELEM binary-show-under-1m: MISSING") {
		t.Fatal("binary-show-under-1m checkbox missing")
	}
	if strings.Contains(resp.Output, "SKIP no-sub-1m-binaries") {
		t.Skip("no sub-1M binaries on machine to toggle")
	}
	if !strings.Contains(resp.Output, "CHECK binary-filter-toggle: pass") {
		t.Fatal("checking binary-show-under-1m should show more rows when small binaries exist")
	}
}
```