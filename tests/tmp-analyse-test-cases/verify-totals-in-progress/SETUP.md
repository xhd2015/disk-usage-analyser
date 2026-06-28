# Scenario

**Feature**: BuildProgressPayload accumulates global totals

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- `server.BuildProgressPayload` is a function that adds totalSize and reclaimableSize to progress event payloads

## Steps
1. Call BuildProgressPayload with accumulated values from previous locations
2. Verify totalSize = accumulatedTotal + curSize
3. Verify reclaimableSize = accumulatedReclaimable + (curSize if rebootSafe else 0)

```go
import (
	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "totals-in-progress"
	_ = server.TmpLocation{Path: "x", Label: "x"}
	return nil
}

```
