## Expected

- `display` starts with `…` or `...`.
- `display` ends with `/my-repo/node_modules` (full suffix preserved).
- `display` length is at most `maxVisibleChars`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("path display harness failed: %v\n%s", err, resp.Output)
	}
	display, ok := resp.JSON["display"].(string)
	if !ok {
		t.Fatalf("display missing in JSON: %v", resp.JSON)
	}
	maxChars, ok := resp.JSON["maxVisibleChars"].(float64)
	if !ok {
		t.Fatalf("maxVisibleChars missing in JSON: %v", resp.JSON)
	}
	if !(strings.HasPrefix(display, "…") || strings.HasPrefix(display, "...")) {
		t.Fatalf("display %q want ellipsis prefix", display)
	}
	if !strings.HasSuffix(display, "/my-repo/node_modules") {
		t.Fatalf("display %q want suffix /my-repo/node_modules", display)
	}
	if len(display) > int(maxChars) {
		t.Fatalf("display length %d exceeds maxVisibleChars %d", len(display), int(maxChars))
	}
}
```