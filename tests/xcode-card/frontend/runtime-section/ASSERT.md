---
label: ui-automation
explanation: Playwright waits for xcode scan then checks runtime-section
---

## Expected

- After scan, if `XCODE_RUNTIME_SECTION: present` then runtime-row-0, runtime-label-0, runtime-count-0, runtime-size-0 must be present.
- If `XCODE_RUNTIME_SECTION: absent`, `XCODE_RUNTIME_GRACEFUL: true` is logged (no simctl / no runtimes).
- SKIP when Xcode card not detected or simctl unavailable.

## Side Effects

- None.

## Errors

- Fails when section present but rows missing.

## Exit Code

- 0 when runtime section present with rows, or gracefully absent.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP xcode-runtime") {
		t.Skip("Xcode card not detected on this machine")
	}
	if strings.Contains(resp.Output, "SKIP xcode-runtime-no-simctl") {
		t.Skip("simctl runtimes not available on this machine")
	}
	if strings.Contains(resp.Output, "XCODE_RUNTIME_SECTION: present") {
		for _, elem := range []string{"runtime-row-0", "runtime-label-0", "runtime-count-0", "runtime-size-0"} {
			if strings.Contains(resp.Output, "ELEM "+elem+": MISSING") {
				t.Fatalf("expected %s when runtime section present", elem)
			}
		}
		return
	}
	if !strings.Contains(resp.Output, "XCODE_RUNTIME_GRACEFUL: true") {
		t.Fatalf("expected runtime section or graceful absence\n%s", resp.Output)
	}
}
```