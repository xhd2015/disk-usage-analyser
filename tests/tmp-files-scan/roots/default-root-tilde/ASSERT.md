## Expected

- Exactly one binary is reported.
- `ScanResult.Roots` contains the fixture home directory.
- Stdout uses `~/Projects/default-app/...`.

## Side Effects

- No real home directory is scanned.

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
		t.Fatalf("expected one binary, got %#v", resp.Result)
	}
	if len(resp.Result.Roots) != 1 || resp.Result.Roots[0] != req.HomeDir {
		t.Fatalf("expected root %q, got %#v", req.HomeDir, resp.Result.Roots)
	}
	if !strings.Contains(resp.Stdout, "~/Projects/default-app/bin/default-app") {
		t.Fatalf("expected tilde path in stdout:\n%s", resp.Stdout)
	}
}
```
