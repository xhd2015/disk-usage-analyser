## Expected

- Two immediate-child rows: `link-file` (file symlink) and `link-dir` (dir symlink).
- `link-file` row: symlinks `1f+0d`, size 0 (target not followed).
- `link-dir` row: symlinks `0f+1d`, size 0 (target not followed).
- Summary unchanged: `1f+1d` symlinks, size 8192 (real files only), no hardlink metrics.

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
	if len(result.Rows) != 2 {
		t.Fatalf("expected two immediate-child rows, got %d: %#v", len(result.Rows), result.Rows)
	}
	linkFilePath := filepath.Join(req.FixtureDir, "link-file")
	linkDirPath := filepath.Join(req.FixtureDir, "link-dir")
	linkFile := rowByPath(result, linkFilePath)
	if linkFile == nil {
		t.Fatalf("missing link-file row in %#v", result.Rows)
	}
	if linkFile.SymlinkFiles != 1 || linkFile.SymlinkDirs != 0 {
		t.Fatalf("link-file symlinks: expected 1f+0d, got %df+%dd", linkFile.SymlinkFiles, linkFile.SymlinkDirs)
	}
	if linkFile.Size != 0 {
		t.Fatalf("link-file size: expected 0, got %d", linkFile.Size)
	}
	linkDir := rowByPath(result, linkDirPath)
	if linkDir == nil {
		t.Fatalf("missing link-dir row in %#v", result.Rows)
	}
	if linkDir.SymlinkFiles != 0 || linkDir.SymlinkDirs != 1 {
		t.Fatalf("link-dir symlinks: expected 0f+1d, got %df+%dd", linkDir.SymlinkFiles, linkDir.SymlinkDirs)
	}
	if linkDir.Size != 0 {
		t.Fatalf("link-dir size: expected 0, got %d", linkDir.Size)
	}
	summary := result.Summary
	if summary.SymlinkFiles != 1 || summary.SymlinkDirs != 1 {
		t.Fatalf("expected 1f+1d symlinks, got files=%d dirs=%d", summary.SymlinkFiles, summary.SymlinkDirs)
	}
	if summary.Size != 8192 {
		t.Fatalf("expected size 8192 (real files only), got %d", summary.Size)
	}
	if summary.HardlinkExtra != 0 || summary.HardlinkSize != 0 || summary.SharedHardlinkSize != 0 || summary.SharedCloneSize != 0 {
		t.Fatalf("expected zero hardlink metrics, got extra=%d size=%d shared_hardlink=%d shared_clone=%d",
			summary.HardlinkExtra, summary.HardlinkSize, summary.SharedHardlinkSize, summary.SharedCloneSize)
	}
}
```