## Expected

- One immediate-child row for `node_modules` (project root has no other subdirs).
- Summary `BunSharedSize` is 4096 (one unique matching clone key, not 8192 from double scan).
- Summary `BunSharedSize` is less than or equal to summary `Size` (bun invariant on project root).
- `node_modules` row `BunSharedSize` is 4096 and `BunSharedSize <= Size`.

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
	if len(result.Rows) != 1 {
		t.Fatalf("expected one immediate-child row, got %d: %#v", len(result.Rows), result.Rows)
	}
	nmPath := filepath.Join(req.FixtureDir, "node_modules")
	nm := rowByPath(result, nmPath)
	if nm == nil {
		t.Fatalf("missing node_modules row in %#v", result.Rows)
	}
	if nm.BunSharedSize != file4K {
		t.Fatalf("node_modules bun_shared: expected %d (dedup once), got %d", file4K, nm.BunSharedSize)
	}
	if nm.BunSharedSize > nm.Size {
		t.Fatalf("node_modules bun_shared (%d) must be <= size (%d)", nm.BunSharedSize, nm.Size)
	}
	summary := result.Summary
	if summary.BunSharedSize != file4K {
		t.Fatalf("project root bun_shared: expected %d (dedup once), got %d", file4K, summary.BunSharedSize)
	}
	if summary.BunSharedSize > summary.Size {
		t.Fatalf("project root bun_shared (%d) must be <= size (%d)", summary.BunSharedSize, summary.Size)
	}
}
```