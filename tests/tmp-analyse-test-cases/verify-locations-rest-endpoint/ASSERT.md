## Expected
- Response Content-Type contains "application/json" (not "text/event-stream")
- Response body is valid JSON parsing to an array of TmpLocation objects
- Array has at least 26 entries (5 core + 21 software)
- At least one trash, cache, log, temp category entry present
- Software categories like go, npm, docker present
- Each location has Detected field (boolean)
- Core locations have Detected=true
- Total detected + not-detected equals total location count
- Response is immediate (no scanning delay)

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(resp.SSEOutput, "application/json") {
		t.Fatalf("expected Content-Type application/json, got: %s", resp.SSEOutput)
	}

	if len(resp.Locations) < 26 {
		t.Fatalf("expected at least 26 locations, got %d", len(resp.Locations))
	}

	if resp.CategoryCount["trash"] == 0 {
		t.Fatal("expected at least one trash location")
	}
	if resp.CategoryCount["cache"] == 0 {
		t.Fatal("expected at least one cache location")
	}
	if resp.CategoryCount["log"] == 0 {
		t.Fatal("expected at least one log location")
	}
	if resp.CategoryCount["temp"] == 0 {
		t.Fatal("expected at least one temp location")
	}
	if resp.CategoryCount["go"] == 0 {
		t.Fatal("expected at least one go location")
	}
	if resp.CategoryCount["docker"] == 0 {
		t.Fatal("expected at least one docker location")
	}

	if resp.DetectedCount+resp.NotDetectedCount != len(resp.Locations) {
		t.Fatalf("detected(%d) + notDetected(%d) != total(%d)", resp.DetectedCount, resp.NotDetectedCount, len(resp.Locations))
	}

	coreCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true}
	for _, loc := range resp.Locations {
		if loc.Path == "" {
			t.Fatal("location has empty Path")
		}
		if loc.Label == "" {
			t.Fatalf("location has empty Label")
		}
		if loc.Category == "" {
			t.Fatalf("location has empty Category")
		}
		if coreCategories[loc.Category] && !loc.Detected {
			t.Fatalf("core location %s must have Detected=true", loc.Path)
		}
	}
}
```
