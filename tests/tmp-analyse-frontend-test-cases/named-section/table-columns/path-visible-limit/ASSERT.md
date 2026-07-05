## Expected

- Harness returns `limit: 56` from `PATH_VISIBLE_CHAR_LIMIT`.

## Errors

- Harness error or limit not equal to 56.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("path visible limit harness failed: %v\n%s", err, resp.Output)
	}
	limit, ok := resp.JSON["limit"].(float64)
	if !ok {
		t.Fatalf("limit missing in JSON: %v", resp.JSON)
	}
	if int(limit) != 56 {
		t.Fatalf("PATH_VISIBLE_CHAR_LIMIT = %d, want 56", int(limit))
	}
}
```