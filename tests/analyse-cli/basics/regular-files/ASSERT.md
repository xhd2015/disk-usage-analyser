## Expected

- One subdirectory row for `sub` with size 2048.
- Summary size is 3072 (1024 + 2048).
- All link columns are zero.

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
		t.Fatalf("expected one subdirectory row, got %d: %#v", len(result.Rows), result.Rows)
	}
	subPath := filepath.Join(req.FixtureDir, "sub")
	sub := result.Rows[0]
	if sub.Path != subPath {
		t.Fatalf("subdir path: expected %q, got %q", subPath, sub.Path)
	}
	if sub.Size != 2048 {
		t.Fatalf("subdir size: expected 2048, got %d", sub.Size)
	}
	if sub.SymlinkFiles != 0 || sub.SymlinkDirs != 0 || sub.HardlinkExtra != 0 {
		t.Fatalf("subdir link metrics should be zero: %#v", sub)
	}
	summary := result.Summary
	if summary.Size != 3072 {
		t.Fatalf("summary size: expected 3072, got %d", summary.Size)
	}
	if summary.HardlinkExtra != 0 || summary.HardlinkSize != 0 || summary.SharedHardlinkSize != 0 || summary.SharedCloneSize != 0 {
		t.Fatalf("summary hardlink metrics should be zero: %#v", summary)
	}
}
```