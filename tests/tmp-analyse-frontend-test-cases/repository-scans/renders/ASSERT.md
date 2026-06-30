## Expected

- `section-repository-scans-heading` is visible.
- `worktrees-section`, `worktrees-scan-btn`, `worktrees-stop-btn` are present.
- `binaries-section`, `binaries-scan-btn`, `binaries-stop-btn` are present.
- No element is reported as MISSING.

## Errors

- playwright-debug exits non-zero when elements are missing.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	required := []string{
		"ELEM section-repository-scans-heading",
		"ELEM worktrees-section",
		"ELEM worktrees-scan-btn",
		"ELEM binaries-section",
		"ELEM binaries-scan-btn",
	}
	for _, prefix := range required {
		if !strings.Contains(resp.Output, prefix) {
			t.Fatalf("missing output line: %s", prefix)
		}
	}
	if strings.Contains(resp.Output, "MISSING") {
		t.Fatalf("missing elements in output:\n%s", resp.Output)
	}
}
```