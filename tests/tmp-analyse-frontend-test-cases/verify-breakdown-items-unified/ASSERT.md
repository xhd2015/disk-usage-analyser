## Expected
- Go card has unified `breakdown-items` wrapper containing all items (no primary/extra distinction)
- Go card shows `breakdown-row-0` (primary path: ~/go/pkg/mod) and `breakdown-row-1` (extra path: ~/Library/Caches/go-build)
- Xcode card shows `breakdown-row-0` through `breakdown-row-4` (DerivedData + four ExtraPaths)
- Each row uses flexbox layout with `display:flex` and `justify-content:space-between`
- Multi-item cards do NOT have standalone `card-path` elements
- Single-item cards (npm, bun, docker, gradle) still have `card-path` centered text
- Single-item cards do NOT have `breakdown-items` wrapper
- Cards that are not detected (e.g. Xcode on non-macOS) are gracefully skipped

```go
import (
	"strconv"
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// === Go Card (multi-item) ===
	goExistsLine := findLine(resp.Output, "CARD_EXISTS go:")
	if strings.Contains(goExistsLine, "true") {
		// Must have breakdown-items wrapper
		line := findLine(resp.Output, "ELEM go-breakdown-items")
		if line == "" || strings.Contains(line, "MISSING") {
			t.Fatal("expected Go card to have breakdown-items wrapper")
		}

		// Row 0: primary path
		for _, label := range []string{
			"go-breakdown-row-0",
			"go-breakdown-label-0",
			"go-breakdown-size-0",
		} {
			line := findLine(resp.Output, "ELEM "+label)
			if line == "" || strings.Contains(line, "MISSING") {
				t.Fatalf("expected %s to exist", label)
			}
		}

		// Row 1: extra path
		for _, label := range []string{
			"go-breakdown-row-1",
			"go-breakdown-label-1",
			"go-breakdown-size-1",
		} {
			line := findLine(resp.Output, "ELEM "+label)
			if line == "" || strings.Contains(line, "MISSING") {
				t.Fatalf("expected %s to exist", label)
			}
		}

		// Must NOT have standalone card-path
		if !strings.Contains(resp.Output, "NO_STANDALONE_PATH go: true") {
			t.Fatal("expected Go card to NOT have standalone card-path (all paths are in breakdown-items)")
		}

		// Row 0 layout: flexbox
		displayLine := findLine(resp.Output, "ROW_LAYOUT go-row0-display")
		if !strings.Contains(displayLine, "flex") {
			t.Fatalf("expected Go row-0 display=flex, got: %s", displayLine)
		}
		justifyLine := findLine(resp.Output, "ROW_LAYOUT go-row0-justify")
		if !strings.Contains(justifyLine, "space-between") {
			t.Fatalf("expected Go row-0 justify-content=space-between, got: %s", justifyLine)
		}
	}

	// === Xcode Card (multi-item, macOS only) ===
	xcodeExistsLine := findLine(resp.Output, "CARD_EXISTS xcode:")
	if strings.Contains(xcodeExistsLine, "true") {
		line := findLine(resp.Output, "ELEM xcode-breakdown-items")
		if line == "" || strings.Contains(line, "MISSING") {
			t.Fatal("expected Xcode card to have breakdown-items wrapper")
		}

		for i := 0; i <= 4; i++ {
			for _, suffix := range []string{"breakdown-row", "breakdown-label", "breakdown-size"} {
				label := "xcode-" + suffix + "-" + strconv.Itoa(i)
				line := findLine(resp.Output, "ELEM "+label)
				if line == "" || strings.Contains(line, "MISSING") {
					t.Fatalf("expected %s to exist", label)
				}
			}
		}

		if !strings.Contains(resp.Output, "NO_STANDALONE_PATH xcode: true") {
			t.Fatal("expected Xcode card to NOT have standalone card-path")
		}

		displayLine := findLine(resp.Output, "ROW_LAYOUT xcode-row0-display")
		if !strings.Contains(displayLine, "flex") {
			t.Fatalf("expected Xcode row-0 display=flex, got: %s", displayLine)
		}
		justifyLine := findLine(resp.Output, "ROW_LAYOUT xcode-row0-justify")
		if !strings.Contains(justifyLine, "space-between") {
			t.Fatalf("expected Xcode row-0 justify-content=space-between, got: %s", justifyLine)
		}
	}

	// === Single-item cards: keep card-path, no breakdown-items ===
	for _, cat := range []string{"npm", "bun", "docker", "gradle"} {
		existsLine := findLine(resp.Output, "CARD_EXISTS "+cat)
		if strings.Contains(existsLine, "false") {
			continue
		}

		// Must have card-path
		cpLine := findLine(resp.Output, "HAS_CARD_PATH "+cat)
		if !strings.Contains(cpLine, "true") {
			t.Fatalf("expected single-item card %s to have card-path", cat)
		}

		// Must NOT have breakdown-items
		bdLine := findLine(resp.Output, "NO_BREAKDOWN_ITEMS "+cat)
		if !strings.Contains(bdLine, "true") {
			t.Fatalf("expected single-item card %s to NOT have breakdown-items", cat)
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
