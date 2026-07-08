## Expected

- `maxDepth` is 1.
- Root children: `dir/` (80 B) and `top.bin` (20 B); no nested `nested.bin` node.
- `dir.size` includes nested bytes; `totalSize` is 100.

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
	if resp.TreeResult.MaxDepth != 1 {
		t.Fatalf("maxDepth: expected 1, got %d", resp.TreeResult.MaxDepth)
	}
	if resp.TreeResult.TotalSize != 100 {
		t.Fatalf("totalSize: expected 100, got %d", resp.TreeResult.TotalSize)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 root children, got %d: %#v", len(root.Children), root.Children)
	}
	dir := treeChildByName(root.Children, "dir")
	if dir == nil || dir.Size != 80 || !dir.IsDir {
		t.Fatalf("dir: expected dir 80 B, got %#v", dir)
	}
	if len(dir.Children) != 0 {
		t.Fatalf("dir children: expected 0 at maxDepth 1, got %d", len(dir.Children))
	}
	top := treeChildByName(root.Children, "top.bin")
	if top == nil || top.Size != 20 || top.IsDir {
		t.Fatalf("top.bin: expected file 20 B, got %#v", top)
	}
}
```