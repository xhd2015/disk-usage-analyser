## Expected
- Row 0 frozen at 3.1GB; row 1 partial 0.8GB; card total 3.9GB.

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
	if p["size"].(int64) != 3900000000 { t.Fatalf("expected 3.1+0.8=3.9GB card total, got %d", p["size"]) }
	if p["breakdownIndex"].(int) != 1 { t.Fatal("expected active row index 1") }
}
```
