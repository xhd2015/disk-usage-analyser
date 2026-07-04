## Expected

- One immediate-child row for `project` (no top-level `node_modules` row).
- `project` row `PnpmSharedSize` is 4096 (aggregated from nested `node_modules`).
- `project` row `Size` is 4096.
- Summary `PnpmSharedSize` is 4096.

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
	projectPath := filepath.Join(req.FixtureDir, "project")
	project := rowByPath(result, projectPath)
	if project == nil {
		t.Fatalf("missing project row in %#v", result.Rows)
	}
	nmPath := filepath.Join(req.FixtureDir, "node_modules")
	if rowByPath(result, nmPath) != nil {
		t.Fatalf("unexpected top-level node_modules row")
	}
	if project.Size != file4K {
		t.Fatalf("project size: expected %d, got %d", file4K, project.Size)
	}
	if project.PnpmSharedSize != file4K {
		t.Fatalf("project pnpm_shared: expected %d, got %d", file4K, project.PnpmSharedSize)
	}
	summary := result.Summary
	if summary.PnpmSharedSize != file4K {
		t.Fatalf("summary pnpm_shared: expected %d, got %d", file4K, summary.PnpmSharedSize)
	}
}
```