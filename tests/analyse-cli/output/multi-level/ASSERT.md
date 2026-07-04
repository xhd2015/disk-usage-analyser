## Expected

- Two immediate-child rows: `alpha`, `gamma` (no `alpha/beta` grandchild row).
- Row sizes are subtree totals: `alpha`=2048, `gamma`=1024; summary=3072.
- Summary path equals the analysed root.

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
	type rowExpect struct {
		path string
		size int64
	}
	wantRows := []rowExpect{
		{filepath.Join(req.FixtureDir, "alpha"), 2048},
		{filepath.Join(req.FixtureDir, "gamma"), 1024},
	}
	if len(result.Rows) != len(wantRows) {
		t.Fatalf("expected %d rows, got %d: %#v", len(wantRows), len(result.Rows), result.Rows)
	}
	for i, want := range wantRows {
		got := result.Rows[i]
		if got.Path != want.path {
			t.Fatalf("row %d path: expected %q, got %q", i, want.path, got.Path)
		}
		if got.Size != want.size {
			t.Fatalf("row %d size: expected %d, got %d", i, want.size, got.Size)
		}
	}
	if result.Summary.Size != 3072 {
		t.Fatalf("summary size: expected 3072, got %d", result.Summary.Size)
	}
	if result.Summary.Path != req.FixtureDir {
		t.Fatalf("summary path: expected %q, got %q", req.FixtureDir, result.Summary.Path)
	}
}
```