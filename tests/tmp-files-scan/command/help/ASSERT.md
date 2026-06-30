## Expected

- Exit code is 0.
- Help text includes the approved command and all documented options.

## Side Effects

- No filesystem scan is performed.

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
		t.Fatalf("unexpected help error: %v", resp.Err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", resp.ExitCode)
	}
	for _, want := range []string{
		"disk-usage-analyser tmp-files scan [OPTIONS]",
		"--go-binaries",
		"--root",
		"--max-depth",
		"--json",
		"-v, --verbose",
		"-h, --help",
	} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("help output missing %q:\n%s", want, resp.Stdout)
		}
	}
}
```
