## Expected
- Locations list is not empty
- Contains at least one "trash" category location (e.g., /Users/testuser/.Trash)
- Contains at least one "cache" category location (e.g., /Users/testuser/Library/Caches)
- Contains at least one "log" category location (e.g., /Users/testuser/Library/Logs)
- Contains at least one "temp" category location
- The user trash location has RebootSafe=true
- Every location has a non-empty Path, Label, and Category
- Every location has a valid Category (one of "trash", "temp", "cache", "log")

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) == 0 {
		t.Fatal("expected at least one location, got none")
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
	validCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true}
	hasTrashPath := false
	for _, loc := range resp.Locations {
		if loc.Path == "" {
			t.Fatal("location has empty Path")
		}
		if loc.Label == "" {
			t.Fatalf("location %s has empty Label", loc.Path)
		}
		if !validCategories[loc.Category] {
			t.Fatalf("location %s has invalid category: %s", loc.Path, loc.Category)
		}
		if strings.Contains(loc.Path, "/.Trash") {
			hasTrashPath = true
			if !loc.RebootSafe {
				t.Fatalf("trash location %s should be RebootSafe=true", loc.Path)
			}
		}
	}
	if !hasTrashPath {
		t.Fatal("expected at least one location with /.Trash in path")
	}
}
```
