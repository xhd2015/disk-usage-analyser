## Expected

- Exit code 0.
- Tree is focused on the `big` subtree (shows `deep.bin` as child of focused root).
- **No** `TOP` match section.
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
	if !strings.Contains(out, "deep.bin") {
		t.Fatalf("--at big should show deep.bin in focused tree:\n%s", out)
	}
	// Global siblings should not appear as top-level tree peers of the focus
	if strings.Contains(out, "mid.bin") {
		t.Fatalf("--at alone focused tree should not list mid.bin as a peer:\n%s", out)
	}
	if strings.Contains(out, "huge.bin") {
		t.Fatalf("--at alone focused tree should not list huge.bin as a peer:\n%s", out)
	}
	stdoutHasNoTopSection(t, out)
	stdoutEndsWithBlankLine(t, out)
}
```
