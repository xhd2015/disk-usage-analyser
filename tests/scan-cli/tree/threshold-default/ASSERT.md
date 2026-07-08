## Expected

- `threshold` is 1 MiB (`1048576`).
- Only `large.bin` appears in `tree.children`; `small.bin` is omitted from display.
- `totalSize` includes both files (`2 MiB + 512 KiB`).

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
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	root := assertRootTree(t, resp.TreeResult, req.FixtureDir)
	if resp.TreeResult.Threshold != 1<<20 {
		t.Fatalf("threshold: expected %d, got %d", 1<<20, resp.TreeResult.Threshold)
	}
	wantTotal := int64(2<<20 + 512<<10)
	if resp.TreeResult.TotalSize != wantTotal {
		t.Fatalf("totalSize: expected %d, got %d", wantTotal, resp.TreeResult.TotalSize)
	}
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 visible child, got %d: %#v", len(root.Children), root.Children)
	}
	if root.Children[0].Name != "large.bin" {
		t.Fatalf("visible child: expected large.bin, got %q", root.Children[0].Name)
	}
	if treeChildByName(root.Children, "small.bin") != nil {
		t.Fatal("small.bin must be omitted from tree display")
	}
}
```