## Expected
- Go card has a clickable cleanup indicator
- Clicking it opens a popover with at least these suggestions:
  - `go clean -cache` (removes ~/Library/Caches/go-build)
  - `go clean -modcache` (removes ~/go/pkg/mod)
- Each suggestion indicates recoverability (rebuilt/re-downloaded automatically)

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP go-popover-test") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "SKIP go-popover-test") {
		t.Log("go not detected, skipping popover content check")
		return
	}

	checks := map[string]string{
		"CLEANUP_GO_CACHE":      "true",
		"CLEANUP_GO_MODCACHE":   "true",
		"CLEANUP_GO_RECOVERABLE": "true",
	}

	for key, expected := range checks {
		line := findLine(resp.Output, key)
		if line == "" {
			t.Fatalf("missing expected output line: %s", key)
		}
		if !strings.Contains(line, expected) {
			t.Fatalf("line %q: expected to contain %q", line, expected)
		}
	}
}

func findLine(output, prefix string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return line
		}
	}
	return ""
}
```
