## Preconditions
- DiscoverLocations includes swap at /private/var/vm/ as a core system location
- Swap has category="swap", rebootSafe=true, and a field indicating it is not reclaimable

## Steps
1. Call DiscoverLocations with a test home directory
2. Find the swap location in the results
3. Verify its properties

```go
import (
	"disk-usage-analyser/server"
)

func Run(t *testing.T, req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	resp := &Response{
		Locations: locations,
	}
	for _, loc := range locations {
		if loc.Category == "swap" {
			resp.Size = loc.Size
			resp.FileCount = loc.FileCount
			resp.SSEOutput = loc.Label
			break
		}
	}
	return resp, nil
}
```
