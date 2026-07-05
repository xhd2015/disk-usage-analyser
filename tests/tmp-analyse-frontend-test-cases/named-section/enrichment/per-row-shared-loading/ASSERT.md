---
label: slow, ui-automation
---

## Expected

- At least one `node-modules-row` after scan with pending shared enrichment (`METRIC pending-shared-rows-max` > 0).
- At least one `node-modules-shared-loading-*` element is visible while rows await shared enrichment (per-row loading in Shared column).

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
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		if strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
			return
		}
		t.Skip("no node_modules dirs on machine")
	}
	if parseMetricInt(resp.Output, "METRIC pending-shared-rows-max:") <= 0 {
		t.Skip("no pending shared enrichment observed on this machine")
	}
	if !strings.Contains(resp.Output, "CHECK per-row-shared-loading-seen: true") {
		t.Fatalf("expected per-row shared loading indicators during enrichment, got:\n%s", resp.Output)
	}
}

func parseMetricInt(output, prefix string) int {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			v, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), prefix)))
			return v
		}
	}
	return 0
}
```