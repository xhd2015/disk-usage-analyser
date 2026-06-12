## Expected
- Go card has extra-breakdown-row-0 element with display:flex and justify-content:space-between
- Go row contains both label and size children
- Xcode card has extra-breakdown-row-0 element with display:flex and justify-content:space-between
- Xcode row contains both label and size children

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	for _, cat := range []string{"go", "xcode"} {
		// Row element exists
		rowLine := findLine(resp.Output, "ELEM "+cat+"-breakdown-row-0")
		if rowLine == "" || strings.Contains(rowLine, "MISSING") {
			t.Fatalf("expected extra-breakdown-row-0 for %s", cat)
		}

		// Flexbox layout
		displayLine := findLine(resp.Output, "ROW_STYLE "+cat+"-display")
		if !strings.Contains(displayLine, "flex") {
			t.Fatalf("expected %s row display=flex, got: %s", cat, displayLine)
		}

		justifyLine := findLine(resp.Output, "ROW_STYLE "+cat+"-justify")
		if !strings.Contains(justifyLine, "space-between") {
			t.Fatalf("expected %s row justify-content=space-between, got: %s", cat, justifyLine)
		}

		// Has child elements
		labelLine := findLine(resp.Output, "ROW_CHILDREN "+cat+"-has-label")
		if !strings.Contains(labelLine, "true") {
			t.Fatalf("expected %s row to contain label element", cat)
		}

		sizeLine := findLine(resp.Output, "ROW_CHILDREN "+cat+"-has-size")
		if !strings.Contains(sizeLine, "true") {
			t.Fatalf("expected %s row to contain size element", cat)
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
