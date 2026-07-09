## Expected

- Exit code 0.
- Tree section lists `a.bin`, `b.bin`, `c.bin`.
- `TOP 2` section with two matches: `a.bin` then `b.bin` (root skipped).
- `MIN: 1B` in summary (not `THRESHOLD:`).
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
	if strings.Contains(out, "THRESHOLD:") {
		t.Fatalf("must use MIN: not THRESHOLD:\n%s", out)
	}
	if !strings.Contains(out, "MIN:") {
		t.Fatalf("missing MIN: summary:\n%s", out)
	}
	for _, name := range []string{"a.bin", "b.bin", "c.bin"} {
		if !strings.Contains(out, name) {
			t.Fatalf("tree/output missing %q:\n%s", name, out)
		}
	}
	if !strings.Contains(out, "TOP 2") {
		t.Fatalf("expected TOP 2:\n%s", out)
	}
	idx := strings.Index(out, "TOP 2")
	section := out[idx:]
	aPath := filepath.Join(req.FixtureDir, "a.bin")
	bPath := filepath.Join(req.FixtureDir, "b.bin")
	if !strings.Contains(section, aPath) || !strings.Contains(section, bPath) {
		t.Fatalf("TOP section should list %q and %q:\n%s", aPath, bPath, out)
	}
	// a.bin should appear before b.bin in the match section
	if strings.Index(section, aPath) > strings.Index(section, bPath) {
		t.Fatalf("expected a.bin ranked above b.bin:\n%s", section)
	}
	lines := 0
	for _, line := range strings.Split(section, "\n")[1:] {
		if strings.TrimSpace(line) == "" {
			break
		}
		lines++
	}
	if lines != 2 {
		t.Fatalf("expected 2 match lines, got %d:\n%s", lines, section)
	}
	stdoutEndsWithBlankLine(t, out)
}
```
