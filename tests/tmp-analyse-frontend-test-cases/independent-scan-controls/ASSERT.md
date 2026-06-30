---
label: slow, ui-automation
---

## Expected

- Page-level `start-scan-btn` remains visible while worktrees scan is in progress.
- Clicking page-level Start Scan initiates cache scan (`scanning-badge` or progress) independently.
- Script logs `PAGE_START_VISIBLE during wt scan: true`.

## Errors

- Repository scans must not disable page-level scan controls.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if !strings.Contains(resp.Output, "PAGE_START_VISIBLE during wt scan: true") {
		t.Fatal("page-level start scan should stay available during worktrees scan")
	}
	if !strings.Contains(resp.Output, "PAGE_SCANNING: true") {
		t.Fatal("page-level scan should start independently")
	}
}
```