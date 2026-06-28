## Expected
- Mock docker CLI returns filtered Images + Build Cache.

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
	if len(resp.RuntimeItems) != 2 { t.Fatalf("expected 2 items from mock docker, got %d", len(resp.RuntimeItems)) }
}
```
