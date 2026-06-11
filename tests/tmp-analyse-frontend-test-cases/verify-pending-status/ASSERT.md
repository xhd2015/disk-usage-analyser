## Expected
- Initially: 0 pending, 0 scanning, 0 done
- After click (EARLY): pending-badge count > 0 (cards show pending)
- Mid scan: scanning-badge count > 0 (some cards being scanned)
- Final: scanning-badge count == 0, pending-badge count == 0
- Final: done-badge count > 0 (cards completed)
- HAS_DONE_BADGES is true

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if !strings.Contains(resp.Output, "EARLY pending-badge count:") {
		t.Fatal("expected EARLY pending-badge count line")
	}
	if strings.Contains(resp.Output, "EARLY pending-badge count: 0") {
		t.Fatal("expected EARLY pending-badge count > 0 (cards should show pending on scan start)")
	}
	if !strings.Contains(resp.Output, "MID scanning-badge count:") {
		t.Fatal("expected MID scanning-badge count line")
	}
	if !strings.Contains(resp.Output, "HAS_DONE_BADGES: true") {
		t.Fatal("expected HAS_DONE_BADGES: true after scan")
	}
	if !strings.Contains(resp.Output, "FINAL scanning-badge count: 0") {
		t.Fatal("expected FINAL scanning-badge count: 0")
	}
	if !strings.Contains(resp.Output, "FINAL pending-badge count: 0") {
		t.Fatal("expected FINAL pending-badge count: 0")
	}
}
```
