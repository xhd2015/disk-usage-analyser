## Expected

- Stdout contains a hit line with size, `macho`, tilde path, and repo annotation.
- The final line is a summary.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No error is returned.

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
	if resp.Err != nil {
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected one hit line and summary, got %d lines:\n%s", len(lines), resp.Stdout)
	}
	hitLine := lines[0]
	for _, want := range []string{" B", "macho", "~/Projects/human-app/bin/human-app", "(repo: ~/Projects/human-app)"} {
		if !strings.Contains(hitLine, want) {
			t.Fatalf("human hit line missing %q:\n%s", want, hitLine)
		}
	}
	if !strings.HasPrefix(lines[1], "Found 1 binaries, total ") {
		t.Fatalf("bad summary line: %q", lines[1])
	}
}
```
