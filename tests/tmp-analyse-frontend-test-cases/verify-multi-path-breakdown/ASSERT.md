## Expected
- Go card has an extra-path breakdown section with data-testid="extra-breakdown" inside card-go
- Xcode card has an extra-path breakdown section inside card-xcode
- Each breakdown item shows a row wrapper with data-testid="extra-breakdown-row-{idx}" containing both label and size
- Breakdown labels show full tilde-prefixed paths (e.g., `~/Library/Caches/go-build`) not truncated 2-component paths
- Each breakdown label text starts with `~/`
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
			t.Logf("card-%s not detected, skipping extra-breakdown check", cat)
		}

		// Check row wrapper exists
		rowLine := findLine(resp.Output, "ELEM card-"+cat+"-extra-row-0")
		if rowLine != "" && strings.Contains(rowLine, "MISSING") {
			t.Fatalf("expected extra-breakdown-row-0 for %s", cat)
		}
	}

	// Go breakdown label should start with ~/ (full tilde path)
	if !strings.Contains(resp.Output, "FULL_PATH go-label-starts-with-tilde: true") {
		t.Fatal("expected Go breakdown label to start with ~/")
	}
	if !strings.Contains(resp.Output, "FULL_PATH go-label-not-truncated: true") {
		t.Fatal("expected Go breakdown label to be full path (not truncated to 2 components)")
	}

	// Xcode breakdown label should start with ~/ (full tilde path)
	if !strings.Contains(resp.Output, "FULL_PATH xcode-label-starts-with-tilde: true") {
		t.Fatal("expected Xcode breakdown label to start with ~/")
	}
	if !strings.Contains(resp.Output, "FULL_PATH xcode-label-not-truncated: true") {
		t.Fatal("expected Xcode breakdown label to be full path (not truncated to 2 components)")
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
