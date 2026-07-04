## Expected

- Analyse root path basename is `node_modules`.
- One immediate-child row for `pkg` (not a `node_modules` row).
- `pkg` row `BunSharedSize` is 4096 (cache-backed bytes attributed to that package).
- Summary `BunSharedSize` is 4096 (full root-walk upfront scan).
- Sum of immediate-child `BunSharedSize` equals summary `BunSharedSize` (breakdown reconciles with summary).

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
	if len(result.Rows) != 1 {
		t.Fatalf("expected one immediate-child row, got %d: %#v", len(result.Rows), result.Rows)
	}
	pkgPath := filepath.Join(req.FixtureDir, "pkg")
	pkg := rowByPath(result, pkgPath)
	if pkg == nil {
		t.Fatalf("missing pkg row in %#v", result.Rows)
	}
	if pkg.BunSharedSize != file4K {
		t.Fatalf("pkg bun_shared: expected %d (cache clone attributed to package), got %d", file4K, pkg.BunSharedSize)
	}
	var childSum int64
	for _, row := range result.Rows {
		childSum += row.BunSharedSize
	}
	summary := result.Summary
	if summary.BunSharedSize != file4K {
		t.Fatalf("summary bun_shared: expected full scan %d, got %d (child sum %d)", file4K, summary.BunSharedSize, childSum)
	}
	if childSum != summary.BunSharedSize {
		t.Fatalf("breakdown sum bun_shared: expected %d to match summary %d (got child sum %d)", summary.BunSharedSize, summary.BunSharedSize, childSum)
	}
}
```