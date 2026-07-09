## Expected

- Exit code 0.
- Human explain succeeds (`KIND: go-build-cache`, HOW TO PURGE / RAW present).
- Stdout contains **no** ANSI escape sequences (`\x1b[` / CSI).
- `$ ` prefixes may still be present (prompt is independent of color).

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d (err=%v stderr=%q)", resp.ExitCode, resp.Err, resp.Stderr)
	}
	assertKindLine(t, resp.Stdout, "go-build-cache")
	assertNoANSI(t, resp.Stdout)
	assertNoANSI(t, resp.Stderr)
	// Prompt still expected on human command lines when product implements $ prefix.
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertHowToPurgeHasDollarCommand(t, resp.Stdout, "go")
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
