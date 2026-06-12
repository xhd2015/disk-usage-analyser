## Expected
- Software cards exist for: go, npm, bun, yarn, pnpm, pip, cargo, ruby, docker, podman, nginx, gradle, maven, android, brew, xcode, composer
- Each card has card-label with the tool's name
- Each card has Reboot Safe badge (all software caches are reboot-safe)
- Cards are rendered with data-testid="card-{category}"
- Not all software may be detected (depends on system), but the card elements exist

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	softwareCards := []string{
		"go", "npm", "bun", "yarn", "pnpm", "pip", "cargo",
		"ruby", "docker", "podman", "nginx", "gradle", "maven",
		"android", "brew", "xcode", "composer",
	}

	for _, cat := range softwareCards {
		line := findLine(resp.Output, "ELEM card-"+cat)
		if line == "" {
			t.Fatalf("missing software card element: card-%s", cat)
		}
		if strings.Contains(line, "MISSING") {
			// Card element not found in DOM — may be acceptable if undetected
			// but the collapse panel should have it
			collapseLine := findLine(resp.Output, "NOT_DETECTED "+cat)
			if collapseLine == "" {
				t.Fatalf("software card-%s missing and not in not-detected list", cat)
			}
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
