# Scenario

**Feature**: Go and Xcode have ExtraPaths; single-path tools have none

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- Go and Xcode locations each span multiple directories on disk, requiring ExtraPaths
- Other software locations scan only a single directory

## Steps
1. Set req.HomeDir to "/Users/testuser"
2. Call DiscoverLocations
3. Verify Go has ExtraPaths containing the build cache path
4. Verify Xcode has ExtraPaths containing the CoreSimulator path
5. Verify all single-path software locations have no ExtraPaths

## Context
- Go: primary Path = ~/go/pkg/mod, ExtraPaths = [~/Library/Caches/go-build]
- Xcode: primary Path = ~/Library/Developer/Xcode/DerivedData, ExtraPaths = [~/Library/Developer/CoreSimulator/Devices]
- All other 15 software locations have empty ExtraPaths

```go
import (
	"path/filepath"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "multi-path-locations"
	req.HomeDir = "/Users/testuser"
	return nil
}

```
