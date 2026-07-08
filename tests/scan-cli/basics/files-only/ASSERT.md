## Expected

- Two file children: `a.txt` 100 B, `b.txt` 200 B.
- `totalSize` is 300.
- Both children have `isDir=false`.

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
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d: %#v", len(root.Children), root.Children)
	}
	a := treeChildByName(root.Children, "a.txt")
	if a == nil || a.Size != 100 || a.IsDir {
		t.Fatalf("a.txt: expected file 100 B, got %#v", a)
	}
	b := treeChildByName(root.Children, "b.txt")
	if b == nil || b.Size != 200 || b.IsDir {
		t.Fatalf("b.txt: expected file 200 B, got %#v", b)
	}
	if resp.TreeResult.TotalSize != 300 {
		t.Fatalf("totalSize: expected 300, got %d", resp.TreeResult.TotalSize)
	}
}
```