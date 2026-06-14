## Preconditions
- TmpLocation has a Reclaimable field (defaults to rebootSafe for existing locations)
- BuildSummary uses Reclaimable field instead of rebootSafe
- Swap location has Reclaimable=false and rebootSafe=true

## Steps
1. Create locations including a swap entry with Reclaimable=false, rebootSafe=true
2. Call BuildSummary
3. Verify swap's size is included in TotalSize but excluded from ReclaimableSize

```go
import (
	"disk-usage-analyser/server"
)

func Run(t *testing.T, req *Request) (*Response, error) {
	locations := []server.TmpLocation{
		{Path: "/Users/x/.Trash", Label: "User Trash", Category: "trash", Size: 1000, FileCount: 5, RebootSafe: true, Reclaimable: true},
		{Path: "/private/var/vm/", Label: "Swap", Category: "swap", Size: 2048, FileCount: 3, RebootSafe: true, Reclaimable: false},
		{Path: "/tmp", Label: "System Temp", Category: "temp", Size: 500, FileCount: 2, RebootSafe: false, Reclaimable: false},
	}

	summary := server.BuildSummary(locations)
	return &Response{
		TotalSize:       summary.TotalSize,
		ReclaimableSize: summary.ReclaimableSize,
	}, nil
}
```
