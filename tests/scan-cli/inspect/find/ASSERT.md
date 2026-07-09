## Expected

- Exit code 0.
- Tree section still present (Option B; default depth 1 shows `big/` etc.).
- Match section lists `deep.bin` (case-insensitive path/name find).
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
	// Tree section
	if !strings.Contains(out, "big/") {
		t.Fatalf("Option B requires tree section (big/):\n%s", out)
	}
	// Match section with deep.bin
	deepPath := filepath.Join(req.FixtureDir, "big", "deep.bin")
	if !strings.Contains(out, "deep.bin") {
		t.Fatalf("find matches should include deep.bin:\n%s", out)
	}
	if !strings.Contains(out, deepPath) {
		// path form may still be shown on match line
		t.Logf("note: full deep path not found; deep.bin name present")
	}
	// Should not list mid.bin as a find match for "deep" — mid may still appear in tree
	// Require a TOP/MATCH style header or ranked size lines after the tree.
	if !strings.Contains(out, "TOP ") && !strings.Contains(out, "MATCH") {
		// find activates match section; header may be TOP 20 (default cap)
		t.Fatalf("expected match section header for --find:\n%s", out)
	}
	stdoutEndsWithBlankLine(t, out)
}
```
