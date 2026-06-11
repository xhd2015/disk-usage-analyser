## Preconditions
- Four pre-scanned TmpLocation entries are provided:
  - User Trash: size=1000, rebootSafe=true
  - User Caches: size=500, rebootSafe=true
  - System Temp: size=2000, rebootSafe=false
  - User Logs: size=300, rebootSafe=true

## Steps
1. Create four TmpLocation entries with known sizes and rebootSafe flags
2. Call BuildSummary with these locations
3. TotalSize should be 1000+500+2000+300 = 3800
4. ReclaimableSize should be 1000+500+300 = 1800 (only rebootSafe=true items)

```go
import (
	"disk-usage-analyser/server"
)

func Run(t *testing.T, req *Request) (*Response, error) {
	locations := []server.TmpLocation{
		{Path: "/Users/x/.Trash", Label: "User Trash", Category: "trash", Size: 1000, FileCount: 5, RebootSafe: true},
		{Path: "/Users/x/Library/Caches", Label: "User Caches", Category: "cache", Size: 500, FileCount: 10, RebootSafe: true},
		{Path: "/tmp", Label: "System Temp", Category: "temp", Size: 2000, FileCount: 3, RebootSafe: false},
		{Path: "/Users/x/Library/Logs", Label: "User Logs", Category: "log", Size: 300, FileCount: 8, RebootSafe: true},
	}
	summary := server.BuildSummary(locations)
	return &Response{
		TotalSize:       summary.TotalSize,
		ReclaimableSize: summary.ReclaimableSize,
	}, nil
}
```
