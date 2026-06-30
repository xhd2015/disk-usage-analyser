## Expected

- Stderr is empty without `-v`.
- No binaries are reported from inside the skipped root.
- Summary reports zero hits.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No fatal error is returned for a remote-backed root.

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
		t.Fatalf("remote root scan should succeed, got: %v", resp.Err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Stderr != "" {
		t.Fatalf("expected empty stderr without -v, got:\n%s", resp.Stderr)
	}
	if resp.Result == nil {
		t.Fatal("expected scan result")
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("remote root should yield zero hits, got %#v", resp.Result.Binaries)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries") {
		t.Fatalf("expected zero summary, got:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "cloud-app") {
		t.Fatalf("cloud-storage hit should not be reported:\n%s", resp.Stdout)
	}
}
```