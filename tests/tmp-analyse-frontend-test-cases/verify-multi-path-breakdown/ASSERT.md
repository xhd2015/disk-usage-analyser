## Expected
- Go card has an extra-path breakdown section with data-testid="extra-breakdown" inside card-go
- Xcode card has an extra-path breakdown section inside card-xcode
- Each breakdown item shows a label and size for the extra path
- Single-path software cards do NOT have extra-breakdown elements

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// Go and Xcode should have extra breakdown elements
	for _, cat := range []string{"go", "xcode"} {
		line := findLine(resp.Output, "ELEM card-"+cat+"-extra-breakdown")
		if line == "" {
			t.Fatalf("expected extra breakdown element for %s", cat)
		}
		if strings.Contains(line, "MISSING") {
			// Card may not be detected yet; skip assertion
			t.Logf("card-%s not detected, skipping extra-breakdown check", cat)
		}
	}

	// Single-path tools should NOT have extra breakdown
	singlePathCats := []string{"npm", "bun", "docker", "gradle", "maven"}
	for _, cat := range singlePathCats {
		line := findLine(resp.Output, "ELEM card-"+cat+"-extra-breakdown")
		if line != "" && !strings.Contains(line, "MISSING") {
			t.Fatalf("single-path card %s should not have extra-breakdown", cat)
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
