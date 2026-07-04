## Expected

- One immediate-child row for `node_modules`.
- `node_modules` row `PnpmSharedSize` is 0 on non-darwin.
- Summary `PnpmSharedSize` is 0.

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
	if nm.PnpmSharedSize != 0 {
		t.Fatalf("non-darwin pnpm_shared: expected 0, got %d", nm.PnpmSharedSize)
	}
	if result.Summary.PnpmSharedSize != 0 {
		t.Fatalf("summary pnpm_shared on non-darwin: expected 0, got %d", result.Summary.PnpmSharedSize)
	}
}
```