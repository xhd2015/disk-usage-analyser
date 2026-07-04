## Expected

- Summary `Size` is 4096.
- `HardlinkExtra` is 2.
- `HardlinkSize` is 4096.
- `SharedHardlinkSize` is 12288 (`4096 × 3`).
- `SharedCloneSize` is 0.

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
	if summary.HardlinkExtra != 2 {
		t.Fatalf("expected hardlinks 2, got %d", summary.HardlinkExtra)
	}
	if summary.HardlinkSize != file4K {
		t.Fatalf("expected hardlink_size %d, got %d", file4K, summary.HardlinkSize)
	}
	if summary.SharedHardlinkSize != file4K*3 {
		t.Fatalf("expected shared_hardlink %d, got %d", file4K*3, summary.SharedHardlinkSize)
	}
	if summary.SharedCloneSize != 0 {
		t.Fatalf("expected shared_clone 0, got %d", summary.SharedCloneSize)
	}
}
```