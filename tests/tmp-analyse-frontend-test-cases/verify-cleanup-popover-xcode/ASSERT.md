## Expected
- Xcode card has a clickable cleanup indicator
- Clicking it opens a popover with at least:
  - DerivedData cleanup suggestion
  - Simulator cleanup suggestion with `simctl shutdown all` and `simctl delete all`
- Each suggestion indicates recoverability

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP xcode-popover-test") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "SKIP xcode-popover-test") {
		t.Log("xcode not detected, skipping popover content check")
		return
	}

	checks := map[string]string{
		"CLEANUP_XCODE_DERIVED":    "true",
		"CLEANUP_XCODE_SIMULATORS": "true",
		"CLEANUP_XCODE_RECOVERABLE": "true",
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
