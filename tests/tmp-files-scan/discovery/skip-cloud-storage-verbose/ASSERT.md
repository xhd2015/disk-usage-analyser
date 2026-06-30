## Expected

- Stderr contains a remote-backed filesystem skip warning mentioning CloudStorage.
- The local repository hit is still reported.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No fatal error is returned.

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
		t.Fatalf("verbose cloud-storage scan should succeed, got: %v", resp.Err)
	}
	stderr := strings.ToLower(resp.Stderr)
	for _, want := range []string{"warning", "remote-backed filesystem", "cloudstorage"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr missing %q; got:\n%s", want, resp.Stderr)
		}
	}
	if !strings.Contains(resp.Stdout, "local-app") {
		t.Fatalf("local hit missing:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "cloud-app") {
		t.Fatalf("cloud-storage hit should be skipped:\n%s", resp.Stdout)
	}
}
```