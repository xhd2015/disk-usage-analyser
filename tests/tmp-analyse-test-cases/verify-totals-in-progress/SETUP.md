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
	_ = server.TmpLocation{Path: "x", Label: "x"}
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	// Accumulated after Loc A (rebootSafe) with 1000 bytes
	// Now scanning Loc B (rebootSafe), curSize=200
	p1 := server.BuildProgressPayload("Loc B", 200, 10, 1000, 1000, true)
	// totalSize: 1000+200=1200, reclaimableSize: 1000+200=1200

	// Now scanning Loc C (not rebootSafe), curSize=500
	// After Loc B completes: accumulated=1200, reclaimable=1200
	p2 := server.BuildProgressPayload("Loc C", 500, 20, 1200, 1200, false)
	// totalSize: 1200+500=1700, reclaimableSize: 1200+0=1200

	resp := &Response{}
	resp.Size = p1["totalSize"].(int64)
	resp.FileCount = p1["reclaimableSize"].(int64)
	resp.TotalSize = p2["totalSize"].(int64)
	resp.ReclaimableSize = p2["reclaimableSize"].(int64)

	return resp, nil
}
```
