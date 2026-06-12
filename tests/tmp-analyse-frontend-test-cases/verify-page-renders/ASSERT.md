## Expected
- `ELEM page-heading` contains "Tmp Files Analyse" and visible=true
- `ELEM start-scan-btn` is visible=true
- `ELEM stop-scan-btn` is visible=false (hidden before scan starts)
- `ELEM summary-bar` is present
- `ELEM total-size` is present
- `ELEM reclaimable-size` is present
- `ELEM section-core-heading` is present and visible
- `ELEM section-software-heading` is present and visible
- Core cards: `card-trash`, `card-temp`, `card-cache`, `card-log` are visible
- Software section exists with `section-software` element
- Software cards exist: `card-go`, `card-npm`, `card-docker`, `card-bun`
- All cards have `card-label`, `card-size`, `reboot-safe-badge` elements
- `ELEM collapse-not-detected` is present (collapsed by default)
- No element is reported as MISSING

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	checks := map[string]string{
		"ELEM page-heading":            "Tmp Files Analyse",
		"ELEM start-scan-btn":          "visible=true",
		"ELEM stop-scan-btn":           "visible=false",
		"ELEM summary-bar":             "visible=true",
		"ELEM total-size":              "",
		"ELEM reclaimable-size":        "",
		"ELEM section-core-heading":    "visible=true",
		"ELEM section-software-heading":"visible=true",
		"ELEM card-trash":              "visible=true",
		"ELEM card-temp":               "visible=true",
		"ELEM card-cache":              "visible=true",
		"ELEM card-log":                "visible=true",
		"ELEM section-software":        "visible=true",
		"ELEM card-go":                 "visible=true",
		"ELEM card-npm":                "visible=true",
		"ELEM card-docker":             "visible=true",
		"ELEM card-bun":                "visible=true",
		"ELEM collapse-not-detected":   "",
	}

	for prefix, expected := range checks {
		line := findLine(resp.Output, prefix)
		if line == "" {
			t.Fatalf("missing expected output line: %s", prefix)
		}
		if expected != "" && !strings.Contains(line, expected) {
			t.Fatalf("line %q: expected to contain %q", line, expected)
		}
		if strings.Contains(line, "MISSING") {
			t.Fatalf("element missing: %s", line)
		}
	}

	if !strings.Contains(resp.Output, "COUNT all-cards") {
		t.Fatal("expected COUNT lines for card elements")
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
