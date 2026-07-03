---
label: slow, ui-automation
---

## Expected

- `vendor-scan-btn` remains visible while node_modules scan is in progress.
- Clicking vendor scan initiates its own scan independently.
- Script logs `VENDOR_BTN_VISIBLE during nm scan: true` and `VENDOR_SCANNING: true`.

## Errors

- Named section scans must not disable each other.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if !strings.Contains(resp.Output, "VENDOR_BTN_VISIBLE during nm scan: true") {
		t.Fatal("vendor scan button should stay available during node_modules scan")
	}
	if !strings.Contains(resp.Output, "VENDOR_SCANNING: true") {
		t.Fatal("vendor scan should start independently")
	}
}
```
