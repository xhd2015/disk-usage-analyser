---
label: slow, ui-automation
---

## Expected

- Repo group totals in DOM order are monotonic non-increasing (largest repo first).
- Within each repo, node_modules row sizes are monotonic non-increasing.
- Script logs `CHECK named-sort-desc: pass`.

## Errors

- SKIP when fewer than 2 node_modules rows (cannot verify sort order).

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
	if strings.Contains(resp.Output, "SKIP insufficient-named-rows") {
		t.Skip("need at least 2 node_modules rows to verify sort order")
	}
	if strings.Contains(resp.Output, "SORT named-repo-totals: not-desc") {
		t.Fatal("node_modules repo totals must be sorted DESC")
	}
	if strings.Contains(resp.Output, "SORT named-child-sizes: not-desc") {
		t.Fatal("node_modules rows within repos must be sorted DESC")
	}
	if !strings.Contains(resp.Output, "CHECK named-sort-desc: pass") {
		t.Fatal("expected named sort DESC check to pass")
	}
}
```
