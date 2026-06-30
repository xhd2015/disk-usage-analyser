## Expected

- Repo order: `/b` (20 MB), `/c` (10 MB), `/a` (5 MB).
- Sizes array is monotonic non-increasing.

## Errors

- Harness import failure until `repositoryScansLayout.ts` exists.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["worktreeRepoOrder"])
	assertOrder(t, order, []string{"/b", "/c", "/a"})
	sizes := jsonIntSlice(t, resp.JSON["worktreeRepoSizes"])
	assertMonotonicDesc(t, sizes)
	if sizes[0] != 20971520 {
		t.Fatalf("largest repo should be 20 MB, got %d", sizes[0])
	}
}
```