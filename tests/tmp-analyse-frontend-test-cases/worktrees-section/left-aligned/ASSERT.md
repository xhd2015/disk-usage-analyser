---
label: slow, ui-automation
---

## Expected

- `worktrees-tree` is present after scan completes.
- Computed `text-align` on `worktrees-tree` is `left`.

## Errors

- SKIP when worktrees tree is missing (no git repos on machine).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM worktrees-tree: MISSING") {
		t.Skip("worktrees tree missing on this machine")
	}
	alignLine := findLine(resp.Output, "TREE_STYLE text-align:")
	if !strings.Contains(alignLine, "left") {
		t.Fatalf("expected worktrees-tree text-align left, got: %s", alignLine)
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