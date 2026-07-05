---
label: slow, ui-automation
explanation: Playwright hover + AX poll; compile and link dev server
---

## Expected

- After scan, when node_modules rows exist:
  - First path cell has truncation styling (prefix ellipsis display or `overflow: hidden` / `text-overflow: ellipsis`).
  - `data-full-path` on the path span holds the complete untruncated path (length > 20).
  - Visible `textContent` may be prefix-truncated and need not equal the full path.
  - Hovering the path shows an Ant Design tooltip containing the full path from `data-full-path`.
- When no node_modules: `node-modules-empty-state` is shown instead.

## Errors

- Path cell must not be MISSING when rows exist.
- Full-path attribute, ellipsis, or tooltip checks must not report FAIL.

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
	if strings.Contains(resp.Output, "PATH_ELLIPSIS: fail") {
		t.Fatal("path cell missing truncation styling")
	}
	if strings.Contains(resp.Output, "PATH_FULL_ATTR: fail") {
		t.Fatal("data-full-path missing or too short")
	}
	if strings.Contains(resp.Output, "PATH_TOOLTIP: fail") {
		t.Fatal("tooltip did not show full path on hover")
	}
	for _, token := range []string{"PATH_VISIBLE_TEXT:", "PATH_FULL_ATTR_TEXT:"} {
		if !strings.Contains(resp.Output, token) {
			t.Fatalf("expected %s in output", token)
		}
	}
}
```