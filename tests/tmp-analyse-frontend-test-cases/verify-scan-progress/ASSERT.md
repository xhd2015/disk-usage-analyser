## Expected
- SSE request to /api/tmp-analyse is fired
- Card sizes change during scan (MID_SIZE shows intermediate values)
- After scan completes, at least one FINAL_SIZE is non-zero
- HAS_NONZERO is true

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if !strings.Contains(resp.Output, "SSE_REQUEST") {
		t.Fatal("expected SSE request to /api/tmp-analyse")
	}
	if !strings.Contains(resp.Output, "MID_SIZE:") {
		t.Fatal("expected MID_SIZE lines (intermediate progress)")
	}
	if !strings.Contains(resp.Output, "FINAL_SIZE:") {
		t.Fatal("expected FINAL_SIZE lines")
	}
	if !strings.Contains(resp.Output, "HAS_NONZERO: true") {
		t.Fatal("expected at least one card to have non-zero final size")
	}
}
```
