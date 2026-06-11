## Expected
- `BUTTON scan-started stop-scan-btn: visible=true` (scan started)
- `BUTTON after stop start-scan-btn: visible=true` (scan stopped, start button shown)
- `BUTTON after stop stop-scan-btn: visible=false` (stop button hidden)
- No FAIL lines in output

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "FAIL:") {
		t.Fatalf("test failed:\n%s", resp.Output)
	}

	if !strings.Contains(resp.Output, "BUTTON scan-started stop-scan-btn: visible=true") {
		t.Fatal("expected stop button visible after starting scan")
	}
	if !strings.Contains(resp.Output, "BUTTON after stop start-scan-btn: visible=true") {
		t.Fatal("expected start button visible after stopping scan")
	}
	if !strings.Contains(resp.Output, "BUTTON after stop stop-scan-btn: visible=false") {
		t.Fatal("expected stop button hidden after stopping scan")
	}
}
```
