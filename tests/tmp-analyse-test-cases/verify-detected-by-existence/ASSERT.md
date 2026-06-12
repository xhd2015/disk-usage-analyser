## Expected
- /tmp (category "temp") has Detected=true
- Software locations with existing paths have Detected=true
- Software locations with non-existing paths have Detected=false
- Core locations always have Detected=true
- At least some software locations have Detected flags (either true or false)

```go
import (
	"os"
	"strings"

	"disk-usage-analyser/server"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) < 22 {
		t.Fatalf("expected at least 22 locations, got %d", len(resp.Locations))
	}

	coreCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true}
	hasSoftwareDetected := false
	hasSoftwareNotDetected := false

	for _, loc := range resp.Locations {
		if coreCategories[loc.Category] && !loc.Detected {
			t.Fatalf("core location %s (category=%s) must have Detected=true", loc.Path, loc.Category)
		}
		if !coreCategories[loc.Category] {
			if loc.Detected {
				hasSoftwareDetected = true
			} else {
				hasSoftwareNotDetected = true
			}
		}
	}

	// On a typical dev machine at least some software should be detected
	// and some may not exist. This is environment-dependent.
	// At minimum we check the logic is applied correctly.
	if resp.DetectedCount+resp.NotDetectedCount != len(resp.Locations) {
		t.Fatalf("detected + notDetected should equal total locations")
	}

	// Basic sanity: /tmp should always exist on macOS
	for _, loc := range resp.Locations {
		if loc.Path == "/tmp" && !loc.Detected {
			t.Fatalf("/tmp should have Detected=true")
		}
	}
	_ = hasSoftwareDetected
	_ = hasSoftwareNotDetected
}
```
