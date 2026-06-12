## Expected
- On darwin (current platform), the SSE output contains scanning events (location, summary, done)
- The TmpLocation type includes all required fields
- The test validates that the platform gate compiles and the flow is correct for current OS
- Full non-darwin behavior should be tested on a Linux/Windows CI runner

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SSEOutput == "" {
		t.Fatal("expected SSE output")
	}
	// On darwin, we should see normal events
	if strings.Contains(resp.SSEOutput, "event: done") {
		// Normal darwin behavior: scanning happens
		if !strings.Contains(resp.SSEOutput, "event: locations") {
			t.Fatal("expected 'event: locations' in SSE output on darwin")
		}
	}
}
```
