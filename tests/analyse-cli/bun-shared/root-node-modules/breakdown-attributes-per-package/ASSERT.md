## Expected

- Analyse root basename is `node_modules`.
- Two immediate-child rows: `pkg-a` (cache clone) and `pkg-b` (regular file).
- `pkg-a` row `BunSharedSize` is 4096 (cache-backed bytes attributed to that package).
- `pkg-b` row `BunSharedSize` is 0.
- Summary `BunSharedSize` is 4096.
- Sum of immediate-child `BunSharedSize` equals summary `BunSharedSize` so the breakdown column is not all zeros while the summary shows the full total.

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
	if pkgA.BunSharedSize != file4K {
		t.Fatalf("pkg-a bun_shared: expected %d (cache clone attributed to package), got %d", file4K, pkgA.BunSharedSize)
	}
	if pkgB.BunSharedSize != 0 {
		t.Fatalf("pkg-b bun_shared: expected 0, got %d", pkgB.BunSharedSize)
	}
	var childSum int64
	for _, row := range result.Rows {
		childSum += row.BunSharedSize
	}
	summary := result.Summary
	if summary.BunSharedSize != file4K {
		t.Fatalf("summary bun_shared: expected %d, got %d", file4K, summary.BunSharedSize)
	}
	if childSum != summary.BunSharedSize {
		t.Fatalf("breakdown sum bun_shared: expected %d to match summary %d (got child sum %d — breakdown shows 0B per row)", summary.BunSharedSize, summary.BunSharedSize, childSum)
	}
}
```