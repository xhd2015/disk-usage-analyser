## Expected

- Analyse root basename is `node_modules`.
- Two immediate-child rows: `pkg-a` (store clone) and `pkg-b` (regular file).
- `pkg-a` row `PnpmSharedSize` is 4096 (store-backed bytes attributed to that package).
- `pkg-b` row `PnpmSharedSize` is 0.
- Summary `PnpmSharedSize` is 4096.
- Sum of immediate-child `PnpmSharedSize` equals summary `PnpmSharedSize`.

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
	if filepath.Base(result.Root) != "node_modules" {
		t.Fatalf("analyse root basename: expected node_modules, got %q", result.Root)
	}
	if len(result.Rows) != 2 {
		t.Fatalf("expected two immediate-child rows, got %d: %#v", len(result.Rows), result.Rows)
	}
	pkgA := rowByPath(result, filepath.Join(req.FixtureDir, "pkg-a"))
	pkgB := rowByPath(result, filepath.Join(req.FixtureDir, "pkg-b"))
	if pkgA == nil || pkgB == nil {
		t.Fatalf("missing pkg rows in %#v", result.Rows)
	}
	if pkgA.PnpmSharedSize != file4K {
		t.Fatalf("pkg-a pnpm_shared: expected %d (store clone attributed to package), got %d", file4K, pkgA.PnpmSharedSize)
	}
	if pkgB.PnpmSharedSize != 0 {
		t.Fatalf("pkg-b pnpm_shared: expected 0, got %d", pkgB.PnpmSharedSize)
	}
	var childSum int64
	for _, row := range result.Rows {
		childSum += row.PnpmSharedSize
	}
	summary := result.Summary
	if summary.PnpmSharedSize != file4K {
		t.Fatalf("summary pnpm_shared: expected %d, got %d", file4K, summary.PnpmSharedSize)
	}
	if childSum != summary.PnpmSharedSize {
		t.Fatalf("breakdown sum pnpm_shared: expected %d to match summary %d (got child sum %d — breakdown shows 0B per row)", summary.PnpmSharedSize, summary.PnpmSharedSize, childSum)
	}
}
```