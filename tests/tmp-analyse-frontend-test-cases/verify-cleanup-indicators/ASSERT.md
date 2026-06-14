## Expected
- Every section card has a cleanup-indicator element visible
- Categories: trash, temp, cache, log, go, npm, bun, yarn, pnpm, pip, cargo, ruby, docker, podman, nginx, gradle, maven, android, brew, xcode, composer
- No indicator is MISSING

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	allCategories := []string{
		"trash", "temp", "cache", "log",
		"go", "npm", "bun", "yarn", "pnpm", "pip", "cargo",
		"ruby", "docker", "podman", "nginx", "gradle", "maven",
		"android", "brew", "xcode", "composer",
	}

	for _, cat := range allCategories {
		prefix := "ELEM cleanup-indicator-" + cat
		line := findLine(resp.Output, prefix)
		if line == "" {
			t.Fatalf("missing cleanup indicator for category: %s", cat)
		}
		if strings.Contains(line, "MISSING") {
			t.Fatalf("cleanup indicator missing for category: %s", cat)
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
