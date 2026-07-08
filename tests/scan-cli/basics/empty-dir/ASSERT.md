## Expected

- Root `tree.children` is empty.
- `totalSize` is 0.
- `path` is the absolute fixture directory.

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
	if len(root.Children) != 0 {
		t.Fatalf("expected no children, got %d: %#v", len(root.Children), root.Children)
	}
	if resp.TreeResult.TotalSize != 0 {
		t.Fatalf("expected totalSize 0, got %d", resp.TreeResult.TotalSize)
	}
}
```