## Expected

- Exit code 0.
- `MAX-DEPTH: 2` in summary.
- Tree includes `deep.bin` under `big/`.
- Trailing blank line.

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
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d (err=%v)\n%s", resp.ExitCode, resp.Err, resp.Stdout)
	}
	out := resp.Stdout
	if !strings.Contains(out, "MAX-DEPTH: 2") {
		t.Fatalf("expected MAX-DEPTH: 2:\n%s", out)
	}
	if !strings.Contains(out, "deep.bin") {
		t.Fatalf("max-depth 2 must show deep.bin:\n%s", out)
	}
	if !strings.Contains(out, "big/") {
		t.Fatalf("tree should still show big/:\n%s", out)
	}
	stdoutEndsWithBlankLine(t, out)
}
```
