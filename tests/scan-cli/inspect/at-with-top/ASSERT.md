## Expected

- Exit code 0.
- Focused tree section present (`deep.bin`).
- Match section present (`TOP 5` or similar) with at least one match under `big`.
- Trailing blank line.

## Exit Code

- 0

```go
import (
	"path/filepath"
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
	if !strings.Contains(out, "deep.bin") {
		t.Fatalf("focused tree should include deep.bin:\n%s", out)
	}
	if !strings.Contains(out, "TOP ") {
		t.Fatalf("expected TOP match section when --at is combined with --top:\n%s", out)
	}
	deepPath := filepath.Join(req.FixtureDir, "big", "deep.bin")
	if !strings.Contains(out, deepPath) && !strings.Contains(out, filepath.Join(req.FixtureDir, "big")) {
		t.Fatalf("match section should reference focused subtree paths:\n%s", out)
	}
	stdoutEndsWithBlankLine(t, out)
}
```
