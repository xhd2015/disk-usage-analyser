## Expected
- First row scanning: card size equals active partial only.

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
	if p["size"].(int64) != 500000000 { t.Fatalf("expected card size=active row only, got %d", p["size"]) }
	if p["breakdownIndex"].(int) != 0 { t.Fatal("expected breakdownIndex=0") }
}
```
