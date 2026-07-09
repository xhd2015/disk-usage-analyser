## Expected

- Exit code 0.
- Summary: `PATH:`, `TOTAL:`, `MIN: 0B` (or `0`), `MAX-DEPTH: 1`, `SOURCE:` with the JSON path.
- Tree shows depth-1 children (`big/`, `mid.bin`, `tiny.bin`) and **not** `deep.bin`.
- No TOP match section.
- No `THRESHOLD:` label.
- Stdout ends with a trailing blank line.

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
		t.Fatalf("expected exit 0, got %d (err=%v)\nstdout:\n%s", resp.ExitCode, resp.Err, resp.Stdout)
	}
	out := resp.Stdout
	for _, want := range []string{"PATH:", "TOTAL:", "MIN:", "MAX-DEPTH:", "SOURCE:"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing summary %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "MIN: 0B") && !strings.Contains(out, "MIN: 0") {
		t.Fatalf("inspect default min should be 0; got:\n%s", out)
	}
	if !strings.Contains(out, "MAX-DEPTH: 1") {
		t.Fatalf("inspect default max-depth should be 1:\n%s", out)
	}
	if !strings.Contains(out, req.FixtureDir) {
		t.Fatalf("PATH should include scan root %q:\n%s", req.FixtureDir, out)
	}
	if !strings.Contains(out, "tree.json") {
		t.Fatalf("SOURCE should mention tree.json:\n%s", out)
	}
	for _, want := range []string{"huge.bin", "mid.bin", "big/"} {
		if !strings.Contains(out, want) {
			t.Fatalf("tree missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "deep.bin") {
		t.Fatalf("depth-1 view must omit deep.bin:\n%s", out)
	}
	if strings.Contains(out, "THRESHOLD:") {
		t.Fatalf("must use MIN: not THRESHOLD:\n%s", out)
	}
	stdoutHasNoTopSection(t, out)
	stdoutEndsWithBlankLine(t, out)
}
```
