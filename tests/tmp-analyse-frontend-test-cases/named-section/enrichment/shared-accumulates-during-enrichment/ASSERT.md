---
label: slow, ui-automation
---

## Expected

- When at least two `named_enriched` SSE events fire while rows remain, shared cells resolve incrementally (`METRIC shared-row-updates` >= 2, `CHECK shared-accumulating: true`).

## Errors

- SKIP when no node_modules dirs or fewer than two enrichments on machine.

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
	enriched := parseMetricInt(resp.Output, "METRIC enriched-sse-events:")
	if enriched < 2 {
		t.Skip("fewer than 2 named_enriched events on this machine")
	}
	if !strings.Contains(resp.Output, "CHECK shared-accumulating: true") {
		t.Fatalf("expected shared column to update incrementally during enrichment (>=2 row updates), got:\n%s", resp.Output)
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