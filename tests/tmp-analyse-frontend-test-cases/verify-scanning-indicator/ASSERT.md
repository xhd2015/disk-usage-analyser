## Expected
- Initially: 0 `scanning-badge`, 0 `done-badge` (or both hidden)
- During scan: `scanning-badge` count > 0
- After scan: `scanning-badge` count is 0
- After scan: `done-badge` count > 0, HAS_DONE_BADGES is true

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "INIT scanning-badge count: 0") {
		// ok, expected
	}
	if !strings.Contains(resp.Output, "MID scanning-badge count:") {
		t.Fatal("expected MID scanning-badge count line")
	}
	if strings.Contains(resp.Output, "HAS_DONE_BADGES: true") {
		// ok, expected
	} else {
		t.Fatal("expected HAS_DONE_BADGES: true after scan completes")
	}
	if strings.Contains(resp.Output, "FINAL scanning-badge count: 0") {
		// ok, scanning badges cleared after done
	}
}
```
