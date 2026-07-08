## Expected

- Only `visible.bin` appears in `tree.children`.
- `totalSize` is `2 MiB + 500` (hidden bytes counted at root).

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
	wantTotal := int64(2<<20 + 500)
	if resp.TreeResult.TotalSize != wantTotal {
		t.Fatalf("totalSize: expected %d, got %d", wantTotal, resp.TreeResult.TotalSize)
	}
	if len(root.Children) != 1 || root.Children[0].Name != "visible.bin" {
		t.Fatalf("expected only visible.bin in tree, got %#v", root.Children)
	}
	if treeChildByName(root.Children, "hidden.bin") != nil {
		t.Fatal("hidden.bin must be omitted from tree display")
	}
}
```