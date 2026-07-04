## Expected

- `Result.Rows` is empty (no subdirectories).
- Summary `Size` is 0.
- All symlink and hardlink metrics on summary are zero.

## Errors

- No error is returned.

```go
import (
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
	if len(result.Rows) != 0 {
		t.Fatalf("expected no subdirectory rows, got %d: %#v", len(result.Rows), result.Rows)
	}
	summary := result.Summary
	if summary.Size != 0 {
		t.Fatalf("expected summary size 0, got %d", summary.Size)
	}
	if summary.SymlinkFiles != 0 || summary.SymlinkDirs != 0 {
		t.Fatalf("expected zero symlinks, got files=%d dirs=%d", summary.SymlinkFiles, summary.SymlinkDirs)
	}
	if summary.HardlinkExtra != 0 || summary.HardlinkSize != 0 || summary.SharedHardlinkSize != 0 || summary.SharedCloneSize != 0 {
		t.Fatalf("expected zero hardlink metrics, got extra=%d size=%d shared_hardlink=%d shared_clone=%d",
			summary.HardlinkExtra, summary.HardlinkSize, summary.SharedHardlinkSize, summary.SharedCloneSize)
	}
	if summary.Path != req.FixtureDir {
		t.Fatalf("summary path: expected %q, got %q", req.FixtureDir, summary.Path)
	}
}
```