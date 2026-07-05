---
label: slow, ui-automation
explanation: Playwright scan wait; compile and link dev server
---

## Expected

- After scan, `node-modules-filter-git`, `node-modules-filter-package-json`, and `node-modules-filter-pm` are present.
- Script logs `FILTER_DEFAULT: git=all packageJson=all pm=all`.
- When no node_modules on machine, empty state is acceptable.

## Errors

- Any filter control missing when rows exist.
- Default selection is not `all` for any control.

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
		return
	}
	if strings.Contains(resp.Output, "ELEM node-modules-filter-git: MISSING") {
		t.Fatal("node-modules-filter-git control missing")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-filter-package-json: MISSING") {
		t.Fatal("node-modules-filter-package-json control missing")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-filter-pm: MISSING") {
		t.Fatal("node-modules-filter-pm control missing")
	}
	if !strings.Contains(resp.Output, "FILTER_DEFAULT: git=all packageJson=all pm=all") {
		t.Fatalf("expected default filter selections logged as all\nOutput:\n%s", resp.Output)
	}
}
```