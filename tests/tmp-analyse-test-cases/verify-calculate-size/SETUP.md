# Scenario

**Feature**: CalculateSize sums file sizes correctly

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- A mock filesystem exists with three files: "a.txt" (100 bytes), "b.txt" (200 bytes), "sub/c.txt" (300 bytes)

## Steps
1. Create a mock filesystem using testing/fstest.MapFS with three files of known sizes
2. Set req.FS to the mock filesystem
3. Calculate the total size by summing all file sizes: 100 + 200 + 300 = 600 bytes

```go
import (
	"io/fs"
	"testing/fstest"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "calculate-size"
	req.FS = fstest.MapFS{
		"a.txt":     &fstest.MapFile{Data: make([]byte, 100)},
		"b.txt":     &fstest.MapFile{Data: make([]byte, 200)},
		"sub/c.txt": &fstest.MapFile{Data: make([]byte, 300)},
	}
	return nil
}

```
