## Expected

- Only `bin/keep` is reported.
- No hit path contains `vendor/` or `node_modules/`.

## Errors

- No harness error is returned.

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Binaries) != 1 {
		t.Fatalf("expected only kept hit, got %#v", resp.Binaries)
	}
	if !strings.Contains(resp.Binaries[0].Path, "bin/keep") {
		t.Fatalf("kept binary missing: %#v", resp.Binaries[0])
	}
	for _, forbidden := range []string{"vendor", "node_modules"} {
		if strings.Contains(resp.Binaries[0].Path, forbidden) {
			t.Fatalf("ignored path appeared: %q", resp.Binaries[0].Path)
		}
	}
}
```
