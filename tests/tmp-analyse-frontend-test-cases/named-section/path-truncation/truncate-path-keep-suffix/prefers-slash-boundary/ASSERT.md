## Expected

- `display` starts with `…` or `...` followed immediately by `/` (slash-aligned prefix cut).
- `display` ends with `/acme/widget/node_modules`.

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
	afterEllipsis := strings.TrimPrefix(strings.TrimPrefix(display, "..."), "…")
	if !strings.HasPrefix(afterEllipsis, "/") {
		t.Fatalf("display %q want slash immediately after ellipsis prefix", display)
	}
	if !strings.HasSuffix(display, "/acme/widget/node_modules") {
		t.Fatalf("display %q want suffix /acme/widget/node_modules", display)
	}
}
```