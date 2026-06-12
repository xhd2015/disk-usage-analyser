## Preconditions
- Go and Xcode locations each have ExtraPaths that should be scanned alongside the primary path
- After scanning, ExtraSizes and ExtraCounts arrays must be populated with results

## Steps
1. Create a mock filesystem with a primary path and extra paths containing known-size files
2. Call ScanWithProgress on the primary path, then on extra paths
3. Simulate the multi-path scan logic that populates ExtraSizes and ExtraCounts
4. Verify ExtraSizes and ExtraCounts have correct length and values

## Context
- For Go: primary = go/pkg/mod (e.g. 1000 bytes), extra = go-build cache (e.g. 500 bytes)
- ExtraSizes[i] corresponds to ExtraPaths[i]; same index for ExtraCounts
- Single-path tools have nil/empty ExtraSizes and ExtraCounts

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
	}
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	primaryFS := req.FS.(fs.FS)
	extraFS := fstest.MapFS{
		"x.txt": &fstest.MapFile{Data: make([]byte, 500)},
		"y.txt": &fstest.MapFile{Data: make([]byte, 300)},
	}

	primarySize, primaryCount, err := server.CalculateSize(primaryFS, ".")
	if err != nil {
		return nil, err
	}

	extraSize, extraCount, err := server.CalculateSize(extraFS, ".")
	if err != nil {
		return nil, err
	}

	return &Response{
		Size:        primarySize,
		FileCount:   primaryCount,
		ExtraSizes:  []int64{extraSize},
		ExtraCounts: []int64{extraCount},
	}, nil
}
```
