## Expected
- The swap card shows a non-reclaimable badge or indicator
- The indicator text indicates swap cannot be reclaimed (e.g., "OS Managed", "Not Reclaimable")
- This distinguishes swap from regular reclaimable items

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP swap-non-reclaimable") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if strings.Contains(resp.Output, "SKIP swap-non-reclaimable") {
		t.Log("swap card not found, skipping non-reclaimable check")
		return
	}

	line := findLine(resp.Output, "SWAP_NON_RECLAIMABLE_EXISTS")
	if line == "" || !strings.Contains(line, "true") {
		t.Fatal("expected non-reclaimable badge on swap card")
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
