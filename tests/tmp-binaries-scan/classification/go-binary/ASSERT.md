## Expected

- Exactly one binary hit with `kind=go`.
- `typeDesc` is non-empty.
- Path contains `bin/go-app`.

## Errors

- No harness error is returned.

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Binaries) != 1 {
		t.Fatalf("expected 1 hit, got %d: %#v", len(resp.Binaries), resp.Binaries)
	}
	hit := resp.Binaries[0]
	if hit.Kind != "go" {
		t.Fatalf("kind = %q, want go", hit.Kind)
	}
	if hit.TypeDesc == "" {
		t.Fatal("expected non-empty typeDesc")
	}
	if !strings.Contains(hit.Path, "bin/go-app") {
		t.Fatalf("unexpected path: %q", hit.Path)
	}
}
```
