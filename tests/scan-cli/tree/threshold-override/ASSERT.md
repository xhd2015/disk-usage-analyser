## Expected

- `threshold` is 10 MiB.
- Only `huge.bin` appears in `tree.children`; `medium.bin` is omitted.
- `totalSize` is 20 MiB (both files counted).

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
	if resp.TreeResult.Threshold != 10<<20 {
		t.Fatalf("threshold: expected %d, got %d", 10<<20, resp.TreeResult.Threshold)
	}
	wantTotal := int64(20 << 20)
	if resp.TreeResult.TotalSize != wantTotal {
		t.Fatalf("totalSize: expected %d, got %d", wantTotal, resp.TreeResult.TotalSize)
	}
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 visible child, got %d: %#v", len(root.Children), root.Children)
	}
	if root.Children[0].Name != "huge.bin" || root.Children[0].Size != 15<<20 {
		t.Fatalf("visible child: expected huge.bin 15 MiB, got %#v", root.Children[0])
	}
	if treeChildByName(root.Children, "medium.bin") != nil {
		t.Fatal("medium.bin must be omitted from tree display")
	}
}
```