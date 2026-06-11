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
	req.FS = fstest.MapFS{
		"a.txt":     &fstest.MapFile{Data: make([]byte, 100)},
		"b.txt":     &fstest.MapFile{Data: make([]byte, 200)},
		"sub/c.txt": &fstest.MapFile{Data: make([]byte, 300)},
		"sub/d.txt": &fstest.MapFile{Data: make([]byte, 400)},
	}
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	fsys := req.FS.(fs.FS)
	var progressSizes []int64
	var progressCounts []int64

	size, count, err := server.ScanWithProgress(fsys, ".", func(s int64, c int64) {
		progressSizes = append(progressSizes, s)
		progressCounts = append(progressCounts, c)
	})
	if err != nil {
		return nil, err
	}

	// Store progress and final values
	resp := &Response{Size: size, FileCount: count}
	// Encode progress data into size as metadata for assertion
	resp.TotalSize = int64(len(progressSizes))
	if len(progressSizes) > 0 {
		resp.ReclaimableSize = progressSizes[len(progressSizes)-1]
	}
	return resp, nil
}
```
