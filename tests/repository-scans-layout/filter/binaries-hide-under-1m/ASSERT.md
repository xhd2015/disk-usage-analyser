## Expected

- `/tiny-only` repo entirely hidden (total < 1 MiB after leaf filter).
- `/mixed` repo visible with only `/mixed/big` (2 MB); `/mixed/small` (800 KB) hidden.
- `visibleBinaryCount` = 1.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["binaryRepoOrder"])
	assertOrder(t, order, []string{"/mixed"})
	count, ok := resp.JSON["visibleBinaryCount"].(float64)
	if !ok || int(count) != 1 {
		t.Fatalf("visibleBinaryCount want 1, got %v", resp.JSON["visibleBinaryCount"])
	}
}
```