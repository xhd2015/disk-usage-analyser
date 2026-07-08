---
label: ui-automation
explanation: Playwright checks xcode breakdown-row-0..4 on /tmp-analyse
---

## Expected

- When Xcode card is detected, `breakdown-row-0` through `breakdown-row-4` all exist.
- Row 4 label includes `DocumentationCache`.
- SKIP when Xcode card not detected.

## Side Effects

- None.

## Errors

- Missing any breakdown row 0–4 fails the test.

## Exit Code

- 0 on PASS or SKIP; non-zero when rows missing.

```go
import (
	"strconv"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP xcode-breakdown-five") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP xcode-breakdown-five") {
		t.Skip("Xcode card not detected on this machine")
	}
	for i := 0; i <= 4; i++ {
		rowKey := "ELEM xcode-breakdown-row-" + strconv.Itoa(i) + ": MISSING"
		if strings.Contains(resp.Output, rowKey) {
			t.Fatalf("missing breakdown row %d\n%s", i, resp.Output)
		}
	}
	if !strings.Contains(resp.Output, "FULL_PATH xcode-row4-documentation-cache: true") {
		t.Fatal("expected row 4 label to include DocumentationCache")
	}
}
```