---
label: slow, ui-automation
explanation: Playwright hover, click copy, clipboard read; compile and link dev server
---

## Expected

- After scan, when node_modules rows exist:
  - Hovering the first path cell shows a tooltip with `data-testid` copy button (`node-modules-path-copy-*`).
  - Clicking the copy button writes the full path (from `data-full-path`) to the clipboard.
- When no node_modules: `node-modules-empty-state` is shown instead.

## Errors

- Path cell or copy button must not be MISSING when rows exist.
- Clipboard check must not report FAIL.

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
	if strings.Contains(resp.Output, "ELEM node-modules-path-copy: MISSING") {
		t.Fatal("tooltip copy button missing")
	}
	if strings.Contains(resp.Output, "PATH_CLIPBOARD: fail") {
		t.Fatal("clipboard does not contain full path after copy click")
	}
	for _, token := range []string{"PATH_FULL_ATTR_TEXT:", "PATH_CLIPBOARD_TEXT:"} {
		if !strings.Contains(resp.Output, token) {
			t.Fatalf("expected %s in output", token)
		}
	}
}
```