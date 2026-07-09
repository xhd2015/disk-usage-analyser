## Expected

- Exit code 0.
- Help text documents `--min` and `--max-depth`.
- Help text does **not** document `--threshold`.

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
		t.Fatalf("expected exit code 0, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	for _, want := range []string{"--min", "--max-depth"} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("help output missing %q:\n%s", want, resp.Stdout)
		}
	}
	if strings.Contains(resp.Stdout, "--threshold") {
		t.Fatalf("help must not document removed --threshold:\n%s", resp.Stdout)
	}
}
```
