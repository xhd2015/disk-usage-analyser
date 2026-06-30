---
label: slow, ui-automation
---

## Expected

- `binaries-tree` is present after scan completes.
- Computed `text-align` on `binaries-tree` is `left`.

## Errors

- SKIP when binaries tree is missing (no scan results on machine).

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
		t.Skip("binaries tree missing on this machine")
	}
	alignLine := findLine(resp.Output, "TREE_STYLE text-align:")
	if !strings.Contains(alignLine, "left") {
		t.Fatalf("expected binaries-tree text-align left, got: %s", alignLine)
	}
}

func findLine(output, prefix string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return line
		}
	}
	return ""
}
```