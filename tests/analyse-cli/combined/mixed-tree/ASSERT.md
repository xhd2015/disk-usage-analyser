## Expected

- Two immediate-child rows: `pkg` (real dir) and `link-pkg` (dir symlink).
- `pkg` subtree: size 4608, symlinks `1f+0d`, hardlinks extra=1, hardlink_size=4096, shared_hardlink=8192, shared_clone=0.
- `link-pkg` row: dir symlink (no follow); symlinks `1f+1d` (link itself + `pkg/link-readme` via target).
- Summary unchanged: symlinks `2f+1d`, hardlinks extra=1, size 4608.

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
	pkgPath := filepath.Join(req.FixtureDir, "pkg")
	linkPkgPath := filepath.Join(req.FixtureDir, "link-pkg")
	pkg := rowByPath(result, pkgPath)
	if pkg == nil {
		t.Fatalf("missing pkg row in %#v", result.Rows)
	}
	if pkg.SymlinkFiles != 1 || pkg.SymlinkDirs != 0 {
		t.Fatalf("pkg symlinks: expected 1f+0d, got %df+%dd", pkg.SymlinkFiles, pkg.SymlinkDirs)
	}
	if pkg.HardlinkExtra != 1 || pkg.HardlinkSize != file4K || pkg.SharedHardlinkSize != file4K*2 || pkg.SharedCloneSize != 0 {
		t.Fatalf("pkg hardlink metrics: got extra=%d size=%d shared_hardlink=%d shared_clone=%d",
			pkg.HardlinkExtra, pkg.HardlinkSize, pkg.SharedHardlinkSize, pkg.SharedCloneSize)
	}
	if pkg.Size != 512+file4K {
		t.Fatalf("pkg size: expected %d, got %d", 512+file4K, pkg.Size)
	}
	linkPkg := rowByPath(result, linkPkgPath)
	if linkPkg == nil {
		t.Fatalf("missing link-pkg row in %#v", result.Rows)
	}
	if linkPkg.SymlinkFiles != 1 || linkPkg.SymlinkDirs != 1 {
		t.Fatalf("link-pkg symlinks: expected 1f+1d, got %df+%dd", linkPkg.SymlinkFiles, linkPkg.SymlinkDirs)
	}
	if linkPkg.HardlinkExtra != 0 || linkPkg.HardlinkSize != 0 || linkPkg.SharedHardlinkSize != 0 || linkPkg.SharedCloneSize != 0 {
		t.Fatalf("link-pkg hardlink metrics should be zero: %#v", linkPkg)
	}
	if linkPkg.Size != 0 {
		t.Fatalf("link-pkg size: expected 0 (symlink, no follow), got %d", linkPkg.Size)
	}

	summary := result.Summary
	if summary.SymlinkFiles != 2 || summary.SymlinkDirs != 1 {
		t.Fatalf("summary symlinks: expected 2f+1d, got %df+%dd", summary.SymlinkFiles, summary.SymlinkDirs)
	}
	if summary.HardlinkExtra != 1 || summary.HardlinkSize != file4K || summary.SharedHardlinkSize != file4K*2 || summary.SharedCloneSize != 0 {
		t.Fatalf("summary hardlink metrics: got extra=%d size=%d shared_hardlink=%d shared_clone=%d",
			summary.HardlinkExtra, summary.HardlinkSize, summary.SharedHardlinkSize, summary.SharedCloneSize)
	}
	if summary.Size != 512+file4K {
		t.Fatalf("summary size: expected %d, got %d", 512+file4K, summary.Size)
	}
}
```