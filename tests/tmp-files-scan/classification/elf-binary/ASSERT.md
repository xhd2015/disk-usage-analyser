## Expected

- One hit is returned with `Kind == "elf"`.
- The hit is streamed to stdout.

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
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one elf hit, got %#v", resp.Result)
	}
	hit := resp.Result.Binaries[0]
	if hit.Kind != "elf" {
		t.Fatalf("expected elf kind, got %#v", hit)
	}
	if !strings.Contains(resp.Stdout, "elf") || !strings.Contains(resp.Stdout, "~/Projects/elf-app/bin/elf-app") {
		t.Fatalf("expected streamed elf line, got:\n%s", resp.Stdout)
	}
}
```
