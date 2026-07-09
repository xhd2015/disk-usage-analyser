## Expected

- Exit code 0.
- Tree section present.
- Match section lists `.bin` nodes (`huge.bin`, `mid.bin`, `deep.bin`) ranked by size.
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
	if !strings.Contains(out, "huge.bin") {
		t.Fatalf("Option B requires tree section with huge.bin:\n%s", out)
	}
	if !strings.Contains(out, "TOP ") && !strings.Contains(out, "MATCH") {
		t.Fatalf("expected match section for --suffix:\n%s", out)
	}
	for _, name := range []string{"huge.bin", "mid.bin", "deep.bin"} {
		if !strings.Contains(out, name) {
			t.Fatalf("suffix matches should include %q:\n%s", name, out)
		}
	}
	hugePath := filepath.Join(req.FixtureDir, "huge.bin")
	if !strings.Contains(out, hugePath) {
		t.Fatalf("expected huge.bin path in matches: %q\n%s", hugePath, out)
	}
	stdoutEndsWithBlankLine(t, out)
}
```
