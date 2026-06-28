## Expected
- Completed row 0 size retained; active row 1 adds partial.

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
	last := resp.ProgressSequence[len(resp.ProgressSequence)-1]
	if last["size"].(int64) != 3300000000 { t.Fatalf("expected retained completed+active=3.3GB, got %d", last["size"]) }
}
```
