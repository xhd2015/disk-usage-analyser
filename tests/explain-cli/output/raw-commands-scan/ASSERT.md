## Expected

- Exit code 0.
- Output contains a `RAW COMMANDS` section.
- Runnable RAW lines use **`$ `** prefix; group comments (`# …`) do not.
- At least one line is `$ disk-usage-analyser scan` (token sequence present with `$ `).
- The explained absolute path appears in the scan command line or nearby RAW COMMANDS block.
- System/ecosystem commands may also appear; they are informational only.

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
		t.Fatalf("expected exit 0, got %d (err=%v stderr=%q)", resp.ExitCode, resp.Err, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "RAW COMMANDS") {
		t.Fatalf("missing RAW COMMANDS section:\n%s", resp.Stdout)
	}
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)

	foundScanWithDollar := false
	foundScanWithPath := false
	for _, line := range strings.Split(resp.Stdout, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "$ ") && strings.Contains(trim, "disk-usage-analyser scan") {
			foundScanWithDollar = true
			if strings.Contains(trim, req.TargetPath) {
				foundScanWithPath = true
			}
		}
	}
	if !foundScanWithDollar {
		t.Fatalf("missing \"$ disk-usage-analyser scan\" line:\n%s", resp.Stdout)
	}
	if !foundScanWithPath {
		// Soft: path may be only on PATH: line; still require path somewhere.
		if !strings.Contains(resp.Stdout, req.TargetPath) {
			t.Fatalf("missing path %q in output:\n%s", req.TargetPath, resp.Stdout)
		}
	}
	assertNoANSI(t, resp.Stdout) // default auto + non-TTY buffer
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
