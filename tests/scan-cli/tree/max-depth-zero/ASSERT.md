## Expected

- `maxDepth` is 0.
- Tree shows full chain `a` → `b` → `c` → `deep.bin` at depth 4.
- `totalSize` is 42.

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
	if resp.TreeResult.MaxDepth != 0 {
		t.Fatalf("maxDepth: expected 0, got %d", resp.TreeResult.MaxDepth)
	}
	if resp.TreeResult.TotalSize != 42 {
		t.Fatalf("totalSize: expected 42, got %d", resp.TreeResult.TotalSize)
	}
	a := root.Children[0]
	if a.Name != "a" || !a.IsDir {
		t.Fatalf("a: expected dir, got %#v", a)
	}
	b := a.Children[0]
	if b.Name != "b" || !b.IsDir {
		t.Fatalf("b: expected dir, got %#v", b)
	}
	c := b.Children[0]
	if c.Name != "c" || !c.IsDir {
		t.Fatalf("c: expected dir, got %#v", c)
	}
	deep := c.Children[0]
	if deep.Name != "deep.bin" || deep.Size != 42 || deep.IsDir || deep.Depth != 4 {
		t.Fatalf("deep.bin: expected file 42 B depth 4, got %#v", deep)
	}
}
```