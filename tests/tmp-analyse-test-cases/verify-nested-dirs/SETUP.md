# Scenario

**Feature**: CalculateSize recursively sums nested dirs

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- A mock filesystem exists with deeply nested structure:
  - "a.txt" (10 bytes)
  - "d1/b.txt" (20 bytes)
  - "d1/d2/c.txt" (30 bytes)
  - "d1/d2/d3/d.txt" (40 bytes)

## Steps
1. Create a mock filesystem with nested directories
2. Set req.FS to the mock filesystem
3. CalculateSize should recursively sum all files = 10 + 20 + 30 + 40 = 100 bytes, 4 files

```go
import (
	"io/fs"
	"testing/fstest"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "nested-dirs"
	req.FS = fstest.MapFS{
		"a.txt":            &fstest.MapFile{Data: make([]byte, 10)},
		"d1/b.txt":         &fstest.MapFile{Data: make([]byte, 20)},
		"d1/d2/c.txt":      &fstest.MapFile{Data: make([]byte, 30)},
		"d1/d2/d3/d.txt":   &fstest.MapFile{Data: make([]byte, 40)},
	}
	return nil
}

```
