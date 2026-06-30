## Expected

- One hit is returned with `Kind == "go"`.
- The type description is still populated from file detection.

## Side Effects

- A temporary Go fixture binary is built under the test repository.

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
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one go hit, got %#v", resp.Result)
	}
	hit := resp.Result.Binaries[0]
	if hit.Kind != "go" {
		t.Fatalf("expected go kind from buildinfo precedence, got %#v", hit)
	}
	if hit.TypeDesc == "" {
		t.Fatalf("expected non-empty type description, got %#v", hit)
	}
	if !strings.Contains(resp.Stdout, "go") || !strings.Contains(resp.Stdout, "~/Projects/go-app/bin/go-app") {
		t.Fatalf("expected streamed go line, got:\n%s", resp.Stdout)
	}
}
```
