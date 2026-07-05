## Expected

- Only `/repos/pkg/with/node_modules` remains visible (`visibleCount` = 1).

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
	assertOrder(t, paths, []string{"/repos/pkg/with/node_modules"})
	count, ok := resp.JSON["visibleCount"].(float64)
	if !ok || int(count) != 1 {
		t.Fatalf("visibleCount want 1, got %v", resp.JSON["visibleCount"])
	}
}
```