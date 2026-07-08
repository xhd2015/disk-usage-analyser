## Expected

- Two immediate children: `subdir` dir 150 B, `root.txt` file 50 B.
- `totalSize` is 200.

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
	sub := treeChildByName(root.Children, "subdir")
	if sub == nil || sub.Size != 150 || !sub.IsDir {
		t.Fatalf("subdir: expected dir 150 B, got %#v", sub)
	}
	rootFile := treeChildByName(root.Children, "root.txt")
	if rootFile == nil || rootFile.Size != 50 || rootFile.IsDir {
		t.Fatalf("root.txt: expected file 50 B, got %#v", rootFile)
	}
	if resp.TreeResult.TotalSize != 200 {
		t.Fatalf("totalSize: expected 200, got %d", resp.TreeResult.TotalSize)
	}
}
```