## Expected

- `maxDepth` is 2.
- Tree shows root → `d1` → `d2`; `leaf.bin` is not displayed.
- `d2.size` and `totalSize` are 1000 (deeper bytes counted).

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
	if resp.TreeResult.MaxDepth != 2 {
		t.Fatalf("maxDepth: expected 2, got %d", resp.TreeResult.MaxDepth)
	}
	if resp.TreeResult.TotalSize != 1000 {
		t.Fatalf("totalSize: expected 1000, got %d", resp.TreeResult.TotalSize)
	}
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 root child, got %d", len(root.Children))
	}
	d1 := root.Children[0]
	if d1.Name != "d1" || d1.Size != 1000 || !d1.IsDir {
		t.Fatalf("d1: expected dir 1000 B, got %#v", d1)
	}
	if len(d1.Children) != 1 {
		t.Fatalf("d1 children: expected 1, got %d", len(d1.Children))
	}
	d2 := d1.Children[0]
	if d2.Name != "d2" || d2.Size != 1000 || !d2.IsDir {
		t.Fatalf("d2: expected dir 1000 B, got %#v", d2)
	}
	if len(d2.Children) != 0 {
		t.Fatalf("d2 children: expected 0 (leaf.bin beyond maxDepth), got %d: %#v", len(d2.Children), d2.Children)
	}
}
```