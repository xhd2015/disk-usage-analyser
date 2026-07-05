## Expected

- Untracked and missing-`gitTracked` hits remain visible.
- Tracked hit is hidden (`visibleCount` = 2).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("column filters harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["repoOrder"])
	assertOrder(t, order, []string{"/repos/mixed"})
	paths := jsonStringSlice(t, resp.JSON["visiblePaths"])
	assertOrder(t, paths, []string{
		"/repos/mixed/untracked/node_modules",
		"/repos/mixed/unknown-git/node_modules",
	})
	count, ok := resp.JSON["visibleCount"].(float64)
	if !ok || int(count) != 2 {
		t.Fatalf("visibleCount want 2, got %v", resp.JSON["visibleCount"])
	}
}
```