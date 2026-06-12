## Expected
- Go location has exactly 1 ExtraPath pointing to the go build cache
- Xcode location has exactly 1 ExtraPath pointing to CoreSimulator Devices
- All other 15 single-path software locations have empty ExtraPaths (length 0)

```go
import (
	"path/filepath"
	"strings"

	"disk-usage-analyser/server"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) != 17 {
		t.Fatalf("expected 17 software locations, got %d", len(resp.Locations))
	}

	multiPathCats := map[string]int{"go": 1, "xcode": 1}
	for _, loc := range resp.Locations {
		expectedExtraCount, isMulti := multiPathCats[loc.Category]
		if isMulti {
			if len(loc.ExtraPaths) != expectedExtraCount {
				t.Fatalf("location %s: expected %d ExtraPaths, got %d", loc.Category, expectedExtraCount, len(loc.ExtraPaths))
			}
			if loc.Category == "go" {
				if !strings.Contains(loc.ExtraPaths[0], "Caches/go-build") && !strings.Contains(loc.ExtraPaths[0], "Cache/go-build") {
					t.Fatalf("Go ExtraPaths[0] should contain go-build cache path, got: %s", loc.ExtraPaths[0])
				}
			}
			if loc.Category == "xcode" {
				if !strings.Contains(loc.ExtraPaths[0], "CoreSimulator") {
					t.Fatalf("Xcode ExtraPaths[0] should contain CoreSimulator path, got: %s", loc.ExtraPaths[0])
				}
			}
		} else {
			if len(loc.ExtraPaths) != 0 {
				t.Fatalf("location %s: expected 0 ExtraPaths, got %d", loc.Category, len(loc.ExtraPaths))
			}
		}
	}

	// Also verify ExtraPaths were captured in the response
	if len(resp.ExtraPaths) < 2 {
		t.Fatalf("expected at least 2 ExtraPaths (one for Go, one for Xcode), got %d", len(resp.ExtraPaths))
	}
}
```
