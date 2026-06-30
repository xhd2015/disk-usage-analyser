## Expected

- Final order: `/new-large` (25 MB), `/mid` (10 MB), `/small` (5 MB).
- Largest size is first after insert.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["worktreeRepoOrder"])
	assertOrder(t, order, []string{"/new-large", "/mid", "/small"})
	sizes := jsonIntSlice(t, resp.JSON["worktreeRepoSizes"])
	assertMonotonicDesc(t, sizes)
	if sizes[0] != 26214400 {
		t.Fatalf("new larger item should be first at 25 MB, got %d", sizes[0])
	}
}
```