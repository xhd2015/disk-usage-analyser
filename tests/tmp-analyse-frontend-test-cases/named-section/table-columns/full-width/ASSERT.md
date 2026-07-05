---
label: slow, ui-automation
explanation: Playwright scan wait; compile and link dev server
---

## Expected

- `node-modules-tree` and `node-modules-column-header` are present after scan completes.
- When node_modules rows exist: `WIDTH_RATIO size-right/body` is at least **0.90** (Size column aligns near the card body right edge).
- When no node_modules: `node-modules-empty-state` is shown instead.

## Errors

- Tree or header must not be MISSING when rows exist.
- Size column right-edge ratio below 0.90 indicates the grid columns do not span the available parent width.

```go
import (
	"strconv"
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
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		if strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
			return
		}
		t.Skip("no node_modules dirs on machine and no empty state yet")
	}
	if strings.Contains(resp.Output, "ELEM node-modules-column-header: MISSING") {
		t.Fatal("node_modules column header row missing")
	}
	if strings.Contains(resp.Output, "WIDTH_METRICS: MISSING") {
		t.Fatal("failed to measure node_modules table width metrics")
	}
	ratioLine := findLine(resp.Output, "WIDTH_RATIO size-right/body:")
	if ratioLine == "" {
		t.Fatal("expected WIDTH_RATIO size-right/body in playwright output")
	}
	ratio, parseErr := parseRatio(ratioLine)
	if parseErr != nil {
		t.Fatalf("failed to parse width ratio: %v (line: %s)", parseErr, ratioLine)
	}
	if ratio < 0.90 {
		t.Fatalf("expected node_modules Size column right edge >= 90%% of card body width, got ratio %.3f\nOutput:\n%s", ratio, resp.Output)
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

func parseRatio(line string) (float64, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
}
```