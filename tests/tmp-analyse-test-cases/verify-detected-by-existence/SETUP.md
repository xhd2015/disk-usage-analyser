## Preconditions
- DiscoverLocations must check if a directory actually exists on disk to set Detected=true
- Core locations (Trash, Caches, Logs, System Temp, System Tmp) always have Detected=true regardless

## Steps
1. Get the real home directory via os.UserHomeDir()
2. Call DiscoverLocations with the real home directory
3. Verify that locations with well-known existing paths (like /tmp) have Detected=true
4. Verify that locations with obviously non-existent paths (like ~/.composer/cache on a clean system) have Detected=false

## Context
- On macOS, /tmp always exists so all System Tmp entries have Detected=true
- The test relies on /tmp existing and a software path like ~/Library/android being absent or present realistically
- Core locations may have Detected=true even if the path doesn't exist (they are always considered detected)

```go
import (
	"os"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get home directory")
	}
	req.HomeDir = homeDir
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	categoryCount := make(map[string]int)
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		categoryCount[loc.Category]++
		if loc.Detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}
	return &Response{
		Locations:        locations,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
		CategoryCount:    categoryCount,
	}, nil
}
```
