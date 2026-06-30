## Expected

- No binary hits are returned.
- Summary says `Found 0 binaries`.
- Text filenames do not appear as hits.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No error is returned.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	if resp.Result == nil || len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected no hits, got %#v", resp.Result)
	}
	if strings.Contains(resp.Stdout, "README.txt") || strings.Contains(resp.Stdout, "main.go") {
		t.Fatalf("text files appeared in stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries") {
		t.Fatalf("missing zero summary:\n%s", resp.Stdout)
	}
}
```
