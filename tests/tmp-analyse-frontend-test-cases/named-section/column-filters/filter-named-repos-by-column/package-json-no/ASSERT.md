## Expected

- False and missing `hasPackageJson` hits remain visible (`visibleCount` = 2).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("column filters harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["repoOrder"])
	assertOrder(t, order, []string{"/repos/pkg"})
	paths := jsonStringSlice(t, resp.JSON["visiblePaths"])
	assertOrder(t, paths, []string{
		"/repos/pkg/without/node_modules",
		"/repos/pkg/unknown/node_modules",
	})
	count, ok := resp.JSON["visibleCount"].(float64)
	if !ok || int(count) != 2 {
		t.Fatalf("visibleCount want 2, got %v", resp.JSON["visibleCount"])
	}
}
```