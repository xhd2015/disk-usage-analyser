## Expected

- Repo order: `/large` (15 MB total) before `/small` (3 MB total).
- Repo totals monotonic non-increasing.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["binaryRepoOrder"])
	assertOrder(t, order, []string{"/large", "/small"})
	totals := jsonIntSlice(t, resp.JSON["binaryRepoTotals"])
	assertMonotonicDesc(t, totals)
	if totals[0] != 15728640 {
		t.Fatalf("largest repo total should be 15 MB, got %d", totals[0])
	}
}
```