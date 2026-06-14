## Expected
- npm card is present (if npm detected on the system)
- After scan, npm card shows EITHER:
  - breakdown-items with subdirectory rows, OR
  - a single card-path (if ~/.npm has no subdirectories)
- The UI correctly renders whatever the backend provides (dynamic behavior)

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP npm-breakdown") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "SKIP npm-breakdown") {
		t.Log("npm not detected or scan button unavailable, skipping npm breakdown check")
		return
	}

	hasBreakdown := strings.Contains(resp.Output, "NPM_HAS_BREAKDOWN: true")
	hasSinglePath := strings.Contains(resp.Output, "NPM_SINGLE_PATH: true")

	if !hasBreakdown && !hasSinglePath {
		t.Fatal("npm card after scan should show either breakdown items or single path, found neither")
	}

	if hasBreakdown {
		line := findLine(resp.Output, "NPM_BREAKDOWN_COUNT")
		if line == "" {
			t.Fatal("expected NPM_BREAKDOWN_COUNT when breakdown items present")
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
