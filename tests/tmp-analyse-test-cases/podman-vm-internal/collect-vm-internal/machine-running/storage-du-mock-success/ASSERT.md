## Expected

- Mock SSH returns fixtures → 2 `TmpVmStorageItem` rows.
- `MachineRunning` is true.
- `TotalSize` equals primary storage du bytes (1735229168).

## Side Effects

- None (mock runner).

## Errors

- Collection must not error; `VmInternal` populated.

## Exit Code

- Test passes when two labeled items present.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatalf("expected successful collection, got: %v", resp.Err)
	}
	if resp.VmInternal == nil {
		t.Fatal("expected non-nil VmInternal")
	}
	if !resp.VmInternal.MachineRunning {
		t.Fatal("expected MachineRunning true")
	}
	if len(resp.VmInternal.Items) != 2 {
		t.Fatalf("expected 2 VmStorageItems, got %d", len(resp.VmInternal.Items))
	}
	if resp.VmInternal.Items[0].Label != "Container storage" {
		t.Fatalf("expected first label Container storage, got %q", resp.VmInternal.Items[0].Label)
	}
	if resp.VmInternal.Items[1].Label != "Overlay layers" {
		t.Fatalf("expected second label Overlay layers, got %q", resp.VmInternal.Items[1].Label)
	}
	if resp.VmInternal.Items[0].Size != 1735229168 {
		t.Fatalf("expected storage size 1735229168, got %d", resp.VmInternal.Items[0].Size)
	}
	if resp.VmInternal.Items[1].Size != 1900000000 {
		t.Fatalf("expected overlay size 1900000000, got %d", resp.VmInternal.Items[1].Size)
	}
	if resp.VmInternal.TotalSize != 1735229168 {
		t.Fatalf("expected TotalSize 1735229168, got %d", resp.VmInternal.TotalSize)
	}
}
```