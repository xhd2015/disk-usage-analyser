---
label: slow, ui-automation
explanation: Playwright scan + path cell inspection; compile and link dev server
---

## Expected

- After scan, when node_modules rows exist with a long path (`data-full-path` length > 40):
  - Visible path text starts with `…` or `...` (prefix hidden).
  - Visible path text ends with `node_modules` (suffix fully visible).
  - `data-full-path` attribute on the path span contains the complete untruncated path including `node_modules`.
- When no node_modules: `node-modules-empty-state` is shown instead.
- When all paths are short (fit without truncation): leaf skips gracefully.

## Errors

- Path cell must not be MISSING when rows exist.
- Prefix ellipsis, suffix, or full-path attribute checks must not report FAIL.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		if strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
			return
		}
		t.Skip("no node_modules dirs on machine and no empty state yet")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-path: MISSING") {
		t.Fatal("node_modules path cell missing")
	}
	if strings.Contains(resp.Output, "PATH_LONG: none") {
		t.Skip("no long node_modules path on machine to verify prefix truncation")
	}
	if strings.Contains(resp.Output, "PATH_PREFIX_ELLIPSIS: fail") {
		t.Fatal("visible path does not start with ellipsis prefix")
	}
	if strings.Contains(resp.Output, "PATH_SUFFIX_NODE_MODULES: fail") {
		t.Fatal("visible path does not end with node_modules suffix")
	}
	if strings.Contains(resp.Output, "PATH_FULL_ATTR: fail") {
		t.Fatal("data-full-path missing or incomplete")
	}
	for _, token := range []string{"PATH_VISIBLE_TEXT:", "PATH_FULL_ATTR_TEXT:"} {
		if !strings.Contains(resp.Output, token) {
			t.Fatalf("expected %s in output", token)
		}
	}
}
```