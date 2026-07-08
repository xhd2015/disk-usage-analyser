## Expected

- Mock `simctl runtime list -j` plus walkable mountPath returns non-empty items.
- Single runtime: Type `iOS 18.5`, Size 1024, Reclaimable 1024, ActiveCount 1.

## Side Effects

- None beyond temp dir in test harness.

## Errors

- None.

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
	if resp.CollectFailed {
		t.Fatal("expected successful collection")
	}
	if len(resp.RuntimeItems) == 0 {
		t.Fatal("expected non-empty runtime items from mock success")
	}
	item := resp.RuntimeItems[0]
	if item.Type != "iOS 18.5" {
		t.Fatalf("expected Type iOS 18.5, got %s", item.Type)
	}
	if item.Size != 1024 {
		t.Fatalf("expected Size 1024, got %d", item.Size)
	}
	if item.Reclaimable != 1024 {
		t.Fatalf("expected Reclaimable 1024, got %d", item.Reclaimable)
	}
	if item.ActiveCount != 1 {
		t.Fatalf("expected ActiveCount 1, got %d", item.ActiveCount)
	}
}
```