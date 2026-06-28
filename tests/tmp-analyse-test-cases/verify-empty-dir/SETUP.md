# Scenario

**Feature**: CalculateSize handles empty directories

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- An empty mock filesystem exists with no files

## Steps
1. Create an empty mock filesystem using testing/fstest.MapFS
2. Set req.FS to the empty filesystem
3. CalculateSize should return size=0, fileCount=0, and no error

```go
import (
	"io/fs"
	"testing/fstest"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "empty-dir"
	req.FS = fstest.MapFS{}
	return nil
}

```
