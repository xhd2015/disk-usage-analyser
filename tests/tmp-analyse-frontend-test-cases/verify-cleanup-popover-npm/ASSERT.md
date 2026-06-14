## Expected
- npm card has a clickable cleanup indicator
- Clicking it opens a popover with at least these suggestions:
  - `npm cache clean --force` (removes ~/.npm/_cacache/ content)
  - `npm cache verify` (verifies cache integrity)
- Each suggestion includes a description about what files are removed
- Each suggestion indicates recoverability (re-downloaded automatically)

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP npm-popover-test") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "SKIP npm-popover-test") {
		t.Log("npm not detected, skipping popover content check")
		return
	}

	checks := map[string]string{
		"CLEANUP_NPM_CACHE_CLEAN":   "true",
		"CLEANUP_NPM_CACHE_VERIFY":  "true",
		"CLEANUP_NPM_RECOVERABLE":   "true",
		"CLEANUP_NPM_HAS_DESCRIPTION": "true",
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
