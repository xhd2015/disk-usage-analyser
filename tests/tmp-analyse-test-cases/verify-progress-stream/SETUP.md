# Scenario

**Feature**: ScanWithProgress fires progress callbacks

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- A mock filesystem with 4 files: a.txt (100), b.txt (200), sub/c.txt (300), sub/d.txt (400)
- total size = 1000, file count = 4

## Steps
1. Create a mock FS with 4 files
2. Call ScanWithProgress with a callback that records each progress call
3. Verify the callback was called, sizes are cumulative, final values match 1000 bytes and 4 files

```go
import (
	"io/fs"
	"testing/fstest"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "progress-stream"
	req.FS = fstest.MapFS{
		"a.txt":     &fstest.MapFile{Data: make([]byte, 100)},
		"b.txt":     &fstest.MapFile{Data: make([]byte, 200)},
		"sub/c.txt": &fstest.MapFile{Data: make([]byte, 300)},
		"sub/d.txt": &fstest.MapFile{Data: make([]byte, 400)},
	}
	return nil
}

```
