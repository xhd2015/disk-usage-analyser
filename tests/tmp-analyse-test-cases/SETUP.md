## Preconditions
- The server package provides TmpLocation and TmpSummary types with fields: Path, Label, Category, Size, FileCount, RebootSafe (for TmpLocation) and Locations, TotalSize, ReclaimableSize (for TmpSummary)
- TmpLocation also has Detected (bool), ExtraPaths ([]string), ExtraSizes ([]int64), ExtraCounts ([]int64) fields
- The server package provides DiscoverLocations(homeDir string) []TmpLocation, CalculateSize(fsys fs.FS, root string) (int64, int64, error), and BuildSummary(locations []TmpLocation) TmpSummary
- The server package provides HandleTmpAnalyse(w http.ResponseWriter, r *http.Request) for SSE streaming
- On non-darwin platforms, HandleTmpAnalyse sends an unsupported_platform event

## Steps
1. Verify that the server types compile and have expected fields
2. Each leaf test overrides Run to exercise the specific function under test

## Context
- Tests exercise the tmp-analyse feature of disk-usage-analyser
- Mock filesystems (testing/fstest.MapFS) are used for deterministic file size tests
- SSE format tests use httptest to capture event stream output
- DiscoverLocations returns 5 core macOS locations plus 17 software cache/log locations

```go
import (
	"disk-usage-analyser/server"
)

type Request struct {
	HomeDir string
	FS      interface{} // set to fs.FS by leaf setups
}

type Response struct {
	Locations       []server.TmpLocation
	Size            int64
	FileCount       int64
	TotalSize       int64
	ReclaimableSize int64
	SSEOutput       string
	CategoryCount   map[string]int
	DetectedCount   int
	NotDetectedCount int
	ExtraPaths      []string
	ExtraSizes      []int64
	ExtraCounts     []int64
}

func Setup(t *testing.T, req *Request) error {
	_ = server.TmpLocation{Path: "x", Label: "x", Category: "trash", Detected: true}
	_ = server.TmpSummary{}
	return nil
}
```
