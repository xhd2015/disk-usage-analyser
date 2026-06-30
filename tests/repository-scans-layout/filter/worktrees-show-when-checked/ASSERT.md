## Expected

- Both repos visible: `/small-main`, `/big-main`.
- Both linked worktrees under `/big-main` visible.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("layout harness failed: %v\n%s", err, resp.Output)
	}
	order := jsonStringSlice(t, resp.JSON["worktreeRepoOrder"])
	if len(order) != 2 {
		t.Fatalf("expected 2 repos, got %d: %v", len(order), order)
	}
	linkedRaw, ok := resp.JSON["linkedOrder"].(map[string]interface{})
	if !ok {
		t.Fatal("linkedOrder missing")
	}
	linked := jsonStringSlice(t, linkedRaw["/big-main"])
	if len(linked) != 2 {
		t.Fatalf("expected 2 linked worktrees, got %d", len(linked))
	}
}
```