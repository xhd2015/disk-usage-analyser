---
label: slow, ui-automation
---

## Expected

- Repo group totals in DOM order are monotonic non-increasing (largest repo first).
- Within each repo, binary row sizes are monotonic non-increasing.
- Script logs `CHECK binary-sort-desc: pass`.

## Errors

- SKIP when fewer than 2 binary rows (cannot verify sort order).

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
	if strings.Contains(resp.Output, "SKIP insufficient-binary-rows") {
		t.Skip("need at least 2 binary rows to verify sort order")
	}
	if strings.Contains(resp.Output, "SORT binary-repo-totals: not-desc") {
		t.Fatal("binary repo totals must be sorted DESC")
	}
	if strings.Contains(resp.Output, "SORT binary-child-sizes: not-desc") {
		t.Fatal("binary rows within repos must be sorted DESC")
	}
	if !strings.Contains(resp.Output, "CHECK binary-sort-desc: pass") {
		t.Fatal("expected binary sort DESC check to pass")
	}
}
```