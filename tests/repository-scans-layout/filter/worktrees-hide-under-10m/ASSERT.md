## Expected

- Only `/big-main` repo visible (main checkout 20 MB).
- Linked children: only `/big-main/wt-big` (12 MB); `/big-main/wt-small` (3 MB) hidden.
- `/small-main` (5 MB main) hidden entirely.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["worktreeRepoOrder"])
	assertOrder(t, order, []string{"/big-main"})
	linkedRaw, ok := resp.JSON["linkedOrder"].(map[string]interface{})
	if !ok {
		t.Fatal("linkedOrder missing")
	}
	linked := jsonStringSlice(t, linkedRaw["/big-main"])
	assertOrder(t, linked, []string{"/big-main/wt-big"})
	count, ok := resp.JSON["visibleWorktreeRepoCount"].(float64)
	if !ok || int(count) != 1 {
		t.Fatalf("visibleWorktreeRepoCount want 1, got %v", resp.JSON["visibleWorktreeRepoCount"])
	}
}
```