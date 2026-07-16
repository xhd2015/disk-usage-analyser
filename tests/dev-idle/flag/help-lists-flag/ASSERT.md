## Expected

- Root help is printed.
- Help text contains `--dev-idle-life`.

## Side Effects

- No web server start.

## Exit Code

- 0 (nil error from run).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ServerWasStarted {
		t.Fatal("root -h must not start web server")
	}
	if resp.Err != nil {
		t.Fatalf("root -h must return nil error, got: %v", resp.Err)
	}
	out := resp.Stdout + resp.Stderr
	if !strings.Contains(out, "--dev-idle-life") {
		t.Fatalf("root help missing --dev-idle-life:\n%s", out)
	}
}
```