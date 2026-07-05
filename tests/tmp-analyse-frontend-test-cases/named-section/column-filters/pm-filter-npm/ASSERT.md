---
label: slow, ui-automation
explanation: Playwright scan wait and PM filter interaction
---

## Expected

- When mixed package managers exist, selecting PM=npm leaves only `npm` PM cells visible.
- Script logs `CHECK pm-filter-npm: pass`.
- Skips when machine lacks mixed PMs or has no node_modules.

## Errors

- `node-modules-filter-pm` missing when rows exist.
- Non-npm PM text visible after filter.

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
	if strings.Contains(resp.Output, "SKIP no-pm-mix") {
		t.Skip("machine lacks mixed package managers in node_modules rows")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-filter-pm: MISSING") {
		t.Fatal("node-modules-filter-pm control missing")
	}
	if strings.Contains(resp.Output, "CHECK pm-filter-npm: pass") {
		return
	}
	t.Fatalf("PM=npm should hide non-npm rows when mixed PMs exist\nOutput:\n%s", resp.Output)
}
```