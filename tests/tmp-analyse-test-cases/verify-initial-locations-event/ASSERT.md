## Expected
- The first SSE event type is "locations" (not "location" singular)
- The data payload is valid JSON that parses as an array
- The array has at least 22 entries (5 core + 17 software)
- The array contains at least one entry with Detected=true
- The first event is sent before any "event: location" (singular) or "event: progress" events

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SSEOutput != "locations" {
		t.Fatalf("expected first event type to be 'locations', got '%s'", resp.SSEOutput)
	}
	if resp.DetectedCount+resp.NotDetectedCount < 22 {
		t.Fatalf("expected at least 22 locations in first event, got %d", resp.DetectedCount+resp.NotDetectedCount)
	}
	if resp.DetectedCount == 0 {
		t.Fatal("expected at least one location with Detected=true in initial locations event")
	}
}
```
