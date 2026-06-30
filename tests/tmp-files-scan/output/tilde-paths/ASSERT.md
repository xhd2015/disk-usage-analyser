## Expected

- Stdout contains `~/Projects/tilde-app/bin/tilde-app`.
- Stdout contains `(repo: ~/Projects/tilde-app)`.
- Stdout does not contain the raw fixture home directory.
- Hit `Path` and `RepoPath` use tilde display form.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No error is returned.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one hit, got %#v", resp.Result)
	}
	hit := resp.Result.Binaries[0]
	if hit.Path != "~/Projects/tilde-app/bin/tilde-app" || hit.RepoPath != "~/Projects/tilde-app" {
		t.Fatalf("expected tilde hit paths, got %#v", hit)
	}
	if !strings.Contains(resp.Stdout, "~/Projects/tilde-app/bin/tilde-app") || !strings.Contains(resp.Stdout, "(repo: ~/Projects/tilde-app)") {
		t.Fatalf("tilde paths missing from stdout:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, req.HomeDir) {
		t.Fatalf("raw home path leaked into stdout:\n%s", resp.Stdout)
	}
}
```
