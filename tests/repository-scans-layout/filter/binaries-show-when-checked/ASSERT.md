## Expected

- Both repos visible: `/tiny-only`, `/mixed`.
- All three binaries visible (`visibleBinaryCount` = 3).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["binaryRepoOrder"])
	if len(order) != 2 {
		t.Fatalf("expected 2 repos, got %d: %v", len(order), order)
	}
	count, ok := resp.JSON["visibleBinaryCount"].(float64)
	if !ok || int(count) != 3 {
		t.Fatalf("visibleBinaryCount want 3, got %v", resp.JSON["visibleBinaryCount"])
	}
}
```