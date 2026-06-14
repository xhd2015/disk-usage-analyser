## Expected
- A swap card with `data-testid="card-swap"` exists in the System Locations section
- The swap card has label="Swap"
- The swap card has a size display
- The swap card has a reboot-safe badge

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	checks := map[string]string{
		"ELEM card-swap":         "",
		"ELEM card-swap-label":   "",
		"ELEM card-swap-size":    "",
		"SWAP_LABEL_CHECK":       "true",
	}

	for prefix, expected := range checks {
		line := findLine(resp.Output, prefix)
		if line == "" {
			t.Fatalf("missing expected output line: %s", prefix)
		}
		if strings.Contains(line, "MISSING") {
			t.Fatalf("element missing: %s", line)
		}
		if expected != "" && !strings.Contains(line, expected) {
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
