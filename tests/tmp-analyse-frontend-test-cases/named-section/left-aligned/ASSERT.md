---
label: slow, ui-automation
---

## Expected

- `node-modules-tree` is present after scan completes.
- Computed `text-align` on `node-modules-tree` is `left`.

## Errors

- SKIP when node_modules tree is missing (no scan results on machine).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM node-modules-tree: MISSING") {
		t.Skip("node_modules tree missing on this machine")
	}
	alignLine := findLine(resp.Output, "TREE_STYLE text-align:")
	if !strings.Contains(alignLine, "left") {
		t.Fatalf("expected node-modules-tree text-align left, got: %s", alignLine)
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
