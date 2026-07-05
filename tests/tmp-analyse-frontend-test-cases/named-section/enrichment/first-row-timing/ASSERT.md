---
label: slow, ui-automation
---

## Expected

- First `node-modules-row` appears within **10 seconds** of scan click (`CHECK first-row-within-budget: true`).
- Script logs `TIMING first-row-ui-ms` <= 10000.

## Errors

- SKIP when no node_modules dirs on machine.

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
	if strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
		return
	}
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		t.Skip("no node_modules dirs on machine")
	}
	if !strings.Contains(resp.Output, "CHECK first-row-within-budget: true") {
		ms := parseTimingMs(resp.Output, "TIMING first-row-ui-ms:")
		t.Fatalf("expected first row within 10s, got first-row-ui-ms=%d\nOutput:\n%s", ms, resp.Output)
	}
}

func parseTimingMs(output, prefix string) int {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			v, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), prefix)))
			return v
		}
	}
	return -1
}
```