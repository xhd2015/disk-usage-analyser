## Expected
- Server is running and /ping returns "pong"
- The Tmp Files Analyse page renders with heading and start button
- Core category cards (trash, temp, cache, log) are present
- All of this works without an externally pre-started server — the test harness starts its own

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// Server /ping must respond
	if !strings.Contains(resp.Output, "PING_OK: true") {
		t.Fatal("expected /ping to return 'pong' (server not reachable)")
	}

	// Page must render
	if !strings.Contains(resp.Output, "PAGE_HEADING present: true") {
		t.Fatal("expected page heading to be present")
	}

	if !strings.Contains(resp.Output, "START_BTN present: true") {
		t.Fatal("expected start scan button to be present")
	}

	// Core cards must exist
	for _, cat := range []string{"trash", "temp", "cache", "log"} {
		line := findLine(resp.Output, "CARD "+cat+" present")
		if !strings.Contains(line, "true") {
			t.Fatalf("expected core card %s to be present", cat)
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
