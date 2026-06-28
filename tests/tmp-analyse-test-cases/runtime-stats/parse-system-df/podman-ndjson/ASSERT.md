## Expected
- Podman Image row normalized to Images; Build Cache retained.

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
if images.Type != "Images" { t.Fatalf("expected normalized Images, got %s", images.Type) }
if images.TotalCount != 9 { t.Fatalf("expected TotalCount=9, got %d", images.TotalCount) }
if images.Size != 4500000000 { t.Fatalf("expected Size=4.5GB, got %d", images.Size) }
}
```
