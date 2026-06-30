## Expected

- Stderr is empty without `-v`.
- The local repository hit is still reported.
- Binaries inside the CloudStorage tree are not reported.
- The command completes without a fatal walk error.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- Remote-backed paths must not abort the scan.

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
		t.Fatalf("cloud-storage skip scan should succeed, got: %v", resp.Err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", resp.ExitCode)
	}
	if resp.Stderr != "" {
		t.Fatalf("expected empty stderr without -v, got:\n%s", resp.Stderr)
	}
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one local hit, got %#v", resp.Result)
	}
	if !strings.Contains(resp.Stdout, "local-app") {
		t.Fatalf("local hit missing:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "cloud-app") {
		t.Fatalf("cloud-storage hit should be skipped:\n%s", resp.Stdout)
	}
}
```