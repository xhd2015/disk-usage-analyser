## Expected

- One immediate-child row for `node_modules`.
- `node_modules` row `PnpmSharedSize` is 4096 (one unique matching clone key).
- `node_modules` row `Size` is 4096.
- Summary `PnpmSharedSize` is 4096 (aligned with immediate-child rows).

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
	if nm.Size != file4K {
		t.Fatalf("node_modules size: expected %d, got %d", file4K, nm.Size)
	}
	if nm.PnpmSharedSize != file4K {
		t.Fatalf("node_modules pnpm_shared: expected %d, got %d", file4K, nm.PnpmSharedSize)
	}
	summary := result.Summary
	if summary.PnpmSharedSize != file4K {
		t.Fatalf("summary pnpm_shared: expected %d, got %d", file4K, summary.PnpmSharedSize)
	}
}
```