## Expected

- One immediate-child row for `node_modules`.
- `node_modules` row `BunSharedSize` is 0 (clone key not in cache index).
- `node_modules` row `Size` is 4096.
- Summary `BunSharedSize` is 0.

## Errors

- No error is returned.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected analyse error: %v", resp.Err)
	}
	result := resp.Result
	nmPath := filepath.Join(req.FixtureDir, "node_modules")
	nm := rowByPath(result, nmPath)
	if nm == nil {
		t.Fatalf("missing node_modules row in %#v", result.Rows)
	}
	if nm.Size != file4K {
		t.Fatalf("node_modules size: expected %d, got %d", file4K, nm.Size)
	}
	if nm.BunSharedSize != 0 {
		t.Fatalf("node_modules bun_shared: expected 0, got %d", nm.BunSharedSize)
	}
	if result.Summary.BunSharedSize != 0 {
		t.Fatalf("summary bun_shared: expected 0, got %d", result.Summary.BunSharedSize)
	}
}
```