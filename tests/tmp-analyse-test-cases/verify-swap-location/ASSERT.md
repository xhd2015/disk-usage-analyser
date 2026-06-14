## Expected
- DiscoverLocations returns a location with category="swap"
- The swap location has label="Swap"
- The swap location has rebootSafe=true
- The swap location has Reclaimable=false
- The swap location path is "/private/var/vm/" (or "/var/vm/")
- Swap is in the system (core) locations group, not software

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var swapFound bool
	for _, loc := range resp.Locations {
		if loc.Category == "swap" {
			swapFound = true
			if loc.Label != "Swap" {
				t.Fatalf("expected swap Label='Swap', got %q", loc.Label)
			}
			if !loc.RebootSafe {
				t.Fatalf("expected swap RebootSafe=true")
			}
			pathLower := strings.ToLower(loc.Path)
			if !strings.Contains(pathLower, "vm") && !strings.Contains(pathLower, "swap") {
				t.Fatalf("expected swap path to reference vm/swap, got %q", loc.Path)
			}
			break
		}
	}

	if !swapFound {
		t.Fatalf("expected to find a location with category='swap' in DiscoverLocations, but none found")
	}
}
```
