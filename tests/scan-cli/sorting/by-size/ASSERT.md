## Expected

- Five root children in order: `large.txt`, `dir-medium`, `dir-tie`, `tie.txt`, `small.txt`.
- Sizes: 500, 200, 100, 100, 50.
- `dir-tie` precedes `tie.txt` when both are 100 B.
- `totalSize` is 950.

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
	if len(root.Children) != 5 {
		t.Fatalf("expected 5 children, got %d: %#v", len(root.Children), root.Children)
	}
	wantOrder := []struct {
		name  string
		size  int64
		isDir bool
	}{
		{"large.txt", 500, false},
		{"dir-medium", 200, true},
		{"dir-tie", 100, true},
		{"tie.txt", 100, false},
		{"small.txt", 50, false},
	}
	for i, want := range wantOrder {
		got := root.Children[i]
		if got.Name != want.name || got.Size != want.size || got.IsDir != want.isDir {
			t.Fatalf("child[%d]: expected %#v, got %#v", i, want, got)
		}
	}
	assertTreeChildrenSorted(t, root.Children)
	if resp.TreeResult.TotalSize != 950 {
		t.Fatalf("totalSize: expected 950, got %d", resp.TreeResult.TotalSize)
	}
}
```