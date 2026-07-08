## Expected

- One immediate child: `subdir` directory with recursive size 300 B.
- `totalSize` is 300.

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
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child, got %d: %#v", len(root.Children), root.Children)
	}
	sub := root.Children[0]
	if sub.Name != "subdir" || sub.Size != 300 || !sub.IsDir {
		t.Fatalf("subdir: expected dir 300 B, got %#v", sub)
	}
	if resp.TreeResult.TotalSize != 300 {
		t.Fatalf("totalSize: expected 300, got %d", resp.TreeResult.TotalSize)
	}
}
```