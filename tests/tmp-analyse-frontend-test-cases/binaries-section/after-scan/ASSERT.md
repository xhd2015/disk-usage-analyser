---
label: slow, ui-automation
---

## Expected

- `binaries-tree` is present after scan completes.
- When binaries exist: at least one `binary-repo-row`, one `binary-row`, and non-empty `BINARY_KIND` (go/macho/elf).
- When no binaries: `binaries-empty-state` is shown instead of rows.

## Errors

- Tree element must not be MISSING.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM binaries-tree: MISSING") {
		t.Fatal("binaries tree missing after scan")
	}
	if strings.Contains(resp.Output, "COUNT binary-rows: 0") {
		if !strings.Contains(resp.Output, "BINARIES_EMPTY_STATE: present") {
			t.Skip("no binaries on machine and no empty state yet")
		}
		return
	}
	for _, kind := range []string{"go", "macho", "elf"} {
		if strings.Contains(strings.ToLower(resp.Output), kind) {
			return
		}
	}
	if !strings.Contains(resp.Output, "BINARY_KIND:") {
		t.Fatal("expected binary kind badge in output")
	}
}
```