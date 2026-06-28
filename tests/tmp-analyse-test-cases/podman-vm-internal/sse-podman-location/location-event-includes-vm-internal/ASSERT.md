## Expected

- Podman location SSE event includes non-nil `vmInternal` with 2 items.
- `runtimeItems` has 2 filtered entries (Images, Build Cache).
- `MachineRunning` true on vmInternal.

## Side Effects

- None (httptest + mocks).

## Errors

- Run must find Podman location event in SSE stream.

## Exit Code

- Test passes when both vmInternal and runtimeItems populated.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.VmInternal == nil {
		snippet := resp.SSEOutput
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		t.Fatalf("expected vmInternal on Podman location event, SSE snippet: %s", snippet)
	}
	if !resp.VmInternal.MachineRunning {
		t.Fatal("expected MachineRunning true on SSE vmInternal")
	}
	if len(resp.VmInternal.Items) != 2 {
		t.Fatalf("expected 2 vmInternal items on SSE event, got %d", len(resp.VmInternal.Items))
	}
	if len(resp.RuntimeItems) != 2 {
		t.Fatalf("expected 2 runtimeItems on SSE event, got %d", len(resp.RuntimeItems))
	}
}
```