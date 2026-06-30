## Expected

- Linked order under `/repo`: `wt-large` (8 MB), `wt-mid` (5 MB), `wt-small` (2 MB).
- Child sizes monotonic non-increasing.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	linkedRaw, ok := resp.JSON["linkedOrder"].(map[string]interface{})
	if !ok {
		t.Fatal("linkedOrder missing from harness output")
	}
	order := jsonStringSlice(t, linkedRaw["/repo"])
	assertOrder(t, order, []string{"/repo/wt-large", "/repo/wt-mid", "/repo/wt-small"})
	sizesRaw, ok := resp.JSON["linkedSizes"].(map[string]interface{})
	if !ok {
		t.Fatal("linkedSizes missing from harness output")
	}
	sizes := jsonIntSlice(t, sizesRaw["/repo"])
	assertMonotonicDesc(t, sizes)
}
```