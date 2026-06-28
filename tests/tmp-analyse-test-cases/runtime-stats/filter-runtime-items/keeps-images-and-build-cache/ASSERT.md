## Expected
- Containers and Local Volumes excluded; Images and Build Cache kept.

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
	for _, item := range resp.RuntimeItems {
		if item.Type != "Images" && item.Type != "Build Cache" {
			t.Fatalf("unexpected type %s", item.Type)
		}
	}
	if len(resp.RuntimeItems) != 2 { t.Fatalf("expected 2 items, got %d", len(resp.RuntimeItems)) }
}
```
