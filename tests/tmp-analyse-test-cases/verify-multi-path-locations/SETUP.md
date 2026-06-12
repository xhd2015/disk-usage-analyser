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
	req.HomeDir = "/Users/testuser"
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	coreCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true}
	var softwareLocs []server.TmpLocation
	var goLoc, xcodeLoc *server.TmpLocation
	for i, loc := range locations {
		if !coreCategories[loc.Category] {
			softwareLocs = append(softwareLocs, loc)
			if loc.Category == "go" {
				goLoc = &locations[i]
			}
			if loc.Category == "xcode" {
				xcodeLoc = &locations[i]
			}
		}
	}

	extraPaths := []string{}
	if goLoc != nil {
		extraPaths = goLoc.ExtraPaths
	}
	if xcodeLoc != nil {
		extraPaths = append(extraPaths, xcodeLoc.ExtraPaths...)
	}

	return &Response{
		Locations:  softwareLocs,
		ExtraPaths: extraPaths,
	}, nil
}
```
