## Expected
- Cards exist on the page before clicking Start Scan (visible from initial render)
- Software cards are present alongside core cards
- The initial SSE event populates the card list dynamically
- Card labels match the backend's TmpLocation labels (not hardcoded defaults)

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// Cards should be visible before any scan (from initial locations event)
	for _, cat := range []string{"trash", "temp", "cache", "log"} {
		line := findLine(resp.Output, "CARD_BEFORE_SCAN "+cat)
		if line == "" {
			t.Fatalf("card-%s not found before scan", cat)
		}
		if !strings.Contains(line, "visible=true") {
			t.Fatalf("card-%s should be visible before scan", cat)
		}
	}

	// Software cards should also be present before scan
	softwareLine := findLine(resp.Output, "SOFTWARE_CARDS_BEFORE_SCAN")
	if softwareLine == "" {
		t.Fatal("no software cards info found before scan")
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
