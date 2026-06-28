## Expected
- size=4.2GB (accumulated card), breakdownIndex=1, breakdownSize=1.1GB.

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
	p := resp.ProgressPayload
	if p["size"].(int64) != 4200000000 { t.Fatalf("card size should be 3.1+1.1=4.2GB, got %d", p["size"]) }
	if p["breakdownIndex"].(int) != 1 { t.Fatalf("expected breakdownIndex=1, got %v", p["breakdownIndex"]) }
	if p["breakdownSize"].(int64) != 1100000000 { t.Fatalf("expected breakdownSize=1.1GB, got %d", p["breakdownSize"]) }
}
```
