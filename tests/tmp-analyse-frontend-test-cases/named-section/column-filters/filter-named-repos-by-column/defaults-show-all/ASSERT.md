## Expected

- Both repos remain visible.
- All four hits remain visible (`visibleCount` = 4).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("column filters harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["repoOrder"])
	assertOrder(t, order, []string{"/repos/alpha", "/repos/beta"})
	paths := jsonStringSlice(t, resp.JSON["visiblePaths"])
	assertOrder(t, paths, []string{
		"/repos/alpha/node_modules",
		"/repos/alpha/packages/app/node_modules",
		"/repos/beta/node_modules",
		"/repos/beta/vendor/pkg/node_modules",
	})
	count, ok := resp.JSON["visibleCount"].(float64)
	if !ok || int(count) != 4 {
		t.Fatalf("visibleCount want 4, got %v", resp.JSON["visibleCount"])
	}
}
```