# Scenario

**Feature**: Home-relative paths use ~ prefix

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- DiscoverLocations should return tilde-shortened paths for all home-relative directories
- Paths not under the home directory (e.g., /tmp, /usr/local) should remain absolute

## Steps
1. Set req.HomeDir to "/Users/testuser"
2. Call DiscoverLocations with the home directory
3. Return all locations for assertion

```go
import (
	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "tilde-paths"
	req.HomeDir = "/Users/testuser"
	return nil
}

```
