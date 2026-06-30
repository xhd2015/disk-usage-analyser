---
label: slow, ui-automation
---

## Expected

- Clicking `binary-repo-checkbox-*` checks all `binary-checkbox-*` leaves in that repo.
- `binary-selected-total` shows non-empty selected size.
- Script logs `CHECK repo-select-all: true`.

## Errors

- SKIP when fewer than two binary leaves exist in a repo row.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP repo-select-all") {
		t.Skip("insufficient binary fixtures for repo select-all test")
	}
	if !strings.Contains(resp.Output, "CHECK repo-select-all: true") {
		t.Fatalf("repo checkbox did not select all leaves\n%s", resp.Output)
	}
}
```