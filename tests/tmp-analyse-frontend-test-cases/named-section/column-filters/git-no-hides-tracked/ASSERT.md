---
label: slow, ui-automation
explanation: Playwright scan wait and filter interaction
---

## Expected

- When both tracked and untracked rows exist, selecting Git=No removes all tracked rows.
- Script logs `CHECK git-no-filter: pass`.
- Skips when machine has only one git state or no node_modules.

## Errors

- `node-modules-filter-git` missing when rows exist.
- Tracked rows remain visible after Git=No.

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
	if strings.Contains(resp.Output, "SKIP no-git-mix") {
		t.Skip("machine lacks both tracked and untracked node_modules rows")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-filter-git: MISSING") {
		t.Fatal("node-modules-filter-git control missing")
	}
	if strings.Contains(resp.Output, "CHECK git-no-filter: pass") {
		return
	}
	t.Fatalf("Git=No should hide tracked rows when mixed git states exist\nOutput:\n%s", resp.Output)
}
```