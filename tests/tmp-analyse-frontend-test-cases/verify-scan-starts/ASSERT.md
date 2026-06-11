## Expected
- `BUTTON start-scan-btn: visible=true` before click
- `BUTTON stop-scan-btn: visible=false` before click
- After click: `BUTTON after click start-scan-btn: visible=false`
- After click: `BUTTON after click stop-scan-btn: visible=true`
- `SSE_REQUEST` line shows request to `/api/tmp-analyse`
- `CARD_SIZE` lines show non-empty values (sizes were fetched)
- `SUMMARY total-size` and `SUMMARY reclaimable-size` show values

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if !strings.Contains(resp.Output, "BUTTON start-scan-btn: visible=true") {
		t.Fatal("expected start-scan-btn visible=true initially")
	}
	if !strings.Contains(resp.Output, "BUTTON stop-scan-btn: visible=false") {
		t.Fatal("expected stop-scan-btn visible=false initially")
	}
	if !strings.Contains(resp.Output, "BUTTON after click start-scan-btn: visible=false") {
		t.Fatal("expected start-scan-btn hidden after click")
	}
	if !strings.Contains(resp.Output, "BUTTON after click stop-scan-btn: visible=true") {
		t.Fatal("expected stop-scan-btn visible after click")
	}
	if !strings.Contains(resp.Output, "SSE_REQUEST") {
		t.Fatal("expected SSE request to /api/tmp-analyse")
	}
	if !strings.Contains(resp.Output, "CARD_SIZE:") {
		t.Fatal("expected CARD_SIZE lines in output")
	}
	if !strings.Contains(resp.Output, "SUMMARY total-size:") {
		t.Fatal("expected SUMMARY total-size in output")
	}
}
```
