## Expected

- Go card shows scanning-badge before done-badge.
- At least one of breakdown-size-0 or breakdown-size-1 updates to a non-zero value during scan (before done-badge).
- Script logs `CHECK breakdown-0-live-update: true` or `CHECK breakdown-1-live-update: true`.

## Side Effects

- None (read-only UI observation).

## Errors

- SKIP when Go card not detected (non-detected environment).

## Exit Code

- playwright-debug exits 0 on PASS; non-zero when breakdown rows stay at zero during scan.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP breakdown-live") {
		t.Skip("Go card not detected on this machine")
	}
	if strings.Contains(resp.Output, "FAIL: breakdown rows did not update") {
		t.Fatal("breakdown rows did not update during scan")
	}
	if !strings.Contains(resp.Output, "CHECK breakdown-0-live-update: true") &&
		!strings.Contains(resp.Output, "CHECK breakdown-1-live-update: true") {
		t.Fatalf("expected live breakdown size update during scan\n%s", resp.Output)
	}
}
```