---
label: slow, ui-automation
---

## Expected

- After node_modules scan completes, title area shows a running-total element with a formatted size.
- Running-total text contains a digit and a size unit (B, KB, MB, GB, etc.).

## Errors

- Running-total element must not be MISSING.
- Running-total text must contain a formatted size (digit + unit).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM node-modules-running-total: MISSING") {
		t.Fatal("node_modules running-total element missing after scan")
	}
	if strings.Contains(resp.Output, "CHECK node-modules-running-total-size: NO_SIZE") {
		t.Fatal("node_modules running-total does not contain a formatted size")
	}
}
```
