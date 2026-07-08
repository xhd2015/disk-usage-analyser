## Expected

- Parse yields 2 TmpRuntimeItem entries.
- Ready runtime: Type `iOS 18.5`, TotalCount 1, ActiveCount 1, Size 90GB, Reclaimable 90GB.
- Mounting runtime: Type `iOS 17.4`, TotalCount 1, ActiveCount 0, Size 45GB, Reclaimable 0.

## Side Effects

- None (pure parse).

## Errors

- None.

## Exit Code

- Test passes when expectations match.

```go
import (
	"testing"

	"disk-usage-analyser/server"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.RuntimeItems) != 2 {
		t.Fatalf("expected 2 runtime items, got %d", len(resp.RuntimeItems))
	}
	byType := make(map[string]server.TmpRuntimeItem, len(resp.RuntimeItems))
	for _, item := range resp.RuntimeItems {
		byType[item.Type] = item
	}
	ready, ok := byType["iOS 18.5"]
	if !ok {
		t.Fatalf("missing Ready runtime iOS 18.5, got %+v", resp.RuntimeItems)
	}
	if ready.TotalCount != 1 || ready.ActiveCount != 1 {
		t.Fatalf("ready runtime counts: total=%d active=%d", ready.TotalCount, ready.ActiveCount)
	}
	if ready.Size != 90000000000 {
		t.Fatalf("ready runtime size: got %d", ready.Size)
	}
	if ready.Reclaimable != 90000000000 {
		t.Fatalf("ready runtime reclaimable: got %d", ready.Reclaimable)
	}

	mounting, ok := byType["iOS 17.4"]
	if !ok {
		t.Fatalf("missing Mounting runtime iOS 17.4, got %+v", resp.RuntimeItems)
	}
	if mounting.TotalCount != 1 || mounting.ActiveCount != 0 {
		t.Fatalf("mounting runtime counts: total=%d active=%d", mounting.TotalCount, mounting.ActiveCount)
	}
	if mounting.Size != 45000000000 {
		t.Fatalf("mounting runtime size: got %d", mounting.Size)
	}
	if mounting.Reclaimable != 0 {
		t.Fatalf("mounting runtime reclaimable: got %d", mounting.Reclaimable)
	}
}
```