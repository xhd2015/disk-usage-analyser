## Expected

- `display` equals the input path exactly (no ellipsis prefix).

```go
import (
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
	path, ok := resp.JSON["path"].(string)
	if !ok {
		t.Fatalf("path missing in JSON: %v", resp.JSON)
	}
	if display != path {
		t.Fatalf("display %q want unchanged path %q", display, path)
	}
}
```