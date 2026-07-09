## Expected

- Exit code 0.
- Tree section present (depth-1: `huge.bin`, `mid.bin`, `big/`).
- Match section header `TOP 2` with exactly two ranking lines (root skipped).
- Top entries by size: `huge.bin` (400) then `mid.bin` (200).
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
	hasDot := false
	for _, line := range strings.Split(out, "\n") {
		if line == "." {
			hasDot = true
			break
		}
	}
	if !hasDot {
		t.Fatalf("expected tree root line '.':\n%s", out)
	}
	for _, want := range []string{"huge.bin", "mid.bin", "big/"} {
		if !strings.Contains(out, want) {
			t.Fatalf("tree section missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "TOP 2") {
		t.Fatalf("expected TOP 2 header:\n%s", out)
	}
	idx := strings.Index(out, "TOP 2")
	section := out[idx:]
	hugePath := filepath.Join(req.FixtureDir, "huge.bin")
	midPath := filepath.Join(req.FixtureDir, "mid.bin")
	if !strings.Contains(section, hugePath) {
		t.Fatalf("TOP section missing huge.bin path %q:\n%s", hugePath, out)
	}
	if !strings.Contains(section, midPath) {
		t.Fatalf("TOP section missing mid.bin path %q:\n%s", midPath, out)
	}
	if strings.Index(section, hugePath) > strings.Index(section, midPath) {
		t.Fatalf("expected huge.bin ranked above mid.bin:\n%s", section)
	}
	lines := 0
	for _, line := range strings.Split(section, "\n")[1:] {
		if strings.TrimSpace(line) == "" {
			break
		}
		lines++
	}
	if lines != 2 {
		t.Fatalf("expected exactly 2 TOP match lines, got %d in:\n%s", lines, section)
	}
	stdoutEndsWithBlankLine(t, out)
}
```
