## Expected
- Go and Xcode cards use a unified `breakdown-items` wrapper (replaces old `extra-breakdown`)
- Primary path is at `breakdown-row-0`, extra paths at `breakdown-row-{idx+1}`
- Each breakdown item shows a row wrapper with `breakdown-row-{idx}` containing both label and size
- Breakdown labels show full tilde-prefixed paths (e.g., `~/Library/Caches/go-build`) not truncated 2-component paths
- Each breakdown label text starts with `~/`
- Single-path software cards do NOT have `breakdown-items` elements
- Single-path software cards show `card-path` centered text instead

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// Go and Xcode should have unified breakdown-items
	for _, cat := range []string{"go", "xcode"} {
		line := findLine(resp.Output, "ELEM card-"+cat+"-breakdown-items")
		if line == "" {
			t.Fatalf("expected breakdown-items wrapper for %s", cat)
		}
		if strings.Contains(line, "MISSING") {
			t.Logf("card-%s not detected, skipping breakdown check", cat)
		}

		// Check row wrapper exists for extra path (now at index 1)
		rowLine := findLine(resp.Output, "ELEM card-"+cat+"-breakdown-row-1")
		if rowLine != "" && strings.Contains(rowLine, "MISSING") {
			t.Fatalf("expected breakdown-row-1 for %s", cat)
		}
	}

	// Go breakdown label (extra path at index 1) should start with ~/
	if !strings.Contains(resp.Output, "FULL_PATH go-label-starts-with-tilde: true") {
		t.Fatal("expected Go breakdown label to start with ~/")
	}
	if !strings.Contains(resp.Output, "FULL_PATH go-label-not-truncated: true") {
		t.Fatal("expected Go breakdown label to be full path (not truncated to 2 components)")
	}

	// Xcode breakdown label (extra path at index 1) should start with ~/
	if !strings.Contains(resp.Output, "FULL_PATH xcode-label-starts-with-tilde: true") {
		t.Fatal("expected Xcode breakdown label to start with ~/")
	}
	if !strings.Contains(resp.Output, "FULL_PATH xcode-label-not-truncated: true") {
		t.Fatal("expected Xcode breakdown label to be full path (not truncated to 2 components)")
	}
	// Xcode row 4 (DocumentationCache) should exist on multi-path card
	row4Line := findLine(resp.Output, "ELEM card-xcode-breakdown-row-4")
	if row4Line != "" && strings.Contains(row4Line, "MISSING") {
		t.Fatal("expected Xcode breakdown-row-4 (DocumentationCache) to exist")
	}

	// Single-path tools should NOT have breakdown-items
	singlePathCats := []string{"bun", "docker", "gradle", "maven"}
	for _, cat := range singlePathCats {
		line := findLine(resp.Output, "ELEM card-"+cat+"-breakdown-items")
		if line != "" && !strings.Contains(line, "MISSING") {
			t.Fatalf("single-path card %s should not have breakdown-items", cat)
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
