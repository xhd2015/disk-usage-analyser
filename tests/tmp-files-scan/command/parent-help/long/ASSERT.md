## Expected

- Exit code is 0; `resp.Err` is nil.
- Stdout is parent/scan usage help (same content family as empty args / `scan -h`).
- Help documents `tmp-files scan [OPTIONS]` and key scan flags.
- No filesystem scan is performed (no binary/named hit lines or summary).
- Stdout ends with a trailing newline `\n`.

## Side Effects

- No scan of the fixture home directory.

## Errors

- No error is returned (must not treat `--help` as an unknown command).

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
		t.Fatalf("tmp-files --help must not error (got: %v)\nstdout=%q stderr=%q", resp.Err, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0 for parent --help, got %d", resp.ExitCode)
	}
	if resp.Result != nil {
		t.Fatalf("parent --help must not run a scan; got Result=%+v", resp.Result)
	}
	out := resp.Stdout
	if strings.Contains(strings.ToLower(out), "unknown") {
		t.Fatalf("parent --help must not report unknown command:\n%s", out)
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
		if !strings.Contains(out, want) {
			t.Fatalf("parent --help output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Found ") && strings.Contains(out, "binaries") {
		t.Fatalf("parent --help must not emit a scan summary:\n%s", out)
	}
	if !strings.HasSuffix(out, "\n") {
		end := out
		if len(end) > 24 {
			end = end[len(end)-24:]
		}
		if end == "" {
			end = "<empty>"
		}
		t.Fatalf("user-facing help stdout must end with trailing \\n, got ending %q", end)
	}
}
```
