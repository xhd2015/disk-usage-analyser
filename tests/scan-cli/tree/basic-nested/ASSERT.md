## Expected

- `totalSize` is 500.
- Root children: `gamma.txt` (300 B file) and `alpha/` (200 B dir), sorted by size desc.
- `alpha` → `beta` → `deep.bin` (200 B file) nested chain with correct sizes and depths.

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
	if resp.TreeResult.TotalSize != 500 {
		t.Fatalf("totalSize: expected 500, got %d", resp.TreeResult.TotalSize)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 root children, got %d: %#v", len(root.Children), root.Children)
	}
	if root.Children[0].Name != "gamma.txt" || root.Children[0].Size != 300 || root.Children[0].IsDir {
		t.Fatalf("first child: expected gamma.txt file 300 B, got %#v", root.Children[0])
	}
	alpha := treeChildByName(root.Children, "alpha")
	if alpha == nil || alpha.Size != 200 || !alpha.IsDir || alpha.Depth != 1 {
		t.Fatalf("alpha: expected dir 200 B depth 1, got %#v", alpha)
	}
	if len(alpha.Children) != 1 {
		t.Fatalf("alpha children: expected 1, got %d", len(alpha.Children))
	}
	beta := alpha.Children[0]
	if beta.Name != "beta" || beta.Size != 200 || !beta.IsDir || beta.Depth != 2 {
		t.Fatalf("beta: expected dir 200 B depth 2, got %#v", beta)
	}
	if len(beta.Children) != 1 {
		t.Fatalf("beta children: expected 1, got %d", len(beta.Children))
	}
	deep := beta.Children[0]
	if deep.Name != "deep.bin" || deep.Size != 200 || deep.IsDir || deep.Depth != 3 {
		t.Fatalf("deep.bin: expected file 200 B depth 3, got %#v", deep)
	}
}
```