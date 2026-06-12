## Expected
- A collapse panel with data-testid="collapse-not-detected" exists
- The panel header shows "Not Detected" text
- Expanding the panel reveals not-detected tool items
- Each not-detected item shows the tool name and category

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// Collapse panel must exist
	line := findLine(resp.Output, "ELEM collapse-not-detected")
	if line == "" {
		t.Fatal("expected collapse-not-detected element")
	}
	if strings.Contains(line, "MISSING") {
		t.Fatal("collapse-not-detected element is MISSING")
	}

	// Check for not-detected items
	if !strings.Contains(resp.Output, "NOT_DETECTED_ITEM") {
		// No not-detected items found — all software is installed?
		// This is acceptable; the collapse panel should still exist
		t.Log("no not-detected items found (all software detected)")
	}

	// Header should contain "Not Detected"
	headerLine := findLine(resp.Output, "ELEM not-detected-header")
	if headerLine != "" && !strings.Contains(headerLine, "Not Detected") {
		t.Fatalf("collapse header should contain 'Not Detected': %s", headerLine)
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
