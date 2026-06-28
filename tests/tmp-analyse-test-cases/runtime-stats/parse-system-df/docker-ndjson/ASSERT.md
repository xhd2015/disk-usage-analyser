## Expected
- Parse docker NDJSON; filter keeps Images (12 total, 8 active, 8.3GB/1.5GB reclaimable) and Build Cache (34 total).

## Side Effects
- None (pure function or mock CLI).

## Errors
- See leaf scenario for expected error vs graceful-empty behavior.

## Exit Code
- Test passes when expectations match.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.RuntimeItems) != 2 { t.Fatalf("expected 2 filtered items, got %d", len(resp.RuntimeItems)) }
images := resp.RuntimeItems[0]
if images.Type != "Images" { t.Fatalf("expected Images, got %s", images.Type) }
if images.TotalCount != 12 { t.Fatalf("expected TotalCount=12, got %d", images.TotalCount) }
if images.ActiveCount != 8 { t.Fatalf("expected ActiveCount=8, got %d", images.ActiveCount) }
if images.Size != 8300000000 { t.Fatalf("expected Size=8.3GB, got %d", images.Size) }
if images.Reclaimable != 1500000000 { t.Fatalf("expected Reclaimable=1.5GB, got %d", images.Reclaimable) }
cache := resp.RuntimeItems[1]
if cache.Type != "Build Cache" { t.Fatalf("expected Build Cache, got %s", cache.Type) }
if cache.TotalCount != 34 { t.Fatalf("expected TotalCount=34, got %d", cache.TotalCount) }
}
```
