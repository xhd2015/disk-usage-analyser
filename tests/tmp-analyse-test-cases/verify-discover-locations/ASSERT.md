## Expected
- Locations list has at least 22 entries (5 core + 17 software)
- Contains at least one "trash" category location (e.g., /Users/testuser/.Trash)
- Contains at least one "cache" category location (e.g., /Users/testuser/Library/Caches)
- Contains at least one "log" category location (e.g., /Users/testuser/Library/Logs)
- Contains at least one "temp" category location
- Contains software-specific categories (go, npm, bun, docker, etc.)
- The user trash location has RebootSafe=true
- Every location has a non-empty Path, Label, Category, and Detected field is a boolean
- Core locations (trash, cache, log, temp) have Detected=true
- Total of detected + not-detected equals total location count

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) < 22 {
		t.Fatalf("expected at least 22 locations (5 core + 17 software), got %d", len(resp.Locations))
	}
	if resp.CategoryCount["trash"] == 0 {
		t.Fatal("expected at least one trash category location")
	}
	if resp.CategoryCount["cache"] == 0 {
		t.Fatal("expected at least one cache category location")
	}
	if resp.CategoryCount["log"] == 0 {
		t.Fatal("expected at least one log category location")
	}
	if resp.CategoryCount["temp"] == 0 {
		t.Fatal("expected at least one temp category location")
	}
	if resp.DetectedCount+resp.NotDetectedCount != len(resp.Locations) {
		t.Fatalf("detected(%d) + notDetected(%d) != total(%d)", resp.DetectedCount, resp.NotDetectedCount, len(resp.Locations))
	}
	hasTrashPath := false
	coreCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true}
	for _, loc := range resp.Locations {
		if loc.Path == "" {
			t.Fatal("location has empty Path")
		}
		if loc.Label == "" {
			t.Fatalf("location %s has empty Label", loc.Path)
		}
		if loc.Category == "" {
			t.Fatalf("location %s has empty Category", loc.Path)
		}
		if strings.Contains(loc.Path, "/.Trash") {
			hasTrashPath = true
			if !loc.RebootSafe {
				t.Fatalf("trash location %s should be RebootSafe=true", loc.Path)
			}
		}
		if coreCategories[loc.Category] && !loc.Detected {
			t.Fatalf("core location %s (category=%s) must have Detected=true", loc.Path, loc.Category)
		}
	}
	if !hasTrashPath {
		t.Fatal("expected at least one location with /.Trash in path")
	}
}
```
