## Expected

- Summary `Size` is 4096 (each inode counted once).
- `SharedHardlinkSize` is 0 (no hard links).
- `SharedCloneSize` is 12288 (`4096 × 3`).

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
	summary := resp.Result.Summary
	if summary.Size != file4K {
		t.Fatalf("expected size %d, got %d", file4K, summary.Size)
	}
	if summary.SharedHardlinkSize != 0 {
		t.Fatalf("expected shared_hardlink 0, got %d", summary.SharedHardlinkSize)
	}
	if summary.SharedCloneSize != file4K*3 {
		t.Fatalf("expected shared_clone %d, got %d", file4K*3, summary.SharedCloneSize)
	}
}
```