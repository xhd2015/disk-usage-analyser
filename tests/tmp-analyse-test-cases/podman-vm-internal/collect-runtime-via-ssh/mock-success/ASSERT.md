## Expected

- Mock inner `podman system df` → 2 filtered `runtimeItems` (Images, Build Cache).

## Side Effects

- None (mock runner).

## Errors

- Collection must not error.

## Exit Code

- Test passes when two runtime items returned.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatalf("expected successful runtime collection, got: %v", resp.Err)
	}
	if len(resp.RuntimeItems) != 2 {
		t.Fatalf("expected 2 runtime items from inner VM system df, got %d", len(resp.RuntimeItems))
	}
	if resp.RuntimeItems[0].Type != "Images" {
		t.Fatalf("expected first type Images, got %q", resp.RuntimeItems[0].Type)
	}
	if resp.RuntimeItems[1].Type != "Build Cache" {
		t.Fatalf("expected second type Build Cache, got %q", resp.RuntimeItems[1].Type)
	}
}
```