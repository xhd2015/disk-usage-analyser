## Expected

- One immediate-child row for `node_modules` (project root has no other subdirs).
- Summary `PnpmSharedSize` is 4096 (one unique matching clone key, not 8192 from double scan).
- Summary `PnpmSharedSize` is less than or equal to summary `Size` (pnpm invariant on project root).
- `node_modules` row `PnpmSharedSize` is 4096 and `PnpmSharedSize <= Size`.

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
	if nm.PnpmSharedSize != file4K {
		t.Fatalf("node_modules pnpm_shared: expected %d (dedup once), got %d", file4K, nm.PnpmSharedSize)
	}
	if nm.PnpmSharedSize > nm.Size {
		t.Fatalf("node_modules pnpm_shared (%d) must be <= size (%d)", nm.PnpmSharedSize, nm.Size)
	}
	summary := result.Summary
	if summary.PnpmSharedSize != file4K {
		t.Fatalf("project root pnpm_shared: expected %d (dedup once), got %d", file4K, summary.PnpmSharedSize)
	}
	if summary.PnpmSharedSize > summary.Size {
		t.Fatalf("project root pnpm_shared (%d) must be <= size (%d)", summary.PnpmSharedSize, summary.Size)
	}
}
```