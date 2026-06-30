---
label: slow, ui-automation
---

## Expected

- `binary-show-under-1m` checkbox exists and is **unchecked** by default.
- After scan, no visible `binary-row` has parsed size strictly under 1 MiB (1048576 bytes).
- Script logs `CHECK binary-filter-default: pass`.

## Errors

- SKIP when no binaries found on machine.

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
		t.Fatal("binary-show-under-1m checkbox missing from section header")
	}
	if !strings.Contains(resp.Output, "CHECKBOX binary-show-under-1m: unchecked") {
		t.Fatal("binary-show-under-1m should be unchecked by default")
	}
	if strings.Contains(resp.Output, "UNDER_1M_VISIBLE: yes") {
		t.Fatal("binaries under 1 MiB must be hidden when checkbox is unchecked")
	}
	if !strings.Contains(resp.Output, "CHECK binary-filter-default: pass") {
		t.Fatal("expected binary filter default check to pass")
	}
}
```