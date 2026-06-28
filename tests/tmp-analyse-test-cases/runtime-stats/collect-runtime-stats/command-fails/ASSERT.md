## Expected
- CLI failure: empty runtimeItems, nil error.

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
	if resp.CollectFailed { t.Fatal("expected graceful handling when command fails") }
	if len(resp.RuntimeItems) != 0 { t.Fatalf("expected empty slice, got %d", len(resp.RuntimeItems)) }
}
```
