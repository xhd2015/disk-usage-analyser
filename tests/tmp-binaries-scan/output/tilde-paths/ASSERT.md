## Expected

- Every binary hit `path` and `repoPath` start with `~/`.
- No hit exposes the absolute fixture `HomeDir`.

## Errors

- No harness error is returned.

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Binaries) == 0 {
		t.Fatal("expected at least one binary hit")
	}
	for _, hit := range resp.Binaries {
		if !strings.HasPrefix(hit.Path, "~/") {
			t.Fatalf("path %q should start with ~/", hit.Path)
		}
		if !strings.HasPrefix(hit.RepoPath, "~/") {
			t.Fatalf("repoPath %q should start with ~/", hit.RepoPath)
		}
		if strings.Contains(hit.Path, req.HomeDir) {
			t.Fatalf("path exposes absolute home: %q", hit.Path)
		}
	}
}
```
